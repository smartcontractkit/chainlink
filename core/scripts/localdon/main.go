package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

const (
	anvilURL   = "http://localhost:8545"
	chainID    = 31337
	deployerPK = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // Anvil account 0

	apiEmail    = "local-test@chainlink.test"
	apiPassword = "localdon-testing-only-not-a-real-password"
)

var nodeURLs = []string{
	"http://localhost:6688",
	"http://localhost:6689",
	"http://localhost:6690",
	"http://localhost:6691",
}

func main() {
	ctx := context.Background()

	// Connect to Anvil
	client, err := ethclient.Dial(anvilURL)
	must(err, "dial anvil")
	log.Println("Connected to Anvil")

	deployerKey, err := crypto.HexToECDSA(deployerPK)
	must(err, "parse deployer key")
	deployer, err := bind.NewKeyedTransactorWithChainID(deployerKey, big.NewInt(chainID))
	must(err, "create deployer transactor")

	// --- Deploy contracts ---

	// LINK token
	linkAddr, linkTx, linkContract, err := link_token_interface.DeployLinkToken(deployer, client)
	must(err, "deploy LINK")
	waitMined(ctx, client, linkTx.Hash(), "LINK token")
	log.Printf("LINK token: %s", linkAddr)

	// OCR2 Aggregator (zero-address access controllers are fine for local dev)
	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))

	ocrAddr, ocrTx, ocrContract, err := ocr2aggregator.DeployOCR2Aggregator(
		deployer, client,
		linkAddr,
		minAnswer, maxAnswer,
		common.Address{}, // billingAccessController
		common.Address{}, // requesterAccessController
		9, "TEST",
	)
	must(err, "deploy OCR2Aggregator")
	waitMined(ctx, client, ocrTx.Hash(), "OCR2Aggregator")
	log.Printf("OCR2Aggregator: %s", ocrAddr)

	// Fund aggregator with LINK
	tx, err := linkContract.Transfer(deployer, ocrAddr, big.NewInt(1000))
	must(err, "fund aggregator with LINK")
	waitMined(ctx, client, tx.Hash(), "LINK transfer")

	// --- Gather node keys ---

	type nodeInfo struct {
		ethAddr       common.Address
		ocrKeyID      string
		onchainPubKey []byte
		offchainPub   [32]byte
		configPub     [32]byte
		peerID        string
	}

	nodes := make([]nodeInfo, len(nodeURLs))
	for i, url := range nodeURLs {
		nc := newNodeClient(url)
		log.Printf("Node %d (%s): connected", i+1, url)

		// ETH key (auto-created on startup)
		nodes[i].ethAddr = nc.getETHAddress()
		log.Printf("  ETH address: %s", nodes[i].ethAddr)

		// OCR2 key (must create)
		nodes[i].ocrKeyID, nodes[i].onchainPubKey, nodes[i].offchainPub, nodes[i].configPub = nc.createOCR2Key()
		log.Printf("  OCR2 key:    %s", nodes[i].ocrKeyID)

		// P2P key (auto-created on startup)
		nodes[i].peerID = nc.getP2PID()
		log.Printf("  P2P peer ID: %s", nodes[i].peerID)
	}

	// --- Fund node ETH addresses ---

	for i, n := range nodes {
		fundTx, err := fundAddress(ctx, client, deployerKey, n.ethAddr, big.NewInt(1e18)) // 1 ETH
		must(err, fmt.Sprintf("fund node %d", i+1))
		waitMined(ctx, client, fundTx, fmt.Sprintf("fund node %d", i+1))
		log.Printf("Funded node %d: %s", i+1, n.ethAddr)
	}

	// --- SetPayees ---

	transmitters := make([]common.Address, len(nodes))
	for i, n := range nodes {
		transmitters[i] = n.ethAddr
	}

	tx, err = ocrContract.SetPayees(deployer, transmitters, transmitters)
	must(err, "SetPayees")
	waitMined(ctx, client, tx.Hash(), "SetPayees")
	log.Println("SetPayees done")

	// --- SetConfig ---

	oracles := make([]confighelper2.OracleIdentityExtra, len(nodes))
	for i, n := range nodes {
		// PeerID from API includes "p2p_" prefix, but OCR config expects bare ID
		peerID := strings.TrimPrefix(n.peerID, "p2p_")
		oracles[i] = confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  n.onchainPubKey,
				TransmitAccount:   ocrtypes2.Account(n.ethAddr.Hex()),
				OffchainPublicKey: n.offchainPub,
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: n.configPub,
		}
	}

	signers, effectiveTransmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err :=
		confighelper2.ContractSetConfigArgsForEthereumIntegrationTest(oracles, 1, 1000000000/100)
	must(err, "ContractSetConfigArgsForEthereumIntegrationTest")

	// Override onchainConfig with median-specific min/max
	onchainConfig, err = median.StandardOnchainConfigCodec{}.Encode(ctx, median.OnchainConfig{
		Min: minAnswer,
		Max: maxAnswer,
	})
	must(err, "encode onchain config")

	tx, err = ocrContract.SetConfig(deployer, signers, effectiveTransmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig)
	must(err, "SetConfig")
	waitMined(ctx, client, tx.Hash(), "SetConfig")
	log.Println("SetConfig done")

	// Get block number for fromBlock in job specs
	blockNum, err := client.BlockNumber(ctx)
	must(err, "get block number")

	// --- Create oracle jobs ---

	bootstrapPeerID := nodes[0].peerID
	for i, url := range nodeURLs {
		nc := newNodeClient(url)
		spec := oracleJobSpec(ocrAddr, nodes[i].ocrKeyID, nodes[i].ethAddr, int64(blockNum), bootstrapPeerID)
		nc.createJob(spec)
		log.Printf("Created oracle job on node %d", i+1)
	}

	log.Println("Setup complete! The DON should start producing rounds shortly.")
	log.Printf("OCR2Aggregator address: %s", ocrAddr)
	log.Printf("Monitor with: cast call %s 'latestAnswer()(int256)' --rpc-url %s", ocrAddr, anvilURL)
}

