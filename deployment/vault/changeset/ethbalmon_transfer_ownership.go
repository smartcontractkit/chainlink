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

type ethBalMonTransferOwnership struct{}

var EthBalMonTransferOwnership cldf.ChangeSetV2[vaulttypes.EthBalMonTransferOwnershipInput] = ethBalMonTransferOwnership{}

func (tw ethBalMonTransferOwnership) VerifyPreconditions(env cldf.Environment, config vaulttypes.EthBalMonTransferOwnershipInput) error {
	return ValidateEthBalMonTransferOwnershipConfig(env.GetContext(), env, config)
}

func (tw ethBalMonTransferOwnership) Apply(e cldf.Environment, config vaulttypes.EthBalMonTransferOwnershipInput) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Generating EthBalMon transferOwnership proposal", "numChains", len(config.Chains))

	evmChains := e.BlockChains.EVMChains()

	var primaryChain cldf_evm.Chain
	for chainSelector := range config.Chains {
		primaryChain = evmChains[chainSelector]
		break
	}

	deps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}
	seqInput := EthBalMonTransferOwnershipSeqInput{
		Chains:     config.Chains,
		MCMSConfig: config.MCMSConfig,
	}
	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, EthBalMonTransferOwnershipSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed on EthBalMonTransferOwnershipSequence sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}

type EthBalMonTransferOwnershipSeqInput struct {
	Chains     map[uint64]vaulttypes.EthBalMonTransferOwnershipChainConfig `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig                               `json:"mcms_config,omitempty"`
}

type EthBalMonTransferOwnershipSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

var EthBalMonTransferOwnershipSequence = operations.NewSequence(
	"ethbalmon-transferownership-sequence",
	semver.MustParse("1.0.0"),
	"Sequence to create transferOwnership EthBalMon batch transaction",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonTransferOwnershipSeqInput) (EthBalMonTransferOwnershipSeqOutput, error) {
		b.Logger.Infow("Starting EthBalMon transferOwnership sequence",
			"chains", len(input.Chains),
		)
		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)
		for chainSelector, chainConfig := range input.Chains {
			opReport, err := operations.ExecuteOperation(b, EthBalMonTransferOwnershipOperation, deps, EthBalMonTransferOwnershipOpInput{
				ChainSelector: chainSelector,
				NewOwner:      chainConfig.NewOwner,
				MCMSConfig:    input.MCMSConfig,
			})
			if err != nil {
				return EthBalMonTransferOwnershipSeqOutput{}, fmt.Errorf("chain %d: failed to generate ownership batch: %w", chainSelector, err)
			}
			opOutput := opReport.Output

			batches = append(batches, opOutput.BatchOperation)
			timelockAddresses[chainSelector] = opOutput.TimelockAddress
			mcmAddressByChain[chainSelector] = opOutput.MCMSAddress
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment, timelockAddresses, mcmAddressByChain, nil, batches, "EthBalMon transferOwnership", ethBalMonProposalTimelockConfig(input.MCMSConfig))

		if err != nil {
			return EthBalMonTransferOwnershipSeqOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}
		b.Logger.Infow("Generated EthBalMon transferOwnership proposal",
			"chains", len(input.Chains), "operations", len(batches))

		return EthBalMonTransferOwnershipSeqOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)

type EthBalMonTransferOwnershipOpInput struct {
	ChainSelector uint64                        `json:"chain_selector"`
	NewOwner      string                        `json:"new_owner"`
	MCMSConfig    *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type EthBalMonTransferOwnershipOpOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
}

var EthBalMonTransferOwnershipOperation = operations.NewOperation(
	"ethbalmon-transferownership-operation",
	semver.MustParse("1.0.0"),
	"Operation to create transferOwnership EthBalMon batch transaction",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonTransferOwnershipOpInput) (EthBalMonTransferOwnershipOpOutput, error) {
		b.Logger.Infow("Starting EthBalMon transferOwnership operation",
			"chainsel", input.ChainSelector,
		)

		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]

		if !ok {
			return EthBalMonTransferOwnershipOpOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		ethBalMonAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.EthBalMonContractType),
		)
		if err != nil {
			return EthBalMonTransferOwnershipOpOutput{},
				fmt.Errorf("failed to get EthBalMon address: %w", err)
		}

		timelockAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return EthBalMonTransferOwnershipOpOutput{},
				fmt.Errorf("failed to get timelock address: %w", err)
		}
		mcmsAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			ethBalMonMCMSContractTypeForProposal(input.MCMSConfig),
		)
		if err != nil {
			return EthBalMonTransferOwnershipOpOutput{},
				fmt.Errorf("failed to get MCMS address: %w", err)
		}

		ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(common.HexToAddress(ethBalMonAddr), chain.Client)
		if err != nil {
			return EthBalMonTransferOwnershipOpOutput{},
				fmt.Errorf("failed to instantiate EthBalanceMonitor at %s: %w", ethBalMonAddr, err)
		}

		transferOwnershipTx, err := ethBalMon.TransferOwnership(cldf.SimTransactOpts(), common.HexToAddress(input.NewOwner))
		if err != nil {
			return EthBalMonTransferOwnershipOpOutput{}, fmt.Errorf("failed to generate transferOwnership calldata on chain %d: %w ", input.ChainSelector, err)
		}
		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.EthBalMonContractType,
						Tags: []string{
							"transferOwnership",
						},
					},
					To:               ethBalMonAddr,
					Data:             transferOwnershipTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		b.Logger.Infow("Generated EthBalMon transferOwnership batch",
			"chainSelector", input.ChainSelector,
			"ethBalMon", ethBalMonAddr,
			"newOwner", input.NewOwner,
		)

		return EthBalMonTransferOwnershipOpOutput{
			ChainSelector:   input.ChainSelector,
			BatchOperation:  batch,
			TimelockAddress: timelockAddr,
			MCMSAddress:     mcmsAddr,
		}, nil
	},
)
