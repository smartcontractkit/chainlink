package features_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/freeport"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers/testoffchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/consumer_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/multiwordconsumer_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/operator"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/forwarders"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	webauth "github.com/smartcontractkit/chainlink/v2/core/web/auth"
)

var oneETH = assets.Eth(*big.NewInt(1000000000000000000))

func TestIntegration_ExternalInitiatorV2(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	ethClient := cltest.NewEthMocksWithStartupAssertions(t)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(10 * time.Millisecond)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient, cltest.UseRealExternalInitiatorManager)
	require.NoError(t, app.Start(testutils.Context(t)))

	var (
		eiName    = "substrate-ei"
		eiSpec    = map[string]interface{}{"foo": "bar"}
		eiRequest = map[string]interface{}{"result": 42}

		jobUUID = uuid.MustParse("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46")

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
			_, err = io.WriteString(w, `{}`)
			require.NoError(t, err)
		}))
		u, _ := url.Parse(bridgeServer.URL)
		err := app.BridgeORM().CreateBridgeType(ctx, &bridges.BridgeType{
			Name: bridges.BridgeName("substrate-adapter1"),
			URL:  models.WebURL(*u),
		})
		require.NoError(t, err)
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

		_, err := webhook.ValidatedWebhookSpec(ctx, tomlSpec, app.GetExternalInitiatorManager())
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

		url := app.Server.URL + "/v2/jobs/" + jobUUID.String() + "/runs"
		bodyBuf := bytes.NewBufferString(body)
		resp, cleanup := cltest.UnauthenticatedPost(t, url, bodyBuf, headers)
		defer cleanup()
		cltest.AssertServerResponse(t, resp, 401)

		cltest.AssertCountStays(t, app.GetDB(), "pipeline_runs", 0)
	})

	t.Run("calling webhook_spec with matching external_initiator_id works", func(t *testing.T) {
		// Simulate request from EI -> Core node
		cltest.AwaitJobActive(t, app.JobSpawner(), jobID, 3*time.Second)

		_ = cltest.CreateJobRunViaExternalInitiatorV2(t, app, jobUUID, *eia, cltest.MustJSONMarshal(t, eiRequest))

		pipelineORM := pipeline.NewORM(app.GetDB(), logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
		bridgeORM := bridges.NewORM(app.GetDB())
		jobORM := job.NewORM(app.GetDB(), pipelineORM, bridgeORM, app.KeyStore, logger.TestLogger(t))

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
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	// set up user
	mockUser := cltest.MustRandomUser(t)
	key, secret := uuid.New().String(), uuid.New().String()
	apiToken := auth.Token{AccessKey: key, Secret: secret}
	orm := app.AuthenticationProvider()
	require.NoError(t, orm.CreateUser(ctx, &mockUser))
	require.NoError(t, orm.SetAuthToken(ctx, &mockUser, &apiToken))

	url := app.Server.URL + "/users"
	headers := make(map[string]string)
	headers[webauth.APIKey] = key
	headers[webauth.APISecret] = secret

	resp, cleanup := cltest.UnauthenticatedGet(t, url, headers)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
}

type OperatorContracts struct {
	user                      *bind.TransactOpts
	multiWordConsumerAddress  common.Address
	singleWordConsumerAddress common.Address
	operatorAddress           common.Address
	linkTokenAddress          common.Address
	linkToken                 *link_token_interface.LinkToken
	multiWord                 *multiwordconsumer_wrapper.MultiWordConsumer
	singleWord                *consumer_wrapper.Consumer
	operator                  *operator.Operator
	sim                       types.Backend
}

func setupOperatorContracts(t *testing.T) OperatorContracts {
	user := evmtestutils.MustNewSimTransactor(t)
	genesisData := gethtypes.GenesisAlloc{
		user.From: {Balance: assets.Ether(1000).ToInt()},
	}
	b := cltest.NewSimulatedBackend(t, genesisData, 2*ethconfig.Defaults.Miner.GasCeil)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(user, b.Client())
	require.NoError(t, err)
	b.Commit()

	operatorAddress, _, operatorContract, err := operator.DeployOperator(user, b.Client(), linkTokenAddress, user.From)
	require.NoError(t, err)
	b.Commit()

	var empty [32]byte
	multiWordConsumerAddress, _, multiWordConsumerContract, err := multiwordconsumer_wrapper.DeployMultiWordConsumer(user, b.Client(), linkTokenAddress, operatorAddress, empty)
	require.NoError(t, err)
	b.Commit()

	singleConsumerAddress, _, singleConsumerContract, err := consumer_wrapper.DeployConsumer(user, b.Client(), linkTokenAddress, operatorAddress, empty)
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
		linkTokenAddress:          linkTokenAddress,
		multiWord:                 multiWordConsumerContract,
		singleWord:                singleConsumerContract,
		operator:                  operatorContract,
		operatorAddress:           operatorAddress,
		sim:                       b,
	}
}

//go:embed singleword-spec-template.yml
var singleWordSpecTemplate string

//go:embed multiword-spec-template.yml
var multiWordSpecTemplate string

// Tests both single and multiple word responses -
// i.e. both fulfillOracleRequest2 and fulfillOracleRequest.
func TestIntegration_DirectRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		eip1559 bool
	}{
		{"legacy mode", false},
		{"eip1559 mode", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			// Simulate a consumer contract calling to obtain ETH quotes in 3 different currencies
			// in a single callback.
			config := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
				c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
			})
			operatorContracts := setupOperatorContracts(t)
			b := operatorContracts.sim
			app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b)

			sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(ctx, testutils.SimulatedChainID)
			require.NoError(t, err)
			authorizedSenders := []common.Address{sendingKeys[0].Address}
			tx, err := operatorContracts.operator.SetAuthorizedSenders(operatorContracts.user, authorizedSenders)
			require.NoError(t, err)
			b.Commit()
			cltest.RequireTxSuccessful(t, b.Client(), tx.Hash())

			// Fund node account with ETH.
			n, err := b.Client().NonceAt(testutils.Context(t), operatorContracts.user.From, nil)
			require.NoError(t, err)
			tx = cltest.NewLegacyTransaction(n, sendingKeys[0].Address, assets.Ether(100).ToInt(), 21000, big.NewInt(1000000000), nil)
			signedTx, err := operatorContracts.user.Signer(operatorContracts.user.From, tx)
			require.NoError(t, err)
			err = b.Client().SendTransaction(testutils.Context(t), signedTx)
			require.NoError(t, err)
			b.Commit()

			err = app.Start(testutils.Context(t))
			require.NoError(t, err)

			mockServerUSD := cltest.NewHTTPMockServer(t, 200, "GET", `{"USD": 614.64}`)
			mockServerEUR := cltest.NewHTTPMockServer(t, 200, "GET", `{"EUR": 507.07}`)
			mockServerJPY := cltest.NewHTTPMockServer(t, 200, "GET", `{"JPY": 63818.86}`)

			nameAndExternalJobID := uuid.New()
			addr := operatorContracts.operatorAddress.Hex()
			spec := fmt.Sprintf(multiWordSpecTemplate, nameAndExternalJobID, addr, nameAndExternalJobID, addr)
			j := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: spec})))
			cltest.AwaitJobActive(t, app.JobSpawner(), j.ID, 5*time.Second)

			var jobID [32]byte
			copy(jobID[:], j.ExternalJobID[:])
			tx, err = operatorContracts.multiWord.SetSpecID(operatorContracts.user, jobID)
			require.NoError(t, err)
			b.Commit()
			cltest.RequireTxSuccessful(t, b.Client(), tx.Hash())

			operatorContracts.user.GasLimit = 1000000
			tx, err = operatorContracts.multiWord.RequestMultipleParametersWithCustomURLs(operatorContracts.user,
				mockServerUSD.URL, "USD",
				mockServerEUR.URL, "EUR",
				mockServerJPY.URL, "JPY",
				big.NewInt(1000),
			)
			require.NoError(t, err)
			b.Commit()
			cltest.RequireTxSuccessful(t, b.Client(), tx.Hash())

			empty := big.NewInt(0)
			assertPricesUint256(t, empty, empty, empty, operatorContracts.multiWord)

			commit, stopBlocks := cltest.Mine(b, 100*time.Millisecond)
			defer stopBlocks()

			pipelineRuns := cltest.WaitForPipelineComplete(t, 0, j.ID, 1, 14, app.JobORM(), testutils.WaitTimeout(t)/2, time.Second)
			pipelineRun := pipelineRuns[0]
			assertPipelineTaskRunsSuccessful(t, pipelineRun.PipelineTaskRuns)
			assertPricesUint256(t, big.NewInt(61464), big.NewInt(50707), big.NewInt(6381886), operatorContracts.multiWord)

			nameAndExternalJobID = uuid.New()
			singleWordSpec := fmt.Sprintf(singleWordSpecTemplate, nameAndExternalJobID, addr, nameAndExternalJobID, addr)
			jobSingleWord := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: singleWordSpec})))
			cltest.AwaitJobActive(t, app.JobSpawner(), jobSingleWord.ID, 5*time.Second)

			var jobIDSingleWord [32]byte
			copy(jobIDSingleWord[:], jobSingleWord.ExternalJobID[:])
			tx, err = operatorContracts.singleWord.SetSpecID(operatorContracts.user, jobIDSingleWord)
			require.NoError(t, err)
			commit()
			cltest.RequireTxSuccessful(t, b.Client(), tx.Hash())
			mockServerUSD2 := cltest.NewHTTPMockServer(t, 200, "GET", `{"USD": 614.64}`)
			tx, err = operatorContracts.singleWord.RequestMultipleParametersWithCustomURLs(operatorContracts.user,
				mockServerUSD2.URL, "USD",
				big.NewInt(1000),
			)
			require.NoError(t, err)
			commit()
			cltest.RequireTxSuccessful(t, b.Client(), tx.Hash())

			pipelineRuns = cltest.WaitForPipelineComplete(t, 0, jobSingleWord.ID, 1, 8, app.JobORM(), testutils.WaitTimeout(t), time.Second)
			pipelineRun = pipelineRuns[0]
			assertPipelineTaskRunsSuccessful(t, pipelineRun.PipelineTaskRuns)
			v, err := operatorContracts.singleWord.CurrentPriceInt(nil)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(61464), v)
		})
	}
}

