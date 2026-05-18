package models

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
)

const PluginName = "liquidityRebalancer"

type PluginConfig struct {
	LiquidityManagerAddress Address          `json:"liquidityManagerAddress"`
	LiquidityManagerNetwork NetworkSelector  `json:"liquidityManagerNetwork,string"`
	ClosePluginTimeoutSec   int              `json:"closePluginTimeoutSec"`
	RebalancerConfig        RebalancerConfig `json:"rebalancerConfig"`
}

type RebalancerConfig struct {
	Type          string   `json:"type"`
	DefaultTarget *big.Int `json:"defaultTarget"`
	// NetworkTargetOverrides is a map of NetworkSelector to big Int amounts
	NetworkTargetOverrides map[NetworkSelector]*big.Int `json:"networkTargetOverrides"`
}

func ValidateRebalancerConfig(config RebalancerConfig) error {
	if config.Type == "" {
		return errors.New("rebalancerType must be provided")
	}

	if !RebalancerIsSupported(config.Type) {
		return fmt.Errorf("rebalancerType %s is not supported, supported types are %+v", config.Type, AllRebalancerTypes)
	}

	return nil
}

const (
	RebalancerTypeTargetAndMin = "target-and-min"
	RebalancerTypeMinLiquidity = "min-liquidity"
	RebalancerTypePingPong     = "ping-pong"
)

var (
	AllRebalancerTypes = []string{
		RebalancerTypePingPong,
		RebalancerTypeMinLiquidity,
		RebalancerTypeTargetAndMin,
	}
)

func RebalancerIsSupported(rebalancerType string) bool {
	return slices.Contains(AllRebalancerTypes, rebalancerType)
}
