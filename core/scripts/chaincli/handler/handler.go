package handler

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	link "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

const (
	defaultChainlinkNodeLogin    = "notreal@fakeemail.ch"
	defaultChainlinkNodePassword = "fj293fbBnlQ!f9vNs~#"
	ethKeysEndpoint              = "/v2/keys/eth"
	ocr2KeysEndpoint             = "/v2/keys/ocr2"
	p2pKeysEndpoint              = "/v2/keys/p2p"
	csaKeysEndpoint              = "/v2/keys/csa"
)

const (
	nodeTOML = `[Log]
JSONConsole = true
Level = 'debug'
[WebServer]
AllowOrigins = '*'
SecureCookies = false
SessionTimeout = '999h0m0s'
[WebServer.TLS]
HTTPSPort = 0
[Feature]
LogPoller = true
[OCR2]
Enabled = true

[Keeper]
TurnLookBack = 0
[[EVM]]
ChainID = '%d'
[[EVM.Nodes]]
Name = 'node-0'
WSURL = '%s'
HTTPURL = '%s'
`
	secretTOML = `
[Mercury.Credentials.cred1]
URL = '%s'
Username = '%s'
Password = '%s'
`
)

// baseHandler is the common handler with a common logic
type baseHandler struct {
	cfg *config.Config

	rpcClient     *rpc.Client
	client        *ethclient.Client
	privateKey    *ecdsa.PrivateKey
	linkToken     *link.LinkToken
	fromAddr      common.Address
	approveAmount *big.Int
}

// NewBaseHandler is the constructor of baseHandler
func NewBaseHandler(cfg *config.Config) *baseHandler {
	// Created a client by the given node address
	rpcClient, err := rpc.Dial(cfg.NodeURL)
	if err != nil {
		log.Fatal("failed to deal with ETH node: ", err)
	}
	nodeClient := ethclient.NewClient(rpcClient)

	// Parse private key
	var fromAddr common.Address
	var privateKey *ecdsa.PrivateKey
	if cfg.PrivateKey != "" {
		d := new(big.Int).SetBytes(common.FromHex(cfg.PrivateKey))
		pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
		privateKey = &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: crypto.S256(),
				X:     pkX,
				Y:     pkY,
			},
			D: d,
		}

		// Init from address
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}
		fromAddr = crypto.PubkeyToAddress(*publicKeyECDSA)
	} else {
		log.Println("WARNING: no PRIVATE_KEY set: cannot use commands that deploy contracts or send transactions")
	}

	// Create link token wrapper
	linkToken, err := link.NewLinkToken(common.HexToAddress(cfg.LinkTokenAddr), nodeClient)
	if err != nil {
		log.Fatal(err)
	}

	approveAmount := big.NewInt(0)
	approveAmount.SetString(cfg.ApproveAmount, 10)

	return &baseHandler{
		cfg:           cfg,
		client:        nodeClient,
		rpcClient:     rpcClient,
		privateKey:    privateKey,
		linkToken:     linkToken,
		fromAddr:      fromAddr,
		approveAmount: approveAmount,
	}
}

func (h *baseHandler) buildTxOpts(ctx context.Context) *bind.TransactOpts {
	nonce, err := h.client.PendingNonceAt(ctx, h.fromAddr)
	if err != nil {
		log.Fatal("PendingNonceAt failed: ", err)
	}

	gasPrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("SuggestGasPrice failed: ", err)
	}

	gasPrice = bigmath.Add(gasPrice, bigmath.Div(gasPrice, big.NewInt(5))) // add 20%

	auth, err := bind.NewKeyedTransactorWithChainID(h.privateKey, big.NewInt(h.cfg.ChainID))
	if err != nil {
		log.Fatal("NewKeyedTransactorWithChainID failed: ", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = h.cfg.GasLimit // in units
	auth.GasPrice = gasPrice

	return auth
}

// Send eth from prefunded account.
// Amount is number of wei.
func (k *Keeper) sendEth(ctx context.Context, to common.Address, amount *big.Int) error {
	txOpts := k.buildTxOpts(ctx)

	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    txOpts.Nonce.Uint64(),
		To:       &to,
		Value:    amount,
		Gas:      txOpts.GasLimit,
		GasPrice: txOpts.GasPrice,
		Data:     nil,
	})
	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(big.NewInt(k.cfg.ChainID)), k.privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign tx: %w", err)
	}

	if err = k.client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if err := k.waitTx(ctx, signedTx); err != nil {
		log.Fatalf("Send ETH failed, error is %s", err.Error())
	}
	log.Println("Send ETH successfully")

	return nil
}

