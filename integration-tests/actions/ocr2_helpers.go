package actions

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/lib/client"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// BuildMedianOCR2Config builds a default OCRv2 config for the given chainlink nodes for a standard median aggregation job
func BuildMedianOCR2Config(
	workerNodes []*nodeclient.ChainlinkK8sClient,
	ocrOffchainOptions contracts.OffchainOptions,
) (*contracts.OCRv2Config, error) {
	S, oracleIdentities, err := GetOracleIdentities(workerNodes)
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
		nil,
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

// GetOracleIdentities retrieves all chainlink nodes' OCR2 config identities with defaul key index
func GetOracleIdentities(chainlinkNodes []*nodeclient.ChainlinkK8sClient) ([]int, []confighelper.OracleIdentityExtra, error) {
	return GetOracleIdentitiesWithKeyIndex(chainlinkNodes, 0)
}

// GetOracleIdentitiesWithKeyIndex retrieves all chainlink nodes' OCR2 config identities by key index
func GetOracleIdentitiesWithKeyIndex(
	chainlinkNodes []*nodeclient.ChainlinkK8sClient,
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
			var ocr2Config nodeclient.OCR2KeyAttributes
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

// CreateOCRv2Jobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRv2Jobs(
	ocrInstances []contracts.OffchainAggregatorV2,
	bootstrapNode *nodeclient.ChainlinkK8sClient,
	workerChainlinkNodes []*nodeclient.ChainlinkK8sClient,
	mockserver *ctfClient.MockserverClient,
	mockServerValue int, // Value to get from the mock server when querying the path
	chainId int64, // EVM chain ID
	forwardingAllowed bool,
	l zerolog.Logger,
) error {
	// Collect P2P ID
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	if err != nil {
		return err
	}
	p2pV2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PIds.Data[0].Attributes.PeerID, bootstrapNode.InternalIP(), 6690)
	mockJuelsPath := "ocr2/juelsPerFeeCoinSource"
	// Set the juelsPerFeeCoinSource config value
	err = mockserver.SetValuePath(mockJuelsPath, mockServerValue)
	if err != nil {
		return err
	}

	// Create the juels bridge for each node only once
	juelsBridge := &nodeclient.BridgeTypeAttributes{
		Name: "juels",
		URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, mockJuelsPath),
	}
	for _, chainlinkNode := range workerChainlinkNodes {
		err = chainlinkNode.MustCreateBridge(juelsBridge)
		if err != nil {
			return fmt.Errorf("failed creating bridge %s on CL node : %w", juelsBridge.Name, err)
		}
	}

	// Initialize map to store job IDs for each chainlink node
	jobIDs := make(map[*nodeclient.ChainlinkK8sClient][]string)

	for _, ocrInstance := range ocrInstances {
		bootstrapSpec := &nodeclient.OCR2TaskJobSpec{
			Name:    fmt.Sprintf("ocr2-bootstrap-%s", ocrInstance.Address()),
			JobType: "bootstrap",
			OCR2OracleSpec: job.OCR2OracleSpec{
				ContractID: ocrInstance.Address(),
				Relay:      "evm",
				RelayConfig: map[string]interface{}{
					"chainID": chainId,
				},
				MonitoringEndpoint:                null.StringFrom(fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, "ocr2")),
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

			nodeContractPairID, err := BuildOCR2NodeContractPairID(chainlinkNode, ocrInstance)
			if err != nil {
				return err
			}
			bta := &nodeclient.BridgeTypeAttributes{
				Name: nodeContractPairID,
				URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, strings.TrimPrefix(nodeContractPairID, "/")),
			}

			err = chainlinkNode.MustCreateBridge(bta)
			if err != nil {
				return fmt.Errorf("failed creating bridge %s on CL node: %w", bta.Name, err)
			}

			ocrSpec := &nodeclient.OCR2TaskJobSpec{
				Name:              fmt.Sprintf("ocr2-%s", uuid.NewString()),
				JobType:           "offchainreporting2",
				MaxTaskDuration:   "1m",
				ObservationSource: nodeclient.ObservationSourceSpecBridge(bta),
				ForwardingAllowed: forwardingAllowed,
				OCR2OracleSpec: job.OCR2OracleSpec{
					PluginType: "median",
					Relay:      "evm",
					RelayConfig: map[string]interface{}{
						"chainID": chainId,
					},
					PluginConfig: map[string]any{
						"juelsPerFeeCoinSource": fmt.Sprintf("\"\"\"%s\"\"\"", nodeclient.ObservationSourceSpecBridge(juelsBridge)),
					},
					ContractConfigTrackerPollInterval: *models.NewInterval(15 * time.Second),
					ContractID:                        ocrInstance.Address(),                   // registryAddr
					OCRKeyBundleID:                    null.StringFrom(nodeOCRKeyId),           // get node ocr2config.ID
					TransmitterID:                     null.StringFrom(nodeTransmitterAddress), // node addr
					P2PV2Bootstrappers:                pq.StringArray{p2pV2Bootstrapper},       // bootstrap node key and address <p2p-key>@bootstrap:6690
				},
			}
			var ocrJob *nodeclient.Job
			ocrJob, err = chainlinkNode.MustCreateJob(ocrSpec)
			if err != nil {
				return fmt.Errorf("creating OCR task job on OCR node have failed: %w", err)
			}
			jobIDs[chainlinkNode] = append(jobIDs[chainlinkNode], ocrJob.Data.ID) // Store each job ID per node
		}
	}
	l.Info().Msg("Verify OCRv2 jobs have been created")
	for chainlinkNode, ids := range jobIDs {
		for _, jobID := range ids {
			err := retry.Do(
				func() error {
					_, resp, err := chainlinkNode.ReadJob(jobID)
					if err != nil {
						return err
					}
					if resp.StatusCode != http.StatusOK {
						return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
					}
					l.Info().
						Str("Node", chainlinkNode.PodName).
						Str("Job ID", jobID).
						Msg("OCRv2 job successfully created")
					return nil
				},
				retry.Attempts(4),
				retry.Delay(time.Second*2),
				retry.OnRetry(func(n uint, err error) {
					l.Debug().
						Str("Node", chainlinkNode.PodName).
						Str("Job ID", jobID).
						Uint("Attempt", n+1).
						Err(err).
						Msg("Retrying job verification")
				}),
			)
			if err != nil {
				l.Error().Err(err).Str("Node", chainlinkNode.PodName).Str("JobID", jobID).Msg("Failed to verify OCRv2 job creation")
			}
		}
	}
	return nil
}

