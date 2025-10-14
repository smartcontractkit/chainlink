package changeset_test

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	state2 "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestTransferToMCMSWithTimelockV2(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	// Setup contracts
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployLinkToken), []uint64{selector}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)

	state, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)

	link, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.TransferToMCMSWithTimelockV2), changeset.TransferToMCMSWithTimelockConfig{
			ContractsByChain: map[uint64][]common.Address{
				selector: {link.LinkToken.Address()},
			},
			MCMSConfig: proposalutils.TimelockConfig{
				MinDelay: 0,
			},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)
	require.Len(t, rt.State().Proposals, 1)
	require.True(t, rt.State().Proposals[0].IsExecuted)

	// We expect now that the link token is owned by the MCMS timelock.
	link, err = changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	o, err := link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Timelock.Address(), o)

	// Try a rollback to the deployer.
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.TransferToDeployer), changeset.TransferToDeployerConfig{
			ContractAddress: link.LinkToken.Address(),
			ChainSel:        selector,
		}),
	)
	require.NoError(t, err)

	o, err = link.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, chain.DeployerKey.From, o)
}

func TestTransferToMCMSWithTimelockV2DataStore(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	// Setup contracts
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployLinkToken), []uint64{selector}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)

	state, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)

	link, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	// Remove LinkToken from AddressBook to simulate datastore only having the contract address
	// TODO: migrate DeployLinkToken to use datastore only
	linkAb := cldf.NewMemoryAddressBookFromMap(map[uint64]map[string]cldf.TypeAndVersion{
		selector: {
			link.LinkToken.Address().String(): cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0),
		},
	})
	err = rt.State().AddressBook.Remove(linkAb)
	require.NoError(t, err)

	// Add link token address to the datastore only
	ds := datastore.NewMemoryDataStore()
	typeAndVersion := cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)

	err = ds.Addresses().Add(datastore.AddressRef{
		Address:       link.LinkToken.Address().String(),
		Type:          datastore.ContractType(typeAndVersion.Type.String()),
		ChainSelector: selector,
		Version:       semver.MustParse(typeAndVersion.Version.String()),
	})
	require.NoError(t, err)
	err = ds.Merge(rt.State().DataStore)
	require.NoError(t, err)
	rt.State().DataStore = ds.Seal()
	newEnv := rt.Environment()
	newEnv.DataStore = ds.Seal()
	// Re create runtime with new environment
	rt = runtime.NewFromEnvironment(newEnv)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.TransferToMCMSWithTimelockV2), changeset.TransferToMCMSWithTimelockConfig{
			ContractsByChain: map[uint64][]common.Address{
				selector: {link.LinkToken.Address()},
			},
			MCMSConfig: proposalutils.TimelockConfig{
				MinDelay: 0,
			},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)
	require.Len(t, rt.State().Proposals, 1)
	require.True(t, rt.State().Proposals[0].IsExecuted)

	// We expect now that the link token is owned by the MCMS timelock.
	addrsDatastore, err := state2.LoadAddressesFromDataStore(rt.State().DataStore, selector, "")
	require.NoError(t, err)
	linkState, err := state2.MaybeLoadLinkTokenChainState(chain, addrsDatastore)
	require.NoError(t, err)

	o, err := linkState.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Timelock.Address(), o)

	// Try a rollback to the deployer.
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.TransferToDeployer), changeset.TransferToDeployerConfig{
			ContractAddress: link.LinkToken.Address(),
			ChainSel:        selector,
		}),
	)
	require.NoError(t, err)

	o, err = linkState.LinkToken.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, chain.DeployerKey.From, o)
}

func TestRenounceTimelockDeployerConfigValidate(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-724")
	t.Parallel()

	selector1 := chain_selectors.TEST_90000001.Selector
	selector2 := chain_selectors.TEST_90000002.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector1, selector2}),
	))
	require.NoError(t, err)

	// Deploy MCMS to selector 1 only, so we have a chain without MCMS
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
			selector1: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	for _, test := range []struct {
		name   string
		config changeset.RenounceTimelockDeployerConfig
		env    cldf.Environment
		err    string
	}{
		{
			name: "valid config",
			env:  rt.Environment(),
			config: changeset.RenounceTimelockDeployerConfig{
				ChainSel: selector1,
			},
		},
		{
			name: "invalid chain selector",
			env:  rt.Environment(),
			config: changeset.RenounceTimelockDeployerConfig{
				ChainSel: 0,
			},
			err: "invalid chain selector: chain selector must be set",
		},
		{
			name: "chain does not exists on env",
			env:  rt.Environment(),
			config: changeset.RenounceTimelockDeployerConfig{
				ChainSel: chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
			},
			err: "chain selector: 16015286601757825753 not found in environment",
		},
		{
			name: "no MCMS deployed",
			env:  rt.Environment(),
			config: changeset.RenounceTimelockDeployerConfig{
				ChainSel: selector2,
			},
			// chain does not match any existing addresses
			err: "timelock not found on chain 5548718428018410741",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.config.Validate(test.env)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRenounceTimelockDeployer(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)

	state, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)

	tl := state.Timelock
	require.NotNil(t, tl)

	adminRole, err := tl.ADMINROLE(nil)
	require.NoError(t, err)

	r, err := tl.GetRoleMemberCount(&bind.CallOpts{}, adminRole)
	require.NoError(t, err)
	require.Equal(t, int64(2), r.Int64())

	// Revoke Deployer
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.RenounceTimelockDeployer), changeset.RenounceTimelockDeployerConfig{
			ChainSel: selector,
		}),
	)
	require.NoError(t, err)

	// Check that the deployer is no longer an admin
	r, err = tl.GetRoleMemberCount(&bind.CallOpts{}, adminRole)
	require.NoError(t, err)
	require.Equal(t, int64(1), r.Int64())

	// Retrieve the admin address
	admin, err := tl.GetRoleMember(&bind.CallOpts{}, adminRole, big.NewInt(0))
	require.NoError(t, err)

	// Check that the admin is the timelock
	require.Equal(t, tl.Address(), admin)
}
