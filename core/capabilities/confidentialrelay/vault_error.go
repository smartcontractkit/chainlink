package confidentialrelay

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// vaultSecretError wraps a per-secret error returned by the vault DON in a
// SecretResponse.Error field. The vault OCR plugin classifies user-caused
// failures (e.g. "key does not exist", "key already exists", invalid public
// key) as *vaulttypes.UserError and surfaces the raw message through
// userFacingError. System failures are replaced with a generic fallback
// (vaulttypes.SecretGetSystemErrorFallback) so their details do not leak.
//
// By the time the relay handler sees the string the Go type is gone (protobuf
// boundary), so translateVaultResponse re-wraps it here. Unwrap returns a
// *vaulttypes.UserError only for user-caused messages; for the system fallback
// it returns nil, so vaulttypes.IsUserError (errors.As) reports false and the
// handler keeps the failure classified as jsonrpc.ErrInternal. The handler then
// maps user errors to jsonrpc.ErrInvalidParams so the enclave receives the
// actual cause and the metrics/logs reflect a user error rather than an
// internal failure.
type vaultSecretError struct {
	namespace string
	key       string
	msg       string
}

func (e *vaultSecretError) Error() string {
	return fmt.Sprintf("vault error for secret %s/%s: %s", e.namespace, e.key, e.msg)
}

func (e *vaultSecretError) Unwrap() error {
	if vaulttypes.IsSecretGetSystemError(e.msg) {
		return nil
	}
	return vaulttypes.NewUserError(e.msg)
}
