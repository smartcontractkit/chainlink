package ccipdata

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_parseLogs(t *testing.T) {
	// generate 100 logs
	logs := make([]logpoller.Log, 100)
	for i := range logs {
		logs[i].LogIndex = int64(i + 1)
		logs[i].BlockNumber = int64(i) * 1000
		logs[i].BlockTimestamp = time.Now()
	}

	parseFn := func(log types.Log) (*uint, error) {
		return &log.Index, nil
	}

	parsedEvents, err := ParseLogs[uint](logs, logger.TestLogger(t), parseFn)
	require.NoError(t, err)
	assert.Len(t, parsedEvents, 100)

	// Make sure everything is parsed according to the parse func
	for i, ev := range parsedEvents {
		assert.Equal(t, i+1, int(ev.Data))
		assert.Equal(t, i*1000, int(ev.BlockNumber))
		assert.Greater(t, ev.BlockTimestampUnixMilli, time.Now().Add(-time.Minute).UnixMilli())
	}
}

func Test_parseLogs_withErrors(t *testing.T) {
	// generate 50 valid logs and 50 errors
	actualErrorCount := 50
	logs := make([]logpoller.Log, actualErrorCount*2)
	for i := range logs {
		logs[i].LogIndex = int64(i + 1)
	}

	// return an error for half of the logs.
	parseFn := func(log types.Log) (*uint, error) {
		if log.Index%2 == 0 {
			return nil, fmt.Errorf("cannot parse %d", log.Index)
		}
		return &log.Index, nil
	}

	log, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	parsedEvents, err := ParseLogs[uint](logs, log, parseFn)
	assert.ErrorContains(t, err, fmt.Sprintf("%d logs were not parsed", len(logs)/2))
	assert.Nil(t, parsedEvents, "No events are returned if there was an error.")

	// logs are written for errors.
	require.Equal(t, actualErrorCount, observed.Len(), "Expect 51 warnings: one for each error and a summary.")
	for i, entry := range observed.All() {
		assert.Equal(t, zapcore.ErrorLevel, entry.Level)
		assert.Contains(t, entry.Message, "Unable to parse log")
		contextMap := entry.ContextMap()
		require.Contains(t, contextMap, "err")
		assert.Contains(t, contextMap["err"], fmt.Sprintf("cannot parse %d", (i+1)*2), "each error should be logged as a warning")
	}
}
