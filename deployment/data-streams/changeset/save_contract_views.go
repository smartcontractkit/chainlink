package changeset

import (
	"errors"
	"fmt"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	dsstate "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"
)

// SaveContractViews saves the contract views to the datastore.
var SaveContractViews = cldf.CreateChangeSet(saveViewsLogic, saveViewsPrecondition)

type SaveContractViewsConfig struct {
	Chains []uint64
}

func (cfg SaveContractViewsConfig) Validate() error {
	if len(cfg.Chains) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func saveViewsPrecondition(_ cldf.Environment, cc SaveContractViewsConfig) error {
	return cc.Validate()
}

func saveViewsLogic(e cldf.Environment, cfg SaveContractViewsConfig) (cldf.ChangesetOutput, error) {
	updatedDataStore := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()
	records, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return cldf.ChangesetOutput{}, errors.New("failed to fetch addresses")
	}

	addressesByChain := utils.AddressRefsToAddressByChain(records)

	envDatastore, err := ds.FromDefault[metadata.SerializedContractMetadata, ds.DefaultMetadata](e.DataStore)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert datastore: %w", err)
	}

	// Operation to Generate + Save Contract Views PER Chain.
	for _, chainSelector := range cfg.Chains {
		chainAddresses, ok := addressesByChain[chainSelector]
		if !ok {
			continue
		}
		chain := e.Chains[chainSelector]

		family, err := chain_selectors.GetSelectorFamily(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get selector family: %w", err)
		}

		switch family {
		case chain_selectors.FamilyEVM:
			cmStore := envDatastore.ContractMetadata()
			chainState, err := dsstate.LoadChainState(e.Logger, chain, chainAddresses, cmStore)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to load chain state: %w", err)
			}

			chainView, _ := chainState.GenerateView(e.GetContext())

			err = saveEvmViewsToDataStore(chainView, cmStore, updatedDataStore, chainSelector)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save views to datastore: %w", err)
			}

		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("unsupported chain selector: %d", chainSelector)
		}
	}

	defaultDs, err := ds.ToDefault(updatedDataStore.Seal())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}

	return cldf.ChangesetOutput{DataStore: defaultDs}, nil
}

func saveEvmViewsToDataStore(chainView view.EVMChainView,
	cmStore ds.ContractMetadataStore[metadata.SerializedContractMetadata],
	updatedDataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	chainSelector uint64) error {
	for address, contractView := range chainView.Configurator {
		err := saveContractViewToDataStore(cmStore, updatedDataStore, chainSelector, address, &contractView)
		if err != nil {
			return fmt.Errorf("failed to save metadata to datastore: %w", err)
		}
	}

	for address, contractView := range chainView.Verifier {
		err := saveContractViewToDataStore(cmStore, updatedDataStore, chainSelector, address, &contractView)
		if err != nil {
			return fmt.Errorf("failed to save metadata to datastore: %w", err)
		}
	}

	for address, contractView := range chainView.FeeManager {
		err := saveContractViewToDataStore(cmStore, updatedDataStore, chainSelector, address, &contractView)
		if err != nil {
			return fmt.Errorf("failed to save metadata to datastore: %w", err)
		}
	}

	for address, contractView := range chainView.RewardManager {
		err := saveContractViewToDataStore(cmStore, updatedDataStore, chainSelector, address, &contractView)
		if err != nil {
			return fmt.Errorf("failed to save metadata to datastore: %w", err)
		}
	}

	for address, contractView := range chainView.VerifierProxy {
		err := saveContractViewToDataStore(cmStore, updatedDataStore, chainSelector, address, &contractView)
		if err != nil {
			return fmt.Errorf("failed to save metadata to datastore: %w", err)
		}
	}

	for address, contractView := range chainView.ChannelConfigStore {
		err := saveContractViewToDataStore(cmStore, updatedDataStore, chainSelector, address, &contractView)
		if err != nil {
			return fmt.Errorf("failed to save metadata to datastore: %w", err)
		}
	}

	return nil
}

func saveContractViewToDataStore[T interfaces.ContractView](
	existingMdStore ds.ContractMetadataStore[metadata.SerializedContractMetadata],
	mutableDatastore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	chainSelector uint64,
	address view.Address,
	view T,
) error {
	cm, err := existingMdStore.Get(ds.NewContractMetadataKey(chainSelector, address))
	if err != nil {
		return fmt.Errorf("failed to get contract metadata: %w", err)
	}
	existingMd, err := metadata.DeserializeMetadata[T](cm.Metadata)
	if err != nil {
		return fmt.Errorf("failed to convert contract metadata: %w", err)
	}

	contractMetadata := metadata.GenericContractMetadata[T]{
		Metadata: existingMd.Metadata,
		View:     view,
	}

	serialized, err := metadata.NewSerializedContractMetadata(contractMetadata)
	if err != nil {
		return fmt.Errorf("failed to serialize contract metadata: %w", err)
	}

	if err = mutableDatastore.ContractMetadata().Upsert(
		ds.ContractMetadata[metadata.SerializedContractMetadata]{
			ChainSelector: chainSelector,
			Address:       address,
			Metadata:      *serialized,
		},
	); err != nil {
		return fmt.Errorf("failed to upsert contract metadata: %w", err)
	}

	return nil
}
