package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

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
	env, err := commonchangeset.Apply(t, env, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployCapabilityRegistry),
			registrySel,
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployOCR3),
			registrySel,
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployForwarder),
			changeset.DeployForwarderRequest{},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployFeedsConsumer),
			&changeset.DeployFeedsConsumerRequest{ChainSelector: registrySel},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			map[uint64]types.MCMSWithTimelockConfigV2{
				registrySel: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		),
	)
	require.NoError(t, err)
	addrs, err := env.ExistingAddresses.AddressesForChain(registrySel)
	require.NoError(t, err)
	timelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(env.Chains[registrySel], addrs)
	require.NoError(t, err)

	_, err = commonchangeset.Apply(t, env,
		map[uint64]*proposalutils.TimelockExecutionContracts{
			registrySel: {Timelock: timelock.Timelock, CallProxy: timelock.CallProxy},
		},
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.AcceptAllOwnershipsProposal),
			&changeset.AcceptAllOwnershipRequest{
				ChainSelector: registrySel,
				MinDelay:      0,
			},
		),
	)
	require.NoError(t, err)
}
