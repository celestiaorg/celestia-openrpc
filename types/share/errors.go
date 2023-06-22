package share

import "errors"

var (
	// ErrNotAvailable is returned whenever DA sampling fails.
	ErrNotAvailable = errors.New("share: data not available")
)
