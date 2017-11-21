package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var r *gin.Engine

func init() {
	r = Router()
}

func TestCreateAssignments(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()

	jsonStr := []byte(`{"version": "1.0.0"}`)
	resp, err := http.Post(server.URL+"/assignments", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should indicate internal server error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["Error saving to database."]}`, string(body), "Repsonse should return JSON")
}
