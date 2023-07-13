package blob

import (
	"crypto/sha256"
	"hash"
	"math/bits"
)

const (
	Size      = sha256.Size
	BlockSize = sha256.BlockSize
)

// New returns a new hash.Hash.
func New() hash.Hash {
	return sha256.New()
}

// Sum returns the SHA256 of the bz.
func Sum(bz []byte) []byte {
	h := sha256.Sum256(bz)
	return h[:]
}

//-------------------------------------------------------------

const (
	TruncatedSize = 20
)

type sha256trunc struct {
	sha256 hash.Hash
}

func (h sha256trunc) Write(p []byte) (n int, err error) {
	return h.sha256.Write(p)
}
func (h sha256trunc) Sum(b []byte) []byte {
	shasum := h.sha256.Sum(b)
	return shasum[:TruncatedSize]
}

func (h sha256trunc) Reset() {
	h.sha256.Reset()
}

func (h sha256trunc) Size() int {
	return TruncatedSize
}

func (h sha256trunc) BlockSize() int {
	return h.sha256.BlockSize()
}

// NewTruncated returns a new hash.Hash.
func NewTruncated() hash.Hash {
	return sha256trunc{
		sha256: sha256.New(),
	}
}

// SumTruncated returns the first 20 bytes of SHA256 of the bz.
func SumTruncated(bz []byte) []byte {
	hash := sha256.Sum256(bz)
	return hash[:TruncatedSize]
}

// TODO: make these have a large predefined capacity
var (
	leafPrefix  = []byte{0}
	innerPrefix = []byte{1}
)

// returns tmhash(<empty>)
func emptyHash() []byte {
	return Sum([]byte{})
}

// returns tmhash(0x00 || leaf)
func leafHash(leaf []byte) []byte {
	return Sum(append(leafPrefix, leaf...))
}

// returns tmhash(0x01 || left || right)
func innerHash(left []byte, right []byte) []byte {
	return Sum(append(innerPrefix, append(left, right...)...))
}

// HashFromByteSlices computes a Merkle tree where the leaves are the byte slice,
// in the provided order. It follows RFC-6962.
func HashFromByteSlices(items [][]byte) []byte {
	switch len(items) {
	case 0:
		return emptyHash()
	case 1:
		return leafHash(items[0])
	default:
		k := getSplitPoint(int64(len(items)))
		left := HashFromByteSlices(items[:k])
		right := HashFromByteSlices(items[k:])
		return innerHash(left, right)
	}
}

// HashFromByteSliceIterative is an iterative alternative to
// HashFromByteSlice motivated by potential performance improvements.
// (#2611) had suggested that an iterative version of
// HashFromByteSlice would be faster, presumably because
// we can envision some overhead accumulating from stack
// frames and function calls. Additionally, a recursive algorithm risks
// hitting the stack limit and causing a stack overflow should the tree
// be too large.
//
// Provided here is an iterative alternative, a test to assert
// correctness and a benchmark. On the performance side, there appears to
// be no overall difference:
//
// BenchmarkHashAlternatives/recursive-4                20000 77677 ns/op
// BenchmarkHashAlternatives/iterative-4                20000 76802 ns/op
//
// On the surface it might seem that the additional overhead is due to
// the different allocation patterns of the implementations. The recursive
// version uses a single [][]byte slices which it then re-slices at each level of the tree.
// The iterative version reproduces [][]byte once within the function and
// then rewrites sub-slices of that array at each level of the tree.
//
// Experimenting by modifying the code to simply calculate the
// hash and not store the result show little to no difference in performance.
//
// These preliminary results suggest:
//
//  1. The performance of the HashFromByteSlice is pretty good
//  2. Go has low overhead for recursive functions
//  3. The performance of the HashFromByteSlice routine is dominated
//     by the actual hashing of data
//
// Although this work is in no way exhaustive, point #3 suggests that
// optimization of this routine would need to take an alternative
// approach to make significant improvements on the current performance.
//
// Finally, considering that the recursive implementation is easier to
// read, it might not be worthwhile to switch to a less intuitive
// implementation for so little benefit.
func HashFromByteSlicesIterative(input [][]byte) []byte {
	items := make([][]byte, len(input))

	for i, leaf := range input {
		items[i] = leafHash(leaf)
	}

	size := len(items)
	for {
		switch size {
		case 0:
			return emptyHash()
		case 1:
			return items[0]
		default:
			rp := 0 // read position
			wp := 0 // write position
			for rp < size {
				if rp+1 < size {
					items[wp] = innerHash(items[rp], items[rp+1])
					rp += 2
				} else {
					items[wp] = items[rp]
					rp++
				}
				wp++
			}
			size = wp
		}
	}
}

// getSplitPoint returns the largest power of 2 less than length
func getSplitPoint(length int64) int64 {
	if length < 1 {
		panic("Trying to split a tree with size < 1")
	}
	uLength := uint(length)
	bitlen := bits.Len(uLength)
	k := int64(1 << uint(bitlen-1))
	if k == length {
		k >>= 1
	}
	return k
}
