package runner

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
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
	t0 := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	runStart := t0.Add(-time.Hour)
	renderDiagnoseProgressLine(&b, 1, 3, 2*time.Second, runStart, t0, true)
	got := b.String()
	require.Contains(t, got, "iter 1/3 (2s)")
	require.Contains(t, got, "1h0m0s")
	require.NotContains(t, got, "·")
	require.NotContains(t, got, "%")
	require.NotContains(t, got, "✅")
	require.NotContains(t, got, "⌛")
	require.NotContains(t, got, "█")
}

func TestRenderDiagnoseProgressLine_noRunWallWhenRunStartZero(t *testing.T) {
	var b strings.Builder
	t0 := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	renderDiagnoseProgressLine(&b, 1, 3, 2*time.Second, time.Time{}, t0, true)
	got := b.String()
	require.Contains(t, got, "iter 1/3 (2s)")
	require.NotContains(t, got, "1h0m0s")
}

func TestRenderDiagnoseProgressLine_notLiveInline(t *testing.T) {
	var b strings.Builder
	t0 := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	renderDiagnoseProgressLine(&b, 1, 3, 2*time.Second, t0, t0, false)
	require.Empty(t, b.String())
}

func TestRenderParallelDiagnoseProgressLine(t *testing.T) {
	var b strings.Builder
	t0 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	now := t0.Add(100 * time.Second)
	p := newParallelDiagnoseProgressAt(10, t0)
	p.bumpCompletedForTest(1)
	p.startAtForTest(1, t0.Add(-40*time.Second))
	p.startAtForTest(3, t0.Add(-10*time.Second))

	renderParallelDiagnoseProgressLine(&b, p, now, true)

	got := b.String()
	require.Contains(t, got, "done 1/10")
	require.Contains(t, got, "iter 2 (2m20s)")
	require.Contains(t, got, "iter 4 (1m50s)")
	require.NotContains(t, got, "active")
	require.NotContains(t, got, "·")
	require.NotContains(t, got, "core/")
}

// Simulates: worker goroutine redraws the next iteration's \r line before the
// receiver prints the previous iteration's table row (unbuffered channel
// completion order). Without ClearInline, Fprintln appends to the progress line.
func TestDiagnoseDigestAfterProgressNeedsClear_mergedWithoutClear(t *testing.T) {
	t.Parallel()
	var stderr strings.Builder
	out := output.NewForTest(false, io.Discard, &stderr, true)
	runStart := time.Date(2020, 1, 1, 11, 0, 0, 0, time.UTC)
	now := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	renderDiagnoseProgressLine(out.HumanStderrWriter(), 5, 20, 2*time.Second, runStart, now, true)
	out.HumanStderr("    4  fail")
	got := stderr.String()
	require.Contains(t, got, "iter 5/20")
	require.Contains(t, got, "    4  fail")
	idxRow := strings.Index(got, "    4")
	require.Positive(t, idxRow)
	require.NotContains(t, got[:idxRow], "\n", "digest must not start on a new line (glitch)")
}

func TestDiagnoseDigestAfterProgressNeedsClear_clearedBeforeHumanStderr(t *testing.T) {
	t.Parallel()
	var stderr strings.Builder
	out := output.NewForTest(false, io.Discard, &stderr, true)
	runStart := time.Date(2020, 1, 1, 11, 0, 0, 0, time.UTC)
	now := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	renderDiagnoseProgressLine(out.HumanStderrWriter(), 5, 20, 2*time.Second, runStart, now, true)
	out.ClearInline()
	out.HumanStderr("    4  fail")
	got := stderr.String()
	require.Contains(t, got, "\r\u001b[K    4  fail\n")
}
