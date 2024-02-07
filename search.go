// grsearch main module - defines the client class
package grsearch

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SearchCmdAble interface {
	FTSearch(ctx context.Context, index string, query string, options *QueryOptions) *QueryCmd
	FTAggregate(ctx context.Context, index string, query string, options *AggregateOptions) *QueryCmd
	FTDropIndex(ctx context.Context, index string, dropDocuments bool) *redis.BoolCmd
	FTCreateIndex(ctx context.Context, index string)
	FTConfigGet(ctx context.Context, keys ...string) *ConfigGetCmd
	FTConfigSet(ctx context.Context, name, value string) *redis.BoolCmd
	FTTagVals(ctx context.Context, index, tag string) *redis.StringSliceCmd
	FTList(ctx context.Context) *redis.StringSliceCmd
	FTInfo(ctx context.Context, index string) *InfoCmd
	FTDictAdd(ctx context.Context, dictionary string, terms ...string) *redis.IntCmd
	FTDictDel(ctx context.Context, dictionary string, terms ...string) *redis.IntCmd
	FTDictDump(ctx context.Context, dictionary string) *redis.StringSliceCmd
	FTSynUpdate(ctx context.Context, index string, group string, terms ...string) *redis.BoolCmd
	FTSynDump(ctx context.Context, index string) *SynonymDumpCmd
	FTAliasAdd(ctx context.Context, alias, index string) *redis.BoolCmd
	FTAliasDel(ctx context.Context, alias string) *redis.BoolCmd
	FTAliasUpdate(ctx context.Context, alias, index string) *redis.BoolCmd
}
