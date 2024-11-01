package src

// This package deploys "offchainreporting2" job specs, which setup the streams trigger
// for the targeted node set
// See https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/integration_test.go#L92
// for how to setup the mercury portion of the streams trigger
//  You can see how all fields are being used here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/helpers_test.go#L314
//  https://github.com/smartcontractkit/infra-k8s/blob/be47098adfb605d79b5bab6aa601bcf443a6c48b/projects/chainlink/files/chainlink-clusters/cl-keystone-cap-one/config.yaml#L1
//  Trigger gets added to the registry here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/relay/evm/evm.go#L360
//  See integration workflow here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/capabilities/integration_tests/workflow.go#L15
import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"math/big"
	"time"

	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mercurytypes "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	datastreamsmercury "github.com/smartcontractkit/chainlink-data-streams/mercury"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"

	verifierContract "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type feed struct {
	id   [32]byte
	name string

	// we create a bridge for each feed
	bridgeName string
	bridgeURL  string
}

func v3FeedID(id [32]byte) [32]byte {
	binary.BigEndian.PutUint16(id[:2], 3)
	return id
}

var feeds = []feed{
	{
		v3FeedID([32]byte{5: 1}),
		"BTC/USD",
		"mock-bridge-btc",
		"http://external-adapter:4001",
	},
	{
		v3FeedID([32]byte{5: 2}),
		"LINK/USD",
		"mock-bridge-link",
		"http://external-adapter:4002",
	},
	{
		v3FeedID([32]byte{5: 3}),
		"NATIVE/USD",
		"mock-bridge-native",
		"http://external-adapter:4003",
	},
}

// See /core/services/ocr2/plugins/mercury/integration_test.go
func setupStreamsTrigger(
	env helpers.Environment,
	nodeSet NodeSet,
	chainID int64,
	p2pPort int64,
	ocrConfigFilePath string,
	artefactsDir string,
) {
	fmt.Printf("Deploying streams trigger for chain %d\n", chainID)
	fmt.Printf("Using OCR config file: %s\n", ocrConfigFilePath)

	fmt.Printf("Deploying Mercury V0.3 contracts\n")
	verifier := deployMercuryV03Contracts(env, artefactsDir)

	fmt.Printf("Generating Mercury OCR config\n")
	ocrConfig := generateMercuryOCR2Config(nodeSet.NodeKeys[1:]) // skip the bootstrap node

	for _, feed := range feeds {
		fmt.Println("Configuring feeds...")
		fmt.Printf("FeedID: %x\n", feed.id)
		fmt.Printf("FeedName: %s\n", feed.name)
		fmt.Printf("BridgeName: %s\n", feed.bridgeName)
		fmt.Printf("BridgeURL: %s\n", feed.bridgeURL)

		latestConfigDetails, err := verifier.LatestConfigDetails(nil, feed.id)
		PanicErr(err)
		latestConfigDigest, err := ocrtypes.BytesToConfigDigest(latestConfigDetails.ConfigDigest[:])
		PanicErr(err)

		digester := mercury.NewOffchainConfigDigester(
			feed.id,
			big.NewInt(chainID),
			verifier.Address(),
			ocrtypes.ConfigDigestPrefixMercuryV02,
		)
		configDigest, err := digester.ConfigDigest(
			context.Background(),
			mercuryOCRConfigToContractConfig(
				ocrConfig,
				latestConfigDetails.ConfigCount,
			),
		)
		PanicErr(err)

		if configDigest.Hex() == latestConfigDigest.Hex() {
			fmt.Printf("Verifier already deployed with the same config (digest: %s), skipping...\n", configDigest.Hex())
		} else {
			fmt.Printf("Verifier contains a different config, updating...\nOld digest: %s\nNew digest: %s\n", latestConfigDigest.Hex(), configDigest.Hex())
			tx, err := verifier.SetConfig(
				env.Owner,
				feed.id,
				ocrConfig.Signers,
				ocrConfig.Transmitters,
				ocrConfig.F,
				ocrConfig.OnchainConfig,
				ocrConfig.OffchainConfigVersion,
				ocrConfig.OffchainConfig,
				nil,
			)
			helpers.ConfirmTXMined(context.Background(), env.Ec, tx, env.ChainID)
			PanicErr(err)
		}

		fmt.Printf("Deploying OCR2 job specs for feed %s\n", feed.name)
		deployOCR2JobSpecsForFeed(nodeSet, verifier, feed, chainID, p2pPort)
	}

	fmt.Println("Finished deploying streams trigger")
}

