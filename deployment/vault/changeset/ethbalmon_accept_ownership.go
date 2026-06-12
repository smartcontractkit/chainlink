package changeset

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/eth_balance_monitor_wrapper"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

type ethBalMonAcceptOwnership struct{}

var EthBalMonAcceptOwnership cldf.ChangeSetV2[vaulttypes.EthBalMonAcceptOwnershipInput] = ethBalMonAcceptOwnership{}

func (tw ethBalMonAcceptOwnership) VerifyPreconditions(env cldf.Environment, config vaulttypes.EthBalMonAcceptOwnershipInput) error {
	return ValidateEthBalMonAcceptOwnershipConfig(env.GetContext(), env, config)
}

func (tw ethBalMonAcceptOwnership) Apply(e cldf.Environment, config vaulttypes.EthBalMonAcceptOwnershipInput) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Generating EthBalMon acceptOwnership proposal", "numChains", len(config.Chains))

	evmChains := e.BlockChains.EVMChains()

	var primaryChain cldf_evm.Chain
	for _, chainSelector := range config.Chains {
		primaryChain = evmChains[chainSelector]
		break
	}

	deps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}
	seqInput := EthBalMonAcceptOwnershipSeqInput{
		Chains:     config.Chains,
		MCMSConfig: config.MCMSConfig,
	}
	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, EthBalMonAcceptOwnershipSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed on EthBalMonAcceptOwnershipSequence sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}

type EthBalMonAcceptOwnershipSeqInput struct {
	Chains     []uint64                      `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type EthBalMonAcceptOwnershipSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

var EthBalMonAcceptOwnershipSequence = operations.NewSequence(
	"ethbalmon-acceptownership-sequence",
	semver.MustParse("1.0.0"),
	"Sequence to create acceptOwnership EthBalMon batch transaction",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonAcceptOwnershipSeqInput) (EthBalMonAcceptOwnershipSeqOutput, error) {
		b.Logger.Infow("Starting EthBalMon acceptOwnership sequence",
			"chains", len(input.Chains),
		)
		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)
		for _, chainSelector := range input.Chains {
			opReport, err := operations.ExecuteOperation(b, EthBalMonAcceptOwnershipOperation, deps, EthBalMonAcceptOwnershipOpInput{
				ChainSelector: chainSelector,
				MCMSConfig:    input.MCMSConfig,
			})
			if err != nil {
				return EthBalMonAcceptOwnershipSeqOutput{}, fmt.Errorf("chain %d: failed to generate acceptOwnership batch: %w", chainSelector, err)
			}
			opOutput := opReport.Output

			batches = append(batches, opOutput.BatchOperation)
			timelockAddresses[chainSelector] = opOutput.TimelockAddress
			mcmAddressByChain[chainSelector] = opOutput.MCMSAddress
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment, timelockAddresses, mcmAddressByChain, nil, batches, "EthBalMon acceptOwnership", ethBalMonProposalTimelockConfig(input.MCMSConfig))
		if err != nil {
			return EthBalMonAcceptOwnershipSeqOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}
		b.Logger.Infow("Generated EthBalMon acceptOwnership proposal",
			"chains", len(input.Chains), "operations", len(batches))

		return EthBalMonAcceptOwnershipSeqOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)

type EthBalMonAcceptOwnershipOpInput struct {
	ChainSelector uint64                        `json:"chain_selector"`
	MCMSConfig    *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type EthBalMonAcceptOwnershipOpOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
}

var EthBalMonAcceptOwnershipOperation = operations.NewOperation(
	"ethbalmon-acceptownership-operation",
	semver.MustParse("1.0.0"),
	"Operation to create acceptOwnership EthBalMon batch transaction",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonAcceptOwnershipOpInput) (EthBalMonAcceptOwnershipOpOutput, error) {
		b.Logger.Infow("Starting EthBalMon acceptOwnership operation",
			"chainsel", input.ChainSelector,
		)

		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return EthBalMonAcceptOwnershipOpOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		ethBalMonAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.EthBalMonContractType),
		)
		if err != nil {
			return EthBalMonAcceptOwnershipOpOutput{},
				fmt.Errorf("failed to get EthBalMon address: %w", err)
		}

		timelockAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return EthBalMonAcceptOwnershipOpOutput{},
				fmt.Errorf("failed to get timelock address: %w", err)
		}

		mcmsAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			ethBalMonMCMSContractTypeForProposal(input.MCMSConfig),
		)
		if err != nil {
			return EthBalMonAcceptOwnershipOpOutput{},
				fmt.Errorf("failed to get MCMS address: %w", err)
		}

		ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(common.HexToAddress(ethBalMonAddr), chain.Client)
		if err != nil {
			return EthBalMonAcceptOwnershipOpOutput{},
				fmt.Errorf("failed to instantiate EthBalanceMonitor at %s: %w", ethBalMonAddr, err)
		}

		acceptOwnershipTx, err := ethBalMon.AcceptOwnership(cldf.SimTransactOpts())
		if err != nil {
			return EthBalMonAcceptOwnershipOpOutput{}, fmt.Errorf("failed to generate acceptOwnership calldata on chain %d: %w", input.ChainSelector, err)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.EthBalMonContractType,
						Tags: []string{
							"acceptOwnership",
						},
					},
					To:               ethBalMonAddr,
					Data:             acceptOwnershipTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		b.Logger.Infow("Generated EthBalMon acceptOwnership batch",
			"chainSelector", input.ChainSelector,
			"ethBalMon", ethBalMonAddr,
		)

		return EthBalMonAcceptOwnershipOpOutput{
			ChainSelector:   input.ChainSelector,
			BatchOperation:  batch,
			TimelockAddress: timelockAddr,
			MCMSAddress:     mcmsAddr,
		}, nil
	},
)
