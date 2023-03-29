package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

// TODO move to evmtest?
var (
	DefaultGasLimit               uint32 = 500000
	DefaultMinimumContractPayment        = assets.NewLinkFromJuels(10_000_000_000_000) // 0.00001 LINK
)
