package tron

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var _ cldf.ChangeSetV2[*DeployForwarderRequest] = DeployForwarder{}

type DeployForwarder struct{}

func (cs DeployForwarder) VerifyPreconditions(env cldf.Environment, req *DeployForwarderRequest) error {
	for _, chainSelector := range req.ChainSelectors {
		_, ok := env.BlockChains.TronChains()[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
	}

	return nil
}

type DeployForwarderRequest struct {
	ChainSelectors []uint64
	Labels         []string
	Qualifier      string
	DeployOptions  *cldf_tron.DeployOptions
}

func (cs DeployForwarder) Apply(env cldf.Environment, req *DeployForwarderRequest) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()
	dataStore := datastore.NewMemoryDataStore()

	for _, chainSelector := range req.ChainSelectors {
		chain := env.BlockChains.TronChains()[chainSelector]
		forwarderResponse, err := DeployKeystoneForwarder(chain, req.DeployOptions, req.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", forwarderResponse.Tv.String(), chain.Selector, forwarderResponse.Address.String())

		err = ab.Save(chain.Selector, forwarderResponse.Address.String(), forwarderResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save KeystoneForwarder: %w", err)
		}

		if err = dataStore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       forwarderResponse.Address.String(),
				Type:          ForwarderContract,
				Version:       semver.MustParse("1.0.0"),
				Qualifier:     req.Qualifier,
				Labels:        datastore.NewLabelSet(forwarderResponse.Tv.Labels.List()...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	return cldf.ChangesetOutput{AddressBook: ab, DataStore: dataStore}, nil
}
