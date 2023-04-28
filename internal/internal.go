package internal

import "encoding/json"

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