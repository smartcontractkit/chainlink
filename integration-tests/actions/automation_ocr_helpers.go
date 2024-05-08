package actions

//revive:disable:dot-imports
import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	ocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func BuildAutoOCR2ConfigVars(
	t *testing.T,
	chainlinkNodes []*client.ChainlinkK8sClient,
	registryConfig contracts.KeeperRegistrySettings,
	registrar string,
	deltaStage time.Duration,
	chainModuleAddress common.Address,
	reorgProtectionEnabled bool,
) (contracts.OCRv2Config, error) {
	return BuildAutoOCR2ConfigVarsWithKeyIndex(t, chainlinkNodes, registryConfig, registrar, deltaStage, 0, common.Address{}, chainModuleAddress, reorgProtectionEnabled)
}

func BuildAutoOCR2ConfigVarsWithKeyIndex(
	t *testing.T,
	chainlinkNodes []*client.ChainlinkK8sClient,
	registryConfig contracts.KeeperRegistrySettings,
	registrar string,
	deltaStage time.Duration,
	keyIndex int,
	registryOwnerAddress common.Address,
	chainModuleAddress common.Address,
	reorgProtectionEnabled bool,
) (contracts.OCRv2Config, error) {
	l := logging.GetTestLogger(t)
	S, oracleIdentities, err := GetOracleIdentitiesWithKeyIndex(chainlinkNodes, keyIndex)
	if err != nil {
		return contracts.OCRv2Config{}, err
	}

	var offC []byte
	var signerOnchainPublicKeys []types.OnchainPublicKey
	var transmitterAccounts []types.Account
	var f uint8
	var offchainConfigVersion uint64
	var offchainConfig []byte

	if registryConfig.RegistryVersion == ethereum.RegistryVersion_2_1 || registryConfig.RegistryVersion == ethereum.RegistryVersion_2_2 {
		offC, err = json.Marshal(ocr2keepers30config.OffchainConfig{
			TargetProbability:    "0.999",
			TargetInRounds:       1,
			PerformLockoutWindow: 3600000, // Intentionally set to be higher than in prod for testing purpose
			GasLimitPerReport:    5_300_000,
			GasOverheadPerUpkeep: 300_000,
			MinConfirmations:     0,
			MaxUpkeepBatchSize:   10,
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
		require.Equal(t, 20, len(signer), "OnChainPublicKey '%v' has wrong length for address", signer)
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		require.True(t, common.IsHexAddress(string(transmitter)), "TransmitAccount '%s' is not a valid Ethereum address", string(transmitter))
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	ocrConfig := contracts.OCRv2Config{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	if registryConfig.RegistryVersion == ethereum.RegistryVersion_2_0 {
		ocrConfig.OnchainConfig = registryConfig.Encode20OnchainConfig(registrar)
	} else if registryConfig.RegistryVersion == ethereum.RegistryVersion_2_1 {
		ocrConfig.TypedOnchainConfig21 = registryConfig.Create21OnchainConfig(registrar, registryOwnerAddress)
	} else if registryConfig.RegistryVersion == ethereum.RegistryVersion_2_2 {
		ocrConfig.TypedOnchainConfig22 = registryConfig.Create22OnchainConfig(registrar, registryOwnerAddress, chainModuleAddress, reorgProtectionEnabled)
	}

	l.Info().Msg("Done building OCR config")
	return ocrConfig, nil
}

// CreateOCRKeeperJobs bootstraps the first node and to the other nodes sends ocr jobs
func CreateOCRKeeperJobs(
	t *testing.T,
	chainlinkNodes []*client.ChainlinkK8sClient,
	registryAddr string,
	chainID int64,
	keyIndex int,
	registryVersion ethereum.KeeperRegistryVersion,
) {
	l := logging.GetTestLogger(t)
	bootstrapNode := chainlinkNodes[0]
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	require.NoError(t, err, "Shouldn't fail reading P2P keys from bootstrap node")
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID

	var contractVersion string
	if registryVersion == ethereum.RegistryVersion_2_2 {
		contractVersion = "v2.1+"
	} else if registryVersion == ethereum.RegistryVersion_2_1 {
		contractVersion = "v2.1"
	} else if registryVersion == ethereum.RegistryVersion_2_0 {
		contractVersion = "v2.0"
	} else {
		require.FailNow(t, fmt.Sprintf("v2.0, v2.1, and v2.2 are the only supported versions, but got something else: %v (iota)", registryVersion))
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
	require.NoError(t, err, "Shouldn't fail creating bootstrap job on bootstrap node")
	// TODO: Use service name returned by chainlink-env once that is available
	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s-node-1:%d", bootstrapP2PId, bootstrapNode.Name(), 6690)

	for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
		nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].EthAddresses()
		require.NoError(t, err, "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
		nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCR2Keys()
		require.NoError(t, err, "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
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
		require.NoError(t, err, "Shouldn't fail creating OCR Task job on OCR node %d err: %+v", nodeIndex+1, err)
	}
	l.Info().Msg("Done creating OCR automation jobs")
}

// DeployAutoOCRRegistryAndRegistrar registry and registrar
func DeployAutoOCRRegistryAndRegistrar(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar) {
	registry := deployRegistry(t, registryVersion, registrySettings, contractDeployer, client, linkToken)
	registrar := deployRegistrar(t, registryVersion, registry, linkToken, contractDeployer, client)

	return registry, registrar
}

func DeployConsumers(t *testing.T, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, linkToken contracts.LinkToken, contractDeployer contracts.ContractDeployer, client blockchain.EVMClient, numberOfUpkeeps int, linkFundsForEachUpkeep *big.Int, upkeepGasLimit uint32, isLogTrigger bool, isMercury bool) ([]contracts.KeeperConsumer, []*big.Int) {
	upkeeps := DeployKeeperConsumers(t, contractDeployer, client, numberOfUpkeeps, isLogTrigger, isMercury)
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(
		t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, isLogTrigger, isMercury,
	)
	return upkeeps, upkeepIds
}

func DeployPerformanceConsumers(
	t *testing.T,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
	upkeepGasLimit uint32,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) ([]contracts.KeeperConsumerPerformance, []*big.Int) {
	upkeeps := DeployKeeperConsumersPerformance(
		t, contractDeployer, client, numberOfUpkeeps, blockRange, blockInterval, checkGasToBurn, performGasToBurn,
	)
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, false, false)
	return upkeeps, upkeepIds
}

func DeployPerformDataCheckerConsumers(
	t *testing.T,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
	upkeepGasLimit uint32,
	expectedData []byte,
) ([]contracts.KeeperPerformDataChecker, []*big.Int) {
	upkeeps := DeployPerformDataChecker(t, contractDeployer, client, numberOfUpkeeps, expectedData)
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, false, false)
	return upkeeps, upkeepIds
}

func deployRegistrar(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registry contracts.KeeperRegistry,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
) contracts.KeeperRegistrar {
	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar, err := contractDeployer.DeployKeeperRegistrar(registryVersion, linkToken.Address(), registrarSettings)
	require.NoError(t, err, "Deploying KeeperRegistrar contract shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for registrar to deploy")
	return registrar
}

func deployRegistry(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	linkToken contracts.LinkToken,
) contracts.KeeperRegistry {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for mock feeds to deploy")

	// Deploy the transcoder here, and then set it to the registry
	transcoder := DeployUpkeepTranscoder(t, contractDeployer, client)
	registry := DeployKeeperRegistry(t, contractDeployer, client,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  transcoder.Address(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        registrySettings,
		},
	)
	return registry
}
