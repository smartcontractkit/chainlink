//go:build integration && db

package opsutils_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestAddEVMCallSequenceToCSOutput_ProposalCombination(t *testing.T) {
	t.Parallel()
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
	)
	env := deployedEnvironment.Env

	// Create initial changeset output with existing proposals to test combination logic
	existingProposal1 := mcmslib.TimelockProposal{
		BaseProposal: mcmslib.BaseProposal{
			Description: "First proposal",
		},
		Operations: []mcmstypes.BatchOperation{
			{
				ChainSelector: mcmstypes.ChainSelector(env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]),
				Transactions: []mcmstypes.Transaction{
					{
						To:               common.HexToAddress("0x1111111111111111111111111111111111111111").String(),
						Data:             []byte("data1"),
						AdditionalFields: json.RawMessage(`{"value": 0}`), // JSON-encoded `{"value": 0}`
					},
				},
			},
		},
	}

	existingProposal2 := mcmslib.TimelockProposal{
		BaseProposal: mcmslib.BaseProposal{
			Description: "Second proposal",
		},
		Operations: []mcmstypes.BatchOperation{
			{
				ChainSelector: mcmstypes.ChainSelector(env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]),
				Transactions: []mcmstypes.Transaction{
					{
						To:               common.HexToAddress("0x1111112222222222222222222222222222222222").String(),
						Data:             []byte("data2"),
						AdditionalFields: json.RawMessage(`{"value": 0}`), // JSON-encoded `{"value": 0}`
					},
				},
			},
		},
	}

	csOutput := cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{
			existingProposal1,
			existingProposal2,
		},
	}

	// Create sequence report with unconfirmed calls to generate a new proposal
	chainSel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]
	seqReport := operations.SequenceReport[string, map[uint64][]opsutils.EVMCallOutput]{
		Report: operations.Report[string, map[uint64][]opsutils.EVMCallOutput]{
			Output: map[uint64][]opsutils.EVMCallOutput{
				chainSel: {
					{
						To:           common.HexToAddress("0x3333333333333333333333333333333333333333"),
						Data:         []byte("new_call_data"),
						ContractType: "TestContract",
						Confirmed:    false, // This will create a new proposal
					},
				},
			},
		},
	}

	mcmsCfg := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second, // No delay for testing
		MCMSAction: mcmstypes.TimelockActionSchedule,
	}

	mcmsDescription := "Third proposal"
	// Load onchain state
	chainState, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)
	t.Logf("mcms state: %+v", chainState.EVMMCMSStateByChain())

	result, err := opsutils.AddEVMCallSequenceToCSOutput(
		env,
		csOutput,
		seqReport,
		nil,
		chainState.EVMMCMSStateByChain(),
		mcmsCfg,
		mcmsDescription,
	)

	require.NoError(t, err)
	assert.Equal(t, seqReport.ExecutionReports, result.Reports)

	// Test the key combination logic:
	// 1. Should have exactly 1 proposal after aggregation
	assert.Len(t, result.MCMSTimelockProposals, 1, "Expected exactly 1 aggregated proposal")

	// 2. Description should be comma-separated combination of all proposals
	aggregatedProposal := result.MCMSTimelockProposals[0]
	expectedDescription := "First proposal, Second proposal, Third proposal"
	assert.Equal(t, expectedDescription, aggregatedProposal.Description,
		"Aggregated proposal should have comma-separated descriptions")

	// 3. Operations should be combined from all proposals
	assert.NotEmpty(t, aggregatedProposal.Operations,
		"Aggregated proposal should contain operations")
}
