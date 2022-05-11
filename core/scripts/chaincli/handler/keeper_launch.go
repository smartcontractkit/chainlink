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
	"net/url"
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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/cmd"
	keeperRegV1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/core/logger"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/web"
)

const (
	defaultChainlinkNodeImage = "smartcontract/chainlink:latest"
	defaultPOSTGRESImage      = "postgres:13"

	defaultChainlinkNodeLogin    = "notreal@fakeemail.ch"
	defaultChainlinkNodePassword = "twochains"
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
	var registry *keeperRegV1.KeeperRegistry
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

	// Approve keeper registry
	approveRegistryTx, err := k.linkToken.Approve(k.buildTxOpts(ctx), registryAddr, k.approveAmount)
	if err != nil {
		log.Fatal(registryAddr.Hex(), ": Approve failed - ", err)
	}
	k.waitTx(ctx, approveRegistryTx)
	log.Println(registryAddr.Hex(), ": KeeperRegistry approved - ", helpers.ExplorerLink(k.cfg.ChainID, approveRegistryTx.Hash()))

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registryAddr, registry, upkeepCount)

	lggr, closeLggr := logger.NewLogger()
	defer closeLggr()

	// Prepare keeper addresses and owners
	var keepers []common.Address
	var owners []common.Address
	for _, startedNode := range startedNodes {
		if startedNode.err != nil {
			log.Println("Failed to start node: ", err)
			continue
		}

		// Create authenticated client
		remoteNodeURL, err := url.Parse(startedNode.url)
		if err != nil {
			log.Fatal(err)
		}
		c := cmd.ClientOpts{RemoteNodeURL: *remoteNodeURL}
		sr := sessions.SessionRequest{Email: defaultChainlinkNodeLogin, Password: defaultChainlinkNodePassword}
		store := &cmd.MemoryCookieStore{}
		tca := cmd.NewSessionCookieAuthenticator(c, store, lggr)
		if _, err = tca.Authenticate(sr); err != nil {
			log.Println("failed to authenticate: ", err)
			continue
		}
		cl := cmd.NewAuthenticatedHTTPClient(lggr, c, tca, sr)

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
	log.Println("Set keepers...")
	setKeepersTx, err := registry.SetKeepers(k.buildTxOpts(ctx), keepers, owners)
	if err != nil {
		log.Fatal("SetKeepers failed: ", err)
	}
	k.waitTx(ctx, setKeepersTx)
	log.Println("Keepers registered:", setKeepersTx.Hash().Hex())

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
		log.Println("Canceling upkeeps...")
		if err = k.cancelAndWithdrawUpkeeps(ctx, registry); err != nil {
			log.Fatal("Failed to cancel upkeeps: ", err)
		}
		log.Println("Upkeeps successfully canceled")
	}
}

// cancelAndWithdrawUpkeeps cancels all upkeeps of the registry and withdraws funds
func (k *Keeper) cancelAndWithdrawUpkeeps(ctx context.Context, registryInstance *keeperRegV1.KeeperRegistry) error {
	count, err := registryInstance.GetUpkeepCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get upkeeps count: %s", err)
	}

	for i := int64(0); i < count.Int64(); i++ {
		var tx *ethtypes.Transaction
		if tx, err = registryInstance.CancelUpkeep(k.buildTxOpts(ctx), big.NewInt(i)); err != nil {
			return fmt.Errorf("failed to cancel upkeep %d: %s", i, err)
		}
		k.waitTx(ctx, tx)

		if tx, err = registryInstance.WithdrawFunds(k.buildTxOpts(ctx), big.NewInt(i), k.fromAddr); err != nil {
			return fmt.Errorf("failed to withdraw upkeep %d: %s", i, err)
		}
		k.waitTx(ctx, tx)

		log.Println("Upkeep successfully canceled and refunded: ", i)
	}

	var tx *ethtypes.Transaction
	if tx, err = registryInstance.RecoverFunds(k.buildTxOpts(ctx)); err != nil {
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
			Name:                     fmt.Sprintf("keeper job - registry %s", registryAddr),
			ContractAddress:          registryAddr,
			FromAddress:              nodeAddr,
			EvmChainID:               int(k.cfg.ChainID),
			MinIncomingConfirmations: 1,
			ObservationSource:        keeper.ExpectedObservationSource,
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
