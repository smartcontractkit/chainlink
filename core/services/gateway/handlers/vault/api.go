package vault

const (
	// Note: any addition to this list should be reflected in
	// HandlerTypeForMethod in handler_factory.go
	MethodSecretsCreate = "vault.secrets.create"
)

type SecretsCreateRequest struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type ResponseBase struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type SecretsCreateResponse struct {
	ResponseBase
	SecretID string `json:"secret_id,omitempty"`
}
