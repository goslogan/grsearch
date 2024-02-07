# grsearch

go-redis based interface to RediSearch & RedisJSON designed to extend the go-redis
API and follow its style as closely as possible.


## Working with search

The **RediSearch** functions are all prefixed with _FT_ and follow the native command syntax as closely as possible although options and parameters for `FT.SEARCH`, `FT.CREATE` and `FT.AGGREGATE` and represented as structs. JSON searching is implemented via the `FTSearchJSON` method as result parsing differs from that needed for hash result parsing. 

Search results are returned in using the `SearchResult` struct. Documents in search results are represented using the `HashQueryResult` and `JSONQueryResult` structs stored in the `Value` property of the results struct.

### Builders

The `IndexBuilder`, `QueryBuilder` and `AggregateBuilder` types provide a fluent interface to the options structs.

```
options := grsearch.NewIndexBuilder().
        On("hash").
		Prefix("account:").
		Schema(grsearch.TagAttribute{
			Name    : "account_id",
			Alias   : "id",
			Sortable: true}).
        Schema(grsearch.TextAttribute{Name: "customer",
		    Sortable: true}).Schema(grsearch.TextAttribute{
		    Name    : "email",
		    Sortable: true}).
        Schema(grsearch.TagAttribute{
		    Name    : "account_owner",
		    Alias   : "owner",
		    Sortable: true}).
        Schema(grsearch.NumericAttribute{
		    Name    : "balance",
		    Sortable: true,
	}).Options()

```

as opposed to


```
options := grsearch.NewIndexOptions()
options.On = "hash"
options.Prefix = []string{"account:"}
options.Schema = []grsearch.SchemaAttribute{
	grsearch.TagAttribute{
		Name    : "account_id",
		Alias   : "id",
		Sortable: true},
    grsearch.TextAttribute{Name: "customer",
		Sortable: true},
	grsearch.TextAttribute{
		Name    : "email",
		Sortable: true},
    grsearch.TagAttribute{
		Name    : "account_owner",
		Alias   : "owner",
		Sortable: true},
	Schema(grsearch.NumericAttribute{
		Name    : "balance",
		Sortable: true,
	}
}
```

### Searching hashes

```
import (
	"context"
	"log"

	"github.com/goslogan/grsearch"
	"github.com/redis/go-redis/v9"
)

client := grsearch.NewClient(&redis.Options{})
ctx := context.Background()

searchResult, err := client.FTSearch(ctx, "customers", "@id:{1128564}", nil).Results()
for id, customer := range searchResult {
	fmt.Println(searchResult.Value[id])
}

```

### Search JSON

JSON searches return a map of `JSONQueryResult`  (keyed by document key name). The Value property is set to the 
result of unmarshalling the string result into a `map[string]interface{}`. 

```

import (
	"context"
	"log"

	"github.com/goslogan/grsearch"
	"github.com/redis/go-redis/v9"
)

client := grsearch.NewClient(&redis.Options{})
ctx := context.Background()

options := grsearch.NewQueryBuilder().
	Return("$..data", "data").
	WithScores().
	Options()

searchResult, err := client.FTSearch(ctx, "jcustomers", "@id:{1128564}", options).Results()
for id, customer := range searchResult {
	fmt.Println(searchResult[id].Value["data"])
}

```

## Working with JSON.


```

import (
	"context"

	"github.com/goslogan/grsearch"
	"github.com/redis/go-redis/v9"
)

import (
	"context"
    "log"

    "github.com/goslogan/grsearch" 
	"github.com/redis/go-redis/v9"
    
)


client := grsearch.NewClient(&redis.Options{})
ctx := context.Background()

if _, err := client.JSONSet(ctx, "helloworld", "$", `{"a": 1, "b": 2, "hello": "world"}`).Result(); err != nil {
    log.Fatalf("Unable to set value: %+v", err)
}

helloVal := client.JSONGet(ctx, "helloworld", "$.hello").Val()


```