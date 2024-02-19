package batchreader

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	type_and_version "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/type_and_version_interface_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_4_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
)

var (
	typeAndVersionABI = abihelpers.MustParseABI(type_and_version.TypeAndVersionInterfaceABI)
)

type EVMTokenPoolBatchedReader struct {
	lggr                logger.Logger
	remoteChainSelector uint64
	offRampAddress      common.Address
	evmBatchCaller      rpclib.EvmBatchCaller

	tokenPoolReaders  map[cciptypes.Address]ccipdata.TokenPoolReader
	tokenPoolReaderMu sync.RWMutex
}

//go:generate mockery --quiet --name TokenPoolBatchedReader --filename token_pool_batched_reader_mock.go --case=underscor
type TokenPoolBatchedReader interface {
	cciptypes.TokenPoolBatchedReader
}

var _ TokenPoolBatchedReader = (*EVMTokenPoolBatchedReader)(nil)

func NewEVMTokenPoolBatchedReader(lggr logger.Logger, remoteChainSelector uint64, offRampAddress cciptypes.Address, evmBatchCaller rpclib.EvmBatchCaller) (*EVMTokenPoolBatchedReader, error) {
	offRampAddrEvm, err := ccipcalc.GenericAddrToEvm(offRampAddress)
	if err != nil {
		return nil, err
	}

	return &EVMTokenPoolBatchedReader{
		lggr:                lggr,
		remoteChainSelector: remoteChainSelector,
		offRampAddress:      offRampAddrEvm,
		evmBatchCaller:      evmBatchCaller,
		tokenPoolReaders:    make(map[cciptypes.Address]ccipdata.TokenPoolReader),
	}, nil
}

func (br *EVMTokenPoolBatchedReader) GetInboundTokenPoolRateLimits(ctx context.Context, tokenPools []cciptypes.Address) ([]cciptypes.TokenBucketRateLimit, error) {
	if len(tokenPools) == 0 {
		return []cciptypes.TokenBucketRateLimit{}, nil
	}

	err := br.loadTokenPoolReaders(ctx, tokenPools)
	if err != nil {
		return nil, err
	}

	tokenPoolReaders := make([]ccipdata.TokenPoolReader, 0, len(tokenPools))
	for _, poolAddress := range tokenPools {
		br.tokenPoolReaderMu.RLock()
		tokenPoolReader, exists := br.tokenPoolReaders[poolAddress]
		br.tokenPoolReaderMu.RUnlock()
		if !exists {
			return nil, fmt.Errorf("token pool %s not found", poolAddress)
		}
		tokenPoolReaders = append(tokenPoolReaders, tokenPoolReader)
	}

	evmCalls := make([]rpclib.EvmCall, 0, len(tokenPoolReaders))
	for _, poolReader := range tokenPoolReaders {
		switch v := poolReader.(type) {
		case *v1_2_0.TokenPool:
			evmCalls = append(evmCalls, v1_2_0.GetInboundTokenPoolRateLimitCall(v.Address(), v.OffRampAddress))
		case *v1_4_0.TokenPool:
			evmCalls = append(evmCalls, v1_4_0.GetInboundTokenPoolRateLimitCall(v.Address(), v.RemoteChainSelector))
		default:
			return nil, fmt.Errorf("unsupported token pool version %T", v)
		}
	}

	return batchCallLatestBlockNumber[cciptypes.TokenBucketRateLimit](ctx, br.evmBatchCaller, evmCalls)
}

// loadTokenPoolReaders loads the token pools into the factory's cache
func (br *EVMTokenPoolBatchedReader) loadTokenPoolReaders(ctx context.Context, tokenPoolAddresses []cciptypes.Address) error {
	var missingTokens []common.Address

	br.tokenPoolReaderMu.RLock()
	for _, poolAddress := range tokenPoolAddresses {
		if _, exists := br.tokenPoolReaders[poolAddress]; !exists {
			evmPoolAddr, err := ccipcalc.GenericAddrToEvm(poolAddress)
			if err != nil {
				return err
			}
			missingTokens = append(missingTokens, evmPoolAddr)
		}
	}
	br.tokenPoolReaderMu.RUnlock()

	// Only continue if there are missing tokens
	if len(missingTokens) == 0 {
		return nil
	}

	typeAndVersions, err := getBatchedTypeAndVersion(ctx, br.evmBatchCaller, missingTokens)
	if err != nil {
		return err
	}

	br.tokenPoolReaderMu.Lock()
	defer br.tokenPoolReaderMu.Unlock()
	for i, tokenPoolAddress := range missingTokens {
		typeAndVersion := typeAndVersions[i]
		poolType, version, err := ccipconfig.ParseTypeAndVersion(typeAndVersion)
		if err != nil {
			return err
		}
		switch version {
		case ccipdata.V1_0_0, ccipdata.V1_1_0, ccipdata.V1_2_0:
			br.tokenPoolReaders[ccipcalc.EvmAddrToGeneric(tokenPoolAddress)] = v1_2_0.NewTokenPool(poolType, tokenPoolAddress, br.offRampAddress)
		case ccipdata.V1_4_0:
			br.tokenPoolReaders[ccipcalc.EvmAddrToGeneric(tokenPoolAddress)] = v1_4_0.NewTokenPool(poolType, tokenPoolAddress, br.remoteChainSelector)
		default:
			return fmt.Errorf("unsupported token pool version %v", version)
		}
	}
	return nil
}

func getBatchedTypeAndVersion(ctx context.Context, evmBatchCaller rpclib.EvmBatchCaller, poolAddresses []common.Address) ([]string, error) {
	var evmCalls []rpclib.EvmCall

	for _, poolAddress := range poolAddresses {
		// Add the typeAndVersion call to the batch
		evmCalls = append(evmCalls, rpclib.NewEvmCall(
			typeAndVersionABI,
			"typeAndVersion",
			poolAddress,
		))
	}

	return batchCallLatestBlockNumber[string](ctx, evmBatchCaller, evmCalls)
}

func batchCallLatestBlockNumber[T any](ctx context.Context, evmBatchCaller rpclib.EvmBatchCaller, evmCalls []rpclib.EvmCall) ([]T, error) {
	results, err := evmBatchCaller.BatchCall(ctx, 0, evmCalls)
	if err != nil {
		return nil, fmt.Errorf("batch call limit: %w", err)
	}

	result, err := rpclib.ParseOutputs[T](results, func(d rpclib.DataAndErr) (T, error) {
		return rpclib.ParseOutput[T](d, 0)
	})
	if err != nil {
		return nil, fmt.Errorf("parse outputs: %w", err)
	}
	return result, nil
}
