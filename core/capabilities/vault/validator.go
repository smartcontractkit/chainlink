package vault

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

var isValidIDComponent = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString

type RequestValidator struct {
	MaxRequestBatchSizeLimiter          limits.BoundLimiter[int]
	MaxCiphertextLengthLimiter          limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierKeyLengthLimiter       limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierOwnerLengthLimiter     limits.BoundLimiter[pkgconfig.Size]
	MaxIdentifierNamespaceLengthLimiter limits.BoundLimiter[pkgconfig.Size]
}

// SecretIdentifierLimits carries resolved identifier length limits for validation.
// When nil, RequestValidator falls back to its configured limiters.
type SecretIdentifierLimits struct {
	MaxOwnerLength     pkgconfig.Size
	MaxNamespaceLength pkgconfig.Size
	MaxKeyLength       pkgconfig.Size
}

func (r *RequestValidator) ValidateCreateSecretsRequest(ctx context.Context, publicKey *tdh2easy.PublicKey, request *vaultcommon.CreateSecretsRequest, skipLabelValidation bool) error {
	return r.validateWriteRequest(ctx, publicKey, request.RequestId, request.EncryptedSecrets, skipLabelValidation, true)
}

func (r *RequestValidator) ValidateUpdateSecretsRequest(ctx context.Context, publicKey *tdh2easy.PublicKey, request *vaultcommon.UpdateSecretsRequest, skipLabelValidation bool) error {
	return r.validateWriteRequest(ctx, publicKey, request.RequestId, request.EncryptedSecrets, skipLabelValidation, true)
}

// ValidateEncryptedSecretsStructure calls validateWriteRequest without the
// owner-scoped ciphertext-size limit, which must be checked separately after
// authorization via ValidateCiphertextSizes.
func (r *RequestValidator) ValidateEncryptedSecretsStructure(ctx context.Context, publicKey *tdh2easy.PublicKey, requestID string, encryptedSecrets []*vaultcommon.EncryptedSecret, skipLabelValidation bool) error {
	return r.validateWriteRequest(ctx, publicKey, requestID, encryptedSecrets, skipLabelValidation, false)
}

// validateWriteRequest performs common validation for CreateSecrets and UpdateSecrets requests.
// It treats publicKey as optional, since it can be nil if the gateway nodes don't have the public key cached yet.
// includeCiphertextSize controls the owner-scoped ciphertext-size check, which must be
// skipped before authorization (see ValidateEncryptedSecretsStructure).
func (r *RequestValidator) validateWriteRequest(ctx context.Context, publicKey *tdh2easy.PublicKey, id string, encryptedSecrets []*vaultcommon.EncryptedSecret, skipLabelValidation bool, includeCiphertextSize bool) error {
	if id == "" {
		return errors.New("request ID must not be empty")
	}
	if err := r.MaxRequestBatchSizeLimiter.Check(ctx, len(encryptedSecrets)); err != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[int]](err); ok {
			return fmt.Errorf("request batch size exceeds maximum of %d: %w", errBoundLimited.Limit, err)
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

		if err := r.ValidateSecretIdentifier(ctx, req.Id.Key, req.Id.Owner, req.Id.Namespace, nil); err != nil {
			return fmt.Errorf("invalid secret identifier at index %d: %w", idx, err)
		}
		if includeCiphertextSize {
			if err := r.ValidateCiphertextSize(ctx, req.Id.Owner, req.EncryptedValue); err != nil {
				return fmt.Errorf("secret encrypted value at index %d is invalid: %w", idx, err)
			}
		}
		if skipLabelValidation {
			if _, err := verifyEncryptedSecret(publicKey, req.EncryptedValue); err != nil {
				return errors.New("Encrypted Secret at index [" + strconv.Itoa(idx) + "] is invalid. Error: " + err.Error())
			}
		} else {
			err := EnsureRightLabelOnSecret(publicKey, req.EncryptedValue, req.Id.Owner)
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

func (r *RequestValidator) ValidateCiphertextSize(ctx context.Context, owner, encryptedValue string) error {
	rawCiphertext, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted value: %w", err)
	}
	// TODO orgID https://smartcontract-it.atlassian.net/browse/CRE-1707
	innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: owner})
	if err := r.MaxCiphertextLengthLimiter.Check(innerCtx, pkgconfig.Size(len(rawCiphertext))*pkgconfig.Byte); err != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[pkgconfig.Size]](err); ok {
			return fmt.Errorf("ciphertext size exceeds maximum allowed size: %s: %w", errBoundLimited.Limit, err)
		}
		return fmt.Errorf("failed to check ciphertext size limit: %w", err)
	}
	return nil
}