func (h *baseHandler) waitDeployment(ctx context.Context, tx *ethtypes.Transaction) {
	if _, err := bind.WaitDeployed(ctx, h.client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err, " ", helpers.ExplorerLink(h.cfg.ChainID, tx.Hash()))
	}
}

func (h *baseHandler) waitTx(ctx context.Context, tx *ethtypes.Transaction) error {
	receipt, err := bind.WaitMined(ctx, h.client, tx)
	if err != nil {
		log.Println("WaitTx failed: ", err)
		return err
	}

	if receipt.Status == ethtypes.ReceiptStatusFailed {
		log.Println("Transaction failed: ", helpers.ExplorerLink(h.cfg.ChainID, tx.Hash()))
		return errors.New("Transaction failed")
	}

	return nil
}

func (h *baseHandler) launchChainlinkNode(ctx context.Context, port int, containerName string, extraTOML string, force bool) (string, func(bool), error) {
	// Create docker client to launch nodes
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return "", nil, fmt.Errorf("failed to create docker client from env: %w", err)
	}

	// Make sure everything works well
	if _, err = dockerClient.Ping(ctx); err != nil {
		return "", nil, fmt.Errorf("failed to ping docker server: %w", err)
	}

	// Pull DB image if needed
	var out io.ReadCloser
	if _, _, err = dockerClient.ImageInspectWithRaw(ctx, h.cfg.PostgresDockerImage); err != nil {
		log.Println("Pulling Postgres docker image...")
		if out, err = dockerClient.ImagePull(ctx, h.cfg.PostgresDockerImage, types.ImagePullOptions{}); err != nil {
			return "", nil, fmt.Errorf("failed to pull Postgres image: %w", err)
		}
		out.Close()
		log.Println("Postgres docker image successfully pulled!")
	}

	// Create network config
	const networkName = "chaincli-local"
	existingNetworks, err := dockerClient.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to list networks: %w", err)
	}

	var found bool
	for _, ntwrk := range existingNetworks {
		if ntwrk.Name == networkName {
			found = true
			break
		}
	}

	if !found {
		if _, err = dockerClient.NetworkCreate(ctx, networkName, types.NetworkCreate{}); err != nil {
			return "", nil, fmt.Errorf("failed to create network: %w", err)
		}
	}

	postgresContainerName := fmt.Sprintf("%s-postgres", containerName)

	// If force flag is on, we check and remove containers with the same name before creating new ones
	if force {
		if err = checkAndRemoveContainer(ctx, dockerClient, postgresContainerName); err != nil {
			return "", nil, fmt.Errorf("failed to remove container: %w", err)
		}
	}

	// Create DB container
	dbContainerResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: h.cfg.PostgresDockerImage,
		Cmd:   []string{"postgres", "-c", `max_connections=1000`},
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=verylongdatabasepassword",
		},
		ExposedPorts: nat.PortSet{"5432": struct{}{}},
	}, nil, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {Aliases: []string{postgresContainerName}},
		},
	}, nil, postgresContainerName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create Postgres container, use --force=true to force removing existing containers: %w", err)
	}

	// Start container
	if err = dockerClient.ContainerStart(ctx, dbContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start DB container: %w", err)
	}
	log.Println("Postgres docker container successfully created and started: ", dbContainerResp.ID)

	time.Sleep(time.Second * 10)

	// If force flag is on, we check and remove containers with the same name before creating new ones
	if force {
		if err = checkAndRemoveContainer(ctx, dockerClient, containerName); err != nil {
			return "", nil, fmt.Errorf("failed to remove container: %w", err)
		}
	}

	// Pull node image if needed
	if _, _, err = dockerClient.ImageInspectWithRaw(ctx, h.cfg.ChainlinkDockerImage); err != nil {
		log.Println("Pulling node docker image...")
		if out, err = dockerClient.ImagePull(ctx, h.cfg.ChainlinkDockerImage, types.ImagePullOptions{}); err != nil {
			return "", nil, fmt.Errorf("failed to pull node image: %w", err)
		}
		out.Close()
		log.Println("Node docker image successfully pulled!")
	}

	// Create temporary file with chainlink node login creds
	apiFile, passwordFile, fileCleanup, err := createCredsFiles()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create creds files: %w", err)
	}

	var baseTOML = fmt.Sprintf(nodeTOML, h.cfg.ChainID, h.cfg.NodeURL, h.cfg.NodeHttpURL)
	tomlFile, tomlFileCleanup, err := createTomlFile(baseTOML)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create toml file: %w", err)
	}
	var secretTOMLStr = fmt.Sprintf(secretTOML, h.cfg.DataStreamsURL, h.cfg.DataStreamsID, h.cfg.DataStreamsKey)
	secretFile, secretTOMLFileCleanup, err := createTomlFile(secretTOMLStr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create secret toml file: %w", err)
	}
	// Create container with mounted files
	portStr := fmt.Sprintf("%d", port)
	nodeContainerResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: h.cfg.ChainlinkDockerImage,
		Cmd:   []string{"-s", "/run/secrets/01-secret.toml", "-c", "/run/secrets/01-config.toml", "local", "n", "-a", "/run/secrets/chainlink-node-api"},
		Env: []string{
			"CL_CONFIG=" + extraTOML,
			"CL_PASSWORD_KEYSTORE=" + defaultChainlinkNodePassword,
			"CL_DATABASE_URL=postgresql://postgres:verylongdatabasepassword@" + postgresContainerName + ":5432/postgres?sslmode=disable",
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
			{
				Type:   mount.TypeBind,
				Source: tomlFile,
				Target: "/run/secrets/01-config.toml",
			},
			{
				Type:   mount.TypeBind,
				Source: secretFile,
				Target: "/run/secrets/01-secret.toml",
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
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {Aliases: []string{containerName}},
		},
	}, nil, containerName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create node container, use --force=true to force removing existing containers: %w", err)
	}

	// Start container
	if err = dockerClient.ContainerStart(ctx, nodeContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start node container: %w", err)
	}

	addr := fmt.Sprintf("http://localhost:%s", portStr)
	log.Println("Node docker container successfully created and started: ", nodeContainerResp.ID, addr)

	if err = waitForNodeReady(ctx, addr); err != nil {
		log.Fatal(err, nodeContainerResp.ID)
	}
	log.Println("Node ready: ", nodeContainerResp.ID)

	return addr, func(writeLogs bool) {
		fileCleanup()
		tomlFileCleanup()
		secretTOMLFileCleanup()

		if writeLogs {
			var rdr io.ReadCloser
			rdr, err2 := dockerClient.ContainerLogs(ctx, nodeContainerResp.ID, types.ContainerLogsOptions{
				ShowStderr: true,
				Timestamps: true,
			})
			if err2 != nil {
				rdr.Close()
				log.Fatal("Failed to collect logs from container: ", err2)
			}

			stdErr, err2 := os.OpenFile(fmt.Sprintf("./%s-stderr.log", nodeContainerResp.ID), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err2 != nil {
				rdr.Close()
				stdErr.Close()
				log.Fatal("Failed to open file: ", err2)
			}

			if _, err2 := stdcopy.StdCopy(io.Discard, stdErr, rdr); err2 != nil {
				rdr.Close()
				stdErr.Close()
				log.Fatal("Failed to write logs to file: ", err2)
			}

			rdr.Close()
			stdErr.Close()
		}

		if err2 := dockerClient.ContainerStop(ctx, nodeContainerResp.ID, container.StopOptions{}); err2 != nil {
			log.Fatal("Failed to stop node container: ", err2)
		}
		if err2 := dockerClient.ContainerRemove(ctx, nodeContainerResp.ID, types.ContainerRemoveOptions{}); err2 != nil {
			log.Fatal("Failed to remove node container: ", err2)
		}

		if err2 := dockerClient.ContainerStop(ctx, dbContainerResp.ID, container.StopOptions{}); err2 != nil {
			log.Fatal("Failed to stop DB container: ", err2)
		}
		if err2 := dockerClient.ContainerRemove(ctx, dbContainerResp.ID, types.ContainerRemoveOptions{}); err2 != nil {
			log.Fatal("Failed to remove DB container: ", err2)
		}
	}, nil
}

