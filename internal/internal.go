package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

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
