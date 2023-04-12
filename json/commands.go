package json

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Set sets the JSON value at the given path in the given key
func Set(ctx context.Context, key, path string, value interface{}) *redis.StatusCmd {

}
