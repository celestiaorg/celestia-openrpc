package node

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc/auth"
)

type API struct {
	// Info returns administrative information about the node.
	Info func(context.Context) (Info, error) `perm:"admin"`
	// Ready returns true once the node's RPC is ready to accept requests.
	Ready func(context.Context) (bool, error)
	// LogLevelSet sets the given component log level to the given level.
	LogLevelSet func(ctx context.Context, name, level string) error `perm:"admin"`
	// AuthVerify returns the permissions assigned to the given token.
	AuthVerify func(ctx context.Context, token string) ([]auth.Permission, error) `perm:"admin"`
	// AuthNew signs and returns a new token with the given permissions.
	AuthNew func(ctx context.Context, perms []auth.Permission) ([]byte, error) `perm:"admin"`
}
