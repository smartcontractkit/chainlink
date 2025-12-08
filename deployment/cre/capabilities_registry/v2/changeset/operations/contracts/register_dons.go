package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type RegisterDonsDeps struct {
	Env      *cldf.Environment
	Strategy strategies.TransactionStrategy
}

type RegisterDonsInput struct {
	Address       string
	ChainSelector uint64
	DONs          []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
	MCMSConfig    *contracts.MCMSConfig
}

type RegisterDonsOutput struct {
	DONs      []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
	Operation *mcmstypes.BatchOperation
}

// RegisterDons is an operation that registers DONs in the V2 Capabilities Registry contract.
var RegisterDons = operations.NewOperation[RegisterDonsInput, RegisterDonsOutput, RegisterDonsDeps](
	"register-dons-op",
	semver.MustParse("1.0.0"),
	"Register DONs in Capabilities Registry",
	func(b operations.Bundle, deps RegisterDonsDeps, input RegisterDonsInput) (RegisterDonsOutput, error) {
		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterDonsOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the CapabilitiesRegistryTransactor contract
		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterDonsOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
		}

		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return capReg.AddDONs(opts, input.DONs)
		})
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterDonsOutput{}, fmt.Errorf("failed to execute AddDONs: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RegisterDons on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully registered %d DONs on chain %d", len(input.DONs), input.ChainSelector)
		}

		return RegisterDonsOutput{
			DONs:      input.DONs,
			Operation: operation,
		}, nil
	},
)
