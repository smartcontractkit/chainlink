package vault

import (
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

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"
	"github.com/smartcontractkit/smdkg/dkgocr"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgrecipientkey"
)

const (
	defaultBatchSize                         = 20
	defaultMaxSecretsPerOwner                = 100
	defaultMaxCiphertextLengthBytes          = 2 * 1024 // 2KB
	defaultMaxIdentifierKeyLengthBytes       = 64
	defaultMaxIdentifierOwnerLengthBytes     = 64
	defaultMaxIdentifierNamespaceLengthBytes = 64

	// The query is empty in this plugin.
	defaultLimitsMaxQueryLength = 100

	// Back of the envelope calculation:
	// - A request can contain 2KB of ciphertext, 192 bytes of metadata (key, owner, namespace),
	// a UUID (16 bytes) plus some overhead = ~2.5KB per request
	// There can be 10 such items in a request, and 20 per batch, so 2.5KB * 10 * 20 = 500KB
	defaultLimitsMaxObservationLength                    = 500 * 1024 // 500KB
	defaultLimitsMaxReportsPlusPrecursorLength           = 500 * 1024 // 500KB
	defaultLimitsMaxReportLength                         = 500 * 1024 // 500KB
	defaultLimitsMaxReportCount                          = 20
	defaultLimitsMaxKeyValueModifiedKeysPlusValuesLength = 1024 * 1024 // 1MB
	defaultLimitsMaxBlobPayloadLength                    = 1024 * 1024 // 1MB
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
	BatchSize                         int
	MaxSecretsPerOwner                int
	MaxCiphertextLengthBytes          int
	MaxIdentifierKeyLengthBytes       int
	MaxIdentifierOwnerLengthBytes     int
	MaxIdentifierNamespaceLengthBytes int
}

func NewReportingPluginFactory(
	lggr logger.Logger,
	store *requests.Store[*vaulttypes.Request],
	db dkgocrtypes.ResultPackageDatabase,
	recipientKey *dkgrecipientkey.Key,
	lazyPublicKey *vaultcap.LazyPublicKey,
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
		lggr:         lggr.Named("VaultReportingPluginFactory"),
		store:        store,
		cfg:          cfg,
		db:           db,
		recipientKey: recipientKey,
	}, nil
}

