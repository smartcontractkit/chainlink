package forwarder

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type DeployOpDeps struct {
	Env *cldf.Environment
}

type DeployOpInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployOpOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // Keeping the address book for backward compatibility, as not everything has been migrated to datastore
}

// DeployKeystoneForwarderOp is an operation that deploys the Keystone Forwarder contract.
var DeployKeystoneForwarderOp = operations.NewOperation[DeployOpInput, DeployOpOutput, DeployOpDeps](
	"deploy-keystone-forwarder-op",
	semver.MustParse("1.0.0"),
	"Deploy KeystoneForwarder Contract",
	func(b operations.Bundle, deps DeployOpDeps, input DeployOpInput) (DeployOpOutput, error) {
		forwarderOutput, err := changeset.DeployForwarder(*deps.Env, changeset.DeployForwarderRequest{ChainSelectors: []uint64{input.ChainSelector}, Qualifier: input.Qualifier})
		if err != nil {
			return DeployOpOutput{}, err
		}
		return DeployOpOutput{
			Addresses:   forwarderOutput.DataStore.Addresses(),
			AddressBook: forwarderOutput.AddressBook, //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
		}, nil
	},
)
