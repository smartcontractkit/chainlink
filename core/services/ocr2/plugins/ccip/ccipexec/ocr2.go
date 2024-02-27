package ccipexec

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/batchreader"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
)

const (
	// exec Report should make sure to cap returned payload to this limit
	MaxExecutionReportLength = 250_000

	// MaxDataLenPerBatch limits the total length of msg data that can be in a batch.
	MaxDataLenPerBatch = 60_000

	// MaximumAllowedTokenDataWaitTimePerBatch defines the maximum time that is allowed
	// for the plugin to wait for token data to be fetched from external providers per batch.
	MaximumAllowedTokenDataWaitTimePerBatch = 2 * time.Second

	// MessagesIterationStep limits number of messages fetched to memory at once when iterating through unexpired CommitRoots
	MessagesIterationStep = 800
)

var (
	_ types.ReportingPluginFactory = &ExecutionReportingPluginFactory{}
	_ types.ReportingPlugin        = &ExecutionReportingPlugin{}
)

type ExecutionPluginStaticConfig struct {
	lggr                     logger.Logger
	onRampReader             ccipdata.OnRampReader
	offRampReader            ccipdata.OffRampReader
	commitStoreReader        ccipdata.CommitStoreReader
	sourcePriceRegistry      ccipdata.PriceRegistryReader
	sourceWrappedNativeToken cciptypes.Address
	tokenDataWorker          tokendata.Worker
	destChainSelector        uint64
	priceRegistryProvider    ccipdataprovider.PriceRegistry
	tokenPoolBatchedReader   batchreader.TokenPoolBatchedReader
	metricsCollector         ccip.PluginMetricsCollector
}

type ExecutionReportingPlugin struct {
	// Misc
	F                int
	lggr             logger.Logger
	offchainConfig   cciptypes.ExecOffchainConfig
	tokenDataWorker  tokendata.Worker
	metricsCollector ccip.PluginMetricsCollector
	// Source
	gasPriceEstimator        prices.GasPriceEstimatorExec
	sourcePriceRegistry      ccipdata.PriceRegistryReader
	sourceWrappedNativeToken cciptypes.Address
	onRampReader             ccipdata.OnRampReader
	// Dest

	commitStoreReader      ccipdata.CommitStoreReader
	destPriceRegistry      ccipdata.PriceRegistryReader
	destWrappedNative      cciptypes.Address
	onchainConfig          cciptypes.ExecOnchainConfig
	offRampReader          ccipdata.OffRampReader
	tokenPoolBatchedReader batchreader.TokenPoolBatchedReader

	// State
	inflightReports *inflightExecReportsContainer
	snoozedRoots    cache.SnoozedRoots
}

func (r *ExecutionReportingPlugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

func (r *ExecutionReportingPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	lggr := r.lggr.Named("ExecutionObservation")
	down, err := r.commitStoreReader.IsDown(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "isDown check errored")
	}
	if down {
		return nil, ccip.ErrCommitStoreIsDown
	}
	// Expire any inflight reports.
	r.inflightReports.expire(lggr)
	inFlight := r.inflightReports.getAll()

	executableObservations, err := r.getExecutableObservations(ctx, lggr, timestamp, inFlight)
	if err != nil {
		return nil, err
	}
	// cap observations which fits MaxObservationLength (after serialized)
	capped := sort.Search(len(executableObservations), func(i int) bool {
		var encoded []byte
		encoded, err = ccip.NewExecutionObservation(executableObservations[:i+1]).Marshal()
		if err != nil {
			// false makes Search keep looking to the right, always including any "erroring" ObservedMessage and allowing us to detect in the bottom
			return false
		}
		return len(encoded) > ccip.MaxObservationLength
	})
	if err != nil {
		return nil, err
	}
	executableObservations = executableObservations[:capped]
	r.metricsCollector.NumberOfMessagesProcessed(ccip.Observation, len(executableObservations))
	lggr.Infow("Observation", "executableMessages", executableObservations)
	// Note can be empty
	return ccip.NewExecutionObservation(executableObservations).Marshal()
}

