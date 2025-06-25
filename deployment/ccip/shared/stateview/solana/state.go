package solana

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog/log"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	solBurnMintTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/burnmint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solLockReleaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/lockrelease_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view"
	solanaview "github.com/smartcontractkit/chainlink/deployment/ccip/view/solana"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// CCIPChainState holds public keys for all the currently deployed CCIP programs
// on a chain. If a key has zero value, it means the program does not exist on the chain.
type CCIPChainState struct {
	// tokens
	LinkToken     solana.PublicKey
	WSOL          solana.PublicKey
	SPL2022Tokens []solana.PublicKey
	SPLTokens     []solana.PublicKey

	// ccip programs
	Router                solana.PublicKey
	FeeQuoter             solana.PublicKey
	OffRamp               solana.PublicKey
	RMNRemote             solana.PublicKey
	BurnMintTokenPools    map[string]solana.PublicKey // metadata id -> BurnMintTokenPool
	LockReleaseTokenPools map[string]solana.PublicKey // metadata id -> LockReleaseTokenPool

	// test programs
	Receiver solana.PublicKey

	// PDAs to avoid redundant lookups
	RouterConfigPDA      solana.PublicKey
	SourceChainStatePDAs map[uint64]solana.PublicKey // deprecated
	DestChainStatePDAs   map[uint64]solana.PublicKey
	TokenPoolLookupTable map[solana.PublicKey]map[test_token_pool.PoolType]map[string]solana.PublicKey // token -> token pool type -> metadata identifier -> lookup table
	FeeQuoterConfigPDA   solana.PublicKey
	OffRampConfigPDA     solana.PublicKey
	OffRampStatePDA      solana.PublicKey
	RMNRemoteConfigPDA   solana.PublicKey
	RMNRemoteCursesPDA   solana.PublicKey
}

func (s CCIPChainState) TokenToTokenProgram(tokenAddress solana.PublicKey) (solana.PublicKey, error) {
	if tokenAddress.Equals(s.LinkToken) || tokenAddress.Equals(s.WSOL) {
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

func (s CCIPChainState) GetRouterInfo() (router, routerConfigPDA solana.PublicKey, err error) {
	if s.Router.IsZero() {
		return solana.PublicKey{}, solana.PublicKey{}, errors.New("router not found in existing state, deploy the router first")
	}
	routerConfigPDA, _, err = state.FindConfigPDA(s.Router)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, fmt.Errorf("failed to find config PDA: %w", err)
	}
	return s.Router, routerConfigPDA, nil
}

func (s CCIPChainState) GetActiveTokenPool(
	poolType solTestTokenPool.PoolType,
	metadata string,
) (solana.PublicKey, cldf.ContractType) {
	switch poolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		if metadata == "" {
			return s.BurnMintTokenPools[shared.CLLMetadata], shared.BurnMintTokenPool
		}
		return s.BurnMintTokenPools[metadata], shared.BurnMintTokenPool
	case solTestTokenPool.LockAndRelease_PoolType:
		if metadata == "" {
			return s.LockReleaseTokenPools[shared.CLLMetadata], shared.LockReleaseTokenPool
		}
		return s.LockReleaseTokenPools[metadata], shared.LockReleaseTokenPool
	default:
		return solana.PublicKey{}, ""
	}
}

