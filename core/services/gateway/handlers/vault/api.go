package vault

const (
	// Note: any addition to this list should be reflected in
	// HandlerTypeForMethod in handler_factory.go
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet = "vault.secrets.get"
)

type SecretsCreateRequest struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Owner string `json:"owner"`
}

type SecretsGetRequest struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
}

type ResponseBase struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type SecretsCreateResponse struct {
	ResponseBase
	SecretID string `json:"secret_id,omitempty"`
}

type SecretsGetResponse struct {
	ResponseBase
	SecretID string `json:"secret_id,omitempty"`
	SecretValue string `json:"secret_value,omitempty"`
}
