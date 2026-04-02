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
	"regexp"
	"slices"
	"sort"
	"time"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"
	"github.com/smartcontractkit/smdkg/dkgocr"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

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

var (
	isValidIDComponent = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
)

type ReportingPluginConfig struct {
	LazyPublicKey *vaultcap.LazyPublicKey
	// Sourced from the DKG DB results package
	PublicKey       *tdh2easy.PublicKey
	PrivateKeyShare *tdh2easy.PrivateShare

	// Sourced from the offchain config
	MaxSecretsPerOwner                limits.BoundLimiter[int]
	MaxCiphertextLengthBytes          limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierKeyLengthBytes       limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierOwnerLengthBytes     limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierNamespaceLengthBytes limits.BoundLimiter[pkgconfig.Size]
	MaxShareLengthBytes               limits.BoundLimiter[pkgconfig.Size]
	MaxRequestBatchSize               limits.BoundLimiter[int]
	MaxBatchSize                      limits.BoundLimiter[int]
	OrgIDAsSecretOwnerEnabled         limits.GateLimiter
}

func NewReportingPluginFactory(
	lggr logger.Logger,
	store *requests.Store[*vaulttypes.Request],
	db dkgocrtypes.ResultPackageDatabase,
	recipientKey *dkgrecipientkey.Key,
	lazyPublicKey *vaultcap.LazyPublicKey,
	limitsFactory limits.Factory,
) (*ReportingPluginFactory, error) {
	if db == nil {
		return nil, errors.New("result package db cannot be nil")
	}

	if recipientKey == nil {
		return nil, errors.New("DKG recipient key cannot be nil when using result package db")
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
	}, nil
}

type ReportingPluginFactory struct {
	lggr          logger.Logger
	store         *requests.Store[*vaulttypes.Request]
	cfg           *ReportingPluginConfig
	db            dkgocrtypes.ResultPackageDatabase
	recipientKey  *dkgrecipientkey.Key
	limitsFactory limits.Factory
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

func initializePluginLimits(ctx context.Context, limitsFactory limits.Factory) (ocr3_1types.ReportingPluginLimits, error) {
	maxQueryBytes, err := cresettings.Default.VaultMaxQuerySizeLimit.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxQuerySizeLimit: %w", err)
	}
	maxObservationBytes, err := cresettings.Default.VaultMaxObservationSizeLimit.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxObservationSizeLimit: %w", err)
	}
	maxReportsPlusPrecursorBytes, err := cresettings.Default.VaultMaxReportsPlusPrecursorSizeLimit.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxReportsPlusPrecursorSizeLimit: %w", err)
	}
	maxReportBytes, err := cresettings.Default.VaultMaxReportSizeLimit.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxReportSizeLimit: %w", err)
	}
	maxReportCount, err := cresettings.Default.VaultMaxReportCount.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxReportCount: %w", err)
	}
	maxKVModifiedKeysPlusValuesBytes, err := cresettings.Default.VaultMaxKeyValueModifiedKeysPlusValuesSizeLimit.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxKeyValueModifiedKeysPlusValuesSizeLimit: %w", err)
	}
	maxKVModifiedKeys, err := cresettings.Default.VaultMaxKeyValueModifiedKeys.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxKeyValueModifiedKeys: %w", err)
	}
	maxBlobPayloadBytes, err := cresettings.Default.VaultMaxBlobPayloadSizeLimit.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxBlobPayloadSizeLimit: %w", err)
	}
	maxPerOracleUnexpiredBlobCumulativePayloadBytes, err := cresettings.Default.VaultMaxPerOracleUnexpiredBlobCumulativePayloadSizeLimit.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxPerOracleUnexpiredBlobCumulativePayloadSizeLimit: %w", err)
	}
	maxPerOracleUnexpiredBlobCount, err := cresettings.Default.VaultMaxPerOracleUnexpiredBlobCount.GetOrDefault(ctx, limitsFactory.Settings)
	if err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("VaultMaxPerOracleUnexpiredBlobCount: %w", err)
	}

	return ocr3_1types.ReportingPluginLimits{
		MaxQueryBytes:                                   int(maxQueryBytes),
		MaxObservationBytes:                             int(maxObservationBytes),
		MaxReportsPlusPrecursorBytes:                    int(maxReportsPlusPrecursorBytes),
		MaxReportBytes:                                  int(maxReportBytes),
		MaxReportCount:                                  maxReportCount,
		MaxKeyValueModifiedKeysPlusValuesBytes:          int(maxKVModifiedKeysPlusValuesBytes),
		MaxKeyValueModifiedKeys:                         maxKVModifiedKeys,
		MaxBlobPayloadBytes:                             int(maxBlobPayloadBytes),
		MaxPerOracleUnexpiredBlobCumulativePayloadBytes: int(maxPerOracleUnexpiredBlobCumulativePayloadBytes),
		MaxPerOracleUnexpiredBlobCount:                  maxPerOracleUnexpiredBlobCount,
	}, nil
}

