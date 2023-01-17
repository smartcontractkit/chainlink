package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/core/cmd"
	registry12 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	registry20 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
)

type startedNodeData struct {
	url     string
	cleanup func(bool)
}

// LaunchAndTest launches keeper registry, chainlink nodes, upkeeps and start performing.
// 1. launch chainlink node using docker image
// 2. get keeper registry instance, deploy if needed
// 3. deploy upkeeps
// 4. create keeper jobs
// 5. fund nodes if needed
// 6. set keepers in the registry
// 7. withdraw funds after tests are done -> TODO: wait until tests are done instead of cancel manually
func (k *Keeper) LaunchAndTest(ctx context.Context, withdraw bool, printLogs bool) {
	lggr, closeLggr := logger.NewLogger()
	logger.Sugared(lggr).ErrorIfFn(closeLggr, "Failed to close logger")

	var extraEvars []string
	if k.cfg.OCR2Keepers {
		extraEvars = []string{
			"FEATURE_OFFCHAIN_REPORTING2=true",
			"FEATURE_LOG_POLLER=true",
			"P2P_NETWORKING_STACK=V2",
			"CHAINLINK_TLS_PORT=0",
			"P2PV2_LISTEN_ADDRESSES=0.0.0.0:8000",
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
			if startedNodes[i].url, startedNodes[i].cleanup, err = k.launchChainlinkNode(ctx, 6688+i, fmt.Sprintf("keeper-%d", i), extraEvars...); err != nil {
				log.Fatal("Failed to start node: ", err)
			}
		}(i)
	}
	wg.Wait()

	// Deploy keeper registry or get an existing one
	upkeepCount, registryAddr, deployer := k.prepareRegistry(ctx)

	// Approve keeper registry
	k.approveFunds(ctx, registryAddr)

	// Prepare keeper addresses and owners
	var keepers []common.Address
	var owners []common.Address
	var cls []cmd.HTTPClient
	for i, startedNode := range startedNodes {
		// Create authenticated client
		var cl cmd.HTTPClient
		var err error
		cl, err = authenticate(startedNode.url, defaultChainlinkNodeLogin, defaultChainlinkNodePassword, lggr)
		if err != nil {
			log.Fatal("Authentication failed, ", err)
		}

		var nodeAddrHex string

		if len(k.cfg.KeeperKeys) > 0 {
			// import key if exists
			var err error
			nodeAddrHex, err = k.addKeyToKeeper(cl, k.cfg.KeeperKeys[i])
			if err != nil {
				log.Fatal("could not add key to keeper", err)
			}
		} else {
			// get node's default wallet address
			var err error
			nodeAddrHex, err = getNodeAddress(cl)
			if err != nil {
				log.Println("Failed to get node addr: ", err)
				continue
			}
		}

		nodeAddr := common.HexToAddress(nodeAddrHex)

		// Create keepers
		if err = k.createKeeperJob(cl, registryAddr.Hex(), nodeAddr.Hex()); err != nil {
			log.Println("Failed to create keeper job: ", err)
			continue
		}

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

		cls = append(cls, cl)
		keepers = append(keepers, nodeAddr)
		owners = append(owners, k.fromAddr)
	}

	if len(keepers) == 0 {
		log.Fatal("no keepers available")
	}

	// Set Keepers
	k.setKeepers(ctx, cls, deployer, keepers, owners)

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registryAddr, deployer, upkeepCount)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
	log.Println("Stopping...")

	// Cleanup resources
	for _, startedNode := range startedNodes {
		if startedNode.cleanup != nil {
			startedNode.cleanup(printLogs)
		}
	}

	// Cancel upkeeps and withdraw funds
	if withdraw {
		log.Println("Canceling upkeeps...")
		switch k.cfg.RegistryVersion {
		case keeper.RegistryVersion_1_1:
			if err := k.cancelAndWithdrawUpkeeps(ctx, big.NewInt(upkeepCount), deployer); err != nil {
				log.Fatal("Failed to cancel upkeeps: ", err)
			}
		case keeper.RegistryVersion_1_2:
			registry, err := registry12.NewKeeperRegistry(
				registryAddr,
				k.client,
			)
			if err != nil {
				log.Fatal("Registry failed: ", err)
			}

			activeUpkeepIds := k.getActiveUpkeepIds(ctx, registry, big.NewInt(0), big.NewInt(0))
			if err := k.cancelAndWithdrawActiveUpkeeps(ctx, activeUpkeepIds, deployer); err != nil {
				log.Fatal("Failed to cancel upkeeps: ", err)
			}
		case keeper.RegistryVersion_2_0:
			registry, err := registry20.NewKeeperRegistry(
				registryAddr,
				k.client,
			)
			if err != nil {
				log.Fatal("Registry failed: ", err)
			}

			activeUpkeepIds := k.getActiveUpkeepIds(ctx, registry, big.NewInt(0), big.NewInt(0))
			if err := k.cancelAndWithdrawActiveUpkeeps(ctx, activeUpkeepIds, deployer); err != nil {
				log.Fatal("Failed to cancel upkeeps: ", err)
			}
		default:
			panic("unexpected registry address")
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
		err = k.createOCR2KeeperJob(client, registryAddr, nodeAddr)
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %s", err)
		}

		return fmt.Errorf("unable to create keeper job: '%v' [%d]", string(body), resp.StatusCode)
	}

	return nil
}

