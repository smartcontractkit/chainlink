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
	Edges             []Edge
}

func NewObservation(liqPerChain []NetworkLiquidity, pendingTransfers []PendingTransfer, edges []Edge) Observation {
	return Observation{
		LiquidityPerChain: liqPerChain,
		PendingTransfers:  pendingTransfers,
		Edges:             edges,
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

type Outcome struct {
	TransfersToReachBalance []Transfer
	PendingTransfers        []PendingTransfer
}

func NewOutcome(transfersToReachBalance []Transfer, pendingTransfers []PendingTransfer) Outcome {
	return Outcome{
		TransfersToReachBalance: transfersToReachBalance,
		PendingTransfers:        pendingTransfers,
	}
}

func (o Outcome) Encode() []byte {
	b, err := json.Marshal(o)
	if err != nil {
		panic(fmt.Errorf("outcome %#v encoding unexpected internal error: %w", o, err))
	}
	return b
}

func DecodeOutcome(b []byte) (Outcome, error) {
	var decodedOutcome Outcome
	err := json.Unmarshal(b, &decodedOutcome)
	return decodedOutcome, err
}
