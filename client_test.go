package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewClient(context.Background(), "http://localhost:26658", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.aBWqglHA-R1u4X1In5HMAqX88V5nDetjA6KflxB0p9U")
	defer client.Close()

	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx, closer := context.WithTimeout(context.Background(), 1*time.Second)
	defer closer()

	resp := client.Share.ProbabilityOfAvailability(ctx)
	assert.NotZero(t, resp)

	info, err := client.Node.Info(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, info.APIVersion)
}
