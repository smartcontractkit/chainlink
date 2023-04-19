package mercury

import (
	"context"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

func newReportingPlugin(t *testing.T) *reportingPlugin {
	return &reportingPlugin{
		f:                       1,
		onchainConfig:           OnchainConfig{Min: big.NewInt(0), Max: big.NewInt(1000)},
		maxFinalizedBlockNumber: newInitialMaxFinalizedBlockNumber(),
		logger:                  logger.Test(t),
	}
}

func Test_ReportingPlugin_shouldReport(t *testing.T) {
	rp := newReportingPlugin(t)
	repts := types.ReportTimestamp{}
	paos := NewParsedAttributedObservations()

	t.Run("reports if all reports have currentBlockNum > validFromBlockNum", func(t *testing.T) {
		for i := range paos {
			paos[i].CurrentBlockNum = 500
			paos[i].ValidFromBlockNum = 499
		}
		shouldReport, err := rp.shouldReport(context.Background(), repts, paos)
		require.NoError(t, err)

		assert.True(t, shouldReport)
	})
	t.Run("does not report if all reports have currentBlockNum == validFromBlockNum", func(t *testing.T) {
		for i := range paos {
			paos[i].CurrentBlockNum = 500
			paos[i].ValidFromBlockNum = 500
		}
		shouldReport, err := rp.shouldReport(context.Background(), repts, paos)
		require.NoError(t, err)

		assert.False(t, shouldReport)
	})
	t.Run("does not report if all reports have currentBlockNum < validFromBlockNum", func(t *testing.T) {
		paos := NewParsedAttributedObservations()
		for i := range paos {
			paos[i].CurrentBlockNum = 499
			paos[i].ValidFromBlockNum = 500
		}
		shouldReport, err := rp.shouldReport(context.Background(), repts, paos)
		require.NoError(t, err)

		assert.False(t, shouldReport)
	})
	t.Run("returns error if it cannot come to consensus about currentBlockNum", func(t *testing.T) {
		paos := NewParsedAttributedObservations()
		for i := range paos {
			paos[i].CurrentBlockNum = 500 + int64(i)
			paos[i].ValidFromBlockNum = 499
		}
		shouldReport, err := rp.shouldReport(context.Background(), repts, paos)
		require.NoError(t, err)

		assert.False(t, shouldReport)
	})
	t.Run("returns error if it cannot come to consensus about validFromBlockNum", func(t *testing.T) {
		paos := NewParsedAttributedObservations()
		for i := range paos {
			paos[i].CurrentBlockNum = 500
			paos[i].ValidFromBlockNum = 499 - int64(i)
		}
		shouldReport, err := rp.shouldReport(context.Background(), repts, paos)
		require.NoError(t, err)

		assert.False(t, shouldReport)
	})
}

var _ DataSource = &mockDataSource{}

type mockDataSource struct{ obs Observation }

func (m mockDataSource) Observe(context.Context) (Observation, error) {
	return m.obs, nil
}

func randBigInt() *big.Int {
	return big.NewInt(rand.Int63())
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func mustDecodeBigInt(b []byte) *big.Int {
	n, err := DecodeValueInt192(b)
	if err != nil {
		panic(err)
	}
	return n
}

func Test_Plugin_Observation(t *testing.T) {
	ctx := context.Background()
	repts := ocrtypes.ReportTimestamp{}
	query := ocrtypes.Query{}

	rp := newReportingPlugin(t)

	t.Run("when all observations are successful", func(t *testing.T) {
		maxFinalizedBlockNumber := int64(rand.Int31())
		rp.maxFinalizedBlockNumber.Store(maxFinalizedBlockNumber)
		obs := Observation{
			BenchmarkPrice: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Bid: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Ask: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			CurrentBlockNum: ObsResult[int64]{
				Val: maxFinalizedBlockNumber + 2,
			},
			CurrentBlockHash: ObsResult[[]byte]{
				Val: randBytes(32),
			},
		}
		rp.dataSource = mockDataSource{obs}

		pbObs, err := rp.Observation(ctx, repts, query)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(pbObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Equal(t, obs.BenchmarkPrice.Val, mustDecodeBigInt(p.BenchmarkPrice))
		assert.Equal(t, obs.Bid.Val, mustDecodeBigInt(p.Bid))
		assert.Equal(t, obs.Ask.Val, mustDecodeBigInt(p.Ask))
		assert.Equal(t, obs.CurrentBlockNum.Val, p.CurrentBlockNum)
		assert.Equal(t, obs.CurrentBlockHash.Val, p.CurrentBlockHash)
		assert.Equal(t, maxFinalizedBlockNumber+1, p.ValidFromBlockNum)
		assert.True(t, p.BenchmarkPriceValid)
		assert.True(t, p.BidValid)
		assert.True(t, p.AskValid)
		assert.True(t, p.CurrentBlockNumValid)
		assert.True(t, p.CurrentBlockHashValid)
		assert.True(t, p.ValidFromBlockNumValid)
	})

	t.Run("when all observations have failed", func(t *testing.T) {
		obs := Observation{
			// Vals should be ignored, this is asserted with .Zero below
			BenchmarkPrice: ObsResult[*big.Int]{
				Val: randBigInt(),
				Err: errors.New("benchmarkPrice exploded"),
			},
			Bid: ObsResult[*big.Int]{
				Val: randBigInt(),
				Err: errors.New("bid exploded"),
			},
			Ask: ObsResult[*big.Int]{
				Val: randBigInt(),
				Err: errors.New("ask exploded"),
			},
			CurrentBlockNum: ObsResult[int64]{
				Err: errors.New("currentBlockNum exploded"),
				Val: rand.Int63(),
			},
			CurrentBlockHash: ObsResult[[]byte]{
				Err: errors.New("currentBlockHash exploded"),
				Val: randBytes(32),
			},
		}
		maxFinalizedBlockNumber := int64(rand.Int31())
		rp.maxFinalizedBlockNumber.Store(maxFinalizedBlockNumber)
		rp.dataSource = mockDataSource{obs}

		pbObs, err := rp.Observation(ctx, repts, query)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(pbObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Zero(t, p.BenchmarkPrice)
		assert.Zero(t, p.Bid)
		assert.Zero(t, p.Ask)
		assert.Zero(t, p.CurrentBlockNum)
		assert.Zero(t, p.CurrentBlockHash)
		assert.Equal(t, maxFinalizedBlockNumber+1, p.ValidFromBlockNum)
		assert.False(t, p.BenchmarkPriceValid)
		assert.False(t, p.BidValid)
		assert.False(t, p.AskValid)
		assert.False(t, p.CurrentBlockNumValid)
		assert.False(t, p.CurrentBlockHashValid)
		assert.True(t, p.ValidFromBlockNumValid)
	})

	t.Run("if maxFinalizedBlockNumber has not been set", func(t *testing.T) {
		rp.maxFinalizedBlockNumber.Store(unfetchedInitialMaxFinalizedBlockNumber)
		obs := Observation{
			BenchmarkPrice: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Bid: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Ask: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			CurrentBlockNum: ObsResult[int64]{
				Val: int64(rand.Int31()),
			},
			CurrentBlockHash: ObsResult[[]byte]{
				Val: randBytes(32),
			},
		}
		rp.dataSource = mockDataSource{obs}

		pbObs, err := rp.Observation(ctx, repts, query)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(pbObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Equal(t, obs.BenchmarkPrice.Val, mustDecodeBigInt(p.BenchmarkPrice))
		assert.Equal(t, obs.Bid.Val, mustDecodeBigInt(p.Bid))
		assert.Equal(t, obs.Ask.Val, mustDecodeBigInt(p.Ask))
		assert.Equal(t, obs.CurrentBlockNum.Val, p.CurrentBlockNum)
		assert.Equal(t, obs.CurrentBlockHash.Val, p.CurrentBlockHash)
		assert.Zero(t, p.ValidFromBlockNum)
		assert.True(t, p.BenchmarkPriceValid)
		assert.True(t, p.BidValid)
		assert.True(t, p.AskValid)
		assert.True(t, p.CurrentBlockNumValid)
		assert.True(t, p.CurrentBlockHashValid)
		assert.False(t, p.ValidFromBlockNumValid)
	})

	t.Run("when some observations have failed", func(t *testing.T) {
		maxFinalizedBlockNumber := int64(rand.Int31())
		rp.maxFinalizedBlockNumber.Store(maxFinalizedBlockNumber)
		obs := Observation{
			BenchmarkPrice: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Bid: ObsResult[*big.Int]{
				Val: randBigInt(),
			},
			Ask: ObsResult[*big.Int]{
				Err: errors.New("ask exploded"),
			},
			CurrentBlockNum: ObsResult[int64]{
				Err: errors.New("currentBlockNum exploded"),
			},
			CurrentBlockHash: ObsResult[[]byte]{
				Val: randBytes(32),
			},
		}
		rp.dataSource = mockDataSource{obs}

		pbObs, err := rp.Observation(ctx, repts, query)
		require.NoError(t, err)

		var p MercuryObservationProto
		require.NoError(t, proto.Unmarshal(pbObs, &p))

		assert.LessOrEqual(t, p.Timestamp, uint32(time.Now().Unix()))
		assert.Equal(t, obs.BenchmarkPrice.Val, mustDecodeBigInt(p.BenchmarkPrice))
		assert.Equal(t, obs.Bid.Val, mustDecodeBigInt(p.Bid))
		assert.Zero(t, p.Ask)
		assert.Zero(t, p.CurrentBlockNum)
		assert.Equal(t, obs.CurrentBlockHash.Val, p.CurrentBlockHash)
		assert.Equal(t, maxFinalizedBlockNumber+1, p.ValidFromBlockNum)
		assert.True(t, p.BenchmarkPriceValid)
		assert.True(t, p.BidValid)
		assert.False(t, p.AskValid)
		assert.False(t, p.CurrentBlockNumValid)
		assert.True(t, p.CurrentBlockHashValid)
		assert.True(t, p.ValidFromBlockNumValid)
	})
}
