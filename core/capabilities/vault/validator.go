package vault

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

var (
	isValidIDComponent = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
)

type RequestValidator struct {
	MaxRequestBatchSizeLimiter          limits.BoundLimiter[int]
	MaxCiphertextLengthLimiter          limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierKeyLengthLimiter       limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierOwnerLengthLimiter     limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierNamespaceLengthLimiter limits.BoundLimiter[pkgconfig.Size]
}

func (r *RequestValidator) ValidateCreateSecretsRequest(ctx context.Context, publicKey *tdh2easy.PublicKey, request *vaultcommon.CreateSecretsRequest, skipLabelValidation bool) error {
	return r.validateWriteRequest(ctx, publicKey, request.RequestId, request.OrgId, request.WorkflowOwner, request.EncryptedSecrets, skipLabelValidation)
}

func (r *RequestValidator) ValidateUpdateSecretsRequest(ctx context.Context, publicKey *tdh2easy.PublicKey, request *vaultcommon.UpdateSecretsRequest, skipLabelValidation bool) error {
	return r.validateWriteRequest(ctx, publicKey, request.RequestId, request.OrgId, request.WorkflowOwner, request.EncryptedSecrets, skipLabelValidation)
}

// validateWriteRequest performs common validation for CreateSecrets and UpdateSecrets requests.
// It treats publicKey as optional, since it can be nil if the gateway nodes don't have the public key cached yet.
func (r *RequestValidator) validateWriteRequest(ctx context.Context, publicKey *tdh2easy.PublicKey, id string, orgID string, workflowOwner string, encryptedSecrets []*vaultcommon.EncryptedSecret, skipLabelValidation bool) error {
	if id == "" {
		return errors.New("request ID must not be empty")
	}
	if err := r.MaxRequestBatchSizeLimiter.Check(ctx, len(encryptedSecrets)); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[int]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("request batch size exceeds maximum of %d", errBoundLimited.Limit)
		}
		return fmt.Errorf("failed to check request batch size limit: %w", err)
	}
	if len(encryptedSecrets) == 0 {
		return errors.New("request batch must contain at least 1 item")
	}

	uniqueIDs := map[string]bool{}
	for idx, req := range encryptedSecrets {
		if req == nil {
			return errors.New("encrypted secret must not be nil at index " + strconv.Itoa(idx))
		}
		if req.Id == nil {
			return errors.New("secret ID must not be nil at index " + strconv.Itoa(idx))
		}

		if req.EncryptedValue == "" {
			return errors.New("secret must have encrypted value set at index " + strconv.Itoa(idx) + ":" + req.Id.String())
		}

		req.Id.Namespace = vaulttypes.NormalizeNamespace(req.Id.Namespace)

		if err := r.ValidateSecretIdentifier(ctx, req.Id.Key, req.Id.Owner, req.Id.Namespace); err != nil {
			return fmt.Errorf("invalid secret identifier at index %d: %w", idx, err)
		}
		if err := r.ValidateCiphertextSize(ctx, req.Id.Owner, req.EncryptedValue); err != nil {
			return fmt.Errorf("secret encrypted value at index %d is invalid: %w", idx, err)
		}
		if skipLabelValidation {
			if _, err := verifyEncryptedSecret(publicKey, req.EncryptedValue); err != nil {
				return errors.New("Encrypted Secret at index [" + strconv.Itoa(idx) + "] is invalid. Error: " + err.Error())
			}
		} else {
			expectedWorkflowOwner := workflowOwner
			if expectedWorkflowOwner == "" && orgID == "" {
				expectedWorkflowOwner = req.Id.Owner
			}
			err := EnsureRightLabelOnSecret(publicKey, req.EncryptedValue, expectedWorkflowOwner, orgID)
			if err != nil {
				return errors.New("Encrypted Secret at index [" + strconv.Itoa(idx) + "] doesn't have owner as the label. Error: " + err.Error())
			}
		}
		_, ok := uniqueIDs[vaulttypes.KeyFor(req.Id)]
		if ok {
			return errors.New("duplicate secret ID found at index " + strconv.Itoa(idx) + ": " + req.Id.String())
		}

		uniqueIDs[vaulttypes.KeyFor(req.Id)] = true
	}

	return nil
}

