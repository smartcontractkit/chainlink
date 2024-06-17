package integration_tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func mockConsensus(t *testing.T, workflowKeys []ocr2key.KeyBundle) capabilities.ConsensusCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"offchain_reporting@1.0.0",
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			obs := req.Inputs.Underlying["observations"]
			report := obs.(*values.List)
			//	rm := map[string]any{
			//		"report": report.Underlying[0],
			//	}

			//rm := map[string]any{
			//	"report": report.Underlying[0],
			//}

			inputs := report.Underlying[0]

			triggerEvent := capabilities.TriggerEvent{}
			if err := inputs.UnwrapTo(&triggerEvent); err != nil {
				return capabilities.CapabilityResponse{}, err
			}

			mercuryReports := []datastreams.FeedReport{}
			err := triggerEvent.Payload.UnwrapTo(&mercuryReports)
			if err != nil {
				return capabilities.CapabilityResponse{}, fmt.Errorf("failed to unwrap mercury reports: %w", err)
			}

			middleReport := mercuryReports[1]

			//	reportHash := computeReportHash(middleReport.FullReport, middleReport.ReportContext)

			//fullReport, err := reporttypes.Decode(middleReport.FullReport)
			//if err != nil {
			//	return capabilities.CapabilityResponse{}, fmt.Errorf("failed to decode full report: %w", err)
			//}

			//fmt.Printf("report: %v\n", mercuryReports)

			//bytes32 completeHash = keccak256(abi.encodePacked(keccak256(rawReport), reportContext));

			reportCtx := ocrTypes.ReportContext{}
			rawCtx := RawReportContext(reportCtx)

			var signatures [][]byte
			for _, key := range workflowKeys {
				sig, err := key.Sign(reportCtx, middleReport.FullReport)
				require.NoError(t, err)

				signatures = append(signatures, sig)
			}

			signedReport := types.SignedReport{
				Report:     middleReport.FullReport,
				Context:    rawCtx,
				Signatures: signatures,
				ID:         []byte("01"),
			}

			rm := map[string]any{
				"report": signedReport,
			}

			rv, err := values.NewMap(rm)
			if err != nil {
				return capabilities.CapabilityResponse{}, err
			}

			//	so the consensus takes the reports and what, resigns existing report or creates new report and resigns it?

			//	the forwarder needs to be configured with the signers of the workflow don, assumably the don f number is that of the
			//	workflow don

			return capabilities.CapabilityResponse{
				Value: rv,
			}, nil
		},
	)
}
func RawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}

/*
func computeReportHash(rawReport, reportContext []byte) []byte {
	crypto.Keccak256()
	reportHash := crypto.Keccak256(rawReport)
	return crypto.Keccak256(reportHash, reportContext)
}*/

type testMercuryCodec struct {
}

func (c testMercuryCodec) UnwrapValid(wrapped values.Value, _ [][]byte, _ int) ([]datastreams.FeedReport, error) {
	dest := []datastreams.FeedReport{}
	err := wrapped.UnwrapTo(&dest)
	return dest, err
}

func (c testMercuryCodec) Wrap(reports []datastreams.FeedReport) (values.Value, error) {
	return values.Wrap(reports)
}