func (r *ExecutionReportingPlugin) getExecutableObservations(ctx context.Context, lggr logger.Logger, timestamp types.ReportTimestamp, inflight []InflightInternalExecutionReport) ([]ccip.ObservedMessage, error) {
	unexpiredReports, err := r.getUnexpiredCommitReports(
		ctx,
		r.commitStoreReader,
		r.onchainConfig.PermissionLessExecutionThresholdSeconds,
		lggr,
	)
	if err != nil {
		return nil, err
	}

	if len(unexpiredReports) == 0 {
		return []ccip.ObservedMessage{}, nil
	}

	getExecTokenData := cache.LazyFunction[execTokenData](func() (execTokenData, error) {
		return r.prepareTokenExecData(ctx)
	})

	for j := 0; j < len(unexpiredReports); {
		unexpiredReportsPart, step := selectReportsToFillBatch(unexpiredReports[j:], MessagesIterationStep)
		j += step

		unexpiredReportsWithSendReqs, err := r.getReportsWithSendRequests(ctx, unexpiredReportsPart)
		if err != nil {
			return nil, err
		}

		getDestPoolRateLimits := cache.LazyFetch(func() (map[cciptypes.Address]*big.Int, error) {
			tokenExecData, err1 := getExecTokenData()
			if err1 != nil {
				return nil, err1
			}
			return r.destPoolRateLimits(ctx, unexpiredReportsWithSendReqs, tokenExecData.sourceToDestTokens)
		})

		for _, unexpiredReport := range unexpiredReportsWithSendReqs {
			r.tokenDataWorker.AddJobsFromMsgs(ctx, unexpiredReport.sendRequestsWithMeta)
		}

		for _, rep := range unexpiredReportsWithSendReqs {
			if ctx.Err() != nil {
				lggr.Warn("Processing of roots killed by context")
				break
			}

			merkleRoot := rep.commitReport.MerkleRoot

			rootLggr := lggr.With("root", hexutil.Encode(merkleRoot[:]),
				"minSeqNr", rep.commitReport.Interval.Min,
				"maxSeqNr", rep.commitReport.Interval.Max,
			)

			if err := rep.validate(); err != nil {
				rootLggr.Errorw("Skipping invalid report", "err", err)
				continue
			}

			// If all messages are already executed and finalized, snooze the root for
			// config.PermissionLessExecutionThresholdSeconds so it will never be considered again.
			if allMsgsExecutedAndFinalized := rep.allRequestsAreExecutedAndFinalized(); allMsgsExecutedAndFinalized {
				rootLggr.Infow("Snoozing root forever since there are no executable txs anymore", "root", hex.EncodeToString(merkleRoot[:]))
				r.snoozedRoots.MarkAsExecuted(merkleRoot)
				continue
			}

			blessed, err := r.commitStoreReader.IsBlessed(ctx, merkleRoot)
			if err != nil {
				return nil, err
			}
			if !blessed {
				rootLggr.Infow("Report is accepted but not blessed")
				continue
			}

			tokenExecData, err := getExecTokenData()
			if err != nil {
				return nil, err
			}

			destPoolRateLimits, err := getDestPoolRateLimits()
			if err != nil {
				return nil, err
			}

			batch := r.buildBatch(
				ctx,
				rootLggr,
				rep,
				inflight,
				tokenExecData.rateLimiterTokenBucket.Tokens,
				tokenExecData.sourceTokenPrices,
				tokenExecData.destTokenPrices,
				tokenExecData.gasPrice,
				tokenExecData.sourceToDestTokens,
				destPoolRateLimits)
			if len(batch) != 0 {
				return batch, nil
			}
			r.snoozedRoots.Snooze(merkleRoot)
		}
	}
	return []ccip.ObservedMessage{}, nil
}

// destPoolRateLimits returns a map that consists of the rate limits of each destination token of the provided reports.
// If a token is missing from the returned map it either means that token was not found or token pool is disabled for this token.
func (r *ExecutionReportingPlugin) destPoolRateLimits(ctx context.Context, commitReports []commitReportWithSendRequests, sourceToDestToken map[cciptypes.Address]cciptypes.Address) (map[cciptypes.Address]*big.Int, error) {
	tokens, err := r.offRampReader.GetTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cached token pools: %w", err)
	}

	dstTokenToPool := make(map[cciptypes.Address]cciptypes.Address)
	dstPoolToToken := make(map[cciptypes.Address]cciptypes.Address)
	dstPoolAddresses := make([]cciptypes.Address, 0)

	for _, msg := range commitReports {
		for _, req := range msg.sendRequestsWithMeta {
			for _, tk := range req.TokenAmounts {
				dstToken, exists := sourceToDestToken[tk.Token]
				if !exists {
					r.lggr.Warnw("token not found on destination chain", "sourceToken", tk)
					continue
				}

				// another message with the same token exists in the report
				// we skip it since we don't want to query for the rate limit twice
				if _, seen := dstTokenToPool[dstToken]; seen {
					continue
				}

				poolAddress, exists := tokens.DestinationPool[dstToken]
				if !exists {
					return nil, fmt.Errorf("pool for token '%s' does not exist", dstToken)
				}

				if tokenAddr, seen := dstPoolToToken[poolAddress]; seen {
					return nil, fmt.Errorf("pool is already seen for token %s", tokenAddr)
				}

				dstTokenToPool[dstToken] = poolAddress
				dstPoolToToken[poolAddress] = dstToken
				dstPoolAddresses = append(dstPoolAddresses, poolAddress)
			}
		}
	}

	rateLimits, err := r.tokenPoolBatchedReader.GetInboundTokenPoolRateLimits(ctx, dstPoolAddresses)
	if err != nil {
		return nil, fmt.Errorf("fetch pool rate limits: %w", err)
	}

	res := make(map[cciptypes.Address]*big.Int, len(dstTokenToPool))
	for i, rateLimit := range rateLimits {
		// if the rate limit is disabled for this token pool then we omit it from the result
		if !rateLimit.IsEnabled {
			continue
		}

		tokenAddr, exists := dstPoolToToken[dstPoolAddresses[i]]
		if !exists {
			return nil, fmt.Errorf("pool to token mapping does not contain %s", dstPoolAddresses[i])
		}
		res[tokenAddr] = rateLimit.Tokens
	}

	return res, nil
}

// Calculates a map that indicates whether a sequence number has already been executed.
// It doesn't matter if the execution succeeded, since we don't retry previous
// attempts even if they failed. Value in the map indicates whether the log is finalized or not.
func (r *ExecutionReportingPlugin) getExecutedSeqNrsInRange(ctx context.Context, min, max uint64) (map[uint64]bool, error) {
	stateChanges, err := r.offRampReader.GetExecutionStateChangesBetweenSeqNums(
		ctx,
		min,
		max,
		int(r.offchainConfig.DestOptimisticConfirmations),
	)
	if err != nil {
		return nil, err
	}
	executedMp := make(map[uint64]bool, len(stateChanges))
	for _, stateChange := range stateChanges {
		executedMp[stateChange.SequenceNumber] = stateChange.Finalized
	}
	return executedMp, nil
}

