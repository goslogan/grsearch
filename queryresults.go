package grstack

// code to process the decoding/parsing of individual values in the query results.

import (
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type QueryResults struct {
	TotalResults int64
	Results      []*QueryResult
	keymap       map[string]int
	Errors       []interface{}
	Format       string
	Attributes   []interface{}
}

type ResultValue interface {
	parse(int, []interface{}) error
}

type QueryResult struct {
	Key         string
	Score       float64
	Explanation interface{}
	Values      ResultValue
}

type HashQueryValue struct {
	Value map[string]string
}

type JSONQueryValue struct {
	Value map[string]string
}

func (r *HashQueryValue) parse(respVersion int, input []interface{}) error {

	if respVersion == 2 {
		results := make(map[string]string, len(input)/2)
		key := ""
		for i := 0; i < len(input); i += 2 {
			key = input[i].(string)
			value := input[i+1].(string)
			results[key] = value
		}
		r.Value = results
		return nil
	} else if respVersion == 3 {
		return nil
	} else {
		return fmt.Errorf("redis: invalid RESP version: %d", respVersion)
	}
}

func (r *HashQueryValue) Scan(dst interface{}) error {
	sCmd := redis.NewMapStringStringResult(r.Value, nil)
	return sCmd.Scan(dst)
}

func (r *JSONQueryValue) parse(respVersion int, input []interface{}) error {

	r.Value = map[string]string{}
	key := input[0].(string)
	value := input[1].(string)

	r.Value[key] = value
	return nil
}

func (r *JSONQueryValue) Scan(path string, to interface{}) error {
	return json.Unmarshal([]byte(r.Value[path]), to)
}

// SetResults stores search results into the struct and builds the
// map used for fast key lookup
func (q *QueryResults) SetResults(r []*QueryResult) {
	q.keymap = map[string]int{}
	for n, v := range r {
		q.keymap[v.Key] = n
	}
	q.Results = r
}

// Key returns the individual result with the
// given key
func (q QueryResults) Key(key string) *QueryResult {
	return q.Results[q.keymap[key]]
}

// Keys returns the redis keys for all of the results
func (q QueryResults) Keys() []string {
	results := make([]string, len(q.keymap))
	for i, k := range q.Results {
		results[i] = k.Key
	}

	return results
}
