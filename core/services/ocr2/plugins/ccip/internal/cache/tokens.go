package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

// NewCachedFeeTokens cache fee tokens returned from PriceRegistry
func NewCachedFeeTokens(
	lp logpoller.LogPoller,
	priceRegistry ccipdata.PriceRegistryReader,
) *CachedChain[[]common.Address] {
	return &CachedChain[[]common.Address]{
		observedEvents:  priceRegistry.FeeTokenEvents(),
		logPoller:       lp,
		address:         []common.Address{priceRegistry.Address()},
		lock:            &sync.RWMutex{},
		value:           []common.Address{},
		lastChangeBlock: 0,
		origin:          &feeTokensOrigin{priceRegistry: priceRegistry},
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
) *CachedChain[CachedTokens] {
	return &CachedChain[CachedTokens]{
		observedEvents:  append(priceRegistry.FeeTokenEvents(), offRamp.TokenEvents()...),
		logPoller:       lp,
		address:         []common.Address{priceRegistry.Address(), offRamp.Address()},
		lock:            &sync.RWMutex{},
		value:           CachedTokens{},
		lastChangeBlock: 0,
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
) *CachedChain[map[common.Address]uint8] {
	return &CachedChain[map[common.Address]uint8]{
		observedEvents:  append(priceRegistryReader.FeeTokenEvents(), offRamp.TokenEvents()...),
		logPoller:       lp,
		address:         []common.Address{priceRegistryReader.Address(), offRamp.Address()},
		lock:            &sync.RWMutex{},
		value:           make(map[common.Address]uint8),
		lastChangeBlock: 0,
		origin: &tokenToDecimals{
			lggr:                lggr,
			priceRegistryReader: priceRegistryReader,
			offRamp:             offRamp,
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
	sourceTokens, err := t.offRamp.GetSupportedTokens(ctx)
	if err != nil {
		return nil, err
	}

	destTokens, err := t.offRamp.GetDestinationTokensFromSourceTokens(ctx, sourceTokens)
	if err != nil {
		return nil, fmt.Errorf("get destination tokens from source tokens: %w", err)
	}

	srcToDstTokenMapping := make(map[common.Address]common.Address, len(sourceTokens))
	for i, sourceToken := range sourceTokens {
		srcToDstTokenMapping[sourceToken] = destTokens[i]
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
	unknownDecimalsTokens := make([]common.Address, 0, len(destTokens))

	for _, token := range destTokens {
		if decimals, exists := t.getCachedDecimals(token); exists {
			mapping[token] = decimals
			continue
		}
		unknownDecimalsTokens = append(unknownDecimalsTokens, token)
	}

	if len(unknownDecimalsTokens) == 0 {
		return mapping, nil
	}

	decimals, err := t.priceRegistryReader.GetTokensDecimals(ctx, unknownDecimalsTokens)
	if err != nil {
		return nil, fmt.Errorf("get tokens decimals: %w", err)
	}
	for i := range decimals {
		t.setCachedDecimals(unknownDecimalsTokens[i], decimals[i])
		mapping[unknownDecimalsTokens[i]] = decimals[i]
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
