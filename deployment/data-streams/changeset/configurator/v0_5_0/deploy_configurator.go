package v0_5_0

import (
	"errors"
	"fmt"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"
)

var DeployConfiguratorChangeset = cldf.CreateChangeSet(deployConfiguratorLogic, deployConfiguratorPrecondition)

type DeployConfiguratorConfig struct {
	ChainsToDeploy []uint64
	Ownership      types.OwnershipSettings
}

func (cc DeployConfiguratorConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cc.Ownership
}

func (cc DeployConfiguratorConfig) Validate() error {
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

func deployConfiguratorLogic(e cldf.Environment, cc DeployConfiguratorConfig) (cldf.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[
		metadata.SerializedContractMetadata,
		ds.DefaultMetadata,
	]()

	err := deploy(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy Configurator", "err", err)
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

func deployConfiguratorPrecondition(_ cldf.Environment, cc DeployConfiguratorConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}

	return nil
}

func deploy(e cldf.Environment, dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata], cc DeployConfiguratorConfig) error {
	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}

		res, err := changeset.DeployContract(e, dataStore, chain, DeployFn(), nil)
		if err != nil {
			return fmt.Errorf("failed to deploy configurator: %w", err)
		}

		contractMetadata := metadata.GenericContractMetadata[v0_5.ConfiguratorView]{
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

func DeployFn() changeset.ContractDeployFn[*configurator.Configurator] {
	return func(chain cldf.Chain) *changeset.ContractDeployment[*configurator.Configurator] {
		ccsAddr, ccsTx, ccs, err := configurator.DeployConfigurator(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return &changeset.ContractDeployment[*configurator.Configurator]{
				Err: err,
			}
		}

		bn, err := chain.Confirm(ccsTx)
		if err != nil {
			return &changeset.ContractDeployment[*configurator.Configurator]{
				Err: err,
			}
		}

		return &changeset.ContractDeployment[*configurator.Configurator]{
			Address:  ccsAddr,
			Block:    bn,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       cldf.NewTypeAndVersion(types.Configurator, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