func (s CCIPChainState) ValidatePoolDeployment(
	e *cldf.Environment,
	poolType solTestTokenPool.PoolType,
	selector uint64,
	tokenPubKey solana.PublicKey,
	validatePoolConfig bool,
	metadata string,
) error {
	chain := e.BlockChains.SolanaChains()[selector]

	var tokenPool solana.PublicKey
	var poolConfigAccount interface{}

	if _, err := s.TokenToTokenProgram(tokenPubKey); err != nil {
		return fmt.Errorf("token %s not found in existing state, deploy the token first", tokenPubKey.String())
	}
	tokenPool, _ = s.GetActiveTokenPool(poolType, metadata)
	if tokenPool.IsZero() {
		return fmt.Errorf("token pool of type %s not found in existing state, deploy the token pool first for chain %d", poolType, chain.Selector)
	}
	switch poolType {
	case solTestTokenPool.BurnAndMint_PoolType:
		poolConfigAccount = solBurnMintTokenPool.State{}
	case solTestTokenPool.LockAndRelease_PoolType:
		poolConfigAccount = solLockReleaseTokenPool.State{}
	default:
		return fmt.Errorf("invalid pool type: %s", poolType)
	}

	if validatePoolConfig {
		poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey, tokenPool)
		if err != nil {
			return fmt.Errorf("failed to get token pool config address (mint: %s, pool: %s): %w", tokenPubKey.String(), tokenPool.String(), err)
		}
		if err := chain.GetAccountDataBorshInto(context.Background(), poolConfigPDA, &poolConfigAccount); err != nil {
			return fmt.Errorf("token pool config not found (mint: %s, pool: %s, type: %s): %w", tokenPubKey.String(), tokenPool.String(), poolType, err)
		}
	}
	return nil
}

func (s CCIPChainState) CommonValidation(e cldf.Environment, selector uint64, tokenPubKey solana.PublicKey) error {
	_, ok := e.BlockChains.SolanaChains()[selector]
	if !ok {
		return fmt.Errorf("chain selector %d not found in environment", selector)
	}
	if tokenPubKey.Equals(s.LinkToken) || tokenPubKey.Equals(s.WSOL) {
		return nil
	}
	if _, err := s.TokenToTokenProgram(tokenPubKey); err != nil {
		return fmt.Errorf("token %s not found in existing state, deploy the token first", tokenPubKey.String())
	}
	return nil
}

func (s CCIPChainState) ValidateRouterConfig(chain cldf_solana.Chain) error {
	_, routerConfigPDA, err := s.GetRouterInfo()
	if err != nil {
		return err
	}
	var routerConfigAccount solRouter.Config
	err = chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, initialize the router first %d", chain.Selector)
	}
	return nil
}

func (s CCIPChainState) ValidateFeeAggregatorConfig(chain cldf_solana.Chain) error {
	if s.GetFeeAggregator(chain).IsZero() {
		return fmt.Errorf("fee aggregator not found in existing state, set the fee aggregator first for chain %d", chain.Selector)
	}
	return nil
}

