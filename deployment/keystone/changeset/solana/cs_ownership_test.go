package solana

import (
	"testing"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test_TransferOwnershipForwarder(t *testing.T) {
	// deploy forwarder
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1, // nodes unused but required in config
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	registrySel := env.AllChainSelectorsSolana()[0]
	ab := cldf.NewMemoryAddressBook()

	env = shouldDeployForwarder(t, env, registrySel, ab)

	result, err := TransferOwnershipForwarder(env, &TransferOwnershipForwarderRequest{
		ChainSel: registrySel,
		MCMSCfg:  proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
	})

	require.NoError(t, err)
	require.NotNil(t, result.MCMSTimelockProposals)
}
