package rhea

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/shared"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arm_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/router"
)

func DeployToNewChain(client *EvmDeploymentConfig) error {
	// Updates client.ARM if any new contracts are deployed
	err := deployARM(client)
	if err != nil {
		return errors.Wrap(err, "arm deployment failed")
	}
	// Updates client.TokenPools if any new contracts are deployed
	err = DeployTokenPools(client)
	if err != nil {
		return errors.Wrap(err, "pool deployment failed")
	}
	// Updates client.ChainConfig.Router if any new contracts are deployed
	err = deployRouter(client)
	if err != nil {
		return errors.Wrap(err, "router deployment failed")
	}
	// Updates client.ChainConfig.UpgradeRouter if any new contracts are deployed
	err = deployUpgradeRouter(client)
	if err != nil {
		return errors.Wrap(err, "upgrade router deployment failed")
	}
	// Update client.PriceRegistry if any new contracts are deployed
	err = deployPriceRegistry(client)
	if err != nil {
		return errors.Wrap(err, "price registry deployment failed")
	}
	return nil
}

func DeployUpgradeRouters(source *EvmDeploymentConfig, dest *EvmDeploymentConfig) error {
	err := deployUpgradeRouter(source)
	if err != nil {
		return errors.Wrap(err, "upgrade router in source chain deployment failed")
	}
	err = deployUpgradeRouter(dest)
	if err != nil {
		return errors.Wrap(err, "upgrade router in dest chain deployment failed")
	}
	return nil
}

