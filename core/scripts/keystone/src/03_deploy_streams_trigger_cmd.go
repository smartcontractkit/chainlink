package src

// This package deploys "offchainreporting2" job specs, which setup the streams trigger
// for the targetted node set
// See https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/integration_test.go#L92
// for how to setup the mercury portion of the streams trigger
//  You can see how all fields are being used here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/helpers_test.go#L314
//  https://github.com/smartcontractkit/infra-k8s/blob/be47098adfb605d79b5bab6aa601bcf443a6c48b/projects/chainlink/files/chainlink-clusters/cl-keystone-cap-one/config.yaml#L1
//  Trigger gets added to the registry here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/relay/evm/evm.go#L360
//  See integration workflow here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/capabilities/integration_tests/workflow.go#L15
//  ^ setup.go provides good insight too
import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"math"
	"math/big"
	"os"
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
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

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
	bridgeUrl  string
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

type deployStreamsTrigger struct{}

func NewDeployStreamsTriggerCommand() *deployStreamsTrigger {
	return &deployStreamsTrigger{}
}

func (g *deployStreamsTrigger) Name() string {
	return "deploy-streams-trigger"
}

func (g *deployStreamsTrigger) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 1337, "chain id")
	ocrConfigFile := fs.String("ocrfile", "ocr_config.json", "path to OCR config file")
	nodeList := fs.String("nodes", "", "Custom node list location")
	publicKeys := fs.String("publickeys", "", "Custom public keys json location")
	force := fs.Bool("force", false, "Force deployment")

	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")

	err := fs.Parse(args)
	if err != nil ||
		*ocrConfigFile == "" || ocrConfigFile == nil ||
		chainID == nil || *chainID == 0 ||
		*ethUrl == "" || ethUrl == nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}

	if *publicKeys == "" {
		*publicKeys = defaultPublicKeys
	}
	if *nodeList == "" {
		*nodeList = defaultNodeList
	}

	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)
	os.Setenv("INSECURE_SKIP_VERIFY", "true")

	env := helpers.SetupEnv(false)

	setupMercuryV03(
		env,
		*nodeList,
		*ocrConfigFile,
		*chainID,
		*publicKeys,
		*force,
	)
}

// See /core/services/ocr2/plugins/mercury/integration_test.go
func setupMercuryV03(env helpers.Environment, nodeListPath string, ocrConfigFilePath string, chainId int64, pubKeysPath string, force bool) {
	fmt.Printf("Deploying streams trigger for chain %d\n", chainId)
	fmt.Printf("Using OCR config file: %s\n", ocrConfigFilePath)
	fmt.Printf("Using node list: %s\n", nodeListPath)
	fmt.Printf("Using public keys: %s\n", pubKeysPath)
	fmt.Printf("Force: %t\n\n", force)

	fmt.Printf("Deploying Mercury V0.3 contracts\n")
	_, _, _, verifier := deployMercuryV03Contracts(env)
	// the 0th index is for the OCR3 capability
	// where the 1st index is for the mercury OCR2 instance
	kbIndex := 1
	nca := downloadNodePubKeys(nodeListPath, chainId, pubKeysPath, kbIndex)
	nodes := downloadNodeAPICredentials(nodeListPath)

	fmt.Printf("Generating OCR3 config\n")
	ocrConfig := generateMercuryOCR2Config(nca)

	for _, feed := range feeds {
		fmt.Println("Configuring feeds...")
		fmt.Printf("FeedID: %x\n", feed.id)
		fmt.Printf("FeedName: %s\n", feed.name)
		fmt.Printf("BridgeName: %s\n", feed.bridgeName)
		fmt.Printf("BridgeURL: %s\n", feed.bridgeUrl)

		fmt.Printf("Setting verifier config\n")
		fmt.Printf("Signers: %v\n", ocrConfig.Signers)
		fmt.Printf("Transmitters: %v\n", ocrConfig.Transmitters)
		fmt.Printf("F: %d\n", ocrConfig.F)
		
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

		fmt.Printf("Deploying OCR2 job specs for feed %s\n", feed.name)
		deployOCR2JobSpecsForFeed(nca, nodes, verifier, feed, chainId, force)
	}

	fmt.Println("Finished deploying streams trigger")

	fmt.Println("Deploying Keystone workflow job")

	feedIds := []string{}
	for _, feed := range feeds {
		feedIds = append(feedIds, fmt.Sprintf("0x%x", feed.id))
	}
	writeTarget := NewEthereumGethTestnetV1WriteCapability().baseCapability.capability
	writeTargetID := fmt.Sprintf("%s@%s", writeTarget.LabelledName, writeTarget.Version)
	workflowConfig := WorkflowJobSpecConfig{
		JobSpecName:          "keystone_workflow",
		WorkflowOwnerAddress: "0x1234567890abcdef1234567890abcdef12345678",
		FeedIDs:              feedIds,
		TargetID:             writeTargetID,
		TargetAddress:        "0x1234567890abcdef1234567890abcdef12345678",
	}
	jobSpecStr := createKeystoneWorkflowJob(workflowConfig)
	for _, n := range nodes {
		api := newNodeAPI(n)

		upsertJob(api, workflowConfig.JobSpecName, jobSpecStr, force)
	}
}

