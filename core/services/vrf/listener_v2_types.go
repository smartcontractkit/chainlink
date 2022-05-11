package vrf

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

// batchFulfillment contains all the information needed in order to
// perform a batch fulfillment operation on the batch VRF coordinator.
type batchFulfillment struct {
	proofs        []batch_vrf_coordinator_v2.VRFTypesProof
	commitments   []batch_vrf_coordinator_v2.VRFTypesRequestCommitment
	totalGasLimit uint64
	runs          []*pipeline.Run
	reqIDs        []*big.Int
	lbs           []log.Broadcast
	maxLinks      []interface{}
}

func newBatchFulfillment(result vrfPipelineResult) *batchFulfillment {
	return &batchFulfillment{
		proofs: []batch_vrf_coordinator_v2.VRFTypesProof{
			batch_vrf_coordinator_v2.VRFTypesProof(result.proof),
		},
		commitments: []batch_vrf_coordinator_v2.VRFTypesRequestCommitment{
			batch_vrf_coordinator_v2.VRFTypesRequestCommitment(result.reqCommitment),
		},
		totalGasLimit: result.gasLimit,
		runs: []*pipeline.Run{
			&result.run,
		},
		reqIDs: []*big.Int{
			result.req.req.RequestId,
		},
		lbs: []log.Broadcast{
			result.req.lb,
		},
		maxLinks: []interface{}{
			result.maxLink,
		},
	}
}

// batchFulfillments manages many batchFulfillment objects.
// It makes organizing many runs into batches that respect the
// batchGasLimit easy via the addRun method.
type batchFulfillments struct {
	fulfillments  []*batchFulfillment
	batchGasLimit uint64
	currIndex     int
}

func newBatchFulfillments(batchGasLimit uint64) *batchFulfillments {
	return &batchFulfillments{
		fulfillments:  []*batchFulfillment{},
		batchGasLimit: batchGasLimit,
		currIndex:     0,
	}
}

// addRun adds the given run to an existing batch, or creates a new
// batch if the batchGasLimit that has been configured was exceeded.
func (b *batchFulfillments) addRun(result vrfPipelineResult) {
	if len(b.fulfillments) == 0 {
		b.fulfillments = append(b.fulfillments, newBatchFulfillment(result))
	} else {
		currBatch := b.fulfillments[b.currIndex]
		if (currBatch.totalGasLimit + result.gasLimit) >= b.batchGasLimit {
			// don't add to curr batch, add new batch and increment index
			b.fulfillments = append(b.fulfillments, newBatchFulfillment(result))
			b.currIndex++
		} else {
			// we're okay on gas, add to current batch
			currBatch.proofs = append(currBatch.proofs, batch_vrf_coordinator_v2.VRFTypesProof(result.proof))
			currBatch.commitments = append(currBatch.commitments, batch_vrf_coordinator_v2.VRFTypesRequestCommitment(result.reqCommitment))
			currBatch.totalGasLimit += result.gasLimit
			currBatch.runs = append(currBatch.runs, &result.run)
			currBatch.reqIDs = append(currBatch.reqIDs, result.req.req.RequestId)
			currBatch.lbs = append(currBatch.lbs, result.req.lb)
			currBatch.maxLinks = append(currBatch.maxLinks, result.maxLink)
		}
	}
}

