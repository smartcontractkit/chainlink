package runner

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSurveyProgress_onTestJSONLine_packageTerminal(t *testing.T) {
	p := newSurveyProgress(2)

	require.False(t, p.onTestJSONLine([]byte(`not json`)))
	require.False(t, p.onTestJSONLine([]byte(`{"Action":"run","Package":"a/b","Test":"TestX"}`)))

	require.True(t, p.onTestJSONLine([]byte(`{"Action":"pass","Package":"a/b"}`)))
	c, tot, _ := p.snapshot()
	require.Equal(t, 1, c)
	require.Equal(t, 2, tot)

	// Duplicate package-level pass must not report a second completion tick.
	require.False(t, p.onTestJSONLine([]byte(`{"Action":"pass","Package":"a/b"}`)))
	c, _, _ = p.snapshot()
	require.Equal(t, 1, c)

	require.True(t, p.onTestJSONLine([]byte(`{"Action":"fail","Package":"c/d"}`)))
	c, _, _ = p.snapshot()
	require.Equal(t, 2, c)
}

func TestSurveyProgress_onTestJSONLine_skipFail(t *testing.T) {
	p := newSurveyProgress(1)
	require.True(t, p.onTestJSONLine([]byte(`{"Action":"skip","Package":"p"}`)))
	c, _, _ := p.snapshot()
	require.Equal(t, 1, c)

	p2 := newSurveyProgress(1)
	require.True(t, p2.onTestJSONLine([]byte(`{"Action":"fail","Package":"p"}`)))
	c2, _, _ := p2.snapshot()
	require.Equal(t, 1, c2)
}

func TestSurveyProgress_lastPkgUpdates(t *testing.T) {
	p := newSurveyProgress(10)
	p.onTestJSONLine([]byte(`{"Action":"run","Package":"x/y","Test":"TestZ"}`))
	_, _, last := p.snapshot()
	require.Equal(t, "x/y", last)
}

func TestEllipsizeRight(t *testing.T) {
	require.Equal(t, "short", ellipsizeRight("short", 10))
	require.Equal(t, "abcdefghij", ellipsizeRight("abcdefghij", 10))
	require.Equal(t, "…hij", ellipsizeRight("abcdefghij", 6))
}

func TestRenderSurveyProgressLine_smoke(t *testing.T) {
	var b strings.Builder
	p := newSurveyProgress(10)
	p.onTestJSONLine([]byte(`{"Action":"pass","Package":"demo/pkg"}`))
	renderSurveyProgressLine(&b, 1, 3, 2*time.Second, p, false)
	require.Contains(t, b.String(), "iter 1/3")
	require.Contains(t, b.String(), "1/10")
}
