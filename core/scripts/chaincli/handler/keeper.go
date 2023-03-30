package handler

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	registrylogic20 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic2_0"
	registry11 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry20 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_counter_wrapper"
	upkeep "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
)

// Keeper is the keepers commands handler
type Keeper struct {
	*baseHandler

	addFundsAmount *big.Int
}

// NewKeeper creates new instance of Keeper
func NewKeeper(cfg *config.Config) *Keeper {
	addFundsAmount := big.NewInt(0)
	addFundsAmount.SetString(cfg.AddFundsAmount, 10)

	return &Keeper{
		baseHandler:    NewBaseHandler(cfg),
		addFundsAmount: addFundsAmount,
	}
}

// DeployKeepers contains a logic to deploy keepers.
func (k *Keeper) DeployKeepers(ctx context.Context) {
	lggr, closeLggr := logger.NewLogger()
	logger.Sugared(lggr).ErrorIfFn(closeLggr, "Failed to close logger")

	keepers, owners := k.keepers()
	upkeepCount, registryAddr, deployer := k.prepareRegistry(ctx)

	// Create Keeper Jobs on Nodes for Registry
	cls := make([]cmd.HTTPClient, len(k.cfg.Keepers))
	for i, keeperAddr := range k.cfg.Keepers {
		url := k.cfg.KeeperURLs[i]
		email := k.cfg.KeeperEmails[i]
		if len(email) == 0 {
			email = defaultChainlinkNodeLogin
		}
		pwd := k.cfg.KeeperPasswords[i]
		if len(pwd) == 0 {
			pwd = defaultChainlinkNodePassword
		}

		cl, err := authenticate(url, email, pwd, lggr)
		if err != nil {
			log.Fatal(err)
		}
		cls[i] = cl

		if err = k.createKeeperJob(cl, k.cfg.RegistryAddress, keeperAddr); err != nil {
			log.Fatal(err)
		}
	}

	// Approve keeper registry
	k.approveFunds(ctx, registryAddr)

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registryAddr, deployer, upkeepCount)

	// Set Keepers on the registry
	k.setKeepers(ctx, cls, deployer, keepers, owners)
}

// DeployRegistry deploys a new keeper registry.
func (k *Keeper) DeployRegistry(ctx context.Context) {
	switch k.cfg.RegistryVersion {
	case keeper.RegistryVersion_1_1:
		k.deployRegistry11(ctx)
	case keeper.RegistryVersion_1_2:
		k.deployRegistry12(ctx)
	case keeper.RegistryVersion_2_0:
		k.deployRegistry20(ctx)
	default:
		panic("unsupported registry version")
	}
}

func (k *Keeper) prepareRegistry(ctx context.Context) (int64, common.Address, keepersDeployer) {
	var upkeepCount int64
	var registryAddr common.Address
	var deployer keepersDeployer
	var keeperRegistry11 *registry11.KeeperRegistry
	var keeperRegistry12 *registry12.KeeperRegistry
	var keeperRegistry20 *registry20.KeeperRegistry
	if k.cfg.RegistryAddress != "" {
		callOpts := bind.CallOpts{
			From:    k.fromAddr,
			Context: ctx,
		}

		// Get existing keeper registry
		switch k.cfg.RegistryVersion {
		case keeper.RegistryVersion_1_1:
			registryAddr, keeperRegistry11 = k.getRegistry11(ctx)
			count, err := keeperRegistry11.GetUpkeepCount(&callOpts)
			if err != nil {
				log.Fatal(registryAddr.Hex(), ": UpkeepCount failed - ", err)
			}
			upkeepCount = count.Int64()
			deployer = &v11KeeperDeployer{keeperRegistry11}
		case keeper.RegistryVersion_1_2:
			registryAddr, keeperRegistry12 = k.getRegistry12(ctx)
			state, err := keeperRegistry12.GetState(&callOpts)
			if err != nil {
				log.Fatal(registryAddr.Hex(), ": failed to getState - ", err)
			}
			upkeepCount = state.State.NumUpkeeps.Int64()
			deployer = &v12KeeperDeployer{keeperRegistry12}
		case keeper.RegistryVersion_2_0:
			registryAddr, keeperRegistry20 = k.getRegistry20(ctx)
			state, err := keeperRegistry20.GetState(&callOpts)
			if err != nil {
				log.Fatal(registryAddr.Hex(), ": failed to getState - ", err)
			}
			upkeepCount = state.State.NumUpkeeps.Int64()
			deployer = &v20KeeperDeployer{KeeperRegistryInterface: keeperRegistry20, cfg: k.cfg}
		default:
			panic(fmt.Errorf("version %s is not supported", k.cfg.RegistryVersion))
		}
	} else {
		// Deploy keeper registry
		switch k.cfg.RegistryVersion {
		case keeper.RegistryVersion_1_1:
			registryAddr, keeperRegistry11 = k.deployRegistry11(ctx)
			deployer = &v11KeeperDeployer{keeperRegistry11}
		case keeper.RegistryVersion_1_2:
			registryAddr, keeperRegistry12 = k.deployRegistry12(ctx)
			deployer = &v12KeeperDeployer{keeperRegistry12}
		case keeper.RegistryVersion_2_0:
			registryAddr, keeperRegistry20 = k.deployRegistry20(ctx)
			deployer = &v20KeeperDeployer{KeeperRegistryInterface: keeperRegistry20, cfg: k.cfg}
		default:
			panic(fmt.Errorf("version %s is not supported", k.cfg.RegistryVersion))
		}
	}

	return upkeepCount, registryAddr, deployer
}

