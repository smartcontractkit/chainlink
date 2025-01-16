package capabilities_test

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/capabilities/components/evmcontracts/capabilities_registry"
	"github.com/smartcontractkit/chainlink/integration-tests/capabilities/components/evmcontracts/forwarder"
	"github.com/smartcontractkit/chainlink/integration-tests/capabilities/components/onchain"

	feeds_consumer_debug "github.com/smartcontractkit/chainlink/integration-tests/capabilities/components/evmcontracts/feed_consumer_debug"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"

	cr_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
	workflow_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

// Copying this to avoid dependency on the core repo
func GetChainType(chainType string) (uint8, error) {
	switch chainType {
	case "evm":
		return 1, nil
	// case Solana:
	// 	return 2, nil
	// case Cosmos:
	// 	return 3, nil
	// case StarkNet:
	// 	return 4, nil
	// case Aptos:
	// 	return 5, nil
	default:
		return 0, fmt.Errorf("unexpected chaintype.ChainType: %#v", chainType)
	}
}

// Copying this to avoid dependency on the core repo
func MarshalMultichainPublicKey(ost map[string]types.OnchainPublicKey) (types.OnchainPublicKey, error) {
	pubKeys := make([][]byte, 0, len(ost))
	for k, pubKey := range ost {
		typ, err := GetChainType(k)
		if err != nil {
			// skipping unknown key type
			continue
		}
		buf := new(bytes.Buffer)
		if err = binary.Write(buf, binary.LittleEndian, typ); err != nil {
			return nil, err
		}
		length := len(pubKey)
		if length < 0 || length > math.MaxUint16 {
			return nil, errors.New("pubKey doesn't fit into uint16")
		}
		if err = binary.Write(buf, binary.LittleEndian, uint16(length)); err != nil {
			return nil, err
		}
		_, _ = buf.Write(pubKey)
		pubKeys = append(pubKeys, buf.Bytes())
	}
	// sort keys based on encoded type to make encoding deterministic
	slices.SortFunc(pubKeys, func(a, b []byte) int { return cmp.Compare(a[0], b[0]) })
	return bytes.Join(pubKeys, nil), nil
}

type WorkflowTestConfig struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

type OCR3Config struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type NodeInfo struct {
	OcrKeyBundleID            string
	TransmitterAddress        string
	PeerID                    string
	Signer                    common.Address
	OffchainPublicKey         [32]byte
	OnchainPublicKey          types.OnchainPublicKey
	ConfigEncryptionPublicKey [32]byte
}

func extractKey(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func getNodesInfo(
	t *testing.T,
	nodes []*clclient.ChainlinkClient,
) (nodesInfo []NodeInfo) {
	nodesInfo = make([]NodeInfo, len(nodes))

	for i, node := range nodes {
		// OCR Keys
		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err)
		nodesInfo[i].OcrKeyBundleID = ocr2Keys.Data[0].ID

		firstOCR2Key := ocr2Keys.Data[0].Attributes
		nodesInfo[i].Signer = common.HexToAddress(extractKey(firstOCR2Key.OnChainPublicKey))

		pubKeys := make(map[string]types.OnchainPublicKey)
		ethOnchainPubKey, err := hex.DecodeString(extractKey(firstOCR2Key.OnChainPublicKey))
		require.NoError(t, err)
		pubKeys["evm"] = ethOnchainPubKey

		multichainPubKey, err := MarshalMultichainPublicKey(pubKeys)
		require.NoError(t, err)
		nodesInfo[i].OnchainPublicKey = multichainPubKey

		offchainPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.OffChainPublicKey))
		require.NoError(t, err)
		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], offchainPublicKeyBytes)
		nodesInfo[i].OffchainPublicKey = offchainPublicKey

		sharedSecretEncryptionPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.ConfigPublicKey))
		require.NoError(t, err)
		var sharedSecretEncryptionPublicKey [32]byte
		copy(sharedSecretEncryptionPublicKey[:], sharedSecretEncryptionPublicKeyBytes)
		nodesInfo[i].ConfigEncryptionPublicKey = sharedSecretEncryptionPublicKey

		// ETH Keys
		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		nodesInfo[i].TransmitterAddress = ethKeys.Data[0].Attributes.Address

		// P2P Keys
		p2pKeys, err := node.MustReadP2PKeys()
		require.NoError(t, err)
		nodesInfo[i].PeerID = p2pKeys.Data[0].Attributes.PeerID
	}

	return nodesInfo
}

