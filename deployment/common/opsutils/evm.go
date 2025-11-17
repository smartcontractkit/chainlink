package opsutils

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// EVMCallInput is the input structure for an EVM call operation.
// Why not pull the chain selector from the chain dependency? Because addresses might be the same across chains and we need to differentiate them.
// This ensures no false report matches between operation runs that have the same call input and address but a different target chain.
type EVMCallInput[IN any] struct {
	// Address is the address of the contract to call.
	Address common.Address `json:"address"`
	// ChainSelector is the selector for the chain on which the contract resides.
	ChainSelector uint64 `json:"chainSelector"`
	// CallInput is the input data for the call.
	CallInput IN `json:"callInput"`
	// NoSend indicates whether or not the transaction should be sent.
	// If true, the transaction data be prepared and returned but not sent.
	NoSend bool `json:"noSend"`
	// GasPrice is a custom gas price to set for the transaction.
	GasPrice uint64 `json:"gasPrice"`
	// GasLimit is a custom gas limit to set for the transaction.
	GasLimit uint64 `json:"gasLimit"`
}

// EVMCallOutput is the output structure for an EVM call operation.
// It contains the transaction and the type of contract that is being called.
type EVMCallOutput struct {
	// To is the address that initiated the transaction.
	To common.Address `json:"to"`
	// Data is the transaction data
	Data []byte `json:"data"`
	// ContractType is the type of contract that is being called.
	ContractType cldf.ContractType `json:"contractType"`
	// Confirmed indicates whether or not the transaction was confirmed.
	Confirmed bool `json:"confirmed"`
}

// NewEVMCallOperation creates a new operation that performs an EVM call.
// Any interfacing with gethwrappers should happen in the call function.
func NewEVMCallOperation[IN any, C any](
	name string,
	version *semver.Version,
	description string,
	abi string,
	contractType cldf.ContractType,
	constructor func(address common.Address, backend bind.ContractBackend) (C, error),
	call func(contract C, opts *bind.TransactOpts, input IN) (*types.Transaction, error),
) *operations.Operation[EVMCallInput[IN], EVMCallOutput, cldf_evm.Chain] {
	return operations.NewOperation(
		name,
		version,
		description,
		func(b operations.Bundle, chain cldf_evm.Chain, input EVMCallInput[IN]) (EVMCallOutput, error) {
			if input.ChainSelector != chain.Selector {
				return EVMCallOutput{}, fmt.Errorf("mismatch between inputted chain selector and selector defined within dependencies: %d != %d", input.ChainSelector, chain.Selector)
			}
			opts := CloneTransactOptsWithGas(chain.DeployerKey, input.GasLimit, input.GasPrice)
			if input.NoSend {
				opts = cldf.SimTransactOpts()
			}
			contract, err := constructor(input.Address, chain.Client)
			if err != nil {
				return EVMCallOutput{}, fmt.Errorf("failed to create contract instance for %s at %s on %s: %w", name, input.Address, chain, err)
			}
			tx, err := call(contract, opts, input.CallInput)
			confirmed := false
			if !input.NoSend {
				// If the call has actually been sent, we need check the call error and confirm the transaction.
				_, err := cldf.ConfirmIfNoErrorWithABI(chain, tx, abi, err)
				if err != nil {
					return EVMCallOutput{}, fmt.Errorf("failed to confirm %s tx against %s on %s: %w", name, input.Address, chain, err)
				}
				b.Logger.Debugw(fmt.Sprintf("Confirmed %s tx against %s on %s", name, input.Address, chain), "hash", tx.Hash().Hex(), "input", input.CallInput)
				confirmed = true
			} else {
				b.Logger.Debugw(fmt.Sprintf("Prepared %s tx against %s on %s", name, input.Address, chain), "input", input.CallInput)
			}
			return EVMCallOutput{
				To:           input.Address,
				Data:         tx.Data(),
				ContractType: contractType,
				Confirmed:    confirmed,
			}, err
		},
	)
}

