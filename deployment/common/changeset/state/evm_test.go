package state

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestMCMSWithTimelockState_GenerateMCMSWithTimelockViewV2(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	chain := env.BlockChains.EVMChains()[selector]

	proposerMcm := deployMCMEvm(t, chain, &mcmstypes.Config{Quorum: 1, Signers: []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000001"),
	}})
	cancellerMcm := deployMCMEvm(t, chain, &mcmstypes.Config{Quorum: 1, Signers: []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000002"),
	}})
	bypasserMcm := deployMCMEvm(t, chain, &mcmstypes.Config{Quorum: 1, Signers: []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000003"),
	}})
	timelock := deployTimelockEvm(t, chain, big.NewInt(1),
		common.HexToAddress("0x0000000000000000000000000000000000000004"),
		[]common.Address{common.HexToAddress("0x0000000000000000000000000000000000000005")},
		[]common.Address{common.HexToAddress("0x0000000000000000000000000000000000000006")},
		[]common.Address{common.HexToAddress("0x0000000000000000000000000000000000000007")},
		[]common.Address{common.HexToAddress("0x0000000000000000000000000000000000000008")},
	)
	callProxy := deployCallProxyEvm(t, chain,
		common.HexToAddress("0x0000000000000000000000000000000000000009"))

	tests := []struct {
		name      string
		contracts *MCMSWithTimelockState
		want      string
		wantErr   string
	}{
		{
			name: "success",
			contracts: &MCMSWithTimelockState{
				ProposerMcm:  proposerMcm,
				CancellerMcm: cancellerMcm,
				BypasserMcm:  bypasserMcm,
				Timelock:     timelock,
				CallProxy:    callProxy,
			},
			want: fmt.Sprintf(`{
				"proposer": {
					"address": "%s",
					"owner":   "%s",
					"config":  {
						"quorum":       1,
						"signers":      ["0x0000000000000000000000000000000000000001"],
						"groupSigners": []
					}
				},
				"canceller": {
					"address": "%s",
					"owner":   "%s",
					"config":  {
						"quorum":       1,
						"signers":      ["0x0000000000000000000000000000000000000002"],
						"groupSigners": []
					}
				},
				"bypasser": {
					"address": "%s",
					"owner":   "%s",
					"config":  {
						"quorum":       1,
						"signers":      ["0x0000000000000000000000000000000000000003"],
						"groupSigners": []
					}
				},
				"timelock": {
					"address": "%s",
					"owner":   "0x0000000000000000000000000000000000000000",
					"membersByRole": {
						"ADMIN_ROLE":     [ "0x0000000000000000000000000000000000000004" ],
						"PROPOSER_ROLE":  [ "0x0000000000000000000000000000000000000005" ],
						"EXECUTOR_ROLE":  [ "0x0000000000000000000000000000000000000006" ],
						"CANCELLER_ROLE": [ "0x0000000000000000000000000000000000000007" ],
						"BYPASSER_ROLE":  [ "0x0000000000000000000000000000000000000008" ]
					}
				},
				"callProxy": {
					"address": "%s",
					"owner":   "0x0000000000000000000000000000000000000000"
				}
			}`, evmAddr(proposerMcm.Address()), evmAddr(chain.DeployerKey.From),
				evmAddr(cancellerMcm.Address()), evmAddr(chain.DeployerKey.From),
				evmAddr(bypasserMcm.Address()), evmAddr(chain.DeployerKey.From),
				evmAddr(timelock.Address()), evmAddr(callProxy.Address())),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := tt.contracts

			got, err := state.GenerateMCMSWithTimelockView()

			if tt.wantErr == "" {
				require.NoError(t, err)
				require.JSONEq(t, tt.want, toJSON(t, &got))
			} else {
				require.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

func TestAddressesForChain(t *testing.T) {
	chainSelector := chain_selectors.ETHEREUM_MAINNET.Selector

	t.Run("environment with AddressBook only", func(t *testing.T) {
		// Create environment with only AddressBook
		addressBook := cldf.NewMemoryAddressBook()
		err := addressBook.Save(chainSelector, "0x1234567890123456789012345678901234567890",
			cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0))
		require.NoError(t, err)

		env := cldf.Environment{
			ExistingAddresses: addressBook,
			DataStore:         nil, // No DataStore
		}

		// Test the merge function
		mergedAddresses, err := AddressesForChain(env, chainSelector, "")
		require.NoError(t, err)

		// Should have address from AddressBook only
		require.Len(t, mergedAddresses, 1)
		require.Contains(t, mergedAddresses, "0x1234567890123456789012345678901234567890")
	})

	t.Run("environment with DataStore only", func(t *testing.T) {
		// Create environment with only DataStore
		dataStore := datastore.NewMemoryDataStore()
		err := dataStore.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       "0xABCDEF1234567890123456789012345678901234",
			Type:          datastore.ContractType(types.RBACTimelock),
			Version:       &deployment.Version1_0_0,
		})
		require.NoError(t, err)

		addressBook := cldf.NewMemoryAddressBook()

		env := cldf.Environment{
			ExistingAddresses: addressBook,
			DataStore:         dataStore.Seal(),
		}

		// Test the merge function
		mergedAddresses, err := AddressesForChain(env, chainSelector, "")
		require.NoError(t, err)

		// Should have address from DataStore only
		require.Len(t, mergedAddresses, 1)
		require.Contains(t, mergedAddresses, "0xABCDEF1234567890123456789012345678901234")
	})
	t.Run("environment with AddressBook and DataStore without qualifier", func(t *testing.T) {
		// Create a mock environment with both AddressBook and DataStore
		addressBook := cldf.NewMemoryAddressBook()
		err := addressBook.Save(chainSelector, "0x1234567890123456789012345678901234567890",
			cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0))
		require.NoError(t, err)

		dataStore := datastore.NewMemoryDataStore()
		err = dataStore.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       "0xABCDEF1234567890123456789012345678901234",
			Type:          datastore.ContractType(types.RBACTimelock),
			Version:       &deployment.Version1_0_0,
			Labels: datastore.NewLabelSet(
				"team:core",
				"environment:production",
				"role:timelock",
			),
		})
		require.NoError(t, err)

		env := cldf.Environment{
			ExistingAddresses: addressBook,
			DataStore:         dataStore.Seal(),
		}

		// Test the merge function
		mergedAddresses, err := AddressesForChain(env, chainSelector, "")
		require.NoError(t, err)

		// Should have addresses from both sources
		require.Len(t, mergedAddresses, 2)
		require.Contains(t, mergedAddresses, "0x1234567890123456789012345678901234567890")
		require.Contains(t, mergedAddresses, "0xABCDEF1234567890123456789012345678901234")

		// Verify that types are correctly preserved
		linkTokenTV := mergedAddresses["0x1234567890123456789012345678901234567890"]
		require.Equal(t, types.LinkToken, linkTokenTV.Type)
		require.Equal(t, deployment.Version1_0_0, linkTokenTV.Version)

		timelockTV := mergedAddresses["0xABCDEF1234567890123456789012345678901234"]
		require.Equal(t, types.RBACTimelock, timelockTV.Type)
		require.Equal(t, deployment.Version1_0_0, timelockTV.Version)

		// Verify labels are preserved in DataStore
		refs := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
		require.Len(t, refs, 1)

		timelockRef := refs[0]
		require.Equal(t, "0xABCDEF1234567890123456789012345678901234", timelockRef.Address)
		require.True(t, timelockRef.Labels.Contains("team:core"))
		require.True(t, timelockRef.Labels.Contains("environment:production"))
		require.True(t, timelockRef.Labels.Contains("role:timelock"))
	})

	t.Run("environment with AddressBook and DataStore with qualifier", func(t *testing.T) {
		dataStore := datastore.NewMemoryDataStore()

		// Add contracts with different qualifiers
		err := dataStore.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       "0x1111111111111111111111111111111111111111",
			Type:          datastore.ContractType(types.RBACTimelock),
			Version:       &deployment.Version1_0_0,
			Qualifier:     "team-a",
			Labels: datastore.NewLabelSet(
				"team:team-a",
				"role:timelock",
			),
		})
		require.NoError(t, err)

		err = dataStore.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       "0x2222222222222222222222222222222222222222",
			Type:          datastore.ContractType(types.RBACTimelock),
			Version:       &deployment.Version1_0_0,
			Qualifier:     "team-b",
			Labels: datastore.NewLabelSet(
				"team:team-b",
				"role:timelock",
			),
		})
		require.NoError(t, err)

		env := cldf.Environment{
			ExistingAddresses: cldf.NewMemoryAddressBook(),
			DataStore:         dataStore.Seal(),
		}

		// Test filtering by qualifier
		mergedAddresses, err := AddressesForChain(env, chainSelector, "team-a")
		require.NoError(t, err)

		// Should only have team-a contract
		require.Len(t, mergedAddresses, 1)
		require.Contains(t, mergedAddresses, "0x1111111111111111111111111111111111111111")
		require.NotContains(t, mergedAddresses, "0x2222222222222222222222222222222222222222")

		// Verify the correct contract type
		timelockTV := mergedAddresses["0x1111111111111111111111111111111111111111"]
		require.Equal(t, types.RBACTimelock, timelockTV.Type)

		// Verify labels are preserved for the filtered contract
		refs := env.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(chainSelector),
			datastore.AddressRefByQualifier("team-a"),
		)
		require.Len(t, refs, 1)

		teamARef := refs[0]
		require.Equal(t, "0x1111111111111111111111111111111111111111", teamARef.Address)
		require.Equal(t, "team-a", teamARef.Qualifier)
		require.True(t, teamARef.Labels.Contains("team:team-a"))
		require.True(t, teamARef.Labels.Contains("role:timelock"))
	})

	t.Run("environment with duplicated addresses in AddressBook and DataStore", func(t *testing.T) {
		const (
			duplicateAddress = "0x1234567890123456789012345678901234567890"
			uniqueAddress    = "0xABCDEF1234567890123456789012345678901234"
		)

		// Create environment with same address in both AddressBook and DataStore
		addressBook := cldf.NewMemoryAddressBook()
		// Add LinkToken to AddressBook
		err := addressBook.Save(chainSelector, duplicateAddress,
			cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0))
		require.NoError(t, err)

		dataStore := datastore.NewMemoryDataStore()

		// Add the SAME address to DataStore but with different type/version and labels
		err = dataStore.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       duplicateAddress,                           // Same address as AddressBook
			Type:          datastore.ContractType(types.RBACTimelock), // Different type from AddressBook LinkToken
			Version:       &deployment.Version1_6_0,                   // Different version
			Labels: datastore.NewLabelSet(
				"team:datastore-team",
				"environment:staging",
				"override:true",
			),
		})
		require.NoError(t, err)

		// Also add a unique DataStore address
		err = dataStore.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       uniqueAddress,
			Type:          datastore.ContractType(types.RBACTimelock),
			Version:       &deployment.Version1_0_0,
			Labels: datastore.NewLabelSet(
				"team:unique-entry",
				"role:timelock",
			),
		})
		require.NoError(t, err)

		env := cldf.Environment{
			ExistingAddresses: addressBook,
			DataStore:         dataStore.Seal(),
		}

		// Test the merge function
		mergedAddresses, err := AddressesForChain(env, chainSelector, "")
		require.NoError(t, err)

		// Should have 2 addresses total (duplicate should be merged, unique should be included)
		require.Len(t, mergedAddresses, 2)
		require.Contains(t, mergedAddresses, duplicateAddress)
		require.Contains(t, mergedAddresses, uniqueAddress)

		// The duplicate address should use DataStore values (DataStore takes precedence)
		duplicateTV := mergedAddresses[duplicateAddress]
		require.Equal(t, types.RBACTimelock, duplicateTV.Type, "DataStore type should override AddressBook type")
		require.Equal(t, deployment.Version1_6_0, duplicateTV.Version, "DataStore version should override AddressBook version")

		// The unique address should have correct type
		uniqueTV := mergedAddresses[uniqueAddress]
		require.Equal(t, types.RBACTimelock, uniqueTV.Type)
		require.Equal(t, deployment.Version1_0_0, uniqueTV.Version)

		// Verify that DataStore labels are preserved for both addresses
		refs := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
		require.Len(t, refs, 2)

		// Find the refs by address
		var duplicateRef, uniqueRef *datastore.AddressRef
		for i := range refs {
			switch refs[i].Address {
			case duplicateAddress:
				duplicateRef = &refs[i]
			case uniqueAddress:
				uniqueRef = &refs[i]
			}
		}

		require.NotNil(t, duplicateRef, "Should find duplicate address in DataStore")
		require.NotNil(t, uniqueRef, "Should find unique address in DataStore")

		// Verify labels are preserved for the duplicate address (which should come from DataStore)
		require.True(t, duplicateRef.Labels.Contains("team:datastore-team"))
		require.True(t, duplicateRef.Labels.Contains("environment:staging"))
		require.True(t, duplicateRef.Labels.Contains("override:true"))

		// Verify labels for the unique address
		require.True(t, uniqueRef.Labels.Contains("team:unique-entry"))
		require.True(t, uniqueRef.Labels.Contains("role:timelock"))
	})
}