const ocr2keeperJobTemplate = `type = "offchainreporting2"
pluginType = "ocr2automation"
relay = "evm"
name = "ocr2-automation"
forwardingAllowed = false
schemaVersion = 1
contractID = "%s"
ocrKeyBundleID = "%s"
transmitterID = "%s"
p2pv2Bootstrappers = [
  "%s"
]

[relayConfig]
chainID = %d

[pluginConfig]
maxServiceWorkers = 100`

// createOCR2KeeperJob creates an ocr2keeper job in the chainlink node by the given address
func (k *Keeper) createOCR2KeeperJob(client cmd.HTTPClient, contractAddr, nodeAddr string) error {
	ocr2KeyConfig, err := getNodeOCR2Config(client)
	if err != nil {
		return fmt.Errorf("failed to get node OCR2 key bundle ID: %s", err)
	}

	request, err := json.Marshal(web.CreateJobRequest{
		TOML: fmt.Sprintf(ocr2keeperJobTemplate,
			contractAddr,            // contractID
			ocr2KeyConfig.ID,        // ocrKeyBundleID
			nodeAddr,                // transmitterID - node wallet address
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %s", err)
		}

		return fmt.Errorf("unable to create ocr2keeper job: '%v' [%d]", string(body), resp.StatusCode)
	}

	return nil
}

// addKeyToKeeper imports the provided ETH sending key to the keeper
func (k *Keeper) addKeyToKeeper(client cmd.HTTPClient, privKeyHex string) (string, error) {
	privkey, err := crypto.HexToECDSA(utils.RemoveHexPrefix(privKeyHex))
	if err != nil {
		log.Fatalf("Failed to decode priv key %s: %v", privKeyHex, err)
	}
	address := crypto.PubkeyToAddress(privkey.PublicKey).Hex()
	log.Printf("importing keeper key %s", address)
	keyJSON, err := ethkey.FromPrivateKey(privkey).ToEncryptedJSON(defaultChainlinkNodePassword, utils.FastScryptParams)
	if err != nil {
		log.Fatalf("Failed to encrypt piv key %s: %v", privKeyHex, err)
	}
	importUrl := url.URL{
		Path: "/v2/keys/evm/import",
	}
	query := importUrl.Query()

	query.Set("oldpassword", defaultChainlinkNodePassword)
	query.Set("evmChainID", fmt.Sprint(k.cfg.ChainID))

	importUrl.RawQuery = query.Encode()
	resp, err := client.Post(importUrl.String(), bytes.NewReader(keyJSON))
	if err != nil {
		log.Fatalf("Failed to import priv key %s: %v", privKeyHex, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read error response body: %s", err)
		}

		return "", fmt.Errorf("unable to create ocr2keeper job: '%v' [%d]", string(body), resp.StatusCode)
	}

	return address, nil
}