// Builds a batch of transactions that can be executed, takes into account
// the available gas, rate limiting, execution state, nonce state, and
// profitability of execution.
func (r *ExecutionReportingPlugin) buildBatch(
	ctx context.Context,
	lggr logger.Logger,
	report commitReportWithSendRequests,
	inflight []InflightInternalExecutionReport,
	aggregateTokenLimit *big.Int,
	sourceTokenPricesUSD map[cciptypes.Address]*big.Int,
	destTokenPricesUSD map[cciptypes.Address]*big.Int,
	gasPrice *big.Int,
	sourceToDestToken map[cciptypes.Address]cciptypes.Address,
	destTokenPoolRateLimits map[cciptypes.Address]*big.Int,
) (executableMessages []ccip.ObservedMessage) {
	inflightSeqNrs, inflightAggregateValue, maxInflightSenderNonces, inflightTokenAmounts, err := inflightAggregates(inflight, destTokenPricesUSD, sourceToDestToken)
	if err != nil {
		lggr.Errorw("Unexpected error computing inflight values", "err", err)
		return []ccip.ObservedMessage{}
	}
	availableGas := uint64(r.offchainConfig.BatchGasLimit)
	expectedNonces := make(map[cciptypes.Address]uint64)
	availableDataLen := MaxDataLenPerBatch
	tokenDataRemainingDuration := MaximumAllowedTokenDataWaitTimePerBatch
	for _, msg := range report.sendRequestsWithMeta {
		msgLggr := lggr.With("messageID", hexutil.Encode(msg.MessageID[:]))
		if msg.Executed {
			msgLggr.Infow("Skipping message already executed", "seqNr", msg.SequenceNumber)
			continue
		}
		if inflightSeqNrs.Contains(msg.SequenceNumber) {
			msgLggr.Infow("Skipping message already inflight", "seqNr", msg.SequenceNumber)
			continue
		}
		if _, ok := expectedNonces[msg.Sender]; !ok {
			// First message in batch, need to populate expected nonce
			if maxInflight, ok := maxInflightSenderNonces[msg.Sender]; ok {
				// Sender already has inflight nonce, populate from there
				expectedNonces[msg.Sender] = maxInflight + 1
			} else {
				// Nothing inflight take from chain.
				// Chain holds existing nonce.
				nonce, err := r.offRampReader.GetSenderNonce(ctx, msg.Sender)
				if err != nil {
					lggr.Errorw("unable to get sender nonce", "err", err, "seqNr", msg.SequenceNumber)
					continue
				}
				expectedNonces[msg.Sender] = nonce + 1
			}
		}
		// Check expected nonce is valid
		if msg.Nonce != expectedNonces[msg.Sender] {
			msgLggr.Warnw("Skipping message invalid nonce", "have", msg.Nonce, "want", expectedNonces[msg.Sender])
			continue
		}

		if !r.isRateLimitEnoughForTokenPool(destTokenPoolRateLimits, msg.TokenAmounts, inflightTokenAmounts, sourceToDestToken) {
			msgLggr.Warnw("Skipping message token pool rate limit hit")
			continue
		}

		msgValue, err := aggregateTokenValue(destTokenPricesUSD, sourceToDestToken, msg.TokenAmounts)
		if err != nil {
			msgLggr.Errorw("Skipping message unable to compute aggregate value", "err", err)
			continue
		}

		// if token limit is smaller than message value skip message
		if tokensLeft, hasCapacity := hasEnoughTokens(aggregateTokenLimit, msgValue, inflightAggregateValue); !hasCapacity {
			msgLggr.Warnw("token limit is smaller than message value", "aggregateTokenLimit", tokensLeft.String(), "msgValue", msgValue.String())
			continue
		}

		tokenData, elapsed, err := r.getTokenDataWithTimeout(ctx, msg, tokenDataRemainingDuration)
		tokenDataRemainingDuration -= elapsed
		if err != nil {
			if errors.Is(err, tokendata.ErrNotReady) {
				msgLggr.Warnw("skipping message due to token data not ready", "err", err)
				continue
			}
			msgLggr.Errorw("skipping message error while getting token data", "err", err)
			continue
		}

		dstWrappedNativePrice, exists := destTokenPricesUSD[r.destWrappedNative]
		if !exists {
			msgLggr.Errorw("token not in dst token prices", "token", r.destWrappedNative)
			continue
		}

		// Fee boosting
		execCostUsd, err := r.gasPriceEstimator.EstimateMsgCostUSD(gasPrice, dstWrappedNativePrice, msg)
		if err != nil {
			msgLggr.Errorw("failed to estimate message cost USD", "err", err)
			return []ccip.ObservedMessage{}
		}

		// calculating the source chain fee, dividing by 1e18 for denomination.
		// For example:
		// FeeToken=link; FeeTokenAmount=1e17 i.e. 0.1 link, price is 6e18 USD/link (1 USD = 1e18),
		// availableFee is 1e17*6e18/1e18 = 6e17 = 0.6 USD

		sourceFeeTokenPrice, exists := sourceTokenPricesUSD[msg.FeeToken]
		if !exists {
			msgLggr.Errorw("token not in source token prices", "token", msg.FeeToken)
			continue
		}

		if len(msg.Data) > availableDataLen {
			msgLggr.Infow("Skipping message, insufficient remaining batch data len",
				"msgDataLen", len(msg.Data), "availableBatchDataLen", availableDataLen)
			continue
		}

		availableFee := big.NewInt(0).Mul(msg.FeeTokenAmount, sourceFeeTokenPrice)
		availableFee = availableFee.Div(availableFee, big.NewInt(1e18))
		availableFeeUsd := waitBoostedFee(time.Since(msg.BlockTimestamp), availableFee, r.offchainConfig.RelativeBoostPerWaitHour)
		if availableFeeUsd.Cmp(execCostUsd) < 0 {
			msgLggr.Infow("Insufficient remaining fee", "availableFeeUsd", availableFeeUsd, "execCostUsd", execCostUsd,
				"sourceBlockTimestamp", msg.BlockTimestamp, "waitTime", time.Since(msg.BlockTimestamp), "boost", r.offchainConfig.RelativeBoostPerWaitHour)
			continue
		}

		messageMaxGas, err := calculateMessageMaxGas(
			msg.GasLimit,
			len(report.sendRequestsWithMeta),
			len(msg.Data),
			len(msg.TokenAmounts),
		)
		if err != nil {
			msgLggr.Errorw("calculate message max gas", "err", err)
			continue
		}

		// Check sufficient gas in batch
		if availableGas < messageMaxGas {
			msgLggr.Infow("Insufficient remaining gas in batch limit", "availableGas", availableGas, "messageMaxGas", messageMaxGas)
			continue
		}
		availableGas -= messageMaxGas
		aggregateTokenLimit.Sub(aggregateTokenLimit, msgValue)
		for _, tk := range msg.TokenAmounts {
			dstToken, exists := sourceToDestToken[tk.Token]
			if !exists {
				msgLggr.Warnw("destination token does not exist", "token", tk.Token)
				continue
			}
			if rl, exists := destTokenPoolRateLimits[dstToken]; exists {
				destTokenPoolRateLimits[dstToken] = rl.Sub(rl, tk.Amount)
			}
		}

		msgLggr.Infow("Adding msg to batch", "seqNr", msg.SequenceNumber, "nonce", msg.Nonce,
			"value", msgValue, "aggregateTokenLimit", aggregateTokenLimit)
		executableMessages = append(executableMessages, ccip.NewObservedMessage(msg.SequenceNumber, tokenData))

		// after message is added to the batch, decrease the available data length
		availableDataLen -= len(msg.Data)

		expectedNonces[msg.Sender] = msg.Nonce + 1
	}

	return executableMessages
}