// ValidateCiphertextSizes checks the owner-scoped ciphertext-size limit for each
// encrypted secret in a write request that already passed structure validation
// (ValidateEncryptedSecretsStructure). It must only be called after
// authorization, with the authorized workflow owner: checking the scoped
// limiter registers a per-owner tenant that spawns a persistent background
// updater, so running it pre-auth would let unauthenticated callers create
// unbounded limiter tenants.
func (r *RequestValidator) ValidateCiphertextSizes(ctx context.Context, owner string, encryptedSecrets []*vaultcommon.EncryptedSecret) error {
	for idx, secret := range encryptedSecrets {
		if secret == nil {
			return errors.New("encrypted secret must not be nil at index " + strconv.Itoa(idx))
		}
		if err := r.ValidateCiphertextSize(ctx, owner, secret.EncryptedValue); err != nil {
			return fmt.Errorf("secret encrypted value at index %d is invalid: %w", idx, err)
		}
	}
	return nil
}

func (r *RequestValidator) ValidateSecretIdentifier(ctx context.Context, idKey, idOwner, idNamespace string, identifierLimits *SecretIdentifierLimits) error {
	if idKey == "" {
		return errors.New("key cannot be empty")
	}
	if idOwner == "" {
		return errors.New("owner cannot be empty")
	}

	if !isValidIDComponent(idKey) || !isValidIDComponent(idOwner) || (idNamespace != "" && !isValidIDComponent(idNamespace)) {
		return errors.New("key, owner and namespace must only contain alphanumeric characters")
	}

	if identifierLimits != nil {
		if err := checkIdentifierComponentLength("owner", idOwner, identifierLimits.MaxOwnerLength); err != nil {
			return err
		}
		if err := checkIdentifierComponentLength("namespace", idNamespace, identifierLimits.MaxNamespaceLength); err != nil {
			return err
		}
		return checkIdentifierComponentLength("key", idKey, identifierLimits.MaxKeyLength)
	}

	// TODO orgID https://smartcontract-it.atlassian.net/browse/CRE-1707
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: idOwner})
	if err := r.MaxIdentifierOwnerLengthLimiter.Check(ctx, pkgconfig.Size(len(idOwner))); err != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[pkgconfig.Size]](err); ok {
			return fmt.Errorf("owner exceeds maximum length of %s: %w", errBoundLimited.Limit, err)
		}
		return fmt.Errorf("failed to check owner length limit: %w", err)
	}

	if err := r.MaxIdentifierNamespaceLengthLimiter.Check(ctx, pkgconfig.Size(len(idNamespace))); err != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[pkgconfig.Size]](err); ok {
			return fmt.Errorf("namespace exceeds maximum length of %s: %w", errBoundLimited.Limit, err)
		}
		return fmt.Errorf("failed to check namespace length limit: %w", err)
	}

	if err := r.MaxIdentifierKeyLengthLimiter.Check(ctx, pkgconfig.Size(len(idKey))); err != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[pkgconfig.Size]](err); ok {
			return fmt.Errorf("key exceeds maximum length of %s: %w", errBoundLimited.Limit, err)
		}
		return fmt.Errorf("failed to check key length limit: %w", err)
	}

	return nil
}

