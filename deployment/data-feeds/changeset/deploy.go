package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
	bundleproxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/bundle_aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func DeployCache(chain deployment.Chain, labels []string) (*types.DeployCacheResponse, error) {
	cacheAddr, tx, cacheContract, err := cache.DeployDataFeedsCache(chain.DeployerKey, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy DataFeedsCache: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm DataFeedsCache: %w", err)
	}

	tvStr, err := cacheContract.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}

	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployCacheResponse{
		Address:  cacheAddr,
		Tx:       tx.Hash(),
		Tv:       tv,
		Contract: cacheContract,
	}
	return resp, nil
}

func DeployAggregatorProxy(chain deployment.Chain, aggregator common.Address, accessController common.Address, labels []string) (*types.DeployProxyResponse, error) {
	proxyAddr, tx, proxyContract, err := proxy.DeployAggregatorProxy(chain.DeployerKey, chain.Client, aggregator, accessController)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy AggregatorProxy: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AggregatorProxy: %w", err)
	}

	// AggregatorProxy contract doesn't implement typeAndVersion interface, so we have to set it manually
	tvStr := "AggregatorProxy 1.0.0"
	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployProxyResponse{
		Address:  proxyAddr,
		Tx:       tx.Hash(),
		Tv:       tv,
		Contract: proxyContract,
	}
	return resp, nil
}

func DeployBundleAggregatorProxy(chain deployment.Chain, aggregator common.Address, owner common.Address, labels []string) (*types.DeployBundleAggregatorProxyResponse, error) {
	proxyAddr, tx, proxyContract, err := bundleproxy.DeployBundleAggregatorProxy(chain.DeployerKey, chain.Client, aggregator, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy BundleAggregatorProxy: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm BundleAggregatorProxy: %w", err)
	}

	tvStr, err := proxyContract.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}
	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	for _, label := range labels {
		tv.Labels.Add(label)
	}

	resp := &types.DeployBundleAggregatorProxyResponse{
		Address:  proxyAddr,
		Tx:       tx.Hash(),
		Tv:       tv,
		Contract: proxyContract,
	}
	return resp, nil
}
