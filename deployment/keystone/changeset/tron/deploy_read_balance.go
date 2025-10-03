package tron

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

var DeployReadBalanceChangeset = cldf.CreateChangeSet(deployReadBalanceLogic, deployReadBalancePrecondition)

// DeployReadBalanceChangeset deploys the ReadBalances contract to the specified chains
// Returns a new addressbook with deployed ReadBalances contracts
func deployReadBalanceLogic(env cldf.Environment, c types.DeployTronConfig) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()
	dataStore := datastore.NewMemoryDataStore()

	for _, chainSelector := range c.ChainsToDeploy {
		chain := env.BlockChains.TronChains()[chainSelector]
		readBalanceResponse, err := DeployReadBalance(chain, c.DeployOptions, c.Labels)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ReadBalances: %w", err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", readBalanceResponse.Tv.String(), chain.Selector, readBalanceResponse.Address.String())

		addr := readBalanceResponse.Address.String()
		isEvm, _ := chain_selectors.IsEvm(chainSelector)
		if isEvm {
			addr = readBalanceResponse.Address.EthAddress().Hex()
		}

		err = ab.Save(chain.Selector, addr, readBalanceResponse.Tv)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ReadBalances: %w", err)
		}

		if err = dataStore.Addresses().Add(
			datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       addr,
				Type:          datastore.ContractType(readBalanceResponse.Tv.Type),
				Version:       semver.MustParse("1.0.0"),
				Qualifier:     c.Qualifier,
				Labels:        datastore.NewLabelSet(readBalanceResponse.Tv.Labels.List()...),
			},
		); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore: %w", err)
		}
	}

	return cldf.ChangesetOutput{AddressBook: ab, DataStore: dataStore}, nil
}

func deployReadBalancePrecondition(env cldf.Environment, c types.DeployTronConfig) error {
	for _, chainSelector := range c.ChainsToDeploy {
		_, ok := env.BlockChains.TronChains()[chainSelector]
		if !ok {
			return errors.New("chain not found in environment")
		}
	}

	return nil
}
