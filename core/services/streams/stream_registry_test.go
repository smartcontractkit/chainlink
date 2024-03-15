package streams

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStream struct {
	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error
}

func (m *mockStream) Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	return m.run, m.trrs, m.err
}

func Test_Registry(t *testing.T) {
	lggr := logger.TestLogger(t)
	runner := &mockRunner{}

	t.Run("Get", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		sr.streams["foo"] = &mockStream{run: &pipeline.Run{ID: 1}}
		sr.streams["bar"] = &mockStream{run: &pipeline.Run{ID: 2}}
		sr.streams["baz"] = &mockStream{run: &pipeline.Run{ID: 3}}

		v, exists := sr.Get("foo")
		assert.True(t, exists)
		assert.Equal(t, sr.streams["foo"], v)

		v, exists = sr.Get("bar")
		assert.True(t, exists)
		assert.Equal(t, sr.streams["bar"], v)

		v, exists = sr.Get("baz")
		assert.True(t, exists)
		assert.Equal(t, sr.streams["baz"], v)

		v, exists = sr.Get("qux")
		assert.Nil(t, v)
		assert.False(t, exists)
	})
	t.Run("Register", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		t.Run("registers new stream", func(t *testing.T) {
			assert.Len(t, sr.streams, 0)
			err := sr.Register("foo", pipeline.Spec{ID: 32, DotDagSource: "source"}, nil)
			require.NoError(t, err)
			assert.Len(t, sr.streams, 1)

			v, exists := sr.Get("foo")
			require.True(t, exists)
			strm := v.(*stream)
			assert.Equal(t, StreamID("foo"), strm.id)
			assert.Equal(t, int32(32), strm.spec.ID)
		})

		t.Run("errors when attempt to re-register a stream with an existing ID", func(t *testing.T) {
			assert.Len(t, sr.streams, 1)
			err := sr.Register("foo", pipeline.Spec{ID: 33, DotDagSource: "source"}, nil)
			require.Error(t, err)
			assert.Len(t, sr.streams, 1)
			assert.EqualError(t, err, "stream already registered for id: \"foo\"")

			v, exists := sr.Get("foo")
			require.True(t, exists)
			strm := v.(*stream)
			assert.Equal(t, StreamID("foo"), strm.id)
			assert.Equal(t, int32(32), strm.spec.ID)
		})
	})
	t.Run("Unregister", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		sr.streams["foo"] = &mockStream{run: &pipeline.Run{ID: 1}}
		sr.streams["bar"] = &mockStream{run: &pipeline.Run{ID: 2}}
		sr.streams["baz"] = &mockStream{run: &pipeline.Run{ID: 3}}

		t.Run("unregisters a stream", func(t *testing.T) {
			assert.Len(t, sr.streams, 3)

			sr.Unregister("foo")

			assert.Len(t, sr.streams, 2)
			_, exists := sr.streams["foo"]
			assert.False(t, exists)
		})
		t.Run("no effect when unregistering a non-existent stream", func(t *testing.T) {
			assert.Len(t, sr.streams, 2)

			sr.Unregister("foo")

			assert.Len(t, sr.streams, 2)
			_, exists := sr.streams["foo"]
			assert.False(t, exists)
		})
	})
}
