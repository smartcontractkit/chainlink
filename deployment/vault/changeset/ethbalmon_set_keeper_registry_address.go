package changeset

import (
	"encoding/json"
	"errors"
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

	seqInput := EthBalMonSetKeeperRegistryAddressSequenceInput{
		Chains:     config.Chains,
		MCMSConfig: config.MCMSConfig,
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
	Chains     map[uint64]vaulttypes.SetKeeperRegistryChainConfig `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig                      `json:"mcms_config,omitempty"`
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

		if len(input.Chains) == 0 {
			return EthBalMonSetKeeperRegistryAddressSequenceOutput{}, errors.New("no chains provided")
		}

		var batches []mcmstypes.BatchOperation
		timelockAddresses := make(map[uint64]string)
		mcmAddressByChain := make(map[uint64]string)

		for chainSelector, chainConfig := range input.Chains {
			opReport, err := operations.ExecuteOperation(
				b,
				SetKeeperRegistryOperation,
				deps,
				SetKeeperRegistryOperationInput{
					ChainSelector:            chainSelector,
					NewKeeperRegistryAddress: chainConfig.NewKeeperRegistryAddress,
					MCMSConfig:               input.MCMSConfig,
				},
			)
			if err != nil {
				return EthBalMonSetKeeperRegistryAddressSequenceOutput{},
					fmt.Errorf("chain %d: failed to generate set keeper registry batch: %w", chainSelector, err)
			}

			opOut := opReport.Output

			batches = append(batches, opOut.BatchOperation)
			timelockAddresses[chainSelector] = opOut.TimelockAddress
			mcmAddressByChain[chainSelector] = opOut.MCMSAddress
		}

		proposal, err := proposeutils.BuildProposalFromBatchesV2(
			deps.Environment,
			timelockAddresses,
			mcmAddressByChain,
			nil,
			batches,
			"EthBalMon SetKeeperRegistryAddress",
			ethBalMonProposalTimelockConfig(input.MCMSConfig),
		)
		if err != nil {
			return EthBalMonSetKeeperRegistryAddressSequenceOutput{},
				fmt.Errorf("failed to build timelock proposal: %w", err)
		}

		b.Logger.Infow("Generated EthBalMon set keeper registry proposal",
			"chains", len(input.Chains),
			"operations", len(batches),
		)

		return EthBalMonSetKeeperRegistryAddressSequenceOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	},
)

type SetKeeperRegistryOperationInput struct {
	ChainSelector            uint64                        `json:"chain_selector"`
	NewKeeperRegistryAddress string                        `json:"new_keeper_registry_address"`
	MCMSConfig               *proposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

type SetKeeperRegistryOperationOutput struct {
	ChainSelector   uint64                   `json:"chain_selector"`
	BatchOperation  mcmstypes.BatchOperation `json:"batch_operation"`
	TimelockAddress string                   `json:"timelock_address"`
	MCMSAddress     string                   `json:"mcms_address"`
}

var SetKeeperRegistryOperation = operations.NewOperation(
	"ethbalmon-set-keeper-registry-op",
	semver.MustParse("1.0.0"),
	"Generate batch operation to set Keeper Registry address on the Ethereum Balance Monitor contract",
	func(
		b operations.Bundle,
		deps VaultDeps,
		input SetKeeperRegistryOperationInput,
	) (SetKeeperRegistryOperationOutput, error) {
		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetKeeperRegistryOperationOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		ethBalMonAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			cldf.ContractType(vaulttypes.EthBalMonContractType),
		)
		if err != nil {
			return SetKeeperRegistryOperationOutput{},
				fmt.Errorf("failed to get EthBalMon address: %w", err)
		}

		timelockAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return SetKeeperRegistryOperationOutput{},
				fmt.Errorf("failed to get timelock address: %w", err)
		}

		mcmsAddr, err := getRequiredContractAddress(
			deps.DataStore,
			input.ChainSelector,
			ethBalMonMCMSContractTypeForProposal(input.MCMSConfig),
		)
		if err != nil {
			return SetKeeperRegistryOperationOutput{},
				fmt.Errorf("failed to get MCMS address: %w", err)
		}

		ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(
			common.HexToAddress(ethBalMonAddr),
			chain.Client,
		)
		if err != nil {
			return SetKeeperRegistryOperationOutput{},
				fmt.Errorf("failed to instantiate EthBalanceMonitor at %s: %w", ethBalMonAddr, err)
		}

		setKeeperRegistryTx, err := ethBalMon.SetKeeperRegistryAddress(
			cldf.SimTransactOpts(),
			common.HexToAddress(input.NewKeeperRegistryAddress),
		)
		if err != nil {
			return SetKeeperRegistryOperationOutput{},
				fmt.Errorf("failed to generate setKeeperRegistryAddress calldata: %w", err)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.EthBalMonContractType,
						Tags: []string{
							"setKeeperRegistryAddress",
						},
					},
					To:               ethBalMonAddr,
					Data:             setKeeperRegistryTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		b.Logger.Infow("Generated EthBalMon set keeper registry batch",
			"chainSelector", input.ChainSelector,
			"ethBalMon", ethBalMonAddr,
			"newKeeperRegistry", input.NewKeeperRegistryAddress,
		)

		return SetKeeperRegistryOperationOutput{
			ChainSelector:   input.ChainSelector,
			BatchOperation:  batch,
			TimelockAddress: timelockAddr,
			MCMSAddress:     mcmsAddr,
		}, nil
	},
)
