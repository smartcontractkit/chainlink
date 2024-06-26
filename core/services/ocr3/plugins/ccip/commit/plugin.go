package commit

import (
	"context"
	"fmt"
	"sort"
	"time"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/ccipocr3/internal/reader"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// Plugin implements the main ocr3 ccip commit plugin logic.
// To learn more about the plugin lifecycle, see the ocr3types.ReportingPlugin interface.
//
// NOTE: If you are changing core plugin logic, you should also update the commit plugin python spec.
type Plugin struct {
	nodeID            commontypes.OracleID
	oracleIDToP2pID   map[commontypes.OracleID]libocrtypes.PeerID
	cfg               cciptypes.CommitPluginConfig
	ccipReader        cciptypes.CCIPReader
	tokenPricesReader cciptypes.TokenPricesReader
	reportCodec       cciptypes.CommitPluginCodec
	msgHasher         cciptypes.MessageHasher
	lggr              logger.Logger

	homeChain reader.HomeChain
}

func NewPlugin(
	_ context.Context,
	nodeID commontypes.OracleID,
	oracleIDToP2pID map[commontypes.OracleID]libocrtypes.PeerID,
	cfg cciptypes.CommitPluginConfig,
	ccipReader cciptypes.CCIPReader,
	tokenPricesReader cciptypes.TokenPricesReader,
	reportCodec cciptypes.CommitPluginCodec,
	msgHasher cciptypes.MessageHasher,
	lggr logger.Logger,
	homeChain reader.HomeChain,
) *Plugin {
	return &Plugin{
		nodeID:            nodeID,
		oracleIDToP2pID:   oracleIDToP2pID,
		cfg:               cfg,
		ccipReader:        ccipReader,
		tokenPricesReader: tokenPricesReader,
		reportCodec:       reportCodec,
		msgHasher:         msgHasher,
		lggr:              lggr,
		homeChain:         homeChain,
	}
}

// Query phase is not used.
func (p *Plugin) Query(_ context.Context, _ ocr3types.OutcomeContext) (types.Query, error) {
	return types.Query{}, nil
}

// Observation phase is used to discover max chain sequence numbers, new messages, gas and token prices.
//
// Max Chain Sequence Numbers:
//
//	It is the sequence number of the last known committed message for each known source chain.
//	If there was a previous outcome we start with the max sequence numbers of the previous outcome.
//	We then read the sequence numbers from the destination chain and override when the on-chain sequence number
//	is greater than previous outcome or when previous outcome did not contain a sequence number for a known source chain.
//
// New Messages:
//
//	We discover new ccip messages only for the chains that the current node is allowed to read from based on the
//	previously discovered max chain sequence numbers. For each chain we scan for new messages
//	in the [max_sequence_number+1, max_sequence_number+1+p.cfg.NewMsgScanBatchSize] range.
//
// Gas Prices:
//
//	We discover the gas prices for each readable source chain.
//
// Token Prices:
//
//	We discover the token prices only for the tokens that are used to pay for ccip fees.
//	The fee tokens are configured in the plugin config.
func (p *Plugin) Observation(
	ctx context.Context, outctx ocr3types.OutcomeContext, _ types.Query,
) (types.Observation, error) {
	supportedChains, err := p.supportedChains()
	if err != nil {
		return types.Observation{}, fmt.Errorf("error finding supported chains by node: %w", err)
	}

	msgBaseDetails := make([]cciptypes.CCIPMsgBaseDetails, 0)
	latestCommittedSeqNumsObservation, err := observeLatestCommittedSeqNums(
		ctx, p.lggr, p.ccipReader, supportedChains, p.cfg.DestChain, p.knownSourceChainsSlice(),
	)
	if err != nil {
		return types.Observation{}, fmt.Errorf("observe latest committed sequence numbers: %w", err)
	}

	var tokenPrices []cciptypes.TokenPrice
	if p.cfg.TokenPricesObserver {
		tokenPrices, err = observeTokenPrices(
			ctx,
			p.tokenPricesReader,
			p.cfg.PricedTokens,
		)
		if err != nil {
			return types.Observation{}, fmt.Errorf("observe token prices: %w", err)
		}
	}

	// Find the gas prices for each source chain.
	var gasPrices []cciptypes.GasPriceChain
	gasPrices, err = observeGasPrices(ctx, p.ccipReader, p.knownSourceChainsSlice())
	if err != nil {
		return types.Observation{}, fmt.Errorf("observe gas prices: %w", err)
	}

	fChain, err := p.homeChain.GetFChain()
	if err != nil {
		return types.Observation{}, fmt.Errorf("get f chain: %w", err)
	}

	// If there's no previous outcome (first round ever), we only observe the latest committed sequence numbers.
	// and on the next round we use those to look for messages.
	if outctx.PreviousOutcome == nil {
		p.lggr.Debugw("first round ever, can't observe new messages yet")
		return cciptypes.NewCommitPluginObservation(
			msgBaseDetails, gasPrices, tokenPrices, latestCommittedSeqNumsObservation, fChain,
		).Encode()
	}

	prevOutcome, err := cciptypes.DecodeCommitPluginOutcome(outctx.PreviousOutcome)
	if err != nil {
		return types.Observation{}, fmt.Errorf("decode commit plugin previous outcome: %w", err)
	}
	p.lggr.Debugw("previous outcome decoded", "outcome", prevOutcome.String())

	// Always observe based on previous outcome. We'll filter out stale messages in the outcome phase.
	newMsgs, err := observeNewMsgs(
		ctx,
		p.lggr,
		p.ccipReader,
		p.msgHasher,
		supportedChains,
		prevOutcome.MaxSeqNums, // TODO: Chainlink common PR to rename.
		p.cfg.NewMsgScanBatchSize,
	)
	if err != nil {
		return types.Observation{}, fmt.Errorf("observe new messages: %w", err)
	}

	p.lggr.Infow("submitting observation",
		"observedNewMsgs", len(newMsgs),
		"gasPrices", len(gasPrices),
		"tokenPrices", len(tokenPrices),
		"latestCommittedSeqNums", latestCommittedSeqNumsObservation,
		"fChain", fChain)

	for _, msg := range newMsgs {
		msgBaseDetails = append(msgBaseDetails, msg.CCIPMsgBaseDetails)
	}

	return cciptypes.NewCommitPluginObservation(
		msgBaseDetails, gasPrices, tokenPrices, latestCommittedSeqNumsObservation, fChain,
	).Encode()

}

func (p *Plugin) ValidateObservation(_ ocr3types.OutcomeContext, _ types.Query, ao types.AttributedObservation) error {
	obs, err := cciptypes.DecodeCommitPluginObservation(ao.Observation)
	if err != nil {
		return fmt.Errorf("decode commit plugin observation: %w", err)
	}

	if err := validateObservedSequenceNumbers(obs.NewMsgs, obs.MaxSeqNums); err != nil {
		return fmt.Errorf("validate sequence numbers: %w", err)
	}

	destSupportedChains, err := p.supportedChains()
	if err != nil {
		return fmt.Errorf("error finding supported chains by node: %w", err)
	}

	err = validateObserverReadingEligibility(obs.NewMsgs, obs.MaxSeqNums, destSupportedChains, p.cfg.DestChain)
	if err != nil {
		return fmt.Errorf("validate observer %d reading eligibility: %w", ao.Observer, err)
	}

	if err := validateObservedTokenPrices(obs.TokenPrices); err != nil {
		return fmt.Errorf("validate token prices: %w", err)
	}

	if err := validateObservedGasPrices(obs.GasPrices); err != nil {
		return fmt.Errorf("validate gas prices: %w", err)
	}

	return nil
}

func (p *Plugin) ObservationQuorum(_ ocr3types.OutcomeContext, _ types.Query) (ocr3types.Quorum, error) {
	// Across all chains we require at least 2F+1 observations.
	return ocr3types.QuorumTwoFPlusOne, nil
}

// Outcome phase is used to construct the final outcome based on the observations of multiple followers.
//
// The outcome contains:
//   - Max Sequence Numbers: The max sequence number for each source chain.
//   - Merkle Roots: One merkle tree root per source chain. The leaves of the tree are the IDs of the observed messages.
//     The merkle root data type contains information about the chain and the sequence numbers range.
func (p *Plugin) Outcome(
	_ ocr3types.OutcomeContext, _ types.Query, aos []types.AttributedObservation,
) (ocr3types.Outcome, error) {
	decodedObservations := make([]cciptypes.CommitPluginObservation, 0)
	for _, ao := range aos {
		obs, err := cciptypes.DecodeCommitPluginObservation(ao.Observation)
		if err != nil {
			return ocr3types.Outcome{}, fmt.Errorf("decode commit plugin observation: %w", err)
		}
		decodedObservations = append(decodedObservations, obs)
	}

	fChains := fChainConsensus(decodedObservations)

	fChainDest, ok := fChains[p.cfg.DestChain]
	if !ok {
		return ocr3types.Outcome{}, fmt.Errorf("missing destination chain %d in fChain config", p.cfg.DestChain)
	}

	maxSeqNums := maxSeqNumsConsensus(p.lggr, fChainDest, decodedObservations)
	p.lggr.Debugw("max sequence numbers consensus", "maxSeqNumsConsensus", maxSeqNums)

	merkleRoots, err := newMsgsConsensus(p.lggr, maxSeqNums, decodedObservations, fChains)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("new messages consensus: %w", err)
	}
	p.lggr.Debugw("new messages consensus", "merkleRoots", merkleRoots)

	tokenPrices, err := tokenPricesConsensus(decodedObservations, fChainDest)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("token prices consensus: %w", err)
	}

	gasPrices := gasPricesConsensus(p.lggr, decodedObservations, fChainDest)
	p.lggr.Debugw("gas prices consensus", "gasPrices", gasPrices)

	outcome := cciptypes.NewCommitPluginOutcome(maxSeqNums, merkleRoots, tokenPrices, gasPrices)
	if outcome.IsEmpty() {
		p.lggr.Debugw("empty outcome")
		return ocr3types.Outcome{}, nil
	}
	p.lggr.Debugw("sending outcome", "outcome", outcome)

	return outcome.Encode()
}

