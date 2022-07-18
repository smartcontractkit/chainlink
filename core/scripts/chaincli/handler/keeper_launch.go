package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/cmd"
	registry12 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/web"
)

const (
	defaultChainlinkNodeImage = "smartcontract/chainlink:1.5.1-root"
	defaultPOSTGRESImage      = "postgres:latest"

	defaultChainlinkNodeLogin    = "notreal@fakeemail.ch"
	defaultChainlinkNodePassword = "fj293fbBnlQ!f9vNs"
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
			if startedNodes[i].url, startedNodes[i].cleanup, err = k.launchChainlinkNode(ctx, 6688+i); err != nil {
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

	lggr, closeLggr := logger.NewLogger()
	defer closeLggr()

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
		cl, err = k.authenticate(startedNode.url, defaultChainlinkNodeLogin, defaultChainlinkNodePassword, lggr)
		if err != nil {
			log.Fatal("Authentication failed, ", err)
		}

		// Get node's wallet address
		var nodeAddrHex string
		if nodeAddrHex, err = k.getNodeAddress(cl); err != nil {
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

// getNodeAddress returns chainlink node's wallet address
func (k *Keeper) getNodeAddress(client cmd.HTTPClient) (string, error) {
	resp, err := client.Get("/v2/keys/eth")
	if err != nil {
		return "", fmt.Errorf("failed to get ETH keys: %s", err)
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	var keys cmd.EthKeyPresenters
	if err = jsonapi.Unmarshal(raw, &keys); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	return keys[0].Address, nil
}

// createKeeperJob creates a keeper job in the chainlink node by the given address
func (k *Keeper) createKeeperJob(client cmd.HTTPClient, registryAddr, nodeAddr string) error {
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
	log.Println("Keeper job has been successfully created in the Chainlink node with address: ", nodeAddr)
	return nil
}

func (k *Keeper) launchChainlinkNode(ctx context.Context, port int) (string, func(), error) {
	// Create docker client to launch nodes
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return "", nil, fmt.Errorf("failed to create docker client from env: %s", err)
	}

	// Make sure everything works well
	if _, err = dockerClient.Ping(ctx); err != nil {
		return "", nil, fmt.Errorf("failed to ping docker server: %s", err)
	}

	// Pull DB image if needed
	var out io.ReadCloser
	if _, _, err = dockerClient.ImageInspectWithRaw(ctx, defaultPOSTGRESImage); err != nil {
		log.Println("Pulling Postgres docker image...")
		if out, err = dockerClient.ImagePull(ctx, defaultPOSTGRESImage, types.ImagePullOptions{}); err != nil {
			return "", nil, fmt.Errorf("failed to pull Postgres image: %s", err)
		}
		out.Close()
		log.Println("Postgres docker image successfully pulled!")
	}

	// Create DB container
	dbContainerResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: defaultPOSTGRESImage,
		Cmd:   []string{"postgres", "-c", `max_connections=1000`},
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=development_password",
		},
		ExposedPorts: nat.PortSet{"5432": struct{}{}},
	}, nil, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create Postgres container: %s", err)
	}

	// Start container
	if err = dockerClient.ContainerStart(ctx, dbContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start DB container: %s", err)
	}
	log.Println("Postgres docker container successfully created and started: ", dbContainerResp.ID)

	dbContainerInspect, err := dockerClient.ContainerInspect(ctx, dbContainerResp.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to inspect Postgres container: %s", err)
	}

	time.Sleep(time.Second * 10)

	// Pull node image if needed
	if _, _, err = dockerClient.ImageInspectWithRaw(ctx, defaultChainlinkNodeImage); err != nil {
		log.Println("Pulling node docker image...")
		if out, err = dockerClient.ImagePull(ctx, defaultChainlinkNodeImage, types.ImagePullOptions{}); err != nil {
			return "", nil, fmt.Errorf("failed to pull node image: %s", err)
		}
		out.Close()
		log.Println("Node docker image successfully pulled!")
	}

	// Create temporary file with chainlink node login creds
	apiFile, passwordFile, fileCleanup, err := k.createCredsFiles()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create creds files: %s", err)
	}

	// Create container with mounted files
	portStr := fmt.Sprintf("%d", port)
	nodeContainerResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: defaultChainlinkNodeImage,
		Cmd:   []string{"local", "n", "-p", "/run/secrets/chainlink-node-password", "-a", "/run/secrets/chainlink-node-api"},
		Env: []string{
			"DATABASE_URL=postgresql://postgres:development_password@" + dbContainerInspect.NetworkSettings.IPAddress + ":5432/postgres?sslmode=disable",
			"ETH_URL=" + k.cfg.NodeURL,
			fmt.Sprintf("ETH_CHAIN_ID=%d", k.cfg.ChainID),
			"LINK_CONTRACT_ADDRESS=" + k.cfg.LinkTokenAddr,
			"DATABASE_BACKUP_MODE=lite",
			"ROOT=/chainlink",
			"LOG_LEVEL=debug",
			"MIN_OUTGOING_CONFIRMATIONS=2",
			"CHAINLINK_TLS_PORT=0",
			"SECURE_COOKIES=false",
			"GAS_ESTIMATOR_MODE=BlockHistory",
			"ALLOW_ORIGINS=*",
			"DATABASE_TIMEOUT=0",
			"KEEPER_CHECK_UPKEEP_GAS_PRICE_FEATURE_ENABLED=true",
		},
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(portStr): {},
		},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: apiFile,
				Target: "/run/secrets/chainlink-node-api",
			},
			{
				Type:   mount.TypeBind,
				Source: passwordFile,
				Target: "/run/secrets/chainlink-node-password",
			},
		},
		PortBindings: nat.PortMap{
			"6688/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: portStr,
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create node container: %s", err)
	}

	// Start container
	if err = dockerClient.ContainerStart(ctx, nodeContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start node container: %s", err)
	}

	addr := fmt.Sprintf("http://localhost:%s", portStr)
	log.Println("Node docker container successfully created and started: ", nodeContainerResp.ID, addr)

	time.Sleep(time.Second * 20)

	return addr, func() {
		fileCleanup()

		if err = dockerClient.ContainerStop(ctx, nodeContainerResp.ID, nil); err != nil {
			log.Fatal("Failed to stop node container: ", err)
		}
		if err = dockerClient.ContainerRemove(ctx, nodeContainerResp.ID, types.ContainerRemoveOptions{}); err != nil {
			log.Fatal("Failed to remove node container: ", err)
		}

		if err = dockerClient.ContainerStop(ctx, dbContainerResp.ID, nil); err != nil {
			log.Fatal("Failed to stop DB container: ", err)
		}
		if err = dockerClient.ContainerRemove(ctx, dbContainerResp.ID, types.ContainerRemoveOptions{}); err != nil {
			log.Fatal("Failed to remove DB container: ", err)
		}
	}, nil
}

// createCredsFiles creates two temporary files with node creds: api and password.
func (k *Keeper) createCredsFiles() (string, string, func(), error) {
	// Create temporary file with chainlink node login creds
	apiFile, err := ioutil.TempFile(os.TempDir(), "chainlink-node-api")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create api file: %s", err)
	}
	_, _ = apiFile.WriteString(defaultChainlinkNodeLogin)
	_, _ = apiFile.WriteString("\n")
	_, _ = apiFile.WriteString(defaultChainlinkNodePassword)

	// Create temporary file with chainlink node password
	passwordFile, err := ioutil.TempFile(os.TempDir(), "chainlink-node-password")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create password file: %s", err)
	}
	_, _ = passwordFile.WriteString(defaultChainlinkNodePassword)

	return apiFile.Name(), passwordFile.Name(), func() {
		os.Remove(apiFile.Name())
		os.Remove(passwordFile.Name())
	}, nil
}
