package share

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/celestiaorg/nmt/namespace"

	appns "github.com/rollkit/celestia-openrpc/types/namespace"
)

// Various reserved namespaces.
var (
	MaxReservedNamespace     = Namespace(appns.MaxReservedNamespace.Bytes())
	ParitySharesNamespace    = Namespace(appns.ParitySharesNamespace.Bytes())
	TailPaddingNamespace     = Namespace(appns.TailPaddingNamespace.Bytes())
	ReservedPaddingNamespace = Namespace(appns.ReservedPaddingNamespace.Bytes())
	TxNamespace              = Namespace(appns.TxNamespace.Bytes())
	PayForBlobNamespace      = Namespace(appns.PayForBlobNamespace.Bytes())
)

// Namespace represents namespace of a Share.
// Consists of version byte and namespace ID.
type Namespace []byte

// NewNamespaceV0 takes a variable size byte slice and creates a version 0 Namespace.
// The byte slice must be <= 10 bytes.
// If it is less than 10 bytes, it will be left padded to size 10 with 0s.
func NewNamespaceV0(id []byte) (Namespace, error) {
	if len(id) > appns.NamespaceVersionZeroIDSize {
		return nil, fmt.Errorf("namespace id must be <= %v, but it was %v bytes", appns.NamespaceVersionZeroIDSize, len(id))
	}

	n := make(Namespace, NamespaceSize)
	// version and zero padding are already set as zero,
	// so simply copying subNID to the end is enough
	copy(n[len(n)-len(id):], id)
	return n, nil
}

// NamespaceFromBytes converts bytes into Namespace and validates it.
func NamespaceFromBytes(b []byte) (Namespace, error) {
	n := Namespace(b)
	return n, n.Validate()
}

// Version reports version of the Namespace.
func (n Namespace) Version() byte {
	return n[appns.NamespaceVersionSize-1]
}

// ID reports ID of the Namespace.
func (n Namespace) ID() namespace.ID {
	return namespace.ID(n[appns.NamespaceVersionSize:])
}

// ToNMT converts the whole Namespace(both Version and ID parts) into NMT's namespace.ID
// NOTE: Once https://github.com/celestiaorg/nmt/issues/206 is closed Namespace should become NNT's
// type.
func (n Namespace) ToNMT() namespace.ID {
	return namespace.ID(n)
}

// ToAppNamespace converts the Namespace to App's definition of Namespace.
// TODO: Unify types between node and app
func (n Namespace) ToAppNamespace() appns.Namespace {
	return appns.Namespace{Version: n.Version(), ID: n.ID()}
}

// Len reports the total length of the namespace.
func (n Namespace) Len() int {
	return len(n)
}

// String stringifies the Namespace.
func (n Namespace) String() string {
	return hex.EncodeToString(n)
}

// Equals compares two Namespaces.
func (n Namespace) Equals(target Namespace) bool {
	return bytes.Equal(n, target)
}

// Validate checks if the namespace is correct.
func (n Namespace) Validate() error {
	if n.Len() != appns.NamespaceSize {
		return fmt.Errorf("invalid namespace length: expected %d, got %d", appns.NamespaceSize, n.Len())
	}
	if n.Version() != appns.NamespaceVersionZero && n.Version() != appns.NamespaceVersionMax {
		return fmt.Errorf("invalid namespace version %v", n.Version())
	}
	if len(n.ID()) != appns.NamespaceIDSize {
		return fmt.Errorf("invalid namespace id length: expected %d, got %d", appns.NamespaceIDSize, n.ID().Size())
	}
	if n.Version() == appns.NamespaceVersionZero && !bytes.HasPrefix(n.ID(), appns.NamespaceVersionZeroPrefix) {
		return fmt.Errorf("invalid namespace id: expect %d leading zeroes", len(appns.NamespaceVersionZeroPrefix))
	}
	return nil
}

// ValidateDataNamespace checks if the Namespace contains real/useful data.
func (n Namespace) ValidateDataNamespace() error {
	if err := n.Validate(); err != nil {
		return err
	}
	if n.Equals(ParitySharesNamespace) || n.Equals(TailPaddingNamespace) {
		return fmt.Errorf("invalid data namespace(%s): parity and tail padding namespace are forbidden", n)
	}
	return nil
}

// ValidateBlobNamespace checks if the Namespace is valid blob namespace.
func (n Namespace) ValidateBlobNamespace() error {
	if err := n.ValidateDataNamespace(); err != nil {
		return err
	}
	if bytes.Compare(n, MaxReservedNamespace) < 1 {
		return fmt.Errorf("invalid blob namespace(%s): reserved namespaces are forbidden", n)
	}
	return nil
}

// IsAboveMax checks if the namespace is above the maximum namespace of the given hash.
func (n Namespace) IsAboveMax(nodeHash []byte) bool {
	return !n.IsLessOrEqual(nodeHash[n.Len() : n.Len()*2])
}

// IsBelowMin checks if the target namespace is below the minimum namespace of the given hash.
func (n Namespace) IsBelowMin(nodeHash []byte) bool {
	return n.IsLess(nodeHash[:n.Len()])
}

// IsOutsideRange checks if the namespace is outside the min-max range of the given hashes.
func (n Namespace) IsOutsideRange(leftNodeHash, rightNodeHash []byte) bool {
	return n.IsBelowMin(leftNodeHash) || n.IsAboveMax(rightNodeHash)
}

// Repeat copies the Namespace t times.
func (n Namespace) Repeat(t int) []Namespace {
	ns := make([]Namespace, t)
	for i := 0; i < t; i++ {
		ns[i] = n
	}
	return ns
}

// IsLess reports if the Namespace is less than the target.
func (n Namespace) IsLess(target Namespace) bool {
	return bytes.Compare(n, target) == -1
}

// IsLessOrEqual reports if the Namespace is less than the target.
func (n Namespace) IsLessOrEqual(target Namespace) bool {
	return bytes.Compare(n, target) < 1
}

// IsGreater reports if the Namespace is greater than the target.
func (n Namespace) IsGreater(target Namespace) bool {
	return bytes.Compare(n, target) == 1
}

// IsGreaterOrEqualThan reports if the Namespace is greater or equal than the target.
func (n Namespace) IsGreaterOrEqualThan(target Namespace) bool {
	return bytes.Compare(n, target) > -1
}
