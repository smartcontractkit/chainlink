package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultapi "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
)

var DefaultNamespace = "main"

type SecretsService interface {
	CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*Response, error)
	UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*Response, error)
	GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*Response, error)
	DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*Response, error)
}

func KeyFor(id *vaultcommon.SecretIdentifier) string {
	namespace := id.Namespace
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return fmt.Sprintf("%s::%s::%s", id.Owner, namespace, id.Key)
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
	return json.Marshal(vaultapi.SignedOCRResponse{
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