func (p *Plugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[[]byte], error) {
	outc, err := cciptypes.DecodeCommitPluginOutcome(outcome)
	if err != nil {
		p.lggr.Errorw("decode commit plugin outcome", "outcome", outcome, "err", err)
		return nil, fmt.Errorf("decode commit plugin outcome: %w", err)
	}

	/*
		todo: Once token/gas prices are implemented, we would want to probably check if outc.MerkleRoots is empty or not
		and only create a report if outc.MerkleRoots is non-empty OR gas/token price timer has expired
	*/

	rep := cciptypes.NewCommitPluginReport(outc.MerkleRoots, outc.TokenPrices, outc.GasPrices)

	encodedReport, err := p.reportCodec.Encode(context.Background(), rep)
	if err != nil {
		return nil, fmt.Errorf("encode commit plugin report: %w", err)
	}

	return []ocr3types.ReportWithInfo[[]byte]{{Report: encodedReport, Info: nil}}, nil
}

func (p *Plugin) ShouldAcceptAttestedReport(
	ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte],
) (bool, error) {
	decodedReport, err := p.reportCodec.Decode(ctx, r.Report)
	if err != nil {
		return false, fmt.Errorf("decode commit plugin report: %w", err)
	}

	isEmpty := decodedReport.IsEmpty()
	if isEmpty {
		p.lggr.Infow("skipping empty report")
		return false, nil
	}

	return true, nil
}

