package vault

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"maps"
	"slices"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

const (
	blobBroadcastTimeout        = 2 * time.Second
	maxConcurrentBlobBroadcasts = 10
)

// marshalPendingQueueBlobPayload encodes pending queue items for OCR3.1 BroadcastBlob.
// Always marshals as PendingQueueBlobItems. Single items are wire-compatible with StoredPendingQueueItem,
// so the non-optimizations unmarshal path can decode them without knowing the new type.
// Batch items wrap each StoredPendingQueueItem (payload + ID) inside an Any.
func marshalPendingQueueBlobPayload(items []*vaultcommon.StoredPendingQueueItem) ([]byte, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("empty pending queue blob payload")
	}
	if len(items) == 1 {
		return proto.Marshal(&vaultcommon.PendingQueueBlobItems{
			Items: []*anypb.Any{items[0].GetItem()},
			Id:    items[0].GetId(),
		})
	}
	anyItems := make([]*anypb.Any, len(items))
	for i, item := range items {
		anyItem, err := anypb.New(item)
		if err != nil {
			return nil, fmt.Errorf("could not wrap StoredPendingQueueItem in Any: %w", err)
		}
		anyItems[i] = anyItem
	}
	return proto.Marshal(&vaultcommon.PendingQueueBlobItems{Items: anyItems, IsBatch: true})
}

// unmarshalPendingQueueBlob decodes a BroadcastBlob payload into one or more StoredPendingQueueItems.
// Batch blobs (is_batch=true) contain each item wrapped as an Any inside PendingQueueBlobItems.
// Non-batch blobs are wire-compatible with StoredPendingQueueItem and unmarshalled directly.
func unmarshalPendingQueueBlob(blob []byte) ([]*vaultcommon.StoredPendingQueueItem, error) {
	pqbi := &vaultcommon.PendingQueueBlobItems{}
	if err := proto.Unmarshal(blob, pqbi); err == nil && pqbi.IsBatch {
		items := make([]*vaultcommon.StoredPendingQueueItem, len(pqbi.GetItems()))
		for i, anyItem := range pqbi.GetItems() {
			item := &vaultcommon.StoredPendingQueueItem{}
			if err := anyItem.UnmarshalTo(item); err != nil {
				return nil, fmt.Errorf("could not unmarshal batch item %d: %w", i, err)
			}
			items[i] = item
		}
		return items, nil
	}
	// Non-batch: PendingQueueBlobItems with is_batch=false is wire-compatible with StoredPendingQueueItem
	single := &vaultcommon.StoredPendingQueueItem{}
	if err := proto.Unmarshal(blob, single); err != nil {
		return nil, err
	}
	return []*vaultcommon.StoredPendingQueueItem{single}, nil
}

// pendingQueueBlobPack is the blob payload side of Observation for the local pending queue.
type pendingQueueBlobPack struct {
	blobPayloads    [][]byte
	blobPayloadIDs  [][]string
	packedItemCount int
	truncated       bool
}

func (r *ReportingPlugin) flushBatch(
	ctx context.Context,
	seqNr uint64,
	currentBatch []*vaultcommon.StoredPendingQueueItem,
	out pendingQueueBlobPack,
	maxBlobBytes int,
	maxBlobHandleCount int,
) (pendingQueueBlobPack, []*vaultcommon.StoredPendingQueueItem, error) {
	if len(currentBatch) == 0 {
		return out, currentBatch, nil
	}
	payload, mErr := marshalPendingQueueBlobPayload(currentBatch)
	if mErr != nil {
		return out, currentBatch, fmt.Errorf("could not marshal pending queue blob payload: %w", mErr)
	}
	if len(payload) > maxBlobBytes {
		return out, currentBatch, fmt.Errorf("pending queue blob payload exceeds max size (%d > %d)", len(payload), maxBlobBytes)
	}
	ids := make([]string, 0, len(currentBatch))
	for _, it := range currentBatch {
		ids = append(ids, it.Id)
		r.lifecycle.RecordBlobBroadcasting(it.Id, seqNr, time.Now())
	}
	out.packedItemCount += len(currentBatch)
	out.blobPayloads = append(out.blobPayloads, payload)
	out.blobPayloadIDs = append(out.blobPayloadIDs, ids)
	currentBatch = currentBatch[:0]
	if len(out.blobPayloads) >= maxBlobHandleCount {
		out.truncated = true
		r.lggr.Warnw("Observed local queue exceeds batch size limit, truncating",
			"queueSize", len(out.blobPayloads),
			"batchSizeLimit", maxBlobHandleCount)
		r.metrics.trackQueueOverflow(ctx, len(out.blobPayloads), maxBlobHandleCount)
	}
	return out, currentBatch, nil
}