func deployMercuryV03Contracts(env helpers.Environment, artefactsDir string) verifierContract.VerifierInterface {
	var confirmDeploy = func(tx *types.Transaction, err error) {
		helpers.ConfirmContractDeployed(context.Background(), env.Ec, tx, env.ChainID)
		PanicErr(err)
	}
	o := LoadOnchainMeta(artefactsDir, env)

	if o.VerifierProxy != nil {
		fmt.Printf("Verifier proxy contract already deployed at %s\n", o.VerifierProxy.Address())
	} else {
		fmt.Printf("Deploying verifier proxy contract\n")
		_, tx, verifierProxy, err := verifier_proxy.DeployVerifierProxy(env.Owner, env.Ec, common.Address{}) // zero address for access controller disables access control
		confirmDeploy(tx, err)
		o.VerifierProxy = verifierProxy
		WriteOnchainMeta(o, artefactsDir)
	}

	if o.Verifier == nil {
		fmt.Printf("Deploying verifier contract\n")
		_, tx, verifier, err := verifierContract.DeployVerifier(env.Owner, env.Ec, o.VerifierProxy.Address())
		confirmDeploy(tx, err)
		o.Verifier = verifier
		WriteOnchainMeta(o, artefactsDir)
	} else {
		fmt.Printf("Verifier contract already deployed at %s\n", o.Verifier.Address().Hex())
	}

	if o.InitializedVerifierAddress != o.Verifier.Address() {
		fmt.Printf("Current initialized verifier address (%s) differs from the new verifier address (%s). Initializing verifier.\n", o.InitializedVerifierAddress.Hex(), o.Verifier.Address().Hex())
		tx, err := o.VerifierProxy.InitializeVerifier(env.Owner, o.Verifier.Address())
		receipt := helpers.ConfirmTXMined(context.Background(), env.Ec, tx, env.ChainID)
		PanicErr(err)
		inited, err := o.VerifierProxy.ParseVerifierInitialized(*receipt.Logs[0])
		PanicErr(err)
		o.InitializedVerifierAddress = inited.VerifierAddress
		WriteOnchainMeta(o, artefactsDir)
	} else {
		fmt.Printf("Verifier %s already initialized\n", o.Verifier.Address().Hex())
	}

	return o.Verifier
}

func deployOCR2JobSpecsForFeed(nodeSet NodeSet, verifier verifierContract.VerifierInterface, feed feed, chainID int64, p2pPort int64) {
	// we assign the first node as the bootstrap node
	for i, n := range nodeSet.NodeKeys {
		// parallel arrays
		api := newNodeAPI(nodeSet.Nodes[i])
		jobSpecName := ""
		jobSpecStr := ""

		upsertBridge(api, feed.bridgeName, feed.bridgeURL)

		if i == 0 {
			// Prepare data for Bootstrap Job
			bootstrapData := MercuryV3BootstrapJobSpecData{
				FeedName:        feed.name,
				VerifierAddress: verifier.Address().Hex(),
				FeedID:          fmt.Sprintf("%x", feed.id),
				ChainID:         chainID,
			}

			// Create Bootstrap Job
			jobSpecName, jobSpecStr = createMercuryV3BootstrapJob(bootstrapData)
		} else {
			// Prepare data for Mercury V3 Job
			mercuryData := MercuryV3JobSpecData{
				FeedName:        "feed-" + feed.name,
				BootstrapHost:   fmt.Sprintf("%s@%s:%d", nodeSet.NodeKeys[0].P2PPeerID, nodeSet.Nodes[0].ServiceName, p2pPort),
				VerifierAddress: verifier.Address().Hex(),
				Bridge:          feed.bridgeName,
				NodeCSAKey:      n.CSAPublicKey,
				FeedID:          fmt.Sprintf("%x", feed.id),
				LinkFeedID:      fmt.Sprintf("%x", feeds[1].id),
				NativeFeedID:    fmt.Sprintf("%x", feeds[2].id),
				OCRKeyBundleID:  n.OCR2BundleID,
				ChainID:         chainID,
			}

			// Create Mercury V3 Job
			jobSpecName, jobSpecStr = createMercuryV3OracleJob(mercuryData)
		}

		upsertJob(api, jobSpecName, jobSpecStr)
	}
}

// Template definitions
const mercuryV3OCR2bootstrapJobTemplate = `
type                              = "bootstrap"
relay                             = "evm"
schemaVersion                     = 1
name                              = "{{ .Name }}"
contractID                        = "{{ .VerifierAddress }}"
feedID                            = "0x{{ .FeedID }}"
contractConfigTrackerPollInterval = "1s"

[relayConfig]
chainID = {{ .ChainID }}
enableTriggerCapability = true
`

