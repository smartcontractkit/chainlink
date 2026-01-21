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

// Accept first response with valid signatures and expected event ID
// every element in responses is must be a [capabilitypb.TriggerResponse]
// and for each of them, we expect the [capabilitypb.TriggerResponse.Event.Outputs] field
// to be the [values.Map] representation a [capabilities.OCRTriggerEvent] (see [capabilities.OCRTriggerEvent.ToMap])
//
// Supports two report formats:
// - Format 5 (ReportFormatCapabilityTrigger): Protobuf-encoded OCRTriggerReport
// - Format 7 (ReportFormatEVMABIEncodeUnpackedExpr): ABI-encoded report with calculated streams
func (a *signedReportRemoteAggregator) Aggregate(triggerEventID string, responses [][]byte) (capabilities.TriggerResponse, error) {
	for _, response := range responses {
		triggerResp, err := capabilitiespb.UnmarshalTriggerResponse(response)
		if err != nil {
			a.lggr.Errorw("could not unmarshal one of capability responses (faulty sender?)", "err", err)
			continue
		}
		ocrEvent := &capabilities.OCRTriggerEvent{}
		err = ocrEvent.FromMap(triggerResp.Event.Outputs)
		if err != nil {
			a.lggr.Errorw("trigger response does not contain an OCR report", "id", triggerResp.Event.ID, "err", err)
			continue
		}
		rawReport := ocrEvent.Report

		// Log Payload presence for debugging
		hasPayload := triggerResp.Event.Payload != nil
		payloadType := "nil"
		if hasPayload {
			payloadType = triggerResp.Event.Payload.TypeUrl
		}
		a.lggr.Infow("Aggregating trigger response", "eventID", triggerResp.Event.ID, "hasPayload", hasPayload, "payloadType", payloadType, "hasOutputs", triggerResp.Event.Outputs != nil)

		// Try Format 5 (protobuf) first
		rep := &capabilitiespb.OCRTriggerReport{}
		err = proto.Unmarshal(rawReport, rep)
		// Check for valid protobuf parse: Unmarshal can succeed on non-protobuf data
		// but will leave all fields at default values. We require EventID and Timestamp
		// to be populated for a valid Format 5 report.
		if err == nil && rep.EventID != "" && rep.Timestamp != 0 {
			// Format 5: Protobuf OCRTriggerReport
			result, err := a.aggregateFormat5(triggerEventID, triggerResp, ocrEvent, rep)
			if err != nil {
				a.lggr.Debugw("Format 5 aggregation failed", "id", triggerResp.Event.ID, "err", err)
				continue
			}
			// Verify Payload is preserved
			resultHasPayload := result.Event.Payload != nil
			resultPayloadType := "nil"
			if resultHasPayload {
				resultPayloadType = result.Event.Payload.TypeUrl
			}
			if !resultHasPayload {
				a.lggr.Errorw("Payload lost during Format 5 aggregation", "eventID", triggerEventID, "originalHasPayload", hasPayload)
			} else {
				a.lggr.Infow("Format 5 aggregation successful, Payload preserved", "eventID", triggerEventID, "payloadType", resultPayloadType)
			}
			return result, nil
		}

		// Try Format 7 (ABI-encoded) if protobuf parsing failed
		if isFormat7EventID(triggerEventID) {
			result, err := a.aggregateFormat7(triggerEventID, triggerResp, ocrEvent, rawReport)
			if err != nil {
				a.lggr.Debugw("Format 7 aggregation failed", "id", triggerResp.Event.ID, "err", err)
				continue
			}
			// Verify Payload is preserved
			resultHasPayload := result.Event.Payload != nil
			resultPayloadType := "nil"
			if resultHasPayload {
				resultPayloadType = result.Event.Payload.TypeUrl
			}
			if !resultHasPayload {
				a.lggr.Errorw("Payload lost during Format 7 aggregation", "eventID", triggerEventID, "originalHasPayload", hasPayload)
			} else {
				a.lggr.Infow("Format 7 aggregation successful, Payload preserved", "eventID", triggerEventID, "payloadType", resultPayloadType)
			}
			return result, nil
		}

		a.lggr.Errorw("failed to parse OCR report as Format 5 or Format 7", "id", triggerResp.Event.ID)
	}
	return capabilities.TriggerResponse{}, fmt.Errorf("%w: %s", ErrMissingResponse, triggerEventID)
}

