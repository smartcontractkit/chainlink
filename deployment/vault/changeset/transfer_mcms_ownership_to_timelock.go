package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// TransferMCMSOwnershipToTimelockChangeset transfers ownership of Bypasser, Canceller, and
// Proposer ManyChainMultiSig contracts to the RBAC Timelock (excluding CallProxy).
// It performs the transfer from KMS/deployer and builds an MCMS proposal for acceptOwnership.
// Use for: (1) migration of existing deployed contracts, (2) after deploy_timelock + set_mcms_config for new chains.
var TransferMCMSOwnershipToTimelockChangeset cldf.ChangeSetV2[types.TransferMCMSOwnershipToTimelockConfig] = transferMCMSOwnershipToTimelockChangeset{}

type transferMCMSOwnershipToTimelockChangeset struct{}

func (t transferMCMSOwnershipToTimelockChangeset) VerifyPreconditions(e cldf.Environment, cfg types.TransferMCMSOwnershipToTimelockConfig) error {
	return ValidateTransferMCMSOwnershipToTimelockConfig(e, cfg)
}

func (t transferMCMSOwnershipToTimelockChangeset) Apply(e cldf.Environment, cfg types.TransferMCMSOwnershipToTimelockConfig) (cldf.ChangesetOutput, error) {
	qualifier := cfg.TimelockIdentifier
	if qualifier == "" {
		qualifier = commonchangeset.DefaultTimelockQualifier
	}

	contractsByChain := make(map[uint64][]common.Address)
	for _, chainSelector := range cfg.ChainSelectors {
		addresses, err := commonstate.AddressesForChain(e, chainSelector, qualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: load addresses: %w", chainSelector, err)
		}

		var addrs []common.Address
		for addr, tv := range addresses {
			switch tv.Type {
			case commontypes.BypasserManyChainMultisig, commontypes.CancellerManyChainMultisig, commontypes.ProposerManyChainMultisig:
				addrs = append(addrs, common.HexToAddress(addr))
			}
		}
		if len(addrs) == 0 {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: no Bypasser/Canceller/Proposer MCMS addresses found in state", chainSelector)
		}
		contractsByChain[chainSelector] = addrs
	}

	mcmsConfig := proposalutils.TimelockConfig{MinDelay: 0}
	if cfg.MCMSConfig != nil {
		mcmsConfig = *cfg.MCMSConfig
	}
	if mcmsConfig.TimelockQualifierPerChain == nil {
		mcmsConfig.TimelockQualifierPerChain = make(map[uint64]string)
	}
	for _, chainSel := range cfg.ChainSelectors {
		mcmsConfig.TimelockQualifierPerChain[chainSel] = qualifier
	}

	return commonchangeset.TransferToMCMSWithTimelockV2(e, commonchangeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain:     contractsByChain,
		MCMSConfig:          mcmsConfig,
		OnlyAcceptOwnership: cfg.OnlyAcceptOwnership,
	})
}
