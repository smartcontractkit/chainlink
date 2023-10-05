package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

// NewCachedFeeTokens cache fee tokens returned from PriceRegistry
func NewCachedFeeTokens(
	lp logpoller.LogPoller,
	priceRegistry ccipdata.PriceRegistryReader,
	optimisticConfirmations int64,
) *CachedChain[[]common.Address] {
	return &CachedChain[[]common.Address]{
		observedEvents:          priceRegistry.FeeTokenEvents(),
		logPoller:               lp,
		address:                 []common.Address{priceRegistry.Address()},
		optimisticConfirmations: optimisticConfirmations,
		lock:                    &sync.RWMutex{},
		value:                   []common.Address{},
		lastChangeBlock:         0,
		origin:                  &feeTokensOrigin{priceRegistry: priceRegistry},
	}
}

type CachedTokens struct {
	SupportedTokens map[common.Address]common.Address
	FeeTokens       []common.Address
}

// NewCachedSupportedTokens cache both fee tokens and supported tokens. Therefore, it uses 4 different events
// when checking for changes in logpoller.LogPoller
func NewCachedSupportedTokens(
	lp logpoller.LogPoller,
	offRamp ccipdata.OffRampReader,
	priceRegistry ccipdata.PriceRegistryReader,
	optimisticConfirmations int64,
) *CachedChain[CachedTokens] {
	return &CachedChain[CachedTokens]{
		observedEvents:          append(priceRegistry.FeeTokenEvents(), offRamp.TokenEvents()...),
		logPoller:               lp,
		address:                 []common.Address{priceRegistry.Address(), offRamp.Address()},
		optimisticConfirmations: optimisticConfirmations,
		lock:                    &sync.RWMutex{},
		value:                   CachedTokens{},
		lastChangeBlock:         0,
		origin: &feeAndSupportedTokensOrigin{
			feeTokensOrigin:       feeTokensOrigin{priceRegistry: priceRegistry},
			supportedTokensOrigin: supportedTokensOrigin{offRamp: offRamp}},
	}
}

func NewTokenToDecimals(
	lggr logger.Logger,
	lp logpoller.LogPoller,
	offRamp ccipdata.OffRampReader,
	priceRegistryReader ccipdata.PriceRegistryReader,
	client evmclient.Client,
	optimisticConfirmations int64,
) *CachedChain[map[common.Address]uint8] {
	return &CachedChain[map[common.Address]uint8]{
		observedEvents:          append(priceRegistryReader.FeeTokenEvents(), offRamp.TokenEvents()...),
		logPoller:               lp,
		address:                 []common.Address{priceRegistryReader.Address(), offRamp.Address()},
		optimisticConfirmations: optimisticConfirmations,
		lock:                    &sync.RWMutex{},
		value:                   make(map[common.Address]uint8),
		lastChangeBlock:         0,
		origin: &tokenToDecimals{
			lggr:                lggr,
			priceRegistryReader: priceRegistryReader,
			offRamp:             offRamp,
			tokenFactory: func(token common.Address) (link_token_interface.LinkTokenInterface, error) {
				return link_token_interface.NewLinkToken(token, client)
			},
		},
	}
}

type supportedTokensOrigin struct {
	offRamp ccipdata.OffRampReader
}

func (t *supportedTokensOrigin) Copy(value map[common.Address]common.Address) map[common.Address]common.Address {
	return copyMap(value)
}

// CallOrigin Generates the source to dest token mapping based on the offRamp.
// NOTE: this queries the offRamp n+1 times, where n is the number of enabled tokens.
func (t *supportedTokensOrigin) CallOrigin(ctx context.Context) (map[common.Address]common.Address, error) {
	srcToDstTokenMapping := make(map[common.Address]common.Address)
	sourceTokens, err := t.offRamp.GetSupportedTokens(ctx)
	if err != nil {
		return nil, err
	}

	seenDestinationTokens := make(map[common.Address]struct{})

	for _, sourceToken := range sourceTokens {
		dst, err1 := t.offRamp.GetDestinationToken(ctx, sourceToken)
		if err1 != nil {
			return nil, err1
		}

		if _, exists := seenDestinationTokens[dst]; exists {
			return nil, fmt.Errorf("offRamp misconfig, destination token %s already exists", dst)
		}

		seenDestinationTokens[dst] = struct{}{}
		srcToDstTokenMapping[sourceToken] = dst
	}
	return srcToDstTokenMapping, nil
}

