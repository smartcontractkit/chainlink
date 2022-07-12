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
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/sessions"
)

// Keeper is the keepers commands handler
type Keeper struct {
	*baseHandler

	addFundsAmount *big.Int
}

// canceller describes the behavior to cancel upkeeps
type canceller interface {
	CancelUpkeep(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error)
	WithdrawFunds(opts *bind.TransactOpts, id *big.Int, to common.Address) (*types.Transaction, error)
	RecoverFunds(opts *bind.TransactOpts) (*types.Transaction, error)
}

// upkeepDeployer contains functions needed to deploy an upkeep
type upkeepDeployer interface {
	RegisterUpkeep(opts *bind.TransactOpts, target common.Address, gasLimit uint32, admin common.Address, checkData []byte) (*types.Transaction, error)
	AddFunds(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error)
}

// keepersDeployer contains functions needed to deploy keepers
type keepersDeployer interface {
	canceller
	upkeepDeployer
	SetKeepers(opts *bind.TransactOpts, keepers []common.Address, payees []common.Address) (*types.Transaction, error)
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
	keepers, owners := k.keepers()
	upkeepCount, registryAddr, deployer := k.prepareRegistry(ctx)

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
	k.approveFunds(ctx, registryAddr)

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registryAddr, deployer, upkeepCount)

	// Set Keepers on the registry
	k.setKeepers(ctx, deployer, keepers, owners)
}

func (k *Keeper) prepareRegistry(ctx context.Context) (int64, common.Address, keepersDeployer) {
	var upkeepCount int64
	var registryAddr common.Address
	var deployer keepersDeployer
	var keeperRegistry11 *registry11.KeeperRegistry
	var keeperRegistry12 *registry12.KeeperRegistry
	isVersion12 := k.cfg.RegistryVersion == keeper.RegistryVersion_1_2
	if k.cfg.RegistryAddress != "" {
		callOpts := bind.CallOpts{
			Pending: false,
			From:    k.fromAddr,
			Context: ctx,
		}
		// Get existing keeper registry
		if isVersion12 {
			registryAddr, keeperRegistry12 = k.getRegistry2(ctx)
			state, err := keeperRegistry12.GetState(&callOpts)
			if err != nil {
				log.Fatal(registryAddr.Hex(), ": failed to getState - ", err)
			}
			upkeepCount = state.State.NumUpkeeps.Int64()
			deployer = keeperRegistry12
		} else {
			registryAddr, keeperRegistry11 = k.getRegistry1(ctx)
			count, err := keeperRegistry11.GetUpkeepCount(&callOpts)
			if err != nil {
				log.Fatal(registryAddr.Hex(), ": UpkeepCount failed - ", err)
			}
			upkeepCount = count.Int64()
			deployer = keeperRegistry11
		}
	} else {
		// Deploy keeper registry
		upkeepCount = 0
		if isVersion12 {
			registryAddr, deployer = k.deployRegistry2(ctx)
		} else {
			registryAddr, deployer = k.deployRegistry1(ctx)
		}
	}

	log.Println("Upkeep Count: ", upkeepCount)
	return upkeepCount, registryAddr, deployer
}

func (k *Keeper) approveFunds(ctx context.Context, registryAddr common.Address) {
	// Approve keeper registry
	approveRegistryTx, err := k.linkToken.Approve(k.buildTxOpts(ctx), registryAddr, k.approveAmount)
	if err != nil {
		log.Fatal(registryAddr.Hex(), ": Approve failed - ", err)
	}
	k.waitTx(ctx, approveRegistryTx)
	log.Println(registryAddr.Hex(), ": KeeperRegistry approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveRegistryTx.Hash()))
}

// deployRegistry2 deploys a version 1.2 keeper registry
func (k *Keeper) deployRegistry2(ctx context.Context) (common.Address, *registry12.KeeperRegistry) {
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
	log.Println("KeeperRegistry deployed:", registryAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryTx.Hash()))
	return registryAddr, registryInstance
}

