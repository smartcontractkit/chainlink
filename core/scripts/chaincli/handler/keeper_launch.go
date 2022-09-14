package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	registry12 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/web"
)

type startedNodeData struct {
	url     string
	err     error
	cleanup func()
}

// LaunchAndTest launches keeper registry, chainlink nodes, upkeeps and start performing.
// 1. launch chainlink node using docker image
// 2. get keeper registry instance, deploy if needed
// 3. deploy upkeeps
// 4. create keeper jobs
// 5. fund nodes if needed
// 6. set keepers in the registry
// 7. withdraw funds after tests are done -> TODO: wait until tests are done instead of cancel manually
func (k *Keeper) LaunchAndTest(ctx context.Context, withdraw bool) {
	lggr, closeLggr := logger.NewLogger()
	defer closeLggr()

	var extraEvars []string
	if k.cfg.OCR2Keepers {
		extraEvars = []string{
			"FEATURE_OFFCHAIN_REPORTING2=true",
			"FEATURE_LOG_POLLER=true",
			"P2P_NETWORKING_STACK=V2",
			"CHAINLINK_TLS_PORT=0",
			fmt.Sprintf("P2PV2_LISTEN_ADDRESSES=127.0.0.1:%d", 8080),
		}
	}

	// Run chainlink nodes and create jobs
	startedNodes := make([]startedNodeData, k.cfg.KeepersCount)
	var wg sync.WaitGroup
	for i := 0; i < k.cfg.KeepersCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			startedNodes[i] = startedNodeData{}

			// Run chainlink node
			var err error
			if startedNodes[i].url, startedNodes[i].cleanup, err = k.launchChainlinkNode(ctx, 6688+i, extraEvars...); err != nil {
				startedNodes[i].err = fmt.Errorf("failed to launch chainlink node: %s", err)
				return
			}
		}(i)
	}
	wg.Wait()

	// Deploy keeper registry or get an existing one
	upkeepCount, registryAddr, deployer := k.prepareRegistry(ctx)

	// Approve keeper registry
	k.approveFunds(ctx, registryAddr)

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registryAddr, deployer, upkeepCount)

	// Prepare keeper addresses and owners
	var keepers []common.Address
	var owners []common.Address
	for _, startedNode := range startedNodes {
		if startedNode.err != nil {
			log.Println("Failed to start node: ", startedNode.err)
			continue
		}

		// Create authenticated client
		var cl cmd.HTTPClient
		var err error
		cl, err = authenticate(startedNode.url, defaultChainlinkNodeLogin, defaultChainlinkNodePassword, lggr)
		if err != nil {
			log.Fatal("Authentication failed, ", err)
		}

		// Get node's wallet address
		var nodeAddrHex string
		if nodeAddrHex, err = getNodeAddress(cl); err != nil {
			log.Println("Failed to get node addr: ", err)
			continue
		}
		nodeAddr := common.HexToAddress(nodeAddrHex)

		// Create keepers
		if err = k.createKeeperJob(cl, registryAddr.Hex(), nodeAddr.Hex()); err != nil {
			log.Println("Failed to create keeper job: ", err)
			continue
		}
		log.Println("Keeper job has been successfully created in the Chainlink node with address ", startedNode.url)

		// Fund node if needed
		fundAmt, ok := (&big.Int{}).SetString(k.cfg.FundNodeAmount, 10)
		if !ok {
			log.Printf("failed to parse FUND_CHAINLINK_NODE: %s", k.cfg.FundNodeAmount)
			continue
		}
		if fundAmt.Cmp(big.NewInt(0)) != 0 {
			if err = k.sendEth(ctx, nodeAddr, fundAmt); err != nil {
				log.Println("Failed to fund chainlink node: ", err)
				continue
			}
		}

		keepers = append(keepers, nodeAddr)
		owners = append(owners, k.fromAddr)
	}

	// Set Keepers
	k.setKeepers(ctx, deployer, keepers, owners)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
	log.Println("Stopping...")

	// Cleanup resources
	for _, startedNode := range startedNodes {
		if startedNode.err == nil && startedNode.cleanup != nil {
			startedNode.cleanup()
		}
	}

	// Cancel upkeeps and withdraw funds
	if withdraw {
		isVersion12 := k.cfg.RegistryVersion == keeper.RegistryVersion_1_2
		log.Println("Canceling upkeeps...")
		if isVersion12 {
			registry, err := registry12.NewKeeperRegistry(
				registryAddr,
				k.client,
			)
			if err != nil {
				log.Fatal("Registry failed: ", err)
			}
			activeUpkeepIds := k.getActiveUpkeepIds(ctx, registry, big.NewInt(0), big.NewInt(0))

			if err = k.cancelAndWithdrawActiveUpkeeps(ctx, activeUpkeepIds, deployer); err != nil {
				log.Fatal("Failed to cancel upkeeps: ", err)
			}
		} else {
			if err := k.cancelAndWithdrawUpkeeps(ctx, big.NewInt(upkeepCount), deployer); err != nil {
				log.Fatal("Failed to cancel upkeeps: ", err)
			}
		}
		log.Println("Upkeeps successfully canceled")
	}
}

