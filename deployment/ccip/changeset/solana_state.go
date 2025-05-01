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

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view"
	viewSolana "github.com/smartcontractkit/chainlink/deployment/ccip/view/solana"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	TokenPool     deployment.ContractType = "TokenPool"
	Receiver      deployment.ContractType = "Receiver"
	SPL2022Tokens deployment.ContractType = "SPL2022Tokens"
	SPLTokens     deployment.ContractType = "SPLTokens"
	WSOL          deployment.ContractType = "WSOL"
	CCIPCommon    deployment.ContractType = "CCIPCommon"
	// for PDAs from AddRemoteChainToSolana
	RemoteSource deployment.ContractType = "RemoteSource"
	RemoteDest   deployment.ContractType = "RemoteDest"

	// Tokenpool lookup table
	TokenPoolLookupTable deployment.ContractType = "TokenPoolLookupTable"
)

// SolCCIPChainState holds public keys for all the currently deployed CCIP programs
// on a chain. If a key has zero value, it means the program does not exist on the chain.
type SolCCIPChainState struct {
	// tokens
	LinkToken     solana.PublicKey
	WSOL          solana.PublicKey
	SPL2022Tokens []solana.PublicKey
	SPLTokens     []solana.PublicKey

	// ccip programs
	Router               solana.PublicKey
	FeeQuoter            solana.PublicKey
	OffRamp              solana.PublicKey
	BurnMintTokenPool    solana.PublicKey
	LockReleaseTokenPool solana.PublicKey
	RMNRemote            solana.PublicKey

	// test programs
	Receiver solana.PublicKey

	// PDAs to avoid redundant lookups
	RouterConfigPDA      solana.PublicKey
	SourceChainStatePDAs map[uint64]solana.PublicKey // deprecated
	DestChainStatePDAs   map[uint64]solana.PublicKey
	TokenPoolLookupTable map[solana.PublicKey]solana.PublicKey
	FeeQuoterConfigPDA   solana.PublicKey
	OffRampConfigPDA     solana.PublicKey
	OffRampStatePDA      solana.PublicKey
	RMNRemoteConfigPDA   solana.PublicKey
	RMNRemoteCursesPDA   solana.PublicKey
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
		SPLTokens:            make([]solana.PublicKey, 0),
		WSOL:                 solana.SolMint,
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
			receiverVersion, ok := versions[OffRamp]
			// if we have an receiver version, we need to make sure it's a newer version
			if ok {
				// if the version is not newer, skip this address
				if receiverVersion.GreaterThan(&tvStr.Version) {
					log.Debug().Str("address", address).Str("type", string(tvStr.Type)).Msg("Skipping receiver address, already loaded newer version")
					continue
				}
			}
			pub := solana.MustPublicKeyFromBase58(address)
			state.Receiver = pub
		case SPL2022Tokens:
			pub := solana.MustPublicKeyFromBase58(address)
			state.SPL2022Tokens = append(state.SPL2022Tokens, pub)
		case SPLTokens:
			pub := solana.MustPublicKeyFromBase58(address)
			state.SPLTokens = append(state.SPLTokens, pub)
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
		case BurnMintTokenPool:
			pub := solana.MustPublicKeyFromBase58(address)
			state.BurnMintTokenPool = pub
		case LockReleaseTokenPool:
			pub := solana.MustPublicKeyFromBase58(address)
			state.LockReleaseTokenPool = pub
		case RMNRemote:
			pub := solana.MustPublicKeyFromBase58(address)
			state.RMNRemote = pub
			rmnRemoteConfigPDA, _, err := solState.FindRMNRemoteConfigPDA(state.RMNRemote)
			if err != nil {
				return state, err
			}
			state.RMNRemoteConfigPDA = rmnRemoteConfigPDA
			rmnRemoteCursesPDA, _, err := solState.FindRMNRemoteCursesPDA(state.RMNRemote)
			if err != nil {
				return state, err
			}
			state.RMNRemoteCursesPDA = rmnRemoteCursesPDA
		default:
			continue
		}
		versions[tvStr.Type] = tvStr.Version
	}
	return state, nil
}

func (s SolCCIPChainState) TokenToTokenProgram(tokenAddress solana.PublicKey) (solana.PublicKey, error) {
	if tokenAddress.Equals(s.LinkToken) {
		return solana.Token2022ProgramID, nil
	}
	if tokenAddress.Equals(s.WSOL) {
		return solana.TokenProgramID, nil
	}
	for _, spl2022Token := range s.SPL2022Tokens {
		if spl2022Token.Equals(tokenAddress) {
			return solana.Token2022ProgramID, nil
		}
	}
	for _, splToken := range s.SPLTokens {
		if splToken.Equals(tokenAddress) {
			return solana.TokenProgramID, nil
		}
	}
	return solana.PublicKey{}, fmt.Errorf("token program not found for token address %s", tokenAddress.String())
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

func ValidateOwnershipSolana(
	e *deployment.Environment,
	chain deployment.SolChain,
	mcms bool,
	programID solana.PublicKey,
	contractType deployment.ContractType,
	tokenAddress solana.PublicKey, // for token pools only
) error {
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	timelockSignerPDA := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	config, _, err := solState.FindConfigPDA(programID)
	if err != nil {
		return fmt.Errorf("failed to find config PDA: %w", err)
	}
	switch contractType {
	case Router:
		programData := solRouter.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commoncs.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for router: %w", err)
		}
	case OffRamp:
		programData := ccip_offramp.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commoncs.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for offramp: %w", err)
		}
	case FeeQuoter:
		programData := fee_quoter.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commoncs.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for feequoter: %w", err)
		}
	case BurnMintTokenPool:
		programData := solTestTokenPool.State{}
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddress, programID)
		err = chain.GetAccountDataBorshInto(e.GetContext(), poolConfigPDA, &programData)
		if err != nil {
			return nil
		}
		if err := commoncs.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Config.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for burnmint_token_pool: %w", err)
		}
	case LockReleaseTokenPool:
		programData := solTestTokenPool.State{}
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddress, programID)
		err = chain.GetAccountDataBorshInto(e.GetContext(), poolConfigPDA, &programData)
		if err != nil {
			return nil
		}
		if err := commoncs.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Config.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for lockrelease_token_pool: %w", err)
		}
	case RMNRemote:
		programData := rmn_remote.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commoncs.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for rmnremote: %w", err)
		}
	default:
		return fmt.Errorf("unsupported contract type: %s", contractType)
	}
	return nil
}

