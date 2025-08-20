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

// EncryptedSecret represents a single encrypted secret in a batch request
type EncryptedSecret struct {
	ID             SecretIdentifier `json:"id"`
	EncryptedValue string           `json:"encrypted_value"`
}

func (e *EncryptedSecret) String() string {
	return "EncryptedSecret{" +
		"ID: " + e.ID.String() +
		", EncryptedValue: " + e.EncryptedValue +
		"}"
}

type CreateSecretsRequest struct {
	RequestID        string            `json:"request_id"`
	EncryptedSecrets []EncryptedSecret `json:"encrypted_secrets"`
}

func (c *CreateSecretsRequest) String() string {
	if len(c.EncryptedSecrets) == 0 {
		return "CreateSecretsRequest{" +
			"RequestID: " + c.RequestID +
			", EncryptedSecrets: []}"
	}

	result := "CreateSecretsRequest{" +
		"RequestID: " + c.RequestID +
		", EncryptedSecrets: ["
	for i, secret := range c.EncryptedSecrets {
		if i > 0 {
			result += ", "
		}
		result += secret.String()
	}
	result += "]}"
	return result
}

type SecretRequest struct {
	ID             SecretIdentifier `json:"id"`
	EncryptionKeys []string         `json:"encryption_keys,omitempty"`
}

type GetSecretsRequest struct {
	Requests []SecretRequest `json:"requests"`
}

func (g *GetSecretsRequest) String() string {
	if len(g.Requests) == 0 {
		return "GetSecretsRequest{Requests: []}"
	}

	result := "GetSecretsRequest{Requests: ["
	for i, request := range g.Requests {
		if i > 0 {
			result += ", "
		}
		result += "SecretRequest{" +
			"ID: " + request.ID.String() +
			", EncryptionKeys: ["
		for j, key := range request.EncryptionKeys {
			if j > 0 {
				result += ", "
			}
			result += key
		}
		result += "]}"
	}
	result += "]}"
	return result
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

type CreateSecretsResponse struct {
	Responses []CreateSecretResponse `json:"responses,omitempty"`
}

type CreateSecretResponse struct {
	ID      SecretIdentifier `json:"id,omitempty"`
	Success bool             `json:"success,omitempty"`
	Error   string           `json:"error,omitempty"`
}

func (r *CreateSecretsResponse) String() string {
	if len(r.Responses) == 0 {
		return "CreateSecretsResponse{Responses: []}"
	}

	result := "CreateSecretsResponse{Responses: ["
	for i, response := range r.Responses {
		if i > 0 {
			result += ", "
		}
		result += "CreateSecretResponse{" +
			"ID: " + response.ID.String() +
			", Success: " + strconv.FormatBool(response.Success) +
			", Error: " + response.Error +
			"}"
	}
	result += "]}"
	return result
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
