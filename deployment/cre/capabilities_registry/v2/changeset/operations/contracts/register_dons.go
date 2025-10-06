package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type RegisterDonsDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type RegisterDonsInput struct {
	Address       string
	ChainSelector uint64
	DONs          []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
	MCMSConfig    *ocr3.MCMSConfig
}

type RegisterDonsOutput struct {
	DONs      []capabilities_registry_v2.CapabilitiesRegistryDONInfo
	Proposals []mcmslib.TimelockProposal
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
		capabilityRegistryTransactor, err := capabilities_registry_v2.NewCapabilitiesRegistryTransactor(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterDonsOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			common.HexToAddress(input.Address),
			RegisterDONsDescription,
		)
		if err != nil {
			return RegisterDonsOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		var resultDONs []capabilities_registry_v2.CapabilitiesRegistryDONInfo

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := capabilityRegistryTransactor.AddDONs(opts, input.DONs)
			if err != nil {
				err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
				return nil, fmt.Errorf("failed to call AddDONs: %w", err)
			}

			// For direct execution, we can get the DONs info
			if input.MCMSConfig == nil {
				// Confirm transaction
				_, err = chain.Confirm(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to confirm AddDONs transaction %s: %w", tx.Hash().String(), err)
				}

				// Get the CapabilitiesRegistryCaller contract
				capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
					common.HexToAddress(input.Address),
					chain.Client,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to create CapabilitiesRegistryCaller: %w", err)
				}

				// Fetch all DONs via generic pagination helper
				donsInfo, err := pkg.GetDONs(nil, capReg)
				if err != nil {
					return nil, fmt.Errorf("failed to call GetDONs: %w", err)
				}

				resultDONs = donsInfo
			}

			return tx, nil
		})
		if err != nil {
			return RegisterDonsOutput{}, fmt.Errorf("failed to execute AddDONs: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RegisterDons on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully registered %d DONs on chain %d", len(resultDONs), input.ChainSelector)
		}

		return RegisterDonsOutput{
			DONs:      resultDONs,
			Proposals: proposals,
		}, nil
	},
)
