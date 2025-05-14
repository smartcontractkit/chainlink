package reward_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var DeployRewardManagerChangeset = cldf.CreateChangeSet(deployRewardManagerLogic, deployRewardManagerPrecondition)

type DeployRewardManager struct {
	LinkTokenAddress common.Address
}

type DeployRewardManagerConfig struct {
	ChainsToDeploy map[uint64]DeployRewardManager
	Ownership      types.OwnershipSettings
}

func (cc DeployRewardManagerConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cc.Ownership
}

func (cc DeployRewardManagerConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for chain := range cc.ChainsToDeploy {
		if err := cldf.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func deployRewardManagerLogic(e cldf.Environment, cc DeployRewardManagerConfig) (cldf.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()
	err := deployRewardManager(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy RewardManager", "err", err)
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

func deployRewardManagerPrecondition(_ cldf.Environment, cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}

	return nil
}

func deployRewardManager(e cldf.Environment,
	dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}

	for chainSel, chainCfg := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}

		res, err := changeset.DeployContract(e, dataStore, chain, RewardManagerDeployFn(chainCfg.LinkTokenAddress), nil)
		if err != nil {
			return err
		}
		contractMetadata := metadata.GenericContractMetadata[v0_5.RewardManagerView]{
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

func RewardManagerDeployFn(linkAddress common.Address) changeset.ContractDeployFn[*rewardManager.RewardManager] {
	return func(chain cldf.Chain) *changeset.ContractDeployment[*rewardManager.RewardManager] {
		ccsAddr, ccsTx, ccs, err := rewardManager.DeployRewardManager(
			chain.DeployerKey,
			chain.Client,
			linkAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*rewardManager.RewardManager]{
				Err: err,
			}
		}
		bn, err := chain.Confirm(ccsTx)
		if err != nil {
			return &changeset.ContractDeployment[*rewardManager.RewardManager]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*rewardManager.RewardManager]{
			Address:  ccsAddr,
			Block:    bn,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       cldf.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
