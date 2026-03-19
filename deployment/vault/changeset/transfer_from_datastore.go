package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

// TransferFromDataStoreInput is the input for the transfer_mcms_ownership_to_timelock pipeline
// when using the catalog-backed changeset. Addresses are read from Environment.DataStore
// (Catalog at pipeline runtime), so deploy_timelock and transfer can run in the same PR.
// Use this for staging_testnet (and later prod_mainnet, prod_oev_mainnet) when running
// both steps in one pipeline.
type TransferFromDataStoreInput struct {
	ChainSelectors     []uint64 `json:"chainSelectors" yaml:"chainSelectors"`
	TimelockIdentifier string   `json:"timelockIdentifier" yaml:"timelockIdentifier"`
}

// TransferMCMSOwnershipFromDataStore runs TransferToMCMSWithTimelockV2 by building the config
// from e.DataStore (the runtime/Catalog-backed datastore). Use instead of the config resolver
// when deploy and transfer run in the same pipeline so the transfer step sees deploy output.
func TransferMCMSOwnershipFromDataStore(e cldf.Environment, input TransferFromDataStoreInput) (cldf.ChangesetOutput, error) {
	if e.DataStore == nil {
		return cldf.ChangesetOutput{}, errors.New("datastore is nil (required for transfer from catalog)")
	}
	if len(input.ChainSelectors) == 0 {
		return cldf.ChangesetOutput{}, errors.New("at least one chain selector is required")
	}

	qualifier := input.TimelockIdentifier
	contractsByChain := make(map[uint64][]common.Address, len(input.ChainSelectors))

	typesOrder := []datastore.ContractType{
		datastore.ContractType(commontypes.BypasserManyChainMultisig.String()),
		datastore.ContractType(commontypes.CancellerManyChainMultisig.String()),
		datastore.ContractType(commontypes.ProposerManyChainMultisig.String()),
	}

	for _, chainSel := range input.ChainSelectors {
		var addrs []common.Address
		for _, ct := range typesOrder {
			filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{
				datastore.AddressRefByChainSelector(chainSel),
				datastore.AddressRefByType(ct),
			}
			if qualifier != "" {
				filters = append(filters, datastore.AddressRefByQualifier(qualifier))
			}
			records := e.DataStore.Addresses().Filter(filters...)
			if len(records) == 0 {
				return cldf.ChangesetOutput{}, fmt.Errorf("no %s found for chain %d (qualifier %q) in datastore", ct, chainSel, qualifier)
			}
			addrs = append(addrs, common.HexToAddress(records[0].Address))
		}
		contractsByChain[chainSel] = addrs
	}

	mcmsConfig := proposalutils.TimelockConfig{MinDelay: 0}
	if mcmsConfig.TimelockQualifierPerChain == nil {
		mcmsConfig.TimelockQualifierPerChain = make(map[uint64]string)
	}
	for _, chainSel := range input.ChainSelectors {
		mcmsConfig.TimelockQualifierPerChain[chainSel] = input.TimelockIdentifier
	}

	cfg := commonchangeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain: contractsByChain,
		MCMSConfig:       mcmsConfig,
	}

	return commonchangeset.TransferToMCMSWithTimelockV2(e, cfg)
}
