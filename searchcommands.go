package grstack

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// FTDropIndex removes an index, optionally dropping documents in the index.
func (c cmdable) FTDropIndex(ctx context.Context, index string, dropDocuments bool) *redis.BoolCmd {
	args := []interface{}{"FT.DROPINDEX", index}
	if dropDocuments {
		args = append(args, "DD")
	}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTCreate creates a new index.
func (c cmdable) FTCreate(ctx context.Context, index string, options *IndexOptions) *redis.BoolCmd {
	args := []interface{}{"FT.CREATE", index}
	args = append(args, options.serialize()...)
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTAggregate runs a search query on an index, and perform saggregate transformations on the results, extracting statistics etc from them
func (c cmdable) FTAggregate(ctx context.Context, index, query string, options *AggregateOptions) *AggregateCmd {
	args := []interface{}{"FT.AGGREGATE", index, query}
	args = append(args, options.serialize()...)
	cmd := NewAggregateCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTSearch queries an index (on hashes)
func (c cmdable) FTSearch(ctx context.Context, index string, query string, qryOptions *QueryOptions) *QueryCmd {
	args := []interface{}{"FT.SEARCH", index, query}
	if qryOptions == nil {
		qryOptions = NewQueryOptions()
	}
	args = append(args, qryOptions.serialize()...)

	cmd := NewQueryCmd(ctx, c, args...)
	cmd.options = qryOptions

	_ = c(ctx, cmd)
	return cmd
}

// FTSearch queries an index on JSON documents
func (c cmdable) FTSearchJSON(ctx context.Context, index string, query string, qryOptions *QueryOptions) *QueryCmd {
	args := []interface{}{"FT.SEARCH", index, query}
	if qryOptions == nil {
		qryOptions = NewQueryOptions()
	}
	qryOptions.json = true
	args = append(args, qryOptions.serialize()...)

	cmd := NewQueryCmd(ctx, c, args...)
	cmd.options = qryOptions

	_ = c(ctx, cmd)
	return cmd
}

// FTConfigGet retrieves public config info from the search config
func (c cmdable) FTConfigGet(ctx context.Context, keys ...string) *ConfigGetCmd {
	args := make([]interface{}, len(keys)+2)
	args[0] = "FT.CONFIG"
	args[1] = "GET"
	for n, arg := range keys {
		args[n+2] = arg
	}

	if len(keys) == 0 {
		args = append(args, "*")
	}

	cmd := NewConfigGetCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTConfigGet sets values in the search config
func (c cmdable) FTConfigSet(ctx context.Context, name, value string) *redis.BoolCmd {
	args := []interface{}{"FT.CONFIG", "SET", name, value}

	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTTagVals returns the distinct values for a given tag
func (c cmdable) FTTagVals(ctx context.Context, index, tag string) *redis.StringSliceCmd {
	args := []interface{}{"FT.TAGVALS", index, tag}

	cmd := redis.NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTList returns a list of all the indexes currently defined
func (c cmdable) FTList(ctx context.Context) *redis.StringSliceCmd {
	cmd := redis.NewStringSliceCmd(ctx)
	_ = c(ctx, cmd)
	return cmd
}

// FTInfo returns information about an index
func (c cmdable) FTInfo(ctx context.Context, index string) *InfoCmd {
	args := []interface{}{"FT.INFO", index}
	cmd := NewInfoCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

/*******************************************************************************
*
* DICTIONARIES
*
*******************************************************************************/

// FTDictAdd adds one more terms to a dictionary
func (c cmdable) FTDictAdd(ctx context.Context, dictionary string, terms ...string) *redis.IntCmd {

	args := make([]interface{}, len(terms)+2)
	args[0] = "FT.DICTADD"
	args[1] = dictionary
	for n, term := range terms {
		args[n+2] = term
	}

	cmd := redis.NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)

	return cmd

}

// FTDictDel removes terms from a dictionary
func (c cmdable) FTDictDel(ctx context.Context, dictionary string, terms ...string) *redis.IntCmd {

	args := make([]interface{}, len(terms)+2)
	args[0] = "FT.DICTDEL"
	args[1] = dictionary
	for n, term := range terms {
		args[n+2] = term
	}

	cmd := redis.NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)

	return cmd
}

// FTDictDump returns a slice containing all the terms in a dictionary
func (c cmdable) FTDictDump(ctx context.Context, dictionary string) *redis.StringSliceCmd {

	args := []interface{}{"FT.DICTDUMP", dictionary}

	cmd := redis.NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)

	return cmd
}

/*******************************************************************************
*
* SYNONYMS
*
*******************************************************************************/

// FTSynUpdate adds to or modifies a synonym group
func (c cmdable) FTSynUpdate(ctx context.Context, index string, group string, terms ...string) *redis.BoolCmd {
	args := make([]interface{}, len(terms)+3)
	args[0] = "FT.SYNUPDATE"
	args[1] = index
	args[2] = group
	for n, term := range terms {
		args[n+2] = term
	}

	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)

	return cmd
}

// FTSynDump returns the contents of synonym map for an index
func (c cmdable) FTSynDump(ctx context.Context, index string) *SynonymDumpCmd {
	args := []interface{}{"FT.SYNDUMP", index}
	cmd := NewSynonymDumpCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

/*******************************************************************************
*
* ALIASES
*
*******************************************************************************/

// FTAliasAdd add an alias to an index.
func (c cmdable) FTAliasAdd(ctx context.Context, alias, index string) *redis.BoolCmd {
	args := []interface{}{"FT.ALIASADD", alias, index}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTAliasDel deletes an alias
func (c cmdable) FTAliasDel(ctx context.Context, alias string) *redis.BoolCmd {
	args := []interface{}{"FT.ALIASDEL", alias}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// FTAliasDel deletes an alias
func (c cmdable) FTAliasUpdate(ctx context.Context, alias, index string) *redis.BoolCmd {
	args := []interface{}{"FT.ALIASUPDATE", alias, index}
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}