// prepareObservationPendingQueueBlobs packs local-queue requests into OCR3.1 blob payloads
// (PendingQueueBlobItems), capped by byte size per blob and by maxBlobHandleCount handles.
func (r *ReportingPlugin) prepareObservationPendingQueueBlobs(
	ctx context.Context,
	seqNr uint64,
	localQueueItems []*vaulttypes.Request,
	pendingQueueHasID map[string]bool,
	maxBlobBytes int,
	maxBlobHandleCount int,
) (pendingQueueBlobPack, error) {
	var out pendingQueueBlobPack
	var currentBatch []*vaultcommon.StoredPendingQueueItem

	for i := 0; i < len(localQueueItems); i++ {
		queueItem := localQueueItems[i]
		if pendingQueueHasID[queueItem.ID()] {
			continue
		}

		anyMsg, err := anypb.New(queueItem.Payload)
		if err != nil {
			return pendingQueueBlobPack{}, fmt.Errorf("could not marshal request payload to Any: %w", err)
		}

		singleItem := &vaultcommon.StoredPendingQueueItem{
			Id:   queueItem.ID(),
			Item: anyMsg,
		}

		candidate := append(slices.Clone(currentBatch), singleItem)
		payload, err := marshalPendingQueueBlobPayload(candidate)
		if err != nil {
			return pendingQueueBlobPack{}, err
		}

		if len(payload) > maxBlobBytes {
			if len(currentBatch) == 0 {
				return pendingQueueBlobPack{}, fmt.Errorf("single pending queue item exceeds max blob payload size (%d > %d)", len(payload), maxBlobBytes)
			}
			// Current batch is full; flush it and retry the same item on the next iteration.
			var ferr error
			flushOut, flushBatch, ferr := r.flushBatch(ctx, seqNr, currentBatch, out, maxBlobBytes, maxBlobHandleCount)
			if ferr != nil {
				return pendingQueueBlobPack{}, ferr
			}
			out = flushOut
			currentBatch = flushBatch
			if out.truncated {
				break
			}
			i--
			continue
		}
		currentBatch = candidate
	}

	var err error
	out, _, err = r.flushBatch(ctx, seqNr, currentBatch, out, maxBlobBytes, maxBlobHandleCount)
	if err != nil {
		return pendingQueueBlobPack{}, err
	}

	return out, nil
}

// prepareLegacyObservationPendingQueueBlobs emits one PendingQueueBlobItems blob per local-queue request.
func (r *ReportingPlugin) prepareLegacyObservationPendingQueueBlobs(
	ctx context.Context,
	seqNr uint64,
	localQueueItems []*vaulttypes.Request,
	pendingQueueHasID map[string]bool,
	maxBlobHandleCount int,
) (pendingQueueBlobPack, error) {
	var out pendingQueueBlobPack

	for _, queueItem := range localQueueItems {
		if pendingQueueHasID[queueItem.ID()] {
			continue
		}

		anyMsg, err := anypb.New(queueItem.Payload)
		if err != nil {
			return pendingQueueBlobPack{}, fmt.Errorf("could not marshal request payload to Any: %w", err)
		}

		item := &vaultcommon.StoredPendingQueueItem{
			Id:   queueItem.ID(),
			Item: anyMsg,
		}

		itemb, err := marshalPendingQueueBlobPayload([]*vaultcommon.StoredPendingQueueItem{item})
		if err != nil {
			return pendingQueueBlobPack{}, fmt.Errorf("could not marshal pending queue item: %w", err)
		}

		if len(out.blobPayloads) >= maxBlobHandleCount {
			out.truncated = true
			r.lggr.Warnw("Observed local queue exceeds batch size limit, truncating",
				"queueSize", len(out.blobPayloads),
				"batchSizeLimit", maxBlobHandleCount)
			r.metrics.trackQueueOverflow(ctx, len(out.blobPayloads), maxBlobHandleCount)
			break
		}

		out.blobPayloads = append(out.blobPayloads, itemb)
		out.blobPayloadIDs = append(out.blobPayloadIDs, []string{item.Id})
		out.packedItemCount++
		r.lifecycle.RecordBlobBroadcasting(item.Id, seqNr, time.Now())
	}

	return out, nil
}

func (r *ReportingPlugin) optimizationsEnabled(ctx context.Context) bool {
	return gateAllows(ctx, r.lggr, r.cfg.VaultOptimizationsEnabled, "VaultOptimizationsEnabled")
}