// AddEVMCallSequenceToCSOutput updates the ChangesetOutput with the results of an EVM call sequence.
// It appends the execution reports from the sequence report to the ChangesetOutput's reports.
// If the sequence execution was successful and MCMS configuration is provided, it adds a proposal to the output.
func AddEVMCallSequenceToCSOutput[IN any](
	e cldf.Environment,
	csOutput cldf.ChangesetOutput,
	seqReport operations.SequenceReport[IN, map[uint64][]EVMCallOutput],
	seqErr error,
	mcmsStateByChain map[uint64]state.MCMSWithTimelockState,
	mcmsCfg *proposalutils.TimelockConfig,
	mcmsDescription string,
) (cldf.ChangesetOutput, error) {
	defer func() { csOutput.Reports = append(csOutput.Reports, seqReport.ExecutionReports...) }()
	if seqErr != nil {
		return csOutput, fmt.Errorf("failed to execute %s: %w", seqReport.Def, seqErr)
	}

	// Return early if MCMS is not being used
	if mcmsCfg == nil {
		return csOutput, nil
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)
	mcmContractByChain := make(map[uint64]string)
	for chainSel, outs := range seqReport.Output {
		for _, out := range outs {
			// If a transaction has already been confirmed, we do not need an operation for it.
			// TODO: Instead of creating 1 batch operation per call, can we batch calls together based on some strategy?
			if out.Confirmed {
				continue
			}
			batchOperation, err := proposalutils.BatchOperationForChain(chainSel, out.To.Hex(), out.Data,
				big.NewInt(0), string(out.ContractType), []string{})
			if err != nil {
				return csOutput, fmt.Errorf("failed to create batch operation for chain with selector %d: %w", chainSel, err)
			}
			batches = append(batches, batchOperation)

			mcmsState, ok := mcmsStateByChain[chainSel]
			if !ok {
				return csOutput, fmt.Errorf("mcms state not found for chain with selector %d", chainSel)
			}
			timelocks[chainSel] = mcmsState.Timelock.Address().Hex()
			inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return csOutput, fmt.Errorf("failed to get inspector for chain with selector %d: %w", chainSel, err)
			}
			mcm, err := mcmsCfg.MCMBasedOnAction(mcmsState)
			if err != nil {
				return csOutput, fmt.Errorf("failed to get MCM contract for chain with selector %d: %w", chainSel, err)
			}
			mcmContractByChain[chainSel] = mcm.Address().Hex()
		}
	}
	if len(batches) == 0 {
		return csOutput, nil
	}

	// Build new proposal from the batches and MCMS configuration.
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmContractByChain,
		inspectors,
		batches,
		mcmsDescription,
		*mcmsCfg,
	)
	if err != nil {
		return csOutput, fmt.Errorf("failed to build mcms proposal: %w", err)
	}

	// Add the new proposal to the ChangesetOutput.
	if csOutput.MCMSTimelockProposals == nil {
		csOutput.MCMSTimelockProposals = make([]mcmslib.TimelockProposal, 1)
	}
	csOutput.MCMSTimelockProposals = append(csOutput.MCMSTimelockProposals, *proposal)
	// Aggregate the proposals into a single proposal.
	// Aggregate the descriptions of all proposals into a single string.
	var builder strings.Builder
	for i, prop := range csOutput.MCMSTimelockProposals {
		builder.WriteString(prop.Description)
		if i < len(csOutput.MCMSTimelockProposals)-1 {
			builder.WriteString(", ")
		}
	}
	aggProposal, err := proposalutils.AggregateProposals(e, mcmsStateByChain, nil, csOutput.MCMSTimelockProposals, builder.String(), mcmsCfg)
	if err != nil {
		return csOutput, fmt.Errorf("failed to aggregate proposals: %w", err)
	}

	csOutput.MCMSTimelockProposals = []mcmslib.TimelockProposal{*aggProposal}
	return csOutput, nil
}

// EVMDeployInput is the input structure for an EVM deploy operation.
type EVMDeployInput[IN any] struct {
	// ChainSelector is the selector for the chain on which the contract will be deployed.
	ChainSelector uint64 `json:"chainSelector"`
	// DeployInput is the input data for the call.
	DeployInput IN `json:"deployInput"`
	// GasPrice is a custom gas price to set for the transaction.
	GasPrice uint64 `json:"gasPrice"`
	// GasLimit is a custom gas limit to set for the transaction.
	GasLimit uint64 `json:"gasLimit"`
	// Qualifier is an optional qualifier for the deployment.
	Qualifier *string `json:"qualifier"`
	// ContractOpts (optional) further configure the deployment with a specific bytecode and version.
	ContractOpts *ContractOpts `json:"contractOpts"`
}

// EVMDeployOutput is the output structure for an EVM deploy operation.
// It contains the new address, the deployment transaction, and the type and version of the contract that was deployed.
type EVMDeployOutput struct {
	// Address is the address of the deployed contract.
	Address common.Address `json:"address"`
	// TypeAndVersion is the type and version of the contract that was deployed.
	TypeAndVersion string `json:"typeAndVersion"`
	// Qualifier is an optional qualifier for the deployment.
	Qualifier *string `json:"qualifier"`
}

// ContractOpts specify the exact bytecode and version of the contract to deploy.
// Deployment operations must define defaults for these options in case users do not provide them.
// These options allow operators to deploy new bytecodes for the same ABI.
type ContractOpts struct {
	Version          *semver.Version
	EVMBytecode      []byte
	ZkSyncVMBytecode []byte
}

