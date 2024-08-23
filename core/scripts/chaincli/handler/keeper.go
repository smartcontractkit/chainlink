package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/streams"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	automationForwarderLogic "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_forwarder_logic"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	registrylogic20 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic2_0"
	registrylogica21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_a_wrapper_2_1"
	registrylogicb21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
	registry11 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry12 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry20 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	registry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/streams_lookup_upkeep_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_counter_wrapper"
	upkeep "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/verifiable_load_streams_lookup_upkeep_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/verifiable_load_upkeep_wrapper"
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

		cl, err := authenticate(ctx, url, email, pwd, lggr)
		if err != nil {
			log.Fatal(err)
		}
		cls[i] = cl

		if err = k.createKeeperJob(ctx, cl, k.cfg.RegistryAddress, keeperAddr); err != nil {
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
func (k *Keeper) DeployRegistry(ctx context.Context, verify bool) {
	if verify {
		if k.cfg.RegistryVersion != keeper.RegistryVersion_2_1 && k.cfg.RegistryVersion != keeper.RegistryVersion_2_0 {
			log.Fatal("keeper registry verification is only supported for version 2.0 and 2.1")
		}
		if k.cfg.ExplorerAPIKey == "" || k.cfg.ExplorerAPIKey == "<explorer-api-key>" || k.cfg.NetworkName == "" || k.cfg.NetworkName == "<network-name>" {
			log.Fatal("please set your explore API key and network name in the .env file to verify the registry contract")
		}

		// Get the current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("failed to get current working directory: %w", err)
		}

		// Check if it is the root directory of chaincli
		if !strings.HasSuffix(currentDir, "core/scripts/chaincli") {
			log.Fatal("please run the command from the root directory of chaincli to verify the registry")
		}
	}

	switch k.cfg.RegistryVersion {
	case keeper.RegistryVersion_1_1:
		k.deployRegistry11(ctx)
	case keeper.RegistryVersion_1_2:
		k.deployRegistry12(ctx)
	case keeper.RegistryVersion_2_0:
		k.deployRegistry20(ctx, verify)
	case keeper.RegistryVersion_2_1:
		k.deployRegistry21(ctx, verify)
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
	var keeperRegistry21 *iregistry21.IKeeperRegistryMaster
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
		case keeper.RegistryVersion_2_1:
			registryAddr, keeperRegistry21 = k.getRegistry21(ctx)
			state, err := keeperRegistry21.GetState(&callOpts)
			if err != nil {
				log.Fatal(registryAddr.Hex(), ": failed to getState - ", err)
			}
			upkeepCount = state.State.NumUpkeeps.Int64()
			deployer = &v21KeeperDeployer{IKeeperRegistryMasterInterface: keeperRegistry21, cfg: k.cfg}
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
			registryAddr, keeperRegistry20 = k.deployRegistry20(ctx, true)
			deployer = &v20KeeperDeployer{KeeperRegistryInterface: keeperRegistry20, cfg: k.cfg}
		case keeper.RegistryVersion_2_1:
			registryAddr, keeperRegistry21 = k.deployRegistry21(ctx, false)
			deployer = &v21KeeperDeployer{IKeeperRegistryMasterInterface: keeperRegistry21, cfg: k.cfg}
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

	if err := k.waitTx(ctx, approveRegistryTx); err != nil {
		log.Fatalf("KeeperRegistry ApproveFunds failed for registryAddr: %s, and approveAmount: %s, error is: %s", k.cfg.RegistryAddress, k.approveAmount, err.Error())
	}

	log.Println(registryAddr.Hex(), ": KeeperRegistry approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveRegistryTx.Hash()))
}

func (k *Keeper) VerifyContract(params ...string) {
	// Change to the contracts directory where the hardhat.config.ts file is located
	if err := k.changeToContractsDirectory(); err != nil {
		log.Fatalf("failed to change to directory where the hardhat.config.ts file is located: %v", err)
	}

	// Append the address and params to the commandArgs slice
	commandArgs := append([]string{}, params...)

	// Format the command string with the commandArgs
	command := fmt.Sprintf(
		"NODE_HTTP_URL='%s' EXPLORER_API_KEY='%s' NETWORK_NAME='%s' pnpm hardhat verify --network env %s",
		k.cfg.NodeHttpURL,
		k.cfg.ExplorerAPIKey,
		k.cfg.NetworkName,
		strings.Join(commandArgs, " "),
	)

	fmt.Println("Running command to verify contract: ", command)
	if err := k.runCommand(command); err != nil {
		log.Println("Contract verification on Explorer failed: ", err)
	}
}

// deployRegistry21 deploys a version 2.1 keeper registry
func (k *Keeper) deployRegistry21(ctx context.Context, verify bool) (common.Address, *iregistry21.IKeeperRegistryMaster) {
	automationForwarderLogicAddr, tx, _, err := automationForwarderLogic.DeployAutomationForwarderLogic(k.buildTxOpts(ctx), k.client)
	if err != nil {
		log.Fatal("Deploy AutomationForwarderLogic failed: ", err)
	}
	k.waitDeployment(ctx, tx)
	log.Println("AutomationForwarderLogic deployed:", automationForwarderLogicAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, tx.Hash()))

	registryLogicBAddr, tx, _, err := registrylogicb21.DeployKeeperRegistryLogicB(
		k.buildTxOpts(ctx),
		k.client,
		k.cfg.Mode,
		common.HexToAddress(k.cfg.LinkTokenAddr),
		common.HexToAddress(k.cfg.LinkETHFeedAddr),
		common.HexToAddress(k.cfg.FastGasFeedAddr),
		automationForwarderLogicAddr,
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	k.waitDeployment(ctx, tx)
	log.Println("KeeperRegistry LogicB 2.1 deployed:", registryLogicBAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, tx.Hash()))

	// verify KeeperRegistryLogicB
	if verify {
		k.VerifyContract(registryLogicBAddr.String(), "0", k.cfg.LinkTokenAddr, k.cfg.LinkETHFeedAddr, k.cfg.FastGasFeedAddr)
		log.Println("KeeperRegistry LogicB 2.1 verified successfully")
	}

	registryLogicAAddr, tx, _, err := registrylogica21.DeployKeeperRegistryLogicA(
		k.buildTxOpts(ctx),
		k.client,
		registryLogicBAddr,
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	k.waitDeployment(ctx, tx)
	log.Println("KeeperRegistry LogicA 2.1 deployed:", registryLogicAAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, tx.Hash()))

	// verify KeeperRegistryLogicA
	if verify {
		k.VerifyContract(registryLogicAAddr.String(), registryLogicBAddr.String())
		log.Println("KeeperRegistry LogicA 2.1 verified successfully")
	}

	registryAddr, deployKeeperRegistryTx, _, err := registry21.DeployKeeperRegistry(
		k.buildTxOpts(ctx),
		k.client,
		registryLogicAAddr,
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	k.waitDeployment(ctx, deployKeeperRegistryTx)
	log.Println("KeeperRegistry 2.1 deployed:", registryAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryTx.Hash()))

	registryInstance, err := iregistry21.NewIKeeperRegistryMaster(registryAddr, k.client)
	if err != nil {
		log.Fatal("Failed to attach to deployed contract: ", err)
	}

	// verify KeeperRegistry
	if verify {
		k.VerifyContract(registryAddr.String(), registryLogicAAddr.String())
		log.Println("KeeperRegistry 2.1 verified successfully")
	}

	return registryAddr, registryInstance
}

// deployRegistry20 deploys a version 2.0 keeper registry
func (k *Keeper) deployRegistry20(ctx context.Context, verify bool) (common.Address, *registry20.KeeperRegistry) {
	registryLogicAddr, deployKeeperRegistryLogicTx, _, err := registrylogic20.DeployKeeperRegistryLogic(
		k.buildTxOpts(ctx),
		k.client,
		k.cfg.Mode,
		common.HexToAddress(k.cfg.LinkTokenAddr),
		common.HexToAddress(k.cfg.LinkETHFeedAddr),
		common.HexToAddress(k.cfg.FastGasFeedAddr),
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	k.waitDeployment(ctx, deployKeeperRegistryLogicTx)
	log.Println("KeeperRegistry2.0 Logic deployed:", registryLogicAddr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, deployKeeperRegistryLogicTx.Hash()))

	// verify KeeperRegistryLogic
	if verify {
		k.VerifyContract(registryLogicAddr.String(), "0", k.cfg.LinkTokenAddr, k.cfg.LinkETHFeedAddr, k.cfg.FastGasFeedAddr)
		log.Println("KeeperRegistry Logic 2.0 verified successfully")
	}

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

	// verify KeeperRegistry
	if verify {
		k.VerifyContract(registryAddr.String(), registryLogicAddr.String())
		log.Println("KeeperRegistry 2.0 verified successfully")
	}

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
	case keeper.RegistryVersion_2_1:
		registryAddr, _ = k.getRegistry21(ctx)
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
	}
	log.Println("KeeperRegistry2.0 config not updated: KEEPER_CONFIG_UPDATE=false")
	return registryAddr, keeperRegistry20
}

