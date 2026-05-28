package vault

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func (r *ReportingPlugin) validateSecretIdentifier(ctx context.Context, id *vaultcommon.SecretIdentifier) (*vaultcommon.SecretIdentifier, error) {
	if id == nil {
		return nil, newUserError("secret identifier cannot be nil")
	}

	if err := r.validator.ValidateSecretIdentifier(ctx, id.Key, id.Owner, id.Namespace); err != nil {
		return nil, newUserError(err.Error())
	}

	newID := &vaultcommon.SecretIdentifier{
		Key:       id.Key,
		Owner:     id.Owner,
		Namespace: id.Namespace,
	}

	return newID, nil
}

func shaForProto(msg proto.Message) (string, error) {
	protoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("could not generate sha for proto message: failed to marshal proto: %w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(protoBytes)), nil
}

func shaForObservation(o *vaultcommon.Observation) (string, error) {
	switch o.RequestType {
	case vaultcommon.RequestType_GET_SECRETS:
		cloned := proto.CloneOf(o)
		for _, r := range cloned.GetGetSecretsResponse().Responses {
			if r.GetData() != nil {
				// Exclude the encrypted shares from the sha, as these need to be aggregated later.
				r.GetData().EncryptedDecryptionKeyShares = nil
			}
		}

		return shaForProto(cloned)
	default:
		return shaForProto(o)
	}
}

func (r *ReportingPlugin) checkRequestBatchLimit(ctx context.Context, batchSize int) error {
	if err := r.cfg.MaxRequestBatchSize.Check(ctx, batchSize); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[int]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("max batch size exceeded for request: %w", err)
		}
		// Fail closed here: this could cause a loss of liveness but
		// the current implementation would only return an error that's
		// not a ErrorBoundLimited if the limiter has been closed.
		return errors.New("failed to check batch size")
	}

	return nil
}

func (r *ReportingPlugin) validateObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.Id == "" {
		return errors.New("observation id cannot be empty")
	}

	switch o.RequestType {
	case vaultcommon.RequestType_GET_SECRETS:
		return r.validateGetSecretsObservation(ctx, o)
	case vaultcommon.RequestType_CREATE_SECRETS:
		return r.validateCreateSecretsObservation(ctx, o)
	case vaultcommon.RequestType_UPDATE_SECRETS:
		return r.validateUpdateSecretsObservation(ctx, o)
	case vaultcommon.RequestType_DELETE_SECRETS:
		return r.validateDeleteSecretsObservation(ctx, o)
	case vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS:
		return r.validateListSecretIdentifiersObservation(ctx, o)
	default:
		return errors.New("invalid observation type: " + o.RequestType.String())
	}
}

func (r *ReportingPlugin) validateGetSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	resp := o.GetGetSecretsResponse()
	if resp == nil {
		return errors.New("GetSecrets observation must have a response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(resp.Responses)); err != nil {
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

	for _, rsp := range respMap {
		d := rsp.GetData()
		if d == nil {
			continue
		}

		innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: rsp.Id.Owner})
		for _, ds := range d.GetEncryptedDecryptionKeyShares() {
			if err := validateEncryptedSharesEntry(ds); err != nil {
				return err
			}

			shareSize, err := encryptedShareSizeForLimit(ds)
			if err != nil {
				return err
			}
			if err := r.cfg.MaxShareLengthBytes.Check(innerCtx, pkgconfig.Size(shareSize)*pkgconfig.Byte); err != nil {
				var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
				if errors.As(err, &errBoundLimited) {
					return fmt.Errorf("share provided exceeds maximum size allowed: %w", err)
				}
				return errors.New("failed to check share size")
			}
		}
	}

	return nil
}