func (k *Keeper) approveFunds(ctx context.Context, registryAddr common.Address) {
	if k.approveAmount.Cmp(big.NewInt(0)) == 0 {
		return
	}
	// Approve keeper registry
	approveRegistryTx, err := k.linkToken.Approve(k.buildTxOpts(ctx), registryAddr, k.approveAmount)
	if err != nil {
		log.Fatal(registryAddr.Hex(), ": Approve failed - ", err)
	}
	k.waitTx(ctx, approveRegistryTx)
	log.Println(registryAddr.Hex(), ": KeeperRegistry approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveRegistryTx.Hash()))
}

// deployRegistry20 deploys a version 2.0 keeper registry
func (k *Keeper) deployRegistry20(ctx context.Context) (common.Address, *registry20.KeeperRegistry) {
	registryLogicAddr, deployKeeperRegistryLogicTx, _, err := registrylogic20.DeployKeeperRegistryLogic(
		k.buildTxOpts(ctx),
		k.client,
		0,
		common.HexToAddress(k.cfg.LinkTokenAddr),
		common.HexToAddress(k.cfg.LinkETHFeedAddr),
		common.HexToAddress(k.cfg.FastGasFeedAddr),
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	k.waitDeployment(ctx, deployKeeperRegistryLogicTx)
	log.Println("KeeperRegistry2.0 Logic deployed:", registryLogicAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryLogicTx.Hash()))

	registryAddr, deployKeeperRegistryTx, registryInstance, err := registry20.DeployKeeperRegistry(
		k.buildTxOpts(ctx),
		k.client,
		registryLogicAddr,
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	k.waitDeployment(ctx, deployKeeperRegistryTx)
	log.Println("KeeperRegistry2.0 deployed:", registryAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryTx.Hash()))
	return registryAddr, registryInstance
}

// deployRegistry12 deploys a version 1.2 keeper registry
func (k *Keeper) deployRegistry12(ctx context.Context) (common.Address, *registry12.KeeperRegistry) {
	registryAddr, deployKeeperRegistryTx, registryInstance, err := registry12.DeployKeeperRegistry(
		k.buildTxOpts(ctx),
		k.client,
		common.HexToAddress(k.cfg.LinkTokenAddr),
		common.HexToAddress(k.cfg.LinkETHFeedAddr),
		common.HexToAddress(k.cfg.FastGasFeedAddr),
		*k.getConfigForRegistry12(),
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	k.waitDeployment(ctx, deployKeeperRegistryTx)
	log.Println("KeeperRegistry1.2 deployed:", registryAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryTx.Hash()))
	return registryAddr, registryInstance
}

// deployRegistry11 deploys a version 1.1 keeper registry
func (k *Keeper) deployRegistry11(ctx context.Context) (common.Address, *registry11.KeeperRegistry) {
	registryAddr, deployKeeperRegistryTx, registryInstance, err := registry11.DeployKeeperRegistry(k.buildTxOpts(ctx), k.client,
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
	k.waitDeployment(ctx, deployKeeperRegistryTx)
	log.Println("KeeperRegistry1.1 deployed:", registryAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryTx.Hash()))
	return registryAddr, registryInstance
}

// UpdateRegistry attaches to an existing registry and possibly updates registry config
func (k *Keeper) UpdateRegistry(ctx context.Context) {
	var registryAddr common.Address
	switch k.cfg.RegistryVersion {
	case keeper.RegistryVersion_1_1:
		registryAddr, _ = k.getRegistry11(ctx)
	case keeper.RegistryVersion_1_2:
		registryAddr, _ = k.getRegistry12(ctx)
	case keeper.RegistryVersion_2_0:
		registryAddr, _ = k.getRegistry20(ctx)
	default:
		panic("unexpected registry address")
	}
	log.Println("KeeperRegistry at:", registryAddr)
}

// getRegistry20 attaches to an existing 2.0 registry and possibly updates registry config
func (k *Keeper) getRegistry20(ctx context.Context) (common.Address, *registry20.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry20, err := registry20.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	if k.cfg.RegistryConfigUpdate {
		panic("KeeperRegistry2.0 could not be updated")
	} else {
		log.Println("KeeperRegistry2.0 config not updated: KEEPER_CONFIG_UPDATE=false")
	}
	return registryAddr, keeperRegistry20
}

// getRegistry12 attaches to an existing 1.2 registry and possibly updates registry config
func (k *Keeper) getRegistry12(ctx context.Context) (common.Address, *registry12.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry12, err := registry12.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	if k.cfg.RegistryConfigUpdate {
		transaction, err := keeperRegistry12.SetConfig(k.buildTxOpts(ctx), *k.getConfigForRegistry12())
		if err != nil {
			log.Fatal("Registry config update: ", err)
		}
		k.waitTx(ctx, transaction)
		log.Println("KeeperRegistry config update:", k.cfg.RegistryAddress, "-", helpers.ExplorerLink(k.cfg.ChainID, transaction.Hash()))
	} else {
		log.Println("KeeperRegistry config not updated: KEEPER_CONFIG_UPDATE=false")
	}
	return registryAddr, keeperRegistry12
}

// getRegistry11 attaches to an existing 1.1 registry and possibly updates registry config
func (k *Keeper) getRegistry11(ctx context.Context) (common.Address, *registry11.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry11, err := registry11.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	if k.cfg.RegistryConfigUpdate {
		transaction, err := keeperRegistry11.SetConfig(k.buildTxOpts(ctx),
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
		k.waitTx(ctx, transaction)
		log.Println("KeeperRegistry config update:", k.cfg.RegistryAddress, "-", helpers.ExplorerLink(k.cfg.ChainID, transaction.Hash()))
	} else {
		log.Println("KeeperRegistry config not updated: KEEPER_CONFIG_UPDATE=false")
	}
	return registryAddr, keeperRegistry11
}

// deployUpkeeps deploys upkeeps and funds upkeeps
func (k *Keeper) deployUpkeeps(ctx context.Context, registryAddr common.Address, deployer upkeepDeployer, existingCount int64) {
	fmt.Println()
	log.Println("Deploying upkeeps...")
	var upkeepAddrs []common.Address
	for i := existingCount; i < k.cfg.UpkeepCount+existingCount; i++ {
		fmt.Println()
		// Deploy
		var upkeepAddr common.Address
		var deployUpkeepTx *types.Transaction
		var err error
		if k.cfg.UpkeepAverageEligibilityCadence > 0 {
			upkeepAddr, deployUpkeepTx, _, err = upkeep.DeployUpkeepPerformCounterRestrictive(k.buildTxOpts(ctx), k.client,
				big.NewInt(k.cfg.UpkeepTestRange), big.NewInt(k.cfg.UpkeepAverageEligibilityCadence),
			)
		} else {
			upkeepAddr, deployUpkeepTx, _, err = upkeep_counter_wrapper.DeployUpkeepCounter(k.buildTxOpts(ctx), k.client,
				big.NewInt(k.cfg.UpkeepTestRange), big.NewInt(k.cfg.UpkeepInterval),
			)
		}
		if err != nil {
			log.Fatal(i, ": Deploy Upkeep failed - ", err)
		}
		k.waitDeployment(ctx, deployUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep deployed - ", helpers.ExplorerLink(k.cfg.ChainID, deployUpkeepTx.Hash()))

		// Register
		registerUpkeepTx, err := deployer.RegisterUpkeep(k.buildTxOpts(ctx),
			upkeepAddr, k.cfg.UpkeepGasLimit, k.fromAddr, []byte(k.cfg.UpkeepCheckData),
		)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
		}
		k.waitTx(ctx, registerUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep registered - ", helpers.ExplorerLink(k.cfg.ChainID, registerUpkeepTx.Hash()))

		upkeepAddrs = append(upkeepAddrs, upkeepAddr)
	}

	var err error
	var upkeepGetter activeUpkeepGetter
	switch k.cfg.RegistryVersion {
	case keeper.RegistryVersion_1_1:
		panic("not supported 1.1 registry")
	case keeper.RegistryVersion_1_2:
		upkeepGetter, err = registry12.NewKeeperRegistry(
			registryAddr,
			k.client,
		)
	case keeper.RegistryVersion_2_0:
		upkeepGetter, err = registry20.NewKeeperRegistry(
			registryAddr,
			k.client,
		)
	default:
		panic("unexpected registry address")
	}
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}

	activeUpkeepIds := k.getActiveUpkeepIds(ctx, upkeepGetter, big.NewInt(existingCount), big.NewInt(k.cfg.UpkeepCount))

	for index, upkeepAddr := range upkeepAddrs {
		// Approve
		k.approveFunds(ctx, registryAddr)

		upkeepId := activeUpkeepIds[index]

		// Fund
		addFundsTx, err := deployer.AddFunds(k.buildTxOpts(ctx), upkeepId, k.addFundsAmount)
		if err != nil {
			log.Fatal(upkeepId, upkeepAddr.Hex(), ": AddFunds failed - ", err)
		}
		k.waitTx(ctx, addFundsTx)
		log.Println(upkeepId, upkeepAddr.Hex(), ": Upkeep funded - ", helpers.ExplorerLink(k.cfg.ChainID, addFundsTx.Hash()))
	}
	fmt.Println()
}

// setKeepers set the keeper list for a registry
func (k *Keeper) setKeepers(ctx context.Context, cls []cmd.HTTPClient, deployer keepersDeployer, keepers, owners []common.Address) {
	if len(keepers) > 0 {
		log.Println("Set keepers...")
		setKeepersTx, err := deployer.SetKeepers(k.buildTxOpts(ctx), cls, keepers, owners)
		if err != nil {
			log.Fatal("SetKeepers failed: ", err)
		}
		k.waitTx(ctx, setKeepersTx)
		log.Println("Keepers registered:", helpers.ExplorerLink(k.cfg.ChainID, setKeepersTx.Hash()))
	} else {
		log.Println("No Keepers to register")
	}
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

// createKeeperJobOnExistingNode connect to existing node to create keeper job
func (k *Keeper) createKeeperJobOnExistingNode(urlStr, email, password, registryAddr, nodeAddr string) error {
	lggr, closeLggr := logger.NewLogger()
	logger.Sugared(lggr).ErrorIfFn(closeLggr, "Failed to close logger")

	cl, err := authenticate(urlStr, email, password, lggr)
	if err != nil {
		return err
	}

	if err = k.createKeeperJob(cl, registryAddr, nodeAddr); err != nil {
		log.Println("Failed to create keeper job: ", err)
		return err
	}

	return nil
}

type activeUpkeepGetter interface {
	Address() common.Address
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
}

// getActiveUpkeepIds retrieves active upkeep ids from registry
func (k *Keeper) getActiveUpkeepIds(ctx context.Context, registry activeUpkeepGetter, from, to *big.Int) []*big.Int {
	activeUpkeepIds, _ := registry.GetActiveUpkeepIDs(&bind.CallOpts{
		Pending: false,
		From:    k.fromAddr,
		Context: ctx,
	}, from, to)
	return activeUpkeepIds
}

// getConfigForRegistry12 returns a config object for registry 1.2
func (k *Keeper) getConfigForRegistry12() *registry12.Config {
	return &registry12.Config{
		PaymentPremiumPPB:    k.cfg.PaymentPremiumPBB,
		FlatFeeMicroLink:     k.cfg.FlatFeeMicroLink,
		BlockCountPerTurn:    big.NewInt(k.cfg.BlockCountPerTurn),
		CheckGasLimit:        k.cfg.CheckGasLimit,
		StalenessSeconds:     big.NewInt(k.cfg.StalenessSeconds),
		GasCeilingMultiplier: k.cfg.GasCeilingMultiplier,
		MinUpkeepSpend:       big.NewInt(k.cfg.MinUpkeepSpend),
		MaxPerformGas:        k.cfg.MaxPerformGas,
		FallbackGasPrice:     big.NewInt(k.cfg.FallbackGasPrice),
		FallbackLinkPrice:    big.NewInt(k.cfg.FallbackLinkPrice),
		Transcoder:           common.HexToAddress(k.cfg.Transcoder),
		Registrar:            common.HexToAddress(k.cfg.Registrar),
	}
}
