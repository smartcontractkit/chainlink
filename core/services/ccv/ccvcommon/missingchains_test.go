package ccvcommon

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestNewMissingChainsMonitor_noMissingChains(t *testing.T) {
	t.Parallel()
	m, err := NewMissingChainsMonitor(logger.TestLogger(t), "job", nil, time.Second)
	require.NoError(t, err)
	require.Nil(t, m)
}

func TestMissingChainsMonitor_reportsEachChainRepeatedly(t *testing.T) {
	t.Parallel()
	lggr, logs := logger.TestLoggerObserved(t, zapcore.DPanicLevel)
	missing := []protocol.ChainSelector{111, 222}

	m, err := NewMissingChainsMonitor(lggr, "CCVExecutor/exec-1", missing, 10*time.Millisecond)
	require.NoError(t, err)
	require.NotNil(t, m)

	servicetest.Run(t, m)

	// One report happens at start, so wait for a second round of both chains.
	require.Eventually(t, func() bool {
		return logs.FilterMessageSnippet("has no chain object on this node").Len() >= 2*len(missing)
	}, 5*time.Second, 10*time.Millisecond)

	for _, sel := range missing {
		var selInt64 int64
		if sel <= math.MaxInt64 {
			selInt64 = int64(sel)
		} else {
			t.Fatalf("chain selector %v is not an int64", sel)
		}

		require.NotEmpty(t, logs.FilterField(zapcore.Field{
			Key:     "chainSelector",
			Type:    zapcore.Uint64Type,
			Integer: selInt64,
		}).All(), "expected a critical log for chain selector %d", sel)
	}

	require.NoError(t, m.Ready(), "a missing chain must not make the monitor unhealthy")
	require.NoError(t, m.HealthReport()[m.Name()])
}
