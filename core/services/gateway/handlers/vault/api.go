package vault

const (
	// Note: any addition to this list should be reflected in
	// HandlerTypeForMethod in handler_factory.go
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet    = "vault.secrets.get"
	MethodSecretsUpdate = "vault.secrets.update"
)

type SecretIdentifier struct {
	Key       string `json:"key,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Owner     string `json:"owner,omitempty"`
}

func (s *SecretIdentifier) String() string {
	return "SecretIdentifier{" +
		"Key: " + s.Key +
		", Namespace: " + s.Namespace +
		", Owner: " + s.Owner +
		"}"
}

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
	ID         string   `json:"id,omitempty"`
	Error      string   `json:"error,omitempty"`
	Format     string   `json:"format,omitempty"`
	Context    []byte   `json:"context,omitempty"`
	Signatures [][]byte `json:"signatures,omitempty"`
}

func (r *ResponseBase) String() string {
	return "ResponseBase{" +
		"ID: " + r.ID +
		", Error: " + r.Error +
		", Format: " + r.Format +
		", Context: []byte blob" +
		", Signatures: [][]byte blob" +
		"}"
}

type SecretsCreateResponse struct {
	ResponseBase
	SecretID SecretIdentifier `json:"secret_id,omitempty"`
	Success  bool             `json:"success,omitempty"`
}

func (r *SecretsCreateResponse) String() string {
	return "SecretsCreateResponse{" +
		"ResponseBase: " + r.ResponseBase.String() +
		", SecretID: " + r.SecretID.String() +
		"}"
}

type SecretsGetResponse struct {
	ResponseBase
	SecretID    SecretIdentifier `json:"secret_id,omitempty"`
	SecretValue SecretData       `json:"secret_value,omitempty"`
	Error       string           `json:"error,omitempty"`
}

type SecretData struct {
	EncryptedValue               string             `json:"encrypted_value,omitempty"`
	EncryptedDecryptionKeyShares []*EncryptedShares `json:"encrypted_decryption_key_shares,omitempty"`
}

type EncryptedShares struct {
	Shares        []string `json:"shares,omitempty"`
	EncryptionKey string   `json:"encryption_key,omitempty"`
}
