package share

import (
	"context"

	"github.com/celestiaorg/celestia-openrpc/types/header"
	"github.com/celestiaorg/rsmt2d"
)

type API struct {
	SharesAvailable func(context.Context, *header.ExtendedHeader) error `perm:"read"`
	GetShare        func(
		ctx context.Context,
		eh *header.ExtendedHeader,
		row, col int,
	) (*Share, error) `perm:"read"`
	GetEDS func(
		ctx context.Context,
		eh *header.ExtendedHeader,
	) (*rsmt2d.ExtendedDataSquare, error) `perm:"read"`
	GetSharesByNamespace func(
		ctx context.Context,
		eh *header.ExtendedHeader,
		namespace Namespace,
	) (*NamespacedShares, error) `perm:"read"`
}
