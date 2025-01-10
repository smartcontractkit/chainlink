package ccipcommit

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider"
	db "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdb"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

const (
	// only dynamic field in CommitReport is tokens PriceUpdates, and we don't expect to need to update thousands of tokens in a single tx
	MaxCommitReportLength = 10_000
	// Maximum inflight seq number range before we consider reports to be failing to get included entirely
	// and restart from the chain's minSeqNum. Want to set it high to allow for large throughput,
	// but low enough to minimize wasted revert cost.
	MaxInflightSeqNumGap = 500
	// OnRampMessagesScanLimit is used to limit number of onramp messages scanned in each Observation.
	// Single CommitRoot can contain up to merklemulti.MaxNumberTreeLeaves, so we scan twice that to be safe and still don't hurt DB performance.
	OnRampMessagesScanLimit = merklemulti.MaxNumberTreeLeaves * 2
)

var (
	_ types.ReportingPluginFactory = &CommitReportingPluginFactory{}
	_ types.ReportingPlugin        = &CommitReportingPlugin{}
)

type update struct {
	timestamp time.Time
	value     *big.Int
}

type CommitPluginStaticConfig struct {
	lggr                          logger.Logger
	newReportingPluginRetryConfig ccipdata.RetryConfig
	// Source
	onRampReader        ccipdata.OnRampReader
	sourceChainSelector uint64
	sourceNative        cciptypes.Address
	// Dest
	offRamp               ccipdata.OffRampReader
	commitStore           ccipdata.CommitStoreReader
	destChainSelector     uint64
	priceRegistryProvider ccipdataprovider.PriceRegistry
	// Offchain
	metricsCollector ccip.PluginMetricsCollector
	chainHealthcheck cache.ChainHealthcheck
	priceService     db.PriceService
}

type CommitReportingPlugin struct {
	lggr logger.Logger
	// Source
	onRampReader        ccipdata.OnRampReader
	sourceChainSelector uint64
	sourceNative        cciptypes.Address
	gasPriceEstimator   prices.GasPriceEstimatorCommit
	// Dest
	destChainSelector       uint64
	commitStoreReader       ccipdata.CommitStoreReader
	destPriceRegistryReader ccipdata.PriceRegistryReader
	offchainConfig          cciptypes.CommitOffchainConfig
	offRampReader           ccipdata.OffRampReader
	F                       int
	// Offchain
	metricsCollector ccip.PluginMetricsCollector
	// State
	chainHealthcheck cache.ChainHealthcheck
	// DB
	priceService db.PriceService
}

// Query is not used by the CCIP Commit plugin.
func (r *CommitReportingPlugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

// Observation calculates the sequence number interval ready to be committed and
// the token and gas price updates required. A valid report could contain a merkle
// root and price updates. Price updates should never contain nil values, otherwise
// the observation will be considered invalid and rejected.
func (r *CommitReportingPlugin) Observation(ctx context.Context, epochAndRound types.ReportTimestamp, _ types.Query) (types.Observation, error) {
	lggr := r.lggr.Named("CommitObservation")
	if healthy, err := r.chainHealthcheck.IsHealthy(ctx); err != nil {
		return nil, err
	} else if !healthy {
		return nil, ccip.ErrChainIsNotHealthy
	}

	// Will return 0,0 if no messages are found. This is a valid case as the report could
	// still contain fee updates.
	minSeqNr, maxSeqNr, messageIDs, err := r.calculateMinMaxSequenceNumbers(ctx, lggr)
	if err != nil {
		return nil, err
	}

	// Fetches multi-lane gasPricesUSD and tokenPricesUSD for the same dest chain
	gasPricesUSD, sourceGasPriceUSD, tokenPricesUSD, err := r.observePriceUpdates(ctx)
	if err != nil {
		return nil, err
	}

	lggr.Infow("Observation",
		"minSeqNr", minSeqNr,
		"maxSeqNr", maxSeqNr,
		"gasPricesUSD", gasPricesUSD,
		"tokenPricesUSD", tokenPricesUSD,
		"epochAndRound", epochAndRound,
		"messageIDs", messageIDs,
	)
	r.metricsCollector.NumberOfMessagesBasedOnInterval(ccip.Observation, minSeqNr, maxSeqNr)

	// Even if all values are empty we still want to communicate our observation
	// with the other nodes, therefore, we always return the observed values.
	return ccip.CommitObservation{
		Interval: cciptypes.CommitStoreInterval{
			Min: minSeqNr,
			Max: maxSeqNr,
		},
		TokenPricesUSD:            tokenPricesUSD,
		SourceGasPriceUSD:         sourceGasPriceUSD,
		SourceGasPriceUSDPerChain: gasPricesUSD,
	}.Marshal()
}

// observePriceUpdates fetches latest gas and token prices from DB as long as price reporting is not disabled.
// The prices are aggregated for all lanes for the same destination chain.
func (r *CommitReportingPlugin) observePriceUpdates(
	ctx context.Context,
) (gasPricesUSD map[uint64]*big.Int, sourceGasPriceUSD *big.Int, tokenPricesUSD map[cciptypes.Address]*big.Int, err error) {
	// Do not observe prices if price reporting is disabled. Price reporting will be disabled for lanes that are not leader lanes.
	if r.offchainConfig.PriceReportingDisabled {
		r.lggr.Infow("Price reporting disabled, skipping gas and token price reads")
		return map[uint64]*big.Int{}, nil, map[cciptypes.Address]*big.Int{}, nil
	}

	// Fetches multi-lane gas prices and token prices, for the given dest chain
	gasPricesUSD, tokenPricesUSD, err = r.priceService.GetGasAndTokenPrices(ctx, r.destChainSelector)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get prices from PriceService: %w", err)
	}

	// Set prices to empty maps if nil to be friendlier to JSON encoding
	if gasPricesUSD == nil {
		gasPricesUSD = map[uint64]*big.Int{}
	}
	if tokenPricesUSD == nil {
		tokenPricesUSD = map[cciptypes.Address]*big.Int{}
	}

	// For backwards compatibility with the older release during phased rollout, set the default gas price on this lane
	sourceGasPriceUSD = gasPricesUSD[r.sourceChainSelector]

	return gasPricesUSD, sourceGasPriceUSD, tokenPricesUSD, nil
}

