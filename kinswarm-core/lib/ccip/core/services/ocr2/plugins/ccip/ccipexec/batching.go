package ccipexec

import (
	"context"
	"fmt"
	"math/big"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/statuschecker"
)

// Batching strategies
const (
	BestEffortBatchingStrategyID = uint32(0)
	ZKOverflowBatchingStrategyID = uint32(1)
)

type BatchContext struct {
	report                     commitReportWithSendRequests
	inflight                   []InflightInternalExecutionReport
	inflightAggregateValue     *big.Int
	lggr                       logger.Logger
	availableDataLen           int
	availableGas               uint64
	expectedNonces             map[cciptypes.Address]uint64
	sendersNonce               map[cciptypes.Address]uint64
	sourceTokenPricesUSD       map[cciptypes.Address]*big.Int
	destTokenPricesUSD         map[cciptypes.Address]*big.Int
	gasPrice                   *big.Int
	sourceToDestToken          map[cciptypes.Address]cciptypes.Address
	aggregateTokenLimit        *big.Int
	tokenDataRemainingDuration time.Duration
	tokenDataWorker            tokendata.Worker
	gasPriceEstimator          prices.GasPriceEstimatorExec
	destWrappedNative          cciptypes.Address
	offchainConfig             cciptypes.ExecOffchainConfig
}

type BatchingStrategy interface {
	BuildBatch(ctx context.Context, batchCtx *BatchContext) ([]ccip.ObservedMessage, []messageExecStatus)
	GetBatchingStrategyID() uint32
}

type BestEffortBatchingStrategy struct{}

type ZKOverflowBatchingStrategy struct {
	statuschecker statuschecker.CCIPTransactionStatusChecker
}

func NewBatchingStrategy(batchingStrategyID uint32, statusChecker statuschecker.CCIPTransactionStatusChecker) (BatchingStrategy, error) {
	var batchingStrategy BatchingStrategy
	switch batchingStrategyID {
	case BestEffortBatchingStrategyID:
		batchingStrategy = &BestEffortBatchingStrategy{}
	case ZKOverflowBatchingStrategyID:
		batchingStrategy = &ZKOverflowBatchingStrategy{
			statuschecker: statusChecker,
		}
	default:
		return nil, errors.Errorf("unknown batching strategy ID %d", batchingStrategyID)
	}
	return batchingStrategy, nil
}

func (s *BestEffortBatchingStrategy) GetBatchingStrategyID() uint32 {
	return BestEffortBatchingStrategyID
}

// BestEffortBatchingStrategy is a batching strategy that tries to batch as many messages as possible (up to certain limits).
func (s *BestEffortBatchingStrategy) BuildBatch(
	ctx context.Context,
	batchCtx *BatchContext,
) ([]ccip.ObservedMessage, []messageExecStatus) {
	batchBuilder := newBatchBuildContainer(len(batchCtx.report.sendRequestsWithMeta))
	for _, msg := range batchCtx.report.sendRequestsWithMeta {
		msgLggr := batchCtx.lggr.With("messageID", hexutil.Encode(msg.MessageID[:]), "seqNr", msg.SequenceNumber)
		status, messageMaxGas, tokenData, msgValue, err := performCommonChecks(ctx, batchCtx, msg, msgLggr)

		if err != nil {
			return []ccip.ObservedMessage{}, []messageExecStatus{}
		}

		if status.shouldBeSkipped() {
			batchBuilder.skip(msg, status)
			continue
		}

		updateBatchContext(batchCtx, msg, messageMaxGas, msgValue, msgLggr)
		batchBuilder.addToBatch(msg, tokenData)
	}
	return batchBuilder.batch, batchBuilder.statuses
}

func (bs *ZKOverflowBatchingStrategy) GetBatchingStrategyID() uint32 {
	return ZKOverflowBatchingStrategyID
}