// EffectiveSecretIdentifierLimits returns donLimits adjusted for the given
// owner: a component is raised only when the owner's resolved per-owner limit
// exceeds the configured default (an explicit privileged override), so
// DON-wide settings remain the authoritative baseline while privileged owners
// keep their allowances. Per-owner limits at or below the default are
// superseded by the DON baseline.
func (r *RequestValidator) EffectiveSecretIdentifierLimits(ctx context.Context, owner string, donLimits SecretIdentifierLimits) (SecretIdentifierLimits, error) {
	// TODO orgID https://smartcontract-it.atlassian.net/browse/CRE-1707
	ownerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: owner})
	applyOverride := func(limiter limits.BoundLimiter[pkgconfig.Size], base pkgconfig.Size) (pkgconfig.Size, error) {
		perOwner, err := limiter.Limit(ownerCtx)
		if err != nil {
			return 0, err
		}
		def, err := limiter.Limit(ctx)
		if err != nil {
			return 0, err
		}
		if perOwner > def {
			return max(base, perOwner), nil
		}
		return base, nil
	}

	out := donLimits
	var err error
	if out.MaxOwnerLength, err = applyOverride(r.MaxIdentifierOwnerLengthLimiter, donLimits.MaxOwnerLength); err != nil {
		return SecretIdentifierLimits{}, fmt.Errorf("failed to resolve owner length limit: %w", err)
	}
	if out.MaxNamespaceLength, err = applyOverride(r.MaxIdentifierNamespaceLengthLimiter, donLimits.MaxNamespaceLength); err != nil {
		return SecretIdentifierLimits{}, fmt.Errorf("failed to resolve namespace length limit: %w", err)
	}
	if out.MaxKeyLength, err = applyOverride(r.MaxIdentifierKeyLengthLimiter, donLimits.MaxKeyLength); err != nil {
		return SecretIdentifierLimits{}, fmt.Errorf("failed to resolve key length limit: %w", err)
	}
	return out, nil
}

func (r *RequestValidator) ValidateGetSecretsRequest(ctx context.Context, request *vaultcommon.GetSecretsRequest) error {
	if len(request.Requests) == 0 {
		return errors.New("no GetSecret request specified in request")
	}
	if len(request.Requests) >= vaulttypes.MaxBatchSize {
		return fmt.Errorf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize)
	}

	uniqueIDs := map[string]bool{}
	for idx, req := range request.Requests {
		if req.Id == nil {
			return errors.New("secret ID must have id set at index " + strconv.Itoa(idx))
		}
		if req.Id.Key == "" {
			return errors.New("secret ID must have key set at index " + strconv.Itoa(idx) + ": " + req.Id.String())
		}
		if err := r.ValidateSecretIdentifier(ctx, req.Id.Key, req.Id.Owner, req.Id.Namespace, nil); err != nil {
			return fmt.Errorf("invalid secret identifier at index %d: %w", idx, err)
		}

		_, ok := uniqueIDs[vaulttypes.KeyFor(req.Id)]
		if ok {
			return errors.New("duplicate secret ID found at index " + strconv.Itoa(idx) + ": " + req.Id.String())
		}

		uniqueIDs[vaulttypes.KeyFor(req.Id)] = true
	}

	return nil
}

func (r *RequestValidator) ValidateListSecretIdentifiersRequest(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) error {
	if request.RequestId == "" || request.Owner == "" {
		return errors.New("requestID or owner must not be empty")
	}
	if err := r.ValidateSecretIdentifier(ctx, request.Owner, request.Owner, request.Namespace, nil); err != nil {
		return fmt.Errorf("invalid secret identifier: %w", err)
	}
	return nil
}

