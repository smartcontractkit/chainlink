package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
)

func TestNew_humanAI_discard(t *testing.T) {
	t.Parallel()
	var humanOut, humanErr strings.Builder
	p := New(false, &humanOut, &humanErr, SkipFD)
	require.False(t, p.AIOutput())
	require.False(t, p.LiveInlineProgress())

	p.HumanStdout("hello")
	require.Contains(t, humanOut.String(), "hello")

	pAI := New(true, &humanOut, &humanErr, SkipFD)
	pAI.HumanStdout("quiet")
	require.NotContains(t, humanOut.String(), "quiet")
}

func TestNew_liveInline_requiresTTYAndHuman(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	// Pipe / builder is not a TTY — human mode but no live inline.
	p := New(false, &b, &b, SkipFD)
	require.False(t, p.LiveInlineProgress())

	// AI mode never enables live inline even with fd 1 (often TTY in tests).
	pAI := New(true, &b, &b, 1)
	require.False(t, pAI.LiveInlineProgress())
}

func TestSparseStdoutln_onlyWhenAI(t *testing.T) {
	t.Parallel()
	var out, err strings.Builder
	p := New(false, &out, &err, SkipFD)
	p.SparseStdoutln("x")
	require.Empty(t, out.String())

	pAI := New(true, &out, &err, SkipFD)
	pAI.SparseStdoutln("path")
	require.Contains(t, out.String(), "path")
}

func TestStderrf_always(t *testing.T) {
	t.Parallel()
	var out, err strings.Builder
	p := New(true, &out, &err, SkipFD)
	p.Stderrf("err: %d\n", 1)
	require.Contains(t, err.String(), "err: 1")
}

func TestWarnWriter_discardWhenAI(t *testing.T) {
	t.Parallel()
	var out, err strings.Builder
	p := New(true, &out, &err, SkipFD)
	w := p.WarnWriter()
	n, e := w.Write([]byte("note\n"))
	require.NoError(t, e)
	require.Equal(t, 5, n)
	require.Empty(t, err.String())
}

func TestIfHuman(t *testing.T) {
	t.Parallel()
	var ran bool
	New(true, nil, nil, SkipFD).IfHuman(func() { ran = true })
	require.False(t, ran)
	New(false, nil, nil, SkipFD).IfHuman(func() { ran = true })
	require.True(t, ran)
}

func TestNewFromApp(t *testing.T) {
	t.Parallel()
	p := NewFromApp(&config.App{AIOutput: true})
	require.True(t, p.AIOutput())
}

func TestNew_nilWriters(t *testing.T) {
	t.Parallel()
	require.NotPanics(t, func() {
		p := New(false, nil, nil, SkipFD)
		p.HumanStderr("x")
		p.Stderrf("y")
	})
}
