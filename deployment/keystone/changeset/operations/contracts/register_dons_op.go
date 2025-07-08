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

type RegisterDonsOpDeps struct {
	Env               *cldf.Environment
	RegistryChain     *evm.Chain
	Contract          *capabilities_registry.CapabilitiesRegistry
	DonToCapabilities map[string][]internal.RegisteredCapability
	DonToNodes        map[string][]deployment.Node
	Dons              []internal.DonCapabilities
}

type RegisterDonsOpInput struct {
	RegistryChainSel uint64
	NodeIDToParams   map[string]capabilities_registry.CapabilitiesRegistryNodeParams
	UseMCMS          bool
}

type RegisterDonsOpOutput struct {
	DonInfos       map[string]capabilities_registry.CapabilitiesRegistryDONInfo
	BatchOperation *mcmstypes.BatchOperation
}

var RegisterDonsOp = operations.NewOperation[RegisterDonsOpInput, RegisterDonsOpOutput, RegisterDonsOpDeps](
	"register-dons-op",
	semver.MustParse("1.0.0"),
	"Register Dons in Capabilities Registry",
	func(b operations.Bundle, deps RegisterDonsOpDeps, input RegisterDonsOpInput) (RegisterDonsOpOutput, error) {
		// TODO: annotate nodes with node_operator_id in JD?

		var donsToRegister []internal.DONToRegister
		for _, don := range deps.Dons {
			nodes, ok := deps.DonToNodes[don.Name]
			if !ok {
				return RegisterDonsOpOutput{}, fmt.Errorf("nodes not found for don %s", don.Name)
			}
			f := don.F
			if f == 0 {
				// TODO: fallback to a default value for compatibility - change to error
				f = uint8(len(nodes) / 3)
				b.Logger.Warnw("F not set for don - falling back to default", "don", don.Name, "f", f)
			}
			donsToRegister = append(donsToRegister, internal.DONToRegister{
				Name:  don.Name,
				F:     f,
				Nodes: nodes,
			})
		}

		nodeIDToP2PID := map[string][32]byte{}
		for nodeID, params := range input.NodeIDToParams {
			nodeIDToP2PID[nodeID] = params.P2pId
		}
		// register DONS
		donsResp, err := internal.RegisterDons(b.Logger, internal.RegisterDonsRequest{
			Env:                   deps.Env,
			RegistryChain:         deps.RegistryChain,
			Registry:              deps.Contract,
			RegistryChainSelector: input.RegistryChainSel,
			NodeIDToP2PID:         nodeIDToP2PID,
			DonToCapabilities:     deps.DonToCapabilities,
			DonsToRegister:        donsToRegister,
			UseMCMS:               input.UseMCMS,
		})
		if err != nil {
			return RegisterDonsOpOutput{}, fmt.Errorf("register-dons-op failed: %w", err)
		}
		b.Logger.Infow("registered DONs", "dons", len(donsResp.DonInfos))

		return RegisterDonsOpOutput{
			DonInfos:       donsResp.DonInfos,
			BatchOperation: donsResp.Ops,
		}, nil
	},
)
