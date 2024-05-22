package commit

import (
	"context"
	"fmt"
	"sort"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccipocr3/internal/codec"
	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"
	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/smartcontractkit/ccipocr3/internal/reader"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Plugin implements the main ocr3 ccip commit plugin logic.
// To learn more about the plugin lifecycle, see the ocr3types.ReportingPlugin interface.
type Plugin struct {
	nodeID            commontypes.OracleID
	cfg               model.CommitPluginConfig
	ccipReader        reader.CCIP
	tokenPricesReader reader.TokenPrices
	reportCodec       codec.Commit
	lggr              logger.Logger

	// readableChains is the set of chains that the plugin can read from.
	readableChains mapset.Set[model.ChainSelector]
	// knownSourceChains is the set of chains that the plugin knows about.
	knownSourceChains mapset.Set[model.ChainSelector]
}

// TODO: background service for home chain config polling

func NewPlugin(
	_ context.Context,
	nodeID commontypes.OracleID,
	cfg model.CommitPluginConfig,
	ccipReader reader.CCIP,
	tokenPricesReader reader.TokenPrices,
	reportCodec codec.Commit,
	lggr logger.Logger,
) *Plugin {
	knownSourceChains := mapset.NewSet[model.ChainSelector](cfg.Reads...)
	for _, inf := range cfg.ObserverInfo {
		knownSourceChains = knownSourceChains.Union(mapset.NewSet(inf.Reads...))
	}

	return &Plugin{
		nodeID:            nodeID,
		cfg:               cfg,
		ccipReader:        ccipReader,
		tokenPricesReader: tokenPricesReader,
		reportCodec:       reportCodec,
		lggr:              lggr,

		readableChains:    mapset.NewSet(cfg.Reads...),
		knownSourceChains: knownSourceChains,
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
func (p *Plugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, _ types.Query) (types.Observation, error) {
	maxSeqNumsPerChain, err := observeMaxSeqNums(
		ctx,
		p.lggr,
		p.ccipReader,
		outctx.PreviousOutcome,
		p.readableChains,
		p.cfg.DestChain,
		p.knownSourceChainsSlice(),
	)
	if err != nil {
		return types.Observation{}, fmt.Errorf("observe max sequence numbers per chain: %w", err)
	}

	newMsgs, err := observeNewMsgs(
		ctx,
		p.lggr,
		p.ccipReader,
		p.readableChains,
		maxSeqNumsPerChain,
		p.cfg.NewMsgScanBatchSize,
	)
	if err != nil {
		return types.Observation{}, fmt.Errorf("observe new messages: %w", err)
	}

	var tokenPrices []model.TokenPrice
	if p.cfg.TokenPricesObserver {
		tokenPrices, err = observeTokenPrices(
			ctx,
			p.tokenPricesReader,
			p.cfg.FeeTokens,
		)
		if err != nil {
			return types.Observation{}, fmt.Errorf("observe token prices: %w", err)
		}
	}

	// Find the gas prices for each source chain.
	var gasPrices []model.GasPriceChain
	gasPrices, err = observeGasPrices(ctx, p.ccipReader, p.knownSourceChainsSlice())
	if err != nil {
		return types.Observation{}, fmt.Errorf("observe gas prices: %w", err)
	}

	return model.NewCommitPluginObservation(newMsgs, gasPrices, tokenPrices, maxSeqNumsPerChain).Encode()
}

func (p *Plugin) ValidateObservation(_ ocr3types.OutcomeContext, _ types.Query, ao types.AttributedObservation) error {
	obs, err := model.DecodeCommitPluginObservation(ao.Observation)
	if err != nil {
		return fmt.Errorf("decode commit plugin observation: %w", err)
	}

	if err := validateObservedSequenceNumbers(obs.NewMsgs, obs.MaxSeqNums); err != nil {
		return fmt.Errorf("validate sequence numbers: %w", err)
	}

	if err := validateObserverReadingEligibility(ao.Observer, obs.NewMsgs, p.cfg.ObserverInfo); err != nil {
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
func (p *Plugin) Outcome(_ ocr3types.OutcomeContext, _ types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	fChainDest, ok := p.cfg.FChain[p.cfg.DestChain]
	if !ok {
		return ocr3types.Outcome{}, fmt.Errorf("missing destination chain %d in fChain config", p.cfg.DestChain)
	}

	decodedObservations := make([]model.CommitPluginObservation, 0)
	for _, ao := range aos {
		obs, err := model.DecodeCommitPluginObservation(ao.Observation)
		if err != nil {
			return ocr3types.Outcome{}, fmt.Errorf("decode commit plugin observation: %w", err)
		}
		decodedObservations = append(decodedObservations, obs)
	}

	maxSeqNums, err := maxSeqNumsConsensus(p.lggr, fChainDest, decodedObservations)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("max sequence numbers consensus: %w", err)
	}
	if len(maxSeqNums) == 0 {
		p.lggr.Warnw("no consensus on max sequence numbers")
		return ocr3types.Outcome{}, nil
	}
	p.lggr.Debugw("max sequence numbers consensus", "maxSeqNumsConsensus", maxSeqNums)

	merkleRoots, err := newMsgsConsensus(p.lggr, maxSeqNums, decodedObservations, p.cfg.FChain)
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("new messages consensus: %w", err)
	}
	p.lggr.Debugw("new messages consensus", "merkleRoots", merkleRoots)

	tokenPrices, err := tokenPricesConsensus(decodedObservations, p.cfg.FChain[p.cfg.DestChain])
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("token prices consensus: %w", err)
	}

	gasPrices, err := gasPricesConsensus(p.lggr, decodedObservations, p.cfg.FChain[p.cfg.DestChain])
	if err != nil {
		return ocr3types.Outcome{}, fmt.Errorf("gas prices consensus: %w", err)
	}

	return model.NewCommitPluginOutcome(maxSeqNums, merkleRoots, tokenPrices, gasPrices).Encode()
}

func (p *Plugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[[]byte], error) {
	outc, err := model.DecodeCommitPluginOutcome(outcome)
	if err != nil {
		return nil, fmt.Errorf("decode commit plugin outcome: %w", err)
	}

	// todo: include gas price updates

	/*
		todo: Once token/gas prices are implemented, we would want to probably check if outc.MerkleRoots is empty or not
		and only create a report if outc.MerkleRoots is non-empty OR gas/token price timer has expired
	*/

	rep := model.NewCommitPluginReport(outc.MerkleRoots, outc.TokenPrices, outc.GasPrices)

	encodedReport, err := p.reportCodec.Encode(context.Background(), rep)
	if err != nil {
		return nil, fmt.Errorf("encode commit plugin report: %w", err)
	}

	return []ocr3types.ReportWithInfo[[]byte]{{Report: encodedReport, Info: nil}}, nil
}

func (p *Plugin) ShouldAcceptAttestedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
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

func (p *Plugin) ShouldTransmitAcceptedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	if !p.cfg.Writer {
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

	// todo: if report is stale -> do not transmit
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

func (p *Plugin) knownSourceChainsSlice() []model.ChainSelector {
	knownSourceChainsSlice := p.knownSourceChains.ToSlice()
	sort.Slice(knownSourceChainsSlice, func(i, j int) bool { return knownSourceChainsSlice[i] < knownSourceChainsSlice[j] })
	return slicelib.Filter(knownSourceChainsSlice, func(ch model.ChainSelector) bool { return ch != p.cfg.DestChain })
}

// Interface compatibility checks.
var _ ocr3types.ReportingPlugin[[]byte] = &Plugin{}