func (r *ReportingPlugin) validateCreateSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetCreateSecretsRequest() == nil || o.GetCreateSecretsResponse() == nil {
		return errors.New("CreateSecrets observation must have both request and response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(o.GetCreateSecretsRequest().EncryptedSecrets)); err != nil {
		return err
	}

	if len(o.GetCreateSecretsRequest().EncryptedSecrets) != len(o.GetCreateSecretsResponse().Responses) {
		return errors.New("CreateSecrets request and response must have the same number of items")
	}

	// We disallow duplicate create requests within a single batch request.
	// This prevents users from clobbering their own writes.
	idSet := map[string]bool{}
	for _, s := range o.GetCreateSecretsRequest().EncryptedSecrets {
		if s.Id == nil {
			return errors.New("CreateSecrets request contains nil secret identifier")
		}
		if err := r.validator.ValidateSecretIdentifier(ctx, s.Id.Key, s.Id.Owner, s.Id.Namespace); err != nil {
			return fmt.Errorf("CreateSecrets request contains invalid secret identifier: %w", err)
		}
		_, ok := idSet[vaulttypes.KeyFor(s.Id)]
		if ok {
			return fmt.Errorf("CreateSecrets requests cannot contain duplicate request for a given secret identifier: %s", s.Id)
		}

		idSet[vaulttypes.KeyFor(s.Id)] = true

		if err := r.validator.ValidateCiphertextSize(ctx, s.Id.Owner, s.EncryptedValue); err != nil {
			return fmt.Errorf("CreateSecrets request: %w", err)
		}
	}

	for _, r := range o.GetCreateSecretsResponse().Responses {
		if r.Id == nil {
			return errors.New("CreateSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateUpdateSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetUpdateSecretsRequest() == nil || o.GetUpdateSecretsResponse() == nil {
		return errors.New("UpdateSecrets observation must have both request and response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(o.GetUpdateSecretsRequest().EncryptedSecrets)); err != nil {
		return err
	}

	if len(o.GetUpdateSecretsRequest().EncryptedSecrets) != len(o.GetUpdateSecretsResponse().Responses) {
		return errors.New("UpdateSecrets request and response must have the same number of items")
	}

	// We disallow duplicate update requests within a single batch request.
	// This prevents users from clobbering their own writes.
	idSet := map[string]bool{}
	for _, s := range o.GetUpdateSecretsRequest().EncryptedSecrets {
		if s.Id == nil {
			return errors.New("UpdateSecrets request contains nil secret identifier")
		}
		if err := r.validator.ValidateSecretIdentifier(ctx, s.Id.Key, s.Id.Owner, s.Id.Namespace); err != nil {
			return fmt.Errorf("UpdateSecrets request contains invalid secret identifier: %w", err)
		}
		_, ok := idSet[vaulttypes.KeyFor(s.Id)]
		if ok {
			return fmt.Errorf("UpdateSecrets requests cannot contain duplicate request for a given secret identifier: %s", s.Id)
		}

		idSet[vaulttypes.KeyFor(s.Id)] = true

		if err := r.validator.ValidateCiphertextSize(ctx, s.Id.Owner, s.EncryptedValue); err != nil {
			return fmt.Errorf("UpdateSecrets request: %w", err)
		}
	}

	for _, r := range o.GetUpdateSecretsResponse().Responses {
		if r.Id == nil {
			return errors.New("UpdateSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateDeleteSecretsObservation(ctx context.Context, o *vaultcommon.Observation) error {
	if o.GetDeleteSecretsRequest() == nil || o.GetDeleteSecretsResponse() == nil {
		return errors.New("DeleteSecrets observation must have both request and response")
	}

	if err := r.checkRequestBatchLimit(ctx, len(o.GetDeleteSecretsRequest().Ids)); err != nil {
		return err
	}

	if len(o.GetDeleteSecretsRequest().Ids) != len(o.GetDeleteSecretsResponse().Responses) {
		return errors.New("DeleteSecrets request and response must have the same number of items")
	}

	// We disallow duplicate delete requests within a single batch request.
	// This prevents users from clobbering their own writes.
	idSet := map[string]bool{}
	for _, id := range o.GetDeleteSecretsRequest().Ids {
		if id == nil {
			return errors.New("DeleteSecrets request contains nil secret identifier")
		}
		if err := r.validator.ValidateSecretIdentifier(ctx, id.Key, id.Owner, id.Namespace); err != nil {
			return fmt.Errorf("DeleteSecrets request contains invalid secret identifier: %w", err)
		}
		_, ok := idSet[vaulttypes.KeyFor(id)]
		if ok {
			return fmt.Errorf("DeleteSecrets requests cannot contain duplicate request for a given secret identifier: %s", id)
		}

		idSet[vaulttypes.KeyFor(id)] = true
	}

	for _, r := range o.GetDeleteSecretsResponse().Responses {
		if r.Id == nil {
			return errors.New("DeleteSecrets response contains nil secret identifier")
		}
	}

	return nil
}

func (r *ReportingPlugin) validateListSecretIdentifiersObservation(ctx context.Context, o *vaultcommon.Observation) error {
	listReq := o.GetListSecretIdentifiersRequest()
	listResp := o.GetListSecretIdentifiersResponse()
	if listReq == nil || listResp == nil {
		return errors.New("ListSecretIdentifiers observation must have both request and response")
	}

	// Passing in owner as key since Validate requires a non-empty key but list secret doesn't have a key
	if err := r.validator.ValidateSecretIdentifier(ctx, listReq.Owner, listReq.Owner, listReq.Namespace); err != nil {
		return fmt.Errorf("ListSecretIdentifiers request contains invalid secret identifier: %w", err)
	}

	if listResp.Success {
		ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: listReq.Owner})
		if err := r.cfg.MaxSecretsPerOwner.Check(ctx, len(listResp.Identifiers)); err != nil {
			var errBoundLimited limits.ErrorBoundLimited[int]
			if errors.As(err, &errBoundLimited) {
				return fmt.Errorf("ListSecretIdentifiers response exceeds maximum number of secrets per owner (have=%d, limit=%d): %w", len(listResp.Identifiers), errBoundLimited.Limit, err)
			}
			return fmt.Errorf("failed to check max secrets per owner limit: %w", err)
		}
	}

	return nil
}
