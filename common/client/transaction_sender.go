package client

import (
	"context"
	"errors"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

var (
	// PromMultiNodeInvariantViolations reports violation of our assumptions
	PromMultiNodeInvariantViolations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "multi_node_invariant_violations",
		Help: "The number of invariant violations",
	}, []string{"network", "chainId", "invariant"})
)

type SendTxResult interface {
	Code() SendTxReturnCode
	Error() error
}

const sendTxQuorum = 0.7

// SendTxRPCClient - defines interface of an RPC used by TransactionSender to broadcast transaction
type SendTxRPCClient[TX any, RESULT SendTxResult] interface {
	// SendTransaction errors returned should include name or other unique identifier of the RPC
	SendTransaction(ctx context.Context, tx TX) RESULT
}

func NewTransactionSender[TX any, RESULT SendTxResult, CHAIN_ID types.ID, RPC SendTxRPCClient[TX, RESULT]](
	lggr logger.Logger,
	chainID CHAIN_ID,
	chainFamily string,
	multiNode *MultiNode[CHAIN_ID, RPC],
	newResult func(err error) RESULT,
	sendTxSoftTimeout time.Duration,
) *TransactionSender[TX, RESULT, CHAIN_ID, RPC] {
	if sendTxSoftTimeout == 0 {
		sendTxSoftTimeout = QueryTimeout / 2
	}
	return &TransactionSender[TX, RESULT, CHAIN_ID, RPC]{
		chainID:           chainID,
		chainFamily:       chainFamily,
		lggr:              logger.Sugared(lggr).Named("TransactionSender").With("chainID", chainID.String()),
		multiNode:         multiNode,
		newResult:         newResult,
		sendTxSoftTimeout: sendTxSoftTimeout,
		chStop:            make(services.StopChan),
	}
}

type TransactionSender[TX any, RESULT SendTxResult, CHAIN_ID types.ID, RPC SendTxRPCClient[TX, RESULT]] struct {
	services.StateMachine
	chainID           CHAIN_ID
	chainFamily       string
	lggr              logger.SugaredLogger
	multiNode         *MultiNode[CHAIN_ID, RPC]
	newResult         func(err error) RESULT
	sendTxSoftTimeout time.Duration // defines max waiting time from first response til responses evaluation

	wg     sync.WaitGroup // waits for all reporting goroutines to finish
	chStop services.StopChan
}

// SendTransaction - broadcasts transaction to all the send-only and primary nodes in MultiNode.
// A returned nil or error does not guarantee that the transaction will or won't be included. Additional checks must be
// performed to determine the final state.
//
// Send-only nodes' results are ignored as they tend to return false-positive responses. Broadcast to them is necessary
// to speed up the propagation of TX in the network.
//
// Handling of primary nodes' results consists of collection and aggregation.
// In the collection step, we gather as many results as possible while minimizing waiting time. This operation succeeds
// on one of the following conditions:
// * Received at least one success
// * Received at least one result and `sendTxSoftTimeout` expired
// * Received results from the sufficient number of nodes defined by sendTxQuorum.
// The aggregation is based on the following conditions:
// * If there is at least one success - returns success
// * If there is at least one terminal error - returns terminal error
// * If there is both success and terminal error - returns success and reports invariant violation
// * Otherwise, returns any (effectively random) of the errors.
func (txSender *TransactionSender[TX, RESULT, CHAIN_ID, RPC]) SendTransaction(ctx context.Context, tx TX) RESULT {
	var result RESULT
	if !txSender.IfStarted(func() {
		txResults := make(chan RESULT)
		txResultsToReport := make(chan RESULT)
		primaryNodeWg := sync.WaitGroup{}

		healthyNodesNum := 0
		err := txSender.multiNode.DoAll(ctx, func(ctx context.Context, rpc RPC, isSendOnly bool) {
			if isSendOnly {
				txSender.wg.Add(1)
				go func(ctx context.Context) {
					ctx, cancel := txSender.chStop.Ctx(context.WithoutCancel(ctx))
					defer cancel()
					defer txSender.wg.Done()
					// Send-only nodes' results are ignored as they tend to return false-positive responses.
					// Broadcast to them is necessary to speed up the propagation of TX in the network.
					_ = txSender.broadcastTxAsync(ctx, rpc, tx)
				}(ctx)
				return
			}

			// Primary Nodes
			healthyNodesNum++
			primaryNodeWg.Add(1)
			go func(ctx context.Context) {
				ctx, cancel := txSender.chStop.Ctx(context.WithoutCancel(ctx))
				defer cancel()
				defer primaryNodeWg.Done()
				r := txSender.broadcastTxAsync(ctx, rpc, tx)
				select {
				case <-ctx.Done():
					txSender.lggr.Debugw("Failed to send tx results", "err", ctx.Err())
					return
				case txResults <- r:
				}

				select {
				case <-ctx.Done():
					txSender.lggr.Debugw("Failed to send tx results to report", "err", ctx.Err())
					return
				case txResultsToReport <- r:
				}
			}(ctx)
		})

		// This needs to be done in parallel so the reporting knows when it's done (when the channel is closed)
		txSender.wg.Add(1)
		go func() {
			defer txSender.wg.Done()
			primaryNodeWg.Wait()
			close(txResultsToReport)
			close(txResults)
		}()

		if err != nil {
			result = txSender.newResult(err)
			return
		}

		txSender.wg.Add(1)
		go txSender.reportSendTxAnomalies(ctx, tx, txResultsToReport)

		result = txSender.collectTxResults(ctx, tx, healthyNodesNum, txResults)
	}) {
		result = txSender.newResult(errors.New("TransactionSender not started"))
	}

	return result
}

