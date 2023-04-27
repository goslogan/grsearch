package grstack

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// JSONArrAppend adds the provided json values to the end of the at path
func (c cmdable) JSONArrAppend(ctx context.Context, key, path string, values ...interface{}) *redis.IntSliceCmd {
	args := []interface{}{"json.arrappend", key, path}
	args = append(args, values...)
	cmd := redis.NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrIndex searches for the first occurrence of a JSON value in an array
func (c cmdable) JSONArrIndex(ctx context.Context, key, path string, value interface{}) *redis.IntSliceCmd {
	args := []interface{}{"json.arrindex", key, path, value}
	cmd := redis.NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrIndexFromTo searches for the first occurrence of a JSON value in an array whilst allowing the start and
// stop options to be provided.
func (c cmdable) JSONArrIndexStartStop(ctx context.Context, key, path string, value interface{}, start, stop int64) *redis.IntSliceCmd {
	args := []interface{}{"json.arrindex", key, path, value, start, stop}
	cmd := redis.NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrInsert inserts the json values into the array at path before the index (shifts to the right)
func (c cmdable) JSONArrInsert(ctx context.Context, key, path string, index int64, values ...interface{}) *redis.IntSliceCmd {
	args := []interface{}{"json.arrinsert", key, path, index}
	args = append(args, values...)
	cmd := redis.NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrLen reports the length of the JSON array at path in key
func (c cmdable) JSONArrLen(ctx context.Context, key, path string) *redis.IntSliceCmd {
	args := []interface{}{"json.arrlen", key, path}
	cmd := redis.NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrPop removes and returns an element from the index in the array
func (c cmdable) JSONArrPop(ctx context.Context, key, path string, index int) *redis.StringSliceCmd {
	args := []interface{}{"json.arrpop", key, path, index}
	cmd := redis.NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONArrTrim trims an array so that it contains only the specified inclusive range of elements
func (c cmdable) JSONArrTrim(ctx context.Context, key, path string, start, stop int) *redis.IntSliceCmd {
	args := []interface{}{"json.arrtrim", key, path, start, stop}
	cmd := redis.NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONClear clears container values (arrays/objects) and set numeric values to 0
func (c cmdable) JSONClear(ctx context.Context, key, path string) *redis.IntCmd {
	args := []interface{}{"json.clear", key, path}
	cmd := redis.NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONDel deletes a value
func (c cmdable) JSONDel(ctx context.Context, key, path string) *redis.IntCmd {
	args := []interface{}{"json.del", key, path}
	cmd := redis.NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONForget deletes a value
func (c cmdable) JSONForget(ctx context.Context, key, path string) *redis.IntCmd {
	args := []interface{}{"json.forget", key, path}
	cmd := redis.NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONGet returns the value at path in JSON serialized form
func (c cmdable) JSONGet(ctx context.Context, key string, paths ...string) *JSONStringCmd {

	args := make([]interface{}, len(paths)+2)
	args[0] = "json.get"
	args[1] = key
	for n, path := range paths {
		args[n+2] = path
	}

	cmd := NewJSONStringCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONMGet returns the values at path from multiple key arguments
// Note - the arguments are reversed when compared with `JSON.MGET` as we want
// to follow the pattern of having the last argument be variable.
func (c cmdable) JSONMGet(ctx context.Context, path string, keys ...string) *JSONStringSliceCmd {

	args := make([]interface{}, len(keys)+1)
	args[0] = "json.mget"
	for n, keys := range keys {
		args[n+1] = keys
	}
	args = append(args, path)

	cmd := NewJSONStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// JSONSet sets the JSON value at the given path in the given key. The value must be something that
// can be marshalled to JSON (using encoding/JSON) unless the argument is a string when we assume that
// it can be passed directly as JSON.
func (c cmdable) JSONSet(ctx context.Context, key, path string, value interface{}) *redis.StatusCmd {
	return c.JSONSetMode(ctx, key, path, value, "")
}

// JSONSetMOde sets the JSON value at the given path in the given key allows the mode to be set
// as well (the mode value must be "XX" or "NX").  The value must be something that can be marshalled to JSON (using encoding/JSON) unless
// the argument is a string when we assume that  it can be passed directly as JSON.
func (c cmdable) JSONSetMode(ctx context.Context, key, path string, value interface{}, mode string) *redis.StatusCmd {

	var bytes []byte
	var err error

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	default:
		bytes, err = json.Marshal(v)
	}

	args := []interface{}{"json.set", key, path, bytes}

	if mode != "" {
		args = append(args, mode)
	}

	cmd := redis.NewStatusCmd(ctx, args...)

	if err != nil {
		cmd.SetErr(err)
	} else {
		_ = c(ctx, cmd)
	}

	return cmd
}
