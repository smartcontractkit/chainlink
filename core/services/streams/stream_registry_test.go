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

		sr.streams[1] = &mockStream{run: &pipeline.Run{ID: 1}}
		sr.streams[2] = &mockStream{run: &pipeline.Run{ID: 2}}
		sr.streams[3] = &mockStream{run: &pipeline.Run{ID: 3}}

		v, exists := sr.Get(1)
		assert.True(t, exists)
		assert.Equal(t, sr.streams[1], v)

		v, exists = sr.Get(2)
		assert.True(t, exists)
		assert.Equal(t, sr.streams[2], v)

		v, exists = sr.Get(3)
		assert.True(t, exists)
		assert.Equal(t, sr.streams[3], v)

		v, exists = sr.Get(4)
		assert.Nil(t, v)
		assert.False(t, exists)
	})
	t.Run("Register", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		t.Run("registers new stream", func(t *testing.T) {
			assert.Len(t, sr.streams, 0)
			err := sr.Register(1, pipeline.Spec{ID: 32, DotDagSource: "source"}, nil)
			require.NoError(t, err)
			assert.Len(t, sr.streams, 1)

			v, exists := sr.Get(1)
			require.True(t, exists)
			strm := v.(*stream)
			assert.Equal(t, StreamID(1), strm.id)
			assert.Equal(t, int32(32), strm.spec.ID)
		})

		t.Run("errors when attempt to re-register a stream with an existing ID", func(t *testing.T) {
			assert.Len(t, sr.streams, 1)
			err := sr.Register(1, pipeline.Spec{ID: 33, DotDagSource: "source"}, nil)
			require.Error(t, err)
			assert.Len(t, sr.streams, 1)
			assert.EqualError(t, err, "stream already registered for id: 1")

			v, exists := sr.Get(1)
			require.True(t, exists)
			strm := v.(*stream)
			assert.Equal(t, StreamID(1), strm.id)
			assert.Equal(t, int32(32), strm.spec.ID)
		})
	})
	t.Run("Unregister", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		sr.streams[1] = &mockStream{run: &pipeline.Run{ID: 1}}
		sr.streams[2] = &mockStream{run: &pipeline.Run{ID: 2}}
		sr.streams[3] = &mockStream{run: &pipeline.Run{ID: 3}}

		t.Run("unregisters a stream", func(t *testing.T) {
			assert.Len(t, sr.streams, 3)

			sr.Unregister(1)

			assert.Len(t, sr.streams, 2)
			_, exists := sr.streams[1]
			assert.False(t, exists)
		})
		t.Run("no effect when unregistering a non-existent stream", func(t *testing.T) {
			assert.Len(t, sr.streams, 2)

			sr.Unregister(1)

			assert.Len(t, sr.streams, 2)
			_, exists := sr.streams[1]
			assert.False(t, exists)
		})
	})
}
