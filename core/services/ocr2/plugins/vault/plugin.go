package vault

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"
	"google.golang.org/protobuf/proto"
)

type ReportingPlugin struct {
	lggr       logger.Logger
	store      *requests.Store[*vaulttypes.Request]
	onchainCfg ocr3types.ReportingPluginConfig
	cfg        *ReportingPluginConfig
	metrics    *pluginMetrics
	validator  *vaultcap.RequestValidator
	lifecycle  *vaultcap.RequestLifecycleTracker

	maxObservationBytes          int
	maxReportsPlusPrecursorBytes int

	// For testing: functions to mock out marshaling/unmarshaling blob handles.
	// The Blob API isn't very test friendly because it uses sum types that belong
	// to an internal package.
	unmarshalBlob func(data []byte) (ocr3_1types.BlobHandle, error)
	marshalBlob   func(handle ocr3_1types.BlobHandle) ([]byte, error)
}

func (r *ReportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return types.Query{}, nil
}

func generateRandomNonce() ([]byte, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return nil, fmt.Errorf("could not generate random nonce: %w", err)
	}

	return nonceBytes, nil
}

func (r *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	start := time.Now()
	defer func() {
		r.lggr.Debugw("observation finished", "seqNr", seqNr, "elapsed", time.Since(start))
	}()

	readKV := NewReadStore(keyValueReader, r.metrics)

	var currentPendingQueueItems []*vaultcommon.StoredPendingQueueItem
	if !gateAllows(ctx, r.lggr, r.cfg.VaultForceEmptyOCRRounds, "VaultForceEmptyOCRRounds") {
		var err error
		currentPendingQueueItems, err = readKV.GetPendingQueue(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not fetch batch of requests: %w", err)
		}
	} else {
		r.lggr.Warnw("VaultForceEmptyOCRRounds is enabled; pending queue is not read this OCR round — store-backed pending observation items are skipped")
	}

	// Avoid log spam by only logging if we have any requests to process.
	if len(currentPendingQueueItems) > 0 {
		r.lggr.Debugw("observation started", "seqNr", seqNr, "batchSize", len(currentPendingQueueItems))
	}

	obspb := &vaultcommon.Observations{}
	optimizations := r.optimizationsEnabled(ctx)

	// First, observe the local queue and broadcast blob payloads so the exact
	// PendingQueueItems + SortNonce wire size is known before packing Observations.
	localQueueItems, ierr := r.store.All()
	if ierr != nil {
		return nil, ierr
	}
	r.metrics.trackLocalQueueSize(ctx, len(localQueueItems))

	// Sort the local queue by ID as we may have to limit its contents
	// later on and we want to maximize the possibility of overlap among
	// honest nodes.
	slices.SortFunc(localQueueItems, func(a, b *vaulttypes.Request) int {
		switch {
		case a.ID() < b.ID():
			return -1
		case a.ID() > b.ID():
			return 1
		default:
			return 0
		}
	})

	pendingQueueHasID := map[string]bool{}
	for _, item := range currentPendingQueueItems {
		pendingQueueHasID[item.Id] = true
	}

	maxBlobBytesSz, err := r.cfg.MaxBlobPayloadBytes.Limit(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch max blob payload size limit: %w", err)
	}
	maxBlobBytes := int(maxBlobBytesSz)

	batchSizeLimit, err := r.cfg.MaxBatchSize.Limit(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch max batch size limit: %w", err)
	}
	maxBlobHandleCount := 2 * batchSizeLimit

	var pack pendingQueueBlobPack
	if optimizations {
		pack, err = r.prepareObservationPendingQueueBlobs(ctx, seqNr, localQueueItems, pendingQueueHasID, maxBlobBytes, maxBlobHandleCount)
		if err != nil {
			return nil, err
		}
		if pack.packedItemCount > 0 {
			r.metrics.trackObservationPendingPack(ctx, pack.packedItemCount, len(pack.blobPayloads))
			r.lggr.Infow("observation packed local items into blob payloads",
				"seqNr", seqNr,
				"packedLocalItemCount", pack.packedItemCount,
				"blobHandleCount", len(pack.blobPayloads),
				"truncated", pack.truncated,
			)
		}
	} else {
		pack, err = r.prepareLegacyObservationPendingQueueBlobs(ctx, seqNr, localQueueItems, pendingQueueHasID, maxBlobHandleCount)
		if err != nil {
			return nil, err
		}
	}

	pendingQueueItems, err := r.broadcastBlobPayloads(ctx, blobBroadcastFetcher, seqNr, pack.blobPayloads, pack.blobPayloadIDs)
	if err != nil {
		return nil, err
	}
	obspb.PendingQueueItems = pendingQueueItems

	// Second, generate a random nonce that we'll use to sort the observations.
	// Each node generates a nonce indepedently, to be concatenated later on.
	nonce, ierr := generateRandomNonce()
	if ierr != nil {
		return nil, fmt.Errorf("could not generate nonce for observation: %w", ierr)
	}
	obspb.SortNonce = nonce

	// Observe store-backed pending queue items after local-queue blob broadcast so blob wire size is known first.
	observedIDs := r.appendPendingQueueObservations(ctx, seqNr, readKV, currentPendingQueueItems, obspb, optimizations)
	if optimizations && len(currentPendingQueueItems) > 0 && len(obspb.Observations) < len(currentPendingQueueItems) {
		r.lggr.Infow("observation: more pending queue items than can be observed",
			"seqNr", seqNr,
			"packedObservationCount", len(obspb.Observations),
			"totalPendingQueueItemCount", len(currentPendingQueueItems),
		)
	}

	obsb, err := proto.MarshalOptions{Deterministic: true}.Marshal(obspb)
	if err != nil {
		return nil, fmt.Errorf("could not marshal observations: %w", err)
	}

	// Avoid log spam by only logging if we have any requests to process.
	if len(currentPendingQueueItems) > 0 {
		r.lggr.Debugw("observation complete", "ids", observedIDs, "batchSize", len(currentPendingQueueItems))
	}
	return types.Observation(obsb), nil
}