func (r *RequestValidator) ValidateDeleteSecretsRequest(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}
	if err := r.MaxRequestBatchSizeLimiter.Check(ctx, len(request.Ids)); err != nil {
		if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[int]](err); ok {
			return fmt.Errorf("request batch size exceeds maximum of %d: %w", errBoundLimited.Limit, err)
		}
		return fmt.Errorf("failed to check request batch size limit: %w", err)
	}
	if len(request.Ids) == 0 {
		return errors.New("request batch must contain at least 1 item")
	}

	uniqueIDs := map[string]bool{}
	for idx, id := range request.Ids {
		if id == nil {
			return errors.New("secret ID must not be nil at index " + strconv.Itoa(idx))
		}
		if err := r.ValidateSecretIdentifier(ctx, id.Key, id.Owner, id.Namespace, nil); err != nil {
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

func checkIdentifierComponentLength(component, value string, limit pkgconfig.Size) error {
	if pkgconfig.Size(len(value)) > limit {
		return fmt.Errorf("%s exceeds maximum length of %s: %w", component, limit, limits.ErrorBoundLimited[pkgconfig.Size]{Limit: limit, Amount: pkgconfig.Size(len(value))})
	}
	return nil
}

// CheckRequestBatchSize validates a request batch size. When maxBatchSize is
// non-nil (e.g. a DON-wide agreed limit), it takes precedence over the
// configured limiter.
func (r *RequestValidator) CheckRequestBatchSize(ctx context.Context, batchSize int, maxBatchSize *int) error {
	if err := r.checkRequestBatchSize(ctx, batchSize, maxBatchSize); err != nil {
		if _, ok := errors.AsType[limits.ErrorBoundLimited[int]](err); ok {
			return fmt.Errorf("max batch size exceeded for request: %w", err)
		}
		return fmt.Errorf("failed to check batch size: %w", err)
	}
	return nil
}

func (r *RequestValidator) checkRequestBatchSize(ctx context.Context, batchSize int, maxBatchSize *int) error {
	if maxBatchSize != nil {
		if batchSize > *maxBatchSize {
			return limits.ErrorBoundLimited[int]{Limit: *maxBatchSize, Amount: batchSize}
		}
		return nil
	}
	return r.MaxRequestBatchSizeLimiter.Check(ctx, batchSize)
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

// Close releases validator limiter resources.
func (r *RequestValidator) Close() error {
	return errors.Join(
		r.MaxRequestBatchSizeLimiter.Close(),
		r.MaxCiphertextLengthLimiter.Close(),
		r.MaxIdentifierKeyLengthLimiter.Close(),
		r.MaxIdentifierOwnerLengthLimiter.Close(),
		r.MaxIdentifierNamespaceLengthLimiter.Close(),
	)
}

// NewRequestValidatorFromLimitsFactory constructs a RequestValidator from CRE limits settings.
func NewRequestValidatorFromLimitsFactory(limitsFactory limits.Factory) (*RequestValidator, error) {
	limiter, err := limits.MakeUpperBoundLimiter(limitsFactory, cresettings.Default.VaultRequestBatchSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("could not create request batch size limiter: %w", err)
	}
	ciphertextLimiter, err := limits.MakeUpperBoundLimiter(limitsFactory, cresettings.Default.PerOwner.VaultCiphertextSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("could not create ciphertext size limiter: %w", err)
	}
	idKeyLengthLimiter, err := limits.MakeUpperBoundLimiter(limitsFactory, cresettings.Default.VaultIdentifierKeySizeLimit)
	if err != nil {
		return nil, fmt.Errorf("could not create identifier key size limiter: %w", err)
	}
	idOwnerLengthLimiter, err := limits.MakeUpperBoundLimiter(limitsFactory, cresettings.Default.VaultIdentifierOwnerSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("could not create identifier owner size limiter: %w", err)
	}
	idNamespaceLengthLimiter, err := limits.MakeUpperBoundLimiter(limitsFactory, cresettings.Default.VaultIdentifierNamespaceSizeLimit)
	if err != nil {
		return nil, fmt.Errorf("could not create identifier namespace size limiter: %w", err)
	}
	validator := NewRequestValidator(limiter, ciphertextLimiter, idKeyLengthLimiter, idOwnerLengthLimiter, idNamespaceLengthLimiter)
	return validator, nil
}

// EnsureRightLabelOnSecret verifies that the TDH2 ciphertext label matches the workflow
// owner label (Ethereum address, left-padded to 32 bytes). owner must be non-empty;
// when the public key is nil, verification is skipped for the same reasons as
// verifyEncryptedSecret.
func EnsureRightLabelOnSecret(publicKey *tdh2easy.PublicKey, secret, owner string) error {
	cipherText, err := verifyEncryptedSecret(publicKey, secret)
	if err != nil {
		return err
	}
	if cipherText == nil {
		return nil
	}
	if owner == "" {
		return errors.New("owner must not be empty for secret label verification")
	}

	expected := vaultutils.WorkflowOwnerToLabel(owner)
	secretLabel := cipherText.Label()
	if secretLabel == expected {
		return nil
	}

	return fmt.Errorf("secret label [%s] does not match workflow owner label [%s]",
		hex.EncodeToString(secretLabel[:]), hex.EncodeToString(expected[:]))
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
