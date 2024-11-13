package capabilities_test

import (
	"bytes"
	"cmp"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	capabilities_registry "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/capabilities_registry_1_1_0"
	feeds_consumer "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/feeds_consumer_1_0_0"
	forwarder "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/ocr3_capability_1_0_0"
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
	var pubKeys [][]byte
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
			return nil, fmt.Errorf("pubKey doesn't fit into uint16")
		}
		if err = binary.Write(buf, binary.LittleEndian, uint16(length)); err != nil { //nolint:gosec
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
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
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
	OcrKeyBundleID     string
	TransmitterAddress string
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
		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err)
		nodesInfo[i].OcrKeyBundleID = ocr2Keys.Data[0].ID

		// eth
		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		nodesInfo[i].TransmitterAddress = ethKeys.Data[0].Attributes.Address
	}

	return nodesInfo
}

func generateOCR3Config(
	t *testing.T,
	nodes []*clclient.ChainlinkClient,
) (config *OCR3Config) {
	oracleIdentities := []confighelper.OracleIdentityExtra{}
	transmissionSchedule := []int{}

	for i, node := range nodes {
		// TODO: Do not provide a bootstrap node to this func
		// We want to skip bootstrap node.
		if i == 0 {
			continue
		}
		transmissionSchedule = append(transmissionSchedule, 0)
		oracleIdentity := confighelper.OracleIdentityExtra{}
		// ocr2
		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err)

		firstOCR2Key := ocr2Keys.Data[0].Attributes

		offchainPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.OffChainPublicKey))
		require.NoError(t, err)
		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], offchainPublicKeyBytes)
		oracleIdentity.OffchainPublicKey = offchainPublicKey

		pubKeys := make(map[string]types.OnchainPublicKey)
		ethOnchainPubKey, err := hex.DecodeString(extractKey(firstOCR2Key.OnChainPublicKey))
		require.NoError(t, err)
		pubKeys["evm"] = ethOnchainPubKey

		// // add aptos key if present
		// if n.AptosOnchainPublicKey != "" {
		// 	aptosPubKey, err := hex.DecodeString(n.AptosOnchainPublicKey)
		// 	if err != nil {
		// 		return Orc2drOracleConfig{}, fmt.Errorf("failed to decode AptosOnchainPublicKey: %w", err)
		// 	}
		// 	pubKeys[string(chaintype.Aptos)] = aptosPubKey
		// }

		multichainPubKey, err := MarshalMultichainPublicKey(pubKeys)
		require.NoError(t, err)
		oracleIdentity.OnchainPublicKey = multichainPubKey

		sharedSecretEncryptionPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.ConfigPublicKey))
		require.NoError(t, err)
		var sharedSecretEncryptionPublicKey [32]byte
		copy(sharedSecretEncryptionPublicKey[:], sharedSecretEncryptionPublicKeyBytes)
		oracleIdentity.ConfigEncryptionPublicKey = sharedSecretEncryptionPublicKey

		// p2p
		p2pKeys, err := node.MustReadP2PKeys()
		require.NoError(t, err)
		oracleIdentity.PeerID = p2pKeys.Data[0].Attributes.PeerID

		// eth
		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		oracleIdentity.TransmitAccount = types.Account(ethKeys.Data[0].Attributes.Address)

		oracleIdentities = append(oracleIdentities, oracleIdentity)
	}

	maxDurationInitialization := 10 * time.Second

	// Generate OCR3 configuration arguments for testing
	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		5*time.Second,              // DeltaProgress: Time between rounds
		5*time.Second,              // DeltaResend: Time between resending unconfirmed transmissions
		5*time.Second,              // DeltaInitial: Initial delay before starting the first round
		2*time.Second,              // DeltaRound: Time between rounds within an epoch
		500*time.Millisecond,       // DeltaGrace: Grace period for delayed transmissions
		1*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
		30*time.Second,             // DeltaStage: Time between stages of the protocol
		uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
		transmissionSchedule,       // TransmissionSchedule: Transmission schedule
		oracleIdentities,           // Oracle identities with their public keys
		nil,                        // Plugin config (empty for now)
		&maxDurationInitialization, // MaxDurationInitialization: ???
		1*time.Second,              // MaxDurationQuery: Maximum duration for querying
		1*time.Second,              // MaxDurationObservation: Maximum duration for observation
		1*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
		1*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
		1,                          // F: Maximum number of faulty oracles
		nil,                        // OnChain config (empty for now)
	)
	require.NoError(t, err)

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