// getRegistry21 attaches to an existing 2.1 registry and possibly updates registry config
func (k *Keeper) getRegistry21(ctx context.Context) (common.Address, *iregistry21.IKeeperRegistryMaster) {
	registryAddr := common.HexToAddress(k.cfg.RegistryAddress)
	keeperRegistry21, err := iregistry21.NewIKeeperRegistryMaster(
		registryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal("Registry failed: ", err)
	}
	if k.cfg.RegistryConfigUpdate {
		panic("KeeperRegistry2.1 could not be updated")
	}
	log.Println("KeeperRegistry2.1 config not updated: KEEPER_CONFIG_UPDATE=false")
	return registryAddr, keeperRegistry21
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

		if err := k.waitTx(ctx, transaction); err != nil {
			log.Fatalf("KeeperRegistry config update failed on registry address: %s, error is: %s", k.cfg.RegistryAddress, err.Error())
		}
		log.Println("KeeperRegistry config update:", k.cfg.RegistryAddress, "-", helpers.ExplorerLink(k.cfg.ChainID, transaction.Hash()))
	}
	log.Println("KeeperRegistry config not updated: KEEPER_CONFIG_UPDATE=false")
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

		if err := k.waitTx(ctx, transaction); err != nil {
			log.Fatalf("KeeperRegistry config update failed on registry address: %s, error is %s", k.cfg.RegistryAddress, err.Error())
		}
		log.Println("KeeperRegistry config update:", k.cfg.RegistryAddress, "-", helpers.ExplorerLink(k.cfg.ChainID, transaction.Hash()))
	}
	log.Println("KeeperRegistry config not updated: KEEPER_CONFIG_UPDATE=false")
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
		var registerUpkeepTx *types.Transaction
		var logUpkeepCounter *log_upkeep_counter_wrapper.LogUpkeepCounter
		var checkData []byte

		switch k.cfg.UpkeepType {
		case config.Conditional:
			checkData = []byte(k.cfg.UpkeepCheckData)
			var err error
			if k.cfg.UpkeepAverageEligibilityCadence > 0 {
				upkeepAddr, deployUpkeepTx, _, err = upkeep.DeployUpkeepPerformCounterRestrictive(
					k.buildTxOpts(ctx),
					k.client,
					big.NewInt(k.cfg.UpkeepTestRange),
					big.NewInt(k.cfg.UpkeepAverageEligibilityCadence),
				)
			} else if k.cfg.VerifiableLoadTest {
				upkeepAddr, deployUpkeepTx, _, err = verifiable_load_upkeep_wrapper.DeployVerifiableLoadUpkeep(
					k.buildTxOpts(ctx),
					k.client,
					common.HexToAddress(k.cfg.Registrar),
					k.cfg.UseArbBlockNumber,
				)
			} else {
				upkeepAddr, deployUpkeepTx, _, err = upkeep_counter_wrapper.DeployUpkeepCounter(
					k.buildTxOpts(ctx),
					k.client,
					big.NewInt(k.cfg.UpkeepTestRange),
					big.NewInt(k.cfg.UpkeepInterval),
				)
			}
			if err != nil {
				log.Fatal(i, ": Deploy Upkeep failed - ", err)
			}
			k.waitDeployment(ctx, deployUpkeepTx)
			log.Println(i, upkeepAddr.Hex(), ": Upkeep deployed - ", helpers.ExplorerLink(k.cfg.ChainID, deployUpkeepTx.Hash()))
			registerUpkeepTx, err = deployer.RegisterUpkeep(k.buildTxOpts(ctx),
				upkeepAddr, k.cfg.UpkeepGasLimit, k.fromAddr, checkData, []byte{},
			)
			if err != nil {
				log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
			}
		case config.Mercury:
			checkData = []byte(k.cfg.UpkeepCheckData)
			var err error
			if k.cfg.VerifiableLoadTest {
				upkeepAddr, deployUpkeepTx, _, err = verifiable_load_streams_lookup_upkeep_wrapper.DeployVerifiableLoadStreamsLookupUpkeep(
					k.buildTxOpts(ctx),
					k.client,
					common.HexToAddress(k.cfg.Registrar),
					k.cfg.UseArbBlockNumber,
				)
			} else {
				upkeepAddr, deployUpkeepTx, _, err = streams_lookup_upkeep_wrapper.DeployStreamsLookupUpkeep(
					k.buildTxOpts(ctx),
					k.client,
					big.NewInt(k.cfg.UpkeepTestRange),
					big.NewInt(k.cfg.UpkeepInterval),
					true,  /* useArbBlock */
					true,  /* staging */
					false, /* verify mercury response */
				)
			}
			if err != nil {
				log.Fatal(i, ": Deploy Upkeep failed - ", err)
			}
			k.waitDeployment(ctx, deployUpkeepTx)
			log.Println(i, upkeepAddr.Hex(), ": Upkeep deployed - ", helpers.ExplorerLink(k.cfg.ChainID, deployUpkeepTx.Hash()))
			registerUpkeepTx, err = deployer.RegisterUpkeep(k.buildTxOpts(ctx),
				upkeepAddr, k.cfg.UpkeepGasLimit, k.fromAddr, checkData, []byte{},
			)
			if err != nil {
				log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
			}
		case config.LogTrigger:
			var err error
			upkeepAddr, deployUpkeepTx, logUpkeepCounter, err = log_upkeep_counter_wrapper.DeployLogUpkeepCounter(
				k.buildTxOpts(ctx),
				k.client,
				big.NewInt(k.cfg.UpkeepTestRange),
			)
			if err != nil {
				log.Fatal(i, ": Deploy Upkeep failed - ", err)
			}
			logTriggerConfigType := abi.MustNewType("tuple(address contractAddress, uint8 filterSelector, bytes32 topic0, bytes32 topic1, bytes32 topic2, bytes32 topic3)")
			logTriggerConfig, err := abi.Encode(map[string]interface{}{
				"contractAddress": upkeepAddr,
				"filterSelector":  0,                                                                    // no indexed topics filtered
				"topic0":          "0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d", // event sig for Trigger()
				"topic1":          "0x",
				"topic2":          "0x",
				"topic3":          "0x",
			}, logTriggerConfigType)
			if err != nil {
				log.Fatal("failed to encode log trigger config", err)
			}
			k.waitDeployment(ctx, deployUpkeepTx)
			log.Println(i, upkeepAddr.Hex(), ": Upkeep deployed - ", helpers.ExplorerLink(k.cfg.ChainID, deployUpkeepTx.Hash()))
			registerUpkeepTx, err = deployer.RegisterUpkeepV2(k.buildTxOpts(ctx),
				upkeepAddr, k.cfg.UpkeepGasLimit, k.fromAddr, 1, []byte{}, logTriggerConfig, []byte{},
			)
			if err != nil {
				log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
			}

			// Start up log trigger cycle
			logUpkeepStartTx, err := logUpkeepCounter.Start(k.buildTxOpts(ctx))
			if err != nil {
				log.Fatal("failed to start log upkeep counter", err)
			}
			if err = k.waitTx(ctx, logUpkeepStartTx); err != nil {
				log.Fatalf("Log upkeep Start() failed for upkeepId: %s, error is %s", upkeepAddr.Hex(), err.Error())
			}
			log.Println(i, upkeepAddr.Hex(), ": Log upkeep successfully started - ", helpers.ExplorerLink(k.cfg.ChainID, logUpkeepStartTx.Hash()))
		default:
			log.Fatal("unexpected upkeep type")
		}

		if err := k.waitTx(ctx, registerUpkeepTx); err != nil {
			log.Fatalf("RegisterUpkeep failed for upkeepId: %s, error is %s", upkeepAddr.Hex(), err.Error())
		}
		log.Println(i, upkeepAddr.Hex(), ": Upkeep registered - ", helpers.ExplorerLink(k.cfg.ChainID, registerUpkeepTx.Hash()))

		upkeepAddrs = append(upkeepAddrs, upkeepAddr)
	}

	var upkeepGetter activeUpkeepGetter
	upkeepCount := big.NewInt(k.cfg.UpkeepCount) // second arg in GetActiveUpkeepIds (on registry)
	{
		var err error
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
		case keeper.RegistryVersion_2_1:
			upkeepGetter, err = iregistry21.NewIKeeperRegistryMaster(
				registryAddr,
				k.client,
			)
		default:
			panic("unexpected registry address")
		}
		if err != nil {
			log.Fatal("Registry failed: ", err)
		}
	}

	activeUpkeepIds := k.getActiveUpkeepIds(ctx, upkeepGetter, big.NewInt(existingCount), upkeepCount)

	for index, upkeepAddr := range upkeepAddrs {
		// Approve
		k.approveFunds(ctx, registryAddr)

		upkeepId := activeUpkeepIds[index]

		// Fund
		addFundsTx, err := deployer.AddFunds(k.buildTxOpts(ctx), upkeepId, k.addFundsAmount)
		if err != nil {
			log.Fatal(upkeepId, upkeepAddr.Hex(), ": AddFunds failed - ", err)
		}

		// Onchain transaction
		if err := k.waitTx(ctx, addFundsTx); err != nil {
			log.Fatalf("AddFunds failed for upkeepId: %s, and upkeepAddr: %s, error is: %s", upkeepId, upkeepAddr.Hex(), err.Error())
		}

		log.Println(upkeepId, upkeepAddr.Hex(), ": Upkeep funded - ", helpers.ExplorerLink(k.cfg.ChainID, addFundsTx.Hash()))
	}

	// set administrative offchain config for mercury upkeeps
	if (k.cfg.UpkeepType == config.Mercury || k.cfg.UpkeepType == config.LogTriggeredFeedLookup) && k.cfg.RegistryVersion == keeper.RegistryVersion_2_1 {
		reg21, err := iregistry21.NewIKeeperRegistryMaster(registryAddr, k.client)
		if err != nil {
			log.Fatalf("cannot create registry 2.1: %v", err)
		}
		v, err := reg21.TypeAndVersion(nil)
		if err != nil {
			log.Fatalf("failed to fetch type and version from registry 2.1: %v", err)
		}
		log.Printf("registry version is %s", v)
		log.Printf("active upkeep ids: %v", activeUpkeepIds)

		adminBytes, err := json.Marshal(streams.UpkeepPrivilegeConfig{
			MercuryEnabled: true,
		})
		if err != nil {
			log.Fatalf("failed to marshal upkeep privilege config: %v", err)
		}

		for _, id := range activeUpkeepIds {
			tx, err2 := reg21.SetUpkeepPrivilegeConfig(k.buildTxOpts(ctx), id, adminBytes)
			if err2 != nil {
				log.Fatalf("failed to upkeep privilege config: %v", err2)
			}
			err2 = k.waitTx(ctx, tx)
			if err2 != nil {
				log.Fatalf("failed to wait for tx: %v", err2)
			}
			log.Printf("upkeep privilege config is set for %s", id.String())

			info, err2 := reg21.GetUpkeep(nil, id)
			if err2 != nil {
				log.Fatalf("failed to fetch upkeep id %s from registry 2.1: %v", id, err2)
			}
			min, err2 := reg21.GetMinBalanceForUpkeep(nil, id)
			if err2 != nil {
				log.Fatalf("failed to fetch upkeep id %s from registry 2.1: %v", id, err2)
			}
			log.Printf("    Balance: %s", info.Balance)
			log.Printf("Min Balance: %s", min.String())
		}
	}

	fmt.Println()
}

// setKeepers set the keeper list for a registry
func (k *Keeper) setKeepers(ctx context.Context, cls []cmd.HTTPClient, deployer keepersDeployer, keepers, owners []common.Address) {
	if len(keepers) > 0 {
		log.Println("Set keepers...")
		opts := k.buildTxOpts(ctx)
		setKeepersTx, err := deployer.SetKeepers(ctx, opts, cls, keepers, owners)
		if err != nil {
			log.Fatal("SetKeepers failed: ", err)
		}

		if err = k.waitTx(ctx, setKeepersTx); err != nil {
			log.Fatalf("SetKeepers failed, error is: %s", err.Error())
		}

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