func setupAppForEthTx(t *testing.T, operatorContracts OperatorContracts) (app *cltest.TestApplication, sendingAddress common.Address, o *observer.ObservedLogs) {
	b := operatorContracts.sim
	lggr, o := logger.TestLoggerObserved(t, zapcore.DebugLevel)

	cfg := configtest.NewGeneralConfigSimulated(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Listener.FallbackPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app = cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, b, lggr)
	b.Commit()

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)

	// Fund node account with ETH.
	n, err := b.Client().NonceAt(testutils.Context(t), operatorContracts.user.From, nil)
	require.NoError(t, err)
	tx := cltest.NewLegacyTransaction(n, sendingKeys[0].Address, assets.Ether(100).ToInt(), 21000, big.NewInt(1000000000), nil)
	signedTx, err := operatorContracts.user.Signer(operatorContracts.user.From, tx)
	require.NoError(t, err)
	err = b.Client().SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	b.Commit()

	err = app.Start(testutils.Context(t))
	require.NoError(t, err)

	testutils.WaitForLogMessage(t, o, "Subscribing to new heads on chain 1337")
	testutils.WaitForLogMessage(t, o, "Subscribed to heads on chain 1337")

	return app, sendingKeys[0].Address, o
}

func TestIntegration_AsyncEthTx(t *testing.T) {
	t.Parallel()
	operatorContracts := setupOperatorContracts(t)
	b := operatorContracts.sim

	t.Run("with FailOnRevert enabled, run succeeds when transaction is successful", func(t *testing.T) {
		app, sendingAddr, o := setupAppForEthTx(t, operatorContracts)
		tomlSpec := `
type            = "webhook"
schemaVersion   = 1
observationSource   = """
	submit_tx  [type=ethtx to="%s"
            data="%s"
            minConfirmations="2"
			failOnRevert=false
			evmChainID="%s"
            from="[\\"%s\\"]"
			]
"""
`
		// This succeeds for whatever reason
		revertingData := "0xdeadbeef"
		tomlSpec = fmt.Sprintf(tomlSpec, operatorContracts.linkTokenAddress.String(), revertingData, testutils.SimulatedChainID.String(), sendingAddr)
		j := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: tomlSpec})))
		cltest.AwaitJobActive(t, app.JobSpawner(), j.ID, testutils.WaitTimeout(t))

		run := cltest.CreateJobRunViaUser(t, app, j.ExternalJobID, "")
		assert.Equal(t, []*string(nil), run.Outputs)
		assert.Equal(t, []*string(nil), run.Errors)

		testutils.WaitForLogMessage(t, o, "Sending transaction")
		gomega.NewWithT(t).Eventually(func() bool {
			b.Commit() // Process new head until tx confirmed, receipt is fetched, and task resumed
			for _, l := range o.All() {
				if strings.Contains(l.Message, "Resume run success") {
					return true
				}
			}
			return false
		}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

		pipelineRuns := cltest.WaitForPipelineComplete(t, 0, j.ID, 1, 1, app.JobORM(), testutils.WaitTimeout(t), time.Second)

		// The run should have succeeded but with the receipt detailing the reverted transaction
		pipelineRun := pipelineRuns[0]
		assertPipelineTaskRunsSuccessful(t, pipelineRun.PipelineTaskRuns)

		outputs := pipelineRun.Outputs.Val.([]interface{})
		require.Len(t, outputs, 1)
		output := outputs[0]
		receipt := output.(map[string]interface{})
		assert.Equal(t, "0x7", receipt["blockNumber"])
		assert.Equal(t, "0x538f", receipt["gasUsed"])
		assert.Equal(t, "0x0", receipt["status"]) // success
	})

	t.Run("with FailOnRevert enabled, run fails with transaction reverted error", func(t *testing.T) {
		app, sendingAddr, o := setupAppForEthTx(t, operatorContracts)
		tomlSpec := `
type            = "webhook"
schemaVersion   = 1
observationSource   = """
	submit_tx  [type=ethtx to="%s"
            data="%s"
            minConfirmations="2"
			failOnRevert=true
			evmChainID="%s"
            from="[\\"%s\\"]"
			]
"""
`
		// This data is a call to link token's `transfer` function and will revert due to insufficient LINK on the sender address
		revertingData := "0xa9059cbb000000000000000000000000526485b5abdd8ae9c6a63548e0215a83e7135e6100000000000000000000000000000000000000000000000db069932ea4fe1400"
		tomlSpec = fmt.Sprintf(tomlSpec, operatorContracts.linkTokenAddress.String(), revertingData, testutils.SimulatedChainID.String(), sendingAddr)
		j := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: tomlSpec})))
		cltest.AwaitJobActive(t, app.JobSpawner(), j.ID, testutils.WaitTimeout(t))

		run := cltest.CreateJobRunViaUser(t, app, j.ExternalJobID, "")
		assert.Equal(t, []*string(nil), run.Outputs)
		assert.Equal(t, []*string(nil), run.Errors)

		testutils.WaitForLogMessage(t, o, "Sending transaction")
		gomega.NewWithT(t).Eventually(func() bool {
			b.Commit() // Process new head until tx confirmed, receipt is fetched, and task resumed
			for _, l := range o.All() {
				if strings.Contains(l.Message, "Resume run success") {
					return true
				}
			}
			return false
		}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

		pipelineRuns := cltest.WaitForPipelineError(t, 0, j.ID, 1, 1, app.JobORM(), testutils.WaitTimeout(t), time.Second)

		// The run should have failed as a revert
		pipelineRun := pipelineRuns[0]
		assertPipelineTaskRunsErrored(t, pipelineRun.PipelineTaskRuns)
	})

	t.Run("with FailOnRevert disabled, run succeeds with output being reverted receipt", func(t *testing.T) {
		app, sendingAddr, o := setupAppForEthTx(t, operatorContracts)
		tomlSpec := `
type            = "webhook"
schemaVersion   = 1
observationSource   = """
	submit_tx  [type=ethtx to="%s"
            data="%s"
            minConfirmations="2"
			failOnRevert=false
			evmChainID="%s"
            from="[\\"%s\\"]"
			]
"""
`
		// This data is a call to link token's `transfer` function and will revert due to insufficient LINK on the sender address
		revertingData := "0xa9059cbb000000000000000000000000526485b5abdd8ae9c6a63548e0215a83e7135e6100000000000000000000000000000000000000000000000db069932ea4fe1400"
		tomlSpec = fmt.Sprintf(tomlSpec, operatorContracts.linkTokenAddress.String(), revertingData, testutils.SimulatedChainID.String(), sendingAddr)
		j := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: tomlSpec})))
		cltest.AwaitJobActive(t, app.JobSpawner(), j.ID, testutils.WaitTimeout(t))

		run := cltest.CreateJobRunViaUser(t, app, j.ExternalJobID, "")
		assert.Equal(t, []*string(nil), run.Outputs)
		assert.Equal(t, []*string(nil), run.Errors)

		testutils.WaitForLogMessage(t, o, "Sending transaction")
		gomega.NewWithT(t).Eventually(func() bool {
			b.Commit() // Process new head until tx confirmed, receipt is fetched, and task resumed
			for _, l := range o.All() {
				if strings.Contains(l.Message, "Resume run success") {
					return true
				}
			}
			return false
		}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

		pipelineRuns := cltest.WaitForPipelineComplete(t, 0, j.ID, 1, 1, app.JobORM(), testutils.WaitTimeout(t), time.Second)

		// The run should have succeeded but with the receipt detailing the reverted transaction
		pipelineRun := pipelineRuns[0]
		assertPipelineTaskRunsSuccessful(t, pipelineRun.PipelineTaskRuns)

		outputs := pipelineRun.Outputs.Val.([]interface{})
		require.Len(t, outputs, 1)
		output := outputs[0]
		receipt := output.(map[string]interface{})
		assert.Equal(t, "0x19", receipt["blockNumber"])
		assert.Equal(t, "0x7a120", receipt["gasUsed"])
		assert.Equal(t, "0x0", receipt["status"])
	})
}