func TestWorkflow(t *testing.T) {
	t.Run("smoke test", func(t *testing.T) {
		in, err := framework.Load[WorkflowTestConfig](t)
		require.NoError(t, err)
		pkey := os.Getenv("PRIVATE_KEY")

		// deploy docker test environment
		bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
		require.NoError(t, err)

		// connect clients
		sc, err := seth.NewClientBuilder().
			WithRpcUrl(bc.Nodes[0].HostWSUrl).
			WithPrivateKeys([]string{pkey}).
			Build()
		require.NoError(t, err)

		capabilitiesRegistryAddress, tx, _, err := capabilities_registry.DeployCapabilitiesRegistry(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed capabilities_registry contract at", capabilitiesRegistryAddress)

		forwarderAddress, tx, _, err := forwarder.DeployKeystoneForwarder(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed forwarder contract at", forwarderAddress)

		// TODO: When the capabilities registry address is provided:
		// - NOPs and nodes are added to the registry.
		// - Nodes are configured to listen to the registry for updates.
		nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, "https://example.com") // TODO: Should not be a thing
		require.NoError(t, err)

		for i, n := range nodeset.CLNodes {
			fmt.Printf("Node %d --> %s\n", i, n.Node.HostURL)
			fmt.Printf("Node P2P %d --> %s\n", i, n.Node.HostP2PURL)
		}

		nodeClients, err := clclient.NewCLDefaultClients(nodeset.CLNodes, framework.L)
		require.NoError(t, err)
		nodesInfo := getNodesInfo(t, nodeClients)

		for i, node := range nodeClients {
			fmt.Println("Node i ", i)
			fmt.Println("Node ", node)
			// First node is a bootstrap node, so we skip it
			if i == 0 {
				continue
			}

			in.NodeSet.NodeSpecs[i].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']

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

				# This is needed for external registry
				[Capabilities.ExternalRegistry]
				Address = '%s'
				NetworkID = 'evm'
				ChainID = '%s'
			`,
				bc.ChainID,
				bc.Nodes[0].DockerInternalWSUrl,
				bc.Nodes[0].DockerInternalHTTPUrl,
				nodesInfo[i].TransmitterAddress,
				forwarderAddress,
				capabilitiesRegistryAddress,
				bc.ChainID,
			)
		}

		nodeset, err = ns.UpgradeNodeSet(in.NodeSet, bc, "https://example.com", 5*time.Second)
		require.NoError(t, err)
		nodeClients, err = clclient.NewCLDefaultClients(nodeset.CLNodes, framework.L)
		require.NoError(t, err)

		ocr3CapabilityAddress, tx, ocr3CapabilityContract, err := ocr3_capability.DeployOCR3Capability(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed ocr3_capability contract at", ocr3CapabilityAddress.Hex())

		feedsConsumerAddress, tx, _, err := feeds_consumer.DeployKeystoneFeedsConsumer(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed feeds_consumer contract at", feedsConsumerAddress.Hex())

		// Add bootstrap spec to the first node
		bootstrapNode := nodeClients[0]

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
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
			`, ocr3CapabilityAddress, bc.ChainID)
			fmt.Println("Creating bootstrap job spec", bootstrapJobSpec)
			r, _, err2 := bootstrapNode.CreateJobRaw(bootstrapJobSpec)
			require.NoError(t, err2)
			require.Equal(t, len(r.Errors), 0)
			fmt.Printf("Response from bootstrap node: %x\n", r)
		}()

		p2pKeys, err := bootstrapNode.MustReadP2PKeys()
		require.NoError(t, err)
		fmt.Println("P2P keys fetched")

		for i, nodeClient := range nodeClients {
			// First node is a bootstrap node, so we skip it
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
					command="/home/capabilities/streams"
				`
				fmt.Println("Creating standard capabilities job spec", scJobSpec)
				response, _, err2 := nodeClient.CreateJobRaw(scJobSpec)
				require.NoError(t, err2)
				require.Equal(t, len(response.Errors), 0)
				fmt.Printf("Response from node %d after streams SC: %x\n", i+1, response)

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
					p2pKeys.Data[0].Attributes.PeerID,
					strings.TrimPrefix(nodeset.CLNodes[0].Node.HostP2PURL, "http://"),
					nodesInfo[i].TransmitterAddress,
					bc.ChainID,
					nodesInfo[i].OcrKeyBundleID,
				)
				fmt.Println("Creating consensus job spec", consensusJobSpec)
				response, _, err2 = nodeClient.CreateJobRaw(consensusJobSpec)
				fmt.Println("err2", err2)
				require.NoError(t, err2)
				require.Equal(t, len(response.Errors), 0)
				fmt.Printf("Response from node %d after consensus job: %x\n", i+1, response)
			}()
		}
		wg.Wait()

		ocr3Config := generateOCR3Config(t, nodeClients)
		// Configure KV store OCR contract
		tx, err = ocr3CapabilityContract.SetConfig(
			sc.NewTXOpts(),
			ocr3Config.Signers,
			ocr3Config.Transmitters,
			ocr3Config.F,
			ocr3Config.OnchainConfig,
			ocr3Config.OffchainConfigVersion,
			ocr3Config.OffchainConfig,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)

		// ✅ Add bootstrap spec
		// ✅ 1. Deploy mock streams capability
		// ✅ 2. Add boostrap job spec
		// ✅ 3. Add OCR3 capability
		// 		- ✅ This starts successfully (search for "OCREndpointV2: Initialized")
		// ✅ 3. Deploy and configure OCR3 contract
		// ✅ 4. Add chain write capabilities
		//  	- Check if they are added (Logs)
		// 5. Deploy capabilities registry
		// 		- Add nodes to registry
		// 		- Add capabilities to registry
		// ✅ 6. Deploy Forwarder
		//      - Configure forwarder
		// ✅ 7. Deploy Feeds Consumer
		// - Add Keystone workflow
	})
}
