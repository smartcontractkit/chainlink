package vault

import (
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// InvalidVaultParamsError marks structural validation failures that must surface as InvalidParamsError.
type InvalidVaultParamsError struct {
	Method string
	Err    error
}

func (e InvalidVaultParamsError) Error() string {
	if prefix := invalidVaultParamsErrorPrefix(e.Method); prefix != "" {
		return prefix + e.Err.Error()
	}
	return e.Err.Error()
}

func (e InvalidVaultParamsError) Unwrap() error {
	return e.Err
}

// IsInvalidVaultParamsError reports whether err is a pre-auth validation failure.
func IsInvalidVaultParamsError(err error) bool {
	var invalidParams InvalidVaultParamsError
	return errors.As(err, &invalidParams)
}

func invalidVaultParamsErrorPrefix(method string) string {
	switch method {
	case vaulttypes.MethodSecretsCreate:
		return "failed to validate create secrets request: "
	case vaulttypes.MethodSecretsUpdate:
		return "failed to validate update secrets request: "
	case vaulttypes.MethodSecretsDelete:
		return "failed to validate delete secrets request: "
	case vaulttypes.MethodSecretsList:
		return "failed to validate list secret identifiers request: "
	default:
		return ""
	}
}
