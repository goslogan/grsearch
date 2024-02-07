package grstack

import (
	"context"
	"fmt"

	"github.com/goslogan/grstack/internal"
	"github.com/redis/go-redis/v9"
)

/*******************************************************************************
 ***** QueryCmd 														  ******
 *******************************************************************************/

type QueryCmd struct {
	redis.Cmd
	val     QueryResults
	options *QueryOptions
	process cmdable // used to initialise iterator
	count   int64   // contains the total number of results if the query was successful
}

// NewQueryCmd returns an initialised query command.
func NewQueryCmd(ctx context.Context, process cmdable, args ...interface{}) *QueryCmd {
	return &QueryCmd{
		process: process,
		Cmd:     *redis.NewCmd(ctx, args...),
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

func (cmd *QueryCmd) Len() int64 {
	if cmd.Err() != nil {
		return 0
	} else {
		return int64(len(cmd.val.Results))
	}
}

func (cmd *QueryCmd) String() string {
	return cmd.Cmd.String()
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
	return NewSearchIterator(ctx, cmd, cmd.process)
}

func (cmd *QueryCmd) postProcess() error {
	if cmd.Err() != nil {
		return cmd.Err()
	}

	rawResults := cmd.Cmd.Val()
	var err error
	var values *QueryResults

	// RESP2 or RESP3?
	switch rawResults.(type) {
	case map[interface{}]interface{}:
		values, err = cmd.postprocessRESP3Response(rawResults)
	case []interface{}:
		values, err = cmd.postprocessRESP2Response(rawResults)
	default:
		return fmt.Errorf("redis: %v is not a valid type search result", rawResults)
	}

	if err != nil {
		return err
	}

	cmd.SetVal(*values)
	return nil
}

func (cmd *QueryCmd) postprocessRESP3Response(baseResponse interface{}) (*QueryResults, error) {
	response, ok := baseResponse.(map[interface{}]interface{})

	if !ok {
		return nil, fmt.Errorf("redis: FT.SEARCH response is not a map")
	}

	output := QueryResults{Results: []*QueryResult{}}
	output.Attributes = response["attributes"].([]interface{})

	if val, ok := response["error"]; ok {
		output.Errors = val.([]interface{})
	}

	if val, ok := response["warning"]; ok {
		output.Warnings = val.([]interface{})
	}

	output.Format = response["format"].(string)
	output.TotalResults = response["total_results"].(int64)

	results := []*QueryResult{}
	for _, r := range response["results"].([]interface{}) {
		rawResult := r.(map[interface{}]interface{})
		current := QueryResult{}
		current.Key = rawResult["id"].(string)

		if cmd.options.WithScores {
			if cmd.options.ExplainScore {
				scoreInfo := rawResult["score"].([]interface{})
				current.Score = scoreInfo[0].(float64)
				current.Explanation = scoreInfo[1]
			} else {
				current.Score = rawResult["score"].(float64)
			}

		}

		var result ResultValue

		if cmd.options.json {
			result = &JSONQueryValue{}
		} else {
			result = &HashQueryValue{}
		}

		if !cmd.options.NoContent {
			err := result.parse(3, rawResult["extra_attributes"])
			if err != nil {
				return nil, err
			}
			current.Values = result
		}
		results = append(results, &current)
	}

	output.SetResults(results)
	return &output, nil

}

func (cmd *QueryCmd) postprocessRESP2Response(baseResponse interface{}) (*QueryResults, error) {

	if _, ok := baseResponse.([]interface{}); !ok {
		return nil, fmt.Errorf("redis: FT.SEARCH response is not a slice")
	}

	output := QueryResults{Format: "STRING"}
	results := []*QueryResult{}

	response := baseResponse.([]interface{})
	output.TotalResults = response[0].(int64)

	for i := 1; i < len(response); i += cmd.options.resultSize() {

		current := &QueryResult{}
		j := 0
		var score float64 = 0
		var explanation []interface{}

		current.Key = response[i+j].(string)
		j++

		if cmd.options.WithScores {
			if cmd.options.ExplainScore {
				scoreData := response[i+j].([]interface{})
				score, _ = internal.Float64(scoreData[0])
				explanation = scoreData[1].([]interface{})
				current.Score = score
				current.Explanation = explanation
			} else {
				current.Score, _ = internal.Float64(response[i+j])
			}
			j++
		}

		if !cmd.options.NoContent {
			if vals, ok := response[i+j].([]interface{}); !ok {
				return nil, fmt.Errorf("redis: response content cannot be parsed")
			} else {
				var result ResultValue
				if cmd.options.json {
					result = &JSONQueryValue{}
				} else {
					result = &HashQueryValue{}
				}
				err := result.parse(2, vals)
				if err != nil {
					return nil, err
				}
				current.Values = result

			}

		}

		results = append(results, current)
		j++

	}

	output.SetResults(results)
	return &output, nil
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

func (cmd *SynonymDumpCmd) postProcess() error {
	r := cmd.Cmd.Val()
	synonymMap := make(map[string][]string)

	switch result := r.(type) {
	case []interface{}: // RESP 2

		for n := 0; n < len(result); n += 2 {
			synonym := result[n].(string)
			groups := make([]string, len(result[n+1].([]interface{})))
			for m, group := range result[n+1].([]interface{}) {
				groups[m] = group.(string)
			}
			synonymMap[synonym] = groups
		}
	case map[interface{}]interface{}: // RESP3
		for k, v := range result {
			groups := make([]string, len(v.([]interface{})))
			for m, group := range v.([]interface{}) {
				groups[m] = group.(string)
			}
			synonymMap[k.(string)] = groups
		}
	default:
		return fmt.Errorf("redis: %v is not a valid result for FT.SYNONYMDUMP", r)
	}

	cmd.SetVal(synonymMap)
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
	redis.Cmd
	val *Info
}

func NewInfoCmd(ctx context.Context, args ...interface{}) *InfoCmd {
	return &InfoCmd{
		Cmd: *redis.NewCmd(ctx, args...),
	}
}

func (c *InfoCmd) SetVal(i *Info) {
	c.val = i
}

func (cmd *InfoCmd) Val() *Info {
	return cmd.val
}

func (cmd *InfoCmd) Result() (*Info, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *InfoCmd) postProcess() error {

	rawResult := cmd.Cmd.Val()
	var mapped map[interface{}]interface{}

	switch v := rawResult.(type) {
	case []interface{}:
		mapped = internal.ToMap(v)
	case map[interface{}]interface{}:
		mapped = v
	default:
		return fmt.Errorf("redis: FT.INFO - invalid response type")
	}

	info := Info{}
	err := info.parse(mapped)

	cmd.SetVal(&info)
	return err
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
	redis.Cmd
	val AggregateResults
}

func NewAggregateCmd(ctx context.Context, args ...interface{}) *AggregateCmd {
	return &AggregateCmd{
		Cmd: *redis.NewCmd(ctx, args...),
	}
}

func (cmd *AggregateCmd) postProcess() error {

	if cmd.Err() != nil {
		return cmd.Err()
	}

	rawResults := cmd.Cmd.Val()
	results := AggregateResults{Results: make([]map[string]interface{}, 0)}

	// RESP2 v RESP3
	switch r := rawResults.(type) {
	case []interface{}:
		results.TotalResults = int64(len(r))
		results.Format = "STRING"
		for _, data := range r {
			result := map[string]interface{}{}
			for k, v := range internal.ToMap(data) {
				result[k.(string)] = v
			}
			results.Results = append(results.Results, result)
		}
	case map[interface{}]interface{}:
		results.TotalResults, _ = internal.Int64(r["total_results"])
		results.Format = r["format"].(string)
		if w, ok := r["error"]; ok {
			results.Warnings = w.([]interface{})
		}
		results.Warnings = r["warning"].([]interface{})
		if e, ok := r["error"]; ok {
			results.Errors = e.([]interface{})
		}
		for _, data := range r["results"].([]interface{}) {
			result := map[string]interface{}{}
			for k, v := range data.(map[interface{}]interface{}) {
				result[k.(string)] = v
			}
			results.Results = append(results.Results, result)
		}
	}

	cmd.SetVal(results)
	return nil
}

func (cmd *AggregateCmd) SetVal(val AggregateResults) {
	cmd.val = val
}

func (cmd *AggregateCmd) Val() AggregateResults {
	return cmd.val
}

func (cmd *AggregateCmd) Result() (AggregateResults, error) {
	return cmd.Val(), cmd.Err()
}

type ExtCmder interface {
	redis.Cmder
	postProcess() error
}
