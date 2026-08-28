package confidentialrelay

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// vaultSecretError wraps a per-secret error returned by the vault DON in a
// SecretResponse.Error field. The vault OCR plugin classifies user-caused
// failures (e.g. "key does not exist") as userError and surfaces the raw
// message through userFacingError. System failures are replaced with a generic
// fallback (vaulttypes.SecretGetSystemErrorFallback) so their details do not
// leak.
//
// By the time the relay handler sees the string the Go type is gone (protobuf
// boundary), so translateVaultResponse re-wraps it here with an explicit isUser
// flag set at construction time. The handler checks the flag via IsUserError to
// map user errors to jsonrpc.ErrInvalidParams and system errors to
// jsonrpc.ErrInternal.
type vaultSecretError struct {
	namespace string
	key       string
	msg       string
	isUser    bool
}

func (e *vaultSecretError) Error() string {
	return fmt.Sprintf("vault error for secret %s/%s: %s", e.namespace, e.key, e.msg)
}

// IsUserError reports whether err is a vaultSecretError marked as a user error.
func IsUserError(err error) bool {
	var vse *vaultSecretError
	if !errors.As(err, &vse) {
		return false
	}
	return vse.isUser
}

// newVaultSecretError creates a vaultSecretError, classifying the message as
// user or system based on whether it matches the vault plugin's system-error
// fallback.
func newVaultSecretError(namespace, key, msg string) *vaultSecretError {
	return &vaultSecretError{
		namespace: namespace,
		key:       key,
		msg:       msg,
		isUser:    !vaulttypes.IsSecretGetSystemError(msg),
	}
}
