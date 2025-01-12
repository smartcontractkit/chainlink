package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSaveExisting(t *testing.T) {
	dummyEnv := deployment.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		Chains: map[uint64]deployment.Chain{
			chainsel.TEST_90000001.Selector: {},
			chainsel.TEST_90000002.Selector: {},
		},
		SolChains: map[uint64]deployment.SolChain{
			chainsel.SOLANA_DEVNET.Selector: {},
		},
	}
	ExistingContracts := ExistingContractsConfig{
		ExistingContracts: []Contract{
			{
				Address: common.BigToAddress(big.NewInt(1)).String(),
				TypeAndVersion: deployment.TypeAndVersion{
					Type:    "dummy1",
					Version: deployment.Version1_5_0,
				},
				ChainSelector: chainsel.TEST_90000001.Selector,
			},
			{
				Address: common.BigToAddress(big.NewInt(2)).String(),
				TypeAndVersion: deployment.TypeAndVersion{
					Type:    "dummy2",
					Version: deployment.Version1_1_0,
				},
				ChainSelector: chainsel.TEST_90000002.Selector,
			},
			{
				Address: solana.NewWallet().PublicKey().String(),
				TypeAndVersion: deployment.TypeAndVersion{
					Type:    "dummy3",
					Version: deployment.Version1_1_0,
				},
				ChainSelector: chainsel.SOLANA_DEVNET.Selector,
			},
		},
	}

	output, err := SaveExistingContracts(dummyEnv, ExistingContracts)
	require.NoError(t, err)
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook))
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.Len(t, addresses, 3)
	addressForChain1, exists := addresses[chainsel.TEST_90000001.Selector]
	require.True(t, exists)
	require.Len(t, addressForChain1, 1)
	addressForSolana, exists := addresses[chainsel.SOLANA_DEVNET.Selector]
	require.True(t, exists)
	require.Len(t, addressForSolana, 1)
}
