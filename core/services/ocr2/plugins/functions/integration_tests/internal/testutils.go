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
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_client_example"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	functionsConfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	DefaultSecretsBytes      = []byte{0xaa, 0xbb, 0xcc}
	DefaultSecretsBase64     = "qrvM"
	DefaultSecretsUrlsBytes  = []byte{0x01, 0x02, 0x03}
	DefaultSecretsUrlsBase64 = "AQID"
	// DefaultThresholdSecretsHex decrypts to the JSON string `{"0x0":"qrvM"}`
	DefaultThresholdSecretsHex = "0x7b225444483243747874223a2265794a48636d393163434936496c41794e5459694c434a44496a6f69533035305a5559325448593056553168543341766148637955584a4b65545a68626b3177527939794f464e78576a59356158646d636a6c4f535430694c434a4d59574a6c62434936496b464251554642515546425155464251554642515546425155464251554642515546425155464251554642515546425155464251554642515545394969776956534936496b4a45536c6c7a51334e7a623055334d6e6444574846474e557056634770585a573157596e565265544d796431526d4d32786c636c705a647a4671536e6c47627a5256615735744e6d773355456855546e6b7962324e746155686f626c51354d564a6a4e6e5230656c70766147644255326372545430694c434a5658324a6863694936496b4a4961544e69627a5a45536d396a4d324d344d6c46614d5852724c325645536b4a484d336c5a556d783555306834576d684954697472623264575a306f33546e4e456232314b5931646853544979616d63305657644f556c526e57465272655570325458706952306c4a617a466e534851314f4430694c434a46496a6f694d7a524956466c354d544e474b307836596e5a584e7a6c314d6d356c655574514e6b397a656e467859335253513239705a315534534652704e4430694c434a47496a6f69557a5132596d6c6952545a584b314176546d744252445677575459796148426862316c6c6330684853556869556c56614e303155556c6f345554306966513d3d222c2253796d43747874223a2253764237652f4a556a552b433358757873384e5378316967454e517759755051623730306a4a6144222c224e6f6e6365223a224d31714b557a6b306b77374767593538227d"
	DefaultArg1                = "arg1"
	DefaultArg2                = "arg2"
)

func ptr[T any](v T) *T { return &v }