func checkAndRemoveContainer(ctx context.Context, dockerClient *client.Client, containerName string) error {
	opts := types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("name", "^/"+regexp.QuoteMeta(containerName)+"$")),
	}

	containers, err := dockerClient.ContainerList(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) > 1 {
		log.Fatal("more than two containers with the same name should not happen")
	} else if len(containers) > 0 {
		if err := dockerClient.ContainerRemove(ctx, containers[0].ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			return fmt.Errorf("failed to remove existing container: %w", err)
		}
		log.Println("successfully removed an existing container with name: ", containerName)
	}

	return nil
}

func waitForNodeReady(ctx context.Context, addr string) error {
	client := &http.Client{}
	defer client.CloseIdleConnections()
	const timeout = 120
	startTime := time.Now().Unix()
	for {
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/health", addr), nil)
		if err != nil {
			return err
		}
		req.Close = true
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		if time.Now().Unix()-startTime > int64(timeout*time.Second) {
			return fmt.Errorf("timed out waiting for node to start, waited %d seconds", timeout)
		}
		time.Sleep(time.Second * 5)
	}
}

// authenticate creates a http client with URL, email and password
func authenticate(ctx context.Context, urlStr, email, password string, lggr logger.Logger) (cmd.HTTPClient, error) {
	remoteNodeURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	c := cmd.ClientOpts{RemoteNodeURL: *remoteNodeURL}
	sr := sessions.SessionRequest{Email: email, Password: password}
	store := &cmd.MemoryCookieStore{}

	tca := cmd.NewSessionCookieAuthenticator(c, store, lggr)
	if _, err = tca.Authenticate(ctx, sr); err != nil {
		log.Println("failed to authenticate: ", err)
		return nil, err
	}

	return cmd.NewAuthenticatedHTTPClient(lggr, c, tca, sr), nil
}

