package vault

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/build"
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
	// TODO(https://smartcontract-it.atlassian.net/browse/PRIV-175): Remove this bypass once we have a Vault DON e2e tests setup in local CRE.
	if build.IsDev() {
		r.lggr.Warnw("bypassing RequestAuthorizer since it is not a production build", "build-mode", build.Mode())
		// returning owner as Owner1, since that's used in vault e2e tests.
		return true, "Owner1", nil
	}
	defer r.clearExpiredAuthorizedRequests()
	digest, err := workflowsyncerv2.DigestForRequest(req)
	if err != nil {
		return false, "", err
	}
	allowlistedRequest := r.fetchAllowlistedItem(r.workflowRegistrySyncer.GetAllowlistedRequests(ctx), digest)
	if allowlistedRequest == nil {
		return false, "", errors.New("request not allowlisted")
	}
	authorizedRequestStr := string(allowlistedRequest.RequestDigest[:]) + "-->" + strconv.FormatUint(uint64(allowlistedRequest.ExpiryTimestamp), 10)
	r.alreadyAuthorizedMutex.Lock()
	defer r.alreadyAuthorizedMutex.Unlock()
	if r.alreadyAuthorizedRequests[authorizedRequestStr] {
		return false, "", errors.New("request already authorized previously")
	}
	currentTimestamp := time.Now().UTC().Unix()
	if currentTimestamp > int64(allowlistedRequest.ExpiryTimestamp) {
		return false, "", errors.New("request authorization expired")
	}
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
