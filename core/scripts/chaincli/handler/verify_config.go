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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/manyminds/api2go/jsonapi"
	ocr2config "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/cmd"
	registry20 "github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/logger"
)

const responseSizeLimit = 10_000_000

type CSAKeyInfo struct {
	NodeName    string         `json:"nodeName"`
	NodeAddress common.Address `json:"nodeAddress"`
	PublicKey   string         `json:"publicKey"`
}

type NodeInfo struct {
	AdminAddress          common.Address   `json:"adminAddress"`
	CSAKeys               []*CSAKeyInfo    `json:"csaKeys"`
	DisplayName           string           `json:"displayName"`
	Ocr2ConfigPublicKey   []string         `json:"ocr2ConfigPublicKey,omitempty"`
	Ocr2ID                []string         `json:"ocr2ID,omitempty"`
	Ocr2OffchainPublicKey []string         `json:"ocr2OffchainPublicKey,omitempty"`
	Ocr2OnchainPublicKey  []string         `json:"ocr2OnchainPublicKey,omitempty"`
	NodeAddress           []common.Address `json:"nodeAddress"`
	OcrSigningAddress     []common.Address `json:"ocrSigningAddress"`
	PayeeAddress          common.Address   `json:"payeeAddress"`
	PeerId                []string         `json:"peerId"`
	Status                string           `json:"status"`
}

type ContractInfo struct {
	ContractAddress string                    `json:"contractAddress"`
	Environment     string                    `json:"environment"`
	Name            string                    `json:"name"`
	Version         string                    `json:"version"`
	RegistryConfig  *AutomationRegistryConfig `json:"registryConfig"`
}

type AutomationRegistryConfig struct {
	OnChainConfig AutomationOnChainConfig        `json:"onChainConfig"`
	PublicConfig  AutomationRegistryPublicConfig `json:"publicConfig"`
	Nodes         []string                       `json:"nodes"`
}

type AutomationOnChainConfig struct {
	PaymentPremiumPPB    int    `json:"paymentPremiumPPB"`
	FlatFeeMicroLink     int    `json:"flatFeeMicroLink"`
	CheckGasLimit        int    `json:"checkGasLimit"`
	StalenessSeconds     int    `json:"stalenessSeconds"`
	GasCeilingMultiplier int    `json:"gasCeilingMultiplier"`
	MinUpkeepSpend       string `json:"minUpkeepSpend"`
	MaxPerformGas        int    `json:"maxPerformGas"`
	MaxCheckDataSize     int    `json:"maxCheckDataSize"`
	MaxPerformDataSize   int    `json:"maxPerformDataSize"`
	FallbackGasPrice     string `json:"fallbackGasPrice"`
	FallbackLinkPrice    string `json:"fallbackLinkPrice"`
	Transcoder           string `json:"transcoder"`
	Registrar            string `json:"registrar"`
}

type AutomationRegistryPublicConfig struct {
	DeltaProgress                           string                `json:"deltaProgress"`
	DeltaResend                             string                `json:"deltaResend"`
	DeltaRound                              string                `json:"deltaRound"`
	DeltaGrace                              string                `json:"deltaGrace"`
	DeltaStage                              string                `json:"deltaStage"`
	F                                       uint32                `json:"f"`
	S                                       []int                 `json:"s"`
	MaxDurationQuery                        string                `json:"maxDurationQuery"`
	MaxDurationObservation                  string                `json:"maxDurationObservation"`
	MaxDurationReport                       string                `json:"maxDurationReport"`
	MaxDurationShouldAcceptFinalizedReport  string                `json:"maxDurationShouldAcceptFinalizedReport"`
	MaxDurationShouldTransmitAcceptedReport string                `json:"maxDurationShouldTransmitAcceptedReport"`
	RMax                                    uint32                `json:"rMax"`
	ReportingPluginConfig                   ReportingPluginConfig `json:"reportingPluginConfig"`
}

type ReportingPluginConfig struct {
	PerformLockoutWindow int    `json:"performLockoutWindow"`
	TargetProbability    string `json:"targetProbability"`
	TargetInRounds       int    `json:"targetInRounds"`
	SamplingJobDuration  int    `json:"samplingJobDuration"`
	MinConfirmations     int    `json:"minConfirmations"`
	GasLimitPerReport    int    `json:"gasLimitPerReport"`
	GasOverheadPerUpkeep int    `json:"gasOverheadPerUpkeep"`
	MaxUpkeepBatchSize   int    `json:"maxUpkeepBatchSize"`
	ReportBlockLag       int    `json:"reportBlockLag"`
}

