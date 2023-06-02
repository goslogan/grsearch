package grstack

// code to process the decoding/parsing of individual values in the query results.

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type QueryResults []*QueryResult

type ResultValue interface {
	parse([]interface{}) error
}

type QueryResult struct {
	Key         string
	Score       float64
	Explanation []interface{}
	Values      ResultValue
}

type HashQueryValue struct {
	Value map[string]string
}

type JSONQueryValue struct {
	Value    map[string]interface{}
	rawValue map[string]string
}

func (r *HashQueryValue) parse(input []interface{}) error {
	results := make(map[string]string, len(input)/2)
	key := ""
	for i := 0; i < len(input); i += 2 {
		key = input[i].(string)
		value := input[i+1].(string)
		results[key] = value
	}
	r.Value = results
	return nil
}

func (r *HashQueryValue) Scan(dst interface{}) error {
	sCmd := redis.NewMapStringStringResult(r.Value, nil)
	return sCmd.Scan(dst)
}

func (r *JSONQueryValue) parse(input []interface{}) error {

	key := input[0].(string)
	rawValue := input[1].(string)
	var result interface{}
	err := json.Unmarshal([]byte(rawValue), &result)

	if r.Value == nil {
		r.rawValue = make(map[string]string)
		r.Value = make(map[string]interface{})
	}

	r.rawValue[key] = rawValue
	r.Value[key] = result
	return err
}

func (r *JSONQueryValue) Scan(path string, to interface{}) error {
	return json.Unmarshal([]byte(r.rawValue[path]), to)
}

// Key returns the individual result with the
// given key
func (q QueryResults) Key(key string) *QueryResult {
	for _, r := range q {
		if r.Key == key {
			return r
		}
	}
	return nil
}

// Keys returns the redis keys for all of the results
func (q QueryResults) Keys() []string {
	results := make([]string, len(q))
	for i, k := range q {
		results[i] = k.Key
	}

	return results
}
