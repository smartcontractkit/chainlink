package vault

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

func (r *ReportingPlugin) generateProtoReport(id string, requestType vaultcommon.RequestType, msg proto.Message) (ocr3types.ReportWithInfo[[]byte], error) {
	if msg == nil {
		return ocr3types.ReportWithInfo[[]byte]{}, errors.New("invalid report: response cannot be nil")
	}

	rpb, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal response to proto: %w", err)
	}

	rip, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vaultcommon.ReportInfo{
		Id:          id,
		RequestType: requestType,
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_PROTOBUF,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return wrapReportWithKeyBundleInfo(rpb, rip)
}

func (r *ReportingPlugin) generateJSONReport(id string, requestType vaultcommon.RequestType, msg proto.Message) (ocr3types.ReportWithInfo[[]byte], error) {
	if msg == nil {
		return ocr3types.ReportWithInfo[[]byte]{}, errors.New("invalid report: response cannot be nil")
	}

	jsonb, err := vaultutils.ToCanonicalJSON(msg)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to convert proto to canonical JSON: %w", err)
	}

	rip, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vaultcommon.ReportInfo{
		Id:          id,
		RequestType: requestType,
		Format:      vaultcommon.ReportFormat_REPORT_FORMAT_JSON,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return wrapReportWithKeyBundleInfo(jsonb, rip)
}

func wrapReportWithKeyBundleInfo(report []byte, reportInfo []byte) (ocr3types.ReportWithInfo[[]byte], error) {
	infos, err := structpb.NewStruct(map[string]any{
		// Use the EVM key bundle to sign the report.
		"keyBundleName": "evm",
		"reportInfo":    reportInfo,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, err
	}

	ip, err := proto.MarshalOptions{Deterministic: true}.Marshal(infos)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, err
	}

	return ocr3types.ReportWithInfo[[]byte]{
		Report: report,
		Info:   ip,
	}, nil
}
