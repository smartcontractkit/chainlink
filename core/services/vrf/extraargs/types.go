package extraargs

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

const functionSignatureLength = 4
const boolAbiType = `[{ "type": "bool" }]`

var extraArgsV1Tag = crypto.Keccak256([]byte("VRF ExtraArgsV1"))[:4]

func FromExtraArgsV1(extraArgs []byte) (nativePayment bool, err error) {
	decodedBool, err := utils.ABIDecode(boolAbiType, extraArgs[functionSignatureLength:])
	if err != nil {
		return false, fmt.Errorf("failed to decode 0x%x to bool", extraArgs[functionSignatureLength:])
	}
	nativePayment, ok := decodedBool[0].(bool)
	if !ok {
		return false, fmt.Errorf("failed to decode 0x%x to bool", extraArgs[functionSignatureLength:])
	}
	return nativePayment, nil
}

func ExtraArgsV1(nativePayment bool) ([]byte, error) {
	encodedArgs, err := utils.ABIEncode(boolAbiType, nativePayment)
	if err != nil {
		return nil, err
	}
	return append(extraArgsV1Tag, encodedArgs...), nil
}
