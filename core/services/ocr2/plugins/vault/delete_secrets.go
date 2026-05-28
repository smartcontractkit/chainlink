package vault

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func (r *ReportingPlugin) observeDeleteSecrets(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.DeleteSecretsRequest)
	o.RequestType = vaultcommon.RequestType_DELETE_SECRETS
	o.Request = &vaultcommon.Observation_DeleteSecretsRequest{
		DeleteSecretsRequest: tp,
	}
	l := r.lggr.With("requestId", tp.RequestId, "requestType", "DeleteSecrets")

	requestsCountForID := map[string]int{}
	for _, sr := range tp.Ids {
		var key string
		if sr == nil {
			key = "<nil>"
		} else {
			key = vaulttypes.KeyFor(sr)
		}
		requestsCountForID[key]++
	}

	resps := []*vaultcommon.DeleteSecretResponse{}
	for _, id := range tp.Ids {
		validatedID, ierr := r.observeDeleteSecretRequest(ctx, reader, id, requestsCountForID)
		if ierr != nil {
			logUserErrorAware(l, "failed to handle delete secret request item", ierr, "id", id)
			errorMsg := userFacingError(ierr, "failed to handle delete secret request")
			resps = append(resps, &vaultcommon.DeleteSecretResponse{
				Id:      id,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			l.Debugw("observed delete secret request item", "id", validatedID)
			resps = append(resps, &vaultcommon.DeleteSecretResponse{
				Id: validatedID,
				// false because it hasn't been processed yet.
				// When the write is handled successfully in StateTransition
				// we'll update this to true.
				Success: false,
			})
		}
	}

	o.Response = &vaultcommon.Observation_DeleteSecretsResponse{
		DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
			Responses: resps,
		},
	}
}

func (r *ReportingPlugin) observeDeleteSecretRequest(ctx context.Context, reader ReadKVStore, identifier *vaultcommon.SecretIdentifier, requestsCountForID map[string]int) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateSecretIdentifier(ctx, identifier)
	if err != nil {
		return id, err
	}

	if requestsCountForID[vaulttypes.KeyFor(identifier)] > 1 {
		return id, newUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}

	ss, err := reader.GetSecret(ctx, id)
	if err != nil {
		return id, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if ss == nil {
		return id, newUserError("key does not exist")
	}

	return id, nil
}

func (r *ReportingPlugin) stateTransitionDeleteSecrets(ctx context.Context, store WriteKVStore, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	first := chosen[0]
	reqID := first.GetDeleteSecretsRequest().RequestId
	// First we'll aggregate the requests.
	// Since the shas for all requests match, we can just take the first entry
	// and sort the requests contained within it.
	req := first.GetDeleteSecretsRequest().Ids
	idToReqs := map[string]*vaultcommon.SecretIdentifier{}
	for _, r := range req {
		idToReqs[vaulttypes.KeyFor(r)] = r
	}

	if !r.optimizationsEnabled(ctx) {
		newReqs := make([]*vaultcommon.SecretIdentifier, 0, len(idToReqs))
		for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
			newReqs = append(newReqs, idToReqs[sreq])
		}

		o.Request = &vaultcommon.Outcome_DeleteSecretsRequest{
			DeleteSecretsRequest: &vaultcommon.DeleteSecretsRequest{
				RequestId: reqID,
				Ids:       newReqs,
			},
		}
	}

	// Next let's aggregate the responses.
	// We do this by taking the first response, and determine if
	// there was a validation error. If not, we write it to the key value store.
	// The responses are sorted by Id.
	resp := first.GetDeleteSecretsResponse()
	idToResps := map[string]*vaultcommon.DeleteSecretResponse{}
	for _, r := range resp.Responses {
		idToResps[vaulttypes.KeyFor(r.Id)] = r
	}

	sortedResps := []*vaultcommon.DeleteSecretResponse{}
	for _, id := range slices.Sorted(maps.Keys(idToResps)) {
		resp := idToResps[id]
		req, found := idToReqs[id]
		if !found {
			r.lggr.Errorw("could not find request for response", "id", id)
			sortedResps = append(sortedResps, &vaultcommon.DeleteSecretResponse{
				Id:      resp.Id,
				Success: false,
				Error:   "internal error: could not find request for response",
			})
			continue
		}
		resp, err := r.stateTransitionDeleteSecretsRequest(ctx, store, req, resp)
		if err != nil {
			logUserErrorAware(r.lggr, "failed to handle delete secret request", err, "id", id, "requestId", reqID)
			errorMsg := userFacingError(err, "failed to handle delete secret request")
			sortedResps = append(sortedResps, &vaultcommon.DeleteSecretResponse{
				Id:      req,
				Success: false,
				Error:   errorMsg,
			})
		} else {
			r.lggr.Debugw("successfully deleted secret in key value store", "method", "DeleteSecrets", "key", vaulttypes.KeyFor(req), "requestId", reqID)
			sortedResps = append(sortedResps, resp)
		}
	}

	o.Response = &vaultcommon.Outcome_DeleteSecretsResponse{
		DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
			Responses: sortedResps,
		},
	}
}

func (r *ReportingPlugin) stateTransitionDeleteSecretsRequest(ctx context.Context, store WriteKVStore, id *vaultcommon.SecretIdentifier, resp *vaultcommon.DeleteSecretResponse) (*vaultcommon.DeleteSecretResponse, error) {
	if resp.GetError() != "" {
		return resp, newUserError(resp.GetError())
	}

	err := store.DeleteSecret(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete secret from key value store: %w", err)
	}

	return &vaultcommon.DeleteSecretResponse{
		Id:      id,
		Success: true,
		Error:   "",
	}, nil
}
