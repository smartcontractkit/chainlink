package vault

import (
	"context"
	"encoding/hex"
	"fmt"
	"maps"
	"slices"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

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
			logUserErrorAware(l, "failed to observe update secret request item", ierr, "id", sr.Id)
			errorMsg := userFacingError(ierr, "failed to handle update secret request")
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

	if !r.optimizationsEnabled(ctx) {
		newReqs := make([]*vaultcommon.EncryptedSecret, 0, len(idToReqs))
		for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
			newReqs = append(newReqs, idToReqs[sreq])
		}

		o.Request = &vaultcommon.Outcome_UpdateSecretsRequest{
			UpdateSecretsRequest: &vaultcommon.UpdateSecretsRequest{
				RequestId:        reqID,
				EncryptedSecrets: newReqs,
			},
		}
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
			logUserErrorAware(r.lggr, "failed to handle update secret request", err, "id", req.Id, "requestID", reqID)
			errorMsg := userFacingError(err, "failed to handle update secret request")
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

	secret, err := store.GetSecret(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret == nil {
		return nil, newUserError("could not write update to key value store: key does not exist")
	}

	err = store.WriteSecret(ctx, req.Id, &vaultcommon.StoredSecret{
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
