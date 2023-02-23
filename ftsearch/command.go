package ftsearch

import (
	"context"
	"fmt"
	"reflect"

	"github.com/redis/go-redis/v9"
)

type QueryCmd struct {
	redis.SliceCmd
	val     map[string]*QueryResult
	options *QueryOptions
}

/*******************************************************************************
 ***** QueryCmd 														  ******
 *******************************************************************************/

// NewQueryCmd returns an initialised query command.
func NewQueryCmd(ctx context.Context, args ...interface{}) *QueryCmd {
	return &QueryCmd{
		SliceCmd: *redis.NewSliceCmd(ctx, args...),
	}
}
func (cmd *QueryCmd) SetVal(val map[string]*QueryResult) {
	cmd.val = val
}

func (cmd *QueryCmd) Val() map[string]*QueryResult {
	return cmd.val
}

func (cmd *QueryCmd) Result() (map[string]*QueryResult, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *QueryCmd) Len() int {
	if cmd.Err() != nil {
		return 0
	} else {
		return len(cmd.val)
	}
}

func (cmd *QueryCmd) Scan(dst interface{}) error {
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if cmd.options.NoContent {
		return fmt.Errorf("ftsearch.Scan - NoContent is set")
	}

	if reflect.TypeOf(dst).Kind() != reflect.Slice {
		return fmt.Errorf("ftsearch.Scan - %T is not a slice", dst)
	}

	list := reflect.ValueOf(dst)

	if list.Len() < cmd.Len() {
		return fmt.Errorf("ftsearch.Scan - %T is not large enough: %d < %d", dst, list.Len(), cmd.Len())
	}

	// Because we don't have access to go-redis internals we fake this
	// as MapStringStringCmds and scan with them.
	n := 0
	for _, result := range cmd.val {
		sCmd := redis.NewMapStringStringCmd(context.Background(), "DUMMY")
		sCmd.SetVal(result.Value)
		item := list.Index(n).Interface()
		if err := sCmd.Scan(item); err != nil {
			return err
		}
		n++
	}

	return nil
}

func (cmd *QueryCmd) String() string {
	return cmd.SliceCmd.String()
}

func (cmd *QueryCmd) toMap(input []interface{}) map[string]string {
	results := make(map[string]string, len(input)/2)
	key := ""
	for i := 0; i < len(input); i += 2 {
		key = input[i].(string)
		value := input[i+1].(string)
		results[key] = value
	}
	return results
}

func (cmd *QueryCmd) parseResult() error {
	if cmd.Err() != nil {
		return cmd.Err()
	}
	rawResults := cmd.SliceCmd.Val()
	resultSize := cmd.options.resultSize
	resultCount := rawResults[0].(int64)
	results := make(map[string]*QueryResult, resultCount)

	for i := 1; i < len(rawResults); i += resultSize {
		j := 0
		var score float64 = 0

		key := rawResults[i+j].(string)
		j++

		if cmd.options.Scores {
			score, _ = rawResults[i+j].(float64)
			j++
		}

		result := QueryResult{
			Score: score,
		}

		if !cmd.options.NoContent {
			result.Value = cmd.toMap(rawResults[i+j].([]interface{}))
			j++
		}

		results[key] = &result
	}

	cmd.SetVal(results)
	return nil
}

/*******************************************************************************
 ***** ConfigGetCmd 													  ******
 *******************************************************************************/

type ConfigGetCmd struct {
	redis.Cmd
	val map[string]string
}

func NewConfigGetCmd(ctx context.Context, args ...interface{}) *ConfigGetCmd {
	return &ConfigGetCmd{
		Cmd: *redis.NewCmd(ctx, args...),
	}
}

func (c *ConfigGetCmd) parseResult() {
	if result, err := c.Slice(); err == nil {
		configs := make(map[string]string, len(result))
		for _, cfg := range result {
			key := cfg.([]interface{})[0].(string)
			if key[0] != '_' {
				if cfg.([]interface{})[1] != nil {
					val := cfg.([]interface{})[1].(string)
					configs[key] = val
				} else {
					configs[key] = ""
				}
			}

		}
		c.SetVal(configs)
	}

}

func (cmd *ConfigGetCmd) SetVal(val map[string]string) {
	cmd.val = val
}

func (cmd *ConfigGetCmd) Val() map[string]string {
	return cmd.val
}

func (cmd *ConfigGetCmd) Result() (map[string]string, error) {
	return cmd.Val(), cmd.Err()
}
