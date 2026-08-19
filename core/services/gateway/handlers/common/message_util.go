package common

import (
	"encoding/json"
	"errors"
	"fmt"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
)

// ValidatedMessageFromResp validates and extracts legacy Gateway Message
// from JSON-RPC results field in the response
func ValidatedMessageFromResp(resp *jsonrpc.Response[json.RawMessage]) (*api.Message, error) {
	if resp.Error != nil {
		return nil, fmt.Errorf("received error, ID: %s", resp.ID)
	}
	if resp.Result == nil {
		return nil, fmt.Errorf("response result is nil, ID: %s", resp.ID)
	}
	var msg api.Message
	err := json.Unmarshal(*resp.Result, &msg)
	if err != nil {
		return nil, err
	}
	msg.Body.MessageId = resp.ID
	err = msg.Validate()
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// ValidatedMessageFromReq validated and extracts a legacy Gateway Message
// from params field of JSON-RPC request
func ValidatedMessageFromReq(req *jsonrpc.Request[json.RawMessage]) (*api.Message, error) {
	if req.Version != "2.0" {
		return nil, errors.New("incorrect jsonrpc version")
	}
	if req.Method == "" {
		return nil, errors.New("empty method field")
	}
	if req.Params == nil {
		return nil, errors.New("missing params attribute")
	}
	var m api.Message
	err := json.Unmarshal(*req.Params, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request params: %w", err)
	}
	m.Body.Method = req.Method
	m.Body.MessageId = req.ID
	err = m.Validate()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ValidatedResponseFromMessage converts a legacy Gateway Message to a JSON-RPC response
func ValidatedResponseFromMessage(msg *api.Message) (*jsonrpc.Response[json.RawMessage], error) {
	if msg == nil {
		return nil, errors.New("nil message")
	}
	if msg.Body.MessageId == "" {
		return nil, errors.New("message ID is empty")
	}
	res, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	rawResult := json.RawMessage(res)
	resp := &jsonrpc.Response[json.RawMessage]{
		Version: "2.0",
		ID:      msg.Body.MessageId,
		Result:  &rawResult,
		Method:  msg.Body.Method,
	}
	return resp, nil
}

// ValidatedRequestFromMessage converts a legacy Gateway Message to a JSON-RPC request
func ValidatedRequestFromMessage(msg *api.Message) (*jsonrpc.Request[json.RawMessage], error) {
	if msg == nil {
		return nil, errors.New("nil message")
	}
	if msg.Body.MessageId == "" {
		return nil, errors.New("message ID is empty")
	}
	if msg.Body.Method == "" {
		return nil, errors.New("method is empty")
	}
	params, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	rawParams := json.RawMessage(params)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      msg.Body.MessageId,
		Method:  msg.Body.Method,
		Params:  &rawParams,
	}
	return req, nil
}