func (r *CommitReportingPlugin) calculateMinMaxSequenceNumbers(ctx context.Context, lggr logger.Logger) (uint64, uint64, []cciptypes.Hash, error) {
	nextSeqNum, err := r.commitStoreReader.GetExpectedNextSequenceNumber(ctx)
	if err != nil {
		return 0, 0, []cciptypes.Hash{}, err
	}

	msgRequests, err := r.onRampReader.GetSendRequestsBetweenSeqNums(ctx, nextSeqNum, nextSeqNum+OnRampMessagesScanLimit, true)
	if err != nil {
		return 0, 0, []cciptypes.Hash{}, err
	}
	if len(msgRequests) == 0 {
		lggr.Infow("No new requests", "nextSeqNum", nextSeqNum)
		return 0, 0, []cciptypes.Hash{}, nil
	}

	messageIDs := make([]cciptypes.Hash, 0, len(msgRequests))
	seqNrs := make([]uint64, 0, len(msgRequests))
	for _, msgReq := range msgRequests {
		seqNrs = append(seqNrs, msgReq.SequenceNumber)
		messageIDs = append(messageIDs, msgReq.MessageID)
	}

	minSeqNr := seqNrs[0]
	maxSeqNr := seqNrs[len(seqNrs)-1]
	if minSeqNr != nextSeqNum {
		// Still report the observation as even partial reports have value e.g. all nodes are
		// missing a single, different log each, they would still be able to produce a valid report.
		lggr.Warnf("Missing sequence number range [%d-%d]", nextSeqNum, minSeqNr)
	}
	if !ccipcalc.ContiguousReqs(lggr, minSeqNr, maxSeqNr, seqNrs) {
		return 0, 0, []cciptypes.Hash{}, errors.New("unexpected gap in seq nums")
	}
	return minSeqNr, maxSeqNr, messageIDs, nil
}

// Gets the latest token price updates based on logs within the heartbeat
// The updates returned by this function are guaranteed to not contain nil values.
func (r *CommitReportingPlugin) getLatestTokenPriceUpdates(ctx context.Context, now time.Time) (map[cciptypes.Address]update, error) {
	tokenPriceUpdates, err := r.destPriceRegistryReader.GetTokenPriceUpdatesCreatedAfter(
		ctx,
		now.Add(-r.offchainConfig.TokenPriceHeartBeat),
		0,
	)
	if err != nil {
		return nil, err
	}

	latestUpdates := make(map[cciptypes.Address]update)
	for _, tokenUpdate := range tokenPriceUpdates {
		priceUpdate := tokenUpdate.TokenPriceUpdate
		// Ordered by ascending timestamps
		timestamp := time.Unix(priceUpdate.TimestampUnixSec.Int64(), 0)
		if priceUpdate.Value != nil && !timestamp.Before(latestUpdates[priceUpdate.Token].timestamp) {
			latestUpdates[priceUpdate.Token] = update{
				timestamp: timestamp,
				value:     priceUpdate.Value,
			}
		}
	}

	return latestUpdates, nil
}

