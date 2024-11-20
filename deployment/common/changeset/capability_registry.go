package changeset

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

type HydrateConfig struct {
	ChainID uint64
}

// HydrateCapabilityRegistry deploys a new capabilities registry contract and hydrates it with the provided data.
func HydrateCapabilityRegistry(v v1_0.CapabilityRegistryView, env deployment.Environment, cfg HydrateConfig) (*deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry], error) {
	chainSelector, err := chainsel.SelectorFromChainId(cfg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain selector from chain id: %w", err)
	}
	chain, ok := env.Chains[chainSelector]
	if !ok {
		return nil, fmt.Errorf("chain with id %d not found", cfg.ChainID)
	}
	deployedContract, err := deployment.DeployContract(
		env.Logger, chain, env.ExistingAddresses,
		func(chain deployment.Chain) deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry] {
			crAddr, tx, cr, err2 := capabilities_registry.DeployCapabilitiesRegistry(
				chain.DeployerKey,
				chain.Client,
			)
			return deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
				Address: crAddr, Contract: cr, Tx: tx, Err: err2,
				Tv: deployment.NewTypeAndVersion("CapabilitiesRegistry", deployment.Version1_0_0),
			}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy contract: %w", err)
	}

	nopsParams := v.NopsToNopsParams()
	for _, nop := range v.Nops {
		nopsParams = append(nopsParams, capabilities_registry.CapabilitiesRegistryNodeOperator{
			Admin: nop.Admin,
			Name:  nop.Name,
		})
	}
	tx, err := deployedContract.Contract.AddNodeOperators(chain.DeployerKey, nopsParams)
	if _, err = deployment.ConfirmIfNoError(chain, tx, keystone.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)); err != nil {
		return nil, fmt.Errorf("failed to add node operators: %w", err)
	}

	capabilitiesParams := v.CapabilitiesToCapabilitiesParams()
	tx, err = deployedContract.Contract.AddCapabilities(chain.DeployerKey, capabilitiesParams)
	if _, err = deployment.ConfirmIfNoError(chain, tx, keystone.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)); err != nil {
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}

	nodesParams, err := v.NodesToNodesParams()
	if err != nil {
		return nil, fmt.Errorf("failed to convert nodes to nodes params: %w", err)
	}
	tx, err = deployedContract.Contract.AddNodes(chain.DeployerKey, nodesParams)
	if _, err = deployment.ConfirmIfNoError(chain, tx, keystone.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)); err != nil {
		return nil, fmt.Errorf("failed to add nodes: %w", err)
	}

	for _, don := range v.Dons {
		cfgs, err := v.CapabilityConfigToCapabilityConfigParams(don)
		if err != nil {
			return nil, fmt.Errorf("failed to convert capability configurations to capability configuration params: %w", err)
		}
		var peerIds [][32]byte
		for _, id := range don.NodeP2PIds {
			peerIds = append(peerIds, id)
		}
		tx, err = deployedContract.Contract.AddDON(chain.DeployerKey, peerIds, cfgs, don.IsPublic, don.AcceptsWorkflows, don.F)
		if _, err = deployment.ConfirmIfNoError(chain, tx, keystone.DecodeErr(capabilities_registry.CapabilitiesRegistryABI, err)); err != nil {
			return nil, fmt.Errorf("failed to add don: %w", err)
		}
	}

	return deployedContract, nil
}
