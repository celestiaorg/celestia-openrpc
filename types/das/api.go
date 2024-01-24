package das

import (
	"context"
)

type API struct {
	SamplingStats func(ctx context.Context) (SamplingStats, error) `perm:"read"`
	WaitCatchUp   func(ctx context.Context) error                      `perm:"read"`
}