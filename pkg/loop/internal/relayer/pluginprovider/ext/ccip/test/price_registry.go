package test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type PriceRegistryReaderEvaluator interface {
	ccip.PriceRegistryReader
	testtypes.Evaluator[ccip.PriceRegistryReader]
}

// PriceRegistryReader is a static implementation of the [ccip.PriceRegistryReader] interface.
// It implements [testtypes.Evaluator[ccip.PriceRegistryReader]] and can be use as a test implementation of the [ccip.PriceRegistryReader] interface.
// when grpc serialization is required
var PriceRegistryReader = staticPriceRegistryReader{
	staticPriceRegistryReaderConfig{
		// Address test data
		addressResponse: "static price registry address",

		// GetFeeTokens test data
		getFeeTokensResponse: []ccip.Address{
			"fee token 1 ",
			"fee token 2",
		},

		// GetGasPriceUpdatesCreatedAfter test data
		getGasPriceUpdatesCreatedAfterRequest: getGasPriceUpdatesCreatedAfterRequest{
			chainSelector: 1,
			ts:            time.Unix(179, 13).UTC(),
			confirmations: 2,
		},
		getGasPriceUpdatesCreatedAfterResponse: []ccip.GasPriceUpdateWithTxMeta{
			{
				TxMeta: ccip.TxMeta{
					BlockTimestampUnixMilli: 1,
					BlockNumber:             1,
					TxHash:                  "gas update tx 1",
					LogIndex:                1,
				},
				GasPriceUpdate: ccip.GasPriceUpdate{
					GasPrice: ccip.GasPrice{
						DestChainSelector: 1,
						Value:             big.NewInt(1),
					},
					TimestampUnixSec: big.NewInt(1),
				},
			},
			{
				TxMeta: ccip.TxMeta{
					BlockTimestampUnixMilli: 2,
					BlockNumber:             2,
					TxHash:                  "gas update 2",
					LogIndex:                2,
				},
				GasPriceUpdate: ccip.GasPriceUpdate{
					GasPrice: ccip.GasPrice{
						DestChainSelector: 2,
						Value:             big.NewInt(2),
					},
					TimestampUnixSec: big.NewInt(2),
				},
			},
		},

		// GetAllGasPriceUpdatesCreatedAfter test data
		getAllGasPriceUpdatesCreatedAfterRequest: getAllGasPriceUpdatesCreatedAfterRequest{
			ts:            time.Unix(189, 15).UTC(),
			confirmations: 3,
		},
		getAllGasPriceUpdatesCreatedAfterResponse: []ccip.GasPriceUpdateWithTxMeta{
			{
				TxMeta: ccip.TxMeta{
					BlockTimestampUnixMilli: 10,
					BlockNumber:             10,
					TxHash:                  "gas update tx 10",
					LogIndex:                10,
				},
				GasPriceUpdate: ccip.GasPriceUpdate{
					GasPrice: ccip.GasPrice{
						DestChainSelector: 10,
						Value:             big.NewInt(10),
					},
					TimestampUnixSec: big.NewInt(10),
				},
			},
			{
				TxMeta: ccip.TxMeta{
					BlockTimestampUnixMilli: 20,
					BlockNumber:             20,
					TxHash:                  "gas update 20",
					LogIndex:                20,
				},
				GasPriceUpdate: ccip.GasPriceUpdate{
					GasPrice: ccip.GasPrice{
						DestChainSelector: 20,
						Value:             big.NewInt(20),
					},
					TimestampUnixSec: big.NewInt(20),
				},
			},
		},

		// GetTokenPriceUpdatesCreatedAfter test data
		getTokenPriceUpdatesCreatedAfterRequest: getTokenPriceUpdatesCreatedAfterRequest{
			ts:            time.Unix(111111111, 17).UTC(),
			confirmations: 2,
		},
		getTokenPriceUpdatesCreatedAfterResponse: []ccip.TokenPriceUpdateWithTxMeta{
			{
				TxMeta: ccip.TxMeta{
					BlockTimestampUnixMilli: 1,
					BlockNumber:             1,
					TxHash:                  "token update 1",
					LogIndex:                1,
				},
				TokenPriceUpdate: ccip.TokenPriceUpdate{
					TokenPrice: ccip.TokenPrice{
						Token: ccip.Address("token 1"),
						Value: big.NewInt(1),
					},
					TimestampUnixSec: big.NewInt(1),
				},
			},
		},

		// GetTokenPrices test data
		getTokenPricesRequest: []ccip.Address{
			"token price request 1",
			"token price request 2",
		},
		getTokenPricesResponse: []ccip.TokenPriceUpdate{
			{
				TokenPrice: ccip.TokenPrice{
					Token: ccip.Address("token price request 1"),
					Value: big.NewInt(1),
				},
				TimestampUnixSec: big.NewInt(1),
			},
			{
				TokenPrice: ccip.TokenPrice{
					Token: ccip.Address("token price request 2"),
					Value: big.NewInt(2),
				},
				TimestampUnixSec: big.NewInt(2),
			},
		},

		// GetTokensDecimals test data
		getTokensDecimalsRequest: []ccip.Address{
			"token decimal 1",
			"token decimal 2",
		},
		getTokensDecimalsResponse: []uint8{7, 11},
	},
}

