package credentials

import (
	"context"
	"net/http"
)

type ClientWithCredentials interface {
	Client(ctx context.Context) *http.Client
}
