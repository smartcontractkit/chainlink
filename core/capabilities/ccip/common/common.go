package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
)

type ExtraDataDecoded struct {
	ExtraArgsDecoded    map[string]any
	DestExecDataDecoded []map[string]any
}

type ExecCallData struct {
	ReportContext [2][32]byte
	Report        []byte
	Info          ccipocr3.ExecuteReportInfo
	ExtraData     ExtraDataDecoded
}

// HashedCapabilityID returns the hashed capability id in a manner equivalent to the capability registry.
func HashedCapabilityID(capabilityLabelledName, capabilityVersion string) (r [32]byte, err error) {
	// TODO: investigate how to avoid parsing the ABI everytime.
	tabi := `[{"type": "string"}, {"type": "string"}]`
	abiEncoded, err := utils.ABIEncode(tabi, capabilityLabelledName, capabilityVersion)
	if err != nil {
		return r, fmt.Errorf("failed to ABI encode capability version and labelled name: %w", err)
	}

	h := crypto.Keccak256(abiEncoded)
	copy(r[:], h)
	return r, nil
}