func (r *RequestValidator) ValidateCiphertextSize(ctx context.Context, owner string, encryptedValue string) error {
	rawCiphertext, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted value: %w", err)
	}
	innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: owner})
	if err := r.MaxCiphertextLengthLimiter.Check(innerCtx, pkgconfig.Size(len(rawCiphertext))*pkgconfig.Byte); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("ciphertext size exceeds maximum allowed size: %s", errBoundLimited.Limit)
		}
		return fmt.Errorf("failed to check ciphertext size limit: %w", err)
	}
	return nil
}

func (r *RequestValidator) ValidateSecretIdentifier(ctx context.Context, idKey string, idOwner string, idNamespace string) error {
	if idKey == "" {
		return errors.New("key cannot be empty")
	}
	if idOwner == "" {
		return errors.New("owner cannot be empty")
	}
	if idNamespace == "" {
		return errors.New("namespace cannot be empty")
	}

	if !isValidIDComponent(idKey) || !isValidIDComponent(idOwner) || !isValidIDComponent(idNamespace) {
		return errors.New("key, owner and namespace must only contain alphanumeric characters")
	}

	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: idOwner})
	if err := r.MaxIdentifierOwnerLengthLimiter.Check(ctx, pkgconfig.Size(len(idOwner))); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("owner exceeds maximum length of %s", errBoundLimited.Limit)
		}
		return fmt.Errorf("failed to check owner length limit: %w", err)
	}

	if err := r.MaxIdentifierNamespaceLengthLimiter.Check(ctx, pkgconfig.Size(len(idNamespace))); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("namespace exceeds maximum length of %s", errBoundLimited.Limit)
		}
		return fmt.Errorf("failed to check namespace length limit: %w", err)
	}

	if err := r.MaxIdentifierKeyLengthLimiter.Check(ctx, pkgconfig.Size(len(idKey))); err != nil {
		var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
		if errors.As(err, &errBoundLimited) {
			return fmt.Errorf("key exceeds maximum length of %s", errBoundLimited.Limit)
		}
		return fmt.Errorf("failed to check key length limit: %w", err)
	}

	return nil
}

func (r *RequestValidator) ValidateGetSecretsRequest(ctx context.Context, request *vaultcommon.GetSecretsRequest) error {
	if len(request.Requests) == 0 {
		return errors.New("no GetSecret request specified in request")
	}
	if len(request.Requests) >= vaulttypes.MaxBatchSize {
		return fmt.Errorf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize)
	}

	for idx, req := range request.Requests {
		if req.Id == nil {
			return errors.New("secret ID must have id set at index " + strconv.Itoa(idx))
		}
		if req.Id.Key == "" {
			return errors.New("secret ID must have key set at index " + strconv.Itoa(idx) + ": " + req.Id.String())
		}
		req.Id.Namespace = vaulttypes.NormalizeNamespace(req.Id.Namespace)
		if err := r.ValidateSecretIdentifier(ctx, req.Id.Key, req.Id.Owner, req.Id.Namespace); err != nil {
			return fmt.Errorf("invalid secret identifier at index %d: %w", idx, err)
		}
	}

	return nil
}

func (r *RequestValidator) ValidateListSecretIdentifiersRequest(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) error {
	request.Namespace = vaulttypes.NormalizeNamespace(request.Namespace)
	if request.RequestId == "" || request.Owner == "" {
		return errors.New("requestID or owner must not be empty")
	}
	if err := r.ValidateSecretIdentifier(ctx, request.Owner, request.Owner, request.Namespace); err != nil {
		return fmt.Errorf("invalid secret identifier: %w", err)
	}
	return nil
}

