package updater

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs-search/ipfs-search/index"
	indexTypes "github.com/ipfs-search/ipfs-search/index/types"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-search/updater"
)

// DefaultUpdater is the default implementation of an Updater.
type DefaultUpdater struct {
	config       *Config
	indexes      []index.Index
	invalidIndex index.Index
	*instr.Instrumentation
}

// getResults represents a subset of indexTypes.Document.
type getResults struct {
	References indexTypes.References `json:"references"`
	LastSeen   time.Time             `json:"last-seen"`
}

// hasReference returns true of a given reference exists, false when it doesn't
func hasReference(refs indexTypes.References, newR *t.Reference) bool {
	// For now, we only support IPFS.
	if newR.Parent.Protocol != t.IPFSProtocol {
		panic(fmt.Sprintf("unsupported protocol: %s", newR.Parent.Protocol))
	}

	for _, r := range refs {
		if r.ParentHash == newR.Parent.ID {
			return true
		}
	}

	return false
}

func appendReference(refs indexTypes.References, r *t.Reference) indexTypes.References {
	return append(refs, indexTypes.Reference{
		ParentHash: r.Parent.ID,
		Name:       r.Name,
	})
}

// Update conditionally updates a given resource.
func (u *DefaultUpdater) Update(ctx context.Context, r *t.ReferencedResource) (updater.UpdateStatus, error) {
	res := new(getResults)
	idx, err := index.MultiGet(ctx, u.indexes, r.ID, res, "last-seen", "references")
	if err != nil {
		return updater.UndefinedStatus, err
	}

	if idx == u.invalidIndex {
		return updater.InvalidStatus, nil
	}

	shouldUpdate := false

	if r.Parent != nil && !hasReference(res.References, r.Reference) {
		res.References = appendReference(res.References, r.Reference)
		shouldUpdate = true
	}

	// TODO: Consider whether to propagate Provider.Date to ReferencedResource and to use
	// it here.
	now := time.Now().UTC()
	if now.Sub(res.LastSeen) > u.config.MinAge {
		res.LastSeen = now
		shouldUpdate = true
	}

	if shouldUpdate {
		if err := idx.Update(ctx, r.ID, res); err != nil {
			return updater.UpdatedStatus, nil
		}

		return updater.InvalidStatus, err
	}

	return updater.UpToDateStatus, nil
}

// New returns a new DefaultUpdater.
func New(config *Config, indexes []index.Index, invalidIndex index.Index, instr *instr.Instrumentation) updater.Updater {
	return &DefaultUpdater{
		config,
		indexes,
		invalidIndex,
		instr,
	}
}
