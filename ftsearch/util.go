package ftsearch

// serializeCountedArgs is used to serialize a string array to
// NAME <count> values. If incZero is true then NAME 0 will be generated
// otherwise empty results will not be generated.
func serializeCountedArgs(name string, incZero bool, args []string) []interface{} {
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