const mercuryV3OCR2OracleJobTemplate = `
type = "offchainreporting2"
schemaVersion = 1
name = "{{ .Name }}"
p2pv2Bootstrappers = ["{{ .BootstrapHost }}"]
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "{{ .VerifierAddress }}"
feedID = "0x{{ .FeedID }}"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "{{ .OCRKeyBundleID }}"
relay = "evm"
pluginType = "mercury"
transmitterID = "{{ .NodeCSAKey }}"
observationSource = """
	price              [type=bridge name="{{ .Bridge }}" timeout="50ms" requestData=""];

	benchmark_price  [type=jsonparse path="result,mid" index=0];
	price -> benchmark_price;

	bid_price [type=jsonparse path="result,bid" index=1];
	price -> bid_price;

	ask_price [type=jsonparse path="result,ask" index=2];
	price -> ask_price;
"""

[relayConfig]
enableTriggerCapability = true
chainID = "{{ .ChainID }}"
`

// Data structures
type MercuryV3BootstrapJobSpecData struct {
	FeedName string
	// Automatically generated from FeedName
	Name            string
	VerifierAddress string
	FeedID          string
	ChainID         int64
}

type MercuryV3JobSpecData struct {
	FeedName string
	// Automatically generated from FeedName
	Name            string
	BootstrapHost   string
	VerifierAddress string
	Bridge          string
	NodeCSAKey      string
	FeedID          string
	LinkFeedID      string
	NativeFeedID    string
	OCRKeyBundleID  string
	ChainID         int64
}

