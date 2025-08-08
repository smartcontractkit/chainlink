package tron

import (
	"context"
	"fmt"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func DeployCache(chain cldf_tron.Chain, labels []string) (*types.DeployTronCacheResponse, error) {

	cacheAddress, txInfo, err := chain.DeployContractAndConfirm(context.Background(), "DataFeedsCache", data_feeds_cache.DataFeedsCacheABI, data_feeds_cache.DataFeedsCacheBin, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkDataFeedsCache: %s, %w", txInfo.ID, err)
	}

	cacheResponse, err := chain.Client.TriggerConstantContract(chain.Address, cacheAddress, "typeAndVersion", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version from %s: %w", cacheAddress, err)
	}

	tv, err := cldf.TypeAndVersionFromString(cacheResponse.ConstantResult[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", cacheResponse.ConstantResult[0], err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployTronCacheResponse{
		Address: cacheAddress,
		Tx:      txInfo.ID,
		Tv:      tv,
	}
	return resp, nil
}
