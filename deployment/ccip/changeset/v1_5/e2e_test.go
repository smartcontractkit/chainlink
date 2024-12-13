package v1_5

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestE2ELegacy(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains:             2,
		NumOfUsersPerChain: 1,
		Nodes:              4,
		Bootstraps:         1,
	})
	// prep
}
