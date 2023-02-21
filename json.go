package utils

import "encoding/json"

var JSON JSONUtil

type JSONUtil struct{}

func (JSONUtil) ParseAndValidate(data string, result any) error {
	if err := json.Unmarshal([]byte(data), result); err != nil {
		return err
	}
	if IsArrayOrSlice(result) {
		return Validate.Struct(Result{List: result})
	}
	return Validate.Struct(result)
}

func (JSONUtil) Parse(data string, result any) error {
	return json.Unmarshal([]byte(data), result)
}
