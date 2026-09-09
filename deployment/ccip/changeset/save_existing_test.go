package changeset

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSaveExisting(t *testing.T) {
	dummyEnv := cldf.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: cldf_chain.NewBlockChains(
			map[uint64]cldf_chain.BlockChain{
				chainsel.TEST_90000001.Selector: cldf_evm.Chain{},
				chainsel.TEST_90000002.Selector: cldf_evm.Chain{},
			}),
	}
	ExistingContracts := ExistingContractsConfig{
		ExistingContracts: []Contract{
			{
				Address: common.BigToAddress(big.NewInt(1)).String(),
				TypeAndVersion: cldf.TypeAndVersion{
					Type:    "dummy1",
					Version: deployment.Version1_5_0,
				},
				ChainSelector: chainsel.TEST_90000001.Selector,
			},
			{
				Address: common.BigToAddress(big.NewInt(2)).String(),
				TypeAndVersion: cldf.TypeAndVersion{
					Type:    "dummy2",
					Version: deployment.Version1_1_0,
				},
				ChainSelector: chainsel.TEST_90000002.Selector,
			},
		},
	}

	output, err := SaveExistingContractsChangeset(dummyEnv, ExistingContracts)
	require.NoError(t, err)
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook)) //nolint:staticcheck // AddressBook is deprecated but still returned by this changeset.
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 2)
	addressForChain1, exists := addresses[chainsel.TEST_90000001.Selector]
	require.True(t, exists)
	require.Len(t, addressForChain1, 1)
}

func TestSaveExistingAddressWithLabels(t *testing.T) {
	dummyEnv := cldf.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: cldf_chain.NewBlockChains(
			map[uint64]cldf_chain.BlockChain{
				chainsel.TEST_90000001.Selector: cldf_evm.Chain{},
				chainsel.TEST_90000002.Selector: cldf_evm.Chain{},
			}),
	}
	dummyType1 := cldf.TypeAndVersion{
		Type:    "dummyType",
		Version: deployment.Version1_5_0,
	}
	dummyType1.AddLabel("label1")
	dummyType1.AddLabel("label2")
	ExistingContracts := ExistingContractsConfig{
		ExistingContracts: []Contract{
			{
				Address:        common.BigToAddress(big.NewInt(1)).String(),
				TypeAndVersion: dummyType1,
				ChainSelector:  chainsel.TEST_90000001.Selector,
			},
		},
	}

	output, err := SaveExistingContractsChangeset(dummyEnv, ExistingContracts)
	require.NoError(t, err)
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook)) //nolint:staticcheck // AddressBook is deprecated but still returned by this changeset.
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 1)
	addressForChain1, exists := addresses[chainsel.TEST_90000001.Selector]
	require.True(t, exists)
	require.Len(t, addressForChain1, 1)
	require.Equal(t, "dummyType 1.5.0 label1 label2", addressForChain1[common.BigToAddress(big.NewInt(1)).String()].String())
}

func TestSaveExistingMCMSAddressWithLabels(t *testing.T) {
	dummyEnv := cldf.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: cldf_chain.NewBlockChains(
			map[uint64]cldf_chain.BlockChain{
				chainsel.TEST_90000001.Selector: cldf_evm.Chain{},
				chainsel.TEST_90000002.Selector: cldf_evm.Chain{},
			}),
	}
	mcmsContractTV := cldf.TypeAndVersion{
		Type:    types.ManyChainMultisig,
		Version: deployment.Version1_0_0,
	}
	mcmsContractTV.AddLabel(types.ProposerRole.String())
	mcmsContractTV.AddLabel(types.BypasserRole.String())
	mcmsContractTV.AddLabel(types.CancellerRole.String())
	ExistingContracts := ExistingContractsConfig{
		ExistingContracts: []Contract{
			{
				Address:        common.BigToAddress(big.NewInt(1)).String(),
				TypeAndVersion: mcmsContractTV,
				ChainSelector:  chainsel.TEST_90000001.Selector,
			},
		},
	}

	output, err := SaveExistingContractsChangeset(dummyEnv, ExistingContracts)
	require.NoError(t, err)
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook)) //nolint:staticcheck // AddressBook is deprecated but still returned by this changeset.
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 1)
	addressForChain1, exists := addresses[chainsel.TEST_90000001.Selector]
	require.True(t, exists)
	require.Len(t, addressForChain1, 1)
	// load mcms state
	mcmsState, err := evmstate.MaybeLoadMCMSWithTimelockChainState(dummyEnv.BlockChains.EVMChains()[chainsel.TEST_90000001.Selector], addressForChain1)
	require.NoError(t, err)
	require.NotNil(t, mcmsState)
	require.NotNil(t, mcmsState.ProposerMcm)
	require.NotNil(t, mcmsState.BypasserMcm)
	require.NotNil(t, mcmsState.CancellerMcm)
}

