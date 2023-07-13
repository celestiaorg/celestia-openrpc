package blob

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	math "math"

	"github.com/celestiaorg/nmt"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"

	"github.com/rollkit/celestia-openrpc/types/appconsts"
	v1 "github.com/rollkit/celestia-openrpc/types/appconsts/v1"
	"github.com/rollkit/celestia-openrpc/types/namespace"
	"github.com/rollkit/celestia-openrpc/types/share"
)

// zeroPadIfNecessary pads the share with trailing zero bytes if the provided
// share has fewer bytes than width. Returns the share unmodified if the
// len(share) is greater than or equal to width.
func zeroPadIfNecessary(share []byte, width int) (padded []byte, bytesOfPadding int) {
	oldLen := len(share)
	if oldLen >= width {
		return share, 0
	}

	missingBytes := width - oldLen
	padByte := []byte{0}
	padding := bytes.Repeat(padByte, missingBytes)
	share = append(share, padding...)
	return share, missingBytes
}

// NewReservedBytes returns a byte slice of length
// appconsts.CompactShareReservedBytes that contains the byteIndex of the first
// unit that starts in a compact share.
func NewReservedBytes(byteIndex uint32) ([]byte, error) {
	if byteIndex >= appconsts.ShareSize {
		return []byte{}, fmt.Errorf("byte index %d must be less than share size %d", byteIndex, appconsts.ShareSize)
	}
	reservedBytes := make([]byte, appconsts.CompactShareReservedBytes)
	binary.BigEndian.PutUint32(reservedBytes, byteIndex)
	return reservedBytes, nil
}

// ParseReservedBytes parses a byte slice of length
// appconsts.CompactShareReservedBytes into a byteIndex.
func ParseReservedBytes(reservedBytes []byte) (uint32, error) {
	if len(reservedBytes) != appconsts.CompactShareReservedBytes {
		return 0, fmt.Errorf("reserved bytes must be of length %d", appconsts.CompactShareReservedBytes)
	}
	byteIndex := binary.BigEndian.Uint32(reservedBytes)
	if appconsts.ShareSize <= byteIndex {
		return 0, fmt.Errorf("byteIndex must be less than share size %d", appconsts.ShareSize)
	}
	return byteIndex, nil
}

// InfoByte is a byte with the following structure: the first 7 bits are
// reserved for version information in big endian form (initially `0000000`).
// The last bit is a "sequence start indicator", that is `1` if this is the
// first share of a sequence and `0` if this is a continuation share.
type InfoByte byte

func NewInfoByte(version uint8, isSequenceStart bool) (InfoByte, error) {
	if version > appconsts.MaxShareVersion {
		return 0, fmt.Errorf("version %d must be less than or equal to %d", version, appconsts.MaxShareVersion)
	}

	prefix := version << 1
	if isSequenceStart {
		return InfoByte(prefix + 1), nil
	}
	return InfoByte(prefix), nil
}

// Version returns the version encoded in this InfoByte. Version is
// expected to be between 0 and appconsts.MaxShareVersion (inclusive).
func (i InfoByte) Version() uint8 {
	version := uint8(i) >> 1
	return version
}

// IsSequenceStart returns whether this share is the start of a sequence.
func (i InfoByte) IsSequenceStart() bool {
	return uint(i)%2 == 1
}

func ParseInfoByte(i byte) (InfoByte, error) {
	isSequenceStart := i%2 == 1
	version := uint8(i) >> 1
	return NewInfoByte(version, isSequenceStart)
}

// NamespacePaddingShare returns a share that acts as padding. Namespace padding
// shares follow a blob so that the next blob may start at an index that
// conforms to blob share commitment rules. The ns parameter provided should
// be the namespace of the blob that precedes this padding in the data square.
func NamespacePaddingShare(ns namespace.Namespace) (share.Share, error) {
	b, err := NewBuilder(ns, appconsts.ShareVersionZero, true).Init()
	if err != nil {
		return share.Share{}, err
	}
	if err := b.WriteSequenceLen(0); err != nil {
		return share.Share{}, err
	}
	padding := bytes.Repeat([]byte{0}, appconsts.FirstSparseShareContentSize)
	b.AddData(padding)

	paddingShare, err := b.Build()
	if err != nil {
		return share.Share{}, err
	}

	return *paddingShare, nil
}

// NamespacePaddingShares returns n namespace padding shares.
func NamespacePaddingShares(ns namespace.Namespace, n int) ([]share.Share, error) {
	var err error
	if n < 0 {
		return nil, errors.New("n must be positive")
	}
	shares := make([]share.Share, n)
	for i := 0; i < n; i++ {
		shares[i], err = NamespacePaddingShare(ns)
		if err != nil {
			return shares, err
		}
	}
	return shares, nil
}

