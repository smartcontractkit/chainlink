package models

const PluginName = "liquidityRebalancer"

type PluginConfig struct {
	LiquidityManagerAddress Address   `json:"liquidityManagerAddress"`
	LiquidityManagerNetwork NetworkID `json:"liquidityManagerNetwork"`
	ClosePluginTimeoutSec   int       `json:"closePluginTimeoutSec"`
}