func generateOCR3Config(
	t *testing.T,
	nodesInfo []NodeInfo,
) (config *OCR3Config) {
	oracleIdentities := []confighelper.OracleIdentityExtra{}
	transmissionSchedule := []int{}

	for _, nodeInfo := range nodesInfo {
		transmissionSchedule = append(transmissionSchedule, 1)
		oracleIdentity := confighelper.OracleIdentityExtra{}
		oracleIdentity.OffchainPublicKey = nodeInfo.OffchainPublicKey
		oracleIdentity.OnchainPublicKey = nodeInfo.OnchainPublicKey
		oracleIdentity.ConfigEncryptionPublicKey = nodeInfo.ConfigEncryptionPublicKey
		oracleIdentity.PeerID = nodeInfo.PeerID
		oracleIdentity.TransmitAccount = types.Account(nodeInfo.TransmitterAddress)
		oracleIdentities = append(oracleIdentities, oracleIdentity)
	}

	maxDurationInitialization := 10 * time.Second

	// Generate OCR3 configuration arguments for testing
	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		20*time.Second,             // DeltaProgress: Time between rounds
		10*time.Second,             // DeltaResend: Time between resending unconfirmed transmissions
		1*time.Second,              // DeltaInitial: Initial delay before starting the first round
		5*time.Second,              // DeltaRound: Time between rounds within an epoch
		1*time.Second,              // DeltaGrace: Grace period for delayed transmissions
		5*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
		10*time.Second,             // DeltaStage: Time between stages of the protocol
		uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
		transmissionSchedule,       // TransmissionSchedule: Transmission schedule
		oracleIdentities,           // Oracle identities with their public keys
		nil,                        // Plugin config (empty for now)
		&maxDurationInitialization, // MaxDurationInitialization: ???
		5*time.Second,              // MaxDurationQuery: Maximum duration for querying
		5*time.Second,              // MaxDurationObservation: Maximum duration for observation
		5*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
		5*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
		1,                          // F: Maximum number of faulty oracles
		nil,                        // OnChain config (empty for now)
	)
	require.NoError(t, err)

	// // values supplied by Alexandr Y
	// signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
	// 	5*time.Second,              // DeltaProgress: Time between rounds
	// 	5*time.Second,              // DeltaResend: Time between resending unconfirmed transmissions
	// 	5*time.Second,              // DeltaInitial: Initial delay before starting the first round
	// 	2*time.Second,              // DeltaRound: Time between rounds within an epoch
	// 	500*time.Millisecond,       // DeltaGrace: Grace period for delayed transmissions
	// 	1*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
	// 	30*time.Second,             // DeltaStage: Time between stages of the protocol
	// 	uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
	// 	transmissionSchedule,       // TransmissionSchedule: Transmission schedule
	// 	oracleIdentities,           // Oracle identities with their public keys
	// 	nil,                        // Plugin config (empty for now)
	// 	&maxDurationInitialization, // MaxDurationInitialization: ???
	// 	1*time.Second,              // MaxDurationQuery: Maximum duration for querying
	// 	1*time.Second,              // MaxDurationObservation: Maximum duration for observation
	// 	1*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
	// 	1*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
	// 	1,                          // F: Maximum number of faulty oracles
	// 	nil,                        // OnChain config (empty for now)
	// )
	// require.NoError(t, err)

	signerAddresses := [][]byte{}
	for _, signer := range signers {
		signerAddresses = append(signerAddresses, signer)
	}

	transmitterAddresses := []common.Address{}
	for _, transmitter := range transmitters {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(transmitter)))
	}

	return &OCR3Config{
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

func GenerateWorkflowIDFromStrings(owner string, name string, workflow []byte, config []byte, secretsURL string) (string, error) {
	ownerWithoutPrefix := owner
	if strings.HasPrefix(owner, "0x") {
		ownerWithoutPrefix = owner[2:]
	}

	ownerb, err := hex.DecodeString(ownerWithoutPrefix)
	if err != nil {
		return "", err
	}

	wid, err := GenerateWorkflowID(ownerb, name, workflow, config, secretsURL)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(wid[:]), nil
}

var (
	versionByte = byte(0)
)

func GenerateWorkflowID(owner []byte, name string, workflow []byte, config []byte, secretsURL string) ([32]byte, error) {
	s := sha256.New()
	_, err := s.Write(owner)
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write([]byte(name))
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write(workflow)
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write([]byte(config))
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write([]byte(secretsURL))
	if err != nil {
		return [32]byte{}, err
	}

	sha := [32]byte(s.Sum(nil))
	sha[0] = versionByte

	return sha, nil
}

func TestWorkflow(t *testing.T) {
	// workflowOwner := "0x00000000000000000000000000000000000000aa"
	// without 0x prefix!
	feedID := "018bfe8840700040000000000000000000000000000000000000000000000000"
	feedBytes := common.HexToHash(feedID)

	t.Run("Keystoen workflow test", func(t *testing.T) {
		in, err := framework.Load[WorkflowTestConfig](t)
		require.NoError(t, err)
		pkey := os.Getenv("PRIVATE_KEY")

		bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
		require.NoError(t, err)

		sc, err := seth.NewClientBuilder().
			WithRpcUrl(bc.Nodes[0].HostWSUrl).
			WithPrivateKeys([]string{pkey}).
			Build()
		require.NoError(t, err)

		capabilitiesRegistryInstance, err := capabilities_registry.Deploy(sc)
		require.NoError(t, err)

		require.NoError(t, capabilitiesRegistryInstance.AddCapabilities(
			[]cr_wrapper.CapabilitiesRegistryCapability{
				{
					LabelledName:   "mock-streams-trigger",
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
					ResponseType:   0, // REPORT
				},
				{
					LabelledName:   "offchain_reporting",
					Version:        "1.0.0",
					CapabilityType: 2, // CONSENSUS
					ResponseType:   0, // REPORT
				},
				{
					LabelledName:   "write_geth-testnet",
					Version:        "1.0.0",
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				{
					LabelledName:   "cron-trigger",
					Version:        "1.0.0",
					CapabilityType: uint8(0), // trigger
				},
				{
					LabelledName:   "custom-compute",
					Version:        "1.0.0",
					CapabilityType: uint8(1), // action
				},
			},
		))

		forwarderInstance, err := forwarder.Deploy(sc)
		fmt.Printf("Deployed forawrder contract: %s", forwarderInstance.Address.String())
		require.NoError(t, err)

		workflowRegistryAddr, tx, workflow_registryInstance, err := workflow_registry.DeployWorkflowRegistry(sc.NewTXOpts(), sc.Client)
		_, decodeErr := sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		_ = workflowRegistryAddr

		workFlowData, err := os.ReadFile("./binary.wasm.br")
		require.NoError(t, err)

		var configData []byte

		workflowName := "PoR"
		donID := uint32(1)
		workflowID, idErr := GenerateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, "")
		require.NoError(t, idErr)

		// 00327e3b13a0f5e4da1a9b980e33fde2602043a06289aa6d4eb65af733cb3be6 is workflowID generated from the above function
		// 2025-01-14 15:11:04 2025-01-14T14:11:04.429Z [ERROR] failed to handle workflow registration: workflowID mismatch: 00b8a8f28d3a29f73c1a052273d4b5ed0951e70eb4fae56f7531539af4aca96f != 00327e3b13a0f5e4da1a9b980e33fde2602043a06289aa6d4eb65af733cb3be6 syncer/workflow_registry.go:485  logger=WorkflowRegistrySyncer stacktrace=github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer.(*workflowRegistry).loadWorkflows

		// workflowID = "00b8a8f28d3a29f73c1a052273d4b5ed0951e70eb4fae56f7531539af4aca96f"
		abi, err := workflow_registry.WorkflowRegistryMetaData.GetAbi()
		require.NoError(t, err)
		sc.ContractStore.AddABI("WorkflowRegistry", *abi)

		allowTx, allowErr := workflow_registryInstance.UpdateAllowedDONs(sc.NewTXOpts(), []uint32{donID}, true)
		_, decodeErr = sc.Decode(allowTx, allowErr)
		require.NoError(t, decodeErr)

		allowAddrTx, allowAddrErr := workflow_registryInstance.UpdateAuthorizedAddresses(sc.NewTXOpts(), []common.Address{sc.MustGetRootKeyAddress()}, true)
		_, decodeErr = sc.Decode(allowAddrTx, allowAddrErr)
		require.NoError(t, decodeErr)

		wrTx, wrErr := workflow_registryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), "https://gist.githubusercontent.com/Tofel/8105d6b8289c253d67d7d0abc60a01d1/raw/a376a30b32b60b8188914ebd092a73f7faec1519/binary.wasm.br", "", "")
		_, decodeErr = sc.Decode(wrTx, wrErr)
		require.NoError(t, decodeErr)

		// TODO: When the capabilities registry address is provided:
		// - NOPs and nodes are added to the registry.
		// - Nodes are configured to listen to the registry for updates.
		nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
		require.NoError(t, err)

		nodeClients, err := clclient.New(nodeset.CLNodes)
		require.NoError(t, err)

		err = onchain.FundNodes(sc, nodeClients, pkey, 5)
		require.NoError(t, err)

		nodesInfo := getNodesInfo(t, nodeClients)

		bootstrapNodeInfo := nodesInfo[0]
		workflowNodesetInfo := nodesInfo[1:]

		// bootstrap node
		// changed: use docker container name instead of 'localhost' for default bootstrappers url
		in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				DefaultBootstrappers = ['%s@localhost:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				DefaultBootstrappers = ['%s@localhost:6690']

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'
			`,
			bootstrapNodeInfo.PeerID,
			bootstrapNodeInfo.PeerID,
			bc.ChainID,
			bc.Nodes[0].DockerInternalWSUrl,
			bc.Nodes[0].DockerInternalHTTPUrl,
		)

		for i := range workflowNodesetInfo {
			in.NodeSet.NodeSpecs[i+1].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:6690']

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'

				[EVM.Workflow]
				FromAddress = '%s'
				ForwarderAddress = '%s'
				GasLimitDefault = 400_000

				[Capabilities.ExternalRegistry]
				Address = '%s'
				NetworkID = 'evm'
				ChainID = '%s'

				[Capabilities.WorkflowRegistry]
				Address = "%s"
				NetworkID = "evm"
				ChainID = "%s"

				[Capabilities.GatewayConnector]
				DonID = "1"
				ChainIDForNodeKey = "%s"
				NodeAddress = '%s'

				[[Capabilities.GatewayConnector.Gateways]]
				Id = "por_gateway"
				URL = "%s"
			`,
				bootstrapNodeInfo.PeerID,
				bootstrapNodeInfo.PeerID,
				bc.ChainID,
				bc.Nodes[0].DockerInternalWSUrl,
				bc.Nodes[0].DockerInternalHTTPUrl,
				workflowNodesetInfo[i].TransmitterAddress,
				forwarderInstance.Address,
				capabilitiesRegistryInstance.Address,
				bc.ChainID,
				workflowRegistryAddr.Hex(),
				bc.ChainID,
				bc.ChainID,
				workflowNodesetInfo[i].TransmitterAddress,
				// assuming that node0 is the bootstrap node
				"ws://node0:5003/node",
			)
		}

		nodeset, err = ns.UpgradeNodeSet(in.NodeSet, bc, 5*time.Second)
		require.NoError(t, err)
		nodeClients, err = clclient.New(nodeset.CLNodes)
		require.NoError(t, err)

		ocr3CapabilityAddress, tx, ocr3CapabilityContract, err := ocr3_capability.DeployOCR3Capability(
			sc.NewTXOpts(),
			sc.Client,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)
		fmt.Println("Deployed ocr3_capability contract at", ocr3CapabilityAddress.Hex())

		_ = feeds_consumer_debug.DeployFeedsConsumerDebug

		// feedsConsumerDebugAddress, tx, feedsConsumerDebugContract, err := feeds_consumer_debug.DeployFeedsConsumerDebug(
		// 	sc.NewTXOpts(),
		// 	sc.Client,
		// )
		// require.NoError(t, err)
		// _, err = bind.WaitMined(context.Background(), sc.Client, tx)
		// require.NoError(t, err)

		// _ = feedsConsumerDebugAddress
		// _ = feedsConsumerDebugContract

		feedsConsumerAddress, tx, feedsConsumerContract, err := feeds_consumer.DeployKeystoneFeedsConsumer(
			sc.NewTXOpts(),
			sc.Client,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// just wait for contracts to be indexed by Blockscout and then verify
		time.Sleep(20 * time.Second)

		err = blockchain.VerifyContract(bc, feedsConsumerAddress.String(),
			"../../contracts",
			"src/v0.8/keystone/KeystoneFeedsConsumer.sol",
			"KeystoneFeedsConsumer",
		)
		require.NoError(t, err)
		err = blockchain.VerifyContract(bc, forwarderInstance.Address.String(),
			"../../contracts",
			"src/v0.8/keystone/KeystoneForwarder.sol",
			"KeystoneForwarder",
		)
		require.NoError(t, err)

		fmt.Println("Deployed feeds_consumer contract at", feedsConsumerAddress.Hex())

		var workflowNameBytes [10]byte
		copy(workflowNameBytes[:], []byte(workflowName))

		tx, err = feedsConsumerContract.SetConfig(
			sc.NewTXOpts(),
			[]common.Address{forwarderInstance.Address},
			[]common.Address{sc.MustGetRootKeyAddress()},
			[][10]byte{workflowNameBytes},
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// Add bootstrap spec to the last node
		bootstrapNode := nodeClients[0]

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			gatewayJobSpec := fmt.Sprintf(`
				type = "gateway"
				schemaVersion = 1
				name = "PoR Gateway"
				forwardingAllowed = false

				[gatewayConfig.ConnectionManagerConfig]
				AuthChallengeLen = 10
				AuthGatewayId = "por_gateway"
				AuthTimestampToleranceSec = 5
				HeartbeatIntervalSec = 20

				[[gatewayConfig.Dons]]
				DonId = "1"
				F = 1
				HandlerName = "web-api-capabilities"
					[gatewayConfig.Dons.HandlerConfig]
					MaxAllowedMessageAgeSec = 1_000

						[gatewayConfig.Dons.HandlerConfig.NodeRateLimiter]
						GlobalBurst = 10
						GlobalRPS = 50
						PerSenderBurst = 10
						PerSenderRPS = 10

					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 1"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 2"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 3"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 4"

				[gatewayConfig.NodeServerConfig]
				HandshakeTimeoutMillis = 1_000
				MaxRequestBytes = 100_000
				Path = "/node"
				Port = 5_003 #this is the port the other nodes will use to connect to the gateway
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000

				[gatewayConfig.UserServerConfig]
				ContentTypeHeader = "application/jsonrpc"
				MaxRequestBytes = 100_000
				Path = "/"
				Port = 5_002
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000

				[gatewayConfig.HTTPClientConfig]
				MaxResponseBytes = 100_000_000
			`,
				// ETH keys of the workflow nodes
				workflowNodesetInfo[0].TransmitterAddress,
				workflowNodesetInfo[1].TransmitterAddress,
				workflowNodesetInfo[2].TransmitterAddress,
				workflowNodesetInfo[3].TransmitterAddress,
			)

			r, _, err2 := bootstrapNode.CreateJobRaw(gatewayJobSpec)
			assert.NoError(t, err2)
			assert.Empty(t, r.Errors)

			bootstrapJobSpec := fmt.Sprintf(`
				type = "bootstrap"
				schemaVersion = 1
				name = "Botostrap"
				contractID = "%s"
				contractConfigTrackerPollInterval = "1s"
				contractConfigConfirmations = 1
				relay = "evm"

				[relayConfig]
				chainID = %s
				providerType = "ocr3-capability"
			`, ocr3CapabilityAddress, bc.ChainID)
			r, _, err3 := bootstrapNode.CreateJobRaw(bootstrapJobSpec)
			assert.NoError(t, err3)
			assert.Empty(t, r.Errors)
		}()

		for i, nodeClient := range nodeClients {
			// Last node is a bootstrap node, so we skip it
			if i == 0 {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				scJobSpec := `
					type = "standardcapabilities"
					schemaVersion = 1
					name = "streams-capabilities"
					command="/home/capabilities/streams-linux-amd64"
				`
				response, _, err2 := nodeClient.CreateJobRaw(scJobSpec)
				assert.NoError(t, err2)
				assert.Empty(t, response.Errors)

				cronJobSpec := `
					type = "standardcapabilities"
					schemaVersion = 1
					name = "cron-capabilities"
					forwardingAllowed = false
					command = "/home/capabilities/cron-linux-amd64"
					config = ""
				`

				response, _, err3 := nodeClient.CreateJobRaw(cronJobSpec)
				assert.NoError(t, err3)
				assert.Empty(t, response.Errors)

				computeJobSpec := `
					type = "standardcapabilities"
					schemaVersion = 1
					name = "compute-capabilities"
					forwardingAllowed = false
					command = "__builtin_custom-compute-action"
					config = """
					NumWorkers = 3
						[rateLimiter]
						globalRPS = 20.0
						globalBurst = 30
						perSenderRPS = 1.0
						perSenderBurst = 5
					"""
				`

				response, _, err4 := nodeClient.CreateJobRaw(computeJobSpec)
				assert.NoError(t, err4)
				assert.Empty(t, response.Errors)

				consensusJobSpec := fmt.Sprintf(`
					type = "offchainreporting2"
					schemaVersion = 1
					name = "Keystone OCR3 Consensus Capability"
					contractID = "%s"
					ocrKeyBundleID = "%s"
					p2pv2Bootstrappers = [
						"%s@%s",
					]
					relay = "evm"
					pluginType = "plugin"
					transmitterID = "%s"

					[relayConfig]
					chainID = "%s"

					[pluginConfig]
					command = "/usr/local/bin/chainlink-ocr3-capability"
					ocrVersion = 3
					pluginName = "ocr-capability"
					providerType = "ocr3-capability"
					telemetryType = "plugin"

					[onchainSigningStrategy]
					strategyName = 'multi-chain'
					[onchainSigningStrategy.config]
					evm = "%s"
					`,
					ocr3CapabilityAddress,
					nodesInfo[i].OcrKeyBundleID,
					bootstrapNodeInfo.PeerID,
					"node0:5001",
					nodesInfo[i].TransmitterAddress,
					bc.ChainID,
					nodesInfo[i].OcrKeyBundleID,
				)
				fmt.Println("consensusJobSpec", consensusJobSpec)
				response, _, err2 = nodeClient.CreateJobRaw(consensusJobSpec)
				assert.NoError(t, err2)
				assert.Empty(t, response.Errors)
			}()
		}
		wg.Wait()

		var nopsToAdd []cr_wrapper.CapabilitiesRegistryNodeOperator
		var nodesToAdd []cr_wrapper.CapabilitiesRegistryNodeParams
		var donNodes [][32]byte
		var signers []common.Address

		for i, node := range nodesInfo {
			if i == 0 {
				continue
			}
			nopsToAdd = append(nopsToAdd, cr_wrapper.CapabilitiesRegistryNodeOperator{
				Admin: common.HexToAddress(node.TransmitterAddress),
				Name:  fmt.Sprintf("NOP %d", i),
			})

			var peerID ragetypes.PeerID
			err = peerID.UnmarshalText([]byte(node.PeerID))
			require.NoError(t, err)

			nodesToAdd = append(nodesToAdd, cr_wrapper.CapabilitiesRegistryNodeParams{
				NodeOperatorId:      uint32(i), //nolint:gosec // disable G115
				Signer:              common.BytesToHash(node.Signer.Bytes()),
				P2pId:               peerID,
				EncryptionPublicKey: [32]byte{1, 2, 3, 4, 5},
				HashedCapabilityIds: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs,
			})

			donNodes = append(donNodes, peerID)
			signers = append(signers, node.Signer)
		}

		// Add NOPs to registry
		tx, err = capabilitiesRegistryInstance.Contract.AddNodeOperators(
			sc.NewTXOpts(),
			nopsToAdd,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// Add nodes to registry
		tx, err = capabilitiesRegistryInstance.Contract.AddNodes(
			sc.NewTXOpts(),
			nodesToAdd,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// Add nodeset to registry
		tx, err = capabilitiesRegistryInstance.Contract.AddDON(
			sc.NewTXOpts(),
			donNodes,
			[]cr_wrapper.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[0],
					Config:       []byte(""),
				},
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[1],
					Config:       []byte(""),
				},
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[2],
					Config:       []byte(""),
				},
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[3],
					Config:       []byte(""),
				},
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[4],
					Config:       []byte(""),
				},
			},
			true,     // is public
			true,     // accepts workflows
			uint8(1), // max number of malicious nodes
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		require.NoError(t, forwarderInstance.SetConfig(
			1,
			1,
			1,
			signers,
		))

		// Wait for OCR listeners to be ready before setting the configuration.
		// If the ConfigSet event is missed, OCR protocol will not start.
		// TODO make it fluent!
		fmt.Println("Waiting 30s for OCR listeners to be ready...")
		time.Sleep(30 * time.Second)
		fmt.Println("Proceeding to set OCR3 configuration.")

		// Configure OCR capability contract
		ocr3Config := generateOCR3Config(t, workflowNodesetInfo)
		tx, err = ocr3CapabilityContract.SetConfig(
			sc.NewTXOpts(),
			ocr3Config.Signers,
			ocr3Config.Transmitters,
			ocr3Config.F,
			ocr3Config.OnchainConfig,
			ocr3Config.OffchainConfigVersion,
			ocr3Config.OffchainConfig,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// It can take a while before the first report is produced, particularly on CI.
		timeout := 10 * time.Minute
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Node sent transaction

		startTime := time.Now()
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("feed did not update, timeout after %s", timeout)
			case <-time.After(10 * time.Second):
				elapsed := time.Since(startTime).Round(time.Second)
				price, _, err := feedsConsumerContract.GetPrice(
					sc.NewCallOpts(),
					feedBytes,
				)
				require.NoError(t, err)

				if price.String() != "0" {
					fmt.Printf("Feed updated after %s - price set, price=%s\n", elapsed, price)
					return
				}
				// ids, prices, timestamps, err := feedsConsumerContract.GetAllFeeds(sc.NewCallOpts())
				// require.NoError(t, err)

				// for i, feedId := range ids {
				// 	fmt.Printf("Feed %s - price=%d, timestamp=%d\n", common.Bytes2Hex(feedId[:]), prices[i], timestamps[i])
				// 	price, _, err := feedsConsumerContract.GetPrice(
				// 		sc.NewCallOpts(),
				// 		feedId,
				// 	)
				// 	require.NoError(t, err)

				// 	if price.String() != "0" {
				// 		fmt.Printf("Feed updated after %s - price set, price=%s\n", elapsed, price)
				// 		return
				// 	}
				// }
				fmt.Printf("Feed not updated yet, waiting for %s\n", elapsed)
			}
		}
	})
}
