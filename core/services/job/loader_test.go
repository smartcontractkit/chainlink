package job_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
)

func TestParseJobSpecsFromTOML(t *testing.T) {
	t.Parallel()

	t.Run("parses valid job specs", func(t *testing.T) {
		tomlConfig := `
[[Jobs]]
type = "cron"
schedule = "CRON_TZ=UTC * * * * * *"
externalJobID = "b3d3c3e3-1234-5678-9abc-def012345678"
name = "test-cron"

[[Jobs]]
type = "directrequest"
schemaVersion = 1
name = "test-dr"
`
		jobSpecs, err := job.ParseJobSpecsFromTOML([]byte(tomlConfig))
		require.NoError(t, err)
		require.Len(t, jobSpecs, 2)

		assert.Equal(t, "cron", jobSpecs[0].Type)
		assert.Contains(t, jobSpecs[0].TOMLString, "schedule")

		assert.Equal(t, "directrequest", jobSpecs[1].Type)
		assert.Contains(t, jobSpecs[1].TOMLString, "schemaVersion")
	})

	t.Run("returns empty for no jobs", func(t *testing.T) {
		tomlConfig := `
[Database]
URL = "postgresql://localhost:5432/chainlink"
`
		jobSpecs, err := job.ParseJobSpecsFromTOML([]byte(tomlConfig))
		require.NoError(t, err)
		assert.Empty(t, jobSpecs)
	})

	t.Run("returns error for invalid TOML", func(t *testing.T) {
		tomlConfig := `
[[Jobs]]
type = "cron
invalid toml
`
		_, err := job.ParseJobSpecsFromTOML([]byte(tomlConfig))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse TOML")
	})

	t.Run("returns error for job missing type", func(t *testing.T) {
		tomlConfig := `
[[Jobs]]
schedule = "* * * * *"
`
		_, err := job.ParseJobSpecsFromTOML([]byte(tomlConfig))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'type' field")
	})
}

