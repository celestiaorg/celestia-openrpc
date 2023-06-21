package blob

import (
	"encoding/json"

	"cosmossdk.io/errors"
	"github.com/celestiaorg/nmt/namespace"
	rpcproto "github.com/rollkit/celestia-openrpc/proto/blob"
	appns "github.com/rollkit/celestia-openrpc/types/namespace"
)

var (
	ErrZeroBlobSize = errors.Register("blob", 11124, "cannot use zero blob size")
)

type Commitment []byte

type Proof rpcproto.Proof

type Blob rpcproto.Blob

// NewBlob creates a new coretypes.Blob from the provided data after performing
// basic stateless checks over it.
func NewBlob(ns appns.Namespace, blob []byte, shareVersion uint8) (*Blob, error) {
	err := ns.ValidateBlobNamespace()
	if err != nil {
		return nil, err
	}

	if len(blob) == 0 {
		return nil, ErrZeroBlobSize
	}

	return &Blob{
		Namespace:        ns.ID,
		Data:             blob,
		ShareVersion:     uint32(shareVersion),
		NamespaceVersion: uint32(ns.Version),
	}, nil
}

// Renamed to avoid conflict with the proto method
func (b *Blob) GetNamespace() namespace.ID {
	return append([]byte{uint8(b.NamespaceVersion)}, b.Namespace...)
}

func (b *Blob) MarshalJSON() ([]byte, error) {
	blob := &rpcproto.Blob{
		Namespace:        b.GetNamespace(),
		Data:             b.Data,
		ShareVersion:     b.ShareVersion,
		NamespaceVersion: b.NamespaceVersion,
		Commitment:       b.Commitment,
	}
	return json.Marshal(blob)
}

func (b *Blob) UnmarshalJSON(data []byte) error {
	var blob rpcproto.Blob
	err := json.Unmarshal(data, &blob)
	if err != nil {
		return err
	}

	b.NamespaceVersion = uint32(blob.Namespace[0])
	b.Namespace = blob.Namespace[1:]
	b.Data = blob.Data
	b.ShareVersion = blob.ShareVersion
	b.Commitment = blob.Commitment
	return nil
}
