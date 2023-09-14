package streams

import (
	"context"
	"errors"
	"maps"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-data-streams/streams"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/assert"
)

type mockORM struct {
	m   map[streams.StreamID]Stream
	err error
}

func (orm *mockORM) LoadStreams(ctx context.Context, lggr logger.Logger, runner Runner, m map[streams.StreamID]Stream) error {
	maps.Copy(m, orm.m)
	return orm.err
}

func Test_StreamCache(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		orm := &mockORM{}
		sc := newStreamCache(orm)
		lggr := logger.TestLogger(t)
		runner := &mockRunner{}

		t.Run("populates cache from database using ORM", func(t *testing.T) {
			assert.Len(t, sc.streams, 0)
			err := sc.Load(testutils.Context(t), lggr, runner)
			assert.NoError(t, err)
			assert.Len(t, sc.streams, 0)

			v, exists := sc.Get("foo")
			assert.Nil(t, v)
			assert.False(t, exists)

			orm.m = make(map[streams.StreamID]Stream)
			orm.m["foo"] = &mockStream{dp: big.NewInt(1)}
			orm.m["bar"] = &mockStream{dp: big.NewInt(2)}
			orm.m["baz"] = &mockStream{dp: big.NewInt(3)}

			err = sc.Load(testutils.Context(t), lggr, runner)
			assert.NoError(t, err)
			assert.Len(t, sc.streams, 3)

			v, exists = sc.Get("foo")
			assert.True(t, exists)
			assert.Equal(t, orm.m["foo"], v)

			v, exists = sc.Get("bar")
			assert.True(t, exists)
			assert.Equal(t, orm.m["bar"], v)

			v, exists = sc.Get("baz")
			assert.True(t, exists)
			assert.Equal(t, orm.m["baz"], v)
		})

		t.Run("returns error if db errors", func(t *testing.T) {
			orm.err = errors.New("something exploded")
			err := sc.Load(testutils.Context(t), lggr, runner)
			assert.EqualError(t, err, "something exploded")
		})
	})
}
