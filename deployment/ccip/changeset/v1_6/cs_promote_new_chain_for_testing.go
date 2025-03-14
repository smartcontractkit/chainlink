package v1_6

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var PromoteNewChainForTestingChangeset = deployment.CreateChangeSet(promoteNewChainForTestingLogic, promoteNewChainForTestingPrecondition)

var (
	PREFIX      = new(big.Int).Lsh(big.NewInt(0x000a), 240) // 0x000a << 240
	PREFIX_MASK = new(big.Int).Lsh(big.NewInt(0xFFFF), 240) // 0xFFFF << 240
)

// PromoteNewChainForTestingConfig is a configuration struct for PromoteNewChainForTestingChangeset.
type PromoteNewChainForTestingConfig struct {
	HomeChainSelector uint64
	HomeConfigType    globals.ConfigType
	NewChain          NewChainDefinition
	RemoteChains      []ChainDefinition
	MCMSConfig        *changeset.MCMSConfig
}

func promoteNewChainForTestingPrecondition(e deployment.Environment, c PromoteNewChainForTestingConfig) error {
	// TODO

	return nil
}

func promoteNewChainForTestingLogic(e deployment.Environment, c PromoteNewChainForTestingConfig) (deployment.ChangesetOutput, error) {
	var allProposals []mcmslib.TimelockProposal

	// 1. Promote the candidates for the commit and exec plugins
	out, err := PromoteCandidateChangeset(e, PromoteCandidateChangesetConfig{
		HomeChainSelector: c.HomeChainSelector,
		MCMS:              c.MCMSConfig,
		PluginInfo: []PromoteCandidatePluginInfo{
			{
				PluginType:           types.PluginTypeCCIPCommit,
				RemoteChainSelectors: []uint64{c.NewChain.Selector},
			},
			{
				PluginType:           types.PluginTypeCCIPExec,
				RemoteChainSelectors: []uint64{c.NewChain.Selector}},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run PromoteCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// 2. Set the OCR3 config on the off ramp on the new chain
	out, err = SetOCR3OffRampChangeset(e, SetOCR3OffRampConfig{
		HomeChainSel:       c.HomeChainSelector,
		RemoteChainSels:    []uint64{c.NewChain.Selector},
		CCIPHomeConfigType: c.HomeConfigType,
		MCMS:               c.MCMSConfig,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetOCR3OffRampChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// 3. Update the fee quoter prices and destinations on the remote chains
	for _, remoteChain := range c.RemoteChains {
		out, err := UpdateFeeQuoterPricesChangeset(e, UpdateFeeQuoterPricesConfig{
			PricesByChain: map[uint64]FeeQuoterPriceUpdatePerSource{
				remoteChain.Selector: FeeQuoterPriceUpdatePerSource{
					TokenPrices: remoteChain.TokenPrices,
					GasPrices:   map[uint64]*big.Int{c.NewChain.Selector: c.NewChain.GasPrice},
				},
			},
			MCMS: c.MCMSConfig,
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterPricesChangeset on chain with selector %d: %w", remoteChain.Selector, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		_, err = UpdateFeeQuoterDestsChangeset(e, UpdateFeeQuoterDestsConfig{
			UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
				c.NewChain.Selector: map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					remoteChain.Selector: c.NewChain.FeeQuoterDestChainConfig,
				},
			},
			MCMS: c.MCMSConfig,
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterDestsChangeset on chain with selector %d: %w", remoteChain.Selector, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)
	}

	// 4. Connect the new chain to the existing chains (use the test router)
	testRouter := true
	connections := make(map[uint64]ConnectionConfig, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		connections[remoteChain.Selector] = ConnectionConfig{
			RMNVerificationDisabled: remoteChain.RMNVerificationDisabled,
			AllowListEnabled:        remoteChain.AllowListEnabled,
		}
	}
	cfg := ConnectNewChainConfig{
		RemoteChains:     connections,
		NewChainSelector: c.NewChain.Selector,
		NewChainConnectionConfig: ConnectionConfig{
			RMNVerificationDisabled: c.NewChain.RMNVerificationDisabled,
			AllowListEnabled:        c.NewChain.AllowListEnabled,
		},
		TestRouter: &testRouter,
		MCMSConfig: c.MCMSConfig,
	}
	err = ConnectNewChainChangeset.VerifyPreconditions(e, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run ConnectNewChainChangeset precondition: %w", err)
	}
	out, err = ConnectNewChainChangeset.Apply(e, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run ConnectNewChainChangeset: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// 5. Aggregate all proposals.
	var batches []mcmstypes.BatchOperation
	for _, proposal := range allProposals {
		batches = append(batches, proposal.Operations...)
	}
	// Store the timelocks, proposers, and inspectors for each chain.
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
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
		fmt.Sprintf("Connect chain with selector %d to other chains", c.NewChain.Selector),
		c.MCMSConfig.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}