func setupOCRContracts(t *testing.T) (*bind.TransactOpts, types.Backend, common.Address, *offchainaggregator.OffchainAggregator, *flags_wrapper.Flags, common.Address) {
	owner := evmtestutils.MustNewSimTransactor(t)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000000", 10) // 1000 eth
	genesisData := gethtypes.GenesisAlloc{
		owner.From: {Balance: sb},
	}
	b := cltest.NewSimulatedBackend(t, genesisData, 2*ethconfig.Defaults.Miner.GasCeil)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(owner, b.Client())
	require.NoError(t, err)
	accessAddress, _, _, err :=
		testoffchainaggregator.DeploySimpleWriteAccessController(owner, b.Client())
	require.NoError(t, err, "failed to deploy test access controller contract")
	b.Commit()

	min, max := new(big.Int), new(big.Int)
	min.Exp(big.NewInt(-2), big.NewInt(191), nil)
	max.Exp(big.NewInt(2), big.NewInt(191), nil)
	max.Sub(max, big.NewInt(1))
	ocrContractAddress, _, ocrContract, err := offchainaggregator.DeployOffchainAggregator(owner, b.Client(),
		1000,             // _maximumGasPrice uint32,
		200,              // _reasonableGasPrice uint32,
		3.6e7,            // 3.6e7 microLINK, or 36 LINK
		1e8,              // _linkGweiPerObservation uint32,
		4e8,              // _linkGweiPerTransmission uint32,
		linkTokenAddress, // _link common.Address,
		min,              // -2**191
		max,              // 2**191 - 1
		accessAddress,
		accessAddress,
		0,
		"TEST")
	require.NoError(t, err)
	_, err = linkContract.Transfer(owner, ocrContractAddress, big.NewInt(1000))
	require.NoError(t, err)

	flagsContractAddress, _, flagsContract, err := flags_wrapper.DeployFlags(owner, b.Client(), owner.From)
	require.NoError(t, err, "failed to deploy flags contract to simulated ethereum blockchain")

	b.Commit()
	return owner, b, ocrContractAddress, ocrContract, flagsContract, flagsContractAddress
}

