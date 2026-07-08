package vault

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

func observationContributionIsErr(o *vaultcommon.Observation) bool {
	return o.GetError() != nil && o.GetError().GetMessage() != ""
}

func observationContributionIsOk(o *vaultcommon.Observation) bool {
	if observationContributionIsErr(o) {
		return false
	}
	switch o.RequestType {
	case vaultcommon.RequestType_GET_SECRETS:
		return o.GetGetSecretsResponse() != nil
	case vaultcommon.RequestType_CREATE_SECRETS:
		return o.GetCreateSecretsResponse() != nil
	case vaultcommon.RequestType_UPDATE_SECRETS:
		return o.GetUpdateSecretsResponse() != nil
	case vaultcommon.RequestType_DELETE_SECRETS:
		return o.GetDeleteSecretsResponse() != nil
	case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
		return o.GetListSecretIdentifiersResponse() != nil
	default:
		return false
	}
}

func observationToErrContribution(o *vaultcommon.Observation, msg string) *vaultcommon.Observation {
	out := &vaultcommon.Observation{
		Id:          o.Id,
		RequestType: o.RequestType,
		Error: &vaultcommon.ObservationError{
			Message: msg,
		},
	}
	return out
}

func (r *ReportingPlugin) checkContribution(
	ctx context.Context,
	readKV ReadKVStore,
	req *vaultcommon.StoredPendingQueueItem,
	o *vaultcommon.Observation,
) error {
	payload, err := req.Item.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("failed to unmarshal request payload: %w", err)
	}

	switch tp := payload.(type) {
	case *vaultcommon.GetSecretsRequest:
		return r.checkGetSecretsContribution(ctx, readKV, tp, o)
	case *vaultcommon.CreateSecretsRequest:
		return r.checkCreateSecretsContribution(ctx, readKV, tp, o)
	case *vaultcommon.UpdateSecretsRequest:
		return r.checkUpdateSecretsContribution(ctx, readKV, tp, o)
	case *vaultcommon.DeleteSecretsRequest:
		return r.checkDeleteSecretsContribution(ctx, readKV, tp, o)
	case *vaultcommon.ListSecretIdentifiersRequest:
		return r.checkListSecretIdentifiersContribution(ctx, readKV, tp, o)
	default:
		return fmt.Errorf("unknown request type %T", payload)
	}
}

func (r *ReportingPlugin) checkGetSecretsContribution(
	ctx context.Context,
	_ ReadKVStore,
	req *vaultcommon.GetSecretsRequest,
	o *vaultcommon.Observation,
) error {
	if err := r.validateGetSecretsRequestPayload(ctx, req); err != nil {
		return err
	}
	resp := o.GetGetSecretsResponse()
	if resp == nil {
		return fmt.Errorf("GetSecrets observation must have a response")
	}
	return validateRequestResponseItemCount(len(req.Requests), len(resp.Responses), "GetSecrets")
}

