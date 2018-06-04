package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestDeleteQueryController_DeleteQuery(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	params := cltest.LoadJSON("../internal/fixtures/web/delete_query_bridges.json")

	resp := cltest.BasicAuthPost(
		app.Server.URL+"/v2/delete_query",
		"application/json",
		bytes.NewBuffer(params),
	)

	cltest.AssertServerResponse(t, resp, 500)

	resp = cltest.BasicAuthPost(
		app.Server.URL+"/v2/delete_query",
		"application/json",
		bytes.NewBuffer([]byte(`{"collection": "tests"}`)),
	)

	for i := 0; i < 5; i++ {
		bt := models.BridgeType{Name: fmt.Sprintf("testbridge%v", i),
			URL:                  cltest.WebURL("http://www.example.com"),
			DefaultConfirmations: 0}
		err := app.AddAdapter(&bt)
		assert.NoError(t, err)
	}
	bt := models.BridgeType{Name: "testRemaining",
		URL: cltest.WebURL("http://www.example.com")}
	err := app.AddAdapter(&bt)
	assert.NoError(t, err)

	resp = cltest.BasicAuthPost(
		app.Server.URL+"/v2/delete_query",
		"application/json",
		bytes.NewBuffer(params),
	)
	cltest.AssertServerResponse(t, resp, 200)

	query := models.DeleteQueryParams{}
	query.Query.Re = json.RawMessage(`{ "name" : "^test" }`)
	found, err := app.Store.AdvancedBridgeSearch(query.Query)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(found))

}
