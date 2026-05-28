package vault

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"slices"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

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
				Id:      validatedID,
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

func (r *ReportingPlugin) observeCreateSecretRequest(ctx context.Context, _ ReadKVStore, secretRequest *vaultcommon.EncryptedSecret, requestsCountForID map[string]int) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateSecretIdentifier(ctx, secretRequest.Id)
	if err != nil {
		return id, err
	}

	if requestsCountForID[vaulttypes.KeyFor(secretRequest.Id)] > 1 {
		return id, newUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}

	if ierr := r.validator.ValidateCiphertextSize(ctx, secretRequest.Id.Owner, secretRequest.EncryptedValue); ierr != nil {
		return id, newUserError(ierr.Error())
	}

	err = vaultcap.EnsureRightLabelOnSecret(r.cfg.PublicKey, secretRequest.EncryptedValue, secretRequest.Id.Owner)
	if err != nil {
		return id, newUserError("failed to verify ciphertext: " + err.Error())
	}

	return id, nil
}

func (r *ReportingPlugin) stateTransitionCreateSecrets(ctx context.Context, store WriteKVStore, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	first := chosen[0]
	reqID := first.GetCreateSecretsRequest().RequestId
	req := first.GetCreateSecretsRequest().EncryptedSecrets
	idToReqs := map[string]*vaultcommon.EncryptedSecret{}
	for _, r := range req {
		idToReqs[vaulttypes.KeyFor(r.Id)] = r
	}

	if !r.optimizationsEnabled(ctx) {
		newReqs := make([]*vaultcommon.EncryptedSecret, 0, len(idToReqs))
		for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
			newReqs = append(newReqs, idToReqs[sreq])
		}

		o.Request = &vaultcommon.Outcome_CreateSecretsRequest{
			CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
				RequestId:        reqID,
				EncryptedSecrets: newReqs,
			},
		}
	}

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
