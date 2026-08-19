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
	codec := api.JsonRPCCodec{}
	jsonRequest, err := jsonrpc.DecodeRequest[json.RawMessage](input, "")
	require.NoError(t, err)
	msg, err := codec.DecodeJSONRequest(jsonRequest)
	require.NoError(t, err)
	msg2, err := codec.DecodeRawRequest(input, "")
	require.NoError(t, err)
	require.Equal(t, "functions_local", msg.Body.DonId)
	require.Equal(t, "aa-bb", msg.Body.MessageId)
	require.Equal(t, "upload", msg.Body.Method)
	require.NotEmpty(t, msg.Body.Payload)
	require.Equal(t, msg.Body.DonId, msg2.Body.DonId)
	require.Equal(t, msg.Body.MessageId, msg2.Body.MessageId)
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

	codec := api.JsonRPCCodec{}
	for _, input := range testCases {
		_, err := codec.DecodeRawRequest([]byte(input), "")
		require.Error(t, err)
	}
}

func TestJsonRPCRequest_Encode(t *testing.T) {
	t.Parallel()

	var msg api.Message
	msg.Body = api.MessageBody{
		MessageId: "aA-bB",
		Receiver:  "0x1234",
		Method:    "upload",
	}
	codec := api.JsonRPCCodec{}
	bytes, err := codec.EncodeLegacyRequest(&msg)
	require.NoError(t, err)

	decoded, err := codec.DecodeRawRequest(bytes, "")
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageId)
	require.Equal(t, "0x1234", decoded.Body.Receiver)
	require.Equal(t, "upload", decoded.Body.Method)
}

func TestJsonRPCResponse_Decode(t *testing.T) {
	t.Parallel()

	input := []byte(`{"jsonrpc": "2.0", "id": "aa-bb", "result": {"body": {"don_id": "functions_local", "payload": {"field": 123}}}}`)
	codec := api.JsonRPCCodec{}
	msg, err := codec.DecodeLegacyResponse(input)
	require.NoError(t, err)
	require.Equal(t, "functions_local", msg.Body.DonId)
	require.Equal(t, "aa-bb", msg.Body.MessageId)
	require.NotEmpty(t, msg.Body.Payload)
}

func TestJsonRPCResponse_Encode(t *testing.T) {
	t.Parallel()

	var msg api.Message
	msg.Body = api.MessageBody{
		MessageId: "aA-bB",
		Receiver:  "0x1234",
		Method:    "upload",
	}
	codec := api.JsonRPCCodec{}
	bytes := codec.EncodeLegacyResponse(&msg)

	decoded, err := codec.DecodeLegacyResponse(bytes)
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageId)
	require.Equal(t, "0x1234", decoded.Body.Receiver)
	require.Equal(t, "upload", decoded.Body.Method)
}
