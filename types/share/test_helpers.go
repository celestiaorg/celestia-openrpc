package share

import (
	"bytes"
	"crypto/rand"
	"sort"

	"github.com/celestiaorg/rsmt2d"
	appns "github.com/rollkit/celestia-openrpc/types/namespace"
	"github.com/stretchr/testify/require"
)

// RandEDS generates EDS filled with the random data with the given size for original square. It
// uses require.TestingT to be able to take both a *testing.T and a *testing.B.
func RandEDS(t require.TestingT, size int) *rsmt2d.ExtendedDataSquare {
	shares := RandShares(t, size*size)
	// recompute the eds
	eds, err := rsmt2d.ComputeExtendedDataSquare(shares, DefaultRSMT2DCodec(), NewConstructor(uint64(size)))
	require.NoError(t, err, "failure to recompute the extended data square")
	return eds
}

// RandShares generate 'total' amount of shares filled with random data. It uses require.TestingT
// to be able to take both a *testing.T and a *testing.B.
func RandShares(t require.TestingT, total int) []Share {
	if total&(total-1) != 0 {
		t.Errorf("total must be power of 2: %d", total)
		t.FailNow()
	}

	shares := make([]Share, total)
	for i := range shares {
		share := make([]byte, ShareSize)
		copy(share[:appns.NamespaceSize], appns.RandomNamespace().Bytes())
		_, err := rand.Read(share[appns.NamespaceSize:])
		require.NoError(t, err)
		shares[i] = share
	}
	sort.Slice(shares, func(i, j int) bool { return bytes.Compare(shares[i], shares[j]) < 0 })

	return shares
}