// getLatestGasPriceUpdate returns the latest gas price updates based on logs within the heartbeat.
// If an update is found, it is not expected to contain a nil value.
func (r *CommitReportingPlugin) getLatestGasPriceUpdate(ctx context.Context, now time.Time) (map[uint64]update, error) {
	gasPriceUpdates, err := r.destPriceRegistryReader.GetAllGasPriceUpdatesCreatedAfter(
		ctx,
		now.Add(-r.offchainConfig.GasPriceHeartBeat),
		0,
	)

	if err != nil {
		return nil, err
	}

	latestUpdates := make(map[uint64]update)
	for _, gasUpdate := range gasPriceUpdates {
		priceUpdate := gasUpdate.GasPriceUpdate
		// Ordered by ascending timestamps
		timestamp := time.Unix(priceUpdate.TimestampUnixSec.Int64(), 0)
		if priceUpdate.Value != nil && !timestamp.Before(latestUpdates[priceUpdate.DestChainSelector].timestamp) {
			latestUpdates[priceUpdate.DestChainSelector] = update{
				timestamp: timestamp,
				value:     priceUpdate.Value,
			}
		}
	}

	r.lggr.Infow("Latest gas price from log poller", "latestUpdates", latestUpdates)
	return latestUpdates, nil
}

func (r *CommitReportingPlugin) Report(ctx context.Context, epochAndRound types.ReportTimestamp, _ types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	now := time.Now()
	lggr := r.lggr.Named("CommitReport")
	if healthy, err := r.chainHealthcheck.IsHealthy(ctx); err != nil {
		return false, nil, err
	} else if !healthy {
		return false, nil, ccip.ErrChainIsNotHealthy
	}

	parsableObservations := ccip.GetParsableObservations[ccip.CommitObservation](lggr, observations)

	intervals, gasPriceObs, tokenPriceObs, err := extractObservationData(lggr, r.F, r.sourceChainSelector, parsableObservations)
	if err != nil {
		return false, nil, err
	}

	agreedInterval, err := calculateIntervalConsensus(intervals, r.F, merklemulti.MaxNumberTreeLeaves)
	if err != nil {
		return false, nil, err
	}

	gasPrices, tokenPrices, err := r.selectPriceUpdates(ctx, now, gasPriceObs, tokenPriceObs)
	if err != nil {
		return false, nil, err
	}
	// If there are no fee updates and the interval is zero there is no report to produce.
	if agreedInterval.Max == 0 && len(gasPrices) == 0 && len(tokenPrices) == 0 {
		lggr.Infow("Empty report, skipping")
		return false, nil, nil
	}

	report, err := r.buildReport(ctx, lggr, agreedInterval, gasPrices, tokenPrices)
	if err != nil {
		return false, nil, err
	}
	encodedReport, err := r.commitStoreReader.EncodeCommitReport(ctx, report)
	if err != nil {
		return false, nil, err
	}
	r.metricsCollector.SequenceNumber(ccip.Report, report.Interval.Max)
	r.metricsCollector.NumberOfMessagesBasedOnInterval(ccip.Report, report.Interval.Min, report.Interval.Max)
	lggr.Infow("Report",
		"merkleRoot", hex.EncodeToString(report.MerkleRoot[:]),
		"minSeqNr", report.Interval.Min,
		"maxSeqNr", report.Interval.Max,
		"gasPriceUpdates", report.GasPrices,
		"tokenPriceUpdates", report.TokenPrices,
		"epochAndRound", epochAndRound,
	)
	return true, encodedReport, nil
}