func TestSaveExistingContractsDualWrite(t *testing.T) {
	t.Parallel()

	dummyEnv := cldf.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: cldf_chain.NewBlockChains(
			map[uint64]cldf_chain.BlockChain{
				chainsel.TEST_90000001.Selector: cldf_evm.Chain{},
			}),
	}
	addr1 := common.BigToAddress(big.NewInt(1)).String()
	addr2 := common.BigToAddress(big.NewInt(2)).String()

	t.Run("records each contract in both registries under its operator qualifier", func(t *testing.T) {
		t.Parallel()
		cfg := ExistingContractsConfig{
			ExistingContracts: []Contract{
				{
					Address: addr1,
					TypeAndVersion: cldf.TypeAndVersion{
						Type:    "dummy1",
						Version: deployment.Version1_5_0,
					},
					ChainSelector: chainsel.TEST_90000001.Selector,
					Qualifier:     "LINK",
				},
				{
					Address: addr2,
					TypeAndVersion: cldf.TypeAndVersion{
						Type:    "dummy2",
						Version: deployment.Version1_1_0,
					},
					ChainSelector: chainsel.TEST_90000001.Selector,
				},
			},
		}

		output, err := SaveExistingContractsChangeset(dummyEnv, cfg)
		require.NoError(t, err)

		addresses, err := output.AddressBook.Addresses() //nolint:staticcheck // AddressBook is deprecated but still returned by this changeset.
		require.NoError(t, err)
		require.Len(t, addresses[chainsel.TEST_90000001.Selector], 2)

		refs, err := output.DataStore.Addresses().Fetch()
		require.NoError(t, err)
		require.Len(t, refs, 2)
		byAddr := make(map[string]datastore.AddressRef, len(refs))
		for _, ref := range refs {
			byAddr[ref.Address] = ref
		}
		require.Equal(t, "LINK", byAddr[addr1].Qualifier)
		require.Empty(t, byAddr[addr2].Qualifier)
	})

	t.Run("re-recording the same contract is idempotent", func(t *testing.T) {
		t.Parallel()
		ec := Contract{
			Address: addr1,
			TypeAndVersion: cldf.TypeAndVersion{
				Type:    "dummy1",
				Version: deployment.Version1_5_0,
			},
			ChainSelector: chainsel.TEST_90000001.Selector,
			Qualifier:     "LINK",
		}
		cfg := ExistingContractsConfig{ExistingContracts: []Contract{ec, ec}}

		output, err := SaveExistingContractsChangeset(dummyEnv, cfg)
		require.NoError(t, err)
		refs, err := output.DataStore.Addresses().Fetch()
		require.NoError(t, err)
		require.Len(t, refs, 1)
	})

	t.Run("rejects a qualifier derived from the imported address", func(t *testing.T) {
		t.Parallel()
		cfg := ExistingContractsConfig{
			ExistingContracts: []Contract{{
				Address: addr1,
				TypeAndVersion: cldf.TypeAndVersion{
					Type:    "dummy1",
					Version: deployment.Version1_5_0,
				},
				ChainSelector: chainsel.TEST_90000001.Selector,
				Qualifier:     strings.ToLower(addr1) + "-dummy1",
			}},
		}

		_, err := SaveExistingContractsChangeset(dummyEnv, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "contains the address being imported")
	})

	t.Run("rejects a key already holding a different contract", func(t *testing.T) {
		t.Parallel()
		tv := cldf.TypeAndVersion{Type: "dummy1", Version: deployment.Version1_5_0}
		cfg := ExistingContractsConfig{
			ExistingContracts: []Contract{
				{Address: addr1, TypeAndVersion: tv, ChainSelector: chainsel.TEST_90000001.Selector, Qualifier: "LINK"},
				{Address: addr2, TypeAndVersion: tv, ChainSelector: chainsel.TEST_90000001.Selector, Qualifier: "LINK"},
			},
		}

		_, err := SaveExistingContractsChangeset(dummyEnv, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "same datastore key")
		require.Contains(t, err.Error(), "provide distinct semantic qualifiers")
		require.Contains(t, err.Error(), "LINK")
	})
}
