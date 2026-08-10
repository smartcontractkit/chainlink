package src

// This package deploys the LLO (Low Latency Oracle) streams-trigger DON for the
// targeted node set. It replaces the legacy Mercury OCR2 streams-trigger.
//
// It:
//   - Deploys the LLO Configurator and ChannelConfigStore contracts.
//   - Generates an LLO OCR3 config and calls Configurator.SetProductionConfig for
//     the streams-trigger DON.
//   - Builds channel-definitions JSON (one channel per feed, referencing that
//     feed's stream id), writes it to the artefacts dir, computes its sha and
//     registers it on the ChannelConfigStore.
//   - Creates stream jobs (one per stream id) + bridges on every node, one LLO
//     oracle job per non-bootstrap node, and a single bootstrap job.
//
// See the LLO integration/helpers tests for the authoritative examples:
//   - core/services/ocr2/plugins/llo/integration_test.go   (setupBlockchain, setProductionConfig, newChannelDefinitionsServer)
//   - core/services/ocr2/plugins/llo/helpers_test.go        (addBootstrapJob, addLLOJob, stream job specs)
import (
	"bytes"
	"context"
	"crypto/ed25519"
	sha3 "crypto/sha3"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	lloevm "github.com/smartcontractkit/chainlink-evm/pkg/llo"
	evm "github.com/smartcontractkit/chainlink-evm/pkg/relay"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// streamsTriggerDonID is the fixed DON id used for the keystone streams-trigger
// LLO DON. It must be stable across runs since it is baked into the on-chain
// config (Configurator config id), the channel-definitions registration and the
// job specs. Any value works as long as it is used consistently; we pick 1.
const streamsTriggerDonID = uint32(1)

// channelDefinitionsFileName is the artefacts-dir file that the generated
// channel-definitions JSON is written to.
const channelDefinitionsFileName = "channel-definitions.json"

// channelDefinitionsURL is the URL that nodes will fetch the channel-definitions
// JSON from. The keystone mock infra already serves feeds under
// http://external-adapter:400x; we host the channel definitions alongside it.
//
// NOTE: This script only WRITES the channel-definitions JSON to the artefacts
// dir and registers its URL + sha on the ChannelConfigStore. It does NOT stand
// up a server for it. Operators MUST host the generated
// channel-definitions.json (written to the artefacts dir) at this URL, otherwise
// the nodes will be unable to fetch the channel definitions and no reports will
// be produced.
const channelDefinitionsURL = "http://external-adapter:4100/channel-definitions.json"

// channelDefinitionsFromBlock is the block from which nodes replay
// ChannelConfigStore logs. 0 replays from genesis, which is safe for a dev
// deployment.
const channelDefinitionsFromBlock = 0

type feed struct {
	id   [32]byte
	name string

	// streamID is the LLO stream id that this feed's stream job produces and
	// that the feed's channel references.
	streamID uint32

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
		1,
		"mock-bridge-btc",
		"http://external-adapter:4001",
	},
	{
		v3FeedID([32]byte{5: 2}),
		"LINK/USD",
		2,
		"mock-bridge-link",
		"http://external-adapter:4002",
	},
	{
		v3FeedID([32]byte{5: 3}),
		"NATIVE/USD",
		3,
		"mock-bridge-native",
		"http://external-adapter:4003",
	},
}

