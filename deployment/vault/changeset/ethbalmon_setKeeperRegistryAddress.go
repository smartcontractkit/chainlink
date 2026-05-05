package changeset

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/eth_balance_monitor_wrapper"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

type setKeeperRegistryAddress struct{}

var SetKeeperRegistryAddress cldf.ChangeSetV2[vaulttypes.EthBalMonSetKeeperRegistryAddressInput] = setKeeperRegistryAddress{}

func (sk setKeeperRegistryAddress) VerifyPreconditions(env cldf.Environment, config vaulttypes.EthBalMonSetKeeperRegistryAddressInput) error {
	return ValidateSetKeeperRegistryAddressConfig(env.GetContext(), env, config)
}

func (sk setKeeperRegistryAddress) Apply(
	e cldf.Environment,
	config vaulttypes.EthBalMonSetKeeperRegistryAddressInput,
) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Generating SetKeeperRegistryAddress proposal for Ethereum Balance Monitor",
		"numChains", len(config.Chains),
	)

	deps := VaultDeps{
		Environment: e,
		DataStore:   e.DataStore,
	}

	seqInput := EthBalMonSetKeeperRegistryAddressSequenceInput{
		Chains: config.Chains,
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		SetKeeperRegistrySequence,
		deps,
		seqInput,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set keeper registry address sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
	}, nil
}

type EthBalMonSetKeeperRegistryAddressSequenceInput struct {
	Chains map[uint64]vaulttypes.SetKeeperRegistryChainConfig `json:"chains"`
}

type EthBalMonSetKeeperRegistryAddressSequenceOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var SetKeeperRegistrySequence = operations.NewSequence(
	"ethbalmon-set-keeper-registry",
	semver.MustParse("1.0.0"),
	"Generate MCMS timelock proposal to set Keeper Registry address on EthBalMon across chains",
	func(
		b operations.Bundle,
		deps VaultDeps,
		input EthBalMonSetKeeperRegistryAddressSequenceInput,
	) (EthBalMonSetKeeperRegistryAddressSequenceOutput, error) {
		b.Logger.Infow("Starting EthBalMon set keeper registry sequence",
			"chains", len(input.Chains),
		)

		opReport, err := operations.ExecuteOperation(
			b,
			SetKeeperRegistryOperation,
			deps,
			SetKeeperRegistryOperationInput{
				Chains: input.Chains,
			},
		)
		if err != nil {
			return EthBalMonSetKeeperRegistryAddressSequenceOutput{},
				fmt.Errorf("failed to generate set keeper registry proposal: %w", err)
		}

		return EthBalMonSetKeeperRegistryAddressSequenceOutput{
			MCMSTimelockProposals: opReport.Output.MCMSTimelockProposals,
		}, nil
	},
)

type SetKeeperRegistryOperationInput struct {
	Chains map[uint64]vaulttypes.SetKeeperRegistryChainConfig `json:"chains"`
}

type SetKeeperRegistryOperationOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var SetKeeperRegistryOperation = operations.NewOperation(
	"ethbalmon-set-keeper-registry-op",
	semver.MustParse("1.0.0"),
	"Generate proposal to set Keeper Registry address on the Ethereum Balance Monitor contract",
	func(
		b operations.Bundle,
		deps VaultDeps,
		input SetKeeperRegistryOperationInput,
	) (SetKeeperRegistryOperationOutput, error) {
		if len(input.Chains) == 0 {
			return SetKeeperRegistryOperationOutput{}, fmt.Errorf("no chains provided")
		}

		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[mcmstypes.ChainSelector]string)
		chainMetadata := make(map[mcmstypes.ChainSelector]mcmstypes.ChainMetadata)

		evmChains := deps.Environment.BlockChains.EVMChains()

		for chainSelector, chainConfig := range input.Chains {
			chain, ok := evmChains[chainSelector]
			if !ok {
				return SetKeeperRegistryOperationOutput{}, fmt.Errorf("chain not found in environment: %d", chainSelector)
			}

			ethBalMonAddr, err := mustGetContractAddress(
				deps.DataStore,
				chainSelector,
				cldf.ContractType(vaulttypes.ETHBALMON_CONTRACT_TYPE),
			)
			if err != nil {
				return SetKeeperRegistryOperationOutput{},
					fmt.Errorf("chain %d: failed to get EthBalMon address: %w", chainSelector, err)
			}

			timelockAddr, err := mustGetContractAddress(
				deps.DataStore,
				chainSelector,
				commontypes.RBACTimelock,
			)
			if err != nil {
				return SetKeeperRegistryOperationOutput{},
					fmt.Errorf("chain %d: failed to get timelock address: %w", chainSelector, err)
			}

			mcmsAddr, err := mustGetContractAddress(
				deps.DataStore,
				chainSelector,
				commontypes.ManyChainMultisig,
			)
			if err != nil {
				return SetKeeperRegistryOperationOutput{},
					fmt.Errorf("chain %d: failed to get MCMS address: %w", chainSelector, err)
			}

			ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(
				common.HexToAddress(ethBalMonAddr),
				chain.Client,
			)
			if err != nil {
				return SetKeeperRegistryOperationOutput{},
					fmt.Errorf("chain %d: failed to instantiate EthBalanceMonitor at %s: %w", chainSelector, ethBalMonAddr, err)
			}

			setKeeperRegistryTx, err := ethBalMon.SetKeeperRegistryAddress(
				cldf.SimTransactOpts(),
				common.HexToAddress(chainConfig.NewKeeperRegistryAddress),
			)
			if err != nil {
				return SetKeeperRegistryOperationOutput{},
					fmt.Errorf("chain %d: failed to generate setKeeperRegistryAddress calldata: %w", chainSelector, err)
			}

			batches = append(batches, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chainSelector),
				Transactions: []mcmstypes.Transaction{
					{
						OperationMetadata: mcmstypes.OperationMetadata{
							ContractType: vaulttypes.ETHBALMON_CONTRACT_TYPE,
							Tags: []string{
								"setKeeperRegistryAddress",
							},
						},
						To:               ethBalMonAddr,
						Data:             setKeeperRegistryTx.Data(),
						AdditionalFields: json.RawMessage(`{"value": 0}`),
					},
				},
			})

			timelockAddresses[mcmstypes.ChainSelector(chainSelector)] = timelockAddr
			chainMetadata[mcmstypes.ChainSelector(chainSelector)] = mcmstypes.ChainMetadata{
				StartingOpCount: 0,
				MCMAddress:      mcmsAddr,
			}
		}

		proposal, err := mcms.NewTimelockProposalBuilder().
			SetVersion("v1").
			SetAction(mcmstypes.TimelockActionBypass).
			SetTimelockAddresses(timelockAddresses).
			SetChainMetadata(chainMetadata).
			SetOperations(batches).
			SetDescription("Set Keeper Registry address on EthBalanceMonitor across chains").
			Build()
		if err != nil {
			return SetKeeperRegistryOperationOutput{}, fmt.Errorf("failed to build timelock proposal: %w", err)
		}

		b.Logger.Infow("Generated EthBalMon set keeper registry proposal",
			"chains", len(input.Chains),
			"operations", len(batches),
		)

		return SetKeeperRegistryOperationOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)
