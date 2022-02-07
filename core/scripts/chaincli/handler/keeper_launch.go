package handler

import (
	"context"
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
)

const (
	defaultChainlinkNodeImage = "smartcontract/chainlink:1.1.0"

	defaultChainlinkNodeLogin    = "test@smartcontract.com"
	defaultChainlinkNodePassword = "!PASsword000!"
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

	log.Println("nodeAddr1", nodeAddr1)
	time.Sleep(time.Minute * 5)

	// Stop container
	//if err = dockerClient.ContainerStop(ctx, containerResp.ID, nil); err != nil {
	//	log.Fatal("Failed to stop container: ", err)
	//}
	//log.Println("Docker container successfully stopped")

	// client := cmd.NewAuthenticatedHTTPClient()
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

	nodeContainerInspect, err := dockerClient.ContainerInspect(ctx, nodeContainerResp.ID)
	if err != nil {
		log.Fatal("Failed to inspect node container: ", err)
	}

	return fmt.Sprintf("http://%s:6688", nodeContainerInspect.NetworkSettings.IPAddress), func() {
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
