package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_client_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	drconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func ptr[T any](v T) *T { return &v }

func SetOracleConfig(t *testing.T, owner *bind.TransactOpts, oracleContract *ocr2dr_oracle.OCR2DROracle, oracles []confighelper2.OracleIdentityExtra, batchSize int) {
	S := make([]int, len(oracles))
	for i := 0; i < len(S); i++ {
		S[i] = 1
	}

	reportingPluginConfigBytes, err := drconfig.EncodeReportingPluginConfig(&drconfig.ReportingPluginConfigWrapper{
		Config: &drconfig.ReportingPluginConfig{
			MaxQueryLengthBytes:       10_000,
			MaxObservationLengthBytes: 10_000,
			MaxReportLengthBytes:      10_000,
			MaxRequestBatchSize:       uint32(batchSize),
			DefaultAggregationMethod:  drconfig.AggregationMethod_AGGREGATION_MODE,
			UniqueReports:             true,
		},
	})
	require.NoError(t, err)

	signersKeys, transmittersAccounts, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress
		1*time.Second,        // deltaResend
		1*time.Second,        // deltaRound
		500*time.Millisecond, // deltaGrace
		2*time.Second,        // deltaStage
		5,                    // RMax (maxRounds)
		S,                    // S (schedule of randomized transmission order)
		oracles,
		reportingPluginConfigBytes,
		200*time.Millisecond, // maxDurationQuery
		200*time.Millisecond, // maxDurationObservation
		200*time.Millisecond, // maxDurationReport
		200*time.Millisecond, // maxDurationAccept
		200*time.Millisecond, // maxDurationTransmit
		1,                    // f (max faulty oracles)
		nil,                  // empty onChain config
	)

	var signers []common.Address
	var transmitters []common.Address
	for i := range signersKeys {
		signers = append(signers, common.BytesToAddress(signersKeys[i]))
		transmitters = append(transmitters, common.HexToAddress(string(transmittersAccounts[i])))
	}

	require.NoError(t, err)
	_, err = oracleContract.SetConfig(
		owner,
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	_, err = oracleContract.DeactivateAuthorizedReceiver(owner)
	require.NoError(t, err)
}

func SetRegistryConfig(t *testing.T, owner *bind.TransactOpts, registryContract *ocr2dr_registry.OCR2DRRegistry, oracleContractAddress common.Address) {
	var maxGasLimit = uint32(450_000)
	var stalenessSeconds = uint32(86_400)
	var gasAfterPaymentCalculation = big.NewInt(21_000 + 5_000 + 2_100 + 20_000 + 2*2_100 - 15_000 + 7_315)
	var weiPerUnitLink = big.NewInt(5000000000000000)
	var gasOverhead = uint32(500_000)
	var requestTimeoutSeconds = uint32(300)

	_, err := registryContract.SetConfig(
		owner,
		maxGasLimit,
		stalenessSeconds,
		gasAfterPaymentCalculation,
		weiPerUnitLink,
		gasOverhead,
		requestTimeoutSeconds,
	)
	require.NoError(t, err)

	var senders = []common.Address{oracleContractAddress}
	_, err = registryContract.SetAuthorizedSenders(
		owner,
		senders,
	)
	require.NoError(t, err)
}

func CreateAndFundSubscriptions(t *testing.T, owner *bind.TransactOpts, linkToken *link_token_interface.LinkToken, registryContractAddress common.Address, registryContract *ocr2dr_registry.OCR2DRRegistry, clientContracts []deployedClientContract) (subscriptionId uint64) {
	_, err := registryContract.CreateSubscription(owner)
	require.NoError(t, err)

	subscriptionID := uint64(1)

	numContracts := len(clientContracts)
	for i := 0; i < numContracts; i++ {
		_, err = registryContract.AddConsumer(owner, subscriptionID, clientContracts[i].Address)
		require.NoError(t, err)
	}

	data, err := utils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	require.NoError(t, err)

	amount := big.NewInt(0).Mul(big.NewInt(int64(numContracts)), big.NewInt(2e18)) // 2 LINK per client
	_, err = linkToken.TransferAndCall(owner, registryContractAddress, amount, data)
	require.NoError(t, err)

	time.Sleep(1000 * time.Millisecond)

	return subscriptionID
}

