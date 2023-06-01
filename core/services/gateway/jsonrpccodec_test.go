package gateway_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

func TestJsonRPCRequest_Decode_Correct(t *testing.T) {
	t.Parallel()

	input := []byte(`{"jsonrpc": "2.0", "id": "aa-bb", "method": "upload", "params": {"body":{"don_id": "functions_local", "payload": {"field": 123}}}}`)
	codec := gateway.JsonRPCCodec{}
	msg, err := codec.DecodeRequest(input)
	require.NoError(t, err)
	require.Equal(t, "functions_local", msg.Body.DonId)
	require.Equal(t, "aa-bb", msg.Body.MessageId)
	require.Equal(t, "upload", msg.Body.Method)
	require.NotEmpty(t, msg.Body.Payload)
}

func TestJsonRPCRequest_Decode_Incorrect(t *testing.T) {
	t.Parallel()

	testCases := map[string]string{
		"missing params":        `{"jsonrpc": "2.0", "id": "abc", "method": "upload"}`,
		"numeric id":            `{"jsonrpc": "2.0", "id": 123, "method": "upload", "params": {}}`,
		"empty method":          `{"jsonrpc": "2.0", "id": "abc", "method": "", "params": {}}`,
		"incorrect rpc version": `{"jsonrpc": "5.1", "id": "abc", "method": "upload", "params": {}}`,
	}

	codec := gateway.JsonRPCCodec{}
	for _, input := range testCases {
		_, err := codec.DecodeRequest([]byte(input))
		require.Error(t, err)
	}
}

func TestJsonRPCRequest_Encode(t *testing.T) {
	t.Parallel()

	var msg gateway.Message
	msg.Body = gateway.MessageBody{
		MessageId: "aA-bB",
		Sender:    "0x1234",
		Method:    "upload",
	}
	codec := gateway.JsonRPCCodec{}
	bytes, err := codec.EncodeRequest(&msg)
	require.NoError(t, err)

	decoded, err := codec.DecodeRequest(bytes)
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageId)
	require.Equal(t, "0x1234", decoded.Body.Sender)
	require.Equal(t, "upload", decoded.Body.Method)
}

func TestJsonRPCResponse_Decode(t *testing.T) {
	t.Parallel()

	input := []byte(`{"jsonrpc": "2.0", "id": "aa-bb", "result": {"body": {"don_id": "functions_local", "payload": {"field": 123}}}}`)
	codec := gateway.JsonRPCCodec{}
	msg, err := codec.DecodeResponse(input)
	require.NoError(t, err)
	require.Equal(t, "functions_local", msg.Body.DonId)
	require.Equal(t, "aa-bb", msg.Body.MessageId)
	require.NotEmpty(t, msg.Body.Payload)
}

func TestJsonRPCResponse_Encode(t *testing.T) {
	t.Parallel()

	var msg gateway.Message
	msg.Body = gateway.MessageBody{
		MessageId: "aA-bB",
		Sender:    "0x1234",
		Method:    "upload",
	}
	codec := gateway.JsonRPCCodec{}
	bytes, err := codec.EncodeResponse(&msg)
	require.NoError(t, err)

	decoded, err := codec.DecodeResponse(bytes)
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageId)
	require.Equal(t, "0x1234", decoded.Body.Sender)
	require.Equal(t, "upload", decoded.Body.Method)
}
