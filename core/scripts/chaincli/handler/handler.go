package handler

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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

	"github.com/smartcontractkit/chainlink/core/cmd"
	link "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/sessions"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

const (
	defaultChainlinkNodeLogin    = "notreal@fakeemail.ch"
	defaultChainlinkNodePassword = "fj293fbBnlQ!f9vNs~#"
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
		log.Fatal("failed to deal with ETH node", err)
	}
	nodeClient := ethclient.NewClient(rpcClient)

	// Parse private key
	d := new(big.Int).SetBytes(common.FromHex(cfg.PrivateKey))
	pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
	privateKey := &ecdsa.PrivateKey{
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
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

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

	gasPrice = bigmath.Add(gasPrice, bigmath.Div(gasPrice, 5)) // add 20%

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
		return fmt.Errorf("failed to sign tx: %s", err)
	}

	if err = k.client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send tx: %s", err)
	}

	k.waitTx(ctx, signedTx)

	return nil
}

func (h *baseHandler) waitDeployment(ctx context.Context, tx *ethtypes.Transaction) {
	if _, err := bind.WaitDeployed(ctx, h.client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err, " ", helpers.ExplorerLink(h.cfg.ChainID, tx.Hash()))
	}
}

func (h *baseHandler) waitTx(ctx context.Context, tx *ethtypes.Transaction) {
	receipt, err := bind.WaitMined(ctx, h.client, tx)
	if err != nil {
		log.Fatal("WaitDeployed failed: ", err)
	}

	if receipt.Status == ethtypes.ReceiptStatusFailed {
		log.Fatal("Transaction failed: ", helpers.ExplorerLink(h.cfg.ChainID, tx.Hash()))
	}
}

