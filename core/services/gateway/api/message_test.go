package api_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
)

func TestMessage_Validate(t *testing.T) {
	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: "abcd",
			Method:    "request",
			DonId:     "donA",
			Receiver:  "0x0000000000000000000000000000000000000000",
			Payload:   []byte("datadata"),
		},
	}
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	err = msg.Sign(privateKey)
	require.NoError(t, err)

	// valid
	require.NoError(t, msg.Validate())

	// missing message ID
	msg.Body.MessageId = ""
	require.Error(t, msg.Validate())
	// message ID ending with null bytes
	msg.Body.MessageId = "myid\x00\x00"
	require.Error(t, msg.Validate())
	msg.Body.MessageId = "abcd"
	require.NoError(t, msg.Validate())

	// missing DON ID
	msg.Body.DonId = ""
	require.Error(t, msg.Validate())
	// DON ID ending with null bytes
	msg.Body.DonId = "mydon\x00\x00"
	require.Error(t, msg.Validate())
	msg.Body.DonId = "donA"
	require.NoError(t, msg.Validate())

	// method name too long
	msg.Body.Method = string(bytes.Repeat([]byte("a"), api.MessageMethodMaxLen+1))
	require.Error(t, msg.Validate())
	// empty method name
	msg.Body.Method = ""
	require.Error(t, msg.Validate())
	// method name ending with null bytes
	msg.Body.Method = "method\x00"
	require.Error(t, msg.Validate())
	msg.Body.Method = "request"
	require.NoError(t, msg.Validate())

	// incorrect receiver
	msg.Body.Receiver = "blah"
	require.Error(t, msg.Validate())
	msg.Body.Receiver = "0x0000000000000000000000000000000000000000"
	require.NoError(t, msg.Validate())

	// invalid signature
	msg.Signature = "0x00"
	require.Error(t, msg.Validate())
}

func TestMessage_MessageSignAndValidateSignature(t *testing.T) {
	t.Parallel()

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: "abcd",
			Method:    "request",
			DonId:     "donA",
			Receiver:  "0x33",
			Payload:   []byte("datadata"),
		},
	}

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Bytes()

	err = msg.Sign(privateKey)
	require.NoError(t, err)
	require.Equal(t, api.MessageSignatureHexEncodedLen, len(msg.Signature))

	// valid
	signer, err := msg.ExtractSigner()
	require.NoError(t, err)
	require.True(t, bytes.Equal(address, signer))

	// invalid
	msg.Body.MessageId = "dbca"
	signer, err = msg.ExtractSigner()
	require.NoError(t, err)
	require.False(t, bytes.Equal(address, signer))
}
