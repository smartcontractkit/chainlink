package ccipevm

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

type usdcAttestationPayload struct {
	Message     []byte
	Attestation []byte
}

func (m usdcAttestationPayload) AbiString() string {
	return `
	[{
		"components": [
			{"name": "message", "type": "bytes"},
			{"name": "attestation", "type": "bytes"}
		],
		"type": "tuple"
	}]`
}

type EVMTokenDataEncoder struct{}

func NewEVMTokenDataEncoder() EVMTokenDataEncoder {
	return EVMTokenDataEncoder{}
}

func (e EVMTokenDataEncoder) EncodeUSDC(_ context.Context, message cciptypes.Bytes, attestation cciptypes.Bytes) (cciptypes.Bytes, error) {
	return abihelpers.EncodeAbiStruct(usdcAttestationPayload{
		Message:     message,
		Attestation: attestation,
	})
}