type ReportingPluginFactory struct {
	lggr         logger.Logger
	store        *requests.Store[*vaulttypes.Request]
	cfg          *ReportingPluginConfig
	db           dkgocrtypes.ResultPackageDatabase
	recipientKey *dkgrecipientkey.Key
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

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	var configProto vaultcommon.ReportingPluginConfig
	if err := proto.Unmarshal(config.OffchainConfig, &configProto); err != nil {
		return nil, ocr3_1types.ReportingPluginInfo{}, fmt.Errorf("could not unmarshal reporting plugin config: %w", err)
	}

	if configProto.BatchSize == 0 {
		configProto.BatchSize = defaultBatchSize
	}

	if configProto.MaxSecretsPerOwner == 0 {
		configProto.MaxSecretsPerOwner = defaultMaxSecretsPerOwner
	}

	if configProto.MaxCiphertextLengthBytes == 0 {
		configProto.MaxCiphertextLengthBytes = defaultMaxCiphertextLengthBytes
	}

	if configProto.MaxIdentifierKeyLengthBytes == 0 {
		configProto.MaxIdentifierKeyLengthBytes = defaultMaxIdentifierKeyLengthBytes
	}

	if configProto.MaxIdentifierOwnerLengthBytes == 0 {
		configProto.MaxIdentifierOwnerLengthBytes = defaultMaxIdentifierOwnerLengthBytes
	}

	if configProto.MaxIdentifierNamespaceLengthBytes == 0 {
		configProto.MaxIdentifierNamespaceLengthBytes = defaultMaxIdentifierNamespaceLengthBytes
	}

	if configProto.LimitsMaxQueryLength == 0 {
		configProto.LimitsMaxQueryLength = defaultLimitsMaxQueryLength
	}

	if configProto.LimitsMaxObservationLength == 0 {
		configProto.LimitsMaxObservationLength = defaultLimitsMaxObservationLength
	}

	if configProto.LimitsMaxReportsPlusPrecursorLength == 0 {
		configProto.LimitsMaxReportsPlusPrecursorLength = defaultLimitsMaxReportsPlusPrecursorLength
	}

	if configProto.LimitsMaxReportLength == 0 {
		configProto.LimitsMaxReportLength = defaultLimitsMaxReportLength
	}

	if configProto.LimitsMaxReportCount == 0 {
		configProto.LimitsMaxReportCount = defaultLimitsMaxReportCount
	}

	if configProto.LimitsMaxKeyValueModifiedKeysPlusValuesLength == 0 {
		configProto.LimitsMaxKeyValueModifiedKeysPlusValuesLength = defaultLimitsMaxKeyValueModifiedKeysPlusValuesLength
	}

	if configProto.LimitsMaxBlobPayloadLength == 0 {
		configProto.LimitsMaxBlobPayloadLength = defaultLimitsMaxBlobPayloadLength
	}

	if configProto.DKGInstanceID == nil {
		return nil, ocr3_1types.ReportingPluginInfo{}, errors.New("DKG instance ID cannot be nil")
	}

	publicKey, privateKeyShare, err := r.getKeyMaterial(ctx, *configProto.DKGInstanceID)
	if err != nil {
		return nil, ocr3_1types.ReportingPluginInfo{}, fmt.Errorf("could not get key material from DB: %w", err)
	}

	r.cfg.LazyPublicKey.Set(publicKey)

	cfg := &ReportingPluginConfig{
		PublicKey:                         publicKey,
		PrivateKeyShare:                   privateKeyShare,
		BatchSize:                         int(configProto.BatchSize),
		MaxSecretsPerOwner:                int(configProto.MaxSecretsPerOwner),
		MaxCiphertextLengthBytes:          int(configProto.MaxCiphertextLengthBytes),
		MaxIdentifierKeyLengthBytes:       int(configProto.MaxIdentifierKeyLengthBytes),
		MaxIdentifierOwnerLengthBytes:     int(configProto.MaxIdentifierOwnerLengthBytes),
		MaxIdentifierNamespaceLengthBytes: int(configProto.MaxIdentifierNamespaceLengthBytes),
	}
	return &ReportingPlugin{
			lggr:       r.lggr.Named("VaultReportingPlugin"),
			store:      r.store,
			cfg:        cfg,
			onchainCfg: config,
		}, ocr3_1types.ReportingPluginInfo{
			Name: "VaultReportingPlugin",
			Limits: ocr3_1types.ReportingPluginLimits{
				MaxQueryLength:                          int(configProto.LimitsMaxQueryLength),
				MaxObservationLength:                    int(configProto.LimitsMaxObservationLength),
				MaxReportsPlusPrecursorLength:           int(configProto.LimitsMaxReportsPlusPrecursorLength),
				MaxReportLength:                         int(configProto.LimitsMaxReportLength),
				MaxReportCount:                          int(configProto.LimitsMaxReportCount),
				MaxKeyValueModifiedKeysPlusValuesLength: int(configProto.LimitsMaxKeyValueModifiedKeysPlusValuesLength),
				MaxBlobPayloadLength:                    int(configProto.LimitsMaxBlobPayloadLength),
			},
		}, nil
}

type ReportingPlugin struct {
	lggr       logger.Logger
	store      *requests.Store[*vaulttypes.Request]
	onchainCfg ocr3types.ReportingPluginConfig
	cfg        *ReportingPluginConfig
}

func (r *ReportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return types.Query{}, nil
}

