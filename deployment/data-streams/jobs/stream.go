package jobs

import (
	"bytes"
	"text/template"

	"github.com/pelletier/go-toml/v2"
)

type Datasource struct {
	BridgeName string
	ReqData    string
}

type ReportField struct {
	ResultPath string
}

type ObservationSource struct {
	Datasources   []Datasource
	AllowedFaults int
	Benchmark     ReportField
	Bid           ReportField
	Ask           ReportField
}

type LLOSpec struct {
	Base

	StreamID          string `toml:"streamID"`
	ObservationSource string `toml:"observationSource,multiline,omitempty"`
}

func (s *LLOSpec) SetObservationSource(obs ObservationSource) error {
	rendered, err := s.buildObservationSource(obs)
	if err != nil {
		return err
	}
	s.ObservationSource = rendered
	return nil
}

func (s *LLOSpec) buildObservationSource(obs ObservationSource) (string, error) {
	var buf bytes.Buffer
	if err := observationTmpl.Execute(&buf, obs); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *LLOSpec) MarshalTOML() ([]byte, error) {
	return toml.Marshal(s)
}

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
}

var pipelineTemplate = `{{range $i, $a := .Datasources}}
{{- $srcNum := inc $i -}}
// data source {{$srcNum}}
ds{{$srcNum}}_payload [type=bridge name="bridge-{{$a.BridgeName}}" timeout="50s" requestData={{$a.ReqData}}];

ds{{$srcNum}}_benchmark [type=jsonparse path="{{$.Benchmark.ResultPath}}"];
ds{{$srcNum}}_bid [type=jsonparse path="{{$.Bid.ResultPath}}"];
ds{{$srcNum}}_ask [type=jsonparse path="{{$.Ask.ResultPath}}"];
{{end -}}

{{range $i, $a := .Datasources}}
{{- $srcNum := inc $i -}}
ds{{$srcNum}}_payload -> ds{{$srcNum}}_benchmark -> benchmark_price;
{{end -}}
benchmark_price [type=median allowedFaults={{.AllowedFaults}} index=0];

{{range $i, $a := .Datasources}}
{{- $srcNum := inc $i -}}
ds{{$srcNum}}_payload -> ds{{$srcNum}}_bid -> bid_price;
{{end -}}
bid_price [type=median allowedFaults={{.AllowedFaults}} index=1];

{{range $i, $a := .Datasources}}
{{- $srcNum := inc $i -}}
ds{{$srcNum}}_payload -> ds{{$srcNum}}_ask -> ask_price;
{{end -}}
ask_price [type=median allowedFaults={{.AllowedFaults}} index=2];
`

var observationTmpl = template.Must(template.New("observationSource").
	Funcs(funcMap).
	Parse(pipelineTemplate),
)
