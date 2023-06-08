package api_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestSignatures_MessageSignAndValidate(t *testing.T) {
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

	signature, err := api.SignMessage(&msg.Body, privateKey)
	require.NoError(t, err)
	require.Equal(t, 65, len(signature))

	msg.Signature = utils.StringToHex(string(signature))
	msg.Body.Sender = utils.StringToHex(string(address))
	require.NoError(t, api.ValidateMessageSignature(msg))
}

func TestSignatures_BytesSignAndValidate(t *testing.T) {
	t.Parallel()

	data := []byte("data_data")

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Bytes()

	signature, err := api.SignData(privateKey, data)
	require.NoError(t, err)
	require.Equal(t, 65, len(signature))

	signer, err := api.ValidateSignature(signature, data)
	require.NoError(t, err)
	require.True(t, bytes.Equal(signer, address))
}