func (c *ContractOpts) Validate(isZkSyncVM bool) error {
	if c.Version == nil {
		return errors.New("version must be defined")
	}
	if isZkSyncVM && len(c.ZkSyncVMBytecode) == 0 {
		return errors.New("zkSyncVM bytecode must be defined")
	}
	if !isZkSyncVM && len(c.EVMBytecode) == 0 {
		return errors.New("evm bytecode must be defined")
	}
	return nil
}

// NewEVMDeployOperation creates a new operation that deploys an EVM contract.
// Any interfacing with gethwrappers should happen in the deploy function.
func NewEVMDeployOperation[IN any](
	name string,
	version *semver.Version,
	description string,
	contractType cldf.ContractType,
	contractMetadata *bind.MetaData,
	defaultContractOpts *ContractOpts,
	makeArgs func(IN) []any,
) *operations.Operation[EVMDeployInput[IN], EVMDeployOutput, cldf_evm.Chain] {
	return operations.NewOperation(
		name,
		version,
		description,
		func(b operations.Bundle, chain cldf_evm.Chain, input EVMDeployInput[IN]) (EVMDeployOutput, error) {
			if input.ChainSelector != chain.Selector {
				return EVMDeployOutput{}, fmt.Errorf("mismatch between inputted chain selector and selector defined within dependencies: %d != %d", input.ChainSelector, chain.Selector)
			}
			if contractMetadata == nil {
				return EVMDeployOutput{}, errors.New("contract metadata must be provided for deployment")
			}
			contractOpts := defaultContractOpts
			if input.ContractOpts != nil {
				contractOpts = input.ContractOpts
			}
			if contractOpts == nil {
				return EVMDeployOutput{}, errors.New("must define ContractOpts for deployment, no defaults provided")
			}
			if err := contractOpts.Validate(chain.IsZkSyncVM); err != nil {
				return EVMDeployOutput{}, fmt.Errorf("invalid ContractOpts: %w", err)
			}
			typeAndVersion := cldf.NewTypeAndVersion(contractType, *contractOpts.Version)
			parsedABI, err := contractMetadata.GetAbi()
			if err != nil {
				return EVMDeployOutput{}, fmt.Errorf("failed to parse ABI for %s: %w", typeAndVersion, err)
			}
			if parsedABI == nil {
				return EVMDeployOutput{}, fmt.Errorf("ABI is nil for %s", typeAndVersion)
			}

			var (
				addr common.Address
				tx   *types.Transaction
			)
			if chain.IsZkSyncVM {
				addr, err = deployZkContract(
					nil,
					contractOpts.ZkSyncVMBytecode,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					parsedABI,
					makeArgs(input.DeployInput)...,
				)
			} else {
				addr, tx, _, err = bind.DeployContract(
					CloneTransactOptsWithGas(chain.DeployerKey, input.GasLimit, input.GasPrice),
					*parsedABI,
					contractOpts.EVMBytecode,
					chain.Client,
					makeArgs(input.DeployInput)...,
				)
			}
			if err != nil {
				b.Logger.Errorw("Failed to deploy contract", "typeAndVersion", typeAndVersion, "chain", chain.String(), "err", err.Error())
				return EVMDeployOutput{}, fmt.Errorf("failed to deploy %s on %s: %w", typeAndVersion, chain, err)
			}
			// ZkSync transactions are confirmed in deployZkContract
			if !chain.IsZkSyncVM {
				_, err := chain.Confirm(tx)
				if err != nil {
					b.Logger.Errorw("Failed to confirm deployment", "typeAndVersion", typeAndVersion, "chain", chain.String(), "err", err.Error())
					return EVMDeployOutput{}, fmt.Errorf("failed to confirm deployment of %s on %s: %w", typeAndVersion, chain, err)
				}
			}
			return EVMDeployOutput{
				Address:        addr,
				TypeAndVersion: typeAndVersion.String(),
				Qualifier:      input.Qualifier,
			}, err
		},
	)
}

// cloneTransactOptsWithGas ensures that we don't impact the transact opts used by other operations.
func CloneTransactOptsWithGas(opts *bind.TransactOpts, gasLimit uint64, gasPrice uint64) *bind.TransactOpts {
	if opts == nil {
		return nil
	}
	newOpts := *opts
	if gasLimit > 0 {
		newOpts.GasLimit = gasLimit
	}
	if gasPrice > 0 {
		newOpts.GasPrice = new(big.Int).SetUint64(gasPrice)
	}
	return &newOpts
}

// GasBoostConfigsForChainMap creates a map of GasBoostConfig pointers for each chain in the provided chainMap.
// If a chain selector exists in gasBoostConfigs, it uses that config; otherwise, it sets nil.
func GasBoostConfigsForChainMap[T any](chainMap map[uint64]T, gasBoostConfigs map[uint64]commontypes.GasBoostConfig) map[uint64]*commontypes.GasBoostConfig {
	cfgs := make(map[uint64]*commontypes.GasBoostConfig, len(chainMap))
	if gasBoostConfigs == nil || chainMap == nil { // in either case, gas boosting should be empty
		return cfgs
	}

	for chainSelector := range chainMap {
		if _, ok := gasBoostConfigs[chainSelector]; ok {
			gasBoostConfig := gasBoostConfigs[chainSelector]
			cfgs[chainSelector] = &gasBoostConfig
		} else {
			cfgs[chainSelector] = nil
		}
	}

	return cfgs
}

