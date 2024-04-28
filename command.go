package grsearch

import (
	"context"
	"fmt"

	"github.com/goslogan/grsearch/internal"
	"github.com/redis/go-redis/v9"
)

/*******************************************************************************
 ***** QueryCmd 														  ******
 *******************************************************************************/

type QueryCmd struct {
	redis.Cmd
	totalResults int64
	keymap       map[string]int
	respData     *RESPData
	val          []*SearchResult
	options      *QueryOptions
	onHash       bool
	process      cmdable // used to initialise iterator
	count        int64   // contains the total number of results if the query was successful
}

type RESPData struct {
	Errors     []interface{}
	Warnings   []interface{}
	Format     string
	Attributes []interface{}
}

// NewQueryCmd returns an initialised query command.
func NewQueryCmd(ctx context.Context, process cmdable, onHash bool, args ...interface{}) *QueryCmd {
	return &QueryCmd{
		process: process,
		onHash:  onHash,
		Cmd:     *redis.NewCmd(ctx, args...),
	}
}
func (cmd *QueryCmd) SetVal(val []*SearchResult) {
	cmd.val = val

	cmd.keymap = map[string]int{}
	for n, v := range val {
		cmd.keymap[v.Key] = n
	}
}

func (cmd *QueryCmd) Val() []*SearchResult {
	return cmd.val
}

