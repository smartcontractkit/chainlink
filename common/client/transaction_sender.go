package client

import (
	"context"
	"errors"
	"fmt"
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

// TxErrorClassifier - defines interface of a function that transforms raw RPC error into the SendTxReturnCode enum
// (e.g. Successful, Fatal, Retryable, etc.)
type TxErrorClassifier[TX any] func(tx TX, err error) SendTxReturnCode

type sendTxResult struct {
	Err        error
	ResultCode SendTxReturnCode
}

const sendTxQuorum = 0.7

// SendTxRPCClient - defines interface of an RPC used by TransactionSender to broadcast transaction
type SendTxRPCClient[TX any] interface {
	// SendTransaction errors returned should include name or other unique identifier of the RPC
	SendTransaction(ctx context.Context, tx TX) error
}

func NewTransactionSender[TX any, CHAIN_ID types.ID, RPC SendTxRPCClient[TX]](
	lggr logger.Logger,
	chainID CHAIN_ID,
	chainFamily string,
	multiNode *MultiNode[CHAIN_ID, RPC],
	txErrorClassifier TxErrorClassifier[TX],
	sendTxSoftTimeout time.Duration,
) *TransactionSender[TX, CHAIN_ID, RPC] {
	if sendTxSoftTimeout == 0 {
		sendTxSoftTimeout = QueryTimeout / 2
	}
	return &TransactionSender[TX, CHAIN_ID, RPC]{
		chainID:           chainID,
		chainFamily:       chainFamily,
		lggr:              logger.Sugared(lggr).Named("TransactionSender").With("chainID", chainID.String()),
		multiNode:         multiNode,
		txErrorClassifier: txErrorClassifier,
		sendTxSoftTimeout: sendTxSoftTimeout,
		chStop:            make(services.StopChan),
	}
}

type TransactionSender[TX any, CHAIN_ID types.ID, RPC SendTxRPCClient[TX]] struct {
	services.StateMachine
	chainID           CHAIN_ID
	chainFamily       string
	lggr              logger.SugaredLogger
	multiNode         *MultiNode[CHAIN_ID, RPC]
	txErrorClassifier TxErrorClassifier[TX]
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
func (txSender *TransactionSender[TX, CHAIN_ID, RPC]) SendTransaction(ctx context.Context, tx TX) (SendTxReturnCode, error) {
	txResults := make(chan sendTxResult)
	txResultsToReport := make(chan sendTxResult)
	primaryNodeWg := sync.WaitGroup{}

	if txSender.State() != "Started" {
		return Retryable, errors.New("TransactionSender not started")
	}

	healthyNodesNum := 0
	err := txSender.multiNode.DoAll(ctx, func(ctx context.Context, rpc RPC, isSendOnly bool) {
		if isSendOnly {
			txSender.wg.Add(1)
			go func() {
				defer txSender.wg.Done()
				// Send-only nodes' results are ignored as they tend to return false-positive responses.
				// Broadcast to them is necessary to speed up the propagation of TX in the network.
				_ = txSender.broadcastTxAsync(ctx, rpc, tx)
			}()
			return
		}

		// Primary Nodes
		healthyNodesNum++
		primaryNodeWg.Add(1)
		go func() {
			defer primaryNodeWg.Done()
			result := txSender.broadcastTxAsync(ctx, rpc, tx)
			select {
			case <-ctx.Done():
				return
			case txResults <- result:
			}

			select {
			case <-ctx.Done():
				return
			case txResultsToReport <- result:
			}
		}()
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
		return Retryable, err
	}

	txSender.wg.Add(1)
	go txSender.reportSendTxAnomalies(tx, txResultsToReport)

	return txSender.collectTxResults(ctx, tx, healthyNodesNum, txResults)
}

func (txSender *TransactionSender[TX, CHAIN_ID, RPC]) broadcastTxAsync(ctx context.Context, rpc RPC, tx TX) sendTxResult {
	txErr := rpc.SendTransaction(ctx, tx)
	txSender.lggr.Debugw("Node sent transaction", "tx", tx, "err", txErr)
	resultCode := txSender.txErrorClassifier(tx, txErr)
	if !slices.Contains(sendTxSuccessfulCodes, resultCode) {
		txSender.lggr.Warnw("RPC returned error", "tx", tx, "err", txErr)
	}
	return sendTxResult{Err: txErr, ResultCode: resultCode}
}

func (txSender *TransactionSender[TX, CHAIN_ID, RPC]) reportSendTxAnomalies(tx TX, txResults <-chan sendTxResult) {
	defer txSender.wg.Done()
	resultsByCode := sendTxResults{}
	// txResults eventually will be closed
	for txResult := range txResults {
		resultsByCode[txResult.ResultCode] = append(resultsByCode[txResult.ResultCode], txResult.Err)
	}

	_, _, criticalErr := aggregateTxResults(resultsByCode)
	if criticalErr != nil {
		txSender.lggr.Criticalw("observed invariant violation on SendTransaction", "tx", tx, "resultsByCode", resultsByCode, "err", criticalErr)
		PromMultiNodeInvariantViolations.WithLabelValues(txSender.chainFamily, txSender.chainID.String(), criticalErr.Error()).Inc()
	}
}

type sendTxResults map[SendTxReturnCode][]error

func aggregateTxResults(resultsByCode sendTxResults) (returnCode SendTxReturnCode, txResult error, err error) {
	severeCode, severeErrors, hasSevereErrors := findFirstIn(resultsByCode, sendTxSevereErrors)
	successCode, successResults, hasSuccess := findFirstIn(resultsByCode, sendTxSuccessfulCodes)
	if hasSuccess {
		// We assume that primary node would never report false positive txResult for a transaction.
		// Thus, if such case occurs it's probably due to misconfiguration or a bug and requires manual intervention.
		if hasSevereErrors {
			const errMsg = "found contradictions in nodes replies on SendTransaction: got success and severe error"
			// return success, since at least 1 node has accepted our broadcasted Tx, and thus it can now be included onchain
			return successCode, successResults[0], errors.New(errMsg)
		}

		// other errors are temporary - we are safe to return success
		return successCode, successResults[0], nil
	}

	if hasSevereErrors {
		return severeCode, severeErrors[0], nil
	}

	// return temporary error
	for code, result := range resultsByCode {
		return code, result[0], nil
	}

	err = fmt.Errorf("expected at least one response on SendTransaction")
	return Retryable, err, err
}

func (txSender *TransactionSender[TX, CHAIN_ID, RPC]) collectTxResults(ctx context.Context, tx TX, healthyNodesNum int, txResults <-chan sendTxResult) (SendTxReturnCode, error) {
	if healthyNodesNum == 0 {
		return Retryable, ErroringNodeError
	}
	ctx, cancel := txSender.chStop.Ctx(ctx)
	defer cancel()
	requiredResults := int(math.Ceil(float64(healthyNodesNum) * sendTxQuorum))
	errorsByCode := sendTxResults{}
	var softTimeoutChan <-chan time.Time
	var resultsCount int
loop:
	for {
		select {
		case <-ctx.Done():
			txSender.lggr.Debugw("Failed to collect of the results before context was done", "tx", tx, "errorsByCode", errorsByCode)
			return Retryable, ctx.Err()
		case result := <-txResults:
			errorsByCode[result.ResultCode] = append(errorsByCode[result.ResultCode], result.Err)
			resultsCount++
			if slices.Contains(sendTxSuccessfulCodes, result.ResultCode) || resultsCount >= requiredResults {
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
	returnCode, result, _ := aggregateTxResults(errorsByCode)
	return returnCode, result
}

func (txSender *TransactionSender[TX, CHAIN_ID, RPC]) Start(ctx context.Context) error {
	return txSender.StartOnce("TransactionSender", func() error {
		return nil
	})
}

func (txSender *TransactionSender[TX, CHAIN_ID, RPC]) Close() error {
	return txSender.StopOnce("TransactionSender", func() error {
		close(txSender.chStop)
		txSender.wg.Wait()
		return nil
	})
}

// findFirstIn - returns the first existing key and value for the slice of keys
func findFirstIn[K comparable, V any](set map[K]V, keys []K) (K, V, bool) {
	for _, k := range keys {
		if v, ok := set[k]; ok {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}
