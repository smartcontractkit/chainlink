package directrequest_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"gotest.tools/assert"
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
		assert.ErrorContains(t, err, "expects a *job.DirectRequestSpec to be present")
	})
}