func newReportingPluginConfigLimiters(factory limits.Factory) (*ReportingPluginConfig, error) {
	maxCiphertextLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultCiphertextSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultCiphertextSizeLimit: %w", err)
	}

	maxIdentifierKeyLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultIdentifierKeySizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultIdentifierKeySizeLimit: %w", err)
	}

	maxIdentifierOwnerLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultIdentifierOwnerSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultIdentifierOwnerSizeLimit: %w", err)
	}

	maxIdentifierNamespaceLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultIdentifierNamespaceSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultIdentifierNamespaceSizeLimit: %w", err)
	}

	maxShareLengthBytesLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultShareSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultShareSizeLimit: %w", err)
	}

	maxRequestBatchSizeLimiter, err := limits.MakeUpperBoundLimiter(factory, cresettings.Default.VaultRequestBatchSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("VaultRequestBatchSizeLimit: %w", err)
	}

	orgIDAsSecretOwnerEnabled, err := limits.MakeGateLimiter(factory, cresettings.Default.VaultOrgIdAsSecretOwnerEnabled)
	if err != nil {
		return nil, fmt.Errorf("VaultOrgIDAsSecretOwnerEnabled: %w", err)
	}

	return &ReportingPluginConfig{
		MaxShareLengthBytes:               maxShareLengthBytesLimiter,
		MaxRequestBatchSize:               maxRequestBatchSizeLimiter,
		MaxCiphertextLengthBytes:          maxCiphertextLengthBytesLimiter,
		MaxIdentifierKeyLengthBytes:       maxIdentifierKeyLengthBytesLimiter,
		MaxIdentifierOwnerLengthBytes:     maxIdentifierOwnerLengthBytesLimiter,
		MaxIdentifierNamespaceLengthBytes: maxIdentifierNamespaceLengthBytesLimiter,
		OrgIDAsSecretOwnerEnabled:         orgIDAsSecretOwnerEnabled,
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

	r.lggr.Debugw("instantiating VaultReportingPlugin with config",
		"maxSecretsPerOwner", logLimit(ctx, r.lggr, cfg.MaxSecretsPerOwner),
		"maxCiphertextLengthBytes", logLimit(ctx, r.lggr, cfg.MaxCiphertextLengthBytes),
		"maxIdentifierKeyLengthBytes", logLimit(ctx, r.lggr, cfg.MaxIdentifierKeyLengthBytes),
		"maxIdentifierOwnerLengthBytes", logLimit(ctx, r.lggr, cfg.MaxIdentifierOwnerLengthBytes),
		"maxIdentifierNamespaceLengthBytes", logLimit(ctx, r.lggr, cfg.MaxIdentifierNamespaceLengthBytes),
		"maxRequestBatchSize", logLimit(ctx, r.lggr, cfg.MaxRequestBatchSize),
		"maxShareLengthBytes", logLimit(ctx, r.lggr, cfg.MaxShareLengthBytes),
		"batchSize", logLimit(ctx, r.lggr, cfg.MaxBatchSize),
	)

	metrics, err := newPluginMetrics(config.ConfigDigest.String())
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not create plugin metrics: %w", err)
	}

	pluginLimits, err := initializePluginLimits(ctx, r.limitsFactory)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo1{}, fmt.Errorf("could not resolve plugin limits: %w", err)
	}

	return &ReportingPlugin{
			lggr:       r.lggr.Named("VaultReportingPlugin"),
			store:      r.store,
			cfg:        cfg,
			metrics:    metrics,
			onchainCfg: config,
			unmarshalBlob: func(data []byte) (ocr3_1types.BlobHandle, error) {
				handle := ocr3_1types.BlobHandle{}
				err := handle.UnmarshalBinary(data)
				return handle, err
			},
			marshalBlob: func(handle ocr3_1types.BlobHandle) ([]byte, error) {
				return handle.MarshalBinary()
			},
		}, ocr3_1types.ReportingPluginInfo1{
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

func (r *ReportingPlugin) orgIDAsSecretOwnerEnabled(ctx context.Context) bool {
	return r.cfg.OrgIDAsSecretOwnerEnabled.AllowErr(ctx) == nil
}

type pendingQueueStore interface {
	WritePendingQueue(ctx context.Context, pending []*vaultcommon.StoredPendingQueueItem) error
}

func (r *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	start := time.Now()
	defer func() {
		r.lggr.Debugw("observation finished", "seqNr", seqNr, "elapsed", time.Since(start))
	}()

	wrappedReadStore := NewKVStoreWrapper(NewReadStore(keyValueReader, r.metrics), r.orgIDAsSecretOwnerEnabled(ctx), r.lggr)

	batch, err := wrappedReadStore.GetPendingQueue(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch of requests: %w", err)
	}
	// Avoid log spam by only logging if we have any requests to process.
	if len(batch) > 0 {
		mbs, _ := r.cfg.MaxBatchSize.Limit(ctx)
		r.lggr.Debugw("observation started", "seqNr", seqNr, "batchSize", mbs)
	}

	ids := []string{}
	obs := []*vaultcommon.Observation{}
	for _, req := range batch {
		o := &vaultcommon.Observation{
			Id: req.Id,
		}
		ids = append(ids, req.Id)

		payload, ierr := req.Item.UnmarshalNew()
		if ierr != nil {
			r.lggr.Errorw("failed to unmarshal request payload", "id", req.Id, "error", ierr)
			continue
		}

		switch tp := payload.(type) {
		case *vaultcommon.GetSecretsRequest:
			r.observeGetSecrets(ctx, wrappedReadStore.WithRequest(tp.OrgId, tp.WorkflowOwner), tp, o)
		case *vaultcommon.CreateSecretsRequest:
			r.observeCreateSecrets(ctx, wrappedReadStore.WithRequest(tp.OrgId, tp.WorkflowOwner), tp, o)
		case *vaultcommon.UpdateSecretsRequest:
			r.observeUpdateSecrets(ctx, wrappedReadStore.WithRequest(tp.OrgId, tp.WorkflowOwner), tp, o)
		case *vaultcommon.DeleteSecretsRequest:
			r.observeDeleteSecrets(ctx, wrappedReadStore.WithRequest(tp.OrgId, tp.WorkflowOwner), tp, o)
		case *vaultcommon.ListSecretIdentifiersRequest:
			r.observeListSecretIdentifiers(ctx, wrappedReadStore.WithRequest(tp.OrgId, tp.WorkflowOwner), tp, o)
		default:
			r.lggr.Errorw("unknown request type, skipping...", "requestType", fmt.Sprintf("%T", payload), "id", req.Id)
			continue
		}

		obs = append(obs, o)
	}

	obspb := &vaultcommon.Observations{
		Observations: obs,
	}

	// First, observe the pending queue that I have.
	// This will get aggregated in the state transition phase
	// to form the DON wide pending queue.
	localQueueItems, ierr := r.store.All()
	if ierr != nil {
		return nil, ierr
	}

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

	// Next, get the current pending queue. We'll use this to dedupe
	// requests when generating an observation for the next state of the
	// pending queue.
	pendingQueue, ierr := wrappedReadStore.GetPendingQueue(ctx)
	if ierr != nil {
		return nil, ierr
	}

	pendingQueueHasID := map[string]bool{}
	for _, item := range pendingQueue {
		pendingQueueHasID[item.Id] = true
	}

	blobPayloads := make([][]byte, 0, len(localQueueItems))
	blobPayloadIDs := make([]string, 0, len(localQueueItems))
	maxObservedLocalQueueItems := 0
	for _, item := range localQueueItems {
		// The item is already in the pending queue. We'll be processing it
		// this round. Let's skip it for now so we don't process duplicates.
		if pendingQueueHasID[item.ID()] {
			continue
		}

		anyMsg, ierr2 := anypb.New(item.Payload)
		if ierr2 != nil {
			return nil, fmt.Errorf("could not marshal request payload to Any: %w", ierr2)
		}

		item := &vaultcommon.StoredPendingQueueItem{
			Id:   item.ID(),
			Item: anyMsg,
		}

		itemb, ierr2 := proto.Marshal(item)
		if ierr2 != nil {
			return nil, fmt.Errorf("could not marshal pending queue item: %w", ierr2)
		}

		if maxObservedLocalQueueItems == 0 {
			l, ierr2 := r.cfg.MaxBatchSize.Limit(ctx)
			if ierr2 != nil {
				return nil, fmt.Errorf("could not fetch max batch size limit: %w", ierr2)
			}
			maxObservedLocalQueueItems = 2 * l
		}

		blobPayloads = append(blobPayloads, itemb)
		blobPayloadIDs = append(blobPayloadIDs, item.Id)

		if len(blobPayloads) >= maxObservedLocalQueueItems {
			r.lggr.Warnw("Observed local queue exceeds batch size limit, truncating",
				"queueSize", len(blobPayloads),
				"batchSizeLimit", maxObservedLocalQueueItems)
			r.metrics.trackQueueOverflow(ctx, len(blobPayloads), maxObservedLocalQueueItems)
			break
		}
	}

	pendingQueueItems, err := r.broadcastBlobPayloads(ctx, blobBroadcastFetcher, seqNr, blobPayloads, blobPayloadIDs)
	if err != nil {
		return nil, err
	}
	obspb.PendingQueueItems = pendingQueueItems

	// Second, generate a random nonce that we'll use to sort the observations.
	// Each node generates a nonce idepedently, to be concatenated later on.
	nonce, ierr := generateRandomNonce()
	if ierr != nil {
		return nil, fmt.Errorf("could not generate nonce for observation: %w", ierr)
	}

	obspb.SortNonce = nonce

	obsb, err := proto.MarshalOptions{Deterministic: true}.Marshal(obspb)
	if err != nil {
		return nil, fmt.Errorf("could not marshal observations: %w", err)
	}

	// Avoid log spam by only logging if we have any requests to process.
	if len(batch) > 0 {
		r.lggr.Debugw("observation complete", "ids", ids, "batchSize", len(batch))
	}
	return types.Observation(obsb), nil
}

// broadcastBlobPayloads broadcasts each payload as a blob in parallel to reduce
// Observation() latency (shortening this phase helps the OCR round finish within
// DeltaProgress). Each call is given a 2-second timeout so that a single slow
// broadcast cannot stall the entire batch. Individual broadcast failures are logged
// and skipped rather than aborting the entire observation, so that one problematic
// payload does not prevent the remaining items from being observed. Context
// cancellation/deadline errors on the parent context are propagated immediately so
// that expired rounds fail fast.
func (r *ReportingPlugin) broadcastBlobPayloads(
	ctx context.Context,
	fetcher ocr3_1types.BlobBroadcastFetcher,
	seqNr uint64,
	payloads [][]byte,
	requestIDs []string,
) ([][]byte, error) {
	results := make([][]byte, len(payloads))

	start := time.Now()
	defer func() {
		r.lggr.Debugw("observation blob broadcast finished", "seqNr", seqNr, "blobCount", len(payloads), "elapsed", time.Since(start))
	}()

	const perBlobTimeout = 2 * time.Second
	var g errgroup.Group
	for i, payload := range payloads {
		g.Go(func() error {
			broadcastCtx, cancel := context.WithTimeout(ctx, perBlobTimeout)
			defer cancel()

			blobHandle, err := fetcher.BroadcastBlob(broadcastCtx, payload, ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: seqNr + 2})
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				r.lggr.Warnw("failed to broadcast pending queue item as blob, skipping",
					"seqNr", seqNr,
					"requestID", requestIDs[i],
					"err", err)
				return nil
			}

			blobHandleBytes, err := r.marshalBlob(blobHandle)
			if err != nil {
				r.lggr.Warnw("failed to marshal blob handle, skipping",
					"seqNr", seqNr,
					"requestID", requestIDs[i],
					"err", err)
				return nil
			}

			results[i] = blobHandleBytes
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

func (r *ReportingPlugin) observeGetSecrets(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.GetSecretsRequest)
	o.RequestType = vaultcommon.RequestType_GET_SECRETS
	o.Request = &vaultcommon.Observation_GetSecretsRequest{
		GetSecretsRequest: tp,
	}
	resps := []*vaultcommon.SecretResponse{}
	for _, secretRequest := range tp.Requests {
		resp, ierr := r.observeGetSecretsRequest(ctx, reader, secretRequest, tp.WorkflowOwner, tp.OrgId)
		if ierr != nil {
			logUserErrorAware(r.lggr, "failed to observe get secret request item", ierr, "id", secretRequest.Id)
			errorMsg := userFacingError(ierr, "failed to handle get secret request")
			resps = append(resps, &vaultcommon.SecretResponse{
				Id: secretRequest.Id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: errorMsg,
				},
			})
		} else {
			r.lggr.Debugw("observed get secret request item", "id", resp.Id)
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

func (s *share) encryptWithKey(pk string) (string, error) {
	publicKey, err := hex.DecodeString(pk)
	if err != nil {
		return "", newUserError("failed to convert public key to bytes: " + err.Error())
	}

	if len(publicKey) != curve25519.PointSize {
		return "", newUserError(fmt.Sprintf("invalid public key size: expected %d bytes, got %d bytes", curve25519.PointSize, len(publicKey)))
	}

	publicKeyLength := [curve25519.PointSize]byte(publicKey)
	encrypted, err := box.SealAnonymous(nil, s.data, &publicKeyLength, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt decryption share: %w", err)
	}

	return hex.EncodeToString(encrypted), nil
}

func generatePlaintextShare(publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, encryptedSecret []byte, workflowOwner string, orgID string) (*share, error) {
	ct := &tdh2easy.Ciphertext{}
	err := ct.UnmarshalVerify(encryptedSecret, publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ciphertext: %w", err)
	}

	es := hex.EncodeToString(encryptedSecret)
	err = vaultcap.EnsureRightLabelOnSecret(publicKey, es, workflowOwner, orgID)
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

func (r *ReportingPlugin) observeGetSecretsRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.SecretRequest, workflowOwner string, orgID string) (*vaultcommon.SecretResponse, error) {
	id, err := r.validateSecretIdentifier(ctx, secretRequest.Id)
	if err != nil {
		return nil, err
	}

	secret, err := reader.GetSecret(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}
	if secret == nil {
		return nil, newUserError("key does not exist")
	}

	if !r.orgIDAsSecretOwnerEnabled(ctx) {
		orgID = ""
	}
	sh, err := generatePlaintextShare(r.cfg.PublicKey, r.cfg.PrivateKeyShare, secret.EncryptedSecret, workflowOwner, orgID)
	if err != nil {
		return nil, err
	}

	shares := []*vaultcommon.EncryptedShares{}
	for _, pk := range secretRequest.EncryptionKeys {
		encShare, err := sh.encryptWithKey(pk)
		if err != nil {
			return nil, err
		}

		shares = append(shares, &vaultcommon.EncryptedShares{
			EncryptionKey: pk,
			Shares: []string{
				encShare,
			},
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

func (r *ReportingPlugin) observeCreateSecrets(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.CreateSecretsRequest)
	o.RequestType = vaultcommon.RequestType_CREATE_SECRETS
	o.Request = &vaultcommon.Observation_CreateSecretsRequest{
		CreateSecretsRequest: tp,
	}
	l := r.lggr.With("requestID", tp.RequestId, "requestType", "CreateSecrets")

	requestsCountForID := map[string]int{}
	for _, sr := range tp.EncryptedSecrets {
		var key string
		// This can happen if a user provides a malformed request.
		// We validate this case away in `handleCreateSecretRequest`,
		// but need to still handle it here to avoid panics.
		if sr.Id == nil {
			key = "<nil>"
		} else {
			key = vaulttypes.KeyFor(sr.Id)
		}
		requestsCountForID[key]++
	}

	resps := []*vaultcommon.CreateSecretResponse{}
	for _, sr := range tp.EncryptedSecrets {
		validatedID, ierr := r.observeCreateSecretRequest(ctx, reader, sr, requestsCountForID, tp.WorkflowOwner, tp.OrgId)
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

func (r *ReportingPlugin) observeCreateSecretRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.EncryptedSecret, requestsCountForID map[string]int, workflowOwner string, orgID string) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateSecretIdentifier(ctx, secretRequest.Id)
	if err != nil {
		return id, err
	}

	if requestsCountForID[vaulttypes.KeyFor(secretRequest.Id)] > 1 {
		return id, newUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}

	if ierr := r.validateCiphertextSize(ctx, secretRequest.Id.Owner, secretRequest.EncryptedValue); ierr != nil {
		return id, newUserError(ierr.Error())
	}

	if !r.orgIDAsSecretOwnerEnabled(ctx) {
		orgID = ""
	}
	err = vaultcap.EnsureRightLabelOnSecret(r.cfg.PublicKey, secretRequest.EncryptedValue, workflowOwner, orgID)
	if err != nil {
		return id, newUserError("failed to verify ciphertext: " + err.Error())
	}

	// Other verifications, such as checking whether the key already exists,
	// or whether we have hit the limit on the number of secrets per owner,
	// are done in the StateTransition phase.
	// This guarantees that we correctly account for changes made in other requests
	// in the batch.
	return id, nil
}

func (r *ReportingPlugin) observeUpdateSecrets(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.UpdateSecretsRequest)
	o.RequestType = vaultcommon.RequestType_UPDATE_SECRETS
	o.Request = &vaultcommon.Observation_UpdateSecretsRequest{
		UpdateSecretsRequest: tp,
	}
	l := r.lggr.With("requestID", tp.RequestId, "requestType", "UpdateSecrets")

	requestsCountForID := map[string]int{}
	for _, sr := range tp.EncryptedSecrets {
		var key string
		// This can happen if a user provides a malformed request.
		// We validate this case away in `handleCreateSecretRequest`,
		// but need to still handle it here to avoid panics.
		if sr.Id == nil {
			key = "<nil>"
		} else {
			key = vaulttypes.KeyFor(sr.Id)
		}
		requestsCountForID[key]++
	}

	resps := []*vaultcommon.UpdateSecretResponse{}
	for _, sr := range tp.EncryptedSecrets {
		validatedID, ierr := r.observeUpdateSecretRequest(ctx, reader, sr, requestsCountForID, tp.WorkflowOwner, tp.OrgId)
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

func (r *ReportingPlugin) observeUpdateSecretRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.EncryptedSecret, requestsCountForID map[string]int, workflowOwner string, orgID string) (*vaultcommon.SecretIdentifier, error) {
	// The checks at this stage are identical since we only check the correctness of the payload
	// at this stage. Checks that are different between update and create, like whether the secret already exists,
	// are handled in the StateTransition phase.
	return r.observeCreateSecretRequest(ctx, reader, secretRequest, requestsCountForID, workflowOwner, orgID)
}

func (r *ReportingPlugin) observeListSecretIdentifiers(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.ListSecretIdentifiersRequest)
	o.RequestType = vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS
	o.Request = &vaultcommon.Observation_ListSecretIdentifiersRequest{
		ListSecretIdentifiersRequest: tp,
	}
	l := r.lggr.With("requestId", tp.RequestId, "requestType", "ListSecretIdentifiers", "owner", tp.Owner)

	resp, err := r.processListSecretIdentifiersRequest(ctx, l, reader, tp)
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

func (r *ReportingPlugin) processListSecretIdentifiersRequest(ctx context.Context, l logger.Logger, reader ReadKVStore, req *vaultcommon.ListSecretIdentifiersRequest) (*vaultcommon.ListSecretIdentifiersResponse, error) {
	if req.Owner == "" {
		return nil, errors.New("invalid request: owner cannot be empty")
	}

	md, err := reader.GetMetadata(ctx, req.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for owner: %w", err)
	}

	if md == nil {
		// No metadata, so the list is empty.
		// The user hasn't added any items to the vault DON yet.
		l.Debugw("successfully read metadata for owner: no metadata found, returning empty list")
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

func (r *ReportingPlugin) observeDeleteSecrets(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.DeleteSecretsRequest)
	o.RequestType = vaultcommon.RequestType_DELETE_SECRETS
	o.Request = &vaultcommon.Observation_DeleteSecretsRequest{
		DeleteSecretsRequest: tp,
	}
	l := r.lggr.With("requestId", tp.RequestId, "requestType", "DeleteSecrets")

	requestsCountForID := map[string]int{}
	for _, sr := range tp.Ids {
		var key string
		// This can happen if a user provides a malformed request.
		// We validate this case away in `handleCreateSecretRequest`,
		// but need to still handle it here to avoid panics.
		if sr == nil {
			key = "<nil>"
		} else {
			key = vaulttypes.KeyFor(sr)
		}
		requestsCountForID[key]++
	}

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
	id, err := r.validateSecretIdentifier(ctx, identifier)
	if err != nil {
		return id, err
	}

	if requestsCountForID[vaulttypes.KeyFor(identifier)] > 1 {
		return id, newUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}

	ss, err := reader.GetSecret(ctx, id)
	if err != nil {
		return id, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if ss == nil {
		return id, newUserError("key does not exist")
	}

	return id, nil
}

func (r *ReportingPlugin) validateSecretIdentifier(ctx context.Context, id *vaultcommon.SecretIdentifier) (*vaultcommon.SecretIdentifier, error) {
	if id == nil {
		return nil, newUserError("invalid secret identifier: cannot be nil")
	}

	if id.Key == "" {
		return nil, newUserError("invalid secret identifier: key cannot be empty")
	}

	if id.Owner == "" {
		return nil, newUserError("invalid secret identifier: owner cannot be empty")
	}

	namespace := id.Namespace
	if namespace == "" {
		namespace = vaulttypes.DefaultNamespace
	}

	if !isValidIDComponent(id.Key) || !isValidIDComponent(id.Owner) || !isValidIDComponent(namespace) {
		return nil, newUserError("invalid secret identifier: key, owner and namespace must only contain alphanumeric characters")
	}

	newID := &vaultcommon.SecretIdentifier{
		Key:       id.Key,
		Owner:     id.Owner,
		Namespace: namespace,
	}

	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: id.Owner})
	if err := r.cfg.MaxIdentifierOwnerLengthBytes.Check(ctx, pkgconfig.Size(len(id.Owner))); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return nil, newUserError(fmt.Sprintf("invalid secret identifier: owner exceeds maximum length of %s", errBoundLimited.Limit))
		}
		return nil, newUserError("failed to check owner length limit: " + err.Error())
	}

	if err := r.cfg.MaxIdentifierNamespaceLengthBytes.Check(ctx, pkgconfig.Size(len(id.Namespace))); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return nil, newUserError(fmt.Sprintf("invalid secret identifier: namespace exceeds maximum length of %s", errBoundLimited.Limit))
		}
		return nil, newUserError("failed to check namespace length limit: " + err.Error())
	}

	if err := r.cfg.MaxIdentifierKeyLengthBytes.Check(ctx, pkgconfig.Size(len(id.Key))); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return nil, newUserError(fmt.Sprintf("invalid secret identifier: key exceeds maximum length of %s", errBoundLimited.Limit))
		}
		return nil, newUserError("failed to check key length limit: " + err.Error())
	}
	return newID, nil
}

func newUserError(msg string) *userError {
	return &userError{msg: msg}
}

type userError struct {
	msg string
}

func (u *userError) Error() string {
	return u.msg
}

func (u *userError) Is(target error) bool {
	_, ok := target.(*userError)
	return ok
}

func userFacingError(err error, fallback string) string {
	if errors.Is(err, &userError{}) {
		return err.Error()
	}

	return fallback
}

func logUserErrorAware(l logger.Logger, msg string, err error, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error", err)
	lggr := l.Helper(1)
	if errors.Is(err, &userError{}) {
		lggr.Debugw(msg, keysAndValues...)
		return
	}

	lggr.Errorw(msg, keysAndValues...)
}

func (r *ReportingPlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	obs := &vaultcommon.Observations{}
	if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
		return errors.New("failed to unmarshal observations: " + err.Error())
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
	// - an observation for each item in the pending queue.
	//   This is because honest nodes will all be reading from
	//   the same deterministic key-value store-based queue.
	// - that all pending queue items can be fetched as blobs.
	wrappedStore := NewKVStoreWrapper(NewReadStore(keyValueReader, r.metrics), r.orgIDAsSecretOwnerEnabled(ctx), r.lggr)
	pendingQueueItems, err := wrappedStore.GetPendingQueue(ctx)
	if err != nil {
		return fmt.Errorf("could not fetch pending queue from store: %w", err)
	}

	if len(idToObs) != len(pendingQueueItems) {
		return errors.New("invalid observation: number of observations doesn't match number of pending requests")
	}

	for _, i := range pendingQueueItems {
		_, seen := idToObs[i.Id]
		if !seen {
			return fmt.Errorf("invalid observation: missing observation for pending request id %s", i.Id)
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

	seen := map[string]bool{}
	for _, i := range obs.PendingQueueItems {
		bh, err := r.unmarshalBlob(i)
		if err != nil {
			return fmt.Errorf("could not unmarshal blob handle from observation pending queue item: %w", err)
		}

		blob, err := blobFetcher.FetchBlob(ctx, bh)
		if err != nil {
			return fmt.Errorf("could not fetch blob for observation pending queue item: %w", err)
		}

		sha := fmt.Sprintf("%x", sha256.Sum256(blob))
		if seen[sha] {
			return errors.New("duplicate item found in pending queue item observation")
		}
		seen[sha] = true
	}

	return nil
}

func (r *ReportingPlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumNMinusF, r.onchainCfg.N, r.onchainCfg.F, aos), nil
}

func shaForProto(msg proto.Message) (string, error) {
	protoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("could not generate sha for proto message: failed to marshal proto: %w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(protoBytes)), nil
}

func shaForObservation(o *vaultcommon.Observation) (string, error) {
	switch o.RequestType {
	case vaultcommon.RequestType_GET_SECRETS:
		cloned := proto.CloneOf(o)
		for _, r := range cloned.GetGetSecretsResponse().Responses {
			if r.GetData() != nil {
				// Exclude the encrypted shares from the sha, as these need to be aggregated later.
				r.GetData().EncryptedDecryptionKeyShares = nil
			}
		}

		return shaForProto(cloned)
	default:
		return shaForProto(o)
	}
}

func (r *ReportingPlugin) checkRequestBatchLimit(ctx context.Context, batchSize int) error {
	if err := r.cfg.MaxRequestBatchSize.Check(ctx, batchSize); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[int]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("max batch size exceeded for request: %w", err)
		}
		// Fail closed here: this could cause a loss of liveness but
		// the current implementation would only return an error that's
		// not a ErrorBoundLimited if the limiter has been closed.
		return errors.New("failed to check batch size")
	}

	return nil
}

func (r *ReportingPlugin) validateObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.Id == "" {
		return errors.New("observation id cannot be empty")
	}

	switch o.RequestType {
	case vaultcommon.RequestType_GET_SECRETS:
		return r.validateGetSecretsObservation(ctx, o)
	case vaultcommon.RequestType_CREATE_SECRETS:
		return r.validateCreateSecretsObservation(ctx, o)
	case vaultcommon.RequestType_UPDATE_SECRETS:
		return r.validateUpdateSecretsObservation(ctx, o)
	case vaultcommon.RequestType_DELETE_SECRETS:
		return r.validateDeleteSecretsObservation(ctx, o)
	case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
		return r.validateListSecretIdentifiersObservation(ctx, o)
	default:
		return errors.New("invalid observation type: " + o.RequestType.String())
	}
}

