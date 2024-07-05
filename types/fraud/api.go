package fraud

import (
	"context"

	"github.com/celestiaorg/go-fraud"
)

type API struct {
	// Subscribe allows to subscribe on a Proof pub sub topic by its type.
	Subscribe func(context.Context, fraud.ProofType) (<-chan *Proof, error) `perm:"read"`
	// Get fetches fraud proofs from the disk by its type.
	Get func(context.Context, fraud.ProofType) ([]Proof, error) `perm:"read"`
}
