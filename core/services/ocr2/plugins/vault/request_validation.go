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
		return nil, vaulttypes.NewUserError("secret identifier cannot be nil")
	}

	namespace := id.Namespace
	if namespace == "" {
		namespace = vaulttypes.DefaultNamespace
	}

	if err := r.checkSecretIdentifier(ctx, id.Key, id.Owner, namespace); err != nil {
		return nil, vaulttypes.NewUserError(err.Error())
	}

	return &vaultcommon.SecretIdentifier{
		Key:       id.Key,
		Owner:     id.Owner,
		Namespace: namespace,
	}, nil
}

// checkSecretIdentifier validates a secret identifier against the active
// DON-wide settings, falling back to node-local configuration when no DON
// settings have been committed. Privileged per-owner limits that exceed the
// configured default are still honored on top of the DON baseline.
func (r *ReportingPlugin) checkSecretIdentifier(ctx context.Context, idKey, idOwner, idNamespace string) error {
	donLimits, err := r.activeSettings.Load().secretIdentifierLimits(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve secret identifier limits: %w", err)
	}
	effective, err := r.validator.EffectiveSecretIdentifierLimits(ctx, idOwner, donLimits)
	if err != nil {
		return fmt.Errorf("failed to resolve secret identifier limits: %w", err)
	}
	return r.validator.ValidateSecretIdentifier(ctx, idKey, idOwner, idNamespace, &effective)
}

// checkRequestBatchSize validates a request batch size against the active
// DON-wide settings, falling back to node-local configuration when no DON
// settings have been committed.
func (r *ReportingPlugin) checkRequestBatchSize(ctx context.Context, batchSize int) error {
	maxBatch := r.activeSettings.Load().maxRequestBatchSize(ctx)
	return r.validator.CheckRequestBatchSize(ctx, batchSize, &maxBatch)
}

func (r *ReportingPlugin) validateDuplicateSecretIdentifierUserError(id *vaultcommon.SecretIdentifier, counts map[string]int) error {
	key := secretIdentifierKey(id)
	if counts[key] > 1 {
		return vaulttypes.NewUserError("duplicate request for secret identifier " + vaulttypes.KeyFor(id))
	}
	return nil
}

func (r *ReportingPlugin) validateEncryptedSecretCiphertextSize(ctx context.Context, owner string, encryptedValue string) error {
	if ierr := r.validator.ValidateCiphertextSize(ctx, owner, encryptedValue); ierr != nil {
		return vaulttypes.NewUserError(ierr.Error())
	}
	return nil
}

func (r *ReportingPlugin) validateEncryptedSecretLabel(owner string, encryptedValue string) error {
	if err := vaultcap.EnsureRightLabelOnSecret(r.cfg.PublicKey, encryptedValue, owner); err != nil {
		return vaulttypes.NewUserError("failed to verify ciphertext: " + err.Error())
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
	if err := r.checkRequestBatchSize(ctx, len(req.Requests)); err != nil {
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

func (r *ReportingPlugin) validateDeleteSecretsRequestPayload(ctx context.Context, req *vaultcommon.DeleteSecretsRequest) error {
	if err := r.checkRequestBatchSize(ctx, len(req.Ids)); err != nil {
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
	return r.checkSecretIdentifier(ctx, req.Owner, req.Owner, req.Namespace)
}

func (r *ReportingPlugin) validateListSecretIdentifiersResponseSize(ctx context.Context, owner string, identifierCount int) error {
	innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: owner})
	return r.cfg.MaxSecretsPerOwner.Check(innerCtx, identifierCount)
}

func decodeEncryptedSecretHex(encryptedValue string) ([]byte, error) {
	encryptedSecret, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return nil, vaulttypes.NewUserError("could not decode secret value: invalid hex: " + err.Error())
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
