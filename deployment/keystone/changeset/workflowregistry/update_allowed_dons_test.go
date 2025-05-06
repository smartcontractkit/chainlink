package workflowregistry_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
)

func TestUpdateAllowedDons(t *testing.T) {
	lggr := logger.Test(t)

	chainSel := chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
	resp := workflowregistry.SetupTestWorkflowRegistry(t, lggr, chainSel)
	registry := resp.Registry

	dons, err := registry.GetAllAllowedDONs(&bind.CallOpts{})
	require.NoError(t, err)

	assert.Empty(t, dons)

	env := deployment.Environment{
		Logger: lggr,
		Chains: map[uint64]deployment.Chain{
			chainSel: resp.Chain,
		},
		ExistingAddresses: resp.AddressBook,
	}

	_, err = workflowregistry.UpdateAllowedDons(
		env,
		&workflowregistry.UpdateAllowedDonsRequest{
			RegistryChainSel: chainSel,
			DonIDs:           []uint32{1},
			Allowed:          true,
		},
	)
	require.NoError(t, err)

	dons, err = registry.GetAllAllowedDONs(&bind.CallOpts{})
	require.NoError(t, err)

	assert.Len(t, dons, 1)
	assert.Equal(t, uint32(1), dons[0])

	_, err = workflowregistry.UpdateAllowedDons(
		env,
		&workflowregistry.UpdateAllowedDonsRequest{
			RegistryChainSel: chainSel,
			DonIDs:           []uint32{1},
			Allowed:          false,
		},
	)
	require.NoError(t, err)

	dons, err = registry.GetAllAllowedDONs(&bind.CallOpts{})
	require.NoError(t, err)

	assert.Empty(t, dons)
}

func Test_UpdateAllowedDons_WithMCMS(t *testing.T) {
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       1,
		UseMCMS:         true,
	})

	req := &workflowregistry.UpdateAllowedDonsRequest{
		RegistryChainSel: te.RegistrySelector,
		DonIDs:           []uint32{1},
		Allowed:          true,
		MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
	}

	out, err := workflowregistry.UpdateAllowedDons(te.Env, req)
	require.NoError(t, err)
	require.Len(t, out.MCMSTimelockProposals, 1)
	require.Nil(t, out.AddressBook)

	capReg := te.OwnedCapabilityRegistry()
	timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
		te.RegistrySelector: {
			Timelock:  capReg.McmsContracts.Timelock,
			CallProxy: capReg.McmsContracts.CallProxy,
		},
	}

	_, err = commonchangeset.Apply(t, te.Env, timelockContracts,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(workflowregistry.UpdateAllowedDons),
			req,
		),
	)
	require.NoError(t, err)
}
