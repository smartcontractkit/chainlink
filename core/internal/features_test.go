package internal_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/consumer_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/multiwordconsumer_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers/testoffchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

var oneETH = assets.Eth(*big.NewInt(1000000000000000000))

func TestIntegration_ExternalInitiatorV2(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()

	cfg := cltest.NewTestGeneralConfig(t)
	cfg.Overrides.FeatureExternalInitiators = null.BoolFrom(true)
	cfg.Overrides.SetTriggerFallbackDBPollInterval(10 * time.Millisecond)

	app, cleanup := cltest.NewApplicationWithConfig(t, cfg, ethClient, cltest.UseRealExternalInitiatorManager)
	defer cleanup()

	require.NoError(t, app.Start())

	var (
		eiName    = "substrate-ei"
		eiSpec    = map[string]interface{}{"foo": "bar"}
		eiRequest = map[string]interface{}{"result": 42}

		jobUUID = uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46")

		expectedCreateJobRequest = map[string]interface{}{
			"jobId":  jobUUID.String(),
			"type":   eiName,
			"params": eiSpec,
		}
	)

	// Setup EI
	var eiURL string
	var eiNotifiedOfCreate bool
	var eiNotifiedOfDelete bool
	{
		mockEI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !eiNotifiedOfCreate {
				require.Equal(t, http.MethodPost, r.Method)

				eiNotifiedOfCreate = true
				defer r.Body.Close()

				var gotCreateJobRequest map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&gotCreateJobRequest)
				require.NoError(t, err)

				require.Equal(t, expectedCreateJobRequest, gotCreateJobRequest)
				w.WriteHeader(http.StatusOK)
			} else {
				require.Equal(t, http.MethodDelete, r.Method)

				eiNotifiedOfDelete = true
				defer r.Body.Close()

				require.Equal(t, fmt.Sprintf("/%v", jobUUID.String()), r.URL.Path)
			}
		}))
		defer mockEI.Close()
		eiURL = mockEI.URL
	}

	// Create the EI record on the Core node
	var eia *auth.Token
	{
		eiCreate := map[string]string{
			"name": eiName,
			"url":  eiURL,
		}
		eiCreateJSON, err := json.Marshal(eiCreate)
		require.NoError(t, err)
		eip := cltest.CreateExternalInitiatorViaWeb(t, app, string(eiCreateJSON))
		eia = &auth.Token{
			AccessKey: eip.AccessKey,
			Secret:    eip.Secret,
		}
	}

	// Create the bridge on the Core node
	var bridgeCalled bool
	{
		bridgeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bridgeCalled = true
			defer r.Body.Close()

			var gotBridgeRequest map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&gotBridgeRequest)
			require.NoError(t, err)

			expectedBridgeRequest := map[string]interface{}{
				"value": float64(42),
			}
			require.Equal(t, expectedBridgeRequest, gotBridgeRequest)

			w.WriteHeader(http.StatusOK)
			require.NoError(t, err)
			io.WriteString(w, `{}`)
		}))
		u, _ := url.Parse(bridgeServer.URL)
		app.Store.CreateBridgeType(&models.BridgeType{
			Name: models.TaskType("substrate-adapter1"),
			URL:  models.WebURL(*u),
		})
		defer bridgeServer.Close()
	}

	// Create the job spec on the Core node
	var jobID int32
	{
		tomlSpec := fmt.Sprintf(`
type            = "webhook"
schemaVersion   = 1
externalJobID           = "%v"
externalInitiators = [
	{
		name = "%s",
		spec = """
	%s
"""
	}
]
observationSource   = """
    parse  [type=jsonparse path="result" data="$(jobRun.requestBody)"]
    submit [type=bridge name="substrate-adapter1" requestData=<{ "value": $(parse) }>]
    parse -> submit
"""
    `, jobUUID, eiName, cltest.MustJSONMarshal(t, eiSpec))

		_, err := webhook.ValidatedWebhookSpec(tomlSpec, app.GetExternalInitiatorManager())
		require.NoError(t, err)
		job := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: tomlSpec})))
		jobID = job.ID
		t.Log("JOB created", job.WebhookSpecID)

		require.Eventually(t, func() bool { return eiNotifiedOfCreate }, 5*time.Second, 10*time.Millisecond, "expected external initiator to be notified of new job")
	}

	t.Run("calling webhook_spec with non-matching external_initiator_id returns unauthorized", func(t *testing.T) {
		eiaWrong := auth.NewToken()
		body := cltest.MustJSONMarshal(t, eiRequest)
		headers := make(map[string]string)
		headers[static.ExternalInitiatorAccessKeyHeader] = eiaWrong.AccessKey
		headers[static.ExternalInitiatorSecretHeader] = eiaWrong.Secret

		url := app.Config.ClientNodeURL() + "/v2/jobs/" + jobUUID.String() + "/runs"
		bodyBuf := bytes.NewBufferString(body)
		resp, cleanup := cltest.UnauthenticatedPost(t, url, bodyBuf, headers)
		defer cleanup()
		cltest.AssertServerResponse(t, resp, 401)

		cltest.AssertCountStays(t, app.Store, &pipeline.Run{}, 0)
	})

	t.Run("calling webhook_spec with matching external_initiator_id works", func(t *testing.T) {
		// Simulate request from EI -> Core node
		cltest.AwaitJobActive(t, app.JobSpawner(), jobID, 3*time.Second)

		_ = cltest.CreateJobRunViaExternalInitiatorV2(t, app, jobUUID, *eia, cltest.MustJSONMarshal(t, eiRequest))

		pipelineORM := pipeline.NewORM(app.Store.DB)
		jobORM := job.NewORM(app.Store.ORM.DB, app.GetChainSet(), pipelineORM, &postgres.NullEventBroadcaster{}, &postgres.NullAdvisoryLocker{}, app.KeyStore)

		runs := cltest.WaitForPipelineComplete(t, 0, jobID, 1, 2, jobORM, 5*time.Second, 300*time.Millisecond)
		require.Len(t, runs, 1)
		run := runs[0]
		require.Len(t, run.PipelineTaskRuns, 2)
		require.Empty(t, run.PipelineTaskRuns[0].Error)
		require.Empty(t, run.PipelineTaskRuns[1].Error)

		assert.True(t, bridgeCalled, "expected bridge server to be called")
	})

	// Delete the job
	{
		cltest.DeleteJobViaWeb(t, app, jobID)
		require.Eventually(t, func() bool { return eiNotifiedOfDelete }, 5*time.Second, 10*time.Millisecond, "expected external initiator to be notified of deleted job")
	}
}

