package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/web"
	"github.com/stretchr/testify/assert"
)

func TestBox_MatchWildcardBoxPath(t *testing.T) {
	t.Parallel()

	boxList := []string{
		"index.html",
		"job_specs/_jobSpecId_/index.html",
		"job_specs/_jobSpecId_/runs/_jobSpecRunId_/routeInfo.json",
	}

	assert.Equal(
		t,
		"index.html",
		web.MatchWildcardBoxPath(boxList, "/", "index.html"),
	)
	assert.Equal(
		t,
		"",
		web.MatchWildcardBoxPath(boxList, "/", "not_found.html"),
	)

	assert.Equal(
		t,
		"job_specs/_jobSpecId_/runs/_jobSpecRunId_/routeInfo.json",
		web.MatchWildcardBoxPath(boxList, "/job_specs/abc123/runs/abc123", "routeInfo.json"),
	)
	assert.Equal(
		t,
		"",
		web.MatchWildcardBoxPath(boxList, "/job_specs/abc123/runs/abc123", "notFound.json"),
	)
}

func TestBox_MatchExactBoxPath(t *testing.T) {
	t.Parallel()

	boxList := []string{"main.js", "page/main.js"}

	assert.Equal(t, "main.js", web.MatchExactBoxPath(boxList, "/main.js"))
	assert.Equal(t, "", web.MatchExactBoxPath(boxList, "/not_found.js"))

	assert.Equal(t, "page/main.js", web.MatchExactBoxPath(boxList, "/page/main.js"))
	assert.Equal(t, "", web.MatchExactBoxPath(boxList, "/ppage/main.js"))
	assert.Equal(t, "", web.MatchExactBoxPath(boxList, "/page/not_found.js"))
}
