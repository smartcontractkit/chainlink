package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// MergeAllDataStores merges all DataStores (after contracts deployments)
func MergeAllDataStores(creEnvironment *cre.Environment, changesetOutputs ...cldf.ChangesetOutput) {
	framework.L.Info().Msg("Merging DataStores (after contracts deployments)...")
	baseDataStore := datastore.NewMemoryDataStore()

	// Merge all subsequent changesets into the base data store
	for i := range changesetOutputs {
		otherDataStore := changesetOutputs[i].DataStore
		mergeErr := baseDataStore.Merge(otherDataStore.Seal())
		if mergeErr != nil {
			panic(errors.Wrap(mergeErr, "DataStores merging failed"))
		}
	}

	mErr := baseDataStore.Merge(creEnvironment.CldfEnvironment.DataStore)
	if mErr != nil {
		panic(errors.Wrap(mErr, "DataStores merging failed"))
	}
	creEnvironment.CldfEnvironment.DataStore = baseDataStore.Seal()
}

func MustGetAddressFromMemoryDataStore(dataStore *datastore.MemoryDataStore, chainSel uint64, contractType string, version *semver.Version, qualifier string) common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		version,
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
	}
	return common.HexToAddress(addrRef.Address)
}

func MightGetAddressFromMemoryDataStore(dataStore *datastore.MemoryDataStore, chainSel uint64, contractType string, version *semver.Version, qualifier string) *common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		version,
		qualifier,
	)

	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		return nil
	}

	return ptr.Ptr(common.HexToAddress(addrRef.Address))
}

func MightGetAddressFromDataStore(dataStore datastore.DataStore, chainSel uint64, contractType string, version *semver.Version, qualifier string) *common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		version,
		qualifier,
	)

	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		return nil
	}
	return ptr.Ptr(common.HexToAddress(addrRef.Address))
}

func MustGetAddressFromDataStore(dataStore datastore.DataStore, chainSel uint64, contractType string, version *semver.Version, qualifier string) string {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		version,
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
	}
	return addrRef.Address
}

func MustGetAddressRefFromDataStore(dataStore datastore.DataStore, chainSel uint64, contractType string, version *semver.Version, qualifier string) datastore.AddressRef {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		version,
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
	}
	return addrRef
}

