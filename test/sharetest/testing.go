package sharetest

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/celestiaorg/rsmt2d"

	"github.com/rollkit/celestia-openrpc/types/namespace"
	"github.com/rollkit/celestia-openrpc/types/share"
)

// RandEDS generates EDS filled with the random data with the given size for original square. It
// uses require.TestingT to be able to take both a *testing.T and a *testing.B.
func RandEDS(size int) (*rsmt2d.ExtendedDataSquare, error) {
	shares, err := RandShares(size * size)
	if err != nil {
		return nil, err
	}
	// recompute the eds
	return rsmt2d.ComputeExtendedDataSquare(shares, share.DefaultRSMT2DCodec(), NewConstructor(uint64(size)))
}

// RandShares generate 'total' amount of shares filled with random data. It uses require.TestingT
// to be able to take both a *testing.T and a *testing.B.
func RandShares(total int) ([]share.Share, error) {
	if total&(total-1) != 0 {
		return nil, fmt.Errorf("total must be power of 2: %d", total)
	}

	var r = rand.New(rand.NewSource(time.Now().Unix())) //nolint:gosec
	shares := make([]share.Share, total)
	for i := range shares {
		shr := make([]byte, share.ShareSize)
		copy(shr[:share.NamespaceSize], RandNamespace())
		_, err := r.Read(shr[share.NamespaceSize:])
		if err != nil {
			return nil, err
		}
		shares[i] = shr
	}
	sort.Slice(shares, func(i, j int) bool { return bytes.Compare(shares[i], shares[j]) < 0 })

	return shares, nil
}

// RandNamespace generates random valid data namespace for testing purposes.
func RandNamespace() share.Namespace {
	var r = rand.New(rand.NewSource(time.Now().Unix())) //nolint:gosec
	rb := make([]byte, namespace.NamespaceVersionZeroIDSize)
	r.Read(rb) // nolint:gosec
	for {
		namespace, _ := share.NewNamespaceV0(rb)
		if err := namespace.ValidateDataNamespace(); err != nil {
			continue
		}
		return namespace
	}
}
