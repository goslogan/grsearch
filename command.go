package grstack

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/goslogan/grstack/internal"
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

func (cmd *QueryCmd) postProcess() error {
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
		var explanation []interface{}

		key := rawResults[i+j].(string)
		j++

		if cmd.options.WithScores {
			if cmd.options.ExplainScore {
				scoreData := rawResults[i+j].([]interface{})
				score = scoreData[0].(float64)
				explanation = scoreData[1].([]interface{})

			} else {
				score, _ = rawResults[i+j].(float64)
			}
			j++
		}

		result := QueryResult{
			Score:       score,
			Explanation: explanation,
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

func (c *ConfigGetCmd) postProcess() error {
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
	return nil
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

/*******************************************************************************
*
* SynDumpCmd
*
*******************************************************************************/

type SynonymDumpCmd struct {
	redis.Cmd
	val map[string][]string
}

func NewSynonymDumpCmd(ctx context.Context, args ...interface{}) *SynonymDumpCmd {
	return &SynonymDumpCmd{
		Cmd: *redis.NewCmd(ctx, args...),
	}
}

func (c *SynonymDumpCmd) postProcess() error {
	if result, err := c.Slice(); err == nil {
		synonymMap := make(map[string][]string)
		for n := 0; n < len(result); n += 2 {
			synonym := result[n].(string)
			groups := make([]string, len(result[n+1].([]interface{})))
			for m, group := range result[n+1].([]interface{}) {
				groups[m] = group.(string)
			}
			synonymMap[synonym] = groups

		}
		c.SetVal(synonymMap)
	}
	return nil
}

func (cmd *SynonymDumpCmd) SetVal(val map[string][]string) {
	cmd.val = val
}

func (cmd *SynonymDumpCmd) Val() map[string][]string {
	return cmd.val
}

func (cmd *SynonymDumpCmd) Result() (map[string][]string, error) {
	return cmd.Val(), cmd.Err()
}

/*******************************************************************************
*
* InfoCmd
*
*******************************************************************************/

type InfoCmd struct {
	redis.SliceCmd
}

func NewInfoCmd(ctx context.Context, args ...interface{}) *InfoCmd {
	return &InfoCmd{
		SliceCmd: *redis.NewSliceCmd(ctx, args...),
	}
}

func (c *InfoCmd) postProcess() error {
	return nil
}

/*******************************************************************************
*
* JSONStringCmd
*
*******************************************************************************/

type JSONStringCmd struct {
	redis.StringCmd
	val []interface{}
}

func NewJSONStringCmd(ctx context.Context, args ...interface{}) *JSONStringCmd {
	return &JSONStringCmd{
		StringCmd: *redis.NewStringCmd(ctx, args...),
	}
}

// given a string containing a JSON array, turn it into
// an array of json.RawMessage objects
func (c *JSONStringCmd) postProcess() error {

	// nil response from JSON.(M)GET (c.StringCmd.err will be "redis: nil")
	if c.StringCmd.Val() == "" && c.StringCmd.Err().Error() == redis.Nil.Error() {
		c.val = nil
		c.SetErr(nil)
		return nil
	}

	if objects, err := internal.ExtractJSONValue(c.StringCmd.Val()); err != nil {
		c.SetErr(err)
		return err
	} else {
		c.SetVal(objects)
		return nil
	}
}

func (cmd *JSONStringCmd) SetVal(val []interface{}) {
	cmd.val = val
}

func (cmd *JSONStringCmd) Val() []interface{} {
	return cmd.val
}

func (cmd *JSONStringCmd) Result() ([]interface{}, error) {
	return cmd.Val(), cmd.Err()
}

// Scan scans the result at position index in the results into the
// destination.
func (cmd *JSONStringCmd) Scan(index int, dst interface{}) error {
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if index < 0 || index >= len(cmd.val) {
		return fmt.Errorf("JSONCmd.Scan - %d is out of range (0..%d)", index, len(cmd.val))
	}

	results := []json.RawMessage{}
	if err := json.Unmarshal([]byte(cmd.StringCmd.Val()), &results); err != nil {
		return err
	} else {
		return json.Unmarshal(results[index], dst)
	}
}

/*******************************************************************************
*
* JSONStringSliceCmd
*
*******************************************************************************/

// TODO: think of a way to implement Scan for this.
type JSONStringSliceCmd struct {
	redis.StringSliceCmd
	val [][]interface{}
}

func NewJSONStringSliceCmd(ctx context.Context, args ...interface{}) *JSONStringSliceCmd {
	return &JSONStringSliceCmd{
		StringSliceCmd: *redis.NewStringSliceCmd(ctx, args...),
	}
}

// postProcess converts an array of bulk string responses into
// an array of arrays of interfaces.
// an array of json.RawMessage objects
func (c *JSONStringSliceCmd) postProcess() error {

	// nil response from JSON.(M)GET (c.StringCmd.err will be "redis: nil")
	if len(c.StringSliceCmd.Val()) == 0 && c.StringSliceCmd.Err().Error() == redis.Nil.Error() {
		c.val = nil
		c.SetErr(nil)
		return nil
	}

	results := [][]interface{}{}

	for _, val := range c.StringSliceCmd.Val() {
		if objects, err := internal.ExtractJSONValue(val); err != nil {
			c.SetErr(err)
			return err
		} else {
			results = append(results, objects)
		}
	}

	c.SetVal(results)
	return nil
}

func (cmd *JSONStringSliceCmd) SetVal(val [][]interface{}) {
	cmd.val = val
}

func (cmd *JSONStringSliceCmd) Val() [][]interface{} {
	return cmd.val
}

func (cmd *JSONStringSliceCmd) Result() ([][]interface{}, error) {
	return cmd.Val(), cmd.Err()
}

/*******************************************************************************
*
* IntSlicePointerCmd
* used to represent a RedisJSON response where the result is either an integer or nil
*
*******************************************************************************/

type IntSlicePointerCmd struct {
	redis.SliceCmd
	val []*int64
}

// NewIntSlicePointerCmd initialises an IntSlicePointerCmd
func NewIntSlicePointerCmd(ctx context.Context, args ...interface{}) *IntSlicePointerCmd {
	return &IntSlicePointerCmd{
		SliceCmd: *redis.NewSliceCmd(ctx, args...),
	}
}

// postProcess converts an array of bulk string responses into
// an array of arrays of interfaces.
// an array of json.RawMessage objects
func (c *IntSlicePointerCmd) postProcess() error {

	if len(c.SliceCmd.Val()) == 0 {
		c.val = nil
		c.SetErr(nil)
		return nil
	}

	results := []*int64{}

	for _, val := range c.SliceCmd.Val() {
		var result int64
		if val == nil {
			results = append(results, nil)
		} else {
			result = val.(int64)
			results = append(results, &result)
		}
	}

	c.SetVal(results)
	return nil
}

func (cmd *IntSlicePointerCmd) SetVal(val []*int64) {
	cmd.val = val
}

func (cmd *IntSlicePointerCmd) Val() []*int64 {
	return cmd.val
}

func (cmd *IntSlicePointerCmd) Result() ([]*int64, error) {
	return cmd.Val(), cmd.Err()
}

type ExtCmder interface {
	redis.Cmder
	postProcess() error
}
