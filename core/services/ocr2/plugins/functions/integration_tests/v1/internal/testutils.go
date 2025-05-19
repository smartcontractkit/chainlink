package testutils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_allow_list"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_client_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
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
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

var nilOpts *bind.CallOpts

func ptr[T any](v T) *T { return &v }

var allowListPrivateKey = "0xae78c8b502571dba876742437f8bc78b689cf8518356c0921393d89caaf284ce"

func SetOracleConfig(t *testing.T, b evmtypes.Backend, owner *bind.TransactOpts, coordinatorContract *functions_coordinator.FunctionsCoordinator, oracles []confighelper2.OracleIdentityExtra, batchSize int, functionsPluginConfig *functionsConfig.ReportingPluginConfig) {
	S := make([]int, len(oracles))
	for i := 0; i < len(S); i++ {
		S[i] = 1
	}

	reportingPluginConfigBytes, err := functionsConfig.EncodeReportingPluginConfig(&functionsConfig.ReportingPluginConfigWrapper{
		Config: functionsPluginConfig,
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
		nil,
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

	_, err = coordinatorContract.SetConfig(
		owner,
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	client.FinalizeLatest(t, b)
}

func CreateAndFundSubscriptions(t *testing.T, b evmtypes.Backend, owner *bind.TransactOpts, linkToken *link_token_interface.LinkToken, routerContractAddress common.Address, routerContract *functions_router.FunctionsRouter, clientContracts []deployedClientContract, allowListContract *functions_allow_list.TermsOfServiceAllowList) (subscriptionID uint64) {
	allowed, err := allowListContract.HasAccess(nilOpts, owner.From, []byte{})
	require.NoError(t, err)
	if !allowed {
		message, err2 := allowListContract.GetMessage(nilOpts, owner.From, owner.From)
		require.NoError(t, err2)
		privateKey, err2 := crypto.HexToECDSA(allowListPrivateKey[2:])
		require.NoError(t, err2)
		flatSignature, err2 := crypto.Sign(message[:], privateKey)
		require.NoError(t, err2)
		var r [32]byte
		copy(r[:], flatSignature[:32])
		var s [32]byte
		copy(s[:], flatSignature[32:64])
		v := flatSignature[65]
		_, err2 = allowListContract.AcceptTermsOfService(owner, owner.From, owner.From, r, s, v)
		require.NoError(t, err2)
		b.Commit()
	}

	_, err = routerContract.CreateSubscription(owner)
	require.NoError(t, err)
	b.Commit()

	subscriptionID = uint64(1)

	numContracts := len(clientContracts)
	for i := 0; i < numContracts; i++ {
		_, err = routerContract.AddConsumer(owner, subscriptionID, clientContracts[i].Address)
		require.NoError(t, err)
		b.Commit()
	}

	data, err := utils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	require.NoError(t, err)

	amount := big.NewInt(0).Mul(big.NewInt(int64(numContracts)), big.NewInt(2e18)) // 2 LINK per client
	_, err = linkToken.TransferAndCall(owner, routerContractAddress, amount, data)
	require.NoError(t, err)
	b.Commit()

	return
}

type deployedClientContract struct {
	Address  common.Address
	Contract *functions_client_example.FunctionsClientExample
}

type Coordinator struct {
	Address  common.Address
	Contract *functions_coordinator.FunctionsCoordinator
}

func StartNewChainWithContracts(t *testing.T, nClients int) (*bind.TransactOpts, evmtypes.Backend, func() common.Hash, func(), Coordinator, Coordinator, []deployedClientContract, common.Address, *functions_router.FunctionsRouter, *link_token_interface.LinkToken, common.Address, *functions_allow_list.TermsOfServiceAllowList) {
	owner := evmtestutils.MustNewSimTransactor(t)
	owner.GasPrice = big.NewInt(int64(DefaultGasPrice))
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10) // 1 eth
	genesisData := types.GenesisAlloc{owner.From: {Balance: sb}}
	b := cltest.NewSimulatedBackend(t, genesisData, 2*ethconfig.Defaults.Miner.GasCeil)
	b.Commit()

	// Deploy LINK token
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(owner, b.Client())
	require.NoError(t, err)
	b.Commit()

	// Deploy mock LINK/ETH price feed
	linkEthFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(owner, b.Client(), 18, big.NewInt(5_000_000_000_000_000))
	require.NoError(t, err)
	b.Commit()

	// Deploy mock LINK/USD price feed
	linkUsdFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(owner, b.Client(), 18, big.NewInt(1_500_00_000))
	require.NoError(t, err)
	b.Commit()

	// Deploy Router contract
	handleOracleFulfillmentSelectorSlice, err := hex.DecodeString("0ca76175")
	require.NoError(t, err)
	var handleOracleFulfillmentSelector [4]byte
	copy(handleOracleFulfillmentSelector[:], handleOracleFulfillmentSelectorSlice[:4])
	functionsRouterConfig := functions_router.FunctionsRouterConfig{
		MaxConsumersPerSubscription:        uint16(100),
		AdminFee:                           big.NewInt(0),
		HandleOracleFulfillmentSelector:    handleOracleFulfillmentSelector,
		MaxCallbackGasLimits:               []uint32{300_000, 500_000, 1_000_000},
		GasForCallExactCheck:               5000,
		SubscriptionDepositMinimumRequests: 10,
		SubscriptionDepositJuels:           big.NewInt(9 * 1e18), // 9 LINK
	}
	routerAddress, _, routerContract, err := functions_router.DeployFunctionsRouter(owner, b.Client(), linkAddr, functionsRouterConfig)
	require.NoError(t, err)
	b.Commit()

	// Deploy Allow List contract
	privateKey, err := crypto.HexToECDSA(allowListPrivateKey[2:])
	proofSignerPublicKey := crypto.PubkeyToAddress(privateKey.PublicKey)
	require.NoError(t, err)
	allowListConfig := functions_allow_list.TermsOfServiceAllowListConfig{
		Enabled:         false, // TODO: true
		SignerPublicKey: proofSignerPublicKey,
	}
	var initialAllowedSenders []common.Address
	var initialBlockedSenders []common.Address
	// The allowlist requires a pointer to the previous allowlist. If none exists, use the null address.
	var nullPreviousAllowlist common.Address
	allowListAddress, _, allowListContract, err := functions_allow_list.DeployTermsOfServiceAllowList(owner, b.Client(), allowListConfig, initialAllowedSenders, initialBlockedSenders, nullPreviousAllowlist)
	require.NoError(t, err)
	b.Commit()

	// Deploy Coordinator contract (matches updateConfig() in FunctionsBilling.sol)
	coordinatorConfig := functions_coordinator.FunctionsBillingConfig{
		FeedStalenessSeconds:                uint32(86_400),
		GasOverheadBeforeCallback:           uint32(325_000),
		GasOverheadAfterCallback:            uint32(50_000),
		RequestTimeoutSeconds:               uint32(300),
		DonFeeCentsUsd:                      uint16(0),
		MaxSupportedRequestDataVersion:      uint16(1),
		FulfillmentGasPriceOverEstimationBP: uint32(1_000),
		FallbackNativePerUnitLink:           big.NewInt(5_000_000_000_000_000),
		MinimumEstimateGasPriceWei:          big.NewInt(1_000_000_000),
		OperationFeeCentsUsd:                uint16(0),
		FallbackUsdPerUnitLink:              uint64(1_400_000_000),
		FallbackUsdPerUnitLinkDecimals:      uint8(8),
		TransmitTxSizeBytes:                 uint16(1764),
	}
	require.NoError(t, err)
	coordinatorAddress, _, coordinatorContract, err := functions_coordinator.DeployFunctionsCoordinator(owner, b.Client(), routerAddress, coordinatorConfig, linkEthFeedAddr, linkUsdFeedAddr)
	require.NoError(t, err)
	b.Commit()
	proposalAddress, _, proposalContract, err := functions_coordinator.DeployFunctionsCoordinator(owner, b.Client(), routerAddress, coordinatorConfig, linkEthFeedAddr, linkUsdFeedAddr)
	require.NoError(t, err)
	b.Commit()

	// Deploy Client contracts
	clientContracts := []deployedClientContract{}
	for i := 0; i < nClients; i++ {
		clientContractAddress, _, clientContract, err := functions_client_example.DeployFunctionsClientExample(owner, b.Client(), routerAddress)
		require.NoError(t, err)
		b.Commit()
		clientContracts = append(clientContracts, deployedClientContract{
			Address:  clientContractAddress,
			Contract: clientContract,
		})
		if i%10 == 0 {
			// Max 10 requests per block
			b.Commit()
		}
	}

	client.FinalizeLatest(t, b)
	commit, stop := cltest.Mine(b, time.Second)

	active := Coordinator{
		Contract: coordinatorContract,
		Address:  coordinatorAddress,
	}
	proposed := Coordinator{
		Contract: proposalContract,
		Address:  proposalAddress,
	}
	return owner, b, commit, stop, active, proposed, clientContracts, routerAddress, routerContract, linkToken, allowListAddress, allowListContract
}

func SetupRouterRoutes(t *testing.T, b evmtypes.Backend, owner *bind.TransactOpts, routerContract *functions_router.FunctionsRouter, coordinatorAddress common.Address, proposedCoordinatorAddress common.Address, allowListAddress common.Address) {
	allowListId, err := routerContract.GetAllowListId(nilOpts)
	require.NoError(t, err)
	var donId [32]byte
	copy(donId[:], DefaultDONId)
	proposedContractSetIds := []([32]byte){allowListId, donId}
	proposedContractSetAddresses := []common.Address{allowListAddress, coordinatorAddress}
	_, err = routerContract.ProposeContractsUpdate(owner, proposedContractSetIds, proposedContractSetAddresses)
	require.NoError(t, err)

	b.Commit()

	_, err = routerContract.UpdateContracts(owner)
	require.NoError(t, err)
	b.Commit()

	// prepare next coordinator
	proposedContractSetIds = []([32]byte){donId}
	proposedContractSetAddresses = []common.Address{proposedCoordinatorAddress}
	_, err = routerContract.ProposeContractsUpdate(owner, proposedContractSetIds, proposedContractSetAddresses)
	require.NoError(t, err)
	b.Commit()
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
	port int,
	b evmtypes.Backend,
	maxGas uint32,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
	ocr2Keystore []byte,
	thresholdKeyShare string,
) *Node {
	ctx := testutils.Context(t)
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(port)))
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true)

		c.Feature.LogPoller = ptr(true)

		c.OCR.Enabled = ptr(false)
		c.OCR2.Enabled = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
		if len(p2pV2Bootstrappers) > 0 {
			c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
		}

		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
		c.EVM[0].Transactions.ForwardersEnabled = ptr(false)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint64(maxGas))
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
		c.EVM[0].GasEstimator.PriceDefault = assets.NewWei(big.NewInt(int64(DefaultGasPrice)))

		if len(thresholdKeyShare) > 0 {
			s.Threshold.ThresholdKeyShare = models.NewSecret(thresholdKeyShare)
		}
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address

	// fund the transmitter address
	n, err := b.Client().NonceAt(testutils.Context(t), owner.From, nil)
	require.NoError(t, err)

	tx := evmtestutils.NewLegacyTransaction(
		n, transmitter,
		assets.Ether(1).ToInt(),
		21000,
		assets.GWei(1).ToInt(),
		nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.Client().SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	b.Commit()

	var kb ocr2key.KeyBundle
	if ocr2Keystore != nil {
		kb, err = app.GetKeyStore().OCR2().Import(ctx, ocr2Keystore, "testPassword")
	} else {
		kb, err = app.GetKeyStore().OCR2().Create(ctx, "evm")
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
		donID                             = "%s"
		contractVersion                   = 1
		contractUpdateCheckFrequencySec   = 1

	`, contractAddress, DefaultDONId))
	require.NoError(t, err)
	err = app.AddJobV2(testutils.Context(t), &job)
	require.NoError(t, err)
	return job
}

func AddOCR2Job(t *testing.T, app *cltest.TestApplication, contractAddress common.Address, keyBundleID string, transmitter common.Address, bridgeURL string) job.Job {
	ctx := testutils.Context(t)
	u, err := url.Parse(bridgeURL)
	require.NoError(t, err)
	require.NoError(t, app.BridgeORM().CreateBridgeType(ctx, &bridges.BridgeType{
		Name: "ea_bridge",
		URL:  models.WebURL(*u),
	}))
	job, err := validate.ValidatedOracleSpecToml(testutils.Context(t), app.Config.OCR2(), app.Config.Insecure(), fmt.Sprintf(`
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
		donID = "%s"
		contractVersion = 1
		minIncomingConfirmations = 3
		requestTimeoutSec = 300
		requestTimeoutCheckFrequencySec = 10
		requestTimeoutBatchLookupSize = 20
		listenerEventHandlerTimeoutSec = 120
		listenerEventsCheckFrequencyMillis = 1000
		maxRequestSizeBytes = 30720
		contractUpdateCheckFrequencySec = 1

			[pluginConfig.decryptionQueueConfig]
			completedCacheTimeoutSec = 300
			maxCiphertextBytes = 10_000
			maxCiphertextIdLength = 100
			maxQueueLength = 100
			decryptRequestTimeoutSec = 100

			[pluginConfig.s4Constraints]
			maxPayloadSizeBytes = 10_1000
			maxSlotsPerUser = 10
	`, contractAddress, keyBundleID, transmitter, DefaultDONId), nil)
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
	if len(request["nodeProvidedSecrets"].(string)) > 0 {
		require.Equal(t, functions.LocationRemote, int(data["secretsLocation"].(float64)))
		require.JSONEq(t, fmt.Sprintf(`{"0x0":"%s"}`, DefaultSecretsBase64), request["nodeProvidedSecrets"].(string))
	}
	args := data["args"].([]interface{})
	require.Len(t, args, 2)
	require.Equal(t, DefaultArg1, args[0].(string))
	require.Equal(t, DefaultArg2, args[1].(string))
	source := data["source"].(string)
	// prepend "0xab" to source and return as result
	return []byte(fmt.Sprintf(`{"result": "success", "statusCode": 200, "data": {"result": "0xab%s", "error": ""}}`, source))
}

func mockEASecretsFetchResponse(t *testing.T, request map[string]any) []byte {
	data := request["data"].(map[string]any)
	require.Equal(t, "fetchThresholdEncryptedSecrets", data["requestType"])
	require.Equal(t, DefaultSecretsUrlsBase64, data["encryptedSecretsUrls"])
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

func CreateFunctionsNodes(
	t *testing.T,
	owner *bind.TransactOpts,
	b evmtypes.Backend,
	routerAddress common.Address,
	nOracleNodes int,
	maxGas int,
	ocr2Keystores [][]byte,
	thresholdKeyShares []string,
) (bootstrapNode *Node, oracleNodes []*cltest.TestApplication, oracleIdentites []confighelper2.OracleIdentityExtra) {
	t.Helper()

	if len(thresholdKeyShares) != 0 && len(thresholdKeyShares) != nOracleNodes {
		require.Fail(t, "thresholdKeyShares must be empty or have length equal to nOracleNodes")
	}
	if len(ocr2Keystores) != 0 && len(ocr2Keystores) != nOracleNodes {
		require.Fail(t, "ocr2Keystores must be empty or have length equal to nOracleNodes")
	}
	if len(ocr2Keystores) != len(thresholdKeyShares) {
		require.Fail(t, "ocr2Keystores and thresholdKeyShares must have the same length")
	}

	bootstrapPort := freeport.GetOne(t)
	bootstrapNode = StartNewNode(t, owner, bootstrapPort, b, uint32(maxGas), nil, nil, "")
	AddBootstrapJob(t, bootstrapNode.App, routerAddress)

	// oracle nodes with jobs, bridges and mock EAs
	ports := freeport.GetN(t, nOracleNodes)
	for i := 0; i < nOracleNodes; i++ {
		var thresholdKeyShare string
		if len(thresholdKeyShares) == 0 {
			thresholdKeyShare = ""
		} else {
			thresholdKeyShare = thresholdKeyShares[i]
		}
		var ocr2Keystore []byte
		if len(ocr2Keystores) == 0 {
			ocr2Keystore = nil
		} else {
			ocr2Keystore = ocr2Keystores[i]
		}
		oracleNode := StartNewNode(t, owner, ports[i], b, uint32(maxGas), []commontypes.BootstrapperLocator{
			{PeerID: bootstrapNode.PeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapPort)}},
		}, ocr2Keystore, thresholdKeyShare)
		oracleNodes = append(oracleNodes, oracleNode.App)
		oracleIdentites = append(oracleIdentites, oracleNode.OracleIdentity)

		ea := StartNewMockEA(t)
		t.Cleanup(ea.Close)

		_ = AddOCR2Job(t, oracleNodes[i], routerAddress, oracleNode.Keybundle.ID(), oracleNode.Transmitter, ea.URL)
	}

	return bootstrapNode, oracleNodes, oracleIdentites
}

func ClientTestRequests(t *testing.T, owner *bind.TransactOpts, b evmtypes.Backend, clientContracts []deployedClientContract, requestLenBytes int, expectedSecrets []byte, subscriptionID uint64, timeout time.Duration) {
	t.Helper()
	var donId [32]byte
	copy(donId[:], []byte(DefaultDONId))
	// send requests
	requestSources := make([][]byte, len(clientContracts))
	rnd := rand.New(rand.NewSource(666))
	for i, client := range clientContracts {
		requestSources[i] = make([]byte, requestLenBytes)
		for j := 0; j < requestLenBytes; j++ {
			requestSources[i][j] = byte(rnd.Uint32() % 256)
		}
		_, err := client.Contract.SendRequest(
			owner,
			hex.EncodeToString(requestSources[i]),
			expectedSecrets,
			[]string{DefaultArg1, DefaultArg2},
			subscriptionID,
			donId,
		)
		require.NoError(t, err)
	}
	client.FinalizeLatest(t, b)

	// validate that all client contracts got correct responses to their requests
	var wg sync.WaitGroup
	for i := 0; i < len(clientContracts); i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			gomega.NewGomegaWithT(t).Eventually(func() [32]byte {
				answer, err := clientContracts[ic].Contract.SLastResponse(nil)
				require.NoError(t, err)
				return answer
			}, timeout, 1*time.Second).Should(gomega.Equal(GetExpectedResponse(requestSources[ic])))
		}()
	}
	wg.Wait()
}
