package handler

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manyminds/api2go/jsonapi"
	ocr2config "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const responseSizeLimit = 10_000_000

type CSAKeyInfo struct {
	NodeName    string `json:"nodeName"`
	NodeAddress string `json:"nodeAddress"`
	PublicKey   string `json:"publicKey"`
}

type NodeInfo struct {
	AdminAddress          common.Address `json:"adminAddress"`
	CSAKeys               []*CSAKeyInfo  `json:"csaKeys"`
	DisplayName           string         `json:"displayName"`
	Ocr2ConfigPublicKey   []string       `json:"ocr2ConfigPublicKey,omitempty"`
	Ocr2ID                []string       `json:"ocr2ID,omitempty"`
	Ocr2OffchainPublicKey []string       `json:"ocr2OffchainPublicKey,omitempty"`
	Ocr2OnchainPublicKey  []string       `json:"ocr2OnchainPublicKey,omitempty"`
	NodeAddress           []string       `json:"nodeAddress"`
	OcrSigningAddress     []string       `json:"ocrSigningAddress"`
	PayeeAddress          common.Address `json:"payeeAddress"`
	PeerId                []string       `json:"peerId"`
	Status                string         `json:"status"`
}

func (h *baseHandler) ScrapeNodes() {
	log, closeLggr := logger.NewLogger()
	logger.Sugared(log).ErrorIfFn(closeLggr, "Failed to close logger")

	ctx := context.Background()
	h.scrapeNodes(ctx, log)
}

