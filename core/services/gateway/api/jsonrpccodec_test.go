package api_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
)

func TestJsonRPCRequest_Decode_Correct(t *testing.T) {
	t.Parallel()

	input := []byte(`{"jsonrpc": "2.0", "id": "aa-bb", "method": "upload", "params": {"body":{"don_id": "functions_local", "payload": {"field": 123}}}}`)
	codec := api.JSONRPCCodec{}
	jsonRequest, err := jsonrpc.DecodeRequest[json.RawMessage](input, "")
	require.NoError(t, err)
	msg, err := codec.DecodeJSONRequest(jsonRequest)
	require.NoError(t, err)
	msg2, err := codec.DecodeRawRequest(input, "")
	require.NoError(t, err)
	require.Equal(t, "functions_local", msg.Body.DonID)
	require.Equal(t, "aa-bb", msg.Body.MessageID)
	require.Equal(t, "upload", msg.Body.Method)
	require.NotEmpty(t, msg.Body.Payload)
	require.Equal(t, msg.Body.DonID, msg2.Body.DonID)
	require.Equal(t, msg.Body.MessageID, msg2.Body.MessageID)
	require.Equal(t, msg.Body.Method, msg2.Body.Method)
	require.Equal(t, msg.Body.Payload, msg2.Body.Payload)
}

func TestJsonRPCRequest_Decode_Incorrect(t *testing.T) {
	t.Parallel()

	testCases := map[string]string{
		"missing params":        `{"jsonrpc": "2.0", "id": "abc", "method": "upload"}`,
		"numeric id":            `{"jsonrpc": "2.0", "id": 123, "method": "upload", "params": {}}`,
		"empty method":          `{"jsonrpc": "2.0", "id": "abc", "method": "", "params": {}}`,
		"incorrect rpc version": `{"jsonrpc": "5.1", "id": "abc", "method": "upload", "params": {}}`,
	}

	codec := api.JSONRPCCodec{}
	for _, input := range testCases {
		_, err := codec.DecodeRawRequest([]byte(input), "")
		require.Error(t, err)
	}
}

func TestJsonRPCRequest_Encode(t *testing.T) {
	t.Parallel()

	var msg api.Message
	msg.Body = api.MessageBody{
		MessageID: "aA-bB",
		Receiver:  "0x1234",
		Method:    "upload",
	}
	codec := api.JSONRPCCodec{}
	bytes, err := codec.EncodeLegacyRequest(&msg)
	require.NoError(t, err)

	decoded, err := codec.DecodeRawRequest(bytes, "")
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageID)
	require.Equal(t, "0x1234", decoded.Body.Receiver)
	require.Equal(t, "upload", decoded.Body.Method)
}

func TestJsonRPCResponse_Decode(t *testing.T) {
	t.Parallel()

	input := []byte(`{"jsonrpc": "2.0", "id": "aa-bb", "result": {"body": {"don_id": "functions_local", "payload": {"field": 123}}}}`)
	codec := api.JSONRPCCodec{}
	msg, err := codec.DecodeLegacyResponse(input)
	require.NoError(t, err)
	require.Equal(t, "functions_local", msg.Body.DonID)
	require.Equal(t, "aa-bb", msg.Body.MessageID)
	require.NotEmpty(t, msg.Body.Payload)
}

func TestJsonRPCResponse_Encode(t *testing.T) {
	t.Parallel()

	var msg api.Message
	msg.Body = api.MessageBody{
		MessageID: "aA-bB",
		Receiver:  "0x1234",
		Method:    "upload",
	}
	codec := api.JSONRPCCodec{}
	bytes := codec.EncodeLegacyResponse(&msg)

	decoded, err := codec.DecodeLegacyResponse(bytes)
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageID)
	require.Equal(t, "0x1234", decoded.Body.Receiver)
	require.Equal(t, "upload", decoded.Body.Method)
}