// getTokenDataWithCappedLatency gets the token data for the provided message.
// Stops and returns an error if more than allowedWaitingTime is passed.
func (r *ExecutionReportingPlugin) getTokenDataWithTimeout(
	ctx context.Context,
	msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta,
	timeout time.Duration,
) ([][]byte, time.Duration, error) {
	if len(msg.TokenAmounts) == 0 {
		return nil, 0, nil
	}

	ctxTimeout, cf := context.WithTimeout(ctx, timeout)
	defer cf()
	tStart := time.Now()
	tokenData, err := r.tokenDataWorker.GetMsgTokenData(ctxTimeout, msg)
	tDur := time.Since(tStart)
	return tokenData, tDur, err
}

func (r *ExecutionReportingPlugin) isRateLimitEnoughForTokenPool(
	destTokenPoolRateLimits map[cciptypes.Address]*big.Int,
	sourceTokenAmounts []cciptypes.TokenAmount,
	inflightTokenAmounts map[cciptypes.Address]*big.Int,
	sourceToDestToken map[cciptypes.Address]cciptypes.Address,
) bool {
	rateLimitsCopy := make(map[cciptypes.Address]*big.Int)
	for destToken, rl := range destTokenPoolRateLimits {
		rateLimitsCopy[destToken] = new(big.Int).Set(rl)
	}

	for sourceToken, amount := range inflightTokenAmounts {
		if destToken, exists := sourceToDestToken[sourceToken]; exists {
			if rl, exists := rateLimitsCopy[destToken]; exists {
				rateLimitsCopy[destToken] = rl.Sub(rl, amount)
			}
		}
	}

	for _, sourceToken := range sourceTokenAmounts {
		destToken, exists := sourceToDestToken[sourceToken.Token]
		if !exists {
			r.lggr.Warnw("dest token not found", "sourceToken", sourceToken.Token)
			continue
		}

		rl, exists := rateLimitsCopy[destToken]
		if !exists {
			r.lggr.Debugw("rate limit not applied to token", "token", destToken)
			continue
		}

		if rl.Cmp(sourceToken.Amount) < 0 {
			r.lggr.Warnw("token pool rate limit reached",
				"token", sourceToken.Token, "destToken", destToken, "amount", sourceToken.Amount, "rateLimit", rl)
			return false
		}
		rateLimitsCopy[destToken] = rl.Sub(rl, sourceToken.Amount)
	}

	return true
}

func hasEnoughTokens(tokenLimit *big.Int, msgValue *big.Int, inflightValue *big.Int) (*big.Int, bool) {
	tokensLeft := big.NewInt(0).Sub(tokenLimit, inflightValue)
	return tokensLeft, tokensLeft.Cmp(msgValue) >= 0
}

func calculateMessageMaxGas(gasLimit *big.Int, numRequests, dataLen, numTokens int) (uint64, error) {
	if !gasLimit.IsUint64() {
		return 0, fmt.Errorf("gas limit %s cannot be casted to uint64", gasLimit)
	}

	gasLimitU64 := gasLimit.Uint64()
	gasOverHeadGas := maxGasOverHeadGas(numRequests, dataLen, numTokens)
	messageMaxGas := gasLimitU64 + gasOverHeadGas

	if messageMaxGas < gasLimitU64 || messageMaxGas < gasOverHeadGas {
		return 0, fmt.Errorf("message max gas overflow, gasLimit=%d gasOverHeadGas=%d", gasLimitU64, gasOverHeadGas)
	}

	return messageMaxGas, nil
}