// createMercuryV3BootstrapJob creates a bootstrap job specification using the provided data.
func createMercuryV3BootstrapJob(data MercuryV3BootstrapJobSpecData) (name string, jobSpecStr string) {
	name = "boot-" + data.FeedName
	data.Name = name

	fmt.Printf("Creating bootstrap job (%s):\nverifier address: %s\nfeed name: %s\nfeed ID: %s\nchain ID: %d\n",
		name, data.VerifierAddress, data.FeedName, data.FeedID, data.ChainID)

	tmpl, err := template.New("bootstrapJob").Parse(mercuryV3OCR2bootstrapJobTemplate)
	PanicErr(err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	PanicErr(err)

	jobSpecStr = buf.String()

	return name, jobSpecStr
}

// createMercuryV3OracleJob creates a Mercury V3 job specification using the provided data.
func createMercuryV3OracleJob(data MercuryV3JobSpecData) (name string, jobSpecStr string) {
	name = "mercury-" + data.FeedName
	data.Name = name
	fmt.Printf("Creating ocr2 job(%s):\nOCR key bundle ID: %s\nverifier address: %s\nbridge: %s\nnodeCSAKey: %s\nfeed name: %s\nfeed ID: %s\nlink feed ID: %s\nnative feed ID: %s\nchain ID: %d\n",
		data.Name, data.OCRKeyBundleID, data.VerifierAddress, data.Bridge, data.NodeCSAKey, data.FeedName, data.FeedID, data.LinkFeedID, data.NativeFeedID, data.ChainID)

	tmpl, err := template.New("mercuryV3Job").Parse(mercuryV3OCR2OracleJobTemplate)
	PanicErr(err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	PanicErr(err)

	jobSpecStr = buf.String()

	return data.Name, jobSpecStr
}

func strToBytes32(str string) [32]byte {
	pkBytes, err := hex.DecodeString(str)
	helpers.PanicErr(err)

	pkBytesFixed := [ed25519.PublicKeySize]byte{}
	n := copy(pkBytesFixed[:], pkBytes)
	if n != ed25519.PublicKeySize {
		fmt.Printf("wrong num elements copied (%s): %d != 32\n", str, n)
		panic("wrong num elements copied")
	}
	return pkBytesFixed
}

func upsertBridge(api *nodeAPI, name string, eaURL string) {
	u, err := url.Parse(eaURL)
	helpers.PanicErr(err)
	url := models.WebURL(*u)
	// Confirmations and MinimumContractPayment are not used, so we can leave them as 0
	b := bridges.BridgeTypeRequest{
		Name: bridges.MustParseBridgeName(name),
		URL:  url,
	}
	payloadb, err := json.Marshal(b)
	helpers.PanicErr(err)
	payload := string(payloadb)

	bridgeActionType := bridgeAction(api, b)
	switch bridgeActionType {
	case shouldCreateBridge:
		fmt.Printf("Creating bridge (%s): %s\n", name, eaURL)
		resp := api.withArg(payload).mustExec(api.methods.CreateBridge)
		resource := mustJSON[presenters.BridgeResource](resp)
		fmt.Printf("Created bridge: %s %s\n", resource.Name, resource.URL)
	case shouldUpdateBridge:
		fmt.Println("Updating existing bridge")
		api.withArgs(name, payload).mustExec(api.methods.UpdateBridge)
		fmt.Println("Updated bridge", name)
	case shouldNoChangeBridge:
		fmt.Println("No changes needed for bridge", name)
	}
}

// create enum for 3 states: create, update, no change
var (
	shouldCreateBridge   = 0
	shouldUpdateBridge   = 1
	shouldNoChangeBridge = 2
)

func bridgeAction(api *nodeAPI, existingBridge bridges.BridgeTypeRequest) int {
	resp, err := api.withArg(existingBridge.Name.String()).exec(api.methods.ShowBridge)
	if err != nil {
		return shouldCreateBridge
	}

	b := mustJSON[presenters.BridgeResource](resp)
	fmt.Printf("Found matching bridge: %s with URL: %s\n", b.Name, b.URL)
	if b.URL == existingBridge.URL.String() {
		return shouldNoChangeBridge
	}
	return shouldUpdateBridge
}

func generateMercuryOCR2Config(nca []NodeKeys) MercuryOCR2Config {
	ctx := context.Background()
	f := uint8(1)
	rawOnchainConfig := mercurytypes.OnchainConfig{
		Min: big.NewInt(0),
		Max: big.NewInt(math.MaxInt64),
	}

	// Values were taken from Data Streams 250ms feeds, given by @austinborn
	rawReportingPluginConfig := datastreamsmercury.OffchainConfig{
		ExpirationWindow: 86400,
		BaseUSDFee:       decimal.NewFromInt(0),
	}

	onchainConfig, err := (datastreamsmercury.StandardOnchainConfigCodec{}).Encode(ctx, rawOnchainConfig)
	helpers.PanicErr(err)
	reportingPluginConfig, err := json.Marshal(rawReportingPluginConfig)
	helpers.PanicErr(err)

	onchainPubKeys := []common.Address{}
	for _, n := range nca {
		onchainPubKeys = append(onchainPubKeys, common.HexToAddress(n.OCR2OnchainPublicKey))
	}

	offchainPubKeysBytes := []ocrtypes.OffchainPublicKey{}
	for _, n := range nca {
		pkBytesFixed := strToBytes32(n.OCR2OffchainPublicKey)
		offchainPubKeysBytes = append(offchainPubKeysBytes, ocrtypes.OffchainPublicKey(pkBytesFixed))
	}

	configPubKeysBytes := []ocrtypes.ConfigEncryptionPublicKey{}
	for _, n := range nca {
		pkBytesFixed := strToBytes32(n.OCR2ConfigPublicKey)
		configPubKeysBytes = append(configPubKeysBytes, ocrtypes.ConfigEncryptionPublicKey(pkBytesFixed))
	}

	identities := []confighelper.OracleIdentityExtra{}
	for index := range nca {
		transmitterAccount := ocrtypes.Account(nca[index].CSAPublicKey)

		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index][:],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].P2PPeerID,
				TransmitAccount:   transmitterAccount,
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	secrets := deployment.XXXGenerateTestOCRSecrets()
	// Values were taken from Data Streams 250ms feeds, given by @austinborn
	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsDeterministic(
		secrets.EphemeralSk,
		secrets.SharedSecret,
		10*time.Second,         // DeltaProgress
		10*time.Second,         // DeltaResend
		400*time.Millisecond,   // DeltaInitial
		5*time.Second,          // DeltaRound
		0,                      // DeltaGrace
		1*time.Second,          // DeltaCertifiedCommitRequest
		0,                      // DeltaStage
		25,                     // rMax
		[]int{len(identities)}, // S
		identities,
		reportingPluginConfig, // reportingPluginConfig []byte,
		nil,                   // maxDurationInitialization *time.Duration,
		0,                     // maxDurationQuery time.Duration,
		250*time.Millisecond,  // Max duration observation
		0,                     // Max duration should accept attested report
		0,                     // Max duration should transmit accepted report
		int(f),                // f
		onchainConfig,
	)
	PanicErr(err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	PanicErr(err)

	offChainTransmitters := make([][32]byte, len(nca))
	for i, n := range nca {
		offChainTransmitters[i] = strToBytes32(n.CSAPublicKey)
	}

	config := MercuryOCR2Config{
		Signers:               signerAddresses,
		Transmitters:          offChainTransmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	return config
}

type MercuryOCR2Config struct {
	Signers               []common.Address
	Transmitters          [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}
