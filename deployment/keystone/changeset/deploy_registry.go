package changeset

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[uint64] = DeployCapabilityRegistry

// Depreciated: use DeployCapabilityRegistryV2 instead
func DeployCapabilityRegistry(env cldf.Environment, registrySelector uint64) (cldf.ChangesetOutput, error) {
	return DeployCapabilityRegistryV2(env, &DeployRequestV2{
		ChainSel: registrySelector,
	})
}

func DeployCapabilityRegistryV2(env cldf.Environment, req *DeployRequestV2) (cldf.ChangesetOutput, error) {
	req.deployFn = kslib.DeployCapabilitiesRegistry
	return deploy(env, req)
}

// DeployRequestV2 is a request to deploy the given deployFn to the given chain
type DeployRequestV2 = struct {
	ChainSel  uint64
	Qualifier string
	Labels    *datastore.LabelSet

	deployFn func(ctx context.Context, chain cldf_evm.Chain, ab cldf.AddressBook) (*kslib.DeployResponse, error)
}

func deploy(env cldf.Environment, req *DeployRequestV2) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.BlockChains.EVMChains()[req.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := cldf.NewMemoryAddressBook()
	resp, err := req.deployFn(env.GetContext(), chain, ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", resp.Tv.String(), chain.Selector, resp.Address.String())

	ds := datastore.NewMemoryDataStore()
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