func deployMercuryV03Contracts(env helpers.Environment) (linkToken *link_token_interface.LinkToken, nativeToken *link_token_interface.LinkToken, verifierProxy *verifier_proxy.VerifierProxy, verifier *verifierContract.Verifier) {
	var confirmDeploy = func(tx *types.Transaction, err error) {
		helpers.ConfirmContractDeployed(context.Background(), env.Ec, tx, env.ChainID)
		PanicErr(err)
	}
	var confirmTx = func(tx *types.Transaction, err error) {
		helpers.ConfirmTXMined(context.Background(), env.Ec, tx, env.ChainID)
		PanicErr(err)
	}

	_, tx, linkToken, err := link_token_interface.DeployLinkToken(env.Owner, env.Ec)
	confirmDeploy(tx, err)

	// Not sure if we actually need to have tokens
	tx, err = linkToken.Transfer(env.Owner, env.Owner.From, big.NewInt(1000))
	confirmTx(tx, err)

	// We reuse the link token for the native token
	_, tx, nativeToken, err = link_token_interface.DeployLinkToken(env.Owner, env.Ec)
	confirmDeploy(tx, err)

	// Not sure if we actually need to have tokens
	tx, err = nativeToken.Transfer(env.Owner, env.Owner.From, big.NewInt(1000))
	confirmTx(tx, err)

	verifierProxyAddr, tx, verifierProxy, err := verifier_proxy.DeployVerifierProxy(env.Owner, env.Ec, common.Address{}) // zero address for access controller disables access control
	confirmDeploy(tx, err)

	verifierAddress, tx, verifier, err := verifierContract.DeployVerifier(env.Owner, env.Ec, verifierProxyAddr)
	confirmDeploy(tx, err)

	tx, err = verifierProxy.InitializeVerifier(env.Owner, verifierAddress)
	confirmTx(tx, err)

	return
}

func deployOCR2JobSpecsForFeed(nca []NodeKeys, nodes []*node, verifier *verifierContract.Verifier, feed feed, chainId int64, force bool) {
	// we assign the first node as the bootstrap node
	for i, n := range nca {
		// parallel arrays
		api := newNodeAPI(nodes[i])
		jobSpecName := ""
		jobSpecStr := ""

		createBridgeIfDoesNotExist(api, feed.bridgeName, feed.bridgeUrl, force)

		if i == 0 {
			// Prepare data for Bootstrap Job
			bootstrapData := MercuryV3BootstrapJobSpecData{
				FeedName:        feed.name,
				VerifierAddress: verifier.Address().Hex(),
				FeedID:          fmt.Sprintf("%x", feed.id),
				ChainID:         chainId,
			}

			// Create Bootstrap Job
			jobSpecName, jobSpecStr = createMercuryV3BootstrapJob(bootstrapData)
		} else {
			// Prepare data for Mercury V3 Job
			mercuryData := MercuryV3JobSpecData{
				FeedName:        fmt.Sprintf("feed-%s", feed.name),
				BootstrapHost:   fmt.Sprintf("%s@%s:%s", nca[0].P2PPeerID, "app-node1", "6690"),
				VerifierAddress: verifier.Address().Hex(),
				Bridge:          feed.bridgeName,
				NodeCSAKey:      n.CSAPublicKey,
				FeedID:          fmt.Sprintf("%x", feed.id),
				LinkFeedID:      fmt.Sprintf("%x", feeds[1].id),
				NativeFeedID:    fmt.Sprintf("%x", feeds[2].id),
				OCRKeyBundleID:  n.OCR2BundleID,
				ChainID:         chainId,
			}

			// Create Mercury V3 Job
			jobSpecName, jobSpecStr = createMercuryV3Job(mercuryData)
		}

		upsertJob(api, jobSpecName, jobSpecStr, force)
	}
}