// ZKOverflowBatchingStrategy is a batching strategy for ZK chains overflowing under certain conditions.
// It is a simple batching strategy that only allows one message to be added to the batch.
// TXM is used to perform the ZK check: if the message failed the check, it will be skipped.
func (bs ZKOverflowBatchingStrategy) BuildBatch(
	ctx context.Context,
	batchCtx *BatchContext,
) ([]ccip.ObservedMessage, []messageExecStatus) {
	batchBuilder := newBatchBuildContainer(len(batchCtx.report.sendRequestsWithMeta))
	inflightSeqNums := getInflightSeqNums(batchCtx.inflight)

	for _, msg := range batchCtx.report.sendRequestsWithMeta {
		msgId := hexutil.Encode(msg.MessageID[:])
		msgLggr := batchCtx.lggr.With("messageID", msgId, "seqNr", msg.SequenceNumber)

		// Check if msg is inflight
		if exists := inflightSeqNums.Contains(msg.SequenceNumber); exists {
			// Message is inflight, skip it
			msgLggr.Infow("Skipping message - already inflight", "message", msgId)
			batchBuilder.skip(msg, SkippedInflight)
			continue
		}
		// Message is not inflight, continue with checks
		// Check if the messsage is overflown using TXM
		statuses, count, err := bs.statuschecker.CheckMessageStatus(ctx, msgId)
		if err != nil {
			batchBuilder.skip(msg, TXMCheckError)
			continue
		}

		msgLggr.Infow("TXM check result", "statuses", statuses, "count", count)

		if len(statuses) == 0 {
			// No status found for message = first time we see it
			msgLggr.Infow("No status found for message - proceeding with checks", "message", msgId)
		} else {
			// Status(es) found for message = check if any of them is final to decide if we should add it to the batch
			hasFatalStatus := false
			for _, s := range statuses {
				if s == types.Fatal {
					msgLggr.Infow("Skipping message - found a fatal TXM status", "message", msgId)
					batchBuilder.skip(msg, TXMFatalStatus)
					hasFatalStatus = true
					break
				}
			}
			if hasFatalStatus {
				continue
			}
			msgLggr.Infow("No fatal status found for message - proceeding with checks", "message", msgId)
		}

		status, messageMaxGas, tokenData, msgValue, err := performCommonChecks(ctx, batchCtx, msg, msgLggr)

		if err != nil {
			return []ccip.ObservedMessage{}, []messageExecStatus{}
		}

		if status.shouldBeSkipped() {
			batchBuilder.skip(msg, status)
			continue
		}

		updateBatchContext(batchCtx, msg, messageMaxGas, msgValue, msgLggr)
		msgLggr.Infow("Adding message to batch", "message", msgId)
		batchBuilder.addToBatch(msg, tokenData)

		// Batch size is limited to 1 for ZK Overflow chains
		break
	}
	return batchBuilder.batch, batchBuilder.statuses
}

