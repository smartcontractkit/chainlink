package contracts

import (
	"fmt"
	"strings"

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

func FindAddressesForChain(addressBook cldf.AddressBook, chainSelector uint64, contractName string) (common.Address, cldf.TypeAndVersion, error) {
	addresses, err := addressBook.AddressesForChain(chainSelector)
	if err != nil {
		return common.Address{}, cldf.TypeAndVersion{}, errors.Wrap(err, "failed to get addresses for chain")
	}

	for addrStr, tv := range addresses {
		if !strings.Contains(tv.String(), contractName) {
			continue
		}

		return common.HexToAddress(addrStr), tv, nil
	}

	return common.Address{}, cldf.TypeAndVersion{}, fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector)
}

// TODO: CRE-742 use datastore
func MustFindAddressesForChain(addressBook cldf.AddressBook, chainSelector uint64, contractName string) common.Address {
	addr, _, err := FindAddressesForChain(addressBook, chainSelector, contractName)
	if err != nil {
		panic(fmt.Errorf("failed to find %s address in the address book for chain %d", contractName, chainSelector))
	}
	return addr
}

// MergeAllDataStores merges all DataStores (after contracts deployments)
func MergeAllDataStores(creEnvironment *cre.Environment, changesetOutputs ...cldf.ChangesetOutput) {
	framework.L.Info().Msg("Merging DataStores (after contracts deployments)...")
	minChangesetsCap := 2
	if len(changesetOutputs) < minChangesetsCap {
		panic(fmt.Errorf("DataStores merging failed: at least %d changesets required", minChangesetsCap))
	}

	// Start with the first changeset's data store
	baseDataStore := changesetOutputs[0].DataStore

	// Merge all subsequent changesets into the base data store
	for i := 1; i < len(changesetOutputs); i++ {
		otherDataStore := changesetOutputs[i].DataStore
		mergeErr := baseDataStore.Merge(otherDataStore.Seal())
		if mergeErr != nil {
			panic(errors.Wrap(mergeErr, "DataStores merging failed"))
		}
	}

	creEnvironment.CldfEnvironment.DataStore = baseDataStore.Seal()
}

func MustGetAddressFromMemoryDataStore(dataStore *datastore.MemoryDataStore, chainSel uint64, contractType string, version string, qualifier string) common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
	}
	return common.HexToAddress(addrRef.Address)
}

func MightGetAddressFromMemoryDataStore(dataStore *datastore.MemoryDataStore, chainSel uint64, contractType string, version string, qualifier string) *common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)

	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		return nil
	}

	return ptr.Ptr(common.HexToAddress(addrRef.Address))
}

func MightGetAddressFromDataStore(dataStore datastore.DataStore, chainSel uint64, contractType string, version string, qualifier string) *common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)

	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		return nil
	}
	return ptr.Ptr(common.HexToAddress(addrRef.Address))
}

func MustGetAddressFromDataStore(dataStore datastore.DataStore, chainSel uint64, contractType string, version string, qualifier string) string {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
	}
	return addrRef.Address
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

func DeployDataFeedsCacheContract(testLogger zerolog.Logger, chainSelector uint64, creEnvironment *cre.Environment) (common.Address, cldf.ChangesetOutput, error) {
	testLogger.Info().Msg("Deploying Data Feeds Cache contract...")
	deployDfConfig := df_changeset_types.DeployConfig{
		ChainsToDeploy: []uint64{chainSelector},
		Labels:         []string{"data-feeds"}, // label required by the changeset
	}

	dfOutput, dfErr := commonchangeset.RunChangeset(df_changeset.DeployCacheChangeset, *creEnvironment.CldfEnvironment, deployDfConfig)
	if dfErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrapf(dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)
	}

	mergeErr := creEnvironment.CldfEnvironment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(mergeErr, "failed to merge address book of Data Feeds Cache contract")
	}
	testLogger.Info().Msgf("Data Feeds Cache contract deployed to %d", chainSelector)

	dataFeedsCacheAddress, _, dataFeedsCacheErr := FindAddressesForChain(
		creEnvironment.CldfEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelector,
		df_changeset.DataFeedsCache.String(),
	)
	if dataFeedsCacheErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrapf(dataFeedsCacheErr, "failed to find Data Feeds Cache contract address on chain %d", chainSelector)
	}
	testLogger.Info().Msgf("Data Feeds Cache contract found on chain %d at address %s", chainSelector, dataFeedsCacheAddress)

	return dataFeedsCacheAddress, dfOutput, nil
}

func DeployReadBalancesContract(testLogger zerolog.Logger, chainSelector uint64, creEnvironment *cre.Environment) (common.Address, cldf.ChangesetOutput, error) {
	testLogger.Info().Msg("Deploying Read Balances contract...")
	deployReadBalanceRequest := &keystone_changeset.DeployRequestV2{ChainSel: chainSelector}
	rbOutput, rbErr := keystone_changeset.DeployBalanceReaderV2(*creEnvironment.CldfEnvironment, deployReadBalanceRequest)
	if rbErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(rbErr, "failed to deploy Read Balances contract")
	}

	mergeErr2 := creEnvironment.CldfEnvironment.ExistingAddresses.Merge(rbOutput.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr2 != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(mergeErr2, "failed to merge address book of Read Balances contract")
	}
	testLogger.Info().Msgf("Read Balances contract deployed to %d", chainSelector)

	readBalancesAddress, _, readContractErr := FindAddressesForChain(
		creEnvironment.CldfEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelector,
		keystone_changeset.BalanceReader.String(),
	)
	if readContractErr != nil {
		return common.Address{}, cldf.ChangesetOutput{}, errors.Wrap(readContractErr, "failed to find Read Balances contract address")
	}
	testLogger.Info().Msgf("Read Balances contract found on chain %d at address %s", chainSelector, readBalancesAddress)

	return readBalancesAddress, rbOutput, nil
}
