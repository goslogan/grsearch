// grstack main module - defines the client class
package grstack

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SearchCmdAble interface {
	redis.Cmdable
	FTSearch(ctx context.Context, index string, query string, options *QueryOptions) *QueryCmd
	DropIndex(ctx context.Context, index string, dropDocuments bool) *redis.BoolCmd
	CreateIndex(ctx context.Context, index string)
}

type Client struct {
	redis.Client
}

// NewClient returns a new search client
func NewClient(options *redis.Options) *Client {
	return &Client{Client: *redis.NewClient(options)}
}
