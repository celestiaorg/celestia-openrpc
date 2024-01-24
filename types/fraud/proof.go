package fraud

import (
	"github.com/celestiaorg/go-fraud"
	"github.com/rollkit/celestia-openrpc/types/header"
)

// Proof embeds the fraud.Proof interface type to provide a concrete type for JSON serialization.
type Proof struct {
	fraud.Proof[*header.ExtendedHeader]
}
