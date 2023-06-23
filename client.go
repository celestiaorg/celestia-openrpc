package client

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"net/http"
)

const AuthKey = "Authorization"

type Client struct {
	Fraud  FraudAPI
	Header HeaderAPI
	State  StateAPI
	Share  ShareAPI
	DAS    DASAPI
	P2P    P2PAPI
	Node   NodeAPI

	closer multiClientCloser
}

// multiClientCloser is a wrapper struct to close clients across multiple namespaces.
type multiClientCloser struct {
	closers []jsonrpc.ClientCloser
}

// register adds a new closer to the multiClientCloser
func (m *multiClientCloser) register(closer jsonrpc.ClientCloser) {
	m.closers = append(m.closers, closer)
}

// closeAll closes all saved clients.
func (m *multiClientCloser) closeAll() {
	for _, closer := range m.closers {
		closer()
	}
}

// Close closes the connections to all namespaces registered on the client.
func (c *Client) Close() {
	c.closer.closeAll()
}

func NewClient(ctx context.Context, addr string, token string) (*Client, error) {
	authHeader := http.Header{AuthKey: []string{fmt.Sprintf("Bearer %s", token)}}

	var client Client

	modules := map[string]interface{}{
		"fraud":  &client.Fraud,
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
		client.closer.register(closer)
	}

	return &client, nil
}
