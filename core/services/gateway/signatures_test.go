package gateway_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestSignatures_MessageSignAndValidate(t *testing.T) {
	t.Parallel()

	msg := &gateway.Message{
		Body: gateway.MessageBody{
			MessageId: "abcd",
			Method:    "request",
			DonId:     "donA",
			Payload:   []byte("datadata"),
		},
	}

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Bytes()

	signature, err := gateway.SignMessage(&msg.Body, privateKey)
	require.NoError(t, err)
	require.Equal(t, 65, len(signature))

	msg.Signature = utils.StringToHex(string(signature))
	msg.Body.Sender = utils.StringToHex(string(address))
	require.NoError(t, gateway.ValidateMessageSignature(msg))
}

func TestSignatures_BytesSignAndValidate(t *testing.T) {
	t.Parallel()

	data := [][]byte{[]byte("data_data")}

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Bytes()

	signature, err := gateway.SignData(data, privateKey)
	require.NoError(t, err)
	require.Equal(t, 65, len(signature))

	require.NoError(t, gateway.ValidateSignature(data, signature, address))
}
