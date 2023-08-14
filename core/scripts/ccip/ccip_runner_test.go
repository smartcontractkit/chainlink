package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/clo"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/csv"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/dione"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/metis/printing"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea/deployment_io"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea/deployments"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/secrets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
)

var (
	// Change these values
	sourceChain        = rhea.ArbitrumGoerli
	destChain          = rhea.Sepolia
	ENV                = dione.StagingBeta
	currentVersion     = ""
	upgradeLaneVersion = ""

	// These will automatically populate or error if the lane doesn't exist
	SOURCE      = laneMapping[ENV][sourceChain][destChain]
	DESTINATION = laneMapping[ENV][destChain][sourceChain]
)

// These functions can be run as a test (prefix with Test) with the following config
// DATABASE_URL
// Use "-v" as a Go tool argument for streaming log output.

// TestCCIP can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
// SEED_KEY   private key used for multi-user tests. Not needed when using the "deploy" command.
// COMMAND    what function to run e.g. "deploy", "setConfig", or "gas".
func TestCCIP(t *testing.T) {
	ownerKey := checkOwnerKeyAndSetupChain(t)
	command := os.Getenv("COMMAND")
	if command == "" {
		t.Log("No command given, skipping ccip-script. This is intended behaviour for automated testing.")
		t.SkipNow()
	}
	// The seed key is used to generate 10 keys from a single key by changing the
	// first character of the given seed with the digits 0-9
	seedKey := os.Getenv("SEED_KEY")
	if seedKey == "" {
		seedKey = ownerKey
	}

	// Configures a client to run tests with using the network defaults and given keys.
	// After updating any contracts be sure to update the network defaults to reflect
	// those changes.
	client := NewCcipClient(t, SOURCE, DESTINATION, ownerKey, seedKey)

	switch command {
	case "ccipSend": // Sends a basic tx with customizable contents
		client.ccipSendBasicTx(t)
	case "deployPingPong":
		rhea.DeployPingPongDapps(t, &SOURCE, &DESTINATION)
	case "startPingPong": // Starts and unpauses the PingPong dapp that is on the source chain.
		client.startPingPong(t)
	case "stopPingPong": // Stops the PingPong dapp by pausing the source chain dapp.
		client.setPingPongPaused(t, true)
	case "printSpecs":
		printing.PrintJobSpecs(ENV, SOURCE, DESTINATION, currentVersion)
	case "setConfig": // Set the config to the commitStore and the offramp
		client.SetOCR2Config(ENV)
		clientOtherWayAround := NewCcipClient(t, DESTINATION, SOURCE, ownerKey, seedKey)
		clientOtherWayAround.SetOCR2Config(ENV)
	case "reemitEvents": // re-set current onchain config to re-emit deployment events
		client.reemitEvents(t, DESTINATION)
	case "setOnRampFeeConfig":
		client.setOnRampFeeConfig(t, &SOURCE)
	case "applyFeeTokensUpdates":
		client.applyFeeTokensUpdates(t, &SOURCE)
	case "batching": // Submit 10 txs. This should result in the txs being batched together
		client.ScalingAndBatching(t)
	case "gas":
		client.TestGasVariousTxs(t)
	case "uncurseSourceARM":
		client.uncurseSourceARM(t)
	case "setConfigARM":
		client.setConfigARM(t)
	case "executeManually":
		client.ExecuteManually(t, &DESTINATION)
	case "wip":
		client.wip(t, &SOURCE, &DESTINATION)
	case "":
		t.Log("No command given, exit successfully")
		t.SkipNow()
	default:
		t.Errorf("Unknown command \"%s\"", command)
	}
}

// TestDeployChain can be run as a test with the following config
// NOTE: deploy chain always runs for all chains
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
func TestRheaDeployChains(t *testing.T) {
	DoForEachChain(t, ENV, func(chain rhea.EvmDeploymentConfig) {
		err := rhea.DeployToNewChain(&chain)
		if err != nil {
			t.Error(err)
		}
		deployment_io.WriteChainConfigToFile(ENV, &chain)
	})
}

