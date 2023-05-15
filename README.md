# grstack

go-redis based interface to RediSearch & RedisJSON designed to extend the go-redis
API and follow its style as closely as possible.


## Example - JSON


```

import (
    "log"

    "github.com/goslogan/grstack" 
	"github.com/redis/go-redis/v9"
    
)

client := grstack.NewClient(&redis.Options{})

if _, err := client.JSONSet(ctx, "helloworld", "$", `{"a": 1, "b": 2, "hello": "world"}`).Result(); err != nil {
    log.Fatalf("Unable to set value: %+v", err)
}

helloVal := client.JSONGet(ctx, "helloworld", "$.hello").Val()

...
```

## Example - search

```
if _, err := client.FTCreateIndex(ctx, "customers", grstack.NewIndexBuilder().
        On("hash").
		Prefix("account:").
		Schema(grstack.TagAttribute{
			Name    : "account_id",
			Alias   : "id",
			Sortable: true}).
        Schema(grstack.TextAttribute{Name: "customer",
		    Sortable: true}).Schema(grstack.TextAttribute{
		    Name    : "email",
		    Sortable: true}).
        Schema(grstack.TagAttribute{
		    Name    : "account_owner",
		    Alias   : "owner",
		    Sortable: true}).
        Schema(grstack.NumericAttribute{
		    Name    : "balance",
		    Sortable: true,
	}).Options()); err != nil {
        log.Fatalf("Unable to create index: %+v", err)
    }

searchResult, err := client.FTSearch(ctx, "customers", "@id:{1128564}").Results()

...

```



