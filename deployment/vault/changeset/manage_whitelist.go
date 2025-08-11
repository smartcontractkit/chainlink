package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var SetWhitelistChangeset = cldf.CreateChangeSet(setWhitelistLogic, setWhitelistPrecondition)

func setWhitelistPrecondition(e cldf.Environment, cfg types.SetWhitelistConfig) error {
	return ValidateSetWhitelistConfig(e, cfg)
}

func setWhitelistLogic(e cldf.Environment, cfg types.SetWhitelistConfig) (cldf.ChangesetOutput, error) {
	lggr := e.Logger

	totalAddresses := 0
	for _, addresses := range cfg.WhitelistByChain {
		totalAddresses += len(addresses)
	}

	lggr.Infow("Setting whitelist state",
		"chains", len(cfg.WhitelistByChain),
		"total_addresses", totalAddresses)

	ds := datastore.NewMemoryDataStore()
	if e.DataStore != nil {
		if err := ds.Merge(e.DataStore); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge existing datastore: %w", err)
		}
	}

	for chainSelector, addresses := range cfg.WhitelistByChain {
		lggr.Infow("Setting whitelist for chain",
			"chain", chainSelector,
			"address_count", len(addresses))

		for _, addr := range addresses {
			lggr.Infow("Whitelist address",
				"chain", chainSelector,
				"address", addr.Address,
				"description", addr.Description)
		}

		whitelistMetadata := types.WhitelistMetadata{
			Addresses: addresses,
		}

		err := ds.ChainMetadata().Upsert(datastore.ChainMetadata{
			ChainSelector: chainSelector,
			Metadata:      whitelistMetadata,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set whitelist chain metadata for chain %d: %w", chainSelector, err)
		}
	}

	lggr.Infow("Whitelist state set successfully")

	return cldf.ChangesetOutput{
		DataStore: ds,
	}, nil
}

// GetWhitelistedAddresses retrieves all whitelisted addresses for given chains from chain metadata
func GetWhitelistedAddresses(e cldf.Environment, chainSelectors []uint64) (map[uint64][]WhitelistEntry, error) {
	whitelist := make(map[uint64][]WhitelistEntry)

	if e.DataStore == nil {
		return nil, errors.New("datastore is nil; whitelist not initialized")
	}

	for _, chainSelector := range chainSelectors {
		whitelistMetadata, err := GetChainWhitelist(e.DataStore, chainSelector)
		if err != nil {
			return nil, err
		}

		var entries []WhitelistEntry
		for _, addr := range whitelistMetadata.Addresses {
			entry := WhitelistEntry{
				Address:   addr.Address,
				Labels:    addr.Labels,
				Qualifier: "whitelist-" + addr.Address,
			}
			entries = append(entries, entry)
		}

		whitelist[chainSelector] = entries
	}

	return whitelist, nil
}

// ValidateWhitelist checks if all addresses in a transfer config are whitelisted
func ValidateWhitelist(e cldf.Environment, cfg types.BatchNativeTransferConfig) ([]types.TransferValidationError, error) {
	var errors []types.TransferValidationError

	chainSelectors := make([]uint64, 0, len(cfg.TransfersByChain))
	for chainSelector := range cfg.TransfersByChain {
		chainSelectors = append(chainSelectors, chainSelector)
	}

	whitelist, err := GetWhitelistedAddresses(e, chainSelectors)
	if err != nil {
		return nil, fmt.Errorf("failed to get whitelist: %w", err)
	}

	for chainSelector, transfers := range cfg.TransfersByChain {
		whitelistedAddrs := make(map[string]bool)
		for _, entry := range whitelist[chainSelector] {
			whitelistedAddrs[entry.Address] = true
		}

		for _, transfer := range transfers {
			if !whitelistedAddrs[transfer.To] {
				errors = append(errors, types.TransferValidationError{
					ChainSelector: chainSelector,
					Address:       transfer.To,
					Error:         "address not in whitelist",
				})
			}
		}
	}

	return errors, nil
}

func GetChainWhitelist(dataStore datastore.DataStore, chainSelector uint64) (*types.WhitelistMetadata, error) {
	chainMetadataKey := datastore.NewChainMetadataKey(chainSelector)
	chainMetadata, err := dataStore.ChainMetadata().Get(chainMetadataKey)
	if err != nil {
		if errors.Is(err, datastore.ErrChainMetadataNotFound) {
			return &types.WhitelistMetadata{Addresses: []types.WhitelistAddress{}}, nil
		}
		return nil, fmt.Errorf("failed to get chain metadata for chain %d: %w", chainSelector, err)
	}

	whitelistMetadata, err := datastore.As[types.WhitelistMetadata](chainMetadata.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to convert chain metadata to whitelist metadata for chain %d: %w", chainSelector, err)
	}

	return &whitelistMetadata, nil
}

func GetChainWhitelistMutable(dataStore datastore.MutableDataStore, chainSelector uint64) (*types.WhitelistMetadata, error) {
	chainMetadataKey := datastore.NewChainMetadataKey(chainSelector)
	chainMetadata, err := dataStore.ChainMetadata().Get(chainMetadataKey)
	if err != nil {
		if errors.Is(err, datastore.ErrChainMetadataNotFound) {
			return &types.WhitelistMetadata{Addresses: []types.WhitelistAddress{}}, nil
		}
		return nil, fmt.Errorf("failed to get chain metadata for chain %d: %w", chainSelector, err)
	}

	whitelistMetadata, err := datastore.As[types.WhitelistMetadata](chainMetadata.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to convert chain metadata to whitelist metadata for chain %d: %w", chainSelector, err)
	}

	return &whitelistMetadata, nil
}

type WhitelistEntry struct {
	Address   string   `json:"address"`
	Labels    []string `json:"labels"`
	Qualifier string   `json:"qualifier"`
}
