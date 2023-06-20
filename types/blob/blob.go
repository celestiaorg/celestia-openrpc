package blob

import (
	"encoding/json"

	"github.com/celestiaorg/nmt/namespace"
	rpcproto "github.com/rollkit/celestia-openrpc/proto/celestia/openrpc"
)

type Blob rpcproto.Blob

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