// RetryDeploymentWithGasBoost is an ExecuteOption that retries EVM deployments with gas boosting.
// It uses the provided GasBoostConfig to adjust the gas limit and gas price on each retry attempt.
func RetryDeploymentWithGasBoost[IN any](cfg *commontypes.GasBoostConfig) operations.ExecuteOption[EVMDeployInput[IN], cldf_evm.Chain] {
	if cfg == nil {
		return withoutRetry[EVMDeployInput[IN], cldf_evm.Chain]()
	}
	c := *cfg

	return operations.WithRetryInput(func(attempt uint, err error, in EVMDeployInput[IN], deps cldf_evm.Chain) EVMDeployInput[IN] {
		gasLimit, gasPrice := GetBoostedGasForAttempt(c, attempt)
		in.GasLimit = gasLimit
		in.GasPrice = gasPrice

		return in
	})
}

// withoutRetry enables us to return an ExecuteOption that does nothing.
func withoutRetry[IN, DEP any]() operations.ExecuteOption[IN, DEP] {
	return func(c *operations.ExecuteConfig[IN, DEP]) {}
}

// RetryCallWithGasBoost is an ExecuteOption that retries EVM calls with gas boosting.
// It uses the provided GasBoostConfig to adjust the gas limit and gas price on each retry attempt.
// If NoSend is true, it will not apply gas boosting since the transaction is never sent.
func RetryCallWithGasBoost[IN any](cfg *commontypes.GasBoostConfig) operations.ExecuteOption[EVMCallInput[IN], cldf_evm.Chain] {
	// Use default retry option if no gas boost config is provided
	if cfg == nil {
		return operations.WithRetry[EVMCallInput[IN], cldf_evm.Chain]()
	}
	c := *cfg

	return operations.WithRetryInput(func(attempt uint, err error, in EVMCallInput[IN], deps cldf_evm.Chain) EVMCallInput[IN] {
		if in.NoSend {
			return in // No gas boost for calls that do not send transactions
		}

		gasLimit, gasPrice := GetBoostedGasForAttempt(c, attempt)
		in.GasLimit = gasLimit
		in.GasPrice = gasPrice

		return in
	})
}

func GetBoostedGasForAttempt(cfg commontypes.GasBoostConfig, attempt uint) (gasLimit uint64, gasPrice uint64) {
	initialGasLimit := uint64(200_000)          // 200k
	gasLimitIncrement := uint64(50_000)         // 50k
	initialGasPrice := uint64(20_000_000_000)   // 20 Gwei
	gasPriceIncrement := uint64(10_000_000_000) // 10 Gwei

	// Override defaults with config values if provided
	if cfg.InitialGasLimit > 0 {
		initialGasLimit = cfg.InitialGasLimit
	}
	if cfg.GasLimitIncrement > 0 {
		gasLimitIncrement = cfg.GasLimitIncrement
	}
	if cfg.InitialGasPrice > 0 {
		initialGasPrice = cfg.InitialGasPrice
	}
	if cfg.GasPriceIncrement > 0 {
		gasPriceIncrement = cfg.GasPriceIncrement
	}

	// initial + attempt * increment
	gasLimit = initialGasLimit + uint64(attempt)*gasLimitIncrement
	gasPrice = initialGasPrice + uint64(attempt)*gasPriceIncrement

	return
}

func deployZkContract(
	deployOpts *accounts.TransactOpts,
	bytecode []byte,
	client *clients.Client,
	wallet *accounts.Wallet,
	parsedABI *abi.ABI,
	args ...any,
) (common.Address, error) {
	var calldata []byte
	var err error
	if len(args) > 0 {
		calldata, err = parsedABI.Pack("", args...)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to pack constructor args: %w", err)
		}
	}

	salt := make([]byte, 32)
	n, err := rand.Read(salt)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to read random bytes: %w", err)
	}
	if n != len(salt) {
		return common.Address{}, fmt.Errorf("failed to read random bytes: expected %d, got %d", len(salt), n)
	}

	txHash, err := wallet.Deploy(deployOpts, accounts.Create2Transaction{
		Bytecode: bytecode,
		Calldata: calldata,
		Salt:     salt,
	})
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy zk contract: %w", err)
	}

	receipt, err := client.WaitMined(context.Background(), txHash)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to confirm zk contract deployment: %w", err)
	}

	return receipt.ContractAddress, nil
}
