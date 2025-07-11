package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type DeployCapabilityRegistryOpDeps struct {
	Env *cldf.Environment
}

type DeployCapabilityRegistryInput struct {
	ChainSelector uint64
}

type DeployCapabilityRegistryOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // Keeping the address book for backward compatibility, as not everything has been migrated to datastore
}

// DeployCapabilityRegistryOp is an operation that deploys the Capability Registry contract.
var DeployCapabilityRegistryOp = operations.NewOperation[DeployCapabilityRegistryInput, DeployCapabilityRegistryOutput, DeployCapabilityRegistryOpDeps](
	"deploy-capability-registry-op",
	semver.MustParse("1.0.0"),
	"Deploy CapabilityRegistry Contract",
	func(b operations.Bundle, deps DeployCapabilityRegistryOpDeps, input DeployCapabilityRegistryInput) (DeployCapabilityRegistryOutput, error) {
		capabilityRegistryOutput, err := changeset.DeployCapabilityRegistryV2(*deps.Env, &changeset.DeployRequestV2{
			ChainSel: input.ChainSelector,
		})
		if err != nil {
			return DeployCapabilityRegistryOutput{}, fmt.Errorf("DeployCapabilityRegistry error: failed to deploy capability registry: %w", err)
		}
		return DeployCapabilityRegistryOutput{
			Addresses:   capabilityRegistryOutput.DataStore.Addresses(),
			AddressBook: capabilityRegistryOutput.AddressBook, //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
		}, nil
	},
)
