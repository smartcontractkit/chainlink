package workflowregistry_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
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

	env := cldf.Environment{
		Logger:            lggr,
		ExistingAddresses: resp.AddressBook,
		BlockChains: cldf_chain.NewBlockChains(
			map[uint64]cldf_chain.BlockChain{
				chainSel: resp.Chain,
			}),
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
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
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
	require.Len(t, out.MCMSTimelockProposals, 1)
	require.Nil(t, out.AddressBook)

	_, err = commonchangeset.Apply(t, te.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(workflowregistry.UpdateAuthorizedAddresses),
			req,
		),
	)
	require.NoError(t, err)
}
