package models_test

import (
	"math/big"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v4"
)

func TestNewInitiatorFromRequest(t *testing.T) {
	t.Parallel()

	job := cltest.NewJob()
	tests := []struct {
		name     string
		initrReq models.InitiatorRequest
		jobSpec  models.JobSpec
		want     models.Initiator
	}{
		{
			name: models.InitiatorWeb,
			initrReq: models.InitiatorRequest{
				Type: models.InitiatorWeb,
			},
			jobSpec: job,
			want: models.Initiator{
				Type:      models.InitiatorWeb,
				JobSpecID: job.ID,
			},
		},
		{
			name: models.InitiatorWeb,
			initrReq: models.InitiatorRequest{
				Type: models.InitiatorFluxMonitor,
				InitiatorParams: models.InitiatorParams{
					IdleTimer: models.IdleTimerConfig{
						Duration: models.MustMakeDuration(5 * time.Second),
					},
					PollTimer: models.PollTimerConfig{
						Period: models.MustMakeDuration(1 * time.Minute),
					},
					Precision:         2,
					Threshold:         5,
					AbsoluteThreshold: 0.01,
				},
			},
			jobSpec: job,
			want: models.Initiator{
				Type:      models.InitiatorFluxMonitor,
				JobSpecID: job.ID,
				InitiatorParams: models.InitiatorParams{
					IdleTimer: models.IdleTimerConfig{
						Duration: models.MustMakeDuration(5 * time.Second),
					},
					PollTimer: models.PollTimerConfig{
						Period: models.MustMakeDuration(1 * time.Minute),
					},
					Precision:         2,
					Threshold:         5,
					AbsoluteThreshold: 0.01,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := models.NewInitiatorFromRequest(
				test.initrReq,
				test.jobSpec,
			)
			assert.Equal(t, test.want, res)
		})
	}
}

func TestInitiatorParams(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	schedule := models.Cron("* * * * *")
	address := common.HexToAddress("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")
	requesters := models.AddressCollection([]common.Address{address})
	big := utils.NewBig(big.NewInt(42))
	topics := models.Topics([][]common.Hash{
		models.TopicsForInitiatorsWhichRequireJobSpecIDTopic[models.InitiatorRunLog],
	})
	json := cltest.JSONFromString(t, `{"foo":42}`)
	duration := models.MustMakeDuration(42 * time.Second)
	time := models.NewAnyTime(time.Now())

	j := cltest.NewJob()
	require.NoError(t, store.CreateJob(&j))

	i := models.Initiator{
		JobSpecID: j.ID,
		InitiatorParams: models.InitiatorParams{
			Schedule:          schedule,
			Time:              time,
			Ran:               true,
			Address:           address,
			Requesters:        requesters,
			Name:              "foo",
			Body:              &json,
			FromBlock:         big,
			ToBlock:           big,
			Topics:            topics,
			JobIDTopicFilter:  models.JobID(uuid.NewV4()),
			RequestData:       json,
			Feeds:             json,
			Threshold:         42.42,
			AbsoluteThreshold: 54, // https://spooniom.com/6-9-42/
			Precision:         42,
			IdleTimer: models.IdleTimerConfig{
				Duration: duration,
			},
			PollTimer: models.PollTimerConfig{
				Period: duration,
			},
		},
	}

	require.NoError(t, store.CreateInitiator(&i))

	saved, err := store.FindInitiator(i.ID)
	require.NoError(t, err)

	assert.Equal(t, schedule, saved.InitiatorParams.Schedule)
	assert.Equal(t, time.Unix(), saved.InitiatorParams.Time.Unix())
	assert.Equal(t, true, saved.InitiatorParams.Ran)
	assert.Equal(t, address, saved.InitiatorParams.Address)
	assert.Equal(t, requesters, saved.InitiatorParams.Requesters)
	assert.Equal(t, "foo", saved.InitiatorParams.Name)
	assert.JSONEq(t, json.String(), saved.InitiatorParams.Body.String())
	assert.Equal(t, big, saved.InitiatorParams.FromBlock)
	assert.Equal(t, big, saved.InitiatorParams.ToBlock)
	assert.Equal(t, topics, saved.InitiatorParams.Topics)
	assert.Equal(t, json, saved.InitiatorParams.RequestData)
	assert.Equal(t, duration, saved.InitiatorParams.IdleTimer.Duration)
	assert.Equal(t, json, saved.InitiatorParams.Feeds)
	assert.Equal(t, float32(42.42), saved.InitiatorParams.Threshold)
	assert.Equal(t, i.InitiatorParams.AbsoluteThreshold,
		saved.InitiatorParams.AbsoluteThreshold)
	assert.Equal(t, int32(42), saved.InitiatorParams.Precision)
	assert.Equal(t, duration, saved.InitiatorParams.PollTimer.Period)

}

func TestNewJobFromRequest(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("CRON_TZ=UTC * * * * 6")
	require.NoError(t, store.CreateJob(&j1))

	jsr := models.JobSpecRequest{
		Initiators: cltest.BuildInitiatorRequests(t, j1.Initiators),
		Tasks:      cltest.BuildTaskRequests(t, j1.Tasks),
		StartAt:    j1.StartAt,
		EndAt:      j1.EndAt,
		MinPayment: assets.NewLink(5),
	}

	j2 := models.NewJobFromRequest(jsr)
	require.NoError(t, store.CreateJob(&j2))

	fetched1, err := store.FindJobSpec(j1.ID)
	assert.NoError(t, err)
	assert.Len(t, fetched1.Initiators, 1)
	assert.Len(t, fetched1.Tasks, 1)
	assert.Nil(t, fetched1.MinPayment)

	fetched2, err := store.FindJobSpec(j2.ID)
	assert.NoError(t, err)
	assert.Len(t, fetched2.Initiators, 1)
	assert.Len(t, fetched2.Tasks, 1)
	assert.Equal(t, assets.NewLink(5), fetched2.MinPayment)
}

func TestJobSpec_Save(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	befCreation := time.Now()
	j1 := cltest.NewJobWithSchedule("* * * * 7")
	aftCreation := time.Now()
	assert.True(t, true, j1.CreatedAt.After(aftCreation), j1.CreatedAt.Before(befCreation))
	assert.False(t, false, j1.CreatedAt.IsZero())

	befInsertion := time.Now()
	assert.NoError(t, store.CreateJob(&j1))
	aftInsertion := time.Now()
	assert.True(t, true, j1.CreatedAt.After(aftInsertion), j1.CreatedAt.Before(befInsertion))

	initr := j1.Initiators[0]

	j2, err := store.FindJobSpec(j1.ID)
	require.NoError(t, err)
	require.Len(t, j2.Initiators, 1)
	assert.Equal(t, initr.Schedule, j2.Initiators[0].Schedule)
}

func TestJobSpec_NewRun(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithSchedule("1 * * * *")
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp", `{"a":1}`)}

	run := cltest.NewJobRun(job)

	assert.Equal(t, job.ID, run.JobSpecID)
	assert.Equal(t, 1, len(run.TaskRuns))

	taskRun := run.TaskRuns[0]
	assert.Equal(t, "noop", taskRun.TaskSpec.Type.String())
	adapter, _ := adapters.For(taskRun.TaskSpec, store.Config, store.ORM, nil)
	assert.NotNil(t, adapter)
	assert.JSONEq(t, `{"type":"NoOp","a":1}`, taskRun.TaskSpec.Params.String())

	assert.Equal(t, job.Initiators[0], run.Initiator)
}

func TestJobEnded(t *testing.T) {
	t.Parallel()

	endAt := cltest.ParseNullableTime(t, "3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		endAt   null.Time
		current time.Time
		want    bool
	}{
		{"no end at", null.TimeFromPtr(nil), endAt.Time, false},
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

	startAt := cltest.ParseNullableTime(t, "3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		startAt null.Time
		current time.Time
		want    bool
	}{
		{"no start at", null.TimeFromPtr(nil), startAt.Time, true},
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
