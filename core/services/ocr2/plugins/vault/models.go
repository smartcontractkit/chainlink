package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
)

type Request struct {
	Payload      proto.Message
	ResponseChan chan *Response

	id         string
	expiryTime time.Time
}

func (r *Request) ID() string {
	return r.id
}

func (r *Request) Copy() *Request {
	newRequest := &Request{
		Payload: proto.Clone(r.Payload),

		// intentionally not copied as we want to keep the reference
		ResponseChan: r.ResponseChan,

		// copied by value
		id:         r.id,
		expiryTime: r.expiryTime,
	}
	return newRequest
}

func (r *Request) ExpiryTime() time.Time {
	return r.expiryTime
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
		ID:    r.id,
		Error: fmt.Sprintf("timeout exceeded: could not process request %s before expiry", r.id),
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

type errResp struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

type payloadResp struct {
	Payload    json.RawMessage `json:"payload"`
	Context    []byte          `json:"__context"`
	Signatures [][]byte        `json:"__signatures"`
}

func (r *Response) ToJSONRPCResult() ([]byte, error) {
	if r.Error != "" {
		return json.Marshal(errResp{Error: r.Error, Success: false})
	}

	return json.Marshal(payloadResp{
		Payload:    r.Payload,
		Context:    r.Context,
		Signatures: r.Signatures,
	})
}

func (r *Response) RequestID() string {
	return r.ID
}
