package changeset

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
	bundleproxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/bundle_aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func DeployCache(chain cldf_evm.Chain, labels []string) (*types.DeployCacheResponse, error) {
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

	tv, err := cldf.TypeAndVersionFromString(tvStr)
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

func DeployAggregatorProxy(chain cldf_evm.Chain, aggregator common.Address, accessController common.Address, labels []string) (*types.DeployProxyResponse, error) {
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
	tv, err := cldf.TypeAndVersionFromString(tvStr)
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

func DeployBundleAggregatorProxy(lggr logger.Logger, chain cldf_evm.Chain, aggregator common.Address, owner common.Address, labels []string) (*types.DeployBundleAggregatorProxyResponse, error) {
	lggr.Debugw("Deploying BundleAggregatorProxy",
		"chainSelector", chain.Selector,
		"aggregator", aggregator.Hex(),
		"owner", owner.Hex(),
		"deployer", chain.DeployerKey.From.Hex())

	proxyAddr, tx, _, err := bundleproxy.DeployBundleAggregatorProxy(chain.DeployerKey, chain.Client, aggregator, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy BundleAggregatorProxy: %w", err)
	}

	lggr.Debugw("BundleAggregatorProxy deploy tx submitted",
		"chainSelector", chain.Selector,
		"txHash", tx.Hash().Hex(),
		"txNonce", tx.Nonce(),
		"predictedAddress", proxyAddr.Hex())

	blockNum, err := chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm BundleAggregatorProxy: %w", err)
	}

	lggr.Debugw("BundleAggregatorProxy deploy tx confirmed",
		"chainSelector", chain.Selector,
		"txHash", tx.Hash().Hex(),
		"blockNumber", blockNum,
		"predictedAddress", proxyAddr.Hex())

	receipt, err := chain.Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt for BundleAggregatorProxy deploy tx: %w", err)
	}

	receiptAddr := receipt.ContractAddress
	if receiptAddr != proxyAddr {
		lggr.Warnw("BundleAggregatorProxy predicted address does not match receipt address",
			"chainSelector", chain.Selector,
			"predictedAddress", proxyAddr.Hex(),
			"receiptAddress", receiptAddr.Hex(),
			"txHash", tx.Hash().Hex(),
			"txNonce", tx.Nonce(),
			"receiptStatus", receipt.Status)
		proxyAddr = receiptAddr
	}

	code, err := chain.Client.CodeAt(context.Background(), proxyAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check code at BundleAggregatorProxy address %s: %w", proxyAddr, err)
	}
	if len(code) == 0 {
		lggr.Errorw("No contract code found at BundleAggregatorProxy address after confirmed deployment",
			"chainSelector", chain.Selector,
			"address", proxyAddr.Hex(),
			"txHash", tx.Hash().Hex(),
			"txNonce", tx.Nonce(),
			"receiptStatus", receipt.Status,
			"blockNumber", receipt.BlockNumber)
		return nil, fmt.Errorf("no contract code at BundleAggregatorProxy address %s (tx %s, status %d)", proxyAddr, tx.Hash(), receipt.Status)
	}

	lggr.Debugw("BundleAggregatorProxy code verified at address",
		"chainSelector", chain.Selector,
		"address", proxyAddr.Hex(),
		"codeSize", len(code))

	proxyContract, err := bundleproxy.NewBundleAggregatorProxy(proxyAddr, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind BundleAggregatorProxy at %s: %w", proxyAddr, err)
	}

	tvStr, err := proxyContract.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version at %s: %w", proxyAddr, err)
	}

	lggr.Debugw("BundleAggregatorProxy typeAndVersion retrieved",
		"chainSelector", chain.Selector,
		"address", proxyAddr.Hex(),
		"typeAndVersion", tvStr)

	tv, err := cldf.TypeAndVersionFromString(tvStr)
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