func (s CCIPChainState) ValidateFeeQuoterConfig(chain cldf_solana.Chain) error {
	if s.FeeQuoter.IsZero() {
		return fmt.Errorf("fee quoter not found in existing state, deploy the fee quoter first for chain %d", chain.Selector)
	}
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(s.FeeQuoter)
	err := chain.GetAccountDataBorshInto(context.Background(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		return fmt.Errorf("fee quoter config not found in existing state, initialize the fee quoter first %d", chain.Selector)
	}
	return nil
}

func (s CCIPChainState) ValidateOffRampConfig(chain cldf_solana.Chain) error {
	if s.OffRamp.IsZero() {
		return fmt.Errorf("offramp not found in existing state, deploy the offramp first for chain %d", chain.Selector)
	}
	var offRampConfig solOffRamp.Config
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(s.OffRamp)
	err := chain.GetAccountDataBorshInto(context.Background(), offRampConfigPDA, &offRampConfig)
	if err != nil {
		return fmt.Errorf("offramp config not found in existing state, initialize the offramp first %d", chain.Selector)
	}
	return nil
}

func (s CCIPChainState) GenerateView(e *cldf.Environment, selector uint64) (view.SolChainView, error) {
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
			tokenView, err := solanaview.GenerateTokenView(e.BlockChains.SolanaChains()[selector], token, program.String())
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
		fqView, err := solanaview.GenerateFeeQuoterView(e.BlockChains.SolanaChains()[selector], s.FeeQuoter, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate fee quoter view %s: %w", s.FeeQuoter, err)
		}
		chainView.FeeQuoter[s.FeeQuoter.String()] = fqView
	}
	if !s.Router.IsZero() {
		routerView, err := solanaview.GenerateRouterView(e.BlockChains.SolanaChains()[selector], s.Router, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate router view %s: %w", s.Router, err)
		}
		chainView.Router[s.Router.String()] = routerView
	}
	if !s.OffRamp.IsZero() {
		offRampView, err := solanaview.GenerateOffRampView(e.BlockChains.SolanaChains()[selector], s.OffRamp, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate offramp view %s: %w", s.OffRamp, err)
		}
		chainView.OffRamp[s.OffRamp.String()] = offRampView
	}
	if !s.RMNRemote.IsZero() {
		rmnRemoteView, err := solanaview.GenerateRMNRemoteView(e.BlockChains.SolanaChains()[selector], s.RMNRemote, remoteChains, allTokens)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate rmn remote view %s: %w", s.RMNRemote, err)
		}
		chainView.RMNRemote[s.RMNRemote.String()] = rmnRemoteView
	}
	for metadata, tokenPool := range s.BurnMintTokenPools {
		if tokenPool.IsZero() {
			continue
		}
		tokenPoolView, err := solanaview.GenerateTokenPoolView(e.BlockChains.SolanaChains()[selector], tokenPool, remoteChains, allTokens, test_token_pool.BurnAndMint_PoolType.String(), metadata)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate burn mint token pool view %s: %w", tokenPool, err)
		}
		chainView.TokenPool[tokenPool.String()] = tokenPoolView
	}
	for metadata, tokenPool := range s.LockReleaseTokenPools {
		if tokenPool.IsZero() {
			continue
		}
		tokenPoolView, err := solanaview.GenerateTokenPoolView(e.BlockChains.SolanaChains()[selector], tokenPool, remoteChains, allTokens, test_token_pool.LockAndRelease_PoolType.String(), metadata)
		if err != nil {
			return chainView, fmt.Errorf("failed to generate lock release token pool view %s: %w", tokenPool, err)
		}
		chainView.TokenPool[tokenPool.String()] = tokenPoolView
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(selector)
	if err != nil {
		return chainView, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	chainView.MCMSWithTimelock, err = solanaview.GenerateMCMSWithTimelockView(e.BlockChains.SolanaChains()[selector], addresses)
	if err != nil {
		e.Logger.Error("failed to generate MCMS with timelock view: %w", err)
		return chainView, nil
	}
	return chainView, nil
}

func (s CCIPChainState) GetFeeAggregator(chain cldf_solana.Chain) solana.PublicKey {
	var config ccip_router.Config
	configPDA, _, _ := state.FindConfigPDA(s.Router)
	err := chain.GetAccountDataBorshInto(context.Background(), configPDA, &config)
	if err != nil {
		return solana.PublicKey{}
	}
	return config.FeeAggregator
}

func FetchOfframpLookupTable(ctx context.Context, chain cldf_solana.Chain, offRampAddress solana.PublicKey) (solana.PublicKey, error) {
	var referenceAddressesAccount ccip_offramp.ReferenceAddresses
	offRampReferenceAddressesPDA, _, _ := state.FindOfframpReferenceAddressesPDA(offRampAddress)
	err := chain.GetAccountDataBorshInto(ctx, offRampReferenceAddressesPDA, &referenceAddressesAccount)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to get offramp reference addresses: %w", err)
	}
	return referenceAddressesAccount.OfframpLookupTable, nil
}

