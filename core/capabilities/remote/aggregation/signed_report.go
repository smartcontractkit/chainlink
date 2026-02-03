package aggregation

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

var (
	ErrMalformedSigner        = errors.New("malformed signer")
	ErrMalformedConfig        = errors.New("malformed config digest")
	ErrInsufficientSignatures = errors.New("insufficient signatures")
	ErrMissingResponse        = errors.New("missing trigger response")
)

type signedReportRemoteAggregator struct {
	allowedSigners        map[common.Address]struct{}
	minRequiredSignatures int
	maxAgeSec             int
	capID                 string
	lggr                  logger.Logger
}

func NewSignedReportRemoteAggregator(allowedSigners [][]byte, minRequiredSignatures int, capID string, maxAgeSec int, lggr logger.Logger) *signedReportRemoteAggregator {
	signersMap := make(map[common.Address]struct{})
	for _, signer := range allowedSigners {
		signersMap[common.BytesToAddress(signer)] = struct{}{}
	}
	lggr = logger.Named(lggr, "SignedReportRemoteAggregator")

	lggr.Infow("created", "allowedSigners", signersMap, "minRequiredSignatures", minRequiredSignatures, "maxAgeSec", maxAgeSec, "capID", capID)
	return &signedReportRemoteAggregator{
		allowedSigners:        signersMap,
		minRequiredSignatures: minRequiredSignatures,
		maxAgeSec:             maxAgeSec,
		capID:                 capID,
		lggr:                  logger.Named(lggr, "SignedReportRemoteAggregator"),
	}
}

// Accept first response with valid signatures and expected event ID.
// Supports two report encodings:
//   - Protobuf (Format 5 / ReportFormatCapabilityTrigger): OCRTriggerReport
//   - ABI (Format 7 / ReportFormatEVMABIEncodeUnpackedExpr): workflow decodes at execution time
func (a *signedReportRemoteAggregator) Aggregate(triggerEventID string, responses [][]byte) (capabilities.TriggerResponse, error) {
	for _, response := range responses {
		triggerResp, err := capabilitiespb.UnmarshalTriggerResponse(response)
		if err != nil {
			a.lggr.Errorw("could not unmarshal capability response", "err", err)
			continue
		}
		ocrEvent := &capabilities.OCRTriggerEvent{}
		err = ocrEvent.FromMap(triggerResp.Event.Outputs)
		if err != nil {
			a.lggr.Errorw("trigger response does not contain an OCR report", "id", triggerResp.Event.ID, "err", err)
			continue
		}
		rawReport := ocrEvent.Report

		if isABIEncodedEventID(triggerEventID) {
			result, err := a.aggregateABI(triggerEventID, triggerResp, ocrEvent, rawReport)
			if err != nil {
				continue
			}
			return result, nil
		}

		rep := &capabilitiespb.OCRTriggerReport{}
		err = proto.Unmarshal(rawReport, rep)
		if err == nil && rep.EventID != "" && rep.Timestamp != 0 {
			result, err := a.aggregateProtobuf(triggerEventID, triggerResp, ocrEvent, rep)
			if err != nil {
				continue
			}
			return result, nil
		}

		a.lggr.Errorw("failed to parse OCR report as protobuf or ABI", "id", triggerResp.Event.ID)
	}
	return capabilities.TriggerResponse{}, fmt.Errorf("%w: %s", ErrMissingResponse, triggerEventID)
}

// aggregateProtobuf handles Format 5 (protobuf-encoded OCRTriggerReport).
func (a *signedReportRemoteAggregator) aggregateProtobuf(
	triggerEventID string,
	triggerResp capabilities.TriggerResponse,
	ocrEvent *capabilities.OCRTriggerEvent,
	rep *capabilitiespb.OCRTriggerReport,
) (capabilities.TriggerResponse, error) {
	if rep.EventID != triggerEventID {
		return capabilities.TriggerResponse{}, fmt.Errorf("unexpected event ID: expected %s, got %s", triggerEventID, rep.EventID)
	}
	timeDiff := time.Since(time.Unix(0, int64(rep.Timestamp))).Abs() //nolint:gosec
	if timeDiff.Nanoseconds() > int64(a.maxAgeSec)*1000000000 {
		return capabilities.TriggerResponse{}, fmt.Errorf("report too old: age %v, maxAge %ds", timeDiff, a.maxAgeSec)
	}
	if err := a.validateSignatures(ocrEvent); err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("invalid signatures: %w", err)
	}
	outputsMap, err := values.FromMapValueProto(rep.Outputs)
	if err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("failed to parse OCR report outputs: %w", err)
	}
	triggerResp.Event.Outputs = outputsMap
	return triggerResp, nil
}

// aggregateABI handles Format 7 (ABI-encoded). The workflow receives RawReport/Timestamp in Outputs and decodes at execution time.
func (a *signedReportRemoteAggregator) aggregateABI(
	triggerEventID string,
	triggerResp capabilities.TriggerResponse,
	ocrEvent *capabilities.OCRTriggerEvent,
	rawReport []byte,
) (capabilities.TriggerResponse, error) {
	timestamp, err := extractABITimestamp(rawReport)
	if err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("ABI timestamp: %w", err)
	}
	timeDiff := time.Since(time.Unix(int64(timestamp), 0)).Abs()
	if timeDiff.Seconds() > float64(a.maxAgeSec) {
		return capabilities.TriggerResponse{}, fmt.Errorf("ABI report too old: age %v, maxAge %ds", timeDiff, a.maxAgeSec)
	}
	donID, err := extractDonIDFromEventID(triggerEventID)
	if err != nil {
		a.lggr.Warnw("extract donID from event ID, using default", "eventID", triggerEventID, "err", err)
		donID = 2
	}
	if err := a.validateSignaturesLegacy(ocrEvent, donID); err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("ABI signature verification: %w", err)
	}
	outputsMap := &valuespb.Map{
		Fields: map[string]*valuespb.Value{
			"RawReport": {Value: &valuespb.Value_BytesValue{BytesValue: rawReport}},
			"Timestamp": {Value: &valuespb.Value_Uint64Value{Uint64Value: uint64(timestamp)}},
		},
	}
	triggerResp.Event.Outputs, err = values.FromMapValueProto(outputsMap)
	if err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("create outputs: %w", err)
	}
	return triggerResp, nil
}

