package vault

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
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

func observationSignalsPendingQueueStall(obs *vaultcommon.Observations) bool {
	return obs != nil && obs.GetPendingQueueStallSignal() == vaultcommon.PendingQueueStallSignal_PENDING_QUEUE_STALL_SIGNAL_STALLED
}

func validatePendingQueueStallSignal(obs *vaultcommon.Observations) error {
	switch obs.GetPendingQueueStallSignal() {
	case vaultcommon.PendingQueueStallSignal_PENDING_QUEUE_STALL_SIGNAL_CONTINUE:
		return nil
	case vaultcommon.PendingQueueStallSignal_PENDING_QUEUE_STALL_SIGNAL_STALLED:
		if len(obs.GetObservations()) != 0 || len(obs.GetPendingQueueItems()) != 0 {
			return errors.New("pending queue stall signal must not include contributions or pending queue items")
		}
		return nil
	default:
		return fmt.Errorf("unknown pending queue stall signal %d", obs.GetPendingQueueStallSignal())
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

// validateContribution is the single source of truth for whether an observation
// contribution is valid for a pending-queue request. It is invoked from the
// Observation self-check, ValidateObservation, and the StateTransition guard so
// the three phases can never disagree on validity: a contribution that one
// phase accepts is accepted by all of them.
//
// The request is read from the pending-queue item (the deterministic source of
// truth shared across nodes); an embedded request, when present, must match it.
// KV-dependent rules (key existence, owner secret counts) remain
// StateTransition-only and are intentionally not checked here.
func (r *ReportingPlugin) validateContribution(
	ctx context.Context,
	_ ReadKVStore,
	pendingItem *vaultcommon.StoredPendingQueueItem,
	o *vaultcommon.Observation,
) error {
	if o.Id == "" {
		return errors.New("request id cannot be empty")
	}

	if observationContributionIsErr(o) {
		if o.Request != nil || o.Response != nil {
			return errors.New("error contribution must not include request or response fields")
		}
		return nil
	}

	if pendingItem == nil {
		return fmt.Errorf("no pending queue item found for request id %s", o.Id)
	}

	payload, err := pendingItem.Item.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("failed to unmarshal pending queue request payload: %w", err)
	}

	if got := requestTypeForPayload(payload); got != o.RequestType {
		return fmt.Errorf("observation request type %s does not match pending queue item type %s", o.RequestType, got)
	}

	switch tp := payload.(type) {
	case *vaultcommon.GetSecretsRequest:
		return r.validateGetSecretsContribution(ctx, tp, o)
	case *vaultcommon.CreateSecretsRequest:
		return r.validateCreateSecretsContribution(ctx, tp, o)
	case *vaultcommon.UpdateSecretsRequest:
		return r.validateUpdateSecretsContribution(ctx, tp, o)
	case *vaultcommon.DeleteSecretsRequest:
		return r.validateDeleteSecretsContribution(ctx, tp, o)
	case *vaultcommon.ListSecretIdentifiersRequest:
		return r.validateListSecretIdentifiersContribution(ctx, tp, o)
	default:
		return fmt.Errorf("unknown request type %T", payload)
	}
}

func (r *ReportingPlugin) validateGetSecretsContribution(ctx context.Context, req *vaultcommon.GetSecretsRequest, o *vaultcommon.Observation) error {
	if embedded := o.GetGetSecretsRequest(); embedded != nil && !proto.Equal(embedded, req) {
		return errors.New("embedded GetSecrets request does not match pending queue request")
	}

	if err := r.validateGetSecretsRequestPayload(ctx, req); err != nil {
		return err
	}

	resp := o.GetGetSecretsResponse()
	if resp == nil {
		return errors.New("GetSecrets observation must have a response")
	}

	if err := validateRequestResponseItemCount(len(req.Requests), len(resp.Responses), "GetSecrets"); err != nil {
		return err
	}

	respMap := map[string]*vaultcommon.SecretResponse{}
	for _, secretResponse := range resp.Responses {
		if secretResponse.Id == nil {
			return errors.New("GetSecrets response contains nil secret identifier")
		}
		if err := r.validator.ValidateSecretIdentifier(ctx, secretResponse.Id.Key, secretResponse.Id.Owner, secretResponse.Id.Namespace); err != nil {
			return fmt.Errorf("GetSecrets response contains invalid secret identifier: %w", err)
		}
		key := vaulttypes.KeyFor(secretResponse.Id)
		if _, ok := respMap[key]; ok {
			return fmt.Errorf("duplicate response found for item %s", key)
		}
		respMap[key] = secretResponse
	}

	return r.validateGetSecretsResponseShares(ctx, req, respMap)
}

func (r *ReportingPlugin) validateCreateSecretsContribution(ctx context.Context, req *vaultcommon.CreateSecretsRequest, o *vaultcommon.Observation) error {
	if embedded := o.GetCreateSecretsRequest(); embedded != nil && !proto.Equal(embedded, req) {
		return errors.New("embedded CreateSecrets request does not match pending queue request")
	}

	if err := r.validator.CheckRequestBatchSize(ctx, len(req.EncryptedSecrets)); err != nil {
		return err
	}

	resp := o.GetCreateSecretsResponse()
	if resp == nil {
		return errors.New("CreateSecrets observation must have a response")
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

	for _, resp := range resp.Responses {
		if resp.Id == nil {
			return errors.New("CreateSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateUpdateSecretsContribution(ctx context.Context, req *vaultcommon.UpdateSecretsRequest, o *vaultcommon.Observation) error {
	if embedded := o.GetUpdateSecretsRequest(); embedded != nil && !proto.Equal(embedded, req) {
		return errors.New("embedded UpdateSecrets request does not match pending queue request")
	}

	if err := r.validator.CheckRequestBatchSize(ctx, len(req.EncryptedSecrets)); err != nil {
		return err
	}

	resp := o.GetUpdateSecretsResponse()
	if resp == nil {
		return errors.New("UpdateSecrets observation must have a response")
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

	for _, resp := range resp.Responses {
		if resp.Id == nil {
			return errors.New("UpdateSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateDeleteSecretsContribution(ctx context.Context, req *vaultcommon.DeleteSecretsRequest, o *vaultcommon.Observation) error {
	if embedded := o.GetDeleteSecretsRequest(); embedded != nil && !proto.Equal(embedded, req) {
		return errors.New("embedded DeleteSecrets request does not match pending queue request")
	}

	if err := r.validateDeleteSecretsRequestPayload(ctx, req); err != nil {
		return err
	}

	resp := o.GetDeleteSecretsResponse()
	if resp == nil {
		return errors.New("DeleteSecrets observation must have a response")
	}

	if err := validateRequestResponseItemCount(len(req.Ids), len(resp.Responses), "DeleteSecrets"); err != nil {
		return err
	}

	for _, resp := range resp.Responses {
		if resp.Id == nil {
			return errors.New("DeleteSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateListSecretIdentifiersContribution(ctx context.Context, req *vaultcommon.ListSecretIdentifiersRequest, o *vaultcommon.Observation) error {
	if embedded := o.GetListSecretIdentifiersRequest(); embedded != nil && !proto.Equal(embedded, req) {
		return errors.New("embedded ListSecretIdentifiers request does not match pending queue request")
	}

	resp := o.GetListSecretIdentifiersResponse()
	if resp == nil {
		return errors.New("ListSecretIdentifiers observation must have a response")
	}

	if !resp.Success {
		if resp.GetError() != "" {
			return fmt.Errorf("%s", resp.GetError())
		}
		return errors.New("ListSecretIdentifiers observation failed")
	}

	if err := r.validateListSecretIdentifiersOwnerWire(ctx, req); err != nil {
		return fmt.Errorf("ListSecretIdentifiers request contains invalid secret identifier: %w", err)
	}

	if err := r.validateListSecretIdentifiersResponseSize(ctx, req.Owner, len(resp.Identifiers)); err != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[int]](err); ok {
			return fmt.Errorf("ListSecretIdentifiers response exceeds maximum number of secrets per owner (have=%d, limit=%d): %w", len(resp.Identifiers), errBoundLimited.Limit, err)
		}
		return fmt.Errorf("failed to check max secrets per owner limit: %w", err)
	}

	return nil
}

func consensusObservationError(errObs []*vaultcommon.Observation, f int) string {
	const fallback = "request is not valid"

	counts := make(map[string]int)
	for _, o := range errObs {
		msg := o.GetError().GetMessage()
		if msg != "" {
			counts[msg]++
		}
	}

	var consensusMsg string
	var maxCount int
	for msg, count := range counts {
		if count >= f+1 && count > maxCount {
			consensusMsg = msg
			maxCount = count
		}
	}

	if consensusMsg != "" {
		return consensusMsg
	}
	return fallback
}

func classifyContributions(obs []*vaultcommon.Observation) (ok, err []*vaultcommon.Observation) {
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

// validateEncryptedShareSize validates the structure and size limit of a single
// encrypted share entry. Returns nil if the share is valid.
func (r *ReportingPlugin) validateEncryptedShareSize(ctx context.Context, es *vaultcommon.EncryptedShares) error {
	if err := validateEncryptedSharesEntry(es); err != nil {
		return err
	}
	shareSize, err := encryptedShareSizeForLimit(es)
	if err != nil {
		return err
	}
	if err := r.cfg.MaxShareLengthBytes.Check(ctx, pkgconfig.Size(shareSize)*pkgconfig.Byte); err != nil {
		if _, ok := errors.AsType[limits.ErrorBoundLimited[pkgconfig.Size]](err); ok {
			return fmt.Errorf("share provided exceeds maximum size allowed: %w", err)
		}
		return errors.New("failed to check share size")
	}
	return nil
}
