package actions

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func CreateOCRv2JobsLocal(
	ocrInstances []contracts.OffchainAggregatorV2,
	bootstrapNode *client.ChainlinkClient,
	workerChainlinkNodes []*client.ChainlinkClient,
	mockAdapter *test_env.Killgrave,
	mockAdapterPath string, // Path on the mock server for the Chainlink nodes to query
	mockAdapterValue int, // Value to get from the mock server when querying the path
	chainId uint64, // EVM chain ID
	forwardingAllowed bool,
	enableChainReaderAndCodec bool,
) error {
	// Collect P2P ID
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	if err != nil {
		return err
	}
	p2pV2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PIds.Data[0].Attributes.PeerID, bootstrapNode.InternalIP(), 6690)
	// Set the value for the jobs to report on
	err = mockAdapter.SetAdapterBasedIntValuePath(mockAdapterPath, []string{http.MethodGet, http.MethodPost}, mockAdapterValue)
	if err != nil {
		return err
	}
	// Set the juelsPerFeeCoinSource config value
	err = mockAdapter.SetAdapterBasedIntValuePath(fmt.Sprintf("%s/juelsPerFeeCoinSource", mockAdapterPath), []string{http.MethodGet, http.MethodPost}, mockAdapterValue)
	if err != nil {
		return err
	}

	for _, ocrInstance := range ocrInstances {
		bootstrapSpec := &client.OCR2TaskJobSpec{
			Name:    fmt.Sprintf("ocr2_bootstrap-%s", uuid.NewString()),
			JobType: "bootstrap",
			OCR2OracleSpec: job.OCR2OracleSpec{
				ContractID: ocrInstance.Address(),
				Relay:      "evm",
				RelayConfig: map[string]interface{}{
					"chainID": chainId,
				},
				MonitoringEndpoint:                null.StringFrom(fmt.Sprintf("%s/%s", mockAdapter.InternalEndpoint, mockAdapterPath)),
				ContractConfigTrackerPollInterval: *models.NewInterval(15 * time.Second),
			},
		}
		_, err := bootstrapNode.MustCreateJob(bootstrapSpec)
		if err != nil {
			return fmt.Errorf("creating bootstrap job have failed: %w", err)
		}

		for _, chainlinkNode := range workerChainlinkNodes {
			nodeTransmitterAddress, err := chainlinkNode.PrimaryEthAddress()
			if err != nil {
				return fmt.Errorf("getting primary ETH address from OCR node have failed: %w", err)
			}
			nodeOCRKeys, err := chainlinkNode.MustReadOCR2Keys()
			if err != nil {
				return fmt.Errorf("getting OCR keys from OCR node have failed: %w", err)
			}
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			bta := &client.BridgeTypeAttributes{
				Name: fmt.Sprintf("%s-%s", mockAdapterPath, uuid.NewString()),
				URL:  fmt.Sprintf("%s/%s", mockAdapter.InternalEndpoint, mockAdapterPath),
			}
			juelsBridge := &client.BridgeTypeAttributes{
				Name: fmt.Sprintf("juels-%s", uuid.NewString()),
				URL:  fmt.Sprintf("%s/%s/juelsPerFeeCoinSource", mockAdapter.InternalEndpoint, mockAdapterPath),
			}
			err = chainlinkNode.MustCreateBridge(bta)
			if err != nil {
				return fmt.Errorf("creating bridge on CL node failed: %w", err)
			}
			err = chainlinkNode.MustCreateBridge(juelsBridge)
			if err != nil {
				return fmt.Errorf("creating bridge on CL node failed: %w", err)
			}

			ocrSpec := &client.OCR2TaskJobSpec{
				Name:              fmt.Sprintf("ocr2-%s", uuid.NewString()),
				JobType:           "offchainreporting2",
				MaxTaskDuration:   "1m",
				ObservationSource: client.ObservationSourceSpecBridge(bta),
				ForwardingAllowed: forwardingAllowed,
				OCR2OracleSpec: job.OCR2OracleSpec{
					PluginType: "median",
					Relay:      "evm",
					RelayConfig: map[string]interface{}{
						"chainID": chainId,
					},
					PluginConfig: map[string]any{
						"juelsPerFeeCoinSource": fmt.Sprintf("\"\"\"%s\"\"\"", client.ObservationSourceSpecBridge(juelsBridge)),
					},
					ContractConfigTrackerPollInterval: *models.NewInterval(15 * time.Second),
					ContractID:                        ocrInstance.Address(),                   // registryAddr
					OCRKeyBundleID:                    null.StringFrom(nodeOCRKeyId),           // get node ocr2config.ID
					TransmitterID:                     null.StringFrom(nodeTransmitterAddress), // node addr
					P2PV2Bootstrappers:                pq.StringArray{p2pV2Bootstrapper},       // bootstrap node key and address <p2p-key>@bootstrap:6690
				},
			}
			if enableChainReaderAndCodec {
				ocrSpec.OCR2OracleSpec.RelayConfig["chainReader"] = evmtypes.ChainReaderConfig{
					Contracts: map[string]evmtypes.ChainContractReader{
						"median": {
							ContractPollingFilter: evmtypes.ContractPollingFilter{
								GenericEventNames: []string{"LatestRoundRequested"},
							},
							ContractABI: `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"requester","type":"address"},{"indexed":false,"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"indexed":false,"internalType":"uint32","name":"epoch","type":"uint32"},{"indexed":false,"internalType":"uint8","name":"round","type":"uint8"}],"name":"RoundRequested","type":"event"},{"inputs":[],"name":"latestTransmissionDetails","outputs":[{"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"internalType":"uint32","name":"epoch","type":"uint32"},{"internalType":"uint8","name":"round","type":"uint8"},{"internalType":"int192","name":"latestAnswer_","type":"int192"},{"internalType":"uint64","name":"latestTimestamp_","type":"uint64"}],"stateMutability":"view","type":"function"}]`,
							Configs: map[string]*evmtypes.ChainReaderDefinition{
								"LatestTransmissionDetails": {
									ChainSpecificName: "latestTransmissionDetails",
									OutputModifications: codec.ModifiersConfig{
										&codec.EpochToTimeModifierConfig{
											Fields: []string{"LatestTimestamp_"},
										},
										&codec.RenameModifierConfig{
											Fields: map[string]string{
												"LatestAnswer_":    "LatestAnswer",
												"LatestTimestamp_": "LatestTimestamp",
											},
										},
									},
								},
								"LatestRoundRequested": {
									ChainSpecificName: "RoundRequested",
									ReadType:          evmtypes.Event,
								},
							},
						},
					},
				}
				ocrSpec.OCR2OracleSpec.RelayConfig["codec"] = evmtypes.CodecConfig{
					Configs: map[string]evmtypes.ChainCodecConfig{
						"MedianReport": {
							TypeABI: `[{"Name": "Timestamp","Type": "uint32"},{"Name": "Observers","Type": "bytes32"},{"Name": "Observations","Type": "int192[]"},{"Name": "JuelsPerFeeCoin","Type": "int192"}]`,
						},
					},
				}
			}

			_, err = chainlinkNode.MustCreateJob(ocrSpec)
			if err != nil {
				return fmt.Errorf("creating OCR task job on OCR node have failed: %w", err)
			}
		}
	}
	return nil
}

