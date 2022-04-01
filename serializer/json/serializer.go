package json

import (
	"encoding/json"
	"nickel/core/errors"
)

type JsonSerializer struct {
}

func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

func (j *JsonSerializer) Decode(input []byte, out interface{}) error {
	if err := json.Unmarshal(input, out); err != nil {
		return errors.Wrap(errors.Serialization, "not was possible to decode the data", err)
	}
	return nil
}

func (j *JsonSerializer) Encode(input interface{}) ([]byte, error) {
	rawMessage, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(errors.Serialization, "not was possible to encode the data", err)
	}
	return rawMessage, err
}
