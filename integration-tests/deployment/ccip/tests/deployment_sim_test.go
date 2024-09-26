package tests

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployCCIPContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	DeployCCIPContractsTest(t, e)
}

func TestJobSpecGeneration(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	JobSpecGenerationTest(t, e)
}

func Test0002_InitialDeployInSimulatedBE(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv := NewMemoryEnvironment(t, lggr, 3)
	InitialDeployTest(t, tenv)
}

func TestAddChainInboundMemory(t *testing.T) {
	// 4 chains where the 4th is added after initial deployment.
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 4)
	AddChainInboundTest(t, e)
}

func TestAddLane(t *testing.T) {
	// TODO: The offchain code doesn't yet support partial lane
	// enablement, need to address then re-enable this test.
	t.Skip()
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 3)
	AddLaneTest(t, e)
}
