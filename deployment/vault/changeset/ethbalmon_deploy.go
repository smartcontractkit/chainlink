package changeset

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/eth_balance_monitor_wrapper"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	vaulttypes "github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

const defaultEthBalMonMinWaitPeriodSeconds uint64 = 60

var DeployEthBalMonChangeSet cldf.ChangeSetV2[vaulttypes.DeployEthBalMonInput] = deployEthBalMon{}

type deployEthBalMon struct {
}

func effectiveMinWaitPeriodSeconds(v uint64) uint64 {
	if v == 0 {
		return defaultEthBalMonMinWaitPeriodSeconds
	}
	return v
}

func getRequiredContractAddress(store ds.DataStore, chainSelector uint64, contractType cldf.ContractType) (string, error) {
	addr, err := GetContractAddress(store, chainSelector, contractType)
	if err != nil {
		return "", fmt.Errorf("failed to get contract address for type %s on chain %d: %w", contractType, chainSelector, err)
	}
	if addr == "" {
		return "", fmt.Errorf("empty contract address for type %s on chain %d", contractType, chainSelector)
	}
	return addr, nil
}

func (d deployEthBalMon) VerifyPreconditions(env cldf.Environment, config vaulttypes.DeployEthBalMonInput) error {
	return ValidateDeployEthBalMonConfig(env.GetContext(), env, config)
}