func oracleJobSpec(contractAddr common.Address, ocrKeyBundleID string, transmitterAddr common.Address, fromBlock int64, bootstrapPeerID string) string {
	// PeerID for bootstrapper must be bare (no "p2p_" prefix)
	barePeerID := strings.TrimPrefix(bootstrapPeerID, "p2p_")
	return fmt.Sprintf(`
type               = "offchainreporting2"
relay              = "evm"
schemaVersion      = 1
pluginType         = "median"
name               = "localdon-median"
contractID         = "%s"
ocrKeyBundleID     = "%s"
transmitterID      = "%s"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
p2pv2Bootstrappers = ["%s@localdon-node-1:6690"]

observationSource = """
    price [type=memo value="42000000000"];
    price_multiply [type=multiply times=1];
    price -> price_multiply;
"""

[relayConfig]
chainID = "%d"
fromBlock = %d

[pluginConfig]
juelsPerFeeCoinSource = """
    price [type=memo value="1000000000000000000"];
    price_multiply [type=multiply times=1];
    price -> price_multiply;
"""
gasPriceSubunitsSource = """
    price [type=memo value="1000000000"];
    price_multiply [type=multiply times=1];
    price -> price_multiply;
"""

[pluginConfig.juelsPerFeeCoinCache]
updateInterval = "1m"
`, contractAddr.Hex(), ocrKeyBundleID, transmitterAddr.Hex(), barePeerID, chainID, fromBlock)
}

// --- Anvil helpers ---

func waitMined(ctx context.Context, client *ethclient.Client, txHash common.Hash, label string) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			if receipt.Status == 0 {
				log.Fatalf("%s: transaction reverted", label)
			}
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func fundAddress(ctx context.Context, client *ethclient.Client, fromKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) (common.Hash, error) {
	nonce, err := client.PendingNonceAt(ctx, crypto.PubkeyToAddress(fromKey.PublicKey))
	if err != nil {
		return common.Hash{}, err
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	tx := types.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), fromKey)
	if err != nil {
		return common.Hash{}, err
	}
	return signedTx.Hash(), client.SendTransaction(ctx, signedTx)
}

// --- Chainlink node HTTP client ---

type nodeClient struct {
	url    string
	client *http.Client
}

