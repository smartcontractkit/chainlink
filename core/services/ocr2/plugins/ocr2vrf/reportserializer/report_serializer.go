package reportserializer

import (
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	types "github.com/smartcontractkit/ocr2vrf/types"
)

type reportSerializer struct {
	e ocr2vrf.EthereumReportSerializer
}

var _ types.ReportSerializer = (*reportSerializer)(nil)

// NewReportSerializer provides a serialization component for sending byte-encoded reports on-chain.
func NewReportSerializer(encryptionGroup kyber.Group) types.ReportSerializer {
	return &reportSerializer{
		e: ocr2vrf.EthereumReportSerializer{
			G: encryptionGroup,
		},
	}
}

// SerializeReport serializes an abstract report into abi-encoded bytes.
func (serializer *reportSerializer) SerializeReport(r types.AbstractReport) ([]byte, error) {

	packed, err := serializer.e.SerializeReport(r)

	if err != nil {
		return nil, errors.Wrap(err, "serialize report")
	}

	return packed, nil
}

// DeserializeReport deserializes a serialized byte array into a report.
func (serializer *reportSerializer) DeserializeReport(reportBytes []byte) (types.AbstractReport, error) {
	// Leaving unimplemented, as serialization here is not two-way. The object that is abi-encoded on-chain is a BeaconReport, not an AbstractReport.
	// So, the AbstractReport is first converted to a BeaconReport before the encoding. Converting an AbstractReport to a BeaconReport requires
	// the removal of some fields, so when it is converted back to a BeaconReport and then deserialized, those fields are missing.
	// Consequently, either the returned object from this function will be an abstract report
	// that has had some fields removed/zeroed, or the return type will be changed to a BeaconReport, which cannot be re-serialized.
	//
	// Also, the need for off-chain deserialization is not currently clear.
	panic("implement me")
}

// MaxReportLength gives the max length of a report to be transmitted on-chain.
func (serializer *reportSerializer) MaxReportLength() uint {
	return 150_000 // TODO: change this.
}

// ReportLength provides the expected report length of a report.
func (serializer *reportSerializer) ReportLength(a types.AbstractReport) uint {
	s, err := serializer.SerializeReport(a)
	if err != nil {
		return 0
	}
	return uint(len(s))
}