const finalityDepth int = 4

func CommitWithFinality(b *backends.SimulatedBackend) {
	for i := 0; i < finalityDepth; i++ {
		b.Commit()
	}
}

type deployedClientContract struct {
	Address  common.Address
	Contract *ocr2dr_client_example.OCR2DRClientExample
}

func StartNewChainWithContracts(t *testing.T, nClients int) (*bind.TransactOpts, *backends.SimulatedBackend, *time.Ticker, common.Address, *ocr2dr_oracle.OCR2DROracle, []deployedClientContract, common.Address, *ocr2dr_registry.OCR2DRRegistry, *link_token_interface.LinkToken) {
	owner := testutils.MustNewSimTransactor(t)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10) // 1 eth
	genesisData := core.GenesisAlloc{owner.From: {Balance: sb}}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
	b.Commit()

	// Deploy contracts
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(owner, b)
	require.NoError(t, err)

	linkEthFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(owner, b, 0, big.NewInt(5021530000000000))
	require.NoError(t, err)

	ocrContractAddress, _, ocrContract, err := ocr2dr_oracle.DeployOCR2DROracle(owner, b)
	require.NoError(t, err)

	registryAddress, _, registryContract, err := ocr2dr_registry.DeployOCR2DRRegistry(owner, b, linkAddr, linkEthFeedAddr, ocrContractAddress)
	require.NoError(t, err)

	_, err = ocrContract.SetRegistry(owner, registryAddress)
	require.NoError(t, err)

	clientContracts := []deployedClientContract{}
	for i := 0; i < nClients; i++ {
		clientContractAddress, _, clientContract, err := ocr2dr_client_example.DeployOCR2DRClientExample(owner, b, ocrContractAddress)
		require.NoError(t, err)
		clientContracts = append(clientContracts, deployedClientContract{
			Address:  clientContractAddress,
			Contract: clientContract,
		})
		if i%10 == 0 {
			// Max 10 requests per block
			b.Commit()
		}
	}
	CommitWithFinality(b)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			b.Commit()
		}
	}()
	return owner, b, ticker, ocrContractAddress, ocrContract, clientContracts, registryAddress, registryContract, linkToken
}

type Node struct {
	App            *cltest.TestApplication
	PeerID         string
	Transmitter    common.Address
	Keybundle      ocr2key.KeyBundle
	Config         config.GeneralConfig
	OracleIdentity confighelper2.OracleIdentityExtra
}

