package services_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
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
	eth.RegisterSubscription("logs", make(chan []types.Log))
	eth.RegisterSubscription("logs", make(chan []types.Log))

	err := nl.Start()
	assert.Nil(t, err)

	assert.True(t, eth.AllCalled())
}

func TestNotificationListenerAddJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	initrAddress := cltest.NewEthAddress()

	tests := []struct {
		name       string
		initType   string
		logAddress common.Address
		wantCount  int
		data       hexutil.Bytes
	}{
		{"basic eth log", "ethlog", initrAddress, 1, hexutil.Bytes{}},
		{"non-matching eth log", "ethlog", cltest.NewEthAddress(), 0, hexutil.Bytes{}},
		{"basic cllog", "runlog", initrAddress, 1, cltest.StringToRunLogPayload(`{"value":"100"}`)},
		{"cllog non-matching", "runlog", cltest.NewEthAddress(), 0, hexutil.Bytes{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RegisterTestingT(t)

			nl := services.NotificationListener{Store: store}
			defer nl.Stop()
			err := nl.Start()
			assert.Nil(t, err)

			eth := cltest.MockEthOnStore(store)
			logChan := make(chan []types.Log, 1)
			eth.RegisterSubscription("logs", logChan)

			j := cltest.NewJob()
			j.Initiators = []models.Initiator{{
				Type:    test.initType,
				Address: initrAddress,
			}}
			assert.Nil(t, store.SaveJob(j))

			nl.AddJob(j)

			logChan <- []types.Log{{
				Address: test.logAddress,
				Data:    test.data,
				Topics:  []common.Hash{common.HexToHash("0x00"), common.HexToHash("0x01"), common.HexToHash("0x22")},
			}}
			<-time.After(100 * time.Millisecond)

			cltest.WaitForRuns(t, j, store, test.wantCount)

			assert.True(t, eth.AllCalled())
		})
	}
}

func outputFromFixture(path string) models.JSON {
	res := gjson.Get(string(cltest.LoadJSON(path)), "params.result.0")
	var out models.JSON
	if err := json.Unmarshal([]byte(res.String()), &out); err != nil {
		panic(err)
	}
	return out
}

func TestStoreFormatLogJSON(t *testing.T) {
	t.Parallel()

	var clData models.JSON
	clDataFixture := `{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000001","functionSelector":"76005c26"}`
	assert.Nil(t, json.Unmarshal([]byte(clDataFixture), &clData))

	hwLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	exampleLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          types.Log
		initr       models.Initiator
		wantErrored bool
		wantOutput  models.JSON
	}{
		{"example ethLog", exampleLog, models.Initiator{Type: "ethlog"}, false,
			outputFromFixture("../internal/fixtures/eth/subscription_logs.json")},
		{"hello world ethLog", hwLog, models.Initiator{Type: "ethlog"}, false,
			outputFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")},
		{"hello world runLog", hwLog, models.Initiator{Type: "runlog"}, false,
			clData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := services.FormatLogJSON(test.initr, test.el)
			assert.JSONEq(t, strings.ToLower(test.wantOutput.String()), strings.ToLower(output.String()))
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}
