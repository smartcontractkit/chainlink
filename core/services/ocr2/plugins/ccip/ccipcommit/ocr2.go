package ccipcommit

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
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
	lggr logger.Logger
	// Source
	onRampReader        ccipdata.OnRampReader
	sourceChainSelector uint64
	sourceNative        common.Address
	// Dest
	offRamp               ccipdata.OffRampReader
	commitStore           ccipdata.CommitStoreReader
	destChainSelector     uint64
	priceRegistryProvider ccipdataprovider.PriceRegistry
	// Offchain
	priceGetter      pricegetter.PriceGetter
	metricsCollector ccip.PluginMetricsCollector
}

type CommitReportingPlugin struct {
	lggr logger.Logger
	// Source
	onRampReader        ccipdata.OnRampReader
	sourceChainSelector uint64
	sourceNative        common.Address
	gasPriceEstimator   prices.GasPriceEstimatorCommit
	// Dest
	commitStoreReader       ccipdata.CommitStoreReader
	destPriceRegistryReader ccipdata.PriceRegistryReader
	offchainConfig          ccipdata.CommitOffchainConfig
	offRampReader           ccipdata.OffRampReader
	F                       int
	// Offchain
	priceGetter      pricegetter.PriceGetter
	metricsCollector ccip.PluginMetricsCollector
	// State
	inflightReports *inflightCommitReportsContainer
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
	// If the commit store is down the protocol should halt.
	down, err := r.commitStoreReader.IsDown(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "isDown check errored")
	}
	if down {
		return nil, ccip.ErrCommitStoreIsDown
	}
	r.inflightReports.expire(lggr)

	// Will return 0,0 if no messages are found. This is a valid case as the report could
	// still contain fee updates.
	minSeqNr, maxSeqNr, err := r.calculateMinMaxSequenceNumbers(ctx, lggr)
	if err != nil {
		return nil, err
	}

	feeTokens, bridgeableTokens, err := ccipcommon.GetDestinationTokens(ctx, r.offRampReader, r.destPriceRegistryReader)
	if err != nil {
		return nil, fmt.Errorf("get destination tokens: %w", err)
	}
	destTokens := ccipcommon.FlattenUniqueSlice(feeTokens, bridgeableTokens)

	sourceGasPriceUSD, tokenPricesUSD, err := r.generatePriceUpdates(ctx, lggr, destTokens)
	if err != nil {
		return nil, err
	}

	lggr.Infow("Observation",
		"minSeqNr", minSeqNr,
		"maxSeqNr", maxSeqNr,
		"sourceGasPriceUSD", sourceGasPriceUSD,
		"tokenPricesUSD", tokenPricesUSD,
		"epochAndRound", epochAndRound)
	r.metricsCollector.NumberOfMessagesBasedOnInterval(ccip.Observation, minSeqNr, maxSeqNr)

	// Even if all values are empty we still want to communicate our observation
	// with the other nodes, therefore, we always return the observed values.
	return ccip.CommitObservation{
		Interval: ccipdata.CommitStoreInterval{
			Min: minSeqNr,
			Max: maxSeqNr,
		},
		TokenPricesUSD:    tokenPricesUSD,
		SourceGasPriceUSD: sourceGasPriceUSD,
	}.Marshal()
}

func (r *CommitReportingPlugin) calculateMinMaxSequenceNumbers(ctx context.Context, lggr logger.Logger) (uint64, uint64, error) {
	nextInflightMin, _, err := r.nextMinSeqNum(ctx, lggr)
	if err != nil {
		return 0, 0, err
	}

	msgRequests, err := r.onRampReader.GetSendRequestsBetweenSeqNums(ctx, nextInflightMin, nextInflightMin+OnRampMessagesScanLimit, true)
	if err != nil {
		return 0, 0, err
	}
	if len(msgRequests) == 0 {
		lggr.Infow("No new requests", "minSeqNr", nextInflightMin)
		return 0, 0, nil
	}
	seqNrs := make([]uint64, 0, len(msgRequests))
	for _, msgReq := range msgRequests {
		seqNrs = append(seqNrs, msgReq.Data.SequenceNumber)
	}

	minSeqNr := seqNrs[0]
	maxSeqNr := seqNrs[len(seqNrs)-1]
	if minSeqNr != nextInflightMin {
		// Still report the observation as even partial reports have value e.g. all nodes are
		// missing a single, different log each, they would still be able to produce a valid report.
		lggr.Warnf("Missing sequence number range [%d-%d]", nextInflightMin, minSeqNr)
	}
	if !ccipcalc.ContiguousReqs(lggr, minSeqNr, maxSeqNr, seqNrs) {
		return 0, 0, errors.New("unexpected gap in seq nums")
	}
	return minSeqNr, maxSeqNr, nil
}

