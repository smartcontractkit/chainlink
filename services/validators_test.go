package services_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidateJob(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []byte
		want  error
	}{
		{"base case", cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json"), nil},
		{
			"error in job",
			cltest.LoadJSON("../internal/fixtures/web/invalid_endat_job.json"),
			models.NewJSONAPIErrorsWith("StartAt cannot be before EndAt"),
		},
		{
			"error in runat initr",
			cltest.LoadJSON("../internal/fixtures/web/run_at_wo_time_job.json"),
			models.NewJSONAPIErrorsWith("RunAt must have a time"),
		},
		{
			"error in task",
			cltest.LoadJSON("../internal/fixtures/web/nonexistent_task_job.json"),
			models.NewJSONAPIErrorsWith("idonotexist is not a supported adapter type"),
		},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var j models.JobSpec
			assert.NoError(t, json.Unmarshal(test.input, &j))
			result := services.ValidateJob(j, store)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateAdapter(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	tt := models.BridgeType{}
	tt.Name = models.MustNewTaskType("solargridreporting")
	tt.URL = cltest.WebURL("https://denergy.eth")
	assert.NoError(t, store.Save(&tt))

	tests := []struct {
		description string
		name        string
		want        error
	}{
		{
			"existing external adapter",
			"solargridreporting",
			models.NewJSONAPIErrorsWith("Adapter solargridreporting already exists"),
		},
		{
			"existing core adapter",
			"ethtx",
			models.NewJSONAPIErrorsWith("Adapter ethtx already exists"),
		},
		{
			"no adapter name",
			"",
			models.NewJSONAPIErrorsWith("No name specified"),
		},
		{
			"invalid adapter name",
			"invalid/adapter",
			models.NewJSONAPIErrorsWith("Task Type validation: name invalid/adapter contains invalid characters"),
		},
		{"new external adapter", "gdaxprice", nil},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			bt := &models.BridgeType{Name: models.TaskType(test.name)}
			result := services.ValidateAdapter(bt, store)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateInitiator(t *testing.T) {
	t.Parallel()
	startAt := time.Now()
	endAt := startAt.Add(time.Second)
	job := cltest.NewJob()
	job.StartAt = cltest.NullableTime(startAt)
	job.EndAt = cltest.NullableTime(endAt)
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"web", `{"type":"web"}`, false},
		{"ethlog", `{"type":"ethlog"}`, false},
		{"runlog", `{"type":"runlog"}`, false},
		{"runat", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, utils.ISO8601UTC(startAt)), false},
		{"runat w/o time", `{"type":"runat"}`, true},
		{"runat w time before start at", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, startAt.Add(-1*time.Second).Unix()), true},
		{"runat w time after end at", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, endAt.Add(time.Second).Unix()), true},
		{"cron", `{"type":"cron","params": {"schedule":"* * * * * *"}}`, false},
		{"cron w/o schedule", `{"type":"cron"}`, true},
		{"non-existent initiator", `{"type":"doesntExist"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var initr models.Initiator
			assert.NoError(t, json.Unmarshal([]byte(test.input), &initr))
			result := services.ValidateInitiator(initr, job)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}

func TestValidateServiceAgreement(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	_, err := store.KeyStore.NewAccount("password") // matches correct_password.txt
	assert.NoError(t, err)
	err = store.KeyStore.Unlock("password")
	assert.NoError(t, err)
	defer cleanup()

	account, err := store.KeyStore.GetAccount()
	assert.NoError(t, err)

	oracles := []string{account.Address.Hex()}

	basic := cltest.EasyJSONFromFixture("../internal/fixtures/web/hello_world_agreement.json")
	basic = basic.Add("oracles", oracles)

	tests := []struct {
		name      string
		input     cltest.EasyJSON
		wantError bool
	}{
		{"basic", basic, false},
		{"no payment", basic.Delete("payment"), true},
		{"less than minimum payment", basic.Add("payment", "1"), true},
		{"less than minimum expiration", basic.Add("expiration", 1), true},
		{"without being listed as an oracle", basic.Add("oracles", []string{}), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sa, err := cltest.ServiceAgreementFromString(test.input.String())
			assert.NoError(t, err)

			result := services.ValidateServiceAgreement(sa, store)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}

func TestValidateMinimumContractPayment_tasksExist(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig()
	defer cleanup()
	config.MinimumContractPayment = *assets.NewLink(1)

	s, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	bta := cltest.NewBridgeType("bridgeA")
	bta.MinimumContractPayment = *assets.NewLink(2)
	assert.NoError(t, s.Save(&bta))

	btb := cltest.NewBridgeType("bridgeB")
	btb.MinimumContractPayment = *assets.NewLink(3)
	assert.NoError(t, s.Save(&btb))

	job := models.NewJob()
	job.Tasks = []models.TaskSpec{
		cltest.NewTask("ethtx"),
		cltest.NewTask("bridgea"),
		cltest.NewTask("bridgeb"),
	}
	assert.NoError(t, s.Save(&job))

	tests := []struct {
		name      string
		amount    assets.Link
		wantValid bool
	}{
		{"less than", *assets.NewLink(5), false},
		{"equal", *assets.NewLink(6), true},
		{"greater than", *assets.NewLink(7), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid, err := services.ValidateMinimumContractPayment(s, job, test.amount)
			assert.NoError(t, err)
			assert.Equal(t, test.wantValid, valid)
		})
	}
}

func TestValidateMinimumContractPayment_taskDoesNotExist(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig()
	defer cleanup()
	config.MinimumContractPayment = *assets.NewLink(1)

	s, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	job := models.NewJob()
	job.Tasks = []models.TaskSpec{
		cltest.NewTask("idontexist"),
	}
	assert.NoError(t, s.Save(&job))

	valid, err := services.ValidateMinimumContractPayment(s, job, *assets.NewLink(1))
	assert.Error(t, err)
	assert.Equal(t, false, valid)
}
