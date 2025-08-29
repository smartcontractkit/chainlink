package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	nautilus "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/utils"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	workflowsyncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
)

type RequestAuthorizer struct {
	workflowRegistrySyncer    workflowsyncerv2.WorkflowRegistrySyncer
	alreadyAuthorizedRequests map[string]bool
	alreadyAuthorizedMutex    sync.Mutex
	lggr                      logger.Logger
}

func (r *RequestAuthorizer) AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (bool, string, error) {
	defer r.clearExpiredAuthorizedRequests()
	digest, err := r.digestForRequest(req)
	if err != nil {
		return false, "", err
	}
	allowlistedRequest := r.fetchAllowlistedItem(r.workflowRegistrySyncer.GetAllowlistedRequests(ctx), digest)
	if allowlistedRequest == nil {
		return false, "", errors.New("request not allowlisted")
	}
	authorizedRequestStr := string(allowlistedRequest.RequestDigest[:]) + "-->" + string(allowlistedRequest.ExpiryTimestamp)
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

func (r *RequestAuthorizer) clearExpiredAuthorizedRequests() {
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

func (r *RequestAuthorizer) fetchAllowlistedItem(allowListedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, digest string) *workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest {
	for _, item := range allowListedRequests {
		if string(item.RequestDigest[:]) == digest {
			return &item
		}
	}
	return nil
}

func (r *RequestAuthorizer) digestForRequest(req jsonrpc.Request[json.RawMessage]) (string, error) {
	var seed any
	switch req.Method {
	case MethodSecretsCreate:
		var createSecretsRequests vaultcommon.CreateSecretsRequest
		if err := json.Unmarshal(*req.Params, &createSecretsRequests); err != nil {
			return "", errors.New("error unmarshaling create secrets request: " + err.Error())
		}
		seed = vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: createSecretsRequests.EncryptedSecrets,
		}
	case MethodSecretsUpdate:
		var updateSecretsRequests vaultcommon.UpdateSecretsRequest
		if err := json.Unmarshal(*req.Params, &updateSecretsRequests); err != nil {
			return "", errors.New("error unmarshaling update secrets request: " + err.Error())
		}
		seed = vaultcommon.UpdateSecretsRequest{
			EncryptedSecrets: updateSecretsRequests.EncryptedSecrets,
		}
	case MethodSecretsDelete:
		var deleteSecretsRequests vaultcommon.DeleteSecretsRequest
		if err := json.Unmarshal(*req.Params, &deleteSecretsRequests); err != nil {
			return "", errors.New("error unmarshaling delete secrets request: " + err.Error())
		}
		seed = vaultcommon.DeleteSecretsRequest{
			Ids: deleteSecretsRequests.Ids,
		}
	default:
		return "", fmt.Errorf("unauthorized method: %s", req.Method)
	}

	return nautilus.CalculateRequestDigest(&seed), nil
}

func NewRequestAuthorizer(lggr logger.Logger, workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer) RequestAuthorizer {
	return RequestAuthorizer{
		workflowRegistrySyncer:    workflowRegistrySyncer,
		lggr:                      logger.Named(lggr, "VaultRequestAuthorizer"),
		alreadyAuthorizedRequests: make(map[string]bool),
	}
}
