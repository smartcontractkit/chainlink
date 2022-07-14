package handler

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/cmd"
	registry11 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/upkeep_counter_wrapper"
	upkeep "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/sessions"
)

// Keeper is the keepers commands handler
type Keeper struct {
	*baseHandler

	addFundsAmount *big.Int
}

// NewKeeper is the constructor of Keeper
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
	keepers, owners := k.keepers()
	k.deployKeepers(ctx, keepers, owners)
}

func (k *Keeper) deployKeepers(ctx context.Context, keepers []common.Address, owners []common.Address) common.Address {
	var registry *registry11.KeeperRegistry
	var registryAddr common.Address
	var upkeepCount int64
	if k.cfg.RegistryAddress != "" {
		// Get existing keeper registry
		registryAddr, registry = k.GetRegistry(ctx)
		callOpts := bind.CallOpts{
			Pending: false,
			From:    k.fromAddr,
			Context: ctx,
		}
		count, err := registry.GetUpkeepCount(&callOpts)
		if err != nil {
			log.Fatal(registryAddr.Hex(), ": UpkeepCount failed - ", err)
		}
		upkeepCount = count.Int64()
	} else {
		// Deploy keeper registry
		registryAddr, registry = k.deployRegistry(ctx)
		upkeepCount = 0
	}
	log.Println("Upkeep Count: ", upkeepCount)

	// Create Keeper Jobs on Nodes for Registry
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
		err := k.createKeeperJobOnExistingNode(url, email, pwd, registryAddr.Hex(), keeperAddr)
		if err != nil {
			log.Printf("Keeper Job not created for keeper %d: %s %s\n", i, url, keeperAddr)
			log.Println("Please create it manually")
		}
	}

	// Approve keeper registry
	approveRegistryTx, err := k.linkToken.Approve(k.buildTxOpts(ctx), registryAddr, k.approveAmount)
	if err != nil {
		log.Fatal(registryAddr.Hex(), ": Approve failed - ", err)
	}
	k.waitTx(ctx, approveRegistryTx)
	log.Println(registryAddr.Hex(), ": KeeperRegistry approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveRegistryTx.Hash()))

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registryAddr, registry, upkeepCount)

	// Set Keepers
	if len(keepers) > 0 {
		log.Println("Set keepers...")
		setKeepersTx, err := registry.SetKeepers(k.buildTxOpts(ctx), keepers, owners)
		if err != nil {
			log.Fatal("SetKeepers failed: ", err)
		}
		k.waitTx(ctx, setKeepersTx)
		log.Println("Keepers registered:", setKeepersTx.Hash().Hex())
	} else {
		log.Println("No Keepers to register")
	}

	return registryAddr
}

func (k *Keeper) deployRegistry(ctx context.Context) (common.Address, *registry11.KeeperRegistry) {
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
	log.Println("KeeperRegistry deployed:", registryAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryTx.Hash()))
	return registryAddr, registryInstance
}

// GetRegistry is used to attach to an existing registry
func (k *Keeper) GetRegistry(ctx context.Context) (common.Address, *registry11.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	registryInstance, err := registry11.NewKeeperRegistry(
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
		k.waitTx(ctx, transaction)
		log.Println("KeeperRegistry config update:", k.cfg.RegistryAddress, "-", helpers.ExplorerLink(k.cfg.ChainID, transaction.Hash()))
	} else {
		log.Println("KeeperRegistry config not updated: KEEPER_CONFIG_UPDATE=false")
	}
	return registryAddr, registryInstance
}

// getRegistry2 attaches to an existing 1.2 registry and possibly updates registry config
func (k *Keeper) getRegistry2(ctx context.Context) (common.Address, *registry12.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry12, err := registry12.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	log.Println("KeeperRegistry at:", k.cfg.RegistryAddress)
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

// getRegistry1 attaches to an existing 1.1 registry and possibly updates registry config
func (k *Keeper) getRegistry1(ctx context.Context) (common.Address, *registry11.KeeperRegistry) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry11, err := registry11.NewKeeperRegistry(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	log.Println("KeeperRegistry at:", k.cfg.RegistryAddress)
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

// deployUpkeeps deploys N amount of upkeeps and register them in the keeper registry deployed above
func (k *Keeper) deployUpkeeps(ctx context.Context, registryAddr common.Address, registryInstance *registry11.KeeperRegistry, existingCount int64) {
	fmt.Println()
	log.Println("Deploying upkeeps...")
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

		// Approve
		approveUpkeepTx, err := k.linkToken.Approve(k.buildTxOpts(ctx), registryAddr, k.approveAmount)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": Approve failed - ", err)
		}
		k.waitTx(ctx, approveUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveUpkeepTx.Hash()))

		// Register
		registerUpkeepTx, err := registryInstance.RegisterUpkeep(k.buildTxOpts(ctx),
			upkeepAddr, k.cfg.UpkeepGasLimit, k.fromAddr, []byte(k.cfg.UpkeepCheckData),
		)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
		}
		k.waitTx(ctx, registerUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep registered - ", helpers.ExplorerLink(k.cfg.ChainID, registerUpkeepTx.Hash()))

		// Fund
		addFundsTx, err := registryInstance.AddFunds(k.buildTxOpts(ctx), big.NewInt(int64(i)), k.addFundsAmount)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": AddFunds failed - ", err)
		}
		k.waitTx(ctx, addFundsTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep funded - ", helpers.ExplorerLink(k.cfg.ChainID, addFundsTx.Hash()))
	}
	fmt.Println()
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
	remoteNodeURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	c := cmd.ClientOpts{RemoteNodeURL: *remoteNodeURL}
	sr := sessions.SessionRequest{Email: email, Password: password}
	store := &cmd.MemoryCookieStore{}
	lggr, close := logger.NewLogger()
	defer close()
	tca := cmd.NewSessionCookieAuthenticator(c, store, lggr)
	if _, err := tca.Authenticate(sr); err != nil {
		log.Println("failed to authenticate: ", err)
		return err
	}
	cl := cmd.NewAuthenticatedHTTPClient(lggr, c, tca, sr)

	if err := k.createKeeperJob(cl, registryAddr, nodeAddr); err != nil {
		log.Println("Failed to create keeper job: ", err)
		return err
	}
	return nil
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