func isABIEncodedEventID(eventID string) bool {
	return strings.HasSuffix(eventID, "_f7")
}

func extractDonIDFromEventID(eventID string) (uint32, error) {
	var donID uint64
	var timestamp uint64
	_, err := fmt.Sscanf(eventID, "streams_%d_%d_f7", &donID, &timestamp)
	if err != nil {
		return 0, fmt.Errorf("parse event ID %q: %w", eventID, err)
	}
	if donID > math.MaxUint32 {
		return 0, fmt.Errorf("donID %d exceeds uint32", donID)
	}
	return uint32(donID), nil
}

func extractABITimestamp(report []byte) (uint32, error) {
	if len(report) < 192 {
		return 0, fmt.Errorf("ABI report too short: %d bytes", len(report))
	}
	return binary.BigEndian.Uint32(report[92:96]), nil
}

// validateSignaturesLegacy validates ABI/Format 7 reports (legacy Mercury signing with donID in ExtraHash).
func (a *signedReportRemoteAggregator) validateSignaturesLegacy(event *capabilities.OCRTriggerEvent, donID uint32) error {
	digest, err := ocr2types.BytesToConfigDigest(event.ConfigDigest)
	if err != nil {
		return errors.Join(ErrMalformedConfig, err)
	}
	reportCtx, err := legacyReportContext(digest, event.SeqNr, donID)
	if err != nil {
		return fmt.Errorf("legacy report context: %w", err)
	}
	fullHash := reportToSigDataLegacy(reportCtx, event.Report)
	validated := map[common.Address]struct{}{}
	for _, sig := range event.Sigs {
		signerPubkey, err2 := crypto.SigToPub(fullHash, sig.Signature)
		if err2 != nil {
			continue
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)
		if _, ok := a.allowedSigners[signerAddr]; !ok {
			continue
		}
		validated[signerAddr] = struct{}{}
		if len(validated) >= a.minRequiredSignatures {
			break
		}
	}
	if len(validated) < a.minRequiredSignatures {
		return fmt.Errorf("%w (legacy): got %d, needed %d", ErrInsufficientSignatures, len(validated), a.minRequiredSignatures)
	}
	return nil
}

const lloPluginVersion uint32 = 1

func legacyReportContext(cd ocr2types.ConfigDigest, seqNr uint64, donID uint32) (ocr2types.ReportContext, error) {
	epoch, round, err := seqNrToEpochAndRound(seqNr)
	if err != nil {
		return ocr2types.ReportContext{}, err
	}
	return ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{ConfigDigest: cd, Epoch: epoch, Round: round},
		ExtraHash:       lloExtraHash(donID),
	}, nil
}

func seqNrToEpochAndRound(seqNr uint64) (epoch uint32, round uint8, err error) {
	if seqNr/256 > math.MaxUint32 {
		return 0, 0, fmt.Errorf("epoch overflows uint32: %d", seqNr)
	}
	epoch = uint32(seqNr / 256) //nolint:gosec
	round = uint8(seqNr % 256)  //nolint:gosec
	return epoch, round, nil
}

func lloExtraHash(donID uint32) common.Hash {
	combined := uint64(donID)<<32 | uint64(lloPluginVersion)
	return common.BigToHash(new(big.Int).SetUint64(combined))
}

func reportToSigDataLegacy(reportCtx ocr2types.ReportContext, report ocr2types.Report) []byte {
	rawCtx := evmutil.RawReportContext(reportCtx)
	sigData := crypto.Keccak256(report)
	sigData = append(sigData, rawCtx[0][:]...)
	sigData = append(sigData, rawCtx[1][:]...)
	sigData = append(sigData, rawCtx[2][:]...)
	return crypto.Keccak256(sigData)
}

func (a *signedReportRemoteAggregator) validateSignatures(event *capabilities.OCRTriggerEvent) error {
	digest, err := ocr2types.BytesToConfigDigest(event.ConfigDigest)
	if err != nil {
		return errors.Join(ErrMalformedConfig, err)
	}
	fullHash := ocr2key.ReportToSigData3(digest, event.SeqNr, event.Report)
	validated := map[common.Address]struct{}{}
	for _, sig := range event.Sigs {
		signerPubkey, err2 := crypto.SigToPub(fullHash, sig.Signature)
		if err2 != nil {
			return errors.Join(ErrMalformedSigner, err2)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)
		if _, ok := a.allowedSigners[signerAddr]; !ok {
			a.lggr.Warnw("invalid signer", "signerAddr", signerAddr)
			continue
		}
		validated[signerAddr] = struct{}{}
		if len(validated) >= a.minRequiredSignatures {
			break // early exit
		}
	}
	if len(validated) < a.minRequiredSignatures {
		return fmt.Errorf("%w: got %d, needed %d", ErrInsufficientSignatures, len(validated), a.minRequiredSignatures)
	}
	return nil
}
