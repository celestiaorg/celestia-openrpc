package fraud

import (
	"context"

	"github.com/celestiaorg/go-fraud"
)

type API struct {
	Subscribe func(context.Context, fraud.ProofType) (<-chan *Proof, error) `perm:"read"`
	Get       func(context.Context, fraud.ProofType) ([]Proof, error)       `perm:"read"`
}
