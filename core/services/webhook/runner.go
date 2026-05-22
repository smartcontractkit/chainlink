package webhook

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// JobRunner executes webhook jobs triggered by external initiators.
type JobRunner interface {
	RunJob(ctx context.Context, jobUUID uuid.UUID, requestBody string, meta jsonserializable.JSONSerializable) (int64, error)
}

// DeprecatedJobRunner rejects webhook job runs after the job type has been removed.
type DeprecatedJobRunner struct{}

var _ JobRunner = DeprecatedJobRunner{}

func (DeprecatedJobRunner) RunJob(context.Context, uuid.UUID, string, jsonserializable.JSONSerializable) (int64, error) {
	return 0, fmt.Errorf("cannot run job of type %q: %w", job.Webhook, job.ErrJobTypeRemoved)
}

// ErrJobNotExists is returned when no webhook job is registered for the given external job ID.
// Deprecated: webhook jobs are no longer executed; kept for API compatibility in run endpoints.
var ErrJobNotExists = errors.New("job does not exist")
