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
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/multiwordconsumer_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
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

	cfg, cleanup := cltest.NewConfig(t)
	defer cleanup()
	cfg.Set("FEATURE_EXTERNAL_INITIATORS", true)
	cfg.Set("TRIGGER_FALLBACK_DB_POLL_INTERVAL", "10ms")

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
		jobORM := job.NewORM(app.Store.ORM.DB, app.Store.Config, pipelineORM, &postgres.NullEventBroadcaster{}, &postgres.NullAdvisoryLocker{})

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
	buf := bytes.NewBufferString(`{"ethGasPriceDefault":15000000}`)

	resp, cleanup := cltest.UnauthenticatedPatch(t, url, buf, headers)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
}

func assertPricesBytes32(t *testing.T, usd, eur, jpy []byte, consumer *multiwordconsumer_wrapper.MultiWordConsumer) {
	var tmp [32]byte
	copy(tmp[:], usd)
	haveUsd, err := consumer.Usd(nil)
	require.NoError(t, err)
	assert.Equal(t, tmp[:], haveUsd[:])
	copy(tmp[:], eur)
	haveEur, err := consumer.Eur(nil)
	require.NoError(t, err)
	assert.Equal(t, tmp[:], haveEur[:])
	copy(tmp[:], jpy)
	haveJpy, err := consumer.Jpy(nil)
	require.NoError(t, err)
	assert.Equal(t, tmp[:], haveJpy[:])
}

func setupMultiWordContracts(t *testing.T) (*bind.TransactOpts, common.Address, common.Address, *link_token_interface.LinkToken, *multiwordconsumer_wrapper.MultiWordConsumer, *operator_wrapper.Operator, *backends.SimulatedBackend) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	user := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10)
	genesisData := core.GenesisAlloc{
		user.From: {Balance: sb}, // 1 eth
	}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(user, b)
	require.NoError(t, err)
	b.Commit()

	operatorAddress, _, operatorContract, err := operator_wrapper.DeployOperator(user, b, linkTokenAddress, user.From)
	require.NoError(t, err)
	b.Commit()

	var empty [32]byte
	consumerAddress, _, consumerContract, err := multiwordconsumer_wrapper.DeployMultiWordConsumer(user, b, linkTokenAddress, operatorAddress, empty)
	require.NoError(t, err)
	b.Commit()

	// The consumer contract needs to have link in it to be able to pay
	// for the data request.
	_, err = linkContract.Transfer(user, consumerAddress, big.NewInt(1000))
	require.NoError(t, err)
	return user, consumerAddress, operatorAddress, linkContract, consumerContract, operatorContract, b
}

func setupOCRContracts(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *offchainaggregator.OffchainAggregator) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	owner := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10) // 1 eth
	genesisData := core.GenesisAlloc{
		owner.From: {Balance: sb},
	}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
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
	b.Commit()
	return owner, b, ocrContractAddress, ocrContract
}

func setupNode(t *testing.T, owner *bind.TransactOpts, port int, dbName string, b *backends.SimulatedBackend) (*cltest.TestApplication, string, common.Address, ocrkey.EncryptedKeyBundle, func()) {
	config, _, ormCleanup := heavyweight.FullTestORM(t, fmt.Sprintf("%s%d", dbName, port), true)
	config.Dialect = dialects.PostgresWithoutLock
	app, appCleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	_, _, err := app.GetKeyStore().OCR().GenerateEncryptedP2PKey()
	require.NoError(t, err)
	p2pIDs := app.GetKeyStore().OCR().DecryptedP2PKeys()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].MustGetPeerID().Raw()

	app.Config.Set("P2P_PEER_ID", peerID)
	app.Config.Set("P2P_LISTEN_PORT", port)
	app.Config.Set("ETH_HEAD_TRACKER_MAX_BUFFER_SIZE", 100)
	app.Config.Set("MIN_OUTGOING_CONFIRMATIONS", 1)
	app.Config.Set("CHAINLINK_DEV", true) // Disables ocr spec validation so we can have fast polling for the test.

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

	_, kb, err := app.GetKeyStore().OCR().GenerateEncryptedOCRKeyBundle()
	require.NoError(t, err)
	return app, peerID, transmitter, kb, func() {
		ormCleanup()
		appCleanup()
	}
}

