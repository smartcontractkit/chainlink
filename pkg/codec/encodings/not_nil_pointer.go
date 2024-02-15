package encodings

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type NotNilPointer struct {
	Elm TypeCodec
}

func (n *NotNilPointer) Encode(value any, into []byte) ([]byte, error) {
	rValue := reflect.ValueOf(value)
	if rValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("%w: expected pointer, got %T", types.ErrInvalidType, value)
	}

	if rValue.IsNil() {
		return nil, fmt.Errorf("%w: pointer is nil", types.ErrInvalidType)
	}

	return n.Elm.Encode(rValue.Elem().Interface(), into)
}

func (n *NotNilPointer) Decode(encoded []byte) (any, []byte, error) {
	val, remaining, err := n.Elm.Decode(encoded)
	if err != nil {
		return nil, nil, err
	}
	ret := reflect.New(n.Elm.GetType())
	reflect.Indirect(ret).Set(reflect.ValueOf(val))
	return ret.Interface(), remaining, nil
}

func (n *NotNilPointer) GetType() reflect.Type {
	return reflect.PointerTo(n.Elm.GetType())
}

func (n *NotNilPointer) Size(numItems int) (int, error) {
	return n.Elm.Size(numItems)
}

func (n *NotNilPointer) FixedSize() (int, error) {
	return n.Elm.FixedSize()
}

var _ TypeCodec = (*NotNilPointer)(nil)