type feeTokensOrigin struct {
	priceRegistry ccipdata.PriceRegistryReader
}

func (t *feeTokensOrigin) Copy(value []common.Address) []common.Address {
	return copyArray(value)
}

func (t *feeTokensOrigin) CallOrigin(ctx context.Context) ([]common.Address, error) {
	return t.priceRegistry.GetFeeTokens(ctx)
}

func copyArray(source []common.Address) []common.Address {
	dst := make([]common.Address, len(source))
	copy(dst, source)
	return dst
}

type feeAndSupportedTokensOrigin struct {
	feeTokensOrigin       feeTokensOrigin
	supportedTokensOrigin supportedTokensOrigin
}

func (t *feeAndSupportedTokensOrigin) Copy(value CachedTokens) CachedTokens {
	return CachedTokens{
		SupportedTokens: t.supportedTokensOrigin.Copy(value.SupportedTokens),
		FeeTokens:       t.feeTokensOrigin.Copy(value.FeeTokens),
	}
}

func (t *feeAndSupportedTokensOrigin) CallOrigin(ctx context.Context) (CachedTokens, error) {
	supportedTokens, err := t.supportedTokensOrigin.CallOrigin(ctx)
	if err != nil {
		return CachedTokens{}, err
	}
	feeToken, err := t.feeTokensOrigin.CallOrigin(ctx)
	if err != nil {
		return CachedTokens{}, err
	}
	return CachedTokens{
		SupportedTokens: supportedTokens,
		FeeTokens:       feeToken,
	}, nil
}

func copyMap[M ~map[K]V, K comparable, V any](m M) M {
	cpy := make(M)
	for k, v := range m {
		cpy[k] = v
	}
	return cpy
}

type tokenToDecimals struct {
	lggr                logger.Logger
	offRamp             ccipdata.OffRampReader
	priceRegistryReader ccipdata.PriceRegistryReader
	tokenFactory        func(address common.Address) (link_token_interface.LinkTokenInterface, error)
	tokenDecimals       sync.Map
}

func (t *tokenToDecimals) Copy(value map[common.Address]uint8) map[common.Address]uint8 {
	return copyMap(value)
}

// CallOrigin Generates the token to decimal mapping for dest tokens and fee tokens.
// NOTE: this queries token decimals n times, where n is the number of tokens whose decimals are not already cached.
func (t *tokenToDecimals) CallOrigin(ctx context.Context) (map[common.Address]uint8, error) {
	destTokens, err := getDestinationAndFeeTokens(ctx, t.offRamp, t.priceRegistryReader)
	if err != nil {
		return nil, err
	}

	mapping := make(map[common.Address]uint8, len(destTokens))
	for _, token := range destTokens {
		if decimals, exists := t.getCachedDecimals(token); exists {
			mapping[token] = decimals
			continue
		}

		tokenContract, err := t.tokenFactory(token)
		if err != nil {
			return nil, err
		}

		decimals, err := tokenContract.Decimals(&bind.CallOpts{Context: ctx})
		if err != nil {
			return nil, fmt.Errorf("get token %s decimals: %w", token, err)
		}

		t.setCachedDecimals(token, decimals)
		mapping[token] = decimals
	}
	return mapping, nil
}

func getDestinationAndFeeTokens(ctx context.Context, offRamp ccipdata.OffRampReader, priceRegistry ccipdata.PriceRegistryReader) ([]common.Address, error) {
	destTokens, err := offRamp.GetDestinationTokens(ctx)
	if err != nil {
		return nil, err
	}

	feeTokens, err := priceRegistry.GetFeeTokens(ctx)
	if err != nil {
		return nil, err
	}

	for _, feeToken := range feeTokens {
		if !slices.Contains(destTokens, feeToken) {
			destTokens = append(destTokens, feeToken)
		}
	}

	return destTokens, nil
}

func (t *tokenToDecimals) getCachedDecimals(token common.Address) (uint8, bool) {
	rawVal, exists := t.tokenDecimals.Load(token.String())
	if !exists {
		return 0, false
	}

	decimals, isUint8 := rawVal.(uint8)
	if !isUint8 {
		return 0, false
	}

	return decimals, true
}

func (t *tokenToDecimals) setCachedDecimals(token common.Address, decimals uint8) {
	t.tokenDecimals.Store(token.String(), decimals)
}
