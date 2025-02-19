package changeset

import (
	"errors"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog/log"

	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	OfframpAddressLookupTable deployment.ContractType = "OfframpAddressLookupTable"
	TokenPool                 deployment.ContractType = "TokenPool"
	Receiver                  deployment.ContractType = "Receiver"
	SPL2022Tokens             deployment.ContractType = "SPL2022Tokens"
	WSOL                      deployment.ContractType = "WSOL"
	FeeAggregator             deployment.ContractType = "FeeAggregator"
	// for PDAs from AddRemoteChainToSolana
	RemoteSource deployment.ContractType = "RemoteSource"
	RemoteDest   deployment.ContractType = "RemoteDest"

	// Tokenpool lookup table
	TokenPoolLookupTable deployment.ContractType = "TokenPoolLookupTable"
)

// SolCCIPChainState holds public keys for all the currently deployed CCIP programs
// on a chain. If a key has zero value, it means the program does not exist on the chain.
type SolCCIPChainState struct {
	LinkToken                 solana.PublicKey
	Router                    solana.PublicKey
	OfframpAddressLookupTable solana.PublicKey
	Receiver                  solana.PublicKey // for tests only
	SPL2022Tokens             []solana.PublicKey
	TokenPool                 solana.PublicKey
	WSOL                      solana.PublicKey
	FeeQuoter                 solana.PublicKey
	OffRamp                   solana.PublicKey
	FeeAggregator             solana.PublicKey

	// PDAs to avoid redundant lookups
	RouterConfigPDA      solana.PublicKey
	SourceChainStatePDAs map[uint64]solana.PublicKey // deprecated
	DestChainStatePDAs   map[uint64]solana.PublicKey
	TokenPoolLookupTable map[solana.PublicKey]solana.PublicKey
	FeeQuoterConfigPDA   solana.PublicKey
	OffRampConfigPDA     solana.PublicKey
	OffRampStatePDA      solana.PublicKey
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
		SPL2022Tokens:        make([]solana.PublicKey, 0),
		TokenPoolLookupTable: make(map[solana.PublicKey]solana.PublicKey),
	}
	for address, tvStr := range addresses {
		switch tvStr.Type {
		case commontypes.LinkToken:
			pub := solana.MustPublicKeyFromBase58(address)
			state.LinkToken = pub
		case Router:
			pub := solana.MustPublicKeyFromBase58(address)
			state.Router = pub
			routerConfigPDA, _, err := solState.FindConfigPDA(state.Router)
			if err != nil {
				return state, err
			}
			state.RouterConfigPDA = routerConfigPDA
		case OfframpAddressLookupTable:
			pub := solana.MustPublicKeyFromBase58(address)
			state.OfframpAddressLookupTable = pub
		case Receiver:
			pub := solana.MustPublicKeyFromBase58(address)
			state.Receiver = pub
		case SPL2022Tokens:
			pub := solana.MustPublicKeyFromBase58(address)
			state.SPL2022Tokens = append(state.SPL2022Tokens, pub)
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
		case TokenPoolLookupTable:
			lookupTablePubKey := solana.MustPublicKeyFromBase58(address)
			// Labels should only have one entry
			for tokenPubKeyStr := range tvStr.Labels {
				tokenPubKey := solana.MustPublicKeyFromBase58(tokenPubKeyStr)
				state.TokenPoolLookupTable[tokenPubKey] = lookupTablePubKey
			}
		case FeeQuoter:
			pub := solana.MustPublicKeyFromBase58(address)
			state.FeeQuoter = pub
			feeQuoterConfigPDA, _, err := solState.FindFqConfigPDA(state.FeeQuoter)
			if err != nil {
				return state, err
			}
			state.FeeQuoterConfigPDA = feeQuoterConfigPDA
		case OffRamp:
			pub := solana.MustPublicKeyFromBase58(address)
			state.OffRamp = pub
			offRampConfigPDA, _, err := solState.FindOfframpConfigPDA(state.OffRamp)
			if err != nil {
				return state, err
			}
			state.OffRampConfigPDA = offRampConfigPDA
			offRampStatePDA, _, err := solState.FindOfframpStatePDA(state.OffRamp)
			if err != nil {
				return state, err
			}
			state.OffRampStatePDA = offRampStatePDA
		case FeeAggregator:
			pub := solana.MustPublicKeyFromBase58(address)
			state.FeeAggregator = pub
		default:
			log.Warn().Str("address", address).Str("type", string(tvStr.Type)).Msg("Unknown address type")
			continue
		}
	}
	state.WSOL = solana.SolMint
	return state, nil
}

func (c SolCCIPChainState) OnRampBytes() ([]byte, error) {
	if !c.Router.IsZero() {
		return c.Router.Bytes(), nil
	}
	return nil, errors.New("no onramp found in the state")
}
