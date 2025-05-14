package fee_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var DeployFeeManagerChangeset = cldf.CreateChangeSet(deployFeeManagerLogic, deployFeeManagerPrecondition)

type DeployFeeManager struct {
	LinkTokenAddress     common.Address
	NativeTokenAddress   common.Address
	VerifierProxyAddress common.Address
	RewardManagerAddress common.Address
}

type DeployFeeManagerConfig struct {
	ChainsToDeploy map[uint64]DeployFeeManager
	Ownership      types.OwnershipSettings
}

func (cc DeployFeeManagerConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cc.Ownership
}

func (cc DeployFeeManagerConfig) Validate() error {
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

func deployFeeManagerLogic(e cldf.Environment, cc DeployFeeManagerConfig) (cldf.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()
	err := deployFeeManager(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy FeeManager", "err", err)
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

func deployFeeManagerPrecondition(_ cldf.Environment, cc DeployFeeManagerConfig) error {
	return cc.Validate()
}

func deployFeeManager(e cldf.Environment,
	dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	cc DeployFeeManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployFeeManagerConfig: %w", err)
	}

	for chainSel, chainCfg := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}

		res, err := changeset.DeployContract(e, dataStore, chain, FeeManagerDeployFn(chainCfg), nil)
		if err != nil {
			return fmt.Errorf("failed to deploy FeeManager: %w", err)
		}

		contractMetadata := metadata.GenericContractMetadata[v0_5.FeeManagerView]{
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

// FeeManagerDeployFn returns a function that deploys a FeeManager contract.
func FeeManagerDeployFn(cfg DeployFeeManager) changeset.ContractDeployFn[*fee_manager_v0_5_0.FeeManager] {
	return func(chain cldf.Chain) *changeset.ContractDeployment[*fee_manager_v0_5_0.FeeManager] {
		ccsAddr, ccsTx, ccs, err := fee_manager_v0_5_0.DeployFeeManager(
			chain.DeployerKey,
			chain.Client,
			cfg.LinkTokenAddress,
			cfg.NativeTokenAddress,
			cfg.VerifierProxyAddress,
			cfg.RewardManagerAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*fee_manager_v0_5_0.FeeManager]{
				Err: err,
			}
		}
		bn, err := chain.Confirm(ccsTx)
		if err != nil {
			return &changeset.ContractDeployment[*fee_manager_v0_5_0.FeeManager]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*fee_manager_v0_5_0.FeeManager]{
			Address:  ccsAddr,
			Block:    bn,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       cldf.NewTypeAndVersion(types.FeeManager, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