func (d deployEthBalMon) Apply(e cldf.Environment, config vaulttypes.DeployEthBalMonInput) (cldf.ChangesetOutput, error) {
	logger := e.Logger
	logger.Infow("Deploying Ethereum Balances Monitor",
		"numChains", len(config.Chains))

	evmChains := e.BlockChains.EVMChains()

	// Pick a deterministic primary chain for deps (only needed for direct execution, not MCMS).
	var (
		primaryChainSelector uint64
		primaryChainSet      bool
	)
	for chainSelector := range config.Chains {
		if !primaryChainSet || chainSelector < primaryChainSelector {
			primaryChainSelector = chainSelector
			primaryChainSet = true
		}
	}

	primaryChain := evmChains[primaryChainSelector]

	deps := VaultDeps{
		Auth:        primaryChain.DeployerKey,
		Chain:       primaryChain,
		Environment: e,
		DataStore:   e.DataStore,
	}
	seqInput := DeployEthBalMonSequenceInput{
		Chains:     config.Chains,
		MCMSConfig: config.MCMSConfig,
	}

	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, DeployEthBalMonSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ethereum balance monitor contract sequence: %w", err)
	}

	logger.Infow("ethbalmon contract deployed successfully",
		"chains", len(config.Chains))

	seqOut := seqReport.Output
	memoryDataStore := ds.NewMemoryDataStore()
	contractsByChain := make(map[uint64]string)

	for _, chainOut := range seqOut.Chains {
		contractsByChain[chainOut.ChainSelector] = chainOut.ContractAddress

		addressRef := ds.AddressRef{
			ChainSelector: chainOut.ChainSelector,
			Address:       chainOut.ContractAddress,
			Type:          ds.ContractType(vaulttypes.EthBalMonContractType),
			Version:       semver.MustParse("1.0.0"),
			Qualifier:     fmt.Sprintf("%s:%s", vaulttypes.EthBalMonContractType, chainOut.ContractAddress),
			Labels: ds.NewLabelSet(
				vaulttypes.EthBalMonContractType,
				"EthBalMonV1_0_0",
			),
		}

		contractMetadata := ds.ContractMetadata{
			ChainSelector: chainOut.ChainSelector,
			Address:       chainOut.ContractAddress,
			Metadata: map[string]any{
				"deployTxHash":            chainOut.DeployTxHash,
				"deployBlockNumber":       chainOut.DeployBlockNumber,
				"keeperRegistryAddress":   chainOut.KeeperRegistryAddress,
				"minWaitPeriodSeconds":    chainOut.MinWaitPeriodSeconds,
				"timelockAddress":         chainOut.TimelockAddress,
				"mcmsAddress":             chainOut.MCMSAddress,
				"transferOwnershipTxHash": chainOut.TransferOwnershipTxHash,
			},
		}

		if err := memoryDataStore.Addresses().Add(addressRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add address ref for chain %d: %w", chainOut.ChainSelector, err)
		}
		if err := memoryDataStore.ContractMetadata().Add(contractMetadata); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add contract metadata for chain %d: %w", chainOut.ChainSelector, err)
		}
	}

	proposal, err := BuildAcceptOwnershipTimelockProposal(
		e,
		AcceptOwnershipProposalInput{
			ContractsByChain: contractsByChain,
			Description:      "Accept ownership of EthBalanceMonitor across chains",
			MCMSConfig:       deployEthBalMonAcceptOwnershipTimelockConfig(config.MCMSConfig),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build accept ownership proposal: %w", err)
	}

	logger.Infow("Ethereum Balance Monitor deployment completed successfully",
		"chains", len(seqOut.Chains),
	)
	return cldf.ChangesetOutput{
		DataStore:             memoryDataStore,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}

// ================================================
// ================================================
// Deploy Ethereum Balance Monitor SEQUENCE
// ================================================
// ================================================

type DeployEthBalMonSequenceInput struct {
	Chains     map[uint64]vaulttypes.DeployEthBalMonChainConfig `json:"chains"`
	MCMSConfig *proposalutils.TimelockConfig                    `json:"mcms_config,omitempty"`
}
type DeployEthBalMonPerChainOutput struct {
	ChainSelector           uint64
	ContractAddress         string
	DeployTxHash            string
	DeployBlockNumber       uint64
	KeeperRegistryAddress   string
	MinWaitPeriodSeconds    uint64
	TimelockAddress         string
	MCMSAddress             string
	TransferOwnershipTxHash string
}

type DeployEthBalMonSequenceOutput struct {
	Chains []DeployEthBalMonPerChainOutput
}

var DeployEthBalMonSequence = operations.NewSequence(
	"deploy-ethbalmon-sequence",
	semver.MustParse("1.0.0"),
	"Deploy ethereum balance monitor contracts and transfer ownership",
	func(b operations.Bundle, deps VaultDeps, input DeployEthBalMonSequenceInput) (DeployEthBalMonSequenceOutput, error) {
		b.Logger.Infow("Starting deploy ethbalmon contract sequence",
			"chains", len(input.Chains),
		)
		out := DeployEthBalMonSequenceOutput{
			Chains: []DeployEthBalMonPerChainOutput{},
		}

		evmChains := deps.Environment.BlockChains.EVMChains()
		for chainSelector, chainConfig := range input.Chains {
			_, ok := evmChains[chainSelector]
			if !ok {
				return DeployEthBalMonSequenceOutput{}, fmt.Errorf("chain not found in environment: %d", chainSelector)
			}
			var rawMinWait uint64
			if chainConfig.SetMinWaitPeriodSeconds != nil {
				rawMinWait = *chainConfig.SetMinWaitPeriodSeconds
			}
			minWait := effectiveMinWaitPeriodSeconds(rawMinWait)
			timelockAddr, err := getRequiredContractAddress(deps.DataStore, chainSelector, commontypes.RBACTimelock)
			if err != nil {
				return DeployEthBalMonSequenceOutput{}, fmt.Errorf("chain %d: failed to get timelock address: %w", chainSelector, err)
			}
			mcmsAddr, err := getRequiredContractAddress(
				deps.DataStore,
				chainSelector,
				ethBalMonMCMSContractTypeForAction(deployEthBalMonAcceptOwnershipMCMSAction(input.MCMSConfig)),
			)
			if err != nil {
				return DeployEthBalMonSequenceOutput{}, fmt.Errorf("chain %d: failed to get mcms address: %w", chainSelector, err)
			}
			deployReport, err := operations.ExecuteOperation(
				b,
				DeployEthBalMonContractOperation,
				deps,
				DeployEthBalMonContractInput{
					ChainSelector:         chainSelector,
					KeeperRegistryAddress: chainConfig.SetKeeperRegistryAddress,
					MinWaitPeriodSeconds:  minWait,
				},
			)
			if err != nil {
				return DeployEthBalMonSequenceOutput{}, fmt.Errorf("chain %d: deploy operation failed: %w", chainSelector, err)
			}
			deployOut := deployReport.Output
			transferReport, err := operations.ExecuteOperation(
				b,
				TransferOwnershipOperation,
				deps,
				TransferEthBalMonOwnershipInput{
					ChainSelector:   chainSelector,
					ContractAddress: deployOut.ContractAddress,
					TimelockAddress: timelockAddr,
				},
			)
			if err != nil {
				return DeployEthBalMonSequenceOutput{}, fmt.Errorf("chain %d: transfer ownership operation failed: %w", chainSelector, err)
			}
			transferOut := transferReport.Output

			out.Chains = append(out.Chains, DeployEthBalMonPerChainOutput{
				ChainSelector:           chainSelector,
				ContractAddress:         deployOut.ContractAddress,
				DeployTxHash:            deployOut.TxHash,
				DeployBlockNumber:       deployOut.BlockNumber,
				KeeperRegistryAddress:   deployOut.KeeperRegistryAddress,
				MinWaitPeriodSeconds:    deployOut.MinWaitPeriodSeconds,
				TimelockAddress:         timelockAddr,
				MCMSAddress:             mcmsAddr,
				TransferOwnershipTxHash: transferOut.TxHash,
			})
		}
		return out, nil
	},
)

// ================================================
// ================================================
// Deploy Ethereum Balance Monitor OPERATION
// ================================================
// ================================================

type DeployEthBalMonContractInput struct {
	ChainSelector         uint64 `json:"chain_selector"`
	KeeperRegistryAddress string `json:"keeper_registry_address"`
	MinWaitPeriodSeconds  uint64 `json:"min_wait_period_seconds"`
}

type DeployEthBalMonContractOutput struct {
	ChainSelector         uint64 `json:"chain_selector"`
	ContractAddress       string `json:"contract_address"`
	TxHash                string `json:"tx_hash"`
	BlockNumber           uint64 `json:"block_number"`
	KeeperRegistryAddress string `json:"keeper_registry_address"`
	MinWaitPeriodSeconds  uint64 `json:"min_wait_period_seconds"`
}

var DeployEthBalMonContractOperation = operations.NewOperation(
	"deploy-ethbalmon-contract",
	semver.MustParse("1.0.0"),
	"Deploy the Ethereum Balance Monitor contract",
	func(b operations.Bundle, deps VaultDeps, input DeployEthBalMonContractInput) (DeployEthBalMonContractOutput, error) {
		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployEthBalMonContractOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		keeperRegistryAddress := common.HexToAddress(input.KeeperRegistryAddress)

		b.Logger.Infow("Deploying EthBalanceMonitor",
			"chainSelector", input.ChainSelector,
			"keeperRegistryAddress", keeperRegistryAddress.Hex(),
			"minWaitPeriodSeconds", input.MinWaitPeriodSeconds,
		)

		ethBalMonAddr, tx, _, err := eth_balance_monitor_wrapper.DeployEthBalanceMonitor(
			chain.DeployerKey,
			chain.Client,
			keeperRegistryAddress,
			new(big.Int).SetUint64(input.MinWaitPeriodSeconds),
		)
		if err != nil {
			return DeployEthBalMonContractOutput{}, fmt.Errorf("failed to deploy EthBalanceMonitor: %w", err)
		}

		blockNumber, err := chain.Confirm(tx)
		if err != nil {
			return DeployEthBalMonContractOutput{}, fmt.Errorf("failed to confirm deploy tx %s: %w", tx.Hash().Hex(), err)
		}

		out := DeployEthBalMonContractOutput{
			ChainSelector:         input.ChainSelector,
			ContractAddress:       ethBalMonAddr.Hex(),
			TxHash:                tx.Hash().Hex(),
			BlockNumber:           blockNumber,
			KeeperRegistryAddress: keeperRegistryAddress.Hex(),
			MinWaitPeriodSeconds:  input.MinWaitPeriodSeconds,
		}

		b.Logger.Infow("EthBalanceMonitor deployed successfully",
			"chainSelector", input.ChainSelector,
			"contractAddress", out.ContractAddress,
			"txHash", out.TxHash,
			"blockNumber", out.BlockNumber,
		)

		return out, nil

	},
)

// ================================================
// ================================================
// Transfer ownership out of KMS OPERATION
// ================================================
// ================================================

type TransferEthBalMonOwnershipInput struct {
	ChainSelector   uint64 `json:"chain_selector"`
	ContractAddress string `json:"contract_address"`
	TimelockAddress string `json:"timelock_address"`
}

type TransferEthBalMonOwnershipOutput struct {
	ChainSelector   uint64 `json:"chain_selector"`
	ContractAddress string `json:"contract_address"`
	TimelockAddress string `json:"timelock_address"`
	TxHash          string `json:"tx_hash"`
}

var TransferOwnershipOperation = operations.NewOperation(
	"transfer-ownership",
	semver.MustParse("1.0.0"),
	"Transfer contract ownership out of KMS to Timelock",
	func(b operations.Bundle, deps VaultDeps, input TransferEthBalMonOwnershipInput) (TransferEthBalMonOwnershipOutput, error) {
		chain, ok := deps.Environment.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return TransferEthBalMonOwnershipOutput{}, fmt.Errorf("chain not found in environment: %d", input.ChainSelector)
		}

		ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(
			common.HexToAddress(input.ContractAddress),
			chain.Client,
		)
		if err != nil {
			return TransferEthBalMonOwnershipOutput{}, fmt.Errorf("failed to instantiate EthBalanceMonitor at %s: %w", input.ContractAddress, err)
		}

		b.Logger.Infow("Transferring EthBalanceMonitor ownership",
			"chainSelector", input.ChainSelector,
			"contractAddress", input.ContractAddress,
			"timelockAddress", input.TimelockAddress,
		)

		tx, err := ethBalMon.TransferOwnership(
			chain.DeployerKey,
			common.HexToAddress(input.TimelockAddress),
		)
		if err != nil {
			return TransferEthBalMonOwnershipOutput{}, fmt.Errorf("failed to transfer ownership: %w", err)
		}

		if _, err := chain.Confirm(tx); err != nil {
			return TransferEthBalMonOwnershipOutput{}, fmt.Errorf("failed to confirm transfer ownership tx %s: %w", tx.Hash().Hex(), err)
		}

		out := TransferEthBalMonOwnershipOutput{
			ChainSelector:   input.ChainSelector,
			ContractAddress: input.ContractAddress,
			TimelockAddress: input.TimelockAddress,
			TxHash:          tx.Hash().Hex(),
		}

		b.Logger.Infow("EthBalanceMonitor ownership transferred successfully",
			"chainSelector", input.ChainSelector,
			"contractAddress", input.ContractAddress,
			"timelockAddress", input.TimelockAddress,
			"txHash", out.TxHash,
		)

		return out, nil

	},
)