// deployRegistry1 deploys a version 1.1 keeper registry
func (k *Keeper) deployRegistry1(ctx context.Context) (common.Address, *registry11.KeeperRegistry) {
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

// GetRegistry attaches to an existing registry and possibly updates registry config
func (k *Keeper) GetRegistry(ctx context.Context) {
	isVersion12 := k.cfg.RegistryVersion == keeper.RegistryVersion_1_2
	if isVersion12 {
		k.getRegistry2(ctx)
	} else {
		k.getRegistry1(ctx)
	}
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

// deployUpkeeps deploys upkeeps and funds upkeeps
func (k *Keeper) deployUpkeeps(ctx context.Context, registryAddr common.Address, deployer upkeepDeployer, existingCount int64) {
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
		k.approveFunds(ctx, registryAddr)

		// Register
		registerUpkeepTx, err := deployer.RegisterUpkeep(k.buildTxOpts(ctx),
			upkeepAddr, k.cfg.UpkeepGasLimit, k.fromAddr, []byte(k.cfg.UpkeepCheckData),
		)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
		}
		k.waitTx(ctx, registerUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep registered - ", helpers.ExplorerLink(k.cfg.ChainID, registerUpkeepTx.Hash()))

		// Fund
		addFundsTx, err := deployer.AddFunds(k.buildTxOpts(ctx), big.NewInt(i), k.addFundsAmount)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": AddFunds failed - ", err)
		}
		k.waitTx(ctx, addFundsTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep funded - ", helpers.ExplorerLink(k.cfg.ChainID, addFundsTx.Hash()))
	}
	fmt.Println()
}

// setKeepers set the keeper list for a registry
func (k *Keeper) setKeepers(ctx context.Context, deployer keepersDeployer, keepers, owners []common.Address) {
	if len(keepers) > 0 {
		log.Println("Set keepers...")
		setKeepersTx, err := deployer.SetKeepers(k.buildTxOpts(ctx), keepers, owners)
		if err != nil {
			log.Fatal("SetKeepers failed: ", err)
		}
		k.waitTx(ctx, setKeepersTx)
		log.Println("Keepers registered:", setKeepersTx.Hash().Hex())
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
	lggr, close := logger.NewLogger()
	defer close()

	cl, err := k.authenticate(urlStr, email, password, lggr)
	if err != nil {
		return err
	}

	if err = k.createKeeperJob(cl, registryAddr, nodeAddr); err != nil {
		log.Println("Failed to create keeper job: ", err)
		return err
	}
	return nil
}

// authenticate creates a http client with URL, email and password
func (k *Keeper) authenticate(urlStr, email, password string, lggr logger.Logger) (cmd.HTTPClient, error) {
	remoteNodeURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	c := cmd.ClientOpts{RemoteNodeURL: *remoteNodeURL}
	sr := sessions.SessionRequest{Email: email, Password: password}
	store := &cmd.MemoryCookieStore{}

	tca := cmd.NewSessionCookieAuthenticator(c, store, lggr)
	if _, err = tca.Authenticate(sr); err != nil {
		log.Println("failed to authenticate: ", err)
		return nil, err
	}

	return cmd.NewAuthenticatedHTTPClient(lggr, c, tca, sr), nil
}

// getActiveUpkeepIds retrieves active upkeep ids from registry 1.2
func (k *Keeper) getActiveUpkeepIds(ctx context.Context, registry *registry12.KeeperRegistry) []*big.Int {
	activeUpkeepIds, err := registry.GetActiveUpkeepIDs(&bind.CallOpts{
		Pending: false,
		From:    k.fromAddr,
		Context: ctx,
	}, big.NewInt(0), big.NewInt(0))
	if err != nil {
		log.Fatal(registry.Address().Hex(), ": failed to get active upkeep Ids - ", err)
	}
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