// helper struct to hold the commitReport and the related send requests
type commitReportWithSendRequests struct {
	commitReport         cciptypes.CommitStoreReport
	sendRequestsWithMeta []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
}

func (r *commitReportWithSendRequests) validate() error {
	// make sure that number of messages is the expected
	if exp := int(r.commitReport.Interval.Max - r.commitReport.Interval.Min + 1); len(r.sendRequestsWithMeta) != exp {
		return errors.Errorf(
			"unexpected missing sendRequestsWithMeta in committed root %x have %d want %d", r.commitReport.MerkleRoot, len(r.sendRequestsWithMeta), exp)
	}

	return nil
}

func (r *commitReportWithSendRequests) allRequestsAreExecutedAndFinalized() bool {
	for _, req := range r.sendRequestsWithMeta {
		if !req.Executed || !req.Finalized {
			return false
		}
	}
	return true
}

// checks if the send request fits the commit report interval
func (r *commitReportWithSendRequests) sendReqFits(sendReq cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) bool {
	return sendReq.SequenceNumber >= r.commitReport.Interval.Min &&
		sendReq.SequenceNumber <= r.commitReport.Interval.Max
}

// getReportsWithSendRequests returns the target reports with populated send requests.
func (r *ExecutionReportingPlugin) getReportsWithSendRequests(
	ctx context.Context,
	reports []cciptypes.CommitStoreReport,
) ([]commitReportWithSendRequests, error) {
	if len(reports) == 0 {
		return nil, nil
	}

	// find interval from all the reports
	intervalMin := reports[0].Interval.Min
	intervalMax := reports[0].Interval.Max
	for _, report := range reports[1:] {
		if report.Interval.Max > intervalMax {
			intervalMax = report.Interval.Max
		}
		if report.Interval.Min < intervalMin {
			intervalMin = report.Interval.Min
		}
	}

	// use errgroup to fetch send request logs and executed sequence numbers in parallel
	eg := &errgroup.Group{}

	var sendRequests []cciptypes.EVM2EVMMessageWithTxMeta
	eg.Go(func() error {
		// We don't need to double-check if logs are finalized because we already checked that in the Commit phase.
		sendReqs, err := r.onRampReader.GetSendRequestsBetweenSeqNums(ctx, intervalMin, intervalMax, false)
		if err != nil {
			return err
		}
		sendRequests = sendReqs
		return nil
	})

	var executedSeqNums map[uint64]bool
	eg.Go(func() error {
		// get executed sequence numbers
		executedMp, err := r.getExecutedSeqNrsInRange(ctx, intervalMin, intervalMax)
		if err != nil {
			return err
		}
		executedSeqNums = executedMp
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	reportsWithSendReqs := make([]commitReportWithSendRequests, len(reports))
	for i, report := range reports {
		reportsWithSendReqs[i] = commitReportWithSendRequests{
			commitReport:         report,
			sendRequestsWithMeta: make([]cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, 0, report.Interval.Max-report.Interval.Min+1),
		}
	}

	for _, sendReq := range sendRequests {
		// if value exists in the map then it's executed
		// if value exists, and it's true then it's considered finalized
		finalized, executed := executedSeqNums[sendReq.SequenceNumber]

		reqWithMeta := cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
			EVM2EVMMessage: sendReq.EVM2EVMMessage,
			BlockTimestamp: time.UnixMilli(sendReq.BlockTimestampUnixMilli),
			Executed:       executed,
			Finalized:      finalized,
			LogIndex:       uint(sendReq.LogIndex),
			TxHash:         sendReq.TxHash,
		}

		// attach the msg to the appropriate reports
		for i := range reportsWithSendReqs {
			if reportsWithSendReqs[i].sendReqFits(reqWithMeta) {
				reportsWithSendReqs[i].sendRequestsWithMeta = append(reportsWithSendReqs[i].sendRequestsWithMeta, reqWithMeta)
			}
		}
	}

	return reportsWithSendReqs, nil
}

func aggregateTokenValue(destTokenPricesUSD map[cciptypes.Address]*big.Int, sourceToDest map[cciptypes.Address]cciptypes.Address, tokensAndAmount []cciptypes.TokenAmount) (*big.Int, error) {
	sum := big.NewInt(0)
	for i := 0; i < len(tokensAndAmount); i++ {
		price, ok := destTokenPricesUSD[sourceToDest[tokensAndAmount[i].Token]]
		if !ok {
			return nil, errors.Errorf("do not have price for source token %v", tokensAndAmount[i].Token)
		}
		sum.Add(sum, new(big.Int).Quo(new(big.Int).Mul(price, tokensAndAmount[i].Amount), big.NewInt(1e18)))
	}
	return sum, nil
}

// Assumes non-empty report. Messages to execute can span more than one report, but are assumed to be in order of increasing
// sequence number.
func (r *ExecutionReportingPlugin) buildReport(ctx context.Context, lggr logger.Logger, observedMessages []ccip.ObservedMessage) ([]byte, error) {
	if err := validateSeqNumbers(ctx, r.commitStoreReader, observedMessages); err != nil {
		return nil, err
	}
	commitReport, err := getCommitReportForSeqNum(ctx, r.commitStoreReader, observedMessages[0].SeqNr)
	if err != nil {
		return nil, err
	}
	lggr.Infow("Building execution report", "observations", observedMessages, "merkleRoot", hexutil.Encode(commitReport.MerkleRoot[:]), "report", commitReport)

	sendReqsInRoot, _, tree, err := getProofData(ctx, r.onRampReader, commitReport.Interval)
	if err != nil {
		return nil, err
	}

	// cap messages which fits MaxExecutionReportLength (after serialized)
	capped := sort.Search(len(observedMessages), func(i int) bool {
		report, err2 := buildExecutionReportForMessages(sendReqsInRoot, tree, commitReport.Interval, observedMessages[:i+1])
		if err2 != nil {
			r.lggr.Errorw("build execution report", "err", err2)
			return false
		}

		encoded, err2 := r.offRampReader.EncodeExecutionReport(report)
		if err2 != nil {
			// false makes Search keep looking to the right, always including any "erroring" ObservedMessage and allowing us to detect in the bottom
			return false
		}
		return len(encoded) > MaxExecutionReportLength
	})

	execReport, err := buildExecutionReportForMessages(sendReqsInRoot, tree, commitReport.Interval, observedMessages[:capped])
	if err != nil {
		return nil, err
	}

	r.metricsCollector.NumberOfMessagesProcessed(ccip.Report, len(execReport.Messages))
	encodedReport, err := r.offRampReader.EncodeExecutionReport(execReport)
	if err != nil {
		return nil, err
	}

	if capped < len(observedMessages) {
		lggr.Warnf(
			"Capping report to fit MaxExecutionReportLength: msgsCount %d -> %d, bytes %d, bytesLimit %d",
			len(observedMessages), capped, len(encodedReport), MaxExecutionReportLength,
		)
	}
	// Double check this verifies before sending.
	valid, err := r.commitStoreReader.VerifyExecutionReport(ctx, execReport)
	if err != nil {
		return nil, errors.Wrap(err, "unable to verify")
	}
	if !valid {
		return nil, errors.New("root does not verify")
	}
	return encodedReport, nil
}

func (r *ExecutionReportingPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	lggr := r.lggr.Named("ExecutionReport")
	parsableObservations := ccip.GetParsableObservations[ccip.ExecutionObservation](lggr, observations)
	// Need at least F+1 observations
	if len(parsableObservations) <= r.F {
		lggr.Warn("Non-empty observations <= F, need at least F+1 to continue")
		return false, nil, nil
	}

	observedMessages, err := calculateObservedMessagesConsensus(parsableObservations, r.F)
	if err != nil {
		return false, nil, err
	}
	if len(observedMessages) == 0 {
		return false, nil, nil
	}

	report, err := r.buildReport(ctx, lggr, observedMessages)
	if err != nil {
		return false, nil, err
	}
	lggr.Infow("Report", "executableObservations", observedMessages)
	return true, report, nil
}

type tallyKey struct {
	seqNr         uint64
	tokenDataHash [32]byte
}

type tallyVal struct {
	tally     int
	tokenData [][]byte
}

func calculateObservedMessagesConsensus(observations []ccip.ExecutionObservation, f int) ([]ccip.ObservedMessage, error) {
	tally := make(map[tallyKey]tallyVal)
	for _, obs := range observations {
		for seqNr, msgData := range obs.Messages {
			tokenDataHash, err := hashlib.BytesOfBytesKeccak(msgData.TokenData)
			if err != nil {
				return nil, fmt.Errorf("bytes of bytes keccak: %w", err)
			}

			key := tallyKey{seqNr: seqNr, tokenDataHash: tokenDataHash}
			if val, ok := tally[key]; ok {
				tally[key] = tallyVal{tally: val.tally + 1, tokenData: msgData.TokenData}
			} else {
				tally[key] = tallyVal{tally: 1, tokenData: msgData.TokenData}
			}
		}
	}

	// We might have different token data for the same sequence number.
	// For that purpose we want to keep the token data with the most occurrences.
	seqNumTally := make(map[uint64]tallyVal)

	// order tally keys to make looping over the entries deterministic
	tallyKeys := make([]tallyKey, 0, len(tally))
	for key := range tally {
		tallyKeys = append(tallyKeys, key)
	}
	sort.Slice(tallyKeys, func(i, j int) bool {
		return hex.EncodeToString(tallyKeys[i].tokenDataHash[:]) < hex.EncodeToString(tallyKeys[j].tokenDataHash[:])
	})

	for _, key := range tallyKeys {
		tallyInfo := tally[key]
		existingTally, exists := seqNumTally[key.seqNr]
		if tallyInfo.tally > f && (!exists || tallyInfo.tally > existingTally.tally) {
			seqNumTally[key.seqNr] = tallyInfo
		}
	}

	finalSequenceNumbers := make([]ccip.ObservedMessage, 0, len(seqNumTally))
	for seqNr, tallyInfo := range seqNumTally {
		finalSequenceNumbers = append(finalSequenceNumbers, ccip.NewObservedMessage(seqNr, tallyInfo.tokenData))
	}
	// buildReport expects sorted sequence numbers (tally map is non-deterministic).
	sort.Slice(finalSequenceNumbers, func(i, j int) bool {
		return finalSequenceNumbers[i].SeqNr < finalSequenceNumbers[j].SeqNr
	})
	return finalSequenceNumbers, nil
}

func (r *ExecutionReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("ShouldAcceptFinalizedReport")
	execReport, err := r.offRampReader.DecodeExecutionReport(report)
	if err != nil {
		lggr.Errorw("Unable to decode report", "err", err)
		return false, err
	}
	lggr = lggr.With("messageIDs", ccipcommon.GetMessageIDsAsHexString(execReport.Messages))

	// If the first message is executed already, this execution report is stale, and we do not accept it.
	stale, err := r.isStaleReport(ctx, execReport.Messages)
	if err != nil {
		return false, err
	}
	if stale {
		lggr.Info("Execution report is stale")
		return false, nil
	}
	// Else just assume in flight
	if err = r.inflightReports.add(lggr, execReport.Messages); err != nil {
		return false, err
	}
	lggr.Info("Accepting finalized report")
	return true, nil
}

func (r *ExecutionReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("ShouldTransmitAcceptedReport")
	execReport, err := r.offRampReader.DecodeExecutionReport(report)
	if err != nil {
		lggr.Errorw("Unable to decode report", "err", err)
		return false, nil
	}
	lggr = lggr.With("messageIDs", ccipcommon.GetMessageIDsAsHexString(execReport.Messages))

	// If report is not stale we transmit.
	// When the executeTransmitter enqueues the tx for tx manager,
	// we mark it as execution_sent, removing it from the set of inflight messages.
	stale, err := r.isStaleReport(ctx, execReport.Messages)
	if err != nil {
		return false, err
	}
	if stale {
		lggr.Info("Execution report is stale")
		return false, nil
	}

	lggr.Info("Transmitting finalized report")
	return true, err
}

func (r *ExecutionReportingPlugin) isStaleReport(ctx context.Context, messages []cciptypes.EVM2EVMMessage) (bool, error) {
	if len(messages) == 0 {
		return true, fmt.Errorf("messages are empty")
	}

	// If the first message is executed already, this execution report is stale.
	// Note the default execution state, including for arbitrary seq number not yet committed
	// is ExecutionStateUntouched.
	msgState, err := r.offRampReader.GetExecutionState(ctx, messages[0].SequenceNumber)
	if err != nil {
		return true, err
	}
	if state := cciptypes.MessageExecutionState(msgState); state == cciptypes.ExecutionStateFailure || state == cciptypes.ExecutionStateSuccess {
		return true, nil
	}

	return false, nil
}

func (r *ExecutionReportingPlugin) Close() error {
	return nil
}

func inflightAggregates(
	inflight []InflightInternalExecutionReport,
	destTokenPrices map[cciptypes.Address]*big.Int,
	sourceToDest map[cciptypes.Address]cciptypes.Address,
) (mapset.Set[uint64], *big.Int, map[cciptypes.Address]uint64, map[cciptypes.Address]*big.Int, error) {
	inflightSeqNrs := mapset.NewSet[uint64]()
	inflightAggregateValue := big.NewInt(0)
	maxInflightSenderNonces := make(map[cciptypes.Address]uint64)
	inflightTokenAmounts := make(map[cciptypes.Address]*big.Int)

	for _, rep := range inflight {
		for _, message := range rep.messages {
			inflightSeqNrs.Add(message.SequenceNumber)
			msgValue, err := aggregateTokenValue(destTokenPrices, sourceToDest, message.TokenAmounts)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			inflightAggregateValue.Add(inflightAggregateValue, msgValue)
			maxInflightSenderNonce, ok := maxInflightSenderNonces[message.Sender]
			if !ok || message.Nonce > maxInflightSenderNonce {
				maxInflightSenderNonces[message.Sender] = message.Nonce
			}

			for _, tk := range message.TokenAmounts {
				if rl, exists := inflightTokenAmounts[tk.Token]; exists {
					inflightTokenAmounts[tk.Token] = rl.Add(rl, tk.Amount)
				} else {
					inflightTokenAmounts[tk.Token] = new(big.Int).Set(tk.Amount)
				}
			}
		}
	}
	return inflightSeqNrs, inflightAggregateValue, maxInflightSenderNonces, inflightTokenAmounts, nil
}

// getTokensPrices returns token prices of the given price registry,
// price values are USD per 1e18 of smallest token denomination, in base units 1e18 (e.g. 5$ = 5e18 USD per 1e18 units).
// this function is used for price registry of both source and destination chains.
func getTokensPrices(ctx context.Context, priceRegistry ccipdata.PriceRegistryReader, tokens []cciptypes.Address) (map[cciptypes.Address]*big.Int, error) {
	prices := make(map[cciptypes.Address]*big.Int)

	fetchedPrices, err := priceRegistry.GetTokenPrices(ctx, tokens)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get token prices of %v", tokens)
	}

	// price registry should always return a price per token ordered by input tokens
	if len(fetchedPrices) != len(tokens) {
		return nil, fmt.Errorf("token prices length exp=%d actual=%d", len(tokens), len(fetchedPrices))
	}

	for i, token := range tokens {
		// price of a token can never be zero
		if fetchedPrices[i].Value.BitLen() == 0 {
			return nil, fmt.Errorf("price of token %s is zero (price registry=%s)", token, priceRegistry.Address())
		}

		// price registry should not report different price for the same token
		price, exists := prices[token]
		if exists && fetchedPrices[i].Value.Cmp(price) != 0 {
			return nil, fmt.Errorf("price registry reported different prices (%s and %s) for the same token %s",
				fetchedPrices[i].Value, price, token)
		}

		prices[token] = fetchedPrices[i].Value
	}

	return prices, nil
}