// calculateIntervalConsensus compresses a set of intervals into one interval
// taking into account f which is the maximum number of faults across the whole DON.
// OCR itself won't call Report unless there are 2*f+1 observations
// https://github.com/smartcontractkit/libocr/blob/master/offchainreporting2/internal/protocol/report_generation_follower.go#L415
// and f of those observations may be either unparseable or adversarially set values. That means
// we'll either have f+1 parsed honest values here, 2f+1 parsed values with f adversarial values or somewhere
// in between.
// rangeLimit is the maximum range of the interval. If the interval is larger than this, it will be truncated. Zero means no limit.
func calculateIntervalConsensus(intervals []cciptypes.CommitStoreInterval, f int, rangeLimit uint64) (cciptypes.CommitStoreInterval, error) {
	// To understand min/max selection here, we need to consider an adversary that controls f values
	// and is intentionally trying to stall the protocol or influence the value returned. For simplicity
	// consider f=1 and n=4 nodes. In that case adversary may try to bias the min or max high/low.
	// We could end up (2f+1=3) with sorted_mins=[1,1,1e9] or [-1e9,1,1] as examples. Selecting
	// sorted_mins[f] ensures:
	// - At least one honest node has seen this value, so adversary cannot bias the value lower which
	// would cause reverts
	// - If an honest oracle reports sorted_min[f] which happens to be stale i.e. that oracle
	// has a delayed view of the chain, then the report will revert onchain but still succeed upon retry
	// - We minimize the risk of naturally hitting the error condition minSeqNum > maxSeqNum due to oracles
	// delayed views of the chain (would be an issue with taking sorted_mins[-f])
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Min < intervals[j].Min
	})
	minSeqNum := intervals[f].Min

	// The only way a report could have a minSeqNum of 0 is when there are no messages to report
	// and the report is potentially still valid for gas fee updates.
	if minSeqNum == 0 {
		return cciptypes.CommitStoreInterval{Min: 0, Max: 0}, nil
	}
	// Consider a similar example to the sorted_mins one above except where they are maxes.
	// We choose the more "conservative" sorted_maxes[f] so:
	// - We are ensured that at least one honest oracle has seen the max, so adversary cannot set it lower and
	// cause the maxSeqNum < minSeqNum errors
	// - If an honest oracle reports sorted_max[f] which happens to be stale i.e. that oracle
	// has a delayed view of the source chain, then we simply lose a little bit of throughput.
	// - If we were to pick sorted_max[-f] i.e. the maximum honest node view (a more "aggressive" setting in terms of throughput),
	// then an adversary can continually send high values e.g. imagine we have observations from all 4 nodes
	// [honest 1, honest 1, honest 2, malicious 2], in this case we pick 2, but it's not enough to be able
	// to build a report since the first 2 honest nodes are unaware of message 2.
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Max < intervals[j].Max
	})
	maxSeqNum := intervals[f].Max
	if maxSeqNum < minSeqNum {
		// If the consensus report is invalid for onchain acceptance, we do not vote for it as
		// an early termination step.
		return cciptypes.CommitStoreInterval{}, errors.New("max seq num smaller than min")
	}

	// If the range is too large, truncate it.
	if rangeLimit > 0 && maxSeqNum-minSeqNum+1 > rangeLimit {
		maxSeqNum = minSeqNum + rangeLimit - 1
	}

	return cciptypes.CommitStoreInterval{
		Min: minSeqNum,
		Max: maxSeqNum,
	}, nil
}

// extractObservationData extracts observation fields into their own slices
// and filters out observation data that are invalid
func extractObservationData(lggr logger.Logger, f int, sourceChainSelector uint64, observations []ccip.CommitObservation) (intervals []cciptypes.CommitStoreInterval, gasPrices map[uint64][]*big.Int, tokenPrices map[cciptypes.Address][]*big.Int, err error) {
	// We require at least f+1 observations to reach consensus. Checking to ensure there are at least f+1 parsed observations.
	if len(observations) <= f {
		return nil, nil, nil, fmt.Errorf("not enough observations to form consensus: #obs=%d, f=%d", len(observations), f)
	}

	gasPriceObservations := make(map[uint64][]*big.Int)
	tokenPriceObservations := make(map[cciptypes.Address][]*big.Int)
	for _, obs := range observations {
		intervals = append(intervals, obs.Interval)

		for selector, price := range obs.SourceGasPriceUSDPerChain {
			if price != nil {
				gasPriceObservations[selector] = append(gasPriceObservations[selector], price)
			}
		}
		// During phased rollout, NOPs running old release only report SourceGasPriceUSD.
		// An empty `SourceGasPriceUSDPerChain` with a non-nil `SourceGasPriceUSD` can only happen with old release.
		if len(obs.SourceGasPriceUSDPerChain) == 0 && obs.SourceGasPriceUSD != nil {
			gasPriceObservations[sourceChainSelector] = append(gasPriceObservations[sourceChainSelector], obs.SourceGasPriceUSD)
		}

		for token, price := range obs.TokenPricesUSD {
			if price != nil {
				tokenPriceObservations[token] = append(tokenPriceObservations[token], price)
			}
		}
	}

	// Price is dropped if there are not enough valid observations. With a threshold of 2*(f-1) + 1, we achieve a balance between safety and liveness.
	// During phased-rollout where some honest nodes may not have started observing the token yet, it requires 5 malicious node with 1 being the leader to successfully alter price.
	// During regular operation, it requires 3 malicious nodes with 1 being the leader to temporarily delay price update for the token.
	priceReportingThreshold := 2*(f-1) + 1

	gasPrices = make(map[uint64][]*big.Int)
	for selector, perChainPriceObservations := range gasPriceObservations {
		if len(perChainPriceObservations) < priceReportingThreshold {
			lggr.Warnf("Skipping chain with selector %d due to not enough valid observations: #obs=%d, f=%d, threshold=%d", selector, len(perChainPriceObservations), f, priceReportingThreshold)
			continue
		}
		gasPrices[selector] = perChainPriceObservations
	}

	tokenPrices = make(map[cciptypes.Address][]*big.Int)
	for token, perTokenPriceObservations := range tokenPriceObservations {
		if len(perTokenPriceObservations) < priceReportingThreshold {
			lggr.Warnf("Skipping token %s due to not enough valid observations: #obs=%d, f=%d, threshold=%d", string(token), len(perTokenPriceObservations), f, priceReportingThreshold)
			continue
		}
		tokenPrices[token] = perTokenPriceObservations
	}

	return intervals, gasPrices, tokenPrices, nil
}

