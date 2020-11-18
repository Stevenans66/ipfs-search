package ipfs

import (
	"fmt"
	"net/url"

	ipfs "github.com/ipfs/go-ipfs-api"
	unixfs "github.com/ipfs/go-unixfs"
	unixfs_pb "github.com/ipfs/go-unixfs/pb"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

// IPFS implements the Protocol interface for the Interplanery Filesystem. It is concurrency-safe.
type IPFS struct {
	config *Config

	gatewayURL *url.URL
	shell      *ipfs.Shell

	*instr.Instrumentation
}

// absolutePath returns the absolute (CID-only) path for a resource.
func absolutePath(*t.ReferencedResource) string {
	return fmt.Sprintf("/ipfs/%s", r.ID)
}

// namedPath returns the (escaped/raw) path for a resource.
// If a reference is available, it is used to generate the filename to facilitate content
// type detection (e.g. /ipfs/<parent_hash>/my_file.jpg instead of /ipfs/<file_hash>/).
func namedPath(r *t.ReferencedResource) string {
	if ref := r.Reference; ref.Name != "" {
		return fmt.Sprintf("/ipfs/%s/%s", ref.Parent.ID, url.PathEscape(ref.Name))
	}

	return absolutePath(r)
}

// GatewayURL returns the URL to request a resource from the gateway.
// If a reference is available, it is used to generate the filename to facilitate content
// type detection (e.g. /ipfs/<parent_hash>/my_file.jpg instead of /ipfs/<file_hash>/).
func (i *IPFS) GatewayURL(r *t.ReferencedResource) string {
	url, err := i.gatewayURL.Parse(namedPath(r))

	if err != nil {
		panic(fmt.Sprintf("error generating GatewayURL: %v", err))
	}

	return url
}

type statResult struct {
	Hash string
	Type string
	Size int64 // unixfs size
}

func typeFromPb(pbType unixfs_pb.Data_DataType) t.ResourceType {
	switch pbType {
	case unixfs.TRaw, unixfs.TFile:
		return t.FileType
	case unixfs.THAMTShard, unixfs.TDirectory, unixfs.TMetadata:
		return t.DirectoryType
	default:
		return t.UnsupportedType
	}
}

func typeFromString(t string) t.ResourceType {
	switch t {
	case "file":
		return t.FileType
	case "directory":
		return t.DirectoryType
	default:
		return t.UnsupportedType
	}
}

// Stat returns a ReferencedResource with Type and Size populated.
func (i *IPFS) Stat(ctx context.Context, r *t.Resource) (*t.ReferencedResource, error) {
	ctx, cancel := context.WithDeadline(ctx, i.config.StatTimeout)
	defer cancel()

	const cmd = "files/stat"

	path := absolutePath(r)
	req := c.Request(cmd, path)

	if err := req.Exec(ctx, &statResult); err != nil {
		return nil, err
	}

	return &t.ReferencedResource{
		Type: typeFromString(stat.Type),
		Size: uint64(stat.Size),
	}
}

// Ls returns a channel with ReferencedResource's with Type and Size populated.
func (i *IPFS) Ls(context.Context, *t.Resource, chan<- t.ReferencedResource) error {

}

// New returns a new IPFS protocol.
func New(config *Config, client *http.Client, instr *instr.Instrumentation) {
	i := &IPFS{
		config,
		instr,
	}

	// Initialize gatewayURL
	gatewayURL, err := url.Parse(config.IPFSGatewayURL)
	if err != nil {
		panic(fmt.Sprintf("could not parse IPFS Gateway URL, error: %v", err))
	}

	if !i.gatewayURL.IsAbs() {
		panic(fmt.Sprintf("gateway URL is not absolute: %s", i.gatewayURL))
	}

	// Store gatewayURL for later reference
	i.gatewayURL = gatewayURL

	// Create IPFS shell
	i.shell = ipfs.NewShellWithClient(config.IPFSAPIURL, client)

	return i
}
