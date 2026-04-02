package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	syncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
)

type RequestAuthorizer interface {
	AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (isAuthorized bool, owner string, err error)
}
type requestAuthorizer struct {
	allowlistedRequestsSyncer syncerv2.AllowlistedRequestsSyncer
	replayGuard               *DigestReplayGuard
	lggr                      logger.Logger
	sleep                     func(time.Duration)
}

const (
	allowlistReadRetryCount    = 3
	allowlistReadRetryInterval = 3 * time.Second
)

// AuthorizeRequest authorizes a request based on the request digest and the allowlisted requests.
// It does NOT check if the request method is allowed.
func (r *requestAuthorizer) AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (isAuthorized bool, owner string, err error) {
	defer r.replayGuard.ClearExpired()
	r.lggr.Infow("AuthorizeRequest", "method", req.Method, "requestID", req.ID)
	requestDigest, err := req.Digest()
	if err != nil {
		r.lggr.Infow("AuthorizeRequest failed to create digest", "method", req.Method, "requestID", req.ID)
		return false, "", err
	}
	requestDigestBytes, err := hex.DecodeString(requestDigest)
	if err != nil {
		r.lggr.Infow("AuthorizeRequest failed to decode digest", "method", req.Method, "requestID", req.ID)
		return false, "", err
	}
	requestDigestBytes32 := [32]byte(requestDigestBytes)
	if r.allowlistedRequestsSyncer == nil {
		r.lggr.Errorw("AuthorizeRequest allowlistedRequestsSyncer is nil", "method", req.Method, "requestID", req.ID)
		return false, "", errors.New("internal error: allowlistedRequestsSyncer is nil")
	}
	allowlistedRequest, _ := r.fetchAllowlistedItemWithRetry(ctx, req.Method, req.ID, requestDigest, requestDigestBytes32)
	if allowlistedRequest == nil {
		return false, "", errors.New("request not allowlisted")
	}

	if time.Now().UTC().Unix() > int64(allowlistedRequest.ExpiryTimestamp) {
		authorizedRequestStr := string(allowlistedRequest.RequestDigest[:])
		r.lggr.Infow("AuthorizeRequest expired authorization", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr)
		return false, "", errors.New("request authorization expired")
	}

	digestKey := string(allowlistedRequest.RequestDigest[:])
	if err := r.replayGuard.CheckAndRecord(digestKey, int64(allowlistedRequest.ExpiryTimestamp)); err != nil {
		r.lggr.Infow("AuthorizeRequest already authorized previously", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", digestKey)
		return false, "", err
	}

	r.lggr.Infow("AuthorizeRequest success in auth", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", digestKey)
	return true, allowlistedRequest.Owner.Hex(), nil
}

func (r *requestAuthorizer) fetchAllowlistedItemWithRetry(ctx context.Context, method string, requestID interface{}, requestDigest string, digest [32]byte) (*workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, []string) {
	var allowedRequestsStrs []string
	for attempt := 0; attempt <= allowlistReadRetryCount; attempt++ {
		allowedRequests := r.allowlistedRequestsSyncer.GetAllowlistedRequests(ctx)
		allowedRequestsStrs = make([]string, 0, len(allowedRequests))
		for _, rr := range allowedRequests {
			allowedReqStr := fmt.Sprintf("Owner: %s, RequestDigest: %s, ExpiryTimestamp: %d", rr.Owner.Hex(), hex.EncodeToString(rr.RequestDigest[:]), rr.ExpiryTimestamp)
			allowedRequestsStrs = append(allowedRequestsStrs, allowedReqStr)
		}
		r.lggr.Infow("AuthorizeRequest GetAllowlistedRequests", "method", method, "requestID", requestID, "attempt", attempt+1, "allowedRequests", allowedRequestsStrs)

		allowlistedRequest := r.fetchAllowlistedItem(allowedRequests, digest)
		if allowlistedRequest != nil {
			return allowlistedRequest, allowedRequestsStrs
		}

		if attempt == allowlistReadRetryCount {
			break
		}

		r.lggr.Warnw("AuthorizeRequest request not found in allowlist, retrying",
			"method", method,
			"requestID", requestID,
			"digestHexStr", requestDigest,
			"attempt", attempt+1,
			"retryInterval", allowlistReadRetryInterval,
			"allowedRequestsStrs", allowedRequestsStrs)

		select {
		case <-ctx.Done():
			r.lggr.Warnw("AuthorizeRequest allowlist retry canceled",
				"method", method,
				"requestID", requestID,
				"digestHexStr", requestDigest,
				"attempt", attempt+1)
			return nil, allowedRequestsStrs
		default:
		}

		r.sleep(allowlistReadRetryInterval)
	}

	r.lggr.Infow("AuthorizeRequest fetchAllowlistedItem request not allowlisted",
		"method", method,
		"requestID", requestID,
		"digestHexStr", requestDigest,
		"allowedRequestsStrs", allowedRequestsStrs)
	return nil, allowedRequestsStrs
}

func (r *requestAuthorizer) fetchAllowlistedItem(allowListedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, digest [32]byte) *workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest {
	for _, item := range allowListedRequests {
		if item.RequestDigest == digest {
			return &item
		}
	}
	return nil
}

func NewRequestAuthorizer(lggr logger.Logger, allowlistedRequestsSyncer syncerv2.AllowlistedRequestsSyncer) *requestAuthorizer {
	return &requestAuthorizer{
		allowlistedRequestsSyncer: allowlistedRequestsSyncer,
		lggr:                      logger.Named(lggr, "VaultRequestAuthorizer"),
		replayGuard:               NewDigestReplayGuard(),
		sleep:                     time.Sleep,
	}
}
