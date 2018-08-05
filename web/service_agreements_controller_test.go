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
	js := cltest.FindJob(app.Store, sa.JobSpecID)
	assert.Equal(t, "0x85820c5ec619a1f517ee6cfeff545ec0ca1a90206e1a38c47f016d4137e801dd", js.Digest)

	assert.Equal(t, big.NewInt(1), sa.Encumbrance.Payment)
	assert.Equal(t, big.NewInt(2), sa.Encumbrance.Expiration)
	assert.Equal(t, "0x220072871b41155e7e1a6c45246a6d18a8a25350917d2c6c6c49d5d79a6af5bf", sa.ID)
}
