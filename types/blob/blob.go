package blob

import (
	"bytes"
	"encoding/json"

	rpcproto "github.com/rollkit/celestia-openrpc/proto/blob/rollkit"
	"github.com/rollkit/celestia-openrpc/types/namespace"
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
