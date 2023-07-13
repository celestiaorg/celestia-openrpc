package blob

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/celestiaorg/nmt"

	"github.com/rollkit/celestia-openrpc/types/appconsts"
	"github.com/rollkit/celestia-openrpc/types/share"
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
type Proof []*nmt.Proof

func (p Proof) Len() int { return len(p) }

type jsonProof struct {
	Start int      `json:"start"`
	End   int      `json:"end"`
	Nodes [][]byte `json:"nodes"`
}

func (p *Proof) MarshalJSON() ([]byte, error) {
	proofs := make([]jsonProof, 0, p.Len())
	for _, pp := range *p {
		proofs = append(proofs, jsonProof{
			Start: pp.Start(),
			End:   pp.End(),
			Nodes: pp.Nodes(),
		})
	}

	return json.Marshal(proofs)
}

func (p *Proof) UnmarshalJSON(data []byte) error {
	var proofs []jsonProof
	err := json.Unmarshal(data, &proofs)
	if err != nil {
		return err
	}

	nmtProofs := make([]*nmt.Proof, len(proofs))
	for i, jProof := range proofs {
		nmtProof := nmt.NewInclusionProof(jProof.Start, jProof.End, jProof.Nodes, NMTIgnoreMaxNamespace)
		nmtProofs[i] = &nmtProof
	}

	*p = nmtProofs
	return nil
}

// Blob represents any application-specific binary data that anyone can submit to Celestia.
type Blob struct {
	Namespace        []byte `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Data             []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	ShareVersion     uint32 `protobuf:"varint,3,opt,name=share_version,json=shareVersion,proto3" json:"share_version,omitempty"`
	NamespaceVersion uint32 `protobuf:"varint,4,opt,name=namespace_version,json=namespaceVersion,proto3" json:"namespace_version,omitempty"`
	Commitment       []byte `protobuf:"bytes,5,opt,name=commitment,proto3" json:"commitment,omitempty"`
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

	_, _ = CreateCommitment(nil)

	return &Blob{
		Namespace:        namespace,
		Data:             data,
		ShareVersion:     uint32(shareVersion),
		NamespaceVersion: 0,
		Commitment:       []byte{},
	}, nil
}

type jsonBlob struct {
	Namespace    share.Namespace `json:"namespace"`
	Data         []byte          `json:"data"`
	ShareVersion uint32          `json:"share_version"`
	Commitment   Commitment      `json:"commitment"`
}

func (b *Blob) MarshalJSON() ([]byte, error) {
	blob := &jsonBlob{
		Namespace:    b.Namespace,
		Data:         b.Data,
		ShareVersion: b.ShareVersion,
		Commitment:   b.Commitment,
	}
	return json.Marshal(blob)
}

func (b *Blob) UnmarshalJSON(data []byte) error {
	var blob jsonBlob
	err := json.Unmarshal(data, &blob)
	if err != nil {
		return err
	}

	b.NamespaceVersion = uint32(blob.Namespace.Version())
	b.Data = blob.Data
	b.ShareVersion = blob.ShareVersion
	b.Commitment = blob.Commitment
	b.Namespace = blob.Namespace
	return nil
}