// LoadChainStateSolana Loads all state for a SolChain into state
func LoadChainStateSolana(chain cldf_solana.Chain, addresses map[string]cldf.TypeAndVersion) (CCIPChainState, error) {
	solState := CCIPChainState{
		SourceChainStatePDAs:  make(map[uint64]solana.PublicKey),
		DestChainStatePDAs:    make(map[uint64]solana.PublicKey),
		BurnMintTokenPools:    make(map[string]solana.PublicKey),
		LockReleaseTokenPools: make(map[string]solana.PublicKey),
		SPL2022Tokens:         make([]solana.PublicKey, 0),
		SPLTokens:             make([]solana.PublicKey, 0),
		WSOL:                  solana.SolMint,
		TokenPoolLookupTable:  make(map[solana.PublicKey]map[test_token_pool.PoolType]map[string]solana.PublicKey),
	}
	// Most programs upgraded in place, but some are not so we always want to
	// load the latest version
	versions := make(map[cldf.ContractType]semver.Version)
	for address, tvStr := range addresses {
		switch tvStr.Type {
		case types.LinkToken:
			pub := solana.MustPublicKeyFromBase58(address)
			solState.LinkToken = pub
		case shared.Router:
			pub := solana.MustPublicKeyFromBase58(address)
			solState.Router = pub
			routerConfigPDA, _, err := state.FindConfigPDA(solState.Router)
			if err != nil {
				return solState, err
			}
			solState.RouterConfigPDA = routerConfigPDA
		case shared.Receiver:
			receiverVersion, ok := versions[shared.OffRamp]
			// if we have an receiver version, we need to make sure it's a newer version
			if ok {
				// if the version is not newer, skip this address
				if receiverVersion.GreaterThan(&tvStr.Version) {
					log.Debug().Str("address", address).Str("type", string(tvStr.Type)).Msg("Skipping receiver address, already loaded newer version")
					continue
				}
			}
			pub := solana.MustPublicKeyFromBase58(address)
			solState.Receiver = pub
		case shared.SPL2022Tokens:
			pub := solana.MustPublicKeyFromBase58(address)
			solState.SPL2022Tokens = append(solState.SPL2022Tokens, pub)
		case shared.SPLTokens:
			pub := solana.MustPublicKeyFromBase58(address)
			solState.SPLTokens = append(solState.SPLTokens, pub)
		case shared.RemoteSource:
			pub := solana.MustPublicKeyFromBase58(address)
			// Labels should only have one entry
			for selStr := range tvStr.Labels {
				selector, err := strconv.ParseUint(selStr, 10, 64)
				if err != nil {
					return solState, err
				}
				solState.SourceChainStatePDAs[selector] = pub
			}
		case shared.RemoteDest:
			pub := solana.MustPublicKeyFromBase58(address)
			// Labels should only have one entry
			for selStr := range tvStr.Labels {
				selector, err := strconv.ParseUint(selStr, 10, 64)
				if err != nil {
					return solState, err
				}
				solState.DestChainStatePDAs[selector] = pub
			}
		case shared.TokenPoolLookupTable:
			lookupTablePubKey := solana.MustPublicKeyFromBase58(address)
			var poolType *test_token_pool.PoolType
			var tokenPubKey solana.PublicKey
			var poolMetadata string
			for label := range tvStr.Labels {
				maybeTokenPubKey, err := solana.PublicKeyFromBase58(label)
				if err == nil {
					tokenPubKey = maybeTokenPubKey
				} else {
					switch label {
					case test_token_pool.BurnAndMint_PoolType.String():
						t := test_token_pool.BurnAndMint_PoolType
						poolType = &t
					case test_token_pool.LockAndRelease_PoolType.String():
						t := test_token_pool.LockAndRelease_PoolType
						poolType = &t
					default:
						poolMetadata = label
					}
				}
			}
			if poolMetadata == "" {
				poolMetadata = shared.CLLMetadata
			}
			if poolType == nil {
				t := test_token_pool.BurnAndMint_PoolType
				poolType = &t
			}
			if solState.TokenPoolLookupTable[tokenPubKey] == nil {
				solState.TokenPoolLookupTable[tokenPubKey] = make(map[test_token_pool.PoolType]map[string]solana.PublicKey)
			}
			if solState.TokenPoolLookupTable[tokenPubKey][*poolType] == nil {
				solState.TokenPoolLookupTable[tokenPubKey][*poolType] = make(map[string]solana.PublicKey)
			}
			solState.TokenPoolLookupTable[tokenPubKey][*poolType][poolMetadata] = lookupTablePubKey
		case shared.FeeQuoter:
			pub := solana.MustPublicKeyFromBase58(address)
			solState.FeeQuoter = pub
			feeQuoterConfigPDA, _, err := state.FindFqConfigPDA(solState.FeeQuoter)
			if err != nil {
				return solState, err
			}
			solState.FeeQuoterConfigPDA = feeQuoterConfigPDA
		case shared.OffRamp:
			offRampVersion, ok := versions[shared.OffRamp]
			// if we have an offramp version, we need to make sure it's a newer version
			if ok {
				// if the version is not newer, skip this address
				if offRampVersion.GreaterThan(&tvStr.Version) {
					log.Debug().Str("address", address).Str("type", string(tvStr.Type)).Msg("Skipping offramp address, already loaded newer version")
					continue
				}
			}
			pub := solana.MustPublicKeyFromBase58(address)
			solState.OffRamp = pub
			offRampConfigPDA, _, err := state.FindOfframpConfigPDA(solState.OffRamp)
			if err != nil {
				return solState, err
			}
			solState.OffRampConfigPDA = offRampConfigPDA
			offRampStatePDA, _, err := state.FindOfframpStatePDA(solState.OffRamp)
			if err != nil {
				return solState, err
			}
			solState.OffRampStatePDA = offRampStatePDA
		case shared.BurnMintTokenPool:
			pub := solana.MustPublicKeyFromBase58(address)
			if len(tvStr.Labels) == 0 {
				solState.BurnMintTokenPools[shared.CLLMetadata] = pub
			}
			// Labels should only have one entry
			for metadataStr := range tvStr.Labels {
				solState.BurnMintTokenPools[metadataStr] = pub
			}
		case shared.LockReleaseTokenPool:
			pub := solana.MustPublicKeyFromBase58(address)
			if len(tvStr.Labels) == 0 {
				solState.LockReleaseTokenPools[shared.CLLMetadata] = pub
			}
			// Labels should only have one entry
			for metadataStr := range tvStr.Labels {
				solState.LockReleaseTokenPools[metadataStr] = pub
			}
		case shared.RMNRemote:
			pub := solana.MustPublicKeyFromBase58(address)
			solState.RMNRemote = pub
			rmnRemoteConfigPDA, _, err := state.FindRMNRemoteConfigPDA(solState.RMNRemote)
			if err != nil {
				return solState, err
			}
			solState.RMNRemoteConfigPDA = rmnRemoteConfigPDA
			rmnRemoteCursesPDA, _, err := state.FindRMNRemoteCursesPDA(solState.RMNRemote)
			if err != nil {
				return solState, err
			}
			solState.RMNRemoteCursesPDA = rmnRemoteCursesPDA
		default:
			continue
		}
		versions[tvStr.Type] = tvStr.Version
	}
	return solState, nil
}

