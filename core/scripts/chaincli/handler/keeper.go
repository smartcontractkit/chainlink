package handler

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	keeper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	upkeep "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

// Keeper is the keepers commands handler
type Keeper struct {
	*baseHandler

	approveAmount  *big.Int
	addFundsAmount *big.Int
}

// NewKeeper is the constructor of Keeper
func NewKeeper(cfg *config.Config) *Keeper {
	approveAmount := big.NewInt(0)
	approveAmount.SetString(cfg.ApproveAmount, 10)

	addFundsAmount := big.NewInt(0)
	addFundsAmount.SetString(cfg.AddFundsAmount, 10)

	return &Keeper{
		baseHandler:    newBaseHandler(cfg),
		approveAmount:  approveAmount,
		addFundsAmount: addFundsAmount,
	}
}

// DeployKeepers contains a logic to deploy keepers.
func (k *Keeper) DeployKeepers(ctx context.Context) {
	var registry *keeper.KeeperRegistry
	var registryAddr common.Address
	if k.cfg.RegistryAddress != "" {
		// Get existing keeper registry
		registryAddr, registry = k.GetRegistry(ctx)
	} else {
		// Deploy keeper registry
		registryAddr, registry = k.deployRegistry(ctx)
	}

	// Approve keeper registry
	approveRegistryTx, err := k.linkToken.Approve(k.buildTxOpts(ctx), registryAddr, k.approveAmount)
	if err != nil {
		log.Fatal(registryAddr.Hex(), ": Approve failed - ", err)
	}
	waitTx(ctx, k.client, approveRegistryTx)
	log.Println(registryAddr.Hex(), ": KeeperRegistry approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveRegistryTx.Hash()))

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registry)

	// Set Keepers
	log.Println("Set keepers...")
	keepers, owners := k.keepers()
	setKeepersTx, err := registry.SetKeepers(k.buildTxOpts(ctx), keepers, owners)
	if err != nil {
		log.Fatal("SetKeepers failed: ", err)
	}
	waitTx(ctx, k.client, setKeepersTx)
	log.Println("Keepers registered:", setKeepersTx.Hash().Hex())
}

func (k *Keeper) deployRegistry(ctx context.Context) (common.Address, *keeper.KeeperRegistry) {
	registryAddr, deployKeeperRegistryTx, registryInstance, err := keeper.DeployKeeperRegistry(k.buildTxOpts(ctx), k.client,
		common.HexToAddress(k.cfg.LinkTokenAddr),
		common.HexToAddress(k.cfg.LinkETHFeedAddr),
		common.HexToAddress(k.cfg.FastGasFeedAddr),
		k.cfg.PaymentPremiumPBB,
		k.cfg.FlatFeeMicroLink,
		big.NewInt(k.cfg.BlockCountPerTurn),
		k.cfg.CheckGasLimit,
		big.NewInt(k.cfg.StalenessSeconds),
		k.cfg.GasCeilingMultiplier,
		big.NewInt(k.cfg.FallbackGasPrice),
		big.NewInt(k.cfg.FallbackLinkPrice),
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	waitDeployment(ctx, k.client, deployKeeperRegistryTx)
	log.Println("KeeperRegistry deployed:", registryAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryTx.Hash()))
	return registryAddr, registryInstance
}

func (k *Keeper) GetRegistry(ctx context.Context) (common.Address, *keeper.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	registryInstance, err := keeper.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	log.Println("KeeperRegistry at:", k.cfg.RegistryAddress)
	if k.cfg.RegistryConfigUpdate {
		transaction, err := registryInstance.SetConfig(k.buildTxOpts(ctx),
			k.cfg.PaymentPremiumPBB,
			k.cfg.FlatFeeMicroLink,
			big.NewInt(k.cfg.BlockCountPerTurn),
			k.cfg.CheckGasLimit,
			big.NewInt(k.cfg.StalenessSeconds),
			k.cfg.GasCeilingMultiplier,
			big.NewInt(k.cfg.FallbackGasPrice),
			big.NewInt(k.cfg.FallbackLinkPrice))
		if err != nil {
			log.Fatal("Registry config update: ", err)
		}
		waitTx(ctx, k.client, transaction)
		log.Println("KeeperRegistry config update:", k.cfg.RegistryAddress, "-", helpers.ExplorerLink(k.cfg.ChainID, transaction.Hash()))
	} else {
		log.Println("KeeperRegistry config not updated: KEEPER_CONFIG_UPDATE=false")
	}
	return registryAddr, registryInstance
}

func (k *Keeper) keepers() ([]common.Address, []common.Address) {
	var addrs []common.Address
	var fromAddrs []common.Address
	for _, addr := range k.cfg.Keepers {
		addrs = append(addrs, common.HexToAddress(addr))
		fromAddrs = append(fromAddrs, k.fromAddr)
	}
	return addrs, fromAddrs
}

// deployUpkeeps deploys N amount of upkeeps and register them in the keeper registry deployed above
func (k *Keeper) deployUpkeeps(ctx context.Context, registryInstance *keeper.KeeperRegistry) {
	fmt.Println()
	log.Println("Deploying upkeeps...")
	for i := int64(0); i < k.cfg.UpkeepCount; i++ {
		fmt.Println()
		// Deploy
		upkeepAddr, deployUpkeepTx, _, err := upkeep.DeployUpkeepPerformCounterRestrictive(k.buildTxOpts(ctx), k.client,
			big.NewInt(k.cfg.UpkeepTestRange), big.NewInt(k.cfg.UpkeepAverageEligibilityCadence),
		)
		if err != nil {
			log.Fatal(i, ": DeployAbi failed - ", err)
		}
		waitDeployment(ctx, k.client, deployUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep deployed - ", helpers.ExplorerLink(k.cfg.ChainID, deployUpkeepTx.Hash()))

		// Approve
		approveUpkeepTx, err := k.linkToken.Approve(k.buildTxOpts(ctx), upkeepAddr, k.approveAmount)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": Approve failed - ", err)
		}
		waitTx(ctx, k.client, approveUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveUpkeepTx.Hash()))

		// Register
		registerUpkeepTx, err := registryInstance.RegisterUpkeep(k.buildTxOpts(ctx),
			upkeepAddr, k.cfg.UpkeepGasLimit, k.fromAddr, []byte(k.cfg.UpkeepCheckData),
		)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
		}
		waitTx(ctx, k.client, registerUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep registered - ", helpers.ExplorerLink(k.cfg.ChainID, registerUpkeepTx.Hash()))

		// Fund
		addFundsTx, err := registryInstance.AddFunds(k.buildTxOpts(ctx), big.NewInt(int64(i)), k.addFundsAmount)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": AddFunds failed - ", err)
		}
		waitTx(ctx, k.client, addFundsTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep funded - ", helpers.ExplorerLink(k.cfg.ChainID, addFundsTx.Hash()))
	}
	fmt.Println()
}

func (k *Keeper) buildTxOpts(ctx context.Context) *bind.TransactOpts {
	nonce, err := k.client.PendingNonceAt(ctx, k.fromAddr)
	if err != nil {
		log.Fatal("PendingNonceAt failed: ", err)
	}

	gasPrice, err := k.client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("SuggestGasPrice failed: ", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(k.privateKey, big.NewInt(k.cfg.ChainID))
	if err != nil {
		log.Fatal("NewKeyedTransactorWithChainID failed: ", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = k.cfg.GasLimit // in units
	auth.GasPrice = gasPrice

	return auth
}