func TestIntegration_AuthToken(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMockCalls := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()

	require.NoError(t, app.Start())

	// set up user
	mockUser := cltest.MustRandomUser()
	apiToken := auth.Token{AccessKey: cltest.APIKey, Secret: cltest.APISecret}
	require.NoError(t, mockUser.SetAuthToken(&apiToken))
	require.NoError(t, app.Store.SaveUser(&mockUser))

	url := app.Config.ClientNodeURL() + "/v2/config"
	headers := make(map[string]string)
	headers[web.APIKey] = cltest.APIKey
	headers[web.APISecret] = cltest.APISecret
	buf := bytes.NewBufferString(`{"ethGasPriceDefault":150000000000}`)

	resp, cleanup := cltest.UnauthenticatedPatch(t, url, buf, headers)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
}

type OperatorContracts struct {
	user                      *bind.TransactOpts
	multiWordConsumerAddress  common.Address
	singleWordConsumerAddress common.Address
	operatorAddress           common.Address
	linkToken                 *link_token_interface.LinkToken
	multiWord                 *multiwordconsumer_wrapper.MultiWordConsumer
	singleWord                *consumer_wrapper.Consumer
	operator                  *operator_wrapper.Operator
	sim                       *backends.SimulatedBackend
}

func setupOperatorContracts(t *testing.T) OperatorContracts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	user := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10)
	genesisData := core.GenesisAlloc{
		user.From: {Balance: sb}, // 1 eth
	}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(user, b)
	require.NoError(t, err)
	b.Commit()

	operatorAddress, _, operatorContract, err := operator_wrapper.DeployOperator(user, b, linkTokenAddress, user.From)
	require.NoError(t, err)
	b.Commit()

	var empty [32]byte
	multiWordConsumerAddress, _, multiWordConsumerContract, err := multiwordconsumer_wrapper.DeployMultiWordConsumer(user, b, linkTokenAddress, operatorAddress, empty)
	require.NoError(t, err)
	b.Commit()

	singleConsumerAddress, _, singleConsumerContract, err := consumer_wrapper.DeployConsumer(user, b, linkTokenAddress, operatorAddress, empty)
	require.NoError(t, err)
	b.Commit()

	// The consumer contract needs to have link in it to be able to pay
	// for the data request.
	_, err = linkContract.Transfer(user, multiWordConsumerAddress, big.NewInt(1000))
	require.NoError(t, err)
	_, err = linkContract.Transfer(user, singleConsumerAddress, big.NewInt(1000))
	require.NoError(t, err)

	return OperatorContracts{
		user:                      user,
		multiWordConsumerAddress:  multiWordConsumerAddress,
		singleWordConsumerAddress: singleConsumerAddress,
		linkToken:                 linkContract,
		multiWord:                 multiWordConsumerContract,
		singleWord:                singleConsumerContract,
		operator:                  operatorContract,
		operatorAddress:           operatorAddress,
		sim:                       b,
	}
}

