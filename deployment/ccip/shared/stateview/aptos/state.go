// Package aptos is a thin bridge to github.com/smartcontractkit/chainlink-aptos/deployment/state,
// kept so the deprecated deployment/ccip/changeset/aptos tree compiles until it is deleted.
// New code should import the chainlink-aptos module directly.
package aptos

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	aptosstate "github.com/smartcontractkit/chainlink-aptos/deployment/state"
)

// CCIPChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type CCIPChainState = aptosstate.CCIPChainState

// LoadOnchainStateAptos loads chain state for Aptos chains from env
func LoadOnchainStateAptos(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	return aptosstate.LoadOnchainState(env)
}

var (
	GetOfframpDynamicConfig = aptosstate.GetOfframpDynamicConfig
	FindAptosAddress        = aptosstate.FindAptosAddress
)
