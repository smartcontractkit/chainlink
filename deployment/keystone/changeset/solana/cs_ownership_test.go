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
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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
	solChain := env.BlockChains.SolanaChains()[solSel]

	mcfg := commontypes.MCMSWithTimelockConfigV2{
		Canceller:        proposalutils.SingleGroupMCMSV2(t),
		Proposer:         proposalutils.SingleGroupMCMSV2(t),
		Bypasser:         proposalutils.SingleGroupMCMSV2(t),
		TimelockMinDelay: big.NewInt(0),
	}

	env = shouldDeployForwarder(t, env, solSel)

	ds := datastore.NewMemoryDataStore()
	solChain.ProgramsPath = getProgramsPath()
	_, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolanaV2(env, ds, solChain, mcfg)
	require.NoError(t, err)

	err = ds.Merge(env.DataStore)
	require.NoError(t, err)

	env.DataStore = ds.Seal()
	out, err := TransferOwnershipForwarder{}.Apply(env, &TransferOwnershipForwarderRequest{
		ChainSel:  solSel,
		MCMSCfg:   proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
		Qualifier: testQualifier,
		Version:   "1.0.0",
	})
	require.NoError(t, err)
	require.Len(t, out.MCMSTimelockProposals, 1)
}