// ensure type implements interface
var _ PriceRegistryReaderEvaluator = staticPriceRegistryReader{}

type staticPriceRegistryReaderConfig struct {
	addressResponse      ccip.Address
	getFeeTokensResponse []ccip.Address
	// handle GetGasPriceUpdatesCreatedAfter
	getGasPriceUpdatesCreatedAfterRequest  getGasPriceUpdatesCreatedAfterRequest
	getGasPriceUpdatesCreatedAfterResponse []ccip.GasPriceUpdateWithTxMeta
	// handle GetAllGasPriceUpdatesCreatedAfter
	getAllGasPriceUpdatesCreatedAfterRequest  getAllGasPriceUpdatesCreatedAfterRequest
	getAllGasPriceUpdatesCreatedAfterResponse []ccip.GasPriceUpdateWithTxMeta
	// handle GetTokenPriceUpdatesCreatedAfter
	getTokenPriceUpdatesCreatedAfterRequest  getTokenPriceUpdatesCreatedAfterRequest
	getTokenPriceUpdatesCreatedAfterResponse []ccip.TokenPriceUpdateWithTxMeta
	// handle GetTokenPrices
	getTokenPricesRequest  []ccip.Address
	getTokenPricesResponse []ccip.TokenPriceUpdate
	// handle GetTokensDecimals
	getTokensDecimalsRequest  []ccip.Address
	getTokensDecimalsResponse []uint8
}

type getGasPriceUpdatesCreatedAfterRequest struct {
	chainSelector uint64
	ts            time.Time
	confirmations int
}

type getAllGasPriceUpdatesCreatedAfterRequest struct {
	ts            time.Time
	confirmations int
}

type getTokenPriceUpdatesCreatedAfterRequest struct {
	ts            time.Time
	confirmations int
}

type staticPriceRegistryReader struct {
	staticPriceRegistryReaderConfig
}

