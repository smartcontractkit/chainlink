package data_feeds

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestGnosisToMCMS(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	ab := deployment.NewMemoryAddressBook()
	for _, chain := range e.AllChainSelectors() {
		_, err := ccipdeployment.DeployMCMSContracts(lggr, e.Chains[chain], ab, ccipdeployment.MCMSConfig{
			Admin:     ccipdeployment.SingleGroupMCMS(t),
			Canceller: ccipdeployment.SingleGroupMCMS(t),
			Bypasser:  ccipdeployment.SingleGroupMCMS(t),
			Proposer:  ccipdeployment.SingleGroupMCMS(t),
			Executors: []common.Address{e.Chains[chain].DeployerKey.From},
		})
		require.NoError(t, err)
	}
	proposals, err := BuildGnosisProposals(e, ab)
	require.NoError(t, err)
	t.Log(proposals)
	// TODO: Sign in memory, then apply. Ditto for MCMS accept.
}
