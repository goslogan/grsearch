package grstack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goslogan/grstack/internal"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

/*******************************************************************************
 ***** QueryCmd 														  ******
 *******************************************************************************/

type QueryCmd struct {
	redis.SliceCmd
	val     QueryResults
	options *QueryOptions
	process cmdable // used to initialise iterator
	count   int64   // Contains the total number of results if the query was successful
}

// NewQueryCmd returns an initialised query command.
func NewQueryCmd(ctx context.Context, process cmdable, args ...interface{}) *QueryCmd {
	return &QueryCmd{
		process:  process,
		SliceCmd: *redis.NewSliceCmd(ctx, args...),
	}
}
func (cmd *QueryCmd) SetVal(val QueryResults) {
	cmd.val = val
}

func (cmd *QueryCmd) Val() QueryResults {
	return cmd.val
}

func (cmd *QueryCmd) Result() (QueryResults, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *QueryCmd) Len() int {
	if cmd.Err() != nil {
		return 0
	} else {
		return len(cmd.val)
	}
}

func (cmd *QueryCmd) String() string {
	return cmd.SliceCmd.String()
}

func (cmd *QueryCmd) SetCount(count int64) {
	cmd.count = count
}

// Count returns the total number of results from a successful query.
func (cmd *QueryCmd) Count() int64 {
	return cmd.count
}

// Iterator returns an iterator for the search.
func (cmd *QueryCmd) Iterator(ctx context.Context) *SearchIterator {
	return NewSearchIterator(ctx, cmd)
}

func (cmd *QueryCmd) postProcess() error {
	if cmd.Err() != nil {
		return cmd.Err()
	}
	rawResults := cmd.SliceCmd.Val()
	resultSize := cmd.options.resultSize()
	resultCount := rawResults[0].(int64)
	results := make([]*QueryResult, 0)

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
			Key:         key,
			Score:       score,
			Explanation: explanation,
			Values:      nil,
		}

		if !cmd.options.NoContent {

			if cmd.options.json {
				result.Values = &JSONQueryValue{}
			} else {
				result.Values = &HashQueryValue{}
			}

			if err := result.Values.parse(rawResults[i+j].([]interface{})); err != nil {
				return err
			}
		}

		results = append(results, &result)
		j++

	}

	cmd.SetCount(resultCount)
	cmd.SetVal(QueryResults(results))
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

var _ redis.Cmder = (*SynonymDumpCmd)(nil)

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
	redis.MapStringInterfaceCmd
	val *Info
}

func NewInfoCmd(ctx context.Context, args ...interface{}) *InfoCmd {
	return &InfoCmd{
		MapStringInterfaceCmd: *redis.NewMapStringInterfaceCmd(ctx, args...),
	}
}

func (c *InfoCmd) SetVal(i *Info) {
	c.val = i
}

func (c *InfoCmd) Val() *Info {
	return c.val
}

func (c *InfoCmd) Result() (*Info, error) {
	return c.Val(), c.Err()
}

func (c *InfoCmd) postProcess() error {
	info := Info{}
	config := mapstructure.DecoderConfig{
		DecodeHook:           mapstructure.ComposeDecodeHookFunc(internal.StringToDurationHookFunc(), internal.StringToMapHookFunc()),
		WeaklyTypedInput:     true,
		Result:               &info,
		IgnoreUntaggedFields: true,
	}
	if decoder, err := mapstructure.NewDecoder(&config); err != nil {
		return err
	} else if err := decoder.Decode(c.MapStringInterfaceCmd.Val()); err != nil {
		return err
	}

	info.Index = *NewIndexOptions()
	info.Index.parseInfo(c.MapStringInterfaceCmd.Val())

	c.SetVal(&info)
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

/*******************************************************************************
*
* AggregateCmd
* used to manage the results from FT.AGGREGATE calls
*
*******************************************************************************/

type AggregateCmd struct {
	redis.SliceCmd
	val []map[string]string
}

func NewAggregateCmd(ctx context.Context, args ...interface{}) *AggregateCmd {
	return &AggregateCmd{
		SliceCmd: *redis.NewSliceCmd(ctx, args...),
	}
}

func (c *AggregateCmd) postProcess() error {
	if len(c.SliceCmd.Val()) == 0 {
		c.val = nil
		c.SetErr(nil)
		return nil
	}

	results := make([]map[string]string, len(c.SliceCmd.Val())-1)

	for n, entry := range c.SliceCmd.Val() {

		if n > 0 {
			row := entry.([]interface{})
			asStrings := map[string]string{}
			for m := 0; m < len(row); m += 2 {
				asStrings[row[m].(string)] = row[m+1].(string)
			}
			results[n-1] = asStrings
		}
	}

	c.SetVal(results)
	return nil
}

func (cmd *AggregateCmd) SetVal(val []map[string]string) {
	cmd.val = val
}

func (cmd *AggregateCmd) Val() []map[string]string {
	return cmd.val
}

func (cmd *AggregateCmd) Result() ([]map[string]string, error) {
	return cmd.Val(), cmd.Err()
}

type ExtCmder interface {
	redis.Cmder
	postProcess() error
}
