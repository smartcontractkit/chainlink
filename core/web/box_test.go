package web_test

import (
	"os"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/stretchr/testify/assert"
)

func TestBox_MatchWildcardBoxPath(t *testing.T) {
	t.Parallel()

	rootIndex := "index.html"
	jobSpecRunShowRouteInfo := strings.Join(
		[]string{"job_specs", "_jobSpecId_", "runs", "_jobSpecRunId_", "routeInfo.json"},
		(string)(os.PathSeparator),
	)
	boxList := []string{rootIndex, jobSpecRunShowRouteInfo}

	assert.Equal(
		t,
		rootIndex,
		web.MatchWildcardBoxPath(boxList, "/", "index.html"),
	)
	assert.Equal(
		t,
		"",
		web.MatchWildcardBoxPath(boxList, "/", "not_found.html"),
	)

	assert.Equal(
		t,
		jobSpecRunShowRouteInfo,
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

	main := "main.js"
	pageMain := strings.Join([]string{"page", "main.js"}, (string)(os.PathSeparator))
	boxList := []string{main, pageMain}

	assert.Equal(t, main, web.MatchExactBoxPath(boxList, "/main.js"))
	assert.Equal(t, "", web.MatchExactBoxPath(boxList, "/not_found.js"))

	assert.Equal(t, pageMain, web.MatchExactBoxPath(boxList, "/page/main.js"))
	assert.Equal(t, "", web.MatchExactBoxPath(boxList, "/ppage/main.js"))
	assert.Equal(t, "", web.MatchExactBoxPath(boxList, "/page/not_found.js"))
}
