package aggregation

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
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
// Supports protobuf (Format 5 / ReportFormatCapabilityTrigger): OCRTriggerReport.
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

		rep := &capabilitiespb.OCRTriggerReport{}
		err = proto.Unmarshal(rawReport, rep)
		if err != nil || rep.EventID == "" || rep.Timestamp == 0 {
			a.lggr.Errorw("failed to parse OCR report as protobuf", "id", triggerResp.Event.ID)
			continue
		}
		if rep.EventID != triggerEventID {
			a.lggr.Warnw("unexpected event ID", "expected", triggerEventID, "got", rep.EventID)
			continue
		}
		timeDiff := time.Since(time.Unix(0, int64(rep.Timestamp))).Abs() //nolint:gosec
		if timeDiff.Nanoseconds() > int64(a.maxAgeSec)*1000000000 {
			a.lggr.Warnw("aggregation report too old", "age", timeDiff, "maxAge", a.maxAgeSec, "reportTimestamp", rep.Timestamp)
			continue
		}
		if err := a.validateSignatures(ocrEvent); err != nil {
			a.lggr.Errorw("invalid signatures", "err", err)
			continue
		}
		outputsMap, err := values.FromMapValueProto(rep.Outputs)
		if err != nil {
			a.lggr.Errorw("failed to parse OCR report outputs", "err", err)
			continue
		}
		triggerResp.Event.Outputs = outputsMap
		return triggerResp, nil
	}
	return capabilities.TriggerResponse{}, fmt.Errorf("%w: %s", ErrMissingResponse, triggerEventID)
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