// selectPriceUpdates filters out gas and token price updates that are already inflight
func (r *CommitReportingPlugin) selectPriceUpdates(ctx context.Context, now time.Time, gasPriceObs map[uint64][]*big.Int, tokenPriceObs map[cciptypes.Address][]*big.Int) ([]cciptypes.GasPrice, []cciptypes.TokenPrice, error) {
	// If price reporting is disabled, there is no need to select price updates.
	if r.offchainConfig.PriceReportingDisabled {
		return nil, nil, nil
	}

	latestGasPrice, err := r.getLatestGasPriceUpdate(ctx, now)
	if err != nil {
		return nil, nil, err
	}

	latestTokenPrices, err := r.getLatestTokenPriceUpdates(ctx, now)
	if err != nil {
		return nil, nil, err
	}

	return r.calculatePriceUpdates(ctx, gasPriceObs, tokenPriceObs, latestGasPrice, latestTokenPrices)
}

// Note priceUpdates must be deterministic.
// The provided gasPriceObs and tokenPriceObs should not contain nil values.
// The returned latestGasPrice and latestTokenPrices should not contain nil values.
func (r *CommitReportingPlugin) calculatePriceUpdates(ctx context.Context, gasPriceObs map[uint64][]*big.Int, tokenPriceObs map[cciptypes.Address][]*big.Int, latestGasPrice map[uint64]update, latestTokenPrices map[cciptypes.Address]update) ([]cciptypes.GasPrice, []cciptypes.TokenPrice, error) {
	var tokenPriceUpdates []cciptypes.TokenPrice
	// Token prices are mostly heartbeat driven. To maximize heartbeat batching, the price inclusion rule is as follows:
	// If any token requires heartbeat update, include all token prices in the report.
	// Otherwise, only include token prices that exceed deviation threshold.
	needTokenHeartbeat := false
	for token := range tokenPriceObs {
		latestTokenPrice, exists := latestTokenPrices[token]
		if !exists || time.Since(latestTokenPrice.timestamp) >= r.offchainConfig.TokenPriceHeartBeat {
			r.lggr.Infow("Token requires heartbeat update", "token", token)
			needTokenHeartbeat = true
			break
		}
	}

	for token, tokenPriceObservations := range tokenPriceObs {
		medianPrice := ccipcalc.BigIntSortedMiddle(tokenPriceObservations)

		if needTokenHeartbeat {
			r.lggr.Debugw("Token price update included due to heartbeat", "token", token, "newPrice", medianPrice)
			tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
				Token: token,
				Value: medianPrice,
			})
			continue
		}

		latestTokenPrice, exists := latestTokenPrices[token]
		if exists {
			if ccipcalc.Deviates(medianPrice, latestTokenPrice.value, int64(r.offchainConfig.TokenPriceDeviationPPB)) {
				r.lggr.Debugw("Token price update included due to deviation",
					"token", token, "newPrice", medianPrice, "existingPrice", latestTokenPrice.value)
				tokenPriceUpdates = append(tokenPriceUpdates, cciptypes.TokenPrice{
					Token: token,
					Value: medianPrice,
				})
			}
		}
	}

	// Determinism required.
	sort.Slice(tokenPriceUpdates, func(i, j int) bool {
		return tokenPriceUpdates[i].Token < tokenPriceUpdates[j].Token
	})

	var gasPriceUpdate []cciptypes.GasPrice
	// Gas prices are mostly heartbeat driven. To maximize heartbeat batching, the price inclusion rule is as follows:
	// If any source chain gas price requires heartbeat update, include all gas prices in the report.
	// Otherwise, only include gas prices that exceed deviation threshold.
	needGasHeartbeat := false
	for chainSelector := range gasPriceObs {
		latestGasPrice, exists := latestGasPrice[chainSelector]
		if !exists || latestGasPrice.value == nil || time.Since(latestGasPrice.timestamp) >= r.offchainConfig.GasPriceHeartBeat {
			r.lggr.Infow("Chain gas price requires heartbeat update", "chainSelector", chainSelector)
			needGasHeartbeat = true
			break
		}
	}

	for chainSelector, gasPriceObservations := range gasPriceObs {
		newGasPrice, err := r.gasPriceEstimator.Median(ctx, gasPriceObservations) // Compute the median price
		if err != nil {
			return nil, nil, fmt.Errorf("failed to calculate median gas price for chain selector %d: %w", chainSelector, err)
		}

		if needGasHeartbeat {
			r.lggr.Debugw("Gas price update included due to heartbeat", "chainSelector", chainSelector)
			gasPriceUpdate = append(gasPriceUpdate, cciptypes.GasPrice{
				DestChainSelector: chainSelector,
				Value:             newGasPrice,
			})
			continue
		}

		latestGasPrice, exists := latestGasPrice[chainSelector]
		if exists && latestGasPrice.value != nil {
			gasPriceDeviated, err := r.gasPriceEstimator.Deviates(ctx, newGasPrice, latestGasPrice.value)
			if err != nil {
				return nil, nil, err
			}
			if gasPriceDeviated {
				r.lggr.Debugw("Gas price update included due to deviation",
					"chainSelector", chainSelector, "newPrice", newGasPrice, "existingPrice", latestGasPrice.value)
				gasPriceUpdate = append(gasPriceUpdate, cciptypes.GasPrice{
					DestChainSelector: chainSelector,
					Value:             newGasPrice,
				})
			}
		}
	}

	sort.Slice(gasPriceUpdate, func(i, j int) bool {
		return gasPriceUpdate[i].DestChainSelector < gasPriceUpdate[j].DestChainSelector
	})

	return gasPriceUpdate, tokenPriceUpdates, nil
}

