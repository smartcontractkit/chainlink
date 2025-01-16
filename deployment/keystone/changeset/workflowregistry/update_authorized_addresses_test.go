package workflowregistry_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
)

func TestUpdateAuthorizedAddresses(t *testing.T) {
	lggr := logger.Test(t)

	chainSel := chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
	resp := workflowregistry.SetupTestWorkflowRegistry(t, lggr, chainSel)
	registry := resp.Registry

	authorizedAddresses, err := registry.GetAllAuthorizedAddresses(&bind.CallOpts{})
	require.NoError(t, err)

	assert.Empty(t, authorizedAddresses)

	env := deployment.Environment{
		Logger: lggr,
		Chains: map[uint64]deployment.Chain{
			chainSel: resp.Chain,
		},
		ExistingAddresses: resp.AddressBook,
	}

	addr := "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	_, err = workflowregistry.UpdateAuthorizedAddresses(
		env,
		&workflowregistry.UpdateAuthorizedAddressesRequest{
			RegistryChainSel: chainSel,
			Addresses:        []string{addr},
			Allowed:          true,
		},
	)
	require.NoError(t, err)

	authorizedAddresses, err = registry.GetAllAuthorizedAddresses(&bind.CallOpts{})
	require.NoError(t, err)

	assert.Len(t, authorizedAddresses, 1)
	assert.Equal(t, authorizedAddresses[0], common.HexToAddress(addr))

	_, err = workflowregistry.UpdateAuthorizedAddresses(
		env,
		&workflowregistry.UpdateAuthorizedAddressesRequest{
			RegistryChainSel: chainSel,
			Addresses:        []string{addr},
			Allowed:          false,
		},
	)
	require.NoError(t, err)

	authorizedAddresses, err = registry.GetAllAuthorizedAddresses(&bind.CallOpts{})
	require.NoError(t, err)

	assert.Empty(t, authorizedAddresses)
}

func Test_UpdateAuthorizedAddresses_WithMCMS(t *testing.T) {
	te := test.SetupTestEnv(t, test.TestConfig{
		WFDonConfig:     test.DonConfig{N: 4},
		AssetDonConfig:  test.DonConfig{N: 4},
		WriterDonConfig: test.DonConfig{N: 4},
		NumChains:       1,
		UseMCMS:         true,
	})

	addr := "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	req := &workflowregistry.UpdateAuthorizedAddressesRequest{
		RegistryChainSel: te.RegistrySelector,
		Addresses:        []string{addr},
		Allowed:          true,
		MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
	}

	out, err := workflowregistry.UpdateAuthorizedAddresses(te.Env, req)
	require.NoError(t, err)
	require.Len(t, out.Proposals, 1)
	require.Nil(t, out.AddressBook)

	contracts := te.ContractSets()[te.RegistrySelector]
	timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
		te.RegistrySelector: {
			Timelock:  contracts.Timelock,
			CallProxy: contracts.CallProxy,
		},
	}

	_, err = commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(workflowregistry.UpdateAuthorizedAddresses),
			Config:    req,
		},
	})
	require.NoError(t, err)
}