// cancelAndWithdrawActiveUpkeeps cancels all active upkeeps and withdraws funds for registry 1.2
func (k *Keeper) cancelAndWithdrawActiveUpkeeps(ctx context.Context, activeUpkeepIds []*big.Int, canceller canceller) error {
	var err error
	for i := 0; i < len(activeUpkeepIds); i++ {
		var tx *ethtypes.Transaction
		upkeepId := activeUpkeepIds[i]
		if tx, err = canceller.CancelUpkeep(k.buildTxOpts(ctx), upkeepId); err != nil {
			return fmt.Errorf("failed to cancel upkeep %s: %s", upkeepId.String(), err)
		}
		k.waitTx(ctx, tx)

		if tx, err = canceller.WithdrawFunds(k.buildTxOpts(ctx), upkeepId, k.fromAddr); err != nil {
			return fmt.Errorf("failed to withdraw upkeep %s: %s", upkeepId.String(), err)
		}
		k.waitTx(ctx, tx)

		log.Printf("Upkeep %s successfully canceled and refunded: ", upkeepId.String())
	}

	var tx *ethtypes.Transaction
	if tx, err = canceller.RecoverFunds(k.buildTxOpts(ctx)); err != nil {
		return fmt.Errorf("failed to recover funds: %s", err)
	}
	k.waitTx(ctx, tx)

	return nil
}

// cancelAndWithdrawUpkeeps cancels all upkeeps for 1.1 registry and withdraws funds
func (k *Keeper) cancelAndWithdrawUpkeeps(ctx context.Context, upkeepCount *big.Int, canceller canceller) error {
	var err error
	for i := int64(0); i < upkeepCount.Int64(); i++ {
		var tx *ethtypes.Transaction
		if tx, err = canceller.CancelUpkeep(k.buildTxOpts(ctx), big.NewInt(i)); err != nil {
			return fmt.Errorf("failed to cancel upkeep %d: %s", i, err)
		}
		k.waitTx(ctx, tx)

		if tx, err = canceller.WithdrawFunds(k.buildTxOpts(ctx), big.NewInt(i), k.fromAddr); err != nil {
			return fmt.Errorf("failed to withdraw upkeep %d: %s", i, err)
		}
		k.waitTx(ctx, tx)

		log.Println("Upkeep successfully canceled and refunded: ", i)
	}

	var tx *ethtypes.Transaction
	if tx, err = canceller.RecoverFunds(k.buildTxOpts(ctx)); err != nil {
		return fmt.Errorf("failed to recover funds: %s", err)
	}
	k.waitTx(ctx, tx)

	return nil
}

// createKeeperJob creates a keeper job in the chainlink node by the given address
func (k *Keeper) createKeeperJob(client cmd.HTTPClient, registryAddr, nodeAddr string) error {
	var err error
	if k.cfg.OCR2Keepers {
		err = k.createOCR2KeeperJob(client, registryAddr)
	} else {
		err = k.createLegacyKeeperJob(client, registryAddr, nodeAddr)
	}
	if err != nil {
		return err
	}

	log.Println("Keeper job has been successfully created in the Chainlink node with address: ", nodeAddr)

	return nil
}

// createLegacyKeeperJob creates a legacy keeper job in the chainlink node by the given address
func (k *Keeper) createLegacyKeeperJob(client cmd.HTTPClient, registryAddr, nodeAddr string) error {
	request, err := json.Marshal(web.CreateJobRequest{
		TOML: testspecs.GenerateKeeperSpec(testspecs.KeeperSpecParams{
			Name:            fmt.Sprintf("keeper job - registry %s", registryAddr),
			ContractAddress: registryAddr,
			FromAddress:     nodeAddr,
			EvmChainID:      int(k.cfg.ChainID),
		}).Toml(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %s", err)
	}

	resp, err := client.Post("/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("failed to create keeper job: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %s", err)
		}

		return fmt.Errorf("unable to create keeper job: '%v' [%d]", string(body), resp.StatusCode)
	}

	return nil
}

const ocr2keeperJobTemplate = `type = "offchainreporting2"
pluginType = "ocr2keeper"
relay = "evm"
name = "ocr2"
schemaVersion = 1
externalJobID = "%s"
maxTaskDuration = "1s"
contractID = "%s"
ocrKeyBundleID = "%s"
transmitterID = "%s"
p2pv2Bootstrappers = [
  "%s"
]

[relayConfig]
chainID = %d

[pluginConfig]
maxQueryLength = 2000
maxObservationLength = 2000
maxReportLength = 2000`

// createOCR2KeeperJob creates an ocr2keeper job in the chainlink node by the given address
func (k *Keeper) createOCR2KeeperJob(client cmd.HTTPClient, contractAddr string) error {
	// TODO: Fetch ocrKeyBundleID and transmitterID
	/*resp, err := client.Get("/v2/keys/eth")
	if err != nil {
		return err
	}

	respRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode, string(respRaw))

	resp, err = client.Get("/v2/keys/ocr2")
	if err != nil {
		return err
	}
	respRaw, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode, string(respRaw))
	return*/

	request, err := json.Marshal(web.CreateJobRequest{
		TOML: fmt.Sprintf(ocr2keeperJobTemplate,
			uuid.New().String(), // externalJobID: UUID
			contractAddr,        // contractID
			"aa53dde3867589b0df01a80429ec641b7a2b963fb1f3381a9207769aed7a0acd", // ocrKeyBundleID
			"",                      // transmitterID - node wallet address
			k.cfg.BootstrapNodeAddr, // bootstrap node key and address
			k.cfg.ChainID,           // chainID
		),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %s", err)
	}

	resp, err := client.Post("/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("failed to create ocr2keeper job: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %s", err)
		}

		return fmt.Errorf("unable to create ocr2keeper job: '%v' [%d]", string(body), resp.StatusCode)
	}

	return nil
}
