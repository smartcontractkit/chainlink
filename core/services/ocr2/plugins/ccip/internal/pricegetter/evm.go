package pricegetter

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
)

const latestRoundDataMethodName = "latestRoundData"

func init() {
	// Ensure existence of latestRoundData method on the Aggregator contract.
	aggregatorABI, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		panic(err)
	}
	if _, ok := aggregatorABI.Methods[latestRoundDataMethodName]; !ok {
		panic(fmt.Errorf("method %s not found on ABI: %+v", latestRoundDataMethodName, aggregatorABI.Methods))
	}
}

type DynamicPriceGetterClient struct {
	BatchCaller rpclib.EvmBatchCaller
}

func NewDynamicPriceGetterClient(batchCaller rpclib.EvmBatchCaller) DynamicPriceGetterClient {
	return DynamicPriceGetterClient{
		BatchCaller: batchCaller,
	}
}

type DynamicPriceGetter struct {
	cfg           config.DynamicPriceGetterConfig
	evmClients    map[uint64]DynamicPriceGetterClient
	aggregatorAbi abi.ABI
}

func NewDynamicPriceGetterConfig(configJson string) (config.DynamicPriceGetterConfig, error) {
	priceGetterConfig := config.DynamicPriceGetterConfig{}
	err := json.Unmarshal([]byte(configJson), &priceGetterConfig)
	if err != nil {
		return config.DynamicPriceGetterConfig{}, fmt.Errorf("parse dynamic price getter config: %w", err)
	}
	err = priceGetterConfig.Validate()
	if err != nil {
		return config.DynamicPriceGetterConfig{}, fmt.Errorf("validate price getter config: %w", err)
	}
	return priceGetterConfig, nil
}

// NewDynamicPriceGetter build a DynamicPriceGetter from a configuration and a map of chain ID to batch callers.
// A batch caller should be provided for all retrieved prices.
func NewDynamicPriceGetter(cfg config.DynamicPriceGetterConfig, evmClients map[uint64]DynamicPriceGetterClient) (*DynamicPriceGetter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating dynamic price getter config: %w", err)
	}
	aggregatorAbi, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		return nil, fmt.Errorf("parse offchainaggregator abi: %w", err)
	}
	priceGetter := DynamicPriceGetter{cfg, evmClients, aggregatorAbi}
	return &priceGetter, nil
}

// TokenPricesUSD implements the PriceGetter interface.
// It returns static prices stored in the price getter, and batch calls to aggregators (on per chain) for aggregator-based prices.
func (d *DynamicPriceGetter) TokenPricesUSD(ctx context.Context, tokens []cciptypes.Address) (map[cciptypes.Address]*big.Int, error) {
	prices := make(map[cciptypes.Address]*big.Int, len(tokens))

	batchCallsPerChain := make(map[uint64][]rpclib.EvmCall)

	// required to maintain the order of the batched rpc calls for mapping the results
	batchCallsTokensOrder := make(map[uint64][]common.Address)

	evmAddrs, err := ccipcalc.GenericAddrsToEvm(tokens...)
	if err != nil {
		return nil, err
	}

	for _, tk := range evmAddrs {
		if aggCfg, isAgg := d.cfg.AggregatorPrices[tk]; isAgg {
			// Batch calls for aggregator-based token prices (one per chain).
			batchCallsPerChain[aggCfg.ChainID] = append(batchCallsPerChain[aggCfg.ChainID], rpclib.NewEvmCall(
				d.aggregatorAbi,
				latestRoundDataMethodName,
				aggCfg.AggregatorContractAddress,
			))
			batchCallsTokensOrder[aggCfg.ChainID] = append(batchCallsTokensOrder[aggCfg.ChainID], tk)
		} else if staticCfg, isStatic := d.cfg.StaticPrices[tk]; isStatic {
			// Fill static prices.
			prices[ccipcalc.EvmAddrToGeneric(tk)] = staticCfg.Price
		} else {
			return nil, fmt.Errorf("no price resolution rule for token %s", tk.Hex())
		}
	}

	for chainID, batchCalls := range batchCallsPerChain {
		client, exists := d.evmClients[chainID]
		if !exists {
			return nil, fmt.Errorf("evm caller for chain %d not found", chainID)
		}

		evmCaller := client.BatchCaller
		tokensOrder := batchCallsTokensOrder[chainID]

		resultsPerChain, err := evmCaller.BatchCall(ctx, 0, batchCalls)
		if err != nil {
			return nil, fmt.Errorf("batch call: %w", err)
		}

		// latestRoundData function has multiple outputs, we want the second one (idx=1)
		latestRounds, err := rpclib.ParseOutputs[*big.Int](resultsPerChain, func(d rpclib.DataAndErr) (*big.Int, error) {
			return rpclib.ParseOutput[*big.Int](d, 1)
		})
		if err != nil {
			return nil, fmt.Errorf("parse outputs: %w", err)
		}

		for i := range tokensOrder {
			// Prices are already in wei (10e18) when coming from aggregator, no conversion needed.
			prices[ccipcalc.EvmAddrToGeneric(tokensOrder[i])] = latestRounds[i]
		}
	}

	return prices, nil
}
