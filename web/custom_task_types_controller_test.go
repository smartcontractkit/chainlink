package web_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomTaskType(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/task_types",
		"application/json",
		bytes.NewBuffer(cltest.LoadJSON("../internal/fixtures/web/create_random_number_task_type.json")),
	)
	cltest.CheckStatusCode(t, resp, 200)
	ttID := cltest.JobJSONFromResponse(resp.Body).ID

	tt := &models.CustomTaskType{}
	assert.Nil(t, app.Store.One("ID", ttID, tt))
	assert.Equal(t, ttID, tt.ID)
	assert.Equal(t, "randomnumber", tt.Name)
	assert.Equal(t, "https://example.lvh.me/randomNumber", tt.URL.String())
}
