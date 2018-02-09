package services_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestNotificationListenerStart(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	nl := services.NotificationListener{Store: store}
	defer nl.Stop()

	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	eth.RegisterSubscription("logs", make(chan strpkg.EthNotification))
	eth.RegisterSubscription("logs", make(chan strpkg.EthNotification))

	nl.Start()

	assert.True(t, eth.AllCalled())
}

func TestNotificationListenerAddJob(t *testing.T) {
	t.Parallel()

	initrAddress := cltest.NewEthAddress()
	eventTopic := "0x06f4bf36b4e011a5c499cef1113c2d166800ce4013f6c2509cab1a0e92b83fb2"
	nonce := "0x402e9676507fe53f4bbb5d771a6597c663bf72ae82bb707356a2614f83fd5b15"
	jobID := "ced8524ac7c14c05aae482e02a098b87"
	jobIDHex := "0x6365643835323461633763313463303561616534383265303261303938623837"

	tests := []struct {
		name       string
		initType   string
		logAddress common.Address
		wantCount  int
		wantJSON   models.JSON
		data       hexutil.Bytes
	}{
		{"basic eth log", "ethlog", initrAddress, 1,
			cltest.JSONFromString(
				`{"address":"%v","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","blockNumber":"0x0","data":"0x","logIndex":"0x0","topics":["%v","%v","%v"],"transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0"}`,
				strings.ToLower(initrAddress.String()), eventTopic, nonce, jobIDHex), hexutil.Bytes{}},
		{"non-matching eth log", "ethlog", cltest.NewEthAddress(), 0, models.JSON{}, hexutil.Bytes{}},
		{"basic cllog", "chainlinklog", initrAddress, 1,
			cltest.JSONFromString(
				`{"value":"100","address":"%v","functionId":"76005c26","dataPrefix":"%v"}`, initrAddress.String(), nonce),
			cltest.StringToRunLogPayload(`{"value":"100"}`)},
		{"cllog non-matching", "chainlinklog", cltest.NewEthAddress(), 0, models.JSON{}, hexutil.Bytes{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RegisterTestingT(t)

			store, cleanup := cltest.NewStore()
			defer cleanup()
			nl := services.NotificationListener{Store: store}
			defer nl.Stop()
			nl.Start()

			eth := cltest.MockEthOnStore(store)
			logChan := make(chan strpkg.EthNotification, 1)
			eth.RegisterSubscription("logs", logChan)

			j := cltest.NewJob()
			j.Initiators = []models.Initiator{{
				Type:    test.initType,
				Address: initrAddress,
			}}
			j.ID = jobID
			assert.Nil(t, store.SaveJob(j))

			nl.AddJob(j)

			logChan <- cltest.NewEthNotification(strpkg.EventLog{
				Address: test.logAddress,
				Data:    test.data,
				Topics: []hexutil.Bytes{
					cltest.StringToBytes(eventTopic),
					cltest.StringToBytes(nonce),
					cltest.StringToBytes(jobIDHex),
				},
			})

			if test.wantCount == 0 {
				Consistently(func() []models.JobRun {
					jrs, _ := store.JobRunsFor(j)
					return jrs
				}).Should(HaveLen(test.wantCount))
			} else {
				var jrs []models.JobRun
				Eventually(func() []models.JobRun {
					jrs, _ = store.JobRunsFor(j)
					return jrs
				}).Should(HaveLen(test.wantCount))
				jr := jrs[0]
				params := jr.TaskRuns[0].Task.Params
				assert.JSONEq(t, test.wantJSON.String(), params.String())
			}

			assert.True(t, eth.AllCalled())
		})
	}
}

func outputFromFixture(path string) models.JSON {
	fix := cltest.JSONFromFixture(path)
	res := fix.Get("params.result")
	var out models.JSON
	if err := json.Unmarshal([]byte(res.String()), &out); err != nil {
		panic(err)
	}
	return out
}

func TestStoreFormatLogOutput(t *testing.T) {
	t.Parallel()

	var clData models.JSON
	clDataFixture := `{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000001","functionId":"76005c26"}`
	assert.Nil(t, json.Unmarshal([]byte(clDataFixture), &clData))

	hwEvent := cltest.EventLogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	exampleEvent := cltest.EventLogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          strpkg.EventLog
		initr       models.Initiator
		wantErrored bool
		wantOutput  models.JSON
	}{
		{"example ethLog", exampleEvent, models.Initiator{Type: "ethlog"}, false,
			outputFromFixture("../internal/fixtures/eth/subscription_logs.json")},
		{"hello world ethLog", hwEvent, models.Initiator{Type: "ethlog"}, false,
			outputFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")},
		{"hello world chainlinkLog", hwEvent, models.Initiator{Type: "chainlinklog"}, false,
			clData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := services.FormatLogOutput(test.initr, test.el)
			assert.JSONEq(t, test.wantOutput.String(), output.String())
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}