func (txSender *TransactionSender[TX, RESULT, CHAIN_ID, RPC]) broadcastTxAsync(ctx context.Context, rpc RPC, tx TX) RESULT {
	result := rpc.SendTransaction(ctx, tx)
	txSender.lggr.Debugw("Node sent transaction", "tx", tx, "err", result.Error())
	if !slices.Contains(sendTxSuccessfulCodes, result.Code()) && ctx.Err() == nil {
		txSender.lggr.Warnw("RPC returned error", "tx", tx, "err", result.Error())
	}
	return result
}

func (txSender *TransactionSender[TX, RESULT, CHAIN_ID, RPC]) reportSendTxAnomalies(ctx context.Context, tx TX, txResults <-chan RESULT) {
	defer txSender.wg.Done()
	resultsByCode := sendTxResults[RESULT]{}
	// txResults eventually will be closed
	for txResult := range txResults {
		resultsByCode[txResult.Code()] = append(resultsByCode[txResult.Code()], txResult)
	}

	_, criticalErr := aggregateTxResults[RESULT](resultsByCode)
	if criticalErr != nil && ctx.Err() == nil {
		txSender.lggr.Criticalw("observed invariant violation on SendTransaction", "tx", tx, "resultsByCode", resultsByCode, "err", criticalErr)
		PromMultiNodeInvariantViolations.WithLabelValues(txSender.chainFamily, txSender.chainID.String(), criticalErr.Error()).Inc()
	}
}

type sendTxResults[RESULT any] map[SendTxReturnCode][]RESULT

func aggregateTxResults[RESULT any](resultsByCode sendTxResults[RESULT]) (result RESULT, criticalErr error) {
	severeErrors, hasSevereErrors := findFirstIn(resultsByCode, sendTxSevereErrors)
	successResults, hasSuccess := findFirstIn(resultsByCode, sendTxSuccessfulCodes)
	if hasSuccess {
		// We assume that primary node would never report false positive txResult for a transaction.
		// Thus, if such case occurs it's probably due to misconfiguration or a bug and requires manual intervention.
		if hasSevereErrors {
			const errMsg = "found contradictions in nodes replies on SendTransaction: got success and severe error"
			// return success, since at least 1 node has accepted our broadcasted Tx, and thus it can now be included onchain
			return successResults[0], errors.New(errMsg)
		}

		// other errors are temporary - we are safe to return success
		return successResults[0], nil
	}

	if hasSevereErrors {
		return severeErrors[0], nil
	}

	// return temporary error
	for _, r := range resultsByCode {
		return r[0], nil
	}

	criticalErr = errors.New("expected at least one response on SendTransaction")
	return result, criticalErr
}

func (txSender *TransactionSender[TX, RESULT, CHAIN_ID, RPC]) collectTxResults(ctx context.Context, tx TX, healthyNodesNum int, txResults <-chan RESULT) RESULT {
	if healthyNodesNum == 0 {
		return txSender.newResult(ErroringNodeError)
	}
	requiredResults := int(math.Ceil(float64(healthyNodesNum) * sendTxQuorum))
	errorsByCode := sendTxResults[RESULT]{}
	var softTimeoutChan <-chan time.Time
	var resultsCount int
loop:
	for {
		select {
		case <-ctx.Done():
			txSender.lggr.Debugw("Failed to collect of the results before context was done", "tx", tx, "errorsByCode", errorsByCode)
			return txSender.newResult(ctx.Err())
		case r := <-txResults:
			errorsByCode[r.Code()] = append(errorsByCode[r.Code()], r)
			resultsCount++
			if slices.Contains(sendTxSuccessfulCodes, r.Code()) || resultsCount >= requiredResults {
				break loop
			}
		case <-softTimeoutChan:
			txSender.lggr.Debugw("Send Tx soft timeout expired - returning responses we've collected so far", "tx", tx, "resultsCount", resultsCount, "requiredResults", requiredResults)
			break loop
		}

		if softTimeoutChan == nil {
			tm := time.NewTimer(txSender.sendTxSoftTimeout)
			softTimeoutChan = tm.C
			// we are fine with stopping timer at the end of function
			//nolint
			defer tm.Stop()
		}
	}

	// ignore critical error as it's reported in reportSendTxAnomalies
	result, _ := aggregateTxResults(errorsByCode)
	txSender.lggr.Debugw("Collected results", "errorsByCode", errorsByCode, "result", result)
	return result
}

func (txSender *TransactionSender[TX, RESULT, CHAIN_ID, RPC]) Start(ctx context.Context) error {
	return txSender.StartOnce("TransactionSender", func() error {
		return nil
	})
}

func (txSender *TransactionSender[TX, RESULT, CHAIN_ID, RPC]) Close() error {
	return txSender.StopOnce("TransactionSender", func() error {
		txSender.lggr.Debug("Closing TransactionSender")
		close(txSender.chStop)
		txSender.wg.Wait()
		return nil
	})
}

// findFirstIn - returns the first existing key and value for the slice of keys
func findFirstIn[K comparable, V any](set map[K]V, keys []K) (V, bool) {
	for _, k := range keys {
		if v, ok := set[k]; ok {
			return v, true
		}
	}
	var zeroV V
	return zeroV, false
}
