package memory

import (
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_aptos_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos/provider"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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

func createAptosChainConfig(chainID string, chain cldf_aptos.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "localnet"
	chainConfig["NetworkNameFull"] = "aptos-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
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

func fundNodesAptos(t *testing.T, aptosChain cldf_aptos.Chain, nodes []*Node) {
	for _, node := range nodes {
		aptoskeys, err := node.App.GetKeyStore().Aptos().GetAll()
		require.NoError(t, err)
		require.Len(t, aptoskeys, 1)
		transmitter := aptoskeys[0]
		transmitterAccountAddress := aptos.AccountAddress{}
		require.NoError(t, transmitterAccountAddress.ParseStringRelaxed(transmitter.Account()))
		FundAptosAccount(t, aptosChain.DeployerSigner, transmitterAccountAddress, 100*1e8, aptosChain.Client)
	}
}

func FundAptosAccount(t *testing.T, signer aptos.TransactionSigner, to aptos.AccountAddress, amount uint64, client aptos.AptosRpcClient) {
	toBytes, err := bcs.Serialize(&to)
	require.NoError(t, err)
	amountBytes, err := bcs.SerializeU64(amount)
	require.NoError(t, err)
	payload := aptos.TransactionPayload{Payload: &aptos.EntryFunction{
		Module: aptos.ModuleId{
			Address: aptos.AccountOne,
			Name:    "aptos_account",
		},
		Function: "transfer",
		Args: [][]byte{
			toBytes,
			amountBytes,
		},
	}}
	tx, err := client.BuildSignAndSubmitTransaction(signer, payload)
	require.NoError(t, err)
	res, err := client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, res.Success, res.VmStatus)
	sender := signer.AccountAddress()
	t.Logf("Funded account %s from %s with %f APT", to.StringLong(), sender.StringLong(), float64(amount)/1e8)
}
