package grsearch

// code to process the decoding/parsing of individual values in the query results.

import (
	"encoding/json"
	"fmt"
)

type QueryResults struct {
	TotalResults int64
	keymap       map[string]int
	Errors       []interface{}
	Warnings     []interface{}
	Format       string
	Attributes   []interface{}
	Results      []*Result
}

type Result struct {
	Key         string
	Score       float64
	Explanation interface{}
	Values      map[string]string
	respVersion int
}

func parseHashResult(respVersion int, source interface{}) (*Result, error) {

	r := Result{}

	if respVersion == 2 {
		input := source.([]interface{})
		results := make(map[string]string, len(input)/2)
		key := ""
		for i := 0; i < len(input); i += 2 {
			key = input[i].(string)
			value := input[i+1].(string)
			results[key] = value
		}
		r.Values = results
		return &r, nil
	} else if respVersion == 3 {
		input := source.(map[interface{}]interface{})
		results := make(map[string]string)
		for k, v := range input {
			results[k.(string)] = v.(string)
		}
		r.Values = results
		return &r, nil
	} else {
		return nil, fmt.Errorf("redis: invalid RESP version: %d", respVersion)
	}
}

func parseJSONResult(respVersion int, source interface{}) (*Result, error) {

	r := Result{}

	r.Values = map[string]string{}
	r.respVersion = respVersion

	if respVersion == 2 {
		input := source.([]interface{})
		key := input[0].(string)
		value := input[1].(string)
		r.Values[key] = value
	} else if respVersion == 3 {
		input := source.(map[interface{}]interface{})
		for k, v := range input {
			r.Values[k.(string)] = v.(string)
		}
	} else {
		return nil, fmt.Errorf("redis: %d is not a valid RESP version", respVersion)
	}
	return &r, nil
}

// SetResults stores search results into the struct and builds the
// map used for fast key lookup
func (q *QueryResults) SetResults(r []*Result) {
	q.keymap = map[string]int{}
	for n, v := range r {
		q.keymap[v.Key] = n
	}
	q.Results = r
}

// Key returns the individual result with the
// given key
func (q QueryResults) Key(key string) *Result {
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

// UnMarshal is a simple wrapper around encoding/JSON to simplify the retrieval of JSON results.
func (q *Result) UnMarshal(key string, target interface{}) error {
	value, ok := q.Values[key]

	if ok {
		return json.Unmarshal([]byte(value), target)
	} else {
		return fmt.Errorf("redis search result value does  not exist: %s", value)
	}

}