func BuildMedianOCR2ConfigLocal(workerNodes []*client.ChainlinkClient, ocrOffchainOptions contracts.OffchainOptions) (*contracts.OCRv2Config, error) {
	S, oracleIdentities, err := GetOracleIdentitiesWithKeyIndexLocal(workerNodes, 0)
	if err != nil {
		return nil, err
	}
	signerKeys, transmitterAccounts, f_, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		30*time.Second,   // deltaProgress time.Duration,
		30*time.Second,   // deltaResend time.Duration,
		10*time.Second,   // deltaRound time.Duration,
		20*time.Second,   // deltaGrace time.Duration,
		20*time.Second,   // deltaStage time.Duration,
		3,                // rMax uint8,
		S,                // s []int,
		oracleIdentities, // oracles []OracleIdentityExtra,
		median.OffchainConfig{
			AlphaReportInfinite: false,
			AlphaReportPPB:      1,
			AlphaAcceptInfinite: false,
			AlphaAcceptPPB:      1,
			DeltaC:              time.Minute * 30,
		}.Encode(), // reportingPluginConfig []byte,
		5*time.Second, // maxDurationQuery time.Duration,
		5*time.Second, // maxDurationObservation time.Duration,
		5*time.Second, // maxDurationReport time.Duration,
		5*time.Second, // maxDurationShouldAcceptFinalizedReport time.Duration,
		5*time.Second, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,             // f int,
		nil,           // The median reporting plugin has an empty onchain config
	)
	if err != nil {
		return nil, err
	}

	// Convert signers to addresses
	var signerAddresses []common.Address
	for _, signer := range signerKeys {
		signerAddresses = append(signerAddresses, common.BytesToAddress(signer))
	}

	// Convert transmitters to addresses
	var transmitterAddresses []common.Address
	for _, account := range transmitterAccounts {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(account)))
	}

	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(ocrOffchainOptions.MinimumAnswer, ocrOffchainOptions.MaximumAnswer)

	return &contracts.OCRv2Config{
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     f_,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        []byte(fmt.Sprintf("0x%s", offchainConfig)),
	}, err
}

