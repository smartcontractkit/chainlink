package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type CSAKeyInfo struct {
	NodeName    string `json:"nodeName"`
	NodeAddress string `json:"nodeAddress"`
	PublicKey   string `json:"publicKey"`
}

func (ci *CSAKeyInfo) Equals(ci2 *CSAKeyInfo) bool {
	return ci.PublicKey == ci2.PublicKey && ci.NodeAddress == ci.NodeAddress
}

type NodeInfo struct {
	AdminAddress          common.Address `json:"adminAddress"`
	CSAKeys               []*CSAKeyInfo  `json:"csaKeys"`
	DisplayName           string         `json:"displayName"`
	Ocr2ConfigPublicKey   []string       `json:"ocr2ConfigPublicKey"`
	Ocr2ID                []string       `json:"ocr2ID"`
	Ocr2OffchainPublicKey []string       `json:"ocr2OffchainPublicKey"`
	Ocr2OnchainPublicKey  []string       `json:"ocr2OnchainPublicKey"`
	NodeAddress           []string       `json:"nodeAddress"`
	OcrSigningAddress     []string       `json:"ocrSigningAddress"`
	PayeeAddress          common.Address `json:"payeeAddress"`
	PeerId                []string       `json:"peerId"`
	Status                string         `json:"status"`
}

func (node NodeInfo) Equals(ni NodeInfo, log logger.Logger) bool {
	diffs := 0

	if len(node.CSAKeys) != len(ni.CSAKeys) {
		log.Errorf("CSA Keys length differs. The node returns %d but weiwatcher has %d", len(node.CSAKeys), len(ni.CSAKeys))
	}
	for i, ci := range node.CSAKeys {
		if !ci.Equals(ni.CSAKeys[i]) {
			diffs++
			log.Errorf("CSA Info differs. The node returns %s but weiwatcher has %s", ci, ni.CSAKeys[i])
		}
	}

	if !cmp.Equal(node.Ocr2ID, ni.Ocr2ID) {
		diffs++
		log.Errorf("OCR2 ID differs. The node returns %s but weiwatcher has %s", node.Ocr2ID, ni.Ocr2ID)
	}

	if !cmp.Equal(node.NodeAddress, ni.NodeAddress) {
		diffs++
		log.Errorf("Node address differs. The node returns %s but weiwatcher has %s", node.NodeAddress, ni.NodeAddress)
	}

	// preprocess the Peer ID from node bc it has p2p_ prefix
	var peerIds []string
	for _, pid := range node.PeerId {
		peerIds = append(peerIds, pid[4:])
	}
	if !cmp.Equal(peerIds, ni.PeerId) {
		diffs++
		log.Errorf("Peer Id differs. The node returns %s but weiwatcher has %s", node.PeerId, ni.PeerId)
	}

	if !cmp.Equal(node.Ocr2OffchainPublicKey, ni.Ocr2OffchainPublicKey) {
		diffs++
		log.Errorf("OCR2 Offchain Public Key differs. The node returns %s but weiwatcher has %s", node.Ocr2OffchainPublicKey, ni.Ocr2OffchainPublicKey)
	}

	if !cmp.Equal(node.Ocr2OnchainPublicKey, ni.Ocr2OnchainPublicKey) {
		diffs++
		log.Errorf("OCR2 Onchain Public Key differs. The node returns %s but weiwatcher has %s", node.Ocr2OnchainPublicKey, ni.Ocr2OnchainPublicKey)
	}

	if !cmp.Equal(node.Ocr2ConfigPublicKey, ni.Ocr2ConfigPublicKey) {
		diffs++
		log.Errorf("OCR2 Config Public Key differs. The node returns %s but weiwatcher has %s", node.Ocr2ConfigPublicKey, ni.Ocr2ConfigPublicKey)
	}

	return diffs == 0
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
	for i, cl := range cls {
		wg.Add(1)
		go h.scrapeNodeInfo(&wg, i, cl, nodes, log)
	}
	wg.Wait()

	// if node info is not in RDD and weiwatchers, don't proceed further
	if !h.cfg.VerifyNodes {
		return
	}
	nodeInfos := h.fetchNodeInfos(ctx, log)
	cnt := 0
	for _, ni := range nodeInfos {
		nodeAddr := ni.NodeAddress[0]
		node := nodes[nodeAddr]
		if node == nil {
			continue
		}
		cnt++

		log.Infof("start comparing data for node %s", nodeAddr)
		if node.Equals(ni, log) {
			log.Infof("node %s info is correct", nodeAddr)
		} else {
			log.Errorf("node %s info differs between the node instance and weiwatcher", nodeAddr)
		}
	}

	if cnt != len(nodes) {
		log.Infof("there are %d nodes provisioned , but .env is missing %d nodes", len(nodes), len(nodes)-cnt)
	}
}

func (h *baseHandler) fetchNodeInfos(ctx context.Context, log logger.Logger) []NodeInfo {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.cfg.NodeConfigURL, nil)
	if err != nil {
		log.Fatalf("failed to build a GET request: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to make a GET request: %s", err)
	}
	defer resp.Body.Close()

	var nodeInfos []NodeInfo
	if err := json.NewDecoder(resp.Body).Decode(&nodeInfos); err != nil {
		log.Fatalf("failed to read response: %s", err)
	}

	return nodeInfos
}

func (h *baseHandler) scrapeNodeInfo(wg *sync.WaitGroup, i int, cl cmd.HTTPClient, nodes map[string]*NodeInfo, log logger.Logger) {
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
	// this assumes the nodes are not multichain nodes and have only 1 node address assigned.
	// for a multichain node, we can pass in a chain id and filter `ethKeys` array based on the chain id
	// in terms of CSA keys, we need to wait for RTSP to support multichain nodes, which may involve creating one
	// CSA key for each chain. but this is still pending so assume only 1 CSA key on a node for now.
	csaKey := &CSAKeyInfo{
		NodeAddress: nodeAddresses[0],
		PublicKey:   csaKeys[0].PubKey,
	}
	ni.CSAKeys = []*CSAKeyInfo{csaKey}

	ni.OcrSigningAddress = []string{common.HexToAddress(strings.TrimPrefix(ocr2Config.OnchainPublicKey, "ocr2on_evm_")).Hex()}
	err = writeJSON(ni, strconv.Itoa(i)+".json")
	if err != nil {
		panic(fmt.Errorf("failed to write node info to JSON: %v", err))
	}

	nodes[nodeAddresses[0]] = ni
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
