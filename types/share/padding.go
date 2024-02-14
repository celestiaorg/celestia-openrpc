package share

import (
	"bytes"
	"errors"

	"github.com/celestia/celestia-openrpc/types/appconsts"
	"github.com/celestia/celestia-openrpc/types/namespace"
)

// NamespacePaddingShare returns a share that acts as padding. Namespace padding
// shares follow a blob so that the next blob may start at an index that
// conforms to blob share commitment rules. The ns parameter provided should
// be the namespace of the blob that precedes this padding in the data square.
func NamespacePaddingShare(ns namespace.Namespace) (Share, error) {
	b, err := NewBuilder(ns, appconsts.ShareVersionZero, true).Init()
	if err != nil {
		return Share{}, err
	}
	if err := b.WriteSequenceLen(0); err != nil {
		return Share{}, err
	}
	padding := bytes.Repeat([]byte{0}, appconsts.FirstSparseShareContentSize)
	b.AddData(padding)

	paddingShare, err := b.Build()
	if err != nil {
		return Share{}, err
	}

	return *paddingShare, nil
}

// NamespacePaddingShares returns n namespace padding shares.
func NamespacePaddingShares(ns namespace.Namespace, n int) ([]Share, error) {
	var err error
	if n < 0 {
		return nil, errors.New("n must be positive")
	}
	shares := make([]Share, n)
	for i := 0; i < n; i++ {
		shares[i], err = NamespacePaddingShare(ns)
		if err != nil {
			return shares, err
		}
	}
	return shares, nil
}
