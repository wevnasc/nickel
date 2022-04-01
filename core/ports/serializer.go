package ports

type SerializerPort interface {
	Decode(input []byte, out interface{}) error
	Encode(input interface{}) ([]byte, error)
}
