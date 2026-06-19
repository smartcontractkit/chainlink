package changeset_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestConfigureOCR3_VerifyPreconditions_ExtraSignerFamilies(t *testing.T) {
	base := changeset.ConfigureOCR3Input{
		ContractChainSelector: 1,
		ContractQualifier:     "q",
		DON: contracts.DonNodeSet{
			Name:    "don",
			NodeIDs: []string{"node-1"},
		},
		OracleConfig: &ocr3.OracleConfig{},
	}

	t.Run("accepts stellar", func(t *testing.T) {
		in := base
		in.ExtraSignerFamilies = []string{chain_selectors.FamilyStellar}
		require.NoError(t, changeset.ConfigureOCR3{}.VerifyPreconditions(cldf.Environment{}, in))
	})

	t.Run("accepts aptos and solana", func(t *testing.T) {
		in := base
		in.ExtraSignerFamilies = []string{chain_selectors.FamilyAptos, chain_selectors.FamilySolana}
		require.NoError(t, changeset.ConfigureOCR3{}.VerifyPreconditions(cldf.Environment{}, in))
	})

	t.Run("rejects unknown family", func(t *testing.T) {
		in := base
		in.ExtraSignerFamilies = []string{"bogus"}
		require.Error(t, changeset.ConfigureOCR3{}.VerifyPreconditions(cldf.Environment{}, in))
	})
}

func TestConfigureOCR3(t *testing.T) {
	h := test.NewTestHarness(t)

	testCases := []struct {
		Name  string
		Input changeset.ConfigureOCR3Input
	}{
		{
			Name: "Consensus Capability OCR3 Deployment",
			Input: changeset.ConfigureOCR3Input{
				ContractChainSelector: h.RegistrySelector,
				ContractQualifier:     "test-ocr3",
				DON: contracts.DonNodeSet{
					Name:    "test-don",                      // This should match the DON created in SetupEnvV2
					NodeIDs: h.Runtime.Environment().NodeIDs, // Use all available node IDs
				},
				OracleConfig: &ocr3.OracleConfig{
					MaxFaultyOracles:     1,
					TransmissionSchedule: []int{len(h.Runtime.Environment().NodeIDs)}, // Single entry with number of nodes
					ConsensusCapOffchainConfig: &ocr3.ConsensusCapOffchainConfig{
						MaxQueryLengthBytes: 1000000,
					},
				},
			},
		},
		{
			Name: "Chain Capability OCR3 Deployment",
			Input: changeset.ConfigureOCR3Input{
				ContractChainSelector: h.RegistrySelector,
				ContractQualifier:     "test-chain-capability-ocr3",
				DON: contracts.DonNodeSet{
					Name:    "test-don",                      // This should match the DON created in SetupEnvV2
					NodeIDs: h.Runtime.Environment().NodeIDs, // Use all available node IDs
				},
				OracleConfig: &ocr3.OracleConfig{
					MaxFaultyOracles:     1,
					TransmissionSchedule: []int{len(h.Runtime.Environment().NodeIDs)}, // Single entry with number of nodes
					ChainCapOffchainConfig: &ocr3.ChainCapOffchainConfig{
						MaxQueryLengthBytes:       1,
						MaxObservationLengthBytes: 2,
						MaxReportLengthBytes:      3,
						MaxOutcomeLengthBytes:     4,
						MaxReportCount:            5,
						MaxBatchSize:              6,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			task := runtime.ChangesetTask(changeset.DeployOCR3{}, changeset.DeployOCR3Input{
				ChainSelector: tc.Input.ContractChainSelector,
				Qualifier:     tc.Input.ContractQualifier,
			})

			err := h.Runtime.Exec(task)
			require.NoError(t, err)

			output := h.Runtime.State().Outputs[task.ID()]
			addresses, err := output.DataStore.Addresses().Fetch()
			require.NoError(t, err, "should fetch addresses without error")
			require.Len(t, addresses, 1, "expected exactly one deployed contract")
			deployedAddress := addresses[0]

			err = h.Runtime.Exec(
				runtime.ChangesetTask(changeset.ConfigureOCR3{}, tc.Input),
			)
			require.NoError(t, err, "ConfigureOCR3 should not return an error")

			// Further verify the deployed contract by connecting to it
			ocr3Contract, err := ocr3_capability.NewOCR3Capability(
				common.HexToAddress(deployedAddress.Address),
				h.Runtime.Environment().BlockChains.EVMChains()[h.RegistrySelector].Client,
			)
			require.NoError(t, err, "failed to create OCR3 contract instance")
			require.NotNil(t, ocr3Contract, "OCR3 contract instance should not be nil")

			// Get ConfigSet events to verify configuration details
			configIterator, err := ocr3Contract.FilterConfigSet(&bind.FilterOpts{})
			require.NoError(t, err, "failed to filter ConfigSet events")

			t.Cleanup(func() {
				configIterator.Close()
			})

			// There should be exactly one ConfigSet event
			require.True(t, configIterator.Next(), "should have at least one ConfigSet event")
			configEvent := configIterator.Event

			// Assert the fault tolerance parameter
			require.Equal(t, uint8(1), configEvent.F, "F should be 1")

			// Assert the number of signers/transmitters matches expected nodes
			require.Len(t, configEvent.Signers, 5, "should have 5 signers")
			require.Len(t, configEvent.Transmitters, 5, "should have 5 transmitters")

			// Assert no more ConfigSet events
			require.False(t, configIterator.Next(), "should have exactly one ConfigSet event")
		})
	}
}
