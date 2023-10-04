package v2

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

// batchFulfillment contains all the information needed in order to
// perform a batch fulfillment operation on the batch VRF coordinator.
type batchFulfillment struct {
	proofs        []VRFProof
	commitments   []RequestCommitment
	totalGasLimit uint32
	runs          []*pipeline.Run
	reqIDs        []*big.Int
	lbs           []log.Broadcast
	maxFees       []*big.Int
	txHashes      []common.Hash
	fromAddress   common.Address
	version       vrfcommon.Version
}

func newBatchFulfillment(result vrfPipelineResult, fromAddress common.Address, version vrfcommon.Version) *batchFulfillment {
	return &batchFulfillment{
		proofs: []VRFProof{
			result.proof,
		},
		commitments: []RequestCommitment{
			result.reqCommitment,
		},
		totalGasLimit: result.gasLimit,
		runs: []*pipeline.Run{
			result.run,
		},
		reqIDs: []*big.Int{
			result.req.req.RequestID(),
		},
		lbs: []log.Broadcast{
			result.req.lb,
		},
		maxFees: []*big.Int{
			result.maxFee,
		},
		txHashes: []common.Hash{
			result.req.req.Raw().TxHash,
		},
		fromAddress: fromAddress,
		version:     version,
	}
}

// batchFulfillments manages many batchFulfillment objects.
// It makes organizing many runs into batches that respect the
// batchGasLimit easy via the addRun method.
type batchFulfillments struct {
	fulfillments  []*batchFulfillment
	batchGasLimit uint32
	currIndex     int
	version       vrfcommon.Version
}

func newBatchFulfillments(batchGasLimit uint32, version vrfcommon.Version) *batchFulfillments {
	return &batchFulfillments{
		fulfillments:  []*batchFulfillment{},
		batchGasLimit: batchGasLimit,
		currIndex:     0,
		version:       version,
	}
}

// addRun adds the given run to an existing batch, or creates a new
// batch if the batchGasLimit that has been configured was exceeded.
func (b *batchFulfillments) addRun(result vrfPipelineResult, fromAddress common.Address) {
	if len(b.fulfillments) == 0 {
		b.fulfillments = append(b.fulfillments, newBatchFulfillment(result, fromAddress, b.version))
	} else {
		currBatch := b.fulfillments[b.currIndex]
		if (currBatch.totalGasLimit + result.gasLimit) >= b.batchGasLimit {
			// don't add to curr batch, add new batch and increment index
			b.fulfillments = append(b.fulfillments, newBatchFulfillment(result, fromAddress, b.version))
			b.currIndex++
		} else {
			// we're okay on gas, add to current batch
			currBatch.proofs = append(currBatch.proofs, result.proof)
			currBatch.commitments = append(currBatch.commitments, result.reqCommitment)
			currBatch.totalGasLimit += result.gasLimit
			currBatch.runs = append(currBatch.runs, result.run)
			currBatch.reqIDs = append(currBatch.reqIDs, result.req.req.RequestID())
			currBatch.lbs = append(currBatch.lbs, result.req.lb)
			currBatch.maxFees = append(currBatch.maxFees, result.maxFee)
			currBatch.txHashes = append(currBatch.txHashes, result.req.req.Raw().TxHash)
		}
	}
}

