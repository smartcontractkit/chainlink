package changeset

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/eth_balance_monitor_wrapper"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"
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
		Chains: config.Chains,
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
	Chains map[uint64]vaulttypes.EthBalMonTransferOwnershipChainConfig `json:"chains"`
}

type EthBalMonTransferOwnershipSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

var EthBalMonTransferOwnershipSequence = operations.NewSequence(
	"ethbalmon-transferownership-operation",
	semver.MustParse("1.0.0"),
	"Sequence to create transferOwnership EthBalMon batch transaction",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonTransferOwnershipSeqInput) (EthBalMonTransferOwnershipSeqOutput, error) {
		b.Logger.Infow("Starting EthBalMon transferOwnership sequence",
			"chains", len(input.Chains),
		)
		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)
		inspectorPerChain := make(map[uint64]mcmssdk.Inspector)
		for chainSelector, chainConfig := range input.Chains {
			opReport, err := operations.ExecuteOperation(b, EthBalMonTransferOwnershipOperation, deps, EthBalMonTransferOwnershipOpInput{
				ChainSelector: chainSelector,
				NewOwner:      chainConfig.NewOwner,
			})
			opOutput := opReport.Output
			if err != nil {
				return EthBalMonTransferOwnershipSeqOutput{}, fmt.Errorf("chain %d: failed to generate ownership batch: %w", chainSelector, err)
			}
			batches = append(batches, opOutput.BatchOperation)
			timelockAddresses[chainSelector] = opOutput.TimelockAddress
			mcmAddressByChain[chainSelector] = opOutput.MCMSAddress
			inspectorPerChain[chainSelector] = opOutput.Inspector
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(deps.Environment, timelockAddresses, mcmAddressByChain, inspectorPerChain, batches, "EthBalMon transferOwnership", proposalutils.TimelockConfig{
			MinDelay: 0,
		})

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
	ChainSelector uint64 `json:"chain_selector"`
	NewOwner      string `json:"new_owner"`
}

type EthBalMonTransferOwnershipOpOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
	Inspector       *mcmsevmsdk.Inspector    `json:"inspector"`
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

		ethBalMonAddr, err := mustGetContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.ETHBALMON_CONTRACT_TYPE),
		)
		if err != nil {
			return EthBalMonTransferOwnershipOpOutput{},
				fmt.Errorf("failed to get EthBalMon address: %w", err)
		}

		timelockAddr, err := mustGetContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return EthBalMonTransferOwnershipOpOutput{},
				fmt.Errorf("failed to get timelock address: %w", err)
		}
		mcmsAddr, err := mustGetContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.BypasserManyChainMultisig,
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
						ContractType: vaulttypes.ETHBALMON_CONTRACT_TYPE,
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

		chainInspector := mcmsevmsdk.NewInspector(chain.Client)

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
			Inspector:       chainInspector,
		}, nil
	},
)