// SetOCR2AdapterResponse sets a single adapter response that correlates with an ocr contract and a chainlink node
// used for OCR2 tests
func SetOCR2AdapterResponse(
	response int,
	ocrInstance contracts.OffchainAggregatorV2,
	chainlinkNode *nodeclient.ChainlinkK8sClient,
	mockserver *ctfClient.MockserverClient,
) error {
	nodeContractPairID, err := BuildOCR2NodeContractPairID(chainlinkNode, ocrInstance)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/%s", nodeContractPairID)
	err = mockserver.SetValuePath(path, response)
	if err != nil {
		return fmt.Errorf("setting mockserver value path failed: %w", err)
	}
	return nil
}

// SetOCR2AllAdapterResponsesToTheSameValue sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to the same response
// used for OCR2 tests
func SetOCR2AllAdapterResponsesToTheSameValue(
	response int,
	ocrInstances []contracts.OffchainAggregatorV2,
	chainlinkNodes []*nodeclient.ChainlinkK8sClient,
	mockserver *ctfClient.MockserverClient,
) error {
	eg := &errgroup.Group{}
	for _, o := range ocrInstances {
		ocrInstance := o
		for _, n := range chainlinkNodes {
			node := n
			eg.Go(func() error {
				return SetOCR2AdapterResponse(response, ocrInstance, node, mockserver)
			})
		}
	}
	return eg.Wait()
}

// BuildOCR2NodeContractPairID builds a UUID based on a related pair of a Chainlink node and OCRv2 contract
func BuildOCR2NodeContractPairID(node *nodeclient.ChainlinkK8sClient, ocrInstance contracts.OffchainAggregatorV2) (string, error) {
	if node == nil {
		return "", fmt.Errorf("chainlink node is nil")
	}
	if ocrInstance == nil {
		return "", fmt.Errorf("OCR Instance is nil")
	}
	nodeAddress, err := node.PrimaryEthAddress()
	if err != nil {
		return "", fmt.Errorf("getting chainlink node's primary ETH address failed: %w", err)
	}
	shortNodeAddr := nodeAddress[2:12]
	shortOCRAddr := ocrInstance.Address()[2:12]
	return strings.ToLower(fmt.Sprintf("node_%s_contract_%s", shortNodeAddr, shortOCRAddr)), nil
}
