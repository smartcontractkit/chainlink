package v1_6

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var AddCandidatesForNewChainChangeset = deployment.CreateChangeSet(addCandidatesForNewChainLogic, addCandidatesForNewChainPrecondition)

type ChainDefinition struct {
	// RMNVerificationDisabled is true if we do not want the RMN to bless messages from this chain.
	RMNVerificationDisabled bool
	// AllowListEnabled is true if we want an allowlist to dictate who can send messages to this chain.
	AllowListEnabled bool
	// Selector is the chain selector of this chain.
	Selector uint64
	// GasPrice defines the USD price (18 decimals) per unit gas for this chain as a destination.
	GasPrice *big.Int
	// TokenPrices define the USD price (18 decimals) per 1e18 of the smallest token denomination for various tokens on this chain.
	TokenPrices map[common.Address]*big.Int
	// FeeQuoterDestChainConfig is the configuration on a fee quoter for this destination.
	FeeQuoterDestChainConfig fee_quoter.FeeQuoterDestChainConfig
}

type NewChainDefinition struct {
	ChainDefinition
	ChainContractParams
	ExistingContracts []commoncs.Contract
	ConfigOnHome      ChainConfig
	OCRParams         CCIPOCRParams
}

// AddCandidatesForNewChainConfig is a configuration struct for AddCandidatesForNewChainChangeset.
type AddCandidatesForNewChainConfig struct {
	HomeChainSelector    uint64
	HomeConfigType       globals.ConfigType
	FeedChainSelector    uint64
	NewChain             NewChainDefinition
	RemoteChains         []ChainDefinition
	MCMSDeploymentConfig commontypes.MCMSWithTimelockConfigV2
	MCMSConfig           *changeset.MCMSConfig
}

func addCandidatesForNewChainPrecondition(e deployment.Environment, c AddCandidatesForNewChainConfig) error {
	// TODO

	return nil
}

func addCandidatesForNewChainLogic(e deployment.Environment, c AddCandidatesForNewChainConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	var allProposals []mcmslib.TimelockProposal

	// 1. Deploy the prerequisite contracts to the new chain
	err := runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return changeset.DeployPrerequisitesChangeset(e, changeset.DeployPrerequisiteConfig{
			Configs: []changeset.DeployPrerequisiteConfigPerChain{
				changeset.DeployPrerequisiteConfigPerChain{
					ChainSelector: c.NewChain.Selector,
				},
			},
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployPrerequisitesChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 2. Save existing contracts
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return commoncs.SaveExistingContractsChangeset(e, commoncs.ExistingContractsConfig{
			ExistingContracts: c.NewChain.ExistingContracts,
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SaveExistingContractsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 3. Deploy MCMS contracts
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return commoncs.DeployMCMSWithTimelockV2(e, map[uint64]commontypes.MCMSWithTimelockConfigV2{
			c.NewChain.Selector: c.MCMSDeploymentConfig,
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployMCMSWithTimelockV2 on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 4. Deploy chain contracts to the new chain
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return DeployChainContractsChangeset(e, DeployChainContractsConfig{
			HomeChainSelector: c.HomeChainSelector,
			ContractParamsPerChain: map[uint64]ChainContractParams{
				c.NewChain.Selector: c.NewChain.ChainContractParams,
			},
		})
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployChainContractsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 5. Update the fee quoter prices on the new chain
	gasPrices := make(map[uint64]*big.Int, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		gasPrices[remoteChain.Selector] = remoteChain.GasPrice
	}
	_, err = UpdateFeeQuoterPricesChangeset(e, UpdateFeeQuoterPricesConfig{
		PricesByChain: map[uint64]FeeQuoterPriceUpdatePerSource{
			c.NewChain.Selector: FeeQuoterPriceUpdatePerSource{
				TokenPrices: c.NewChain.TokenPrices,
				GasPrices:   gasPrices,
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterPricesChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 6. Update the fee quoter destinations on the new chain
	destChainConfigs := make(map[uint64]fee_quoter.FeeQuoterDestChainConfig, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		destChainConfigs[remoteChain.Selector] = remoteChain.FeeQuoterDestChainConfig
	}
	_, err = UpdateFeeQuoterDestsChangeset(e, UpdateFeeQuoterDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			c.NewChain.Selector: destChainConfigs,
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterDestsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 7. Fetch the next DON ID from the capabilities registry
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	donID, err := state.Chains[c.HomeChainSelector].CapabilityRegistry.GetNextDONId(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get next DON ID: %w", err)
	}

	// 8. Add new chain config to the home chain
	out, err := UpdateChainConfigChangeset(e, UpdateChainConfigConfig{
		HomeChainSelector: c.HomeChainSelector,
		RemoteChainAdds: map[uint64]ChainConfig{
			c.NewChain.Selector: c.NewChain.ConfigOnHome,
		},
		MCMS: c.MCMSConfig,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateChainConfigChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// 9. Add the DON to the registry and set candidate for the commit plugin
	out, err = AddDonAndSetCandidateChangeset(e, AddDonAndSetCandidateChangesetConfig{
		SetCandidateConfigBase: SetCandidateConfigBase{
			HomeChainSelector: c.HomeChainSelector,
			FeedChainSelector: c.FeedChainSelector,
			MCMS:              c.MCMSConfig,
		},
		PluginInfo: SetCandidatePluginInfo{
			PluginType: types.PluginTypeCCIPCommit,
			OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
				c.NewChain.Selector: c.NewChain.OCRParams,
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run AddDonAndSetCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// 10. Set the candidate for the exec plugin
	out, err = SetCandidateChangeset(e, SetCandidateChangesetConfig{
		SetCandidateConfigBase: SetCandidateConfigBase{
			HomeChainSelector: c.HomeChainSelector,
			FeedChainSelector: c.FeedChainSelector,
			MCMS:              c.MCMSConfig,
		},
		PluginInfo: []SetCandidatePluginInfo{
			{
				PluginType: types.PluginTypeCCIPExec,
				OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
					c.NewChain.Selector: c.NewChain.OCRParams,
				},
			},
		},
		DonIDOverrides: map[uint64]uint32{
			c.NewChain.Selector: donID,
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// TODO:
	// Add RMN deployment on the new chain?
	//  - SetRMNRemoteConfigChangeset
	//  - SetRMNRemoteOnRMNProxyChangeset

	// 11. Aggregate all proposals.
	var batches []mcmstypes.BatchOperation
	for _, proposal := range allProposals {
		batches = append(batches, proposal.Operations...)
	}
	// Store the timelocks, proposers, and inspectors for each chain.
	timelocks := make(map[uint64]string)
	proposers := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)
	for _, op := range batches {
		chainSel := uint64(op.ChainSelector)
		timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().Hex()
		proposers[chainSel] = state.Chains[chainSel].ProposerMcm.Address().Hex()
		inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get MCMS inspector for chain with selector %d: %w", chainSel, err)
		}
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		fmt.Sprintf("Deploy and add candidates for chain with selector %d", c.NewChain.Selector),
		c.MCMSConfig.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{AddressBook: newAddresses, MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func runAndSaveAddresses(fn func() (deployment.ChangesetOutput, error), new deployment.AddressBook, existing deployment.AddressBook) error {
	output, err := fn()
	if err != nil {
		return fmt.Errorf("failed to run changeset: %w", err)
	}
	err = new.Merge(output.AddressBook)
	if err != nil {
		return fmt.Errorf("failed to update new address book: %w", err)
	}
	err = existing.Merge(output.AddressBook)
	if err != nil {
		return fmt.Errorf("failed to update existing address book: %w", err)
	}

	return nil
}
