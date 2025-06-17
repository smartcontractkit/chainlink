package batchreader

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"

	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/type_and_version"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_4_0"
)

var (
	typeAndVersionABI = abihelpers.MustParseABI(type_and_version.ITypeAndVersionABI)
)

type EVMTokenPoolBatchedReader struct {
	lggr                logger.Logger
	remoteChainSelector uint64
	offRampAddress      common.Address
	evmBatchCaller      rpclib.EvmBatchCaller

	tokenPoolReaders  map[cciptypes.Address]ccipdata.TokenPoolReader
	tokenPoolReaderMu sync.RWMutex
}

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

	results, err := br.evmBatchCaller.BatchCall(ctx, 0, evmCalls)
	if err != nil {
		return nil, fmt.Errorf("batch call limit: %w", err)
	}

	resultsParsed, err := rpclib.ParseOutputs[cciptypes.TokenBucketRateLimit](results, func(d rpclib.DataAndErr) (cciptypes.TokenBucketRateLimit, error) {
		return rpclib.ParseOutput[cciptypes.TokenBucketRateLimit](d, 0)
	})
	if err != nil {
		return nil, fmt.Errorf("parse outputs: %w", err)
	}
	return resultsParsed, nil
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

	results, err := evmBatchCaller.BatchCall(ctx, 0, evmCalls)
	if err != nil {
		return nil, fmt.Errorf("batch call limit: %w", err)
	}

	result, err := rpclib.ParseOutputs[string](results, func(d rpclib.DataAndErr) (string, error) {
		tAndV, err1 := rpclib.ParseOutput[string](d, 0)
		if err1 != nil {
			// typeAndVersion method do not exist for 1.0 pools. We are going to get an ErrEmptyOutput in that case.
			// Some chains, like the simulated chains, will simply revert with "execution reverted"
			if errors.Is(err1, rpclib.ErrEmptyOutput) || ccipcommon.IsTxRevertError(err1) {
				return "LegacyPool " + ccipdata.V1_0_0, nil
			}
			return "", err1
		}

		return tAndV, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse outputs: %w", err)
	}
	return result, nil
}

func (br *EVMTokenPoolBatchedReader) Close() error {
	return nil
}
