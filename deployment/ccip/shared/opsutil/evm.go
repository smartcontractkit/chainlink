package opsutil

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
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
	Address common.Address
	// ChainSelector is the selector for the chain on which the contract resides.
	ChainSelector uint64
	// CallInput is the input data for the call.
	CallInput IN
	// NoSend indicates whether or not the transaction should be sent.
	// If true, the transaction data be prepared and returned but not sent.
	NoSend bool
}

// EVMCallOutput is the output structure for an EVM call operation.
// It contains the transaction and the type of contract that is being called.
type EVMCallOutput struct {
	// Tx is the transaction formed by the call.
	Tx *types.Transaction
	// ContractType is the type of contract that is being called.
	ContractType cldf.ContractType
	// Confirmed indicates whether or not the transaction was confirmed.
	Confirmed bool
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
				confirmed = true
			}
			return EVMCallOutput{
				Tx:           tx,
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
			// If a transaction has already been confirmed, we do not need an operation for it.
			// TODO: Instead of creating 1 batch operation per call, can we batch calls together based on some strategy?
			if out.Confirmed {
				continue
			}
			batchOperation, err := proposalutils.BatchOperationForChain(chainSel, out.Tx.To().Hex(), out.Tx.Data(),
				big.NewInt(0), string(out.ContractType), []string{})
			if err != nil {
				return csOutput, fmt.Errorf("failed to create batch operation for chain with selector %d: %w", chainSel, err)
			}
			batches = append(batches, batchOperation)

			if state.Chains[chainSel].Timelock == nil {
				return csOutput, fmt.Errorf("timelock not found for chain with selector %d", chainSel)
			}
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
		return csOutput, fmt.Errorf("failed to build mcms proposal: %w", err)
	}

	csOutput.MCMSTimelockProposals = []mcmslib.TimelockProposal{*proposal}
	return csOutput, nil
}