func (h *baseHandler) Verify() {
	log, closeLggr := logger.NewLogger()
	logger.Sugared(log).ErrorIfFn(closeLggr, "Failed to close logger")

	ctx := context.Background()
	oracleIdentities := h.verifyNodes(ctx, log)
	h.verifyOnchain(ctx, oracleIdentities, log)
}

func (h *baseHandler) verifyOnchain(ctx context.Context, oracleIdentities []ocr2config.OracleIdentityExtra, log logger.Logger) {
	contractInfos := fetchContractInfos(ctx, h.cfg.AutomationURL, log)
	var ci *ContractInfo
	for _, c := range contractInfos {
		if c.ContractAddress == h.cfg.RegistryAddress {
			ci = &c
			break
		}
	}
	if ci == nil {
		log.Fatalf("Failed to find a registry contract at %s in Automation contract URL provided %s", h.cfg.RegistryAddress, h.cfg.AutomationURL)
	}

	blockNum, err := h.client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("Failed to get block number: %v", err)
	}

	registryAddr := common.HexToAddress(h.cfg.RegistryAddress)
	keeperRegistry20, err := registry20.NewKeeperRegistry(
		registryAddr,
		h.client,
	)
	if err != nil {
		log.Fatal("Registry failed: %v", err)
	}

	log.Warn(blockNum)
	log.Warn(blockNum - h.cfg.LookBackDepth)
	iter, err := keeperRegistry20.FilterConfigSet(&bind.FilterOpts{
		Start:   blockNum - h.cfg.LookBackDepth,
		End:     &blockNum,
		Context: ctx,
	})
	if err != nil {
		log.Fatal("failed to retrieve ConfigSet events: %v", err)
	}

	if iter.Next() {
		log.Warn(iter.Event.F)
		log.Warn(iter.Event.PreviousConfigBlockNumber)
	}

	gs, err := keeperRegistry20.GetState(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		log.Fatalf("failed to call GetState on %s", registryAddr.Hex())
	}

	rddCfg := ci.RegistryConfig.OnChainConfig
	onchainCfg := gs.Config
	diffs := 0
	if onchainCfg.PaymentPremiumPPB != uint32(rddCfg.PaymentPremiumPPB) {
		diffs++
		log.Errorf("PaymentPremiumPPB differs. The onchain config has %d but RDD has %d", onchainCfg.PaymentPremiumPPB, rddCfg.PaymentPremiumPPB)
	}
	if onchainCfg.FlatFeeMicroLink != uint32(rddCfg.FlatFeeMicroLink) {
		diffs++
		log.Errorf("FlatFeeMicroLink differs. The onchain config has %d but RDD has %d", onchainCfg.FlatFeeMicroLink, rddCfg.FlatFeeMicroLink)
	}
	if onchainCfg.CheckGasLimit != uint32(rddCfg.CheckGasLimit) {
		diffs++
		log.Errorf("CheckGasLimit differs. The onchain config has %d but RDD has %d", onchainCfg.CheckGasLimit, rddCfg.CheckGasLimit)
	}
	if onchainCfg.StalenessSeconds.Uint64() != uint64(rddCfg.StalenessSeconds) {
		diffs++
		log.Errorf("StalenessSeconds differs. The onchain config has %d but RDD has %d", onchainCfg.StalenessSeconds.Uint64(), rddCfg.StalenessSeconds)
	}
	if onchainCfg.GasCeilingMultiplier != uint16(rddCfg.GasCeilingMultiplier) {
		diffs++
		log.Errorf("GasCeilingMultiplier differs. The onchain config has %d but RDD has %d", onchainCfg.GasCeilingMultiplier, rddCfg.GasCeilingMultiplier)
	}
	if onchainCfg.MinUpkeepSpend.String() != rddCfg.MinUpkeepSpend {
		diffs++
		log.Errorf("MinUpkeepSpend differs. The onchain config has %s but RDD has %s", onchainCfg.MinUpkeepSpend.String(), rddCfg.MinUpkeepSpend)
	}
	if onchainCfg.MaxPerformGas != uint32(rddCfg.MaxPerformGas) {
		diffs++
		log.Errorf("MaxPerformGas differs. The onchain config has %d but RDD has %d", onchainCfg.MaxPerformGas, rddCfg.MaxPerformGas)
	}
	if onchainCfg.MaxCheckDataSize != uint32(rddCfg.MaxCheckDataSize) {
		diffs++
		log.Errorf("MaxCheckDataSize differs. The onchain config has %d but RDD has %d", onchainCfg.MaxCheckDataSize, rddCfg.MaxCheckDataSize)
	}
	if onchainCfg.MaxPerformDataSize != uint32(rddCfg.MaxPerformDataSize) {
		diffs++
		log.Errorf("MaxPerformDataSize differs. The onchain config has %d but RDD has %d", onchainCfg.MaxPerformDataSize, rddCfg.MaxPerformDataSize)
	}
	if onchainCfg.FallbackGasPrice.String() != rddCfg.FallbackGasPrice {
		diffs++
		log.Errorf("FallbackGasPrice differs. The onchain config has %s but RDD has %s", onchainCfg.FallbackGasPrice, rddCfg.FallbackGasPrice)
	}
	if onchainCfg.FallbackLinkPrice.String() != rddCfg.FallbackLinkPrice {
		diffs++
		log.Errorf("FallbackLinkPrice differs. The onchain config has %s but RDD has %s", onchainCfg.FallbackLinkPrice.String(), rddCfg.FallbackLinkPrice)
	}
	if onchainCfg.Transcoder.String() != rddCfg.Transcoder {
		diffs++
		log.Errorf("Transcoder differs. The onchain config has %s but RDD has %s", onchainCfg.Transcoder.String(), rddCfg.Transcoder)
	}
	if onchainCfg.Registrar.String() != rddCfg.Registrar {
		diffs++
		log.Errorf("Registrar differs. The onchain config has %d but RDD has %d", onchainCfg.Registrar.String(), rddCfg.Registrar)
	}

	if diffs == 0 {
		log.Infof("Onchain config for contract at %s matches the RDD file.", h.cfg.RegistryAddress)
	}
}

