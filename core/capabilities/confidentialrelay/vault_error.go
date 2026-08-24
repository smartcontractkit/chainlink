package confidentialrelay

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// vaultSecretError wraps a per-secret error returned by the vault DON in a
// SecretResponse.Error field. The vault OCR plugin classifies user-caused
// failures (e.g. "key does not exist", "key already exists", invalid public
// key) as *vaulttypes.UserError and surfaces the raw message through
// userFacingError. By the time the relay handler sees the string the Go type
// is gone (protobuf boundary), so translateVaultResponse re-wraps it in this
// type to carry the vault plugin's classification across.
//
// Unwrap returns a *vaulttypes.UserError so vaulttypes.IsUserError (errors.As)
// can detect it through this wrapper. The handler then maps it to
// jsonrpc.ErrInvalidParams (a client error) instead of ErrInternal, so the
// enclave receives the actual cause and the metrics/logs reflect a user error
// rather than an internal failure.
type vaultSecretError struct {
	namespace string
	key       string
	msg       string
}

func (e *vaultSecretError) Error() string {
	return fmt.Sprintf("vault error for secret %s/%s: %s", e.namespace, e.key, e.msg)
}

func (e *vaultSecretError) Unwrap() error {
	return vaulttypes.NewUserError(e.msg)
}
