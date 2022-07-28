package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxAttemptsController_Index_Success(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	borm := app.TxmORM()
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	_, from := cltest.MustInsertRandomKey(t, app.KeyStore.Eth(), 0)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 0, 1, from)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 1, 2, from)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 2, 3, from)

	resp, cleanup := client.Get("/v2/tx_attempts?size=2")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	var links jsonapi.Links
	var attempts []presenters.EthTxResource
	body := cltest.ParseResponseBody(t, resp)

	require.NoError(t, web.ParsePaginatedResponse(body, &attempts, &links))
	assert.NotEmpty(t, links["next"].Href)
	assert.Empty(t, links["prev"].Href)
	require.Len(t, attempts, 2)
	assert.Equal(t, "3", attempts[0].SentAt, "expected tx attempts order by sentAt descending")
	assert.Equal(t, "2", attempts[1].SentAt, "expected tx attempts order by sentAt descending")
}

func TestTxAttemptsController_Index_Error(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	resp, cleanup := client.Get("/v2/tx_attempts?size=TrainingDay")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, 422)
}
