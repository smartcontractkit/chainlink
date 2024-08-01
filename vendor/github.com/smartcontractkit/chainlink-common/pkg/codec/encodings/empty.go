package encodings

import (
	"reflect"
)

type Empty struct{}

func (Empty) Encode(_ any, into []byte) ([]byte, error) {
	return into, nil
}

func (Empty) Decode(encoded []byte) (any, []byte, error) {
	return struct{}{}, encoded, nil
}

func (Empty) GetType() reflect.Type {
	return reflect.TypeOf(struct{}{})
}

func (Empty) Size(numItems int) (int, error) {
	return 0, nil
}

func (Empty) FixedSize() (int, error) {
	return 0, nil
}

var _ TypeCodec = Empty{}