// Evaluate implements types_test.Evaluator.
func (s staticPriceRegistryReader) Evaluate(ctx context.Context, other ccip.PriceRegistryReader) error {
	// Address test case
	gotAddress, err := other.Address(ctx)
	if err != nil {
		return fmt.Errorf("got error on Address: %w", err)
	}
	if gotAddress != s.addressResponse {
		return fmt.Errorf("unexpected Address: want %s, got %s", s.addressResponse, gotAddress)
	}

	// GetFeeTokens test case
	gotFeeTokens, err := other.GetFeeTokens(ctx)
	if err != nil {
		return fmt.Errorf("got error on GetFeeTokens: %w", err)
	}
	if len(gotFeeTokens) != len(s.getFeeTokensResponse) {
		return fmt.Errorf("unexpected number of fee tokens: want %d, got %d", len(s.getFeeTokensResponse), len(gotFeeTokens))
	}

	// GetGasPriceUpdatesCreatedAfter test case
	gotGasPriceUpdates, err := other.GetGasPriceUpdatesCreatedAfter(ctx, s.getGasPriceUpdatesCreatedAfterRequest.chainSelector, s.getGasPriceUpdatesCreatedAfterRequest.ts, s.getGasPriceUpdatesCreatedAfterRequest.confirmations)
	if err != nil {
		return fmt.Errorf("got error on GetGasPriceUpdatesCreatedAfter: %w", err)
	}
	if len(gotGasPriceUpdates) != len(s.getGasPriceUpdatesCreatedAfterResponse) {
		return fmt.Errorf("unexpected number of gas price updates: want %d, got %d", len(s.getGasPriceUpdatesCreatedAfterResponse), len(gotGasPriceUpdates))
	}

	// GetAllGasPriceUpdatesCreatedAfter test case
	gotAllGasPriceUpdates, err := other.GetAllGasPriceUpdatesCreatedAfter(ctx, s.getAllGasPriceUpdatesCreatedAfterRequest.ts, s.getAllGasPriceUpdatesCreatedAfterRequest.confirmations)
	if err != nil {
		return fmt.Errorf("got error on GetAllGasPriceUpdatesCreatedAfter: %w", err)
	}
	if len(gotAllGasPriceUpdates) != len(s.getAllGasPriceUpdatesCreatedAfterResponse) {
		return fmt.Errorf("unexpected number of gas price updates: want %d, got %d", len(s.getAllGasPriceUpdatesCreatedAfterResponse), len(gotAllGasPriceUpdates))
	}

	// GetTokenPriceUpdatesCreatedAfter test case
	gotTokenPriceUpdates, err := other.GetTokenPriceUpdatesCreatedAfter(ctx, s.getTokenPriceUpdatesCreatedAfterRequest.ts, s.getTokenPriceUpdatesCreatedAfterRequest.confirmations)
	if err != nil {
		return fmt.Errorf("got error on GetTokenPriceUpdatesCreatedAfter: %w", err)
	}
	if len(gotTokenPriceUpdates) != len(s.getTokenPriceUpdatesCreatedAfterResponse) {
		return fmt.Errorf("unexpected number of token price updates: want %d, got %d", len(s.getTokenPriceUpdatesCreatedAfterResponse), len(gotTokenPriceUpdates))
	}

	// GetTokenPrices test case
	gotTokenPrices, err := other.GetTokenPrices(ctx, s.getTokenPricesRequest)
	if err != nil {
		return fmt.Errorf("got error on GetTokenPrices: %w", err)
	}
	if len(gotTokenPrices) != len(s.getTokenPricesResponse) {
		return fmt.Errorf("unexpected number of token prices: want %d, got %d", len(s.getTokenPricesResponse), len(gotTokenPrices))
	}

	// GetTokensDecimals test case
	gotTokensDecimals, err := other.GetTokensDecimals(ctx, s.getTokensDecimalsRequest)
	if err != nil {
		return fmt.Errorf("got error on GetTokensDecimals: %w", err)
	}
	if len(gotTokensDecimals) != len(s.getTokensDecimalsResponse) {
		return fmt.Errorf("unexpected number of token decimals: want %d, got %d", len(s.getTokensDecimalsResponse), len(gotTokensDecimals))
	}

	return nil
}

// Address implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) Address(ctx context.Context) (ccip.Address, error) {
	return s.addressResponse, nil
}

// Close implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) Close() error {
	return nil
}

// GetFeeTokens implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) GetFeeTokens(ctx context.Context) ([]ccip.Address, error) {
	return s.getFeeTokensResponse, nil
}