// TestRheaDeployLane can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
func TestRheaDeployLane(t *testing.T) {
	key := checkOwnerKeyAndSetupChain(t)
	rhea.DeployLanes(t, &SOURCE, &DESTINATION)
	deployment_io.PrettyPrintLanes(ENV, &SOURCE, &DESTINATION)

	client := NewCcipClient(t, SOURCE, DESTINATION, key, key)
	client.SetOCR2Config(ENV)
	clientOtherWayAround := NewCcipClient(t, DESTINATION, SOURCE, key, key)
	clientOtherWayAround.SetOCR2Config(ENV)

	client.startPingPong(t)

	env := ENV
	if SOURCE.ChainConfig.EvmChainId == 1337 || DESTINATION.ChainConfig.EvmChainId == 1337 {
		env = dione.Prod_Swift
	}
	don := dione.NewDON(env, logger.TestLogger(t))
	don.ClearAllJobs(ccip.ChainName(int64(SOURCE.ChainConfig.EvmChainId)), ccip.ChainName(int64(DESTINATION.ChainConfig.EvmChainId)))
	don.AddTwoWaySpecs(SOURCE, DESTINATION)

	// Sometimes jobs don't get added correctly. This script looks for missing jobs
	// and attempts to add them.
	don.AddMissingSpecs(DESTINATION, SOURCE, "")
	don.AddMissingSpecs(SOURCE, DESTINATION, "")
}

// TestRheaDeployUpgradeLane performs first part of blue-green deployment based on the configuration stored under rhea.EvmDeploymentConfig.UpgradeLaneConfig
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
// IMPORTANT: There has to be additional router deployed and its address saved to rhea.EVMChainConfig.UpgradeRouter
//
// Entire process is visualized https://docs.google.com/presentation/d/10Bb9e3naEqqMfbyR_ML343TJd7dGHkQFd66G0I2J5y4/edit#slide=id.g1df2dd79c37_0_0
// Upgrade deployment includes:
// * onRamp/offRamp/commitStore deployment if needed (onRamp has set prevOnRamp to the onRamp that is currently deployed)
// * registering onRamp and offRamp components in upgrade routers
// * deploy ping pong dapp
// * setOCR2Config with proper addresses (but job specs are not updated at this point)
// All contracts' addresses deployed at this point are saved under rhea.EvmDeploymentConfig.UpgradeLaneConfig
func TestRheaDeployUpgradeLane(t *testing.T) {
	key := checkOwnerKeyAndSetupChain(t)
	err := rhea.DeployUpgradeRouters(&SOURCE, &DESTINATION)
	if err != nil {
		t.Error(err)
		return
	}
	rhea.DeployUpgradeLanes(t, &SOURCE, &DESTINATION)
	// Print all deployed contracts
	deployment_io.PrettyPrintLanes(ENV, &SOURCE, &DESTINATION)

	// Set OCR2 & Dynamic configs pointing to upgrade routers
	client := NewUpgradeLaneCcipClient(t, SOURCE, DESTINATION, key, key)
	client.SetOCR2Config(ENV)
	client.SetDynamicConfigOnRamp(t)

	clientOtherWayAround := NewUpgradeLaneCcipClient(t, DESTINATION, SOURCE, key, key)
	clientOtherWayAround.SetOCR2Config(ENV)
	clientOtherWayAround.SetDynamicConfigOnRamp(t)

	// At this point PingPong should properly work using Upgrade Routers and Upgrade Lane
	client.startPingPong(t)
	// Please go to /deployments/ENV/upgrade_line/SOURCE and /deployments/ENV/upgrade_line/DESTINATION
	// and replace current UpgradeLaneConfig in ENV.go file with the json content
	// Also remember to set UpgradeLaneConfig.DeployCommitStore and UpgradeLaneConfig.DeployRamp to false

	// Add new job specs for Upgrade Lanes
	env := ENV
	if SOURCE.ChainConfig.EvmChainId == 1337 || DESTINATION.ChainConfig.EvmChainId == 1337 {
		env = dione.Prod_Swift
	}

	don := dione.NewDON(env, logger.TestLogger(t))
	don.ClearAllLaneJobsByVersion(ccip.ChainName(int64(SOURCE.ChainConfig.EvmChainId)), ccip.ChainName(int64(DESTINATION.ChainConfig.EvmChainId)), upgradeLaneVersion)
	don.AddTwoWaySpecsByVersion(SOURCE.OnlyEvmConfig(), SOURCE.UpgradeLaneConfig, DESTINATION.OnlyEvmConfig(), DESTINATION.UpgradeLaneConfig, upgradeLaneVersion)

	// Sometimes jobs don't get added correctly. This script looks for missing jobs
	// and attempts to add them.
	don.AddMissingSpecsByLanes(DESTINATION.OnlyEvmConfig(), DESTINATION.UpgradeLaneConfig, SOURCE.OnlyEvmConfig(), SOURCE.UpgradeLaneConfig, upgradeLaneVersion)
	don.AddMissingSpecsByLanes(SOURCE.OnlyEvmConfig(), SOURCE.UpgradeLaneConfig, DESTINATION.OnlyEvmConfig(), DESTINATION.UpgradeLaneConfig, upgradeLaneVersion)
}