func (r *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	// Note: this could mean that we end up processing more than `batchSize` requests
	// in the aggregate, since all nodes will fetch `batchSize` requests and they aren't
	// guaranteed to fetch the same requests.
	batch, err := r.store.FirstN(r.cfg.BatchSize)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch of requests: %w", err)
	}
	// Avoid log spam by only logging if we have any requests to process.
	if len(batch) > 0 {
		r.lggr.Debugw("observation started", "seqNr", seqNr, "batchSize", r.cfg.BatchSize)
	}

	ids := []string{}
	obs := []*vaultcommon.Observation{}
	for _, req := range batch {
		o := &vaultcommon.Observation{
			Id: req.ID(),
		}
		ids = append(ids, req.ID())

		switch req.Payload.(type) {
		case *vaultcommon.GetSecretsRequest:
			r.observeGetSecrets(ctx, NewReadStore(keyValueReader), req.Payload, o)
		case *vaultcommon.CreateSecretsRequest:
			r.observeCreateSecrets(ctx, NewReadStore(keyValueReader), req.Payload, o)
		case *vaultcommon.UpdateSecretsRequest:
			r.observeUpdateSecrets(ctx, NewReadStore(keyValueReader), req.Payload, o)
		case *vaultcommon.DeleteSecretsRequest:
			r.observeDeleteSecrets(ctx, NewReadStore(keyValueReader), req.Payload, o)
		case *vaultcommon.ListSecretIdentifiersRequest:
			r.observeListSecretIdentifiers(ctx, NewReadStore(keyValueReader), req.Payload, o)
		default:
			r.lggr.Errorw("unknown request type, skipping...", "requestType", fmt.Sprintf("%T", req.Payload), "id", req.ID())
			continue
		}

		obs = append(obs, o)
	}

	obsb, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vaultcommon.Observations{
		Observations: obs,
	})
	if err != nil {
		return nil, fmt.Errorf("could not marshal observations: %w", err)
	}

	// Avoid log spam by only logging if we have any requests to process.
	if len(batch) > 0 {
		r.lggr.Debugw("observation complete", "ids", ids, "batchSize", len(batch))
	}
	return types.Observation(obsb), nil
}

