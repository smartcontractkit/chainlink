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
		app.Server.URL+"/v2/bridge_types",
		"application/json",
		bytes.NewBuffer(cltest.LoadJSON("../internal/fixtures/web/create_random_number_bridge_type.json")),
	)
	cltest.CheckStatusCode(t, resp, 200)
	btName := cltest.ParseCommonJSON(resp.Body).Name

	bt := &models.BridgeType{}
	assert.Nil(t, app.Store.One("Name", btName, bt))
	assert.Equal(t, "randomnumber", bt.Name)
	assert.Equal(t, "https://example.com/randomNumber", bt.URL.String())
}
