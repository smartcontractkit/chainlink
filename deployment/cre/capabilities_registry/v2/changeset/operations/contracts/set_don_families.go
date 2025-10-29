package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type SetDONFamiliesDeps struct {
	Env                  *cldf.Environment
	CapabilitiesRegistry *capabilities_registry_v2.CapabilitiesRegistry
	MCMSContracts        *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type SetDONFamiliesInput struct {
	DonName            string
	AddToFamilies      []string
	RemoveFromFamilies []string

	RegistryChainSel uint64

	MCMSConfig *ocr3.MCMSConfig
}

func (i *SetDONFamiliesInput) Validate() error {
	if i.DonName == "" {
		return errors.New("must specify DonName")
	}

	if len(i.AddToFamilies) == 0 && len(i.RemoveFromFamilies) == 0 {
		return errors.New("must specify at least one family to add or remove")
	}

	return nil
}

type SetDONFamiliesOutput struct {
	DonInfo   capabilities_registry_v2.CapabilitiesRegistryDONInfo
	Proposals []mcmslib.TimelockProposal
}

var SetDONFamilies = operations.NewOperation[SetDONFamiliesInput, SetDONFamiliesOutput, SetDONFamiliesDeps](
	"set-don-families-op",
	semver.MustParse("1.0.0"),
	"Set DON Families in Capabilities Registry",
	func(b operations.Bundle, deps SetDONFamiliesDeps, input SetDONFamiliesInput) (SetDONFamiliesOutput, error) {
		if err := input.Validate(); err != nil {
			return SetDONFamiliesOutput{}, err
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
		if !ok {
			return SetDONFamiliesOutput{}, cldf.ErrChainNotFound
		}

		// Fetch the DON to get the ID. We don't want callers using the ID, since the name is more user-friendly.
		don, err := deps.CapabilitiesRegistry.GetDONByName(&bind.CallOpts{}, input.DonName)
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return SetDONFamiliesOutput{}, fmt.Errorf("failed to call GetDONByName: %w", err)
		}

		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			deps.CapabilitiesRegistry.Address(),
			SetDONFamiliesDescription,
		)
		if err != nil {
			return SetDONFamiliesOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		var resultDon capabilities_registry_v2.CapabilitiesRegistryDONInfo

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := deps.CapabilitiesRegistry.SetDONFamilies(opts, don.Id, input.AddToFamilies, input.RemoveFromFamilies)
			if err != nil {
				err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
				return nil, fmt.Errorf("failed to call SetDONFamilies: %w", err)
			}

			// For direct execution, we can confirm and get the updated DON info
			if input.MCMSConfig == nil {
				// Confirm transaction
				_, err = chain.Confirm(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to confirm SetDONFamilies transaction %s: %w", tx.Hash().String(), err)
				}

				ctx := b.GetContext()
				_, err = bind.WaitMined(ctx, chain.Client, tx)
				if err != nil {
					return nil, fmt.Errorf("failed to mine SetDONFamilies transaction %s: %w", tx.Hash().String(), err)
				}

				latestDON, err := deps.CapabilitiesRegistry.GetDON(&bind.CallOpts{}, don.Id)
				if err != nil {
					err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
					return nil, fmt.Errorf("failed to call GetDONByName: %w", err)
				}

				// Get the updated DON info
				resultDon = latestDON
			}

			return tx, nil
		})
		if err != nil {
			return SetDONFamiliesOutput{}, fmt.Errorf("failed to execute SetDONFamilies: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetDONFamilies '%s' on chain %d", input.DonName, input.RegistryChainSel)
		} else {
			deps.Env.Logger.Infof("Successfully set DON families '%s' on chain %d", input.DonName, input.RegistryChainSel)
		}

		return SetDONFamiliesOutput{
			DonInfo:   resultDon,
			Proposals: proposals,
		}, nil
	},
)