func (s SolCCIPChainState) GetRouterInfo() (router, routerConfigPDA solana.PublicKey, err error) {
	if s.Router.IsZero() {
		return solana.PublicKey{}, solana.PublicKey{}, errors.New("router not found in existing state, deploy the router first")
	}
	routerConfigPDA, _, err = solState.FindConfigPDA(s.Router)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, fmt.Errorf("failed to find config PDA: %w", err)
	}
	return s.Router, routerConfigPDA, nil
}

func FindReceiverTargetAccount(receiverID solana.PublicKey) solana.PublicKey {
	receiverTargetAccount, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, receiverID)
	return receiverTargetAccount
}

func (s SolCCIPChainState) GenerateView(solChain deployment.SolChain) (view.SolChainView, error) {
	chainView := view.NewSolChain()
	var remoteChains []uint64
	for selector := range s.DestChainStatePDAs {
		remoteChains = append(remoteChains, selector)
	}
	var allTokens []solana.PublicKey
	allTokens = append(allTokens, s.LinkToken)
	allTokens = append(allTokens, s.WSOL)
	allTokens = append(allTokens, s.SPL2022Tokens...)
	allTokens = append(allTokens, s.SPLTokens...)
	for _, token := range allTokens {
		if !token.IsZero() {
			program, err := s.TokenToTokenProgram(token)
			if err != nil {
				return chainView, fmt.Errorf("failed to find token program for token %s: %w", token, err)
			}
			tokenView, err := viewSolana.GenerateTokenView(solChain, token, program.String())
			if err != nil {
				return chainView, fmt.Errorf("failed to generate token view for token %s: %w", token, err)
			}
			if token.Equals(s.LinkToken) {
				chainView.LinkToken = tokenView
			} else {
				chainView.Tokens[token.String()] = tokenView
			}
		}
	}
	if !s.FeeQuoter.IsZero() {
		fqView, err := viewSolana.GenerateFeeQuoterView(solChain, s.FeeQuoter, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate fee quoter view %s: %w", s.FeeQuoter, err)
		}
		chainView.FeeQuoter[s.FeeQuoter.String()] = fqView
	}
	if !s.Router.IsZero() {
		routerView, err := viewSolana.GenerateRouterView(solChain, s.Router, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate router view %s: %w", s.Router, err)
		}
		chainView.Router[s.Router.String()] = routerView
	}
	if !s.OffRamp.IsZero() {
		offRampView, err := viewSolana.GenerateOffRampView(solChain, s.OffRamp, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate offramp view %s: %w", s.OffRamp, err)
		}
		chainView.OffRamp[s.OffRamp.String()] = offRampView
	}
	if !s.RMNRemote.IsZero() {
		rmnRemoteView, err := viewSolana.GenerateRMNRemoteView(solChain, s.RMNRemote, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate rmn remote view %s: %w", s.RMNRemote, err)
		}
		chainView.RMNRemote[s.RMNRemote.String()] = rmnRemoteView
	}
	if !s.BurnMintTokenPool.IsZero() {
		tokenPoolView, err := viewSolana.GenerateTokenPoolView(solChain, s.BurnMintTokenPool, remoteChains, allTokens, "BurnMintTokenPool")
		if err != nil {
			return chainView, fmt.Errorf("failed to generate burn mint token pool view %s: %w", s.BurnMintTokenPool, err)
		}
		chainView.TokenPool[s.BurnMintTokenPool.String()] = tokenPoolView
	}
	if !s.LockReleaseTokenPool.IsZero() {
		tokenPoolView, err := viewSolana.GenerateTokenPoolView(solChain, s.LockReleaseTokenPool, remoteChains, allTokens, "LockReleaseTokenPool")
		if err != nil {
			return chainView, fmt.Errorf("failed to generate lock release token pool view %s: %w", s.LockReleaseTokenPool, err)
		}
		chainView.TokenPool[s.LockReleaseTokenPool.String()] = tokenPoolView
	}
	return chainView, nil
}

func (s SolCCIPChainState) GetFeeAggregator(chain deployment.SolChain) solana.PublicKey {
	var config solRouter.Config
	configPDA, _, _ := solState.FindConfigPDA(s.Router)
	err := chain.GetAccountDataBorshInto(context.Background(), configPDA, &config)
	if err != nil {
		return solana.PublicKey{}
	}
	return config.FeeAggregator
}