func setupNode(t *testing.T, owner *bind.TransactOpts, portV2 int,
	b types.Backend, overrides func(c *chainlink.Config, s *chainlink.Secrets),
) (*cltest.TestApplication, string, common.Address, ocrkey.KeyV2) {
	ctx := testutils.Context(t)
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(portV2)))
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.OCR.Enabled = ptr(true)
		c.OCR2.Enabled = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())

		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", portV2)}
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

		// GracePeriod < ObservationTimeout
		c.EVM[0].OCR.ObservationGracePeriod = commonconfig.MustNewDuration(100 * time.Millisecond)

		if overrides != nil {
			overrides(c, s)
		}
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)
	transmitter := sendingKeys[0].Address

	// Fund the transmitter address with some ETH
	n, err := b.Client().NonceAt(testutils.Context(t), owner.From, nil)
	require.NoError(t, err)

	tx := cltest.NewLegacyTransaction(n, transmitter, assets.Ether(100).ToInt(), 21000, big.NewInt(1000000000), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.Client().SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	b.Commit()

	key, err := app.GetKeyStore().OCR().Create(ctx)
	require.NoError(t, err)
	return app, p2pKey.PeerID().Raw(), transmitter, key
}

func setupForwarderEnabledNode(t *testing.T, owner *bind.TransactOpts, portV2 int, b types.Backend, overrides func(c *chainlink.Config, s *chainlink.Secrets)) (*cltest.TestApplication, string, common.Address, common.Address, ocrkey.KeyV2) {
	ctx := testutils.Context(t)
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(portV2)))
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.OCR.Enabled = ptr(true)
		c.OCR2.Enabled = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", portV2)}
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

		c.EVM[0].Transactions.ForwardersEnabled = ptr(true)

		if overrides != nil {
			overrides(c, s)
		}
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
	require.NoError(t, err)
	transmitter := sendingKeys[0].Address

	// Fund the transmitter address with some ETH
	n, err := b.Client().NonceAt(testutils.Context(t), owner.From, nil)
	require.NoError(t, err)

	tx := cltest.NewLegacyTransaction(n, transmitter, assets.Ether(100).ToInt(), 21000, big.NewInt(1000000000), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.Client().SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	b.Commit()

	key, err := app.GetKeyStore().OCR().Create(ctx)
	require.NoError(t, err)

	// deploy a forwarder
	forwarder, _, authorizedForwarder, err := authorized_forwarder.DeployAuthorizedForwarder(owner, b.Client(), common.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB"), owner.From, common.Address{}, []byte{})
	require.NoError(t, err)

	// set EOA as an authorized sender for the forwarder
	_, err = authorizedForwarder.SetAuthorizedSenders(owner, []common.Address{transmitter})
	require.NoError(t, err)
	b.Commit()

	// add forwarder address to be tracked in db
	forwarderORM := forwarders.NewORM(app.GetDB())
	chainID, err := b.Client().ChainID(testutils.Context(t))
	require.NoError(t, err)
	_, err = forwarderORM.CreateForwarder(testutils.Context(t), forwarder, ubig.Big(*chainID))
	require.NoError(t, err)

	return app, p2pKey.PeerID().Raw(), transmitter, forwarder, key
}

func TestIntegration_OCR(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809; passes local but fails CI")
	tests.SkipShort(t, "long test")
	t.Parallel()
	tests := []struct {
		id      int
		name    string
		eip1559 bool
	}{
		{1, "legacy mode", false},
		{2, "eip1559 mode", true},
	}

	numOracles := 4
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			bootstrapNodePortV2 := freeport.GetOne(t)
			g := gomega.NewWithT(t)
			owner, b, ocrContractAddress, ocrContract, flagsContract, flagsContractAddress := setupOCRContracts(t)

			// Note it's plausible these ports could be occupied on a CI machine.
			// May need a port randomize + retry approach if we observe collisions.
			appBootstrap, bootstrapPeerID, _, _ := setupNode(t, owner, bootstrapNodePortV2, b, nil)
			var (
				oracles      []confighelper.OracleIdentityExtra
				transmitters []common.Address
				keys         []ocrkey.KeyV2
				apps         []*cltest.TestApplication
			)
			ports := freeport.GetN(t, numOracles)
			for i := 0; i < numOracles; i++ {
				app, peerID, transmitter, key := setupNode(t, owner, ports[i], b, func(c *chainlink.Config, s *chainlink.Secrets) {
					c.EVM[0].FlagsContractAddress = ptr(types.EIP55AddressFromAddress(flagsContractAddress))
					c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(test.eip1559)

					c.P2P.V2.DefaultBootstrappers = &[]ocrcommontypes.BootstrapperLocator{
						{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePortV2)}},
					}
				})

				keys = append(keys, key)
				apps = append(apps, app)
				transmitters = append(transmitters, transmitter)

				oracles = append(oracles, confighelper.OracleIdentityExtra{
					OracleIdentity: confighelper.OracleIdentity{
						OnChainSigningAddress: ocrtypes.OnChainSigningAddress(key.OnChainSigning.Address()),
						TransmitAddress:       transmitter,
						OffchainPublicKey:     key.PublicKeyOffChain(),
						PeerID:                peerID,
					},
					SharedSecretEncryptionPublicKey: key.PublicKeyConfig(),
				})
			}

			stopBlocks := utils.FiniteTicker(time.Second, func() {
				b.Commit()
			})
			defer stopBlocks()

			_, err := ocrContract.SetPayees(owner,
				transmitters,
				transmitters,
			)
			require.NoError(t, err)
			b.Commit()
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

			err = appBootstrap.Start(testutils.Context(t))
			require.NoError(t, err)

			jb, err := ocr.ValidatedOracleSpecToml(appBootstrap.Config, appBootstrap.GetRelayers().LegacyEVMChains(), fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "boot"
contractAddress    = "%s"
evmChainID		   = "%s"
isBootstrapPeer    = true
`, ocrContractAddress, testutils.SimulatedChainID.String()))
			require.NoError(t, err)
			jb.Name = null.NewString("boot", true)
			err = appBootstrap.AddJobV2(testutils.Context(t), &jb)
			require.NoError(t, err)

			// Raising flags to initiate hibernation
			_, err = flagsContract.RaiseFlag(owner, ocrContractAddress)
			require.NoError(t, err, "failed to raise flag for ocrContractAddress")
			_, err = flagsContract.RaiseFlag(owner, evmutils.ZeroAddress)
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
			for i := 0; i < numOracles; i++ {
				err = apps[i].Start(testutils.Context(t))
				require.NoError(t, err)

				// Since this API speed is > ObservationTimeout we should ignore it and still produce values.
				slowServers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					time.Sleep(5 * time.Second)
					res.WriteHeader(http.StatusOK)
					_, err := res.Write([]byte(`{"data":10}`))
					require.NoError(t, err)
				}))
				t.Cleanup(slowServers[i].Close)
				servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					b, err := io.ReadAll(req.Body)
					require.NoError(t, err)
					var m bridges.BridgeMetaDataJSON
					require.NoError(t, json.Unmarshal(b, &m))
					if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
						metaLock.Lock()
						delete(expectedMeta, m.Meta.LatestAnswer.String())
						metaLock.Unlock()
					}
					res.WriteHeader(http.StatusOK)
					_, err = res.Write([]byte(`{"data":10}`))
					require.NoError(t, err)
				}))
				t.Cleanup(servers[i].Close)
				u, _ := url.Parse(servers[i].URL)
				err := apps[i].BridgeORM().CreateBridgeType(testutils.Context(t), &bridges.BridgeType{
					Name: bridges.BridgeName(fmt.Sprintf("bridge%d", i)),
					URL:  models.WebURL(*u),
				})
				require.NoError(t, err)

				// Note we need: observationTimeout + observationGracePeriod + DeltaGrace (500ms) < DeltaRound (1s)
				// So 200ms + 200ms + 500ms < 1s
				jb, err := ocr.ValidatedOracleSpecToml(apps[i].Config, apps[i].GetRelayers().LegacyEVMChains(), fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "web oracle spec"
contractAddress    = "%s"
evmChainID		   = "%s"
isBootstrapPeer    = false
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
`, ocrContractAddress, testutils.SimulatedChainID.String(), keys[i].ID(), transmitters[i], fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i))
				require.NoError(t, err)
				jb.Name = null.NewString("testocr", true)
				err = apps[i].AddJobV2(testutils.Context(t), &jb)
				require.NoError(t, err)
				jids = append(jids, jb.ID)
			}

			// Assert that all the OCR jobs get a run with valid values eventually.
			for i := 0; i < numOracles; i++ {
				// Want at least 2 runs so we see all the metadata.
				pr := cltest.WaitForPipelineComplete(t, i, jids[i],
					2, 7, apps[i].JobORM(), time.Minute, time.Second)
				jb, err := pr[0].Outputs.MarshalJSON()
				require.NoError(t, err)
				assert.Equal(t, []byte(fmt.Sprintf("[\"%d\"]", 10*i)), jb, "pr[0] %+v pr[1] %+v", pr[0], pr[1])
				require.NoError(t, err)
			}

			// 4 oracles reporting 0, 10, 20, 30. Answer should be 20 (results[4/2]).
			g.Eventually(func() string {
				answer, err := ocrContract.LatestAnswer(nil)
				require.NoError(t, err)
				return answer.String()
			}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal("20"))

			ctx := testutils.Context(t)
			for _, app := range apps {
				jobs, _, err := app.JobORM().FindJobs(ctx, 0, 1000)
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
			metaLock.Lock()
			defer metaLock.Unlock()
			assert.Empty(t, expectedMeta, "expected metadata %v", expectedMeta)
		})
	}
}

