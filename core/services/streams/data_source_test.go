package streams

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-data-streams/streams"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockStream struct {
	dp  DataPoint
	err error
}

func (m *mockStream) Observe(ctx context.Context) (DataPoint, error) {
	return m.dp, m.err
}

func Test_DataSource(t *testing.T) {
	lggr := logger.TestLogger(t)
	sc := newStreamCache(nil)
	ds := NewDataSource(lggr, sc)
	ctx := testutils.Context(t)

	streamIDs := make(map[streams.StreamID]struct{})
	streamIDs[streams.StreamID("ETH/USD")] = struct{}{}
	streamIDs[streams.StreamID("BTC/USD")] = struct{}{}
	streamIDs[streams.StreamID("LINK/USD")] = struct{}{}

	t.Run("Observe", func(t *testing.T) {
		t.Run("returns errors if no streams are defined", func(t *testing.T) {
			vals, err := ds.Observe(ctx, streamIDs)
			assert.NoError(t, err)

			assert.Equal(t, streams.StreamValues{
				"BTC/USD":  streams.ObsResult[*big.Int]{Val: nil, Err: ErrMissingStream{id: "BTC/USD"}},
				"ETH/USD":  streams.ObsResult[*big.Int]{Val: nil, Err: ErrMissingStream{id: "ETH/USD"}},
				"LINK/USD": streams.ObsResult[*big.Int]{Val: nil, Err: ErrMissingStream{id: "LINK/USD"}},
			}, vals)
		})
		t.Run("observes each stream with success and returns values matching map argument", func(t *testing.T) {
			sc.streams["ETH/USD"] = &mockStream{
				dp: big.NewInt(2181),
			}
			sc.streams["BTC/USD"] = &mockStream{
				dp: big.NewInt(40602),
			}
			sc.streams["LINK/USD"] = &mockStream{
				dp: big.NewInt(15),
			}

			vals, err := ds.Observe(ctx, streamIDs)
			assert.NoError(t, err)

			assert.Equal(t, streams.StreamValues{
				"BTC/USD":  streams.ObsResult[*big.Int]{Val: big.NewInt(40602), Err: nil},
				"ETH/USD":  streams.ObsResult[*big.Int]{Val: big.NewInt(2181), Err: nil},
				"LINK/USD": streams.ObsResult[*big.Int]{Val: big.NewInt(15), Err: nil},
			}, vals)
		})
		t.Run("observes each stream and returns success/errors", func(t *testing.T) {
			sc.streams["ETH/USD"] = &mockStream{
				dp:  big.NewInt(2181),
				err: errors.New("something exploded"),
			}
			sc.streams["BTC/USD"] = &mockStream{
				dp: big.NewInt(40602),
			}
			sc.streams["LINK/USD"] = &mockStream{
				err: errors.New("something exploded 2"),
			}

			vals, err := ds.Observe(ctx, streamIDs)
			assert.NoError(t, err)

			assert.Equal(t, streams.StreamValues{
				"BTC/USD":  streams.ObsResult[*big.Int]{Val: big.NewInt(40602), Err: nil},
				"ETH/USD":  streams.ObsResult[*big.Int]{Val: big.NewInt(2181), Err: errors.New("something exploded")},
				"LINK/USD": streams.ObsResult[*big.Int]{Val: nil, Err: errors.New("something exploded 2")},
			}, vals)
		})
	})
}
