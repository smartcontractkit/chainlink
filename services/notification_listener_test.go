package services_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	RegisterTestingT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	initrAddress := cltest.NewEthAddress()

	tests := []struct {
		name       string
		initType   string
		logAddress common.Address
		want       int
	}{
		{"basic eth log", "ethlog", initrAddress, 1},
		{"non-matching eth log", "ethlog", cltest.NewEthAddress(), 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
			assert.Nil(t, store.SaveJob(j))

			nl.AddJob(j)

			logChan <- cltest.NewEthNotification(strpkg.EventLog{
				Address: test.logAddress,
			})
			<-time.After(100 * time.Millisecond)

			jrs, err := store.JobRunsFor(j)
			assert.Nil(t, err)
			assert.Equal(t, test.want, len(jrs))

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
	clDataFixture := `{"url":"https://etherprice.com/api","path":["recent","usd"]}`
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
