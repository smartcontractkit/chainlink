package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/crypto"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
)

const (
	privateKey = "6c358b4f16344f03cfce12ebf7b768301bbe6a8977c98a2a2d76699f8bc56161"
)

func unsignedMessage() api.Message {
	return api.Message{
		Body: api.MessageBody{
			Method:    "testMethod",
			MessageId: "msg-123",
			DonId:     "test_don",
		},
	}
}

func TestValidatedMessageFromReq(t *testing.T) {
	validMsg := unsignedMessage()
	key, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err)
	err = validMsg.Sign(key)
	require.NoError(t, err)
	params, err := json.Marshal(validMsg)
	require.NoError(t, err)
	rawParams := json.RawMessage(params)

	t.Run("valid request", func(t *testing.T) {
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Method:  "testMethod",
			Params:  &rawParams,
		}
		msg, err := ValidatedMessageFromReq(req)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, "testMethod", msg.Body.Method)
		require.Equal(t, "msg-123", msg.Body.MessageId)
	})

	t.Run("invalid message", func(t *testing.T) {
		invalidMsg := unsignedMessage()
		invalidParams, err := json.Marshal(invalidMsg)
		require.NoError(t, err)
		rawInvalidParams := json.RawMessage(invalidParams)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Method:  "testMethod",
			Params:  &rawInvalidParams,
		}
		_, err = ValidatedMessageFromReq(req)
		require.Error(t, err)
		require.EqualError(t, err, "invalid hex-encoded signature length")
	})

	t.Run("incorrect jsonrpc version", func(t *testing.T) {
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "1.0",
			ID:      "msg-123",
			Method:  "testMethod",
			Params:  &rawParams,
		}
		msg, err := ValidatedMessageFromReq(req)
		require.Nil(t, msg)
		require.EqualError(t, err, "incorrect jsonrpc version")
	})

	t.Run("empty method field", func(t *testing.T) {
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Method:  "",
			Params:  &rawParams,
		}
		msg, err := ValidatedMessageFromReq(req)
		require.Nil(t, msg)
		require.EqualError(t, err, "empty method field")
	})

	t.Run("missing params attribute", func(t *testing.T) {
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Method:  "testMethod",
			Params:  nil,
		}
		msg, err := ValidatedMessageFromReq(req)
		require.Nil(t, msg)
		require.EqualError(t, err, "missing params attribute")
	})

	t.Run("invalid params json", func(t *testing.T) {
		rawParams := json.RawMessage([]byte(`{invalid json}`))
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Method:  "testMethod",
			Params:  &rawParams,
		}
		msg, err := ValidatedMessageFromReq(req)
		require.Nil(t, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to unmarshal request params")
	})
}

func TestValidatedMessageFromResp(t *testing.T) {
	validMsg := unsignedMessage()
	key, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err)
	err = validMsg.Sign(key)
	require.NoError(t, err)
	result, err := json.Marshal(validMsg)
	require.NoError(t, err)
	rawResult := json.RawMessage(result)

	t.Run("valid response", func(t *testing.T) {
		resp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Result:  &rawResult,
		}
		msg, err := ValidatedMessageFromResp(resp)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, "testMethod", msg.Body.Method)
		require.Equal(t, "msg-123", msg.Body.MessageId)
	})

	t.Run("response with error", func(t *testing.T) {
		resp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Error:   &jsonrpc.WireError{Code: 1, Message: "some error"},
		}
		msg, err := ValidatedMessageFromResp(resp)
		require.Nil(t, msg)
		require.ErrorContains(t, err, "received error")
	})

	t.Run("nil result", func(t *testing.T) {
		resp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Result:  nil,
		}
		msg, err := ValidatedMessageFromResp(resp)
		require.Nil(t, msg)
		require.ErrorContains(t, err, "response result is nil")
	})

	t.Run("invalid result json", func(t *testing.T) {
		rawResult := json.RawMessage([]byte(`{invalid json}`))
		resp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Result:  &rawResult,
		}
		msg, err := ValidatedMessageFromResp(resp)
		require.Nil(t, msg)
		require.Error(t, err)
	})

	t.Run("invalid message", func(t *testing.T) {
		invalidMsg := unsignedMessage()
		result, err := json.Marshal(invalidMsg)
		require.NoError(t, err)
		rawResult := json.RawMessage(result)
		resp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      "msg-123",
			Result:  &rawResult,
		}
		msg, err := ValidatedMessageFromResp(resp)
		require.Nil(t, msg)
		require.Error(t, err)
	})
}

func TestValidatedResponseFromMessage(t *testing.T) {
	validMsg := unsignedMessage()

	t.Run("valid message", func(t *testing.T) {
		resp, err := ValidatedResponseFromMessage(&validMsg)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "2.0", resp.Version)
		require.Equal(t, "msg-123", resp.ID)
		require.NotNil(t, resp.Result)
		var msg api.Message
		err = json.Unmarshal(*resp.Result, &msg)
		require.NoError(t, err)
		require.Equal(t, validMsg.Body.Method, msg.Body.Method)
		require.Equal(t, validMsg.Body.MessageId, msg.Body.MessageId)
		require.Equal(t, validMsg.Body.DonId, msg.Body.DonId)
	})

	t.Run("nil message", func(t *testing.T) {
		resp, err := ValidatedResponseFromMessage(nil)
		require.Nil(t, resp)
		require.EqualError(t, err, "nil message")
	})

	t.Run("empty message ID", func(t *testing.T) {
		msg := validMsg
		msg.Body.MessageId = ""
		resp, err := ValidatedResponseFromMessage(&msg)
		require.Nil(t, resp)
		require.EqualError(t, err, "message ID is empty")
	})
}

func TestValidatedRequestFromMessage(t *testing.T) {
	validMsg := unsignedMessage()

	t.Run("valid message", func(t *testing.T) {
		req, err := ValidatedRequestFromMessage(&validMsg)
		require.NoError(t, err)
		require.NotNil(t, req)
		require.Equal(t, "2.0", req.Version)
		require.Equal(t, validMsg.Body.MessageId, req.ID)
		require.Equal(t, validMsg.Body.Method, req.Method)
		require.NotNil(t, req.Params)
		var msg api.Message
		err = json.Unmarshal(*req.Params, &msg)
		require.NoError(t, err)
		require.Equal(t, validMsg.Body.Method, msg.Body.Method)
		require.Equal(t, validMsg.Body.MessageId, msg.Body.MessageId)
		require.Equal(t, validMsg.Body.DonId, msg.Body.DonId)
	})

	t.Run("nil message", func(t *testing.T) {
		req, err := ValidatedRequestFromMessage(nil)
		require.Nil(t, req)
		require.EqualError(t, err, "nil message")
	})

	t.Run("empty message ID", func(t *testing.T) {
		msg := validMsg
		msg.Body.MessageId = ""
		req, err := ValidatedRequestFromMessage(&msg)
		require.Nil(t, req)
		require.EqualError(t, err, "message ID is empty")
	})

	t.Run("empty method", func(t *testing.T) {
		msg := validMsg
		msg.Body.Method = ""
		req, err := ValidatedRequestFromMessage(&msg)
		require.Nil(t, req)
		require.EqualError(t, err, "method is empty")
	})
}