// TestRheaPromoteUpgradeLaneDeployment promotes the deployment from upgrade lane
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
//
// Deployment promotion applies following steps:
// - enables upgrade offRamp by registering it in dest router and updating OCR2 config
// - enabled upgrade onRamp by setting dynamicConfig and setting it in source router
// - populating new jobs specs to the NOPs
//
// In the meantime it:
// - populates addresses from rhea.EvmDeploymentConfig.UpgradeLaneConfig to rhea.EvmDeploymentConfig.LaneConfig
// - sets addresses in rhea.EvmDeploymentConfig.UpgradeLaneConfig to 0 (this is a sign that there is no pending deployment)
func TestRheaPromoteUpgradeLaneDeployment(t *testing.T) {
	key := checkOwnerKeyAndSetupChain(t)

	rhea.EnableUpgradeOffRamps(t, &SOURCE, &DESTINATION, func(src *rhea.EvmDeploymentConfig, dst *rhea.EvmDeploymentConfig) {
		client := NewCcipClient(t, *src, *dst, key, key)
		client.SetOCR2Config(ENV)
	})
	rhea.EnableUpgradeOnRamps(t, &SOURCE, &DESTINATION, func(src *rhea.EvmDeploymentConfig, dst *rhea.EvmDeploymentConfig) {
		// Points to UpgradeLanes and regular Router
		client := NewCcipClientByLane(t, src.OnlyEvmConfig(), src.UpgradeLaneConfig, dst.OnlyEvmConfig(), dst.UpgradeLaneConfig, key, key)
		client.SetDynamicConfigOnRamp(t)
	})
	deployment_io.PrettyPrintLanes(ENV, &SOURCE, &DESTINATION)
}

// TestRheaPostUpgradeDeploymentClean removed job specs with currentVersion. MAKE SURE you bumped this value before running this script
func TestRheaPostUpgradeDeploymentClean(t *testing.T) {
	checkOwnerKeyAndSetupChain(t)

	// Remove job specs from previous deployment
	env := ENV
	if SOURCE.ChainConfig.EvmChainId == 1337 || DESTINATION.ChainConfig.EvmChainId == 1337 {
		env = dione.Prod_Swift
	}

	don := dione.NewDON(env, logger.TestLogger(t))
	don.ClearAllLaneJobsByVersion(ccip.ChainName(int64(SOURCE.ChainConfig.EvmChainId)), ccip.ChainName(int64(DESTINATION.ChainConfig.EvmChainId)), currentVersion)
}

// TestDione can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
func TestDione(t *testing.T) {
	checkOwnerKeyAndSetupChain(t)

	env := ENV
	if SOURCE.ChainConfig.EvmChainId == 1337 || DESTINATION.ChainConfig.EvmChainId == 1337 {
		env = dione.Prod_Swift
	}

	don := dione.NewDON(env, logger.TestLogger(t))
	don.ClearAllJobs(ccip.ChainName(int64(SOURCE.ChainConfig.EvmChainId)), ccip.ChainName(int64(DESTINATION.ChainConfig.EvmChainId)))
	don.AddTwoWaySpecs(SOURCE, DESTINATION)

	// Sometimes jobs don't get added correctly. This script looks for missing jobs
	// and attempts to add them.
	don.AddMissingSpecs(DESTINATION, SOURCE, "")
	don.AddMissingSpecs(SOURCE, DESTINATION, "")
}