func (h *baseHandler) scrapeNodes(ctx context.Context, log logger.Logger) {
	log.Warn("This scrapes node address, peer ID, CSA node address, CSA public key, OCR2 ID, OCR2 config pub key, OCR2 onchain pub key, and OCR2 offchain pub key.")
	log.Warn("This does NOT scrape for payee address, admin address etc. Please verify that manually.")
	cls := make([]cmd.HTTPClient, len(h.cfg.KeeperURLs))
	for i := range h.cfg.KeeperURLs {
		url := h.cfg.KeeperURLs[i]
		email := h.cfg.KeeperEmails[i]
		if len(email) == 0 {
			email = defaultChainlinkNodeLogin
		}
		pwd := h.cfg.KeeperPasswords[i]
		if len(pwd) == 0 {
			pwd = defaultChainlinkNodePassword
		}

		cl, err := authenticate(url, email, pwd, log)
		if err != nil {
			log.Fatal(err)
		}
		cls[i] = cl
	}

	nodes := map[string]*NodeInfo{}
	var wg sync.WaitGroup
	oracleIdentities := make([]ocr2config.OracleIdentityExtra, len(cls))
	for i, cl := range cls {
		wg.Add(1)
		go func(i int, cl cmd.HTTPClient) {
			defer wg.Done()

			// get node addresses
			resp, err := nodeRequest(cl, "/v2/keys/eth")
			if err != nil {
				log.Fatalf("failed to get ETH keys: %s", err)
			}
			var ethKeys cmd.EthKeyPresenters
			if err = jsonapi.Unmarshal(resp, &ethKeys); err != nil {
				log.Fatalf("failed to unmarshal response body: %s", err)
			}
			var nodeAddresses []string
			for index := range ethKeys {
				nodeAddresses = append(nodeAddresses, common.HexToAddress(ethKeys[index].Address).Hex())
			}
			csaKey := &CSAKeyInfo{
				NodeAddress: nodeAddresses[0],
			}
			ni := &NodeInfo{
				NodeAddress: nodeAddresses,
			}

			// get node ocr2 config
			ocr2Config, err := getNodeOCR2Config(cl)
			if err != nil {
				log.Fatalf("failed to get node OCR2 config: %s", err)
			}
			ni.Ocr2ID = []string{ocr2Config.ID}
			ni.Ocr2OnchainPublicKey = []string{ocr2Config.OnchainPublicKey}
			ni.Ocr2OffchainPublicKey = []string{ocr2Config.OffChainPublicKey}
			ni.Ocr2ConfigPublicKey = []string{ocr2Config.ConfigPublicKey}

			// get node p2p config
			resp, err = nodeRequest(cl, "/v2/keys/p2p")
			if err != nil {
				log.Fatalf("failed to get p2p keys: %s", err)
			}
			var p2pKeys cmd.P2PKeyPresenters
			if err = jsonapi.Unmarshal(resp, &p2pKeys); err != nil {
				log.Fatalf("failed to unmarshal response body: %s", err)
			}
			ni.PeerId = []string{p2pKeys[0].PeerID}

			// get node csa config
			resp, err = nodeRequest(cl, "/v2/keys/csa")
			if err != nil {
				log.Fatalf("failed to get CSA keys: %s", err)
			}
			var csaKeys cmd.CSAKeyPresenters
			if err = jsonapi.Unmarshal(resp, &csaKeys); err != nil {
				log.Fatalf("failed to unmarshal response body: %s", err)
			}
			csaKey.PublicKey = csaKeys[0].PubKey

			ni.CSAKeys = []*CSAKeyInfo{csaKey}
			nodes[nodeAddresses[0]] = ni

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			if err != nil {
				panic(fmt.Errorf("failed to decode %s: %v", ocr2Config.OffChainPublicKey, err))
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			if n != ed25519.PublicKeySize {
				panic(fmt.Errorf("wrong num elements copied"))
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				panic(fmt.Errorf("failed to decode %s: %v", ocr2Config.ConfigPublicKey, err))
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				panic(fmt.Errorf("wrong num elements copied"))
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnchainPublicKey, "ocr2on_evm_"))
			if err != nil {
				panic(fmt.Errorf("failed to decode %s: %v", ocr2Config.OnchainPublicKey, err))
			}

			oracleIdentities[i] = ocr2config.OracleIdentityExtra{
				OracleIdentity: ocr2config.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeys[0].PeerID,
					TransmitAccount:   ocr2types.Account(ni.NodeAddress[0]),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}

			ni.OcrSigningAddress = []string{common.HexToAddress(strings.TrimPrefix(ocr2Config.OnchainPublicKey, "ocr2on_evm_")).Hex()}
			err = writeJSON(ni, strconv.Itoa(i)+".json")
			if err != nil {
				panic(fmt.Errorf("failed to write node info to JSON: %v", err))
			}

		}(i, cl)
	}
	wg.Wait()

	// if node info is not in RDD and weiwatchers, don't proceed further
	if !h.cfg.VerifyNodes {
		return
	}
	nodeInfos := fetchNodeInfos(ctx, h.cfg.NodeConfigURL, log)
	cnt := 0
	for _, ni := range nodeInfos {
		nodeAddr := ni.NodeAddress[0]
		node := nodes[nodeAddr]
		if node == nil {
			continue
		}
		cnt++
		diffs := 0

		log.Infof("start comparing data for node %s", nodeAddr)
		if node.CSAKeys[0].NodeAddress != ni.CSAKeys[0].NodeAddress {
			diffs++
			log.Errorf("CSA Node Address differs. The node returns %s but weiwatcher has %s", node.CSAKeys[0].NodeAddress, ni.CSAKeys[0].NodeAddress)
		}
		if node.CSAKeys[0].PublicKey != ni.CSAKeys[0].PublicKey {
			diffs++
			log.Errorf("CSA Public Key differs. The node returns %s but weiwatcher has %s", node.CSAKeys[0].PublicKey, ni.CSAKeys[0].PublicKey)
		}

		if !reflect.DeepEqual(node.Ocr2ID, ni.Ocr2ID) {
			diffs++
			log.Errorf("OCR2 ID differs. The node returns %s but weiwatcher has %s", node.Ocr2ID, ni.Ocr2ID)
		}

		if !reflect.DeepEqual(node.NodeAddress, ni.NodeAddress) {
			diffs++
			log.Errorf("Node address differs. The node returns %s but weiwatcher has %s", node.NodeAddress, ni.NodeAddress)
		}

		// preprocess the Peer ID from node bc it has p2p_ prefix
		var peerIds []string
		for _, pid := range node.PeerId {
			peerIds = append(peerIds, pid[4:])
		}
		if !reflect.DeepEqual(peerIds, ni.PeerId) {
			diffs++
			log.Errorf("Peer Id differs. The node returns %s but weiwatcher has %s", node.PeerId, ni.PeerId)
		}

		if !reflect.DeepEqual(node.Ocr2OffchainPublicKey, ni.Ocr2OffchainPublicKey) {
			diffs++
			log.Errorf("OCR2 Offchain Public Key differs. The node returns %s but weiwatcher has %s", node.Ocr2OffchainPublicKey, ni.Ocr2OffchainPublicKey)
		}

		if !reflect.DeepEqual(node.Ocr2OnchainPublicKey, ni.Ocr2OnchainPublicKey) {
			diffs++
			log.Errorf("OCR2 Onchain Public Key differs. The node returns %s but weiwatcher has %s", node.Ocr2OnchainPublicKey, ni.Ocr2OnchainPublicKey)
		}

		if !reflect.DeepEqual(node.Ocr2ConfigPublicKey, ni.Ocr2ConfigPublicKey) {
			diffs++
			log.Errorf("OCR2 Config Public Key differs. The node returns %s but weiwatcher has %s", node.Ocr2ConfigPublicKey, ni.Ocr2ConfigPublicKey)
		}
		if diffs == 0 {
			log.Infof("node %s info is correct", nodeAddr)
		}
	}

	if cnt != len(nodes) {
		log.Infof("there are %d nodes provisioned , but .env is missing %d nodes", len(nodes), len(nodes)-cnt)
	}
}

func JSONMarshalWithoutEscape(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func writeJSON(data interface{}, path string) error {
	dataBytes, err := JSONMarshalWithoutEscape(data)
	if err != nil {
		return err
	}

	return os.WriteFile(path, dataBytes, 0644)
}

func fetchNodeInfos(ctx context.Context, automationNodesConfigURL string, log logger.Logger) []NodeInfo {
	nodeBytes := fetchRawBytes(ctx, automationNodesConfigURL, log)

	var nodeInfos []NodeInfo
	if err := json.Unmarshal(nodeBytes, &nodeInfos); err != nil {
		log.Fatalf("failed to unmarshal response: %s", err)
	}

	return nodeInfos
}

func fetchRawBytes(ctx context.Context, url string, log logger.Logger) []byte {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("failed to build a GET request: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to make a GET request: %s", err)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, responseSizeLimit))
	if err != nil {
		log.Fatalf("failed to read response: %s", err)
	}

	return body
}
