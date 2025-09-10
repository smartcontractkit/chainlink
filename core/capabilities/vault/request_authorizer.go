package vault

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
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
	digest, err := r.digestForRequest(req)
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

func (r *requestAuthorizer) digestForRequest(req jsonrpc.Request[json.RawMessage]) ([32]byte, error) {
	var seed any
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		var createSecretsRequests vaultcommon.CreateSecretsRequest
		if err := json.Unmarshal(*req.Params, &createSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling create secrets request: " + err.Error())
		}
		seed = vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: createSecretsRequests.EncryptedSecrets,
		}
	case vaulttypes.MethodSecretsUpdate:
		var updateSecretsRequests vaultcommon.UpdateSecretsRequest
		if err := json.Unmarshal(*req.Params, &updateSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling update secrets request: " + err.Error())
		}
		seed = vaultcommon.UpdateSecretsRequest{
			EncryptedSecrets: updateSecretsRequests.EncryptedSecrets,
		}
	case vaulttypes.MethodSecretsList:
		var listSecretsRequests vaultcommon.ListSecretIdentifiersRequest
		if err := json.Unmarshal(*req.Params, &listSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling list secrets request: " + err.Error())
		}
		seed = vaultcommon.ListSecretIdentifiersRequest{
			Owner:     listSecretsRequests.Owner,
			Namespace: listSecretsRequests.Namespace,
		}
	case vaulttypes.MethodSecretsDelete:
		var deleteSecretsRequests vaultcommon.DeleteSecretsRequest
		if err := json.Unmarshal(*req.Params, &deleteSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling delete secrets request: " + err.Error())
		}
		seed = vaultcommon.DeleteSecretsRequest{
			Ids: deleteSecretsRequests.Ids,
		}
	default:
		return [32]byte{}, fmt.Errorf("unauthorized method: %s", req.Method)
	}

	return CalculateRequestDigest(seed), nil
}

// CalculateRequestDigest creates a SHA256 digest of the request for integrity verification
// This function is shared between client (JWT generation) and server (JWT validation)
func CalculateRequestDigest(req any) [32]byte {
	var data []byte
	if m, ok := req.(proto.Message); ok {
		// Use protobuf canonical serialization
		serialized, err := proto.Marshal(m)
		if err == nil {
			data = serialized
		} else {
			// fallback to string representation if marshal fails
			data = []byte(fmt.Sprintf("%v", req))
		}
	} else if s, ok := req.(fmt.Stringer); ok {
		data = []byte(s.String())
	} else {
		data = []byte(fmt.Sprintf("%v", req))
	}

	hash := sha256.Sum256(data)
	return hash
}

func NewRequestAuthorizer(lggr logger.Logger, workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer) *requestAuthorizer {
	return &requestAuthorizer{
		workflowRegistrySyncer:    workflowRegistrySyncer,
		lggr:                      logger.Named(lggr, "VaultRequestAuthorizer"),
		alreadyAuthorizedRequests: make(map[string]bool),
	}
}
