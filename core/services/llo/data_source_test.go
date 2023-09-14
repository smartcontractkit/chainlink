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

// func (m *mockRegistry) Register(streamID StreamID, spec pipeline.Spec, rrs ResultRunSaver) error {
//     return nil
// }
// func (m *mockRegistry) Unregister(streamID StreamID) {}

func makeStreamWithSingleResult[T any](res T, err error) *mockStream {
	return &mockStream{
		trrs: []pipeline.TaskRunResult{pipeline.TaskRunResult{Result: pipeline.Result{Value: res}}},
		err:  err,
	}
}

func Test_DataSource(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := &mockRegistry{make(map[streams.StreamID]*mockStream)}
	ds := NewDataSource(lggr, reg)
	ctx := testutils.Context(t)

	streamIDs := make(map[streams.StreamID]struct{})
	streamIDs[streams.StreamID("ETH/USD")] = struct{}{}
	streamIDs[streams.StreamID("BTC/USD")] = struct{}{}
	streamIDs[streams.StreamID("LINK/USD")] = struct{}{}

	t.Run("Observe", func(t *testing.T) {
		t.Run("returns errors if no streams are defined", func(t *testing.T) {
			vals, err := ds.Observe(ctx, streamIDs)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				"BTC/USD":  llo.ObsResult[*big.Int]{Val: nil, Err: ErrMissingStream{id: "BTC/USD"}},
				"ETH/USD":  llo.ObsResult[*big.Int]{Val: nil, Err: ErrMissingStream{id: "ETH/USD"}},
				"LINK/USD": llo.ObsResult[*big.Int]{Val: nil, Err: ErrMissingStream{id: "LINK/USD"}},
			}, vals)
		})
		t.Run("observes each stream with success and returns values matching map argument", func(t *testing.T) {
			reg.streams["ETH/USD"] = makeStreamWithSingleResult[*big.Int](big.NewInt(2181), nil)
			reg.streams["BTC/USD"] = makeStreamWithSingleResult[*big.Int](big.NewInt(40602), nil)
			reg.streams["LINK/USD"] = makeStreamWithSingleResult[*big.Int](big.NewInt(15), nil)

			vals, err := ds.Observe(ctx, streamIDs)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				"BTC/USD":  llo.ObsResult[*big.Int]{Val: big.NewInt(40602), Err: nil},
				"ETH/USD":  llo.ObsResult[*big.Int]{Val: big.NewInt(2181), Err: nil},
				"LINK/USD": llo.ObsResult[*big.Int]{Val: big.NewInt(15), Err: nil},
			}, vals)
		})
		t.Run("observes each stream and returns success/errors", func(t *testing.T) {
			reg.streams["ETH/USD"] = makeStreamWithSingleResult[*big.Int](big.NewInt(2181), errors.New("something exploded"))
			reg.streams["BTC/USD"] = makeStreamWithSingleResult[*big.Int](big.NewInt(40602), nil)
			reg.streams["LINK/USD"] = makeStreamWithSingleResult[*big.Int](nil, errors.New("something exploded 2"))

			vals, err := ds.Observe(ctx, streamIDs)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				"BTC/USD":  llo.ObsResult[*big.Int]{Val: big.NewInt(40602), Err: nil},
				"ETH/USD":  llo.ObsResult[*big.Int]{Val: big.NewInt(2181), Err: errors.New("something exploded")},
				"LINK/USD": llo.ObsResult[*big.Int]{Val: nil, Err: errors.New("something exploded 2")},
			}, vals)
		})
	})
}
