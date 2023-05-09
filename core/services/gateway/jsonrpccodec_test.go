package gateway_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

func TestJsonRPCRequest_Decode(t *testing.T) {
	t.Parallel()

	input := []byte(`{"jsonrpc": "2.0", "id": "aa-bb", "method": "upload", "params": {"body":{"don_id": "functions_local", "payload": {"field": 123}}}}`)
	msg, err := gateway.DecodeRequest(input)
	require.NoError(t, err)
	require.Equal(t, "functions_local", msg.Body.DonId)
	require.Equal(t, "aa-bb", msg.Body.MessageId)
	require.Equal(t, "upload", msg.Body.Method)
	require.NotEmpty(t, msg.Body.Payload)
}

func TestJsonRPCRequest_Encode(t *testing.T) {
	t.Parallel()

	var msg gateway.Message
	msg.Body = gateway.MessageBody{
		MessageId: "aA-bB",
		Sender:    "0x1234",
		Method:    "upload",
	}
	bytes, err := gateway.EncodeRequest(&msg)
	require.NoError(t, err)

	decoded, err := gateway.DecodeRequest(bytes)
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageId)
	require.Equal(t, "0x1234", decoded.Body.Sender)
	require.Equal(t, "upload", decoded.Body.Method)
}

func TestJsonRPCResponse_Decode(t *testing.T) {
	t.Parallel()

	input := []byte(`{"jsonrpc": "2.0", "id": "aa-bb", "result": {"body": {"don_id": "functions_local", "payload": {"field": 123}}}}`)
	msg, err := gateway.DecodeResponse(input)
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
	bytes, err := gateway.EncodeResponse(&msg)
	require.NoError(t, err)

	decoded, err := gateway.DecodeResponse(bytes)
	require.NoError(t, err)
	require.Equal(t, "aA-bB", decoded.Body.MessageId)
	require.Equal(t, "0x1234", decoded.Body.Sender)
	require.Equal(t, "upload", decoded.Body.Method)
}
