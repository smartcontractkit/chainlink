package ccipexec

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/batchreader"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/statuschecker"
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
	MessagesIterationStep = 1024
)

var (
	_ types.ReportingPluginFactory = &ExecutionReportingPluginFactory{}
	_ types.ReportingPlugin        = &ExecutionReportingPlugin{}
)

type ExecutionPluginStaticConfig struct {
	lggr                          logger.Logger
	onRampReader                  ccipdata.OnRampReader
	offRampReader                 ccipdata.OffRampReader
	commitStoreReader             ccipdata.CommitStoreReader
	sourcePriceRegistryProvider   ccipdataprovider.PriceRegistry
	sourceWrappedNativeToken      cciptypes.Address
	tokenDataWorker               tokendata.Worker
	destChainSelector             uint64
	priceRegistryProvider         ccipdataprovider.PriceRegistry // destination price registry provider.
	tokenPoolBatchedReader        batchreader.TokenPoolBatchedReader
	metricsCollector              ccip.PluginMetricsCollector
	chainHealthcheck              cache.ChainHealthcheck
	newReportingPluginRetryConfig ccipdata.RetryConfig
	txmStatusChecker              statuschecker.CCIPTransactionStatusChecker
}

type ExecutionReportingPlugin struct {
	// Misc
	F                int
	lggr             logger.Logger
	offchainConfig   cciptypes.ExecOffchainConfig
	tokenDataWorker  tokendata.Worker
	metricsCollector ccip.PluginMetricsCollector
	batchingStrategy BatchingStrategy

	// Source
	gasPriceEstimator           prices.GasPriceEstimatorExec
	sourcePriceRegistry         ccipdata.PriceRegistryReader
	sourcePriceRegistryProvider ccipdataprovider.PriceRegistry
	sourcePriceRegistryLock     sync.RWMutex
	sourceWrappedNativeToken    cciptypes.Address
	onRampReader                ccipdata.OnRampReader

	// Dest
	destChainSelector      uint64
	commitStoreReader      ccipdata.CommitStoreReader
	destPriceRegistry      ccipdata.PriceRegistryReader
	destWrappedNative      cciptypes.Address
	onchainConfig          cciptypes.ExecOnchainConfig
	offRampReader          ccipdata.OffRampReader
	tokenPoolBatchedReader batchreader.TokenPoolBatchedReader

	// State
	inflightReports  *inflightExecReportsContainer
	commitRootsCache cache.CommitsRootsCache
	chainHealthcheck cache.ChainHealthcheck
}

func (r *ExecutionReportingPlugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

func (r *ExecutionReportingPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	lggr := r.lggr.Named("ExecutionObservation")
	if healthy, err := r.chainHealthcheck.IsHealthy(ctx); err != nil {
		return nil, err
	} else if !healthy {
		return nil, ccip.ErrChainIsNotHealthy
	}

	// Ensure that the source price registry is synchronized with the onRamp.
	if err := r.ensurePriceRegistrySynchronization(ctx); err != nil {
		return nil, fmt.Errorf("ensuring price registry synchronization: %w", err)
	}

	// Expire any inflight reports.
	r.inflightReports.expire(lggr)
	inFlight := r.inflightReports.getAll()

	executableObservations, err := r.getExecutableObservations(ctx, lggr, inFlight)
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

func (r *ExecutionReportingPlugin) getExecutableObservations(ctx context.Context, lggr logger.Logger, inflight []InflightInternalExecutionReport) ([]ccip.ObservedMessage, error) {
	unexpiredReports, err := r.commitRootsCache.RootsEligibleForExecution(ctx)
	if err != nil {
		return nil, err
	}
	r.metricsCollector.UnexpiredCommitRoots(len(unexpiredReports))

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
				r.commitRootsCache.MarkAsExecuted(merkleRoot)
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

			batch, msgExecStates := r.buildBatch(
				ctx,
				inflight,
				rootLggr,
				rep,
				tokenExecData.rateLimiterTokenBucket.Tokens,
				tokenExecData.sourceTokenPrices,
				tokenExecData.destTokenPrices,
				tokenExecData.gasPrice,
				tokenExecData.sourceToDestTokens)
			if len(batch) != 0 {
				lggr.Infow("Execution batch created", "batchSize", len(batch), "messageStates", msgExecStates)
				return batch, nil
			}
			r.commitRootsCache.Snooze(merkleRoot)
		}
	}
	return []ccip.ObservedMessage{}, nil
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
		executedMp[stateChange.SequenceNumber] = stateChange.TxMeta.IsFinalized()
	}
	return executedMp, nil
}

