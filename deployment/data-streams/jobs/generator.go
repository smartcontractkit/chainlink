package jobs

import (
	"fmt"

	"github.com/google/uuid"
)

// JobSpecGenerator knows how to generate job specs for each stream type.
// It doesn't cover bootstrap and LLO job specs.
type JobSpecGenerator interface {
	GenerateJobSpec(ssc StreamSpecConfig, externalJobID uuid.UUID) (*StreamJobSpec, error)
}

func GeneratorForStreamType(st StreamType) (JobSpecGenerator, error) {
	switch st {
	case StreamTypeQuote:
		return &QuoteStreamJobSpecGenerator{}, nil
	case StreamTypeMedian:
		return MedianStreamJobSpecGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported stream type: %s", st)
	}
}

type QuoteStreamJobSpecGenerator struct{}

func (QuoteStreamJobSpecGenerator) GenerateJobSpec(ssc StreamSpecConfig, externalJobID uuid.UUID) (spec *StreamJobSpec, err error) {
	if externalJobID == uuid.Nil {
		externalJobID = uuid.New()
	}
	spec = &StreamJobSpec{
		BaseJobSpec: BaseJobSpec{
			Name:          fmt.Sprintf("%s | %d", ssc.Name, ssc.StreamID),
			Type:          JobSpecTypeStream,
			SchemaVersion: 1,
			ExternalJobID: externalJobID,
		},
		StreamID: ssc.StreamID,
	}

	datasources := make([]Datasource, len(ssc.APIs))
	params := ssc.EARequestParams
	for i, api := range ssc.APIs {
		datasources[i] = Datasource{
			BridgeName: api,
			ReqData:    fmt.Sprintf(`"{\"data\":{\"endpoint\":\"%s\",\"from\":\"%s\",\"to\":\"%s\"}}"`, params.Endpoint, params.From, params.To),
		}
	}

	base := BaseObservationSource{
		Datasources:   datasources,
		AllowedFaults: len(datasources) - 1,
	}
	err = spec.SetObservationSource(base, ssc.ReportFields)

	return spec, err
}

type MedianStreamJobSpecGenerator struct{}

func (MedianStreamJobSpecGenerator) GenerateJobSpec(ssc StreamSpecConfig, externalJobID uuid.UUID) (spec *StreamJobSpec, err error) {
	// Quote and Median generate their job specs in the same way.
	return QuoteStreamJobSpecGenerator{}.GenerateJobSpec(ssc, externalJobID)
}
