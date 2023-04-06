package types

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
)

type MockTxManager = txmmocks.TxManager[*Address, *TxHash, *BlockHash]

func MustGetABI(json string) abi.ABI {
	abi, err := abi.JSON(strings.NewReader(json))
	if err != nil {
		panic("could not parse ABI: " + err.Error())
	}
	return abi
}
