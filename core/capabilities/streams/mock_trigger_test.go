package streams

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	triggercfg "github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/streams"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const testFeedID triggercfg.FeedId = "0x0000000000000000000000000000000000000000000000000000000000000001"

func TestMockTriggerServiceStartClose(t *testing.T) {
	t.Parallel()

	svc, err := NewMockTriggerService(10, logger.TestLogger(t))
	require.NoError(t, err)
	require.NoError(t, svc.Start(context.Background()))

	done := make(chan error, 1)
	go func() {
		done <- svc.Close()
	}()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for mock trigger service to close")
	}
}

func TestMockTriggerServiceRegisterEmitUnregister(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	svc, err := NewMockTriggerService(200, logger.TestLogger(t))
	require.NoError(t, err)
	require.NoError(t, svc.MercuryTriggerService.Start(ctx))
	t.Cleanup(func() {
		require.NoError(t, svc.MercuryTriggerService.Close())
	})

	req := newMockTriggerRequest(t, "workflow-a", "trigger-a", 200)
	ch, err := svc.RegisterTrigger(ctx, req)
	require.NoError(t, err)

	svc.subscribersMu.Lock()
	require.Equal(t, []triggercfg.FeedId{testFeedID}, svc.subscribers[req.Metadata.WorkflowID])
	svc.subscribersMu.Unlock()

	require.NoError(t, svc.ProcessReport([]datastreams.FeedReport{newSignedMockFeedReport(t, svc, testFeedID)}))

	var resp commoncap.TriggerResponse
	select {
	case resp = <-ch:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for mock trigger event")
	}

	triggerEvent := datastreams.StreamsTriggerEvent{}
	require.NoError(t, resp.Event.Outputs.UnwrapTo(&triggerEvent))
	require.Equal(t, MockTriggerCapabilityID, resp.Event.TriggerType)
	require.Len(t, triggerEvent.Payload, 1)
	require.Len(t, triggerEvent.Payload[0].Signatures, svc.meta.MinRequiredSignatures)
	require.Len(t, triggerEvent.Metadata.Signers, svc.meta.MinRequiredSignatures)
	require.Equal(t, svc.meta.MinRequiredSignatures, triggerEvent.Metadata.MinRequiredSignatures)

	require.NoError(t, svc.UnregisterTrigger(ctx, req))

	svc.subscribersMu.Lock()
	require.Empty(t, svc.subscribers)
	svc.subscribersMu.Unlock()

	select {
	case _, ok := <-ch:
		require.False(t, ok)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for subscriber channel to close")
	}
}

func TestMockTriggerServiceUnregisterKeepsSubscriberStateOnError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	svc, err := NewMockTriggerService(200, logger.TestLogger(t))
	require.NoError(t, err)
	require.NoError(t, svc.MercuryTriggerService.Start(ctx))
	t.Cleanup(func() {
		require.NoError(t, svc.MercuryTriggerService.Close())
	})

	req := newMockTriggerRequest(t, "workflow-a", "trigger-a", 200)
	_, err = svc.RegisterTrigger(ctx, req)
	require.NoError(t, err)

	badReq := req
	badReq.TriggerID = "trigger-b"
	err = svc.UnregisterTrigger(ctx, badReq)
	require.Error(t, err)

	svc.subscribersMu.Lock()
	require.Equal(t, []triggercfg.FeedId{testFeedID}, svc.subscribers[req.Metadata.WorkflowID])
	svc.subscribersMu.Unlock()
}

func newMockTriggerRequest(t *testing.T, workflowID, triggerID string, maxFrequencyMs uint64) commoncap.TriggerRegistrationRequest {
	t.Helper()

	cfg, err := values.WrapMap(triggercfg.TriggerConfig{
		FeedIds:        []triggercfg.FeedId{testFeedID},
		MaxFrequencyMs: maxFrequencyMs,
	})
	require.NoError(t, err)

	return commoncap.TriggerRegistrationRequest{
		TriggerID: triggerID,
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID,
		},
		Config: cfg,
	}
}

func newSignedMockFeedReport(t *testing.T, svc *MockTriggerService, feedID triggercfg.FeedId) datastreams.FeedReport {
	t.Helper()

	timestamp := time.Now().Unix()
	reportCtx := ocrTypes.ReportContext{
		ReportTimestamp: ocrTypes.ReportTimestamp{Epoch: uint32(baseTimestamp + 1)},
	}

	report := datastreams.FeedReport{
		FeedID:               string(feedID),
		ReportContext:        rawReportContext(reportCtx),
		ObservationTimestamp: timestamp,
	}
	fullReport, err := newReport(svc.lggr, common.HexToHash(string(feedID)), big.NewInt(123456), timestamp)
	require.NoError(t, err)
	report.FullReport = fullReport

	sigData := append(crypto.Keccak256(report.FullReport), report.ReportContext...)
	hash := crypto.Keccak256(sigData)
	for _, signer := range svc.signers {
		sig, err := crypto.Sign(hash, signer)
		require.NoError(t, err)
		report.Signatures = append(report.Signatures, sig)
	}

	return report
}
