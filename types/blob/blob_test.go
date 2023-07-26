package blob

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rollkit/celestia-openrpc/types/appconsts"
	appns "github.com/rollkit/celestia-openrpc/types/namespace"
	"github.com/rollkit/celestia-openrpc/types/share"
)

func TestBlobMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		blobJSON string
		blob     *Blob
	}{
		{
			"valid blob",
			`{"namespace":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAQIDBAUGBwg=","data":"aGVsbG8gd29ybGQ=","share_version":0,"commitment":"I6VBbcCIpcliy0hYTCLdX13m18ImVdABclJupNGueko="}`,
			&Blob{
				Namespace:        append(bytes.Repeat([]byte{0x00}, 21), []byte{1, 2, 3, 4, 5, 6, 7, 8}...),
				Data:             []byte("hello world"),
				ShareVersion:     uint32(appconsts.ShareVersionZero),
				NamespaceVersion: uint32(appns.NamespaceVersionZero),
				Commitment:       []byte{0x23, 0xa5, 0x41, 0x6d, 0xc0, 0x88, 0xa5, 0xc9, 0x62, 0xcb, 0x48, 0x58, 0x4c, 0x22, 0xdd, 0x5f, 0x5d, 0xe6, 0xd7, 0xc2, 0x26, 0x55, 0xd0, 0x1, 0x72, 0x52, 0x6e, 0xa4, 0xd1, 0xae, 0x7a, 0x4a},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespace, err := share.NewBlobNamespaceV0([]byte{1, 2, 3, 4, 5, 6, 7, 8})
			require.NoError(t, err)
			require.NotEmpty(t, namespace)

			blob := &Blob{}
			err = blob.UnmarshalJSON([]byte(tt.blobJSON))
			require.NoError(t, err)

			require.Equal(t, tt.blob.ShareVersion, blob.ShareVersion)
			require.Equal(t, tt.blob.NamespaceVersion, blob.NamespaceVersion)
			require.Equal(t, tt.blob.Namespace, blob.Namespace)
			require.Equal(t, tt.blob.Data, blob.Data)
			require.Equal(t, tt.blob.Commitment, blob.Commitment)

			blobJSON, err := blob.MarshalJSON()
			require.NoError(t, err)
			require.Equal(t, tt.blobJSON, string(blobJSON))
		})
	}
}
