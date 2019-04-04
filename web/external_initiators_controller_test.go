package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
)

func TestExternalInitiatorsController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	client := app.NewHTTPClient()

	resp, cleanup := client.Post("/v2/external_initiators", nil)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 201)
}

func TestExternalInitiatorsController_Delete(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	client := app.NewHTTPClient()

	resp, cleanup := client.Delete("/v2/external_initiators")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 202)
}