func performCommonChecks(
	ctx context.Context,
	batchCtx *BatchContext,
	msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta,
	msgLggr logger.Logger,
) (messageStatus, uint64, [][]byte, *big.Int, error) {
	if msg.Executed {
		msgLggr.Infow("Skipping message - already executed")
		return AlreadyExecuted, 0, nil, nil, nil
	}

	if len(msg.Data) > batchCtx.availableDataLen {
		msgLggr.Infow("Skipping message - insufficient remaining batch data length", "msgDataLen", len(msg.Data), "availableBatchDataLen", batchCtx.availableDataLen)
		return InsufficientRemainingBatchDataLength, 0, nil, nil, nil
	}

	messageMaxGas, err1 := calculateMessageMaxGas(
		msg.GasLimit,
		len(batchCtx.report.sendRequestsWithMeta),
		len(msg.Data),
		len(msg.TokenAmounts),
	)
	if err1 != nil {
		msgLggr.Errorw("Skipping message - message max gas calculation error", "err", err1)
		return MessageMaxGasCalcError, 0, nil, nil, nil
	}

	// Check sufficient gas in batch
	if batchCtx.availableGas < messageMaxGas {
		msgLggr.Infow("Skipping message - insufficient remaining batch gas limit", "availableGas", batchCtx.availableGas, "messageMaxGas", messageMaxGas)
		return InsufficientRemainingBatchGas, 0, nil, nil, nil
	}

	if _, ok := batchCtx.expectedNonces[msg.Sender]; !ok {
		nonce, ok1 := batchCtx.sendersNonce[msg.Sender]
		if !ok1 {
			msgLggr.Errorw("Skipping message - missing nonce", "sender", msg.Sender)
			return MissingNonce, 0, nil, nil, nil
		}
		batchCtx.expectedNonces[msg.Sender] = nonce + 1
	}

	// Check expected nonce is valid for sequenced messages.
	// Sequenced messages have non-zero nonces.
	if msg.Nonce > 0 && msg.Nonce != batchCtx.expectedNonces[msg.Sender] {
		msgLggr.Warnw("Skipping message - invalid nonce", "have", msg.Nonce, "want", batchCtx.expectedNonces[msg.Sender])
		return InvalidNonce, 0, nil, nil, nil
	}

	msgValue, err1 := aggregateTokenValue(batchCtx.lggr, batchCtx.destTokenPricesUSD, batchCtx.sourceToDestToken, msg.TokenAmounts)
	if err1 != nil {
		msgLggr.Errorw("Skipping message - aggregate token value compute error", "err", err1)
		return AggregateTokenValueComputeError, 0, nil, nil, nil
	}

	// if token limit is smaller than message value skip message
	if tokensLeft, hasCapacity := hasEnoughTokens(batchCtx.aggregateTokenLimit, msgValue, batchCtx.inflightAggregateValue); !hasCapacity {
		msgLggr.Warnw("Skipping message - aggregate token limit exceeded", "aggregateTokenLimit", tokensLeft.String(), "msgValue", msgValue.String())
		return AggregateTokenLimitExceeded, 0, nil, nil, nil
	}

	tokenData, elapsed, err1 := getTokenDataWithTimeout(ctx, msg, batchCtx.tokenDataRemainingDuration, batchCtx.tokenDataWorker)
	batchCtx.tokenDataRemainingDuration -= elapsed
	if err1 != nil {
		if errors.Is(err1, tokendata.ErrNotReady) {
			msgLggr.Warnw("Skipping message - token data not ready", "err", err1)
			return TokenDataNotReady, 0, nil, nil, nil
		}
		msgLggr.Errorw("Skipping message - token data fetch error", "err", err1)
		return TokenDataFetchError, 0, nil, nil, nil
	}

	dstWrappedNativePrice, exists := batchCtx.destTokenPricesUSD[batchCtx.destWrappedNative]
	if !exists {
		msgLggr.Errorw("Skipping message - token not in destination token prices", "token", batchCtx.destWrappedNative)
		return TokenNotInDestTokenPrices, 0, nil, nil, nil
	}

	// calculating the source chain fee, dividing by 1e18 for denomination.
	// For example:
	// FeeToken=link; FeeTokenAmount=1e17 i.e. 0.1 link, price is 6e18 USD/link (1 USD = 1e18),
	// availableFee is 1e17*6e18/1e18 = 6e17 = 0.6 USD
	sourceFeeTokenPrice, exists := batchCtx.sourceTokenPricesUSD[msg.FeeToken]
	if !exists {
		msgLggr.Errorw("Skipping message - token not in source token prices", "token", msg.FeeToken)
		return TokenNotInSrcTokenPrices, 0, nil, nil, nil
	}

	// Fee boosting
	execCostUsd, err1 := batchCtx.gasPriceEstimator.EstimateMsgCostUSD(batchCtx.gasPrice, dstWrappedNativePrice, msg)
	if err1 != nil {
		msgLggr.Errorw("Failed to estimate message cost USD", "err", err1)
		return "", 0, nil, nil, errors.New("failed to estimate message cost USD")
	}

	availableFee := big.NewInt(0).Mul(msg.FeeTokenAmount, sourceFeeTokenPrice)
	availableFee = availableFee.Div(availableFee, big.NewInt(1e18))
	availableFeeUsd := waitBoostedFee(time.Since(msg.BlockTimestamp), availableFee, batchCtx.offchainConfig.RelativeBoostPerWaitHour)
	if availableFeeUsd.Cmp(execCostUsd) < 0 {
		msgLggr.Infow(
			"Skipping message - insufficient remaining fee",
			"availableFeeUsd", availableFeeUsd,
			"execCostUsd", execCostUsd,
			"sourceBlockTimestamp", msg.BlockTimestamp,
			"waitTime", time.Since(msg.BlockTimestamp),
			"boost", batchCtx.offchainConfig.RelativeBoostPerWaitHour,
		)
		return InsufficientRemainingFee, 0, nil, nil, nil
	}

	return SuccesfullyValidated, messageMaxGas, tokenData, msgValue, nil
}

