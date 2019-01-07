package services_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			"zero initiators",
			cltest.LoadJSON("../internal/fixtures/web/zero_initiators.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one initiator only",
			cltest.LoadJSON("../internal/fixtures/web/initiator_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one task only",
			cltest.LoadJSON("../internal/fixtures/web/task_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
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

	bt := models.BridgeType{}
	bt.Name = models.MustNewTaskType("solargridreporting")
	bt.URL = cltest.WebURL("https://denergy.eth")
	assert.NoError(t, store.SaveBridgeType(&bt))

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

	account, err := store.KeyStore.GetFirstAccount()
	assert.NoError(t, err)

	oracles := []string{account.Address.Hex()}

	basic := cltest.EasyJSONFromFixture("../internal/fixtures/web/hello_world_agreement.json")
	basic = basic.Add("oracles", oracles)
	threeDays, _ := time.ParseDuration("72h")
	basic = basic.Add("endAt", time.Now().Add(threeDays))

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
		{"past allowed end at", basic.Add("endAt", "3000-06-19T22:17:19Z"), true},
		{"before allowed end at", basic.Add("endAt", "2018-06-19T22:17:19Z"), true},
		{"more than one initiator should fail",
			basic.Add("initiators",
				[]models.Initiator{
					{0, "", models.InitiatorServiceAgreementExecutionLog,
						models.InitiatorParams{}},
					{0, "", models.InitiatorWeb, models.InitiatorParams{}},
				}),
			true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sa, err := cltest.ServiceAgreementFromString(test.input.String())
			require.NoError(t, err)

			result := services.ValidateServiceAgreement(sa, store)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}