// ----- helpers -----

func toJSON[T any](t *testing.T, value T) string {
	t.Helper()

	bytes, err := json.Marshal(value)
	require.NoError(t, err)

	return string(bytes)
}

func deployMCMEvm(
	t *testing.T, chain cldf_evm.Chain, config *mcmstypes.Config,
) *bindings.ManyChainMultiSig {
	t.Helper()

	_, tx, contract, err := bindings.DeployManyChainMultiSig(chain.DeployerKey, chain.Client)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	groupQuorums, groupParents, signerAddresses, signerGroups, err := mcmsevmsdk.ExtractSetConfigInputs(config)
	require.NoError(t, err)
	tx, err = contract.SetConfig(chain.DeployerKey, signerAddresses, signerGroups, groupQuorums, groupParents, false)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	return contract
}

func deployTimelockEvm(
	t *testing.T, chain cldf_evm.Chain, minDelay *big.Int, admin common.Address,
	proposers, executors, cancellers, bypassers []common.Address,
) *bindings.RBACTimelock {
	t.Helper()
	_, tx, contract, err := bindings.DeployRBACTimelock(
		chain.DeployerKey, chain.Client, minDelay, admin, proposers, executors, cancellers, bypassers)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	return contract
}

func deployCallProxyEvm(
	t *testing.T, chain cldf_evm.Chain, target common.Address,
) *bindings.CallProxy {
	t.Helper()
	_, tx, contract, err := bindings.DeployCallProxy(chain.DeployerKey, chain.Client, target)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	return contract
}

func evmAddr(addr common.Address) string {
	return strings.ToLower(addr.Hex())
}
