package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestKeypath(t *testing.T) {
	t.Parallel()

	t.Run("can be constructed from a period-delimited string", func(t *testing.T) {
		kp, err := pipeline.NewKeypathFromString("")
		assert.NoError(t, err)
		assert.Equal(t, pipeline.Keypath{}, kp)

		kp, err = pipeline.NewKeypathFromString("foo")
		assert.NoError(t, err)
		assert.Equal(t, pipeline.Keypath{[]string{"foo"}}, kp)

		kp, err = pipeline.NewKeypathFromString("foo.bar")
		assert.NoError(t, err)
		assert.Equal(t, pipeline.Keypath{[]string{"foo", "bar"}}, kp)

		kp, err = pipeline.NewKeypathFromString("a.b.c.d.e")
		assert.NoError(t, err)
		assert.Equal(t, pipeline.Keypath{[]string{"a", "b", "c", "d", "e"}}, kp)
	})

	t.Run("wrong keypath", func(t *testing.T) {
		wrongKeyPath := []string{
			".",
			"..",
			"x.",
			".y",
			"x.y.",
			"x.y..z",
		}

		for _, keypath := range wrongKeyPath {
			t.Run(keypath, func(t *testing.T) {
				_, err := pipeline.NewKeypathFromString(keypath)
				assert.ErrorIs(t, err, pipeline.ErrWrongKeypath)
			})
		}
	})
}