type pendingQueueStore interface {
	WritePendingQueue(ctx context.Context, pending []*vaultcommon.StoredPendingQueueItem) error
}

// appendPendingQueueObservations appends one Observation per store-backed pending queue item.
// When applyWireCap is true, stops before obspb exceeds maxObservationBytes.
func (r *ReportingPlugin) appendPendingQueueObservations(
	ctx context.Context,
	seqNr uint64,
	readKV *KVStore,
	currentPendingQueueItems []*vaultcommon.StoredPendingQueueItem,
	obspb *vaultcommon.Observations,
	applyWireCap bool,
) []string {
	ids := make([]string, 0, len(currentPendingQueueItems))
	for _, req := range currentPendingQueueItems {
		o := &vaultcommon.Observation{
			Id: req.Id,
		}

		payload, err := req.Item.UnmarshalNew()
		if err != nil {
			r.lggr.Errorw("failed to unmarshal request payload", "id", req.Id, "error", err)
			continue
		}

		switch tp := payload.(type) {
		case *vaultcommon.GetSecretsRequest:
			r.observeGetSecrets(ctx, readKV, tp, o)
		case *vaultcommon.CreateSecretsRequest:
			r.observeCreateSecrets(ctx, readKV, tp, o)
		case *vaultcommon.UpdateSecretsRequest:
			r.observeUpdateSecrets(ctx, readKV, tp, o)
		case *vaultcommon.DeleteSecretsRequest:
			r.observeDeleteSecrets(ctx, readKV, tp, o)
		case *vaultcommon.ListSecretIdentifiersRequest:
			r.observeListSecretIdentifiers(ctx, readKV, tp, o)
		default:
			r.lggr.Errorw("unknown request type, skipping...", "requestType", fmt.Sprintf("%T", payload), "id", req.Id)
			continue
		}

		obspb.Observations = append(obspb.Observations, o)
		if applyWireCap && proto.Size(obspb) > r.maxObservationBytes {
			obspb.Observations = obspb.Observations[:len(obspb.Observations)-1]
			r.lggr.Warnw("observation proto would exceed max observation bytes; stopping pending-queue observation pack",
				"seqNr", seqNr,
				"id", req.Id,
				"maxObservationBytes", r.maxObservationBytes,
				"packedObservationCount", len(obspb.Observations),
				"pendingQueueItemCount", len(currentPendingQueueItems),
			)
			break
		}
		ids = append(ids, req.Id)
		r.lifecycle.RecordObservedOutcome(req.Id, seqNr, time.Now())
	}
	return ids
}

