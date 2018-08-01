package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestServiceAgreementsController_Create(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	sa := cltest.FixtureCreateServiceAgreementViaWeb(t, app, "../internal/fixtures/web/hello_world_job.json")
	assert.NotEqual(t, "", sa.ID)
}