func TestLoadJobsFromConfig(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)

	t.Run("loads valid job successfully", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		externalJobID := uuid.New()

		jobSpec := job.JobSpecConfig{
			Type:       "cron",
			TOMLString: `type = "cron"`,
		}

		// Mock validator
		validators := map[string]job.ValidatorFunc{
			"cron": func(tomlString string) (job.Job, error) {
				return job.Job{
					Type:          job.Cron,
					ExternalJobID: externalJobID,
					Name:          mustNewNullString("test-job"),
				}, nil
			},
		}

		// Mock ORM calls
		mockORM.On("FindJobByExternalJobID", ctx, externalJobID).
			Return(job.Job{}, errors.New("not found"))
		mockORM.On("CreateJob", ctx, mustHaveExternalJobID(t, externalJobID)).
			Return(nil)

		err := job.LoadJobsFromConfig(ctx, []job.JobSpecConfig{jobSpec}, mockORM, validators, lggr)
		require.NoError(t, err)
		mockORM.AssertExpectations(t)
	})

	t.Run("skips existing job", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		externalJobID := uuid.New()

		jobSpec := job.JobSpecConfig{
			Type:       "cron",
			TOMLString: `type = "cron"`,
		}

		validators := map[string]job.ValidatorFunc{
			"cron": func(tomlString string) (job.Job, error) {
				return job.Job{
					Type:          job.Cron,
					ExternalJobID: externalJobID,
					Name:          mustNewNullString("test-job"),
				}, nil
			},
		}

		// Mock that job already exists
		existingJob := job.Job{
			ID:            123,
			ExternalJobID: externalJobID,
		}
		mockORM.On("FindJobByExternalJobID", ctx, externalJobID).
			Return(existingJob, nil)

		err := job.LoadJobsFromConfig(ctx, []job.JobSpecConfig{jobSpec}, mockORM, validators, lggr)
		require.NoError(t, err)
		mockORM.AssertExpectations(t)
		mockORM.AssertNotCalled(t, "CreateJob")
	})

	t.Run("continues on validator error", func(t *testing.T) {
		mockORM := mocks.NewORM(t)

		jobSpecs := []job.JobSpecConfig{
			{Type: "cron", TOMLString: `invalid`},
			{Type: "directrequest", TOMLString: `valid`},
		}

		externalJobID := uuid.New()
		validators := map[string]job.ValidatorFunc{
			"cron": func(tomlString string) (job.Job, error) {
				return job.Job{}, errors.New("validation failed")
			},
			"directrequest": func(tomlString string) (job.Job, error) {
				return job.Job{
					Type:          job.DirectRequest,
					ExternalJobID: externalJobID,
				}, nil
			},
		}

		mockORM.On("FindJobByExternalJobID", ctx, externalJobID).
			Return(job.Job{}, errors.New("not found"))
		mockORM.On("CreateJob", ctx, mustHaveExternalJobID(t, externalJobID)).
			Return(nil)

		// Should not error, should continue
		err := job.LoadJobsFromConfig(ctx, jobSpecs, mockORM, validators, lggr)
		require.NoError(t, err)
		mockORM.AssertExpectations(t)
	})

	t.Run("returns error for unregistered validator", func(t *testing.T) {
		mockORM := mocks.NewORM(t)

		jobSpec := job.JobSpecConfig{
			Type:       "unknown",
			TOMLString: `type = "unknown"`,
		}

		validators := map[string]job.ValidatorFunc{}

		err := job.LoadJobsFromConfig(ctx, []job.JobSpecConfig{jobSpec}, mockORM, validators, lggr)
		require.NoError(t, err) // Function continues on error
	})

	t.Run("generates ExternalJobID if missing", func(t *testing.T) {
		mockORM := mocks.NewORM(t)

		jobSpec := job.JobSpecConfig{
			Type:       "cron",
			TOMLString: `type = "cron"`,
		}

		validators := map[string]job.ValidatorFunc{
			"cron": func(tomlString string) (job.Job, error) {
				return job.Job{
					Type:          job.Cron,
					ExternalJobID: uuid.Nil, // Missing
				}, nil
			},
		}

		mockORM.On("FindJobByExternalJobID", ctx, mustNotBeNil(t)).
			Return(job.Job{}, errors.New("not found"))
		mockORM.On("CreateJob", ctx, mustNotHaveNilExternalJobID(t)).
			Return(nil)

		err := job.LoadJobsFromConfig(ctx, []job.JobSpecConfig{jobSpec}, mockORM, validators, lggr)
		require.NoError(t, err)
		mockORM.AssertExpectations(t)
	})

	t.Run("handles empty job list", func(t *testing.T) {
		mockORM := mocks.NewORM(t)
		validators := map[string]job.ValidatorFunc{}

		err := job.LoadJobsFromConfig(ctx, []job.JobSpecConfig{}, mockORM, validators, lggr)
		require.NoError(t, err)
	})
}

// Test helpers
func mustNewNullString(s string) null.String {
	return null.StringFrom(s)
}

// Mock helper that matches any job with a specific ExternalJobID
func mustHaveExternalJobID(t *testing.T, expectedID uuid.UUID) interface{} {
	return mock.MatchedBy(func(j *job.Job) bool {
		if j.ExternalJobID != expectedID {
			t.Logf("Expected ExternalJobID %s, got %s", expectedID, j.ExternalJobID)
			return false
		}
		return true
	})
}

// Mock helper that matches any job with a non-nil ExternalJobID
func mustNotHaveNilExternalJobID(t *testing.T) interface{} {
	return mock.MatchedBy(func(j *job.Job) bool {
		if j.ExternalJobID == uuid.Nil {
			t.Log("Expected non-nil ExternalJobID, got nil")
			return false
		}
		return true
	})
}

// Mock helper that matches any non-nil UUID
func mustNotBeNil(t *testing.T) interface{} {
	return mock.MatchedBy(func(id uuid.UUID) bool {
		if id == uuid.Nil {
			t.Log("Expected non-nil UUID, got nil")
			return false
		}
		return true
	})
}
