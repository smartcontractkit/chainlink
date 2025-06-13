package opsutils

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
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

// EVMDeployInput is the input structure for an EVM deploy operation.
type EVMDeployInput[IN any] struct {
	// ChainSelector is the selector for the chain on which the contract will be deployed.
	ChainSelector uint64 `json:"chainSelector"`
	// DeployInput is the input data for the call.
	DeployInput IN `json:"deployInput"`
}

// EVMDeployOutput is the output structure for an EVM deploy operation.
// It contains the new address, the deployment transaction, and the type and version of the contract that was deployed.
type EVMDeployOutput struct {
	// Address is the address of the deployed contract.
	Address common.Address `json:"address"`
	// TypeAndVersion is the type and version of the contract that was deployed.
	TypeAndVersion string `json:"typeAndVersion"`
}

// VMDeployers defines the various deployer functions available for EVM-based chains.
// Currently, it defines an EVM deployer and a ZksyncVM deployer, but can be extended.
type VMDeployers[IN any] struct {
	DeployEVM      func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput IN) (common.Address, *types.Transaction, error)
	DeployZksyncVM func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, deployInput IN) (common.Address, error)
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
			opts := chain.DeployerKey
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

// NewEVMDeployOperation creates a new operation that deploys an EVM contract.
// Any interfacing with gethwrappers should happen in the deploy function.
func NewEVMDeployOperation[IN any](
	name string,
	version *semver.Version,
	description string,
	typeAndVersion cldf.TypeAndVersion,
	deployers VMDeployers[IN],
) *operations.Operation[EVMDeployInput[IN], EVMDeployOutput, cldf_evm.Chain] {
	return operations.NewOperation(
		name,
		version,
		description,
		func(b operations.Bundle, chain cldf_evm.Chain, input EVMDeployInput[IN]) (EVMDeployOutput, error) {
			if input.ChainSelector != chain.Selector {
				return EVMDeployOutput{}, fmt.Errorf("mismatch between inputted chain selector and selector defined within dependencies: %d != %d", input.ChainSelector, chain.Selector)
			}
			var (
				addr common.Address
				tx   *types.Transaction
				err  error
			)
			if chain.IsZkSyncVM {
				addr, err = deployers.DeployZksyncVM(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
					input.DeployInput,
				)
			} else {
				addr, tx, err = deployers.DeployEVM(
					chain.DeployerKey,
					chain.Client,
					input.DeployInput,
				)
			}
			if err != nil {
				b.Logger.Errorw("Failed to deploy contract", "typeAndVersion", typeAndVersion, "chain", chain.String(), "err", err.Error())
				return EVMDeployOutput{}, fmt.Errorf("failed to deploy %s on %s: %w", typeAndVersion, chain, err)
			}
			// Non-ZkSyncVM chains require manual confirmation of the deployment transaction.
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
