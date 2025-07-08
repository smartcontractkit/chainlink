package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type AddCapabilitiesOpDeps struct {
	Chain             evm.Chain
	Contract          *capabilities_registry.CapabilitiesRegistry
	DonToCapabilities map[string][]internal.RegisteredCapability
}

type AddCapabilitiesOpInput struct {
	UseMCMS bool
}

type AddCapabilitiesOpOutput struct {
	BatchOperation *mcmstypes.BatchOperation
}

var AddCapabilitiesOp = operations.NewOperation[AddCapabilitiesOpInput, AddCapabilitiesOpOutput, AddCapabilitiesOpDeps](
	"add-capabilities-op",
	semver.MustParse("1.0.0"),
	"Add Capabilities to Capabilities Registry",
	func(b operations.Bundle, deps AddCapabilitiesOpDeps, input AddCapabilitiesOpInput) (AddCapabilitiesOpOutput, error) {
		var capabilities []capabilities_registry.CapabilitiesRegistryCapability
		for _, don := range deps.DonToCapabilities {
			for _, donCap := range don {
				capabilities = append(capabilities, donCap.CapabilitiesRegistryCapability)
			}
		}
		batchOp, err := internal.AddCapabilities(b.Logger, deps.Contract, deps.Chain, capabilities, input.UseMCMS)
		if err != nil {
			return AddCapabilitiesOpOutput{}, fmt.Errorf("add-capabilities-op failed: %w", err)
		}
		b.Logger.Info("Added capabilities to Capabilities Registry")

		return AddCapabilitiesOpOutput{BatchOperation: batchOp}, nil
	},
)