// Tests both single and multiple word responses -
// i.e. both fulfillOracleRequest2 and fulfillOracleRequest.
func TestIntegration_DirectRequest(t *testing.T) {
	// Simulate a consumer contract calling to obtain ETH quotes in 3 different currencies
	// in a single callback.
	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.SetTriggerFallbackDBPollInterval(100 * time.Millisecond)
	operatorContracts := setupOperatorContracts(t)
	b := operatorContracts.sim
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	defer cleanup()

	sendingKeys, err := app.KeyStore.Eth().SendingKeys()
	require.NoError(t, err)
	authorizedSenders := []common.Address{sendingKeys[0].Address.Address()}
	tx, err := operatorContracts.operator.SetAuthorizedSenders(operatorContracts.user, authorizedSenders)
	require.NoError(t, err)
	b.Commit()
	cltest.RequireTxSuccessful(t, b, tx.Hash())

	// Fund node account with ETH.
	n, err := b.NonceAt(context.Background(), operatorContracts.user.From, nil)
	require.NoError(t, err)
	tx = types.NewTransaction(n, sendingKeys[0].Address.Address(), big.NewInt(1000000000000000000), 21000, big.NewInt(1000000000), nil)
	signedTx, err := operatorContracts.user.Signer(operatorContracts.user.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	b.Commit()

	err = app.Start()
	require.NoError(t, err)

	mockServerUSD, cleanup := cltest.NewHTTPMockServer(t, 200, "GET", `{"USD": 614.64}`)
	defer cleanup()
	mockServerEUR, cleanup := cltest.NewHTTPMockServer(t, 200, "GET", `{"EUR": 507.07}`)
	defer cleanup()
	mockServerJPY, cleanup := cltest.NewHTTPMockServer(t, 200, "GET", `{"JPY": 63818.86}`)
	defer cleanup()

	spec := string(cltest.MustReadFile(t, "../testdata/tomlspecs/multiword-response-spec.toml"))
	spec = strings.ReplaceAll(spec, "0x613a38AC1659769640aaE063C651F48E0250454C", operatorContracts.operatorAddress.Hex())
	j := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: spec})))
	cltest.AwaitJobActive(t, app.JobSpawner(), j.ID, 5*time.Second)

	var jobID [32]byte
	copy(jobID[:], j.ExternalJobID.Bytes())
	tx, err = operatorContracts.multiWord.SetSpecID(operatorContracts.user, jobID)
	require.NoError(t, err)
	b.Commit()
	cltest.RequireTxSuccessful(t, b, tx.Hash())

	operatorContracts.user.GasLimit = 1000000
	tx, err = operatorContracts.multiWord.RequestMultipleParametersWithCustomURLs(operatorContracts.user,
		mockServerUSD.URL, "USD",
		mockServerEUR.URL, "EUR",
		mockServerJPY.URL, "JPY",
		big.NewInt(1000),
	)
	require.NoError(t, err)
	b.Commit()
	cltest.RequireTxSuccessful(t, b, tx.Hash())

	empty := big.NewInt(0)
	assertPricesUint256(t, empty, empty, empty, operatorContracts.multiWord)

	stopBlocks := finiteTicker(100*time.Millisecond, func() {
		triggerAllKeys(t, app)
		b.Commit()
	})
	defer stopBlocks()

	pipelineRuns := cltest.WaitForPipelineComplete(t, 0, j.ID, 1, 14, app.JobORM(), 10*time.Second, 100*time.Millisecond)
	pipelineRun := pipelineRuns[0]
	cltest.AssertPipelineTaskRunsSuccessful(t, pipelineRun.PipelineTaskRuns)
	assertPricesUint256(t, big.NewInt(61464), big.NewInt(50707), big.NewInt(6381886), operatorContracts.multiWord)

	// Do a single word request
	singleWordSpec := string(cltest.MustReadFile(t, "../testdata/tomlspecs/direct-request-spec-cbor.toml"))
	singleWordSpec = strings.ReplaceAll(singleWordSpec, "0x613a38AC1659769640aaE063C651F48E0250454C", operatorContracts.operatorAddress.Hex())
	jobSingleWord := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: singleWordSpec})))
	cltest.AwaitJobActive(t, app.JobSpawner(), jobSingleWord.ID, 5*time.Second)

	var jobIDSingleWord [32]byte
	copy(jobIDSingleWord[:], jobSingleWord.ExternalJobID.Bytes())
	tx, err = operatorContracts.singleWord.SetSpecID(operatorContracts.user, jobIDSingleWord)
	require.NoError(t, err)
	b.Commit()
	cltest.RequireTxSuccessful(t, b, tx.Hash())
	mockServerUSD2, cleanup := cltest.NewHTTPMockServer(t, 200, "GET", `{"USD": 614.64}`)
	defer cleanup()
	tx, err = operatorContracts.singleWord.RequestMultipleParametersWithCustomURLs(operatorContracts.user,
		mockServerUSD2.URL, "USD",
		big.NewInt(1000),
	)
	require.NoError(t, err)
	b.Commit()
	cltest.RequireTxSuccessful(t, b, tx.Hash())

	pipelineRuns = cltest.WaitForPipelineComplete(t, 0, jobSingleWord.ID, 1, 8, app.JobORM(), 5*time.Second, 100*time.Millisecond)
	pipelineRun = pipelineRuns[0]
	cltest.AssertPipelineTaskRunsSuccessful(t, pipelineRun.PipelineTaskRuns)
	v, err := operatorContracts.singleWord.CurrentPriceInt(nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(61464), v)
}

