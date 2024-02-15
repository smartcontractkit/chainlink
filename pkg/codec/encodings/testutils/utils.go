package testutils

import (
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type TestTypeCodec struct {
	Value any
	Bytes []byte
	Err   error
}

func (t *TestTypeCodec) Size(size int) (int, error) {
	return len(t.Bytes) * size, t.Err
}

func (t *TestTypeCodec) FixedSize() (int, error) {
	return len(t.Bytes), t.Err
}

func (t *TestTypeCodec) Encode(_ any, into []byte) ([]byte, error) {
	all := append(into, t.Bytes...)
	return all, t.Err
}

func (t *TestTypeCodec) Decode(encoded []byte) (any, []byte, error) {
	if len(encoded) < len(t.Bytes) {
		return nil, nil, types.ErrInvalidEncoding
	}

	return t.Value, encoded[len(t.Bytes):], t.Err
}

func (t *TestTypeCodec) GetType() reflect.Type {
	return reflect.TypeOf(t.Value)
}

var _ encodings.TypeCodec = &TestTypeCodec{}
