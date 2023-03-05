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

// FTList returns a list of all the indexes currently defined
func (c *Client) FTList(ctx context.Context) *redis.StringSliceCmd {
	cmd := redis.NewStringSliceCmd(ctx)
	_ = c.Process(ctx, cmd)
	return cmd
}

// FTInfo returns information about an index
func (c *Client) FTInfo(ctx context.Context, index string) *InfoCmd {
	args := []interface{}{"ft.info", index}
	cmd := NewInfoCmd(ctx, args)

	if err := c.Process(ctx, cmd); err == nil {
		cmd.parseResult()
	}

	return cmd
}

/*******************************************************************************
*
* DICTIONARIES
*
*******************************************************************************/

// FTDictAdd adds one more terms to a dictionary
func (c *Client) FTDictAdd(ctx context.Context, dictionary string, terms ...string) *redis.IntCmd {

	args := make([]interface{}, len(terms)+2)
	args[0] = "ft.dictadd"
	args[1] = dictionary
	for n, term := range terms {
		args[n+2] = term
	}

	cmd := redis.NewIntCmd(ctx, args...)
	_ = c.Process(ctx, cmd)

	return cmd

}

// FTDictDel removes terms from a dictionary
func (c *Client) FTDictDel(ctx context.Context, dictionary string, terms ...string) *redis.IntCmd {

	args := make([]interface{}, len(terms)+2)
	args[0] = "ft.dictdel"
	args[1] = dictionary
	for n, term := range terms {
		args[n+2] = term
	}

	cmd := redis.NewIntCmd(ctx, args...)
	_ = c.Process(ctx, cmd)

	return cmd
}

// FTDictDump returns a slice containing all the terms in a dictionary
func (c *Client) FTDictDump(ctx context.Context, dictionary string) *redis.StringSliceCmd {

	args := []interface{}{"ft.dictdump", dictionary}

	cmd := redis.NewStringSliceCmd(ctx, args...)
	_ = c.Process(ctx, cmd)

	return cmd
}

/*******************************************************************************
*
* SYNONYMS
*
*******************************************************************************/

// FTSynUpdate adds to or modifies a synonym group
func (c *Client) FTSynUpdate(ctx context.Context, index string, group string, terms ...string) *redis.BoolCmd {
	args := make([]interface{}, len(terms)+3)
	args[0] = "ft.synupdate"
	args[1] = index
	args[2] = group
	for n, term := range terms {
		args[n+2] = term
	}

	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)

	return cmd
}

// FTSynDump returns the contents of synonym map for an index
func (c *Client) FTSynDump(ctx context.Context, index string) *SynonymDumpCmd {
	args := []interface{}{"ft.syndump", index}
	cmd := NewSynonymDumpCmd(ctx, args...)
	if err := c.Process(ctx, cmd); err == nil {
		cmd.parseResult()
	}
	return cmd
}

/*******************************************************************************
*
* ALIASES
*
*******************************************************************************/

// FTAliasAdd add an alias to an index.
func (c *Client) FTAliasAdd(ctx context.Context, alias, index string) *redis.BoolCmd {
	args := []interface{}{"ft.aliasadd", alias, index}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}

// FTAliasDel deletes an alias
func (c *Client) FTAliasDel(ctx context.Context, alias string) *redis.BoolCmd {
	args := []interface{}{"ft.aliasdel", alias}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}

// FTAliasDel deletes an alias
func (c *Client) FTAliasUpdate(ctx context.Context, alias, index string) *redis.BoolCmd {
	args := []interface{}{"ft.aliasupdate", alias, index}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}
