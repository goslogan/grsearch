// Package grsearch is an extension to [github.com/redis/go-redis], implementing support for
// [RedisJSON] and [RediSearch]. It attempts to follow the syntactic style of go-redis as closely as possible.
//
// [RediSearch]: https://redis.io/docs/stack/search/
// [RedisJSON]: https://redis.io/docs/stack/json/
package grsearch

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Client represents the wrapped go-redis client. New clients should
// be created using [NewClient]
type Client struct {
	redis.Client
	cmdable
}

type cmdable func(ctx context.Context, cmd redis.Cmder) error

// NewClient returns a new search client using the same options as the standard
// go-redis client.
func NewClient(options *redis.Options) *Client {
	client := &Client{Client: *redis.NewClient(options)}
	client.cmdable = client.Process
	return client
}

// FromRedisClient builds a client from an existing redis client
func FromRedisCLient(redisClient *redis.Client) *Client {
	client := &Client{Client: *redisClient}
	client.cmdable = client.Process
	return client
}

func (c *Client) Process(ctx context.Context, cmd redis.Cmder) error {
	err := c.Client.Process(ctx, cmd)
	if c, ok := cmd.(ExtCmder); ok {
		err = c.postProcess()
		if err != nil {
			cmd.SetErr(err)
		}
	}

	return err
}
