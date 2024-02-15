package encodings_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings/testutils"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestNotNilPointer(t *testing.T) {
	t.Parallel()
	anyValue := "foo"
	elm := &testutils.TestTypeCodec{
		Value: anyValue,
		Bytes: []byte("foo"),
	}

	nnp := encodings.NotNilPointer{Elm: elm}

	t.Run("Encode delegates", func(t *testing.T) {
		anyPrefix := []byte("bar")
		actual, err := nnp.Encode(&anyValue, anyPrefix)
		require.NoError(t, err)
		require.Equal(t, append(anyPrefix, elm.Bytes...), actual)
	})

	t.Run("Encode returns an error if the pointer is nil", func(t *testing.T) {
		_, err := nnp.Encode((*string)(nil), nil)
		require.Error(t, err)
	})

	t.Run("Encode returns an error element is not a pointer", func(t *testing.T) {
		_, err := nnp.Encode(anyValue, nil)
		require.Error(t, err)
	})

	t.Run("Decode delegates decoding returning a pointer to the value", func(t *testing.T) {
		anySuffix := []byte("fooz")
		actual, remaining, err := nnp.Decode(append(elm.Bytes, anySuffix...))
		require.NoError(t, err)
		require.Equal(t, &anyValue, actual)
		require.Equal(t, anySuffix, remaining)
	})

	t.Run("Decode returns errors from decoder", func(t *testing.T) {
		anyError := fmt.Errorf("%w: nope", types.ErrInvalidEncoding)
		errNnp := &encodings.NotNilPointer{Elm: &testutils.TestTypeCodec{Value: "foo", Err: anyError}}
		_, _, err := errNnp.Decode([]byte("foo"))
		require.Equal(t, anyError, err)
	})

	t.Run("GetType returns pointer to element's GetType", func(t *testing.T) {
		assert.Equal(t, reflect.TypeOf(&anyValue), nnp.GetType())
	})

	t.Run("Size delegates", func(t *testing.T) {
		size, err := nnp.Size(100)
		require.NoError(t, err)
		assert.Equal(t, 100*len(elm.Bytes), size)
	})

	t.Run("FixedSize delegates", func(t *testing.T) {
		size, err := nnp.FixedSize()
		require.NoError(t, err)
		assert.Equal(t, len(elm.Bytes), size)
	})
}
