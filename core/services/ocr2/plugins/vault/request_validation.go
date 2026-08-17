package vault

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func secretIdentifierKey(id *vaultcommon.SecretIdentifier) string {
	if id == nil {
		return "<nil>"
	}
	return vaulttypes.KeyFor(id)
}

func buildSecretIdentifierCounts(ids []*vaultcommon.SecretIdentifier) map[string]int {
	counts := make(map[string]int, len(ids))
	for _, id := range ids {
		counts[secretIdentifierKey(id)]++
	}
	return counts
}

func buildEncryptedSecretIdentifierCounts(secrets []*vaultcommon.EncryptedSecret) map[string]int {
	counts := make(map[string]int, len(secrets))
	for _, s := range secrets {
		var key string
		if s.Id == nil {
			key = "<nil>"
		} else {
			key = vaulttypes.KeyFor(s.Id)
		}
		counts[key]++
	}
	return counts
}

func buildGetSecretsRequestIdentifierCounts(requests []*vaultcommon.SecretRequest) map[string]int {
	counts := make(map[string]int, len(requests))
	for _, sr := range requests {
		counts[secretIdentifierKey(sr.Id)]++
	}
	return counts
}

func validateRequestResponseItemCount(requestCount, responseCount int, method string) error {
	if requestCount != responseCount {
		return fmt.Errorf("%s request and response must have the same number of items", method)
	}
	return nil
}

func (r *ReportingPlugin) validateSecretIdentifier(ctx context.Context, id *vaultcommon.SecretIdentifier) (*vaultcommon.SecretIdentifier, error) {
	if id == nil {
		return nil, newUserError("secret identifier cannot be nil")
	}

	namespace := id.Namespace
	if namespace == "" {
		namespace = vaulttypes.DefaultNamespace
	}

	if err := r.validator.ValidateSecretIdentifier(ctx, id.Key, id.Owner, namespace); err != nil {
		return nil, newUserError(err.Error())
	}

	return &vaultcommon.SecretIdentifier{
		Key:       id.Key,
		Owner:     id.Owner,
		Namespace: namespace,
	}, nil
}

func (r *ReportingPlugin) validateDuplicateSecretIdentifierUserError(id *vaultcommon.SecretIdentifier, counts map[string]int) error {
	key := secretIdentifierKey(id)
	if counts[key] > 1 {
		return newUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}
	return nil
}

func (r *ReportingPlugin) validateEncryptedSecretCiphertextSize(ctx context.Context, owner string, encryptedValue string) error {
	if ierr := r.validator.ValidateCiphertextSize(ctx, owner, encryptedValue); ierr != nil {
		return newUserError(ierr.Error())
	}
	return nil
}

func (r *ReportingPlugin) validateEncryptedSecretLabel(owner string, encryptedValue string) error {
	if err := vaultcap.EnsureRightLabelOnSecret(r.cfg.PublicKey, encryptedValue, owner); err != nil {
		return newUserError("failed to verify ciphertext: " + err.Error())
	}
	return nil
}

// validateEncryptedSecretPayload checks identifier, batch dupes, ciphertext size, and TDH2 label.
// KV-dependent rules (key exists, owner secret count) are applied at StateTransition.
func (r *ReportingPlugin) validateEncryptedSecretPayload(
	ctx context.Context,
	secret *vaultcommon.EncryptedSecret,
	requestsCountForID map[string]int,
) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateSecretIdentifier(ctx, secret.Id)
	if err != nil {
		return id, err
	}

	if err := r.validateDuplicateSecretIdentifierUserError(secret.Id, requestsCountForID); err != nil {
		return id, err
	}

	if err := r.validateEncryptedSecretCiphertextSize(ctx, secret.Id.Owner, secret.EncryptedValue); err != nil {
		return id, err
	}

	if err := r.validateEncryptedSecretLabel(secret.Id.Owner, secret.EncryptedValue); err != nil {
		return id, err
	}

	// Other verifications, such as checking whether the key already exists,
	// or whether we have hit the limit on the number of secrets per owner,
	// are done in the StateTransition phase.
	// This guarantees that we correctly account for changes made in other requests
	// in the batch.
	return id, nil
}

func (r *ReportingPlugin) validateGetSecretsRequestItem(
	ctx context.Context,
	secretRequest *vaultcommon.SecretRequest,
	requestsCountForID map[string]int,
) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateSecretIdentifier(ctx, secretRequest.Id)
	if err != nil {
		return nil, err
	}

	if err := r.validateDuplicateSecretIdentifierUserError(secretRequest.Id, requestsCountForID); err != nil {
		return nil, err
	}

	return id, nil
}

