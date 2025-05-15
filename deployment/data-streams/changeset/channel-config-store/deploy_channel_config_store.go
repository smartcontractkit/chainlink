package channel_config_store

import (
	"errors"
	"fmt"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

// DeployChannelConfigStoreChangeset deploys ChannelConfigStore to the chains specified in the config.
var DeployChannelConfigStoreChangeset = cldf.CreateChangeSet(deployChannelConfigStoreLogic, deployChannelConfigStorePrecondition)

func (cc DeployChannelConfigStoreConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cc.Ownership
}

type DeployChannelConfigStoreConfig struct {
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy []uint64
	Ownership      types.OwnershipSettings
}

func (cc DeployChannelConfigStoreConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for _, chain := range cc.ChainsToDeploy {
		if err := cldf.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func deployChannelConfigStoreLogic(e cldf.Environment, cc DeployChannelConfigStoreConfig) (cldf.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()
	err := deploy(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy ChannelConfigStore", "err", err)
		return cldf.ChangesetOutput{}, cldf.MaybeDataErr(err)
	}

	records, err := dataStore.Addresses().Fetch()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch addresses: %w", err)
	}
	proposals, err := mcmsutil.GetTransferOwnershipProposals(e, cc, records)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership to MCMS: %w", err)
	}

	sealedDS, err := ds.ToDefault(dataStore.Seal())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}
	return cldf.ChangesetOutput{
		DataStore:             sealedDS,
		MCMSTimelockProposals: proposals,
	}, nil
}

func deployChannelConfigStorePrecondition(_ cldf.Environment, cc DeployChannelConfigStoreConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployChannelConfigStoreConfig: %w", err)
	}

	return nil
}

func deploy(e cldf.Environment, dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata], cc DeployChannelConfigStoreConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployChannelConfigStoreConfig: %w", err)
	}
	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		res, err := changeset.DeployContract(e, dataStore, chain, channelConfigStoreDeployFn(), nil)
		if err != nil {
			return err
		}

		contractMetadata := metadata.GenericContractMetadata[v0_5.ChannelConfigStoreView]{
			Metadata: metadata.CommonContractMetadata{
				DeployBlock: res.Block,
			},
		}

		serialized, err := metadata.NewSerializedContractMetadata(contractMetadata)
		if err != nil {
			return fmt.Errorf("failed to serialize contract metadata: %w", err)
		}

		if err = dataStore.ContractMetadata().Upsert(
			ds.ContractMetadata[metadata.SerializedContractMetadata]{
				ChainSelector: chain.Selector,
				Address:       res.Address.String(),
				Metadata:      *serialized,
			},
		); err != nil {
			return fmt.Errorf("failed to upser contract metadata: %w", err)
		}
	}

	return nil
}

// channelConfigStoreDeployFn returns a function that deploys a ChannelConfigStore contract.
func channelConfigStoreDeployFn() changeset.ContractDeployFn[*channel_config_store.ChannelConfigStore] {
	return func(chain cldf.Chain) *changeset.ContractDeployment[*channel_config_store.ChannelConfigStore] {
		ccsAddr, ccsTx, ccs, err := channel_config_store.DeployChannelConfigStore(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return &changeset.ContractDeployment[*channel_config_store.ChannelConfigStore]{
				Err: err,
			}
		}
		bn, err := chain.Confirm(ccsTx)
		if err != nil {
			return &changeset.ContractDeployment[*channel_config_store.ChannelConfigStore]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*channel_config_store.ChannelConfigStore]{
			Address:  ccsAddr,
			Block:    bn,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       cldf.NewTypeAndVersion(types.ChannelConfigStore, deployment.Version1_0_0),
			Err:      nil,
		}
	}
}