type Builder struct {
	namespace      namespace.Namespace
	shareVersion   uint8
	isFirstShare   bool
	isCompactShare bool
	rawShareData   []byte
}

func min[T constraints.Integer](i, j T) T {
	if i < j {
		return i
	}
	return j
}

// BlobMinSquareSize returns the minimum square size that can contain shareCount
// number of shares.
func BlobMinSquareSize(shareCount int) int {
	return RoundUpPowerOfTwo(int(math.Ceil(math.Sqrt(float64(shareCount)))))
}

// SubTreeWidth determines the maximum number of leaves per subtree in the share
// commitment over a given blob. The input should be the total number of shares
// used by that blob. The reasoning behind this algorithm is discussed in depth
// in ADR013
// (celestia-app/docs/architecture/adr-013-non-interative-default-rules-for-zero-padding).
func SubTreeWidth(shareCount, subtreeRootThreshold int) int {
	// per ADR013, we use a predetermined threshold to determine width of sub
	// trees used to create share commitments
	s := (shareCount / subtreeRootThreshold)

	// round up if the width is not an exact multiple of the threshold
	if shareCount%subtreeRootThreshold != 0 {
		s++
	}

	// use a power of two equal to or larger than the multiple of the subtree
	// root threshold
	s = RoundUpPowerOfTwo(s)

	// use the minimum of the subtree width and the min square size, this
	// gurarantees that a valid value is returned
	return min(s, BlobMinSquareSize(shareCount))
}

func NewEmptyBuilder() *Builder {
	return &Builder{
		rawShareData: make([]byte, 0, appconsts.ShareSize),
	}
}

// Init() needs to be called right after this method
func NewBuilder(ns namespace.Namespace, shareVersion uint8, isFirstShare bool) *Builder {
	return &Builder{
		namespace:      ns,
		shareVersion:   shareVersion,
		isFirstShare:   isFirstShare,
		isCompactShare: isCompactShare(ns),
	}
}

