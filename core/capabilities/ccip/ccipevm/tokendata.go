package ccipevm

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/utils/abihelpers"
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

func (e EVMTokenDataEncoder) EncodeUSDC(_ context.Context, message, attestation ccipocr3.Bytes) (ccipocr3.Bytes, error) {
	return abihelpers.EncodeAbiStruct(usdcAttestationPayload{
		Message:     message,
		Attestation: attestation,
	})
}
