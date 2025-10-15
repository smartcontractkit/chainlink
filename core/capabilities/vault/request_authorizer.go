package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	workflowsyncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
)

type RequestAuthorizer interface {
	AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (isAuthorized bool, owner string, err error)
}
type requestAuthorizer struct {
	workflowRegistrySyncer    workflowsyncerv2.WorkflowRegistrySyncer
	alreadyAuthorizedRequests map[string]int64
	alreadyAuthorizedMutex    sync.Mutex
	lggr                      logger.Logger
}

// AuthorizeRequest authorizes a request based on the request digest and the allowlisted requests.
// It does NOT check if the request method is allowed.
func (r *requestAuthorizer) AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (isAuthorized bool, owner string, err error) {
	defer r.clearExpiredAuthorizedRequests()
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
	if r.workflowRegistrySyncer == nil {
		r.lggr.Errorw("AuthorizeRequest workflowRegistrySyncer is nil", "method", req.Method, "requestID", req.ID)
		return false, "", errors.New("internal error: workflowRegistrySyncer is nil")
	}
	allowedRequests := r.workflowRegistrySyncer.GetAllowlistedRequests(ctx)
	requestDigests := make([]string, 0, len(allowedRequests))
	for _, allowedRequest := range allowedRequests {
		requestDigests = append(requestDigests, hex.EncodeToString(allowedRequest.RequestDigest[:]))
	}
	r.lggr.Infow("AuthorizeRequest GetAllowlistedRequests", "method", req.Method, "requestID", req.ID, "allowedRequests", allowedRequests, "requestDigestHexStrs", requestDigests)
	allowlistedRequest := r.fetchAllowlistedItem(allowedRequests, requestDigestBytes32)
	if allowlistedRequest == nil {
		r.lggr.Infow("AuthorizeRequest fetchAllowlistedItem request not allowlisted",
			"method", req.Method,
			"requestID", req.ID,
			"digestHexStr", requestDigest,
			"allowedRequestDigestHexStrs", requestDigests)
		return false, "", errors.New("request not allowlisted")
	}
	authorizedRequestStr := string(allowlistedRequest.RequestDigest[:])

	r.alreadyAuthorizedMutex.Lock()
	defer r.alreadyAuthorizedMutex.Unlock()
	if r.alreadyAuthorizedRequests[authorizedRequestStr] > 0 {
		r.lggr.Infow("AuthorizeRequest already authorized previously", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr)
		return false, "", errors.New("request already authorized previously")
	}
	if time.Now().UTC().Unix() > int64(allowlistedRequest.ExpiryTimestamp) {
		r.lggr.Infow("AuthorizeRequest expired authorization", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr)
		return false, "", errors.New("request authorization expired")
	}
	r.lggr.Infow("AuthorizeRequest success in auth", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr)
	r.alreadyAuthorizedRequests[authorizedRequestStr] = int64(allowlistedRequest.ExpiryTimestamp)
	return true, allowlistedRequest.Owner.Hex(), nil
}

func (r *requestAuthorizer) clearExpiredAuthorizedRequests() {
	r.alreadyAuthorizedMutex.Lock()
	defer r.alreadyAuthorizedMutex.Unlock()
	for request, expiry := range r.alreadyAuthorizedRequests {
		if time.Now().UTC().Unix() > expiry {
			delete(r.alreadyAuthorizedRequests, request)
		}
	}
}

func (r *requestAuthorizer) fetchAllowlistedItem(allowListedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, digest [32]byte) *workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest {
	for _, item := range allowListedRequests {
		if item.RequestDigest == digest {
			return &item
		}
	}
	return nil
}

func NewRequestAuthorizer(lggr logger.Logger, workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer) *requestAuthorizer {
	return &requestAuthorizer{
		workflowRegistrySyncer:    workflowRegistrySyncer,
		lggr:                      logger.Named(lggr, "VaultRequestAuthorizer"),
		alreadyAuthorizedRequests: make(map[string]int64),
	}
}
