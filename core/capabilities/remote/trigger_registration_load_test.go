package remote_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// donFaultTolerance is the DON F value (uint8) used in load tests; v must be in [0, 255].
func donFaultTolerance(t *testing.T, v int) uint8 {
	t.Helper()
	require.GreaterOrEqual(t, v, 0)
	require.LessOrEqual(t, v, 255, "F value out of uint8 range for test DON")
	return uint8(v) //nolint:gosec // G115: range-checked above
}

// noopDispatcher is a types.Dispatcher that discards all sends.
type noopDispatcher struct {
	services.StateMachine
}

func (d *noopDispatcher) Send(_ p2ptypes.PeerID, _ *remotetypes.MessageBody) error { return nil }
func (d *noopDispatcher) SetReceiver(_ string, _ uint32, _ remotetypes.Receiver) error {
	return nil
}
func (d *noopDispatcher) RemoveReceiver(_ string, _ uint32) {}
func (d *noopDispatcher) SetReceiverForMethod(_ string, _ uint32, _ string, _ remotetypes.Receiver) error {
	return nil
}
func (d *noopDispatcher) RemoveReceiverForMethod(_ string, _ uint32, _ string) {}
func (d *noopDispatcher) Start(_ context.Context) error                        { return nil }
func (d *noopDispatcher) Close() error                                         { return nil }
func (d *noopDispatcher) Ready() error                                         { return nil }
func (d *noopDispatcher) HealthReport() map[string]error                       { return nil }
func (d *noopDispatcher) Name() string                                         { return "noopDispatcher" }

// --- helpers ---

func generateWorkflowID(i int) string {
	return fmt.Sprintf("%064x", i)
}

// noopTrigger implements commoncap.TriggerCapability with no-op operations
// and large buffered channels so it never blocks during load tests.
type noopTrigger struct {
	info commoncap.CapabilityInfo
}

func (t *noopTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return t.info, nil
}

func (t *noopTrigger) RegisterTrigger(_ context.Context, _ commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	ch := make(chan commoncap.TriggerResponse, 1)
	return ch, nil
}

func (t *noopTrigger) UnregisterTrigger(_ context.Context, _ commoncap.TriggerRegistrationRequest) error {
	return nil
}

func (t *noopTrigger) AckEvent(_ context.Context, _, _, _ string) error {
	return nil
}

// --- Test: Subscriber registrationLoop traffic volume ---

// --- Test: Publisher sendRegistrationChecks traffic volume ---

// --- Benchmark: Publisher processing rate for duplicate registrations ---

func BenchmarkRegistrationProcessing(b *testing.B) {
	ctx := context.Background()
	lggr, _ := logger.New()
	dispatcher := &noopDispatcher{}

	capDonMembers := make([]p2ptypes.PeerID, 4)
	for i := range capDonMembers {
		pid := utils.MustNewPeerID()
		if err := capDonMembers[i].UnmarshalText([]byte(pid)); err != nil {
			b.Fatal(err)
		}
	}
	capDon := commoncap.DON{ID: 1, Members: capDonMembers, F: 1}

	workflowDonMembers := make([]p2ptypes.PeerID, 7)
	for i := range workflowDonMembers {
		pid := utils.MustNewPeerID()
		if err := workflowDonMembers[i].UnmarshalText([]byte(pid)); err != nil {
			b.Fatal(err)
		}
	}
	workflowDon := commoncap.DON{ID: 2, Members: workflowDonMembers, F: 2}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Benchmark Trigger",
	}
	underlying := &noopTrigger{info: capInfo}

	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     time.Hour,
		RegistrationExpiry:      2 * time.Hour,
		MinResponsesToAggregate: 1,
		MessageExpiry:           time.Hour,
		MaxBatchSize:            100,
		BatchCollectionPeriod:   time.Second,
	}
	workflowDONs := map[uint32]commoncap.DON{workflowDon.ID: workflowDon}
	publisher := remote.NewTriggerPublisher(capInfo.ID, "LogTrigger", dispatcher, lggr)
	if err := publisher.SetConfig(cfg, underlying, capDon, workflowDONs); err != nil {
		b.Fatal(err)
	}
	if err := publisher.Start(ctx); err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { _ = publisher.Close() })

	// Pre-register one trigger so subsequent Receive calls hit the fast
	// "already exists" path.
	sender := workflowDonMembers[0]
	triggerReq := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger_0",
		Metadata:  commoncap.RequestMetadata{WorkflowID: generateWorkflowID(0)},
	}
	marshaled, err := pb.MarshalTriggerRegistrationRequest(triggerReq)
	if err != nil {
		b.Fatal(err)
	}
	regMsg := &remotetypes.MessageBody{
		Sender:      sender[:],
		Method:      remotetypes.MethodRegisterTrigger,
		CallerDonId: workflowDon.ID,
		Payload:     marshaled,
	}
	publisher.Receive(ctx, regMsg)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		publisher.Receive(ctx, regMsg)
	}

	b.StopTimer()
	elapsed := b.Elapsed()
	rate := float64(b.N) / elapsed.Seconds()
	b.ReportMetric(rate, "msgs/sec")
}

// --- Test: Subscriber RLock hold duration during registration loop ---
