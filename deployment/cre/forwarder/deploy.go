package forwarder

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/sync/errgroup"

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

// DeployOp is an operation that deploys the Keystone Forwarder contract.
var DeployOp = operations.NewOperation[DeployOpInput, DeployOpOutput, DeployOpDeps](
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

type DeploySequenceDeps struct {
	Env *cldf.Environment // The environment in which the Keystone Forwarders will be deployed
}

type DeploySequenceInput struct {
	Targets   []uint64 // The target chains for the Keystone Forwarders
	Qualifier string   // The qualifier for the forwarder deployment
}

type DeploySequenceOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // The address book containing the deployed Keystone Forwarders
	Datastore   datastore.DataStore
}

var DeploySequence = operations.NewSequence[DeploySequenceInput, DeploySequenceOutput, DeploySequenceDeps](
	"deploy-keystone-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Forwarders",
	func(b operations.Bundle, deps DeploySequenceDeps, input DeploySequenceInput) (DeploySequenceOutput, error) {
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()
		contractErrGroup := &errgroup.Group{}
		for _, target := range input.Targets {
			contractErrGroup.Go(func() error {
				r, err := operations.ExecuteOperation(b, DeployOp, DeployOpDeps(deps), DeployOpInput{
					ChainSelector: target,
					Qualifier:     input.Qualifier,
				})
				if err != nil {
					return err
				}
				err = ab.Merge(r.Output.AddressBook)
				if err != nil {
					return fmt.Errorf("failed to save Keystone Forwarder address on address book for target %d: %w", target, err)
				}
				addrs, err := r.Output.Addresses.Fetch()
				if err != nil {
					return fmt.Errorf("failed to fetch Keystone Forwarder addresses for target %d: %w", target, err)
				}
				for _, addr := range addrs {
					if addrRefErr := as.AddressRefStore.Add(addr); addrRefErr != nil {
						return fmt.Errorf("failed to save Keystone Forwarder address on datastore for target %d: %w", target, addrRefErr)
					}
				}

				return nil
			})
		}
		if err := contractErrGroup.Wait(); err != nil {
			return DeploySequenceOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to deploy Keystone contracts: %w", err)
		}
		return DeploySequenceOutput{AddressBook: ab, Addresses: as.Addresses(), Datastore: as.Seal()}, nil
	},
)