func setupStreamsTrigger(
	env helpers.Environment,
	nodeSet NodeSet,
	chainID int64,
	p2pPort int64,
	ocrConfigFilePath string,
	artefactsDir string,
) {
	fmt.Printf("Deploying streams trigger (LLO) for chain %d\n", chainID)
	fmt.Printf("Using OCR config file: %s\n", ocrConfigFilePath)

	fmt.Printf("Deploying LLO contracts (Configurator + ChannelConfigStore)\n")
	configuratorAddr, configStoreAddr := deployLLOContracts(env, artefactsDir)

	fmt.Printf("Setting LLO OCR3 config on the Configurator\n")
	ocrConfig := generateLLOOCR3Config(nodeSet.NodeKeys[1:]) // skip the bootstrap node
	donIDPadded := lloevm.DonIDToBytes32(streamsTriggerDonID)
	tx, err := configuratorContract(env, configuratorAddr).SetProductionConfig(
		env.Owner,
		donIDPadded,
		ocrConfig.Signers,
		ocrConfig.OffchainTransmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	)
	PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), env.Ec, tx, env.ChainID)

	fmt.Printf("Building and registering channel definitions\n")
	setChannelDefinitions(env, configStoreAddr, artefactsDir)

	fmt.Printf("Deploying LLO job specs\n")
	deployLLOJobSpecs(nodeSet, configuratorAddr, configStoreAddr, chainID, p2pPort)

	fmt.Println("Finished deploying streams trigger (LLO)")
}

func configuratorContract(env helpers.Environment, addr common.Address) *configurator.Configurator {
	c, err := configurator.NewConfigurator(addr, env.Ec)
	PanicErr(err)
	return c
}

// deployLLOContracts deploys (or reuses) the LLO Configurator and
// ChannelConfigStore contracts and returns their addresses.
func deployLLOContracts(env helpers.Environment, artefactsDir string) (configuratorAddr common.Address, configStoreAddr common.Address) {
	var confirmDeploy = func(tx *types.Transaction, err error) {
		helpers.ConfirmContractDeployed(context.Background(), env.Ec, tx, env.ChainID)
		PanicErr(err)
	}
	o := LoadOnchainMeta(artefactsDir, env)

	if o.Configurator != nil {
		fmt.Printf("Configurator contract already deployed at %s\n", o.Configurator.Address().Hex())
	} else {
		fmt.Printf("Deploying Configurator contract\n")
		_, tx, c, err := configurator.DeployConfigurator(env.Owner, env.Ec)
		confirmDeploy(tx, err)
		o.Configurator = c
		WriteOnchainMeta(o, artefactsDir)
	}

	if o.ChannelConfigStore != nil {
		fmt.Printf("ChannelConfigStore contract already deployed at %s\n", o.ChannelConfigStore.Address().Hex())
	} else {
		fmt.Printf("Deploying ChannelConfigStore contract\n")
		_, tx, ccs, err := channel_config_store.DeployChannelConfigStore(env.Owner, env.Ec)
		confirmDeploy(tx, err)
		o.ChannelConfigStore = ccs
		WriteOnchainMeta(o, artefactsDir)
	}

	return o.Configurator.Address(), o.ChannelConfigStore.Address()
}

// buildChannelDefinitions builds the LLO channel-definitions map from the feeds
// list: one channel per feed, referencing that feed's stream id.
func buildChannelDefinitions() llotypes.ChannelDefinitions {
	defs := make(llotypes.ChannelDefinitions, len(feeds))
	for i, f := range feeds {
		// One channel per feed. The channel id is deterministic (feed index+1).
		channelID := llotypes.ChannelID(i + 1)
		// Legacy Mercury v0.3 compatible report format (benchmark/bid/ask quote).
		opts := llotypes.ChannelOpts(fmt.Appendf(nil,
			`{"baseUSDFee":"0","expirationWindow":86400,"feedId":"0x%x","multiplier":"1000000000000000000"}`,
			f.id,
		))
		defs[channelID] = llotypes.ChannelDefinition{
			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
			Streams: []llotypes.Stream{
				{StreamID: f.streamID, Aggregator: llotypes.AggregatorQuote},
			},
			Opts: opts,
		}
	}
	return defs
}