func (b *Builder) Init() (*Builder, error) {
	if b.isCompactShare {
		if err := b.prepareCompactShare(); err != nil {
			return nil, err
		}
	} else {
		if err := b.prepareSparseShare(); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (b *Builder) AvailableBytes() int {
	return appconsts.ShareSize - len(b.rawShareData)
}

func (b *Builder) ImportRawShare(rawBytes []byte) *Builder {
	b.rawShareData = rawBytes
	return b
}

func (b *Builder) AddData(rawData []byte) (rawDataLeftOver []byte) {
	// find the len left in the pending share
	pendingLeft := appconsts.ShareSize - len(b.rawShareData)

	// if we can simply add the tx to the share without creating a new
	// pending share, do so and return
	if len(rawData) <= pendingLeft {
		b.rawShareData = append(b.rawShareData, rawData...)
		return nil
	}

	// if we can only add a portion of the rawData to the pending share,
	// then we add it and add the pending share to the finalized shares.
	chunk := rawData[:pendingLeft]
	b.rawShareData = append(b.rawShareData, chunk...)

	// We need to finish this share and start a new one
	// so we return the leftover to be written into a new share
	return rawData[pendingLeft:]
}

func (b *Builder) Build() (*share.Share, error) {
	return share.NewShare(b.rawShareData)
}

// IsEmptyShare returns true if no data has been written to the share
func (b *Builder) IsEmptyShare() bool {
	expectedLen := appconsts.NamespaceSize + appconsts.ShareInfoBytes
	if b.isCompactShare {
		expectedLen += appconsts.CompactShareReservedBytes
	}
	if b.isFirstShare {
		expectedLen += appconsts.SequenceLenBytes
	}
	return len(b.rawShareData) == expectedLen
}

func (b *Builder) ZeroPadIfNecessary() (bytesOfPadding int) {
	b.rawShareData, bytesOfPadding = zeroPadIfNecessary(b.rawShareData, appconsts.ShareSize)
	return bytesOfPadding
}

// isEmptyReservedBytes returns true if the reserved bytes are empty.
func (b *Builder) isEmptyReservedBytes() (bool, error) {
	indexOfReservedBytes := b.indexOfReservedBytes()
	reservedBytes, err := ParseReservedBytes(b.rawShareData[indexOfReservedBytes : indexOfReservedBytes+appconsts.CompactShareReservedBytes])
	if err != nil {
		return false, err
	}
	return reservedBytes == 0, nil
}

// indexOfReservedBytes returns the index of the reserved bytes in the share.
func (b *Builder) indexOfReservedBytes() int {
	if b.isFirstShare {
		// if the share is the first share, the reserved bytes follow the namespace, info byte, and sequence length
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes + appconsts.SequenceLenBytes
	}
	// if the share is not the first share, the reserved bytes follow the namespace and info byte
	return appconsts.NamespaceSize + appconsts.ShareInfoBytes
}

// indexOfInfoBytes returns the index of the InfoBytes.
func (b *Builder) indexOfInfoBytes() int {
	// the info byte is immediately after the namespace
	return appconsts.NamespaceSize
}

// MaybeWriteReservedBytes will be a no-op if the reserved bytes
// have already been populated. If the reserved bytes are empty, it will write
// the location of the next unit of data to the reserved bytes.
func (b *Builder) MaybeWriteReservedBytes() error {
	if !b.isCompactShare {
		return errors.New("this is not a compact share")
	}

	empty, err := b.isEmptyReservedBytes()
	if err != nil {
		return err
	}
	if !empty {
		return nil
	}

	byteIndexOfNextUnit := len(b.rawShareData)
	reservedBytes, err := NewReservedBytes(uint32(byteIndexOfNextUnit))
	if err != nil {
		return err
	}

	indexOfReservedBytes := b.indexOfReservedBytes()
	// overwrite the reserved bytes of the pending share
	for i := 0; i < appconsts.CompactShareReservedBytes; i++ {
		b.rawShareData[indexOfReservedBytes+i] = reservedBytes[i]
	}
	return nil
}

// writeSequenceLen writes the sequence length to the first share.
func (b *Builder) WriteSequenceLen(sequenceLen uint32) error {
	if b == nil {
		return errors.New("the builder object is not initialized (is nil)")
	}
	if !b.isFirstShare {
		return errors.New("not the first share")
	}
	sequenceLenBuf := make([]byte, appconsts.SequenceLenBytes)
	binary.BigEndian.PutUint32(sequenceLenBuf, sequenceLen)

	for i := 0; i < appconsts.SequenceLenBytes; i++ {
		b.rawShareData[appconsts.NamespaceSize+appconsts.ShareInfoBytes+i] = sequenceLenBuf[i]
	}

	return nil
}

// FlipSequenceStart flips the sequence start indicator of the share provided
func (b *Builder) FlipSequenceStart() {
	infoByteIndex := b.indexOfInfoBytes()

	// the sequence start indicator is the last bit of the info byte so flip the
	// last bit
	b.rawShareData[infoByteIndex] = b.rawShareData[infoByteIndex] ^ 0x01
}

func (b *Builder) prepareCompactShare() error {
	shareData := make([]byte, 0, appconsts.ShareSize)
	infoByte, err := NewInfoByte(b.shareVersion, b.isFirstShare)
	if err != nil {
		return err
	}
	placeholderSequenceLen := make([]byte, appconsts.SequenceLenBytes)
	placeholderReservedBytes := make([]byte, appconsts.CompactShareReservedBytes)

	shareData = append(shareData, b.namespace.Bytes()...)
	shareData = append(shareData, byte(infoByte))

	if b.isFirstShare {
		shareData = append(shareData, placeholderSequenceLen...)
	}

	shareData = append(shareData, placeholderReservedBytes...)

	b.rawShareData = shareData

	return nil
}

func (b *Builder) prepareSparseShare() error {
	shareData := make([]byte, 0, appconsts.ShareSize)
	infoByte, err := NewInfoByte(b.shareVersion, b.isFirstShare)
	if err != nil {
		return err
	}
	placeholderSequenceLen := make([]byte, appconsts.SequenceLenBytes)

	shareData = append(shareData, b.namespace.Bytes()...)
	shareData = append(shareData, byte(infoByte))

	if b.isFirstShare {
		shareData = append(shareData, placeholderSequenceLen...)
	}

	b.rawShareData = shareData
	return nil
}

func isCompactShare(ns namespace.Namespace) bool {
	return ns.IsTx() || ns.IsPayForBlob()
}

// RoundDownPowerOfTwo returns the next power of two less than or equal to input.
func RoundDownPowerOfTwo[I constraints.Integer](input I) (I, error) {
	if input <= 0 {
		return 0, fmt.Errorf("input %v must be positive", input)
	}
	roundedUp := RoundUpPowerOfTwo(input)
	if roundedUp == input {
		return roundedUp, nil
	}
	return roundedUp / 2, nil
}

// RoundUpPowerOfTwo returns the next power of two greater than or equal to input.
func RoundUpPowerOfTwo[I constraints.Integer](input I) I {
	var result I = 1
	for result < input {
		result = result << 1
	}
	return result
}

// SubtreeRootThreshold works as a target upper bound for the number of subtree
// roots in the share commitment. If a blob contains more shares than this
// number, then the height of the subtree roots will increase by one so that the
// number of subtree roots in the share commitment decreases by a factor of two.
// This step is repeated until the number of subtree roots is less than the
// SubtreeRootThreshold.
//
// The rationale for this value is described in more detail in ADR-013.
func SubtreeRootThreshold(_ uint64) int {
	return v1.SubtreeRootThreshold
}

// SquareSizeUpperBound is the maximum original square width possible
// for a version of the state machine. The maximum is decided through
// governance. See `DefaultGovMaxSquareSize`.
func SquareSizeUpperBound(_ uint64) int {
	return v1.SquareSizeUpperBound
}

var (
	DefaultSubtreeRootThreshold = SubtreeRootThreshold(appconsts.LatestVersion)
	DefaultSquareSizeUpperBound = SquareSizeUpperBound(appconsts.LatestVersion)
)

// SparseShareSplitter lazily splits blobs into shares that will eventually be
// included in a data square. It also has methods to help progressively count
// how many shares the blobs written take up.
type SparseShareSplitter struct {
	shares []share.Share
}

func NewSparseShareSplitter() *SparseShareSplitter {
	return &SparseShareSplitter{}
}

// Write writes the provided blob to this sparse share splitter. It returns an
// error or nil if no error is encountered.
func (sss *SparseShareSplitter) Write(blob Blob) error {
	if !slices.Contains(appconsts.SupportedShareVersions, uint8(blob.ShareVersion)) {
		return fmt.Errorf("unsupported share version: %d", blob.ShareVersion)
	}

	rawData := blob.Data
	blobNamespace, err := namespace.New(uint8(blob.NamespaceVersion), blob.Namespace)
	if err != nil {
		return err
	}

	// First share
	b, err := NewBuilder(blobNamespace, uint8(blob.ShareVersion), true).Init()
	if err != nil {
		return err
	}
	if err := b.WriteSequenceLen(uint32(len(rawData))); err != nil {
		return err
	}

	for rawData != nil {

		rawDataLeftOver := b.AddData(rawData)
		if rawDataLeftOver == nil {
			// Just call it on the latest share
			b.ZeroPadIfNecessary()
		}

		share, err := b.Build()
		if err != nil {
			return err
		}
		sss.shares = append(sss.shares, *share)

		b, err = NewBuilder(blobNamespace, uint8(blob.ShareVersion), false).Init()
		if err != nil {
			return err
		}
		rawData = rawDataLeftOver
	}

	return nil
}

// WriteNamespacePaddingShares adds padding shares with the namespace of the
// last written share. This is useful to follow the non-interactive default
// rules. This function assumes that at least one share has already been
// written.
func (sss *SparseShareSplitter) WriteNamespacePaddingShares(count int) error {
	if count < 0 {
		return errors.New("cannot write negative namespaced shares")
	}
	if count == 0 {
		return nil
	}
	if len(sss.shares) == 0 {
		return errors.New("cannot write namespace padding shares on an empty SparseShareSplitter")
	}
	lastBlob := sss.shares[len(sss.shares)-1]
	lastBlobNs, err := lastBlob.Namespace()
	if err != nil {
		return err
	}
	nsPaddingShares, err := NamespacePaddingShares(lastBlobNs, count)
	if err != nil {
		return err
	}
	sss.shares = append(sss.shares, nsPaddingShares...)

	return nil
}

// Export finalizes and returns the underlying shares.
func (sss *SparseShareSplitter) Export() []share.Share {
	return sss.shares
}

// Count returns the current number of shares that will be made if exporting.
func (sss *SparseShareSplitter) Count() int {
	return len(sss.shares)
}

// SplitBlobs splits the provided blobs into shares.
func SplitBlobs(blobs ...Blob) ([]share.Share, error) {
	writer := NewSparseShareSplitter()
	for _, blob := range blobs {
		if err := writer.Write(blob); err != nil {
			return nil, err
		}
	}
	return writer.Export(), nil
}

// CreateCommitment generates the share commitment for a given blob.
// See [data square layout rationale] and [blob share commitment rules].
//
// [data square layout rationale]: ../../specs/src/specs/data_square_layout.md
// [blob share commitment rules]: ../../specs/src/specs/data_square_layout.md#blob-share-commitment-rules
func CreateCommitment(blob *Blob) ([]byte, error) {
	shares, err := SplitBlobs(*blob)
	if err != nil {
		return nil, err
	}

	// the commitment is the root of a merkle mountain range with max tree size
	// determined by the number of roots required to create a share commitment
	// over that blob. The size of the tree is only increased if the number of
	// subtree roots surpasses a constant threshold.
	subTreeWidth := SubTreeWidth(len(shares), DefaultSubtreeRootThreshold)
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
			treeSize, err := RoundDownPowerOfTwo(totalSize)
			if err != nil {
				return treeSizes, err
			}
			treeSizes = append(treeSizes, treeSize)
			totalSize = totalSize - treeSize
		}
	}

	return treeSizes, nil
}
