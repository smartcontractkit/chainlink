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

type ConfigureCapabilitiesRegistryDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type ConfigureCapabilitiesRegistryInput struct {
	RegistryChainSel uint64
	RegistryRef      datastore.AddressRefKey
	MCMSConfig       *crecontracts.MCMSConfig
	Description      string
	// Deprecated: Use RegistryRef
	// TODO(PRODCRE-1030): Remove support for ContractAddress
	ContractAddress string
	Nops            []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	Nodes           []contracts.NodesInput
	Capabilities    []contracts.RegisterableCapability
	DONs            []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
}

func (c ConfigureCapabilitiesRegistryInput) Validate() error {
	if c.ContractAddress == "" && c.RegistryRef == nil {
		return errors.New("must set either registry ref or contract address")
	}
	if c.RegistryRef != nil && c.ContractAddress != "" {
		return errors.New("cannot set both address and registry ref")
	}
	return nil
}

type ConfigureCapabilitiesRegistryOutput struct {
	Capabilities          []contracts.RegisterableCapability
	Nops                  []capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams
	Nodes                 []capabilities_registry_v2.CapabilitiesRegistryNodeParams
	DONs                  []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
	MCMSTimelockProposals []mcmslib.TimelockProposal
}

var ConfigureCapabilitiesRegistry = operations.NewSequence(
	"configure-capabilities-registry",
	semver.MustParse("1.0.0"),
	"Configures the capabilities registry by registering node operators, nodes, dons and capabilities",
	func(b operations.Bundle, deps ConfigureCapabilitiesRegistryDeps, input ConfigureCapabilitiesRegistryInput) (ConfigureCapabilitiesRegistryOutput, error) {
		var allOperations []types.BatchOperation

		addr := input.ContractAddress
		if input.RegistryRef != nil {
			registryAddressRef, err := deps.Env.DataStore.Addresses().Get(input.RegistryRef)
			if err != nil {
				return ConfigureCapabilitiesRegistryOutput{}, fmt.Errorf("failed to get registry address: %w", err)
			}
			addr = registryAddressRef.Address
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
		if !ok {
			return ConfigureCapabilitiesRegistryOutput{}, fmt.Errorf("chain with selector %d not found", input.RegistryChainSel)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			common.HexToAddress(addr),
			contracts.ConfigureCapabilitiesRegistryDescription,
		)
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Register Node Operators
		registerNopsReport, err := operations.ExecuteOperation(b, contracts.RegisterNops, contracts.RegisterNopsDeps{
			Env:      deps.Env,
			Strategy: strategy,
		}, contracts.RegisterNopsInput{
			ChainSelector: input.RegistryChainSel,
			Address:       addr,
			Nops:          input.Nops,
			MCMSConfig:    input.MCMSConfig,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}
		if registerNopsReport.Output.Operation != nil {
			allOperations = append(allOperations, *registerNopsReport.Output.Operation)
		}

		// Register capabilities
		registerCapabilitiesReport, err := operations.ExecuteOperation(b, contracts.RegisterCapabilities, contracts.RegisterCapabilitiesDeps{
			Env:      deps.Env,
			Strategy: strategy,
		}, contracts.RegisterCapabilitiesInput{
			ChainSelector: input.RegistryChainSel,
			Address:       addr,
			Capabilities:  input.Capabilities,
			MCMSConfig:    input.MCMSConfig,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}
		if registerCapabilitiesReport.Output.Operation != nil {
			allOperations = append(allOperations, *registerCapabilitiesReport.Output.Operation)
		}

		// Register Nodes
		registerNodesReport, err := operations.ExecuteOperation(b, contracts.RegisterNodes, contracts.RegisterNodesDeps{
			Env:      deps.Env,
			Strategy: strategy,
		}, contracts.RegisterNodesInput{
			ChainSelector:     input.RegistryChainSel,
			Address:           addr,
			Nodes:             input.Nodes,
			MCMSConfig:        input.MCMSConfig,
			AllNOPsInContract: registerNopsReport.Output.AllContractExpectedNOPs,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}
		if registerNodesReport.Output.Operation != nil {
			allOperations = append(allOperations, *registerNodesReport.Output.Operation)
		}

		// Register DONs
		registerDONsReport, err := operations.ExecuteOperation(b, contracts.RegisterDons, contracts.RegisterDonsDeps{
			Env:      deps.Env,
			Strategy: strategy,
		}, contracts.RegisterDonsInput{
			ChainSelector: input.RegistryChainSel,
			Address:       addr,
			DONs:          input.DONs,
			MCMSConfig:    input.MCMSConfig,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}
		if registerDONsReport.Output.Operation != nil {
			allOperations = append(allOperations, *registerDONsReport.Output.Operation)
		}

		var proposals []mcmslib.TimelockProposal

		if len(allOperations) > 0 {
			proposal, mErr := strategy.BuildProposal(allOperations)
			if mErr != nil {
				return ConfigureCapabilitiesRegistryOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", mErr)
			}

			proposals = append(proposals, *proposal)
		}

		return ConfigureCapabilitiesRegistryOutput{
			Nops:                  registerNopsReport.Output.Nops,
			Nodes:                 registerNodesReport.Output.Nodes,
			Capabilities:          registerCapabilitiesReport.Output.Capabilities,
			DONs:                  registerDONsReport.Output.DONs,
			MCMSTimelockProposals: proposals,
		}, nil
	},
)