// TestDionePopulateNodeKeys
// 1. gets the keys from the nodes based upon ENV (OCR2Keys EthKeys PeerId) using json/credentials/ for auth
// 2. writes the node keys into a file in json/nodes/
func TestDionePopulateNodeKeys(t *testing.T) {
	checkOwnerKey(t)

	env := ENV
	if SOURCE.ChainConfig.EvmChainId == 1337 || DESTINATION.ChainConfig.EvmChainId == 1337 {
		env = dione.Prod_Swift
	}

	don := dione.NewDON(env, logger.TestLogger(t))
	don.LoadCurrentNodeParams()
	don.WriteToFile()
}

// TestUpdateAllLanes
// 1. updates all the available lanes with new offramp, onramp, commit store
// 2. creates new jobs
// 3. set ocrConfig for both
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
func TestUpdateAllLanes(t *testing.T) {
	ownerKey := checkOwnerKey(t)
	if _, ok := laneMapping[ENV]; !ok {
		t.Error("set environment not supported")
	}

	env := ENV
	if SOURCE.ChainConfig.EvmChainId == 1337 || DESTINATION.ChainConfig.EvmChainId == 1337 {
		env = dione.Prod_Swift
	}

	don := dione.NewDON(env, logger.TestLogger(t))

	// Potential todo: remove old deployment artifact permissions
	// Optimizations:
	// 		Concurrent chain contracts deployment before any lange deployment
	// 		Concurrent lane contract deployment for non-intersecting lanes
	// 		Concurrent lane contract deployment within a bidirectional deploy
	// 		Not waiting for mining, self incrementing the nonce

	// 		Downsides: less control and worse retry experience
	// 			As failures should be very rare this is probably worth it
	upgradeLane := func(source, dest rhea.EvmDeploymentConfig) {
		if !source.LaneConfig.DeploySettings.DeployLane {
			source.Logger.Warnf("Please set \"DeployRamp and DeployCommitStore\" to true for the given EvmChainConfigs and make sure "+
				"the right ones are set. Source: %d, Dest %d", source.ChainConfig.EvmChainId, dest.ChainConfig.EvmChainId)
			return
		}
		if !dest.LaneConfig.DeploySettings.DeployLane {
			dest.Logger.Warnf("Please set \"DeployRamp and DeployCommitStore\" to true for the given EvmChainConfigs and make sure "+
				"the right ones are set. Source: %d, Dest %d", dest.ChainConfig.EvmChainId, source.ChainConfig.EvmChainId)
			return
		}
		if source.ChainConfig.DeploySettings.DeployRouter || dest.ChainConfig.DeploySettings.DeployRouter {
			dest.Logger.Warnf("Routers should never be set to true Source: %d, Dest %d", dest.ChainConfig.EvmChainId, source.ChainConfig.EvmChainId)
			return
		}
		// Removes any old job specs
		don.ClearAllJobs(ccip.ChainName(int64(source.ChainConfig.EvmChainId)), ccip.ChainName(int64(dest.ChainConfig.EvmChainId)))
		// Deploys the new contracts and updates `source` and `dest`
		rhea.DeployLanes(t, &source, &dest)
		// Prints the new config and writes them to file
		deployment_io.PrettyPrintLanes(ENV, &source, &dest)
		// Add new job specs
		don.AddTwoWaySpecs(source, dest)
		// Set the OCR2 config on the source contracts
		client := NewCcipClient(t, source, dest, ownerKey, ownerKey)
		client.SetOCR2Config(ENV)
		// Set the OCR2 config on the destination contracts
		client = NewCcipClient(t, dest, source, ownerKey, ownerKey)
		client.SetOCR2Config(ENV)
		// Starts the ping pong dapp
		client.startPingPong(t)
	}

	// This script only deploys new lane contracts. Please deploy any new chain contracts
	// and update the config before running this.

	DoForEachBidirectionalLane(t, ENV, upgradeLane)
}

// How to add tokens in 3 steps
// Add token to config
// **	If the token is new add it to `models.go` and set its symbol, decimals and price
// **	Add it to the chain config in e.g. prod.go
// **	Leave the pool address empty
// ** 	Depending on the pool type fill in the token address or not (wrapped doesn't have a token so leave it empty)
// **   Set DeployTokenPools to `true` for chains that need the pool deployed
//
// Run `TestRheaDeployChains` to deploy the new pools
// ** 	Run output should be written to console & ./json/deployments/env/chain/....
// ** 	Modify the chain config to include the new info
// **   Set DeployTokenPools back to `false` where changed
//
// Run TestSyncTokens
// ** 	This should set the correct config on each ramp and token pool based on previous steps
func TestSyncTokens(t *testing.T) {
	ownerKey := checkOwnerKey(t)
	DoForEachLane(t, ENV, func(source rhea.EvmDeploymentConfig, destination rhea.EvmDeploymentConfig) {
		client := NewCcipClient(t, source, destination, ownerKey, ownerKey)
		client.SyncTokenPools()
	})
}

