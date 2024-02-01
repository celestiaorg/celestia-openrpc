package clientbuilder

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/filecoin-project/go-jsonrpc"
)

const AuthKey = "Authorization"

// MultiClientCloser is a wrapper struct to close clients across multiple namespaces.
type MultiClientCloser struct {
	closers []jsonrpc.ClientCloser
}

// Register adds a new closer to the multiClientCloser
func (m *MultiClientCloser) Register(closer jsonrpc.ClientCloser) {
	m.closers = append(m.closers, closer)
}

// CloseAll closes all saved clients.
func (m *MultiClientCloser) CloseAll() {
	for _, closer := range m.closers {
		closer()
	}
}

func NewClient(ctx context.Context, addr string, token string, client interface{}) (interface{}, error) {
	var authHeader http.Header
	if token != "" {
		authHeader = http.Header{AuthKey: []string{fmt.Sprintf("Bearer %s", token)}}
	}

	v := reflect.ValueOf(client).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := strings.ToLower(field.Name)
		module := v.Field(i).Addr().Interface()
		_, err := jsonrpc.NewClient(ctx, addr, name, module, authHeader)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
