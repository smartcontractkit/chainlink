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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/web"
)

const (
	defaultChainlinkNodeImage = "smartcontract/chainlink:1.1.0"

	defaultChainlinkNodeLogin    = "test@smartcontract.com"
	defaultChainlinkNodePassword = "!PASsword000!"

	keeperJobTemplate = `
type            		 	= "keeper"
schemaVersion   		 	= 3
name            		 	= "CHAINCLI: Load testing"
contractAddress 		 	= "%s"
fromAddress     		 	= "%s"
evmChainID      		 	= %d
minIncomingConfirmations	= 1
externalJobID   		 	= "123e4567-e89b-12d3-a456-426655440002"


observationSource = """
encode_check_upkeep_tx   [type=ethabiencode
                          abi="checkUpkeep(uint256 id, address from)"
                          data="{\\"id\\":$(jobSpec.upkeepID),\\"from\\":$(jobSpec.fromAddress)}"]
check_upkeep_tx          [type=ethcall
                          failEarly=true
                          extractRevertReason=true
                          evmChainID="$(jobSpec.evmChainID)"
                          contract="$(jobSpec.contractAddress)"
                          gas="$(jobSpec.checkUpkeepGasLimit)"
                          gasPrice="$(jobSpec.gasPrice)"
                          gasTipCap="$(jobSpec.gasTipCap)"
                          gasFeeCap="$(jobSpec.gasFeeCap)"
                          data="$(encode_check_upkeep_tx)"]
decode_check_upkeep_tx   [type=ethabidecode
                          abi="bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
encode_perform_upkeep_tx [type=ethabiencode
                          abi="performUpkeep(uint256 id, bytes calldata performData)"
                          data="{\\"id\\": $(jobSpec.upkeepID),\\"performData\\":$(decode_check_upkeep_tx.performData)}"]
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          from="[$(jobSpec.fromAddress)]"
                          evmChainID="$(jobSpec.evmChainID)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\\"jobID\\":$(jobSpec.jobID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx
"""`
)

// LaunchAndTest launches keeper registry, chainlink nodes, upkeeps and start performing.
// 1. launch chainlink node using docker image
// 2. return node's wallet address
// 3. fund node by ETH and LINK
// 4. deploy keeper registry and upkeeps
// 5. create keeper job using keeper registry address
// 6. set keepers in the registry
// 7. wait until tests are done
func (k *Keeper) LaunchAndTest(ctx context.Context) {
	nodeAddr1, cleanup1 := k.launchChainlinkNode(ctx)
	defer cleanup1()

	if err := k.createKeepers(ctx, nodeAddr1); err != nil {
		log.Fatal("failed to create keepers: ", err)
	}

	log.Println("nodeAddr1", nodeAddr1)
	time.Sleep(time.Minute * 5)

	// Stop container
	//if err = dockerClient.ContainerStop(ctx, containerResp.ID, nil); err != nil {
	//	log.Fatal("Failed to stop container: ", err)
	//}
	//log.Println("Docker container successfully stopped")

	// client := cmd.NewAuthenticatedHTTPClient()
}

type cfg struct {
	nodeURL string
}

func (c cfg) ClientNodeURL() string    { return c.nodeURL }
func (c cfg) InsecureSkipVerify() bool { return true }

func (k *Keeper) createKeepers(ctx context.Context, nodeURL string) error {
	sr := sessions.SessionRequest{Email: defaultChainlinkNodeLogin, Password: defaultChainlinkNodePassword}
	store := &cmd.MemoryCookieStore{}
	tca := cmd.NewSessionCookieAuthenticator(cfg{nodeURL: nodeURL}, store, logger.NewLogger())
	if _, err := tca.Authenticate(sr); err != nil {
		return err
	}

	cl := cmd.NewAuthenticatedHTTPClient(cfg{nodeURL: nodeURL}, tca, sr)

	resp, err := cl.Get("/v2/keys/eth")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var keys cmd.EthKeyPresenters
	if err = jsonapi.Unmarshal(raw, &keys); err != nil {
		return err
	}

	keeperRegistryAddr := k.deployKeepers(ctx)

	request, err := json.Marshal(web.CreateJobRequest{
		TOML: fmt.Sprintf(keeperJobTemplate, keeperRegistryAddr.Hex(), keys[0].Address, k.cfg.ChainID),
	})
	if err != nil {
		return err
	}

	resp, err = cl.Post("/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("unable to create keeper job: '%v' [%d]", string(body), resp.StatusCode)
	}

	log.Println("Keeper job has been successfully created")

	return nil
}

func (k *Keeper) launchChainlinkNode(ctx context.Context) (string, func()) {
	// Create docker client to launch nodes
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Failed to create docker client from env: ", err)
	}

	// Make sure everything works well
	if _, err = dockerClient.Ping(ctx); err != nil {
		log.Fatal("Failed to ping docker server: ", err)
	}

	log.Println("Docker client successfully created")

	// Pull DB image if needed
	out, err := dockerClient.ImagePull(ctx, "postgres:13", types.ImagePullOptions{})
	if err != nil {
		log.Fatal("Failed to pull DB image: ", err)
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
		log.Fatal("Failed to create DB container: ", err)
	}
	log.Println("DB docker container successfully created")

	// Start container
	if err = dockerClient.ContainerStart(ctx, dbContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal("Failed to start DB container: ", err)
	}

	dbContainerInspect, err := dockerClient.ContainerInspect(ctx, dbContainerResp.ID)
	if err != nil {
		log.Fatal("Failed to inspect DB container: ", err)
	}

	// Pull node image if needed
	out, err = dockerClient.ImagePull(ctx, defaultChainlinkNodeImage, types.ImagePullOptions{})
	if err != nil {
		log.Fatal("Failed to pull node image: ", err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// Create temporary file with chainlink node login creds
	apiFile, err := ioutil.TempFile(os.TempDir(), "chainlink-node-api")
	if err != nil {
		log.Fatal("Failed to create api file: ", err)
	}
	apiFile.WriteString(defaultChainlinkNodeLogin)
	apiFile.WriteString("\n")
	apiFile.WriteString(defaultChainlinkNodePassword)

	// Create temporary file with chainlink node password
	passwordFile, err := ioutil.TempFile(os.TempDir(), "chainlink-node-password")
	if err != nil {
		log.Fatal("Failed to create password file: ", err)
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
			"LOG_LEVEL=info",
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
			"6688": {},
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
					HostPort: "6688",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		log.Fatal("Failed to create node container: ", err)
	}
	log.Println("Node docker container successfully created")

	// Start container
	if err = dockerClient.ContainerStart(ctx, nodeContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal("Failed to start node container: ", err)
	}

	time.Sleep(time.Second * 30)

	return fmt.Sprintf("http://localhost:%d", 6688), func() {
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
	}
}
