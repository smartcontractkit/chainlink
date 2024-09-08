package ccipcalc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

func EvmAddrsToGeneric(evmAddrs ...common.Address) []cciptypes.Address {
	res := make([]cciptypes.Address, 0, len(evmAddrs))
	for _, addr := range evmAddrs {
		res = append(res, cciptypes.Address(addr.String()))
	}
	return res
}

func EvmAddrToGeneric(evmAddr common.Address) cciptypes.Address {
	return cciptypes.Address(evmAddr.String())
}

func GenericAddrsToEvm(genericAddrs ...cciptypes.Address) ([]common.Address, error) {
	evmAddrs := make([]common.Address, 0, len(genericAddrs))
	for _, addr := range genericAddrs {
		if !common.IsHexAddress(string(addr)) {
			return nil, fmt.Errorf("%s not an evm address", addr)
		}
		evmAddrs = append(evmAddrs, common.HexToAddress(string(addr)))
	}
	return evmAddrs, nil
}

func GenericAddrToEvm(genAddr cciptypes.Address) (common.Address, error) {
	evmAddrs, err := GenericAddrsToEvm(genAddr)
	if err != nil {
		return common.Address{}, err
	}
	return evmAddrs[0], nil
}

func HexToAddress(h string) cciptypes.Address {
	return cciptypes.Address(common.HexToAddress(h).String())
}
