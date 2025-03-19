package aggregation

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
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
	return &signedReportRemoteAggregator{
		allowedSigners:        signersMap,
		minRequiredSignatures: minRequiredSignatures,
		maxAgeSec:             maxAgeSec,
		capID:                 capID,
		lggr:                  lggr,
	}
}

// Accept first response with valid signatures and expected event ID
func (a *signedReportRemoteAggregator) Aggregate(triggerEventID string, responses [][]byte) (capabilities.TriggerResponse, error) {
	for _, response := range responses {
		triggerResp, err := capabilitiespb.UnmarshalTriggerResponse(response)
		if err != nil {
			a.lggr.Errorw("could not unmarshal one of capability responses (faulty sender?)", "error", err)
			continue
		}
		if triggerResp.Event.OCREvent == nil || len(triggerResp.Event.OCREvent.Report) == 0 {
			a.lggr.Errorw("trigger response does not contain an OCR report", "id", triggerResp.Event.ID)
			continue
		}

		rawReport := triggerResp.Event.OCREvent.Report
		rep := &capabilitiespb.OCRTriggerReport{}
		err = proto.Unmarshal(rawReport, rep)
		if err != nil {
			a.lggr.Errorw("failed to parse OCR report", "id", triggerResp.Event.ID)
			continue
		}

		if rep.EventID != triggerEventID {
			a.lggr.Debugw("unexpected event ID", "expected", triggerEventID, "got", rep.EventID)
			continue
		}

		timeDiff := int64(rep.Timestamp) - time.Now().UnixNano()
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		if timeDiff > int64(a.maxAgeSec)*1000000000 { // nanoseconds
			a.lggr.Debugw("report too old", "age", timeDiff)
			continue
		}

		if err := a.validateSignatures(triggerResp.Event.OCREvent); err != nil {
			a.lggr.Errorw("invalid signatures", "err", err)
			continue
		}
		// Replace "Outputs" field with the one extracted from the OCR report and drop the binary report
		outputsMap, err := values.FromMapValueProto(rep.Outputs)
		if err != nil {
			a.lggr.Errorw("failed to parse OCR report outputs", "err", err)
			continue
		}
		triggerResp.Event.Outputs = outputsMap
		triggerResp.Event.OCREvent.Report = nil
		return triggerResp, nil
	}
	return capabilities.TriggerResponse{}, errors.New("no valid response found")
}

func (a *signedReportRemoteAggregator) validateSignatures(event *capabilities.OCRTriggerEvent) error {
	digest, err := ocr2types.BytesToConfigDigest(event.ConfigDigest)
	if err != nil {
		return fmt.Errorf("malformed config digest: %w", err)
	}
	fullHash := ocr2key.ReportToSigData3(digest, event.SeqNr, event.Report)
	validated := map[common.Address]struct{}{}
	for _, sig := range event.Sigs {
		signerPubkey, err2 := crypto.SigToPub(fullHash, sig.Signature)
		if err2 != nil {
			return fmt.Errorf("malformed signer: %w", err2)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)
		if _, ok := a.allowedSigners[signerAddr]; !ok {
			a.lggr.Debugw("invalid signer", "signerAddr", signerAddr)
			continue
		}
		validated[signerAddr] = struct{}{}
		if len(validated) >= a.minRequiredSignatures {
			break // early exit
		}
	}
	if len(validated) < a.minRequiredSignatures {
		return fmt.Errorf("not enough valid signatures %d, needed %d", len(validated), a.minRequiredSignatures)
	}
	return nil
}
