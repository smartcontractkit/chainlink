package vault

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"sync"
	"time"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"
	"github.com/smartcontractkit/smdkg/dkgocr"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/dkgrecipientkey"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	blobBroadcastTimeout        = 2 * time.Second
	maxConcurrentBlobBroadcasts = 10
	sortNonceLength             = 32
)

type ReportingPluginConfig struct {
	LazyPublicKey *vaultcap.LazyPublicKey
	// Sourced from the DKG DB results package
	PublicKey       *tdh2easy.PublicKey
	PrivateKeyShare *tdh2easy.PrivateShare

	// Sourced from the offchain config
	MaxSecretsPerOwner              limits.BoundLimiter[int]
	MaxShareLengthBytes             limits.BoundLimiter[pkgconfig.Size]
	MaxBatchSize                    limits.BoundLimiter[int]
	MaxPendingQueueWriteSize        limits.BoundLimiter[int]
	MaxBlobPayloadBytes             limits.BoundLimiter[pkgconfig.Size]
	VaultForceEmptyOCRRounds        limits.GateLimiter
	VaultPendingQueueStallThreshold limits.BoundLimiter[int]
}

func NewReportingPluginFactory(
	lggr logger.Logger,
	store *requests.Store[*vaulttypes.Request],
	db dkgocrtypes.ResultPackageDatabase,
	recipientKey *dkgrecipientkey.Key,
	lazyPublicKey *vaultcap.LazyPublicKey,
	limitsFactory limits.Factory,
	lifecycle *vaultcap.RequestLifecycleTracker,
) (*ReportingPluginFactory, error) {
	if db == nil {
		return nil, errors.New("result package db cannot be nil")
	}

	if recipientKey == nil {
		return nil, errors.New("DKG recipient key cannot be nil when using result package db")
	}

	if lifecycle == nil {
		return nil, errors.New("request lifecycle tracker cannot be nil")
	}

	cfg := &ReportingPluginConfig{
		LazyPublicKey: lazyPublicKey,
	}

	return &ReportingPluginFactory{
		lggr:          lggr.Named("VaultReportingPluginFactory"),
		store:         store,
		cfg:           cfg,
		db:            db,
		recipientKey:  recipientKey,
		limitsFactory: limitsFactory,
		lifecycle:     lifecycle,
	}, nil
}

type ReportingPluginFactory struct {
	lggr          logger.Logger
	store         *requests.Store[*vaulttypes.Request]
	cfg           *ReportingPluginConfig
	db            dkgocrtypes.ResultPackageDatabase
	recipientKey  *dkgrecipientkey.Key
	limitsFactory limits.Factory
	lifecycle     *vaultcap.RequestLifecycleTracker
}

func (r *ReportingPluginFactory) getKeyMaterial(ctx context.Context, instanceID string) (publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, err error) {
	pack, err := r.db.ReadResultPackage(ctx, dkgocrtypes.InstanceID(instanceID))
	if err != nil {
		return nil, nil, fmt.Errorf("could not read result package from db: %w", err)
	}
	if pack == nil {
		return nil, nil, fmt.Errorf("no result package found in db for instance ID %s", instanceID)
	}
	rP := dkgocr.NewResultPackage()
	err = rP.UnmarshalBinary(pack.ReportWithResultPackage)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal result package: %w", err)
	}

	tdh2PubKey, err := tdh2shim.TDH2PublicKeyFromDKGResult(rP)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get tdh2 public key from DKG result: %w", err)
	}
	publicKey, err = tdh2ToTDH2EasyPK(tdh2PubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not convert to tdh2easy public key: %w", err)
	}

	tdh2PrivateKeyShare, err := tdh2shim.TDH2PrivateShareFromDKGResult(rP, r.recipientKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get tdh2 private key share from DKG result: %w", err)
	}
	privateKeyShare, err = tdh2ToTDH2EasyKS(tdh2PrivateKeyShare)
	if err != nil {
		return nil, nil, fmt.Errorf("could not convert to tdh2easy private key share: %w", err)
	}

	return publicKey, privateKeyShare, nil
}

const dkgPollInterval = 2 * time.Second

// pollForKeyMaterial polls the DKG result package database until the key
// material for the given instance ID is available or the context is cancelled.
// This avoids returning an immediate error when the DKG protocol hasn't
// completed yet, which would trigger libocr's exponential backoff (up to 2
// minutes between retries). By polling here within the MaxDurationInitialization
// window, the vault oracle can start as soon as the DKG result is written.
func (r *ReportingPluginFactory) pollForKeyMaterial(ctx context.Context, instanceID string) (publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, err error) {
	for {
		publicKey, privateKeyShare, err = r.getKeyMaterial(ctx, instanceID)
		if err == nil {
			return publicKey, privateKeyShare, nil
		}

		r.lggr.Debugw("DKG result package not yet available, will retry", "instanceID", instanceID, "error", err)

		select {
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("context cancelled while waiting for DKG key material (instanceID=%s): %w", instanceID, err)
		case <-time.After(dkgPollInterval):
		}
	}
}

func newReportingPluginConfigLimiters(factory limits.Factory) (*ReportingPluginConfig, error) {
	maxShareLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultShareSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultShareSizeLimit: %w", err)
	}

	vaultForceEmptyOCRRounds, err := limits.MakeGateLimiter(factory, cresettings.Default.VaultForceEmptyOCRRounds)
	if err != nil {
		return nil, fmt.Errorf("VaultForceEmptyOCRRounds: %w", err)
	}

	vaultPendingQueueStallThreshold, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultPendingQueueStallThreshold)
	if err != nil {
		return nil, fmt.Errorf("VaultPendingQueueStallThreshold: %w", err)
	}

	maxBlobPayloadBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultMaxBlobPayloadSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultMaxBlobPayloadSizeLimit: %w", err)
	}

	maxPendingQueueWriteSizeLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultPendingQueueWriteSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultPendingQueueWriteSizeLimit: %w", err)
	}

	return &ReportingPluginConfig{
		MaxShareLengthBytes:             maxShareLengthBytesLimiter,
		MaxBlobPayloadBytes:             maxBlobPayloadBytesLimiter,
		MaxPendingQueueWriteSize:        maxPendingQueueWriteSizeLimiter,
		VaultForceEmptyOCRRounds:        vaultForceEmptyOCRRounds,
		VaultPendingQueueStallThreshold: vaultPendingQueueStallThreshold,
	}, nil
}