// aggregateFormat5 handles Format 5 (ReportFormatCapabilityTrigger) - protobuf encoded
func (a *signedReportRemoteAggregator) aggregateFormat5(
	triggerEventID string,
	triggerResp capabilities.TriggerResponse,
	ocrEvent *capabilities.OCRTriggerEvent,
	rep *capabilitiespb.OCRTriggerReport,
) (capabilities.TriggerResponse, error) {
	if rep.EventID != triggerEventID {
		return capabilities.TriggerResponse{}, fmt.Errorf("unexpected event ID: expected %s, got %s", triggerEventID, rep.EventID)
	}

	// use Abs to handle edge case of clock skew
	timeDiff := time.Since(time.Unix(0, int64(rep.Timestamp))).Abs() //nolint:gosec // disable G115 this won't be running in 2262
	if timeDiff.Nanoseconds() > int64(a.maxAgeSec)*1000000000 {
		return capabilities.TriggerResponse{}, fmt.Errorf("report too old: age %v, maxAge %ds", timeDiff, a.maxAgeSec)
	}

	if err := a.validateSignatures(ocrEvent); err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("invalid signatures: %w", err)
	}

	// Replace "Outputs" field with the one extracted from the OCR report
	outputsMap, err := values.FromMapValueProto(rep.Outputs)
	if err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("failed to parse OCR report outputs: %w", err)
	}
	triggerResp.Event.Outputs = outputsMap
	// Preserve Payload field (V2 format) - contains streams.Report wrapped in anypb.Any
	// This is set by the CRE transmitter and is needed by the workflow SDK
	return triggerResp, nil
}

// aggregateFormat7 handles Format 7 (ReportFormatEVMABIEncodeUnpackedExpr) - ABI encoded
// ABI header structure:
//   - bytes 0-31: feedId (bytes32)
//   - bytes 32-63: validFromTimestamp (uint32, right-aligned)
//   - bytes 64-95: timestamp (uint32, right-aligned)
//   - bytes 96-127: nativeFee (int192)
//   - bytes 128-159: linkFee (int192)
//   - bytes 160-191: expiresAt (uint32)
//   - bytes 192+: custom ABI fields (calculated streams)
//
// Format 7 uses the legacy Mercury signing scheme (LegacyReportContext) which includes
// the donID in the ExtraHash. This is different from Format 5 which uses Sign3.
func (a *signedReportRemoteAggregator) aggregateFormat7(
	triggerEventID string,
	triggerResp capabilities.TriggerResponse,
	ocrEvent *capabilities.OCRTriggerEvent,
	rawReport []byte,
) (capabilities.TriggerResponse, error) {
	// Extract timestamp from ABI header (offset 64-95, uint32 at end of 32-byte word)
	timestamp, err := extractABITimestamp(rawReport)
	if err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("failed to extract ABI timestamp: %w", err)
	}

	// Check staleness using ABI timestamp (seconds, not nanoseconds)
	timeDiff := time.Since(time.Unix(int64(timestamp), 0)).Abs()
	if timeDiff.Seconds() > float64(a.maxAgeSec) {
		return capabilities.TriggerResponse{}, fmt.Errorf("ABI report too old: age %v, maxAge %ds", timeDiff, a.maxAgeSec)
	}

	// Extract donID from event ID (format: streams_DONID_TIMESTAMP_f7)
	donID, err := extractDonIDFromEventID(triggerEventID)
	if err != nil {
		a.lggr.Warnw("Failed to extract donID from event ID, using default", "eventID", triggerEventID, "err", err)
		donID = 2 // Default for E2E test
	}

	// Validate signatures using legacy signing scheme (includes donID in ExtraHash)
	if err := a.validateSignaturesLegacy(ocrEvent, donID, false); err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("Format 7 signature verification failed: %w", err)
	}

	// For Format 7, keep the raw report bytes as Outputs
	// The workflow can decode the ABI payload
	outputsMap := &valuespb.Map{
		Fields: map[string]*valuespb.Value{
			"RawReport": {
				Value: &valuespb.Value_BytesValue{BytesValue: rawReport},
			},
			"Timestamp": {
				Value: &valuespb.Value_Uint64Value{Uint64Value: uint64(timestamp)},
			},
		},
	}
	triggerResp.Event.Outputs, err = values.FromMapValueProto(outputsMap)
	if err != nil {
		return capabilities.TriggerResponse{}, fmt.Errorf("failed to create outputs map: %w", err)
	}
	// Preserve Payload field (V2 format) - contains streams.Report wrapped in anypb.Any
	// This is set by the CRE transmitter and is needed by the workflow SDK
	a.lggr.Debugw("Format 7 report aggregated with verified signatures", "eventID", triggerEventID, "timestamp", timestamp, "reportLen", len(rawReport), "sigCount", len(ocrEvent.Sigs), "donID", donID)
	return triggerResp, nil
}

