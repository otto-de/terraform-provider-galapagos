package credentials

import (
	"context"
	"net/http"
)

type ClientWithCredentials interface {
	Client(ctx context.Context) *http.Client
}

type boundClient struct {
	client *http.Client
}

func (b *boundClient) Client(context.Context) *http.Client {
	return b.client
}

func WithClient(client *http.Client) ClientWithCredentials {
	return &boundClient{
		client: client,
	}
}