func (r *ReportingPlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	obs := &vaultcommon.Observations{}
	if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
		return errors.New("failed to unmarshal observations: " + err.Error())
	}

	if len(ao.Observation) > r.maxObservationBytes {
		return fmt.Errorf("invalid observation: wire size %d exceeds max observation bytes %d", len(ao.Observation), r.maxObservationBytes)
	}

	idToObs := map[string]*vaultcommon.Observation{}
	for _, o := range obs.Observations {
		err := r.validateObservation(ctx, o)
		if err != nil {
			return errors.New("invalid observation: " + err.Error())
		}

		_, seen := idToObs[o.Id]
		if seen {
			return errors.New("invalid observation: a single observation cannot contain duplicate observations for the same request id")
		}

		idToObs[o.Id] = o
	}

	// We expect
	// - that every observation id corresponds to an item in the pending queue.
	//   This is because honest nodes may omit tail items when the full Observations proto would exceed the
	//   max observation byte limit.
	// - that all pending queue items can be fetched as blobs.
	readKV := NewReadStore(keyValueReader, r.metrics)
	var pendingQueueItems []*vaultcommon.StoredPendingQueueItem
	if !gateAllows(ctx, r.lggr, r.cfg.VaultForceEmptyOCRRounds, "VaultForceEmptyOCRRounds") {
		var err error
		pendingQueueItems, err = readKV.GetPendingQueue(ctx)
		if err != nil {
			return fmt.Errorf("could not fetch pending queue from store: %w", err)
		}
	} else {
		r.lggr.Warnw("VaultForceEmptyOCRRounds is enabled; pending queue is not read this OCR round — store-backed pending observation items are skipped")
	}

	pendingIDs := map[string]bool{}
	for _, i := range pendingQueueItems {
		pendingIDs[i.Id] = true
	}
	for id := range idToObs {
		if !pendingIDs[id] {
			return fmt.Errorf("invalid observation: observation id %s is not present in the pending queue", id)
		}
	}

	l, err := r.cfg.MaxBatchSize.Limit(ctx)
	if err != nil {
		return fmt.Errorf("could not fetch max batch size limit: %w", err)
	}

	// The Observation method enforces a max pending queue batch size of 2x the batch size.
	// We can therefore reject any observation with a higher number of observations as invalid.
	maxBatchSize := 2 * l
	if len(obs.PendingQueueItems) > maxBatchSize {
		return fmt.Errorf("invalid observation: too many pending queue items provided, have %d, want max %d", len(obs.PendingQueueItems), maxBatchSize)
	}

	seenItem := map[string]bool{}
	for _, i := range obs.PendingQueueItems {
		bh, err := r.unmarshalBlob(i)
		if err != nil {
			return fmt.Errorf("could not unmarshal blob handle from observation pending queue item: %w", err)
		}

		blob, err := blobFetcher.FetchBlob(ctx, bh)
		if err != nil {
			return fmt.Errorf("could not fetch blob for observation pending queue item: %w", err)
		}

		items, err := unmarshalPendingQueueBlob(blob)
		if err != nil {
			return fmt.Errorf("could not decode pending queue blob: %w", err)
		}
		for _, pit := range items {
			sha, err := shaForProto(pit)
			if err != nil {
				return fmt.Errorf("could not compute sha for pending queue item: %w", err)
			}
			if seenItem[sha] {
				return errors.New("duplicate item found in pending queue item observation")
			}
			seenItem[sha] = true
		}
	}

	return nil
}

