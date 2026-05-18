package models

import (
	"encoding/json"
	"fmt"
	"math/big"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type NetworkLiquidity struct {
	Network   NetworkSelector
	Liquidity *ubig.Big
}

func (n NetworkLiquidity) String() string {
	return fmt.Sprintf("NetworkLiquidity{Network: %d, Liquidity: %s}", n.Network, n.Liquidity.String())
}

func NewNetworkLiquidity(chain NetworkSelector, liq *big.Int) NetworkLiquidity {
	return NetworkLiquidity{
		Network:   chain,
		Liquidity: ubig.New(liq),
	}
}

type Observation struct {
	// LiquidityPerChain is the liquidity per chain that is known in the rebalancer graph.
	LiquidityPerChain []NetworkLiquidity
	// ResolvedTransfers are the resolved versions of the proposed transfers in the last outcome.
	ResolvedTransfers []Transfer
	// PendingTransfers are transfers that are in one of the TransferStatus states.
	PendingTransfers []PendingTransfer
	// InflightTransfers are the transfers that are currently in flight and have not been included onchain.
	InflightTransfers []Transfer
	// Edges are the edges of the rebalancer graph.
	Edges []Edge
	// ConfigDigests contains the config digests for each chain and rebalancer.
	ConfigDigests []ConfigDigestWithMeta
}

func NewObservation(
	liqPerChain []NetworkLiquidity,
	resolvedTransfers []Transfer,
	pendingTransfers []PendingTransfer,
	inflightTransfers []Transfer,
	edges []Edge,
	configDigests []ConfigDigestWithMeta,
) Observation {
	return Observation{
		LiquidityPerChain: liqPerChain,
		PendingTransfers:  pendingTransfers,
		InflightTransfers: inflightTransfers,
		ResolvedTransfers: resolvedTransfers,
		Edges:             edges,
		ConfigDigests:     configDigests,
	}
}

func (o Observation) Encode() ([]byte, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("observation %#v encoding unexpected internal error: %w", o, err)
	}
	return b, nil
}

func DecodeObservation(b []byte) (Observation, error) {
	var obs Observation
	err := json.Unmarshal(b, &obs)
	return obs, err
}

type Outcome struct {
	// These are transfers proposed by the rebalancing algorithm to reach a balanced state
	// in terms of liquidity.
	// These are not yet ready to execute by the plugin because they have not been resolved.
	// Proposed transfers are only resolved in the Observation stage of OCR3.
	ProposedTransfers []ProposedTransfer

	// These are transfers that have been proposed by the rebalancing algorithm and have been
	// resolved in the last observation.
	// Since these are "send" operations, they are ready to execute onchain.
	ResolvedTransfers []Transfer

	// These are transfers that are in one of the TransferStatus states.
	// Depending on their state they may be ready to execute onchain.
	PendingTransfers []PendingTransfer

	// ConfigDigests contains the config digests for each chain and rebalancer.
	ConfigDigests []ConfigDigestWithMeta
}

func NewOutcome(
	proposedTransfers []ProposedTransfer,
	resolvedTransfers []Transfer,
	pendingTransfers []PendingTransfer,
	configDigests []ConfigDigestWithMeta,
) Outcome {
	return Outcome{
		ProposedTransfers: proposedTransfers,
		ResolvedTransfers: resolvedTransfers,
		PendingTransfers:  pendingTransfers,
		ConfigDigests:     configDigests,
	}
}

func (o Outcome) Encode() ([]byte, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("outcome %#v encoding unexpected internal error: %w", o, err)
	}
	return b, nil
}

func DecodeOutcome(b []byte) (Outcome, error) {
	var decodedOutcome Outcome
	err := json.Unmarshal(b, &decodedOutcome)
	return decodedOutcome, err
}

type ConfigDigestWithMeta struct {
	Digest     ConfigDigest
	NetworkSel NetworkSelector
}
