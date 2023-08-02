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
	msg.Body.MessageId = "abcd"

	// missing DON ID
	msg.Body.DonId = ""
	require.Error(t, msg.Validate())
	msg.Body.DonId = "donA"

	// method too long
	msg.Body.Method = string(bytes.Repeat([]byte("a"), api.MessageMethodMaxLen+1))
	require.Error(t, msg.Validate())
	msg.Body.Method = "request"

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
			Payload:   []byte("datadata"),
		},
	}

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Bytes()

	err = msg.Sign(privateKey)
	require.NoError(t, err)
	require.Equal(t, api.MessageSignatureHexEncodedLen, len(msg.Signature))

	signer, err := msg.ValidateSignature()
	require.NoError(t, err)
	require.True(t, bytes.Equal(address, signer))
}
