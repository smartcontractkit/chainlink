package changeset

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"
	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

type testFixture struct {
	rt                        *runtime.Runtime
	selector                  uint64
	workflowRegistryAddress   string
	workflowRegistryQualifier string
}

// setupTest sets up an runtime with a single EVM simulated chain and a deployed WorkflowRegistry.
func setupTest(t *testing.T) *testFixture {
	var (
		qualifier = "test-workflow-registry-v2"
		selector  = chainselectors.GETH_TESTNET.Selector
	)

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(DeployWorkflowRegistry{}, DeployWorkflowRegistryInput{
			ChainSelector: selector,
			Qualifier:     qualifier,
		}),
	)
	require.NoError(t, err, "failed to deploy WorkflowRegistry")

	workflowRegistryAddress := rt.State().DataStore.Addresses().Filter(
		datastore.AddressRefByQualifier("test-workflow-registry-v2"),
	)[0].Address

	// UpdateAllowedSigners to include the allowed owner key
	allowedSigner := rt.Environment().BlockChains.EVMChains()[selector].DeployerKey.From
	err = rt.Exec(
		runtime.ChangesetTask(UpdateAllowedSigners{}, UpdateAllowedSignersInput{
			ChainSelector:             selector,
			WorkflowRegistryQualifier: qualifier,
			Signers: []common.Address{
				allowedSigner,
			},
			Allowed: new(true),
		}),
	)
	require.NoError(t, err, "failed to update allowed signers")
	t.Logf("Updated allowed signers to include key %s", allowedSigner)

	// SetDONLimit to allow workflows to be created
	err = rt.Exec(
		runtime.ChangesetTask(SetDONLimit{}, SetDONLimitInput{
			ChainSelector:             selector,
			WorkflowRegistryQualifier: qualifier,
			DONFamily:                 "zone-a",
			DONLimit:                  new(uint32(100)),
			UserDefaultLimit:          new(uint32(100)),
			MCMSConfig:                nil,
		},
		),
	)
	require.NoError(t, err, "failed to set DON limit")

	return &testFixture{
		rt:                        rt,
		selector:                  selector,
		workflowRegistryAddress:   workflowRegistryAddress,
		workflowRegistryQualifier: qualifier,
	}
}

func setupTestWithMCMS(t *testing.T) *testFixture {
	fixture := setupTest(t)
	handOverToMCMS(t, fixture)

	return fixture
}

// handOverToMCMS deploys MCMS and makes the timelock the real owner of the registry. Tests that
// still need to send deployer transactions to the registry must do so before calling it.
func handOverToMCMS(t *testing.T, fixture *testFixture) {
	t.Helper()

	err := fixture.rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(mcmschangesets.DeployMCMSWithTimelockV2), map[uint64]cldfproposalutils.MCMSWithTimelockConfig{
			fixture.selector: cldftesthelpers.SingleGroupTimelockConfig(t),
		}),
	)
	require.NoError(t, err, "failed to deploy MCMS")

	// Deployer-side setup is done in setupTest; now the timelock really takes ownership.
	test.TransferOwnershipToMCMS(t, fixture.rt, fixture.selector)
}

// linkTimelockAsOwner links the timelock as a workflow owner, through a real MCMS proposal that is
// signed and executed. The registry only accepts user workflow calls from a linked owner, so a test
// that sends them through MCMS needs this first.
func linkTimelockAsOwner(t *testing.T, fixture *testFixture) {
	t.Helper()

	mcmsConfig := &crecontracts.MCMSConfig{
		MinDelay:                  0, // the simulated chain does not advance time on its own
		TimelockQualifierPerChain: map[uint64]string{fixture.selector: ""},
	}
	env := fixture.rt.Environment()
	chain := env.BlockChains.EVMChains()[fixture.selector]
	mcmsContracts, err := strategies.GetMCMSContracts(env, fixture.selector, *mcmsConfig)
	require.NoError(t, err)
	timelock := mcmsContracts.Timelock.Address()

	validity, proof, signature := generateAndSignOwnershipProof(
		t,
		common.HexToAddress(fixture.workflowRegistryAddress),
		timelock.Hex(),
		chain,
		chain.DeployerKey.From.Hex(),
		"123",
		"12",
		"WorkflowRegistry 2.0.0",
		0, // 0 for linking
	)
	err = fixture.rt.Exec(
		runtime.ChangesetTask(UserLinkOwner{}, UserLinkOwnerInput{
			ValidityTimestamp:         validity,
			Proof:                     common.Bytes2Hex(proof.Bytes()),
			Signature:                 common.Bytes2Hex(signature),
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			MCMSConfig:                mcmsConfig,
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
	)
	require.NoError(t, err, "failed to link the timelock as a workflow owner")
}
