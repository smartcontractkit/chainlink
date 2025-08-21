package vault

const (
	// Note: any addition to this list should be reflected in
	// HandlerTypeForMethod in handler_factory.go
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet    = "vault.secrets.get"
	MethodSecretsUpdate = "vault.secrets.update"
)
