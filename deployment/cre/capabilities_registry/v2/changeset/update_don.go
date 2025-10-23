package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ cldf.ChangeSetV2[UpdateDONInput] = UpdateDON{}

type UpdateDONInput struct {
	RegistryQualifier string `json:"registryQualifier" yaml:"registryQualifier"`
	RegistryChainSel  uint64 `json:"registryChainSel" yaml:"registryChainSel"`

	// DONName to update, this is required
	DONName string `json:"donName" yaml:"donName"`
	// NewDonName is optional
	NewDonName string `json:"newDonName" yaml:"newDonName"`

	CapabilityConfigs []contracts.CapabilityConfig `json:"capabilityConfigs" yaml:"capabilityConfigs"` // if Config subfield is nil, a default config is used

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool `json:"force" yaml:"force"`

	MCMSConfig *ocr3.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type UpdateDON struct{}

func (u UpdateDON) VerifyPreconditions(_ cldf.Environment, config UpdateDONInput) error {
	if config.DONName == "" {
		return errors.New("must provide a non-empty DONName")
	}

	return nil
}

func (u UpdateDON) Apply(e cldf.Environment, config UpdateDONInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.RegistryChainSel, emptyQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(config.RegistryChainSel, config.RegistryQualifier)

	chain, ok := e.BlockChains.EVMChains()[config.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", config.RegistryChainSel)
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

	don, nodes, err := sequences.GetDonNodes(config.DONName, capReg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get DON %s nodes: %w", config.DONName, err)
	}

	p2pIDs := make([]p2pkey.PeerID, 0)
	for _, node := range nodes {
		p2pIDs = append(p2pIDs, node.P2pId)
	}

	updateDonReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.UpdateDON,
		contracts.UpdateDONDeps{
			Env:                  &e,
			CapabilitiesRegistry: capReg,
			MCMSContracts:        mcmsContracts,
		},
		contracts.UpdateDONInput{
			ChainSelector:     config.RegistryChainSel,
			P2PIDs:            p2pIDs,
			CapabilityConfigs: config.CapabilityConfigs,
			DonName:           config.DONName,
			NewDonName:        config.NewDonName,
			F:                 don.F,
			IsPrivate:         !don.IsPublic,
			Force:             config.Force,
			MCMSConfig:        config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to update DON %s: %w", config.DONName, err)
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{updateDonReport.ToGenericReport()},
		MCMSTimelockProposals: updateDonReport.Output.Proposals,
	}, nil
}
