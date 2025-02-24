package jobs

import (
	"github.com/pelletier/go-toml/v2"
)

type Datasource struct {
	BridgeName string
	ReqData    string
}

type ReportFieldLLO struct {
	ResultPath string
}

type Pipeline interface {
	Render() (string, error)
}

type BaseObservationSource struct {
	Datasources   []Datasource
	AllowedFaults int
	Benchmark     ReportFieldLLO
}

type QuoteObservationSource struct {
	BaseObservationSource
	Bid ReportFieldLLO
	Ask ReportFieldLLO
}

type MedianObservationSource struct {
	BaseObservationSource
}

func renderObservationTemplate(fname string, obs any) (string, error) {
	return renderTemplate(fname, obs)
}

func (src QuoteObservationSource) Render() (string, error) {
	return renderObservationTemplate("osrc_mercury_v1_quote.go.tmpl", src)
}

func (src MedianObservationSource) Render() (string, error) {
	return renderObservationTemplate("osrc_mercury_v1_median.go.tmpl", src)
}

type StreamJobSpec struct {
	Base

	StreamID          string `toml:"streamID"`
	ObservationSource string `toml:"observationSource,multiline,omitempty"`
}

func (s *StreamJobSpec) SetObservationSource(obs Pipeline) error {
	rendered, err := obs.Render()
	if err != nil {
		return err
	}
	s.ObservationSource = rendered
	return nil
}

func (s *StreamJobSpec) MarshalTOML() ([]byte, error) {
	return toml.Marshal(s)
}
