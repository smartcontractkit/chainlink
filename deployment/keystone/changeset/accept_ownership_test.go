package changeset_test

import (
	"math/big"
	"testing"

	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

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
			Config:    registrySel,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config: map[uint64]types.MCMSWithTimelockConfig{
				registrySel: {
					Canceller:         commonchangeset.SingleGroupMCMS(t),
					Bypasser:          commonchangeset.SingleGroupMCMS(t),
					Proposer:          commonchangeset.SingleGroupMCMS(t),
					TimelockExecutors: env.AllDeployerKeys(),
					TimelockMinDelay:  big.NewInt(0),
				},
			},
		},
	})
	require.NoError(t, err)
	addrs, err := env.ExistingAddresses.AddressesForChain(registrySel)
	require.NoError(t, err)
	timelock, err := commonchangeset.MaybeLoadMCMSWithTimelockState(env.Chains[registrySel], addrs)
	require.NoError(t, err)

	_, err = commonchangeset.ApplyChangesets(t, env, map[uint64]*owner_helpers.RBACTimelock{
		registrySel: timelock.Timelock,
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
