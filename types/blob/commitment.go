package blob

import (
	"crypto/sha256"

	"github.com/celestiaorg/nmt"

	"github.com/rollkit/celestia-openrpc/types/appconsts"
	"github.com/rollkit/celestia-openrpc/types/share"
)

// CreateCommitment generates the share commitment for a given blob.
// See [data square layout rationale] and [blob share commitment rules].
//
// [data square layout rationale]: ../../specs/src/specs/data_square_layout.md
// [blob share commitment rules]: ../../specs/src/specs/data_square_layout.md#blob-share-commitment-rules
// NOTE: We assume Blob.Namespace contains both id and version bytes
func CreateCommitment(blob *Blob) ([]byte, error) {
	shares, err := SplitBlobs(*blob)
	if err != nil {
		return nil, err
	}

	// the commitment is the root of a merkle mountain range with max tree size
	// determined by the number of roots required to create a share commitment
	// over that blob. The size of the tree is only increased if the number of
	// subtree roots surpasses a constant threshold.
	subTreeWidth := share.SubTreeWidth(len(shares), appconsts.DefaultSubtreeRootThreshold)
	treeSizes, err := merkleMountainRangeSizes(uint64(len(shares)), uint64(subTreeWidth))
	if err != nil {
		return nil, err
	}
	leafSets := make([][][]byte, len(treeSizes))
	cursor := uint64(0)
	for i, treeSize := range treeSizes {
		leafSets[i] = share.ToBytes(shares[cursor : cursor+treeSize])
		cursor = cursor + treeSize
	}

	// create the commitments by pushing each leaf set onto an nmt
	subTreeRoots := make([][]byte, len(leafSets))
	for i, set := range leafSets {
		// create the nmt todo(evan) use nmt wrapper
		tree := nmt.New(sha256.New(), nmt.NamespaceIDSize(appconsts.NamespaceSize), nmt.IgnoreMaxNamespace(true))
		for _, leaf := range set {
			namespace := blob.Namespace
			// the namespace must be added again here even though it is already
			// included in the leaf to ensure that the hash will match that of
			// the nmt wrapper (pkg/wrapper). Each namespace is added to keep
			// the namespace in the share, and therefore the parity data, while
			// also allowing for the manual addition of the parity namespace to
			// the parity data.
			nsLeaf := make([]byte, 0)
			nsLeaf = append(nsLeaf, namespace...)
			nsLeaf = append(nsLeaf, leaf...)

			err = tree.Push(nsLeaf)
			if err != nil {
				return nil, err
			}
		}
		// add the root
		root, err := tree.Root()
		if err != nil {
			return nil, err
		}
		subTreeRoots[i] = root
	}
	return HashFromByteSlices(subTreeRoots), nil
}

func CreateCommitments(blobs []*Blob) ([][]byte, error) {
	commitments := make([][]byte, len(blobs))
	for i, blob := range blobs {
		commitment, err := CreateCommitment(blob)
		if err != nil {
			return nil, err
		}
		commitments[i] = commitment
	}
	return commitments, nil
}

// merkleMountainRangeSizes returns the sizes (number of leaf nodes) of the
// trees in a merkle mountain range constructed for a given totalSize and
// maxTreeSize.
//
// https://docs.grin.mw/wiki/chain-state/merkle-mountain-range/
// https://github.com/opentimestamps/opentimestamps-server/blob/master/doc/merkle-mountain-range.md
// TODO: potentially rename function because this doesn't return heights
func merkleMountainRangeSizes(totalSize, maxTreeSize uint64) ([]uint64, error) {
	var treeSizes []uint64

	for totalSize != 0 {
		switch {
		case totalSize >= maxTreeSize:
			treeSizes = append(treeSizes, maxTreeSize)
			totalSize = totalSize - maxTreeSize
		case totalSize < maxTreeSize:
			treeSize, err := share.RoundDownPowerOfTwo(totalSize)
			if err != nil {
				return treeSizes, err
			}
			treeSizes = append(treeSizes, treeSize)
			totalSize = totalSize - treeSize
		}
	}

	return treeSizes, nil
}
