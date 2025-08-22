package vault

import (
	"encoding/json"
	"strconv"
)

const (
	// Note: any addition to this list should be reflected in
	// HandlerTypeForMethod in handler_factory.go
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet    = "vault.secrets.get"
	MethodSecretsUpdate = "vault.secrets.update"
	MethodSecretsDelete = "vault.secrets.delete"
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

// SignedResponse is a structure that represents a signed response from the Vault DON.
// It should be validated by the client before use.
// The Payload field contains the actual response data, while Context and Signatures
// are used for signature verification and context information.
type SignedResponse struct {
	Payload    json.RawMessage `json:"payload"`
	Context    []byte          `json:"__context"`
	Signatures [][]byte        `json:"__signatures"`
}

type ResponseBase struct {
	ID       string         `json:"id,omitempty"`
	Error    string         `json:"error,omitempty"`
	Response SignedResponse `json:"response,omitempty"`
}

func (r *ResponseBase) String() string {
	return "ResponseBase{" +
		"ID: " + r.ID +
		", Error: " + r.Error +
		", Response: SignedResponse{" +
		", Payload: " + string(r.Response.Payload) +
		", Context: []byte blob" +
		", Signatures: [][]byte blob}" +
		"}"
}

type SecretsCreateResponse struct {
	SecretID SecretIdentifier `json:"secret_id,omitempty"`
	Success  bool             `json:"success,omitempty"`
}

func (r *SecretsCreateResponse) String() string {
	return "SecretsCreateResponse{" +
		", SecretID: " + r.SecretID.String() +
		", Success: " + strconv.FormatBool(r.Success) +
		"}"
}

type SecretsGetResponse struct {
	SecretID    SecretIdentifier `json:"secret_id,omitempty"`
	SecretValue SecretData       `json:"secret_value,omitempty"`
	Error       string           `json:"error,omitempty"`
}

func (r *SecretsGetResponse) String() string {
	return "SecretsGetResponse{" +
		", SecretID: " + r.SecretID.String() +
		", SecretValue: <val>" +
		", Error: " + r.Error +
		"}"
}

type SecretData struct {
	EncryptedValue               string             `json:"encrypted_value,omitempty"`
	EncryptedDecryptionKeyShares []*EncryptedShares `json:"encrypted_decryption_key_shares,omitempty"`
}

type EncryptedShares struct {
	Shares        []string `json:"shares,omitempty"`
	EncryptionKey string   `json:"encryption_key,omitempty"`
}