// ======================================================
// ======================================================
// Operation 3: Build accept ownership batch
// ======================================================
// ======================================================

type AcceptOwnershipProposalInput struct {
	ContractsByChain map[uint64]string
	Description      string
	MCMSConfig       proposalutils.TimelockConfig
}

func BuildAcceptOwnershipTimelockProposal(
	e cldf.Environment,
	input AcceptOwnershipProposalInput,
) (*mcms.TimelockProposal, error) {
	if len(input.ContractsByChain) == 0 {
		return nil, errors.New("no contracts provided to build accept ownership proposal")
	}

	var batches []mcmstypes.BatchOperation
	timelockAddresses := make(map[uint64]string)
	mcmAddressByChain := make(map[uint64]string)

	for chainSelector, contractAddr := range input.ContractsByChain {
		chain, ok := e.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return nil, fmt.Errorf("chain not found in environment: %d", chainSelector)
		}

		timelockAddr, err := getRequiredContractAddress(
			e.DataStore,
			chainSelector,
			commontypes.RBACTimelock,
		)
		if err != nil {
			return nil, fmt.Errorf("chain %d: %w", chainSelector, err)
		}

		mcmsAddr, err := getRequiredContractAddress(
			e.DataStore,
			chainSelector,
			ethBalMonMCMSContractTypeForProposal(&input.MCMSConfig),
		)
		if err != nil {
			return nil, fmt.Errorf("chain %d: %w", chainSelector, err)
		}

		ethBalMon, err := eth_balance_monitor_wrapper.NewEthBalanceMonitor(
			common.HexToAddress(contractAddr),
			chain.Client,
		)
		if err != nil {
			return nil, fmt.Errorf("chain %d: failed to instantiate EthBalanceMonitor at %s: %w", chainSelector, contractAddr, err)
		}

		acceptOwnershipTx, err := ethBalMon.AcceptOwnership(cldf.SimTransactOpts())
		if err != nil {
			return nil, fmt.Errorf("chain %d: failed to generate acceptOwnership calldata: %w", chainSelector, err)
		}

		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					OperationMetadata: mcmstypes.OperationMetadata{
						ContractType: vaulttypes.EthBalMonContractType,
						Tags:         []string{"acceptOwnership"},
					},
					To:               contractAddr,
					Data:             acceptOwnershipTx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		})

		timelockAddresses[chainSelector] = timelockAddr
		mcmAddressByChain[chainSelector] = mcmsAddr
	}

	description := input.Description
	if description == "" {
		description = "Accept ownership of EthBalanceMonitor across chains"
	}

	tlCfg := input.MCMSConfig
	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		e,
		timelockAddresses,
		mcmAddressByChain,
		nil,
		batches,
		description,
		tlCfg,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build timelock proposal: %w", err)
	}

	return proposal, nil
}