// buildReport assumes there is at least one message in reqs.
func (r *CommitReportingPlugin) buildReport(ctx context.Context, lggr logger.Logger, interval cciptypes.CommitStoreInterval, gasPrices []cciptypes.GasPrice, tokenPrices []cciptypes.TokenPrice) (cciptypes.CommitStoreReport, error) {
	// If no messages are needed only include fee updates
	if interval.Min == 0 {
		return cciptypes.CommitStoreReport{
			TokenPrices: tokenPrices,
			GasPrices:   gasPrices,
			MerkleRoot:  [32]byte{},
			Interval:    interval,
		}, nil
	}

	// Logs are guaranteed to be in order of seq num, since these are finalized logs only
	// and the contract's seq num is auto-incrementing.
	sendRequests, err := r.onRampReader.GetSendRequestsBetweenSeqNums(ctx, interval.Min, interval.Max, true)
	if err != nil {
		return cciptypes.CommitStoreReport{}, err
	}
	if len(sendRequests) == 0 {
		lggr.Warn("No messages found in interval",
			"minSeqNr", interval.Min,
			"maxSeqNr", interval.Max)
		return cciptypes.CommitStoreReport{}, fmt.Errorf("tried building a tree without leaves")
	}

	leaves := make([][32]byte, 0, len(sendRequests))
	var seqNrs []uint64
	for _, req := range sendRequests {
		leaves = append(leaves, req.Hash)
		seqNrs = append(seqNrs, req.SequenceNumber)
	}
	if !ccipcalc.ContiguousReqs(lggr, interval.Min, interval.Max, seqNrs) {
		return cciptypes.CommitStoreReport{}, errors.Errorf("do not have full range [%v, %v] have %v", interval.Min, interval.Max, seqNrs)
	}
	tree, err := merklemulti.NewTree(hashutil.NewKeccak(), leaves)
	if err != nil {
		return cciptypes.CommitStoreReport{}, err
	}

	return cciptypes.CommitStoreReport{
		GasPrices:   gasPrices,
		TokenPrices: tokenPrices,
		MerkleRoot:  tree.Root(),
		Interval:    interval,
	}, nil
}

