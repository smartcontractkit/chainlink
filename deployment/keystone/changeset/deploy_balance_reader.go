package changeset

import (
	"fmt"
	"maps"
	"slices"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[DeployBalanceReaderRequest] = DeployBalanceReader

type DeployBalanceReaderRequest struct {
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
}

// DeployBalanceReader deploys the BalanceReader contract to all chains in the environment
// callers must merge the output addressbook with the existing one
// Deprecated: use DeployBalanceReaderV2 instead
func DeployBalanceReader(env cldf.Environment, cfg DeployBalanceReaderRequest) (cldf.ChangesetOutput, error) {
	out := cldf.ChangesetOutput{
		AddressBook: cldf.NewMemoryAddressBook(),
		DataStore:   datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata](),
	}

	selectors := cfg.ChainSelectors
	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(env.BlockChains.EVMChains()))
	}
	for _, sel := range selectors {
		req := &DeployRequestV2{
			ChainSel: sel,
			deployFn: internal.DeployBalanceReader,
		}
		csOut, err := DeployBalanceReaderV2(env, req)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy BalanceReader to chain selector %d: %w", sel, err)
		}
		if err := out.AddressBook.Merge(csOut.AddressBook); err != nil { //nolint:staticcheck // TODO CRE-400
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for chain selector %d: %w", sel, err)
		}
		if err := out.DataStore.Merge(csOut.DataStore.Seal()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge datastore for chain selector %d: %w", sel, err)
		}
	}

	return out, nil
}

func DeployBalanceReaderV2(env cldf.Environment, req *DeployRequestV2) (cldf.ChangesetOutput, error) {
	req.deployFn = internal.DeployBalanceReader
	return deploy(env, req)
}
