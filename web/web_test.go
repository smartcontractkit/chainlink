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

	jsonStr := []byte(`{"subtasks":[{"adapterType": "httpJSON", "adapterParams": {"endpoint": "https://bitstamp.net/api/ticker/", "fields": ["last"]}}], "schedule": "* * * * *","version":"1.0.0"}`)
	resp, err := http.Post(server.URL+"/assignments", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"id":1}`, string(body), "Repsonse should return JSON")
}

func TestCreateInvalidAssignments(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()

	jsonStr := []byte(`{"subtasks":[{"adapterType": "ethereumBytes32", "adapterParams": {}}], "schedule": "* * * * *","version":"1.0.0"}`)
	resp, err := http.Post(server.URL+"/assignments", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["\"ethereumBytes32\" is not a supported adapter type."]}`, string(body), "Repsonse should return JSON")
}
