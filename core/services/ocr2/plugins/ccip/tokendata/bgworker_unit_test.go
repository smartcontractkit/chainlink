package tokendata

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newResultsCache(t *testing.T) {
	ctx := context.Background()

	t.Run("add and get", func(t *testing.T) {
		c := newResultsCache(ctx, time.Hour, time.Hour)
		c.add(123, []msgResult{{}, {}, {}})
		v, exists := c.get(123)
		assert.True(t, exists)
		assert.Equal(t, []msgResult{{}, {}, {}}, v)
	})

	t.Run("expired", func(t *testing.T) {
		c := newResultsCache(ctx, time.Millisecond, time.Millisecond)
		c.add(123, []msgResult{{}, {}, {}})
		time.Sleep(10 * time.Millisecond)
		_, exists := c.get(123)
		assert.False(t, exists)
	})

}
