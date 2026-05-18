package graph

import (
	"math/big"
	"reflect"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

type Vertex struct {
	NetworkSelector  models.NetworkSelector
	LiquidityManager models.Address
}

type XChainLiquidityManagerData struct {
	RemoteLiquidityManagerAddress models.Address
	LocalBridgeAdapterAddress     models.Address
	RemoteTokenAddress            models.Address
}

func (d XChainLiquidityManagerData) Equals(other XChainLiquidityManagerData) bool {
	return d.RemoteLiquidityManagerAddress == other.RemoteLiquidityManagerAddress &&
		d.LocalBridgeAdapterAddress == other.LocalBridgeAdapterAddress &&
		d.RemoteTokenAddress == other.RemoteTokenAddress
}

type Data struct {
	Liquidity               *big.Int
	TokenAddress            models.Address
	LiquidityManagerAddress models.Address
	XChainLiquidityManagers map[models.NetworkSelector]XChainLiquidityManagerData
	ConfigDigest            models.ConfigDigest
	NetworkSelector         models.NetworkSelector
	MinimumLiquidity        *big.Int
	TargetLiquidity         *big.Int
}

func (d Data) Equals(other Data) bool {
	return d.Liquidity.Cmp(other.Liquidity) == 0 &&
		d.TokenAddress == other.TokenAddress &&
		d.LiquidityManagerAddress == other.LiquidityManagerAddress &&
		reflect.DeepEqual(d.XChainLiquidityManagers, other.XChainLiquidityManagers) &&
		d.ConfigDigest == other.ConfigDigest &&
		d.NetworkSelector == other.NetworkSelector
}

func (d Data) Clone() Data {
	xChainRebalancers := make(map[models.NetworkSelector]XChainLiquidityManagerData)
	for k, v := range d.XChainLiquidityManagers {
		xChainRebalancers[k] = v
	}
	liqManagerAddr := models.Address{}
	copy(liqManagerAddr[:], d.LiquidityManagerAddress[:])
	tokenAddr := models.Address{}
	copy(tokenAddr[:], d.TokenAddress[:])
	liq := d.Liquidity
	if liq == nil {
		liq = big.NewInt(0)
	}
	minLiq := d.MinimumLiquidity
	if minLiq == nil {
		minLiq = big.NewInt(0)
	}
	targetLiq := d.TargetLiquidity
	if targetLiq == nil {
		targetLiq = big.NewInt(0)
	}
	return Data{
		Liquidity:               big.NewInt(0).Set(liq),
		TokenAddress:            tokenAddr,
		LiquidityManagerAddress: liqManagerAddr,
		XChainLiquidityManagers: xChainRebalancers,
		ConfigDigest:            d.ConfigDigest.Clone(),
		NetworkSelector:         d.NetworkSelector,
		MinimumLiquidity:        big.NewInt(0).Set(minLiq),
		TargetLiquidity:         big.NewInt(0).Set(targetLiq),
	}
}