func (r *ReportingPlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumNMinusF, r.onchainCfg.N, r.onchainCfg.F, aos), nil
}

func (r *ReportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	writeKV := NewWriteStore(keyValueReadWriter, r.metrics)

	marshalledObs := map[uint8]*vaultcommon.Observations{}
	for _, ao := range aos {
		obs := &vaultcommon.Observations{}
		if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
			// Note: this shouldn't happen as all observations are validated in ValidateObservation.
			r.lggr.Errorw("failed to unmarshal observations", "error", err, "observation", ao.Observation)
			continue
		}

		marshalledObs[uint8(ao.Observer)] = obs
	}

	// ---
	// Phase 1: Process requests from the pending queue by aggregating observations.
	// ---

	// obsMap is a map from observation id -> list of observations across oracles.
	obsMap := map[string][]*vaultcommon.Observation{}
	oidsToReqIDs := map[uint8][]string{} // for debugging only
	for _, ao := range aos {
		observer := uint8(ao.Observer)
		obs := marshalledObs[observer]
		for _, o := range obs.Observations {
			if _, ok := obsMap[o.Id]; !ok {
				obsMap[o.Id] = []*vaultcommon.Observation{}
			}
			obsMap[o.Id] = append(obsMap[o.Id], o)

			if _, ok := oidsToReqIDs[observer]; !ok {
				oidsToReqIDs[observer] = []string{}
			}
			oidsToReqIDs[observer] = append(oidsToReqIDs[observer], o.Id)
		}
	}

	r.lggr.Debugw("stateTransition started", "oracleIDsToRequestIDs", oidsToReqIDs)

	os := &vaultcommon.Outcomes{
		Outcomes: []*vaultcommon.Outcome{},
	}

outcomePackLoop:
	for _, id := range slices.Sorted(maps.Keys(obsMap)) {
		obs := obsMap[id]
		// For each observation we've received for a given Id,
		// we'll sha it and store it in `shaToObs`.
		// This means that each entry in `shaToObs` will contain a list of all
		// of the entries matching a given sha.
		shaToObs := map[string][]*vaultcommon.Observation{}
		for _, ob := range obs {
			sha, err := shaForObservation(ob)
			if err != nil {
				r.lggr.Errorw("failed to compute sha for observation", "error", err, "observation", ob)
				continue
			}
			shaToObs[sha] = append(shaToObs[sha], ob)
		}

		// Now let's identify the "chosen" observation.
		// We do this by checking if which sha has 2F+1 observations.
		// Once we have it, we can break, as mathematically only one
		// sha can reach at least 2F+1 observaions.
		chosen := []*vaultcommon.Observation{}
		for _, sha := range slices.Sorted(maps.Keys(shaToObs)) {
			obs := shaToObs[sha]

			o := obs[0]
			switch {
			case o.RequestType == vaultcommon.RequestType_GET_SECRETS && len(obs) >= 2*r.onchainCfg.F+1:
				// GetRequests required 2F+1 observations because we need exactly T=F+1 shares to reconstruct the secret.
				// Since F shares can be fault, that means T+F=2F+1 shares are required, necessitating 2F+1 observations.
				if r.optimizationsEnabled(ctx) {
					chosen = shaToObs[sha][:2*r.onchainCfg.F+1]
				} else {
					chosen = shaToObs[sha]
				}
				r.lggr.Debugw("sufficient observations for sha", "sha", sha, "requestType", "GetSecrets", "count", len(obs), "threshold", 2*r.onchainCfg.F+1, "id", id)
			case o.RequestType != vaultcommon.RequestType_GET_SECRETS && len(obs) >= r.onchainCfg.F+1:
				// F+1 means that at least 1 honest node has provided this observation, so that's enough for all other request
				// types.
				// Technically we could have two shas with F+1 observations. If that happens we'll pick the last one.
				// This is deterministic since we're sorting by shas above.
				chosen = shaToObs[sha]
				r.lggr.Debugw("sufficient observations for sha", "sha", sha, "count", len(obs), "threshold", r.onchainCfg.F+1, "id", id)
			}
		}

		if len(chosen) == 0 {
			shaToObsCount := map[string]int{}
			for sha, obs := range shaToObs {
				shaToObsCount[sha] = len(obs)
			}
			r.lggr.Warnw("insufficient observations found for id", "id", id, "shaToObsCount", shaToObsCount)
			continue
		}

		// The shas are the same so the requests will have
		// the same Id and Type.
		first := chosen[0]
		o := &vaultcommon.Outcome{
			Id:          first.Id,
			RequestType: first.RequestType,
		}
		switch first.RequestType {
		case vaultcommon.RequestType_GET_SECRETS:
			r.stateTransitionGetSecrets(ctx, chosen, o)
		case vaultcommon.RequestType_CREATE_SECRETS:
			r.stateTransitionCreateSecrets(ctx, writeKV, chosen, o)
		case vaultcommon.RequestType_UPDATE_SECRETS:
			r.stateTransitionUpdateSecrets(ctx, writeKV, chosen, o)
		case vaultcommon.RequestType_DELETE_SECRETS:
			r.stateTransitionDeleteSecrets(ctx, writeKV, chosen, o)
		case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
			r.stateTransitionListSecretIdentifiers(ctx, writeKV, chosen, o)
		default:
			r.lggr.Debugw("unknown request type, skipping...", "requestType", first.RequestType, "id", id)
			continue
		}

		os.Outcomes = append(os.Outcomes, o)
		if r.optimizationsEnabled(ctx) && proto.Size(os) > r.maxReportsPlusPrecursorBytes {
			os.Outcomes = os.Outcomes[:len(os.Outcomes)-1]
			r.lggr.Warnw("state transition: more observations than can be included in response",
				"id", id,
				"maxReportsPlusPrecursorBytes", r.maxReportsPlusPrecursorBytes,
				"packedOutcomeCount", len(os.Outcomes),
				"scheduledRequestIDs", len(obsMap),
			)
			break outcomePackLoop
		}
	}

	// ---
	// Phase 2: Process the pending queue.
	// ---
	err := r.stateTransitionPendingQueue(ctx, seqNr, writeKV, marshalledObs, blobFetcher)
	if err != nil {
		return ocr3_1types.ReportsPlusPrecursor{}, fmt.Errorf("could not process pending queue during state transition: %w", err)
	}

	nowST := time.Now()
	for _, out := range os.Outcomes {
		r.lifecycle.RecordStateTransitionOutcome(out.Id, seqNr, nowST)
	}

	ospb, err := proto.MarshalOptions{Deterministic: true}.Marshal(os)
	if err != nil {
		return ocr3_1types.ReportsPlusPrecursor{}, fmt.Errorf("could not marshal outcomes: %w", err)
	}

	if len(os.Outcomes) > 0 {
		r.lggr.Debugw("State transition complete", "count", len(os.Outcomes), "err", err)
	}
	return ocr3_1types.ReportsPlusPrecursor(ospb), nil
}

