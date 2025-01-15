package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
	env, err := commonchangeset.ApplyChangesets(t, env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployCapabilityRegistry),
			Config:    registrySel,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployOCR3),
			Config:    registrySel,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployForwarder),
			Config:    changeset.DeployForwarderRequest{},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployFeedsConsumer),
			Config:    &changeset.DeployFeedsConsumerRequest{ChainSelector: registrySel},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config: map[uint64]types.MCMSWithTimelockConfig{
				registrySel: proposalutils.SingleGroupTimelockConfig(t),
			},
		},
	})
	require.NoError(t, err)
	addrs, err := env.ExistingAddresses.AddressesForChain(registrySel)
	require.NoError(t, err)
	timelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(env.Chains[registrySel], addrs)
	require.NoError(t, err)

	_, err = commonchangeset.ApplyChangesets(t, env, map[uint64]*proposalutils.TimelockExecutionContracts{
		registrySel: &proposalutils.TimelockExecutionContracts{
			Timelock:  timelock.Timelock,
			CallProxy: timelock.CallProxy,
		},
	}, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.AcceptAllOwnershipsProposal),
			Config: &changeset.AcceptAllOwnershipRequest{
				ChainSelector: registrySel,
				MinDelay:      0,
			},
		},
	})
	require.NoError(t, err)
}