func (r *ReportingPlugin) validateGetSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetGetSecretsRequest() == nil || o.GetGetSecretsResponse() == nil {
		return errors.New("GetSecrets observation must have both request and response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(o.GetGetSecretsRequest().Requests)); err != nil {
		return err
	}

	if len(o.GetGetSecretsRequest().Requests) != len(o.GetGetSecretsResponse().Responses) {
		return errors.New("GetSecrets request and response must have the same number of items")
	}

	// check for that we have an entry per encrypted key in the request
	// we should have max 1 share per observation per encrypted key
	req, resp := o.GetGetSecretsRequest(), o.GetGetSecretsResponse()
	reqMap := map[string]*vaultcommon.SecretRequest{}
	for _, r := range req.Requests {
		if r.Id == nil {
			return errors.New("GetSecrets request contains nil secret identifier")
		}
		key := vaulttypes.KeyFor(r.Id)
		if _, ok := reqMap[key]; ok {
			return fmt.Errorf("duplicate request found for item %s", key)
		}
		reqMap[key] = r
	}

	respMap := map[string]*vaultcommon.SecretResponse{}
	for _, r := range resp.Responses {
		if r.Id == nil {
			return errors.New("GetSecrets response contains nil secret identifier")
		}
		key := vaulttypes.KeyFor(r.Id)
		if _, ok := respMap[key]; ok {
			return fmt.Errorf("duplicate response found for item %s", key)
		}
		respMap[key] = r
	}

	if len(reqMap) != len(respMap) {
		return errors.New("observation doesn't contain matching number of requests and responses")
	}

	for _, rq := range reqMap {
		key := vaulttypes.KeyFor(rq.Id)
		rsp, ok := respMap[key]
		if !ok {
			return fmt.Errorf("missing response for request with id %s", key)
		}

		d := rsp.GetData()
		if d != nil {
			decryptionShares := d.GetEncryptedDecryptionKeyShares()
			if len(rq.EncryptionKeys) != len(d.GetEncryptedDecryptionKeyShares()) {
				return errors.New("observation must contain a share per encryption key provided")
			}

			innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: rq.Id.Owner})
			for _, ds := range decryptionShares {
				if len(ds.Shares) != 1 {
					return errors.New("observation must have exactly 1 share per encryption key")
				}

				share := ds.Shares[0]
				if err := r.cfg.MaxShareLengthBytes.Check(innerCtx, pkgconfig.Size(len(share))*pkgconfig.Byte); err != nil {
					var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
					if errors.As(err, &errBoundLimited) {
						return fmt.Errorf("share provided exceeds maximum size allowed: %w", err)
					}
					return errors.New("failed to check share size")
				}
			}
		}
	}

	return nil
}

