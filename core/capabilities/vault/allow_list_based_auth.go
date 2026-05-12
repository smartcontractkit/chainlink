package vault

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	workflowsyncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
)

const (
	// The workflow registry syncer polls every 12s by default. Keep the
	// retry window comfortably above that so newly allowlisted requests
	// can propagate to every node before auth gives up.
	allowListBasedAuthRetryCount    = 10
	allowListBasedAuthRetryInterval = 3 * time.Second
)

type allowListBasedAuth struct {
	workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer
	lggr                   logger.Logger
	retryCount             int
	retryInterval          time.Duration
	// vaultBase64EncodingEnabled, when non-nil and AllowErr(ctx)==nil, indicates the
	// gateway forwards ciphertext as base64. Allowlist digests are always registered
	// against the client (hex) wire form, so we normalize before Digest().
	vaultBase64EncodingEnabled limits.GateLimiter
}

// AuthorizeRequest authorizes a request using AllowListBasedAuth.
// It does NOT check if the request method is allowed.
func (r *allowListBasedAuth) AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	r.lggr.Debugw("AllowListBasedAuth authorizing request", "method", req.Method, "requestID", req.ID)
	reqForDigest := req
	if r.vaultBase64EncodingEnabled != nil && r.vaultBase64EncodingEnabled.AllowErr(ctx) == nil {
		normalized, normErr := vaultRequestEncryptedSecretsBase64ToHexForAllowlistDigest(req)
		if normErr != nil {
			r.lggr.Debugw("AllowListBasedAuth failed to normalize ciphertext for digest", "method", req.Method, "requestID", req.ID, "error", normErr)
			return nil, normErr
		}
		reqForDigest = normalized
	}
	requestDigest, err := reqForDigest.Digest()
	if err != nil {
		r.lggr.Debugw("AllowListBasedAuth failed to create digest", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, err
	}
	requestDigestBytes, err := hex.DecodeString(requestDigest)
	if err != nil {
		r.lggr.Debugw("AllowListBasedAuth failed to decode digest", "method", req.Method, "requestID", req.ID, "requestDigest", requestDigest, "error", err)
		return nil, err
	}
	requestDigestBytes32 := [32]byte(requestDigestBytes)
	if r.workflowRegistrySyncer == nil {
		r.lggr.Errorw("AllowListBasedAuth workflowRegistrySyncer is nil", "method", req.Method, "requestID", req.ID)
		return nil, errors.New("internal error: workflowRegistrySyncer is nil")
	}
	allowlistedRequest, allowedRequestsStrs, err := r.findAllowlistedItemWithRetry(ctx, req, requestDigest, requestDigestBytes32)
	if err != nil {
		return nil, err
	}
	if allowlistedRequest == nil {
		r.lggr.Debugw("AllowListBasedAuth request digest not allowlisted",
			"method", req.Method,
			"requestID", req.ID,
			"digestHexStr", requestDigest,
			"allowedRequestsStrs", allowedRequestsStrs)
		return nil, errors.New("request not allowlisted")
	}

	if time.Now().UTC().Unix() > int64(allowlistedRequest.ExpiryTimestamp) {
		authorizedRequestStr := string(allowlistedRequest.RequestDigest[:])
		r.lggr.Debugw("AllowListBasedAuth authorization expired", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr, "expiryTimestamp", allowlistedRequest.ExpiryTimestamp)
		return nil, errors.New("request authorization expired")
	}

	digestKey := string(allowlistedRequest.RequestDigest[:])
	r.lggr.Debugw("AllowListBasedAuth authorization succeeded", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", digestKey, "owner", allowlistedRequest.Owner.Hex(), "expiryTimestamp", allowlistedRequest.ExpiryTimestamp)
	return &AuthResult{
		workflowOwner: allowlistedRequest.Owner.Hex(),
		digest:        digestKey,
		expiresAt:     int64(allowlistedRequest.ExpiryTimestamp),
	}, nil
}

func (r *allowListBasedAuth) findAllowlistedItemWithRetry(ctx context.Context, req jsonrpc.Request[json.RawMessage], requestDigest string, requestDigestBytes32 [32]byte) (*workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, []string, error) {
	for attempt := 0; attempt <= r.retryCount; attempt++ {
		allowedRequests := r.workflowRegistrySyncer.GetAllowlistedRequests(ctx)
		allowedRequestsStrs := make([]string, 0, len(allowedRequests))
		for _, rr := range allowedRequests {
			allowedReqStr := fmt.Sprintf("AuthorizedOwner: %s, RequestDigest: %s, ExpiryTimestamp: %d", rr.Owner.Hex(), hex.EncodeToString(rr.RequestDigest[:]), rr.ExpiryTimestamp)
			allowedRequestsStrs = append(allowedRequestsStrs, allowedReqStr)
		}
		r.lggr.Debugw("AllowListBasedAuth loaded allowlisted requests", "method", req.Method, "requestID", req.ID, "attempt", attempt+1, "allowedRequests", allowedRequestsStrs)

		allowlistedRequest := r.fetchAllowlistedItem(allowedRequests, requestDigestBytes32)
		if allowlistedRequest != nil {
			return allowlistedRequest, allowedRequestsStrs, nil
		}
		if attempt == r.retryCount {
			return nil, allowedRequestsStrs, nil
		}

		r.lggr.Debugw("AllowListBasedAuth request digest not yet allowlisted, retrying",
			"method", req.Method,
			"requestID", req.ID,
			"digestHexStr", requestDigest,
			"attempt", attempt+1,
			"maxAttempts", r.retryCount+1,
			"retryInterval", r.retryInterval)
		if err := sleepWithContext(ctx, r.retryInterval); err != nil {
			r.lggr.Debugw("AllowListBasedAuth retry canceled", "method", req.Method, "requestID", req.ID, "error", err)
			return nil, nil, err
		}
	}

	return nil, nil, nil // unreachable: loop always returns
}

func (r *allowListBasedAuth) fetchAllowlistedItem(allowListedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, digest [32]byte) *workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest {
	for _, item := range allowListedRequests {
		if item.RequestDigest == digest {
			return &item
		}
	}
	return nil
}

// NewAllowListBasedAuth creates the allowlist-backed Vault auth mechanism.
// vaultBase64EncodingEnabled may be nil; when set and limits allow, ciphertext in
// create/update params is decoded from base64 to hex before computing the allowlist digest.
func NewAllowListBasedAuth(lggr logger.Logger, workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer, vaultBase64EncodingEnabled limits.GateLimiter) *allowListBasedAuth {
	return &allowListBasedAuth{
		workflowRegistrySyncer:     workflowRegistrySyncer,
		lggr:                       logger.Named(lggr, "VaultAllowListBasedAuth"),
		retryCount:                 allowListBasedAuthRetryCount,
		retryInterval:              allowListBasedAuthRetryInterval,
		vaultBase64EncodingEnabled: vaultBase64EncodingEnabled,
	}
}

func vaultRequestEncryptedSecretsBase64ToHexForAllowlistDigest(req jsonrpc.Request[json.RawMessage]) (jsonrpc.Request[json.RawMessage], error) {
	if req.Params == nil {
		return req, nil
	}
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		parsed := &vaultcommon.CreateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return jsonrpc.Request[json.RawMessage]{}, err
		}
		if err := convertEncryptedSecretsBase64WireToHexStrings(parsed.EncryptedSecrets); err != nil {
			return jsonrpc.Request[json.RawMessage]{}, err
		}
		b, err := json.Marshal(parsed)
		if err != nil {
			return jsonrpc.Request[json.RawMessage]{}, err
		}
		raw := json.RawMessage(b)
		out := req
		out.Params = &raw
		return out, nil
	case vaulttypes.MethodSecretsUpdate:
		parsed := &vaultcommon.UpdateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return jsonrpc.Request[json.RawMessage]{}, err
		}
		if err := convertEncryptedSecretsBase64WireToHexStrings(parsed.EncryptedSecrets); err != nil {
			return jsonrpc.Request[json.RawMessage]{}, err
		}
		b, err := json.Marshal(parsed)
		if err != nil {
			return jsonrpc.Request[json.RawMessage]{}, err
		}
		raw := json.RawMessage(b)
		out := req
		out.Params = &raw
		return out, nil
	default:
		return req, nil
	}
}

func convertEncryptedSecretsBase64WireToHexStrings(secrets []*vaultcommon.EncryptedSecret) error {
	for i, item := range secrets {
		if item == nil || item.EncryptedValue == "" {
			continue
		}
		raw, err := base64.StdEncoding.DecodeString(item.EncryptedValue)
		if err != nil {
			return fmt.Errorf("encrypted secret index %d: %w", i, err)
		}
		item.EncryptedValue = hex.EncodeToString(raw)
	}
	return nil
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
