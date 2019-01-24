package models_test

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunLogEvent_JSON(t *testing.T) {
	t.Parallel()

	clData := cltest.JSONFromString(`{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000017","functionSelector":"0x76005c26"}`)
	hwLog := cltest.LogFromFixture("../../internal/fixtures/eth/subscription_logs_hello_world.json")
	tests := []struct {
		name        string
		el          models.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"hello world", hwLog, false, clData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initr := models.Initiator{Type: models.InitiatorRunLog}
			le := models.InitiatorLogEvent{Initiator: initr, Log: test.el}.LogRequest()
			output, err := le.JSON()
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.NoError(t, err)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestEthLogEvent_JSON(t *testing.T) {
	t.Parallel()

	hwLog := cltest.LogFromFixture("../../internal/fixtures/eth/subscription_logs_hello_world.json")
	exampleLog := cltest.LogFromFixture("../../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          models.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"example", exampleLog, false, cltest.JSONResultFromFixture("../../internal/fixtures/eth/subscription_logs.json")},
		{"hello world", hwLog, false, cltest.JSONResultFromFixture("../../internal/fixtures/eth/subscription_logs_hello_world.json")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initr := models.Initiator{Type: models.InitiatorEthLog}
			le := models.InitiatorLogEvent{Initiator: initr, Log: test.el}.LogRequest()
			output, err := le.JSON()
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestRequestLogEvent_Requester(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		logFactory (func(string, common.Address, common.Address, int, string) models.Log)
		input      common.Hash
		want       common.Address
	}{
		{
			"runlog basic",
			cltest.NewRunLog,
			common.HexToHash("0x00000000000000000000000059b15a7ae74c803cc151ffe63042faa826c96eee"),
			common.HexToAddress("0x59b15a7ae74c803cc151ffe63042faa826c96eee"),
		},
		{
			"salog basic",
			cltest.NewServiceAgreementExecutionLog,
			common.HexToHash("0x00000000000000000000000059b15a7ae74c803cc151ffe63042faa826c96eee"),
			common.HexToAddress("0x59b15a7ae74c803cc151ffe63042faa826c96eee"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rl := cltest.NewRunLog("id", cltest.NewAddress(), cltest.NewAddress(), 0, "{}")
			rl.Topics[models.RequestLogTopicRequester] = test.input
			le := models.RunLogEvent{models.InitiatorLogEvent{Log: rl}}

			assert.Equal(t, test.want, le.Requester())
		})
	}
}

func TestRequestLogEvent_Validate(t *testing.T) {
	t.Parallel()

	job := cltest.NewJob()
	job.ID = "4a1eb0e8df314cb894024a38991cff0f"

	noRequesters := []common.Address{}
	permittedAddr := cltest.NewAddress()
	unpermittedAddr := cltest.NewAddress()
	requesterList := []common.Address{permittedAddr}

	tests := []struct {
		name                string
		eventLogTopic       common.Hash
		jobIDTopic          common.Hash
		initiatorRequesters []common.Address
		requesterAddress    common.Address
		want                bool
	}{
		{"wrong jobid", models.RunLogTopic, cltest.StringToHash("wrongjob"), noRequesters, unpermittedAddr, false},
		{"proper hex jobid", models.RunLogTopic, cltest.StringToHash(job.ID), noRequesters, unpermittedAddr, true},
		{"incorrect encoded jobid", models.RunLogTopic, common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"), noRequesters, unpermittedAddr, true},
		{"correct requester", models.RunLogTopic, cltest.StringToHash(job.ID), requesterList, permittedAddr, true},
		{"incorrect requester", models.RunLogTopic, cltest.StringToHash(job.ID), requesterList, unpermittedAddr, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Any log factory works since we overwrite topics.
			log := cltest.NewRunLog(
				job.ID, cltest.NewAddress(),
				test.requesterAddress, 1, "{}")

			log.Topics = []common.Hash{
				test.eventLogTopic,
				test.jobIDTopic,
				test.requesterAddress.Hash(),
				{},
			}

			logRequest := models.InitiatorLogEvent{
				JobSpec: job,
				Log:     log,
				Initiator: models.Initiator{
					Type: models.InitiatorRunLog,
					InitiatorParams: models.InitiatorParams{
						Requesters: test.initiatorRequesters,
					},
				},
			}.LogRequest()

			assert.Equal(t, test.want, logRequest.Validate())
		})
	}
}

func TestStartRunOrSALogSubscription_ValidateSenders(t *testing.T) {
	requester := cltest.NewAddress()

	tests := []struct {
		name       string
		job        models.JobSpec
		requester  common.Address
		logFactory (func(string, common.Address, common.Address, int, string) models.Log)
		wantStatus models.RunStatus
	}{
		{
			"runlog contains valid requester",
			first(cltest.NewJobWithRunLogInitiator()),
			requester,
			cltest.NewRunLog,
			models.RunStatusCompleted,
		},
		{
			"runlog has wrong requester",
			first(cltest.NewJobWithRunLogInitiator()),
			cltest.NewAddress(),
			cltest.NewRunLog,
			models.RunStatusErrored,
		},
		{
			"salog contains valid requester",
			first(cltest.NewJobWithSALogInitiator()),
			requester,
			cltest.NewServiceAgreementExecutionLog,
			models.RunStatusCompleted,
		},
		{
			"salog has wrong requester",
			first(cltest.NewJobWithSALogInitiator()),
			cltest.NewAddress(),
			cltest.NewServiceAgreementExecutionLog,
			models.RunStatusErrored,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, _ := cltest.NewConfigWithPrivateKey("../../internal/fixtures/keys/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json")
			app, cleanup := cltest.NewApplicationWithConfigAndUnlockedAccount(config)
			defer cleanup()

			eth := app.MockEthClient()
			logs := make(chan models.Log, 1)
			eth.Context("app.Start()", func(eth *cltest.EthMock) {
				eth.Register("eth_getBlockByNumber", models.BlockHeader{})
				eth.Register("eth_getTransactionCount", "0x1")
				eth.RegisterSubscription("logs", logs)
			})
			assert.NoError(t, app.Start())

			js := test.job
			initr := js.Initiators[0]
			initr.Requesters = []common.Address{requester}
			_, err := services.NewInitiatorSubscription(initr, js, app.Store, nil, services.ReceiveLogRequest)
			assert.NoError(t, err)

			logs <- test.logFactory(js.ID, cltest.NewAddress(), test.requester, 1, `{}`)
			eth.EventuallyAllCalled(t)

			gomega.NewGomegaWithT(t).Eventually(func() models.RunStatus {
				runs, err := app.Store.JobRunsFor(js.ID)
				require.NoError(t, err)
				return runs[0].Status
			}).Should(gomega.Equal(test.wantStatus))
		})
	}
}

func first(a models.JobSpec, b interface{}) models.JobSpec {
	return a
}
