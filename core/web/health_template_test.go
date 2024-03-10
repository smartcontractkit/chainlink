package web

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

var (
	//go:embed testdata/health.html
	healthHTML string

	//go:embed testdata/health.txt
	healthTXT string
)

func checks() []presenters.Check {
	const passing, failing = HealthStatusPassing, HealthStatusFailing
	return []presenters.Check{
		{Name: "foo", Status: passing},
		{Name: "foo.bar", Status: failing, Output: "example error message"},
		{Name: "foo.bar.1", Status: passing},
		{Name: "foo.bar.1.A", Status: passing},
		{Name: "foo.bar.1.B", Status: passing},
		{Name: "foo.bar.2", Status: failing, Output: `error:
this is a multi-line error:
new line:
original error`},
		{Name: "foo.bar.2.A", Status: failing, Output: "failure!"},
		{Name: "foo.bar.2.B", Status: passing},
		{Name: "foo.baz", Status: passing},
	}
	//TODO truncated error
}

func Test_checkTree_WriteHTMLTo(t *testing.T) {
	ct := newCheckTree(checks())
	var b bytes.Buffer
	require.NoError(t, ct.WriteHTMLTo(&b))
	got := b.String()
	require.Equalf(t, healthHTML, got, "got: %s", got)
}

func Test_writeTextTo(t *testing.T) {
	var b bytes.Buffer
	require.NoError(t, writeTextTo(&b, checks()))
	got := b.String()
	require.Equalf(t, healthTXT, got, "got: %s", got)
}
