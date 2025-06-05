package solana

import (
	"math/big"
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func Test_TransferOwnershipForwarder(t *testing.T) {
	skipInCI(t)
	// deploy forwarder
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1, // nodes unused but required in config
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	solSel := env.BlockChains.ListChainSelectors(cldfchain.WithFamily(chain_selectors.FamilySolana))[0]
	ab := cldf.NewMemoryAddressBook()

	mcfg := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		solSel: {
			Canceller:        proposalutils.SingleGroupMCMSV2(t),
			Proposer:         proposalutils.SingleGroupMCMSV2(t),
			Bypasser:         proposalutils.SingleGroupMCMSV2(t),
			TimelockMinDelay: big.NewInt(0),
		},
	}

	env = shouldDeployForwarder(t, env, solSel, ab)

	_, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, env.BlockChains.SolanaChains()[solSel], env.ExistingAddresses, mcfg[solSel])
	require.NoError(t, err)

	result, err := TransferOwnershipForwarder(env, &TransferOwnershipForwarderRequest{
		ChainSel: solSel,
		MCMSCfg:  proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
	})

	require.NoError(t, err)
	require.NotNil(t, result.MCMSTimelockProposals)
}
