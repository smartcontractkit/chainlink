package changeset

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gagliardetto/solana-go"

	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	AddressLookupTable deployment.ContractType = "AddressLookupTable"
	TokenPool          deployment.ContractType = "TokenPool"
	Receiver           deployment.ContractType = "Receiver"
	SPL2022Tokens      deployment.ContractType = "SPL2022Tokens"
	WSOL               deployment.ContractType = "WSOL"
	// for PDAs from AddRemoteChainToSolana
	RemoteSource deployment.ContractType = "RemoteSource"
	RemoteDest   deployment.ContractType = "RemoteDest"
)

// SolChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type SolCCIPChainState struct {
	LinkToken          solana.PublicKey
	Router             solana.PublicKey
	Timelock           solana.PublicKey
	AddressLookupTable solana.PublicKey // for chain writer
	Receiver           solana.PublicKey // for tests only
	SPL2022Tokens      []solana.PublicKey
	TokenPool          solana.PublicKey
	WSOL               solana.PublicKey
	// PDAs to avoid redundant lookups
	RouterStatePDA       solana.PublicKey
	RouterConfigPDA      solana.PublicKey
	SourceChainStatePDAs map[uint64]solana.PublicKey
	DestChainStatePDAs   map[uint64]solana.PublicKey
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
	state := SolCCIPChainState{
		SourceChainStatePDAs: make(map[uint64]solana.PublicKey),
		DestChainStatePDAs:   make(map[uint64]solana.PublicKey),
	}
	var spl2022Tokens []solana.PublicKey
	for address, tvStr := range addresses {
		switch tvStr.Type {
		case commontypes.LinkToken:
			pub := solana.MustPublicKeyFromBase58(address)
			state.LinkToken = pub
		case Router:
			pub := solana.MustPublicKeyFromBase58(address)
			state.Router = pub
			routerStatePDA, _, err := solState.FindStatePDA(state.Router)
			if err != nil {
				return state, err
			}
			state.RouterStatePDA = routerStatePDA
			routerConfigPDA, _, err := solState.FindConfigPDA(state.Router)
			if err != nil {
				return state, err
			}
			state.RouterConfigPDA = routerConfigPDA
		case AddressLookupTable:
			pub := solana.MustPublicKeyFromBase58(address)
			state.AddressLookupTable = pub
		case Receiver:
			pub := solana.MustPublicKeyFromBase58(address)
			state.Receiver = pub
		case SPL2022Tokens:
			pub := solana.MustPublicKeyFromBase58(address)
			spl2022Tokens = append(spl2022Tokens, pub)
		case TokenPool:
			pub := solana.MustPublicKeyFromBase58(address)
			state.TokenPool = pub
		case RemoteSource:
			pub := solana.MustPublicKeyFromBase58(address)
			// Labels should only have one entry
			for selStr := range tvStr.Labels {
				selector, err := strconv.ParseUint(selStr, 10, 64)
				if err != nil {
					return state, err
				}
				state.SourceChainStatePDAs[selector] = pub
			}
		case RemoteDest:
			pub := solana.MustPublicKeyFromBase58(address)
			// Labels should only have one entry
			for selStr := range tvStr.Labels {
				selector, err := strconv.ParseUint(selStr, 10, 64)
				if err != nil {
					return state, err
				}
				state.DestChainStatePDAs[selector] = pub
			}
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	state.WSOL = solana.SolMint
	state.SPL2022Tokens = spl2022Tokens
	return state, nil
}
