package v1_6

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
)

var _ cldf.ChangeSet[ccipseq.DeployChainContractsConfig] = DeployChainContractsChangeset

// DeployChainContracts deploys all new CCIP v1.6 or later contracts for the given chains.
// It returns the new addresses for the contracts.
// DeployChainContractsChangeset is idempotent. If there is an error, it will return the successfully deployed addresses and the error so that the caller can call the
// changeset again with the same input to retry the failed deployment.
// Caller should update the environment's address book with the returned addresses.
// Points to note :
// In case of migrating from legacy ccip to 1.6, the previous RMN address should be set while deploying RMNRemote.
// if there is no existing RMN address found, RMNRemote will be deployed with 0x0 address for previous RMN address
// which will set RMN to 0x0 address immutably in RMNRemote.
func DeployChainContractsChangeset(env cldf.Environment, c ccipseq.DeployChainContractsConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses, err := deployChainContractsForChains(env, c.HomeChainSelector, c.ContractParamsPerChain)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return cldf.ChangesetOutput{AddressBook: newAddresses}, cldf.MaybeDataErr(err)
	}
	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

func ValidateHomeChainState(e cldf.Environment, homeChainSel uint64, existingState stateview.CCIPOnChainState) error {
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return err
	}
	capReg := existingState.Chains[homeChainSel].CapabilityRegistry
	if capReg == nil {
		e.Logger.Errorw("Failed to get capability registry")
		return errors.New("capability registry not found")
	}
	cr, err := capReg.GetHashedCapabilityId(
		&bind.CallOpts{}, shared.CapabilityLabelledName, shared.CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return err
	}
	if cr != shared.CCIPCapabilityID {
		return fmt.Errorf("unexpected mismatch between calculated ccip capability id (%s) and expected ccip capability id constant (%s)",
			hexutil.Encode(cr[:]),
			hexutil.Encode(shared.CCIPCapabilityID[:]))
	}
	capability, err := capReg.GetCapability(nil, shared.CCIPCapabilityID)
	if err != nil {
		e.Logger.Errorw("Failed to get capability", "err", err)
		return err
	}
	ccipHome, err := ccip_home.NewCCIPHome(capability.ConfigurationContract, e.BlockChains.EVMChains()[homeChainSel].Client)
	if err != nil {
		e.Logger.Errorw("Failed to get ccip config", "err", err)
		return err
	}
	if ccipHome.Address() != existingState.Chains[homeChainSel].CCIPHome.Address() {
		return errors.New("ccip home address mismatch")
	}
	rmnHome := existingState.Chains[homeChainSel].RMNHome
	if rmnHome == nil {
		e.Logger.Errorw("Failed to get rmn home", "err", err)
		return errors.New("rmn home not found")
	}
	return nil
}

func deployChainContractsForChains(
	e cldf.Environment,
	homeChainSel uint64,
	contractParamsPerChain map[uint64]ccipseq.ChainContractParams) (cldf.AddressBook, error) {
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return nil, err
	}

	err = ValidateHomeChainState(e, homeChainSel, existingState)
	if err != nil {
		return nil, err
	}
	deps := opsutil.ConfigureDependencies{
		Env:          e,
		CurrentState: existingState,
		AddressBook:  cldf.NewMemoryAddressBook(),
	}
	_, err = operations.ExecuteSequence(e.OperationsBundle, ccipseq.DeployChainContractsSeq, deps, ccipseq.DeployChainContractsConfig{
		HomeChainSelector:      homeChainSel,
		ContractParamsPerChain: contractParamsPerChain,
	})
	if err != nil {
		e.Logger.Errorw("Failed to deploy chain contracts", "err", err)
		return nil, err
	}

	return deps.AddressBook, nil
}
