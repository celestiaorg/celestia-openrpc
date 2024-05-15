package header

import (
	"context"

	libhead "github.com/celestiaorg/go-header"
	"github.com/celestiaorg/go-header/sync"
)

type API struct {
	LocalHead func(context.Context) (*ExtendedHeader, error) `perm:"read"`
	GetByHash func(
		ctx context.Context,
		hash libhead.Hash,
	) (*ExtendedHeader, error) `perm:"read"`
	GetRangeByHeight func(
		context.Context,
		*ExtendedHeader,
		uint64,
	) ([]*ExtendedHeader, error) `perm:"read"`
	GetByHeight   func(context.Context, uint64) (*ExtendedHeader, error)    `perm:"read"`
	WaitForHeight func(context.Context, uint64) (*ExtendedHeader, error)    `perm:"read"`
	SyncState     func(ctx context.Context) (sync.State, error)             `perm:"read"`
	SyncWait      func(ctx context.Context) error                           `perm:"read"`
	NetworkHead   func(ctx context.Context) (*ExtendedHeader, error)        `perm:"read"`
	Subscribe     func(ctx context.Context) (<-chan *ExtendedHeader, error) `perm:"read"`
}
