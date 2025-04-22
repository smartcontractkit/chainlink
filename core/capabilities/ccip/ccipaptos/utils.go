package ccipaptos

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

func abiEncodeUint32(data uint32) ([]byte, error) {
	return utils.ABIEncode(`[{ "type": "uint32" }]`, data)
}

func abiDecodeUint32(data []byte) (uint32, error) {
	raw, err := utils.ABIDecode(`[{ "type": "uint32" }]`, data)
	if err != nil {
		return 0, fmt.Errorf("abi decode uint32: %w", err)
	}

	val := *abi.ConvertType(raw[0], new(uint32)).(*uint32)
	return val, nil
}