func (lsn *listenerV2) processBatch(
	l logger.Logger,
	subID uint64,
	fromAddress common.Address,
	startBalanceNoReserveLink *big.Int,
	maxCallbackGasLimit uint64,
	batch *batchFulfillment,
) (processedRequestIDs []string) {
	start := time.Now()

	// Enqueue a single batch tx for requests that we're able to fulfill based on whether
	// they passed simulation or not.
	payload, err := batchCoordinatorV2ABI.Pack("fulfillRandomWords", batch.proofs, batch.commitments)
	if err != nil {
		// should never happen
		l.Errorw("Failed to pack batch fulfillRandomWords payload",
			"err", err, "proofs", batch.proofs, "commitments", batch.commitments)
		return
	}

	// Bump the total gas limit by a bit so that we account for the overhead of the batch
	// contract's calling.
	totalGasLimitBumped := batchFulfillmentGasEstimate(
		uint64(len(batch.proofs)),
		maxCallbackGasLimit,
		float64(lsn.job.VRFSpec.BatchFulfillmentGasMultiplier),
	)
	ll := l.With("numRequestsInBatch", len(batch.reqIDs),
		"requestIDs", batch.reqIDs,
		"batchSumGasLimit", batch.totalGasLimit,
		"linkBalance", startBalanceNoReserveLink,
		"totalGasLimitBumped", totalGasLimitBumped,
		"gasMultiplier", lsn.job.VRFSpec.BatchFulfillmentGasMultiplier,
	)
	ll.Info("Enqueuing batch fulfillment")
	var ethTX txmgr.EthTx
	err = lsn.q.Transaction(func(tx pg.Queryer) error {
		if err = lsn.pipelineRunner.InsertFinishedRuns(batch.runs, true, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "inserting finished pipeline runs")
		}

		if err = lsn.logBroadcaster.MarkManyConsumed(batch.lbs, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "mark logs consumed")
		}

		maxLinkStr := bigmath.Accumulate(batch.maxLinks).String()
		reqIDHashes := []common.Hash{}
		for _, reqID := range batch.reqIDs {
			reqIDHashes = append(reqIDHashes, common.BytesToHash(reqID.Bytes()))
		}
		ethTX, err = lsn.txm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:      fromAddress,
			ToAddress:        lsn.batchCoordinator.Address(),
			EncodedPayload:   payload,
			GasLimit:         totalGasLimitBumped,
			MinConfirmations: null.Uint32From(uint32(lsn.cfg.MinRequiredOutgoingConfirmations())),
			Strategy:         txmgr.NewSendEveryStrategy(),
			Meta: &txmgr.EthTxMeta{
				RequestIDs: reqIDHashes,
				MaxLink:    &maxLinkStr,
				SubID:      &subID,
			},
		}, pg.WithQueryer(tx))

		return errors.Wrap(err, "create batch fulfillment eth transaction")
	})
	if err != nil {
		ll.Errorw("Error enqueuing batch fulfillments, requeuing requests", "err", err)
		return
	}
	ll.Infow("Enqueued fulfillment", "ethTxID", ethTX.ID)

	// mark requests as processed since the fulfillment has been successfully enqueued
	// to the txm.
	for _, reqID := range batch.reqIDs {
		processedRequestIDs = append(processedRequestIDs, reqID.String())
		incProcessedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2)
	}

	ll.Infow("Successfully enqueued batch", "duration", time.Since(start))

	return
}

// getUnconsumed returns the requests in the given slice that are not expired
// and not marked consumed in the log broadcaster.
func (lsn *listenerV2) getUnconsumed(l logger.Logger, reqs []pendingRequest) (unconsumed []pendingRequest, processed []string) {
	for _, req := range reqs {
		// Check if we can ignore the request due to its age.
		if time.Now().UTC().Sub(req.utcTimestamp) >= lsn.job.VRFSpec.RequestTimeout {
			l.Infow("Request too old, dropping it",
				"reqID", req.req.RequestId.String(),
				"txHash", req.req.Raw.TxHash)
			lsn.markLogAsConsumed(req.lb)
			processed = append(processed, req.req.RequestId.String())
			incDroppedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2, reasonAge)
			continue
		}

		// This check to see if the log was consumed needs to be in the same
		// goroutine as the mark consumed to avoid processing duplicates.
		consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(req.lb)
		if err != nil {
			// Do not process for now, retry on next iteration.
			l.Errorw("Could not determine if log was already consumed",
				"reqID", req.req.RequestId.String(),
				"txHash", req.req.Raw.TxHash,
				"error", err)
		} else if consumed {
			processed = append(processed, req.req.RequestId.String())
		} else {
			unconsumed = append(unconsumed, req)
		}
	}
	return
}

func batchFulfillmentGasEstimate(
	batchSize uint64,
	maxCallbackGasLimit uint64,
	gasMultiplier float64,
) uint64 {
	return uint64(
		gasMultiplier * float64((maxCallbackGasLimit+400_000)+batchSize*BatchFulfillmentIterationGasCost),
	)
}
