package blob

import "github.com/celestiaorg/celestia-openrpc/types/share"

// SplitBlobs splits the provided blobs into shares.
func SplitBlobs(blobs ...Blob) ([]share.Share, error) {
	writer := share.NewSparseShareSplitter()
	for _, blob := range blobs {
		if err := writer.Write(blob.NamespaceVersion, blob.Namespace, blob.Data); err != nil {
			return nil, err
		}
	}
	return writer.Export(), nil
}
