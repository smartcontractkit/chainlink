package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

var _ deployment.ChangeSet[ChainInboundChangesetConfig] = NewChainInboundChangeset

type ChainInboundChangesetConfig struct {
	HomeChainSelector    uint64
	NewChainSelector     uint64
	SourceChainSelectors []uint64
}

func (c ChainInboundChangesetConfig) Validate() error {
	if c.HomeChainSelector == 0 {
		return fmt.Errorf("HomeChainSelector must be set")
	}
	if c.NewChainSelector == 0 {
		return fmt.Errorf("NewChainSelector must be set")
	}
	if len(c.SourceChainSelectors) == 0 {
		return fmt.Errorf("SourceChainSelectors must be set")
	}
	return nil
}

// NewChainInboundChangeset generates a proposal
// to connect the new chain to the existing chains.
func NewChainInboundChangeset(
	e deployment.Environment,
	cfg ChainInboundChangesetConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// Generate proposal which enables new destination (from test router) on all source chains.
	var batches []timelock.BatchChainOperation
	for _, source := range cfg.SourceChainSelectors {
		enableOnRampDest, err := state.Chains[source].OnRamp.ApplyDestChainConfigUpdates(deployment.SimTransactOpts(), []onramp.OnRampDestChainConfigArgs{
			{
				DestChainSelector: cfg.NewChainSelector,
				Router:            state.Chains[source].TestRouter.Address(),
			},
		})
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		enableFeeQuoterDest, err := state.Chains[source].FeeQuoter.ApplyDestChainConfigUpdates(
			deployment.SimTransactOpts(),
			[]fee_quoter.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: cfg.NewChainSelector,
					DestChainConfig:   DefaultFeeQuoterDestChainConfig(),
				},
			})
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(source),
			Batch: []mcms.Operation{
				{
					// Enable the source in on ramp
					To:    state.Chains[source].OnRamp.Address(),
					Data:  enableOnRampDest.Data(),
					Value: big.NewInt(0),
				},
				{
					To:    state.Chains[source].FeeQuoter.Address(),
					Data:  enableFeeQuoterDest.Data(),
					Value: big.NewInt(0),
				},
			},
		})
	}

	addChainOp, err := applyChainConfigUpdatesOp(e, state, cfg.HomeChainSelector, []uint64{cfg.NewChainSelector})
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
		Batch: []mcms.Operation{
			addChainOp,
		},
	})

	var (
		timelocksPerChain = make(map[uint64]common.Address)
		proposerMCMSes    = make(map[uint64]*gethwrappers.ManyChainMultiSig)
	)
	for _, chain := range append(cfg.SourceChainSelectors, cfg.HomeChainSelector) {
		timelocksPerChain[chain] = state.Chains[chain].Timelock.Address()
		proposerMCMSes[chain] = state.Chains[chain].ProposerMcm
	}
	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to set new chains",
		0,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

func applyChainConfigUpdatesOp(
	e deployment.Environment,
	state CCIPOnChainState,
	homeChainSel uint64,
	chains []uint64,
) (mcms.Operation, error) {
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return mcms.Operation{}, err
	}
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return mcms.Operation{}, err
	}
	var chainConfigUpdates []ccip_home.CCIPHomeChainConfigArgs
	for _, chainSel := range chains {
		chainConfig := setupConfigInfo(chainSel, nodes.NonBootstraps().PeerIDs(),
			nodes.DefaultF(), encodedExtraChainConfig)
		chainConfigUpdates = append(chainConfigUpdates, chainConfig)
	}

	addChain, err := state.Chains[homeChainSel].CCIPHome.ApplyChainConfigUpdates(
		deployment.SimTransactOpts(),
		nil,
		chainConfigUpdates,
	)
	if err != nil {
		return mcms.Operation{}, err
	}
	return mcms.Operation{
		To:    state.Chains[homeChainSel].CCIPHome.Address(),
		Data:  addChain.Data(),
		Value: big.NewInt(0),
	}, nil
}
