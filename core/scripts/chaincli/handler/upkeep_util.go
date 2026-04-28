package handler

import (
	"math"
	"math/big"
	"strings"

	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/accounts/abi"
	registry11 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_2"
)

const (
	zeroPrefix   = "0x"
	upkeepPrefix = "UPx"
)

var (
	registry11ABI = mustParseABI(registry11.KeeperRegistryABI)
	registry12ABI = mustParseABI(registry12.KeeperRegistryABI)
)

func mustParseABI(json string) abi.ABI {
	a, err := abi.JSON(strings.NewReader(json))
	if err != nil {
		panic(err)
	}
	return *a
}

// LeastSignificant32 returns the least significant 32 bits of the input as a uint64.
func LeastSignificant32(num *big.Int) uint64 {
	max32 := big.NewInt(math.MaxUint32)
	return big.NewInt(0).And(num, max32).Uint64()
}

// ParseUpkeepId parses the upkeep id input string to a big int pointer.
func ParseUpkeepId(upkeepIdStr string) (*big.Int, bool) {
	if strings.HasPrefix(upkeepIdStr, upkeepPrefix) {
		upkeepIdStr = zeroPrefix + upkeepIdStr[len(upkeepPrefix):]
	}

	upkeepId, ok := ethmath.ParseBig256(upkeepIdStr)
	if !ok {
		return ethmath.ParseBig256(zeroPrefix + upkeepIdStr)
	}
	return upkeepId, ok
}
