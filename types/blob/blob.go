package blob

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/celestiaorg/nmt"
	"github.com/celestiaorg/nmt/namespace"
	rpcproto "github.com/rollkit/celestia-openrpc/proto/celestia/openrpc"
)

const (
	// NMTIgnoreMaxNamespace is currently used value for IgnoreMaxNamespace option in NMT.
	// IgnoreMaxNamespace defines whether the largest possible namespace.ID MAX_NID should be 'ignored'.
	// If set to true, this allows for shorter proofs in particular use-cases.
	NMTIgnoreMaxNamespace = true
)

var (
	ErrBlobNotFound = errors.New("blob: not found")
	ErrInvalidProof = errors.New("blob: invalid proof")
)

// Proof is a collection of nmt.Proofs that verifies the inclusion of the data.
type Proof []*nmt.Proof

func (p Proof) Len() int { return len(p) }

// equal is a temporary method that compares two proofs.
// should be removed in BlobService V1.
func (p Proof) equal(input Proof) error {
	if p.Len() != input.Len() {
		return ErrInvalidProof
	}

	for i, proof := range p {
		pNodes := proof.Nodes()
		inputNodes := input[i].Nodes()
		for i, node := range pNodes {
			if !bytes.Equal(node, inputNodes[i]) {
				return ErrInvalidProof
			}
		}

		if proof.Start() != input[i].Start() || proof.End() != input[i].End() {
			return ErrInvalidProof
		}

		if !bytes.Equal(proof.LeafHash(), input[i].LeafHash()) {
			return ErrInvalidProof
		}

	}
	return nil
}

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
		nmtProof := nmt.NewInclusionProof(jProof.Start, jProof.End, jProof.Nodes, NMTIgnoreMaxNamespace) // ipld.NMTIgnoreMaxNamespace
		nmtProofs[i] = &nmtProof
	}

	*p = nmtProofs
	return nil
}

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

// Blob represents any application-specific binary data that anyone can submit to Celestia.
type Blob struct {
	rpcproto.Blob `json:"blob"`

	Commitment Commitment `json:"commitment"`
}

// Namespace returns blob's namespace.
func (b *Blob) Namespace() namespace.ID {
	return append([]byte{uint8(b.NamespaceVersion)}, b.NamespaceId...)
}

type jsonBlob struct {
	Namespace    namespace.ID `json:"namespace"`
	Data         []byte       `json:"data"`
	ShareVersion uint32       `json:"share_version"`
	Commitment   Commitment   `json:"commitment"`
}

func (b *Blob) MarshalJSON() ([]byte, error) {
	blob := &jsonBlob{
		Namespace:    b.Namespace(),
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

	b.Blob.NamespaceVersion = uint32(blob.Namespace[0])
	b.Blob.NamespaceId = blob.Namespace[1:]
	b.Blob.Data = blob.Data
	b.Blob.ShareVersion = blob.ShareVersion
	b.Commitment = blob.Commitment
	return nil
}
