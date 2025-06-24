package api

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

type JsonRPCCodec struct {
}

var _ Codec = (*JsonRPCCodec)(nil)

func (*JsonRPCCodec) DecodeRequest(msgBytes []byte, jwtToken string) (*Message, error) {
	var request jsonrpc2.Request
	jsonRpcHandler := jsonrpc2.Handler{}
	request, err := jsonRpcHandler.DecodeRequest(msgBytes, jwtToken)
	if err != nil {
		return nil, err
	}
	var msg Message
	err = json.Unmarshal(request.Params, &msg)
	if err != nil {
		return nil, err
	}
	msg.Body.MessageId = request.ID
	msg.Body.Method = request.Method
	return &msg, nil
}

func (*JsonRPCCodec) EncodeRequest(msg *Message) ([]byte, error) {
	params, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	request := jsonrpc2.Request{
		Version: jsonrpc2.JsonRpcVersion,
		ID:      msg.Body.MessageId,
		Method:  msg.Body.Method,
		Params:  params,
	}
	return json.Marshal(request)
}

func (*JsonRPCCodec) DecodeResponse(msgBytes []byte) (*Message, error) {
	var response jsonrpc2.Response
	err := json.Unmarshal(msgBytes, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("received non-empty error field: %v", response.Error)
	}
	var msg Message
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response.Result, &msg)
	msg.Body.MessageId = response.ID
	return &msg, nil
}

func (*JsonRPCCodec) EncodeResponse(msg *Message) ([]byte, error) {
	result, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	response := jsonrpc2.Response{
		Version: jsonrpc2.JsonRpcVersion,
		ID:      msg.Body.MessageId,
		Result:  result,
	}
	return json.Marshal(response)
}

func (*JsonRPCCodec) EncodeNewErrorResponse(id string, code int, message string, data []byte) ([]byte, error) {
	response := jsonrpc2.Response{
		Version: jsonrpc2.JsonRpcVersion,
		ID:      id,
		Error: &jsonrpc2.WireError{
			Code:    int64(code),
			Message: message,
			Data:    (*json.RawMessage)(&data),
		},
	}
	return json.Marshal(response)
}
