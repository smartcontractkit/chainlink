package mercury

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