func (r *ExecutionReportingPlugin) getUnexpiredCommitReports(
	ctx context.Context,
	commitStoreReader ccipdata.CommitStoreReader,
	permissionExecutionThreshold time.Duration,
	lggr logger.Logger,
) ([]cciptypes.CommitStoreReport, error) {
	acceptedReports, err := commitStoreReader.GetAcceptedCommitReportsGteTimestamp(
		ctx,
		time.Now().Add(-permissionExecutionThreshold),
		0,
	)
	if err != nil {
		return nil, err
	}

	var reports []cciptypes.CommitStoreReport
	for _, acceptedReport := range acceptedReports {
		reports = append(reports, acceptedReport.CommitStoreReport)
	}

	notSnoozedReports := make([]cciptypes.CommitStoreReport, 0)
	for _, report := range reports {
		if r.snoozedRoots.IsSnoozed(report.MerkleRoot) {
			lggr.Debugw("Skipping snoozed root",
				"minSeqNr", report.Interval.Min,
				"maxSeqNr", report.Interval.Max,
				"root", hex.EncodeToString(report.MerkleRoot[:]),
			)
			continue
		}
		notSnoozedReports = append(notSnoozedReports, report)
	}

	r.metricsCollector.UnexpiredCommitRoots(len(notSnoozedReports))
	lggr.Infow("Unexpired roots", "all", len(reports), "notSnoozed", len(notSnoozedReports))
	return notSnoozedReports, nil
}

