package web_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateBridgeType(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/task_types",
		"application/json",
		bytes.NewBuffer(cltest.LoadJSON("../internal/fixtures/web/create_random_number_task_type.json")),
	)
	cltest.CheckStatusCode(t, resp, 200)
	btID := cltest.JobJSONFromResponse(resp.Body).ID

	bt := &models.BridgeType{}
	assert.Nil(t, app.Store.One("ID", btID, bt))
	assert.Equal(t, btID, bt.ID)
	assert.Equal(t, "randomnumber", bt.Name)
	assert.Equal(t, "https://example.com/randomNumber", bt.URL.String())
}
