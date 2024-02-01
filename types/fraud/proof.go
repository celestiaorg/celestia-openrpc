package fraud

import (
	"github.com/celestiaorg/celestia-openrpc/types/header"
	"github.com/celestiaorg/go-fraud"
)

// Proof embeds the fraud.Proof interface type to provide a concrete type for JSON serialization.
type Proof struct {
	fraud.Proof[*header.ExtendedHeader]
}