func (h *baseHandler) launchChainlinkNode(ctx context.Context, port int, containerName string, extraEnvVars ...string) (string, func(bool), error) {
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
	if _, _, err = dockerClient.ImageInspectWithRaw(ctx, h.cfg.PostgresDockerImage); err != nil {
		log.Println("Pulling Postgres docker image...")
		if out, err = dockerClient.ImagePull(ctx, h.cfg.PostgresDockerImage, types.ImagePullOptions{}); err != nil {
			return "", nil, fmt.Errorf("failed to pull Postgres image: %s", err)
		}
		out.Close()
		log.Println("Postgres docker image successfully pulled!")
	}

	// Create network config
	const networkName = "chaincli-local"
	existingNetworks, err := dockerClient.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to list networks: %s", err)
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
			return "", nil, fmt.Errorf("failed to create network: %s", err)
		}
	}

	// Create DB container
	postgresContainerName := fmt.Sprintf("%s-postgres", containerName)
	dbContainerResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: h.cfg.PostgresDockerImage,
		Cmd:   []string{"postgres", "-c", `max_connections=1000`},
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=development_password",
		},
		ExposedPorts: nat.PortSet{"5432": struct{}{}},
	}, nil, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {Aliases: []string{postgresContainerName}},
		},
	}, nil, postgresContainerName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create Postgres container: %s", err)
	}

	// Start container
	if err = dockerClient.ContainerStart(ctx, dbContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start DB container: %s", err)
	}
	log.Println("Postgres docker container successfully created and started: ", dbContainerResp.ID)

	time.Sleep(time.Second * 10)

	// Pull node image if needed
	if _, _, err = dockerClient.ImageInspectWithRaw(ctx, h.cfg.ChainlinkDockerImage); err != nil {
		log.Println("Pulling node docker image...")
		if out, err = dockerClient.ImagePull(ctx, h.cfg.ChainlinkDockerImage, types.ImagePullOptions{}); err != nil {
			return "", nil, fmt.Errorf("failed to pull node image: %s", err)
		}
		out.Close()
		log.Println("Node docker image successfully pulled!")
	}

	// Create temporary file with chainlink node login creds
	apiFile, passwordFile, fileCleanup, err := createCredsFiles()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create creds files: %s", err)
	}

	// Create container with mounted files
	portStr := fmt.Sprintf("%d", port)
	nodeContainerResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: h.cfg.ChainlinkDockerImage,
		Cmd:   []string{"local", "n", "-p", "/run/secrets/chainlink-node-password", "-a", "/run/secrets/chainlink-node-api"},
		Env: append([]string{
			"DATABASE_URL=postgresql://postgres:development_password@" + postgresContainerName + ":5432/postgres?sslmode=disable",
			"ETH_URL=" + h.cfg.NodeURL,
			fmt.Sprintf("ETH_CHAIN_ID=%d", h.cfg.ChainID),
			"LINK_CONTRACT_ADDRESS=" + h.cfg.LinkTokenAddr,
			"DATABASE_BACKUP_MODE=lite",
			"SKIP_DATABASE_PASSWORD_COMPLEXITY_CHECK=true",
			"LOG_LEVEL=debug",
			"CHAINLINK_TLS_PORT=0",
			"SECURE_COOKIES=false",
			"ALLOW_ORIGINS=*",
		}, extraEnvVars...),
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
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {Aliases: []string{containerName}},
		},
	}, nil, containerName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create node container: %s", err)
	}

	// Start container
	if err = dockerClient.ContainerStart(ctx, nodeContainerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start node container: %s", err)
	}

	addr := fmt.Sprintf("http://localhost:%s", portStr)
	log.Println("Node docker container successfully created and started: ", nodeContainerResp.ID, addr)

	if err = waitForNodeReady(addr); err != nil {
		log.Fatal(err, nodeContainerResp.ID)
	}
	log.Println("Node ready: ", nodeContainerResp.ID)

	return addr, func(writeLogs bool) {
		fileCleanup()

		if writeLogs {
			var rdr io.ReadCloser
			rdr, err := dockerClient.ContainerLogs(ctx, nodeContainerResp.ID, types.ContainerLogsOptions{
				ShowStderr: true,
				Timestamps: true,
			})
			if err != nil {
				rdr.Close()
				log.Fatal("Failed to collect logs from container: ", err)
			}

			stdErr, err := os.OpenFile(fmt.Sprintf("./%s-stderr.log", nodeContainerResp.ID), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				rdr.Close()
				stdErr.Close()
				log.Fatal("Failed to open file: ", err)
			}

			if _, err := stdcopy.StdCopy(io.Discard, stdErr, rdr); err != nil {
				rdr.Close()
				stdErr.Close()
				log.Fatal("Failed to write logs to file: ", err)
			}

			rdr.Close()
			stdErr.Close()
		}

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

func waitForNodeReady(addr string) error {
	client := &http.Client{}
	defer client.CloseIdleConnections()
	const timeout = 120
	startTime := time.Now().Unix()
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/health", addr), nil)
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
func authenticate(urlStr, email, password string, lggr logger.Logger) (cmd.HTTPClient, error) {
	remoteNodeURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	c := cmd.ClientOpts{RemoteNodeURL: *remoteNodeURL}
	sr := sessions.SessionRequest{Email: email, Password: password}
	store := &cmd.MemoryCookieStore{}

	tca := cmd.NewSessionCookieAuthenticator(c, store, lggr)
	if _, err = tca.Authenticate(sr); err != nil {
		log.Println("failed to authenticate: ", err)
		return nil, err
	}

	return cmd.NewAuthenticatedHTTPClient(lggr, c, tca, sr), nil
}

func nodeRequest(client cmd.HTTPClient, path string) ([]byte, error) {
	resp, err := client.Get(path)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error from client: %s", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read response body: %s", err)
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
func getNodeAddress(client cmd.HTTPClient) (string, error) {
	resp, err := nodeRequest(client, "/v2/keys/eth")
	if err != nil {
		return "", fmt.Errorf("failed to get ETH keys: %s", err)
	}

	var keys cmd.EthKeyPresenters
	if err = jsonapi.Unmarshal(resp, &keys); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	return keys[0].Address, nil
}

// getNodeOCR2Config returns chainlink node's OCR2 bundle key ID
func getNodeOCR2Config(client cmd.HTTPClient) (*cmd.OCR2KeyBundlePresenter, error) {
	resp, err := nodeRequest(client, "/v2/keys/ocr2")
	if err != nil {
		return nil, fmt.Errorf("failed to get OCR2 keys: %s", err)
	}

	var keys cmd.OCR2KeyBundlePresenters
	if err = jsonapi.Unmarshal(resp, &keys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
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
func getP2PKeyID(client cmd.HTTPClient) (string, error) {
	resp, err := nodeRequest(client, "/v2/keys/p2p")
	if err != nil {
		return "", fmt.Errorf("failed to get P2P keys: %s", err)
	}

	var keys cmd.P2PKeyPresenters
	if err = jsonapi.Unmarshal(resp, &keys); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	return keys[0].ID, nil
}

// createCredsFiles creates two temporary files with node creds: api and password.
func createCredsFiles() (string, string, func(), error) {
	// Create temporary file with chainlink node login creds
	apiFile, err := os.CreateTemp("", "chainlink-node-api")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create api file: %s", err)
	}
	_, _ = apiFile.WriteString(defaultChainlinkNodeLogin)
	_, _ = apiFile.WriteString("\n")
	_, _ = apiFile.WriteString(defaultChainlinkNodePassword)

	// Create temporary file with chainlink node password
	passwordFile, err := os.CreateTemp("", "chainlink-node-password")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create password file: %s", err)
	}
	_, _ = passwordFile.WriteString(defaultChainlinkNodePassword)

	return apiFile.Name(), passwordFile.Name(), func() {
		os.RemoveAll(apiFile.Name())
		os.RemoveAll(passwordFile.Name())
	}, nil
}
