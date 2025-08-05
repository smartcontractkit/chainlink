package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type RegisterNodesOpDeps struct {
	Env               *cldf.Environment
	RegistryChain     *evm.Chain
	Contract          *capabilities_registry.CapabilitiesRegistry
	DonToCapabilities map[string][]internal.RegisteredCapability
	NopsToNodeIDs     map[capabilities_registry.CapabilitiesRegistryNodeOperator][]string
	DonToNodes        map[string][]deployment.Node
}

type RegisterNodesOpInput struct {
	RegistryChainSel uint64
	Nops             []*capabilities_registry.CapabilitiesRegistryNodeOperatorAdded
	UseMCMS          bool
}

type RegisterNodesOpOutput struct {
	NodeIDToParams map[string]capabilities_registry.CapabilitiesRegistryNodeParams
	BatchOperation *mcmstypes.BatchOperation
}

var RegisterNodesOp = operations.NewOperation[RegisterNodesOpInput, RegisterNodesOpOutput, RegisterNodesOpDeps](
	"register-nodes-op",
	semver.MustParse("1.0.0"),
	"Register Nodes in Capabilities Registry",
	func(b operations.Bundle, deps RegisterNodesOpDeps, input RegisterNodesOpInput) (RegisterNodesOpOutput, error) {
		nodesResp, err := internal.RegisterNodes(b.Logger, &internal.RegisterNodesRequest{
			Env:                   deps.Env,
			RegistryChainSelector: input.RegistryChainSel,
			Registry:              deps.Contract,
			RegistryChain:         deps.RegistryChain,
			NopToNodeIDs:          deps.NopsToNodeIDs,
			DonToNodes:            deps.DonToNodes,
			DonToCapabilities:     deps.DonToCapabilities,
			Nops:                  input.Nops,
			UseMCMS:               input.UseMCMS,
		})
		if err != nil {
			return RegisterNodesOpOutput{}, fmt.Errorf("register-nodes-op failed: %w", err)
		}
		b.Logger.Infow("registered nodes", "nodes", nodesResp.NodeIDToParams)

		return RegisterNodesOpOutput{
			NodeIDToParams: nodesResp.NodeIDToParams,
			BatchOperation: nodesResp.Ops,
		}, nil
	},
)
