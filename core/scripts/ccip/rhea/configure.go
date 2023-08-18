package rhea

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/shared"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/router"
)

var (
	zeroAddress = common.HexToAddress("0x0")
)

func setOffRampOnTokenPools(t *testing.T, client EvmConfig, lane *EVMLaneConfig) {
	for token, tokenConfig := range client.ChainConfig.SupportedTokens {
		if tokenConfig.TokenPoolType == FeeTokenOnly {
			continue
		}
		pool, err := lock_release_token_pool.NewLockReleaseTokenPool(tokenConfig.Pool, client.Client)
		shared.RequireNoError(t, err)

		rampUpdate := lock_release_token_pool.TokenPoolRampUpdate{
			Ramp:    lane.OffRamp,
			Allowed: true,
			RateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(token.Multiplier(), big.NewInt(RATE_LIMIT_CAPACITY_DOLLAR)),
				Rate:      new(big.Int).Mul(token.Multiplier(), big.NewInt(RATE_LIMIT_RATE_DOLLAR)),
			},
		}

		// Configure offramp address on pool
		tx, err := pool.ApplyRampUpdates(client.Owner, []lock_release_token_pool.TokenPoolRampUpdate{}, []lock_release_token_pool.TokenPoolRampUpdate{rampUpdate})
		shared.RequireNoError(t, err)
		err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
		shared.RequireNoError(t, err)
		client.Logger.Infof("Offramp pool configured with offramp address: %s", helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	}
}

func setPriceRegistryPrices(t *testing.T, client *EvmDeploymentConfig, destChainSelector uint64) {
	priceRegistry, err := price_registry.NewPriceRegistry(client.ChainConfig.PriceRegistry, client.Client)
	shared.RequireNoError(t, err)

	priceUpdates := price_registry.InternalPriceUpdates{
		TokenPriceUpdates: []price_registry.InternalTokenPriceUpdate{},
		DestChainSelector: destChainSelector,
		// Set 1e18 units of gas to $2k, being fairly reasonable for eth
		// These values will get auto updated by the DON
		UsdPerUnitGas: big.NewInt(2000e9), // $2000 per eth * 1gwei = 2000e9
	}

	for _, tokenConfig := range client.ChainConfig.SupportedTokens {

		priceUpdates.TokenPriceUpdates = append(priceUpdates.TokenPriceUpdates, price_registry.InternalTokenPriceUpdate{
			SourceToken: tokenConfig.Token,
			// Use price per 1e18 units when updating PriceRegistry to embed decimal info into price
			UsdPerToken: GetPricePer1e18Units(tokenConfig.Price, tokenConfig.Decimals),
		})
	}

	tx, err := priceRegistry.UpdatePrices(client.Owner, priceUpdates)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
}

func attachUpgradeOnRampsToRouter(t *testing.T, sourceClient *EvmDeploymentConfig, destChainSelector uint64) {
	if sourceClient.UpgradeLaneConfig.OnRamp == zeroAddress {
		sourceClient.Logger.Infof("There is no OnRamp on upgrade lane (pending deployment). Skipping")
		return
	}

	routerAddress := sourceClient.ChainConfig.Router
	upgradeOnRampAddress := sourceClient.UpgradeLaneConfig.OnRamp
	// Source Router --xxx--> OnRamp
	// Source Router ---> Upgrade OnRamp
	setOnRampRouter(t, routerAddress, upgradeOnRampAddress, sourceClient, destChainSelector)

	upgradeRouterAddress := sourceClient.ChainConfig.UpgradeRouter
	// Upgrade Source Router ---> nil
	setOnRampRouter(t, upgradeRouterAddress, zeroAddress, sourceClient, destChainSelector)

	// Updating configuration after finishing deployment
	sourceClient.LaneConfig.OnRamp = upgradeOnRampAddress
	sourceClient.UpgradeLaneConfig.OnRamp = zeroAddress
	sourceClient.LaneConfig.DeploySettings.DeployedAtBlock = sourceClient.UpgradeLaneConfig.DeploySettings.DeployedAtBlock
}