func setupOCRContracts(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *offchainaggregator.OffchainAggregator, *flags_wrapper.Flags, common.Address) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	owner := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10) // 1 eth
	genesisData := core.GenesisAlloc{
		owner.From: {Balance: sb},
	}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(owner, b)
	require.NoError(t, err)
	accessAddress, _, _, err :=
		testoffchainaggregator.DeploySimpleWriteAccessController(owner, b)
	require.NoError(t, err, "failed to deploy test access controller contract")
	b.Commit()

	min, max := new(big.Int), new(big.Int)
	min.Exp(big.NewInt(-2), big.NewInt(191), nil)
	max.Exp(big.NewInt(2), big.NewInt(191), nil)
	max.Sub(max, big.NewInt(1))
	ocrContractAddress, _, ocrContract, err := offchainaggregator.DeployOffchainAggregator(owner, b,
		1000,             // _maximumGasPrice uint32,
		200,              //_reasonableGasPrice uint32,
		3.6e7,            // 3.6e7 microLINK, or 36 LINK
		1e8,              // _linkGweiPerObservation uint32,
		4e8,              // _linkGweiPerTransmission uint32,
		linkTokenAddress, //_link common.Address,
		min,              // -2**191
		max,              // 2**191 - 1
		accessAddress,
		accessAddress,
		0,
		"TEST")
	require.NoError(t, err)
	_, err = linkContract.Transfer(owner, ocrContractAddress, big.NewInt(1000))
	require.NoError(t, err)

	flagsContractAddress, _, flagsContract, err := flags_wrapper.DeployFlags(owner, b, owner.From)
	require.NoError(t, err, "failed to deploy flags contract to simulated ethereum blockchain")

	b.Commit()
	return owner, b, ocrContractAddress, ocrContract, flagsContract, flagsContractAddress
}

