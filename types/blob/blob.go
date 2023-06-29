package blob

import (
	"encoding/json"

	"cosmossdk.io/errors"
	"github.com/celestiaorg/nmt/namespace"

	appns "github.com/rollkit/celestia-openrpc/types/namespace"
)

var (
	ErrZeroBlobSize = errors.Register("blob", 11124, "cannot use zero blob size")
)

type Commitment []byte

type Proof struct {
	Start uint32   `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	End   uint32   `protobuf:"varint,2,opt,name=end,proto3" json:"end,omitempty"`
	Nodes [][]byte `protobuf:"bytes,3,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

type Blob struct {
	Namespace        []byte `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Data             []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	ShareVersion     uint32 `protobuf:"varint,3,opt,name=share_version,json=shareVersion,proto3" json:"share_version,omitempty"`
	NamespaceVersion uint32 `protobuf:"varint,4,opt,name=namespace_version,json=namespaceVersion,proto3" json:"namespace_version,omitempty"`
	Commitment       []byte `protobuf:"bytes,5,opt,name=commitment,proto3" json:"commitment,omitempty"`
}

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
	blob := &Blob{
		Namespace:        b.GetNamespace(),
		Data:             b.Data,
		ShareVersion:     b.ShareVersion,
		NamespaceVersion: b.NamespaceVersion,
		Commitment:       b.Commitment,
	}
	return json.Marshal(blob)
}

func (b *Blob) UnmarshalJSON(data []byte) error {
	var blob Blob
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
