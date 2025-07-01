package contracts

import (
	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployKeystoneForwardersSequenceDeps struct {
	Env *cldf.Environment // The environment in which the Keystone Forwarders will be deployed
}

type DeployKeystoneForwardersInput struct {
	Targets []uint64 // The target chains for the Keystone Forwarders
}

type DeployKeystoneForwardersOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // The address book containing the deployed Keystone Forwarders
}

var DeployKeystoneForwardersSequence = operations.NewSequence[DeployKeystoneForwardersInput, DeployKeystoneForwardersOutput, DeployKeystoneForwardersSequenceDeps](
	"deploy-keystone-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Forwarders",
	func(b operations.Bundle, deps DeployKeystoneForwardersSequenceDeps, input DeployKeystoneForwardersInput) (DeployKeystoneForwardersOutput, error) {
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()
		contractErrGroup := &errgroup.Group{}
		for _, target := range input.Targets {
			contractErrGroup.Go(func() error {
				r, err := operations.ExecuteOperation(b, DeployKeystoneForwarderOp, DeployForwarderOpDeps(deps), DeployForwarderOpInput{
					ChainSelector: target,
				})
				if err != nil {
					return err
				}
				err = ab.Merge(r.Output.AddressBook)
				if err != nil {
					return pkgerrors.Wrapf(err, "failed to save Keystone Forwarder address on address book for target %d", target)
				}
				addrs, err := r.Output.Addresses.Fetch()
				if err != nil {
					return pkgerrors.Wrapf(err, "failed to fetch Keystone Forwarder addresses for target %d", target)
				}
				for _, addr := range addrs {
					if addrRefErr := as.AddressRefStore.Add(addr); addrRefErr != nil {
						return pkgerrors.Wrapf(addrRefErr, "failed to save Keystone Forwarder address on datastore for target %d", target)
					}
				}

				return nil
			})
		}
		if err := contractErrGroup.Wait(); err != nil {
			return DeployKeystoneForwardersOutput{AddressBook: ab, Addresses: as.Addresses()}, pkgerrors.Wrap(err, "failed to deploy Keystone contracts")
		}
		return DeployKeystoneForwardersOutput{AddressBook: ab, Addresses: as.Addresses()}, nil
	},
)
