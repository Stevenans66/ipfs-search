package protocols

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

type Protocol interface {
	SupportedProtocols() []t.Protocol
	GatewayURL(*t.ReferencedResource) (string, error)
	Stat(context.Context, *t.Resource) (*t.ReferencedResource, error)
	Ls(context.Context, *t.Resource, chan<- t.ReferencedResource) error
}