func (r *CommitReportingPlugin) nextMinSeqNum(ctx context.Context, lggr logger.Logger) (inflightMin, onChainMin uint64, err error) {
	nextMinOnChain, err := r.commitStoreReader.GetExpectedNextSequenceNumber(ctx)
	if err != nil {
		return 0, 0, err
	}
	// There are several scenarios to consider here for nextMin and inflight intervals.
	// 1. nextMin=2, inflight=[[2,3],[4,5]]. Node is waiting for [2,3] and [4,5] to be included, should return 6 to build on top.
	// 2. nextMin=2, inflight=[[4,5]]. [2,3] is expired but not yet visible onchain (means our cache expiry
	// was too low). In this case still want to return 6.
	// 3. nextMin=2, inflight=[] but other nodes have inflight=[2,3]. Say our node restarted and lost its cache. In this case
	// we still return the chain's nextMin, other oracles will ignore our observation. Other nodes however will build [4,5]
	// and then we'll add that to our cache in ShouldAcceptFinalizedReport, putting us into the previous position at which point
	// we can start contributing again.
	// 4. nextMin=4, inflight=[[2,3],[4,5]]. We see the onchain update, but haven't expired from our cache yet. Should happen
	// regularly and we just return 6.
	// 5. nextMin=2, inflight=[[4,5]]. [2,3] failed to get onchain for some reason. We'll return 6 and continue building even though
	// subsequent reports will revert, but eventually they will all expire OR we'll hit MaxInflightSeqNumGap and forcibly
	// expire them all. This scenario can also occur if there is a reorg which reorders the reports such that one reverts.
	maxInflight := r.inflightReports.maxInflightSeqNr()
	if (maxInflight > nextMinOnChain) && ((maxInflight - nextMinOnChain) > MaxInflightSeqNumGap) {
		r.inflightReports.reset(lggr)
		return nextMinOnChain, nextMinOnChain, nil
	}
	return mathutil.Max(nextMinOnChain, maxInflight+1), nextMinOnChain, nil
}

// All prices are USD ($1=1e18) denominated. All prices must be not nil.
// Return token prices should contain the exact same tokens as in tokenDecimals.
func (r *CommitReportingPlugin) generatePriceUpdates(
	ctx context.Context,
	lggr logger.Logger,
	destTokens []common.Address,
) (sourceGasPriceUSD prices.GasPrice, tokenPricesUSD map[common.Address]*big.Int, err error) {
	// Include wrapped native in our token query as way to identify the source native USD price.
	// notice USD is in 1e18 scale, i.e. $1 = 1e18
	queryTokens := ccipcommon.FlattenUniqueSlice([]common.Address{r.sourceNative}, destTokens)
	sort.Slice(queryTokens, func(i, j int) bool { return queryTokens[i].String() < queryTokens[j].String() }) // make the query deterministic

	rawTokenPricesUSD, err := r.priceGetter.TokenPricesUSD(ctx, queryTokens)
	if err != nil {
		return nil, nil, err
	}
	lggr.Infow("Raw token prices", "rawTokenPrices", rawTokenPricesUSD)

	// make sure that we got prices for all the tokens of our query
	for _, token := range queryTokens {
		if rawTokenPricesUSD[token] == nil {
			return nil, nil, errors.Errorf("missing token price: %+v", token)
		}
	}

	sourceNativePriceUSD, exists := rawTokenPricesUSD[r.sourceNative]
	if !exists {
		return nil, nil, fmt.Errorf("missing source native (%s) price", r.sourceNative)
	}

	destTokensDecimals, err := r.destPriceRegistryReader.GetTokensDecimals(ctx, destTokens)
	if err != nil {
		return nil, nil, fmt.Errorf("get tokens decimals: %w", err)
	}

	tokenPricesUSD = make(map[common.Address]*big.Int, len(rawTokenPricesUSD))
	for i, token := range destTokens {
		tokenPricesUSD[token] = calculateUsdPer1e18TokenAmount(rawTokenPricesUSD[token], destTokensDecimals[i])
	}

	sourceGasPrice, err := r.gasPriceEstimator.GetGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}
	if sourceGasPrice == nil {
		return nil, nil, errors.Errorf("missing gas price")
	}
	sourceGasPriceUSD, err = r.gasPriceEstimator.DenoteInUSD(sourceGasPrice, sourceNativePriceUSD)
	if err != nil {
		return nil, nil, err
	}

	lggr.Infow("Observing gas price", "observedGasPriceWei", sourceGasPrice, "observedGasPriceUSD", sourceGasPriceUSD)
	lggr.Infow("Observing token prices", "tokenPrices", tokenPricesUSD, "sourceNativePriceUSD", sourceNativePriceUSD)
	return sourceGasPriceUSD, tokenPricesUSD, nil
}