func TestIntegration_OCR_ForwarderFlow(t *testing.T) {
	t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
	tests.SkipShort(t, "long test")
	t.Parallel()
	numOracles := 4
	t.Run("ocr_forwarder_flow", func(t *testing.T) {
		bootstrapNodePortV2 := freeport.GetOne(t)
		g := gomega.NewWithT(t)
		owner, b, ocrContractAddress, ocrContract, flagsContract, flagsContractAddress := setupOCRContracts(t)

		// Note it's plausible these ports could be occupied on a CI machine.
		// May need a port randomize + retry approach if we observe collisions.
		appBootstrap, bootstrapPeerID, _, _ := setupNode(t, owner, bootstrapNodePortV2, b, nil)

		var (
			oracles             []confighelper.OracleIdentityExtra
			transmitters        []common.Address
			forwardersContracts []common.Address
			keys                []ocrkey.KeyV2
			apps                []*cltest.TestApplication
		)
		ports := freeport.GetN(t, numOracles)
		for i := 0; i < numOracles; i++ {
			app, peerID, transmitter, forwarder, key := setupForwarderEnabledNode(t, owner, ports[i], b, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Feature.LogPoller = ptr(true)
				c.EVM[0].FlagsContractAddress = ptr(types.EIP55AddressFromAddress(flagsContractAddress))
				c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
				c.P2P.V2.DefaultBootstrappers = &[]ocrcommontypes.BootstrapperLocator{
					{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePortV2)}},
				}
			})

			keys = append(keys, key)
			apps = append(apps, app)
			forwardersContracts = append(forwardersContracts, forwarder)
			transmitters = append(transmitters, transmitter)

			oracles = append(oracles, confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnChainSigningAddress: ocrtypes.OnChainSigningAddress(key.OnChainSigning.Address()),
					TransmitAddress:       forwarder,
					OffchainPublicKey:     key.PublicKeyOffChain(),
					PeerID:                peerID,
				},
				SharedSecretEncryptionPublicKey: key.PublicKeyConfig(),
			})
		}

		stopBlocks := utils.FiniteTicker(time.Second, func() {
			b.Commit()
		})
		defer stopBlocks()

		_, err := ocrContract.SetPayees(owner,
			forwardersContracts,
			transmitters,
		)
		require.NoError(t, err)
		b.Commit()

		signers, effectiveTransmitters, threshold, encodedConfigVersion, encodedConfig, err := confighelper.ContractSetConfigArgsForIntegrationTest(
			oracles,
			1,
			1000000000/100, // threshold PPB
		)
		require.NoError(t, err)
		require.Equal(t, effectiveTransmitters, forwardersContracts)
		_, err = ocrContract.SetConfig(owner,
			signers,
			effectiveTransmitters, // forwarder Addresses
			threshold,
			encodedConfigVersion,
			encodedConfig,
		)
		require.NoError(t, err)
		b.Commit()

		err = appBootstrap.Start(testutils.Context(t))
		require.NoError(t, err)

		// set forwardingAllowed = true
		jb, err := ocr.ValidatedOracleSpecToml(appBootstrap.Config, appBootstrap.GetRelayers().LegacyEVMChains(), fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "boot"
contractAddress    = "%s"
evmChainID		   = "%s"
forwardingAllowed  = true
isBootstrapPeer    = true
`, ocrContractAddress, testutils.SimulatedChainID.String()))
		require.NoError(t, err)
		jb.Name = null.NewString("boot", true)
		err = appBootstrap.AddJobV2(testutils.Context(t), &jb)
		require.NoError(t, err)

		// Raising flags to initiate hibernation
		_, err = flagsContract.RaiseFlag(owner, ocrContractAddress)
		require.NoError(t, err, "failed to raise flag for ocrContractAddress")
		_, err = flagsContract.RaiseFlag(owner, evmutils.ZeroAddress)
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
		for i := 0; i < numOracles; i++ {
			err = apps[i].Start(testutils.Context(t))
			require.NoError(t, err)

			// Since this API speed is > ObservationTimeout we should ignore it and still produce values.
			slowServers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				time.Sleep(5 * time.Second)
				res.WriteHeader(http.StatusOK)
				_, err := res.Write([]byte(`{"data":10}`))
				require.NoError(t, err)
			}))
			t.Cleanup(slowServers[i].Close)
			servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				b, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				var m bridges.BridgeMetaDataJSON
				require.NoError(t, json.Unmarshal(b, &m))
				if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
					metaLock.Lock()
					delete(expectedMeta, m.Meta.LatestAnswer.String())
					metaLock.Unlock()
				}
				res.WriteHeader(http.StatusOK)
				_, err = res.Write([]byte(`{"data":10}`))
				require.NoError(t, err)
			}))
			t.Cleanup(servers[i].Close)
			u, _ := url.Parse(servers[i].URL)
			err := apps[i].BridgeORM().CreateBridgeType(testutils.Context(t), &bridges.BridgeType{
				Name: bridges.BridgeName(fmt.Sprintf("bridge%d", i)),
				URL:  models.WebURL(*u),
			})
			require.NoError(t, err)

			// Note we need: observationTimeout + observationGracePeriod + DeltaGrace (500ms) < DeltaRound (1s)
			// So 200ms + 200ms + 500ms < 1s
			// forwardingAllowed = true
			jb, err := ocr.ValidatedOracleSpecToml(apps[i].Config, apps[i].GetRelayers().LegacyEVMChains(), fmt.Sprintf(`
type               = "offchainreporting"
schemaVersion      = 1
name               = "web oracle spec"
contractAddress    = "%s"
evmChainID         = "%s"
forwardingAllowed  = true
isBootstrapPeer    = false
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
`, ocrContractAddress, testutils.SimulatedChainID.String(), keys[i].ID(), transmitters[i], fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i))
			require.NoError(t, err)
			jb.Name = null.NewString("testocr", true)
			err = apps[i].AddJobV2(testutils.Context(t), &jb)
			require.NoError(t, err)
			jids = append(jids, jb.ID)
		}

		// Assert that all the OCR jobs get a run with valid values eventually.
		for i := 0; i < numOracles; i++ {
			// Want at least 2 runs so we see all the metadata.
			pr := cltest.WaitForPipelineComplete(t, i, jids[i],
				2, 7, apps[i].JobORM(), time.Minute, time.Second)
			jb, err := pr[0].Outputs.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, []byte(fmt.Sprintf("[\"%d\"]", 10*i)), jb, "pr[0] %+v pr[1] %+v", pr[0], pr[1])
			require.NoError(t, err)
		}

		// 4 oracles reporting 0, 10, 20, 30. Answer should be 20 (results[4/2]).
		g.Eventually(func() string {
			answer, err := ocrContract.LatestAnswer(nil)
			require.NoError(t, err)
			return answer.String()
		}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal("20"))

		ctx := testutils.Context(t)
		for _, app := range apps {
			jobs, _, err := app.JobORM().FindJobs(ctx, 0, 1000)
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
		metaLock.Lock()
		defer metaLock.Unlock()
		assert.Empty(t, expectedMeta, "expected metadata %v", expectedMeta)
	})
}

func TestIntegration_BlockHistoryEstimator(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	var initialDefaultGasPrice int64 = 5_000_000_000
	maxGasPrice := assets.NewWeiI(10 * initialDefaultGasPrice)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.BlockHistory.CheckInclusionBlocks = ptr[uint16](0)
		c.EVM[0].GasEstimator.PriceDefault = assets.NewWeiI(initialDefaultGasPrice)
		c.EVM[0].GasEstimator.Mode = ptr("BlockHistory")
		c.EVM[0].RPCBlockQueryDelay = ptr[uint16](0)
		c.EVM[0].GasEstimator.BlockHistory.BlockHistorySize = ptr[uint16](2)
		c.EVM[0].FinalityDepth = ptr[uint32](3)
	})

	ethClient := cltest.NewEthMocks(t)
	ethClient.On("ConfiguredChainID").Return(big.NewInt(client.NullClientChainID)).Maybe()
	chchNewHeads := make(chan evmtestutils.RawSub[*types.Head], 1)

	db := pgtest.NewSqlxDB(t)
	kst := cltest.NewKeyStore(t, db)
	require.NoError(t, kst.Unlock(ctx, cltest.Password))

	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		KeyStore:       kst.Eth(),
		DB:             db,
		Client:         ethClient,
	})

	b41 := types.Block{
		Number:       41,
		Hash:         evmutils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(41_000_000_000, 41_500_000_000),
	}
	b42 := types.Block{
		Number:       42,
		Hash:         evmutils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(44_000_000_000, 45_000_000_000),
	}
	b43 := types.Block{
		Number:       43,
		Hash:         evmutils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(48_000_000_000, 49_000_000_000, 31_000_000_000),
	}

	evmChainID := ubig.New(evmtest.MustGetDefaultChainID(t, cfg.EVMConfigs()))
	h40 := types.Head{Hash: evmutils.NewHash(), Number: 40, EVMChainID: evmChainID}
	h41 := types.Head{Hash: b41.Hash, ParentHash: h40.Hash, Number: 41, EVMChainID: evmChainID}
	h42 := types.Head{Hash: b42.Hash, ParentHash: h41.Hash, Number: 42, EVMChainID: evmChainID}

	mockEth := &clienttest.MockEth{EthClient: ethClient}
	ethClient.On("SubscribeToHeads", mock.Anything).
		Return(
			func(ctx context.Context) (<-chan *types.Head, ethereum.Subscription, error) {
				ch := make(chan *types.Head)
				sub := mockEth.NewSub(t)
				chchNewHeads <- evmtestutils.NewRawSub(ch, sub.Err())
				return ch, sub, nil
			},
		)
	// Nonce syncer
	ethClient.On("PendingNonceAt", mock.Anything, mock.Anything).Maybe().Return(uint64(0), nil)

	// BlockHistoryEstimator boot calls
	ethClient.On("HeadByNumber", mock.Anything, mock.AnythingOfType("*big.Int")).Return(&h42, nil)
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2a" &&
			b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x29"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b42
		elems[1].Result = &b41
	})

	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("ConfiguredChainID", mock.Anything).Return(*evmtest.MustGetDefaultChainID(t, cfg.EVMConfigs()), nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(oneETH.ToInt(), nil)
	// HeadTracker backfill
	ethClient.On("HeadByHash", mock.Anything, h40.Hash).Return(&h40, nil).Maybe()
	ethClient.On("HeadByHash", mock.Anything, h41.Hash).Return(&h41, nil).Maybe()
	ethClient.On("HeadByHash", mock.Anything, h42.Hash).Return(&h42, nil).Maybe()

	for _, re := range legacyChains.Slice() {
		servicetest.Run(t, re)
	}
	var newHeads evmtestutils.RawSub[*types.Head]
	select {
	case newHeads = <-chchNewHeads:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for app to subscribe")
	}

	chain := evmtest.MustGetDefaultChain(t, legacyChains)
	estimator := chain.GasEstimator()
	gasPrice, gasLimit, err := estimator.GetFee(testutils.Context(t), nil, 500_000, maxGasPrice, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(500000), gasLimit)
	assert.Equal(t, "41.5 gwei", gasPrice.GasPrice.String())
	assert.Equal(t, initialDefaultGasPrice, chain.Config().EVM().GasEstimator().PriceDefault().Int64()) // unchanged

	// BlockHistoryEstimator new blocks
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 && b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2b"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b43
	})
	ethClient.On("Close").Return().Once()

	// Simulate one new head and check the gas price got updated
	h43 := cltest.Head(43)
	h43.ParentHash = h42.Hash
	newHeads.TrySend(h43)

	require.Eventually(t, func() bool {
		gasPrice, _, err := estimator.GetFee(testutils.Context(t), nil, 500000, maxGasPrice, nil, nil)
		require.NoError(t, err)
		return gasPrice.GasPrice.String() == "45 gwei"
	}, testutils.WaitTimeout(t), cltest.DBPollingInterval)
}

func triggerAllKeys(t *testing.T, app *cltest.TestApplication) {
	for _, chain := range app.GetRelayers().LegacyEVMChains().Slice() {
		keys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.Context(t), chain.ID())
		require.NoError(t, err)
		for _, k := range keys {
			chain.TxManager().Trigger(k.Address)
		}
	}
}

func assertPricesUint256(t *testing.T, usd, eur, jpy *big.Int, consumer *multiwordconsumer_wrapper.MultiWordConsumer) {
	haveUsd, err := consumer.UsdInt(nil)
	require.NoError(t, err)
	assert.Equal(t, usd.Cmp(haveUsd), 0)
	haveEur, err := consumer.EurInt(nil)
	require.NoError(t, err)
	assert.Equal(t, eur.Cmp(haveEur), 0)
	haveJpy, err := consumer.JpyInt(nil)
	require.NoError(t, err)
	assert.Equal(t, jpy.Cmp(haveJpy), 0)
}

func ptr[T any](v T) *T { return &v }

func assertPipelineTaskRunsSuccessful(t testing.TB, runs []pipeline.TaskRun) {
	t.Helper()
	for i, run := range runs {
		require.True(t, run.Error.IsZero(), "pipeline.Task run failed (idx: %v, dotID: %v, error: '%v')", i, run.GetDotID(), run.Error.ValueOrZero())
	}
}

func assertPipelineTaskRunsErrored(t testing.TB, runs []pipeline.TaskRun) {
	t.Helper()
	for i, run := range runs {
		require.False(t, run.Error.IsZero(), "expected pipeline.Task run to have failed, but it succeeded (idx: %v, dotID: %v, output: '%v')", i, run.GetDotID(), run.Output)
	}
}
