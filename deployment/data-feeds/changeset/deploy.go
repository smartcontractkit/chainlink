package changeset

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sethvargo/go-retry"

	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
	bundleproxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/bundle_aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func waitForContractCode(ctx context.Context, client cldf_evm.OnchainClient, tx *ethtypes.Transaction) (common.Address, error) {
	receipt, err := client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get tx receipt: %w", err)
	}
	addr := receipt.ContractAddress

	err = retry.Do(ctx, retry.WithMaxDuration(90*time.Second, retry.WithCappedDuration(30*time.Second, retry.NewFibonacci(5*time.Second))), func(ctx context.Context) error {
		code, err := client.CodeAt(ctx, addr, receipt.BlockNumber)
		if err != nil {
			return retry.RetryableError(err)
		}
		if len(code) == 0 {
			return retry.RetryableError(fmt.Errorf("no contract code at %s (block %s) yet", addr, receipt.BlockNumber))
		}
		return nil
	})
	return addr, err
}

func DeployCache(chain cldf_evm.Chain, labels []string) (*types.DeployCacheResponse, error) {
	cacheAddr, tx, cacheContract, err := cache.DeployDataFeedsCache(chain.DeployerKey, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy DataFeedsCache: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm DataFeedsCache: %w", err)
	}

	if _, err := waitForContractCode(context.Background(), chain.Client, tx); err != nil {
		return nil, fmt.Errorf("failed to verify DataFeedsCache deployment: %w", err)
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

	if _, err := waitForContractCode(context.Background(), chain.Client, tx); err != nil {
		return nil, fmt.Errorf("failed to verify AggregatorProxy deployment: %w", err)
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

	deployedAddr, err := waitForContractCode(context.Background(), chain.Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify BundleAggregatorProxy deployment: %w", err)
	}

	if deployedAddr != proxyAddr {
		lggr.Warnw("BundleAggregatorProxy predicted address does not match deployed address",
			"chainSelector", chain.Selector,
			"predictedAddress", proxyAddr.Hex(),
			"deployedAddress", deployedAddr.Hex(),
			"txHash", tx.Hash().Hex())
		proxyAddr = deployedAddr
	}

	lggr.Debugw("BundleAggregatorProxy deployed and code verified",
		"chainSelector", chain.Selector,
		"address", proxyAddr.Hex())

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