// setChannelDefinitions builds the channel-definitions JSON, writes it to the
// artefacts dir, computes its sha (matching newChannelDefinitionsServer in the
// LLO tests) and registers the URL + sha on the ChannelConfigStore.
func setChannelDefinitions(env helpers.Environment, configStoreAddr common.Address, artefactsDir string) {
	defs := buildChannelDefinitions()

	// Match newChannelDefinitionsServer: MarshalIndent(defs, "", "  ") then sha3.Sum256.
	channelDefinitionsJSON, err := json.MarshalIndent(defs, "", "  ")
	PanicErr(err)
	sha := sha3.Sum256(channelDefinitionsJSON)

	ensureArtefactsDir(artefactsDir)
	outPath := filepath.Join(artefactsDir, channelDefinitionsFileName)
	err = os.WriteFile(outPath, channelDefinitionsJSON, 0600)
	PanicErr(err)

	// NOTE: The channel-definitions JSON has been written to outPath but is NOT
	// served by this script. Operators MUST host this exact file (byte for byte,
	// so the sha matches) at channelDefinitionsURL, otherwise nodes cannot fetch
	// the channel definitions.
	fmt.Printf("Channel definitions written to %s\n", outPath)
	fmt.Printf("NOTE: host this file at %s (sha: 0x%x)\n", channelDefinitionsURL, sha)

	ccs, err := channel_config_store.NewChannelConfigStore(configStoreAddr, env.Ec)
	PanicErr(err)
	tx, err := ccs.SetChannelDefinitions(env.Owner, streamsTriggerDonID, channelDefinitionsURL, sha)
	PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), env.Ec, tx, env.ChainID)
}

func deployLLOJobSpecs(nodeSet NodeSet, configuratorAddr common.Address, configStoreAddr common.Address, chainID int64, p2pPort int64) {
	// we assign the first node as the bootstrap node
	for i, n := range nodeSet.NodeKeys {
		api := newNodeAPI(nodeSet.Nodes[i])

		// Every node needs the bridges + stream jobs so it can produce
		// observations for each stream.
		for _, f := range feeds {
			upsertBridge(api, f.bridgeName, f.bridgeURL)
		}
		for _, f := range feeds {
			name, spec := createStreamJob(StreamJobSpecData{
				FeedName: f.name,
				StreamID: f.streamID,
				Bridge:   f.bridgeName,
			})
			upsertJob(api, name, spec)
		}

		var jobSpecName, jobSpecStr string
		if i == 0 {
			// Bootstrap job
			jobSpecName, jobSpecStr = createLLOBootstrapJob(LLOBootstrapJobSpecData{
				DonID:                             streamsTriggerDonID,
				ConfiguratorAddress:               configuratorAddr.Hex(),
				ChannelDefinitionsContractAddress: configStoreAddr.Hex(),
				ChannelDefinitionsFromBlock:       channelDefinitionsFromBlock,
				ChainID:                           chainID,
			})
		} else {
			// LLO oracle job
			jobSpecName, jobSpecStr = createLLOJob(LLOJobSpecData{
				DonID:                             streamsTriggerDonID,
				BootstrapHost:                     fmt.Sprintf("%s@%s:%d", nodeSet.NodeKeys[0].P2PPeerID, nodeSet.Nodes[0].ServiceName, p2pPort),
				ConfiguratorAddress:               configuratorAddr.Hex(),
				ChannelDefinitionsContractAddress: configStoreAddr.Hex(),
				ChannelDefinitionsFromBlock:       channelDefinitionsFromBlock,
				NodeCSAKey:                        n.CSAPublicKey,
				OCRKeyBundleID:                    n.OCR2BundleID,
				ChainID:                           chainID,
			})
		}
		upsertJob(api, jobSpecName, jobSpecStr)
	}
}

// Template definitions

// streamJobTemplate produces a single quote stream (benchmark/bid/ask) from the
// feed's bridge. See addQuoteStreamJob in the LLO helpers tests.
const streamJobTemplate = `
type = "stream"
schemaVersion = 1
name = "{{ .Name }}"
streamID = {{ .StreamID }}
observationSource = """
	price              [type=bridge name="{{ .Bridge }}" timeout="50ms" requestData=""];

	benchmark_price  [type=jsonparse path="result,mid" index=0];
	price -> benchmark_price;

	bid_price [type=jsonparse path="result,bid" index=1];
	price -> bid_price;

	ask_price [type=jsonparse path="result,ask" index=2];
	price -> ask_price;
"""
`

