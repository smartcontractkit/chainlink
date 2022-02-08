package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/cmd"
	keeper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/web"
)

const (
	defaultChainlinkNodeImage = "smartcontract/chainlink:1.1.0"

	defaultChainlinkNodeLogin    = "test@smartcontract.com"
	defaultChainlinkNodePassword = "!PASsword000!"
)

type cfg struct {
	nodeURL string
}

func (c cfg) ClientNodeURL() string    { return c.nodeURL }
func (c cfg) InsecureSkipVerify() bool { return true }

// LaunchAndTest launches keeper registry, chainlink nodes, upkeeps and start performing.
// 1. launch chainlink node using docker image
// 2. return node's wallet address
// 3. fund node by ETH and LINK
// 4. deploy keeper registry and upkeeps
// 5. create keeper job using keeper registry address
// 6. set keepers in the registry
// 7. wait until tests are done
func (k *Keeper) LaunchAndTest(ctx context.Context) {
	var registry *keeper.KeeperRegistry
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

	// Fund registry
	if err = k.sendEth(ctx, registryAddr, 5); err != nil {
		log.Fatal(registryAddr.Hex(), ": Fund failed - ", err)
	}

	// Run chainlink nodes and create jobs
	nodesCount := 2
	nodeAddrs := make([]common.Address, nodesCount)
	cleanups := make([]func(), nodesCount)
	errs := make([]error, nodesCount)
	var wg sync.WaitGroup
	for i := 0; i < nodesCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Run chainlink node
			var err error
			var nodeURL string
			if nodeURL, cleanups[i], err = k.launchChainlinkNode(ctx, 6688+i); err != nil {
				errs[i] = fmt.Errorf("failed to launch chainlink node: %s", err)
				return
			}

			// Create authenticated client
			c := cfg{nodeURL: nodeURL}
			sr := sessions.SessionRequest{Email: defaultChainlinkNodeLogin, Password: defaultChainlinkNodePassword}
			store := &cmd.MemoryCookieStore{}
			tca := cmd.NewSessionCookieAuthenticator(c, store, logger.NewLogger())
			if _, err = tca.Authenticate(sr); err != nil {
				errs[i] = fmt.Errorf("failed to authenticate: %s", err)
				return
			}
			cl := cmd.NewAuthenticatedHTTPClient(c, tca, sr)

			// Get node's wallet address
			nodeAddr, err := k.getNodeAddress(cl)
			if err != nil {
				errs[i] = fmt.Errorf("failed to get node addr: %s", err)
				return
			}
			nodeAddrs[i] = common.HexToAddress(nodeAddr)

			// Create keepers
			if err = k.createKeeperJob(cl, registryAddr.Hex(), nodeAddr); err != nil {
				errs[i] = fmt.Errorf("failed to create keeper job: %s", err)
				return
			}
		}(i)
	}
	wg.Wait()

	// Make sure there were no errors
	for _, err = range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Cleanup resources
	defer func() {
		for _, cleanup := range cleanups {
			go cleanup()
		}
	}()

	// Deploy Upkeeps
	k.deployUpkeeps(ctx, registryAddr, registry, upkeepCount)

	// Prepare keeper addresses and owners
	var owners []common.Address
	for _, nodeAddr := range nodeAddrs {
		// Fund node
		if err = k.sendEth(ctx, nodeAddr, 5); err != nil {
			log.Fatal("Failed to fund chainlink node: ", err)
		}

		owners = append(owners, k.fromAddr)
	}

	// Set Keepers
	log.Println("Set keepers...")
	setKeepersTx, err := registry.SetKeepers(k.buildTxOpts(ctx), nodeAddrs, owners)
	if err != nil {
		log.Fatal("SetKeepers failed: ", err)
	}
	k.waitTx(ctx, setKeepersTx)
	log.Println("Keepers registered:", setKeepersTx.Hash().Hex())

	time.Sleep(time.Minute * 5)

	// Stop container
	//if err = dockerClient.ContainerStop(ctx, containerResp.ID, nil); err != nil {
	//	log.Fatal("Failed to stop container: ", err)
	//}
	//log.Println("Docker container successfully stopped")

	// client := cmd.NewAuthenticatedHTTPClient()
}

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

func (k *Keeper) createKeeperJob(client cmd.HTTPClient, contractAddr, nodeAddr string) error {
	request, err := json.Marshal(web.CreateJobRequest{
		TOML: testspecs.GenerateKeeperSpec(testspecs.KeeperSpecParams{
			ContractAddress:          contractAddr,
			FromAddress:              nodeAddr,
			EvmChainID:               int(k.cfg.ChainID),
			MinIncomingConfirmations: 1,
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

	log.Println("Keeper job has been successfully created")

	return nil
}

func (k *Keeper) launchChainlinkNode(ctx context.Context, port int) (string, func(), error) {
	portStr := fmt.Sprintf("%d", port)

	// Create docker client to launch nodes
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return "", nil, fmt.Errorf("failed to create docker client from env: %s", err)
	}

	// Make sure everything works well
	if _, err = dockerClient.Ping(ctx); err != nil {
		return "", nil, fmt.Errorf("failed to ping docker server: %s", err)
	}

	log.Println("Docker client successfully created")

	// Pull DB image if needed
	out, err := dockerClient.ImagePull(ctx, "postgres:13", types.ImagePullOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to pull DB image: %s", err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// Create DB container
	dbContainerResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Hostname: "postgres1",
		Image:    "postgres:13",
		Cmd:      []string{"postgres", "-c", `max_connections=1000`},
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=development_password",
		},
		ExposedPorts: nat.PortSet{"5432": struct{}{}},
	}, nil, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create DB container: %s", err)
	}
	log.Println("DB docker container successfully created")

	// Start container
	if err = dockerClient.ContainerStart(ctx, dbContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start DB container: %s", err)
	}

	dbContainerInspect, err := dockerClient.ContainerInspect(ctx, dbContainerResp.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to inspect DB container: %s", err)
	}

	// Pull node image if needed
	out, err = dockerClient.ImagePull(ctx, defaultChainlinkNodeImage, types.ImagePullOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to pull node image: %s", err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// Create temporary file with chainlink node login creds
	apiFile, err := ioutil.TempFile(os.TempDir(), "chainlink-node-api")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create api file: %s", err)
	}
	apiFile.WriteString(defaultChainlinkNodeLogin)
	apiFile.WriteString("\n")
	apiFile.WriteString(defaultChainlinkNodePassword)

	// Create temporary file with chainlink node password
	passwordFile, err := ioutil.TempFile(os.TempDir(), "chainlink-node-password")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create password file: %s", err)
	}
	passwordFile.WriteString(defaultChainlinkNodePassword)

	// Create container with mounted files
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
			"FEATURE_UI_CSA_KEYS=true",
			"FEATURE_UI_FEEDS_MANAGER=true",
		},
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(portStr): {},
		},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: apiFile.Name(),
				Target: "/run/secrets/chainlink-node-api",
			},
			{
				Type:   mount.TypeBind,
				Source: passwordFile.Name(),
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
	log.Println("Node docker container successfully created")

	// Start container
	if err = dockerClient.ContainerStart(ctx, nodeContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start node container: %s", err)
	}

	time.Sleep(time.Second * 10)

	return fmt.Sprintf("http://localhost:%s", portStr), func() {
		os.Remove(apiFile.Name())
		os.Remove(passwordFile.Name())

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
