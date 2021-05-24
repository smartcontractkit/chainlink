package webhook_test

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	webhookmocks "github.com/smartcontractkit/chainlink/core/services/webhook/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWebhookDelegate(t *testing.T) {
	var (
		spec = &job.Job{
			ID:            123,
			Type:          job.Webhook,
			Name:          null.StringFrom("sergtoshi stevemoto"),
			SchemaVersion: 1,
			ExternalJobID: uuid.NewV4(),
			WebhookSpec:   &job.WebhookSpec{},
			PipelineSpec:  &pipeline.Spec{},
		}

		requestBody = "foo"
		vars        = pipeline.NewVarsFrom(map[string]interface{}{
			"jobSpec": map[string]interface{}{
				"databaseID":    spec.ID,
				"externalJobID": spec.ExternalJobID,
				"name":          spec.Name.ValueOrZero(),
			},
			"jobRun": map[string]interface{}{
				"requestBody": requestBody,
				"meta":        map[string]interface{}{},
			},
		})
		meta      = pipeline.JSONSerializable{Val: "bar"}
		runner    = new(pipelinemocks.Runner)
		eiManager = new(webhookmocks.ExternalInitiatorManager)
		delegate  = webhook.NewDelegate(runner, eiManager)
	)

	services, err := delegate.ServicesForSpec(*spec)
	require.NoError(t, err)
	require.Len(t, services, 1)
	service := services[0]

	// Should error before service is started
	_, err = delegate.WebhookJobRunner().RunJob(context.Background(), spec.ExternalJobID, requestBody, meta)
	require.Error(t, err)
	require.Equal(t, webhook.ErrJobNotExists, errors.Cause(err))

	// Should succeed after service is started upon a successful run
	err = service.Start()
	require.NoError(t, err)

	runner.On("ExecuteAndInsertFinishedRun", mock.Anything, *spec.PipelineSpec, vars, meta, mock.Anything, true).
		Return(int64(123), pipeline.FinalResult{}, nil).Once()

	runID, err := delegate.WebhookJobRunner().RunJob(context.Background(), spec.ExternalJobID, requestBody, meta)
	require.NoError(t, err)
	require.Equal(t, int64(123), runID)

	// Should error after service is started upon a failed run
	expectedErr := errors.New("foo bar")

	runner.On("ExecuteAndInsertFinishedRun", mock.Anything, *spec.PipelineSpec, vars, meta, mock.Anything, true).
		Return(int64(0), pipeline.FinalResult{}, expectedErr).Once()

	_, err = delegate.WebhookJobRunner().RunJob(context.Background(), spec.ExternalJobID, requestBody, meta)
	require.Equal(t, expectedErr, errors.Cause(err))

	// Should error after service is stopped
	err = service.Close()
	require.NoError(t, err)

	_, err = delegate.WebhookJobRunner().RunJob(context.Background(), spec.ExternalJobID, requestBody, meta)
	require.Equal(t, webhook.ErrJobNotExists, errors.Cause(err))

	runner.AssertExpectations(t)
}