func FindSolanaAddress(tv cldf.TypeAndVersion, addresses map[string]cldf.TypeAndVersion) solana.PublicKey {
	for address, tvStr := range addresses {
		if tv.String() == tvStr.String() {
			pub := solana.MustPublicKeyFromBase58(address)
			return pub
		}
	}
	return solana.PublicKey{}
}

func ValidateOwnershipSolana(
	e *cldf.Environment,
	chain cldf_solana.Chain,
	mcms bool,
	programID solana.PublicKey,
	contractType cldf.ContractType,
	tokenAddress solana.PublicKey, // for token pools only
) error {
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	timelockSignerPDA := commonstate.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	config, _, err := state.FindConfigPDA(programID)
	if err != nil {
		return fmt.Errorf("failed to find config PDA: %w", err)
	}
	switch contractType {
	case shared.Router:
		programData := ccip_router.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commonchangeset.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for router: %w", err)
		}
	case shared.OffRamp:
		programData := ccip_offramp.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commonchangeset.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for offramp: %w", err)
		}
	case shared.FeeQuoter:
		programData := fee_quoter.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commonchangeset.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for feequoter: %w", err)
		}
	case shared.BurnMintTokenPool:
		programData := test_token_pool.State{}
		poolConfigPDA, _ := tokens.TokenPoolConfigAddress(tokenAddress, programID)
		err = chain.GetAccountDataBorshInto(e.GetContext(), poolConfigPDA, &programData)
		if err != nil {
			return nil
		}
		if err := commonchangeset.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Config.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for burnmint_token_pool: %w", err)
		}
	case shared.LockReleaseTokenPool:
		programData := test_token_pool.State{}
		poolConfigPDA, _ := tokens.TokenPoolConfigAddress(tokenAddress, programID)
		err = chain.GetAccountDataBorshInto(e.GetContext(), poolConfigPDA, &programData)
		if err != nil {
			return nil
		}
		if err := commonchangeset.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Config.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for lockrelease_token_pool: %w", err)
		}
	case shared.RMNRemote:
		programData := rmn_remote.Config{}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return fmt.Errorf("failed to get account data: %w", err)
		}
		if err := commonchangeset.ValidateOwnershipSolanaCommon(mcms, chain.DeployerKey.PublicKey(), timelockSignerPDA, programData.Owner); err != nil {
			return fmt.Errorf("failed to validate ownership for rmnremote: %w", err)
		}
	default:
		return fmt.Errorf("unsupported contract type: %s", contractType)
	}
	return nil
}