func attachUpgradeOffRampsToRouter(t *testing.T, destinationClient *EvmDeploymentConfig, sourceChainSelector uint64) {
	if destinationClient.UpgradeLaneConfig.OffRamp == zeroAddress {
		destinationClient.Logger.Infof("There is no OffRamp on upgrade lane (pending deployment). Skipping")
		return
	}

	routerAddress := destinationClient.ChainConfig.Router
	upgradeOffRampAddress := destinationClient.UpgradeLaneConfig.OffRamp
	// Upgrade OffRamp ----> Router
	setOffRampOnRouter(t, routerAddress, upgradeOffRampAddress, destinationClient, sourceChainSelector)

	upgradeRouterAddress := destinationClient.ChainConfig.UpgradeRouter
	// Upgrade OffRamp --XXX--> Upgrade Router
	removeOffRampFromRouter(t, upgradeRouterAddress, upgradeOffRampAddress, destinationClient, sourceChainSelector)

	// Updating configuration after finishing deployment
	destinationClient.LaneConfig.OffRamp = upgradeOffRampAddress
	destinationClient.UpgradeLaneConfig.OffRamp = zeroAddress

	// Updating CommitStore config as well
	if destinationClient.UpgradeLaneConfig.CommitStore != zeroAddress {
		destinationClient.LaneConfig.CommitStore = destinationClient.UpgradeLaneConfig.CommitStore
		destinationClient.UpgradeLaneConfig.CommitStore = zeroAddress
	}
}

func setOnRampRouter(t *testing.T, routerAddress common.Address, onRampAddress common.Address, sourceClient *EvmDeploymentConfig, destChainSelector uint64) {
	sourceClient.Logger.Infof("Setting the onRamp on the Router")
	routerContract, err := router.NewRouter(routerAddress, sourceClient.Client)
	shared.RequireNoError(t, err)
	sourceClient.Logger.Infof("Registering new onRamp")
	tx, err := routerContract.ApplyRampUpdates(sourceClient.Owner, []router.RouterOnRamp{
		{DestChainSelector: destChainSelector, OnRamp: onRampAddress}}, nil, nil)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(sourceClient.Logger, sourceClient.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
}

func setOnRampOnTokenPools(t *testing.T, sourceClient *EvmDeploymentConfig, onRampAddress common.Address) {
	for _, tokenConfig := range sourceClient.ChainConfig.SupportedTokens {
		if tokenConfig.TokenPoolType == FeeTokenOnly {
			continue
		}
		pool, err := lock_release_token_pool.NewLockReleaseTokenPool(tokenConfig.Pool, sourceClient.Client)
		shared.RequireNoError(t, err)

		rampUpdate := lock_release_token_pool.TokenPoolRampUpdate{
			Ramp:    onRampAddress,
			Allowed: true,
			RateLimiterConfig: lock_release_token_pool.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
				Rate:      new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)),
			},
		}

		// Configure offramp address on pool
		tx, err := pool.ApplyRampUpdates(sourceClient.Owner, []lock_release_token_pool.TokenPoolRampUpdate{rampUpdate}, []lock_release_token_pool.TokenPoolRampUpdate{})
		shared.RequireNoError(t, err)
		err = shared.WaitForMined(sourceClient.Logger, sourceClient.Client, tx.Hash(), true)
		shared.RequireNoError(t, err)
		sourceClient.Logger.Infof("Onramp pool configured with offramp address: %s", helpers.ExplorerLink(int64(sourceClient.ChainConfig.EvmChainId), tx.Hash()))
	}
}

func setOffRampOnRouter(t *testing.T, routerAddress common.Address, offRampAddress common.Address, client *EvmDeploymentConfig, sourceChainSelector uint64) {
	client.Logger.Infof("Setting the offRamp on the Router")
	routerContract, err := router.NewRouter(routerAddress, client.Client)
	shared.RequireNoError(t, err)

	offRamps, err := routerContract.GetOffRamps(&bind.CallOpts{})
	shared.RequireNoError(t, err)
	for _, offRamp := range offRamps {
		if offRamp.OffRamp == offRampAddress {
			client.Logger.Infof("OffRamp already configured on router. Skipping")
			return
		}
	}

	tx, err := routerContract.ApplyRampUpdates(client.Owner, nil, nil, []router.RouterOffRamp{
		{SourceChainSelector: sourceChainSelector, OffRamp: offRampAddress}})
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
}

func removeOffRampFromRouter(t *testing.T, routerAddress common.Address, offRampAddress common.Address, client *EvmDeploymentConfig, sourceChainSelector uint64) {
	client.Logger.Infof("Removing the offRamp freom the Router")
	routerContract, err := router.NewRouter(routerAddress, client.Client)
	shared.RequireNoError(t, err)

	offRamps, err := routerContract.GetOffRamps(&bind.CallOpts{})
	shared.RequireNoError(t, err)

	offRampRegistered := false
	for _, offRamp := range offRamps {
		if offRamp.OffRamp == offRampAddress {
			offRampRegistered = true
		}
	}

	if !offRampRegistered {
		client.Logger.Infof("OffRamp not configured on router. Skipping")
		return
	}

	tx, err := routerContract.ApplyRampUpdates(client.Owner, nil, []router.RouterOffRamp{
		{SourceChainSelector: sourceChainSelector, OffRamp: client.LaneConfig.OffRamp}}, nil)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
}