func (r *ReportingPlugin) validateDeleteSecretsRequestItem(
	ctx context.Context,
	identifier *vaultcommon.SecretIdentifier,
	requestsCountForID map[string]int,
) (*vaultcommon.SecretIdentifier, error) {
	id, err := r.validateSecretIdentifier(ctx, identifier)
	if err != nil {
		return id, err
	}

	if err := r.validateDuplicateSecretIdentifierUserError(identifier, requestsCountForID); err != nil {
		return id, err
	}

	return id, nil
}

func (r *ReportingPlugin) validateGetSecretsRequestPayload(ctx context.Context, req *vaultcommon.GetSecretsRequest) error {
	if err := r.validator.CheckRequestBatchSize(ctx, len(req.Requests)); err != nil {
		return err
	}

	counts := buildGetSecretsRequestIdentifierCounts(req.Requests)
	for _, sr := range req.Requests {
		if sr.Id == nil {
			return errors.New("GetSecrets request contains nil secret identifier")
		}
		if _, err := r.validateGetSecretsRequestItem(ctx, sr, counts); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReportingPlugin) validateEncryptedSecretsRequestWire(
	ctx context.Context,
	secrets []*vaultcommon.EncryptedSecret,
	method string,
) error {
	if err := r.validator.CheckRequestBatchSize(ctx, len(secrets)); err != nil {
		return err
	}

	// We disallow duplicate create/update requests within a single batch request.
	// This prevents users from clobbering their own writes.
	idSet := map[string]bool{}
	for _, s := range secrets {
		if s.Id == nil {
			return fmt.Errorf("%s request contains nil secret identifier", method)
		}
		if err := r.validator.ValidateSecretIdentifier(ctx, s.Id.Key, s.Id.Owner, s.Id.Namespace); err != nil {
			return fmt.Errorf("%s request contains invalid secret identifier: %w", method, err)
		}
		key := vaulttypes.KeyFor(s.Id)
		if idSet[key] {
			return fmt.Errorf("%s requests cannot contain duplicate request for a given secret identifier: %s", method, s.Id)
		}
		idSet[key] = true

		if err := r.validator.ValidateCiphertextSize(ctx, s.Id.Owner, s.EncryptedValue); err != nil {
			return fmt.Errorf("%s request: %w", method, err)
		}
	}
	return nil
}

func (r *ReportingPlugin) validateDeleteSecretsRequestWire(ctx context.Context, ids []*vaultcommon.SecretIdentifier) error {
	if err := r.validator.CheckRequestBatchSize(ctx, len(ids)); err != nil {
		return err
	}

	idSet := map[string]bool{}
	for _, id := range ids {
		if id == nil {
			return errors.New("DeleteSecrets request contains nil secret identifier")
		}
		if err := r.validator.ValidateSecretIdentifier(ctx, id.Key, id.Owner, id.Namespace); err != nil {
			return fmt.Errorf("DeleteSecrets request contains invalid secret identifier: %w", err)
		}
		key := vaulttypes.KeyFor(id)
		if idSet[key] {
			return fmt.Errorf("DeleteSecrets requests cannot contain duplicate request for a given secret identifier: %s", id)
		}
		idSet[key] = true
	}
	return nil
}

func (r *ReportingPlugin) validateDeleteSecretsRequestPayload(ctx context.Context, req *vaultcommon.DeleteSecretsRequest) error {
	if err := r.validator.CheckRequestBatchSize(ctx, len(req.Ids)); err != nil {
		return err
	}

	counts := buildSecretIdentifierCounts(req.Ids)
	for _, id := range req.Ids {
		if id == nil {
			return errors.New("DeleteSecrets request contains nil secret identifier")
		}
		if _, err := r.validateDeleteSecretsRequestItem(ctx, id, counts); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReportingPlugin) validateListSecretIdentifiersOwnerNonempty(req *vaultcommon.ListSecretIdentifiersRequest) error {
	if req.Owner == "" {
		return errors.New("invalid request: owner cannot be empty")
	}
	return nil
}

func (r *ReportingPlugin) validateListSecretIdentifiersOwnerWire(ctx context.Context, req *vaultcommon.ListSecretIdentifiersRequest) error {
	return r.validator.ValidateSecretIdentifier(ctx, req.Owner, req.Owner, req.Namespace)
}

func (r *ReportingPlugin) validateListSecretIdentifiersResponseSize(ctx context.Context, owner string, identifierCount int) error {
	innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: owner})
	return r.cfg.MaxSecretsPerOwner.Check(innerCtx, identifierCount)
}

func decodeEncryptedSecretHex(encryptedValue string) ([]byte, error) {
	encryptedSecret, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return nil, newUserError("could not decode secret value: invalid hex: " + err.Error())
	}
	return encryptedSecret, nil
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
