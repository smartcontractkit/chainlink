package client

import (
	"context"
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type TransactionSender[TX any] interface {
	SendTransaction(ctx context.Context, tx TX) (SendTxReturnCode, error)
}

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
	SendTransaction(ctx context.Context, tx TX) error
}

func NewTransactionSender[TX any, CHAIN_ID types.ID, RPC SendTxRPCClient[TX]](
	lggr logger.Logger,
	chainID CHAIN_ID,
	chainFamily string,
	multiNode *MultiNode[CHAIN_ID, RPC],
	txErrorClassifier TxErrorClassifier[TX],
	sendTxSoftTimeout time.Duration,
) TransactionSender[TX] {
	if sendTxSoftTimeout == 0 {
		sendTxSoftTimeout = QueryTimeout / 2
	}
	return &transactionSender[TX, CHAIN_ID, RPC]{
		chainID:           chainID,
		chainFamily:       chainFamily,
		lggr:              logger.Sugared(lggr).Named("TransactionSender").With("chainID", chainID.String()),
		multiNode:         multiNode,
		txErrorClassifier: txErrorClassifier,
		sendTxSoftTimeout: sendTxSoftTimeout,
	}
}

type transactionSender[TX any, CHAIN_ID types.ID, RPC SendTxRPCClient[TX]] struct {
	services.StateMachine
	chainID           CHAIN_ID
	chainFamily       string
	lggr              logger.SugaredLogger
	multiNode         *MultiNode[CHAIN_ID, RPC]
	txErrorClassifier TxErrorClassifier[TX]
	sendTxSoftTimeout time.Duration // defines max waiting time from first response til responses evaluation

	// TODO: add start/ stop methods. Start doesn't need to do much.
	// TODO: Stop should stop sending transactions, and close chStop to stop collecting results, reporting/ etc.
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
func (txSender *transactionSender[TX, CHAIN_ID, RPC]) SendTransaction(ctx context.Context, tx TX) (SendTxReturnCode, error) {
	txResults := make(chan sendTxResult, len(txSender.multiNode.primaryNodes))
	txResultsToReport := make(chan sendTxResult, len(txSender.multiNode.primaryNodes))
	primaryWg := sync.WaitGroup{}

	var err error
	ok := txSender.multiNode.IfNotStopped(func() {
		err = txSender.multiNode.DoAll(ctx, func(ctx context.Context, rpc RPC, isSendOnly bool) bool {
			if isSendOnly {
				// Use multiNode wg to ensure transactions are done sending before multinode shuts down
				txSender.multiNode.wg.Add(1)
				fmt.Println("Calling send only rpc SendTransaction()")
				go func() {
					defer txSender.multiNode.wg.Done()
					// Send-only nodes' results are ignored as they tend to return false-positive responses.
					// Broadcast to them is necessary to speed up the propagation of TX in the network.
					_ = txSender.broadcastTxAsync(ctx, rpc, tx)
				}()
				return true
			}

			// Primary Nodes
			primaryWg.Add(1)
			go func() {
				defer primaryWg.Done()
				result := txSender.broadcastTxAsync(ctx, rpc, tx)
				txResultsToReport <- result
				txResults <- result
			}()
			return true
		})
		if err != nil {
			primaryWg.Wait()
			close(txResultsToReport)
			close(txResults)
			return
		}

		// This needs to be done in parallel so the reporting knows when it's done (when the channel is closed)
		txSender.multiNode.wg.Add(1)
		go func() {
			defer txSender.multiNode.wg.Done()
			primaryWg.Wait()
			close(txResultsToReport)
			close(txResults)
		}()

		txSender.multiNode.wg.Add(1)
		go txSender.reportSendTxAnomalies(tx, txResultsToReport)
	})
	if !ok {
		return 0, fmt.Errorf("aborted while broadcasting tx - MultiNode is stopped: %w", context.Canceled)
	}
	if err != nil {
		return 0, err
	}

	return txSender.collectTxResults(ctx, tx, len(txSender.multiNode.primaryNodes), txResults)
}

func (txSender *transactionSender[TX, CHAIN_ID, RPC]) broadcastTxAsync(ctx context.Context, rpc RPC, tx TX) sendTxResult {
	txErr := rpc.SendTransaction(ctx, tx)
	txSender.lggr.Debugw("Node sent transaction", "tx", tx, "err", txErr)
	resultCode := txSender.txErrorClassifier(tx, txErr)
	if !slices.Contains(sendTxSuccessfulCodes, resultCode) {
		txSender.lggr.Warnw("RPC returned error", "tx", tx, "err", txErr)
	}
	return sendTxResult{Err: txErr, ResultCode: resultCode}
}

func (txSender *transactionSender[TX, CHAIN_ID, RPC]) reportSendTxAnomalies(tx TX, txResults <-chan sendTxResult) {
	defer txSender.multiNode.wg.Done()
	resultsByCode := sendTxErrors{}
	// txResults eventually will be closed
	for txResult := range txResults {
		resultsByCode[txResult.ResultCode] = append(resultsByCode[txResult.ResultCode], txResult.Err)
	}

	_, _, criticalErr := aggregateTxResults(resultsByCode)
	if criticalErr != nil {
		txSender.lggr.Criticalw("observed invariant violation on SendTransaction", "tx", tx, "resultsByCode", resultsByCode, "err", criticalErr)
		txSender.SvcErrBuffer.Append(criticalErr)
		PromMultiNodeInvariantViolations.WithLabelValues(txSender.chainFamily, txSender.chainID.String(), criticalErr.Error()).Inc()
	}
}

type sendTxErrors map[SendTxReturnCode][]error

func aggregateTxResults(resultsByCode sendTxErrors) (returnCode SendTxReturnCode, txResult error, err error) {
	// TODO: Modify this to return the corresponding returnCode with the error
	severeCode, severeErrors, hasSevereErrors := findFirstIn(resultsByCode, sendTxSevereErrors)
	successCode, successResults, hasSuccess := findFirstIn(resultsByCode, sendTxSuccessfulCodes)
	if hasSuccess {
		// We assume that primary node would never report false positive txResult for a transaction.
		// Thus, if such case occurs it's probably due to misconfiguration or a bug and requires manual intervention.
		if hasSevereErrors {
			const errMsg = "found contradictions in nodes replies on SendTransaction: got success and severe error"
			// return success, since at least 1 node has accepted our broadcasted Tx, and thus it can now be included onchain
			return successCode, successResults[0], fmt.Errorf(errMsg)
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
	return 0, err, err
}

func (txSender *transactionSender[TX, CHAIN_ID, RPC]) collectTxResults(ctx context.Context, tx TX, healthyNodesNum int, txResults <-chan sendTxResult) (SendTxReturnCode, error) {
	if healthyNodesNum == 0 {
		// TODO: Should we return fatal here, retryable, or 0?
		return 0, ErroringNodeError
	}
	// combine context and stop channel to ensure we stop, when signal received
	ctx, cancel := txSender.chStop.Ctx(ctx)
	defer cancel()
	requiredResults := int(math.Ceil(float64(healthyNodesNum) * sendTxQuorum))
	errorsByCode := sendTxErrors{}
	var softTimeoutChan <-chan time.Time
	var resultsCount int
loop:
	for {
		select {
		case <-ctx.Done():
			txSender.lggr.Debugw("Failed to collect of the results before context was done", "tx", tx, "errorsByCode", errorsByCode)
			return 0, ctx.Err()
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
