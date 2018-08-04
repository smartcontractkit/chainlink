package web_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestServiceAgreementsController_Create(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	sa := cltest.FixtureCreateServiceAgreementViaWeb(t, app, "../internal/fixtures/web/hello_world_agreement.json")
	assert.NotEqual(t, "", sa.ID)
	cltest.FindJob(app.Store, sa.JobSpecID)
	assert.Equal(t, big.NewInt(1), sa.Encumbrance.Payment)
}
