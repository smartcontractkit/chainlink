package workflows

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

// WFYamlSpec generates a validate yaml spec for a workflow for the given name and owner
func WFYamlSpec(t *testing.T, name, owner string) string {
	t.Helper()
	// we use a template to generate the yaml spec so that we can omitting the name and owner fields as needed
	var wfTmpl = `
{{if .Name}}
name: {{.Name}}
{{end}}
{{if .Owner}}
owner: {{.Owner}}
{{end}}
triggers:
- id: trigger_test@1.0.0
  config: {}

consensus:
  - id: offchain_reporting@1.0.0
    ref: offchain_reporting_1
    config: {}

targets:
  - id: write_polygon_mainnet@1.0.0
    ref: write_polygon_mainnet_1 
    config: {}
`
	type cfg struct {
		Name, Owner string
	}
	c := cfg{Name: name, Owner: owner}

	tm := template.Must(template.New("yaml").Parse(wfTmpl))

	buf := new(bytes.Buffer)
	require.NoError(t, tm.Execute(buf, c))
	return buf.String()
}
