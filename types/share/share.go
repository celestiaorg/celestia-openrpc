package share

import (
	"github.com/celestiaorg/nmt"

	"github.com/rollkit/celestia-openrpc/types/core"
	appns "github.com/rollkit/celestia-openrpc/types/namespace"
)

// Root represents root commitment to multiple Shares.
// In practice, it is a commitment to all the Data in a square.
type Root = core.DataAvailabilityHeader

// NamespacedShares represents all shares with proofs within a specific namespace of an EDS.
type NamespacedShares []NamespacedRow

// NamespacedRow represents all shares with proofs within a specific namespace of a single EDS row.
type NamespacedRow struct {
	Shares []Share
	Proof  *nmt.Proof
}

// Share contains the raw share data without the corresponding namespace.
// NOTE: Alias for the byte is chosen to keep maximal compatibility, especially with rsmt2d.
// Ideally, we should define reusable type elsewhere and make everyone(Core, rsmt2d, ipld) to rely
// on it.
type Share = []byte

// GetNamespace gets the namespace GetNamespace from the share.
func GetNamespace(s Share) Namespace {
	return s[:appns.NamespaceSize]
}

// GetData gets data from the share.
func GetData(s Share) []byte {
	return s[appns.NamespaceSize:]
}
