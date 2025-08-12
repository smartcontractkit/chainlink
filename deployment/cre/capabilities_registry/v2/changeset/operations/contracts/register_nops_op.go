package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type RegisterNopsOpDeps struct {
	Env *cldf.Environment
}

type RegisterNopsOpInput struct {
	Address       string
	ChainSelector uint64
	Nops          []capabilities_registry_v2.CapabilitiesRegistryNodeOperator
}

type RegisterNopsOpOutput struct {
	Nops []*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded
}

// RegisterNopsOp is an operation that registers node operators in the V2 Capabilities Registry contract.
var RegisterNopsOp = operations.NewOperation[RegisterNopsOpInput, RegisterNopsOpOutput, RegisterNopsOpDeps](
	"register-nops-op",
	semver.MustParse("1.0.0"),
	"Register Node Operators in Capabilities Registry",
	func(b operations.Bundle, deps RegisterNopsOpDeps, input RegisterNopsOpInput) (RegisterNopsOpOutput, error) {
		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterNopsOpOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the CapabilitiesRegistryTransactor contract
		capabilityRegistryTransactor, err := capabilities_registry_v2.NewCapabilitiesRegistryTransactor(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterNopsOpOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
		}

		tx, err := capabilityRegistryTransactor.AddNodeOperators(chain.DeployerKey, input.Nops)
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterNopsOpOutput{}, fmt.Errorf("failed to call AddNodeOperators: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return RegisterNopsOpOutput{}, fmt.Errorf("failed to confirm AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
		}

		ctx := b.GetContext()
		receipt, err := bind.WaitMined(ctx, chain.Client, tx)
		if err != nil {
			return RegisterNopsOpOutput{}, fmt.Errorf("failed to mine AddNodeOperators confirm transaction %s: %w", tx.Hash().String(), err)
		}
		if len(receipt.Logs) != len(input.Nops) {
			return RegisterNopsOpOutput{}, fmt.Errorf("expected %d log entries for AddNodeOperators, got %d", len(input.Nops), len(receipt.Logs))
		}

		// Get the CapabilitiesRegistryFilterer contract
		capabilityRegistryFilterer, err := capabilities_registry_v2.NewCapabilitiesRegistryFilterer(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterNopsOpOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryFilterer: %w", err)
		}

		resp := RegisterNopsOpOutput{
			Nops: make([]*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded, 0, len(receipt.Logs)),
		}
		// Parse the logs to get the added node operators
		for i, log := range receipt.Logs {
			if log == nil {
				continue
			}

			o, err := capabilityRegistryFilterer.ParseNodeOperatorAdded(*log)
			if err != nil {
				return RegisterNopsOpOutput{}, fmt.Errorf("failed to parse log %d for operator added: %w", i, err)
			}
			resp.Nops = append(resp.Nops, o)
		}

		return resp, nil
	},
)
