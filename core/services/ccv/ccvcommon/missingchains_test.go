package ccvcommon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestNewMissingChainsMonitor_noMissingChains(t *testing.T) {
	m, err := NewMissingChainsMonitor(logger.TestLogger(t), "job", nil, time.Second)
	require.NoError(t, err)
	require.Nil(t, m)
}

func TestMissingChainsMonitor_reportsEachChainRepeatedly(t *testing.T) {
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
		require.NotEmpty(t, logs.FilterField(zapcore.Field{
			Key:     "chainSelector",
			Type:    zapcore.Uint64Type,
			Integer: int64(sel),
		}).All(), "expected a critical log for chain selector %d", sel)
	}

	require.NoError(t, m.Ready(), "a missing chain must not make the monitor unhealthy")
	require.NoError(t, m.HealthReport()[m.Name()])
}
