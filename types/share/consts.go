package share

import "github.com/celestiaorg/rsmt2d"

var (
	// DefaultRSMT2DCodec is the default codec creator used for data erasure.
	DefaultRSMT2DCodec = rsmt2d.NewLeoRSCodec
)

const (
	// ShareSize is the size of a share in bytes.
	ShareSize = 512
)