func (r *RequestValidator) ValidateDeleteSecretsRequest(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}
	if len(request.Ids) >= vaulttypes.MaxBatchSize {
		return errors.New("request batch size exceeds maximum of " + strconv.Itoa(vaulttypes.MaxBatchSize))
	}

	uniqueIDs := map[string]bool{}
	for idx, id := range request.Ids {
		if id == nil {
			return errors.New("secret ID must not be nil at index " + strconv.Itoa(idx))
		}
		id.Namespace = vaulttypes.NormalizeNamespace(id.Namespace)
		if err := r.ValidateSecretIdentifier(ctx, id.Key, id.Owner, id.Namespace); err != nil {
			return fmt.Errorf("invalid secret identifier at index %d: %w", idx, err)
		}

		_, ok := uniqueIDs[vaulttypes.KeyFor(id)]
		if ok {
			return errors.New("duplicate secret ID found at index " + strconv.Itoa(idx) + ": " + id.String())
		}

		uniqueIDs[vaulttypes.KeyFor(id)] = true
	}
	return nil
}

func NewRequestValidator(
	maxRequestBatchSizeLimiter limits.BoundLimiter[int],
	maxCiphertextLengthLimiter limits.BoundLimiter[pkgconfig.Size],
	maxIdentifierKeyLengthLimiter limits.BoundLimiter[pkgconfig.Size],
	maxIdentifierOwnerLengthLimiter limits.BoundLimiter[pkgconfig.Size],
	maxIdentifierNamespaceLengthLimiter limits.BoundLimiter[pkgconfig.Size],
) *RequestValidator {
	return &RequestValidator{
		MaxRequestBatchSizeLimiter:          maxRequestBatchSizeLimiter,
		MaxCiphertextLengthLimiter:          maxCiphertextLengthLimiter,
		MaxIdentifierKeyLengthLimiter:       maxIdentifierKeyLengthLimiter,
		MaxIdentifierOwnerLengthLimiter:     maxIdentifierOwnerLengthLimiter,
		MaxIdentifierNamespaceLengthLimiter: maxIdentifierNamespaceLengthLimiter,
	}
}

// EnsureRightLabelOnSecret verifies that the TDH2 ciphertext label matches either the
// workflowOwner (Ethereum address, left-padded) or the orgID (SHA256 hash). Either
// parameter can be empty to skip that check. The function succeeds if the label matches
// at least one non-empty owner.
func EnsureRightLabelOnSecret(publicKey *tdh2easy.PublicKey, secret string, workflowOwner string, orgID string) error {
	cipherText, err := verifyEncryptedSecret(publicKey, secret)
	if err != nil {
		return err
	}
	if cipherText == nil {
		return nil
	}
	secretLabel := cipherText.Label()
	expectedLabels := make([]string, 0, 2)

	if workflowOwner != "" {
		expected := vaultutils.WorkflowOwnerToLabel(workflowOwner)
		expectedLabels = append(expectedLabels, hex.EncodeToString(expected[:]))
		if secretLabel == expected {
			return nil
		}
	}

	if orgID != "" {
		expected := vaultutils.OrgIDToLabel(orgID)
		expectedLabels = append(expectedLabels, hex.EncodeToString(expected[:]))
		if secretLabel == expected {
			return nil
		}
	}

	return errors.New("secret label [" + hex.EncodeToString(secretLabel[:]) + "] does not match any of the provided owner labels; expectedLabels=[" + strings.Join(expectedLabels, ", ") + "]")
}

func verifyEncryptedSecret(publicKey *tdh2easy.PublicKey, secret string) (*tdh2easy.Ciphertext, error) {
	cipherBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, errors.New("failed to decode encrypted value:" + err.Error())
	}
	if publicKey == nil {
		// Public key can be nil if gateway cache isn't populated yet (immediately after gateway reboots).
		// Ok to not validate in such cases, since this validation also runs on Vault Nodes.
		return nil, nil
	}

	cipherText := &tdh2easy.Ciphertext{}
	if err := cipherText.UnmarshalVerify(cipherBytes, publicKey); err != nil {
		return nil, errors.New("failed to verify encrypted value:" + err.Error())
	}
	return cipherText, nil
}