func (lsn *listenerV2) processBatch(
	l logger.Logger,
	subID *big.Int,
	startBalanceNoReserveLink *big.Int,
	maxCallbackGasLimit uint32,
	batch *batchFulfillment,
	fromAddress common.Address,
) (processedRequestIDs []string) {
	start := time.Now()
	ctx, cancel := lsn.chStop.NewCtx()
	defer cancel()

	// Enqueue a single batch tx for requests that we're able to fulfill based on whether
	// they passed simulation or not.
	var (
		payload           []byte
		err               error
		txMetaSubID       *uint64
		txMetaGlobalSubID *string
	)

	if batch.version == vrfcommon.V2 {
		payload, err = batchCoordinatorV2ABI.Pack("fulfillRandomWords", ToV2Proofs(batch.proofs), ToV2Commitments(batch.commitments))
		if err != nil {
			// should never happen
			l.Errorw("Failed to pack batch fulfillRandomWords payload",
				"err", err, "proofs", batch.proofs, "commitments", batch.commitments)
			return
		}
		txMetaSubID = ptr(subID.Uint64())
	} else if batch.version == vrfcommon.V2Plus {
		payload, err = batchCoordinatorV2PlusABI.Pack("fulfillRandomWords", ToV2PlusProofs(batch.proofs), ToV2PlusCommitments(batch.commitments))
		if err != nil {
			// should never happen
			l.Errorw("Failed to pack batch fulfillRandomWords payload",
				"err", err, "proofs", batch.proofs, "commitments", batch.commitments)
			return
		}
		txMetaGlobalSubID = ptr(subID.String())
	} else {
		panic("batch version should be v2 or v2plus")
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
		"fromAddress", fromAddress,
		"linkBalance", startBalanceNoReserveLink,
		"totalGasLimitBumped", totalGasLimitBumped,
		"gasMultiplier", lsn.job.VRFSpec.BatchFulfillmentGasMultiplier,
	)
	ll.Info("Enqueuing batch fulfillment")
	var ethTX txmgr.Tx
	err = lsn.q.Transaction(func(tx pg.Queryer) error {
		if err = lsn.pipelineRunner.InsertFinishedRuns(batch.runs, true, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "inserting finished pipeline runs")
		}

		if err = lsn.logBroadcaster.MarkManyConsumed(batch.lbs, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "mark logs consumed")
		}

		maxLink, maxEth := accumulateMaxLinkAndMaxEth(batch)
		txHashes := []common.Hash{}
		copy(txHashes, batch.txHashes)
		reqIDHashes := []common.Hash{}
		for _, reqID := range batch.reqIDs {
			reqIDHashes = append(reqIDHashes, common.BytesToHash(reqID.Bytes()))
		}
		ethTX, err = lsn.txm.CreateTransaction(ctx, txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      lsn.batchCoordinator.Address(),
			EncodedPayload: payload,
			FeeLimit:       totalGasLimitBumped,
			Strategy:       txmgrcommon.NewSendEveryStrategy(),
			Meta: &txmgr.TxMeta{
				RequestIDs:      reqIDHashes,
				MaxLink:         &maxLink,
				MaxEth:          &maxEth,
				SubID:           txMetaSubID,
				GlobalSubID:     txMetaGlobalSubID,
				RequestTxHashes: txHashes,
			},
		})

		return errors.Wrap(err, "create batch fulfillment eth transaction")
	})
	if err != nil {
		ll.Errorw("Error enqueuing batch fulfillments, requeuing requests", "err", err)
		return
	}
	ll.Infow("Enqueued fulfillment", "ethTxID", ethTX.GetID())

	// mark requests as processed since the fulfillment has been successfully enqueued
	// to the txm.
	for _, reqID := range batch.reqIDs {
		processedRequestIDs = append(processedRequestIDs, reqID.String())
		vrfcommon.IncProcessedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, vrfcommon.V2)
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
				"reqID", req.req.RequestID().String(),
				"txHash", req.req.Raw().TxHash)
			lsn.markLogAsConsumed(req.lb)
			processed = append(processed, req.req.RequestID().String())
			vrfcommon.IncDroppedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, vrfcommon.V2, vrfcommon.ReasonAge)
			continue
		}

		// This check to see if the log was consumed needs to be in the same
		// goroutine as the mark consumed to avoid processing duplicates.
		consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(req.lb)
		if err != nil {
			// Do not process for now, retry on next iteration.
			l.Errorw("Could not determine if log was already consumed",
				"reqID", req.req.RequestID().String(),
				"txHash", req.req.Raw().TxHash,
				"error", err)
		} else if consumed {
			processed = append(processed, req.req.RequestID().String())
		} else {
			unconsumed = append(unconsumed, req)
		}
	}
	return
}

func batchFulfillmentGasEstimate(
	batchSize uint64,
	maxCallbackGasLimit uint32,
	gasMultiplier float64,
) uint32 {
	return uint32(
		gasMultiplier * float64((uint64(maxCallbackGasLimit)+400_000)+batchSize*BatchFulfillmentIterationGasCost),
	)
}

func accumulateMaxLinkAndMaxEth(batch *batchFulfillment) (maxLinkStr string, maxEthStr string) {
	maxLink := big.NewInt(0)
	maxEth := big.NewInt(0)
	for i := range batch.commitments {
		if batch.commitments[i].VRFVersion == vrfcommon.V2 {
			// v2 always bills in link
			maxLink.Add(maxLink, batch.maxFees[i])
		} else {
			// v2plus can bill in link or eth, depending on the commitment
			if batch.commitments[i].NativePayment() {
				maxEth.Add(maxEth, batch.maxFees[i])
			} else {
				maxLink.Add(maxLink, batch.maxFees[i])
			}
		}
	}
	return maxLink.String(), maxEth.String()
}
