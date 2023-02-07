package utils

import "encoding/json"

var JSON JSONUtil

type JSONUtil struct{}

func (JSONUtil) ParseAndValidate(data string, result any) error {
	if err := json.Unmarshal([]byte(data), result); err != nil {
		return err
	}
	return Validate.Struct(result)
}
