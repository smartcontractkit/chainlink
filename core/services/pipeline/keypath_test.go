package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestKeypath(t *testing.T) {
	t.Parallel()

	t.Run("can be constructed from a period-delimited string with 2 or fewer parts", func(t *testing.T) {
		t.Parallel()

		kp, err := pipeline.NewKeypathFromString("")
		assert.NoError(t, err)
		assert.Equal(t, pipeline.Keypath{}, kp)

		kp, err = pipeline.NewKeypathFromString("foo")
		assert.NoError(t, err)
		assert.Equal(t, pipeline.Keypath{NumParts: 1, Part0: "foo"}, kp)

		kp, err = pipeline.NewKeypathFromString("foo.bar")
		assert.NoError(t, err)
		assert.Equal(t, pipeline.Keypath{NumParts: 2, Part0: "foo", Part1: "bar"}, kp)
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
			t.Run(keypath, func(t *testing.T) {
				t.Parallel()

				_, err := pipeline.NewKeypathFromString(keypath)
				assert.ErrorIs(t, err, pipeline.ErrWrongKeypath)
			})
		}
	})
}
