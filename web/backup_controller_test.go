package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
)

func TestBackupController_Show(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/backup")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)
}
