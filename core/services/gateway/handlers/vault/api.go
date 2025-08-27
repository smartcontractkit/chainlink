package vault

import (
	"encoding/json"
	"fmt"
)

const (
	// Note: any addition to this list should be reflected in
	// HandlerTypeForMethod in handler_factory.go
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet    = "vault.secrets.get"
	MethodSecretsUpdate = "vault.secrets.update"
	MethodSecretsDelete = "vault.secrets.delete"

	MaxBatchSize = 10
)

// SignedOCRResponse is the response format for OCR signed reports, as returned by the Vault DON.
// External clients should verify that the signatures match the payload and context, before trusting this response.
// Only after validating, clients should decode the payload for further processing.
// If however the Error field is non-empty, it indicates there was an error talking to the Vault DON.
type SignedOCRResponse struct {
	Error      string          `json:"error"`
	Payload    json.RawMessage `json:"payload"`
	Context    []byte          `json:"__context"`
	Signatures [][]byte        `json:"__signatures"`
}

func (r *SignedOCRResponse) String() string {
	return fmt.Sprintf("SignedOCRResponse { Error: %s, Payload: %s, Context: <[%d]byte blob>, Signatures: <[%d][]byte blob>}", r.Error, string(r.Payload), len(r.Context), len(r.Signatures))
}
