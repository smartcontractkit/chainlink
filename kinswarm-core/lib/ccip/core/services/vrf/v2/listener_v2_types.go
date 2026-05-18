package v2

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	heaps "github.com/theodesp/go-heaps"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

type errPossiblyInsufficientFunds struct{}

func (errPossiblyInsufficientFunds) Error() string {
	return "Simulation errored, possibly insufficient funds. Request will remain unprocessed until funds are available"
}

type errBlockhashNotInStore struct{}

func (errBlockhashNotInStore) Error() string {
	return "Blockhash not in store"
}

type errProofVerificationFailed struct{}

func (errProofVerificationFailed) Error() string {
	return "Proof verification failed"
}

type fulfilledReqV2 struct {
	blockNumber uint64
	reqID       string
}

func (a fulfilledReqV2) Compare(b heaps.Item) int {
	a1 := a
	a2 := b.(fulfilledReqV2)
	switch {
	case a1.blockNumber > a2.blockNumber:
		return 1
	case a1.blockNumber < a2.blockNumber:
		return -1
	default:
		return 0
	}
}

type pendingRequest struct {
	confirmedAtBlock uint64
	req              RandomWordsRequested
	utcTimestamp     time.Time

	// used for exponential backoff when retrying
	attempts int
	lastTry  time.Time
}

type vrfPipelineResult struct {
	err error
	// maxFee indicates how much juels (link) or wei (ether) would be paid for the VRF request
	// if it were to be fulfilled at the maximum gas price (i.e gas lane gas price).
	maxFee *big.Int
	// fundsNeeded indicates a "minimum balance" in juels or wei that must be held in the
	// subscription's account in order to fulfill the request.
	fundsNeeded   *big.Int
	run           *pipeline.Run
	payload       string
	gasLimit      uint64
	req           pendingRequest
	proof         VRFProof
	reqCommitment RequestCommitment
}

// batchFulfillment contains all the information needed in order to
// perform a batch fulfillment operation on the batch VRF coordinator.
type batchFulfillment struct {
	proofs        []VRFProof
	commitments   []RequestCommitment
	totalGasLimit uint64
	runs          []*pipeline.Run
	reqIDs        []*big.Int
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
		if (currBatch.totalGasLimit + result.gasLimit) >= uint64(b.batchGasLimit) {
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
	err = sqlutil.TransactDataSource(ctx, lsn.ds, nil, func(tx sqlutil.DataSource) error {
		if err = lsn.pipelineRunner.InsertFinishedRuns(ctx, tx, batch.runs, true); err != nil {
			return fmt.Errorf("inserting finished pipeline runs: %w", err)
		}

		maxLink, maxEth := accumulateMaxLinkAndMaxEth(batch)
		var (
			txHashes    []common.Hash
			reqIDHashes []common.Hash
		)
		copy(txHashes, batch.txHashes)
		for _, reqID := range batch.reqIDs {
			reqIDHashes = append(reqIDHashes, common.BytesToHash(reqID.Bytes()))
		}
		ethTX, err = lsn.chain.TxManager().CreateTransaction(ctx, txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      lsn.batchCoordinator.Address(),
			EncodedPayload: payload,
			FeeLimit:       uint64(totalGasLimitBumped),
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
		if err != nil {
			return fmt.Errorf("create batch fulfillment eth transaction: %w", err)
		}

		return nil
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

// getReadyAndExpired filters out requests that are expired from the given pendingRequest slice
// and returns requests that are ready for processing.
func (lsn *listenerV2) getReadyAndExpired(l logger.Logger, reqs []pendingRequest) (ready []pendingRequest, expired []string) {
	for _, req := range reqs {
		// Check if we can ignore the request due to its age.
		if time.Now().UTC().Sub(req.utcTimestamp) >= lsn.job.VRFSpec.RequestTimeout {
			l.Infow("Request too old, dropping it",
				"reqID", req.req.RequestID().String(),
				"txHash", req.req.Raw().TxHash)
			expired = append(expired, req.req.RequestID().String())
			vrfcommon.IncDroppedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, vrfcommon.V2, vrfcommon.ReasonAge)
			continue
		}
		// we always check if the requests are already fulfilled prior to trying to fulfill them again
		ready = append(ready, req)
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