type execTokenData struct {
	rateLimiterTokenBucket cciptypes.TokenBucketRateLimit
	sourceTokenPrices      map[cciptypes.Address]*big.Int
	destTokenPrices        map[cciptypes.Address]*big.Int
	sourceToDestTokens     map[cciptypes.Address]cciptypes.Address
	gasPrice               *big.Int
}

// prepareTokenExecData gather all the pre-execution data needed for token execution into a single lazy call.
// This is done to avoid fetching the data multiple times for each message. Additionally, most of the RPC calls
// within that function is cached, so it should be relatively fast and not require any RPC batching.
func (r *ExecutionReportingPlugin) prepareTokenExecData(ctx context.Context) (execTokenData, error) {
	// This could result in slightly different values on each call as
	// the function returns the allowed amount at the time of the last block.
	// Since this will only increase over time, the highest observed value will
	// always be the lower bound of what would be available on chain
	// since we already account for inflight txs.
	rateLimiterTokenBucket, err := r.offRampReader.CurrentRateLimiterState(ctx)
	if err != nil {
		return execTokenData{}, err
	}

	sourceFeeTokens, err := r.sourcePriceRegistry.GetFeeTokens(ctx)
	if err != nil {
		return execTokenData{}, fmt.Errorf("get source fee tokens: %w", err)
	}
	sourceTokensPrices, err := getTokensPrices(
		ctx,
		r.sourcePriceRegistry,
		ccipcommon.FlattenUniqueSlice(
			sourceFeeTokens,
			[]cciptypes.Address{r.sourceWrappedNativeToken},
		),
	)
	if err != nil {
		return execTokenData{}, err
	}

	destFeeTokens, destBridgedTokens, err := ccipcommon.GetDestinationTokens(ctx, r.offRampReader, r.destPriceRegistry)
	if err != nil {
		return execTokenData{}, fmt.Errorf("get destination tokens: %w", err)
	}
	destTokenPrices, err := getTokensPrices(
		ctx,
		r.destPriceRegistry,
		ccipcommon.FlattenUniqueSlice(
			destFeeTokens,
			destBridgedTokens,
			[]cciptypes.Address{r.destWrappedNative},
		),
	)
	if err != nil {
		return execTokenData{}, err
	}

	sourceToDestTokens, err := r.offRampReader.GetSourceToDestTokensMapping(ctx)
	if err != nil {
		return execTokenData{}, err
	}

	gasPrice, err := r.gasPriceEstimator.GetGasPrice(ctx)
	if err != nil {
		return execTokenData{}, err
	}

	return execTokenData{
		rateLimiterTokenBucket: rateLimiterTokenBucket,
		sourceTokenPrices:      sourceTokensPrices,
		sourceToDestTokens:     sourceToDestTokens,
		destTokenPrices:        destTokenPrices,
		gasPrice:               gasPrice,
	}, nil
}

