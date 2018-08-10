package services_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"net/url"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
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
		{"error in job", cltest.LoadJSON("../internal/fixtures/web/invalid_endat_job.json"),
			errors.New(`job validation: startat cannot be before endat`)},
		{"error in runat initr", cltest.LoadJSON("../internal/fixtures/web/run_at_wo_time_job.json"),
			errors.New(`job validation: initiator validation: runat must have a time`)},
		{"error in task", cltest.LoadJSON("../internal/fixtures/web/nonexistent_task_job.json"),
			errors.New(`job validation: task validation: idonotexist is not a supported adapter type`)},
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
	u, err := url.Parse("https://denergy.eth")
	assert.NoError(t, err)
	tt.URL = models.WebURL{URL: u}
	assert.NoError(t, store.Save(&tt))

	tests := []struct {
		description string
		name        string
		want        error
	}{
		{"existing external adapter", "solargridreporting",
			errors.New("adapter validation: adapter solargridreporting exists")},
		{"existing core adapter", "ethtx",
			errors.New("adapter validation: adapter ethtx exists")},
		{"no adapter name", "",
			errors.New("adapter validation: no name specified")},
		{"invalid adapter name", "invalid/adapter",
			errors.New("adapter validation error: Task Type validation: name invalid/adapter contains invalid characters")},
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
		{"runat", fmt.Sprintf(`{"type":"runat","time":"%v"}`, utils.ISO8601UTC(startAt)), false},
		{"runat w/o time", `{"type":"runat"}`, true},
		{"runat w time before start at", fmt.Sprintf(`{"type":"runat","time":"%v"}`, startAt.Add(-1*time.Second).Unix()), true},
		{"runat w time after end at", fmt.Sprintf(`{"type":"runat","time":"%v"}`, endAt.Add(time.Second).Unix()), true},
		{"cron", `{"type":"cron","schedule":"* * * * * *"}`, false},
		{"cron w/o schedule", `{"type":"cron"}`, true},
		{"non-existent initiator", `{"type":"doesntExist"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var initr models.Initiator
			assert.NoError(t, json.Unmarshal([]byte(test.input), &initr))
			result := services.ValidateInitiator(initr, job)
			if test.wantError {
				assert.Error(t, result)
			} else {
				assert.NoError(t, result)
			}
		})
	}
}

func TestValidateServiceAgreement(t *testing.T) {
	t.Parallel()

	testConfig, cleanup := cltest.NewConfig()
	config := testConfig.Config
	defer cleanup()

	basic := cltest.EasyJSONFromFixture("../internal/fixtures/web/hello_world_agreement.json")
	tests := []struct {
		name      string
		input     cltest.EasyJSON
		wantError bool
	}{
		{"basic", basic, false},
		{"no payment", basic.Delete("payment"), true},
		{"less than minimum payment", basic.Add("payment", 1), true},
		{"less than minimum expiration", basic.Add("expiration", 1), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sa := cltest.ServiceAgreementFromString(test.input.String())
			result := services.ValidateServiceAgreement(sa, config)

			cltest.ErrorPresence(t, test.wantError, result)
		})
	}
}
