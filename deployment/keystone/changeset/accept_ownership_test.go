package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
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
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	env, err := commonchangeset.Apply(t, env, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(changeset.DeployCapabilityRegistry),
		registrySel,
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(changeset.DeployOCR3),
		registrySel,
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(changeset.DeployForwarder),
		changeset.DeployForwarderRequest{},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(changeset.DeployFeedsConsumer),
		&changeset.DeployFeedsConsumerRequest{ChainSelector: registrySel},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		map[uint64]types.MCMSWithTimelockConfigV2{
			registrySel: proposalutils.SingleGroupTimelockConfigV2(t),
		},
	))
	require.NoError(t, err)

	_, err = commonchangeset.Apply(t, env,
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
