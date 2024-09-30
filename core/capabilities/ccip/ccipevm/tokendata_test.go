package ccipevm

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

func Test_EVMTokenDataEncoder(t *testing.T) {
	var empty usdcAttestationPayload
	encoder := NewEVMTokenDataEncoder()

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
			message:     []byte("message"),
			attestation: []byte("attestation"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encoder.EncodeUSDC(tests.Context(t), tc.message, tc.attestation)
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
