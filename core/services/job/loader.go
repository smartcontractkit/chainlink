// Package job provides job management functionality including configuration-based job loading.
// This package supports loading job specifications from TOML configuration files at node startup,
// allowing jobs to be defined declaratively alongside other node configuration.
package job

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// JobSpecConfig represents a job specification from configuration file.
// It contains the job type and the raw TOML string for type-specific validation.
type JobSpecConfig struct {
	Type       string `toml:"type"`       // Job type (e.g., "cron", "directrequest", "webhook")
	TOMLString string `toml:"-"`          // Raw TOML representation for validation
}

// ValidatorFunc is a function that validates a TOML job spec string and returns a Job.
// This pattern breaks import cycles by allowing validators to be injected from higher-level packages.
// Example: cron.ValidatedCronSpec can be passed as the validator for "cron" type jobs.
type ValidatorFunc func(tomlString string) (Job, error)

// LoadJobsFromConfig loads job specifications from configuration and creates them in the database.
// Jobs are only created if they don't already exist (checked by ExternalJobID).
// If a job fails to load, an error is logged and processing continues with remaining jobs.
//
// Parameters:
//   - ctx: Context for database operations
//   - jobSpecs: Job specifications parsed from configuration
//   - orm: Database ORM for job operations
//   - validators: Map of job type to validation function (e.g., "cron" -> cron.ValidatedCronSpec)
//   - lggr: Logger for operation tracking
//
// Returns nil on success. Individual job failures are logged but don't fail the entire operation.
func LoadJobsFromConfig(ctx context.Context, jobSpecs []JobSpecConfig, orm ORM, validators map[string]ValidatorFunc, lggr logger.Logger) error {
	if len(jobSpecs) == 0 {
		lggr.Debug("No jobs specified in configuration")
		return nil
	}

	lggr.Infow("Loading jobs from configuration", "count", len(jobSpecs))

	successCount := 0
	failCount := 0

	for i, jobSpec := range jobSpecs {
		if err := loadSingleJob(ctx, jobSpec, orm, validators, lggr); err != nil {
			failCount++
			lggr.Errorw("Failed to load job from config",
				"index", i,
				"type", jobSpec.Type,
				"err", err,
			)
			// Continue with other jobs
			continue
		}
		successCount++
	}

	lggr.Infow("Completed loading jobs from configuration",
		"total", len(jobSpecs),
		"success", successCount,
		"failed", failCount,
	)

	return nil
}

func loadSingleJob(ctx context.Context, jobSpec JobSpecConfig, orm ORM, validators map[string]ValidatorFunc, lggr logger.Logger) error {
	// Get validator for this job type
	validator, ok := validators[jobSpec.Type]
	if !ok {
		return fmt.Errorf("no validator registered for job type: %s", jobSpec.Type)
	}

	// Parse job using type-specific validator
	job, err := validator(jobSpec.TOMLString)
	if err != nil {
		return fmt.Errorf("failed to validate job spec: %w", err)
	}

	// Check if job already exists by ExternalJobID
	if job.ExternalJobID == uuid.Nil {
		lggr.Warnw("Job spec missing ExternalJobID, generating one", "type", job.Type, "name", job.Name.ValueOrZero())
		job.ExternalJobID = uuid.New()
	}

	existing, err := orm.FindJobByExternalJobID(ctx, job.ExternalJobID)
	if err == nil && existing.ID != 0 {
		lggr.Infow("Job already exists, skipping creation",
			"externalJobID", job.ExternalJobID,
			"jobID", existing.ID,
			"name", existing.Name.ValueOrZero(),
		)
		return nil
	}

	// Create job in database
	if err := orm.CreateJob(ctx, &job); err != nil {
		return fmt.Errorf("failed to create job in database: %w", err)
	}

	lggr.Infow("Created job from configuration",
		"jobID", job.ID,
		"externalJobID", job.ExternalJobID,
		"name", job.Name.ValueOrZero(),
		"type", job.Type,
	)

	return nil
}

// ParseJobSpecsFromTOML parses a TOML configuration and extracts job specifications.
// This helper function is used to load jobs from the main node configuration file.
//
// The expected format is:
//
//	[[Jobs]]
//	type = "cron"
//	schedule = "CRON_TZ=UTC * * * * * *"
//	externalJobID = "b3d3c3e3-1234-5678-9abc-def012345678"
//
//	[[Jobs]]
//	type = "directrequest"
//	# ... other job fields
//
// Returns an empty slice if no [[Jobs]] section exists in the configuration.
func ParseJobSpecsFromTOML(tomlBytes []byte) ([]JobSpecConfig, error) {
	type ConfigWithJobs struct {
		Jobs []map[string]interface{} `toml:"Jobs"`
	}

	var config ConfigWithJobs
	if err := toml.Unmarshal(tomlBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	if len(config.Jobs) == 0 {
		return nil, nil
	}

	jobSpecs := make([]JobSpecConfig, 0, len(config.Jobs))
	for i, jobMap := range config.Jobs {
		// Extract type
		jobType, ok := jobMap["type"].(string)
		if !ok {
			return nil, fmt.Errorf("job %d missing 'type' field", i)
		}

		// Convert map back to TOML string for validation
		tomlBytes, err := toml.Marshal(jobMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal job %d: %w", i, err)
		}

		jobSpecs = append(jobSpecs, JobSpecConfig{
			Type:       jobType,
			TOMLString: string(tomlBytes),
		})
	}

	return jobSpecs, nil
}