// getTokenDataWithCappedLatency gets the token data for the provided message.
// Stops and returns an error if more than allowedWaitingTime is passed.
func getTokenDataWithTimeout(
	ctx context.Context,
	msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta,
	timeout time.Duration,
	tokenDataWorker tokendata.Worker,
) ([][]byte, time.Duration, error) {
	if len(msg.TokenAmounts) == 0 {
		return nil, 0, nil
	}

	ctxTimeout, cf := context.WithTimeout(ctx, timeout)
	defer cf()
	tStart := time.Now()
	tokenData, err := tokenDataWorker.GetMsgTokenData(ctxTimeout, msg)
	tDur := time.Since(tStart)
	return tokenData, tDur, err
}

func getProofData(
	ctx context.Context,
	sourceReader ccipdata.OnRampReader,
	interval cciptypes.CommitStoreInterval,
) (sendReqsInRoot []cciptypes.EVM2EVMMessageWithTxMeta, leaves [][32]byte, tree *merklemulti.Tree[[32]byte], err error) {
	// We don't need to double-check if logs are finalized because we already checked that in the Commit phase.
	sendReqs, err := sourceReader.GetSendRequestsBetweenSeqNums(ctx, interval.Min, interval.Max, false)
	if err != nil {
		return nil, nil, nil, err
	}

	if err1 := validateSendRequests(sendReqs, interval); err1 != nil {
		return nil, nil, nil, err1
	}

	leaves = make([][32]byte, 0, len(sendReqs))
	for _, req := range sendReqs {
		leaves = append(leaves, req.Hash)
	}
	tree, err = merklemulti.NewTree(hashutil.NewKeccak(), leaves)
	if err != nil {
		return nil, nil, nil, err
	}
	return sendReqs, leaves, tree, nil
}

func validateSendRequests(sendReqs []cciptypes.EVM2EVMMessageWithTxMeta, interval cciptypes.CommitStoreInterval) error {
	if len(sendReqs) == 0 {
		return fmt.Errorf("could not find any requests in the provided interval %v", interval)
	}

	gotInterval := cciptypes.CommitStoreInterval{
		Min: sendReqs[0].SequenceNumber,
		Max: sendReqs[0].SequenceNumber,
	}

	for _, req := range sendReqs[1:] {
		if req.SequenceNumber < gotInterval.Min {
			gotInterval.Min = req.SequenceNumber
		}
		if req.SequenceNumber > gotInterval.Max {
			gotInterval.Max = req.SequenceNumber
		}
	}

	if (gotInterval.Min != interval.Min) || (gotInterval.Max != interval.Max) {
		return fmt.Errorf("interval %v is not the expected %v", gotInterval, interval)
	}
	return nil
}

func getInflightSeqNums(inflight []InflightInternalExecutionReport) mapset.Set[uint64] {
	seqNums := mapset.NewSet[uint64]()
	for _, report := range inflight {
		for _, msg := range report.messages {
			seqNums.Add(msg.SequenceNumber)
		}
	}
	return seqNums
}

