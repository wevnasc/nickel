package json

import (
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