func (h *baseHandler) verifyNodes(ctx context.Context, log logger.Logger) []ocr2config.OracleIdentityExtra {
	log.Warn("This verifies node address, peer ID, CSA node address, CSA public key, OCR2 ID, OCR2 config pub key, OCR2 onchain pub key, and OCR2 offchain pub key.")
	log.Warn("This does NOT check for payee address. Please verify that manually.")
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
			var nodeAddresses []common.Address
			for index := range ethKeys {
				nodeAddresses = append(nodeAddresses, common.HexToAddress(ethKeys[index].Address))
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
			nodes[nodeAddresses[0].Hex()] = ni

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
					TransmitAccount:   ocr2types.Account(ni.NodeAddress[0].String()),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}

			web3.
				ni.OcrSigningAddress = []common.Address{common.HexToAddress(strings.TrimPrefix(ocr2Config.OnchainPublicKey, "ocr2on_evm_"))}
			err = writeJSON(ni, strconv.Itoa(i)+".json")
			if err != nil {
				panic(fmt.Errorf("failed to write node info to JSON: %v", err))
			}

		}(i, cl)
	}
	wg.Wait()

	nodeInfos := fetchNodeInfos(ctx, h.cfg.NodeConfigURL, log)
	cnt := 0
	for _, ni := range nodeInfos {
		nodeAddr := ni.NodeAddress[0]
		node := nodes[nodeAddr.Hex()]
		if node == nil {
			continue
		}
		cnt++
		diffs := 0

		log.Infof("start comparing data for node %s", nodeAddr.String())
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
			log.Infof("node %s info is correct", nodeAddr.String())
		}
	}

	if cnt != len(nodes) {
		log.Info("there are %d nodes provisioned, but weiwatcher is missing %d nodes", len(nodes), len(nodes)-cnt)
	}

	return oracleIdentities
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

func fetchContractInfos(ctx context.Context, automationContractsConfigURL string, log logger.Logger) []ContractInfo {
	contractBytes := fetchRawBytes(ctx, automationContractsConfigURL, log)

	var contractInfos []ContractInfo
	if err := json.Unmarshal(contractBytes, &contractInfos); err != nil {
		log.Fatalf("failed to unmarshal response: %s", err)
	}

	return contractInfos
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
