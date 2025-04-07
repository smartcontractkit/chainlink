package aggregation_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/proto"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	// "github.com/smartcontractkit/chainlink/common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

func TestSignedReportAggregator_Aggregate(t *testing.T) {
	t.Parallel()

	// Setup test keys
	kb1, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)
	addr1 := common.BytesToAddress(kb1.PublicKey())

	kb2, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)
	addr2 := common.BytesToAddress(kb2.PublicKey())

	kb3, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)
	addr3 := common.BytesToAddress(kb3.PublicKey())

	// kb4 is an unauthorized signer
	kb4, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)

	// Setup test data
	eventID := "test-event-123"
	configDigest := ocr2types.ConfigDigest{}
	for i := 0; i < len(configDigest); i++ {
		configDigest[i] = byte(i)
	}
	seqNr := uint64(123)
	currentTime := uint64(time.Now().UnixNano()) //nolint:gosec // disable G115

	// Create a valid outputs map
	outputsMap, err := values.NewMap(map[string]values.Value{
		"result": values.NewString("success"),
		"data":   values.NewInt64(42),
	})
	require.NoError(t, err)
	outputsProto := values.ProtoMap(outputsMap) // outputsMap.ProtoMap()

	// Create a valid OCR report
	validReport := &capabilitiespb.OCRTriggerReport{
		EventID:   eventID,
		Timestamp: currentTime,
		Outputs:   outputsProto,
	}

	reportBytes, err := proto.Marshal(validReport)
	require.NoError(t, err)

	sig1, err := kb1.Sign3(configDigest, seqNr, reportBytes)
	require.NoError(t, err)
	sig2, err := kb2.Sign3(configDigest, seqNr, reportBytes)
	require.NoError(t, err)
	unauthorizedSig, err := kb4.Sign3(configDigest, seqNr, reportBytes) // Unauthorized signature
	require.NoError(t, err)

	// Setup logger
	lggr, observedLogs := logger.TestObserved(t, zapcore.WarnLevel)

	// Setup aggregator with 2 required signatures
	allowedSigners := [][]byte{
		addr1.Bytes(),
		addr2.Bytes(),
		addr3.Bytes(),
	}
	minRequiredSignatures := 2
	maxAgeSec := 10
	capID := "test-capability"

	aggregator := aggregation.NewSignedReportRemoteAggregator(
		allowedSigners,
		minRequiredSignatures,
		capID,
		maxAgeSec,
		lggr,
	)

	// Test cases
	// NOTE: we are checking the logs for errors, so we need to clear the logs before invocation of the aggregator
	t.Run("happy path - valid response with enough signatures", func(t *testing.T) {
		// Prepare valid OCR event with enough valid signatures
		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Report:       reportBytes,
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: sig1},
				{Signature: sig2},
			},
		}
		m, err := ocrEvent.ToMap()
		require.NoError(t, err)
		// Create trigger event with OCR report
		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: m,
		}

		// Create TriggerResponse
		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		// Marshal response
		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Expected flattened response
		expectedResp := capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				ID:      eventID,
				Outputs: outputsMap,
			},
		}

		// Call Aggregate
		clearLogger(t, observedLogs)
		result, err := aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.NoError(t, err)

		// Verify response is flattened correctly
		assert.Equal(t, expectedResp.Event.ID, result.Event.ID)
		assert.Equal(t, expectedResp.Event.Outputs, result.Event.Outputs)
	})

	t.Run("error - unmarshallable response", func(t *testing.T) {
		// Invalid bytes that can't be unmarshalled
		invalidBytes := []byte("not a valid response")

		// Call Aggregate with invalid bytes
		clearLogger(t, observedLogs)
		_, err := aggregator.Aggregate(eventID, [][]byte{invalidBytes})
		require.ErrorIs(t, err, aggregation.ErrMissingResponse)
		gotLog := observedLogs.FilterMessage("could not unmarshal one of capability responses (faulty sender?)")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - response without OCR event", func(t *testing.T) {
		// Create event without OCR data
		triggerEvent := capabilities.TriggerEvent{
			ID: eventID,
			// No OCREvent
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with response without OCR event
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("trigger response does not contain an OCR report")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - response with malformed Output map", func(t *testing.T) {
		// Create event with malformed Output map
		o, err := values.NewMap(map[string]any{
			"UnexpectedKey": "unexpected value",
		})
		require.NoError(t, err)
		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: o,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with response without OCR event
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("trigger response does not contain an OCR report")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - unparsable OCR report", func(t *testing.T) {
		// Create an OCR event with invalid report bytes
		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Report:       []byte("not a valid report"),
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: sig1},
				{Signature: sig2},
			},
		}
		m, err := ocrEvent.ToMap()
		require.NoError(t, err)
		triggerEvent := capabilities.TriggerEvent{
			ID: eventID,

			Outputs: m,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with invalid report bytes
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("failed to parse OCR report")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - mismatched event ID", func(t *testing.T) {
		// Create a valid report but with wrong event ID
		wrongIDReport := &capabilitiespb.OCRTriggerReport{
			EventID:   "wrong-event-id",
			Timestamp: currentTime,
			Outputs:   outputsProto,
		}

		wrongReportBytes, err := proto.Marshal(wrongIDReport)
		require.NoError(t, err)

		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Report:       wrongReportBytes,
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: sig1},
				{Signature: sig2},
			},
		}
		m, err := ocrEvent.ToMap()
		require.NoError(t, err)

		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: m,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with mismatched event ID
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("unexpected event ID")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - report too old", func(t *testing.T) {
		// Create an old report (beyond maxAgeSec)
		oldTime := currentTime - uint64((maxAgeSec+5)*1000000000) //nolint:gosec // disable G115

		oldReport := &capabilitiespb.OCRTriggerReport{
			EventID:   eventID,
			Timestamp: oldTime,
			Outputs:   outputsProto,
		}

		oldReportBytes, err := proto.Marshal(oldReport)
		require.NoError(t, err)

		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Report:       oldReportBytes,
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: sig1},
				{Signature: sig2},
			},
		}
		m, err := ocrEvent.ToMap()
		require.NoError(t, err)

		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: m,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with old report
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("aggregation report too old")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - not enough valid signatures", func(t *testing.T) {
		// Only one valid signature when two are required
		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Report:       reportBytes,
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: sig1},
				// Missing second valid signature
			},
		}
		m, err := ocrEvent.ToMap()
		require.NoError(t, err)
		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: m,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with not enough valid signatures
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("invalid signatures")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - signatures from unauthorized signers", func(t *testing.T) {
		// One valid signature and one unauthorized
		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Report:       reportBytes,
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: sig1},
				{Signature: unauthorizedSig}, // Unauthorized signer
			},
		}

		m, err := ocrEvent.ToMap()
		require.NoError(t, err)

		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: m,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with unauthorized signer
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("invalid signer")
		assert.Equal(t, 1, gotLog.Len())
	})

	t.Run("error - malformed config digest", func(t *testing.T) {
		// Invalid config digest
		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: []byte("invalid config digest"), // Wrong length
			SeqNr:        seqNr,
			Report:       reportBytes,
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: sig1},
				{Signature: sig2},
			},
		}

		m, err := ocrEvent.ToMap()
		require.NoError(t, err)
		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: m,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with malformed config digest
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("invalid signatures")
		assert.Equal(t, 1, gotLog.Len())
		gotLogLine := gotLog.All()[0]

		assert.True(t, logContainErr(t, gotLogLine, aggregation.ErrMalformedConfig), "expected error to be contained in log")
	})

	t.Run("error - malformed signature", func(t *testing.T) {
		// Invalid signature format
		ocrEvent := &capabilities.OCRTriggerEvent{
			ConfigDigest: configDigest[:],
			SeqNr:        seqNr,
			Report:       reportBytes,
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signature: []byte("invalid signature")},
				{Signature: sig2},
			},
		}
		m, err := ocrEvent.ToMap()
		require.NoError(t, err)

		triggerEvent := capabilities.TriggerEvent{
			ID:      eventID,
			Outputs: m,
		}

		triggerResp := capabilities.TriggerResponse{
			Event: triggerEvent,
		}

		respBytes, err := capabilitiespb.MarshalTriggerResponse(triggerResp)
		require.NoError(t, err)

		// Call Aggregate with malformed signature
		clearLogger(t, observedLogs)
		_, err = aggregator.Aggregate(eventID, [][]byte{respBytes})
		require.Error(t, err)
		gotLog := observedLogs.FilterMessage("invalid signatures")
		assert.Equal(t, 1, gotLog.Len())
		gotLogLine := gotLog.All()[0]
		assert.True(t, logContainErr(t, gotLogLine, aggregation.ErrMalformedSigner), "expected error to be contained in log")
	})

}

// TODO this seems useful in our common logging package
func logContainErr(t testing.TB, log observer.LoggedEntry, err error) bool {
	gotErr := extractErr(t, log)
	if gotErr == nil {
		t.Logf("expected error, got nil")
		return false
	}
	return errors.Is(gotErr, err)
}

func extractErr(t testing.TB, log observer.LoggedEntry) error {
	t.Helper()
	const errKey = "err"
	ctxMap := log.ContextMap()
	if ctxMap == nil {
		t.Logf("expected context map, got %v", log)
		return nil
	}
	if _, ok := ctxMap[errKey]; !ok {
		t.Logf("expected %s key in context map, got %v", errKey, ctxMap)
		return nil
	}
	// don't use ctxMap because it erases the type
	var f zapcore.Field
	for _, c := range log.Context {
		if c.Key == errKey {
			f = c
			break
		}
	}
	gotErr, ok := f.Interface.(error)
	if !ok {
		t.Logf("expected error type, got %T", ctxMap[errKey])
		return nil
	}
	return gotErr
}

func clearLogger(t testing.TB, lggr *observer.ObservedLogs) {
	t.Helper()
	lggr.TakeAll()
}