func SetOracleConfig(t *testing.T, owner *bind.TransactOpts, oracleContract *ocr2dr_oracle.OCR2DROracle, oracles []confighelper2.OracleIdentityExtra, batchSize int) {
	S := make([]int, len(oracles))
	for i := 0; i < len(S); i++ {
		S[i] = 1
	}

	reportingPluginConfigBytes, err := functionsConfig.EncodeReportingPluginConfig(&functionsConfig.ReportingPluginConfigWrapper{
		Config: &functionsConfig.ReportingPluginConfig{
			MaxQueryLengthBytes:       10_000,
			MaxObservationLengthBytes: 10_000,
			MaxReportLengthBytes:      10_000,
			MaxRequestBatchSize:       uint32(batchSize),
			DefaultAggregationMethod:  functionsConfig.AggregationMethod_AGGREGATION_MODE,
			UniqueReports:             true,
			ThresholdPluginConfig: &functionsConfig.ThresholdReportingPluginConfig{
				MaxQueryLengthBytes:       10_000,
				MaxObservationLengthBytes: 10_000,
				MaxReportLengthBytes:      10_000,
				RequestCountLimit:         100,
				RequestTotalBytesLimit:    1_000,
				RequireLocalRequestCheck:  true,
			},
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
	ocr2Keystore []byte,
	thresholdKeyShare string,
) *Node {
	p2pKey, err := p2pkey.NewV2()
	require.NoError(t, err)
	config, _ := heavyweight.FullTestDBV2(t, fmt.Sprintf("%s%d", dbName, port), func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true)

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

		if len(thresholdKeyShare) > 0 {
			s.Threshold.ThresholdKeyShare = models.NewSecret(thresholdKeyShare)
		}
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

	var kb ocr2key.KeyBundle
	if ocr2Keystore != nil {
		kb, err = app.GetKeyStore().OCR2().Import(ocr2Keystore, "testPassword")
	} else {
		kb, err = app.GetKeyStore().OCR2().Create("evm")
	}
	require.NoError(t, err)

	err = app.Start(testutils.Context(t))
	require.NoError(t, err)

	return &Node{
		App:         app,
		PeerID:      p2pKey.PeerID().Raw(),
		Transmitter: transmitter,
		Keybundle:   kb,
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
		name                              = "functions-bootstrap"
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
	job, err := validate.ValidatedOracleSpecToml(app.Config.OCR2(), app.Config.Insecure(), fmt.Sprintf(`
		type               = "offchainreporting2"
		name               = "functions-node"
		schemaVersion      = 1
		relay              = "evm"
		contractID         = "%s"
		ocrKeyBundleID     = "%s"
		transmitterID      = "%s"
		contractConfigConfirmations = 1
		contractConfigTrackerPollInterval = "1s"
		pluginType         = "functions"
		observationSource  = """
			run_computation    [type="bridge" name="ea_bridge" requestData="{\\"note\\": \\"observationSource is unused but the bridge is required\\"}"]
			run_computation
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

			[pluginConfig.decryptionQueueConfig]
			completedCacheTimeoutSec = 300
			maxCiphertextBytes = 10_000
			maxCiphertextIdLength = 100
			maxQueueLength = 100
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
		var responsePayload []byte
		if jsonMap["endpoint"].(string) == "lambda" {
			responsePayload = mockEALambdaExecutionResponse(t, jsonMap)
		} else if jsonMap["endpoint"].(string) == "fetcher" {
			responsePayload = mockEASecretsFetchResponse(t, jsonMap)
		} else {
			require.Fail(t, "unknown external adapter endpoint '%s'", jsonMap["endpoint"].(string))
		}
		res.WriteHeader(http.StatusOK)
		_, err = res.Write(responsePayload)
		require.NoError(t, err)
	}))
}

func mockEALambdaExecutionResponse(t *testing.T, request map[string]any) []byte {
	data := request["data"].(map[string]any)
	require.Equal(t, functions.LanguageJavaScript, int(data["language"].(float64)))
	require.Equal(t, functions.LocationInline, int(data["codeLocation"].(float64)))
	require.Equal(t, functions.LocationRemote, int(data["secretsLocation"].(float64)))
	if data["secrets"] != DefaultSecretsBase64 && request["nodeProvidedSecrets"] != fmt.Sprintf(`{"0x0":"%s"}`, DefaultSecretsBase64) {
		assert.Fail(t, "expected secrets or nodeProvidedSecrets to be '%s'", DefaultSecretsBase64)
	}
	args := data["args"].([]interface{})
	require.Equal(t, 2, len(args))
	require.Equal(t, DefaultArg1, args[0].(string))
	require.Equal(t, DefaultArg2, args[1].(string))
	source := data["source"].(string)
	// prepend "0xab" to source and return as result
	return []byte(fmt.Sprintf(`{"result": "success", "statusCode": 200, "data": {"result": "0xab%s", "error": ""}}`, source))
}

func mockEASecretsFetchResponse(t *testing.T, request map[string]any) []byte {
	data := request["data"].(map[string]any)
	require.Equal(t, "fetchThresholdEncryptedSecrets", data["requestType"], "expected requestType to be 'fetchThresholdEncryptedSecrets' but got '%s'", data["requestType"])
	require.Equal(t, DefaultSecretsUrlsBase64, data["encryptedSecretsUrls"], "expected encryptedSecretsUrls to be '%s' but got '%s'", DefaultSecretsUrlsBase64, data["encryptedSecretsUrls"])
	return []byte(fmt.Sprintf(`{"result": "success", "statusCode": 200, "data": {"result": "%s", "error": ""}}`, DefaultThresholdSecretsHex))
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
