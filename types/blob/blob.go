package blob

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/celestiaorg/nmt"

	"github.com/celestiaorg/celestia-openrpc/types/appconsts"
	"github.com/celestiaorg/celestia-openrpc/types/share"

	"github.com/celestiaorg/go-square/blob"
	"github.com/celestiaorg/go-square/inclusion"
	"github.com/celestiaorg/go-square/merkle"
)

const (
	// NMTIgnoreMaxNamespace is currently used value for IgnoreMaxNamespace option in NMT.
	// IgnoreMaxNamespace defines whether the largest possible Namespace MAX_NID should be 'ignored'.
	// If set to true, this allows for shorter proofs in particular use-cases.
	NMTIgnoreMaxNamespace = true
)

var (
	ErrBlobNotFound = errors.New("blob: not found")
	ErrInvalidProof = errors.New("blob: invalid proof")
)

// Commitment is a Merkle Root of the subtree built from shares of the Blob.
// It is computed by splitting the blob into shares and building the Merkle subtree to be included
// after Submit.
type Commitment []byte

func (com Commitment) String() string {
	return string(com)
}

// Equal ensures that commitments are the same
func (com Commitment) Equal(c Commitment) bool {
	return bytes.Equal(com, c)
}

// Proof is a collection of nmt.Proofs that verifies the inclusion of the data.
// Proof proves the WHOLE namespaced data for the particular row.
// TODO (@vgonkivs): rework `Proof` in order to prove a particular blob.
// https://github.com/celestiaorg/celestia-node/issues/2303
type Proof []*nmt.Proof

func (p Proof) Len() int { return len(p) }

// Blob represents any application-specific binary data that anyone can submit to Celestia.
type Blob struct {
	blob.Blob `json:"blob"`

	Commitment Commitment `json:"commitment"`

	// the celestia-node's namespace type
	// this is to avoid converting to and from app's type
	namespace share.Namespace

	// index represents the index of the blob's first share in the EDS.
	// Only retrieved, on-chain blobs will have the index set. Default is -1.
	index int
}

// NewBlobV0 constructs a new blob from the provided Namespace and data.
// The blob will be formatted as v0 shares.
func NewBlobV0(namespace share.Namespace, data []byte) (*Blob, error) {
	return NewBlob(appconsts.ShareVersionZero, namespace, data)
}

// NewBlob constructs a new blob from the provided Namespace, data and share version.
func NewBlob(shareVersion uint8, namespace share.Namespace, data []byte) (*Blob, error) {
	if len(data) == 0 || len(data) > appconsts.DefaultMaxBytes {
		return nil, fmt.Errorf("blob data must be > 0 && <= %d, but it was %d bytes", appconsts.DefaultMaxBytes, len(data))
	}
	if err := namespace.ValidateForBlob(); err != nil {
		return nil, err
	}

	blob := blob.Blob{
		NamespaceId:      namespace.ID(),
		Data:             data,
		ShareVersion:     uint32(shareVersion),
		NamespaceVersion: uint32(namespace.Version()),
	}

	com, err := inclusion.CreateCommitment(&blob, merkle.HashFromByteSlices, appconsts.DefaultSubtreeRootThreshold)
	if err != nil {
		return nil, err
	}
	//nolint:govet
	return &Blob{Blob: blob, Commitment: com, namespace: namespace, index: -1}, nil
}

// DefaultGasPrice returns the default gas price, letting node automatically
// determine the Fee based on the passed blob sizes.
func DefaultGasPrice() float64 {
	return -1.0
}

type jsonBlob struct {
	Namespace    share.Namespace `json:"namespace"`
	Data         []byte          `json:"data"`
	ShareVersion uint32          `json:"share_version"`
	Commitment   Commitment      `json:"commitment"`
	Index        int             `json:"index"`
}

func (b *Blob) MarshalJSON() ([]byte, error) {
	ns, err := share.NamespaceFromBytes(b.Namespace().Bytes())
	if err != nil {
		return nil, err
	}
	blob := &jsonBlob{
		Namespace:    ns,
		Data:         b.Data,
		ShareVersion: b.ShareVersion,
		Commitment:   b.Commitment,
		Index:        b.index,
	}
	return json.Marshal(blob)
}

func (b *Blob) UnmarshalJSON(data []byte) error {
	var blob jsonBlob
	err := json.Unmarshal(data, &blob)
	if err != nil {
		return err
	}

	b.Blob.NamespaceVersion = uint32(blob.Namespace.Version())
	b.Blob.NamespaceId = blob.Namespace.ID()
	b.Blob.Data = blob.Data
	b.Blob.ShareVersion = blob.ShareVersion
	b.Commitment = blob.Commitment
	b.namespace = blob.Namespace
	b.index = blob.Index
	return nil
}