func newNodeClient(nodeURL string) *nodeClient {
	jar, _ := cookiejar.New(nil)
	nc := &nodeClient{
		url:    nodeURL,
		client: &http.Client{Jar: jar},
	}
	nc.login()
	return nc
}

func (nc *nodeClient) login() {
	body, _ := json.Marshal(map[string]string{
		"email":    apiEmail,
		"password": apiPassword,
	})
	resp, err := nc.client.Post(nc.url+"/sessions", "application/json", bytes.NewReader(body))
	must(err, "login")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("login to %s failed (%d): %s", nc.url, resp.StatusCode, b)
	}
}

func (nc *nodeClient) get(path string) json.RawMessage {
	resp, err := nc.client.Get(nc.url + path)
	must(err, "GET "+path)
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("GET %s failed (%d): %s", path, resp.StatusCode, b)
	}
	return b
}

func (nc *nodeClient) post(path string, payload interface{}) json.RawMessage {
	var body io.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		body = bytes.NewReader(b)
	}
	resp, err := nc.client.Post(nc.url+path, "application/json", body)
	must(err, "POST "+path)
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("POST %s failed (%d): %s", path, resp.StatusCode, b)
	}
	return b
}

// getETHAddress returns the first ETH key address from the node.
func (nc *nodeClient) getETHAddress() common.Address {
	raw := nc.get("/v2/keys/eth")
	var resp struct {
		Data []struct {
			Attributes struct {
				Address string `json:"address"`
			} `json:"attributes"`
		} `json:"data"`
	}
	must(json.Unmarshal(raw, &resp), "parse ETH keys")
	if len(resp.Data) == 0 {
		log.Fatal("no ETH keys found")
	}
	return common.HexToAddress(resp.Data[0].Attributes.Address)
}

// createOCR2Key creates a new OCR2 EVM key bundle and returns its components.
func (nc *nodeClient) createOCR2Key() (id string, onchainPub []byte, offchainPub [32]byte, configPub [32]byte) {
	raw := nc.post("/v2/keys/ocr2/evm", nil)
	var resp struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				OnchainPublicKey  string `json:"onchainPublicKey"`
				OffchainPublicKey string `json:"offchainPublicKey"`
				ConfigPublicKey   string `json:"configPublicKey"`
			} `json:"attributes"`
		} `json:"data"`
	}
	must(json.Unmarshal(raw, &resp), "parse OCR2 key")
	id = resp.Data.ID

	// Keys come back as "ocr2on_evm_<hex>", "ocr2off_evm_<hex>", "ocr2cfg_evm_<hex>"
	onchainPub = mustDecodeOCR2Key(resp.Data.Attributes.OnchainPublicKey, "ocr2on_evm_")

	offBytes := mustDecodeOCR2Key(resp.Data.Attributes.OffchainPublicKey, "ocr2off_evm_")
	copy(offchainPub[:], offBytes)

	cfgBytes := mustDecodeOCR2Key(resp.Data.Attributes.ConfigPublicKey, "ocr2cfg_evm_")
	copy(configPub[:], cfgBytes)

	return
}

// getP2PID returns the P2P peer ID for this node.
func (nc *nodeClient) getP2PID() string {
	raw := nc.get("/v2/keys/p2p")
	var resp struct {
		Data []struct {
			Attributes struct {
				PeerID string `json:"peerId"`
			} `json:"attributes"`
		} `json:"data"`
	}
	must(json.Unmarshal(raw, &resp), "parse P2P keys")
	if len(resp.Data) == 0 {
		log.Fatal("no P2P keys found")
	}
	return resp.Data[0].Attributes.PeerID
}

// createJob creates a job from a TOML spec.
func (nc *nodeClient) createJob(tomlSpec string) {
	nc.post("/v2/jobs", map[string]string{"toml": tomlSpec})
}

// --- Utility ---

func mustDecodeOCR2Key(value, prefix string) []byte {
	hexStr := strings.TrimPrefix(value, prefix)
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, err := hex.DecodeString(hexStr)
	must(err, "decode OCR2 key hex: "+value)
	return b
}

func must(err error, label string) {
	if err != nil {
		log.Fatalf("%s: %v", label, err)
	}
}