func (r *ReportingPlugin) checkCreateSecretsContribution(
	ctx context.Context,
	readKV ReadKVStore,
	req *vaultcommon.CreateSecretsRequest,
	o *vaultcommon.Observation,
) error {
	if err := r.validator.CheckRequestBatchSize(ctx, len(req.EncryptedSecrets)); err != nil {
		return err
	}
	resp := o.GetCreateSecretsResponse()
	if resp == nil {
		return fmt.Errorf("CreateSecrets observation must have a response")
	}
	if err := validateRequestResponseItemCount(len(req.EncryptedSecrets), len(resp.Responses), "CreateSecrets"); err != nil {
		return err
	}
	requestsCountForID := buildEncryptedSecretIdentifierCounts(req.EncryptedSecrets)
	for _, sr := range req.EncryptedSecrets {
		if _, err := r.validateEncryptedSecretPayload(ctx, sr, requestsCountForID); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReportingPlugin) checkUpdateSecretsContribution(
	ctx context.Context,
	readKV ReadKVStore,
	req *vaultcommon.UpdateSecretsRequest,
	o *vaultcommon.Observation,
) error {
	if err := r.validator.CheckRequestBatchSize(ctx, len(req.EncryptedSecrets)); err != nil {
		return err
	}
	resp := o.GetUpdateSecretsResponse()
	if resp == nil {
		return fmt.Errorf("UpdateSecrets observation must have a response")
	}
	if err := validateRequestResponseItemCount(len(req.EncryptedSecrets), len(resp.Responses), "UpdateSecrets"); err != nil {
		return err
	}
	requestsCountForID := buildEncryptedSecretIdentifierCounts(req.EncryptedSecrets)
	for _, sr := range req.EncryptedSecrets {
		if _, err := r.validateEncryptedSecretPayload(ctx, sr, requestsCountForID); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReportingPlugin) checkDeleteSecretsContribution(
	ctx context.Context,
	readKV ReadKVStore,
	req *vaultcommon.DeleteSecretsRequest,
	o *vaultcommon.Observation,
) error {
	if err := r.validateDeleteSecretsRequestPayload(ctx, req); err != nil {
		return err
	}
	resp := o.GetDeleteSecretsResponse()
	if resp == nil {
		return fmt.Errorf("DeleteSecrets observation must have a response")
	}
	return validateRequestResponseItemCount(len(req.Ids), len(resp.Responses), "DeleteSecrets")
}

func (r *ReportingPlugin) checkListSecretIdentifiersContribution(
	ctx context.Context,
	_ ReadKVStore,
	req *vaultcommon.ListSecretIdentifiersRequest,
	o *vaultcommon.Observation,
) error {
	resp := o.GetListSecretIdentifiersResponse()
	if resp == nil {
		return fmt.Errorf("ListSecretIdentifiers observation must have a response")
	}
	if !resp.Success {
		if resp.GetError() != "" {
			return fmt.Errorf("%s", resp.GetError())
		}
		return fmt.Errorf("ListSecretIdentifiers observation failed")
	}
	if err := r.validateListSecretIdentifiersOwnerWire(ctx, req); err != nil {
		return err
	}
	return r.validateListSecretIdentifiersResponseSize(ctx, req.Owner, len(resp.Identifiers))
}

func combineObservationErrors(errObs []*vaultcommon.Observation) string {
	const fallback = "request is not valid"

	seen := make(map[string]struct{}, len(errObs))
	messages := make([]string, 0, len(errObs))
	for _, o := range errObs {
		msg := o.GetError().GetMessage()
		if msg == "" {
			continue
		}
		if _, dup := seen[msg]; dup {
			continue
		}
		seen[msg] = struct{}{}
		messages = append(messages, msg)
	}
	if len(messages) == 0 {
		return fallback
	}
	return strings.Join(messages, "; ")
}

func classifyContributions(obs []*vaultcommon.Observation) (ok []*vaultcommon.Observation, err []*vaultcommon.Observation) {
	for _, o := range obs {
		switch {
		case observationContributionIsErr(o):
			err = append(err, o)
		case observationContributionIsOk(o):
			ok = append(ok, o)
		}
	}
	return ok, err
}

// observerOkCoverage counts the distinct pending-queue ids for which the observer contributed
// an Ok observation. pendingIDs scopes coverage to the current queue (Byzantine ids outside the
// queue are ignored); pass nil to count all Ok contributions. Used to attribute prefix divergence
// to a specific oracle — a node consistently reporting lower coverage than peers is withholding
// or truncating its observation prefix, which stalls head-of-queue quorum under include-invalid.
func observerOkCoverage(obs *vaultcommon.Observations, pendingIDs map[string]bool) int {
	if obs == nil {
		return 0
	}
	seen := map[string]bool{}
	for _, o := range obs.Observations {
		if !observationContributionIsOk(o) {
			continue
		}
		if pendingIDs != nil && !pendingIDs[o.Id] {
			continue
		}
		seen[o.Id] = true
	}
	return len(seen)
}

// coverageSpread returns max-min of per-observer Ok prefix coverage. A non-zero spread means
// oracles disagree on how much of the pending queue they observed — the head-of-queue stall
// signature under include-invalid.
func coverageSpread(coverages []int) int {
	if len(coverages) == 0 {
		return 0
	}
	minC, maxC := coverages[0], coverages[0]
	for _, c := range coverages[1:] {
		if c < minC {
			minC = c
		}
		if c > maxC {
			maxC = c
		}
	}
	return maxC - minC
}

func requestTypeForPayload(payload proto.Message) vaultcommon.RequestType {
	switch payload.(type) {
	case *vaultcommon.GetSecretsRequest:
		return vaultcommon.RequestType_GET_SECRETS
	case *vaultcommon.CreateSecretsRequest:
		return vaultcommon.RequestType_CREATE_SECRETS
	case *vaultcommon.UpdateSecretsRequest:
		return vaultcommon.RequestType_UPDATE_SECRETS
	case *vaultcommon.DeleteSecretsRequest:
		return vaultcommon.RequestType_DELETE_SECRETS
	case *vaultcommon.ListSecretIdentifiersRequest:
		return vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS
	default:
		return vaultcommon.RequestType_UNKNOWN
	}
}

func buildRejectedOutcome(id string, payload proto.Message, requestType vaultcommon.RequestType, errMsg string) *vaultcommon.Outcome {
	o := &vaultcommon.Outcome{
		Id:          id,
		RequestType: requestType,
	}
	switch requestType {
	case vaultcommon.RequestType_GET_SECRETS:
		o.Response = &vaultcommon.Outcome_GetSecretsResponse{
			GetSecretsResponse: &vaultcommon.GetSecretsResponse{
				Responses: rejectedGetSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_CREATE_SECRETS:
		o.Response = &vaultcommon.Outcome_CreateSecretsResponse{
			CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{
				Responses: rejectedCreateSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_UPDATE_SECRETS:
		o.Response = &vaultcommon.Outcome_UpdateSecretsResponse{
			UpdateSecretsResponse: &vaultcommon.UpdateSecretsResponse{
				Responses: rejectedUpdateSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_DELETE_SECRETS:
		o.Response = &vaultcommon.Outcome_DeleteSecretsResponse{
			DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
				Responses: rejectedDeleteSecretsResponses(payload, errMsg),
			},
		}
	case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
		o.Response = &vaultcommon.Outcome_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: &vaultcommon.ListSecretIdentifiersResponse{
				Success: false,
				Error:   errMsg,
			},
		}
	}
	return o
}

func rejectedGetSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.SecretResponse {
	req, ok := payload.(*vaultcommon.GetSecretsRequest)
	if !ok || len(req.Requests) == 0 {
		return []*vaultcommon.SecretResponse{{
			Result: &vaultcommon.SecretResponse_Error{Error: errMsg},
		}}
	}
	resps := make([]*vaultcommon.SecretResponse, len(req.Requests))
	for i, sr := range req.Requests {
		resps[i] = &vaultcommon.SecretResponse{
			Id: sr.Id,
			Result: &vaultcommon.SecretResponse_Error{
				Error: errMsg,
			},
		}
	}
	return resps
}

func rejectedCreateSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.CreateSecretResponse {
	req, ok := payload.(*vaultcommon.CreateSecretsRequest)
	if !ok || len(req.EncryptedSecrets) == 0 {
		return []*vaultcommon.CreateSecretResponse{{
			Success: false,
			Error:   errMsg,
		}}
	}
	resps := make([]*vaultcommon.CreateSecretResponse, len(req.EncryptedSecrets))
	for i, sr := range req.EncryptedSecrets {
		resps[i] = &vaultcommon.CreateSecretResponse{
			Id:      sr.Id,
			Success: false,
			Error:   errMsg,
		}
	}
	return resps
}

func rejectedUpdateSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.UpdateSecretResponse {
	req, ok := payload.(*vaultcommon.UpdateSecretsRequest)
	if !ok || len(req.EncryptedSecrets) == 0 {
		return []*vaultcommon.UpdateSecretResponse{{
			Success: false,
			Error:   errMsg,
		}}
	}
	resps := make([]*vaultcommon.UpdateSecretResponse, len(req.EncryptedSecrets))
	for i, sr := range req.EncryptedSecrets {
		resps[i] = &vaultcommon.UpdateSecretResponse{
			Id:      sr.Id,
			Success: false,
			Error:   errMsg,
		}
	}
	return resps
}

func rejectedDeleteSecretsResponses(payload proto.Message, errMsg string) []*vaultcommon.DeleteSecretResponse {
	req, ok := payload.(*vaultcommon.DeleteSecretsRequest)
	if !ok || len(req.Ids) == 0 {
		return []*vaultcommon.DeleteSecretResponse{{
			Success: false,
			Error:   errMsg,
		}}
	}
	resps := make([]*vaultcommon.DeleteSecretResponse, len(req.Ids))
	for i, id := range req.Ids {
		resps[i] = &vaultcommon.DeleteSecretResponse{
			Id:      id,
			Success: false,
			Error:   errMsg,
		}
	}
	return resps
}
