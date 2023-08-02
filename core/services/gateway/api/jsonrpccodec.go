package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Wrapping/unwrapping Message objects into JSON RPC ones folllowing https://www.jsonrpc.org/specification
type JsonRPCRequest struct {
	Version string   `json:"jsonrpc"`
	Id      string   `json:"id"`
	Method  string   `json:"method"`
	Params  *Message `json:"params,omitempty"`
}

type JsonRPCResponse struct {
	Version string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Result  *Message      `json:"result,omitempty"`
	Error   *JsonRPCError `json:"error,omitempty"`
}

// JSON-RPC error can only be sent to users. It is not used for messages between Gateways and Nodes.
type JsonRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type JsonRPCCodec struct {
}

var _ Codec = (*JsonRPCCodec)(nil)

func (*JsonRPCCodec) DecodeRequest(msgBytes []byte) (*Message, error) {
	var request JsonRPCRequest
	err := json.Unmarshal(msgBytes, &request)
	if err != nil {
		return nil, err
	}
	if request.Version != "2.0" {
		return nil, errors.New("incorrect jsonrpc version")
	}
	if request.Method == "" {
		return nil, errors.New("empty method field")
	}
	if request.Params == nil {
		return nil, errors.New("missing params attribute")
	}
	request.Params.Body.MessageId = request.Id
	request.Params.Body.Method = request.Method
	return request.Params, nil
}

func (*JsonRPCCodec) EncodeRequest(msg *Message) ([]byte, error) {
	request := JsonRPCRequest{
		Version: "2.0",
		Id:      msg.Body.MessageId,
		Method:  msg.Body.Method,
		Params:  msg,
	}
	return json.Marshal(request)
}

func (*JsonRPCCodec) DecodeResponse(msgBytes []byte) (*Message, error) {
	var response JsonRPCResponse
	err := json.Unmarshal(msgBytes, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("received non-empty error field: %v", response.Error)
	}
	if response.Result != nil {
		response.Result.Body.MessageId = response.Id
	}
	return response.Result, nil
}

func (*JsonRPCCodec) EncodeResponse(msg *Message) ([]byte, error) {
	response := JsonRPCResponse{
		Version: "2.0",
		Id:      msg.Body.MessageId,
		Result:  msg,
	}
	return json.Marshal(response)
}

func (*JsonRPCCodec) EncodeNewErrorResponse(id string, code int, message string, data []byte) ([]byte, error) {
	response := JsonRPCResponse{
		Version: "2.0",
		Id:      id,
		Error: &JsonRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	return json.Marshal(response)
}
