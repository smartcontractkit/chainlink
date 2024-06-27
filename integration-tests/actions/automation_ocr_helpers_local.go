package actions

//revive:disable:dot-imports
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	ocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2keepers20config "github.com/smartcontractkit/ocr2keepers/pkg/v2/config"
	ocr2keepers30config "github.com/smartcontractkit/ocr2keepers/pkg/v3/config"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

func BuildAutoOCR2ConfigVarsLocal(
	l zerolog.Logger,
	chainlinkNodes []*client.ChainlinkClient,
	registryConfig contracts.KeeperRegistrySettings,
	registrar string,
	deltaStage time.Duration,
	registryOwnerAddress common.Address,
) (contracts.OCRv2Config, error) {
	return BuildAutoOCR2ConfigVarsWithKeyIndexLocal(l, chainlinkNodes, registryConfig, registrar, deltaStage, 0, registryOwnerAddress)
}

func BuildAutoOCR2ConfigVarsWithKeyIndexLocal(
	l zerolog.Logger,
	chainlinkNodes []*client.ChainlinkClient,
	registryConfig contracts.KeeperRegistrySettings,
	registrar string,
	deltaStage time.Duration,
	keyIndex int,
	registryOwnerAddress common.Address,
) (contracts.OCRv2Config, error) {
	S, oracleIdentities, err := GetOracleIdentitiesWithKeyIndexLocal(chainlinkNodes, keyIndex)
	if err != nil {
		return contracts.OCRv2Config{}, err
	}

	var offC []byte
	var signerOnchainPublicKeys []types.OnchainPublicKey
	var transmitterAccounts []types.Account
	var f uint8
	var offchainConfigVersion uint64
	var offchainConfig []byte

	if registryConfig.RegistryVersion == ethereum.RegistryVersion_2_1 {
		offC, err = json.Marshal(ocr2keepers30config.OffchainConfig{
			TargetProbability:    "0.999",
			TargetInRounds:       1,
			PerformLockoutWindow: 3600000, // Intentionally set to be higher than in prod for testing purpose
			GasLimitPerReport:    5_300_000,
			GasOverheadPerUpkeep: 300_000,
			MinConfirmations:     0,
			MaxUpkeepBatchSize:   1,
		})
		if err != nil {
			return contracts.OCRv2Config{}, err
		}

		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = ocr3.ContractSetConfigArgsForTests(
			10*time.Second,        // deltaProgress time.Duration,
			15*time.Second,        // deltaResend time.Duration,
			500*time.Millisecond,  // deltaInitial time.Duration,
			1000*time.Millisecond, // deltaRound time.Duration,
			200*time.Millisecond,  // deltaGrace time.Duration,
			300*time.Millisecond,  // deltaCertifiedCommitRequest time.Duration
			deltaStage,            // deltaStage time.Duration,
			24,                    // rMax uint64,
			S,                     // s []int,
			oracleIdentities,      // oracles []OracleIdentityExtra,
			offC,                  // reportingPluginConfig []byte,
			20*time.Millisecond,   // maxDurationQuery time.Duration,
			20*time.Millisecond,   // maxDurationObservation time.Duration, // good to here
			1200*time.Millisecond, // maxDurationShouldAcceptAttestedReport time.Duration,
			20*time.Millisecond,   // maxDurationShouldTransmitAcceptedReport time.Duration,
			1,                     // f int,
			nil,                   // onchainConfig []byte,
		)
		if err != nil {
			return contracts.OCRv2Config{}, err
		}
	} else {
		offC, err = json.Marshal(ocr2keepers20config.OffchainConfig{
			TargetProbability:    "0.999",
			TargetInRounds:       1,
			PerformLockoutWindow: 3600000, // Intentionally set to be higher than in prod for testing purpose
			GasLimitPerReport:    5_300_000,
			GasOverheadPerUpkeep: 300_000,
			SamplingJobDuration:  3000,
			MinConfirmations:     0,
			MaxUpkeepBatchSize:   1,
		})
		if err != nil {
			return contracts.OCRv2Config{}, err
		}

		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = ocr2.ContractSetConfigArgsForTests(
			10*time.Second,        // deltaProgress time.Duration,
			15*time.Second,        // deltaResend time.Duration,
			3000*time.Millisecond, // deltaRound time.Duration,
			200*time.Millisecond,  // deltaGrace time.Duration,
			deltaStage,            // deltaStage time.Duration,
			24,                    // rMax uint8,
			S,                     // s []int,
			oracleIdentities,      // oracles []OracleIdentityExtra,
			offC,                  // reportingPluginConfig []byte,
			20*time.Millisecond,   // maxDurationQuery time.Duration,
			20*time.Millisecond,   // maxDurationObservation time.Duration,
			1200*time.Millisecond, // maxDurationReport time.Duration,
			20*time.Millisecond,   // maxDurationShouldAcceptFinalizedReport time.Duration,
			20*time.Millisecond,   // maxDurationShouldTransmitAcceptedReport time.Duration,
			1,                     // f int,
			nil,                   // onchainConfig []byte,
		)
		if err != nil {
			return contracts.OCRv2Config{}, err
		}
	}

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		if len(signer) != 20 {
			return contracts.OCRv2Config{}, fmt.Errorf("OnChainPublicKey '%v' has wrong length for address", signer)
		}
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		if !common.IsHexAddress(string(transmitter)) {
			return contracts.OCRv2Config{}, fmt.Errorf("TransmitAccount '%s' is not a valid Ethereum address", string(transmitter))
		}
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	onchainConfig, err := registryConfig.EncodeOnChainConfig(registrar, registryOwnerAddress)
	if err != nil {
		return contracts.OCRv2Config{}, err
	}

	l.Info().Msg("Done building OCR config")
	return contracts.OCRv2Config{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}

// CreateOCRKeeperJobs bootstraps the first node and to the other nodes sends ocr jobs
func CreateOCRKeeperJobsLocal(
	l zerolog.Logger,
	chainlinkNodes []*client.ChainlinkClient,
	registryAddr string,
	chainID int64,
	keyIndex int,
	registryVersion ethereum.KeeperRegistryVersion,
) error {
	bootstrapNode := chainlinkNodes[0]
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	if err != nil {
		l.Error().Err(err).Msg("Shouldn't fail reading P2P keys from bootstrap node")
		return err
	}
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID

	var contractVersion string
	if registryVersion == ethereum.RegistryVersion_2_1 {
		contractVersion = "v2.1"
	} else if registryVersion == ethereum.RegistryVersion_2_0 {
		contractVersion = "v2.0"
	} else {
		return errors.New("v2.0 and v2.1 are the only supported versions")
	}

	bootstrapSpec := &client.OCR2TaskJobSpec{
		Name:    "ocr2 bootstrap node " + registryAddr,
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: registryAddr,
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID": int(chainID),
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
	_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
	if err != nil {
		l.Error().Err(err).Msg("Shouldn't fail creating bootstrap job on bootstrap node")
		return err
	}

	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, bootstrapNode.InternalIP(), 6690)
	for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
		nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].EthAddresses()
		if err != nil {
			l.Error().Err(err).Msgf("Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
			return err
		}
		nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCR2Keys()
		if err != nil {
			l.Error().Err(err).Msgf("Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
			return err
		}
		var nodeOCRKeyId []string
		for _, key := range nodeOCRKeys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				nodeOCRKeyId = append(nodeOCRKeyId, key.ID)
				break
			}
		}

		autoOCR2JobSpec := client.OCR2TaskJobSpec{
			Name:    "ocr2 " + registryAddr,
			JobType: "offchainreporting2",
			OCR2OracleSpec: job.OCR2OracleSpec{
				PluginType: "ocr2automation",
				Relay:      "evm",
				RelayConfig: map[string]interface{}{
					"chainID": int(chainID),
				},
				PluginConfig: map[string]interface{}{
					"mercuryCredentialName": "\"cred1\"",
					"contractVersion":       "\"" + contractVersion + "\"",
				},
				ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
				ContractID:                        registryAddr,                                      // registryAddr
				OCRKeyBundleID:                    null.StringFrom(nodeOCRKeyId[0]),                  // get node ocr2config.ID
				TransmitterID:                     null.StringFrom(nodeTransmitterAddress[keyIndex]), // node addr
				P2PV2Bootstrappers:                pq.StringArray{P2Pv2Bootstrapper},                 // bootstrap node key and address <p2p-key>@bootstrap:8000
			},
		}

		_, err = chainlinkNodes[nodeIndex].MustCreateJob(&autoOCR2JobSpec)
		if err != nil {
			l.Error().Err(err).Msgf("Shouldn't fail creating OCR Task job on OCR node %d err: %+v", nodeIndex+1, err)
			return err
		}

	}
	l.Info().Msg("Done creating OCR automation jobs")
	return nil
}
