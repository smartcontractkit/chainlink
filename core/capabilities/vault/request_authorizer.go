package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
	workflowsyncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
)

type RequestAuthorizer interface {
	AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (isAuthorized bool, owner string, err error)
}
type requestAuthorizer struct {
	workflowRegistrySyncer    workflowsyncerv2.WorkflowRegistrySyncer
	alreadyAuthorizedRequests map[string]bool
	alreadyAuthorizedMutex    sync.Mutex
	lggr                      logger.Logger
}

func (r *requestAuthorizer) AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (isAuthorized bool, owner string, err error) {
	defer r.clearExpiredAuthorizedRequests()
	r.lggr.Infow("AuthorizeRequest", "method", req.Method, "requestID", req.ID)
	digest, err := vaultutils.DigestForRequest(req)
	if err != nil {
		r.lggr.Infow("AuthorizeRequest failed to create digest", "method", req.Method, "requestID", req.ID)
		return false, "", err
	}
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
	allowlistedRequest := r.fetchAllowlistedItem(allowedRequests, digest)
	if allowlistedRequest == nil {
		r.lggr.Infow("AuthorizeRequest fetchAllowlistedItem request not allowlisted", "method", req.Method, "requestID", req.ID, "digestHexStr", hex.EncodeToString(digest[:]), "allowedRequestDigestHexStrs", requestDigests)
		return false, "", errors.New("request not allowlisted")
	}
	authorizedRequestStr := string(allowlistedRequest.RequestDigest[:]) + "-->" + strconv.FormatUint(uint64(allowlistedRequest.ExpiryTimestamp), 10)

	r.alreadyAuthorizedMutex.Lock()
	defer r.alreadyAuthorizedMutex.Unlock()
	if r.alreadyAuthorizedRequests[authorizedRequestStr] {
		r.lggr.Infow("AuthorizeRequest already authorized previously", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr)
		return false, "", errors.New("request already authorized previously")
	}
	currentTimestamp := time.Now().UTC().Unix()
	if currentTimestamp > int64(allowlistedRequest.ExpiryTimestamp) {
		r.lggr.Infow("AuthorizeRequest expired authorization", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr)
		return false, "", errors.New("request authorization expired")
	}
	r.lggr.Infow("AuthorizeRequest success in auth", "method", req.Method, "requestID", req.ID, "authorizedRequestStr", authorizedRequestStr)
	r.alreadyAuthorizedRequests[authorizedRequestStr] = true
	return true, allowlistedRequest.Owner.Hex(), nil
}

func (r *requestAuthorizer) clearExpiredAuthorizedRequests() {
	r.alreadyAuthorizedMutex.Lock()
	defer r.alreadyAuthorizedMutex.Unlock()
	for request := range r.alreadyAuthorizedRequests {
		expiryStr := strings.Split(request, "-->")[1]
		expiry, err := strconv.Atoi(expiryStr)
		if err != nil {
			panic("could not parse expiry timestamp: " + err.Error())
		}
		if time.Now().UTC().Unix() > int64(expiry) {
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
		alreadyAuthorizedRequests: make(map[string]bool),
	}
}
