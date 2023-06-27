package share

import (
	"github.com/celestiaorg/nmt"
	"github.com/celestiaorg/rsmt2d"

	"github.com/rollkit/celestia-openrpc/types/core"
)

var (
	// DefaultRSMT2DCodec is the default codec creator used for data erasure.
	DefaultRSMT2DCodec = rsmt2d.NewLeoRSCodec
)

const (
	// NamespaveVersionSize is the size of a namespace version in bytes.
	NamespaceVersionSize = 1

	// NamespaceIDSize is the size of a namespace ID in bytes.
	NamespaceIDSize = 28

	// NamespaceSize is the size of a namespace (version + ID) in bytes.
	NamespaceSize = NamespaceVersionSize + NamespaceIDSize

	// ShareSize is the size of a share in bytes.
	ShareSize = 512
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
	return s[:NamespaceSize]
}

// GetData gets data from the share.
func GetData(s Share) []byte {
	return s[NamespaceSize:]
}
