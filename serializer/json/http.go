package json

import (
	"io"
	"io/ioutil"
	"nickel/core/ports"
)

func DecodeBody(serializer ports.SerializerPort, body io.ReadCloser, out interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return serializer.Decode(data, out)
}
