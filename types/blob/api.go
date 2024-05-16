package blob

import (
	"context"

	"github.com/celestiaorg/celestia-openrpc/types/share"
)

type API struct {
	// Submit sends Blobs and reports the height in which they were included.
	// Allows sending multiple Blobs atomically synchronously.
	// Uses default wallet registered on the Node.
	Submit func(context.Context, []*Blob, float64) (uint64, error) `perm:"write"`
	// Get retrieves the blob by commitment under the given namespace and height.
	Get func(context.Context, uint64, share.Namespace, Commitment) (*Blob, error) `perm:"read"`
	// GetAll returns all blobs at the given height under the given namespaces.
	GetAll func(context.Context, uint64, []share.Namespace) ([]*Blob, error) `perm:"read"`
	// GetProof retrieves proofs in the given namespaces at the given height by commitment.
	GetProof func(context.Context, uint64, share.Namespace, Commitment) (*Proof, error) `perm:"read"`
	// Included checks whether a blob's given commitment(Merkle subtree root) is included at
	// given height and under the namespace.
	Included func(context.Context, uint64, share.Namespace, *Proof, Commitment) (bool, error) `perm:"read"`
}
