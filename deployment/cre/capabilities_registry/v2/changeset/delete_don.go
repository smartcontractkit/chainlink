package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	opscontracts "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

var _ cldf.ChangeSetV2[DeleteDONsInput] = DeleteDONs{}

type DeleteDONsInput struct {
	RegistryQualifier string                   `json:"registryQualifier" yaml:"registryQualifier"`
	RegistryChainSel  uint64                   `json:"registryChainSel" yaml:"registryChainSel"`
	DonNames          []string                 `json:"donNames" yaml:"donNames"`
	MCMSConfig        *crecontracts.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type DeleteDONs struct{}

func (d DeleteDONs) VerifyPreconditions(_ cldf.Environment, cfg DeleteDONsInput) error {
	if len(cfg.DonNames) == 0 {
		return errors.New("must provide at least one DON name")
	}
	for _, n := range cfg.DonNames {
		if n == "" {
			return errors.New("donNames cannot contain an empty string")
		}
	}
	return nil
}

func (d DeleteDONs) Apply(e cldf.Environment, cfg DeleteDONsInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if cfg.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, cfg.RegistryChainSel, *cfg.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(cfg.RegistryChainSel, cfg.RegistryQualifier)

	chain, ok := e.BlockChains.EVMChains()[cfg.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", cfg.RegistryChainSel)
	}

	registryAddressRef, err := e.DataStore.Addresses().Get(registryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get registry address: %w", err)
	}

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(registryAddressRef.Address), chain.Client,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
	}

	// Validate each DON exists up-front for clear failures
	for _, name := range cfg.DonNames {
		if _, getErr := capReg.GetDONByName(nil, name); getErr != nil {
			return cldf.ChangesetOutput{}, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, getErr)
		}
	}

	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		cfg.MCMSConfig,
		mcmsContracts,
		capReg.Address(),
		"Delete DON(s)",
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	deleteReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		opscontracts.DeleteDON,
		opscontracts.DeleteDONDeps{
			Env:                  &e,
			Strategy:             strategy,
			CapabilitiesRegistry: capReg,
		},
		opscontracts.DeleteDONInput{
			ChainSelector: cfg.RegistryChainSel,
			DonNames:      cfg.DonNames,
			MCMSConfig:    cfg.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to delete DONs %v: %w", cfg.DonNames, err)
	}

	var proposals []mcmslib.TimelockProposal
	if deleteReport.Output.Operation != nil {
		proposal, mcmsErr := strategy.BuildProposal([]types.BatchOperation{*deleteReport.Output.Operation})
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal for DeleteDON on chain %d: %w", cfg.RegistryChainSel, mcmsErr)
		}
		proposals = append(proposals, *proposal)
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{deleteReport.ToGenericReport()},
		MCMSTimelockProposals: proposals,
	}, nil
}
