package json

import (
	"encoding/json"
)

type JsonSerializer struct {
}

func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

func (j *JsonSerializer) Decode(input []byte, out interface{}) error {
	if err := json.Unmarshal(input, out); err != nil {
		return err
	}
	return nil
}

func (j *JsonSerializer) Encode(input interface{}) ([]byte, error) {
	rawMessage, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return rawMessage, err
}
