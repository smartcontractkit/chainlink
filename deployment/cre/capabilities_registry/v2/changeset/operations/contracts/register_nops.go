package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type RegisterNopsDeps struct {
	Env      *cldf.Environment
	Strategy strategies.TransactionStrategy
}

type RegisterNopsInput struct {
	Address       string
	ChainSelector uint64
	Nops          []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	MCMSConfig    *contracts.MCMSConfig
}

type RegisterNopsOutput struct {
	Nops      []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	Operation *mcmstypes.BatchOperation
}

// RegisterNops is an operation that registers node operators in the V2 Capabilities Registry contract.
var RegisterNops = operations.NewOperation[RegisterNopsInput, RegisterNopsOutput, RegisterNopsDeps](
	"register-nops-op",
	semver.MustParse("1.0.0"),
	"Register Node Operators in Capabilities Registry",
	func(b operations.Bundle, deps RegisterNopsDeps, input RegisterNopsInput) (RegisterNopsOutput, error) {
		if len(input.Nops) == 0 {
			// The contract allows to pass an empty array of NOPs.
			return RegisterNopsOutput{
				Nops: []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{},
			}, nil
		}

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterNopsOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the NewCapabilitiesRegistry contract
		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterNopsOutput{}, fmt.Errorf("failed to create NewCapabilitiesRegistry: %w", err)
		}

		dedupedNOPs, err := dedupNOPs(deps.Env.Logger, input.Nops, capReg)
		if err != nil {
			return RegisterNopsOutput{}, fmt.Errorf("failed to dedupe NOPs: %w", err)
		}

		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return capReg.AddNodeOperators(opts, dedupedNOPs)
		})
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterNopsOutput{}, fmt.Errorf("failed to execute AddNodeOperators: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RegisterNops on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully registered %d node operators on chain %d", len(dedupedNOPs), input.ChainSelector)
		}

		return RegisterNopsOutput{
			Nops:      dedupedNOPs,
			Operation: operation,
		}, nil
	},
)

func dedupNOPs(lggr logger.Logger, inputNOPs []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams, capReg *capabilities_registry_v2.CapabilitiesRegistry) ([]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams, error) {
	contractNOPs, err := pkg.GetNodeOperators(nil, capReg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodes from contract: %w", err)
	}
	contractNOPsMap := make(map[string]struct{})
	for _, nop := range contractNOPs {
		contractNOPsMap[nop.Name] = struct{}{}
	}

	var dedupedNOPs []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	for i, nop := range inputNOPs {
		if _, exists := contractNOPsMap[nop.Name]; exists {
			lggr.Infof("NOP with name %s already registered in contract, skipping", nop.Name)
			continue
		}

		dedupedNOPs = append(dedupedNOPs, inputNOPs[i])
	}

	return dedupedNOPs, nil
}