func setupNode(t *testing.T, owner *bind.TransactOpts, port int, dbName string, b *backends.SimulatedBackend) (*cltest.TestApplication, string, common.Address, ocrkey.KeyV2, *configtest.TestGeneralConfig, func()) {
	config, _, ormCleanup := heavyweight.FullTestORM(t, fmt.Sprintf("%s%d", dbName, port), true, true)

	app, appCleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	_, err := app.GetKeyStore().P2P().Create()
	require.NoError(t, err)
	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()

	config.Overrides.P2PPeerID = &peerID
	config.Overrides.P2PListenPort = null.IntFrom(int64(port))
	config.Overrides.Dev = null.BoolFrom(true) // Disables ocr spec validation so we can have fast polling for the test.

	sendingKeys, err := app.KeyStore.Eth().SendingKeys()
	require.NoError(t, err)
	transmitter := sendingKeys[0].Address.Address()

	// Fund the transmitter address with some ETH
	n, err := b.NonceAt(context.Background(), owner.From, nil)
	require.NoError(t, err)

	tx := types.NewTransaction(n, transmitter, big.NewInt(1000000000000000000), 21000, big.NewInt(1000000000), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	b.Commit()

	key, err := app.GetKeyStore().OCR().Create()
	require.NoError(t, err)
	return app, peerID.Raw(), transmitter, key, config, func() {
		ormCleanup()
		appCleanup()
	}
}

func TestIntegration_OCR(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)
	owner, b, ocrContractAddress, ocrContract, flagsContract, flagsContractAddress := setupOCRContracts(t)

	// Note it's plausible these ports could be occupied on a CI machine.
	// May need a port randomize + retry approach if we observe collisions.
	appBootstrap, bootstrapPeerID, _, _, _, cleanup := setupNode(t, owner, 19999, "bootstrap", b)
	defer cleanup()

	var (
		oracles      []confighelper.OracleIdentityExtra
		transmitters []common.Address
		keys         []ocrkey.KeyV2
		apps         []*cltest.TestApplication
	)
	for i := 0; i < 4; i++ {
		app, peerID, transmitter, key, cfg, cleanup := setupNode(t, owner, 20000+i, fmt.Sprintf("oracle%d", i), b)
		defer cleanup()
		// We want to quickly poll for the bootstrap node to come up, but if we poll too quickly
		// we'll flood it with messages and slow things down. 5s is about how long it takes the
		// bootstrap node to come up.
		cfg.Overrides.SetOCRBootstrapCheckInterval(5 * time.Second)
		// GracePeriod < ObservationTimeout
		cfg.Overrides.SetOCRObservationGracePeriod(100 * time.Millisecond)
		cfg.Overrides.GlobalFlagsContractAddress = null.StringFrom(flagsContractAddress.String())

		keys = append(keys, key)
		apps = append(apps, app)
		transmitters = append(transmitters, transmitter)

		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnChainSigningAddress: ocrtypes.OnChainSigningAddress(key.OnChainSigning.Address()),
				TransmitAddress:       transmitter,
				OffchainPublicKey:     ocrtypes.OffchainPublicKey(key.PublicKeyOffChain()),
				PeerID:                peerID,
			},
			SharedSecretEncryptionPublicKey: ocrtypes.SharedSecretEncryptionPublicKey(key.PublicKeyConfig()),
		})
	}

	stopBlocks := finiteTicker(time.Second, func() {
		b.Commit()
	})
	defer stopBlocks()

	_, err := ocrContract.SetPayees(owner,
		transmitters,
		transmitters,
	)
	require.NoError(t, err)
	signers, transmitters, threshold, encodedConfigVersion, encodedConfig, err := confighelper.ContractSetConfigArgsForIntegrationTest(
		oracles,
		1,
		1000000000/100, // threshold PPB
	)
	require.NoError(t, err)
	_, err = ocrContract.SetConfig(owner,
		signers,
		transmitters,
		threshold,
		encodedConfigVersion,
		encodedConfig,
	)
	require.NoError(t, err)
	b.Commit()

	err = appBootstrap.Start()
	require.NoError(t, err)

	ocrJob, err := offchainreporting.ValidatedOracleSpecToml(appBootstrap.GetChainSet(), fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "boot"
contractAddress    = "%s"
isBootstrapPeer    = true
`, ocrContractAddress))
	require.NoError(t, err)
	_, err = appBootstrap.AddJobV2(context.Background(), ocrJob, null.NewString("boot", true))
	require.NoError(t, err)

	// Raising flags to initiate hibernation
	_, err = flagsContract.RaiseFlag(owner, ocrContractAddress)
	require.NoError(t, err, "failed to raise flag for ocrContractAddress")
	_, err = flagsContract.RaiseFlag(owner, utils.ZeroAddress)
	require.NoError(t, err, "failed to raise flag for ZeroAddress")

	b.Commit()

	var jids []int32
	var servers, slowServers = make([]*httptest.Server, 4), make([]*httptest.Server, 4)
	// We expect metadata of:
	//  latestAnswer:nil // First call
	//  latestAnswer:0
	//  latestAnswer:10
	//  latestAnswer:20
	//  latestAnswer:30
	var metaLock sync.Mutex
	expectedMeta := map[string]struct{}{
		"0": {}, "10": {}, "20": {}, "30": {},
	}
	for i := 0; i < 4; i++ {
		err = apps[i].Start()
		require.NoError(t, err)

		// Since this API speed is > ObservationTimeout we should ignore it and still produce values.
		slowServers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(5 * time.Second)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"data":10}`))
		}))
		defer slowServers[i].Close()
		servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			b, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			var m models.BridgeMetaDataJSON
			require.NoError(t, json.Unmarshal(b, &m))
			if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
				metaLock.Lock()
				delete(expectedMeta, m.Meta.LatestAnswer.String())
				metaLock.Unlock()
			}
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"data":10}`))
		}))
		defer servers[i].Close()
		u, _ := url.Parse(servers[i].URL)
		apps[i].Store.CreateBridgeType(&models.BridgeType{
			Name: models.TaskType(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		})

		// Note we need: observationTimeout + observationGracePeriod + DeltaGrace (500ms) < DeltaRound (1s)
		// So 200ms + 200ms + 500ms < 1s
		ocrJob, err := offchainreporting.ValidatedOracleSpecToml(apps[i].GetChainSet(), fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "web oracle spec"
contractAddress    = "%s"
isBootstrapPeer    = false
p2pBootstrapPeers  = [
    "/ip4/127.0.0.1/tcp/19999/p2p/%s"
]
keyBundleID        = "%s"
transmitterAddress = "%s"
observationTimeout = "100ms"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
observationSource = """
    // data source 1
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="data"];
    ds1_multiply [type=multiply times=%d];

    // data source 2
    ds2          [type=http method=GET url="%s"];
    ds2_parse    [type=jsonparse path="data"];
    ds2_multiply [type=multiply times=%d];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
"""
`, ocrContractAddress, bootstrapPeerID, keys[i].ID(), transmitters[i], fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i))
		require.NoError(t, err)
		jb, err := apps[i].AddJobV2(context.Background(), ocrJob, null.NewString("testocr", true))
		require.NoError(t, err)
		jids = append(jids, jb.ID)
	}

	// Assert that all the OCR jobs get a run with valid values eventually.
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Want at least 2 runs so we see all the metadata.
			pr := cltest.WaitForPipelineComplete(t, ic, jids[ic], 2, 0, apps[ic].JobORM(), 1*time.Minute, 1*time.Second)
			jb, err := pr[0].Outputs.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, []byte(fmt.Sprintf("[\"%d\"]", 10*ic)), jb)
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	// 4 oracles reporting 0, 10, 20, 30. Answer should be 20 (results[4/2]).
	g.Eventually(func() string {
		answer, err := ocrContract.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.String()
	}, 10*time.Second, 200*time.Millisecond).Should(gomega.Equal("20"))

	for _, app := range apps {
		jobs, _, err := app.JobORM().JobsV2(0, 1000)
		require.NoError(t, err)
		// No spec errors
		for _, j := range jobs {
			ignore := 0
			for i := range j.JobSpecErrors {
				// Non-fatal timing related error, ignore for testing.
				if strings.Contains(j.JobSpecErrors[i].Description, "leader's phase conflicts tGrace timeout") {
					ignore++
				}
			}
			require.Len(t, j.JobSpecErrors, ignore)
		}
	}
	assert.Len(t, expectedMeta, 0, "expected metadata %v", expectedMeta)

	wg.Add(5)
	go func() {
		defer wg.Done()
		require.NoError(t, appBootstrap.Stop())
	}()
	for i := range apps {
		app := apps[i]
		go func() {
			defer wg.Done()
			require.NoError(t, app.Stop())
		}()
	}
	wg.Wait()
}