// Builds a batch of transactions that can be executed, takes into account
// the available gas, rate limiting, execution state, nonce state, and
// profitability of execution.
func (r *ExecutionReportingPlugin) buildBatch(
	ctx context.Context,
	inflight []InflightInternalExecutionReport,
	lggr logger.Logger,
	report commitReportWithSendRequests,
	aggregateTokenLimit *big.Int,
	sourceTokenPricesUSD map[cciptypes.Address]*big.Int,
	destTokenPricesUSD map[cciptypes.Address]*big.Int,
	gasPrice *big.Int,
	sourceToDestToken map[cciptypes.Address]cciptypes.Address,
) ([]ccip.ObservedMessage, []messageExecStatus) {
	// We assume that next observation will start after previous epoch transmission so nonces should be already updated onchain.
	// Worst case scenario we will try to process the same message again, and it will be skipped but protocol would progress anyway.
	// We don't use inflightCache here to avoid cases in which inflight cache keeps progressing but due to transmission failures
	// previous reports are not included onchain. That can lead to issues with IncorrectNonce skips,
	// because we enforce sequential processing per sender (per sender's nonce ordering is enforced by Offramp contract)
	sendersNonce, err := r.offRampReader.ListSenderNonces(ctx, report.uniqueSenders())
	if err != nil {
		lggr.Errorw("Fetching senders nonce", "err", err)
		return []ccip.ObservedMessage{}, []messageExecStatus{}
	}

	inflightAggregateValue, err := getInflightAggregateRateLimit(lggr, inflight, destTokenPricesUSD, sourceToDestToken)
	if err != nil {
		lggr.Errorw("Unexpected error computing inflight values", "err", err)
		return []ccip.ObservedMessage{}, nil
	}

	batchCtx := &BatchContext{
		report,
		inflight,
		inflightAggregateValue,
		lggr,
		MaxDataLenPerBatch,
		uint64(r.offchainConfig.BatchGasLimit),
		make(map[cciptypes.Address]uint64),
		sendersNonce,
		sourceTokenPricesUSD,
		destTokenPricesUSD,
		gasPrice,
		sourceToDestToken,
		aggregateTokenLimit,
		MaximumAllowedTokenDataWaitTimePerBatch,
		r.tokenDataWorker,
		r.gasPriceEstimator,
		r.destWrappedNative,
		r.offchainConfig,
		r.destChainSelector,
	}

	return r.batchingStrategy.BuildBatch(ctx, batchCtx)
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

		encoded, err2 := r.offRampReader.EncodeExecutionReport(ctx, report)
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

	encodedReport, err := r.offRampReader.EncodeExecutionReport(ctx, execReport)
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
	if len(execReport.Messages) > 0 {
		r.metricsCollector.NumberOfMessagesProcessed(ccip.Report, len(execReport.Messages))
		r.metricsCollector.SequenceNumber(ccip.Report, execReport.Messages[len(execReport.Messages)-1].SequenceNumber)
	}
	return encodedReport, nil
}

// Returns required number of observations to reach consensus
func (r *ExecutionReportingPlugin) getConsensusThreshold() int {
	// Default consensus threshold is F+1
	consensusThreshold := r.F + 1
	if r.batchingStrategy.GetBatchingStrategyID() == ZKOverflowBatchingStrategyID {
		// For batching strategy 1, consensus threshold is 2F+1
		// This is because chains that can overflow need to reach consensus during the inflight cache period
		// to avoid 2 transmissions round of an overflown message.
		consensusThreshold = 2*r.F + 1
	}
	return consensusThreshold
}

