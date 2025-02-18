package changeset

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog/log"

	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	TokenPool     deployment.ContractType = "TokenPool"
	Receiver      deployment.ContractType = "Receiver"
	SPL2022Tokens deployment.ContractType = "SPL2022Tokens"
	WSOL          deployment.ContractType = "WSOL"
	FeeAggregator deployment.ContractType = "FeeAggregator"
	// for PDAs from AddRemoteChainToSolana
	RemoteSource deployment.ContractType = "RemoteSource"
	RemoteDest   deployment.ContractType = "RemoteDest"

	// Tokenpool lookup table
	TokenPoolLookupTable deployment.ContractType = "TokenPoolLookupTable"
)

// SolCCIPChainState holds public keys for all the currently deployed CCIP programs
// on a chain. If a key has zero value, it means the program does not exist on the chain.
type SolCCIPChainState struct {
	LinkToken     solana.PublicKey
	Router        solana.PublicKey
	Receiver      solana.PublicKey // for tests only
	SPL2022Tokens []solana.PublicKey
	TokenPool     solana.PublicKey
	WSOL          solana.PublicKey
	FeeQuoter     solana.PublicKey
	OffRamp       solana.PublicKey
	FeeAggregator solana.PublicKey
	// PDAs to avoid redundant lookups
	RouterConfigPDA      solana.PublicKey
	SourceChainStatePDAs map[uint64]solana.PublicKey // deprecated
	DestChainStatePDAs   map[uint64]solana.PublicKey
	TokenPoolLookupTable map[solana.PublicKey]solana.PublicKey
	FeeQuoterConfigPDA   solana.PublicKey
	OffRampConfigPDA     solana.PublicKey
	OffRampStatePDA      solana.PublicKey
}

func FetchOfframpLookupTable(ctx context.Context, chain deployment.SolChain, offRampAddress solana.PublicKey) (solana.PublicKey, error) {
	var referenceAddressesAccount solOffRamp.ReferenceAddresses
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(offRampAddress)
	err := chain.GetAccountDataBorshInto(ctx, offRampReferenceAddressesPDA, &referenceAddressesAccount)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to get offramp reference addresses: %w", err)
	}
	return referenceAddressesAccount.OfframpLookupTable, nil
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
	// Most programs upgraded in place, but some are not so we always want to
	// load the latest version
	versions := make(map[deployment.ContractType]semver.Version)
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
			offRampVersion, ok := versions[OffRamp]
			// if we have an offramp version, we need to make sure it's a newer version
			if ok {
				// if the version is not newer, skip this address
				if offRampVersion.GreaterThan(&tvStr.Version) {
					log.Debug().Str("address", address).Str("type", string(tvStr.Type)).Msg("Skipping offramp address, already loaded newer version")
					continue
				}
			}
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
		existingVersion, ok := versions[tvStr.Type]
		// This shouldn't happen, so we want to log it
		if ok {
			log.Warn().Str("existingVersion", existingVersion.String()).Str("type", string(tvStr.Type)).Msg("Duplicate address type found")
		}
		versions[tvStr.Type] = tvStr.Version
	}
	state.WSOL = solana.SolMint
	return state, nil
}

func FindSolanaAddress(tv deployment.TypeAndVersion, addresses map[string]deployment.TypeAndVersion) solana.PublicKey {
	for address, tvStr := range addresses {
		if tv.String() == tvStr.String() {
			pub := solana.MustPublicKeyFromBase58(address)
			return pub
		}
	}
	return solana.PublicKey{}
}

func (c SolCCIPChainState) OnRampBytes() ([]byte, error) {
	if !c.Router.IsZero() {
		return c.Router.Bytes(), nil
	}
	return nil, errors.New("no onramp found in the state")
}