func logLimit[N limits.Number](ctx context.Context, lggr logger.Logger, limiter limits.BoundLimiter[N]) N {
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: "DUMMY-OWNER-FOR-LOGGING"})
	limit, err := limiter.Limit(ctx)
	if err != nil {
		lggr.Errorw("could not fetch limit", "error", err)
	}
	return limit
}

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	var configProto vaultcommon.ReportingPluginConfig
	if err := proto.Unmarshal(config.OffchainConfig, &configProto); err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not unmarshal reporting plugin config: %w", err)
	}

	cfg, err := newReportingPluginConfigLimiters(r.limitsFactory)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create reporting plugin config limiters: %w", err)
	}

	maxSecretsPerOwnerLimit := cresettings.Default.PerOwner.VaultSecretsLimit
	if configProto.MaxSecretsPerOwner != 0 {
		maxSecretsPerOwnerLimit.DefaultValue = int(configProto.MaxSecretsPerOwner)
	}

	cfg.MaxSecretsPerOwner, err = limits.MakeUpperBoundLimiter(r.limitsFactory, maxSecretsPerOwnerLimit)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create max secrets per owner limiter: %w", err)
	}

	cfg.MaxBatchSize, err = limits.MakeUpperBoundLimiter(r.limitsFactory, cresettings.Default.VaultPluginBatchSizeLimit)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create max batch size limiter: %w", err)
	}

	if configProto.DKGInstanceID == nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, errors.New("DKG instance ID cannot be nil")
	}

	r.lggr.Debugw("fetching key material for instance id", "instanceID", *configProto.DKGInstanceID)
	publicKey, privateKeyShare, err := r.pollForKeyMaterial(ctx, *configProto.DKGInstanceID)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not get key material from DB: %w", err)
	}

	r.cfg.LazyPublicKey.Set(publicKey)

	cfg.PublicKey = publicKey
	cfg.PrivateKeyShare = privateKeyShare

	metrics, err := newPluginMetrics(config.ConfigDigest.String())
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create plugin metrics: %w", err)
	}

	pluginLimits, err := initializePluginLimits(ctx, r.limitsFactory)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not resolve plugin limits: %w", err)
	}

	validator, err := vaultcap.NewRequestValidatorFromLimitsFactory(r.limitsFactory)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create request validator: %w", err)
	}

	r.lggr.Debugw("instantiating VaultReportingPlugin with config",
		"maxSecretsPerOwner", logLimit(ctx, r.lggr, cfg.MaxSecretsPerOwner),
		"maxCiphertextLengthBytes", logLimit(ctx, r.lggr, validator.MaxCiphertextLengthLimiter),
		"maxIdentifierKeyLengthBytes", logLimit(ctx, r.lggr, validator.MaxIdentifierKeyLengthLimiter),
		"maxIdentifierOwnerLengthBytes", logLimit(ctx, r.lggr, validator.MaxIdentifierOwnerLengthLimiter),
		"maxIdentifierNamespaceLengthBytes", logLimit(ctx, r.lggr, validator.MaxIdentifierNamespaceLengthLimiter),
		"maxRequestBatchSize", logLimit(ctx, r.lggr, validator.MaxRequestBatchSizeLimiter),
		"maxShareLengthBytes", logLimit(ctx, r.lggr, cfg.MaxShareLengthBytes),
		"batchSize", logLimit(ctx, r.lggr, cfg.MaxBatchSize),
		"maxPendingQueueWriteSize", logLimit(ctx, r.lggr, cfg.MaxPendingQueueWriteSize),
		"maxBlobPayloadBytes", logLimit(ctx, r.lggr, cfg.MaxBlobPayloadBytes),
		"maxQueryBytes", pluginLimits.MaxQueryBytes,
		"maxObservationBytes", pluginLimits.MaxObservationBytes,
		"maxReportsPlusPrecursorBytes", pluginLimits.MaxReportsPlusPrecursorBytes,
		"maxReportBytes", pluginLimits.MaxReportBytes,
		"maxReportCount", pluginLimits.MaxReportCount,
		"maxKeyValueModifiedKeysPlusValuesBytes", pluginLimits.MaxKeyValueModifiedKeysPlusValuesBytes,
		"maxKeyValueModifiedKeys", pluginLimits.MaxKeyValueModifiedKeys,
		"reportingPluginLimitsMaxBlobPayloadBytes", pluginLimits.MaxBlobPayloadBytes,
		"maxPerOracleUnexpiredBlobCumulativePayloadBytes", pluginLimits.MaxPerOracleUnexpiredBlobCumulativePayloadBytes,
		"maxPerOracleUnexpiredBlobCount", pluginLimits.MaxPerOracleUnexpiredBlobCount,
	)

	r.lifecycle.SetConfigDigest(config.ConfigDigest.String())

	plugin := &ReportingPlugin{
		lggr:                         r.lggr.Named("VaultReportingPlugin"),
		store:                        r.store,
		cfg:                          cfg,
		metrics:                      metrics,
		onchainCfg:                   config,
		validator:                    validator,
		lifecycle:                    r.lifecycle,
		maxObservationBytes:          pluginLimits.MaxObservationBytes,
		maxReportsPlusPrecursorBytes: pluginLimits.MaxReportsPlusPrecursorBytes,
		unmarshalBlob: func(data []byte) (ocr3_1types.BlobHandle, error) {
			handle := ocr3_1types.BlobHandle{}
			err := handle.UnmarshalBinary(data)
			return handle, err
		},
		marshalBlob: func(handle ocr3_1types.BlobHandle) ([]byte, error) {
			return handle.MarshalBinary()
		},
	}
	return plugin, ocr3_1types.ReportingPluginInfo1{
		Name:   "VaultReportingPlugin",
		Limits: pluginLimits,
	}, nil
}

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

	pendingQueueStallTracker pendingQueueStallTracker
}

type pendingQueueStallTracker struct {
	mu    sync.Mutex
	seqNr uint64
	count int
}

func (t *pendingQueueStallTracker) record(seqNr uint64) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.seqNr != seqNr {
		t.seqNr = seqNr
		t.count = 0
		return t.count
	}
	t.count++
	return t.count
}

func (t *pendingQueueStallTracker) getCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.count
}

func countPendingQueueStallSignals(aos []types.AttributedObservation) int {
	count := 0
	for _, ao := range aos {
		obs := &vaultcommon.Observations{}
		if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
			continue
		}
		if observationSignalsPendingQueueStall(obs) {
			count++
		}
	}
	return count
}

func countPendingQueueStallSignalsInMap(obsByObserver map[uint8]*vaultcommon.Observations) int {
	count := 0
	for _, obs := range obsByObserver {
		if observationSignalsPendingQueueStall(obs) {
			count++
		}
	}
	return count
}

func (r *ReportingPlugin) purgeStalledPendingQueue(
	ctx context.Context,
	l logger.Logger,
	store pendingQueueStore,
	stallSignalCount int,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	if err := store.WritePendingQueue(ctx, nil); err != nil {
		return ocr3_1types.ReportsPlusPrecursor{}, fmt.Errorf("could not purge stalled pending queue: %w", err)
	}
	r.metrics.trackPendingQueueWrittenSize(ctx, 0)
	r.metrics.trackPendingQueuePurge(ctx)
	l.Warnw("purged stalled pending queue after f+1 stall signals",
		"stallSignalCount", stallSignalCount,
		"threshold", r.onchainCfg.F+1,
	)
	ospb, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vaultcommon.Outcomes{})
	if err != nil {
		return ocr3_1types.ReportsPlusPrecursor{}, fmt.Errorf("could not marshal empty outcomes after pending queue purge: %w", err)
	}
	return ocr3_1types.ReportsPlusPrecursor(ospb), nil
}

func (r *ReportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return types.Query{}, nil
}

func generateRandomNonce() ([]byte, error) {
	nonceBytes := make([]byte, sortNonceLength)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return nil, fmt.Errorf("could not generate random nonce: %w", err)
	}

	return nonceBytes, nil
}