func GetOracleIdentitiesWithKeyIndexLocal(
	chainlinkNodes []*client.ChainlinkClient,
	keyIndex int,
) ([]int, []confighelper.OracleIdentityExtra, error) {
	S := make([]int, len(chainlinkNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(chainlinkNodes))
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(chainlinkNodes))
	eg := &errgroup.Group{}
	for i, cl := range chainlinkNodes {
		index, chainlinkNode := i, cl
		eg.Go(func() error {
			addresses, err := chainlinkNode.EthAddresses()
			if err != nil {
				return err
			}
			ocr2Keys, err := chainlinkNode.MustReadOCR2Keys()
			if err != nil {
				return err
			}
			var ocr2Config client.OCR2KeyAttributes
			for _, key := range ocr2Keys.Data {
				if key.Attributes.ChainType == string(chaintype.EVM) {
					ocr2Config = key.Attributes
					break
				}
			}

			keys, err := chainlinkNode.MustReadP2PKeys()
			if err != nil {
				return err
			}
			p2pKeyID := keys.Data[0].Attributes.PeerID

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			if err != nil {
				return err
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("wrong number of elements copied")
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("wrong number of elements copied")
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_"))
			if err != nil {
				return err
			}

			sharedSecretEncryptionPublicKeys[index] = configPkBytesFixed
			oracleIdentities[index] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   types.Account(addresses[keyIndex]),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[index] = 1
			log.Debug().
				Interface("OnChainPK", onchainPkBytes).
				Interface("OffChainPK", offchainPkBytesFixed).
				Interface("ConfigPK", configPkBytesFixed).
				Str("PeerID", p2pKeyID).
				Str("Address", addresses[keyIndex]).
				Msg("Oracle identity")
			return nil
		})
	}

	return S, oracleIdentities, eg.Wait()
}

// DeleteJobs will delete ALL jobs from the nodes
func DeleteJobs(nodes []*client.ChainlinkClient) error {
	for _, node := range nodes {
		if node == nil {
			return fmt.Errorf("found a nil chainlink node in the list of chainlink nodes while tearing down: %v", nodes)
		}
		jobs, _, err := node.ReadJobs()
		if err != nil {
			return fmt.Errorf("error reading jobs from chainlink node: %w", err)
		}
		for _, maps := range jobs.Data {
			if _, ok := maps["id"]; !ok {
				return fmt.Errorf("error reading job id from chainlink node's jobs %+v", jobs.Data)
			}
			id := maps["id"].(string)
			_, err2 := node.DeleteJob(id)
			if err2 != nil {
				return fmt.Errorf("error deleting job from chainlink node: %w", err)
			}
		}
	}
	return nil
}

// DeleteBridges will delete ALL bridges from the nodes
func DeleteBridges(nodes []*client.ChainlinkClient) error {
	for _, node := range nodes {
		if node == nil {
			return fmt.Errorf("found a nil chainlink node in the list of chainlink nodes while tearing down: %v", nodes)
		}

		bridges, _, err := node.ReadBridges()
		if err != nil {
			return err
		}
		for _, b := range bridges.Data {
			_, err = node.DeleteBridge(b.Attributes.Name)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