// GetGasPriceUpdatesCreatedAfter implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confirmations int) ([]ccip.GasPriceUpdateWithTxMeta, error) {
	// Check request
	if s.getGasPriceUpdatesCreatedAfterRequest.chainSelector != chainSelector {
		return nil, fmt.Errorf("unexpected chainSelector: want %d, got %d", s.getGasPriceUpdatesCreatedAfterRequest.chainSelector, chainSelector)
	}
	if s.getGasPriceUpdatesCreatedAfterRequest.ts != ts {
		return nil, fmt.Errorf("unexpected ts: want %s, got %s", s.getGasPriceUpdatesCreatedAfterRequest.ts, ts)
	}
	if s.getGasPriceUpdatesCreatedAfterRequest.confirmations != confirmations {
		return nil, fmt.Errorf("unexpected confirmations: want %d, got %d", s.getGasPriceUpdatesCreatedAfterRequest.confirmations, confirmations)
	}

	return s.getGasPriceUpdatesCreatedAfterResponse, nil
}

// GetAllGasPriceUpdatesCreatedAfter implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) GetAllGasPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confirmations int) ([]ccip.GasPriceUpdateWithTxMeta, error) {
	if s.getAllGasPriceUpdatesCreatedAfterRequest.ts != ts {
		return nil, fmt.Errorf("unexpected ts: want %s, got %s", s.getAllGasPriceUpdatesCreatedAfterRequest.ts, ts)
	}
	if s.getAllGasPriceUpdatesCreatedAfterRequest.confirmations != confirmations {
		return nil, fmt.Errorf("unexpected confirmations: want %d, got %d", s.getAllGasPriceUpdatesCreatedAfterRequest.confirmations, confirmations)
	}

	return s.getAllGasPriceUpdatesCreatedAfterResponse, nil
}

// GetTokenPriceUpdatesCreatedAfter implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confirmations int) ([]ccip.TokenPriceUpdateWithTxMeta, error) {
	// Check request
	if s.getTokenPriceUpdatesCreatedAfterRequest.ts != ts {
		return nil, fmt.Errorf("unexpected ts: want %s, got %s", s.getTokenPriceUpdatesCreatedAfterRequest.ts, ts)
	}
	if s.getTokenPriceUpdatesCreatedAfterRequest.confirmations != confirmations {
		return nil, fmt.Errorf("unexpected confirmations: want %d, got %d", s.getTokenPriceUpdatesCreatedAfterRequest.confirmations, confirmations)
	}

	return s.getTokenPriceUpdatesCreatedAfterResponse, nil
}

// GetTokenPrices implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) GetTokenPrices(ctx context.Context, wantedTokens []ccip.Address) ([]ccip.TokenPriceUpdate, error) {
	// Check request
	if len(s.getTokenPricesRequest) != len(wantedTokens) {
		return nil, fmt.Errorf("unexpected number of tokens: want %d, got %d", len(s.getTokenPricesRequest), len(wantedTokens))
	}
	for i, token := range wantedTokens {
		if s.getTokenPricesRequest[i] != token {
			return nil, fmt.Errorf("unexpected token: want %s, got %s", s.getTokenPricesRequest[i], token)
		}
	}

	return s.getTokenPricesResponse, nil
}

// GetTokensDecimals implements ccip.PriceRegistryReader.
func (s staticPriceRegistryReader) GetTokensDecimals(ctx context.Context, tokenAddresses []ccip.Address) ([]uint8, error) {
	// Check request
	if len(s.getTokensDecimalsRequest) != len(tokenAddresses) {
		return nil, fmt.Errorf("unexpected number of tokens: want %d, got %d", len(s.getTokensDecimalsRequest), len(tokenAddresses))
	}
	for i, token := range tokenAddresses {
		if s.getTokensDecimalsRequest[i] != token {
			return nil, fmt.Errorf("unexpected token: want %s, got %s", s.getTokensDecimalsRequest[i], token)
		}
	}

	return s.getTokensDecimalsResponse, nil
}
