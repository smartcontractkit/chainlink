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
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestMCMSWithTimelockState_GenerateMCMSWithTimelockViewV2(t *testing.T) {
	envConfig := memory.MemoryEnvironmentConfig{Chains: 1}
	env := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, envConfig)
	chain := env.BlockChains.EVMChains()[env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]]

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
}

func TestAddressesForChainOnlyAddressBook(t *testing.T) {
	chainSelector := chain_selectors.ETHEREUM_MAINNET.Selector

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
}

func TestAddressesForChainOnlyDataStore(t *testing.T) {
	chainSelector := chain_selectors.ETHEREUM_MAINNET.Selector

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