// Input price is USD per full token, with 18 decimal precision
// Result price is USD per 1e18 of smallest token denomination, with 18 decimal precision
// Example: 1 USDC = 1.00 USD per full token, each full token is 6 decimals -> 1 * 1e18 * 1e18 / 1e6 = 1e30
func calculateUsdPer1e18TokenAmount(price *big.Int, decimals uint8) *big.Int {
	tmp := big.NewInt(0).Mul(price, big.NewInt(1e18))
	return tmp.Div(tmp, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
}

// Gets the latest token price updates based on logs within the heartbeat
// The updates returned by this function are guaranteed to not contain nil values.
func (r *CommitReportingPlugin) getLatestTokenPriceUpdates(ctx context.Context, now time.Time, checkInflight bool) (map[common.Address]update, error) {
	tokenPriceUpdates, err := r.destPriceRegistryReader.GetTokenPriceUpdatesCreatedAfter(
		ctx,
		now.Add(-r.offchainConfig.TokenPriceHeartBeat),
		0,
	)
	if err != nil {
		return nil, err
	}

	latestUpdates := make(map[common.Address]update)
	for _, tokenUpdate := range tokenPriceUpdates {
		priceUpdate := tokenUpdate.Data
		// Ordered by ascending timestamps
		timestamp := time.Unix(priceUpdate.TimestampUnixSec.Int64(), 0)
		if priceUpdate.Value != nil && !timestamp.Before(latestUpdates[priceUpdate.Token].timestamp) {
			latestUpdates[priceUpdate.Token] = update{
				timestamp: timestamp,
				value:     priceUpdate.Value,
			}
		}
	}
	if !checkInflight {
		return latestUpdates, nil
	}

	// todo this comparison is faulty, as a previously-sent update's onchain timestamp can be higher than inflight timestamp
	// to properly fix, need a solution to map from onchain request to offchain timestamp
	// leaving it as is, as token prices are updated infrequently, so this should not cause many issues
	latestInflightTokenPriceUpdates := r.inflightReports.latestInflightTokenPriceUpdates()
	for inflightToken, latestInflightUpdate := range latestInflightTokenPriceUpdates {
		if latestInflightUpdate.timestamp.After(latestUpdates[inflightToken].timestamp) && latestInflightUpdate.value != nil {
			latestUpdates[inflightToken] = latestInflightUpdate
		}
	}

	return latestUpdates, nil
}

// getLatestGasPriceUpdate returns the latest gas price update based on logs within the heartbeat.
// If an update is found, it is not expected to contain a nil value. If no updates found, empty update with nil value is returned.
func (r *CommitReportingPlugin) getLatestGasPriceUpdate(ctx context.Context, now time.Time, checkInflight bool) (gasUpdate update, error error) {
	if checkInflight {
		latestInflightGasPriceUpdates := r.inflightReports.latestInflightGasPriceUpdates()
		if inflightUpdate, exists := latestInflightGasPriceUpdates[r.sourceChainSelector]; exists {
			gasUpdate = inflightUpdate
			r.lggr.Infow("Latest gas price from inflight", "gasPriceUpdateVal", gasUpdate.value, "gasPriceUpdateTs", gasUpdate.timestamp)

			// Gas price can fluctuate frequently, many updates may be in flight.
			// If there is gas price update inflight, use it as source of truth, no need to check onchain.
			return gasUpdate, nil
		}
	}

	// If there are no price updates inflight, check latest prices onchain
	gasPriceUpdates, err := r.destPriceRegistryReader.GetGasPriceUpdatesCreatedAfter(
		ctx,
		r.sourceChainSelector,
		now.Add(-r.offchainConfig.GasPriceHeartBeat),
		0,
	)
	if err != nil {
		return update{}, err
	}

	for _, priceUpdate := range gasPriceUpdates {
		// Ordered by ascending timestamps
		timestamp := time.Unix(priceUpdate.Data.TimestampUnixSec.Int64(), 0)
		if !timestamp.Before(gasUpdate.timestamp) {
			gasUpdate = update{
				timestamp: timestamp,
				value:     priceUpdate.Data.Value,
			}
		}
	}

	r.lggr.Infow("Latest gas price from log poller", "gasPriceUpdateVal", gasUpdate.value, "gasPriceUpdateTs", gasUpdate.timestamp)
	return gasUpdate, nil
}

func (r *CommitReportingPlugin) Report(ctx context.Context, epochAndRound types.ReportTimestamp, _ types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	now := time.Now()
	lggr := r.lggr.Named("CommitReport")
	parsableObservations := ccip.GetParsableObservations[ccip.CommitObservation](lggr, observations)

	feeTokens, bridgeableTokens, err := ccipcommon.GetDestinationTokens(ctx, r.offRampReader, r.destPriceRegistryReader)
	if err != nil {
		return false, nil, fmt.Errorf("get destination tokens: %w", err)
	}
	destTokens := ccipcommon.FlattenUniqueSlice(feeTokens, bridgeableTokens)

	// Filters out parsable but faulty observations
	validObservations, err := validateObservations(ctx, lggr, destTokens, r.F, parsableObservations)
	if err != nil {
		return false, nil, err
	}

	var intervals []ccipdata.CommitStoreInterval
	for _, obs := range validObservations {
		intervals = append(intervals, obs.Interval)
	}

	agreedInterval, err := calculateIntervalConsensus(intervals, r.F, merklemulti.MaxNumberTreeLeaves)
	if err != nil {
		return false, nil, err
	}

	latestGasPrice, err := r.getLatestGasPriceUpdate(ctx, now, true)
	if err != nil {
		return false, nil, err
	}

	latestTokenPrices, err := r.getLatestTokenPriceUpdates(ctx, now, true)
	if err != nil {
		return false, nil, err
	}

	tokenPrices, gasPrices, err := r.calculatePriceUpdates(validObservations, latestGasPrice, latestTokenPrices)
	if err != nil {
		return false, nil, err
	}
	// If there are no fee updates and the interval is zero there is no report to produce.
	if len(tokenPrices) == 0 && len(gasPrices) == 0 && agreedInterval.Max == 0 {
		lggr.Infow("Empty report, skipping")
		return false, nil, nil
	}

	report, err := r.buildReport(ctx, lggr, agreedInterval, gasPrices, tokenPrices)
	if err != nil {
		return false, nil, err
	}
	encodedReport, err := r.commitStoreReader.EncodeCommitReport(report)
	if err != nil {
		return false, nil, err
	}
	r.metricsCollector.NumberOfMessagesBasedOnInterval(ccip.Report, report.Interval.Min, report.Interval.Max)
	lggr.Infow("Report",
		"merkleRoot", hex.EncodeToString(report.MerkleRoot[:]),
		"minSeqNr", report.Interval.Min,
		"maxSeqNr", report.Interval.Max,
		"tokenPriceUpdates", report.TokenPrices,
		"gasPriceUpdates", report.GasPrices,
		"epochAndRound", epochAndRound,
	)
	return true, encodedReport, nil
}

// validateObservations validates the given observations.
// An observation is rejected if any of its gas price or token price is nil. With current CommitObservation implementation, prices
// are checked to ensure no nil values before adding to Observation, hence an observation that contains nil values comes from a faulty node.
func validateObservations(ctx context.Context, lggr logger.Logger, destTokens []common.Address, f int, observations []ccip.CommitObservation) (validObs []ccip.CommitObservation, err error) {
	for _, obs := range observations {
		// If gas price is reported as nil, the observation is faulty, skip the observation.
		if obs.SourceGasPriceUSD == nil {
			lggr.Warnw("Skipping observation due to nil SourceGasPriceUSD")
			continue
		}

		// If observed number of token prices does not match number of supported tokens on dest chain, skip the observation.
		if len(destTokens) != len(obs.TokenPricesUSD) {
			lggr.Warnw("Skipping observation due to token count mismatch", "expecting", len(destTokens), "got", len(obs.TokenPricesUSD))
			continue
		}
		// If any of the observed token prices is reported as nil, or not supported on dest chain, the observation is faulty, skip the observation.
		// Printing all faulty prices instead of short-circuiting to make log more informative.
		skipObservation := false
		for token, price := range obs.TokenPricesUSD {
			if price == nil {
				lggr.Warnw("Nil value in observed TokenPricesUSD", "token", token.Hex())
				skipObservation = true
			}
			if !slices.Contains(destTokens, token) {
				lggr.Warnw("Unsupported token in observed TokenPricesUSD", "token", token.Hex())
				skipObservation = true
			}
		}
		if skipObservation {
			lggr.Warnw("Skipping observation due to invalid TokenPricesUSD")
			continue
		}

		validObs = append(validObs, obs)
	}

	// We require at least f+1 valid observations. This corresponds to the scenario where f of the 2f+1 are faulty.
	if len(validObs) <= f {
		return nil, errors.Errorf("Not enough valid observations to form consensus: #obs=%d, f=%d", len(validObs), f)
	}

	return validObs, nil
}

// calculateIntervalConsensus compresses a set of intervals into one interval
// taking into account f which is the maximum number of faults across the whole DON.
// OCR itself won't call Report unless there are 2*f+1 observations
// https://github.com/smartcontractkit/libocr/blob/master/offchainreporting2/internal/protocol/report_generation_follower.go#L415
// and f of those observations may be either unparseable or adversarially set values. That means
// we'll either have f+1 parsed honest values here, 2f+1 parsed values with f adversarial values or somewhere
// in between.
// rangeLimit is the maximum range of the interval. If the interval is larger than this, it will be truncated. Zero means no limit.
func calculateIntervalConsensus(intervals []ccipdata.CommitStoreInterval, f int, rangeLimit uint64) (ccipdata.CommitStoreInterval, error) {
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
		return ccipdata.CommitStoreInterval{Min: 0, Max: 0}, nil
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
		return ccipdata.CommitStoreInterval{}, errors.New("max seq num smaller than min")
	}

	// If the range is too large, truncate it.
	if rangeLimit > 0 && maxSeqNum-minSeqNum+1 > rangeLimit {
		maxSeqNum = minSeqNum + rangeLimit - 1
	}

	return ccipdata.CommitStoreInterval{
		Min: minSeqNum,
		Max: maxSeqNum,
	}, nil
}

