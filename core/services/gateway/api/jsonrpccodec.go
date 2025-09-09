package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

type JsonRPCCodec struct {
}

var _ Codec = (*JsonRPCCodec)(nil)

func (j *JsonRPCCodec) DecodeRawRequest(msgBytes []byte, jwtToken string) (*Message, error) {
	jsonRequest, err := jsonrpc2.DecodeRequest[json.RawMessage](msgBytes, jwtToken)
	if err != nil {
		return nil, err
	}
	return j.DecodeJSONRequest(jsonRequest)
}

func (*JsonRPCCodec) DecodeJSONRequest(request jsonrpc2.Request[json.RawMessage]) (*Message, error) {
	var msg Message
	err := json.Unmarshal(*request.Params, &msg)
	if err != nil {
		return nil, err
	}
	msg.Body.MessageId = request.ID
	msg.Body.Method = request.Method
	return &msg, nil
}

func (*JsonRPCCodec) EncodeLegacyRequest(msg *Message) ([]byte, error) {
	request := jsonrpc2.Request[Message]{
		Version: jsonrpc2.JsonRpcVersion,
		ID:      msg.Body.MessageId,
		Method:  msg.Body.Method,
		Params:  msg,
	}
	return json.Marshal(request)
}

func (*JsonRPCCodec) DecodeLegacyResponse(msgBytes []byte) (*Message, error) {
	var response jsonrpc2.Response[Message]
	err := json.Unmarshal(msgBytes, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("received non-empty error field: %v", response.Error)
	}
	if response.Result == nil {
		return nil, errors.New("received empty result field")
	}

	response.Result.Body.MessageId = response.ID
	return response.Result, nil
}

func (*JsonRPCCodec) EncodeLegacyResponse(msg *Message) []byte {
	response := jsonrpc2.Response[Message]{
		Version: jsonrpc2.JsonRpcVersion,
		ID:      msg.Body.MessageId,
		Result:  msg,
		Method:  msg.Body.Method,
	}
	rawMsg, err := json.Marshal(response)
	if err != nil {
		return fatalError(err)
	}
	return rawMsg
}

func (*JsonRPCCodec) EncodeNewErrorResponse(id string, code int64, message string, data []byte) []byte {
	response := jsonrpc2.Response[json.RawMessage]{
		Version: jsonrpc2.JsonRpcVersion,
		ID:      id,
		Error: &jsonrpc2.WireError{
			Code:    code,
			Message: message,
			Data:    (*json.RawMessage)(&data),
		},
	}
	rawErrMsg, err := json.Marshal(response)
	if err != nil {
		return fatalError(err)
	}
	return rawErrMsg
}

func fatalError(err error) []byte {
	return []byte("fatal error: " + err.Error())
}