func (p *Plugin) ShouldTransmitAcceptedReport(
	ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte],
) (bool, error) {
	isWriter, err := p.supportsDestChain()
	if err != nil {
		return false, fmt.Errorf("can't know if it's a writer: %w", err)
	}
	if !isWriter {
		p.lggr.Debugw("not a writer, skipping report transmission")
		return false, nil
	}

	decodedReport, err := p.reportCodec.Decode(ctx, r.Report)
	if err != nil {
		return false, fmt.Errorf("decode commit plugin report: %w", err)
	}

	p.lggr.Debugw("transmitting report",
		"roots", len(decodedReport.MerkleRoots),
		"tokenPriceUpdates", len(decodedReport.PriceUpdates.TokenPriceUpdates),
		"gasPriceUpdates", len(decodedReport.PriceUpdates.GasPriceUpdates),
	)

	// todo: if report is stale -> do not transmit (check the spec for the exact condition)
	return true, nil
}

func (p *Plugin) Close() error {
	timeout := 10 * time.Second
	ctx, cf := context.WithTimeout(context.Background(), timeout)
	defer cf()

	if err := p.ccipReader.Close(ctx); err != nil {
		return fmt.Errorf("close ccip reader: %w", err)
	}
	return nil
}

func (p *Plugin) knownSourceChainsSlice() []cciptypes.ChainSelector {
	knownSourceChains, err := p.homeChain.GetKnownCCIPChains()
	if err != nil {
		p.lggr.Errorw("error getting known chains", "err", err)
		return nil
	}
	knownSourceChainsSlice := knownSourceChains.ToSlice()
	sort.Slice(
		knownSourceChainsSlice,
		func(i, j int) bool { return knownSourceChainsSlice[i] < knownSourceChainsSlice[j] },
	)
	return slicelib.Filter(knownSourceChainsSlice, func(ch cciptypes.ChainSelector) bool { return ch != p.cfg.DestChain })
}

func (p *Plugin) supportedChains() (mapset.Set[cciptypes.ChainSelector], error) {
	p2pID, exists := p.oracleIDToP2pID[p.nodeID]
	if !exists {
		return nil, fmt.Errorf("oracle ID %d not found in oracleIDToP2pID", p.nodeID)
	}
	supportedChains, err := p.homeChain.GetSupportedChainsForPeer(p2pID)
	if err != nil {
		p.lggr.Warnw("error getting supported chains", err)
		return mapset.NewSet[cciptypes.ChainSelector](), fmt.Errorf("error getting supported chains: %w", err)
	}

	return supportedChains, nil
}

func (p *Plugin) supportsDestChain() (bool, error) {
	destChainConfig, err := p.homeChain.GetChainConfig(p.cfg.DestChain)
	if err != nil {
		return false, fmt.Errorf("get chain config: %w", err)
	}
	return destChainConfig.SupportedNodes.Contains(p.oracleIDToP2pID[p.nodeID]), nil
}

// Interface compatibility checks.
var _ ocr3types.ReportingPlugin[[]byte] = &Plugin{}