// TestPrintNodeBalances can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
func TestPrintNodeBalances(t *testing.T) {
	checkOwnerKeyAndSetupChain(t)

	don := dione.NewOfflineDON(ENV, logger.TestLogger(t))

	printing.PrintNodeBalances(&SOURCE, don.GetSendingKeys(SOURCE.ChainConfig.EvmChainId))
	printing.PrintNodeBalances(&DESTINATION, don.GetSendingKeys(DESTINATION.ChainConfig.EvmChainId))
}

func TestFundNodes(t *testing.T) {
	key := checkOwnerKeyAndSetupChain(t)

	don := dione.NewOfflineDON(ENV, logger.TestLogger(t))
	don.FundNodeKeys(&SOURCE, key, big.NewInt(4e18), big.NewInt(4e18))
}

func TestFundPingPong(t *testing.T) {
	minimumBalance := new(big.Int).Mul(big.NewInt(20), big.NewInt(1e18))

	DoForEachBidirectionalLane(t, ENV, func(source rhea.EvmDeploymentConfig, destination rhea.EvmDeploymentConfig) {
		FundPingPong(t, source, minimumBalance)
		FundPingPong(t, destination, minimumBalance)
	})
}

// TestPrintAllNodeBalancesPerEnv can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
// It will print the node balances for all chains where the given `env` is deployed
func TestPrintAllNodeBalancesPerEnv(t *testing.T) {
	ownerKey := checkOwnerKey(t)

	for _, source := range chainMapping[ENV] {
		source.SetupChain(t, ownerKey)
		don := dione.NewOfflineDON(ENV, logger.TestLogger(t))
		printing.PrintNodeBalances(&source, don.GetSendingKeys(source.ChainConfig.EvmChainId))
	}
}

// TestFundAllNodesPerEnv can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
// It will fund the node balances for all chains where the given `env` is deployed
func TestFundAllNodesPerEnv(t *testing.T) {
	ownerKey := checkOwnerKey(t)
	for _, source := range chainMapping[ENV] {
		source.SetupChain(t, ownerKey)
		don := dione.NewOfflineDON(ENV, logger.TestLogger(t))
		don.FundNodeKeys(&source, ownerKey, big.NewInt(5e18), big.NewInt(4e18))
	}
}

// TestWriteNodesWalletsToCSV according to set ENV it writes a CSV file in csv/node-wallets/ directory
// with all the node wallets for the given ENV per chain
func TestWriteNodesWalletsToCSV(t *testing.T) {
	headers := []string{"Environment", "Chain Name", "Chain Id", "Address"}
	path := "csv/node-wallets"
	fileName := fmt.Sprintf("%s-%s.csv", string(ENV), time.Now().Format("2006-01-02 15:04:05"))
	filePath := fmt.Sprintf("%s/%s", path, fileName)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("failed to create directory %s: %s", path, err)
	}
	csv.PrepareCsvFile(filePath, headers)
	don := dione.NewOfflineDON(ENV, logger.TestLogger(t))
	DoForEachChain(t, ENV, func(chain rhea.EvmDeploymentConfig) {
		records := don.GetAllNodesWallets(chain.ChainConfig.EvmChainId)
		csv.AppendToFile(filePath, records, ccip.ChainName(int64(chain.ChainConfig.EvmChainId)), ENV)
	})
}

func checkOwnerKeyAndSetupChain(t *testing.T) string {
	ownerKey := checkOwnerKey(t)
	SOURCE.SetupChain(t, ownerKey)
	DESTINATION.SetupChain(t, ownerKey)

	return ownerKey
}

