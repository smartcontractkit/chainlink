package llo

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

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

func makeStreamWithSingleResult[T any](runID int64, res T, err error) *mockStream {
	return &mockStream{
		run:  &pipeline.Run{ID: runID},
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

func (m *mockOpts) VerboseLogging() bool { return true }
func (m *mockOpts) SeqNr() uint64        { return 1042 }
func (m *mockOpts) OutCtx() ocr3types.OutcomeContext {
	return ocr3types.OutcomeContext{SeqNr: 1042, PreviousOutcome: ocr3types.Outcome([]byte("foo"))}
}
func (m *mockOpts) ConfigDigest() ocr2types.ConfigDigest {
	return ocr2types.ConfigDigest{6, 5, 4}
}

type mockTelemeter struct {
	mu                     sync.Mutex
	v3PremiumLegacyPackets []v3PremiumLegacyPacket
}

type v3PremiumLegacyPacket struct {
	run      *pipeline.Run
	trrs     pipeline.TaskRunResults
	streamID uint32
	opts     llo.DSOpts
	val      llo.StreamValue
	err      error
}

var _ Telemeter = &mockTelemeter{}

func (m *mockTelemeter) EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v3PremiumLegacyPackets = append(m.v3PremiumLegacyPackets, v3PremiumLegacyPacket{run, trrs, streamID, opts, val, err})
}

func Test_DataSource(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := &mockRegistry{make(map[streams.StreamID]*mockStream)}
	ds := newDataSource(lggr, reg, NullTelemeter)
	ctx := testutils.Context(t)
	opts := &mockOpts{}

	t.Run("Observe", func(t *testing.T) {
		t.Run("doesn't set any values if no streams are defined", func(t *testing.T) {
			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, opts)
			assert.NoError(t, err)

			assert.Equal(t, makeStreamValues(), vals)
		})
		t.Run("observes each stream with success and returns values matching map argument", func(t *testing.T) {
			reg.streams[1] = makeStreamWithSingleResult[*big.Int](1, big.NewInt(2181), nil)
			reg.streams[2] = makeStreamWithSingleResult[*big.Int](2, big.NewInt(40602), nil)
			reg.streams[3] = makeStreamWithSingleResult[*big.Int](3, big.NewInt(15), nil)

			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, opts)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				2: llo.ToDecimal(decimal.NewFromInt(40602)),
				1: llo.ToDecimal(decimal.NewFromInt(2181)),
				3: llo.ToDecimal(decimal.NewFromInt(15)),
			}, vals)
		})
		t.Run("observes each stream and returns success/errors", func(t *testing.T) {
			reg.streams[1] = makeStreamWithSingleResult[*big.Int](1, big.NewInt(2181), errors.New("something exploded"))
			reg.streams[2] = makeStreamWithSingleResult[*big.Int](2, big.NewInt(40602), nil)
			reg.streams[3] = makeStreamWithSingleResult[*big.Int](3, nil, errors.New("something exploded 2"))

			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, opts)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				2: llo.ToDecimal(decimal.NewFromInt(40602)),
				1: nil,
				3: nil,
			}, vals)
		})

		t.Run("records telemetry", func(t *testing.T) {
			tm := &mockTelemeter{}
			ds.t = tm

			reg.streams[1] = makeStreamWithSingleResult[*big.Int](100, big.NewInt(2181), nil)
			reg.streams[2] = makeStreamWithSingleResult[*big.Int](101, big.NewInt(40602), nil)
			reg.streams[3] = makeStreamWithSingleResult[*big.Int](102, big.NewInt(15), nil)

			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, opts)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				2: llo.ToDecimal(decimal.NewFromInt(40602)),
				1: llo.ToDecimal(decimal.NewFromInt(2181)),
				3: llo.ToDecimal(decimal.NewFromInt(15)),
			}, vals)

			require.Len(t, tm.v3PremiumLegacyPackets, 3)
			m := make(map[int]v3PremiumLegacyPacket)
			for _, pkt := range tm.v3PremiumLegacyPackets {
				m[int(pkt.run.ID)] = pkt
			}
			pkt := m[100]
			assert.Equal(t, 100, int(pkt.run.ID))
			assert.Len(t, pkt.trrs, 1)
			assert.Equal(t, 1, int(pkt.streamID))
			assert.Equal(t, opts, pkt.opts)
			assert.Equal(t, "2181", pkt.val.(*llo.Decimal).String())
			assert.Nil(t, pkt.err)
		})
	})
}
