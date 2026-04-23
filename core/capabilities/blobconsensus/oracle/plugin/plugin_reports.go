package plugin

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	ocrtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

const ReportMetaDataPrependLength = 109

const InfoRequestID = "requestID"
const InfoConsensusFailureMessage = "failureMessage"
const InfoConsensusFailureCode = "failureCode"
const InfoKeyBundleName = "keyBundleName"

func (r *reportingPlugin) Reports(ctx context.Context, seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[[]byte], error) {
	requestsOutcome := &oracletypes.Outcome{}
	err := proto.Unmarshal(outcome, requestsOutcome)
	if err != nil {
		return nil, err
	}

	var reports []ocr3types.ReportPlus[[]byte]
	var successIDs, failureIDs []string

	// Create a report for each outcome
	for _, reqOutcome := range requestsOutcome.Outcomes {
		switch v := reqOutcome.GetOutcome().(type) {
		case *oracletypes.ConsensusOutcome_Success:
			successOutcome := v.Success
			r.lggr.Debugw("received successful consensus outcome", "seqNr", seqNr, "requestID", successOutcome.Metadata.RequestId)
			reqMetadata := successOutcome.Metadata
			var report []byte
			switch reqMetadata.RequestType {
			case oracletypes.RequestType_VALUE_CONSENSUS:
				report = successOutcome.Outcome
			case oracletypes.RequestType_REPORT_GENERATION:
				// If the request type is report extract the report from the values.Value before signing it
				serialisedValue := successOutcome.Outcome
				value := &valuespb.Value{}
				if err := proto.Unmarshal(serialisedValue, value); err != nil {
					return nil, fmt.Errorf("failed to unmarshal value for request %s: %w", reqMetadata.RequestId, err)
				}

				report = value.GetBytesValue()
				if report == nil {
					return nil, fmt.Errorf("failed to get report bytes for request %s", reqMetadata.RequestId)
				}
			}

			meta := ocrtypes.Metadata{
				Version:          1,
				ExecutionID:      reqMetadata.WorkflowExecutionId,
				Timestamp:        uint32(successOutcome.Timestamp.AsTime().Unix()), // nolint
				DONID:            reqMetadata.WorkflowDonId,
				DONConfigVersion: reqMetadata.WorkflowDonConfigVersion,
				WorkflowID:       reqMetadata.WorkflowId,
				WorkflowName:     reqMetadata.WorkflowName,
				WorkflowOwner:    reqMetadata.WorkflowOwner,
				ReportID:         reqMetadata.ReportId,
			}

			metadataPrepend, err := meta.Encode()
			if err != nil {
				return nil, fmt.Errorf("failed to encode metadata for request %s: %w", reqMetadata.RequestId, err)
			}

			reportWithMetaData := append(metadataPrepend, report...)

			// Check if the report is too large to transmit
			if len(reportWithMetaData) > r.maxReportLengthBytes {
				r.lggr.Errorw("report is too large to transmit", "seqNr", seqNr, "requestID", reqMetadata.RequestId,
					"reportSize", len(reportWithMetaData), "maxReportLengthBytes", r.maxReportLengthBytes)
				failureMsg := fmt.Sprintf(
					"report too large: the report for this request is %d bytes which exceeds the maximum allowed size of %d bytes; reduce the size of the data being returned",
					len(reportWithMetaData), r.maxReportLengthBytes)
				info, err := createFailedConsensusReportInfo(reqMetadata.RequestId, reqMetadata.KeyBundleId, failureMsg,
					oracletypes.ConsensusFailureCode_REPORT_TOO_LARGE)
				if err != nil {
					return nil, fmt.Errorf("failed to create report info for oversized report %s: %w", reqMetadata.RequestId, err)
				}
				reports = append(reports, ocr3types.ReportPlus[[]byte]{
					ReportWithInfo: ocr3types.ReportWithInfo[[]byte]{
						Report: []byte{},
						Info:   info,
					},
					TransmissionScheduleOverride: nil,
				})
				continue
			}

			info, err := createSuccessfulConsensusReportInfo(reqMetadata)
			if err != nil {
				return nil, fmt.Errorf("failed to create report info for successful consensus request %s: %w", reqMetadata.RequestId, err)
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: ocr3types.ReportWithInfo[[]byte]{
					Report: reportWithMetaData,
					Info:   info,
				},
				TransmissionScheduleOverride: nil,
			})
			successIDs = append(successIDs, reqMetadata.RequestId)
		case *oracletypes.ConsensusOutcome_Failure:
			failedOutcome := v.Failure
			r.lggr.Debugw("received failed consensus outcome", "seqNr", seqNr, "requestID", failedOutcome.RequestID)
			info, err := createFailedConsensusReportInfo(failedOutcome.RequestID, failedOutcome.KeyBundleId, failedOutcome.FailureMessage,
				failedOutcome.Code)
			if err != nil {
				return nil, fmt.Errorf("failed to create report info for failed consensus outcome %s: %w", failedOutcome.RequestID, err)
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: ocr3types.ReportWithInfo[[]byte]{
					Report: []byte{},
					Info:   info,
				},
				TransmissionScheduleOverride: nil,
			})
			failureIDs = append(failureIDs, failedOutcome.RequestID)
		default:
			r.lggr.Warnw("received unknown consensus outcome type", "seqNr", seqNr, "outcome", outcome)
		}

		if len(reports) == r.maxNumberOfReports {
			r.lggr.Warnw("maximum number of reports reached, stopping further report generation for this round", "seqNr", seqNr, "maxNumberOfReports", r.maxNumberOfReports)
			break
		}
	}

	r.lggr.Debugw("consensus plugin reports complete", "seqNr", seqNr, "numReports", len(reports), "successIDs", successIDs, "failureIDs", failureIDs)
	return reports, nil
}

// The report info is created as a map else the OCR3OnchainKeyringMultiChainAdapter will not work.
// OCR3OnchainKeyringMultiChainAdapter (in core) requires that the key bundle id is added to the map with the key
// "keyBundleName".
func createSuccessfulConsensusReportInfo(reqMetadata *oracletypes.RequestMetaData) ([]byte, error) {
	infos, err := structpb.NewStruct(map[string]any{
		InfoKeyBundleName: reqMetadata.KeyBundleId,
		InfoRequestID:     reqMetadata.RequestId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create structpb for report info: %w", err)
	}

	infoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(infos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return infoBytes, nil
}

func createFailedConsensusReportInfo(requestID string, keyBundleID string, failureMessage string,
	failureCode oracletypes.ConsensusFailureCode) ([]byte, error) {
	infos, err := structpb.NewStruct(map[string]any{
		InfoKeyBundleName:           keyBundleID,
		InfoRequestID:               requestID,
		InfoConsensusFailureMessage: failureMessage,
		InfoConsensusFailureCode:    failureCode.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create structpb for report info: %w", err)
	}

	infoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(infos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return infoBytes, nil
}