// TestCLO prepares chains and lanes env according to set SOURCE, DESTINATION, ENV and run any provided function at the end
// It uses set configuration by selected ENV and overrides any variables provided by calling CLO API configuration
// You must set additional env variables FMS_AUTH_TOKEN, CLO_QUERY_URL for CLO requests
func TestCLO(t *testing.T) {
	ownerKey := checkOwnerKeyAndSetupChain(t)
	seedKey := os.Getenv("SEED_KEY")
	if seedKey == "" {
		t.Error("must set seed key")
	}

	// Set configuration queried from CLO, laneID is lane id from CLO API
	sourceContracts, destContracts := clo.GetTargetChainsContracts(t, SOURCE.ChainConfig.EvmChainId, DESTINATION.ChainConfig.EvmChainId)
	clo.SetChainConfig(sourceContracts, &SOURCE)
	clo.SetChainConfig(destContracts, &DESTINATION)
	legA, legB := clo.GetTargetLaneConfig(t, SOURCE.ChainConfig.EvmChainId, DESTINATION.ChainConfig.EvmChainId, "10")
	clo.SetLaneConfig(legA, &SOURCE, &DESTINATION)
	clo.SetLaneConfig(legB, &DESTINATION, &SOURCE)

	deployment_io.PrettyPrintLanes(ENV, &SOURCE, &DESTINATION)
	client := NewCcipClient(t, SOURCE, DESTINATION, ownerKey, seedKey)
	// Add any function after it for pulled configuration ex:
	client.startPingPong(t)
}

// This ALWAYS uses the production env
func Test__PROD__SetAllowListAllLanes(t *testing.T) {
	ownerKey := checkOwnerKey(t)

	// Simply comment out the lanes that are not needed.
	allProdLanes := []*rhea.EvmDeploymentConfig{
		&deployments.Prod_SepoliaToOptimismGoerli,
		&deployments.Prod_SepoliaToAvaxFuji,
		&deployments.Prod_SepoliaToArbitrumGoerli,
		&deployments.Prod_SepoliaToPolygonMumbai,
		// Quorum allowList is turned off for now, do not uncomment
		//&deployments.Prod_SepoliaToQuorum,

		&deployments.Prod_AvaxFujiToSepolia,
		&deployments.Prod_AvaxFujiToOptimismGoerli,
		&deployments.Prod_AvaxFujiToPolygonMumbai,

		&deployments.Prod_OptimismGoerliToAvaxFuji,
		&deployments.Prod_OptimismGoerliToSepolia,

		&deployments.Prod_ArbitrumGoerliToSepolia,

		&deployments.Prod_PolygonMumbaiToSepolia,
		&deployments.Prod_PolygonMumbaiToAvaxFuji,
	}

	for _, lane := range allProdLanes {
		lane.SetupChain(t, ownerKey)
		client := CCIPClient{Source: NewSourceClient(t, lane.OnlyEvmConfig(), lane.LaneConfig)}
		client.Source.Owner = rhea.GetOwner(t, ownerKey, client.Source.ChainId, lane.ChainConfig.GasSettings)

		client.setAllowList(t)
	}
}

// TestUpdateLaneARMAddress can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
// It applies ARM address in all contracts that require it. ARM address is taken from rhea.EVMChainConfig.ARM
// SourceChain.ARM is applied to
// * OnRamp (via SetDynamicConfig)
//
// DestinationChain.ARM is applied to:
// * OffRamp (via setOCR2Config)
// * CommitStore (via setOCR2Config)
func TestUpdateLaneARMAddress(t *testing.T) {
	key := checkOwnerKeyAndSetupChain(t)

	client := NewCcipClient(t, SOURCE, DESTINATION, key, key)
	client.SetDynamicConfigOnRamp(t)
	client.SetOCR2Config(ENV)

	clientOtherWayAround := NewCcipClient(t, DESTINATION, SOURCE, key, key)
	clientOtherWayAround.SetDynamicConfigOnRamp(t)
	clientOtherWayAround.SetOCR2Config(ENV)
}

func TestFinalityTags(t *testing.T) {
	checkOwnerKeyAndSetupChain(t)

	// Ensure that HeaderByBlockNumber using finality tag works.
	finalityTagChains := []uint64{420, 43113, 421613, 11155111, 1, 10, 43114, 42161}
	for _, chainID := range finalityTagChains {
		client, err := ethclient.Dial(secrets.GetRPC(chainID))
		require.NoError(t, err)
		f, err := client.HeaderByNumber(context.Background(), big.NewInt(rpc.FinalizedBlockNumber.Int64()))
		require.NoError(t, err, "chainID: %d", chainID)
		fmt.Println(f.Number.String())
	}
}