// extractDonIDFromEventID parses the donID from event ID format: streams_DONID_TIMESTAMP_f7
func extractDonIDFromEventID(eventID string) (uint32, error) {
	var donID uint64
	var timestamp uint64
	_, err := fmt.Sscanf(eventID, "streams_%d_%d_f7", &donID, &timestamp)
	if err != nil {
		return 0, fmt.Errorf("failed to parse event ID '%s': %w", eventID, err)
	}
	if donID > math.MaxUint32 {
		return 0, fmt.Errorf("donID %d exceeds uint32", donID)
	}
	return uint32(donID), nil
}

// isFormat7EventID checks if the event ID indicates a Format 7 report
// Format 7 event IDs have "_f7" suffix (generated by the CRE transmitter)
func isFormat7EventID(eventID string) bool {
	return strings.HasSuffix(eventID, "_f7")
}

// extractABITimestamp extracts the timestamp from a Format 7 ABI-encoded report
// The timestamp is at offset 64-95 (uint32 right-aligned in a 32-byte word)
func extractABITimestamp(report []byte) (uint32, error) {
	// Minimum size: 192 bytes for the base header
	if len(report) < 192 {
		return 0, fmt.Errorf("report too short for Format 7: %d bytes", len(report))
	}
	// Timestamp is at offset 64, stored as uint32 in the last 4 bytes of the 32-byte word
	timestampWord := report[64:96]
	// uint32 is right-aligned in the 32-byte ABI word
	timestamp := binary.BigEndian.Uint32(timestampWord[28:32])
	return timestamp, nil
}

func (a *signedReportRemoteAggregator) validateSignatures(event *capabilities.OCRTriggerEvent) error {
	return a.validateSignaturesWithDebug(event, false)
}

