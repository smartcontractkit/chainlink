package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestKeypath(t *testing.T) {
	t.Parallel()

	t.Run("can be constructed from a period-delimited string with 2 or fewer parts", func(t *testing.T) {
		t.Parallel()

		kp, err := pipeline.NewKeypathFromString("")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{}, kp)

		kp, err = pipeline.NewKeypathFromString("foo")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{NumParts: 1, Part0: "foo"}, kp)

		kp, err = pipeline.NewKeypathFromString("foo.bar")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{NumParts: 2, Part0: "foo", Part1: "bar"}, kp)
	})

	t.Run("wrong keypath", func(t *testing.T) {
		t.Parallel()

		wrongKeyPath := []string{
			".",
			"..",
			"x.",
			".y",
			"x.y.",
			"x.y.z",
		}

		for _, keypath := range wrongKeyPath {
			_, err := pipeline.NewKeypathFromString(keypath)
			require.ErrorIs(t, err, pipeline.ErrWrongKeypath)
		}
	})
}
