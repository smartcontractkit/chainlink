package example

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDisableRetryExampleChangeset(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	changesetInput := operations.EmptyInput{}
	_, err := DisableRetryExampleChangeset{}.Apply(e, changesetInput)
	require.ErrorContains(t, err, "operation failed")
}

func TestUpdateInputExampleChangeset(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	changesetInput := operations.EmptyInput{}
	_, err := UpdateInputExampleChangeset{}.Apply(e, changesetInput)
	require.NoError(t, err)
}
