package updater

import (
	"context"
	t "github.com/ipfs-search/ipfs-search/types"
)

// UpdateStatus reflects the status of an update.
type UpdateStatus uint8

const (
	// UndefinedStatus represents an undefined status.
	UndefinedStatus = iota
	// NotFoundStatus implies that the document was not found.
	NotFoundStatus
	// InvalidStatus implies that the document was previously found to be invalid and cannot be indexed.
	InvalidStatus
	// UpdatedStatus implies that the document has been updated.
	UpdatedStatus
	// UpToDateStatus implies that the document was already up to date.
	UpToDateStatus
)

func (s UpdateStatus) String() string {
	switch s {
	case UndefinedStatus:
		return "undefined"
	case NotFoundStatus:
		return "not found"
	case InvalidStatus:
		return "invalid"
	case UpdatedStatus:
		return "updated"
	case UpToDateStatus:
		return "up to date"
	default:
		panic("Invalid value for UpdateStatus.")
	}
}

// Updater conditionally updates indexed resources.
type Updater interface {
	Update(context.Context, *t.ReferencedResource) (UpdateStatus, error)
}