// broadcastBlobPayloads broadcasts each payload as a blob in parallel to reduce
// Observation() latency (shortening this phase helps the OCR round finish within
// DeltaProgress). Each call is given a 2-second timeout so that a single slow
// broadcast cannot stall the entire batch. No more than 10 broadcasts are allowed
// in flight at a time. Individual broadcast failures are logged and skipped rather
// than aborting the entire observation, so that one problematic payload does not
// prevent the remaining items from being observed. Context cancellation/deadline
// errors on the parent context are propagated immediately so that expired rounds
// fail fast.
func (r *ReportingPlugin) broadcastBlobPayloads(
	ctx context.Context,
	fetcher ocr3_1types.BlobBroadcastFetcher,
	seqNr uint64,
	payloads [][]byte,
	requestIDs [][]string,
) ([][]byte, error) {
	results := make([][]byte, len(payloads))

	start := time.Now()
	defer func() {
		r.lggr.Debugw("observation blob broadcast finished", "seqNr", seqNr, "blobCount", len(payloads), "elapsed", time.Since(start))
	}()

	var g errgroup.Group
	g.SetLimit(maxConcurrentBlobBroadcasts)
	for i, payload := range payloads {
		requestID := requestIDs[i]
		g.Go(func() error {
			broadcastCtx, cancel := context.WithTimeout(ctx, blobBroadcastTimeout)
			defer cancel()

			blobHandle, err := fetcher.BroadcastBlob(broadcastCtx, payload, ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: seqNr + 2})
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				r.lggr.Warnw("failed to broadcast pending queue item as blob, skipping",
					"seqNr", seqNr,
					"requestID", requestID[0],
					"err", err)
				return nil
			}

			blobHandleBytes, err := r.marshalBlob(blobHandle)
			if err != nil {
				r.lggr.Warnw("failed to marshal blob handle, skipping",
					"seqNr", seqNr,
					"requestID", requestID[0],
					"err", err)
				return nil
			}

			results[i] = blobHandleBytes
			for _, id := range requestID {
				r.lifecycle.RecordBlobBroadcasted(id, seqNr, time.Now())
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	filtered := make([][]byte, 0, len(results))
	for _, item := range results {
		if item != nil {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (r *ReportingPlugin) stateTransitionPendingQueue(ctx context.Context, seqNr uint64, store pendingQueueStore, obs map[uint8]*vaultcommon.Observations, blobFetcher ocr3_1types.BlobFetcher) error {
	// Step 1: Create a map of id -> sha -> count.
	idToShaToCount := map[string]map[string]int{}
	oidsToIDs := map[uint8][]string{} // for debugging only
	shaToItem := map[string]*vaultcommon.StoredPendingQueueItem{}
	for oid, o := range obs {
		shaSeenForOracle := map[string]bool{}
		for _, pqi := range o.PendingQueueItems {
			bh, err := r.unmarshalBlob(pqi)
			if err != nil {
				r.lggr.Errorw("failed to unmarshal blob handle from pending queue item", "error", err, "item", pqi)
				continue
			}

			blob, err := blobFetcher.FetchBlob(ctx, bh)
			if err != nil {
				r.lggr.Errorw("failed to fetch blob for pending queue item", "error", err, "item", pqi)
				continue
			}

			items, err := unmarshalPendingQueueBlob(blob)
			if err != nil {
				r.lggr.Errorw("failed to unmarshal blob into pending queue item(s)", "error", err, "item", pqi)
				continue
			}

			for _, i := range items {
				oidsToIDs[oid] = append(oidsToIDs[oid], i.Id)

				sha, err := shaForProto(i)
				if err != nil {
					r.lggr.Errorw("failed to compute sha for pending queue item", "error", err, "item", pqi)
					continue
				}

				if shaSeenForOracle[sha] {
					r.lggr.Warnw("duplicate sha found for oracle, skipping...", "oracleID", oid, "sha", sha, "item", pqi, "blobHandle", bh)
					continue
				}

				shaSeenForOracle[sha] = true

				shaToItem[sha] = i

				if _, ok := idToShaToCount[i.Id]; !ok {
					idToShaToCount[i.Id] = map[string]int{}
				}
				idToShaToCount[i.Id][sha]++
			}
		}
	}

	r.lggr.Debugw("processing pending queue", "oracleIDsToPendingQueueIDs", oidsToIDs)

	// Step 2: Generate the aggregated pending queue.
	// Any observation that has been seen F+1 times is kept.
	keptItems := []*vaultcommon.StoredPendingQueueItem{}
	// We don't need to sort here since keptItems are sorted later.
	for id, shaToCount := range idToShaToCount {
		maxCount := 0
		chosenSha := ""

		// Identify the sha with the most count.
		// We sort the sha to ensure deterministic iteration within an ID.
		// This can matter in a tie-breaker situation where two items
		// have the same count.
		for _, sha := range slices.Sorted(maps.Keys(shaToCount)) {
			count := shaToCount[sha]

			if count > maxCount {
				maxCount = count
				chosenSha = sha
			}
		}

		if maxCount >= r.onchainCfg.F+1 {
			keptItems = append(keptItems, shaToItem[chosenSha])
		} else {
			r.lggr.Warnw("pending queue item did not reach F+1 consensus, skipping...", "maxCount", maxCount, "id", id, "idToShaToCount", idToShaToCount, "F", r.onchainCfg.F)
		}
	}

	// Step 3: Generate the salt that we'll use to sort the list deterministically.
	salt := []byte{}
	for _, oid := range slices.Sorted(maps.Keys(obs)) {
		salt = append(salt, obs[oid].SortNonce...)
	}

	// Step 4: Sort the kept items by sha(id || salt)
	// The salt ensures that items are ordered randomly each time, preventing
	// front-running and dishonest nodes from manipulating the order of items in the pending queue.
	slices.SortFunc(keptItems, func(i *vaultcommon.StoredPendingQueueItem, j *vaultcommon.StoredPendingQueueItem) int {
		return bytes.Compare(sortKey(i.Id, salt), sortKey(j.Id, salt))
	})

	r.metrics.trackPendingQueueWrittenSize(ctx, len(keptItems))
	r.lggr.Infow("pending queue items persisted to storage", "seqNr", seqNr, "writtenCount", len(keptItems))

	now := time.Now()
	for _, it := range keptItems {
		r.lifecycle.RecordWrittenToPendingQueue(ctx, it.Id, seqNr, now)
	}

	return store.WritePendingQueue(ctx, keptItems)
}

func sortKey(id string, nonce []byte) []byte {
	h := sha256.New()
	h.Write([]byte(id))
	h.Write(nonce)
	return h.Sum(nil)
}