// lloBootstrapJobTemplate follows addBootstrapJob in the LLO helpers tests:
// providerType = "llo", contractID = configurator address, relayConfig carries
// the donID and channel definitions contract fields.
const lloBootstrapJobTemplate = `
type                              = "bootstrap"
relay                             = "evm"
schemaVersion                     = 1
name                              = "{{ .Name }}"
contractID                        = "{{ .ConfiguratorAddress }}"
contractConfigTrackerPollInterval = "1s"

[relayConfig]
chainID = {{ .ChainID }}
enableTriggerCapability = true
lloDonID = {{ .DonID }}
providerType = "llo"
channelDefinitionsContractAddress = "{{ .ChannelDefinitionsContractAddress }}"
channelDefinitionsContractFromBlock = {{ .ChannelDefinitionsFromBlock }}
`

// lloJobTemplate follows addLLOJob in the LLO helpers tests: pluginType = "llo",
// no feedID, transmitterID = node CSA pubkey hex, [pluginConfig] with serverURL,
// serverPubKey, donID and channel definitions contract fields.
const lloJobTemplate = `
type = "offchainreporting2"
schemaVersion = 1
name = "{{ .Name }}"
p2pv2Bootstrappers = ["{{ .BootstrapHost }}"]
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "{{ .ConfiguratorAddress }}"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "{{ .OCRKeyBundleID }}"
relay = "evm"
pluginType = "llo"
transmitterID = "{{ .NodeCSAKey }}"

[pluginConfig]
servers = { "{{ .ServerURL }}" = "{{ .ServerPubKey }}" }
donID = {{ .DonID }}
channelDefinitionsContractAddress = "{{ .ChannelDefinitionsContractAddress }}"
channelDefinitionsContractFromBlock = {{ .ChannelDefinitionsFromBlock }}

[relayConfig]
enableTriggerCapability = true
chainID = {{ .ChainID }}
lloDonID = {{ .DonID }}
`

// Data structures

type StreamJobSpecData struct {
	FeedName string
	// Automatically generated from FeedName
	Name     string
	StreamID uint32
	Bridge   string
}

type LLOBootstrapJobSpecData struct {
	// Automatically generated
	Name                              string
	DonID                             uint32
	ConfiguratorAddress               string
	ChannelDefinitionsContractAddress string
	ChannelDefinitionsFromBlock       int64
	ChainID                           int64
}

type LLOJobSpecData struct {
	// Automatically generated
	Name                              string
	DonID                             uint32
	BootstrapHost                     string
	ConfiguratorAddress               string
	ChannelDefinitionsContractAddress string
	ChannelDefinitionsFromBlock       int64
	NodeCSAKey                        string
	OCRKeyBundleID                    string
	ChainID                           int64

	// LLO transmits reports to a data-streams server over the mercury
	// transmitter protocol. These MUST be configured for a real deployment.
	//
	// NOTE: There is no default streams server in the keystone dev infra, so we
	// leave placeholders. Operators MUST set ServerURL / ServerPubKey to a real
	// data-streams server before reports can be transmitted.
	ServerURL    string
	ServerPubKey string
}

const (
	// NOTE: placeholders for the data-streams (mercury) transmitter server. See
	// LLOJobSpecData above. Operators MUST replace these with a real server.
	placeholderServerURL    = "streams-server:1234"
	placeholderServerPubKey = "0000000000000000000000000000000000000000000000000000000000000000"
)