func nodeRequest(ctx context.Context, client cmd.HTTPClient, path string) ([]byte, error) {
	resp, err := client.Get(ctx, path)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error from client: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read response body: %w", err)
	}

	type errorDetail struct {
		Detail string `json:"detail"`
	}

	type errorResp struct {
		Errors []errorDetail `json:"errors"`
	}

	var errs errorResp
	if err := json.Unmarshal(raw, &errs); err == nil && len(errs.Errors) > 0 {
		return []byte{}, fmt.Errorf("error returned from api: %s", errs.Errors[0].Detail)
	}

	return raw, nil
}

// getNodeAddress returns chainlink node's wallet address
func getNodeAddress(ctx context.Context, client cmd.HTTPClient) (string, error) {
	resp, err := nodeRequest(ctx, client, ethKeysEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get ETH keys: %w", err)
	}

	var keys cmd.EthKeyPresenters
	if err = jsonapi.Unmarshal(resp, &keys); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return keys[0].Address, nil
}

// getNodeOCR2Config returns chainlink node's OCR2 bundle key ID
func getNodeOCR2Config(ctx context.Context, client cmd.HTTPClient) (*cmd.OCR2KeyBundlePresenter, error) {
	resp, err := nodeRequest(ctx, client, ocr2KeysEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get OCR2 keys: %w", err)
	}

	var keys cmd.OCR2KeyBundlePresenters
	if err = jsonapi.Unmarshal(resp, &keys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	var evmKey cmd.OCR2KeyBundlePresenter
	for _, key := range keys {
		if key.ChainType == string(chaintype.EVM) {
			evmKey = key
			break
		}
	}

	return &evmKey, nil
}

// getP2PKeyID returns chainlink node's P2P key ID
func getP2PKeyID(ctx context.Context, client cmd.HTTPClient) (string, error) {
	resp, err := nodeRequest(ctx, client, p2pKeysEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get P2P keys: %w", err)
	}

	var keys cmd.P2PKeyPresenters
	if err = jsonapi.Unmarshal(resp, &keys); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return keys[0].ID, nil
}

// createCredsFiles creates two temporary files with node creds: api and password.
func createCredsFiles() (string, string, func(), error) {
	// Create temporary file with chainlink node login creds
	apiFile, err := os.CreateTemp("", "chainlink-node-api")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create api file: %w", err)
	}
	_, _ = apiFile.WriteString(defaultChainlinkNodeLogin)
	_, _ = apiFile.WriteString("\n")
	_, _ = apiFile.WriteString(defaultChainlinkNodePassword)

	// Create temporary file with chainlink node password
	passwordFile, err := os.CreateTemp("", "chainlink-node-password")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create password file: %w", err)
	}
	_, _ = passwordFile.WriteString(defaultChainlinkNodePassword)

	return apiFile.Name(), passwordFile.Name(), func() {
		os.RemoveAll(apiFile.Name())
		os.RemoveAll(passwordFile.Name())
	}, nil
}

// createTomlFile creates temporary file with TOML config
func createTomlFile(tomlString string) (string, func(), error) {
	// Create temporary file with chainlink node TOML config
	tomlFile, err := os.CreateTemp("", "chainlink-toml-config")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create toml file: %w", err)
	}
	_, _ = tomlFile.WriteString(tomlString)

	return tomlFile.Name(), func() {
		os.RemoveAll(tomlFile.Name())
	}, nil
}
