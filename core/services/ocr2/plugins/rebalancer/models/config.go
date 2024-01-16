package models

import (
	"errors"
	"fmt"
	"slices"
)

const PluginName = "liquidityRebalancer"

type PluginConfig struct {
	LiquidityManagerAddress Address          `json:"liquidityManagerAddress"`
	LiquidityManagerNetwork NetworkSelector  `json:"liquidityManagerNetwork"`
	ClosePluginTimeoutSec   int              `json:"closePluginTimeoutSec"`
	RebalancerConfig        RebalancerConfig `json:"rebalancerConfig"`
}

type RebalancerConfig struct {
	Type                   string                 `json:"type"`
	RandomRebalancerConfig RandomRebalancerConfig `json:"randomRebalancerConfig"`
}

type RandomRebalancerConfig struct {
	MaxNumTransfers      int  `json:"maxNumTransfers"`
	CheckSourceDestEqual bool `json:"checkSourceDestEqual"`
}

func ValidateRebalancerConfig(config RebalancerConfig) error {
	if config.Type == "" {
		return errors.New("rebalancerType must be provided")
	}

	if !RebalancerIsSupported(config.Type) {
		return fmt.Errorf("rebalancerType %s is not supported, supported types are %+v", config.Type, AllRebalancerTypes)
	}

	if config.Type == RebalancerTypeRandom {
		return validateRandomRebalancerConfig(config.RandomRebalancerConfig)
	}

	return nil
}

func validateRandomRebalancerConfig(cfg RandomRebalancerConfig) error {
	if cfg.MaxNumTransfers <= 0 {
		return errors.New("maxNumTransfers must be positive")
	}

	return nil
}

const (
	RebalancerTypeRandom = "random"
	RebalancerTypeDummy  = "dummy"
)

var (
	AllRebalancerTypes = []string{
		RebalancerTypeRandom,
		RebalancerTypeDummy,
	}
)

func RebalancerIsSupported(rebalancerType string) bool {
	return slices.Contains(AllRebalancerTypes, rebalancerType)
}