func (r *ReportingPlugin) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader) error {
	// Not currently used by the protocol, so we don't implement it.
	return errors.New("not implemented")
}

func (r *ReportingPlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	outcomes := &vaultcommon.Outcomes{}
	err := proto.Unmarshal([]byte(reportsPlusPrecursor), outcomes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal outcomes: %w", err)
	}

	reports := []ocr3types.ReportPlus[[]byte]{}
	for _, o := range outcomes.Outcomes {
		switch o.RequestType {
		case vaultcommon.RequestType_GET_SECRETS:
			rep, err := r.generateProtoReport(o.Id, o.RequestType, o.GetGetSecretsResponse())
			if err != nil {
				r.lggr.Errorw("failed to generate Proto report", "error", err, "id", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_CREATE_SECRETS:
			rep, err := r.generateJSONReport(o.Id, o.RequestType, o.GetCreateSecretsResponse())
			if err != nil {
				r.lggr.Errorw("failed to generate JSON report", "error", err, "id", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_UPDATE_SECRETS:
			rep, err := r.generateJSONReport(o.Id, o.RequestType, o.GetUpdateSecretsResponse())
			if err != nil {
				r.lggr.Errorw("failed to generate JSON report", "error", err, "id", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_DELETE_SECRETS:
			rep, err := r.generateJSONReport(o.Id, o.RequestType, o.GetDeleteSecretsResponse())
			if err != nil {
				r.lggr.Errorw("failed to generate JSON report", "error", err, "id", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
			rep, err := r.generateJSONReport(o.Id, o.RequestType, o.GetListSecretIdentifiersResponse())
			if err != nil {
				r.lggr.Errorw("failed to generate JSON report", "error", err, "id", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		default:
		}
	}

	if len(reports) > 0 {
		r.lggr.Debugw("Reports complete", "count", len(reports))
	}
	return reports, nil
}

func (r *ReportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (r *ReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (r *ReportingPlugin) Close() error {
	return errors.Join(
		r.cfg.MaxSecretsPerOwner.Close(),
		r.cfg.MaxCiphertextLengthBytes.Close(),
		r.cfg.MaxIdentifierKeyLengthBytes.Close(),
		r.cfg.MaxIdentifierOwnerLengthBytes.Close(),
		r.cfg.MaxIdentifierNamespaceLengthBytes.Close(),
		r.cfg.MaxShareLengthBytes.Close(),
		r.cfg.MaxRequestBatchSize.Close(),
		r.cfg.MaxBatchSize.Close(),
		r.cfg.VaultForceEmptyOCRRounds.Close(),
	)
}
