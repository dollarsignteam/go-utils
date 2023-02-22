package utils

import "encoding/json"

// JSON utility instance
var JSON JSONUtil

// JSONUtil is a struct with methods for parsing and validating JSON data.
type JSONUtil struct{}

// ParseAndValidate parses and validates JSON data into the given result struct.
// If the result is an array or slice, it is validated as a list of items.
func (JSONUtil) ParseAndValidate(data string, result any) error {
	if err := json.Unmarshal([]byte(data), result); err != nil {
		return err
	}
	if IsArrayOrSlice(result) {
		return Validate.Struct(Result{List: result})
	}
	return Validate.Struct(result)
}

// Parse parses JSON data into the given result struct.
func (JSONUtil) Parse(data string, result any) error {
	return json.Unmarshal([]byte(data), result)
}
