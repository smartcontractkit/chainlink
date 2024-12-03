package changeset_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestAcceptAllOwnership(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	registrySel := env.AllChainSelectors()[0]
	chCapReg, err := changeset.DeployCapabilityRegistry(env, registrySel)
	require.NoError(t, err)
	require.NotNil(t, chCapReg)
	err = env.ExistingAddresses.Merge(chCapReg.AddressBook)
	require.NoError(t, err)

	chOcr3, err := changeset.DeployOCR3(env, registrySel)
	require.NoError(t, err)
	require.NotNil(t, chOcr3)
	err = env.ExistingAddresses.Merge(chOcr3.AddressBook)
	require.NoError(t, err)

	chForwarder, err := changeset.DeployForwarder(env, registrySel)
	require.NoError(t, err)
	require.NotNil(t, chForwarder)
	err = env.ExistingAddresses.Merge(chForwarder.AddressBook)
	require.NoError(t, err)

	chConsumer, err := changeset.DeployFeedsConsumer(env, &changeset.DeployFeedsConsumerRequest{
		ChainSelector: registrySel,
	})
	require.NoError(t, err)
	require.NotNil(t, chConsumer)
	err = env.ExistingAddresses.Merge(chConsumer.AddressBook)
	require.NoError(t, err)

	chMcms, err := commonchangeset.DeployMCMSWithTimelock(env, map[uint64]types.MCMSWithTimelockConfig{
		registrySel: {
			Canceller:         commonchangeset.SingleGroupMCMS(t),
			Bypasser:          commonchangeset.SingleGroupMCMS(t),
			Proposer:          commonchangeset.SingleGroupMCMS(t),
			TimelockExecutors: env.AllDeployerKeys(),
			TimelockMinDelay:  big.NewInt(0),
		},
	})
	err = env.ExistingAddresses.Merge(chMcms.AddressBook)
	require.NoError(t, err)

	require.NoError(t, err)
	require.NotNil(t, chMcms)

	resp, err := changeset.TransferAllOwnership(env, &changeset.TransferAllOwnershipRequest{
		ChainSelector: registrySel,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Test the changeset
	output, err := changeset.AcceptAllOwnershipsProposal(env, &changeset.AcceptAllOwnershipRequest{
		ChainSelector: registrySel,
		MinDelay:      time.Duration(0),
	})
	require.NoError(t, err)
	require.NotNil(t, output)
	require.Len(t, output.Proposals, 1)
	proposal := output.Proposals[0]
	require.Len(t, proposal.Transactions, 1)
	txs := proposal.Transactions[0]
	require.Len(t, txs.Batch, 4)
}