// marshalPendingQueueBlobPayload encodes pending queue items for OCR3.1 BroadcastBlob.
// Always marshals as PendingQueueBlobItems. Single items are wire-compatible with StoredPendingQueueItem,
// so the non-optimizations unmarshal path can decode them without knowing the new type.
// Batch items wrap each StoredPendingQueueItem (payload + ID) inside an Any.
func marshalPendingQueueBlobPayload(items []*vaultcommon.StoredPendingQueueItem) ([]byte, error) {
	if len(items) == 0 {
		return nil, errors.New("empty pending queue blob payload")
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

func (r *ReportingPlugin) shouldPurgePendingQueue(ctx context.Context) bool {
	if gateAllows(ctx, r.lggr, r.cfg.VaultForceEmptyOCRRounds, "VaultForceEmptyOCRRounds") {
		return true
	}
	stallThreshold, err := r.cfg.VaultPendingQueueStallThreshold.Limit(ctx)
	if err != nil {
		r.lggr.Errorw("could not fetch pending queue stall threshold", "error", err)
		return false
	}
	if stallThreshold == 0 {
		return false
	}
	// The pending queue must never be purged with fewer than 2f+1 stall
	// signals, regardless of the configured stall threshold.
	if stallThreshold < 2*r.onchainCfg.F+1 {
		stallThreshold = 2*r.onchainCfg.F + 1
	}

	stalledObservationCount := r.pendingQueueStallTracker.getCount()
	return stalledObservationCount >= stallThreshold
}

type pendingQueueStore interface {
	WritePendingQueue(ctx context.Context, pending []*vaultcommon.StoredPendingQueueItem) error
}

func (r *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	l := r.roundLggr(seqNr)
	start := time.Now()
	defer func() {
		l.Debugw("observation finished", "elapsed", time.Since(start))
	}()

	obspb := &vaultcommon.Observations{}

	r.pendingQueueStallTracker.record(seqNr)

	// First, generate a random nonce that we'll use to sort the observations.
	// Each node generates a nonce independently, to be concatenated later on.
	nonce, ierr := generateRandomNonce()
	if ierr != nil {
		return nil, fmt.Errorf("could not generate nonce for observation: %w", ierr)
	}
	obspb.SortNonce = nonce

	if shouldPurge := r.shouldPurgePendingQueue(ctx); shouldPurge {
		obspb.PendingQueueStallSignal = vaultcommon.PendingQueueStallSignal_PENDING_QUEUE_STALL_SIGNAL_STALLED
		obsb, err := proto.MarshalOptions{Deterministic: true}.Marshal(obspb)
		if err != nil {
			return nil, fmt.Errorf("could not marshal observations: %w", err)
		}
		r.metrics.trackPendingQueueStallSignal(ctx)
		l.Warnw("pending queue stall threshold reached; signaling queue purge",
			"observationCount", r.pendingQueueStallTracker.getCount(),
		)
		return types.Observation(obsb), nil
	}

	readKV := NewReadStore(keyValueReader, r.metrics)

	var currentPendingQueueItems []*vaultcommon.StoredPendingQueueItem
	var err error
	currentPendingQueueItems, err = readKV.GetPendingQueue(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch of requests: %w", err)
	}

	// Avoid log spam by only logging if we have any requests to process.
	if len(currentPendingQueueItems) > 0 {
		l.Debugw("observation started", "batchSize", len(currentPendingQueueItems))
	}

	// Second, observe the local queue and broadcast blob payloads so the exact
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

	pendingQueueItems, err := r.broadcastBlobPayloads(ctx, blobBroadcastFetcher, seqNr, pack.blobPayloads, pack.blobPayloadIDs)
	if err != nil {
		return nil, err
	}
	obspb.PendingQueueItems = pendingQueueItems

	// Observe store-backed pending queue items after local-queue blob broadcast so blob wire size is known first.
	observedIDs := r.appendPendingQueueObservations(ctx, seqNr, readKV, currentPendingQueueItems, obspb)
	if len(currentPendingQueueItems) > 0 && len(obspb.Observations) < len(currentPendingQueueItems) {
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
		l.Debugw("observation complete", "ids", observedIDs, "batchSize", len(currentPendingQueueItems))
	}
	return types.Observation(obsb), nil
}

func (r *ReportingPlugin) observePendingQueueItem(
	ctx context.Context,
	seqNr uint64,
	readKV *KVStore,
	req *vaultcommon.StoredPendingQueueItem,
) (*vaultcommon.Observation, error) {
	o := &vaultcommon.Observation{
		Id: req.Id,
	}

	payload, err := req.Item.UnmarshalNew()
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request payload: %w", err)
	}

	switch tp := payload.(type) {
	case *vaultcommon.GetSecretsRequest:
		r.observeGetSecrets(ctx, seqNr, req.Id, readKV, tp, o)
	case *vaultcommon.CreateSecretsRequest:
		r.observeCreateSecrets(ctx, seqNr, req.Id, tp, o)
	case *vaultcommon.UpdateSecretsRequest:
		r.observeUpdateSecrets(ctx, seqNr, req.Id, tp, o)
	case *vaultcommon.DeleteSecretsRequest:
		r.observeDeleteSecrets(ctx, seqNr, req.Id, readKV, tp, o)
	case *vaultcommon.ListSecretIdentifiersRequest:
		r.observeListSecretIdentifiers(ctx, seqNr, req.Id, readKV, tp, o)
	default:
		return nil, fmt.Errorf("unknown request type %T", payload)
	}

	return o, nil
}

func (r *ReportingPlugin) validatePendingQueueObservationsPrefix(
	pendingQueueItems []*vaultcommon.StoredPendingQueueItem,
	obs *vaultcommon.Observations,
) error {
	expected := pendingQueueItems

	if len(obs.Observations) > len(expected) {
		return fmt.Errorf("invalid observation: got %d store-backed observations, want at most %d", len(obs.Observations), len(expected))
	}

	for i, o := range obs.Observations {
		if o.Id != expected[i].Id {
			return fmt.Errorf("invalid observation: observation at position %d has id %s, want %s", i, o.Id, expected[i].Id)
		}
	}

	return nil
}

// appendPendingQueueObservations appends one Observation per store-backed pending queue item.
// Stops before obspb exceeds maxObservationBytes.
func (r *ReportingPlugin) appendPendingQueueObservations(
	ctx context.Context,
	seqNr uint64,
	readKV *KVStore,
	currentPendingQueueItems []*vaultcommon.StoredPendingQueueItem,
	obspb *vaultcommon.Observations,
) []string {
	ids := make([]string, 0, len(currentPendingQueueItems))
	for _, req := range currentPendingQueueItems {
		o, err := r.observePendingQueueItem(ctx, seqNr, readKV, req)
		if err != nil {
			o = observationToErrContribution(&vaultcommon.Observation{
				Id: req.Id,
			}, userFacingError(err, "failed to observe pending queue item"))
			if payload, uerr := req.Item.UnmarshalNew(); uerr == nil {
				o.RequestType = requestTypeForPayload(payload)
			}
			r.requestLggr(seqNr, req.Id).Warnw("pending queue item observation failed; emitting error contribution", "error", err)
		} else {
			if cerr := r.validateContribution(ctx, req, o); cerr != nil {
				o = observationToErrContribution(o, userFacingError(cerr, "request is not valid"))
				r.requestLggr(seqNr, req.Id).Warnw("pending queue item failed contribution self-check; emitting error contribution", "error", cerr)
			}
		}

		obspb.Observations = append(obspb.Observations, o)
		if proto.Size(obspb) > r.maxObservationBytes {
			obspb.Observations = obspb.Observations[:len(obspb.Observations)-1]
			r.requestLggr(seqNr, req.Id).Warnw("observation proto would exceed max observation bytes; stopping pending-queue observation pack",
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

func (r *ReportingPlugin) observeGetSecrets(ctx context.Context, seqNr uint64, requestID string, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	l := r.typedRequestLggr(seqNr, requestID, "GetSecrets")
	tp := req.(*vaultcommon.GetSecretsRequest)
	o.RequestType = vaultcommon.RequestType_GET_SECRETS

	requestsCountForID := map[string]int{}
	for _, sr := range tp.Requests {
		var key string
		if sr.Id == nil {
			key = "<nil>"
		} else {
			key = vaulttypes.KeyFor(sr.Id)
		}
		requestsCountForID[key]++
	}

	resps := []*vaultcommon.SecretResponse{}
	for _, secretRequest := range tp.Requests {
		resp, ierr := r.observeGetSecretsRequest(ctx, reader, secretRequest, requestsCountForID)
		if ierr != nil {
			logUserErrorAware(l, "failed to observe get secret request item", ierr, "id", secretRequest.Id)
			errorMsg := userFacingError(ierr, vaulttypes.SecretGetSystemErrorFallback)
			resps = append(resps, &vaultcommon.SecretResponse{
				Id: secretRequest.Id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: errorMsg,
				},
			})
		} else {
			l.Debugw("observed get secret request item", "id", resp.Id)
			resps = append(resps, resp)
		}
	}

	o.Response = &vaultcommon.Observation_GetSecretsResponse{
		GetSecretsResponse: &vaultcommon.GetSecretsResponse{
			Responses: resps,
		},
	}
}

type share struct {
	data []byte
}

func (s *share) encryptWithKeyBinary(pk string) ([]byte, error) {
	publicKey, err := hex.DecodeString(pk)
	if err != nil {
		return nil, vaulttypes.NewUserError("failed to convert public key to bytes: " + err.Error())
	}

	if len(publicKey) != curve25519.PointSize {
		return nil, vaulttypes.NewUserError(fmt.Sprintf("invalid public key size: expected %d bytes, got %d bytes", curve25519.PointSize, len(publicKey)))
	}

	publicKeyLength := [curve25519.PointSize]byte(publicKey)
	encrypted, err := box.SealAnonymous(nil, s.data, &publicKeyLength, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt decryption share: %w", err)
	}

	return encrypted, nil
}

func generatePlaintextShare(publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, encryptedSecret []byte, workflowOwner string) (*share, error) {
	ct := &tdh2easy.Ciphertext{}
	err := ct.UnmarshalVerify(encryptedSecret, publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ciphertext: %w", err)
	}

	es := hex.EncodeToString(encryptedSecret)
	err = vaultcap.EnsureRightLabelOnSecret(publicKey, es, workflowOwner)
	if err != nil {
		return nil, errors.New("failed to verify label on secret. error: " + err.Error())
	}

	s, err := tdh2easy.Decrypt(ct, privateKeyShare)
	if err != nil {
		return nil, fmt.Errorf("could not generate decryption share: %w", err)
	}

	sb, err := s.Marshal()
	if err != nil {
		return nil, errors.New("could not marshal decryption share")
	}

	return &share{data: sb}, nil
}

func (r *ReportingPlugin) observeGetSecretsRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.SecretRequest, requestsCountForID map[string]int) (*vaultcommon.SecretResponse, error) {
	id, err := r.validateGetSecretsRequestItem(ctx, secretRequest, requestsCountForID)
	if err != nil {
		return nil, err
	}

	secret, err := reader.GetSecret(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}
	if secret == nil {
		return nil, vaulttypes.NewUserError("key does not exist")
	}

	sh, err := generatePlaintextShare(r.cfg.PublicKey, r.cfg.PrivateKeyShare, secret.EncryptedSecret, id.Owner)
	if err != nil {
		return nil, err
	}

	shares := []*vaultcommon.EncryptedShares{}
	for _, pk := range secretRequest.EncryptionKeys {
		encShare, err := sh.encryptWithKeyBinary(pk)
		if err != nil {
			return nil, err
		}

		shares = append(shares, &vaultcommon.EncryptedShares{
			EncryptionKey: pk,
			BinaryShares:  [][]byte{encShare},
		})
	}

	return &vaultcommon.SecretResponse{
		Id: id,
		Result: &vaultcommon.SecretResponse_Data{
			Data: &vaultcommon.SecretData{
				EncryptedValue:               hex.EncodeToString(secret.EncryptedSecret),
				EncryptedDecryptionKeyShares: shares,
			},
		},
	}, nil
}

func (r *ReportingPlugin) observeCreateSecrets(ctx context.Context, seqNr uint64, requestID string, req proto.Message, o *vaultcommon.Observation) {
	l := r.typedRequestLggr(seqNr, requestID, "CreateSecrets")
	tp := req.(*vaultcommon.CreateSecretsRequest)
	o.RequestType = vaultcommon.RequestType_CREATE_SECRETS
	o.Request = &vaultcommon.Observation_CreateSecretsRequest{
		CreateSecretsRequest: tp,
	}

	requestsCountForID := buildEncryptedSecretIdentifierCounts(tp.EncryptedSecrets)

	resps := []*vaultcommon.CreateSecretResponse{}
	for _, sr := range tp.EncryptedSecrets {
		validatedID, ierr := r.validateEncryptedSecretPayload(ctx, sr, requestsCountForID)
		if ierr != nil {
			logUserErrorAware(l, "failed to handle create secret request item", ierr, "id", sr.Id)
			errorMsg := userFacingError(ierr, "failed to handle create secret request")
			resps = append(resps, &vaultcommon.CreateSecretResponse{
				Id:      sr.Id,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			l.Debugw("observed create secret request item", "id", validatedID)
			resps = append(resps, &vaultcommon.CreateSecretResponse{
				Id: validatedID,
				// false because it hasn't been processed yet.
				// When the write is handled successfully in StateTransition
				// we'll update this to true.
				Success: false,
			})
		}
	}

	o.Response = &vaultcommon.Observation_CreateSecretsResponse{
		CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
			Responses: resps,
		},
	}
}

func (r *ReportingPlugin) observeUpdateSecrets(ctx context.Context, seqNr uint64, requestID string, req proto.Message, o *vaultcommon.Observation) {
	l := r.typedRequestLggr(seqNr, requestID, "UpdateSecrets")
	tp := req.(*vaultcommon.UpdateSecretsRequest)
	o.RequestType = vaultcommon.RequestType_UPDATE_SECRETS
	o.Request = &vaultcommon.Observation_UpdateSecretsRequest{
		UpdateSecretsRequest: tp,
	}

	requestsCountForID := buildEncryptedSecretIdentifierCounts(tp.EncryptedSecrets)

	resps := []*vaultcommon.UpdateSecretResponse{}
	for _, sr := range tp.EncryptedSecrets {
		validatedID, ierr := r.validateEncryptedSecretPayload(ctx, sr, requestsCountForID)
		if ierr != nil {
			logUserErrorAware(l, "failed to observe update secret request item", ierr, "id", sr.Id)
			errorMsg := userFacingError(ierr, "failed to handle update secret request")
			resps = append(resps, &vaultcommon.UpdateSecretResponse{
				Id:      sr.Id,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			l.Debugw("observed update secret request item", "id", validatedID)
			resps = append(resps, &vaultcommon.UpdateSecretResponse{
				Id: validatedID,
				// false because it hasn't been processed yet.
				// When the write is handled successfully in StateTransition
				// we'll update this to true.
				Success: false,
			})
		}
	}

	o.Response = &vaultcommon.Observation_UpdateSecretsResponse{
		UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
			Responses: resps,
		},
	}
}

func (r *ReportingPlugin) observeListSecretIdentifiers(ctx context.Context, seqNr uint64, requestID string, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.ListSecretIdentifiersRequest)
	l := r.typedRequestLggr(seqNr, requestID, "ListSecretIdentifiers").With("owner", tp.Owner)
	o.RequestType = vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS
	o.Request = &vaultcommon.Observation_ListSecretIdentifiersRequest{
		ListSecretIdentifiersRequest: tp,
	}

	resp, err := r.processListSecretIdentifiersRequest(ctx, seqNr, requestID, reader, tp)
	if err != nil {
		l.Debugw("failed to process list secret identifiers request", "error", err)
		o.Response = &vaultcommon.Observation_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: &vaultcommon.ListSecretIdentifiersResponse{
				Error:   err.Error(),
				Success: false,
			},
		}
		return
	}

	l.Debugw("observed list secret identifiers request")
	o.Response = &vaultcommon.Observation_ListSecretIdentifiersResponse{
		ListSecretIdentifiersResponse: resp,
	}
}

func (r *ReportingPlugin) processListSecretIdentifiersRequest(ctx context.Context, seqNr uint64, requestID string, reader ReadKVStore, req *vaultcommon.ListSecretIdentifiersRequest) (*vaultcommon.ListSecretIdentifiersResponse, error) {
	if err := r.validateListSecretIdentifiersOwnerNonempty(req); err != nil {
		return nil, err
	}

	md, err := reader.GetMetadata(ctx, req.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for owner: %w", err)
	}

	if md == nil {
		// No metadata, so the list is empty.
		// The user hasn't added any items to the vault DON yet.
		r.typedRequestLggr(seqNr, requestID, "ListSecretIdentifiers").With("owner", req.Owner).Debugw("successfully read metadata for owner: no metadata found, returning empty list")
		return &vaultcommon.ListSecretIdentifiersResponse{Identifiers: []*vaultcommon.SecretIdentifier{}, Success: true}, nil
	}

	sort.Slice(md.SecretIdentifiers, func(i, j int) bool {
		if md.SecretIdentifiers[i].Namespace == md.SecretIdentifiers[j].Namespace {
			return md.SecretIdentifiers[i].Key < md.SecretIdentifiers[j].Key
		}
		return md.SecretIdentifiers[i].Namespace < md.SecretIdentifiers[j].Namespace
	})

	if req.Namespace == "" {
		return &vaultcommon.ListSecretIdentifiersResponse{Identifiers: md.SecretIdentifiers, Success: true}, nil
	}

	si := []*vaultcommon.SecretIdentifier{}
	for _, id := range md.SecretIdentifiers {
		if id.Namespace == req.Namespace {
			si = append(si, id)
		}
	}

	return &vaultcommon.ListSecretIdentifiersResponse{
		Identifiers: si,
		Success:     true,
	}, nil
}

func (r *ReportingPlugin) observeDeleteSecrets(ctx context.Context, seqNr uint64, requestID string, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	l := r.typedRequestLggr(seqNr, requestID, "DeleteSecrets")
	tp := req.(*vaultcommon.DeleteSecretsRequest)
	o.RequestType = vaultcommon.RequestType_DELETE_SECRETS
	o.Request = &vaultcommon.Observation_DeleteSecretsRequest{
		DeleteSecretsRequest: tp,
	}

	requestsCountForID := buildSecretIdentifierCounts(tp.Ids)

	resps := []*vaultcommon.DeleteSecretResponse{}
	for _, id := range tp.Ids {
		validatedID, ierr := r.observeDeleteSecretRequest(ctx, reader, id, requestsCountForID)
		if ierr != nil {
			logUserErrorAware(l, "failed to handle delete secret request item", ierr, "id", id)
			errorMsg := userFacingError(ierr, "failed to handle delete secret request")
			resps = append(resps, &vaultcommon.DeleteSecretResponse{
				Id:      id,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			l.Debugw("observed delete secret request item", "id", validatedID)
			resps = append(resps, &vaultcommon.DeleteSecretResponse{
				Id: validatedID,
				// false because it hasn't been processed yet.
				// When the write is handled successfully in StateTransition
				// we'll update this to true.
				Success: false,
			})
		}
	}

	o.Response = &vaultcommon.Observation_DeleteSecretsResponse{
		DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
			Responses: resps,
		},
	}
}

func (r *ReportingPlugin) observeDeleteSecretRequest(ctx context.Context, reader ReadKVStore, identifier *vaultcommon.SecretIdentifier, requestsCountForID map[string]int) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateDeleteSecretsRequestItem(ctx, identifier, requestsCountForID)
	if err != nil {
		return id, err
	}

	ss, err := reader.GetSecret(ctx, id)
	if err != nil {
		return id, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if ss == nil {
		return id, vaulttypes.NewUserError("key does not exist")
	}

	return id, nil
}

func userFacingError(err error, fallback string) string {
	if vaulttypes.IsUserError(err) {
		return err.Error()
	}

	return fallback
}

func logUserErrorAware(l logger.Logger, msg string, err error, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, "error", err)
	lggr := l.Helper(1)
	if vaulttypes.IsUserError(err) {
		lggr.Debugw(msg, keysAndValues...)
		return
	}

	lggr.Errorw(msg, keysAndValues...)
}

func (r *ReportingPlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	valLggr := r.roundLggr(seqNr).With("oracleID", ao.Observer)
	obs := &vaultcommon.Observations{}
	if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
		valLggr.Debugw("validate observation failed", "error", err)
		return errors.New("failed to unmarshal observations: " + err.Error())
	}

	if len(ao.Observation) > r.maxObservationBytes {
		return fmt.Errorf("invalid observation: wire size %d exceeds max observation bytes %d", len(ao.Observation), r.maxObservationBytes)
	}

	if len(obs.SortNonce) != sortNonceLength {
		return fmt.Errorf("invalid observation: sort nonce must be %d bytes, got %d", sortNonceLength, len(obs.SortNonce))
	}

	if err := validatePendingQueueStallSignal(obs); err != nil {
		return fmt.Errorf("invalid observation: %w", err)
	}

	if observationSignalsPendingQueueStall(obs) {
		return nil
	}

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

	pendingQueueByID := map[string]*vaultcommon.StoredPendingQueueItem{}
	for _, item := range pendingQueueItems {
		pendingQueueByID[item.Id] = item
	}

	idToObs := map[string]*vaultcommon.Observation{}
	for _, o := range obs.Observations {
		err := r.validateContribution(ctx, pendingQueueByID[o.Id], o)
		if err != nil {
			valLggr.Debugw("validate observation failed", "requestID", o.Id, "error", err)
			return errors.New("invalid observation: " + err.Error())
		}

		_, seen := idToObs[o.Id]
		if seen {
			return errors.New("invalid observation: a single observation cannot contain duplicate observations for the same request id")
		}

		idToObs[o.Id] = o
	}

	// We expect
	// - that every request id corresponds to an item in the pending queue.
	//   This is because honest nodes may omit tail items when the full Observations proto would exceed the
	//   max observation byte limit.
	// - that all pending queue items can be fetched as blobs.
	if !gateAllows(ctx, r.lggr, r.cfg.VaultForceEmptyOCRRounds, "VaultForceEmptyOCRRounds") {
		if err := r.validatePendingQueueObservationsPrefix(pendingQueueItems, obs); err != nil {
			return err
		}
	}

	pendingIDs := map[string]bool{}
	for _, i := range pendingQueueItems {
		pendingIDs[i.Id] = true
	}
	for id := range idToObs {
		if !pendingIDs[id] {
			return fmt.Errorf("invalid observation: request id %s is not present in the pending queue", id)
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

// observerOkCoverage counts the distinct pending-queue ids for which the observer contributed
// an Ok observation. pendingIDs scopes coverage to the current queue (Byzantine ids outside the
// queue are ignored); pass nil to count all Ok contributions. Used to attribute prefix divergence
// to a specific oracle — a node consistently reporting lower coverage than peers is withholding
// or truncating its observation prefix, which stalls head-of-queue quorum under include-invalid.
func observerOkCoverage(obs *vaultcommon.Observations, pendingIDs map[string]bool) int {
	if obs == nil {
		return 0
	}
	seen := map[string]bool{}
	for _, o := range obs.Observations {
		if !observationContributionIsOk(o) {
			continue
		}
		if pendingIDs != nil && !pendingIDs[o.Id] {
			continue
		}
		seen[o.Id] = true
	}
	return len(seen)
}

// coverageSpread returns max-min of per-observer Ok prefix coverage. A non-zero spread means
// oracles disagree on how much of the pending queue they observed — the head-of-queue stall
// signature under include-invalid.
func coverageSpread(coverages []int) int {
	if len(coverages) == 0 {
		return 0
	}
	minC, maxC := coverages[0], coverages[0]
	for _, c := range coverages[1:] {
		if c < minC {
			minC = c
		}
		if c > maxC {
			maxC = c
		}
	}
	return maxC - minC
}

func (r *ReportingPlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	if !quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, r.onchainCfg.N, r.onchainCfg.F, aos) {
		return false, nil
	}

	if countPendingQueueStallSignals(aos) >= r.onchainCfg.F+1 {
		return true, nil
	}

	if gateAllows(ctx, r.lggr, r.cfg.VaultForceEmptyOCRRounds, "VaultForceEmptyOCRRounds") {
		return true, nil
	}

	readKV := NewReadStore(keyValueReader, r.metrics)
	pendingQueueItems, err := readKV.GetPendingQueue(ctx)
	if err != nil {
		return false, fmt.Errorf("could not fetch pending queue from store: %w", err)
	}
	if len(pendingQueueItems) == 0 {
		return true, nil
	}

	pendingIDs := map[string]bool{}
	for _, item := range pendingQueueItems {
		pendingIDs[item.Id] = true
	}
	okCount := 0
	errCount := 0
	coverages := make([]int, 0, len(aos))
	coverageByObserver := make(map[uint8]int, len(aos))
	for _, ao := range aos {
		obs := &vaultcommon.Observations{}
		if uerr := proto.Unmarshal([]byte(ao.Observation), obs); uerr != nil {
			continue
		}
		observer := uint8(ao.Observer)
		coverage := observerOkCoverage(obs, pendingIDs)
		coverages = append(coverages, coverage)
		coverageByObserver[observer] = coverage
		if len(obs.Observations) > 0 {
			headObs := obs.Observations[0]
			switch {
			case observationContributionIsErr(headObs):
				errCount++
			case observationContributionIsOk(headObs):
				okCount++
			}
		}
	}
	r.metrics.trackObservationPrefixCoverageSpread(ctx, coverageSpread(coverages))
	for observer, coverage := range coverageByObserver {
		r.metrics.trackObservationPrefixCoverage(ctx, observer, coverage)
	}
	if errCount >= r.onchainCfg.F+1 || okCount >= 2*r.onchainCfg.F+1 {
		return true, nil
	}

	return false, nil
}

func shaForProto(msg proto.Message) (string, error) {
	protoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("could not generate sha for proto message: failed to marshal proto: %w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(protoBytes)), nil
}

func (r *ReportingPlugin) shaForObservation(o *vaultcommon.Observation) (string, error) {
	switch o.RequestType {
	case vaultcommon.RequestType_GET_SECRETS:
		cloned := proto.CloneOf(o)
		for _, rsp := range cloned.GetGetSecretsResponse().Responses {
			if rsp.GetData() != nil {
				for _, es := range rsp.GetData().EncryptedDecryptionKeyShares {
					es.Shares = nil
					es.BinaryShares = nil
				}
			}
		}

		return shaForProto(cloned)
	default:
		return shaForProto(o)
	}
}

// GetSecrets request legitimacy: the request id must appear in >= 2F+1 observations total (across any SHA).
// With at most F Byzantine nodes, 2F+1 observations means at least F+1 honest nodes observed the
// request.
//
// Share sufficiency: encrypted shares cannot be cryptographically validated (ValidateObservation
// checks structure/size only), so a Byzantine node may emit a fake share that still matches an
// honest SHA. We therefore pick the largest same-SHA group of size >= F+1. All honest nodes read
// the same KV and produce the same SHA, so every honest observation (>= F+1 of them) falls in this
// group; a Byzantine-only group cannot reach F+1 (only F Byzantine exist), so the largest
// qualifying group is the honest one.
//
// Share threshold: we return all shares in the chosen group, capped at 2F+1. The group contains
// >= F+1 honest (valid) shares, and at most F shares in it can be Byzantine fakes. Returning up
// to 2F+1 shares guarantees that even after discarding up to F invalid shares, >= F+1 valid shares
// remain — enough to reconstruct (threshold T = F+1).
//
// Returns nil when either the request-legitimacy or share-sufficiency threshold is unmet.
func (r *ReportingPlugin) chooseGetSecretsObservations(totalForID int, shaToObs map[string][]*vaultcommon.Observation) []*vaultcommon.Observation {
	f := r.onchainCfg.F
	if totalForID < 2*f+1 {
		return nil
	}
	var chosen []*vaultcommon.Observation
	for _, sha := range slices.Sorted(maps.Keys(shaToObs)) {
		obs := shaToObs[sha]
		if len(obs) >= f+1 && len(obs) > len(chosen) {
			chosen = obs
		}
	}
	if maxChosen := 2*f + 1; len(chosen) > maxChosen {
		chosen = chosen[:maxChosen]
	}
	return chosen
}

func (r *ReportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	l := r.roundLggr(seqNr)
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

	if stallSignalCount := countPendingQueueStallSignalsInMap(marshalledObs); stallSignalCount >= r.onchainCfg.F+1 {
		return r.purgeStalledPendingQueue(ctx, l, writeKV, stallSignalCount)
	}

	// ---
	// Phase 1: Process requests from the pending queue by aggregating observations.
	// ---

	// obsMap is a map from request id -> list of observations across oracles.
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

	l.Debugw("stateTransition started", "oracleIDsToRequestIDs", oidsToReqIDs)

	os := &vaultcommon.Outcomes{
		Outcomes: []*vaultcommon.Outcome{},
	}

	var pendingQueueItems []*vaultcommon.StoredPendingQueueItem
	var pendingQueueErr error
	pendingQueueItems, pendingQueueErr = writeKV.GetPendingQueue(ctx)
	if pendingQueueErr != nil {
		return ocr3_1types.ReportsPlusPrecursor{}, fmt.Errorf("could not fetch pending queue during state transition: %w", pendingQueueErr)
	}

	pendingQueueByID := map[string]*vaultcommon.StoredPendingQueueItem{}
	for _, item := range pendingQueueItems {
		pendingQueueByID[item.Id] = item
	}

	idsToProcess := make([]string, 0, len(pendingQueueItems))
	for _, item := range pendingQueueItems {
		idsToProcess = append(idsToProcess, item.Id)
	}
	for _, id := range idsToProcess {
		obs, ok := obsMap[id]
		// This can only happen if the pending queue item is not in the obsMap
		// at which point we know any other requests in the pending queue can't be processed so we can break.
		if !ok {
			r.lggr.Warnw("no observations for pending queue item; stopping state transition pending queue processing", "id", id)
			break
		}

		okObs, errObs := classifyContributions(obs)
		f := r.onchainCfg.F
		if len(errObs) >= f+1 {
			requestType := vaultcommon.RequestType_UNKNOWN
			var payload proto.Message
			if item, found := pendingQueueByID[id]; found {
				if p, uerr := item.Item.UnmarshalNew(); uerr == nil {
					payload = p
					requestType = requestTypeForPayload(p)
				}
			}
			rejected := buildRejectedOutcome(id, payload, requestType, consensusObservationError(errObs, f))
			os.Outcomes = append(os.Outcomes, rejected)
			r.lggr.Infow("rejecting invalid pending queue item after f+1 error contributions",
				"seqNr", seqNr,
				"id", id,
				"errCount", len(errObs),
				"threshold", f+1,
			)
			if proto.Size(os) > r.maxReportsPlusPrecursorBytes {
				os.Outcomes = os.Outcomes[:len(os.Outcomes)-1]
				r.lggr.Warnw("state transition: rejected outcome exceeds max reports plus precursor bytes",
					"id", id,
					"maxReportsPlusPrecursorBytes", r.maxReportsPlusPrecursorBytes,
				)
				break
			}
			continue
		}
		if len(okObs) < 2*f+1 {
			r.lggr.Warnw("insufficient ok observations for pending queue item; stopping state transition",
				"id", id,
				"okCount", len(okObs),
				"errCount", len(errObs),
				"threshold", 2*f+1,
			)
			break
		}
		obs = okObs

		// For each observation we've received for a given Id,
		// we'll sha it and store it in `shaToObs`.
		// This means that each entry in `shaToObs` will contain a list of all
		// of the entries matching a given sha.
		shaToObs := map[string][]*vaultcommon.Observation{}
		for _, ob := range obs {
			sha, err := r.shaForObservation(ob)
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
		if len(obs) > 0 && obs[0].RequestType == vaultcommon.RequestType_GET_SECRETS {
			if chosen = r.chooseGetSecretsObservations(len(obs), shaToObs); chosen != nil {
				r.lggr.Debugw("sufficient observations for sha", "requestType", "GetSecrets", "relaxedConsensus", true, "totalForID", len(obs), "count", len(chosen), "threshold", r.onchainCfg.F+1, "id", id)
			}
		} else {
			for _, sha := range slices.Sorted(maps.Keys(shaToObs)) {
				obs := shaToObs[sha]

				if len(obs) >= r.onchainCfg.F+1 {
					// F+1 means that at least 1 honest node has provided this observation, so that's enough for all other request
					// types.
					// Technically we could have two shas with F+1 observations. If that happens we'll pick the last one.
					// This is deterministic since we're sorting by shas above.
					chosen = shaToObs[sha]
					l.Debugw("sufficient observations for sha", "sha", sha, "count", len(obs), "threshold", r.onchainCfg.F+1, "id", id)
				}
			}
		}

		if len(chosen) == 0 {
			shaToObsCount := map[string]int{}
			for sha, obs := range shaToObs {
				shaToObsCount[sha] = len(obs)
			}
			l.Warnw("insufficient observations found for requestID", "requestID", id, "shaToObsCount", shaToObsCount)
			break
		}

		// Defense in depth: re-run the same validateContribution used by the
		// Observation self-check and ValidateObservation on each chosen observation
		// before aggregating it. Honest nodes never fail this guard (chosen
		// observations already passed ValidateObservation); a failure means KV
		// divergence or a Byzantine observation, so we skip the item entirely.
		item := pendingQueueByID[id]
		guardFailed := false
		for _, ob := range chosen {
			if gerr := r.validateContribution(ctx, item, ob); gerr != nil {
				l.Warnw("state transition guard rejected chosen observation; skipping item",
					"seqNr", seqNr, "id", id, "error", gerr)
				guardFailed = true
				break
			}
		}
		if guardFailed {
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
			r.stateTransitionGetSecrets(chosen, o)
		case vaultcommon.RequestType_CREATE_SECRETS:
			r.stateTransitionCreateSecrets(ctx, writeKV, chosen, o)
		case vaultcommon.RequestType_UPDATE_SECRETS:
			r.stateTransitionUpdateSecrets(ctx, writeKV, chosen, o)
		case vaultcommon.RequestType_DELETE_SECRETS:
			r.stateTransitionDeleteSecrets(ctx, writeKV, chosen, o)
		case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
			r.stateTransitionListSecretIdentifiers(chosen, o)
		default:
			l.Debugw("unknown request type, skipping...", "requestType", first.RequestType, "id", id)
			continue
		}

		os.Outcomes = append(os.Outcomes, o)
		if proto.Size(os) > r.maxReportsPlusPrecursorBytes {
			os.Outcomes = os.Outcomes[:len(os.Outcomes)-1]
			l.Warnw("state transition: more observations than can be included in response",
				"requestID", id,
				"maxReportsPlusPrecursorBytes", r.maxReportsPlusPrecursorBytes,
				"packedOutcomeCount", len(os.Outcomes),
				"scheduledRequestIDs", len(obsMap),
			)
			break
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
		l.Debugw("State transition complete", "outcomeCount", len(os.Outcomes))
	}
	return ocr3_1types.ReportsPlusPrecursor(ospb), nil
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
				return fmt.Errorf("failed to fetch blob for pending queue item: %w", err)
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
	slices.SortFunc(keptItems, func(i, j *vaultcommon.StoredPendingQueueItem) int {
		return bytes.Compare(sortKey(i.Id, salt), sortKey(j.Id, salt))
	})

	if err := r.cfg.MaxPendingQueueWriteSize.Check(ctx, len(keptItems)); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[int]
		if !errors.As(err, &errBoundLimited) {
			return fmt.Errorf("failed to check pending queue write size limit: %w", err)
		}
		keptItems = keptItems[:errBoundLimited.Limit]
	}

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

func (r *ReportingPlugin) stateTransitionGetSecrets(chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	// Next, we deal with the responses.
	// For each request, we take the Id of the first observation
	// then aggregate the encrypted shares across all observations.
	// We sort these by Id and use the result as the response.
	idToAggResponse := map[string]*vaultcommon.SecretResponse{}
	for _, resp := range chosen {
		getSecretsResp := resp.GetGetSecretsResponse()
		for _, rsp := range getSecretsResp.Responses {
			key := vaulttypes.KeyFor(rsp.Id)
			mergedResp, ok := idToAggResponse[key]
			if !ok {
				idToAggResponse[key] = &vaultcommon.SecretResponse{
					Id:     rsp.Id,
					Result: rsp.Result,
				}
				continue
			}

			if rsp.GetData() != nil {
				r.aggregateGetSecretsShares(mergedResp, rsp)
			}
		}
	}

	sortedResponses := []*vaultcommon.SecretResponse{}
	for _, k := range slices.Sorted(maps.Keys(idToAggResponse)) {
		sortedResponses = append(sortedResponses, idToAggResponse[k])
	}

	o.Response = &vaultcommon.Outcome_GetSecretsResponse{
		GetSecretsResponse: &vaultcommon.GetSecretsResponse{
			Responses: sortedResponses,
		},
	}
}

func (r *ReportingPlugin) aggregateGetSecretsShares(
	mergedResp *vaultcommon.SecretResponse,
	rsp *vaultcommon.SecretResponse,
) {
	data := mergedResp.GetData()
	if len(data.EncryptedDecryptionKeyShares) == 0 {
		data.EncryptedDecryptionKeyShares = []*vaultcommon.EncryptedShares{}
	}

	keyToShares := map[string]*vaultcommon.EncryptedShares{}
	for _, s := range data.EncryptedDecryptionKeyShares {
		keyToShares[s.EncryptionKey] = s
	}

	for _, existing := range rsp.GetData().EncryptedDecryptionKeyShares {
		if shares, ok := keyToShares[existing.EncryptionKey]; ok {
			appendEncryptedShareEntry(shares, existing)
		} else {
			// This shouldn't happen -- this is because we're aggregating
			// requests that have a matching sha (excluding the decryption share).
			// Accordingly, we can assume that the request has been made with the same
			// set of encryption keys.
			r.lggr.Errorw("unexpected encryption key in response", "id", rsp.Id, "encryptionKey", existing.EncryptionKey)
		}
	}
}

func (r *ReportingPlugin) stateTransitionCreateSecrets(ctx context.Context, store WriteKVStore, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	first := chosen[0]
	reqID := first.GetCreateSecretsRequest().RequestId
	// First we'll aggregate the requests.
	// Since the shas for all requests match, we can just take the first entry
	// and sort the requests contained within it.
	req := first.GetCreateSecretsRequest().EncryptedSecrets
	idToReqs := map[string]*vaultcommon.EncryptedSecret{}
	for _, r := range req {
		idToReqs[vaulttypes.KeyFor(r.Id)] = r
	}

	// Next let's aggregate the responses.
	// We do this by taking the first response, and determine if
	// there was a validation error. If not, we write it to the key value store.
	// The responses are sorted by Id.
	resp := first.GetCreateSecretsResponse()
	idToResps := map[string]*vaultcommon.CreateSecretResponse{}
	for _, r := range resp.Responses {
		idToResps[vaulttypes.KeyFor(r.Id)] = r
	}

	sortedResps := []*vaultcommon.CreateSecretResponse{}
	for _, id := range slices.Sorted(maps.Keys(idToResps)) {
		resp := idToResps[id]
		req, found := idToReqs[id]
		if !found {
			// This shouldn't happen, as we've validated that the request and response
			// have the same number of items.
			r.lggr.Errorw("could not find request for response", "id", id, "requestID", reqID)
			sortedResps = append(sortedResps, &vaultcommon.CreateSecretResponse{
				Id:      resp.Id,
				Success: false,
				Error:   "internal error: could not find request for response",
			})
			continue
		}
		resp, err := r.stateTransitionCreateSecretsRequest(ctx, store, req, resp)
		if err != nil {
			logUserErrorAware(r.lggr, "failed to handle create secret request", err, "id", req.Id, "requestID", reqID)
			errorMsg := userFacingError(err, "failed to handle create secret request")
			sortedResps = append(sortedResps, &vaultcommon.CreateSecretResponse{
				Id:      req.Id,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			r.lggr.Debugw("successfully wrote secret to key value store", "method", "CreateSecrets", "key", vaulttypes.KeyFor(req.Id), "requestID", reqID)
			sortedResps = append(sortedResps, resp)
		}
	}

	o.Response = &vaultcommon.Outcome_CreateSecretsResponse{
		CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
			Responses: sortedResps,
		},
	}
}

func (r *ReportingPlugin) stateTransitionCreateSecretsRequest(ctx context.Context, store WriteKVStore, req *vaultcommon.EncryptedSecret, resp *vaultcommon.CreateSecretResponse) (*vaultcommon.CreateSecretResponse, error) {
	if resp.GetError() != "" {
		return resp, vaulttypes.NewUserError(resp.GetError())
	}

	encryptedSecret, err := decodeEncryptedSecretHex(req.EncryptedValue)
	if err != nil {
		return nil, err
	}

	secret, err := store.GetSecret(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret != nil {
		return nil, vaulttypes.NewUserError("could not write to key value store: key already exists")
	}

	count, err := store.GetSecretIdentifiersCountForOwner(ctx, req.Id.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret identifiers count for owner: %w", err)
	}

	// TODO orgID https://smartcontract-it.atlassian.net/browse/CRE-1707
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: req.Id.Owner})
	if ierr := r.cfg.MaxSecretsPerOwner.Check(ctx, count+1); ierr != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[int]](ierr); ok {
			return nil, vaulttypes.NewUserError(fmt.Sprintf("could not write to key value store: owner %s has reached maximum number of secrets (limit=%d)", req.Id.Owner, errBoundLimited.Limit))
		}
		return nil, fmt.Errorf("failed to check max secrets per owner limit: %w", ierr)
	}

	err = store.WriteSecret(ctx, req.Id, &vaultcommon.StoredSecret{
		EncryptedSecret: encryptedSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write secret to key value store: %w", err)
	}

	return &vaultcommon.CreateSecretResponse{
		Id:      req.Id,
		Success: true,
		Error:   "",
	}, nil
}

func (r *ReportingPlugin) stateTransitionUpdateSecrets(ctx context.Context, store WriteKVStore, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	first := chosen[0]
	reqID := first.GetUpdateSecretsRequest().RequestId
	// First we'll aggregate the requests.
	// Since the shas for all requests match, we can just take the first entry
	// and sort the requests contained within it.
	req := first.GetUpdateSecretsRequest().EncryptedSecrets
	idToReqs := map[string]*vaultcommon.EncryptedSecret{}
	for _, r := range req {
		idToReqs[vaulttypes.KeyFor(r.Id)] = r
	}

	// Next let's aggregate the responses.
	// We do this by taking the first response, and determine if
	// there was a validation error. If not, we write it to the key value store.
	// The responses are sorted by Id.
	resp := first.GetUpdateSecretsResponse()
	idToResps := map[string]*vaultcommon.UpdateSecretResponse{}
	for _, r := range resp.Responses {
		idToResps[vaulttypes.KeyFor(r.Id)] = r
	}

	sortedResps := []*vaultcommon.UpdateSecretResponse{}
	for _, id := range slices.Sorted(maps.Keys(idToResps)) {
		resp := idToResps[id]
		req, found := idToReqs[id]
		if !found {
			r.lggr.Errorw("could not find request for response", "id", id, "requestID", reqID)
			sortedResps = append(sortedResps, &vaultcommon.UpdateSecretResponse{
				Id:      resp.Id,
				Success: false,
				Error:   "internal error: could not find request for response",
			})
			continue
		}
		resp, err := r.stateTransitionUpdateSecretsRequest(ctx, store, req, resp)
		if err != nil {
			logUserErrorAware(r.lggr, "failed to handle update secret request", err, "id", req.Id, "requestID", reqID)
			errorMsg := userFacingError(err, "failed to handle update secret request")
			sortedResps = append(sortedResps, &vaultcommon.UpdateSecretResponse{
				Id:      req.Id,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			r.lggr.Debugw("successfully wrote secret to key value store", "method", "UpdateSecrets", "key", vaulttypes.KeyFor(req.Id), "requestID", reqID)
			sortedResps = append(sortedResps, resp)
		}
	}

	o.Response = &vaultcommon.Outcome_UpdateSecretsResponse{
		UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
			Responses: sortedResps,
		},
	}
}

func (r *ReportingPlugin) stateTransitionUpdateSecretsRequest(ctx context.Context, store WriteKVStore, req *vaultcommon.EncryptedSecret, resp *vaultcommon.UpdateSecretResponse) (*vaultcommon.UpdateSecretResponse, error) {
	if resp.GetError() != "" {
		return resp, vaulttypes.NewUserError(resp.GetError())
	}

	encryptedSecret, err := decodeEncryptedSecretHex(req.EncryptedValue)
	if err != nil {
		return nil, err
	}

	secret, err := store.GetSecret(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret == nil {
		return nil, vaulttypes.NewUserError("could not write update to key value store: key does not exist")
	}

	err = store.WriteSecret(ctx, req.Id, &vaultcommon.StoredSecret{
		EncryptedSecret: encryptedSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write secret to key value store: %w", err)
	}

	return &vaultcommon.UpdateSecretResponse{
		Id:      req.Id,
		Success: true,
		Error:   "",
	}, nil
}

func (r *ReportingPlugin) stateTransitionDeleteSecrets(ctx context.Context, store WriteKVStore, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	first := chosen[0]
	reqID := first.GetDeleteSecretsRequest().RequestId
	// First we'll aggregate the requests.
	// Since the shas for all requests match, we can just take the first entry
	// and sort the requests contained within it.
	req := first.GetDeleteSecretsRequest().Ids
	idToReqs := map[string]*vaultcommon.SecretIdentifier{}
	for _, r := range req {
		idToReqs[vaulttypes.KeyFor(r)] = r
	}

	// Next let's aggregate the responses.
	// We do this by taking the first response, and determine if
	// there was a validation error. If not, we write it to the key value store.
	// The responses are sorted by Id.
	resp := first.GetDeleteSecretsResponse()
	idToResps := map[string]*vaultcommon.DeleteSecretResponse{}
	for _, r := range resp.Responses {
		idToResps[vaulttypes.KeyFor(r.Id)] = r
	}

	sortedResps := []*vaultcommon.DeleteSecretResponse{}
	for _, id := range slices.Sorted(maps.Keys(idToResps)) {
		resp := idToResps[id]
		req, found := idToReqs[id]
		if !found {
			r.lggr.Errorw("could not find request for response", "id", id)
			sortedResps = append(sortedResps, &vaultcommon.DeleteSecretResponse{
				Id:      resp.Id,
				Success: false,
				Error:   "internal error: could not find request for response",
			})
			continue
		}
		resp, err := r.stateTransitionDeleteSecretsRequest(ctx, store, req, resp)
		if err != nil {
			logUserErrorAware(r.lggr, "failed to handle delete secret request", err, "id", id, "requestId", reqID)
			errorMsg := userFacingError(err, "failed to handle delete secret request")
			sortedResps = append(sortedResps, &vaultcommon.DeleteSecretResponse{
				Id:      req,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			r.lggr.Debugw("successfully deleted secret in key value store", "method", "DeleteSecrets", "key", vaulttypes.KeyFor(req), "requestId", reqID)
			sortedResps = append(sortedResps, resp)
		}
	}

	o.Response = &vaultcommon.Outcome_DeleteSecretsResponse{
		DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
			Responses: sortedResps,
		},
	}
}

func (r *ReportingPlugin) stateTransitionDeleteSecretsRequest(ctx context.Context, store WriteKVStore, id *vaultcommon.SecretIdentifier, resp *vaultcommon.DeleteSecretResponse) (*vaultcommon.DeleteSecretResponse, error) {
	if resp.GetError() != "" {
		return resp, vaulttypes.NewUserError(resp.GetError())
	}

	err := store.DeleteSecret(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete secret from key value store: %w", err)
	}

	return &vaultcommon.DeleteSecretResponse{
		Id:      id,
		Success: true,
		Error:   "",
	}, nil
}

func (r *ReportingPlugin) stateTransitionListSecretIdentifiers(chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	// All of the logic for the ListSecretIdentifiers request is in the
	// observation phase. This returns the observations in sorted order,
	// so we can just take the first aggregated response and use it as the outcome.
	first := chosen[0]
	o.Response = &vaultcommon.Outcome_ListSecretIdentifiersResponse{
		ListSecretIdentifiersResponse: first.GetListSecretIdentifiersResponse(),
	}
}

func (r *ReportingPlugin) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader) error {
	// Not currently used by the protocol, so we don't implement it.
	return errors.New("not implemented")
}

func (r *ReportingPlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	l := r.roundLggr(seqNr)
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
				l.Errorw("failed to generate Proto report", "error", err, "requestID", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_CREATE_SECRETS:
			createResp := proto.Clone(o.GetCreateSecretsResponse()).(*vaultcommon.CreateSecretsResponse)
			createResp.RequestId = o.Id
			rep, err := r.generateJSONReport(o.Id, o.RequestType, createResp)
			if err != nil {
				l.Errorw("failed to generate JSON report", "error", err, "requestID", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_UPDATE_SECRETS:
			updateResp := proto.Clone(o.GetUpdateSecretsResponse()).(*vaultcommon.UpdateSecretsResponse)
			updateResp.RequestId = o.Id
			rep, err := r.generateJSONReport(o.Id, o.RequestType, updateResp)
			if err != nil {
				l.Errorw("failed to generate JSON report", "error", err, "requestID", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_DELETE_SECRETS:
			deleteResp := proto.Clone(o.GetDeleteSecretsResponse()).(*vaultcommon.DeleteSecretsResponse)
			deleteResp.RequestId = o.Id
			rep, err := r.generateJSONReport(o.Id, o.RequestType, deleteResp)
			if err != nil {
				l.Errorw("failed to generate JSON report", "error", err, "requestID", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
			listResp := proto.Clone(o.GetListSecretIdentifiersResponse()).(*vaultcommon.ListSecretIdentifiersResponse)
			listResp.RequestId = o.Id
			rep, err := r.generateJSONReport(o.Id, o.RequestType, listResp)
			if err != nil {
				l.Errorw("failed to generate JSON report", "error", err, "requestID", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		default:
		}
	}

	if len(reports) > 0 {
		l.Debugw("Reports complete", "reportCount", len(reports))
	}
	return reports, nil
}

func (r *ReportingPlugin) generateProtoReport(id string, requestType vaultcommon.RequestType, msg proto.Message) (ocr3types.ReportWithInfo[[]byte], error) {
	if msg == nil {
		return ocr3types.ReportWithInfo[[]byte]{}, errors.New("invalid report: response cannot be nil")
	}

	rpb, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal response to proto: %w", err)
	}

	rip, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vaultcommon.ReportInfo{
		Id:          id,
		RequestType: requestType,
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_PROTOBUF,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return wrapReportWithKeyBundleInfo(rpb, rip)
}

func (r *ReportingPlugin) generateJSONReport(id string, requestType vaultcommon.RequestType, msg proto.Message) (ocr3types.ReportWithInfo[[]byte], error) {
	if msg == nil {
		return ocr3types.ReportWithInfo[[]byte]{}, errors.New("invalid report: response cannot be nil")
	}

	jsonb, err := vaultutils.ToCanonicalJSON(msg)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to convert proto to canonical JSON: %w", err)
	}

	rip, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vaultcommon.ReportInfo{
		Id:          id,
		RequestType: requestType,
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_JSON,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return wrapReportWithKeyBundleInfo(jsonb, rip)
}

func wrapReportWithKeyBundleInfo(report, reportInfo []byte) (ocr3types.ReportWithInfo[[]byte], error) {
	infos, err := structpb.NewStruct(map[string]any{
		// Use the EVM key bundle to sign the report.
		"keyBundleName": "evm",
		"reportInfo":    reportInfo,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, err
	}

	ip, err := proto.MarshalOptions{Deterministic: true}.Marshal(infos)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, err
	}

	return ocr3types.ReportWithInfo[[]byte]{
		Report: report,
		Info:   ip,
	}, nil
}

func (r *ReportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (r *ReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (r *ReportingPlugin) Close() error {
	return errors.Join(
		r.validator.Close(),
		r.cfg.MaxSecretsPerOwner.Close(),
		r.cfg.MaxShareLengthBytes.Close(),
		r.cfg.MaxBatchSize.Close(),
		r.cfg.MaxPendingQueueWriteSize.Close(),
		r.cfg.MaxBlobPayloadBytes.Close(),
		r.cfg.VaultForceEmptyOCRRounds.Close(),
		r.cfg.VaultPendingQueueStallThreshold.Close(),
	)
}