func IsSolanaProgramOwnedByTimelock(
	e *cldf.Environment,
	chain cldf_solana.Chain,
	chainState CCIPChainState,
	contractType cldf.ContractType,
	tokenAddress solana.PublicKey, // for token pools only
	tokenPoolMetadata string,
) bool {
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return false
	}
	mcmState, err := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return false
	}
	timelockSignerPDA := commonstate.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	switch contractType {
	case shared.Router:
		programData := ccip_router.Config{}
		config, _, err := state.FindConfigPDA(chainState.Router)
		if err != nil {
			return false
		}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return false
		}
		return programData.Owner.Equals(timelockSignerPDA)
	case shared.OffRamp:
		programData := ccip_offramp.Config{}
		config, _, err := state.FindConfigPDA(chainState.OffRamp)
		if err != nil {
			return false
		}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return false
		}
		return programData.Owner.Equals(timelockSignerPDA)
	case shared.FeeQuoter:
		programData := fee_quoter.Config{}
		config, _, err := state.FindConfigPDA(chainState.FeeQuoter)
		if err != nil {
			return false
		}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return false
		}
		return programData.Owner.Equals(timelockSignerPDA)
	case shared.BurnMintTokenPool:
		programData := test_token_pool.State{}
		metadata := shared.CLLMetadata
		if tokenPoolMetadata != "" {
			metadata = tokenPoolMetadata
		}
		poolConfigPDA, _ := tokens.TokenPoolConfigAddress(tokenAddress, chainState.BurnMintTokenPools[metadata])
		err = chain.GetAccountDataBorshInto(e.GetContext(), poolConfigPDA, &programData)
		if err != nil {
			return false
		}
		return programData.Config.Owner.Equals(timelockSignerPDA)
	case shared.LockReleaseTokenPool:
		programData := test_token_pool.State{}
		metadata := shared.CLLMetadata
		if tokenPoolMetadata != "" {
			metadata = tokenPoolMetadata
		}
		poolConfigPDA, _ := tokens.TokenPoolConfigAddress(tokenAddress, chainState.LockReleaseTokenPools[metadata])
		err = chain.GetAccountDataBorshInto(e.GetContext(), poolConfigPDA, &programData)
		if err != nil {
			return false
		}
		return programData.Config.Owner.Equals(timelockSignerPDA)
	case shared.RMNRemote:
		programData := rmn_remote.Config{}
		config, _, err := state.FindConfigPDA(chainState.RMNRemote)
		if err != nil {
			return false
		}
		err = chain.GetAccountDataBorshInto(e.GetContext(), config, &programData)
		if err != nil {
			return false
		}
		return programData.Owner.Equals(timelockSignerPDA)
	default:
		return false
	}
}

func FindReceiverTargetAccount(receiverID solana.PublicKey) solana.PublicKey {
	receiverTargetAccount, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, receiverID)
	return receiverTargetAccount
}
