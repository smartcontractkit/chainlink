package solana

import (
	"fmt"
	"maps"
	"slices"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/solana"
)

var _ cldf.ChangeSet[DeployForwarderRequest] = DeployForwarder

type DeployForwarderRequest struct {
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
}

func DeployForwarder(env deployment.Environment, req DeployForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	out.AddressBook = cldf.NewMemoryAddressBook() //nolint:staticcheck // TODO CRE-400
	out.DataStore = datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata]()

	selectors := req.ChainSelectors
	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(env.Chains))
	}

	for _, sel := range selectors {
		req := &DeployRequest{
			ChainSel: sel,
			deployFn: solana.DeployForwarder,
		}
		csOut, err := deploy(env, req)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", sel, err)
		}
		if err := out.AddressBook.Merge(csOut.AddressBook); err != nil { //nolint:staticcheck // TODO CRE-400
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for chain selector %d: %w", sel, err)
		}
		if err := out.DataStore.Merge(csOut.DataStore.Seal()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge datastore for chain selector %d: %w", sel, err)
		}
	}
	// convert all the addresses to t
	return out, nil
}