func (cmd *QueryCmd) Result() ([]*SearchResult, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *QueryCmd) Len() int64 {
	if cmd.Err() != nil {
		return 0
	} else {
		return int64(len(cmd.val))
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

// RESPData returns the additional data returned with a RESP3 response if set.
func (cmd *QueryCmd) RESP3Data() *RESPData {
	return cmd.respData
}

// SetRESPData stores the additional data returned with a RESP3 response if set.
func (cmd *QueryCmd) SetRESP3Data(data *RESPData) {
	cmd.respData = data
}

// TotalResults returns the total number of possible results for the query (whilst Count returns
// the number of results from a single call to FTSEearch)
func (cmd *QueryCmd) TotalResults() int64 {
	return cmd.totalResults
}

// SetTotalResults store the total number of possible results for the query.
func (cmd *QueryCmd) SetTotalResults(r int64) {
	cmd.totalResults = r
}

// Key returns the individual result with the
// given key
func (cmd *QueryCmd) Key(key string) *SearchResult {
	return cmd.val[cmd.keymap[key]]
}

// Keys returns the redis keys for all of the results
func (cmd *QueryCmd) Keys() []string {
	results := make([]string, len(cmd.keymap))
	for i, k := range cmd.val {
		results[i] = k.Key
	}

	return results
}

func (cmd *QueryCmd) postProcess() error {
	if cmd.Err() != nil {
		return cmd.Err()
	}

	var parser func(int, interface{}) (*SearchResult, error)
	if cmd.onHash {
		parser = parseHashResult
	} else {
		parser = parseJSONResult
	}

	rawResults := cmd.Cmd.Val()
	var err error

	// RESP2 or RESP3?
	switch rawResults.(type) {
	case map[interface{}]interface{}:
		err = cmd.postprocessRESP3Response(rawResults, parser)
	case []interface{}:
		err = cmd.postprocessRESP2Response(rawResults, parser)
	default:
		return fmt.Errorf("redis: %v is not a valid type search result", rawResults)
	}

	if err != nil {
		return err
	}

	return nil
}

func (cmd *QueryCmd) postprocessRESP3Response(baseResponse interface{}, parser func(int, interface{}) (*SearchResult, error)) error {
	response, ok := baseResponse.(map[interface{}]interface{})

	data := RESPData{}

	if !ok {
		return fmt.Errorf("redis: FT.SEARCH response is not a map")
	}

	data.Attributes = response["attributes"].([]interface{})

	if val, ok := response["error"]; ok {
		data.Errors = val.([]interface{})
	}

	if val, ok := response["warning"]; ok {
		data.Warnings = val.([]interface{})
	}

	data.Format = response["format"].(string)

	cmd.SetRESP3Data(&data)
	cmd.SetTotalResults(response["total_results"].(int64))

	results := []*SearchResult{}
	for _, r := range response["results"].([]interface{}) {
		rawResult := r.(map[interface{}]interface{})
		var current *SearchResult

		if cmd.options.NoContent {
			current = &SearchResult{}
		} else {
			var err error
			if current, err = parser(3, rawResult["extra_attributes"]); err != nil {
				return err
			}
		}
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

		results = append(results, current)
	}

	cmd.SetVal(results)
	return nil
}

func (cmd *QueryCmd) postprocessRESP2Response(baseResponse interface{}, parser func(int, interface{}) (*SearchResult, error)) error {

	if _, ok := baseResponse.([]interface{}); !ok {
		return fmt.Errorf("redis: FT.SEARCH response is not a slice")
	}

	data := RESPData{Format: "STRING"}
	cmd.SetRESP3Data(&data)
	results := []*SearchResult{}

	response := baseResponse.([]interface{})
	cmd.SetTotalResults(response[0].(int64))

	for i := 1; i < len(response); i += cmd.options.resultSize() {

		var current *SearchResult
		var score float64 = 0
		var explanation []interface{}
		j := 0
		key := response[i+j].(string)

		j++

		if cmd.options.WithScores {
			if cmd.options.ExplainScore {
				scoreData := response[i+j].([]interface{})
				score, _ = internal.Float64(scoreData[0])
				explanation = scoreData[1].([]interface{})
			} else {
				score, _ = internal.Float64(response[i+j])
				explanation = nil
			}
			j++
		}

		if cmd.options.NoContent {
			current = &SearchResult{}
		} else {
			if vals, ok := response[i+j].([]interface{}); !ok {
				return fmt.Errorf("redis: response content cannot be parsed")
			} else {
				var err error
				if current, err = parser(2, vals); err != nil {
					return err
				}
			}
		}

		current.Key = key
		current.Explanation = explanation
		current.Score = score

		results = append(results, current)
		j++

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
	respData     *RESPData
	val          []map[string]interface{}
	totalResults int64
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

	respData := &RESPData{}
	rawResults := cmd.Cmd.Val()
	results := make([]map[string]interface{}, 0)

	// RESP2 v RESP3
	switch r := rawResults.(type) {
	case []interface{}:
		cmd.SetTotalResults(int64(len(r) - 1))
		respData.Format = "STRING"
		for n := 1; n < len(r); n++ {
			result := map[string]interface{}{}
			for k, v := range internal.ToMap(r[n]) {
				result[k.(string)] = v
			}
			results = append(results, result)
		}
	case map[interface{}]interface{}:
		n, _ := internal.Int64(r["total_results"])
		cmd.SetTotalResults(n)

		respData.Format = r["format"].(string)
		if w, ok := r["error"]; ok {
			respData.Warnings = w.([]interface{})
		}
		if e, ok := r["error"]; ok {
			respData.Warnings = r["warning"].([]interface{})
			respData.Errors = e.([]interface{})
		}
		for _, data := range r["results"].([]interface{}) {
			result := map[string]interface{}{}
			for k, v := range data.(map[interface{}]interface{}) {
				result[k.(string)] = v
			}
			results = append(results, result)
		}
	}
	cmd.SetRESP3Data(respData)
	cmd.SetVal(results)
	return nil
}

func (cmd *AggregateCmd) SetVal(val []map[string]interface{}) {
	cmd.val = val
}

func (cmd *AggregateCmd) Val() []map[string]interface{} {
	return cmd.val
}

func (cmd *AggregateCmd) Result() ([]map[string]interface{}, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *AggregateCmd) SetTotalResults(n int64) {
	cmd.totalResults = n
}

func (cmd *AggregateCmd) TotalResults() int64 {
	return cmd.totalResults
}

// RESPData returns the additional data returned with a RESP3 response if set.
func (cmd *AggregateCmd) RESP3Data() *RESPData {
	return cmd.respData
}

// SetRESPData stores the additional data returned with a RESP3 response if set.
func (cmd *AggregateCmd) SetRESP3Data(data *RESPData) {
	cmd.respData = data
}

type ExtCmder interface {
	redis.Cmder
	postProcess() error
}
