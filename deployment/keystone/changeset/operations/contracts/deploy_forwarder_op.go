package contracts

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type DeployForwarderOpDeps struct {
	Env *cldf.Environment
}

type DeployForwarderOpInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployForwarderOpOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // Keeping the address book for backward compatibility, as not everything has been migrated to datastore
}

// DeployKeystoneForwarderOp is an operation that deploys the Keystone Forwarder contract.
var DeployKeystoneForwarderOp = operations.NewOperation[DeployForwarderOpInput, DeployForwarderOpOutput, DeployForwarderOpDeps](
	"deploy-keystone-forwarder-op",
	semver.MustParse("1.0.0"),
	"Deploy KeystoneForwarder Contract",
	func(b operations.Bundle, deps DeployForwarderOpDeps, input DeployForwarderOpInput) (DeployForwarderOpOutput, error) {
		forwarderOutput, err := changeset.DeployForwarder(*deps.Env, changeset.DeployForwarderRequest{ChainSelectors: []uint64{input.ChainSelector}, Qualifier: input.Qualifier})
		if err != nil {
			return DeployForwarderOpOutput{}, err
		}
		return DeployForwarderOpOutput{
			Addresses:   forwarderOutput.DataStore.Addresses(),
			AddressBook: forwarderOutput.AddressBook, //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
		}, nil
	},
)