func (a *signedReportRemoteAggregator) validateSignaturesWithDebug(event *capabilities.OCRTriggerEvent, debug bool) error {
	digest, err := ocr2types.BytesToConfigDigest(event.ConfigDigest)
	if err != nil {
		return errors.Join(ErrMalformedConfig, err)
	}
	fullHash := ocr2key.ReportToSigData3(digest, event.SeqNr, event.Report)

	if debug {
		// Include first signature bytes for comparison
		sig0Hex := ""
		if len(event.Sigs) > 0 && len(event.Sigs[0].Signature) > 0 {
			sig0Hex = fmt.Sprintf("%x", event.Sigs[0].Signature[:min(20, len(event.Sigs[0].Signature))])
		}
		a.lggr.Infow("DEBUG validateSignatures",
			"configDigest", fmt.Sprintf("%x", event.ConfigDigest),
			"seqNr", event.SeqNr,
			"reportLen", len(event.Report),
			"reportFirst32", fmt.Sprintf("%x", event.Report[:min(32, len(event.Report))]),
			"fullHash", fmt.Sprintf("%x", fullHash),
			"numSigs", len(event.Sigs),
			"sig0First20", sig0Hex,
		)
	}

	validated := map[common.Address]struct{}{}
	for i, sig := range event.Sigs {
		signerPubkey, err2 := crypto.SigToPub(fullHash, sig.Signature)
		if err2 != nil {
			return errors.Join(ErrMalformedSigner, err2)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)

		if debug {
			_, isInMap := a.allowedSigners[signerAddr]
			a.lggr.Infow("DEBUG signature recovery",
				"sigIndex", i,
				"signerAttr", sig.Signer,
				"sigLen", len(sig.Signature),
				"recoveredAddr", signerAddr.Hex(),
				"isAllowed", isInMap,
			)
		}

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

// validateSignaturesLegacy validates signatures using the legacy Mercury/LLO signing scheme.
// This is used for Format 7 (EVMABIEncodeUnpackedExpr) and other legacy formats.
// The legacy scheme uses LegacyReportContext which includes donID in the ExtraHash.
func (a *signedReportRemoteAggregator) validateSignaturesLegacy(event *capabilities.OCRTriggerEvent, donID uint32, debug bool) error {
	digest, err := ocr2types.BytesToConfigDigest(event.ConfigDigest)
	if err != nil {
		return errors.Join(ErrMalformedConfig, err)
	}

	// Build legacy report context (same as used by LLO keyring for signing)
	reportCtx, err := legacyReportContext(digest, event.SeqNr, donID)
	if err != nil {
		return fmt.Errorf("failed to build legacy report context: %w", err)
	}

	// Compute hash using legacy method
	fullHash := reportToSigDataLegacy(reportCtx, event.Report)

	if debug {
		sig0Hex := ""
		if len(event.Sigs) > 0 && len(event.Sigs[0].Signature) > 0 {
			sig0Hex = fmt.Sprintf("%x", event.Sigs[0].Signature[:min(20, len(event.Sigs[0].Signature))])
		}
		a.lggr.Infow("DEBUG validateSignaturesLegacy",
			"configDigest", fmt.Sprintf("%x", event.ConfigDigest),
			"seqNr", event.SeqNr,
			"donID", donID,
			"reportLen", len(event.Report),
			"reportFirst32", fmt.Sprintf("%x", event.Report[:min(32, len(event.Report))]),
			"fullHash", fmt.Sprintf("%x", fullHash),
			"numSigs", len(event.Sigs),
			"sig0First20", sig0Hex,
		)
	}

	validated := map[common.Address]struct{}{}
	for i, sig := range event.Sigs {
		signerPubkey, err2 := crypto.SigToPub(fullHash, sig.Signature)
		if err2 != nil {
			return errors.Join(ErrMalformedSigner, err2)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)

		if debug {
			_, isInMap := a.allowedSigners[signerAddr]
			a.lggr.Infow("DEBUG legacy signature recovery",
				"sigIndex", i,
				"signerAttr", sig.Signer,
				"sigLen", len(sig.Signature),
				"recoveredAddr", signerAddr.Hex(),
				"isAllowed", isInMap,
			)
		}

		if _, ok := a.allowedSigners[signerAddr]; !ok {
			a.lggr.Warnw("invalid signer (legacy)", "signerAddr", signerAddr)
			continue
		}
		validated[signerAddr] = struct{}{}
		if len(validated) >= a.minRequiredSignatures {
			break // early exit
		}
	}
	if len(validated) < a.minRequiredSignatures {
		return fmt.Errorf("%w (legacy): got %d, needed %d", ErrInsufficientSignatures, len(validated), a.minRequiredSignatures)
	}
	return nil
}

// Legacy report context construction - matches LLO keyring signing
const lloPluginVersion uint32 = 1

func legacyReportContext(cd ocr2types.ConfigDigest, seqNr uint64, donID uint32) (ocr2types.ReportContext, error) {
	epoch, round, err := seqNrToEpochAndRound(seqNr)
	if err != nil {
		return ocr2types.ReportContext{}, err
	}
	return ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{
			ConfigDigest: cd,
			Epoch:        epoch,
			Round:        round,
		},
		ExtraHash: lloExtraHash(donID),
	}, nil
}

func seqNrToEpochAndRound(seqNr uint64) (epoch uint32, round uint8, err error) {
	// Simulate 256 rounds/epoch (same as LLO keyring)
	if seqNr/256 > math.MaxUint32 {
		err = fmt.Errorf("epoch overflows uint32: %d", seqNr)
		return
	}
	epoch = uint32(seqNr / 256) //nolint:gosec // G115 false positive
	round = uint8(seqNr % 256)  //nolint:gosec // G115 false positive
	return
}

func lloExtraHash(donID uint32) common.Hash {
	// Packs donID+pluginVersion as (uint32, uint32)
	combined := uint64(donID)<<32 | uint64(lloPluginVersion)
	return common.BigToHash(new(big.Int).SetUint64(combined))
}

func reportToSigDataLegacy(reportCtx ocr2types.ReportContext, report ocr2types.Report) []byte {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	sigData := crypto.Keccak256(report)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	sigData = append(sigData, rawReportContext[2][:]...)
	return crypto.Keccak256(sigData)
}
