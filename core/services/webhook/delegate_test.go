package webhook_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipelinemocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	webhookmocks "github.com/smartcontractkit/chainlink/core/services/webhook/mocks"
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
		meta        = pipeline.JSONSerializable{Val: "bar", Valid: true}
		vars        = map[string]interface{}{
			"jobSpec": map[string]interface{}{
				"databaseID":    spec.ID,
				"externalJobID": spec.ExternalJobID,
				"name":          spec.Name.ValueOrZero(),
			},
			"jobRun": map[string]interface{}{
				"requestBody": requestBody,
				"meta":        meta.Val,
			},
		}
		runner    = new(pipelinemocks.Runner)
		eiManager = new(webhookmocks.ExternalInitiatorManager)
		delegate  = webhook.NewDelegate(runner, eiManager, logger.TestLogger(t))
	)

	services, err := delegate.ServicesForSpec(*spec)
	require.NoError(t, err)
	require.Len(t, services, 1)
	service := services[0]

	// Should error before service is started
	_, err = delegate.WebhookJobRunner().RunJob(testutils.Context(t), spec.ExternalJobID, requestBody, meta)
	require.Error(t, err)
	require.Equal(t, webhook.ErrJobNotExists, errors.Cause(err))

	// Should succeed after service is started upon a successful run
	err = service.Start(testutils.Context(t))
	require.NoError(t, err)

	runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
		Return(false, nil).
		Run(func(args mock.Arguments) {
			run := args.Get(1).(*pipeline.Run)
			run.ID = int64(123)

			require.Equal(t, vars, run.Inputs.Val)
		}).Once()

	runID, err := delegate.WebhookJobRunner().RunJob(testutils.Context(t), spec.ExternalJobID, requestBody, meta)
	require.NoError(t, err)
	require.Equal(t, int64(123), runID)

	// Should error after service is started upon a failed run
	expectedErr := errors.New("foo bar")

	runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
		Return(false, expectedErr).Once()

	_, err = delegate.WebhookJobRunner().RunJob(testutils.Context(t), spec.ExternalJobID, requestBody, meta)
	require.Equal(t, expectedErr, errors.Cause(err))

	// Should error after service is stopped
	err = service.Close()
	require.NoError(t, err)

	_, err = delegate.WebhookJobRunner().RunJob(testutils.Context(t), spec.ExternalJobID, requestBody, meta)
	require.Equal(t, webhook.ErrJobNotExists, errors.Cause(err))
}
