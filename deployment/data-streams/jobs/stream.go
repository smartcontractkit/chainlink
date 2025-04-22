package jobs

import (
	"errors"

	"github.com/pelletier/go-toml/v2"
)

type Datasource struct {
	BridgeName string
	ReqData    string
}

type ReportFieldLLO struct {
	ResultPath string
	// StreamID allows assigning an own stream ID to the report field.
	StreamID *string
}

type Pipeline interface {
	Render() (string, error)
}

type BaseObservationSource struct {
	Datasources   []Datasource
	AllowedFaults int
}

type ReportFields interface {
	GetStreamType() StreamType
}

type QuoteReportFields struct {
	Bid       ReportFieldLLO
	Benchmark ReportFieldLLO
	Ask       ReportFieldLLO
}

func (quote QuoteReportFields) GetStreamType() StreamType {
	return StreamTypeQuote
}

type MedianReportFields struct {
	Benchmark ReportFieldLLO
}

func (median MedianReportFields) GetStreamType() StreamType {
	return StreamTypeMedian
}

type StreamJobSpec struct {
	Base

	StreamID          uint32 `toml:"streamID"`
	ObservationSource string `toml:"observationSource,multiline,omitempty"`
}

func (s *StreamJobSpec) SetObservationSource(base BaseObservationSource, rf ReportFields) error {
	tmpl, err := templateForStreamType(rf.GetStreamType())
	if err != nil {
		return err
	}
	observationSourceData := struct {
		BaseObservationSource
		ReportFields
	}{
		base,
		rf,
	}
	rendered, err := renderTemplate(tmpl, observationSourceData)
	if err != nil {
		return err
	}
	s.ObservationSource = rendered
	return nil
}

func templateForStreamType(st StreamType) (string, error) {
	switch st {
	case StreamTypeQuote:
		return "osrc_mercury_v1_quote.go.tmpl", nil
	case StreamTypeMedian:
		return "osrc_mercury_v1_median.go.tmpl", nil
	default:
		return "", errors.New("unsupported stream type")
	}
}

func (s *StreamJobSpec) MarshalTOML() ([]byte, error) {
	return toml.Marshal(s)
}
