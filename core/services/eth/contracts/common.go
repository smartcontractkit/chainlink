package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func mustGetABI(json string) abi.ABI {
	abi, err := abi.JSON(strings.NewReader(json))
	if err != nil {
		panic("could not parse ABI: " + err.Error())
	}
	return abi
}