func deployARM(client *EvmDeploymentConfig) error {
	if !client.ChainConfig.DeploySettings.DeployARM {
		if client.ChainConfig.ARM.Hex() == "0x0000000000000000000000000000000000000000" || client.ChainConfig.ARMProxy.Hex() == "0x0000000000000000000000000000000000000000" {
			return fmt.Errorf("deploy new arm set to false but no arm (proxy) given in config")
		}
		client.Logger.Infof("Skipping ARM deployment, using ARM on %s, proxy on %s", client.ChainConfig.ARM, client.ChainConfig.ARMProxy)
		return nil
	}

	client.Logger.Infof("Deploying ARM")
	var armAddress common.Address
	var tx *evmtypes.Transaction
	var err error
	armConfig := client.ChainConfig.ARMConfig
	switch armConfig {
	case nil:
		armAddress, tx, _, err = mock_arm_contract.DeployMockARMContract(client.Owner, client.Client)
	default:
		armAddress, tx, _, err = arm_contract.DeployARMContract(client.Owner, client.Client, *armConfig)
	}
	if err != nil {
		return err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	client.Logger.Infof("ARM deployed on %s in tx: %s", armAddress.Hex(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	client.ChainConfig.ARM = armAddress

	client.Logger.Infof("Deploying ARM proxy")
	proxyAddress, _, _, err := arm_proxy_contract.DeployARMProxyContract(client.Owner, client.Client, client.ChainConfig.ARM)
	if err != nil {
		return err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	client.Logger.Infof("ARM proxy deployed on %s in tx: %s", proxyAddress.Hex(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	client.ChainConfig.ARMProxy = proxyAddress

	return nil
}

func DeployTokenPools(client *EvmDeploymentConfig) error {
	for tokenName, tokenConfig := range client.ChainConfig.SupportedTokens {
		if err := deployPool(client, tokenName, tokenConfig); err != nil {
			return errors.Wrapf(err, "failed %s", tokenName)
		}
	}
	return nil
}

func deployPool(client *EvmDeploymentConfig, tokenName Token, tokenConfig EVMBridgedToken) error {
	if tokenConfig.TokenPoolType == FeeTokenOnly {
		client.Logger.Infof("Skipping pool deployment for fee only token")
		return nil
	}
	// Only deploy a new pool if there is no current pool address given
	// and the deploySetting indicate a new pool should be deployed.
	if client.ChainConfig.DeploySettings.DeployTokenPools && tokenConfig.Pool == common.HexToAddress("") {
		client.Logger.Infof("Deploying token pool for %s token", tokenName)
		var poolAddress, tokenAddress common.Address
		var err error
		switch tokenConfig.TokenPoolType {
		case LockRelease:
			poolAddress, err = deployLockReleaseTokenPool(client, tokenName, tokenConfig.Token, tokenConfig.PoolAllowList)
		case BurnMint:
			poolAddress, err = deployBurnMintTokenPool(client, tokenName, tokenConfig.Token, tokenConfig.PoolAllowList)
		case Wrapped:
			tokenAddress, poolAddress, err = deployWrappedTokenPool(client, tokenConfig.Token, tokenName, tokenConfig.PoolAllowList)
			// Since we also deployed the token we need to set it
			tokenConfig.Token = tokenAddress
		default:
			return fmt.Errorf("unknown pool type %s", tokenConfig.TokenPoolType)
		}
		if err != nil {
			return err
		}
		client.ChainConfig.SupportedTokens[tokenName] = EVMBridgedToken{
			Token:         tokenConfig.Token,
			Pool:          poolAddress,
			Price:         tokenConfig.Price,
			Decimals:      tokenConfig.Decimals,
			TokenPoolType: tokenConfig.TokenPoolType,
			PoolAllowList: tokenConfig.PoolAllowList,
		}
		return nil
	}

	// If no pools should be deployed but there is no pool address set fail.
	if tokenConfig.Pool == common.HexToAddress("") {
		return fmt.Errorf("deploy new %s pool set to false but no %s pool given in config", tokenName, tokenConfig.TokenPoolType)
	}
	client.Logger.Infof("Skipping %s Pool deployment, using Pool on %s", tokenName, tokenConfig.Pool)

	return nil
}

func deployLockReleaseTokenPool(client *EvmDeploymentConfig, tokenName Token, tokenAddress common.Address, poolAllowList []common.Address) (common.Address, error) {
	tokenPoolAddress, tx, _, err := lock_release_token_pool.DeployLockReleaseTokenPool(
		client.Owner,
		client.Client,
		tokenAddress,
		poolAllowList,
		client.ChainConfig.ARMProxy)
	if err != nil {
		return common.Address{}, err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return common.Address{}, err
	}
	client.Logger.Infof("Lock/release pool for %s deployed on %s in tx %s", tokenName, tokenPoolAddress, helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	pool, err := lock_release_token_pool.NewLockReleaseTokenPool(tokenPoolAddress, client.Client)
	if err != nil {
		return common.Address{}, err
	}
	err = fillPoolWithTokens(client, pool, tokenAddress, tokenName)
	return tokenPoolAddress, err
}

func deployBurnMintTokenPool(client *EvmDeploymentConfig, tokenName Token, tokenAddress common.Address, poolAllowList []common.Address) (common.Address, error) {
	client.Logger.Infof("Deploying token pool for %s token", tokenName)
	tokenPoolAddress, tx, _, err := burn_mint_token_pool.DeployBurnMintTokenPool(
		client.Owner,
		client.Client,
		tokenAddress,
		poolAllowList,
		client.ChainConfig.ARMProxy)
	if err != nil {
		return common.Address{}, err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return common.Address{}, err
	}
	client.Logger.Infof("Burn/mint pool for %s deployed on %s in tx %s", tokenName, tokenPoolAddress, helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	return tokenPoolAddress, nil
}

func deployWrappedTokenPool(client *EvmDeploymentConfig, tokenAddress common.Address, tokenName Token, poolAllowList []common.Address) (common.Address, common.Address, error) {
	client.Logger.Infof("Deploying token pool for %s token", tokenName)
	if tokenName.Symbol() == "" {
		return common.Address{}, common.Address{}, fmt.Errorf("no token symbol given for wrapped token pool %s", tokenName)
	}

	// Only deploy a new token if there is no current token address given
	if tokenAddress == common.HexToAddress("") {
		newTokenAddress, tx, _, err := burn_mint_erc677.DeployBurnMintERC677(client.Owner, client.Client, string(tokenName), tokenName.Symbol(), tokenName.Decimals(), big.NewInt(0))
		if err != nil {
			return common.Address{}, common.Address{}, err
		}
		tokenAddress = newTokenAddress
		if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
			return common.Address{}, common.Address{}, err
		}
		client.Logger.Infof("New %s token deployed on %s in tx %s", tokenName, tokenAddress, helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	} else {
		client.Logger.Infof("Using existing %s token deployed at", tokenName, tokenAddress)
	}

	poolAddress, err := deployBurnMintTokenPool(client, tokenName, tokenAddress, poolAllowList)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	token, err := burn_mint_erc677.NewBurnMintERC677(tokenAddress, client.Client)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	tx, err := token.GrantMintAndBurnRoles(client.Owner, poolAddress)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return common.Address{}, common.Address{}, err
	}

	return tokenAddress, poolAddress, nil
}

// deployRouter always uses an empty list of offRamps. Ramps should be set in the offRamp deployment step.
func deployRouter(client *EvmDeploymentConfig) error {
	if !client.ChainConfig.DeploySettings.DeployRouter {
		client.Logger.Infof("Skipping Router deployment, using Router on %s", client.ChainConfig.Router)
		return nil
	}

	client.Logger.Infof("Deploying Router")
	nativeFeeToken := common.Address{}
	if client.ChainConfig.WrappedNative != "" {
		nativeFeeToken = client.ChainConfig.SupportedTokens[client.ChainConfig.WrappedNative].Token
	}

	routerAddress, tx, _, err := router.DeployRouter(client.Owner, client.Client, nativeFeeToken, client.ChainConfig.ARMProxy)
	if err != nil {
		return err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	client.ChainConfig.Router = routerAddress

	client.Logger.Infof(fmt.Sprintf("Router deployed on %s in tx %s", routerAddress.String(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash())))
	return nil
}

// deployUpgradeRouter always uses an empty list of offRamps. Ramps should be set in the offRamp deployment step.
func deployUpgradeRouter(client *EvmDeploymentConfig) error {
	if !client.ChainConfig.DeploySettings.DeployUpgradeRouter {
		client.Logger.Infof("Skipping Upgrade Router deployment, using Router on %s", client.ChainConfig.UpgradeRouter)
		return nil
	}

	client.Logger.Infof("Deploying Router")
	nativeFeeToken := common.Address{}
	if client.ChainConfig.WrappedNative != "" {
		nativeFeeToken = client.ChainConfig.SupportedTokens[client.ChainConfig.WrappedNative].Token
	}

	routerAddress, tx, _, err := router.DeployRouter(client.Owner, client.Client, nativeFeeToken, client.ChainConfig.ARMProxy)
	if err != nil {
		return err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	client.ChainConfig.UpgradeRouter = routerAddress

	client.Logger.Infof(fmt.Sprintf("Router deployed on %s in tx %s", routerAddress.String(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash())))
	return nil
}

// deployPriceRegistry Prices is deployed without any feeUpdaters
func deployPriceRegistry(client *EvmDeploymentConfig) error {
	if !client.ChainConfig.DeploySettings.DeployPriceRegistry {
		client.Logger.Infof("Skipping PriceRegistry deployment, using PriceRegistry on %s", client.ChainConfig.PriceRegistry)
		return nil
	}

	feeTokens := make([]common.Address, len(client.ChainConfig.FeeTokens))
	for i, token := range client.ChainConfig.FeeTokens {
		feeTokens[i] = client.ChainConfig.SupportedTokens[token].Token
	}

	client.Logger.Infof("Deploying PriceRegistry")
	priceRegistry, tx, _, err := price_registry.DeployPriceRegistry(
		client.Owner,
		client.Client,
		[]common.Address{},
		feeTokens,
		60*60*24*14, // two weeks
	)
	if err != nil {
		return err
	}
	if err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true); err != nil {
		return err
	}
	client.ChainConfig.PriceRegistry = priceRegistry

	client.Logger.Infof(fmt.Sprintf("PriceRegistry deployed on %s in tx %s", priceRegistry.String(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash())))
	return nil
}
