package types

// Reference to indexed item
type Reference struct {
	Parent *Resource
	Name   string
	Type   ResourceType
	Size   uint64
}

// TODO: (Un)marshall References such that they can be serialized to Elasticsearch
// Was:
// ParentHash string `json:"parent_hash"`
// Name       string `json:"name"`
// The idea is that References serializes to:
// [
// 	{
// 		"parent_hash": <hash>,
// 		"name": <hash>
// 	}
// ]
// Consider using different processing/index data types to allow for decoupling between
// storage serialization and internal representation.

// String shows the name
func (r *Reference) String() string {
	return r.Name
}

// References represents a list of references
type References []Reference

// Contains returns true of a given reference exists, false when it doesn't
func (refs References) Contains(newRef *Reference) bool {
	newP := newRef.Parent

	for _, r := range refs {
		oldP := r.Parent

		if oldP.Protocol != newP.Protocol {
			panic("unmatching protocols in reference")
		}

		if newP.ID == oldP.ID {
			return true
		}
	}

	return false
}

// ReferencedResource is a resource with zero or more references to it.
type ReferencedResource struct {
	*Resource
	*Reference
}
