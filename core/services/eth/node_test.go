package eth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NodeWrapError(t *testing.T) {
	t.Run("handles nil errors", func(t *testing.T) {
		err := wrap(nil, "foo")
		assert.NoError(t, err)
	})

	t.Run("adds extra info to context deadline exceeded errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		err := ctx.Err()

		err = wrap(err, "foo")

		assert.EqualError(t, err, "foo call failed: remote eth node timed out: context deadline exceeded")
	})
}
