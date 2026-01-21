package vault

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

type RequestValidator struct {
	MaxRequestBatchSizeLimiter limits.BoundLimiter[int]
}

func (r *RequestValidator) ValidateCreateSecretsRequest(publicKey *tdh2easy.PublicKey, request *vaultcommon.CreateSecretsRequest) error {
	return r.validateWriteRequest(publicKey, request.RequestId, request.EncryptedSecrets)
}

func (r *RequestValidator) ValidateUpdateSecretsRequest(publicKey *tdh2easy.PublicKey, request *vaultcommon.UpdateSecretsRequest) error {
	return r.validateWriteRequest(publicKey, request.RequestId, request.EncryptedSecrets)
}

// validateWriteRequest performs common validation for CreateSecrets and UpdateSecrets requests
// It treats publicKey as optional, since it can be nil if the gateway nodes don't have the public key cached yet
func (r *RequestValidator) validateWriteRequest(publicKey *tdh2easy.PublicKey, id string, encryptedSecrets []*vaultcommon.EncryptedSecret) error {
	if id == "" {
		return errors.New("request ID must not be empty")
	}
	if err := r.MaxRequestBatchSizeLimiter.Check(context.Background(), len(encryptedSecrets)); err != nil {
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

		if req.Id.Key == "" || req.Id.Namespace == "" || req.Id.Owner == "" {
			return errors.New("secret ID must have key, namespace and owner set at index " + strconv.Itoa(idx) + ":" + req.Id.String())
		}

		if req.EncryptedValue == "" {
			return errors.New("secret must have encrypted value set at index " + strconv.Itoa(idx) + ":" + req.Id.String())
		}
		err := EnsureRightLabelOnSecret(publicKey, req.EncryptedValue, req.Id.Owner)
		if err != nil {
			return errors.New("Encrypted Secret at index [" + strconv.Itoa(idx) + "] doesn't have owner as the label. Error: " + err.Error())
		}
		_, ok := uniqueIDs[vaulttypes.KeyFor(req.Id)]
		if ok {
			return errors.New("duplicate secret ID found at index " + strconv.Itoa(idx) + ": " + req.Id.String())
		}

		uniqueIDs[vaulttypes.KeyFor(req.Id)] = true
	}

	return nil
}

func (r *RequestValidator) ValidateGetSecretsRequest(request *vaultcommon.GetSecretsRequest) error {
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
	}

	return nil
}

func (r *RequestValidator) ValidateListSecretIdentifiersRequest(request *vaultcommon.ListSecretIdentifiersRequest) error {
	if request.RequestId == "" || request.Owner == "" || request.Namespace == "" {
		return errors.New("requestID, owner or namespace must not be empty")
	}
	return nil
}

func (r *RequestValidator) ValidateDeleteSecretsRequest(request *vaultcommon.DeleteSecretsRequest) error {
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
		if id.Key == "" || id.Namespace == "" || id.Owner == "" {
			return errors.New("secret ID must have key, namespace and owner set at index " + strconv.Itoa(idx) + ": " + id.String())
		}

		_, ok := uniqueIDs[vaulttypes.KeyFor(id)]
		if ok {
			return errors.New("duplicate secret ID found at index " + strconv.Itoa(idx) + ": " + id.String())
		}

		uniqueIDs[vaulttypes.KeyFor(id)] = true
	}
	return nil
}

func NewRequestValidator(maxRequestBatchSizeLimiter limits.BoundLimiter[int]) *RequestValidator {
	return &RequestValidator{
		MaxRequestBatchSizeLimiter: maxRequestBatchSizeLimiter,
	}
}

func EnsureRightLabelOnSecret(publicKey *tdh2easy.PublicKey, secret, owner string) error {
	cipherText := &tdh2easy.Ciphertext{}
	cipherBytes, err := hex.DecodeString(secret)
	if err != nil {
		return errors.New("failed to decode encrypted value:" + err.Error())
	}
	if publicKey == nil {
		// Public key can be nil if gateway cache isn't populated yet(immediately after gateway reboots)
		// Ok to not validate in such cases, since this validation also runs on Vault Nodes
		return nil
	}
	err = cipherText.UnmarshalVerify(cipherBytes, publicKey)
	if err != nil {
		return errors.New("failed to verify encrypted value:" + err.Error())
	}
	secretLabel := cipherText.Label()
	ownerAddr := common.HexToAddress(owner)
	var ownerLabel [32]byte
	copy(ownerLabel[12:], ownerAddr.Bytes()) // left-pad with 12 zero
	if secretLabel != ownerLabel {
		return errors.New("secret label [" + hex.EncodeToString(secretLabel[:]) + "] does not match owner label [" + hex.EncodeToString(ownerLabel[:]) + "]")
	}
	return nil
}
