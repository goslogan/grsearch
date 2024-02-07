package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

// SerializeCountedArgs is used to serialize a string array to
// NAME <count> values. If incZero is true then NAME 0 will be generated
// otherwise empty results will not be generated.
func SerializeCountedArgs(name string, incZero bool, args []string) []interface{} {
	if len(args) > 0 || incZero {
		result := make([]interface{}, 2+len(args))

		result[0] = name
		result[1] = len(args)
		for pos, val := range args {
			result[pos+2] = val
		}

		return result
	} else {
		return nil
	}
}

func ExtractJSONValue(val string) ([]interface{}, error) {

	if val == "" {
		return nil, nil
	} else {
		objects := []interface{}{}
		if err := json.Unmarshal([]byte(val), &objects); err != nil {
			return nil, err
		}
		return objects, nil
	}

}

// StringToDurationHookFunc returns a function that decodes strings to
// time.Duration (given that the input is in milliseconds)
func StringToDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {

		if t != reflect.TypeOf(time.Duration(5)) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			return time.ParseDuration(fmt.Sprintf("%dms", data))
		case reflect.Float32, reflect.Float64:
			return time.ParseDuration(fmt.Sprintf("%fms", data))
		case reflect.String:
			return time.ParseDuration(fmt.Sprintf("%sms", data.(string)))
		default:
			return data, nil
		}

	}
}

// SliceToMapHookFunc returns a function that converts a slice to a map[string]interface{}
// if and only if the output is struct
func StringToMapHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {

		if f.Kind() != reflect.Slice {
			return data, nil
		}

		return ToStringMap(data.([]interface{})), nil
	}

}

// ToStringMap simply converts a slice of key/value pairs returned from Redis
// into a map[string]interface{}
func ToStringMap(input []interface{}) map[string]interface{} {
	results := map[string]interface{}{}
	key := ""

	for n, v := range input {
		if n%2 == 0 {
			key = v.(string)
		} else {
			results[key] = v
		}
	}

	return results
}

// ToMap converts an interface containing a slice of key/value pairs returned from redis into a
// map[interface{}]interface{}. If the interface passed in is already a map, return it. Otherwise
// returns an empty interface (buyer beware)
func ToMap(i interface{}) map[interface{}]interface{} {
	results := map[interface{}]interface{}{}
	var key interface{}

	switch input := i.(type) {
	case []interface{}:
		for n, v := range input {
			if n%2 == 0 {
				key = v
			} else {
				results[key] = v
			}
		}
	case map[interface{}]interface{}:
		results = input
	}

	return results
}

// AppendStringArg appends the name and value if value is not empty
func AppendStringArg(args []interface{}, name, value string) []interface{} {
	if value != "" {
		return append(args, name, value)
	} else {
		return args
	}
}

// Float64 takes an interface value and tries to convert it to a float64.
// We only look at types Redis might return
func Float64(arg interface{}) (float64, error) {

	switch v := arg.(type) {
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("redis: unable to convert %v to a float", arg)
	}
}

// Int64 takes an interface value and tries to convert it to an int64. We only
// look at types Redis might return.
func Int64(arg interface{}) (int64, error) {
	switch v := arg.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("redis: unable to convert %v to an integer", arg)
	}
}