func upsertJob(api *nodeAPI, jobSpecName string, jobSpecStr string, upsert bool) {
		jobsResp := api.mustExec(api.methods.ListJobs)
		jobs := mustJSON[[]JobSpec](jobsResp)
		for _, job := range *jobs {
			if job.Name == jobSpecName {
			if !upsert {
				fmt.Printf("Job already exists: %s, skipping..\n", jobSpecName)
				return
			}

					fmt.Printf("Job already exists: %s, replacing..\n", jobSpecName)
					api.withArg(job.Id).mustExec(api.methods.DeleteJob)
					fmt.Printf("Deleted job: %s\n", jobSpecName)
			break
			}
		}

		fmt.Printf("Deploying jobspec: %s\n... \n", jobSpecStr)
		_, err := api.withArg(jobSpecStr).exec(api.methods.CreateJob)
		if err != nil {
			panic(fmt.Sprintf("Failed to deploy job spec: %s Error: %s", jobSpecStr, err))
		}
	}

type WorkflowJobSpecConfig struct {
	JobSpecName          string
	WorkflowOwnerAddress string
	FeedIDs              []string
	TargetID             string
	TargetAddress        string
}

func createKeystoneWorkflowJob(workflowConfig WorkflowJobSpecConfig) string {
	const keystoneWorkflowTemplate = `
type = "workflow"
schemaVersion = 1
name = "{{ .JobSpecName }}"
workflow = """
name: "ccip_kiab" 
owner: '{{ .WorkflowOwnerAddress }}'
triggers:
 - id: streams-trigger@1.0.0
   config:
     maxFrequencyMs: 10000
     feedIds:
{{- range .FeedIDs }}
       - '{{ . }}'
{{- end }}

consensus:
 - id: offchain_reporting@1.0.0
   ref: ccip_feeds
   inputs:
     observations:
       - $(trigger.outputs)
   config:
     report_id: '0001'
     aggregation_method: data_feeds
     aggregation_config:
{{- range .FeedIDs }}
         '{{ . }}':
           deviation: '0.05'
           heartbeat: 1800
{{- end }}
     encoder: EVM
     encoder_config:
       abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports

targets:
 - id: {{ .TargetID }} 
   inputs:
     signed_report: $(ccip_feeds.outputs)
   config:
     address: '{{ .TargetAddress }}'
     deltaStage: 5s
     schedule: oneAtATime

"""
workflowOwner = "{{ .WorkflowOwnerAddress }}"
`

	tmpl, err := template.New("workflow").Parse(keystoneWorkflowTemplate)

	if err != nil {
		panic(err)
	}
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, workflowConfig)
	if err != nil {
		panic(err)
	}

	return renderedTemplate.String()
}

// Template definitions
const bootstrapJobTemplate = `
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

const mercuryV3JobTemplate = `
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
	name = fmt.Sprintf("boot-%s", data.FeedName)
	data.Name = name

	fmt.Printf("Creating bootstrap job (%s):\nverifier address: %s\nfeed name: %s\nfeed ID: %s\nchain ID: %d\n",
		name, data.VerifierAddress, data.FeedName, data.FeedID, data.ChainID)

	tmpl, err := template.New("bootstrapJob").Parse(bootstrapJobTemplate)
	PanicErr(err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	PanicErr(err)

	jobSpecStr = buf.String()

	return name, jobSpecStr
}