func TestIntegration_BlockHistoryEstimator(t *testing.T) {
	t.Parallel()

	var initialDefaultGasPrice int64 = 5000000000

	c := cltest.NewTestGeneralConfig(t)
	c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)

	ethClient, sub, assertMocksCalled := cltest.NewEthMocksWithDefaultChain(t)
	defer assertMocksCalled()
	chchNewHeads := make(chan chan<- *models.Head, 1)

	db := pgtest.NewGormDB(t)
	kst := cltest.NewKeyStore(t, db)
	require.NoError(t, kst.Unlock(cltest.Password))

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, KeyStore: kst.Eth(), Client: ethClient, GeneralConfig: c, ChainCfg: evmtypes.ChainCfg{
		EvmGasPriceDefault:                    utils.NewBigI(initialDefaultGasPrice),
		GasEstimatorMode:                      null.StringFrom("BlockHistory"),
		BlockHistoryEstimatorBlockDelay:       null.IntFrom(0),
		BlockHistoryEstimatorBlockHistorySize: null.IntFrom(2),
		EvmFinalityDepth:                      null.IntFrom(3),
	}})

	b41 := gas.Block{
		Number:       41,
		Hash:         utils.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(41000000000, 41500000000),
	}
	b42 := gas.Block{
		Number:       42,
		Hash:         utils.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(44000000000, 45000000000),
	}
	b43 := gas.Block{
		Number:       43,
		Hash:         utils.NewHash(),
		Transactions: cltest.TransactionsFromGasPrices(48000000000, 49000000000, 31000000000),
	}

	h40 := models.Head{Hash: utils.NewHash(), Number: 40}
	h41 := models.Head{Hash: b41.Hash, ParentHash: h40.Hash, Number: 41}
	h42 := models.Head{Hash: b42.Hash, ParentHash: h41.Hash, Number: 42}

	sub.On("Err").Return(nil)
	sub.On("Unsubscribe").Return(nil).Maybe()

	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) { chchNewHeads <- args.Get(1).(chan<- *models.Head) }).
		Return(sub, nil)
	// Nonce syncer
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)

	// BlockHistoryEstimator boot calls
	ethClient.On("HeadByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(&h42, nil)
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x29" &&
			b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2a"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b41
		elems[1].Result = &b42
	})

	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Return(c.DefaultChainID(), nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)

	require.NoError(t, cc.Start())
	var newHeads chan<- *models.Head
	select {
	case newHeads = <-chchNewHeads:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for app to subscribe")
	}

	chain := evmtest.MustGetDefaultChain(t, cc)
	estimator := chain.TxManager().GetGasEstimator()
	gasPrice, gasLimit, err := estimator.EstimateGas(nil, 500000)
	require.NoError(t, err)
	assert.Equal(t, uint64(500000), gasLimit)
	assert.Equal(t, "41500000000", gasPrice.String())
	assert.Equal(t, initialDefaultGasPrice, chain.Config().EvmGasPriceDefault().Int64()) // unchanged

	// BlockHistoryEstimator new blocks
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2a" &&
			b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x2b"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b43
		elems[1].Result = &b42
	})

	// HeadTracker backfill
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(42)).Return(&h42, nil)
	ethClient.On("HeadByNumber", mock.Anything, big.NewInt(41)).Return(&h41, nil)

	// Simulate one new head and check the gas price got updated
	newHeads <- cltest.Head(43)

	gomega.NewGomegaWithT(t).Eventually(func() string {
		gasPrice, _, err := estimator.EstimateGas(nil, 500000)
		require.NoError(t, err)
		return gasPrice.String()
	}, cltest.DBWaitTimeout, cltest.DBPollingInterval).Should(gomega.Equal("45000000000"))
}

