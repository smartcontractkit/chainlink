package runner

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDiagnoseProgress_onTestJSONLine_packageTerminal(t *testing.T) {
	p := newDiagnoseProgress(2)

	require.False(t, p.onTestJSONLine([]byte(`not json`)))
	require.False(t, p.onTestJSONLine([]byte(`{"Action":"run","Package":"a/b","Test":"TestX"}`)))

	require.True(t, p.onTestJSONLine([]byte(`{"Action":"pass","Package":"a/b"}`)))
	c, tot, _, _ := p.snapshot()
	require.Equal(t, 1, c)
	require.Equal(t, 2, tot)

	// Duplicate package-level pass must not report a second completion tick.
	require.False(t, p.onTestJSONLine([]byte(`{"Action":"pass","Package":"a/b"}`)))
	c, _, _, _ = p.snapshot()
	require.Equal(t, 1, c)

	require.True(t, p.onTestJSONLine([]byte(`{"Action":"fail","Package":"c/d"}`)))
	c, _, _, _ = p.snapshot()
	require.Equal(t, 2, c)
}

func TestDiagnoseProgress_onTestJSONLine_skipFail(t *testing.T) {
	p := newDiagnoseProgress(1)
	require.True(t, p.onTestJSONLine([]byte(`{"Action":"skip","Package":"p"}`)))
	c, _, _, _ := p.snapshot()
	require.Equal(t, 1, c)

	p2 := newDiagnoseProgress(1)
	require.True(t, p2.onTestJSONLine([]byte(`{"Action":"fail","Package":"p"}`)))
	c2, _, _, _ := p2.snapshot()
	require.Equal(t, 1, c2)
}

func TestDiagnoseProgress_lastPkgUpdates(t *testing.T) {
	p := newDiagnoseProgress(10)
	p.onTestJSONLine([]byte(`{"Action":"run","Package":"x/y","Test":"TestZ"}`))
	_, _, last, _ := p.snapshot()
	require.Equal(t, "x/y", last)
}

func TestDiagnoseProgress_pkgOutcomeOnTerminal(t *testing.T) {
	p := newDiagnoseProgress(5)
	p.onTestJSONLine([]byte(`{"Action":"run","Package":"p/q","Test":"TestZ"}`))
	_, _, _, out := p.snapshot()
	require.Empty(t, out)
	p.onTestJSONLine([]byte(`{"Action":"pass","Package":"p/q"}`))
	_, _, last, out := p.snapshot()
	require.Equal(t, "p/q", last)
	require.Equal(t, "pass", out)
}

func TestShortenChainlinkImportPath(t *testing.T) {
	t.Parallel()
	require.Empty(t, shortenChainlinkImportPath(""))
	require.Equal(t, ".", shortenChainlinkImportPath(chainlinkModulePrefix))
	require.Equal(t, "core/foo", shortenChainlinkImportPath(chainlinkModulePrefix+"/core/foo"))
	require.Equal(t, "other.com/pkg", shortenChainlinkImportPath("other.com/pkg"))
}

func TestEllipsizeRight(t *testing.T) {
	require.Equal(t, "short", ellipsizeRight("short", 10))
	require.Equal(t, "abcdefghij", ellipsizeRight("abcdefghij", 10))
	require.Equal(t, "…hij", ellipsizeRight("abcdefghij", 6))
}

func TestRenderDiagnoseProgressLine_smoke(t *testing.T) {
	var b strings.Builder
	p := newDiagnoseProgress(10)
	p.onTestJSONLine([]byte(`{"Action":"pass","Package":"demo/pkg"}`))
	renderDiagnoseProgressLine(&b, 1, 3, 2*time.Second, p, true)
	require.Contains(t, b.String(), "iter 1/3")
	require.Contains(t, b.String(), "1/10 10%")
	require.Contains(t, b.String(), "✅")
	require.NotContains(t, b.String(), "█")
}

func TestRenderDiagnoseProgressLine_inProgressShowsHourglass(t *testing.T) {
	var b strings.Builder
	p := newDiagnoseProgress(10)
	p.onTestJSONLine([]byte(`{"Action":"run","Package":"demo/pkg","Test":"TestX"}`))
	renderDiagnoseProgressLine(&b, 1, 3, 2*time.Second, p, true)
	require.Contains(t, b.String(), "⌛")
	require.NotContains(t, b.String(), "✅")
}

func TestRenderDiagnoseProgressLine_notTTY(t *testing.T) {
	var b strings.Builder
	p := newDiagnoseProgress(10)
	p.onTestJSONLine([]byte(`{"Action":"pass","Package":"demo/pkg"}`))
	renderDiagnoseProgressLine(&b, 1, 3, 2*time.Second, p, false)
	require.Empty(t, b.String())
}
