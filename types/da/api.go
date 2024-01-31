package da

import (
	"context"

	"github.com/rollkit/go-da"
)

type API struct {
	MaxBlobSize func(ctx context.Context) (uint64, error)                                                      `perm:"read"`
	Get         func(ctx context.Context, ids []da.ID, ns da.Namespace) ([]da.Blob, error)                     `perm:"read"`
	GetProofs   func(ctx context.Context, ids []da.ID, ns da.Namespace) ([]da.Proof, error)                    `perm:"read"`
	GetIDs      func(ctx context.Context, height uint64, ns da.Namespace) ([]da.ID, error)                     `perm:"read"`
	Commit      func(ctx context.Context, blobs []da.Blob, ns da.Namespace) ([]da.Commitment, error)           `perm:"read"`
	Validate    func(ctx context.Context, ids []da.ID, proofs []da.Proof, ns da.Namespace) ([]bool, error)     `perm:"read"`
	Submit      func(ctx context.Context, blobs []da.Blob, gasPrice float64, ns da.Namespace) ([]da.ID, error) `perm:"write"`
}
