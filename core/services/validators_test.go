package services_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
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
		{"base case", cltest.MustReadFile(t, "testdata/hello_world_job.json"), nil},
		{
			"error in job",
			cltest.MustReadFile(t, "testdata/invalid_endat_job.json"),
			models.NewJSONAPIErrorsWith("StartAt cannot be before EndAt"),
		},
		{
			"error in runat initr",
			cltest.MustReadFile(t, "testdata/run_at_wo_time_job.json"),
			models.NewJSONAPIErrorsWith("RunAt must have a time"),
		},
		{
			"error in task",
			cltest.MustReadFile(t, "testdata/nonexistent_task_job.json"),
			models.NewJSONAPIErrorsWith("idonotexist is not a supported adapter type"),
		},
		{
			"zero initiators",
			cltest.MustReadFile(t, "testdata/zero_initiators.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one initiator only",
			cltest.MustReadFile(t, "testdata/initiator_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one task only",
			cltest.MustReadFile(t, "testdata/task_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
	}

	store, cleanup := cltest.NewStore(t)
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

func TestValidateJob_DevRejectsSleepAdapter(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	sleepingJob := cltest.NewJobWithWebInitiator()
	sleepingJob.Tasks[0].Type = adapters.TaskTypeSleep

	store.Config.Set("CHAINLINK_DEV", true)
	assert.NoError(t, services.ValidateJob(sleepingJob, store))

	store.Config.Set("CHAINLINK_DEV", false)
	assert.Error(t, services.ValidateJob(sleepingJob, store))
}

func TestValidateAdapter(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Create a duplicate
	bt := models.BridgeType{}
	bt.Name = models.MustNewTaskType("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, store.CreateBridgeType(&bt))

	tests := []struct {
		description string
		name        string
		url         models.WebURL
		want        error
	}{
		{
			"existing external adapter",
			"solargridreporting",
			bt.URL,
			models.NewJSONAPIErrorsWith("Adapter solargridreporting already exists"),
		},
		{
			"existing core adapter",
			"ethtx",
			bt.URL,
			models.NewJSONAPIErrorsWith("Adapter ethtx already exists"),
		},
		{
			"no adapter name",
			"",
			bt.URL,
			models.NewJSONAPIErrorsWith("No name specified"),
		},
		{
			"invalid adapter name",
			"invalid/adapter",
			bt.URL,
			models.NewJSONAPIErrorsWith("Task Type validation: name invalid/adapter contains invalid characters"),
		},
		{
			"invalid adapter url",
			"adapterwithinvalidurl",
			cltest.WebURL(t, "//denergy"),
			models.NewJSONAPIErrorsWith("Invalid URL format"),
		},
		{"new external adapter", "gdaxprice", bt.URL, nil},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			btr := &models.BridgeTypeRequest{Name: models.TaskType(test.name), URL: test.url}
			result := services.ValidateBridgeType(btr, store)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateExternalInitiator(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	//  Add duplicate
	exi := models.ExternalInitiator{Name: "duplicate", URL: cltest.WebURL(t, "https://a.web.url")}
	assert.NoError(t, store.CreateExternalInitiator(&exi))

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"basic", `{"name":"bitcoin","url":"https://test.url"}`, false},
		{"basic w/ underscore", `{"name":"bit_coin","url":"https://test.url"}`, false},
		{"bad url", `{"name":"bitcoin","url":"//test.url"}`, true},
		{"duplicate name", `{"name":"duplicate","url":"https://test.url"}`, true},
		{"invalid name characters", `{"name":"<invalid>","url":"https://test.url"}`, true},
		{"missing name", `{"url":"https://test.url"}`, true},
		{"missing url", `{"name":"bitcoin"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var exr models.ExternalInitiatorRequest

			assert.NoError(t, json.Unmarshal([]byte(test.input), &exr))
			result := services.ValidateExternalInitiator(&exr, store)

			cltest.AssertError(t, test.wantError, result)
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
		{"external", `{"type":"external","params":{"name":"bitcoin"}}`, false},
		{"runlog", `{"type":"runlog"}`, false},
		{"runat", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, utils.ISO8601UTC(startAt)), false},
		{"runat w/o time", `{"type":"runat"}`, true},
		{"runat w time before start at", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, startAt.Add(-1*time.Second).Unix()), true},
		{"runat w time after end at", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, endAt.Add(time.Second).Unix()), true},
		{"cron", `{"type":"cron","params": {"schedule":"* * * * * *"}}`, false},
		{"cron w/o schedule", `{"type":"cron"}`, true},
		{"external w/o name", `{"type":"external"}`, true},
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

	store, cleanup := cltest.NewStore(t)
	_, err := store.KeyStore.NewAccount("password") // matches correct_password.txt
	assert.NoError(t, err)
	err = store.KeyStore.Unlock("password")
	assert.NoError(t, err)
	defer cleanup()

	account, err := store.KeyStore.GetFirstAccount()
	assert.NoError(t, err)

	oracles := []string{account.Address.Hex()}

	basic := string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
	basic = cltest.MustJSONSet(t, basic, "oracles", oracles)
	threeDays, _ := time.ParseDuration("72h")
	basic = cltest.MustJSONSet(t, basic, "endAt", time.Now().Add(threeDays))

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"basic", basic, false},
		{"no payment", cltest.MustJSONDel(t, basic, "payment"), true},
		{"less than minimum payment", cltest.MustJSONSet(t, basic, "payment", "1"), true},
		{"less than minimum expiration", cltest.MustJSONSet(t, basic, "expiration", 1), true},
		{"without being listed as an oracle", cltest.MustJSONSet(t, basic, "oracles", []string{}), true},
		{"past allowed end at", cltest.MustJSONSet(t, basic, "endAt", "3000-06-19T22:17:19Z"), true},
		{"before allowed end at", cltest.MustJSONSet(t, basic, "endAt", "2018-06-19T22:17:19Z"), true},
		{"more than one initiator should fail",
			cltest.MustJSONSet(t, basic, "initiators",
				[]models.Initiator{{
					JobSpecID:       models.NewID(),
					Type:            models.InitiatorServiceAgreementExecutionLog,
					InitiatorParams: models.InitiatorParams{},
				}, {
					JobSpecID:       models.NewID(),
					Type:            models.InitiatorWeb,
					InitiatorParams: models.InitiatorParams{},
				},
				}),
			true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sa, err := cltest.ServiceAgreementFromString(test.input)
			require.NoError(t, err)

			result := services.ValidateServiceAgreement(sa, store)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}
