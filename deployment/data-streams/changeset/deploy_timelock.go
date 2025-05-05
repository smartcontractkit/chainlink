package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
)

var DeployAndTransferMCMSChangeset = cldf.CreateChangeSet(deployAndTransferMcmsLogic, deployAndTransferMcmsPrecondition)

type DeployMCMSConfig struct {
	ChainsToDeploy []uint64
	Ownership      types.OwnershipSettings
	Config         commontypes.MCMSWithTimelockConfigV2
}

func deployAndTransferMcmsLogic(e deployment.Environment, cc DeployMCMSConfig) (deployment.ChangesetOutput, error) {
	cfgByChain := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, chain := range cc.ChainsToDeploy {
		cfgByChain[chain] = cc.Config
	}

	mcmsOut, err := commonChangesets.DeployMCMSWithTimelockV2(e, cfgByChain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS: %w", err)
	}

	// CallProxy has no owner, RBACTimelock has an "admin" setting in place of owner
	transferContracts := []deployment.ContractType{
		commontypes.ProposerManyChainMultisig,
		commontypes.BypasserManyChainMultisig,
		commontypes.CancellerManyChainMultisig,
	}

	var proposals []mcms.TimelockProposal
	if cc.Ownership.ShouldTransfer && cc.Ownership.MCMSProposalConfig != nil {
		for _, contractType := range transferContracts {
			// all MCMS contracts are version 1.0.0 right now
			contractFilter := deployment.NewTypeAndVersion(contractType, deployment.Version1_0_0)
			contractTransfer, err := mcmsutil.TransferToMCMSWithTimelockForTypeAndVersion(e, mcmsOut.AddressBook, contractFilter, *cc.Ownership.MCMSProposalConfig)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: mcmsOut.AddressBook}, fmt.Errorf("failed to transfer %s to MCMS: %w", contractType, err)
			}
			proposals = append(proposals, contractTransfer.MCMSTimelockProposals...)
		}
	}

	return deployment.ChangesetOutput{AddressBook: mcmsOut.AddressBook, MCMSTimelockProposals: proposals}, nil
}

func deployAndTransferMcmsPrecondition(_ deployment.Environment, cc DeployMCMSConfig) error {
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
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}
