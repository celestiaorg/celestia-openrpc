package blob

import (
	"context"

	"github.com/rollkit/celestia-openrpc/types/share"
)

type API struct {
	Submit   func(context.Context, []*Blob, float64) (uint64, error)                               `perm:"write"`
	Get      func(context.Context, uint64, share.Namespace, Commitment) (*Blob, error)        `perm:"read"`
	GetAll   func(context.Context, uint64, []share.Namespace) ([]*Blob, error)                     `perm:"read"`
	GetProof func(context.Context, uint64, share.Namespace, Commitment) (*Proof, error)       `perm:"read"`
	Included func(context.Context, uint64, share.Namespace, *Proof, Commitment) (bool, error) `perm:"read"`
}
