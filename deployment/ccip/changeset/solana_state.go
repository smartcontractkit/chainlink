package changeset

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	AddressLookupTable deployment.ContractType = "AddressLookupTable"
	TokenPool          deployment.ContractType = "TokenPool"
	Receiver           deployment.ContractType = "Receiver"
	Sol2022Tokens      deployment.ContractType = "Sol2022Tokens"
)

// SolChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type SolCCIPChainState struct {
	LinkToken          solana.PublicKey
	Router             solana.PublicKey
	Timelock           solana.PublicKey
	AddressLookupTable solana.PublicKey // for chain writer
	Receiver           solana.PublicKey // for tests only
	// TODO: i dont know how to load a token from its address and type
	// because unlike evm, solana does not store the symbol on chain
	Sol2022Tokens map[TokenSymbol]solana.PublicKey
}

func LoadOnchainStateSolana(e deployment.Environment) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		SolChains: make(map[uint64]SolCCIPChainState),
	}
	for chainSelector, chain := range e.SolChains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
		}
		chainState, err := LoadChainStateSolana(chain, addresses)
		if err != nil {
			return state, err
		}
		state.SolChains[chainSelector] = chainState
	}
	return state, nil
}

// LoadChainStateSolana Loads all state for a SolChain into state
func LoadChainStateSolana(chain deployment.SolChain, addresses map[string]deployment.TypeAndVersion) (SolCCIPChainState, error) {
	var state SolCCIPChainState
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.LinkToken = pub
		case deployment.NewTypeAndVersion(Router, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.Router = pub
		case deployment.NewTypeAndVersion(AddressLookupTable, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.AddressLookupTable = pub
		case deployment.NewTypeAndVersion(Receiver, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.Receiver = pub
		// case deployment.NewTypeAndVersion(Sol2022Tokens, deployment.Version1_0_0).String():
		// 	pub := solana.MustPublicKeyFromBase58(address)
		// 	state.Sol2022Tokens[TokenSymbol(pub)] = pub
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}
