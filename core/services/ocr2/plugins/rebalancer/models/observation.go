package models

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type NetworkLiquidity struct {
	Network   NetworkSelector
	Liquidity *big.Int
}

func (n NetworkLiquidity) String() string {
	return fmt.Sprintf("NetworkLiquidity{Network: %d, Liquidity: %s}", n.Network, n.Liquidity.String())
}

func NewNetworkLiquidity(chain NetworkSelector, liq *big.Int) NetworkLiquidity {
	return NetworkLiquidity{
		Network:   chain,
		Liquidity: liq,
	}
}

type Observation struct {
	LiquidityPerChain []NetworkLiquidity
	PendingTransfers  []PendingTransfer
}

func NewObservation(liqPerChain []NetworkLiquidity, pendingTransfers []PendingTransfer) Observation {
	return Observation{
		LiquidityPerChain: liqPerChain,
		PendingTransfers:  pendingTransfers,
	}
}

func (o Observation) Encode() []byte {
	b, err := json.Marshal(o)
	if err != nil {
		panic(fmt.Errorf("observation %#v encoding unexpected internal error: %w", o, err))
	}
	return b
}

func DecodeObservation(b []byte) (Observation, error) {
	var obs Observation
	err := json.Unmarshal(b, &obs)
	return obs, err
}