// Note priceUpdates must be deterministic.
// The provided latestTokenPrices should not contain nil values.
func (r *CommitReportingPlugin) calculatePriceUpdates(observations []ccip.CommitObservation, latestGasPrice update, latestTokenPrices map[common.Address]update) ([]ccipdata.TokenPrice, []ccipdata.GasPrice, error) {
	priceObservations := make(map[common.Address][]*big.Int)
	var sourceGasObservations []prices.GasPrice

	for _, obs := range observations {
		sourceGasObservations = append(sourceGasObservations, obs.SourceGasPriceUSD)
		// iterate over any token which price is included in observations
		for token, price := range obs.TokenPricesUSD {
			priceObservations[token] = append(priceObservations[token], price)
		}
	}

	var tokenPriceUpdates []ccipdata.TokenPrice
	for token, tokenPriceObservations := range priceObservations {
		medianPrice := ccipcalc.BigIntSortedMiddle(tokenPriceObservations)

		latestTokenPrice, exists := latestTokenPrices[token]
		if exists {
			tokenPriceUpdatedRecently := time.Since(latestTokenPrice.timestamp) < r.offchainConfig.TokenPriceHeartBeat
			tokenPriceNotChanged := !ccipcalc.Deviates(medianPrice, latestTokenPrice.value, int64(r.offchainConfig.TokenPriceDeviationPPB))
			if tokenPriceUpdatedRecently && tokenPriceNotChanged {
				r.lggr.Debugw("price was updated recently, skipping the update",
					"token", token, "newPrice", medianPrice, "existingPrice", latestTokenPrice.value)
				continue // skip the update if we recently had a price update close to the new value
			}
		}

		tokenPriceUpdates = append(tokenPriceUpdates, ccipdata.TokenPrice{
			Token: token,
			Value: medianPrice,
		})
	}

	// Determinism required.
	sort.Slice(tokenPriceUpdates, func(i, j int) bool {
		return bytes.Compare(tokenPriceUpdates[i].Token[:], tokenPriceUpdates[j].Token[:]) == -1
	})

	newGasPrice, err := r.gasPriceEstimator.Median(sourceGasObservations) // Compute the median price
	if err != nil {
		return nil, nil, err
	}
	destChainSelector := r.sourceChainSelector // Assuming plugin lane is A->B, we write to B the gas price of A

	var gasPrices []ccipdata.GasPrice
	// Default to updating so that we update if there are no prior updates.
	shouldUpdate := true
	if latestGasPrice.value != nil {
		gasPriceUpdatedRecently := time.Since(latestGasPrice.timestamp) < r.offchainConfig.GasPriceHeartBeat
		gasPriceDeviated, err := r.gasPriceEstimator.Deviates(newGasPrice, latestGasPrice.value)
		if err != nil {
			return nil, nil, err
		}
		if gasPriceUpdatedRecently && !gasPriceDeviated {
			shouldUpdate = false
		}
	}
	if shouldUpdate {
		// Although onchain interface accepts multi gas updates, we only do 1 gas price per report for now.
		gasPrices = append(gasPrices, ccipdata.GasPrice{DestChainSelector: destChainSelector, Value: newGasPrice})
	}

	return tokenPriceUpdates, gasPrices, nil
}

