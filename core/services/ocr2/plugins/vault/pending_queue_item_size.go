package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

const (
	// maxEncryptionKeyStringSlack bounds the wire size of a single encryption public key
	// string taken from GetSecretsRequest (typically hex-encoded curve25519); actual keys are
	// shorter but this keeps estimates conservative without parsing every key format.
	maxEncryptionKeyStringSlack = 128

	// assumedMaxUserFacingErrorStringBytes is added per Create/Update/Delete batched row on the
	// observation path because outcome-side errors are not length-validated in
	// ValidateObservation. There is NO hard guarantee that runtime errors stay within this
	// slack—if they exceed it, the real observation can be larger than this estimate and OCR
	// may reject the wire payload. Tune alongside internal error templates or add explicit
	// validation if this becomes a practical risk.
	assumedMaxUserFacingErrorStringBytes = 8192

	// observationShellExtraBytes approximates protobuf overhead for Observation.id, enum
	// request_type, and oneof tags beyond the nested request/response payloads.
	observationShellExtraBytes = 80

	// getSecretsResponseDataOverheadBytes covers SecretResponse + SecretData field tags/lengths
	// around the ciphertext and share list (excluding raw ciphertext and share payloads).
	getSecretsResponseDataOverheadBytes = 120

	// getSecretsEncryptedSharesPerKeyWireBytes approximates EncryptedShares (encryption_key field
	// + repeated shares tag overhead) for one key on one node's observation.
	getSecretsEncryptedSharesPerKeyWireBytes = 96

	// getSecretsAggregatedShareListProtobufSlackPerShare is a per-share allowance when an outcome
	// aggregates 2F+1 distinct share strings under one encryption key (repeated field growth).
	getSecretsAggregatedShareListProtobufSlackPerShare = 32

	// secretIdentifierProtoSurroundBytes approximates protobuf overhead around the three string
	// fields of SecretIdentifier beyond the raw UTF-8 payload lengths.
	secretIdentifierProtoSurroundBytes = 48
)

func sizeLimiterMaxBytes(ctx context.Context, l limits.BoundLimiter[pkgconfig.Size]) (int, error) {
	v, err := l.Limit(ctx)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func secretIdentifierUpperBoundBytes(ctx context.Context, id *vaultcommon.SecretIdentifier, cfg *ReportingPluginConfig) (int, error) {
	if id == nil {
		return 32, nil
	}
	inner := contexts.WithCRE(ctx, contexts.CRE{Owner: id.Owner})
	k, err := cfg.MaxIdentifierKeyLengthBytes.Limit(inner)
	if err != nil {
		return 0, err
	}
	o, err := cfg.MaxIdentifierOwnerLengthBytes.Limit(inner)
	if err != nil {
		return 0, err
	}
	n, err := cfg.MaxIdentifierNamespaceLengthBytes.Limit(inner)
	if err != nil {
		return 0, err
	}
	ps := proto.Size(id)
	maxWire := int(k + o + n + secretIdentifierProtoSurroundBytes)
	if ps > maxWire {
		return ps, nil
	}
	return maxWire, nil
}

func observationShellBytes(requestID string) int {
	return observationShellExtraBytes + len(requestID)
}

func unmarshalPendingQueueItemPayload(item *vaultcommon.StoredPendingQueueItem) (proto.Message, error) {
	if item == nil || item.Item == nil {
		return nil, errors.New("nil pending queue item or item payload")
	}
	msg, err := item.Item.UnmarshalNew()
	if err != nil {
		return nil, fmt.Errorf("unmarshal pending queue item: %w", err)
	}
	return msg, nil
}

func pendingQueueItemObservationBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, cfg *ReportingPluginConfig, f int) (int, error) {
	_ = f
	msg, err := unmarshalPendingQueueItemPayload(item)
	if err != nil {
		return 0, err
	}
	switch tp := msg.(type) {
	case *vaultcommon.GetSecretsRequest:
		return getSecretsObservationBytes(ctx, item, tp, cfg)
	case *vaultcommon.CreateSecretsRequest:
		return createSecretsObservationBytes(ctx, item, tp)
	case *vaultcommon.UpdateSecretsRequest:
		return updateSecretsObservationBytes(ctx, item, tp)
	case *vaultcommon.DeleteSecretsRequest:
		return deleteSecretsObservationBytes(ctx, item, tp)
	case *vaultcommon.ListSecretIdentifiersRequest:
		return listSecretIdentifiersObservationBytes(ctx, item, tp, cfg)
	default:
		return 0, fmt.Errorf("unknown pending queue item type %T", msg)
	}
}

func pendingQueueItemOutcomeBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, cfg *ReportingPluginConfig, f int) (int, error) {
	msg, err := unmarshalPendingQueueItemPayload(item)
	if err != nil {
		return 0, err
	}
	switch tp := msg.(type) {
	case *vaultcommon.GetSecretsRequest:
		return getSecretsOutcomeBytes(ctx, item, tp, cfg, f)
	case *vaultcommon.CreateSecretsRequest:
		return createSecretsOutcomeBytes(ctx, item, tp)
	case *vaultcommon.UpdateSecretsRequest:
		return updateSecretsOutcomeBytes(ctx, item, tp)
	case *vaultcommon.DeleteSecretsRequest:
		return deleteSecretsOutcomeBytes(ctx, item, tp)
	case *vaultcommon.ListSecretIdentifiersRequest:
		return listSecretIdentifiersOutcomeBytes(ctx, item, tp, cfg)
	default:
		return 0, fmt.Errorf("unknown pending queue item type %T", msg)
	}
}

func getSecretsObservationBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.GetSecretsRequest, cfg *ReportingPluginConfig) (int, error) {
	ctMax, err := sizeLimiterMaxBytes(ctx, cfg.MaxCiphertextLengthBytes)
	if err != nil {
		return 0, err
	}
	// MaxShareLengthBytes is enforced in validateGetSecretsObservation on the wire form of
	// each share string (prefix + base64-encoded anonymous box ciphertext), not on raw TDH2
	// binary share bytes—so this limiter is the correct bound here.
	shareMax, err := sizeLimiterMaxBytes(ctx, cfg.MaxShareLengthBytes)
	if err != nil {
		return 0, err
	}

	total := observationShellBytes(item.Id)
	for _, sr := range req.Requests {
		if sr == nil || sr.Id == nil {
			return 0, errors.New("get secrets request contains nil secret request or identifier")
		}
		keys := len(sr.EncryptionKeys)
		if keys == 0 {
			keys = 1
		}
		idBytes, err := secretIdentifierUpperBoundBytes(ctx, sr.Id, cfg)
		if err != nil {
			return 0, err
		}
		// Observation does not echo the GetSecretsRequest; only per-item SecretResponse (success path worst case).
		total += idBytes + 2*ctMax + getSecretsResponseDataOverheadBytes
		total += keys * (maxEncryptionKeyStringSlack + shareMax + getSecretsEncryptedSharesPerKeyWireBytes)
	}
	return total, nil
}

func getSecretsOutcomeBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.GetSecretsRequest, cfg *ReportingPluginConfig, f int) (int, error) {
	ctMax, err := sizeLimiterMaxBytes(ctx, cfg.MaxCiphertextLengthBytes)
	if err != nil {
		return 0, err
	}
	shareMax, err := sizeLimiterMaxBytes(ctx, cfg.MaxShareLengthBytes)
	if err != nil {
		return 0, err
	}
	byz := 2*f + 1

	// MaxShareLengthBytes applies to each aggregated share string on the wire (same encoding
	// path as single-node observations); see validateGetSecretsObservation.

	total := observationShellBytes(item.Id)
	for _, sr := range req.Requests {
		if sr == nil || sr.Id == nil {
			return 0, errors.New("get secrets request contains nil secret request or identifier")
		}
		keys := len(sr.EncryptionKeys)
		if keys == 0 {
			keys = 1
		}
		idBytes, err := secretIdentifierUpperBoundBytes(ctx, sr.Id, cfg)
		if err != nil {
			return 0, err
		}
		total += idBytes + 2*ctMax + getSecretsResponseDataOverheadBytes
		total += keys * (maxEncryptionKeyStringSlack + shareMax*byz + getSecretsAggregatedShareListProtobufSlackPerShare*byz + getSecretsEncryptedSharesPerKeyWireBytes)
	}
	return total, nil
}

func createSecretsObservationBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.CreateSecretsRequest) (int, error) {
	_ = ctx
	resps := make([]*vaultcommon.CreateSecretResponse, 0, len(req.EncryptedSecrets))
	for _, es := range req.EncryptedSecrets {
		if es.Id == nil {
			return 0, errors.New("create secrets request contains nil identifier")
		}
		resps = append(resps, &vaultcommon.CreateSecretResponse{Id: es.Id, Success: true, Error: ""})
	}
	obs := &vaultcommon.Observation{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_CREATE_SECRETS,
		Request: &vaultcommon.Observation_CreateSecretsRequest{
			CreateSecretsRequest: proto.Clone(req).(*vaultcommon.CreateSecretsRequest),
		},
		Response: &vaultcommon.Observation_CreateSecretsResponse{
			CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{Responses: resps},
		},
	}
	base := proto.Size(obs)
	return base + len(resps)*assumedMaxUserFacingErrorStringBytes, nil
}

func updateSecretsObservationBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.UpdateSecretsRequest) (int, error) {
	_ = ctx
	resps := make([]*vaultcommon.UpdateSecretResponse, 0, len(req.EncryptedSecrets))
	for _, es := range req.EncryptedSecrets {
		if es.Id == nil {
			return 0, errors.New("update secrets request contains nil identifier")
		}
		resps = append(resps, &vaultcommon.UpdateSecretResponse{Id: es.Id, Success: true, Error: ""})
	}
	obs := &vaultcommon.Observation{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
		Request: &vaultcommon.Observation_UpdateSecretsRequest{
			UpdateSecretsRequest: proto.Clone(req).(*vaultcommon.UpdateSecretsRequest),
		},
		Response: &vaultcommon.Observation_UpdateSecretsResponse{
			UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{Responses: resps},
		},
	}
	base := proto.Size(obs)
	return base + len(resps)*assumedMaxUserFacingErrorStringBytes, nil
}

func deleteSecretsObservationBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.DeleteSecretsRequest) (int, error) {
	_ = ctx
	resps := make([]*vaultcommon.DeleteSecretResponse, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id == nil {
			return 0, errors.New("delete secrets request contains nil identifier")
		}
		resps = append(resps, &vaultcommon.DeleteSecretResponse{Id: id, Success: true, Error: ""})
	}
	obs := &vaultcommon.Observation{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_DELETE_SECRETS,
		Request: &vaultcommon.Observation_DeleteSecretsRequest{
			DeleteSecretsRequest: proto.Clone(req).(*vaultcommon.DeleteSecretsRequest),
		},
		Response: &vaultcommon.Observation_DeleteSecretsResponse{
			DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{Responses: resps},
		},
	}
	base := proto.Size(obs)
	return base + len(resps)*assumedMaxUserFacingErrorStringBytes, nil
}

func listSecretIdentifiersObservationBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.ListSecretIdentifiersRequest, cfg *ReportingPluginConfig) (int, error) {
	// MaxSecretsPerOwner is owner-scoped (PerOwner.VaultSecretsLimit); match validateListSecretIdentifiersObservation.
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: req.Owner})
	maxList, err := cfg.MaxSecretsPerOwner.Limit(ctx)
	if err != nil {
		return 0, err
	}
	identifiers := make([]*vaultcommon.SecretIdentifier, maxList)
	for i := range maxList {
		k, err := cfg.MaxIdentifierKeyLengthBytes.Limit(ctx)
		if err != nil {
			return 0, err
		}
		o, err := cfg.MaxIdentifierOwnerLengthBytes.Limit(ctx)
		if err != nil {
			return 0, err
		}
		n, err := cfg.MaxIdentifierNamespaceLengthBytes.Limit(ctx)
		if err != nil {
			return 0, err
		}
		identifiers[i] = &vaultcommon.SecretIdentifier{
			Key:       strings.Repeat("k", int(k)),
			Owner:     strings.Repeat("o", int(o)),
			Namespace: strings.Repeat("n", int(n)),
		}
	}
	resp := &vaultcommon.ListSecretIdentifiersResponse{
		Identifiers: identifiers,
		Success:     true,
		Error:       "",
	}
	obs := &vaultcommon.Observation{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
		Request: &vaultcommon.Observation_ListSecretIdentifiersRequest{
			ListSecretIdentifiersRequest: proto.Clone(req).(*vaultcommon.ListSecretIdentifiersRequest),
		},
		Response: &vaultcommon.Observation_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: resp,
		},
	}
	return proto.Size(obs), nil
}

func createSecretsOutcomeBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.CreateSecretsRequest) (int, error) {
	_ = ctx
	resps := make([]*vaultcommon.CreateSecretResponse, 0, len(req.EncryptedSecrets))
	for _, es := range req.EncryptedSecrets {
		if es.Id == nil {
			return 0, errors.New("create secrets request contains nil identifier")
		}
		resps = append(resps, &vaultcommon.CreateSecretResponse{
			Id:      es.Id,
			Success: false,
			Error:   strings.Repeat("e", assumedMaxUserFacingErrorStringBytes),
		})
	}
	out := &vaultcommon.Outcome{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_CREATE_SECRETS,
		Response: &vaultcommon.Outcome_CreateSecretsResponse{
			CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{Responses: resps},
		},
	}
	return proto.Size(out), nil
}

func updateSecretsOutcomeBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.UpdateSecretsRequest) (int, error) {
	_ = ctx
	resps := make([]*vaultcommon.UpdateSecretResponse, 0, len(req.EncryptedSecrets))
	for _, es := range req.EncryptedSecrets {
		if es.Id == nil {
			return 0, errors.New("update secrets request contains nil identifier")
		}
		resps = append(resps, &vaultcommon.UpdateSecretResponse{
			Id:      es.Id,
			Success: false,
			Error:   strings.Repeat("e", assumedMaxUserFacingErrorStringBytes),
		})
	}
	out := &vaultcommon.Outcome{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_UPDATE_SECRETS,
		Response: &vaultcommon.Outcome_UpdateSecretsResponse{
			UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{Responses: resps},
		},
	}
	return proto.Size(out), nil
}

func deleteSecretsOutcomeBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.DeleteSecretsRequest) (int, error) {
	_ = ctx
	resps := make([]*vaultcommon.DeleteSecretResponse, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id == nil {
			return 0, errors.New("delete secrets request contains nil identifier")
		}
		resps = append(resps, &vaultcommon.DeleteSecretResponse{
			Id:      id,
			Success: false,
			Error:   strings.Repeat("e", assumedMaxUserFacingErrorStringBytes),
		})
	}
	out := &vaultcommon.Outcome{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_DELETE_SECRETS,
		Response: &vaultcommon.Outcome_DeleteSecretsResponse{
			DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{Responses: resps},
		},
	}
	return proto.Size(out), nil
}

func listSecretIdentifiersOutcomeBytes(ctx context.Context, item *vaultcommon.StoredPendingQueueItem, req *vaultcommon.ListSecretIdentifiersRequest, cfg *ReportingPluginConfig) (int, error) {
	// Outcome omits the request payload, but MaxSecretsPerOwner still resolves per owner tenant.
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: req.Owner})
	maxList, err := cfg.MaxSecretsPerOwner.Limit(ctx)
	if err != nil {
		return 0, err
	}
	identifiers := make([]*vaultcommon.SecretIdentifier, maxList)
	for i := range maxList {
		k, err := cfg.MaxIdentifierKeyLengthBytes.Limit(ctx)
		if err != nil {
			return 0, err
		}
		o, err := cfg.MaxIdentifierOwnerLengthBytes.Limit(ctx)
		if err != nil {
			return 0, err
		}
		n, err := cfg.MaxIdentifierNamespaceLengthBytes.Limit(ctx)
		if err != nil {
			return 0, err
		}
		identifiers[i] = &vaultcommon.SecretIdentifier{
			Key:       strings.Repeat("k", int(k)),
			Owner:     strings.Repeat("o", int(o)),
			Namespace: strings.Repeat("n", int(n)),
		}
	}
	resp := &vaultcommon.ListSecretIdentifiersResponse{
		Identifiers: identifiers,
		Success:     true,
		Error:       "",
	}
	out := &vaultcommon.Outcome{
		Id:          item.Id,
		RequestType: vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS,
		Response: &vaultcommon.Outcome_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: resp,
		},
	}
	return proto.Size(out), nil
}
