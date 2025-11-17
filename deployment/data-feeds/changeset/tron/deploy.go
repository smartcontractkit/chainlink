package tron

import (
	"context"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

const (
	DeploymentBlockLabel = "deployment-block"
	DeploymentHashLabel  = "deployment-hash"
)

func DeployCache(chain cldf_tron.Chain, deployOptions *cldf_tron.DeployOptions, labels []string) (*types.DeployTronResponse, error) {
	cacheAddress, txInfo, err := chain.DeployContractAndConfirm(context.Background(), "DataFeedsCache", cache.DataFeedsCacheABI, cache.DataFeedsCacheBin, nil, deployOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm ChainlinkDataFeedsCache: %+v, %w", txInfo, err)
	}

	cacheResponse, err := chain.Client.TriggerConstantContract(chain.Address, cacheAddress, "typeAndVersion()", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version from %s: %w", cacheAddress, err)
	}

	typeAndVersion, err := changeset.ExtractTypeAndVersion(cacheResponse.ConstantResult[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode type and version from %s: %w", cacheResponse.ConstantResult[0], err)
	}

	tv, err := cldf.TypeAndVersionFromString(typeAndVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", typeAndVersion, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	tv.Labels.Add(fmt.Sprintf("%s: %s", DeploymentHashLabel, txInfo.ID))
	tv.Labels.Add(fmt.Sprintf("%s: %d", DeploymentBlockLabel, txInfo.BlockNumber))

	resp := &types.DeployTronResponse{
		Address: cacheAddress,
		Tx:      txInfo.ID,
		Tv:      tv,
	}
	return resp, nil
}

func DeployAggregatorProxy(chain cldf_tron.Chain, aggregator address.Address, accessController address.Address, deployOptions *cldf_tron.DeployOptions, labels []string) (*types.DeployTronResponse, error) {
	proxyAddress, txInfo, err := chain.DeployContractAndConfirm(context.Background(), "AggregatorProxy", proxy.AggregatorProxyABI, proxy.AggregatorProxyBin, []any{aggregator.EthAddress(), accessController.EthAddress()}, deployOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AggregatorProxy: %+v, %w", txInfo, err)
	}

	// AggregatorProxy contract doesn't implement typeAndVersion interface, so we have to set it manually
	tvStr := "AggregatorProxy 1.0.0"
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	tv.Labels.Add(fmt.Sprintf("%s: %s", DeploymentHashLabel, txInfo.ID))
	tv.Labels.Add(fmt.Sprintf("%s: %d", DeploymentBlockLabel, txInfo.BlockNumber))

	resp := &types.DeployTronResponse{
		Address: proxyAddress,
		Tx:      txInfo.ID,
		Tv:      tv,
	}
	return resp, nil
}