func triggerAllKeys(t *testing.T, app *cltest.TestApplication) {
	keys, err := app.KeyStore.Eth().SendingKeys()
	require.NoError(t, err)
	// FIXME: This is a hack. Remove after https://app.clubhouse.io/chainlinklabs/story/15103/use-in-memory-event-broadcaster-instead-of-postgres-event-broadcaster-in-transactional-tests-so-it-actually-works
	for _, chain := range app.GetChainSet().Chains() {
		for _, k := range keys {
			chain.TxManager().Trigger(k.Address.Address())
		}
	}
}

func assertPricesUint256(t *testing.T, usd, eur, jpy *big.Int, consumer *multiwordconsumer_wrapper.MultiWordConsumer) {
	haveUsd, err := consumer.UsdInt(nil)
	require.NoError(t, err)
	assert.True(t, usd.Cmp(haveUsd) == 0)
	haveEur, err := consumer.EurInt(nil)
	require.NoError(t, err)
	assert.True(t, eur.Cmp(haveEur) == 0)
	haveJpy, err := consumer.JpyInt(nil)
	require.NoError(t, err)
	assert.True(t, jpy.Cmp(haveJpy) == 0)
}

func finiteTicker(period time.Duration, onTick func()) func() {
	tick := time.NewTicker(period)
	chStop := make(chan struct{})
	go func() {
		for {
			select {
			case <-tick.C:
				onTick()
			case <-chStop:
				return
			}
		}
	}()

	// NOTE: tick.Stop does not close the ticker channel,
	// so we still need another way of returning (chStop).
	return func() {
		tick.Stop()
		close(chStop)
	}
}