func TestIntegration_OCR(t *testing.T) {
	t.Parallel()

	owner, b, ocrContractAddress, ocrContract := setupOCRContracts(t)

	// Note it's plausible these ports could be occupied on a CI machine.
	// May need a port randomize + retry approach if we observe collisions.
	appBootstrap, bootstrapPeerID, _, _, cleanup := setupNode(t, owner, 19999, "bootstrap", b)
	defer cleanup()

	var (
		oracles      []confighelper.OracleIdentityExtra
		transmitters []common.Address
		kbs          []ocrkey.EncryptedKeyBundle
		apps         []*cltest.TestApplication
	)
	for i := 0; i < 4; i++ {
		app, peerID, transmitter, kb, cleanup := setupNode(t, owner, 20000+i, fmt.Sprintf("oracle%d", i), b)
		defer cleanup()
		// We want to quickly poll for the bootstrap node to come up, but if we poll too quickly
		// we'll flood it with messages and slow things down. 5s is about how long it takes the
		// bootstrap node to come up.
		app.Config.Set("OCR_BOOTSTRAP_CHECK_INTERVAL", "5s")
		// GracePeriod < ObservationTimeout
		app.Config.Set("OCR_OBSERVATION_GRACE_PERIOD", "100ms")

		kbs = append(kbs, kb)
		apps = append(apps, app)
		transmitters = append(transmitters, transmitter)

		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnChainSigningAddress: ocrtypes.OnChainSigningAddress(kb.OnChainSigningAddress),
				TransmitAddress:       transmitter,
				OffchainPublicKey:     ocrtypes.OffchainPublicKey(kb.OffChainPublicKey),
				PeerID:                peerID,
			},
			SharedSecretEncryptionPublicKey: ocrtypes.SharedSecretEncryptionPublicKey(kb.ConfigPublicKey),
		})
	}

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	go func() {
		for range tick.C {
			b.Commit()
		}
	}()

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
	defer appBootstrap.Stop()

	ocrJob, err := offchainreporting.ValidatedOracleSpecToml(appBootstrap.Config.Config, fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "boot"
contractAddress    = "%s"
isBootstrapPeer    = true
`, ocrContractAddress))
	require.NoError(t, err)
	_, err = appBootstrap.AddJobV2(context.Background(), ocrJob, null.NewString("boot", true))
	require.NoError(t, err)

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
		defer apps[i].Stop()

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
		ocrJob, err := offchainreporting.ValidatedOracleSpecToml(apps[i].Config.Config, fmt.Sprintf(`
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
`, ocrContractAddress, bootstrapPeerID, kbs[i].ID, transmitters[i], fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i))
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
	gomega.NewGomegaWithT(t).Eventually(func() string {
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
}

func TestIntegration_DirectRequest(t *testing.T) {
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()

	httpAwaiter := cltest.NewAwaiter()
	httpServer, assertCalled := cltest.NewHTTPMockServer(
		t,
		http.StatusOK,
		"GET",
		`{"USD": "31982"}`,
		func(header http.Header, _ string) {
			httpAwaiter.ItHappened()
		},
	)
	defer assertCalled()

	ethClient, sub, assertMockCalls := cltest.NewEthMocks(t)
	defer assertMockCalls()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	defer cleanup()

	blocks := cltest.NewBlocks(t, 12)

	sub.On("Err").Return(nil).Maybe()
	sub.On("Unsubscribe").Return(nil).Maybe()

	ethClient.On("HeadByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(blocks.Head(10), nil)

	var headCh chan<- *models.Head
	ethClient.On("SubscribeNewHead", mock.Anything, mock.Anything).Maybe().
		Run(func(args mock.Arguments) {
			headCh = args.Get(1).(chan<- *models.Head)
		}).
		Return(sub, nil)

	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ChainID", mock.Anything).Maybe().Return(app.Store.Config.ChainID(), nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)
	ethClient.On("HeadByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(blocks.Head(0), nil)
	logsCh := cltest.MockSubscribeToLogsCh(ethClient, sub)

	require.NoError(t, app.Start())

	store := app.Store
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	defer eventBroadcaster.Close()

	pipelineORM := pipeline.NewORM(store.DB)
	jobORM := job.NewORM(store.ORM.DB, store.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})

	directRequestSpec := string(cltest.MustReadFile(t, "../testdata/tomlspecs/direct-request-spec.toml"))
	directRequestSpec = strings.Replace(directRequestSpec, "http://example.com", httpServer.URL, 1)
	request := web.CreateJobRequest{TOML: directRequestSpec}
	output, err := json.Marshal(request)
	require.NoError(t, err)
	job := cltest.CreateJobViaWeb(t, app, output)

	eventBroadcaster.Notify(postgres.ChannelJobCreated, "")

	runLog := cltest.NewRunLog(t, job.ExternalIDEncodeStringToTopic(), job.DirectRequestSpec.ContractAddress.Address(), cltest.NewAddress(), 1, `{}`)
	runLog.BlockHash = blocks.Head(1).Hash
	var logs chan<- types.Log
	cltest.CallbackOrTimeout(t, "obtain log channel", func() {
		logs = <-logsCh
	}, 5*time.Second)
	cltest.CallbackOrTimeout(t, "send run log", func() {
		logs <- runLog
	}, 30*time.Second)

	eventBroadcaster.Notify(postgres.ChannelRunStarted, "")
	for i := 0; i < 12; i++ {
		headCh <- blocks.Head(uint64(i))
	}

	httpAwaiter.AwaitOrFail(t)

	runs := cltest.WaitForPipelineComplete(t, 0, job.ID, 1, 3, jobORM, 5*time.Second, 300*time.Millisecond)
	require.Len(t, runs, 1)
	run := runs[0]
	require.Len(t, run.PipelineTaskRuns, 3)
	require.Empty(t, run.PipelineTaskRuns[0].Error)
	require.Empty(t, run.PipelineTaskRuns[1].Error)
	require.Empty(t, run.PipelineTaskRuns[2].Error)
}

func TestIntegration_BlockHistoryEstimator(t *testing.T) {
	t.Parallel()

	var initialDefaultGasPrice int64 = 5000000000

	c, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	c.Set("ETH_GAS_PRICE_DEFAULT", initialDefaultGasPrice)
	c.Set("GAS_ESTIMATOR_MODE", "BlockHistory")
	c.Set("GAS_UPDATER_BLOCK_DELAY", 0)
	c.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", 2)
	// Limit the headtracker backfill depth just so we aren't here all week
	c.Set("ETH_FINALITY_DEPTH", 3)

	ethClient, sub, assertMocksCalled := cltest.NewEthMocks(t)
	defer assertMocksCalled()
	chchNewHeads := make(chan chan<- *models.Head, 1)

	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, c,
		ethClient,
	)
	defer cleanup()

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
	ethClient.On("ChainID", mock.Anything).Return(c.ChainID(), nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)

	require.NoError(t, app.Start())
	var newHeads chan<- *models.Head
	select {
	case newHeads = <-chchNewHeads:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for app to subscribe")
	}

	estimator := app.TxManager.GetGasEstimator()
	gasPrice, gasLimit, err := estimator.EstimateGas(nil, 500000)
	require.NoError(t, err)
	assert.Equal(t, uint64(500000), gasLimit)
	assert.Equal(t, "41500000000", gasPrice.String())
	assert.Equal(t, initialDefaultGasPrice, app.Config.EthGasPriceDefault().Int64()) // unchanged

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
	for _, k := range keys {
		app.TxManager.Trigger(k.Address.Address())
	}
}
