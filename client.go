package grstack

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis.Client
	cmdable
}

type cmdable func(ctx context.Context, cmd redis.Cmder) error

// NewClient returns a new search client
func NewClient(options *redis.Options) *Client {
	client := &Client{Client: *redis.NewClient(options)}
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
