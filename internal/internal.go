package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

// serializeCountedArgs is used to serialize a string array to
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

// AppendStringArg appends the name and value if value is not empty
func AppendStringArg(args []interface{}, name, value string) []interface{} {
	if value != "" {
		return append(args, name, value)
	} else {
		return args
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

		return ToMap(data.([]interface{})), nil
	}

}

// ToMap simply converts a slice of key/value pairs returned from Redis
// into a map[string]interface{}
func ToMap(input []interface{}) map[string]interface{} {
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