func StartNewNode(
	t *testing.T,
	owner *bind.TransactOpts,
	port uint16,
	dbName string,
	b *backends.SimulatedBackend,
	maxGas uint32,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
) *Node {
	p2pKey, err := p2pkey.NewV2()
	require.NoError(t, err)
	config, _ := heavyweight.FullTestDBV2(t, fmt.Sprintf("%s%d", dbName, port), func(c *chainlink.Config, s *chainlink.Secrets) {
		c.DevMode = true

		c.Feature.LogPoller = ptr(true)

		c.OCR.Enabled = ptr(false)
		c.OCR2.Enabled = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.V1.Enabled = ptr(false)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = models.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = models.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
		if len(p2pV2Bootstrappers) > 0 {
			c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
		}

		c.EVM[0].LogPollInterval = models.MustNewDuration(1 * time.Second)
		c.EVM[0].Transactions.ForwardersEnabled = ptr(false)
		c.EVM[0].GasEstimator.LimitDefault = ptr(maxGas)
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address

	// fund the transmitter address
	n, err := b.NonceAt(testutils.Context(t), owner.From, nil)
	require.NoError(t, err)

	tx := types.NewTransaction(
		n, transmitter,
		assets.Ether(1).ToInt(),
		21000,
		assets.GWei(1).ToInt(),
		nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	b.Commit()

	kb, err := app.GetKeyStore().OCR2().Create("evm")
	require.NoError(t, err)

	err = app.Start(testutils.Context(t))
	require.NoError(t, err)

	return &Node{
		App:         app,
		PeerID:      p2pKey.PeerID().Raw(),
		Transmitter: transmitter,
		Keybundle:   kb,
		Config:      config,
		OracleIdentity: confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  kb.PublicKey(),
				TransmitAccount:   ocrtypes2.Account(transmitter.String()),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            p2pKey.PeerID().Raw(),
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		},
	}
}

func AddBootstrapJob(t *testing.T, app *cltest.TestApplication, contractAddress common.Address) job.Job {
	job, err := ocrbootstrap.ValidatedBootstrapSpecToml(fmt.Sprintf(`
		type                              = "bootstrap"
		name                              = "dr-ocr-bootstrap"
		schemaVersion                     = 1
		relay                             = "evm"
		contractConfigConfirmations       = 1
		contractConfigTrackerPollInterval = "1s"
		contractID                        = "%s"

		[relayConfig]
		chainID                           = 1337
		fromBlock                         = 1
	`, contractAddress))
	require.NoError(t, err)
	err = app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	return job
}

func AddOCR2Job(t *testing.T, app *cltest.TestApplication, contractAddress common.Address, keyBundleID string, transmitter common.Address, bridgeURL string) job.Job {
	u, err := url.Parse(bridgeURL)
	require.NoError(t, err)
	require.NoError(t, app.BridgeORM().CreateBridgeType(&bridges.BridgeType{
		Name: "ea_bridge",
		URL:  models.WebURL(*u),
	}))
	job, err := validate.ValidatedOracleSpecToml(app.Config, fmt.Sprintf(`
		type               = "offchainreporting2"
		name               = "dr-ocr-node"
		schemaVersion      = 1
		relay              = "evm"
		contractID         = "%s"
		ocrKeyBundleID     = "%s"
		transmitterID      = "%s"
		contractConfigConfirmations = 1
		contractConfigTrackerPollInterval = "1s"
		maxTaskDuration    = "30s"
		pluginType         = "functions"
		observationSource  = """
			decode_log         [type="ethabidecodelog" abi="OracleRequest(bytes32 indexed requestId, address requestingContract, address requestInitiator, uint64 subscriptionId, address subscriptionOwner, bytes data)" data="$(jobRun.logData)" topics="$(jobRun.logTopics)"]
			decode_cbor        [type="cborparse" data="$(decode_log.data)"]
			run_computation    [type="bridge" name="ea_bridge" requestData="{\\"id\\": $(jobSpec.externalJobID), \\"data\\": $(decode_cbor)}"]
			parse_result       [type=jsonparse data="$(run_computation)" path="data,result"]
			parse_error        [type=jsonparse data="$(run_computation)" path="data,error"]

			decode_log -> decode_cbor -> run_computation -> parse_result -> parse_error
		"""

		[relayConfig]
		chainID = 1337
		fromBlock = 1

		[pluginConfig]
		minIncomingConfirmations = 3
		requestTimeoutSec = 300
		requestTimeoutCheckFrequencySec = 10
		requestTimeoutBatchLookupSize = 20
		listenerEventHandlerTimeoutSec = 120
		maxRequestSizeBytes = 30720
	`, contractAddress, keyBundleID, transmitter))
	require.NoError(t, err)
	err = app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	return job
}

func StartNewMockEA(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		var jsonMap map[string]any
		require.NoError(t, json.Unmarshal(b, &jsonMap))
		source := jsonMap["data"].(map[string]any)["source"].(string)
		res.WriteHeader(http.StatusOK)
		// prepend "0xab" to source and return as result
		res.Write([]byte(fmt.Sprintf(`{"data": {"result": "0xab%s", "error": ""}}`, source)))
	}))
}

// Mock EA prepends 0xab to source and user contract crops the answer to first 32 bytes
func GetExpectedResponse(source []byte) [32]byte {
	var resp [32]byte
	resp[0] = 0xab
	for j := 0; j < 31; j++ {
		if j >= len(source) {
			break
		}
		resp[j+1] = source[j]
	}
	return resp
}
