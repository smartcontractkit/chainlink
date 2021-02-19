package directrequest_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/log/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	broadcaster := new(mocks.Broadcaster)
	runner := new(mocks.PipelineRunner)
	_, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "event_broadcaster", true)
	defer cleanupDB()

	delegate := directrequest.NewDelegate(broadcaster, runner, orm.DB)

	t.Run("Spec without DirectRequestSpec", func(t *testing.T) {
		spec := job.SpecDB{}
		_, err := delegate.ServicesForSpec(spec)
		assert.Error(t, err, "expects a *job.DirectRequestSpec to be present")
	})

	t.Run("Spec with DirectRequestSpec", func(t *testing.T) {
		spec := job.SpecDB{DirectRequestSpec: &job.DirectRequestSpec{}}
		services, err := delegate.ServicesForSpec(spec)
		require.NoError(t, err)
		assert.Len(t, services, 1)
	})
}

func TestDelegate_ServicesListenerHandleLog(t *testing.T) {
	broadcaster := new(mocks.Broadcaster)
	runner := new(mocks.PipelineRunner)
	_, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "event_broadcaster", true)
	defer cleanupDB()

	delegate := directrequest.NewDelegate(broadcaster, runner, orm.DB)

	spec := job.SpecDB{DirectRequestSpec: &job.DirectRequestSpec{}}
	services, err := delegate.ServicesForSpec(spec)
	require.NoError(t, err)
	assert.Len(t, services, 1)
	service := services[0]

	t.Run("Log is an OracleRequest", func(t *testing.T) {
		var listener log.Listener
		broadcaster.On("Register", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			listener = args.Get(1).(log.Listener)
		})
		broadcaster.On("Close").Return(nil)

		log := new(mocks.Broadcast)
		defer log.AssertExpectations(t)

		log.On("WasAlreadyConsumed").Return(false)
		logOracleRequest := contracts.LogOracleRequest{}
		log.On("DecodeLog").Return(logOracleRequest)
		log.On("MarkConsumed").Return(logOracleRequest)

		err := service.Start()
		defer service.Close()
		require.NoError(t, err)

		listener.HandleLog(log, nil)
	})
}
