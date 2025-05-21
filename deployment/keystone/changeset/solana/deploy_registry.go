package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/solana"
)

type DeployRequest = struct {
	ChainSel    uint64
	Qualifier   string
	Labels      *datastore.LabelSet
	BuildConfig *BuildSolanaConfig

	deployFn func(ctx context.Context, chain cldf.SolChain, ab cldf.AddressBook) (*kslib.DeployResponse, error)
}

func deploy(env cldf.Environment, req *DeployRequest) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.SolChains[req.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, errors.New("chain not found in environment")
	}

	ab := cldf.NewMemoryAddressBook()
	resp, err := req.deployFn(env.GetContext(), chain, ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", resp.Tv.String(), chain.Selector, resp.Address.String())

	ds := datastore.NewMemoryDataStore[
		datastore.DefaultMetadata,
		datastore.DefaultMetadata,
	]()

	r := datastore.AddressRef{
		ChainSelector: req.ChainSel,
		Address:       resp.Address.String(),
		Type:          datastore.ContractType(resp.Tv.Type),
		Version:       &resp.Tv.Version,
		Qualifier:     req.Qualifier,
		Labels:        datastore.NewLabelSet(),
	}

	if req.Labels != nil {
		r.Labels = *req.Labels
	}

	// add labels from the response
	for _, l := range resp.Tv.Labels.List() {
		r.Labels.Add(l)
	}

	if err = ds.Addresses().Add(r); err != nil {
		return cldf.ChangesetOutput{DataStore: ds},
			fmt.Errorf("failed to save address ref in datastore: %w", err)
	}

	return cldf.ChangesetOutput{AddressBook: ab, DataStore: ds}, nil
}
