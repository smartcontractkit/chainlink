package adapters

import (
	"encoding/json"
	"fmt"

	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	mcms_utils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// AptosCurseMCMSReader implements changesets.MCMSReader for the Aptos family,
// targeting the CurseMCMS contract rather than the regular MCMS contract.
// It is registered for FamilyAptos so that the fastcurse framework's
// OutputBuilder produces proposals against CurseMCMS.
type AptosCurseMCMSReader struct{}

func (r *AptosCurseMCMSReader) GetChainMetadata(e deployment.Environment, chainSelector uint64, input mcms_utils.Input) (mcmstypes.ChainMetadata, error) {
	chain, ok := e.BlockChains.AptosChains()[chainSelector]
	if !ok {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("aptos chain with selector %d not found", chainSelector)
	}

	stateMap, err := aptos.LoadOnchainStateAptos(e)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to load aptos onchain state: %w", err)
	}
	state, ok := stateMap[chainSelector]
	if !ok {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("aptos chain %d not found in state", chainSelector)
	}
	curseMCMSAddr := state.CurseMCMSAddress

	role, err := proposalutils.GetAptosRoleFromAction(input.TimelockAction)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to get role from action: %w", err)
	}
	inspector := aptosmcms.NewInspectorWithMCMSType(chain.Client, role, aptosmcms.MCMSTypeCurse)

	opCount, err := inspector.GetOpCount(e.GetContext(), curseMCMSAddr.StringLong())
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to get opCount for CurseMCMS at %s on chain %d: %w", curseMCMSAddr.StringLong(), chainSelector, err)
	}

	afm := aptosmcms.AdditionalFieldsMetadata{
		MCMSType: aptosmcms.MCMSTypeCurse,
	}
	afBytes, err := json.Marshal(afm)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to marshal additional fields metadata: %w", err)
	}

	return mcmstypes.ChainMetadata{
		StartingOpCount:  opCount,
		MCMAddress:       curseMCMSAddr.StringLong(),
		AdditionalFields: afBytes,
	}, nil
}

func (r *AptosCurseMCMSReader) GetTimelockRef(e deployment.Environment, chainSelector uint64, input mcms_utils.Input) (datastore.AddressRef, error) {
	stateMap, err := aptos.LoadOnchainStateAptos(e)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("failed to load aptos onchain state: %w", err)
	}
	state, ok := stateMap[chainSelector]
	if !ok {
		return datastore.AddressRef{}, fmt.Errorf("aptos chain %d not found in state", chainSelector)
	}
	return datastore.AddressRef{
		Address:       state.CurseMCMSAddress.StringLong(),
		ChainSelector: chainSelector,
	}, nil
}

func (r *AptosCurseMCMSReader) GetMCMSRef(e deployment.Environment, chainSelector uint64, input mcms_utils.Input) (datastore.AddressRef, error) {
	stateMap, err := aptos.LoadOnchainStateAptos(e)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("failed to load aptos onchain state: %w", err)
	}
	state, ok := stateMap[chainSelector]
	if !ok {
		return datastore.AddressRef{}, fmt.Errorf("aptos chain %d not found in state", chainSelector)
	}
	return datastore.AddressRef{
		Address:       state.CurseMCMSAddress.StringLong(),
		ChainSelector: chainSelector,
	}, nil
}