func (r *CommitReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, reportTimestamp types.ReportTimestamp, report types.Report) (bool, error) {
	parsedReport, err := r.commitStoreReader.DecodeCommitReport(ctx, report)
	if err != nil {
		return false, err
	}
	lggr := r.lggr.Named("CommitShouldAcceptFinalizedReport").With(
		"merkleRoot", parsedReport.MerkleRoot,
		"minSeqNum", parsedReport.Interval.Min,
		"maxSeqNum", parsedReport.Interval.Max,
		"gasPriceUpdates", parsedReport.GasPrices,
		"tokenPriceUpdates", parsedReport.TokenPrices,
		"reportTimestamp", reportTimestamp,
	)
	// Empty report, should not be put on chain
	if parsedReport.MerkleRoot == [32]byte{} && len(parsedReport.GasPrices) == 0 && len(parsedReport.TokenPrices) == 0 {
		lggr.Warn("Empty report, should not be put on chain")
		return false, nil
	}

	if healthy, err1 := r.chainHealthcheck.IsHealthy(ctx); err1 != nil {
		return false, err1
	} else if !healthy {
		return false, ccip.ErrChainIsNotHealthy
	}

	if r.isStaleReport(ctx, lggr, parsedReport, reportTimestamp) {
		lggr.Infow("Rejecting stale report")
		return false, nil
	}

	r.metricsCollector.SequenceNumber(ccip.ShouldAccept, parsedReport.Interval.Max)
	lggr.Infow("Accepting finalized report", "merkleRoot", hexutil.Encode(parsedReport.MerkleRoot[:]))
	return true, nil
}

// ShouldTransmitAcceptedReport checks if the report is stale, if it is it should not be transmitted.
func (r *CommitReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, reportTimestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("CommitShouldTransmitAcceptedReport")
	parsedReport, err := r.commitStoreReader.DecodeCommitReport(ctx, report)
	if err != nil {
		return false, err
	}
	if healthy, err1 := r.chainHealthcheck.IsHealthy(ctx); err1 != nil {
		return false, err1
	} else if !healthy {
		return false, ccip.ErrChainIsNotHealthy
	}
	// If report is not stale we transmit.
	// When the commitTransmitter enqueues the tx for tx manager,
	// we mark it as fulfilled, effectively removing it from the set of inflight messages.
	shouldTransmit := !r.isStaleReport(ctx, lggr, parsedReport, reportTimestamp)

	lggr.Infow("ShouldTransmitAcceptedReport",
		"shouldTransmit", shouldTransmit,
		"reportTimestamp", reportTimestamp)
	return shouldTransmit, nil
}

// isStaleReport checks a report to see if the contents have become stale.
// It does so in four ways:
//  1. if there is a merkle root, check if the sequence numbers match up with onchain data
//  2. if there is no merkle root, check if current price's epoch and round is after onchain epoch and round
//  3. if there is a gas price update check to see if the value is different from the last
//     reported value
//  4. if there are token prices check to see if the values are different from the last
//     reported values.
//
// If there is a merkle root present, staleness is only measured based on the merkle root
// If there is no merkle root but there is a gas update, only this gas update is used for staleness checks.
// If only price updates are included, the price updates are used to check for staleness
// If nothing is included the report is always considered stale.
func (r *CommitReportingPlugin) isStaleReport(ctx context.Context, lggr logger.Logger, report cciptypes.CommitStoreReport, reportTimestamp types.ReportTimestamp) bool {
	// If there is a merkle root, ignore all other staleness checks and only check for sequence number staleness
	if report.MerkleRoot != [32]byte{} {
		return r.isStaleMerkleRoot(ctx, lggr, report.Interval)
	}

	hasGasPriceUpdate := len(report.GasPrices) > 0
	hasTokenPriceUpdates := len(report.TokenPrices) > 0

	// If there is no merkle root, no gas price update and no token price update
	// we don't want to write anything on-chain, so we consider this report stale.
	if !hasGasPriceUpdate && !hasTokenPriceUpdates {
		return true
	}

	// We consider a price update as stale when, there isn't an update or there is an update that is stale.
	gasPriceStale := !hasGasPriceUpdate || r.isStaleGasPrice(ctx, lggr, report.GasPrices)
	tokenPricesStale := !hasTokenPriceUpdates || r.isStaleTokenPrices(ctx, lggr, report.TokenPrices)

	if gasPriceStale && tokenPricesStale {
		return true
	}

	// If report only has price update, check if its epoch and round lags behind the latest onchain
	lastPriceEpochAndRound, err := r.commitStoreReader.GetLatestPriceEpochAndRound(ctx)
	if err != nil {
		// Assume it's a transient issue getting the last report and try again on the next round
		return true
	}

	thisEpochAndRound := ccipcalc.MergeEpochAndRound(reportTimestamp.Epoch, reportTimestamp.Round)
	return lastPriceEpochAndRound >= thisEpochAndRound
}