// selectReportsToFillBatch returns the reports to fill the message limit. Single Commit Root contains exactly (Interval.Max - Interval.Min + 1) messages.
// We keep adding reports until we reach the message limit. Please see the tests for more examples and edge cases.
// unexpiredReports have to be sorted by Interval.Min. Otherwise, the batching logic will not be efficient,
// because it picks messages and execution states based on the report[0].Interval.Min - report[len-1].Interval.Max range.
// Having unexpiredReports not sorted properly will lead to fetching more messages and execution states to the memory than the messagesLimit provided.
// However, logs from LogPoller are returned ordered by (block_number, log_index), so it should preserve the order of Interval.Min.
// Single CommitRoot can have up to 256 messages, with current MessagesIterationStep of 800, it means processing 4 CommitRoots at once.
func selectReportsToFillBatch(unexpiredReports []cciptypes.CommitStoreReport, messagesLimit uint64) ([]cciptypes.CommitStoreReport, int) {
	currentNumberOfMessages := uint64(0)
	var index int

	for index = range unexpiredReports {
		currentNumberOfMessages += unexpiredReports[index].Interval.Max - unexpiredReports[index].Interval.Min + 1
		if currentNumberOfMessages >= messagesLimit {
			break
		}
	}
	index = min(index+1, len(unexpiredReports))
	return unexpiredReports[:index], index
}
