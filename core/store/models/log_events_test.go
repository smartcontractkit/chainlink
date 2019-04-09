package models_test

import (
	"math/big"
	"strings"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRunLog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		log         models.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{
			name:        "old non-commitment",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog0original.json"),
			wantErrored: false,
			wantData: cltest.JSONFromString(t, `{
				"url":"https://etherprice.com/api",
				"path":["recent","usd"],
				"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42",
				"dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000017",
				"functionSelector":"0x76005c26"}`),
		},
		{
			name:        "20190123 fulfillment params",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190123withFulfillmentParams.json"),
			wantErrored: false,
			wantData: cltest.JSONFromString(t, `{
				"url":"https://min-api.cryptocompare.com/data/price?fsym=eth&tsyms=usd,eur,jpy",
				"path":["usd"],
				"address":"0xf25186b5081ff5ce73482ad761db0eb0d25abfbf",
				"dataPrefix":"0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f80000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000009fbda871d559710256a2502a2517b794b482db40042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005c4a7338",
				"functionSelector":"0xeea57e70"}`),
		},
		{
			name:        "20190207 without indexes",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190207withoutIndexes.json"),
			wantErrored: false,
			wantData: cltest.JSONFromString(t, `{
				"url":"https://min-api.cryptocompare.com/data/price?fsym=eth&tsyms=usd,eur,jpy",
				"path":["usd"],
				"address":"0xf25186b5081ff5ce73482ad761db0eb0d25abfbf",
				"dataPrefix":"0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f80000000000000000000000000000000000000000000000000de0b6b3a76400010000000000000000000000009fbda871d559710256a2502a2517b794b482db40042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005c4a7338",
				"functionSelector":"0x4ab0d190"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := models.ParseRunLog(test.log)
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.NoError(t, err)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestEthLogEvent_JSON(t *testing.T) {
	t.Parallel()

	hwLog := cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog0original.json")
	exampleLog := cltest.LogFromFixture(t, "../../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          models.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"example", exampleLog, false, cltest.JSONResultFromFixture(t, "../../internal/fixtures/eth/subscription_logs.json")},
		{"hello world", hwLog, false, cltest.JSONResultFromFixture(t, "../../internal/fixtures/eth/requestLog0original.json")},
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
		{"wrong jobid", models.RunLogTopic0original, cltest.StringToHash("wrongjob"), noRequesters, unpermittedAddr, false},
		{"proper hex jobid", models.RunLogTopic0original, cltest.StringToHash(job.ID), noRequesters, unpermittedAddr, true},
		{"incorrect encoded jobid", models.RunLogTopic0original, common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"), noRequesters, unpermittedAddr, true},
		{"correct requester", models.RunLogTopic0original, cltest.StringToHash(job.ID), requesterList, permittedAddr, true},
		{"incorrect requester", models.RunLogTopic0original, cltest.StringToHash(job.ID), requesterList, unpermittedAddr, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Any log factory works since we overwrite topics.
			log := cltest.NewRunLog(
				t,
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
		logFactory (func(*testing.T, string, common.Address, common.Address, int, string) models.Log)
		wantStatus models.RunStatus
	}{
		{
			"runlog contains valid requester",
			cltest.NewJobWithRunLogInitiator(),
			requester,
			cltest.NewRunLog,
			models.RunStatusCompleted,
		},
		{
			"runlog has wrong requester",
			cltest.NewJobWithRunLogInitiator(),
			cltest.NewAddress(),
			cltest.NewRunLog,
			models.RunStatusErrored,
		},
		{
			"salog contains valid requester",
			cltest.NewJobWithSALogInitiator(),
			requester,
			cltest.NewServiceAgreementExecutionLog,
			models.RunStatusCompleted,
		},
		{
			"salog has wrong requester",
			cltest.NewJobWithSALogInitiator(),
			cltest.NewAddress(),
			cltest.NewServiceAgreementExecutionLog,
			models.RunStatusErrored,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithKey()
			defer cleanup()

			eth := app.MockEthClient()
			logs := make(chan models.Log, 1)
			eth.Context("app.Start()", func(eth *cltest.EthMock) {
				eth.Register("eth_getTransactionCount", "0x1")
				eth.RegisterSubscription("logs", logs)
			})
			assert.NoError(t, app.StartAndConnect())

			js := test.job
			js.Initiators[0].Requesters = []common.Address{requester}
			require.NoError(t, app.AddJob(js))

			logs <- test.logFactory(t, js.ID, cltest.NewAddress(), test.requester, 1, `{}`)
			eth.EventuallyAllCalled(t)

			gomega.NewGomegaWithT(t).Eventually(func() []models.JobRun {
				runs, err := app.Store.JobRunsFor(js.ID)
				require.NoError(t, err)
				return runs
			}).Should(gomega.HaveLen(1))

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

func TestFilterQueryFactory_InitiatorRunLog(t *testing.T) {
	t.Parallel()

	i := models.Initiator{
		Type:      models.InitiatorRunLog,
		JobSpecID: "4a1eb0e8df314cb894024a38991cff0f",
	}
	fromBlock := big.NewInt(42)
	filter, err := models.FilterQueryFactory(i, fromBlock)
	assert.NoError(t, err)

	want := ethereum.FilterQuery{
		FromBlock: fromBlock.Add(fromBlock, big.NewInt(1)),
		Topics: [][]common.Hash{
			{
				models.RunLogTopic20190207withoutIndexes,
				models.RunLogTopic20190123withFullfillmentParams,
				models.RunLogTopic0original,
			}, {
				common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
				common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
			},
		},
	}
	assert.Equal(t, want, filter)
}

func TestRunLogEvent_ContractPayment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		log         models.Log
		wantErrored bool
		want        *assets.Link
	}{
		{
			name:        "old non-commitment",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog0original.json"),
			wantErrored: false,
			want:        assets.NewLink(1),
		},
		{
			name:        "20190123 with fulfillment params",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190123withFulfillmentParams.json"),
			wantErrored: false,
			want:        assets.NewLink(1000000000000000000),
		},
		{
			name:        "20190207 without indexes",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190207withoutIndexes.json"),
			wantErrored: false,
			want:        assets.NewLink(1000000000000000001),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rle := models.RunLogEvent{models.InitiatorLogEvent{Log: test.log}}

			received, err := rle.ContractPayment()

			cltest.AssertError(t, test.wantErrored, err)
			assert.Equal(t, test.want, received)
		})
	}
}

func TestRunLogEvent_Requester(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		log         models.Log
		wantErrored bool
		want        common.Address
	}{
		{
			name:        "old non-commitment",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog0original.json"),
			wantErrored: false,
			want:        common.HexToAddress("0xd352677fcded6c358e03c73ea2a8a2832dffc0a4"),
		},
		{
			name:        "20190123 with fulfillment params",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190123withFulfillmentParams.json"),
			wantErrored: false,
			want:        common.HexToAddress("0x9fbda871d559710256a2502a2517b794b482db41"),
		},
		{
			name:        "20190207 without indexes",
			log:         cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190207withoutIndexes.json"),
			wantErrored: false,
			want:        common.HexToAddress("0x9fbda871d559710256a2502a2517b794b482db40"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rle := models.RunLogEvent{models.InitiatorLogEvent{Log: test.log}}

			received := rle.Requester()

			assert.Equal(t, test.want, received)
		})
	}
}

func TestRunLogEvent_RunRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		log           models.Log
		wantRequestID string
		wantTxHash    string
		wantRequester common.Address
	}{
		{
			name:          "old non-commitment",
			log:           cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog0original.json"),
			wantRequestID: "0x0000000000000000000000000000000000000000000000000000000000000017",
			wantTxHash:    "0xe05b171038320aca6634ce50de669bd0baa337130269c3ce3594ce4d45fc342a",
			wantRequester: common.HexToAddress("0xd352677fcded6c358e03c73ea2a8a2832dffc0a4"),
		},
		{
			name:          "20190123 with fulfillment params",
			log:           cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190123withFulfillmentParams.json"),
			wantRequestID: "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8",
			wantTxHash:    "0x04250548cd0b5d03b3bf1331aa83f32b35879440db31a6008d151260a5f3cc76",
			wantRequester: common.HexToAddress("0x9fbda871d559710256a2502a2517b794b482db41"),
		},
		{
			name:          "20190207 without indexes",
			log:           cltest.LogFromFixture(t, "../../internal/fixtures/eth/requestLog20190207withoutIndexes.json"),
			wantRequestID: "0xc524fafafcaec40652b1f84fca09c231185437d008d195fccf2f51e64b7062f8",
			wantTxHash:    "0x04250548cd0b5d03b3bf1331aa83f32b35879440db31a6008d151260a5f3cc76",
			wantRequester: common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rle := models.RunLogEvent{models.InitiatorLogEvent{Log: test.log}}
			rr, err := rle.RunRequest()
			require.NoError(t, err)

			assert.Equal(t, &test.wantRequestID, rr.RequestID)
			assert.Equal(t, test.wantTxHash, rr.TxHash.Hex())
			assert.Equal(t, &test.wantRequester, rr.Requester)
		})
	}
}