func (r *ExecutionReportingPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	lggr := r.lggr.Named("ExecutionReport")
	if healthy, err := r.chainHealthcheck.IsHealthy(ctx); err != nil {
		return false, nil, err
	} else if !healthy {
		return false, nil, ccip.ErrChainIsNotHealthy
	}
	consensusThreshold := r.getConsensusThreshold()
	lggr.Infof("Consensus threshold set to: %d", consensusThreshold)

	parsableObservations := ccip.GetParsableObservations[ccip.ExecutionObservation](lggr, observations)
	if len(parsableObservations) < consensusThreshold {
		lggr.Warnf("Insufficient observations: only %d received, but need more than %d to proceed", len(parsableObservations), consensusThreshold)
		return false, nil, nil
	}

	observedMessages, err := calculateObservedMessagesConsensus(parsableObservations, consensusThreshold)
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
			tokenDataHash, err := hashutil.BytesOfBytesKeccak(msgData.TokenData)
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
	execReport, err := r.offRampReader.DecodeExecutionReport(ctx, report)
	if err != nil {
		lggr.Errorw("Unable to decode report", "err", err)
		return false, err
	}
	lggr = lggr.With("messageIDs", ccipcommon.GetMessageIDsAsHexString(execReport.Messages))

	if healthy, err1 := r.chainHealthcheck.IsHealthy(ctx); err1 != nil {
		return false, err1
	} else if !healthy {
		return false, ccip.ErrChainIsNotHealthy
	}
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
	if len(execReport.Messages) > 0 {
		r.metricsCollector.SequenceNumber(ccip.ShouldAccept, execReport.Messages[len(execReport.Messages)-1].SequenceNumber)
	}
	lggr.Info("Accepting finalized report")
	return true, nil
}

func (r *ExecutionReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	lggr := r.lggr.Named("ShouldTransmitAcceptedReport")
	execReport, err := r.offRampReader.DecodeExecutionReport(ctx, report)
	if err != nil {
		lggr.Errorw("Unable to decode report", "err", err)
		return false, nil
	}
	lggr = lggr.With("messageIDs", ccipcommon.GetMessageIDsAsHexString(execReport.Messages))

	if healthy, err1 := r.chainHealthcheck.IsHealthy(ctx); err1 != nil {
		return false, err1
	} else if !healthy {
		return false, ccip.ErrChainIsNotHealthy
	}
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
		return true, errors.New("messages are empty")
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

func getInflightAggregateRateLimit(
	lggr logger.Logger,
	inflight []InflightInternalExecutionReport,
	destTokenPrices map[cciptypes.Address]*big.Int,
	sourceToDest map[cciptypes.Address]cciptypes.Address,
) (*big.Int, error) {
	inflightAggregateValue := big.NewInt(0)

	for _, rep := range inflight {
		for _, message := range rep.messages {
			msgValue, err := aggregateTokenValue(lggr, destTokenPrices, sourceToDest, message.TokenAmounts)
			if err != nil {
				return nil, err
			}
			inflightAggregateValue.Add(inflightAggregateValue, msgValue)
		}
	}
	return inflightAggregateValue, nil
}

// getTokensPrices returns token prices of the given price registry,
// price values are USD per 1e18 of smallest token denomination, in base units 1e18 (e.g. 5$ = 5e18 USD per 1e18 units).
// this function is used for price registry of both source and destination chains.
func getTokensPrices(ctx context.Context, priceRegistry ccipdata.PriceRegistryReader, tokens []cciptypes.Address) (map[cciptypes.Address]*big.Int, error) {
	tokenPrices := make(map[cciptypes.Address]*big.Int)

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
			priceRegistryAddress, err := priceRegistry.Address(ctx)
			if err != nil {
				return nil, fmt.Errorf("get price registry address: %w", err)
			}
			return nil, fmt.Errorf("price of token %s is zero (price registry=%s)", token, priceRegistryAddress)
		}

		// price registry should not report different price for the same token
		price, exists := tokenPrices[token]
		if exists && fetchedPrices[i].Value.Cmp(price) != 0 {
			return nil, fmt.Errorf("price registry reported different prices (%s and %s) for the same token %s",
				fetchedPrices[i].Value, price, token)
		}

		tokenPrices[token] = fetchedPrices[i].Value
	}

	return tokenPrices, nil
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

