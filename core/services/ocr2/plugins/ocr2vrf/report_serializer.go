package ocr2vrf

import (
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	types "github.com/smartcontractkit/ocr2vrf/types"
)

type ReportSerializer struct {
	G kyber.Group
}

// Return the serialized byte representation of the report, as
// expected by the onchain machinery.
func (serializer *ReportSerializer) SerializeReport(r types.AbstractReport) ([]byte, error) {

	s := ocr2vrf.ReportSerializer{
		G: serializer.G,
	}
	packed, err := s.SerializeReport(r)

	if err != nil {
		return nil, errors.Wrap(err, "serialize report")
	}

	return packed, nil
}

func (serializer *ReportSerializer) DeserializeReport(reportBytes []byte) (types.BeaconReport, error) {
	s := ocr2vrf.ReportSerializer{
		G: serializer.G,
	}
	r, err := s.DeserializeReport(reportBytes)

	if err != nil {
		return types.BeaconReport{}, errors.Wrap(err, "deserialize report")
	}

	return r, nil
}

// Return the longest possible report which can be passed onchain
func (serializer *ReportSerializer) MaxReportLength() uint {
	return 150_000 // TODO: change this.
}

// Return the predicted length of the output from SerializeReport
func (serializer *ReportSerializer) ReportLength(a types.AbstractReport) uint {
	s, err := serializer.SerializeReport(a)
	if err != nil {
		return 0
	}
	return uint(len(s))
}
