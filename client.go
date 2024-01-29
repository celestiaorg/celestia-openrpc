package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	clientbuilder "github.com/rollkit/celestia-openrpc/builder"
	"github.com/rollkit/celestia-openrpc/types/blob"
	"github.com/rollkit/celestia-openrpc/types/das"
	"github.com/rollkit/celestia-openrpc/types/fraud"
	"github.com/rollkit/celestia-openrpc/types/header"
	"github.com/rollkit/celestia-openrpc/types/node"
	"github.com/rollkit/celestia-openrpc/types/p2p"
	"github.com/rollkit/celestia-openrpc/types/share"
	"github.com/rollkit/celestia-openrpc/types/state"
)

const AuthKey = "Authorization"

type Client struct {
	Fraud  fraud.API
	Blob   blob.API
	Header header.API
	State  state.API
	Share  share.API
	DAS    das.API
	P2P    p2p.API
	Node   node.API

	closer clientbuilder.MultiClientCloser
}

// Close closes the connections to all namespaces registered on the client.
func (c *Client) Close() {
	c.closer.CloseAll()
}

func NewClient(ctx context.Context, addr string, token string) (*Client, error) {
	var authHeader http.Header
	if token != "" {
		authHeader = http.Header{AuthKey: []string{fmt.Sprintf("Bearer %s", token)}}
	}

	var client Client

	modules := map[string]interface{}{
		"fraud":  &client.Fraud,
		"blob":   &client.Blob,
		"header": &client.Header,
		"state":  &client.State,
		"share":  &client.Share,
		"das":    &client.DAS,
		"p2p":    &client.P2P,
		"node":   &client.Node,
	}

	for name, module := range modules {
		closer, err := jsonrpc.NewClient(ctx, addr, name, module, authHeader)
		if err != nil {
			return nil, err
		}
		client.closer.Register(closer)
	}

	return &client, nil
}
