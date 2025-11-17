package changeset_test

import (
	"crypto/ecdsa"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestAcceptAllOwnership(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{registrySel}),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployCapabilityRegistryV2), &changeset.DeployRequestV2{
			ChainSel: registrySel,
		}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployOCR3V2), &changeset.DeployRequestV2{
			ChainSel: registrySel,
		}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployForwarder), changeset.DeployForwarderRequest{}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployFeedsConsumer), &changeset.DeployFeedsConsumerRequest{
			ChainSelector: registrySel,
		}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
			registrySel: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.AcceptAllOwnershipsProposal), &changeset.AcceptAllOwnershipRequest{
			ChainSelector: registrySel,
			MinDelay:      0,
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)
	require.Len(t, rt.State().Proposals, 1)
	require.True(t, rt.State().Proposals[0].IsExecuted)
}
