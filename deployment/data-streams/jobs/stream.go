package jobs

import (
	"errors"

	"github.com/pelletier/go-toml/v2"
)

type Datasource struct {
	BridgeName string
	ReqData    string
}

type BaseObservationSource struct {
	Datasources   []Datasource
	AllowedFaults int
}

// StreamSpecConfig defines the configuration for a data stream specification.
// It allows specifying a custom job spec generator for advanced use cases.
type StreamSpecConfig struct {
	StreamID   uint32
	Name       string
	StreamType StreamType
	// ReportFields should be QuoteReportFields, MedianReportFields, etc., based on the stream type.
	ReportFields    ReportFields
	EARequestParams EARequestParams
	APIs            []string

	// Generator allows us to specify a custom job spec generator. We might want to do that in case we need to modify
	// the way this particular job is generated, or we might want to provide a generator for a custom stream type.
	// If omitted, the default generator for this stream type will be used. If there is no such generator, an error
	// will be returned.
	Generator JobSpecGenerator
}

type EARequestParams struct {
	Endpoint string `json:"endpoint"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type StreamJobSpec struct {
	BaseJobSpec

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
