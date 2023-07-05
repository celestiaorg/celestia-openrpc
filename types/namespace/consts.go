package namespace

import (
	"bytes"
	"math"
)

const (
	// NamespaveVersionSize is the size of a namespace version in bytes.
	NamespaceVersionSize = 1

	// NamespaceIDSize is the size of a namespace ID in bytes.
	NamespaceIDSize = 28

	// NamespaceSize is the size of a namespace (version + ID) in bytes.
	NamespaceSize = NamespaceVersionSize + NamespaceIDSize

	// NamespaceVersionZero is the first namespace version.
	NamespaceVersionZero = uint8(0)

	// NamespaceVersionMax is the max namespace version.
	NamespaceVersionMax = math.MaxUint8

	// NamespaceZeroPrefixSize is the number of `0` bytes that are prefixed to
	// namespace IDs for version 0.
	NamespaceVersionZeroPrefixSize = 18

	// NamespaceVersionZeroIDSize is the number of bytes available for
	// user-specified namespace ID in a namespace ID for version 0.
	NamespaceVersionZeroIDSize = NamespaceIDSize - NamespaceVersionZeroPrefixSize
)

var (
	// NamespaceVersionZeroPrefix is the prefix of a namespace ID for version 0.
	NamespaceVersionZeroPrefix = bytes.Repeat([]byte{0}, NamespaceVersionZeroPrefixSize)

	// TxNamespace is the namespace reserved for transaction data.
	TxNamespace = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1})

	// IntermediateStateRootsNamespace is the namespace reserved for
	// intermediate state root data.
	IntermediateStateRootsNamespace = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 2})

	// PayForBlobNamespace is the namespace reserved for PayForBlobs transactions.
	PayForBlobNamespace = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 4})

	// ReservedPaddingNamespace is the namespace used for padding after all
	// reserved namespaces. In practice this padding is after transactions
	// (ordinary and PFBs) but before blobs.
	ReservedPaddingNamespace = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 255})

	// MaxReservedNamespace is lexicographically the largest namespace that is
	// reserved for protocol use.
	MaxReservedNamespace = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 255})

	// TailPaddingNamespace is the namespace reserved for tail padding. All data
	// with this namespace will be ignored.
	TailPaddingNamespace = Namespace{
		Version: math.MaxUint8,
		ID:      append(bytes.Repeat([]byte{0xFF}, NamespaceIDSize-1), 0xFE),
	}

	// ParitySharesNamespace is the namespace reserved for erasure coded data.
	ParitySharesNamespace = Namespace{
		Version: math.MaxUint8,
		ID:      bytes.Repeat([]byte{0xFF}, NamespaceIDSize),
	}
)
