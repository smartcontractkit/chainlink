package opsutil

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Addressable is an interface that provides a method to get the address of a contract.
type Addressable interface {
	Address() common.Address
}

// EVMCallDeps contains the dependencies required for an EVM call operation.
// It includes the contract bindings, the chain on which the operation is performed, and transaction options.
type EVMCallDeps[C Addressable] struct {
	Contract C
	Chain    cldf.Chain
	Opts     *bind.TransactOpts
}

// EVMCallInput is the input structure for an EVM call operation.
// It contains the address of the contract, the chain selector, and the input data required for the call.
// Why not pull the address from contract bindings and the chain selector from the chain? Because we need both to be included in the report.
// This ensures we don't conflict with other operations that have the same call input but a different target address / target chain.
type EVMCallInput[IN any] struct {
	Address       common.Address
	ChainSelector uint64
	CallInput     IN
}

// EVMCallOutput is the output structure for an EVM call operation.
// It contains the transaction that was sent and the type of contract that was called.
type EVMCallOutput struct {
	Tx           *types.Transaction
	ContractType cldf.ContractType
}

// NewEVMCallOperation creates a new operation that performs an EVM call.
func NewEVMCallOperation[C Addressable, IN any](
	name string,
	version *semver.Version,
	description string,
	abi string,
	call func(contract C, opts *bind.TransactOpts, input IN) (EVMCallOutput, error),
) *operations.Operation[EVMCallInput[IN], EVMCallOutput, EVMCallDeps[C]] {
	return operations.NewOperation(
		name,
		version,
		description,
		func(b operations.Bundle, deps EVMCallDeps[C], input EVMCallInput[IN]) (EVMCallOutput, error) {
			if input.Address != deps.Contract.Address() {
				return EVMCallOutput{}, fmt.Errorf("mismatch between inputted address and address connected to bindings: %s != %s", input.Address, deps.Contract.Address())
			}
			if input.ChainSelector != deps.Chain.Selector {
				return EVMCallOutput{}, fmt.Errorf("mismatch between inputted chain selector and actual chain selector: %d != %d", input.ChainSelector, deps.Chain.Selector)
			}
			out, err := call(deps.Contract, deps.Opts, input.CallInput)
			if err != nil {
				return EVMCallOutput{}, fmt.Errorf("failed to call %s on %s: %w", name, deps.Contract.Address(), err)
			}
			if !deps.Opts.NoSend {
				_, err = cldf.ConfirmIfNoErrorWithABI(deps.Chain, out.Tx, abi, err)
				if err != nil {
					return EVMCallOutput{}, fmt.Errorf("failed to confirm %s tx: %w", name, err)
				}
			}
			return out, err
		},
	)
}

// UpdateCSOutputViaEVMCallSequence updates the ChangesetOutput with the results of an EVM call sequence.
// It appends the execution reports from the sequence report to the ChangesetOutput's reports.
// If the sequence execution was successful and MCMS configuration is provided, it builds a proposal.
func UpdateCSOutputViaEVMCallSequence[IN any](
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	csOutput cldf.ChangesetOutput,
	seqReport operations.SequenceReport[IN, map[uint64][]EVMCallOutput],
	seqErr error,
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
	for chainSel, outs := range seqReport.Output {
		for _, out := range outs {
			batchOperation, err := proposalutils.BatchOperationForChain(chainSel, out.Tx.To().Hex(), out.Tx.Data(),
				big.NewInt(0), string(out.ContractType), []string{})
			if err != nil {
				return csOutput, fmt.Errorf("failed to create batch operation for chain with selector %d: %w", chainSel, err)
			}
			batches = append(batches, batchOperation)

			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().Hex()
			inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return csOutput, fmt.Errorf("failed to get inspector for chain with selector %d: %w", chainSel, err)
			}
		}
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, mcmsCfg)
	if err != nil {
		return csOutput, fmt.Errorf("failed to get mcms contracts by chain: %w", err)
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		batches,
		mcmsDescription,
		*mcmsCfg,
	)
	if err != nil {
		return csOutput, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}

	csOutput.MCMSTimelockProposals = []mcmslib.TimelockProposal{*proposal}
	return csOutput, nil
}
