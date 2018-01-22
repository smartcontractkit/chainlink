package web_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTaskType(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/task_types",
		"application/json",
		bytes.NewBuffer(cltest.LoadJSON("../internal/fixtures/web/create_random_number_task_type.json")),
	)
	assert.Equal(t, 200, resp.StatusCode)
	ttID := cltest.JobJSONFromResponse(resp.Body).ID

	tt := &models.TaskType{}
	assert.Nil(t, app.Store.One("ID", ttID, tt))
	assert.Equal(t, ttID, tt.ID)
	assert.Equal(t, "randomNumber", tt.Name)
	assert.Equal(t, "https://example.smartcontract.com/randomNumber", tt.HandlerURL)
}
