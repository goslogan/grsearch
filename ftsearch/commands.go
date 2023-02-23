package ftsearch

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// FTDropIndex removes an index, optionally dropping documents in the index.
func (c *Client) FTDropIndex(ctx context.Context, index string, dropDocuments bool) *redis.BoolCmd {
	args := []interface{}{"ft.dropindex", index}
	if dropDocuments {
		args = append(args, "DD")
	}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}

// FTCreateIndex creates a new index.
func (c *Client) FTCreateIndex(ctx context.Context, index string, options *IndexOptions) *redis.BoolCmd {
	args := []interface{}{"ft.create", index}
	args = append(args, options.serialize()...)
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}

// FTSearch queries an index.
func (c *Client) FTSearch(ctx context.Context, index string, query string, qryOptions *QueryOptions) *QueryCmd {
	args := []interface{}{"ft.search", index, query}
	if qryOptions == nil {
		qryOptions = NewQueryOptions()
	}
	args = append(args, qryOptions.serialize()...)

	cmd := NewQueryCmd(ctx, args...)
	qryOptions.setResultSize()
	cmd.options = qryOptions

	if err := c.Process(ctx, cmd); err == nil {
		cmd.parseResult()
	}

	return cmd
}

// FTConfigGet retrieves public config info from the search config
func (c *Client) FTConfigGet(ctx context.Context, keys ...string) *ConfigGetCmd {
	args := make([]interface{}, len(keys)+2)
	args[0] = "ft.config"
	args[1] = "get"
	for n, arg := range keys {
		args[n+2] = arg
	}

	if len(keys) == 0 {
		args = append(args, "*")
	}

	cmd := NewConfigGetCmd(ctx, args...)
	if err := c.Process(ctx, cmd); err == nil {
		cmd.parseResult()
	}

	return cmd
}

// FTConfigGet sets values in the search config
func (c *Client) FTConfigSet(ctx context.Context, name, value string) *redis.BoolCmd {
	args := []interface{}{"ft.config", "set", name, value}

	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}

// FTTagVals returns the distinct values for a given tag
func (c *Client) FTTagVals(ctx context.Context, index, tag string) *redis.StringSliceCmd {
	args := []interface{}{"ft.tagvals", index, tag}

	cmd := redis.NewStringSliceCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}