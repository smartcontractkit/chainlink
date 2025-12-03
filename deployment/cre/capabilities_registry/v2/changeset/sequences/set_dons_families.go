package sequences

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type SetDONsFamiliesDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type SetDONsFamiliesInput struct {
	DONsChanges []DONFamiliesChange

	RegistryRef datastore.AddressRefKey

	MCMSConfig *crecontracts.MCMSConfig
}

func (i *SetDONsFamiliesInput) Validate() error {
	if len(i.DONsChanges) == 0 {
		return errors.New("must specify at least one DON change")
	}
	return nil
}

type DONFamiliesChange struct {
	DonName            string   `json:"donName" yaml:"donName"`
	AddToFamilies      []string `json:"addToFamilies" yaml:"addToFamilies"`
	RemoveFromFamilies []string `json:"removeFromFamilies" yaml:"removeFromFamilies"`
}

type SetDONsFamiliesOutput struct {
	DonsInfo  []capabilities_registry_v2.CapabilitiesRegistryDONInfo
	Proposals []mcmslib.TimelockProposal
}

var SetDONsFamilies = operations.NewSequence[SetDONsFamiliesInput, SetDONsFamiliesOutput, SetDONsFamiliesDeps](
	"set-dons-families-seq",
	semver.MustParse("1.0.0"),
	"Set DONs Families in Capabilities Registry",
	func(b operations.Bundle, deps SetDONsFamiliesDeps, input SetDONsFamiliesInput) (SetDONsFamiliesOutput, error) {
		if err := input.Validate(); err != nil {
			return SetDONsFamiliesOutput{}, err
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryRef.ChainSelector()]
		if !ok {
			return SetDONsFamiliesOutput{}, cldf.ErrChainNotFound
		}

		registryAddressRef, err := deps.Env.DataStore.Addresses().Get(input.RegistryRef)
		if err != nil {
			return SetDONsFamiliesOutput{}, fmt.Errorf("failed to get registry address: %w", err)
		}

		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(registryAddressRef.Address), chain.Client,
		)
		if err != nil {
			return SetDONsFamiliesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			capReg.Address(),
			contracts.SetDONFamiliesDescription,
		)
		if err != nil {
			return SetDONsFamiliesOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		var mcmsOperations []types.BatchOperation
		var donsInfo []capabilities_registry_v2.CapabilitiesRegistryDONInfo

		for _, change := range input.DONsChanges {
			report, err := operations.ExecuteOperation(
				b,
				contracts.SetDONFamilies,
				contracts.SetDONFamiliesDeps{
					Env:                  deps.Env,
					CapabilitiesRegistry: capReg,
					Strategy:             strategy,
				},
				contracts.SetDONFamiliesInput{
					DonName:            change.DonName,
					AddToFamilies:      change.AddToFamilies,
					RemoveFromFamilies: change.RemoveFromFamilies,
					MCMSConfig:         input.MCMSConfig,
					RegistryChainSel:   input.RegistryRef.ChainSelector(),
				},
			)
			if err != nil {
				return SetDONsFamiliesOutput{}, fmt.Errorf("failed to set families for DON %s: %w", change.DonName, err)
			}

			donsInfo = append(donsInfo, report.Output.DonInfo)
			if report.Output.Operation != nil {
				mcmsOperations = append(mcmsOperations, *report.Output.Operation)
			}
		}

		var proposals []mcmslib.TimelockProposal

		if len(mcmsOperations) > 0 {
			proposal, err := strategy.BuildProposal(mcmsOperations)
			if err != nil {
				return SetDONsFamiliesOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
			}

			proposals = append(proposals, *proposal)
		}

		return SetDONsFamiliesOutput{
			DonsInfo:  donsInfo,
			Proposals: proposals,
		}, nil
	},
)
