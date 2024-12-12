package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
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

type AddDonAndSetCandidateChangesetConfig struct {
	HomeChainSelector uint64
	FeedChainSelector uint64
	NewChainSelector  uint64
	PluginType        types.PluginType
	NodeIDs           []string
	CCIPOCRParams     CCIPOCRParams
}

func (a AddDonAndSetCandidateChangesetConfig) Validate(e deployment.Environment, state CCIPOnChainState) (deployment.Nodes, error) {
	if a.HomeChainSelector == 0 {
		return nil, fmt.Errorf("HomeChainSelector must be set")
	}
	if a.FeedChainSelector == 0 {
		return nil, fmt.Errorf("FeedChainSelector must be set")
	}
	if a.NewChainSelector == 0 {
		return nil, fmt.Errorf("ocr config chain selector must be set")
	}
	if a.PluginType != types.PluginTypeCCIPCommit &&
		a.PluginType != types.PluginTypeCCIPExec {
		return nil, fmt.Errorf("PluginType must be set to either CCIPCommit or CCIPExec")
	}
	// TODO: validate token config
	if len(a.NodeIDs) == 0 {
		return nil, fmt.Errorf("nodeIDs must be set")
	}
	nodes, err := deployment.NodeInfo(a.NodeIDs, e.Offchain)
	if err != nil {
		return nil, fmt.Errorf("get node info: %w", err)
	}

	// check that chain config is set up for the new chain
	chainConfig, err := state.Chains[a.HomeChainSelector].CCIPHome.GetChainConfig(nil, a.NewChainSelector)
	if err != nil {
		return nil, fmt.Errorf("get all chain configs: %w", err)
	}

	// FChain should never be zero if a chain config is set in CCIPHome
	if chainConfig.FChain == 0 {
		return nil, fmt.Errorf("chain config not set up for new chain %d", a.NewChainSelector)
	}

	err = a.CCIPOCRParams.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid ccip ocr params: %w", err)
	}

	if e.OCRSecrets.IsEmpty() {
		return nil, fmt.Errorf("OCR secrets must be set")
	}

	return nodes, nil
}

// AddDonAndSetCandidateChangeset adds new DON for destination to home chain
// and sets the commit plugin config as candidateConfig for the don.
func AddDonAndSetCandidateChangeset(
	e deployment.Environment,
	cfg AddDonAndSetCandidateChangesetConfig,
) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	nodes, err := cfg.Validate(e, state)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, err)
	}

	newDONArgs, err := internal.BuildOCR3ConfigForCCIPHome(
		e.OCRSecrets,
		state.Chains[cfg.NewChainSelector].OffRamp,
		e.Chains[cfg.NewChainSelector],
		nodes.NonBootstraps(),
		state.Chains[cfg.HomeChainSelector].RMNHome.Address(),
		cfg.CCIPOCRParams.OCRParameters,
		cfg.CCIPOCRParams.CommitOffChainConfig,
		cfg.CCIPOCRParams.ExecuteOffChainConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	latestDon, err := internal.LatestCCIPDON(state.Chains[cfg.HomeChainSelector].CapabilityRegistry)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	commitConfig, ok := newDONArgs[cfg.PluginType]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing commit plugin in ocr3Configs")
	}
	donID := latestDon.Id + 1
	addDonOp, err := newDonWithCandidateOp(
		donID, commitConfig,
		state.Chains[cfg.HomeChainSelector].CapabilityRegistry,
		nodes.NonBootstraps(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	var (
		timelocksPerChain = map[uint64]common.Address{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].Timelock.Address(),
		}
		proposerMCMSes = map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSelector: state.Chains[cfg.HomeChainSelector].ProposerMcm,
		}
	)
	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
			Batch:           []mcms.Operation{addDonOp},
		}},
		"setCandidate for commit and AddDon on new Chain",
		0, // minDelay
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w", err)
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

// newDonWithCandidateOp sets the candidate commit config by calling setCandidate on CCIPHome contract through the AddDON call on CapReg contract
// This should be done first before calling any other UpdateDON calls
// This proposes to set up OCR3 config for the commit plugin for the DON
func newDonWithCandidateOp(
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes,
) (mcms.Operation, error) {
	encodedSetCandidateCall, err := internal.CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		pluginConfig.PluginType,
		pluginConfig,
		[32]byte{},
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("pack set candidate call: %w", err)
	}
	addDonTx, err := capReg.AddDON(deployment.SimTransactOpts(), nodes.PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: internal.CCIPCapabilityID,
			Config:       encodedSetCandidateCall,
		},
	}, false, false, nodes.DefaultF())
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("could not generate add don tx w/ commit config: %w", err)
	}
	return mcms.Operation{
		To:    capReg.Address(),
		Data:  addDonTx.Data(),
		Value: big.NewInt(0),
	}, nil
}
