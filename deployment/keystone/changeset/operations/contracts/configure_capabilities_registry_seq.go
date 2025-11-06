package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type ConfigureCapabilitiesRegistrySeqDeps struct {
	Env  *cldf.Environment
	Dons []internal.DonCapabilities // externally sourced based on the environment
}

type ConfigureCapabilitiesRegistrySeqInput struct {
	RegistryChainSel uint64

	UseMCMS         bool
	ContractAddress *common.Address
}

func (c ConfigureCapabilitiesRegistrySeqInput) Validate() error {
	if c.ContractAddress == nil {
		return errors.New("ContractAddress is not set")
	}
	_, ok := chainsel.ChainBySelector(c.RegistryChainSel)
	if !ok {
		return fmt.Errorf("chain %d not found in environment", c.RegistryChainSel)
	}
	return nil
}

type ConfigureCapabilitiesRegistrySeqOutput struct {
	DonInfos map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

var ConfigureCapabilitiesRegistrySeq = operations.NewSequence[ConfigureCapabilitiesRegistrySeqInput, ConfigureCapabilitiesRegistrySeqOutput, ConfigureCapabilitiesRegistrySeqDeps](
	"configure-capabilities-registry-seq",
	semver.MustParse("1.0.0"),
	"Configure Capabilities Registry",
	func(b operations.Bundle, deps ConfigureCapabilitiesRegistrySeqDeps, input ConfigureCapabilitiesRegistrySeqInput) (ConfigureCapabilitiesRegistrySeqOutput, error) {
		if err := input.Validate(); err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("input validation failed: %w", err)
		}
		for _, don := range deps.Dons {
			if err := don.Validate(); err != nil {
				return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("don validation failed for '%s': %w", don.Name, err)
			}
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
		if !ok {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", input.RegistryChainSel)
		}

		capabilitiesRegistry, err := contracts.GetOwnedContractV2[*capabilities_registry.CapabilitiesRegistry](deps.Env.DataStore.Addresses(), chain, input.ContractAddress.Hex())
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to get capabilities registry contract: %w", err)
		}
		if input.UseMCMS && capabilitiesRegistry.McmsContracts == nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("capabilities registry contract %s is not owned by MCMS", capabilitiesRegistry.Contract.Address())
		}

		donInfos, err := internal.DonInfos(deps.Dons, deps.Env.Offchain)
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to get don infos: %w", err)
		}

		// all the subsequent calls to the registry are in terms of nodes
		// compute the mapping of dons to their nodes for reuse in various registry calls
		donToNodes, err := internal.MapDonsToNodes(donInfos, true, input.RegistryChainSel)
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to map dons to nodes: %w", err)
		}

		// TODO: we can remove this abstractions and refactor the functions that accept them to accept []DonInfos/DonCapabilities
		// they are unnecessary indirection
		donToCapabilities, err := internal.MapDonsToCaps(capabilitiesRegistry.Contract, donInfos)
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to map dons to capabilities: %w", err)
		}
		nopsToNodeIDs, err := internal.NopsToNodes(donInfos, deps.Dons, input.RegistryChainSel)
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to map nops to nodes: %w", err)
		}

		_, err = operations.ExecuteOperation(b, AddCapabilitiesOp, AddCapabilitiesOpDeps{
			Chain:             chain,
			Contract:          capabilitiesRegistry.Contract,
			DonToCapabilities: donToCapabilities,
		}, AddCapabilitiesOpInput{
			UseMCMS: input.UseMCMS,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to add capabilities to registry: %w", err)
		}

		// register node operators
		nopsReport, err := operations.ExecuteOperation(b, RegisterNopsOp, RegisterNopsOpDeps{
			Env:           deps.Env,
			RegistryChain: &chain,
			NopsToNodes:   nopsToNodeIDs,
			Contract:      capabilitiesRegistry.Contract,
		}, RegisterNopsOpInput{
			UseMCMS:          input.UseMCMS,
			RegistryChainSel: input.RegistryChainSel,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to register node operators: %w", err)
		}
		nopsResp := nopsReport.Output

		// register nodes
		nodesReport, err := operations.ExecuteOperation(b, RegisterNodesOp, RegisterNodesOpDeps{
			Env:               deps.Env,
			RegistryChain:     &chain,
			Contract:          capabilitiesRegistry.Contract,
			NopsToNodeIDs:     nopsToNodeIDs,
			DonToNodes:        donToNodes,
			DonToCapabilities: donToCapabilities,
		}, RegisterNodesOpInput{
			RegistryChainSel: input.RegistryChainSel,
			Nops:             nopsResp.Nops,
			UseMCMS:          input.UseMCMS,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to register nodes: %w", err)
		}
		nodesResp := nodesReport.Output

		donsReport, err := operations.ExecuteOperation(b, RegisterDonsOp, RegisterDonsOpDeps{
			Env:               deps.Env,
			RegistryChain:     &chain,
			Contract:          capabilitiesRegistry.Contract,
			Dons:              deps.Dons,
			DonToNodes:        donToNodes,
			DonToCapabilities: donToCapabilities,
		}, RegisterDonsOpInput{
			RegistryChainSel: input.RegistryChainSel,
			NodeIDToParams:   nodesResp.NodeIDToParams,
			UseMCMS:          input.UseMCMS,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, fmt.Errorf("failed to register DONS: %w", err)
		}

		return ConfigureCapabilitiesRegistrySeqOutput{
			DonInfos: donsReport.Output.DonInfos,
		}, nil
	},
)
