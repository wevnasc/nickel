package json

import (
	"bytes"
	"io"
	"io/ioutil"
	"nickel/core/errors"
	"nickel/serializer"
)

func DecodeBody(serializer serializer.Serializer, body io.ReadCloser, out interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(errors.Serialization, "not was possible to decode http body data", err)
	}
	return serializer.Decode(data, out)
}

func EncodeBody(serializer serializer.Serializer, payload interface{}) (io.Reader, error) {
	data, err := serializer.Encode(payload)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(data), nil
}
