package mcms

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
)

var DeployAndTransferMCMSChangeset = cldf.CreateChangeSet(deployAndTransferMcmsLogic, deployAndTransferMcmsPrecondition)

type DeployMCMSConfig struct {
	ChainsToDeploy []uint64
	Ownership      types.OwnershipSettings
	Config         commontypes.MCMSWithTimelockConfigV2
}

func (cc DeployMCMSConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cc.Ownership
}

func deployAndTransferMcmsLogic(e cldf.Environment, cc DeployMCMSConfig) (cldf.ChangesetOutput, error) {
	cfgByChain := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, chain := range cc.ChainsToDeploy {
		cfgByChain[chain] = cc.Config
	}

	mcmsOut, err := commonChangesets.DeployMCMSWithTimelockV2(e, cfgByChain)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS: %w", err)
	}

	// CallProxy has no owner, RBACTimelock has an "admin" setting in place of owner
	transferContracts := []cldf.ContractType{
		commontypes.ProposerManyChainMultisig,
		commontypes.BypasserManyChainMultisig,
		commontypes.CancellerManyChainMultisig,
	}

	// DeployMCMSWithTimelockV2 currently does not use the DataStore
	ds, err := utils.AddressBookToNewDataStore[metadata.SerializedContractMetadata, datastore.DefaultMetadata](mcmsOut.AddressBook) //nolint:staticcheck // won't migrate now
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to address book: %w", err)
	}

	var transferAddresses []datastore.AddressRef
	for _, contractType := range transferContracts {
		addrs := ds.Addresses().Filter(datastore.AddressRefByType(datastore.ContractType(contractType)))
		transferAddresses = append(transferAddresses, addrs...)
	}

	// environment needs the timelock MCMS address to propose the transfer - but it's excluded from the transfer itself
	requiredAddrs := datastore.NewMemoryDataStore[metadata.SerializedContractMetadata, datastore.DefaultMetadata]()
	addrs := ds.Addresses().Filter(datastore.AddressRefByType(datastore.ContractType(commontypes.RBACTimelock)))
	for _, addr := range addrs {
		err := requiredAddrs.Addresses().Add(addr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add address to data store: %w", err)
		}
	}
	requiredAddrsDs, err := datastore.ToDefault(requiredAddrs.Seal())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}
	e.DataStore = requiredAddrsDs.Seal()

	proposals, err := mcmsutil.GetTransferOwnershipProposals(e, cc, transferAddresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership to MCMS: %w", err)
	}

	sealedDS, err := datastore.ToDefault(ds.Seal())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           mcmsOut.AddressBook, //nolint:staticcheck // won't migrate now - kept for backwards compatibility until AddressBook is removed
		DataStore:             sealedDS,
		MCMSTimelockProposals: proposals}, nil
}

func deployAndTransferMcmsPrecondition(_ cldf.Environment, cc DeployMCMSConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployMCMSConfig: %w", err)
	}
	return nil
}

func (cc DeployMCMSConfig) Validate() error {
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
