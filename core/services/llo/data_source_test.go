package llo

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
)

type mockStream struct {
	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error
}

func (m *mockStream) Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	return m.run, m.trrs, m.err
}

type mockRegistry struct {
	streams map[streams.StreamID]*mockStream
}

func (m *mockRegistry) Get(streamID streams.StreamID) (strm streams.Stream, exists bool) {
	strm, exists = m.streams[streamID]
	return
}

func makeStreamWithSingleResult[T any](res T, err error) *mockStream {
	return &mockStream{
		trrs: []pipeline.TaskRunResult{pipeline.TaskRunResult{Task: &pipeline.MemoTask{}, Result: pipeline.Result{Value: res}}},
		err:  err,
	}
}

func makeStreamValues() llo.StreamValues {
	return llo.StreamValues{
		1: nil,
		2: nil,
		3: nil,
	}
}

type mockOpts struct{}

func (m mockOpts) VerboseLogging() bool { return true }
func (m mockOpts) SeqNr() uint64        { return 42 }

func Test_DataSource(t *testing.T) {
	t.Skip("waiting on https://github.com/smartcontractkit/chainlink/pull/13780")
	lggr := logger.TestLogger(t)
	reg := &mockRegistry{make(map[streams.StreamID]*mockStream)}
	ds := newDataSource(lggr, reg)
	ctx := testutils.Context(t)

	t.Run("Observe", func(t *testing.T) {
		t.Run("doesn't set any values if no streams are defined", func(t *testing.T) {
			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, mockOpts{})
			assert.NoError(t, err)

			assert.Equal(t, makeStreamValues(), vals)
		})
		t.Run("observes each stream with success and returns values matching map argument", func(t *testing.T) {
			reg.streams[1] = makeStreamWithSingleResult[*big.Int](big.NewInt(2181), nil)
			reg.streams[2] = makeStreamWithSingleResult[*big.Int](big.NewInt(40602), nil)
			reg.streams[3] = makeStreamWithSingleResult[*big.Int](big.NewInt(15), nil)

			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, mockOpts{})
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{}, vals)
		})
		t.Run("observes each stream and returns success/errors", func(t *testing.T) {
			reg.streams[1] = makeStreamWithSingleResult[*big.Int](big.NewInt(2181), errors.New("something exploded"))
			reg.streams[2] = makeStreamWithSingleResult[*big.Int](big.NewInt(40602), nil)
			reg.streams[3] = makeStreamWithSingleResult[*big.Int](nil, errors.New("something exploded 2"))

			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, mockOpts{})
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{}, vals)
		})
	})
}