// ensurePriceRegistrySynchronization ensures that the source price registry points to the same as the one configured on the onRamp.
// This is required since the price registry address on the onRamp can change over time.
func (r *ExecutionReportingPlugin) ensurePriceRegistrySynchronization(ctx context.Context) error {
	needPriceRegistryUpdate := false
	r.sourcePriceRegistryLock.RLock()
	priceRegistryAddress, err := r.onRampReader.SourcePriceRegistryAddress(ctx)
	if err != nil {
		r.sourcePriceRegistryLock.RUnlock()
		return fmt.Errorf("getting price registry from onramp: %w", err)
	}

	currentPriceRegistryAddress := cciptypes.Address("")
	if r.sourcePriceRegistry != nil {
		currentPriceRegistryAddress, err = r.sourcePriceRegistry.Address(ctx)
		if err != nil {
			return fmt.Errorf("get current priceregistry address: %w", err)
		}
	}

	needPriceRegistryUpdate = r.sourcePriceRegistry == nil || priceRegistryAddress != currentPriceRegistryAddress
	r.sourcePriceRegistryLock.RUnlock()
	if !needPriceRegistryUpdate {
		return nil
	}

	// Update the price registry if required.
	r.sourcePriceRegistryLock.Lock()
	defer r.sourcePriceRegistryLock.Unlock()

	// Price registry address changed or not initialized yet, updating source price registry.
	sourcePriceRegistry, err := r.sourcePriceRegistryProvider.NewPriceRegistryReader(ctx, priceRegistryAddress)
	if err != nil {
		return err
	}
	oldPriceRegistry := r.sourcePriceRegistry
	r.sourcePriceRegistry = sourcePriceRegistry
	// Close the old price registry
	if oldPriceRegistry != nil {
		if err1 := oldPriceRegistry.Close(); err1 != nil {
			r.lggr.Warnw("failed to close old price registry", "err", err1)
		}
	}
	return nil
}

// selectReportsToFillBatch returns the reports to fill the message limit. Single Commit Root contains exactly (Interval.Max - Interval.Min + 1) messages.
// We keep adding reports until we reach the message limit. Please see the tests for more examples and edge cases.
// unexpiredReports have to be sorted by Interval.Min. Otherwise, the batching logic will not be efficient,
// because it picks messages and execution states based on the report[0].Interval.Min - report[len-1].Interval.Max range.
// Having unexpiredReports not sorted properly will lead to fetching more messages and execution states to the memory than the messagesLimit provided.
// However, logs from LogPoller are returned ordered by (block_number, log_index), so it should preserve the order of Interval.Min.
// Single CommitRoot can have up to 256 messages, with current MessagesIterationStep of 1024, it means processing 4 CommitRoots at once.
func selectReportsToFillBatch(unexpiredReports []cciptypes.CommitStoreReport, messagesLimit uint64) ([]cciptypes.CommitStoreReport, int) {
	currentNumberOfMessages := uint64(0)
	nbReports := 0
	for _, report := range unexpiredReports {
		reportMsgCount := report.Interval.Max - report.Interval.Min + 1
		if currentNumberOfMessages+reportMsgCount > messagesLimit {
			break
		}
		currentNumberOfMessages += reportMsgCount
		nbReports++
	}
	return unexpiredReports[:nbReports], nbReports
}