func aggregateTokenValue(lggr logger.Logger, destTokenPricesUSD map[cciptypes.Address]*big.Int, sourceToDest map[cciptypes.Address]cciptypes.Address, tokensAndAmount []cciptypes.TokenAmount) (*big.Int, error) {
	sum := big.NewInt(0)
	for i := 0; i < len(tokensAndAmount); i++ {
		price, ok := destTokenPricesUSD[sourceToDest[tokensAndAmount[i].Token]]
		if !ok {
			// If we don't have a price for the token, we will assume it's worth 0.
			lggr.Infof("No price for token %s, assuming 0", tokensAndAmount[i].Token)
			continue
		}
		sum.Add(sum, new(big.Int).Quo(new(big.Int).Mul(price, tokensAndAmount[i].Amount), big.NewInt(1e18)))
	}
	return sum, nil
}

func updateBatchContext(
	batchCtx *BatchContext,
	msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta,
	messageMaxGas uint64,
	msgValue *big.Int,
	msgLggr logger.Logger) {
	batchCtx.availableGas -= messageMaxGas
	batchCtx.availableDataLen -= len(msg.Data)
	batchCtx.aggregateTokenLimit.Sub(batchCtx.aggregateTokenLimit, msgValue)
	if msg.Nonce > 0 {
		batchCtx.expectedNonces[msg.Sender] = msg.Nonce + 1
	}

	msgLggr.Infow(
		"Message successfully added to execution batch",
		"nonce", msg.Nonce,
		"sender", msg.Sender,
		"value", msgValue,
		"availableAggrTokenLimit", batchCtx.aggregateTokenLimit,
		"availableGas", batchCtx.availableGas,
		"availableDataLen", batchCtx.availableDataLen,
	)
}

func hasEnoughTokens(tokenLimit *big.Int, msgValue *big.Int, inflightValue *big.Int) (*big.Int, bool) {
	tokensLeft := big.NewInt(0).Sub(tokenLimit, inflightValue)
	return tokensLeft, tokensLeft.Cmp(msgValue) >= 0
}

func buildExecutionReportForMessages(
	msgsInRoot []cciptypes.EVM2EVMMessageWithTxMeta,
	tree *merklemulti.Tree[[32]byte],
	commitInterval cciptypes.CommitStoreInterval,
	observedMessages []ccip.ObservedMessage,
) (cciptypes.ExecReport, error) {
	innerIdxs := make([]int, 0, len(observedMessages))
	var messages []cciptypes.EVM2EVMMessage
	var offchainTokenData [][][]byte
	for _, observedMessage := range observedMessages {
		if observedMessage.SeqNr < commitInterval.Min || observedMessage.SeqNr > commitInterval.Max {
			// We only return messages from a single root (the root of the first message).
			continue
		}
		innerIdx := int(observedMessage.SeqNr - commitInterval.Min)
		if innerIdx >= len(msgsInRoot) || innerIdx < 0 {
			return cciptypes.ExecReport{}, fmt.Errorf("invalid inneridx SeqNr=%d IntervalMin=%d msgsInRoot=%d",
				observedMessage.SeqNr, commitInterval.Min, len(msgsInRoot))
		}
		messages = append(messages, msgsInRoot[innerIdx].EVM2EVMMessage)
		offchainTokenData = append(offchainTokenData, observedMessage.TokenData)
		innerIdxs = append(innerIdxs, innerIdx)
	}

	merkleProof, err := tree.Prove(innerIdxs)
	if err != nil {
		return cciptypes.ExecReport{}, err
	}

	// any capped proof will have length <= this one, so we reuse it to avoid proving inside loop, and update later if changed
	return cciptypes.ExecReport{
		Messages:          messages,
		Proofs:            merkleProof.Hashes,
		ProofFlagBits:     abihelpers.ProofFlagsToBits(merkleProof.SourceFlags),
		OffchainTokenData: offchainTokenData,
	}, nil
}

