package vaulttypes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

var DefaultNamespace = "main"

const (
	// MethodSecretsCreate Note: additional methods should be reflected
	// in the `Methods` list below.
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet    = "vault.secrets.get"
	MethodSecretsUpdate = "vault.secrets.update"
	MethodSecretsDelete = "vault.secrets.delete"
	MethodSecretsList   = "vault.secrets.list"
	MethodPublicKeyGet  = "vault.publicKey.get"

	// RequestIDSeparator is used to separate parts(owner, user-provided-requestId) of the request ID.
	RequestIDSeparator = "::"

	// MaxBatchSize is the maximum number of secrets that can be created/updated/deleted in a single request.
	MaxBatchSize = 10
)

var Methods = []string{
	MethodSecretsCreate,
	MethodSecretsUpdate,
	MethodSecretsDelete,
	MethodSecretsList,
	MethodPublicKeyGet,
}

// SignedOCRResponse is the response format for OCR signed reports, as returned by the Vault DON.
// External clients should verify that the signatures match the payload and context, before trusting this response.
// Only after validating, clients should decode the payload for further processing.
// If however the Error field is non-empty, it indicates there was an error talking to the Vault DON.
type SignedOCRResponse struct {
	Error      string          `json:"error"`
	Payload    json.RawMessage `json:"payload"`
	Context    []byte          `json:"context"`
	Signatures [][]byte        `json:"signatures"`
}

func (r *SignedOCRResponse) String() string {
	return fmt.Sprintf("SignedOCRResponse { Error: %s, Payload: %s, Context: <[%d]byte blob>, Signatures: <[%d][]byte blob>}", r.Error, string(r.Payload), len(r.Context), len(r.Signatures))
}

type SecretsService interface {
	CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*Response, error)
	UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*Response, error)
	GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*Response, error)
	DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*Response, error)
	ListSecretIdentifiers(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) (*Response, error)

	GetPublicKey(ctx context.Context, request *vaultcommon.GetPublicKeyRequest) (*vaultcommon.GetPublicKeyResponse, error)
}

type Request struct {
	Payload      proto.Message
	ResponseChan chan *Response

	IDVal         string
	ExpiryTimeVal time.Time
}

func (r *Request) ID() string {
	return r.IDVal
}

func (r *Request) Copy() *Request {
	newRequest := &Request{
		Payload: proto.Clone(r.Payload),

		// intentionally not copied as we want to keep the reference
		ResponseChan: r.ResponseChan,

		// copied by value
		IDVal:         r.IDVal,
		ExpiryTimeVal: r.ExpiryTimeVal,
	}
	return newRequest
}

func (r *Request) ExpiryTime() time.Time {
	return r.ExpiryTimeVal
}

func (r *Request) SendResponse(ctx context.Context, response *Response) {
	select {
	case <-ctx.Done():
		return
	case r.ResponseChan <- response:
	}
}

func (r *Request) SendTimeout(ctx context.Context) {
	r.SendResponse(ctx, &Response{
		ID:    r.IDVal,
		Error: fmt.Sprintf("timeout exceeded: could not process request %s before expiry", r.IDVal),
	})
}

type Response struct {
	ID         string
	Error      string
	Payload    []byte
	Format     string
	Context    []byte
	Signatures [][]byte
}

func (r *Response) ToJSONRPCResult() ([]byte, error) {
	return json.Marshal(SignedOCRResponse{
		Error:      r.Error,
		Payload:    r.Payload,
		Context:    r.Context,
		Signatures: r.Signatures,
	})
}

func (r *Response) RequestID() string {
	return r.ID
}

func (r *Response) String() string {
	return fmt.Sprintf("Response { ID: %s, Error: %s, Payload: %s, Format: %s }", r.ID, r.Error, string(r.Payload), r.Format)
}