func (r *ReportingPlugin) validateCiphertextSize(ctx context.Context, owner string, encryptedValue string) error {
	rawCiphertextB, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return fmt.Errorf("invalid hex encoding for ciphertext: %w", err)
	}
	innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: owner})
	if err := r.cfg.MaxCiphertextLengthBytes.Check(innerCtx, pkgconfig.Size(len(rawCiphertextB))*pkgconfig.Byte); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("ciphertext size exceeds maximum allowed size: %s", errBoundLimited.Limit)
		}
		return errors.New("failed to check ciphertext size")
	}
	return nil
}

func (r *ReportingPlugin) validateCreateSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetCreateSecretsRequest() == nil || o.GetCreateSecretsResponse() == nil {
		return errors.New("CreateSecrets observation must have both request and response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(o.GetCreateSecretsRequest().EncryptedSecrets)); err != nil {
		return err
	}

	if len(o.GetCreateSecretsRequest().EncryptedSecrets) != len(o.GetCreateSecretsResponse().Responses) {
		return errors.New("CreateSecrets request and response must have the same number of items")
	}

	// We disallow duplicate create requests within a single batch request.
	// This prevents users from clobbering their own writes.
	idSet := map[string]bool{}
	for _, s := range o.GetCreateSecretsRequest().EncryptedSecrets {
		if s.Id == nil {
			return errors.New("CreateSecrets request contains nil secret identifier")
		}
		_, ok := idSet[vaulttypes.KeyFor(s.Id)]
		if ok {
			return fmt.Errorf("CreateSecrets requests cannot contain duplicate request for a given secret identifier: %s", s.Id)
		}

		idSet[vaulttypes.KeyFor(s.Id)] = true

		if err := r.validateCiphertextSize(ctx, s.Id.Owner, s.EncryptedValue); err != nil {
			return fmt.Errorf("CreateSecrets request: %w", err)
		}
	}

	for _, r := range o.GetCreateSecretsResponse().Responses {
		if r.Id == nil {
			return errors.New("CreateSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateUpdateSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetUpdateSecretsRequest() == nil || o.GetUpdateSecretsResponse() == nil {
		return errors.New("UpdateSecrets observation must have both request and response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(o.GetUpdateSecretsRequest().EncryptedSecrets)); err != nil {
		return err
	}

	if len(o.GetUpdateSecretsRequest().EncryptedSecrets) != len(o.GetUpdateSecretsResponse().Responses) {
		return errors.New("UpdateSecrets request and response must have the same number of items")
	}

	// We disallow duplicate update requests within a single batch request.
	// This prevents users from clobbering their own writes.
	idSet := map[string]bool{}
	for _, s := range o.GetUpdateSecretsRequest().EncryptedSecrets {
		if s.Id == nil {
			return errors.New("UpdateSecrets request contains nil secret identifier")
		}
		_, ok := idSet[vaulttypes.KeyFor(s.Id)]
		if ok {
			return fmt.Errorf("UpdateSecrets requests cannot contain duplicate request for a given secret identifier: %s", s.Id)
		}

		idSet[vaulttypes.KeyFor(s.Id)] = true

		if err := r.validateCiphertextSize(ctx, s.Id.Owner, s.EncryptedValue); err != nil {
			return fmt.Errorf("UpdateSecrets request: %w", err)
		}
	}

	for _, r := range o.GetUpdateSecretsResponse().Responses {
		if r.Id == nil {
			return errors.New("UpdateSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateDeleteSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetDeleteSecretsRequest() == nil || o.GetDeleteSecretsResponse() == nil {
		return errors.New("DeleteSecrets observation must have both request and response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(o.GetDeleteSecretsRequest().Ids)); err != nil {
		return err
	}

	if len(o.GetDeleteSecretsRequest().Ids) != len(o.GetDeleteSecretsResponse().Responses) {
		return errors.New("DeleteSecrets request and response must have the same number of items")
	}

	// We disallow duplicate delete requests within a single batch request.
	// This prevents users from clobbering their own writes.
	idSet := map[string]bool{}
	for _, r := range o.GetDeleteSecretsRequest().Ids {
		if r == nil {
			return errors.New("DeleteSecrets request contains nil secret identifier")
		}
		_, ok := idSet[vaulttypes.KeyFor(r)]
		if ok {
			return fmt.Errorf("DeleteSecrets requests cannot contain duplicate request for a given secret identifier: %s", r)
		}

		idSet[vaulttypes.KeyFor(r)] = true
	}

	for _, r := range o.GetDeleteSecretsResponse().Responses {
		if r.Id == nil {
			return errors.New("DeleteSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateListSecretIdentifiersObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetListSecretIdentifiersRequest() == nil || o.GetListSecretIdentifiersResponse() == nil {
		return errors.New("ListSecretIdentifiers observation must have both request and response")
	}

	resp := o.GetListSecretIdentifiersResponse()
	if resp.Success {
		ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: o.GetListSecretIdentifiersRequest().Owner})
		if err := r.cfg.MaxSecretsPerOwner.Check(ctx, len(resp.Identifiers)); err != nil {
			var errBoundLimited limits.ErrorBoundLimited[int]
			if errors.As(err, &errBoundLimited) {
				return fmt.Errorf("ListSecretIdentifiers response exceeds maximum number of secrets per owner (have=%d, limit=%d)", len(resp.Identifiers), errBoundLimited.Limit)
			}
			return fmt.Errorf("failed to check max secrets per owner limit: %w", err)
		}
	}

	return nil
}

func (r *ReportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	wrappedStore := NewKVStoreWrapper(NewWriteStore(keyValueReadWriter, r.metrics), r.orgIDAsSecretOwnerEnabled(ctx), r.lggr)

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
				chosen = shaToObs[sha]
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
			os.Outcomes = append(os.Outcomes, o)
		case vaultcommon.RequestType_CREATE_SECRETS:
			req := first.GetCreateSecretsRequest()
			r.stateTransitionCreateSecrets(ctx, wrappedStore.WithRequest(req.OrgId, req.WorkflowOwner), chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		case vaultcommon.RequestType_UPDATE_SECRETS:
			req := first.GetUpdateSecretsRequest()
			r.stateTransitionUpdateSecrets(ctx, wrappedStore.WithRequest(req.OrgId, req.WorkflowOwner), chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		case vaultcommon.RequestType_DELETE_SECRETS:
			req := first.GetDeleteSecretsRequest()
			r.stateTransitionDeleteSecrets(ctx, wrappedStore.WithRequest(req.OrgId, req.WorkflowOwner), chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
			req := first.GetListSecretIdentifiersRequest()
			r.stateTransitionListSecretIdentifiers(ctx, wrappedStore.WithRequest(req.OrgId, req.WorkflowOwner), chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		default:
			r.lggr.Debugw("unknown request type, skipping...", "requestType", first.RequestType, "id", id)
			continue
		}
	}

	// ---
	// Phase 2: Process the pending queue.
	// ---
	err := r.stateTransitionPendingQueue(ctx, wrappedStore, marshalledObs, blobFetcher)
	if err != nil {
		return ocr3_1types.ReportsPlusPrecursor{}, fmt.Errorf("could not process pending queue during state transition: %w", err)
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

func (r *ReportingPlugin) stateTransitionPendingQueue(ctx context.Context, store pendingQueueStore, obs map[uint8]*vaultcommon.Observations, blobFetcher ocr3_1types.BlobFetcher) error {
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

			i := &vaultcommon.StoredPendingQueueItem{}
			err = proto.Unmarshal(blob, i)
			if err != nil {
				r.lggr.Errorw("failed to unmarshal blob into pending queue item", "error", err, "item", pqi)
				continue
			}

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

	// Step 5: Apply batch size and write the latest batch to the store's pending queue.
	if err := r.cfg.MaxBatchSize.Check(ctx, len(keptItems)); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[int]
		if !errors.As(err, &errBoundLimited) {
			return fmt.Errorf("failed to check batch size limit: %w", err)
		}
		keptItems = keptItems[:errBoundLimited.Limit]
	}

	return store.WritePendingQueue(ctx, keptItems)
}

func sortKey(id string, nonce []byte) []byte {
	h := sha256.New()
	h.Write([]byte(id))
	h.Write(nonce)
	return h.Sum(nil)
}

func (r *ReportingPlugin) stateTransitionGetSecrets(ctx context.Context, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	first := chosen[0]
	// First, let's generate the aggregated request.
	// We've validated that all requests with the same sha have the same
	// contents, so we can just sort the SecretRequests by their ID
	// and use that as the aggregated request.
	reqs := first.GetGetSecretsRequest().Requests
	idToReqs := map[string]*vaultcommon.SecretRequest{}
	for _, req := range reqs {
		idToReqs[vaulttypes.KeyFor(req.Id)] = req
	}

	newReqs := []*vaultcommon.SecretRequest{}
	for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
		newReqs = append(newReqs, idToReqs[sreq])
	}

	o.Request = &vaultcommon.Outcome_GetSecretsRequest{
		GetSecretsRequest: &vaultcommon.GetSecretsRequest{
			Requests:      newReqs,
			OrgId:         first.GetGetSecretsRequest().OrgId,
			WorkflowOwner: first.GetGetSecretsRequest().WorkflowOwner,
		},
	}

	// Next, we deal with the responses.
	// For each request, we take the Id of the first observation
	// then aggregate the encrypted shares across all observations.
	// Like with the requests, we sort these by Id and use the result as the response.
	idToAggResponse := map[string]*vaultcommon.SecretResponse{}
	for _, resp := range chosen {
		getSecretsResp := resp.GetGetSecretsResponse()
		for _, rsp := range getSecretsResp.Responses {
			key := vaulttypes.KeyFor(rsp.Id)
			mergedResp, ok := idToAggResponse[key]
			if !ok {
				resp := &vaultcommon.SecretResponse{
					Id:     rsp.Id,
					Result: rsp.Result,
				}
				idToAggResponse[key] = resp
				continue
			}

			if rsp.GetData() != nil {
				data := mergedResp.GetData()

				if len(data.EncryptedDecryptionKeyShares) == 0 {
					data.EncryptedDecryptionKeyShares = []*vaultcommon.EncryptedShares{}
				}

				keyToShares := map[string]*vaultcommon.EncryptedShares{}
				for _, s := range data.EncryptedDecryptionKeyShares {
					keyToShares[s.EncryptionKey] = s
				}

				innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: rsp.Id.Owner})
				for _, existing := range rsp.GetData().EncryptedDecryptionKeyShares {
					if len(existing.Shares) != 1 {
						// This should not happen because we validate against this in ValidateObservation.
						r.lggr.Errorw("exactly 1 share must be provided in the response, skipping", "id", rsp.Id)
						continue
					}
					share := existing.Shares[0]
					if err := r.cfg.MaxShareLengthBytes.Check(innerCtx, pkgconfig.Size(len(share))*pkgconfig.Byte); err != nil {
						var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
						if errors.As(err, &errBoundLimited) {
							r.lggr.Errorw("share exceeds max allowed size, skipping...", "id", rsp.Id, "encryptionKey", existing.EncryptionKey, "err", err)
						} else {
							r.lggr.Errorw("could not check max allowed share size, skipping...", "id", rsp.Id, "encryptionKey", existing.EncryptionKey, "err", err)
						}
						continue
					}

					if shares, ok := keyToShares[existing.EncryptionKey]; ok {
						shares.Shares = append(shares.Shares, share)
					} else {
						// This shouldn't happen -- this is because we're aggregating
						// requests that have a matching sha (excluding the decryption share).
						// Accordingly, we can assume that the request has been made with the same
						// set of encryption keys.
						r.lggr.Errorw("unexpected encryption key in response", "id", rsp.Id, "encryptionKey", existing.EncryptionKey)
					}
				}
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

	newReqs := []*vaultcommon.EncryptedSecret{}
	for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
		newReqs = append(newReqs, idToReqs[sreq])
	}

	o.Request = &vaultcommon.Outcome_CreateSecretsRequest{
		CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
			RequestId:        reqID,
			EncryptedSecrets: newReqs,
			OrgId:            first.GetCreateSecretsRequest().OrgId,
			WorkflowOwner:    first.GetCreateSecretsRequest().WorkflowOwner,
		},
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
		return resp, newUserError(resp.GetError())
	}

	encryptedSecret, err := hex.DecodeString(req.EncryptedValue)
	if err != nil {
		return nil, newUserError("could not decode secret value: invalid hex" + err.Error())
	}

	secret, err := store.GetSecret(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret != nil {
		return nil, newUserError("could not write to key value store: key already exists")
	}

	count, err := store.GetSecretIdentifiersCountForOwner(ctx, req.Id.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret identifiers count for owner: %w", err)
	}

	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: req.Id.Owner})
	if ierr := r.cfg.MaxSecretsPerOwner.Check(ctx, count+1); ierr != nil {
		var errBoundLimited limits.ErrorBoundLimited[int]
		if errors.As(ierr, &errBoundLimited) {
			return nil, newUserError(fmt.Sprintf("could not write to key value store: owner %s has reached maximum number of secrets (limit=%d)", req.Id.Owner, errBoundLimited.Limit))
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

	newReqs := []*vaultcommon.EncryptedSecret{}
	for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
		newReqs = append(newReqs, idToReqs[sreq])
	}

	o.Request = &vaultcommon.Outcome_UpdateSecretsRequest{
		UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
			RequestId:        reqID,
			EncryptedSecrets: newReqs,
			OrgId:            first.GetUpdateSecretsRequest().OrgId,
			WorkflowOwner:    first.GetUpdateSecretsRequest().WorkflowOwner,
		},
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
		return resp, newUserError(resp.GetError())
	}

	encryptedSecret, err := hex.DecodeString(req.EncryptedValue)
	if err != nil {
		return nil, newUserError("could not decode secret value: invalid hex" + err.Error())
	}

	secret, err := store.GetSecret(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret == nil {
		return nil, newUserError("could not write update to key value store: key does not exist")
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

	newReqs := []*vaultcommon.SecretIdentifier{}
	for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
		newReqs = append(newReqs, idToReqs[sreq])
	}

	o.Request = &vaultcommon.Outcome_DeleteSecretsRequest{
		DeleteSecretsRequest: &vaultcommon.DeleteSecretsRequest{
			RequestId:     reqID,
			Ids:           newReqs,
			OrgId:         first.GetDeleteSecretsRequest().OrgId,
			WorkflowOwner: first.GetDeleteSecretsRequest().WorkflowOwner,
		},
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
		return resp, newUserError(resp.GetError())
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

func (r *ReportingPlugin) stateTransitionListSecretIdentifiers(ctx context.Context, store WriteKVStore, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	// All of the logic for the ListSecretIdentifiers request is in the
	// observation phase. This returns the observations in sorted order,
	// so we can just take the first aggregated request and response and
	// use it as the outcome.
	first := chosen[0]
	o.Request = &vaultcommon.Outcome_ListSecretIdentifiersRequest{
		ListSecretIdentifiersRequest: first.GetListSecretIdentifiersRequest(),
	}
	o.Response = &vaultcommon.Outcome_ListSecretIdentifiersResponse{
		ListSecretIdentifiersResponse: first.GetListSecretIdentifiersResponse(),
	}
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

func wrapReportWithKeyBundleInfo(report []byte, reportInfo []byte) (ocr3types.ReportWithInfo[[]byte], error) {
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
		r.cfg.MaxSecretsPerOwner.Close(),
		r.cfg.MaxCiphertextLengthBytes.Close(),
		r.cfg.MaxIdentifierKeyLengthBytes.Close(),
		r.cfg.MaxIdentifierOwnerLengthBytes.Close(),
		r.cfg.MaxIdentifierNamespaceLengthBytes.Close(),
		r.cfg.MaxShareLengthBytes.Close(),
		r.cfg.MaxRequestBatchSize.Close(),
		r.cfg.MaxBatchSize.Close(),
		r.cfg.OrgIDAsSecretOwnerEnabled.Close(),
	)
}