func (r *ReportingPlugin) observeGetSecrets(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.GetSecretsRequest)
	o.RequestType = vaultcommon.RequestType_GET_SECRETS
	o.Request = &vaultcommon.Observation_GetSecretsRequest{
		GetSecretsRequest: tp,
	}
	resps := []*vaultcommon.SecretResponse{}
	for _, secretRequest := range tp.Requests {
		resp, ierr := r.observeGetSecretsRequest(ctx, reader, secretRequest)
		if ierr != nil {
			r.lggr.Errorw("failed to observe get secret request item", "id", secretRequest.Id, "error", ierr)
			errorMsg := "failed to handle get secret request"
			if errors.Is(ierr, &userError{}) {
				errorMsg = ierr.Error()
			}
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

func (r *ReportingPlugin) observeGetSecretsRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.SecretRequest) (*vaultcommon.SecretResponse, error) {
	id, err := r.validateSecretIdentifier(secretRequest.Id)
	if err != nil {
		return nil, err
	}

	secret, err := reader.GetSecret(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret == nil {
		return nil, newUserError("key does not exist")
	}

	ct := &tdh2easy.Ciphertext{}
	err = ct.UnmarshalVerify(secret.EncryptedSecret, r.cfg.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ciphertext: %w", err)
	}

	share, err := tdh2easy.Decrypt(ct, r.cfg.PrivateKeyShare)
	if err != nil {
		return nil, fmt.Errorf("could not generate decryption share: %w", err)
	}

	shareb, err := share.Marshal()
	if err != nil {
		return nil, errors.New("could not marshal decryption share")
	}

	shares := []*vaultcommon.EncryptedShares{}
	for _, pk := range secretRequest.EncryptionKeys {
		publicKey, err := hex.DecodeString(pk)
		if err != nil {
			return nil, newUserError("failed to convert public key to bytes: " + err.Error())
		}

		if len(publicKey) != curve25519.PointSize {
			return nil, newUserError(fmt.Sprintf("invalid public key size: expected %d bytes, got %d bytes", curve25519.PointSize, len(publicKey)))
		}

		publicKeyLength := [curve25519.PointSize]byte(publicKey)
		encrypted, err := box.SealAnonymous(nil, shareb, &publicKeyLength, rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt decryption share: %w", err)
		}

		shares = append(shares, &vaultcommon.EncryptedShares{
			EncryptionKey: pk,
			Shares: []string{
				hex.EncodeToString(encrypted),
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
		validatedID, ierr := r.observeCreateSecretRequest(ctx, reader, sr, requestsCountForID)
		if ierr != nil {
			l.Errorw("failed to handle create secret request item", "id", sr.Id, "error", ierr)
			errorMsg := "failed to handle create secret request"
			if errors.Is(ierr, &userError{}) {
				errorMsg = ierr.Error()
			}
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

func (r *ReportingPlugin) observeCreateSecretRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.EncryptedSecret, requestsCountForID map[string]int) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateSecretIdentifier(secretRequest.Id)
	if err != nil {
		return id, err
	}

	if requestsCountForID[vaulttypes.KeyFor(secretRequest.Id)] > 1 {
		return id, newUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}

	rawCiphertext := secretRequest.EncryptedValue
	rawCiphertextB, err := hex.DecodeString(rawCiphertext)
	if err != nil {
		return id, newUserError("invalid hex encoding for ciphertext: " + err.Error())
	}

	if len(rawCiphertextB) > r.cfg.MaxCiphertextLengthBytes {
		return id, newUserError(fmt.Sprintf("ciphertext size exceeds maximum allowed size: %d bytes", r.cfg.MaxCiphertextLengthBytes))
	}

	ct := &tdh2easy.Ciphertext{}
	err = ct.UnmarshalVerify(rawCiphertextB, r.cfg.PublicKey)
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
		validatedID, ierr := r.observeUpdateSecretRequest(ctx, reader, sr, requestsCountForID)
		if ierr != nil {
			l.Errorw("failed to observe update secret request item", "id", sr.Id, "error", ierr)
			errorMsg := "failed to handle update secret request"
			if errors.Is(ierr, &userError{}) {
				errorMsg = ierr.Error()
			}
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

func (r *ReportingPlugin) observeUpdateSecretRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.EncryptedSecret, requestsCountForID map[string]int) (*vaultcommon.SecretIdentifier, error) {
	// The checks at this stage are identical since we only check the correctness of the payload
	// at this stage. Checks that are different between update and create, like whether the secret already exists,
	// are handled in the StateTransition phase.
	return r.observeCreateSecretRequest(ctx, reader, secretRequest, requestsCountForID)
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

	md, err := reader.GetMetadata(req.Owner)
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
			l.Errorw("failed to handle delete secret request item", "id", id, "error", ierr)
			errorMsg := "failed to handle delete secret request"
			if errors.Is(ierr, &userError{}) {
				errorMsg = ierr.Error()
			}
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
	id, err := r.validateSecretIdentifier(identifier)
	if err != nil {
		return id, err
	}

	if requestsCountForID[vaulttypes.KeyFor(identifier)] > 1 {
		return id, newUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}

	ss, err := reader.GetSecret(id)
	if err != nil {
		return id, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if ss == nil {
		return id, newUserError("key does not exist")
	}

	return id, nil
}

func (r *ReportingPlugin) validateSecretIdentifier(id *vaultcommon.SecretIdentifier) (*vaultcommon.SecretIdentifier, error) {
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

	if len(id.Owner) > r.cfg.MaxIdentifierOwnerLengthBytes {
		return nil, newUserError(fmt.Sprintf("invalid secret identifier: owner exceeds maximum length of %d bytes", r.cfg.MaxIdentifierOwnerLengthBytes))
	}

	if len(id.Namespace) > r.cfg.MaxIdentifierNamespaceLengthBytes {
		return nil, newUserError(fmt.Sprintf("invalid secret identifier: namespace exceeds maximum length of %d bytes", r.cfg.MaxIdentifierNamespaceLengthBytes))
	}

	if len(id.Key) > r.cfg.MaxIdentifierKeyLengthBytes {
		return nil, newUserError(fmt.Sprintf("invalid secret identifier: key exceeds maximum length of %d bytes", r.cfg.MaxIdentifierKeyLengthBytes))
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

func (r *ReportingPlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, keyValueReader ocr3_1types.KeyValueReader, blobFetcher ocr3_1types.BlobFetcher) error {
	obs := &vaultcommon.Observations{}
	if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
		return errors.New("failed to unmarshal observations: " + err.Error())
	}

	seen := map[string]bool{}
	for _, o := range obs.Observations {
		err := validateObservation(o)
		if err != nil {
			return errors.New("invalid observation: " + err.Error())
		}

		_, ok := seen[o.Id]
		if ok {
			return errors.New("invalid observation: a single observation cannot contain duplicate observations for the same request id")
		}

		seen[o.Id] = true
	}

	return nil
}

func (r *ReportingPlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReader ocr3_1types.KeyValueReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, r.onchainCfg.N, r.onchainCfg.F, aos), nil
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

func validateObservation(o *vaultcommon.Observation) error {
	if o.Id == "" {
		return errors.New("observation id cannot be empty")
	}

	switch o.RequestType {
	case vaultcommon.RequestType_GET_SECRETS:
		if o.GetGetSecretsRequest() == nil || o.GetGetSecretsResponse() == nil {
			return errors.New("GetSecrets observation must have both request and response")
		}

		if len(o.GetGetSecretsRequest().Requests) != len(o.GetGetSecretsResponse().Responses) {
			return errors.New("GetSecrets request and response must have the same number of items")
		}
	case vaultcommon.RequestType_CREATE_SECRETS:
		if o.GetCreateSecretsRequest() == nil || o.GetCreateSecretsResponse() == nil {
			return errors.New("CreateSecrets observation must have both request and response")
		}

		if len(o.GetCreateSecretsRequest().EncryptedSecrets) != len(o.GetCreateSecretsResponse().Responses) {
			return errors.New("CreateSecrets request and response must have the same number of items")
		}

		// We disallow duplicate create requests within a single batch request.
		// This prevents users from clobbering their own writes.
		idSet := map[string]bool{}
		for _, r := range o.GetCreateSecretsRequest().EncryptedSecrets {
			_, ok := idSet[vaulttypes.KeyFor(r.Id)]
			if ok {
				return fmt.Errorf("CreateSecrets requests cannot contain duplicate request for a given secret identifier: %s", r.Id)
			}

			idSet[vaulttypes.KeyFor(r.Id)] = true
		}
	case vaultcommon.RequestType_UPDATE_SECRETS:
		if o.GetUpdateSecretsRequest() == nil || o.GetUpdateSecretsResponse() == nil {
			return errors.New("UpdateSecrets observation must have both request and response")
		}

		if len(o.GetUpdateSecretsRequest().EncryptedSecrets) != len(o.GetUpdateSecretsResponse().Responses) {
			return errors.New("UpdateSecrets request and response must have the same number of items")
		}

		// We disallow duplicate update requests within a single batch request.
		// This prevents users from clobbering their own writes.
		idSet := map[string]bool{}
		for _, r := range o.GetUpdateSecretsRequest().EncryptedSecrets {
			_, ok := idSet[vaulttypes.KeyFor(r.Id)]
			if ok {
				return fmt.Errorf("UpdateSecrets requests cannot contain duplicate request for a given secret identifier: %s", r.Id)
			}

			idSet[vaulttypes.KeyFor(r.Id)] = true
		}
	case vaultcommon.RequestType_DELETE_SECRETS:
		if o.GetDeleteSecretsRequest() == nil || o.GetDeleteSecretsResponse() == nil {
			return errors.New("DeleteSecrets observation must have both request and response")
		}

		if len(o.GetDeleteSecretsRequest().Ids) != len(o.GetDeleteSecretsResponse().Responses) {
			return errors.New("DeleteSecrets request and response must have the same number of items")
		}

		// We disallow duplicate delete requests within a single batch request.
		// This prevents users from clobbering their own writes.
		idSet := map[string]bool{}
		for _, r := range o.GetDeleteSecretsRequest().Ids {
			_, ok := idSet[vaulttypes.KeyFor(r)]
			if ok {
				return fmt.Errorf("DeleteSecrets requests cannot contain duplicate request for a given secret identifier: %s", r)
			}

			idSet[vaulttypes.KeyFor(r)] = true
		}
	case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
		if o.GetListSecretIdentifiersRequest() == nil || o.GetListSecretIdentifiersResponse() == nil {
			return errors.New("ListSecretIdentifiers observation must have both request and response")
		}
	default:
		return errors.New("invalid observation type: " + o.RequestType.String())
	}

	return nil
}

func (r *ReportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	store := NewWriteStore(keyValueReadWriter)

	obsMap := map[string][]*vaultcommon.Observation{}
	for _, ao := range aos {
		obs := &vaultcommon.Observations{}
		if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
			// Note: this shouldn't happen as all observations are validated in ValidateObservation.
			r.lggr.Errorw("failed to unmarshal observations", "error", err, "observation", ao.Observation)
			continue
		}

		for _, o := range obs.Observations {
			if _, ok := obsMap[o.Id]; !ok {
				obsMap[o.Id] = []*vaultcommon.Observation{}
			}
			obsMap[o.Id] = append(obsMap[o.Id], o)
		}
	}

	os := &vaultcommon.Outcomes{
		Outcomes: []*vaultcommon.Outcome{},
	}
	for id, obs := range obsMap {
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
		threshold := 2*r.onchainCfg.F + 1
		for sha, obs := range shaToObs {
			if len(obs) >= threshold {
				r.lggr.Debugw("sufficient observations for sha", "sha", sha, "count", len(obs), "threshold", threshold, "id", id)
				chosen = shaToObs[sha]
				break
			}
		}

		if len(chosen) == 0 {
			r.lggr.Warnw("insufficient observations found for id", "id", id, "threshold", threshold)
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
			r.stateTransitionCreateSecrets(ctx, store, chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		case vaultcommon.RequestType_UPDATE_SECRETS:
			r.stateTransitionUpdateSecrets(ctx, store, chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		case vaultcommon.RequestType_DELETE_SECRETS:
			r.stateTransitionDeleteSecrets(ctx, store, chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
			r.stateTransitionListSecretIdentifiers(ctx, store, chosen, o)
			os.Outcomes = append(os.Outcomes, o)
		default:
			r.lggr.Debugw("unknown request type, skipping...", "requestType", first.RequestType, "id", id)
			continue
		}
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
			Requests: newReqs,
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

				for _, existing := range rsp.GetData().EncryptedDecryptionKeyShares {
					if shares, ok := keyToShares[existing.EncryptionKey]; ok {
						shares.Shares = append(shares.Shares, existing.Shares...)
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
			r.lggr.Errorw("failed to handle create secret request", "id", req.Id, "requestID", reqID, "error", err)
			errorMsg := "failed to handle create secret request"
			if errors.Is(err, &userError{}) {
				errorMsg = err.Error()
			}
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

	secret, err := store.GetSecret(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret != nil {
		return nil, newUserError("could not write to key value store: key already exists")
	}

	count, err := store.GetSecretIdentifiersCountForOwner(req.Id.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret identifiers count for owner: %w", err)
	}

	if count+1 > r.cfg.MaxSecretsPerOwner {
		return nil, newUserError(fmt.Sprintf("could not write to key value store: owner %s has reached maximum number of secrets (%d)", req.Id.Owner, r.cfg.MaxSecretsPerOwner))
	}

	err = store.WriteSecret(req.Id, &vaultcommon.StoredSecret{
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
			r.lggr.Errorw("failed to handle update secret request", "id", req.Id, "requestID", reqID, "error", err)
			errorMsg := "failed to handle update secret request"
			if errors.Is(err, &userError{}) {
				errorMsg = err.Error()
			}
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

	secret, err := store.GetSecret(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret == nil {
		return nil, newUserError("could not write update to key value store: key does not exist")
	}

	err = store.WriteSecret(req.Id, &vaultcommon.StoredSecret{
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
			RequestId: reqID,
			Ids:       newReqs,
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
			r.lggr.Errorw("failed to handle delete secret request", "id", id, "requestId", reqID, "error", err)
			errorMsg := "failed to handle delete secret request"
			if errors.Is(err, &userError{}) {
				errorMsg = err.Error()
			}
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

	err := store.DeleteSecret(id)
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

func (r *ReportingPlugin) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueReader) error {
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
	return nil
}