// buildReport assumes there is at least one message in reqs.
func (r *CommitReportingPlugin) buildReport(ctx context.Context, lggr logger.Logger, interval ccipdata.CommitStoreInterval, gasPrices []ccipdata.GasPrice, tokenPrices []ccipdata.TokenPrice) (ccipdata.CommitStoreReport, error) {
	// If no messages are needed only include fee updates
	if interval.Min == 0 {
		return ccipdata.CommitStoreReport{
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
		return ccipdata.CommitStoreReport{}, err
	}
	if len(sendRequests) == 0 {
		lggr.Warn("No messages found in interval",
			"minSeqNr", interval.Min,
			"maxSeqNr", interval.Max)
		return ccipdata.CommitStoreReport{}, fmt.Errorf("tried building a tree without leaves")
	}

	leaves := make([][32]byte, 0, len(sendRequests))
	var seqNrs []uint64
	for _, req := range sendRequests {
		leaves = append(leaves, req.Data.Hash)
		seqNrs = append(seqNrs, req.Data.SequenceNumber)
	}
	if !ccipcalc.ContiguousReqs(lggr, interval.Min, interval.Max, seqNrs) {
		return ccipdata.CommitStoreReport{}, errors.Errorf("do not have full range [%v, %v] have %v", interval.Min, interval.Max, seqNrs)
	}
	tree, err := merklemulti.NewTree(hashlib.NewKeccakCtx(), leaves)
	if err != nil {
		return ccipdata.CommitStoreReport{}, err
	}

	return ccipdata.CommitStoreReport{
		GasPrices:   gasPrices,
		TokenPrices: tokenPrices,
		MerkleRoot:  tree.Root(),
		Interval:    interval,
	}, nil
}

func (r *CommitReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, reportTimestamp types.ReportTimestamp, report types.Report) (bool, error) {
	parsedReport, err := r.commitStoreReader.DecodeCommitReport(report)
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

	if r.isStaleReport(ctx, lggr, parsedReport, true, reportTimestamp) {
		lggr.Infow("Rejecting stale report")
		return false, nil
	}

	epochAndRound := ccipcalc.MergeEpochAndRound(reportTimestamp.Epoch, reportTimestamp.Round)
	if err := r.inflightReports.add(lggr, parsedReport, epochAndRound); err != nil {
		return false, err
	}
	r.metricsCollector.SequenceNumber(parsedReport.Interval.Max)
	lggr.Infow("Accepting finalized report", "merkleRoot", hexutil.Encode(parsedReport.MerkleRoot[:]))
	return true, nil
}

// ShouldTransmitAcceptedReport checks if the report is stale, if it is it should not be transmitted.
func (r *CommitReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, reportTimestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("CommitShouldTransmitAcceptedReport")
	parsedReport, err := r.commitStoreReader.DecodeCommitReport(report)
	if err != nil {
		return false, err
	}
	// If report is not stale we transmit.
	// When the commitTransmitter enqueues the tx for tx manager,
	// we mark it as fulfilled, effectively removing it from the set of inflight messages.
	shouldTransmit := !r.isStaleReport(ctx, lggr, parsedReport, false, reportTimestamp)

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
func (r *CommitReportingPlugin) isStaleReport(ctx context.Context, lggr logger.Logger, report ccipdata.CommitStoreReport, checkInflight bool, reportTimestamp types.ReportTimestamp) bool {
	// If there is a merkle root, ignore all other staleness checks and only check for sequence number staleness
	if report.MerkleRoot != [32]byte{} {
		return r.isStaleMerkleRoot(ctx, lggr, report.Interval, checkInflight)
	}

	hasGasPriceUpdate := len(report.GasPrices) > 0
	hasTokenPriceUpdates := len(report.TokenPrices) > 0

	// If there is no merkle root, no gas price update and no token price update
	// we don't want to write anything on-chain, so we consider this report stale.
	if !hasGasPriceUpdate && !hasTokenPriceUpdates {
		return true
	}
	// Commit plugin currently only supports 1 gas price per report. If report contains more than 1, reject the report.
	if len(report.GasPrices) > 1 {
		lggr.Errorw("Report is stale because it contains more than 1 gas price update", "GasPriceUpdates", report.GasPrices)
		return true
	}

	// We consider a price update as stale when, there isn't an update or there is an update that is stale.
	gasPriceStale := !hasGasPriceUpdate || r.isStaleGasPrice(ctx, lggr, report.GasPrices[0], checkInflight)
	tokenPricesStale := !hasTokenPriceUpdates || r.isStaleTokenPrices(ctx, lggr, report.TokenPrices, checkInflight)

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

func (r *CommitReportingPlugin) isStaleMerkleRoot(ctx context.Context, lggr logger.Logger, reportInterval ccipdata.CommitStoreInterval, checkInflight bool) bool {
	nextInflightMin, nextOnChainMin, err := r.nextMinSeqNum(ctx, lggr)
	if err != nil {
		// Assume it's a transient issue getting the last report and try again on the next round
		return true
	}

	if checkInflight && nextInflightMin != reportInterval.Min {
		// There are sequence numbers missing between the commitStoreReader/inflight txs and the proposed report.
		// The report will fail onchain unless the inflight cache is in an incorrect state. A state like this
		// could happen for various reasons, e.g. a reboot of the node emptying the caches, and should be self-healing.
		// We do not submit a tx and wait for the protocol to self-heal by updating the caches or invalidating
		// inflight caches over time.
		lggr.Errorw("Next inflight min is not equal to the proposed min of the report", "nextInflightMin", nextInflightMin)
		return true
	}

	if !checkInflight && nextOnChainMin > reportInterval.Min {
		// If the next min is already greater than this reports min, this report is stale.
		lggr.Infow("Report is stale because of root", "onchain min", nextOnChainMin, "report min", reportInterval.Min)
		return true
	}

	// If a report has root and valid sequence number, the report should be submitted, regardless of price staleness
	return false
}

func (r *CommitReportingPlugin) isStaleGasPrice(ctx context.Context, lggr logger.Logger, gasPrice ccipdata.GasPrice, checkInflight bool) bool {
	latestGasPrice, err := r.getLatestGasPriceUpdate(ctx, time.Now(), checkInflight)
	if err != nil {
		lggr.Errorw("Report is stale because getLatestGasPriceUpdate failed", "err", err)
		return true
	}

	if latestGasPrice.value != nil {
		gasPriceDeviated, err := r.gasPriceEstimator.Deviates(gasPrice.Value, latestGasPrice.value)
		if err != nil {
			lggr.Errorw("Report is stale because deviation check failed", "err", err)
			return true
		}

		if !gasPriceDeviated {
			lggr.Infow("Report is stale because of gas price",
				"latestGasPriceUpdate", latestGasPrice.value,
				"currentUsdPerUnitGas", gasPrice.Value,
				"destChainSelector", gasPrice.DestChainSelector)
			return true
		}
	}

	return false
}

func (r *CommitReportingPlugin) isStaleTokenPrices(ctx context.Context, lggr logger.Logger, priceUpdates []ccipdata.TokenPrice, checkInflight bool) bool {
	// getting the last price updates without including inflight is like querying
	// current prices onchain, but uses logpoller's data to save on the RPC requests
	latestTokenPriceUpdates, err := r.getLatestTokenPriceUpdates(ctx, time.Now(), checkInflight)
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
