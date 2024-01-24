package node

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc/auth"
)

type API struct {
	Info        func(context.Context) (Info, error)                                `perm:"admin"`
	LogLevelSet func(ctx context.Context, name, level string) error                `perm:"admin"`
	AuthVerify  func(ctx context.Context, token string) ([]auth.Permission, error) `perm:"admin"`
	AuthNew     func(ctx context.Context, perms []auth.Permission) ([]byte, error) `perm:"admin"`
}
