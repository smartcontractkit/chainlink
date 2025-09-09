package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
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

type AppendCapabilitiesOpDeps struct {
	Env               *cldf.Environment
	P2pToCapabilities map[p2pkey.PeerID][]capabilities_registry.CapabilitiesRegistryCapability
	RegistryRef       datastore.AddressRefKey
}

type AppendCapabilitiesOpInput struct {
	RegistryChainSel uint64
	MCMSConfig       *changeset.MCMSConfig
}

type AppendCapabilitiesOpOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var AppendCapabilitiesOp = operations.NewOperation[AppendCapabilitiesOpInput, AppendCapabilitiesOpOutput, AppendCapabilitiesOpDeps](
	"append-capabilities-op",
	semver.MustParse("1.0.0"),
	"Append Capabilities to Capabilities Registry",
	func(b operations.Bundle, deps AppendCapabilitiesOpDeps, input AppendCapabilitiesOpInput) (AppendCapabilitiesOpOutput, error) {
		changesetOutput, err := changeset.AppendNodeCapabilities(*deps.Env, &changeset.AppendNodeCapabilitiesRequest{
			RegistryChainSel:  input.RegistryChainSel,
			P2pToCapabilities: deps.P2pToCapabilities,
			MCMSConfig:        input.MCMSConfig,
			RegistryRef:       deps.RegistryRef,
		})
		if err != nil {
			return AppendCapabilitiesOpOutput{}, fmt.Errorf("append-capabilities-op failed: %w", err)
		}
		b.Logger.Info("Added capabilities to Capabilities Registry")
		if input.MCMSConfig != nil {
			return AppendCapabilitiesOpOutput{MCMSTimelockProposals: changesetOutput.MCMSTimelockProposals}, nil
		}
		return AppendCapabilitiesOpOutput{}, nil
	},
)

type UpdateDonOpDeps struct {
	Env    *cldf.Environment
	P2PIDs []p2pkey.PeerID // this is the unique identifier for the don
	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	RegistryRef       datastore.AddressRefKey
	CapabilityConfigs []internal.CapabilityConfig
}

type UpdateDonOpInput struct {
	RegistryChainSel uint64
	MCMSConfig       *changeset.MCMSConfig
}

type UpdateDonOpOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var UpdateDonOp = operations.NewOperation[UpdateDonOpInput, UpdateDonOpOutput, UpdateDonOpDeps](
	"update-don-op",
	semver.MustParse("1.0.0"),
	"Update Don in Capabilities Registry",
	func(b operations.Bundle, deps UpdateDonOpDeps, input UpdateDonOpInput) (UpdateDonOpOutput, error) {
		changesetOutput, err := changeset.UpdateDon(*deps.Env, &changeset.UpdateDonRequest{
			RegistryChainSel:  input.RegistryChainSel,
			P2PIDs:            deps.P2PIDs,
			MCMSConfig:        input.MCMSConfig,
			RegistryRef:       deps.RegistryRef,
			CapabilityConfigs: deps.CapabilityConfigs,
		})
		if err != nil {
			return UpdateDonOpOutput{}, fmt.Errorf("update-don-op failed: %w", err)
		}
		b.Logger.Info("Added capabilities to Capabilities Registry")
		if input.MCMSConfig != nil {
			return UpdateDonOpOutput{MCMSTimelockProposals: changesetOutput.MCMSTimelockProposals}, nil
		}
		return UpdateDonOpOutput{}, nil
	},
)
