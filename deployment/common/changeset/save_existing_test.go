package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
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
		},
	}

	output, err := SaveExistingContractsChangeset(dummyEnv, ExistingContracts)
	require.NoError(t, err)
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook))
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 2)
	addressForChain1, exists := addresses[chainsel.TEST_90000001.Selector]
	require.True(t, exists)
	require.Len(t, addressForChain1, 1)
}

func TestSaveExistingAddressWithLabels(t *testing.T) {
	dummyEnv := deployment.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		Chains: map[uint64]deployment.Chain{
			chainsel.TEST_90000001.Selector: {},
			chainsel.TEST_90000002.Selector: {},
		},
	}
	dummyType1 := deployment.TypeAndVersion{
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
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook))
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 1)
	addressForChain1, exists := addresses[chainsel.TEST_90000001.Selector]
	require.True(t, exists)
	require.Len(t, addressForChain1, 1)
	require.Equal(t, "dummyType 1.5.0 label1 label2", addressForChain1[common.BigToAddress(big.NewInt(1)).String()].String())
}

func TestSaveExistingMCMSAddressWithLabels(t *testing.T) {
	dummyEnv := deployment.Environment{
		Name:              "dummy",
		Logger:            logger.TestLogger(t),
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		Chains: map[uint64]deployment.Chain{
			chainsel.TEST_90000001.Selector: {},
			chainsel.TEST_90000002.Selector: {},
		},
	}
	mcmsContractTV := deployment.TypeAndVersion{
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
	require.NoError(t, dummyEnv.ExistingAddresses.Merge(output.AddressBook))
	addresses, err := dummyEnv.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, addresses, 1)
	addressForChain1, exists := addresses[chainsel.TEST_90000001.Selector]
	require.True(t, exists)
	require.Len(t, addressForChain1, 1)
	// load mcms state
	mcmsState, err := MaybeLoadMCMSWithTimelockChainState(dummyEnv.Chains[chainsel.TEST_90000001.Selector], addressForChain1)
	require.NoError(t, err)
	require.NotNil(t, mcmsState)
	require.NotNil(t, mcmsState.ProposerMcm)
	require.NotNil(t, mcmsState.BypasserMcm)
	require.NotNil(t, mcmsState.CancellerMcm)
}