func setPriceRegistryUpdater(t *testing.T, client *EvmDeploymentConfig) {
	priceRegistry, err := price_registry.NewPriceRegistry(client.ChainConfig.PriceRegistry, client.Client)
	shared.RequireNoError(t, err)

	tx, err := priceRegistry.ApplyPriceUpdatersUpdates(client.Owner, []common.Address{client.LaneConfig.CommitStore}, []common.Address{})
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
}

func fillPoolWithTokens(client *EvmDeploymentConfig, pool *lock_release_token_pool.LockReleaseTokenPool, tokenAddress common.Address, tokenName Token) error {
	token, err := burn_mint_erc677.NewBurnMintERC677(tokenAddress, client.Client)
	if err != nil {
		return err
	}

	// fill offramp token pool with 0.5 token
	amount := new(big.Int).Div(tokenName.Multiplier(), big.NewInt(2))
	tx, err := token.Approve(client.Owner, pool.Address(), amount)
	if err != nil {
		return err
	}
	client.Logger.Infof("Approving token to the token pool: %s", helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	tx, err = pool.AddLiquidity(client.Owner, amount)
	if err != nil {
		return err
	}
	client.Logger.Infof("Adding liquidity to the token pool: %s", helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	client.Logger.Infof("Pool filled with tokens: %s", helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	return nil
}

func setPoolAllowList(client *EvmDeploymentConfig, poolAddress common.Address, poolAllowList []common.Address, tokenName Token) error {
	pool, err := lock_release_token_pool.NewLockReleaseTokenPool(poolAddress, client.Client)
	if err != nil {
		return err
	}

	isEnabled, err := pool.GetAllowListEnabled(&bind.CallOpts{})
	if err != nil {
		return err
	}

	if !isEnabled {
		if len(poolAllowList) > 0 {
			return fmt.Errorf("%s pool does not have allowList enabled, but has allowList defined in config", tokenName)
		}
		return nil
	}

	currentAllowList, err := pool.GetAllowList(&bind.CallOpts{})
	if err != nil {
		return err
	}

	var toRemove []common.Address
	for _, addr := range currentAllowList {
		if !slices.Contains(poolAllowList, addr) {
			toRemove = append(toRemove, addr)
		}
	}
	var toAdd []common.Address
	for _, addr := range poolAllowList {
		if !slices.Contains(currentAllowList, addr) {
			toAdd = append(toAdd, addr)
		}
	}
	if len(toRemove) == 0 && len(toAdd) == 0 {
		client.Logger.Infof("Nothing to add or remove from allowlist for %s pool", tokenName)
		return nil
	}

	client.Logger.Infof("ApplyAllowListUpdates for %s pool: toRemove=%v, toAdd=%v", tokenName, toRemove, toAdd)

	tx, err := pool.ApplyAllowListUpdates(client.Owner, toRemove, toAdd)
	if err != nil {
		return err
	}
	return shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
}

func FundPingPong(client EvmConfig, lane *EVMLaneConfig, fundingAmount *big.Int, tokenAddress common.Address) error {
	linkToken, err := burn_mint_erc677.NewBurnMintERC677(tokenAddress, client.Client)
	if err != nil {
		return err
	}

	tx, err := linkToken.Transfer(client.Owner, lane.PingPongDapp, fundingAmount)
	if err != nil {
		return err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	client.Logger.Infof("Ping pong funded with %s in tx: %s", fundingAmount.String(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	return nil
}

func UpdateDeployedAt(t *testing.T, source *EvmDeploymentConfig, dest *EvmDeploymentConfig) {
	sourceBlock, err := source.Client.BlockNumber(context.Background())
	require.NoError(t, err)

	source.ChainConfig.DeploySettings.DeployedAtBlock = sourceBlock
	source.LaneConfig.DeploySettings.DeployedAtBlock = sourceBlock

	destBlock, err := dest.Client.BlockNumber(context.Background())
	require.NoError(t, err)

	dest.ChainConfig.DeploySettings.DeployedAtBlock = destBlock
	dest.LaneConfig.DeploySettings.DeployedAtBlock = destBlock
}

func UpdateDeployedAtUpgradeLane(t *testing.T, source *EvmDeploymentConfig, dest *EvmDeploymentConfig) {
	sourceBlock, err := source.Client.BlockNumber(context.Background())
	require.NoError(t, err)

	source.UpgradeLaneConfig.DeploySettings.DeployedAtBlock = sourceBlock

	destBlock, err := dest.Client.BlockNumber(context.Background())
	require.NoError(t, err)

	dest.UpgradeLaneConfig.DeploySettings.DeployedAtBlock = destBlock
}
