package serializer

type Serializer interface {
	Decode(input []byte, out interface{}) error
	Encode(input interface{}) ([]byte, error)
}
