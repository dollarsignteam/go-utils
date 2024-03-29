package utils

import "encoding/json"

// JSON utility instance
var JSON JsonUtil

// JsonUtil is a struct with methods for parsing and validating JSON data.
type JsonUtil struct{}

// ParseAndValidate parses and validates JSON data into the given result struct.
// If the result is an array or slice, it is validated as a list of items.
func (JsonUtil) ParseAndValidate(data string, result any) error {
	if err := json.Unmarshal([]byte(data), result); err != nil {
		return err
	}
	if IsArrayOrSlice(result) {
		return ValidateStruct(Result{List: result})
	}
	return ValidateStruct(result)
}

// Parse parses JSON data into the given result struct.
func (JsonUtil) Parse(data string, result any) error {
	return json.Unmarshal([]byte(data), result)
}