// createMercuryV3Job creates a Mercury V3 job specification using the provided data.
func createMercuryV3Job(data MercuryV3JobSpecData) (name string, jobSpecStr string) {
	name = fmt.Sprintf("mercury-%s", data.FeedName)
	data.Name = name
	fmt.Printf("Creating ocr2 job(%s):\nOCR key bundle ID: %s\nverifier address: %s\nbridge: %s\nnodeCSAKey: %s\nfeed name: %s\nfeed ID: %s\nlink feed ID: %s\nnative feed ID: %s\nchain ID: %d\n",
		data.Name, data.OCRKeyBundleID, data.VerifierAddress, data.Bridge, data.NodeCSAKey, data.FeedName, data.FeedID, data.LinkFeedID, data.NativeFeedID, data.ChainID)

	tmpl, err := template.New("mercuryV3Job").Parse(mercuryV3JobTemplate)
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

func createBridgeIfDoesNotExist(api *nodeAPI, name string, eaURL string, force bool) {
	u, err := url.Parse(eaURL)
	url := models.WebURL(*u)
	// Confirmations and MinimumContractPayment are not used, so we can leave them as 0
	b := bridges.BridgeTypeRequest{
		Name: bridges.MustParseBridgeName(name),
		URL:  url,
	}
	payloadb, err := json.Marshal(b)
	helpers.PanicErr(err)
	payload := string(payloadb)

	fmt.Printf("Creating bridge (%s): %s\n", name, eaURL)
	if doesBridgeExist(api, name) {
		if force {
			fmt.Println("Force flag is set, updating existing bridge")
			api.withArgs(name, payload).mustExec(api.methods.UpdateBridge)
			fmt.Println("Updated bridge", name)
		} else {
			fmt.Println("Bridge", name, "already exists, skipping creation")
			return
		}
	} else {
		resp := api.withArg(payload).mustExec(api.methods.CreateBridge)
		resource := mustJSON[presenters.BridgeResource](resp)
		fmt.Printf("Created bridge: %s %s\n", resource.Name, resource.URL)
	}
}

func doesBridgeExist(api *nodeAPI, name string) bool {
	resp, err := api.withArg(name).exec(api.methods.ShowBridge)

	if err != nil {
		return false
	}

	b := mustJSON[presenters.BridgeResource](resp)
	fmt.Printf("Found bridge: %s with URL: %s\n", b.Name, b.URL)
	return true
}

func generateMercuryOCR2Config(nca []NodeKeys) MercuryOCR2Config {
	f := uint8(1)
	rawOnchainConfig := mercurytypes.OnchainConfig{
		Min: big.NewInt(0),
		Max: big.NewInt(math.MaxInt64),
	}
	rawReportingPluginConfig := datastreamsmercury.OffchainConfig{
		ExpirationWindow: 1,
		BaseUSDFee:       decimal.NewFromInt(100),
	}

	onchainConfig, err := (datastreamsmercury.StandardOnchainConfigCodec{}).Encode(rawOnchainConfig)
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
		transmitterAccount := ocrtypes.Account(fmt.Sprintf("%x", nca[index].CSAPublicKey[:]))

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

	signers, _, _, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTestsMercuryV02(
		2*time.Second,          // DeltaProgress
		20*time.Second,         // DeltaResend
		400*time.Millisecond,   // DeltaInitial
		100*time.Millisecond,   // DeltaRound
		0,                      // DeltaGrace
		300*time.Millisecond,   // DeltaCertifiedCommitRequest
		1*time.Minute,          // DeltaStage
		100,                    // rMax
		[]int{len(identities)}, // S
		identities,
		reportingPluginConfig, // reportingPluginConfig []byte,
		250*time.Millisecond,  // Max duration observation
		int(f),                // f
		onchainConfig,
	)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	PanicErr(err)

	var offChainTransmitters [][32]byte
	for _, n := range nca {
		offChainTransmitters = append(offChainTransmitters, strToBytes32(n.CSAPublicKey))
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
