package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	types2 "github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func BuildAutoOCR2ConfigVars(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	registryConfig contracts.KeeperRegistrySettings,
	registrar string,
	deltaStage time.Duration,
) contracts.OCRConfig {
	return BuildAutoOCR2ConfigVarsWithKeyIndex(t, chainlinkNodes, registryConfig, registrar, deltaStage, 0)
}

func BuildAutoOCR2ConfigVarsWithKeyIndex(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	registryConfig contracts.KeeperRegistrySettings,
	registrar string,
	deltaStage time.Duration,
	keyIndex int,
) contracts.OCRConfig {
	l := utils.GetTestLogger(t)
	S, oracleIdentities := getOracleIdentitiesWithKeyIndex(t, chainlinkNodes, keyIndex)

	signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		10*time.Second,        // deltaProgress time.Duration,
		15*time.Second,        // deltaResend time.Duration,
		3000*time.Millisecond, // deltaRound time.Duration,
		200*time.Millisecond,  // deltaGrace time.Duration,
		deltaStage,            // deltaStage time.Duration,
		24,                    // rMax uint8,
		S,                     // s []int,
		oracleIdentities,      // oracles []OracleIdentityExtra,
		types2.OffchainConfig{
			TargetProbability:    "0.999",
			TargetInRounds:       1,
			PerformLockoutWindow: 3600000, // Intentionally set to be higher than in prod for testing purpose
			GasLimitPerReport:    5_300_000,
			GasOverheadPerUpkeep: 300_000,
			SamplingJobDuration:  3000,
			MinConfirmations:     0,
			MaxUpkeepBatchSize:   1,
		}.Encode(), // reportingPluginConfig []byte,
		20*time.Millisecond,   // maxDurationQuery time.Duration,
		20*time.Millisecond,   // maxDurationObservation time.Duration,
		1200*time.Millisecond, // maxDurationReport time.Duration,
		20*time.Millisecond,   // maxDurationShouldAcceptFinalizedReport time.Duration,
		20*time.Millisecond,   // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                     // f int,
		nil,                   // onchainConfig []byte,
	)
	require.NoError(t, err, "Shouldn't fail ContractSetConfigArgsForTests")

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		require.Equal(t, 20, len(signer), "OnChainPublicKey has wrong length for address")
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		require.True(t, common.IsHexAddress(string(transmitter)), "TransmitAccount is not a valid Ethereum address")
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	onchainConfig, err := registryConfig.EncodeOnChainConfig(registrar)
	require.NoError(t, err, "Shouldn't fail encoding config")

	l.Info().Msg("Done building OCR config")
	return contracts.OCRConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

// CreateOCRKeeperJobs bootstraps the first node and to the other nodes sends ocr jobs
func CreateOCRKeeperJobs(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	registryAddr string,
	chainID int64,
	keyIndex int,
) {
	l := utils.GetTestLogger(t)
	bootstrapNode := chainlinkNodes[0]
	bootstrapNode.RemoteIP()
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	require.NoError(t, err, "Shouldn't fail reading P2P keys from bootstrap node")
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID

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
	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, bootstrapNode.RemoteIP(), 6690)

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
				PluginConfig:                      map[string]interface{}{},
				ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
				ContractID:                        registryAddr,                                      // registryAddr
				OCRKeyBundleID:                    null.StringFrom(nodeOCRKeyId[0]),                  // get node ocr2config.ID
				TransmitterID:                     null.StringFrom(nodeTransmitterAddress[keyIndex]), // node addr
				P2PV2Bootstrappers:                pq.StringArray{P2Pv2Bootstrapper},                 // bootstrap node key and address <p2p-key>@bootstrap:8000
			},
		}

		_, err = chainlinkNodes[nodeIndex].MustCreateJob(&autoOCR2JobSpec)
		require.NoError(t, err, "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
	}
	l.Info().Msg("Done creating OCR automation jobs")
}

// DeployAutoOCRRegistryAndRegistrar registry and registrar
func DeployAutoOCRRegistryAndRegistrar(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	numberOfUpkeeps int,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar) {
	registry := deployRegistry(t, registryVersion, registrySettings, contractDeployer, client, linkToken)

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err := linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfUpkeeps))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrar := deployRegistrar(t, registryVersion, registry, linkToken, contractDeployer, client)

	return registry, registrar
}

func DeployConsumers(
	t *testing.T,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
	upkeepGasLimit uint32,
) ([]contracts.KeeperConsumer, []*big.Int) {
	upkeeps := DeployKeeperConsumers(t, contractDeployer, client, numberOfUpkeeps)
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(
		t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses,
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
	upkeepIds := RegisterUpkeepContracts(
		t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses,
	)
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
	upkeepIds := RegisterUpkeepContracts(
		t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses,
	)
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
