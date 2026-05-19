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

type ethBalMonWithdraw struct{}

var EthBalMonWithdraw cldf.ChangeSetV2[vaulttypes.EthBalMonWithdrawInput] = ethBalMonWithdraw{}

func (w ethBalMonWithdraw) VerifyPreconditions(env cldf.Environment, config vaulttypes.EthBalMonWithdrawInput) error {
	return ValidateEthBalMonWithdrawConfig(env.GetContext(), env, config)
}

func (w ethBalMonWithdraw) Apply(e cldf.Environment, config vaulttypes.EthBalMonWithdrawInput) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Generating EthBalMon withdraw proposal", "numChains", len(config.Chains))

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
	seqInput := EthBalMonWithdrawSeqInput{
		Chains:     config.Chains,
		MCMSConfig: config.MCMSConfig,
	}
	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, EthBalMonWithdrawSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed on EthBalMonWithdrawSequence sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}

type EthBalMonWithdrawSeqInput struct {
	Chains     map[uint64]vaulttypes.EthBalMonWithdrawChainConfig `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig                      `json:"mcms_config,omitempty"`
}

type EthBalMonWithdrawSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal `json:"mcms_timelock_proposals"`
}

var EthBalMonWithdrawSequence = operations.NewSequence(
	"ethbalmon-withdraw-sequence",
	semver.MustParse("1.0.0"),
	"Sequence to create operation for EthBalMon withdraw",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonWithdrawSeqInput) (EthBalMonWithdrawSeqOutput, error) {
		b.Logger.Infow("Starting EthBalMon withdraw sequence",
			"numChains", len(input.Chains),
		)
		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)
		for chainSelector, chainConfig := range input.Chains {
			opReport, err := operations.ExecuteOperation(b, EthBalMonWithdrawOperation, deps, EthBalMonWithdrawOpInput{
				ChainSelector: chainSelector,
				Amount:        chainConfig.Amount,
				Payee:         chainConfig.Payee,
				MCMSConfig:    input.MCMSConfig,
			})
			if err != nil {
				return EthBalMonWithdrawSeqOutput{}, fmt.Errorf("chain %d: failed to generate withdraw batch: %w", chainSelector, err)
			}
			opOutput := opReport.Output

			batches = append(batches, opOutput.BatchOperation)
			timelockAddresses[chainSelector] = opOutput.TimelockAddress
			mcmAddressByChain[chainSelector] = opOutput.MCMSAddress
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(deps.Environment, timelockAddresses, mcmAddressByChain, nil, batches, "EthBalMon Withdraw", ethBalMonProposalTimelockConfig(input.MCMSConfig))

		if err != nil {
			return EthBalMonWithdrawSeqOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}
		b.Logger.Infow("Generated EthBalMon withdraw proposal",
			"chains", len(input.Chains), "operations", len(batches))

		return EthBalMonWithdrawSeqOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)

type EthBalMonWithdrawOpInput struct {
	ChainSelector uint64                        `json:"chain_selector"`
	Amount        *big.Int                      `json:"amount"`
	Payee         string                        `json:"payee"`
	MCMSConfig    *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type EthBalMonWithdrawOpOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
}

var EthBalMonWithdrawOperation = operations.NewOperation(
	"ethbalmon-withdraw-operation",
	semver.MustParse("1.0.0"),
	"Operation to create withdraw EthBalMon batch transaction",
	func(b operations.Bundle, deps VaultDeps, input EthBalMonWithdrawOpInput) (EthBalMonWithdrawOpOutput, error) {
		b.Logger.Infow("Starting EthBalMon withdraw operation",
			"chainsel", input.ChainSelector,
		)

		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]

		if !ok {
			return EthBalMonWithdrawOpOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		ethBalMonAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.EthBalMonContractType),
		)
		if err != nil {
			return EthBalMonWithdrawOpOutput{},
				fmt.Errorf("failed to get EthBalMon address: %w", err)
		}

		timelockAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return EthBalMonWithdrawOpOutput{},
				fmt.Errorf("failed to get timelock address: %w", err)
		}
		mcmsAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			ethBalMonMCMSContractTypeForProposal(input.MCMSConfig),
		)
		if err != nil {
			return EthBalMonWithdrawOpOutput{},
				fmt.Errorf("failed to get MCMS address: %w", err)
		}

		ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(common.HexToAddress(ethBalMonAddr), chain.Client)
		if err != nil {
			return EthBalMonWithdrawOpOutput{},
				fmt.Errorf("failed to instantiate EthBalanceMonitor at %s: %w", ethBalMonAddr, err)
		}

		withdrawTx, err := ethBalMon.Withdraw(cldf.SimTransactOpts(), input.Amount, common.HexToAddress(input.Payee))
		if err != nil {
			return EthBalMonWithdrawOpOutput{}, fmt.Errorf("failed to generate withdraw calldata on chain %d: %w", input.ChainSelector, err)
		}
		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.EthBalMonContractType,
						Tags: []string{
							"withdraw",
						},
					},
					To:               ethBalMonAddr,
					Data:             withdrawTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		b.Logger.Infow("Generated EthBalMon withdraw batch",
			"chainSelector", input.ChainSelector,
			"ethBalMon", ethBalMonAddr,
			"amount", input.Amount,
			"payee", input.Payee,
		)

		return EthBalMonWithdrawOpOutput{
			ChainSelector:   input.ChainSelector,
			BatchOperation:  batch,
			TimelockAddress: timelockAddr,
			MCMSAddress:     mcmsAddr,
		}, nil
	},
)
