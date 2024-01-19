package codec_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
)

func TestMultiModifier(t *testing.T) {
	t.Parallel()

	type testStruct struct{ A int }
	testType := reflect.TypeOf(testStruct{})
	mod1 := codec.NewRenamer(map[string]string{"A": "B"})
	mod2 := codec.NewRenamer(map[string]string{"B": "C"})
	chainMod := codec.MultiModifier{mod1, mod2}

	t.Run("RetypeToOffChain chains modifiers", func(t *testing.T) {
		offChain, err := chainMod.RetypeToOffChain(testType, "")
		require.NoError(t, err)
		m1, err := mod1.RetypeToOffChain(testType, "")
		require.NoError(t, err)
		expected, err := mod2.RetypeToOffChain(m1, "")
		require.NoError(t, err)
		assert.Equal(t, expected, offChain)
	})

	t.Run("TransformToOffChain chains modifiers", func(t *testing.T) {
		_, err := chainMod.RetypeToOffChain(testType, "")
		require.NoError(t, err)

		input := testStruct{A: 100}
		actual, err := chainMod.TransformToOffChain(input, "")
		require.NoError(t, err)

		m1, err := mod1.TransformToOffChain(input, "")
		require.NoError(t, err)
		expected, err := mod2.TransformToOffChain(m1, "")
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("TransformToOnChain chains modifiers", func(t *testing.T) {
		offChainType, err := chainMod.RetypeToOffChain(testType, "")
		require.NoError(t, err)

		input := reflect.New(offChainType).Elem()
		input.FieldByName("C").SetInt(100)
		actual, err := chainMod.TransformToOnChain(input.Interface(), "")
		require.NoError(t, err)

		expected := testStruct{A: 100}
		assert.Equal(t, expected, actual)
	})
}