func NewDataStoreFromExisting(existing datastore.DataStore) (*datastore.MemoryDataStore, error) {
	memoryDatastore := datastore.NewMemoryDataStore()

	mergeErr := memoryDatastore.Merge(existing)
	if mergeErr != nil {
		return nil, fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	return memoryDatastore, nil
}

func ConfigureDataFeedsCache(testLogger zerolog.Logger, input *cre.ConfigureDataFeedsCacheInput) (*cre.ConfigureDataFeedsCacheOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if input.Out != nil && input.Out.UseCache {
		return input.Out, nil
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	if input.AdminAddress != (common.Address{}) {
		setAdminConfig := df_changeset_types.SetFeedAdminConfig{
			ChainSelector: input.ChainSelector,
			CacheAddress:  input.DataFeedsCacheAddress,
			AdminAddress:  input.AdminAddress,
			IsAdmin:       true,
		}
		_, setAdminErr := commonchangeset.RunChangeset(df_changeset.SetFeedAdminChangeset, *input.CldEnv, setAdminConfig)
		if setAdminErr != nil {
			return nil, errors.Wrap(setAdminErr, "failed to set feed admin")
		}
	}

	metadatas := []data_feeds_cache.DataFeedsCacheWorkflowMetadata{}
	for idx := range input.AllowedWorkflowNames {
		metadatas = append(metadatas, data_feeds_cache.DataFeedsCacheWorkflowMetadata{
			AllowedWorkflowName:  df_changeset.HashedWorkflowName(input.AllowedWorkflowNames[idx]),
			AllowedSender:        input.AllowedSenders[idx],
			AllowedWorkflowOwner: input.AllowedWorkflowOwners[idx],
		})
	}

	feeIDs := []string{}
	for _, feedID := range input.FeedIDs {
		feeIDs = append(feeIDs, feedID[:32])
	}

	_, setFeedConfigErr := commonchangeset.RunChangeset(df_changeset.SetFeedConfigChangeset, *input.CldEnv, df_changeset_types.SetFeedDecimalConfig{
		ChainSelector:    input.ChainSelector,
		CacheAddress:     input.DataFeedsCacheAddress,
		DataIDs:          feeIDs,
		Descriptions:     input.Descriptions,
		WorkflowMetadata: metadatas,
	})

	if setFeedConfigErr != nil {
		return nil, errors.Wrap(setFeedConfigErr, "failed to set feed config")
	}

	out := &cre.ConfigureDataFeedsCacheOutput{
		DataFeedsCacheAddress: input.DataFeedsCacheAddress,
		FeedIDs:               input.FeedIDs,
		AllowedSenders:        input.AllowedSenders,
		AllowedWorkflowOwners: input.AllowedWorkflowOwners,
		AllowedWorkflowNames:  input.AllowedWorkflowNames,
	}

	if input.AdminAddress != (common.Address{}) {
		out.AdminAddress = input.AdminAddress
	}

	input.Out = out

	return out, nil
}

func DeployDataFeedsCacheContract(testLogger zerolog.Logger, chainSelector uint64, creEnvironment *cre.Environment) (common.Address, error) {
	testLogger.Info().Msg("Deploying Data Feeds Cache contract...")
	deployDfConfig := df_changeset_types.DeployConfig{
		ChainsToDeploy: []uint64{chainSelector},
		Labels:         []string{"data-feeds"}, // label required by the changeset
	}

	dfOutput, dfErr := commonchangeset.RunChangeset(df_changeset.DeployCacheChangeset, *creEnvironment.CldfEnvironment, deployDfConfig)
	if dfErr != nil {
		return common.Address{}, errors.Wrapf(dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)
	}
	testLogger.Info().Msgf("Data Feeds Cache contract deployed to %d", chainSelector)

	memoryDatastore, mErr := NewDataStoreFromExisting(creEnvironment.CldfEnvironment.DataStore)
	if mErr != nil {
		return common.Address{}, fmt.Errorf("failed to create memory datastore: %w", mErr)
	}
	if dfOutput.DataStore != nil {
		err := memoryDatastore.Merge(dfOutput.DataStore.Seal())
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to merge updated datastore: %w", err)
		}
		creEnvironment.CldfEnvironment.DataStore = memoryDatastore.Seal()
	}

	dataFeedsCacheAddress := MustGetAddressFromMemoryDataStore(memoryDatastore, chainSelector, "DataFeedsCache", semver.MustParse("1.0.0"), "")
	testLogger.Info().Msgf("Data Feeds Cache contract found on chain %d at address %s", chainSelector, dataFeedsCacheAddress)

	return dataFeedsCacheAddress, nil
}

func DeployReadBalancesContract(testLogger zerolog.Logger, chainSelector uint64, creEnvironment *cre.Environment) (common.Address, error) {
	testLogger.Info().Msg("Deploying Read Balances contract...")
	deployReadBalanceRequest := &keystone_changeset.DeployRequestV2{ChainSel: chainSelector}
	rbOutput, rbErr := keystone_changeset.DeployBalanceReaderV2(*creEnvironment.CldfEnvironment, deployReadBalanceRequest)
	if rbErr != nil {
		return common.Address{}, errors.Wrap(rbErr, "failed to deploy Read Balances contract")
	}
	testLogger.Info().Msgf("Read Balances contract deployed to %d", chainSelector)

	memoryDatastore, mErr := NewDataStoreFromExisting(creEnvironment.CldfEnvironment.DataStore)
	if mErr != nil {
		return common.Address{}, fmt.Errorf("failed to create memory datastore: %w", mErr)
	}

	if rbOutput.DataStore != nil {
		err := memoryDatastore.Merge(rbOutput.DataStore.Seal())
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to merge updated datastore: %w", err)
		}
		creEnvironment.CldfEnvironment.DataStore = memoryDatastore.Seal()
	}

	readBalancesAddress := MustGetAddressFromMemoryDataStore(memoryDatastore, chainSelector, "BalanceReader", semver.MustParse("1.0.0"), "")
	testLogger.Info().Msgf("Read Balances contract found on chain %d at address %s", chainSelector, readBalancesAddress)

	return readBalancesAddress, nil
}
