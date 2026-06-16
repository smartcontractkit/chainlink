package changeset

import (
	"encoding/json"
	"fmt"
	"math/big"

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

type ethBalMonSetWatchList struct{}

var EthBalMonSetWatchList cldf.ChangeSetV2[vaulttypes.EthBalMonSetWatchListInput] = ethBalMonSetWatchList{}

func (sw ethBalMonSetWatchList) VerifyPreconditions(env cldf.Environment, config vaulttypes.EthBalMonSetWatchListInput) error {
	return ValidateEthBalMonSetWatchListConfig(env.GetContext(), env, config)
}

func (sw ethBalMonSetWatchList) Apply(e cldf.Environment, config vaulttypes.EthBalMonSetWatchListInput) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Generating EthBalMon setWatchList proposal", "numChains", len(config.Chains))

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

	seqInput := EthBalMonSetWatchListSeqInput{
		Chains:     config.Chains,
		MCMSConfig: config.MCMSConfig,
	}

	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, EthBalMonSetWatchListSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed on EthBalMonSetWatchListSequence sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}

type EthBalMonSetWatchListSeqInput struct {
	Chains     map[uint64]vaulttypes.EthBalMonSetWatchListChainConfig `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig                          `json:"mcms_config,omitempty"`
}

type EthBalMonSetWatchListSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

var EthBalMonSetWatchListSequence = operations.NewSequence(
	"ethbalmon-setWatchList-sequence",
	semver.MustParse("1.0.0"),
	"Sequence to create operations for EthBalMon setWatchList",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonSetWatchListSeqInput) (EthBalMonSetWatchListSeqOutput, error) {
		b.Logger.Infow("Starting EthBalMon setWatchList sequence",
			"chains", len(input.Chains),
		)
		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)

		for chainSelector, chainConfig := range input.Chains {
			opReport, err := operations.ExecuteOperation(b, EthBalMonSetWatchListOperation, deps, EthBalMonSetWatchListOpInput{
				ChainSelector:   chainSelector,
				Addresses:       chainConfig.Addresses,
				MinBalancesWei:  chainConfig.MinBalancesWei,
				TopUpAmountsWei: chainConfig.TopUpAmountsWei,
				MCMSConfig:      input.MCMSConfig,
			})
			if err != nil {
				return EthBalMonSetWatchListSeqOutput{}, fmt.Errorf("chain %d: failed to generate setWatchList batch: %w", chainSelector, err)
			}
			opOutput := opReport.Output

			batches = append(batches, opOutput.BatchOperation)
			timelockAddresses[chainSelector] = opOutput.TimelockAddress
			mcmAddressByChain[chainSelector] = opOutput.MCMSAddress
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment, timelockAddresses, mcmAddressByChain, nil, batches, "EthBalMon SetWatchList", ethBalMonProposalTimelockConfig(input.MCMSConfig))

		if err != nil {
			return EthBalMonSetWatchListSeqOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}
		b.Logger.Infow("Generated EthBalMon setWatchList proposal",
			"chains", len(input.Chains), "operations", len(batches))

		return EthBalMonSetWatchListSeqOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)

type EthBalMonSetWatchListOpInput struct {
	ChainSelector   uint64                        `json:"chain_selector"`
	Addresses       []common.Address              `json:"addresses"`
	MinBalancesWei  []*big.Int                    `json:"min_balance_wei"`
	TopUpAmountsWei []*big.Int                    `json:"topup_amounts_wei"`
	MCMSConfig      *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type EthBalMonSetWatchListOpOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
}

var EthBalMonSetWatchListOperation = operations.NewOperation(
	"ethbalmon-setWatchList-operation",
	semver.MustParse("1.0.0"),
	"Operation to create transaction batch for EthBalMon setWatchList",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonSetWatchListOpInput) (EthBalMonSetWatchListOpOutput, error) {
		b.Logger.Infow("Starting EthBalMon setWatchList operation",
			"chainsel", input.ChainSelector,
			"addresses", len(input.Addresses),
		)
		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]

		if !ok {
			return EthBalMonSetWatchListOpOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		ethBalMonAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.EthBalMonContractType),
		)
		if err != nil {
			return EthBalMonSetWatchListOpOutput{},
				fmt.Errorf("failed to get EthBalMon address: %w", err)
		}

		timelockAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return EthBalMonSetWatchListOpOutput{},
				fmt.Errorf("failed to get timelock address: %w", err)
		}
		mcmsAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			ethBalMonMCMSContractTypeForProposal(input.MCMSConfig),
		)
		if err != nil {
			return EthBalMonSetWatchListOpOutput{},
				fmt.Errorf("failed to get MCMS address: %w", err)
		}

		ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(
			common.HexToAddress(ethBalMonAddr),
			chain.Client,
		)
		if err != nil {
			return EthBalMonSetWatchListOpOutput{},
				fmt.Errorf("failed to instantiate EthBalanceMonitor at %s: %w", ethBalMonAddr, err)
		}

		setWatchListTx, err := ethBalMon.SetWatchList(cldf.SimTransactOpts(), input.Addresses, input.MinBalancesWei, input.TopUpAmountsWei)
		if err != nil {
			return EthBalMonSetWatchListOpOutput{}, fmt.Errorf("failed to generate setWatchList calldata on chain %d: %w ", input.ChainSelector, err)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.EthBalMonContractType,
						Tags: []string{
							"setWatchList",
						},
					},
					To:               ethBalMonAddr,
					Data:             setWatchListTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		b.Logger.Infow("Generated EthBalMon setWatchlist batch",
			"chainSelector", input.ChainSelector,
			"ethBalMon", ethBalMonAddr,
			"newWatchList", input.Addresses,
		)

		return EthBalMonSetWatchListOpOutput{
			ChainSelector:   input.ChainSelector,
			BatchOperation:  batch,
			TimelockAddress: timelockAddr,
			MCMSAddress:     mcmsAddr,
		}, nil
	},
)
