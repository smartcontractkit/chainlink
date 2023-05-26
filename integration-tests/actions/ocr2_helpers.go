package actions

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// DeployOCRv2Contracts deploys a number of OCRv2 contracts and configures them with defaults
func DeployOCRv2Contracts(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkWorkerNodes []*client.Chainlink,
	client blockchain.EVMClient,
) ([]contracts.OffchainAggregatorV2, error) {
	var ocrInstances []contracts.OffchainAggregatorV2
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffchainAggregatorV2(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		if err != nil {
			return nil, fmt.Errorf("OCRv2 instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for OCRv2 contract deployments: %w", err)
			}
		}
	}
	err := client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCRv2 contract deployments: %w", err)
	}

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, node := range chainlinkWorkerNodes {
		addr, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, fmt.Errorf("error getting node's primary ETH address: %w", err)
		}
		transmitters = append(transmitters, addr)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err = ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for setting OCR payees: %w", err)
			}
		}
	}
	return ocrInstances, client.WaitForEvents()
}

func ConfigureOCRv2AggregatorContracts(
	client blockchain.EVMClient,
	contractConfig *contracts.OCRv2Config,
	ocrv2Contracts []contracts.OffchainAggregatorV2,
) error {
	for contractCount, ocrInstance := range ocrv2Contracts {
		// Exclude the first node, which will be used as a bootstrapper
		err := ocrInstance.SetConfig(contractConfig)
		if err != nil {
			return fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return fmt.Errorf("failed to wait for setting OCR config: %w", err)
			}
		}
	}
	return client.WaitForEvents()
}

// BuildMedianOCR2Config builds a default OCRv2 config for the given chainlink nodes for a standard median aggregation job
func BuildMedianOCR2Config(workerNodes []*client.Chainlink) (*contracts.OCRv2Config, error) {
	S, oracleIdentities, err := GetOracleIdentities(workerNodes)
	if err != nil {
		return nil, err
	}
	signerKeys, transmitterAccounts, f_, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
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

	return &contracts.OCRv2Config{
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     f_,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        []byte(fmt.Sprintf("0x%s", offchainConfig)),
	}, nil
}

// GetOracleIdentities retrieves all chainlink nodes' OCR2 config identities with defaul key index
func GetOracleIdentities(chainlinkNodes []*client.Chainlink) ([]int, []confighelper.OracleIdentityExtra, error) {
	return GetOracleIdentitiesWithKeyIndex(chainlinkNodes, 0)
}

// GetOracleIdentitiesWithKeyIndex retrieves all chainlink nodes' OCR2 config identities by key index
func GetOracleIdentitiesWithKeyIndex(
	chainlinkNodes []*client.Chainlink,
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
				return fmt.Errorf("Wrong number of elements copied")
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("Wrong number of elements copied")
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

// CreateOCRJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRv2Jobs(
	ocrInstances []contracts.OffchainAggregatorV2,
	bootstrapNode *client.Chainlink,
	workerChainlinkNodes []*client.Chainlink,
	mockserver *ctfClient.MockserverClient,
	mockServerPath string, // Path on the mock server for the Chainlink nodes to query
	mockServerValue int, // Value to get from the mock server when querying the path
	chainId uint64, // EVM chain ID
) error {
	// Collect P2P ID
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	if err != nil {
		return err
	}
	p2pV2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PIds.Data[0].Attributes.PeerID, bootstrapNode.InternalIP(), 6690)
	// Set the value for the jobs to report on
	err = mockserver.SetValuePath(mockServerPath, mockServerValue)
	if err != nil {
		return err
	}
	// Set the juelsPerFeeCoinSource config value
	err = mockserver.SetValuePath(fmt.Sprintf("%s/juelsPerFeeCoinSource", mockServerPath), mockServerValue)
	if err != nil {
		return err
	}

	for _, ocrInstance := range ocrInstances {
		bootstrapSpec := &client.OCR2TaskJobSpec{
			Name:    "ocr2 bootstrap node",
			JobType: "bootstrap",
			OCR2OracleSpec: job.OCR2OracleSpec{
				ContractID: ocrInstance.Address(),
				Relay:      "evm",
				RelayConfig: map[string]interface{}{
					"chainID": chainId,
				},
				MonitoringEndpoint:                null.StringFrom(fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, mockServerPath)),
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
				Name: mockServerPath,
				URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, mockServerPath),
			}
			juelsBridge := &client.BridgeTypeAttributes{
				Name: "juels",
				URL:  fmt.Sprintf("%s/%s/juelsPerFeeCoinSource", mockserver.Config.ClusterURL, mockServerPath),
			}
			err = chainlinkNode.MustCreateBridge(bta)
			if err != nil {
				return fmt.Errorf("creating bridge job have failed: %w", err)
			}
			err = chainlinkNode.MustCreateBridge(juelsBridge)
			if err != nil {
				return fmt.Errorf("creating bridge job have failed: %w", err)
			}

			ocrSpec := &client.OCR2TaskJobSpec{
				Name:              "ocr2",
				JobType:           "offchainreporting2",
				MaxTaskDuration:   "1m",
				ObservationSource: client.ObservationSourceSpecBridge(bta),
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
			_, err = chainlinkNode.MustCreateJob(ocrSpec)
			if err != nil {
				return fmt.Errorf("creating OCR task job on OCR node have failed: %w", err)
			}
		}
	}
	return nil
}

// StartNewOCR2Round requests a new round from the ocr2 contracts and waits for confirmation
func StartNewOCR2Round(
	roundNumber int64,
	ocrInstances []contracts.OffchainAggregatorV2,
	client blockchain.EVMClient,
	timeout time.Duration,
) error {
	for i := 0; i < len(ocrInstances); i++ {
		err := ocrInstances[i].RequestNewRound()
		if err != nil {
			return fmt.Errorf("requesting new OCR round %d have failed: %w", i+1, err)
		}
		ocrRound := contracts.NewOffchainAggregatorV2RoundConfirmer(ocrInstances[i], big.NewInt(roundNumber), timeout, nil)
		client.AddHeaderEventSubscription(ocrInstances[i].Address(), ocrRound)
		err = client.WaitForEvents()
		if err != nil {
			return fmt.Errorf("failed to wait for event subscriptions of OCR instance %d: %w", i+1, err)
		}
	}
	return nil
}
