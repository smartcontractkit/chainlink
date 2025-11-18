package memory

import (
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_aptos_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos/provider"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func getTestAptosChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.APTOS_LOCALNET.Selector}
}

func generateChainsAptos(t *testing.T, numChains int) []cldf_chain.BlockChain {
	t.Helper()

	testAptosChainSelectors := getTestAptosChainSelectors()
	if len(testAptosChainSelectors) < numChains {
		t.Fatalf("not enough test aptos chain selectors available")
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testAptosChainSelectors[i]

		c, err := cldf_aptos_provider.NewCTFChainProvider(t, selector,
			cldf_aptos_provider.CTFChainProviderConfig{
				Once:              once,
				DeployerSignerGen: cldf_aptos_provider.AccountGenCTFDefault(),
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
		aptosChain := c.(cldf_aptos.Chain)
		err = migrateAccountToFA(t, aptosChain.DeployerSigner, aptosChain.Client)
		require.NoError(t, err)
	}
	return chains
}

func migrateAccountToFA(t *testing.T, signer aptos.TransactionSigner, client aptos.AptosRpcClient) error {
	// Migrate APT Coin to FA, required for CCIP
	payload := aptos.TransactionPayload{
		Payload: &aptos.EntryFunction{
			Module: aptos.ModuleId{
				Address: aptos.AccountOne,
				Name:    "coin",
			},
			Function: "migrate_to_fungible_store",
			ArgTypes: []aptos.TypeTag{
				{
					Value: &aptos.StructTag{
						Address: aptos.AccountOne,
						Module:  "aptos_coin",
						Name:    "AptosCoin",
					},
				},
			},
			Args: nil,
		},
	}

	// This might fail once this function is removed, remove once the node has been upgraded
	res, err := client.BuildSignAndSubmitTransaction(signer, payload)
	require.NoError(t, err)
	tx, err := client.WaitForTransaction(res.Hash)
	require.NoError(t, err)
	require.Truef(t, tx.Success, "Migrating APT to FungibleAsset failed: %v", tx.VmStatus)
	accountAddress := signer.AccountAddress()
	logger.TestLogger(t).Infof("Migrated account %v to Fungible Asset APT", accountAddress.StringLong())
	return err
}
