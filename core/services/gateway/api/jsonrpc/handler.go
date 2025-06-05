package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Wrapping/unwrapping Message objects into JSON RPC ones folllowing https://www.jsonrpc.org/specification
type Request struct {
	Version string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	Auth    string          `json:"auth,omitempty"`
}

type Response struct {
	Version string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

// JSON-RPC error can only be sent to users. It is not used for messages between Gateways and Nodes.
type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Handler struct {
}

func (*Handler) DecodeRequest(requestBytes []byte) (Request, error) {
	var request Request
	err := json.Unmarshal(requestBytes, &request)
	if err != nil {
		return Request{}, err
	}
	if request.Version != "2.0" {
		return Request{}, errors.New("incorrect jsonrpc version")
	}
	if request.Method == "" {
		return Request{}, errors.New("empty method field")
	}
	if request.Params == nil {
		return Request{}, errors.New("missing params attribute")
	}
	return request, nil
}

func (*Handler) EncodeRequest(request *Request) ([]byte, error) {
	return json.Marshal(request)
}

func (*Handler) DecodeResponse(responseBytes []byte) (Response, error) {
	var response Response
	err := json.Unmarshal(responseBytes, &response)
	if err != nil {
		return Response{}, err
	}
	if response.Error != nil {
		return Response{}, fmt.Errorf("received non-empty error field: %v", response.Error)
	}
	return response, nil
}

func (*Handler) EncodeResponse(response *Response) ([]byte, error) {
	return json.Marshal(response)
}

func (r *Request) EncodeErrorReponse(err *Error) ([]byte, error) {
	return json.Marshal(Response{
		Version: "2.0",
		ID:      r.ID,
		Error:   err,
	})
}