func createStreamJob(data StreamJobSpecData) (name string, jobSpecStr string) {
	name = fmt.Sprintf("stream-spec-%d-%s", data.StreamID, data.FeedName)
	data.Name = name

	fmt.Printf("Creating stream job (%s):\nstream ID: %d\nbridge: %s\n", name, data.StreamID, data.Bridge)

	tmpl, err := template.New("streamJob").Parse(streamJobTemplate)
	PanicErr(err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	PanicErr(err)

	return name, buf.String()
}

// createLLOBootstrapJob creates an LLO bootstrap job specification.
func createLLOBootstrapJob(data LLOBootstrapJobSpecData) (name string, jobSpecStr string) {
	name = fmt.Sprintf("boot-streams-llo-%d", data.DonID)
	data.Name = name

	fmt.Printf("Creating bootstrap job (%s):\nconfigurator address: %s\ndon ID: %d\nchannel definitions contract: %s\nchain ID: %d\n",
		name, data.ConfiguratorAddress, data.DonID, data.ChannelDefinitionsContractAddress, data.ChainID)

	tmpl, err := template.New("lloBootstrapJob").Parse(lloBootstrapJobTemplate)
	PanicErr(err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	PanicErr(err)

	return name, buf.String()
}

// createLLOJob creates an LLO oracle job specification.
func createLLOJob(data LLOJobSpecData) (name string, jobSpecStr string) {
	name = fmt.Sprintf("streams-llo-%d", data.DonID)
	data.Name = name
	if data.ServerURL == "" {
		data.ServerURL = placeholderServerURL
	}
	if data.ServerPubKey == "" {
		data.ServerPubKey = placeholderServerPubKey
	}

	fmt.Printf("Creating LLO job (%s):\nOCR key bundle ID: %s\nconfigurator address: %s\nnodeCSAKey: %s\ndon ID: %d\nchannel definitions contract: %s\nchain ID: %d\n",
		data.Name, data.OCRKeyBundleID, data.ConfiguratorAddress, data.NodeCSAKey, data.DonID, data.ChannelDefinitionsContractAddress, data.ChainID)

	tmpl, err := template.New("lloJob").Parse(lloJobTemplate)
	PanicErr(err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	PanicErr(err)

	return data.Name, buf.String()
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

// LLOOCR3Config holds the arguments for Configurator.SetProductionConfig.
type LLOOCR3Config struct {
	Signers               [][]byte
	OffchainTransmitters  [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func generateLLOOCR3Config(nca []NodeKeys) LLOOCR3Config {
	f := uint8(1)

	// LLO onchain config: production config has no predecessor.
	onchainConfig, err := (&lloprotocol.EVMOnchainConfigCodec{}).Encode(lloprotocol.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: nil,
	})
	helpers.PanicErr(err)

	// LLO reporting plugin (offchain) config.
	reportingPluginConfig, err := lloprotocol.OffchainConfig{
		ProtocolVersion:                     1,
		DefaultMinReportIntervalNanoseconds: uint64(time.Second),
		EnableObservationCompression:        true,
	}.Encode()
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

	secrets := focr.XXXGenerateTestOCRSecrets()
	signers, _, _, onchainConfigOut, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsDeterministic(
		secrets.EphemeralSk,
		secrets.SharedSecret,
		2*time.Second,          // DeltaProgress
		20*time.Second,         // DeltaResend
		400*time.Millisecond,   // DeltaInitial
		500*time.Millisecond,   // DeltaRound
		250*time.Millisecond,   // DeltaGrace
		300*time.Millisecond,   // DeltaCertifiedCommitRequest
		1*time.Minute,          // DeltaStage
		100,                    // rMax
		[]int{len(identities)}, // S
		identities,
		reportingPluginConfig, // reportingPluginConfig []byte
		nil,                   // maxDurationInitialization *time.Duration
		0,                     // maxDurationQuery time.Duration
		250*time.Millisecond,  // maxDurationObservation
		0,                     // maxDurationShouldAcceptAttestedReport
		0,                     // maxDurationShouldTransmitAcceptedReport
		int(f),                // f
		onchainConfig,
	)
	PanicErr(err)

	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	PanicErr(err)
	// Configurator.SetProductionConfig expects signers as [][]byte (20-byte addresses).
	onchainSigners := make([][]byte, len(signerAddresses))
	for i, addr := range signerAddresses {
		onchainSigners[i] = addr.Bytes()
	}

	offchainTransmitters := make([][32]byte, len(nca))
	for i, n := range nca {
		offchainTransmitters[i] = strToBytes32(n.CSAPublicKey)
	}

	return LLOOCR3Config{
		Signers:               onchainSigners,
		OffchainTransmitters:  offchainTransmitters,
		F:                     f,
		OnchainConfig:         onchainConfigOut,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}
