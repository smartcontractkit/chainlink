package por

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	lggr := logger.TestLogger(t)
	delegate := NewDelegate(lggr)

	// Test with valid POR spec
	t.Run("ValidPORSpec", func(t *testing.T) {
		jobSpec := job.Job{
			ID:   1,
			Type: job.POR,
			PORSpec: &job.PORSpec{
				// Add any required fields for PORSpec here
			},
		}

		services, err := delegate.ServicesForSpec(context.Background(), jobSpec)
		require.NoError(t, err)
		require.Len(t, services, 1)

		service := services[0]
		require.NotNil(t, service)

		// Test that service implements job.ServiceCtx
		ctx := context.Background()
		err = service.Start(ctx)
		assert.NoError(t, err)

		err = service.Close()
		assert.NoError(t, err)
	})

	// Test with missing POR spec
	t.Run("MissingPORSpec", func(t *testing.T) {
		jobSpec := job.Job{
			ID:      1,
			Type:    job.POR,
			PORSpec: nil, // Missing POR spec
		}

		services, err := delegate.ServicesForSpec(context.Background(), jobSpec)
		require.Error(t, err)
		require.Nil(t, services)
		assert.Contains(t, err.Error(), "POR Delegate expects a *job.PORSpec to be present")
	})
}

func TestDelegate_JobType(t *testing.T) {
	lggr := logger.TestLogger(t)
	delegate := NewDelegate(lggr)

	jobType := delegate.JobType()
	assert.Equal(t, job.POR, jobType)
}

func TestPORPluginService(t *testing.T) {
	lggr := logger.TestLogger(t)
	
	// Create a mock factory (this would be properly instantiated in real usage)
	var factory *por.PorReportingPluginFactory // This is nil for this test, but structure is correct
	
	spec := &job.PORSpec{}
	service := NewPORPluginService(factory, spec, lggr)

	require.NotNil(t, service)
	assert.Equal(t, spec, service.spec)

	// Test service lifecycle
	ctx := context.Background()
	err := service.Start(ctx)
	assert.NoError(t, err)

	err = service.Close()
	assert.NoError(t, err)
}
