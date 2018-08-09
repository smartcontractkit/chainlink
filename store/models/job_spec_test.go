package models_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestJobSpec_Save(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1, initr := cltest.NewJobWithSchedule("* * * * 7")
	assert.NoError(t, store.SaveJob(&j1))

	store.Save(j1)
	j2, err := store.FindJob(j1.ID)
	assert.NoError(t, err)
	assert.Equal(t, initr.Schedule, j2.Initiators[0].Schedule)
}

func TestJobSpec_NewRun(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, initr := cltest.NewJobWithSchedule("1 * * * *")
	job.Tasks = []models.TaskSpec{cltest.NewTask("NoOp", `{"a":1}`)}

	run := job.NewRun(initr)

	assert.Equal(t, job.ID, run.JobID)
	assert.Equal(t, 1, len(run.TaskRuns))

	taskRun := run.TaskRuns[0]
	assert.Equal(t, "noop", taskRun.Task.Type.String())
	adapter, _ := adapters.For(taskRun.Task, store)
	assert.NotNil(t, adapter)
	assert.JSONEq(t, `{"type":"NoOp","a":1}`, taskRun.Task.Params.String())

	assert.Equal(t, initr, run.Initiator)
}

func TestJobEnded(t *testing.T) {
	t.Parallel()

	endAt := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		endAt   null.Time
		current time.Time
		want    bool
	}{
		{"no end at", null.Time{Valid: false}, endAt.Time, false},
		{"before end at", endAt, endAt.Time.Add(-time.Nanosecond), false},
		{"at end at", endAt, endAt.Time, false},
		{"after end at", endAt, endAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.EndAt = test.endAt

			assert.Equal(t, test.want, job.Ended(test.current))
		})
	}
}

func TestJobSpec_Started(t *testing.T) {
	t.Parallel()

	startAt := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		startAt null.Time
		current time.Time
		want    bool
	}{
		{"no start at", null.Time{Valid: false}, startAt.Time, true},
		{"before start at", startAt, startAt.Time.Add(-time.Nanosecond), false},
		{"at start at", startAt, startAt.Time, true},
		{"after start at", startAt, startAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.StartAt = test.startAt

			assert.Equal(t, test.want, job.Started(test.current))
		})
	}
}

func TestTaskSpec_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	tests := []struct {
		name          string
		taskType      string
		confirmations uint64
		json          string
		output        string
	}{
		{"noop", "noop", 0, `{"type":"noOp"}`, `{"type":"noop","confirmations":0}`},
		{
			"httpget",
			"httpget",
			0,
			`{"type":"httpget","url":"http://www.no.com"}`,
			`{"type":"httpget","url":"http://www.no.com","confirmations":0}`,
		},
		{
			"with confirmations",
			"noop",
			10,
			`{"type":"noop","confirmations":10}`,
			`{"type":"noop","confirmations":10}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var task models.TaskSpec
			err := json.Unmarshal([]byte(test.json), &task)
			assert.NoError(t, err)
			assert.Equal(t, test.confirmations, task.Confirmations)

			assert.Equal(t, test.taskType, task.Type.String())
			_, err = adapters.For(task, store)
			assert.NoError(t, err)

			s, err := json.Marshal(task)
			assert.NoError(t, err)
			assert.JSONEq(t, test.output, string(s))
		})
	}
}

func TestNormalizeSpecJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"identitiy", `{"endAt":"2018-07-02T21:51:36Z"}`, `{"endAt":"2018-07-02T21:51:36Z","initiators":null,"startAt":null,"tasks":null}`},
		{"ordering of keys", `{"startAt":"2018-07-02T21:51:36Z","endAt":"2018-07-02T21:51:36Z"}`, `{"endAt":"2018-07-02T21:51:36Z","initiators":null,"startAt":"2018-07-02T21:51:36Z","tasks":null}`},
		{"task type normalization", `{"tasks":[{"type":"noOp"}]}`, `{"endAt":null,"initiators":null,"startAt":null,"tasks":[{"confirmations":0.000000e+00,"type":"noop"}]}`},
		{"initiator type normalization", `{"initiators":[{"type":"runLog"}]}`, `{"endAt":null,"initiators":[{"address":"0x0000000000000000000000000000000000000000","id":0.000000e+00,"jobId":"","time":"0001-01-01T00:00:00Z","type":"runlog"}],"startAt":null,"tasks":null}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := models.NormalizeSpecJSON(test.input)

			assert.Equal(t, test.want, result)
			assert.NoError(t, err)
		})
	}
}

func TestNewTaskType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		errored bool
	}{
		{"basic", "NoOp", "noop", false},
		{"special characters", "-_-", "-_-", false},
		{"invalid character", "NoOp!", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := models.NewTaskType(test.input)

			if test.errored {
				assert.Error(t, err)
			} else {
				assert.Equal(t, models.TaskType(test.want), got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceAgreementRequest_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		wantDigest  string
		wantPayment int64
	}{
		{"basic",
			`{"payment":1,"initiators":[{"type":"web"}],"tasks":[{"type":"httpget","url":"https://bitstamp.net/api/ticker/"},{"type":"jsonparse","path":["last"]},{"type":"ethbytes32"},{"type":"ethtx"}]}`,
			"0x57bf5be3447b9a3f8491b6538b01f828bcfcaf2d685ea90375ed4ec2943f4865", 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreementRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			assert.Equal(t, test.wantDigest, sar.Digest)
			assert.Equal(t, big.NewInt(test.wantPayment), sar.Encumbrance.Payment)
		})
	}
}

func TestNewServiceAgreementFromRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		wantDigest  string
		wantPayment int64
	}{
		{"basic",
			`{"payment":1,"initiators":[{"type":"web"}],"tasks":[{"type":"httpget","url":"https://bitstamp.net/api/ticker/"},{"type":"jsonparse","path":["last"]},{"type":"ethbytes32"},{"type":"ethtx"}]}`,
			"0x57bf5be3447b9a3f8491b6538b01f828bcfcaf2d685ea90375ed4ec2943f4865", 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreementRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			sa, err := models.NewServiceAgreementFromRequest(sar)
			assert.NoError(t, err)
			assert.Equal(t, test.wantDigest, sa.ID)
			assert.Equal(t, big.NewInt(test.wantPayment), sa.Encumbrance.Payment)
		})
	}
}

func TestEncumbrance_ABI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		payment    *big.Int
		expiration *big.Int
		want       string
	}{
		{"basic", big.NewInt(1), big.NewInt(2), "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002"},
		{"empty", nil, nil, "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enc := models.Encumbrance{
				Payment:    test.payment,
				Expiration: test.expiration,
			}

			assert.Equal(t, test.want, enc.ABI())
		})
	}
}