// Validates the given message observations do not exceed the committed sequence numbers
// in the commitStoreReader.
func validateSeqNumbers(serviceCtx context.Context, commitStore ccipdata.CommitStoreReader, observedMessages []ccip.ObservedMessage) error {
	nextMin, err := commitStore.GetExpectedNextSequenceNumber(serviceCtx)
	if err != nil {
		return err
	}
	// observedMessages are always sorted by SeqNr and never empty, so it's safe to take last element
	maxSeqNumInBatch := observedMessages[len(observedMessages)-1].SeqNr

	if maxSeqNumInBatch >= nextMin {
		return errors.Errorf("Cannot execute uncommitted seq num. nextMin %v, seqNums %v", nextMin, observedMessages)
	}
	return nil
}

// Gets the commit report from the saved logs for a given sequence number.
func getCommitReportForSeqNum(ctx context.Context, commitStoreReader ccipdata.CommitStoreReader, seqNum uint64) (cciptypes.CommitStoreReport, error) {
	acceptedReports, err := commitStoreReader.GetCommitReportMatchingSeqNum(ctx, seqNum, 0)
	if err != nil {
		return cciptypes.CommitStoreReport{}, err
	}

	if len(acceptedReports) == 0 {
		return cciptypes.CommitStoreReport{}, errors.Errorf("seq number not committed")
	}

	return acceptedReports[0].CommitStoreReport, nil
}

type messageStatus string

const (
	SuccesfullyValidated                 messageStatus = "successfully_validated"
	AlreadyExecuted                      messageStatus = "already_executed"
	SenderAlreadySkipped                 messageStatus = "sender_already_skipped"
	MessageMaxGasCalcError               messageStatus = "message_max_gas_calc_error"
	InsufficientRemainingBatchDataLength messageStatus = "insufficient_remaining_batch_data_length"
	InsufficientRemainingBatchGas        messageStatus = "insufficient_remaining_batch_gas"
	MissingNonce                         messageStatus = "missing_nonce"
	InvalidNonce                         messageStatus = "invalid_nonce"
	AggregateTokenValueComputeError      messageStatus = "aggregate_token_value_compute_error"
	AggregateTokenLimitExceeded          messageStatus = "aggregate_token_limit_exceeded"
	TokenDataNotReady                    messageStatus = "token_data_not_ready"
	TokenDataFetchError                  messageStatus = "token_data_fetch_error"
	TokenNotInDestTokenPrices            messageStatus = "token_not_in_dest_token_prices"
	TokenNotInSrcTokenPrices             messageStatus = "token_not_in_src_token_prices"
	InsufficientRemainingFee             messageStatus = "insufficient_remaining_fee"
	AddedToBatch                         messageStatus = "added_to_batch"
	TXMCheckError                        messageStatus = "txm_check_error"
	TXMFatalStatus                       messageStatus = "txm_fatal_status"
	SkippedInflight                      messageStatus = "skipped_inflight"
)

func (m messageStatus) shouldBeSkipped() bool {
	return m != SuccesfullyValidated
}

type messageExecStatus struct {
	SeqNr     uint64
	MessageId string
	Status    messageStatus
}

func newMessageExecState(seqNr uint64, messageId cciptypes.Hash, status messageStatus) messageExecStatus {
	return messageExecStatus{
		SeqNr:     seqNr,
		MessageId: hexutil.Encode(messageId[:]),
		Status:    status,
	}
}

type batchBuildContainer struct {
	batch    []ccip.ObservedMessage
	statuses []messageExecStatus
}

func newBatchBuildContainer(capacity int) *batchBuildContainer {
	return &batchBuildContainer{
		batch:    make([]ccip.ObservedMessage, 0, capacity),
		statuses: make([]messageExecStatus, 0, capacity),
	}
}

func (m *batchBuildContainer) skip(msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, status messageStatus) {
	m.addState(msg, status)
}

func (m *batchBuildContainer) addToBatch(msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenData [][]byte) {
	m.addState(msg, AddedToBatch)
	m.batch = append(m.batch, ccip.NewObservedMessage(msg.SequenceNumber, tokenData))
}

func (m *batchBuildContainer) addState(msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, state messageStatus) {
	m.statuses = append(m.statuses, newMessageExecState(msg.SequenceNumber, msg.MessageID, state))
}
