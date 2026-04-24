package ccipevm

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

func Test_EVMTokenDataEncoder(t *testing.T) {
	var empty usdcAttestationPayload
	encoder := NewEVMTokenDataEncoder()

	//https://testnet.snowtrace.io/tx/0xeeb0ad6b26bacd1570a9361724a36e338f4aacf1170dec64399220b7483b7eed/eventlog?chainid=43113
	//https://iris-api-sandbox.circle.com/v1/attestations/0x69fb1b419d648cf6c9512acad303746dc85af3b864af81985c76764aba60bf6b
	realMessage, err := cciptypes.NewBytesFromString("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000f8000000000000000100000006000000000004ac0d000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d00000000000000000000000009f3b8679c73c2fef8b59b4f3444d4e156fb70aa5000000000000000000000000c08835adf4884e51ff076066706e407506826d9d000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc650000000000000000000000004f32ae7f112c26b109357785e5c66dc5d747fbce00000000000000000000000000000000000000000000000000000000000000640000000000000000000000007a4d8f8c18762d362e64b411d7490fba112811cd0000000000000000")
	require.NoError(t, err)
	realAttestation, err := cciptypes.NewBytesFromString("0xee466fbd340596aa56e3e40d249869573e4008d84d795b4f2c3cba8649083d08653d38190d0df7e0ee12ae685df2f806d100a03b3716ab1ff2013c7201f1c2d01c9af959b55a4b52dbd0319eed69ce9ace25259830e0b1bff79faf0c9c5d1b5e6d6304e824d657db38f802bcff3e97d0bd30f2ffc62b62381f52c1668ceaa5a73a1b")
	require.NoError(t, err)

	tt := []struct {
		name        string
		message     []byte
		attestation []byte
	}{
		{
			name:        "empty both fields",
			message:     nil,
			attestation: []byte{},
		},
		{
			name:        "empty attestation",
			message:     []byte("message"),
			attestation: nil,
		},
		{
			name:        "both attestation and message are set",
			message:     realMessage,
			attestation: realAttestation,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encoder.EncodeUSDC(t.Context(), tc.message, tc.attestation)
			require.NoError(t, err)

			decoded, err := abihelpers.ABIDecode(empty.AbiString(), got)
			require.NoError(t, err)

			converted := abi.ConvertType(decoded[0], &empty)
			casted, ok := converted.(*usdcAttestationPayload)
			require.True(t, ok)

			if tc.message == nil {
				require.Empty(t, casted.Message)
			} else {
				require.Equal(t, tc.message, casted.Message)
			}

			if tc.attestation == nil {
				require.Empty(t, casted.Attestation)
			} else {
				require.Equal(t, tc.attestation, casted.Attestation)
			}
		})
	}
}
