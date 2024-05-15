package header

import (
	"context"

	libhead "github.com/celestiaorg/go-header"
	"github.com/celestiaorg/go-header/sync"
)

type API struct {
	// LocalHead returns the ExtendedHeader of the chain head.
	LocalHead func(context.Context) (*ExtendedHeader, error) `perm:"read"`
	// GetByHash returns the header of the given hash from the node's header store.
	GetByHash func(
		ctx context.Context,
		hash libhead.Hash,
	) (*ExtendedHeader, error) `perm:"read"`
	// GetRangeByHeight returns the given range (from:to) of ExtendedHeaders
	// from the node's header store and verifies that the returned headers are
	// adjacent to each other.
	GetRangeByHeight func(
		context.Context,
		*ExtendedHeader,
		uint64,
	) ([]*ExtendedHeader, error) `perm:"read"`
	// GetByHeight returns the ExtendedHeader at the given height if it is
	// currently available.
	GetByHeight func(context.Context, uint64) (*ExtendedHeader, error) `perm:"read"`
	// WaitForHeight blocks until the header at the given height has been processed
	// by the store or context deadline is exceeded.
	WaitForHeight func(context.Context, uint64) (*ExtendedHeader, error) `perm:"read"`
	// SyncState returns the current state of the header Syncer.
	SyncState func(ctx context.Context) (sync.State, error) `perm:"read"`
	// SyncWait blocks until the header Syncer is synced to network head.
	SyncWait func(ctx context.Context) error `perm:"read"`
	// NetworkHead provides the Syncer's view of the current network head.
	NetworkHead func(ctx context.Context) (*ExtendedHeader, error) `perm:"read"`
	// Subscribe to recent ExtendedHeaders from the network.
	Subscribe func(ctx context.Context) (<-chan *ExtendedHeader, error) `perm:"read"`
}