func (r *CommitReportingPlugin) isStaleMerkleRoot(ctx context.Context, lggr logger.Logger, reportInterval cciptypes.CommitStoreInterval) bool {
	nextSeqNum, err := r.commitStoreReader.GetExpectedNextSequenceNumber(ctx)
	if err != nil {
		// Assume it's a transient issue getting the last report and try again on the next round
		return true
	}

	// The report is not stale and correct only if nextSeqNum == reportInterval.Min.
	// Mark it stale if the condition isn't met.
	if nextSeqNum != reportInterval.Min {
		lggr.Infow("The report is stale because of sequence number mismatch with the commit store interval min value",
			"nextSeqNum", nextSeqNum, "reportIntervalMin", reportInterval.Min)
		return true
	}

	lggr.Infow("Report root is not stale", "nextSeqNum", nextSeqNum, "reportIntervalMin", reportInterval.Min)

	// If a report has root and valid sequence number, the report should be submitted, regardless of price staleness
	return false
}

func (r *CommitReportingPlugin) isStaleGasPrice(ctx context.Context, lggr logger.Logger, gasPriceUpdates []cciptypes.GasPrice) bool {
	latestGasPrice, err := r.getLatestGasPriceUpdate(ctx, time.Now())
	if err != nil {
		lggr.Errorw("Gas price is stale because getLatestGasPriceUpdate failed", "err", err)
		return true
	}

	for _, gasPriceUpdate := range gasPriceUpdates {
		latestUpdate, exists := latestGasPrice[gasPriceUpdate.DestChainSelector]
		if !exists || latestUpdate.value == nil {
			lggr.Infow("Found non-stale gas price", "chainSelector", gasPriceUpdate.DestChainSelector, "gasPriceUSd", gasPriceUpdate.Value)
			return false
		}

		gasPriceDeviated, err := r.gasPriceEstimator.Deviates(ctx, gasPriceUpdate.Value, latestUpdate.value)
		if err != nil {
			lggr.Errorw("Gas price is stale because deviation check failed", "err", err)
			return true
		}

		if gasPriceDeviated {
			lggr.Infow("Found non-stale gas price", "chainSelector", gasPriceUpdate.DestChainSelector, "gasPriceUSd", gasPriceUpdate.Value, "latestUpdate", latestUpdate.value)
			return false
		}
		lggr.Infow("Gas price is stale", "chainSelector", gasPriceUpdate.DestChainSelector, "gasPriceUSd", gasPriceUpdate.Value, "latestGasPrice", latestUpdate.value)
	}

	lggr.Infow("All gas prices are stale")
	return true
}

func (r *CommitReportingPlugin) isStaleTokenPrices(ctx context.Context, lggr logger.Logger, priceUpdates []cciptypes.TokenPrice) bool {
	// getting the last price updates without including inflight is like querying
	// current prices onchain, but uses logpoller's data to save on the RPC requests
	latestTokenPriceUpdates, err := r.getLatestTokenPriceUpdates(ctx, time.Now())
	if err != nil {
		return true
	}

	for _, tokenUpdate := range priceUpdates {
		latestUpdate, ok := latestTokenPriceUpdates[tokenUpdate.Token]
		priceEqual := ok && !ccipcalc.Deviates(tokenUpdate.Value, latestUpdate.value, int64(r.offchainConfig.TokenPriceDeviationPPB))

		if !priceEqual {
			lggr.Infow("Found non-stale token price", "token", tokenUpdate.Token, "usdPerToken", tokenUpdate.Value, "latestUpdate", latestUpdate.value)
			return false
		}
		lggr.Infow("Token price is stale", "latestTokenPrice", latestUpdate.value, "usdPerToken", tokenUpdate.Value, "token", tokenUpdate.Token)
	}

	lggr.Infow("All token prices are stale")
	return true
}

func (r *CommitReportingPlugin) Close() error {
	return nil
}
