package remote_test

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
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

// countingDispatcher is a types.Dispatcher implementation that atomically
// counts Send calls by method type.
type countingDispatcher struct {
	services.StateMachine
	mu     sync.RWMutex
	counts map[string]*atomic.Int64
}

func newCountingDispatcher() *countingDispatcher {
	return &countingDispatcher{
		counts: make(map[string]*atomic.Int64),
	}
}

func (d *countingDispatcher) Send(_ p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	method := msgBody.Method
	d.mu.RLock()
	c, ok := d.counts[method]
	d.mu.RUnlock()
	if !ok {
		d.mu.Lock()
		c, ok = d.counts[method]
		if !ok {
			c = &atomic.Int64{}
			d.counts[method] = c
		}
		d.mu.Unlock()
	}
	c.Add(1)
	return nil
}

func (d *countingDispatcher) Count(method string) int64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if c, ok := d.counts[method]; ok {
		return c.Load()
	}
	return 0
}

func (d *countingDispatcher) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.counts = make(map[string]*atomic.Int64)
}

func (d *countingDispatcher) Snapshot() map[string]int64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	snap := make(map[string]int64, len(d.counts))
	for method, c := range d.counts {
		snap[method] = c.Load()
	}
	return snap
}

func (d *countingDispatcher) SetReceiver(_ string, _ uint32, _ remotetypes.Receiver) error {
	return nil
}
func (d *countingDispatcher) RemoveReceiver(_ string, _ uint32) {}
func (d *countingDispatcher) SetReceiverForMethod(_ string, _ uint32, _ string, _ remotetypes.Receiver) error {
	return nil
}
func (d *countingDispatcher) RemoveReceiverForMethod(_ string, _ uint32, _ string) {}
func (d *countingDispatcher) Start(_ context.Context) error                        { return nil }
func (d *countingDispatcher) Close() error                                         { return nil }
func (d *countingDispatcher) Ready() error                                         { return nil }
func (d *countingDispatcher) HealthReport() map[string]error                       { return nil }
func (d *countingDispatcher) Name() string                                         { return "countingDispatcher" }

// formatP2PSnapshot returns a stable, human-readable summary of outbound P2P
// sends by method (for test logs / team evidence).
func formatP2PSnapshot(snap map[string]int64) string {
	if len(snap) == 0 {
		return "(none)"
	}
	methods := make([]string, 0, len(snap))
	for m := range snap {
		methods = append(methods, m)
	}
	sort.Strings(methods)
	var b strings.Builder
	for _, m := range methods {
		fmt.Fprintf(&b, "%s=%d ", m, snap[m])
	}
	return strings.TrimSpace(b.String())
}

// --- helpers ---

func generatePeers(t *testing.T, n int) []p2ptypes.PeerID {
	t.Helper()
	peers := make([]p2ptypes.PeerID, n)
	for i := range n {
		pid := utils.MustNewPeerID()
		require.NoError(t, peers[i].UnmarshalText([]byte(pid)))
	}
	return peers
}

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

func (t *noopTrigger) AckEvent(_ context.Context, _ string, _ string, _ string) error {
	return nil
}

// --- Test: Subscriber registrationLoop traffic volume ---

func TestRegistrationTrafficVolume(t *testing.T) {
	cases := []struct {
		nRegistrations  int
		capDonSize      int
		workflowDonSize int
	}{
		{100, 4, 7},
		{500, 4, 7},
		{1000, 4, 7},
		{2500, 4, 7},
		{5000, 4, 7},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("N=%d_cap=%d_wf=%d", tc.nRegistrations, tc.capDonSize, tc.workflowDonSize), func(t *testing.T) {
			lggr := logger.Test(t)
			dispatcher := newCountingDispatcher()

			capDonMembers := generatePeers(t, tc.capDonSize)
			capDon := commoncap.DON{ID: 1, Members: capDonMembers, F: donFaultTolerance(t, tc.capDonSize/3)}

			workflowDonMembers := generatePeers(t, tc.workflowDonSize)
			workflowDon := commoncap.DON{ID: 2, Members: workflowDonMembers, F: donFaultTolerance(t, (tc.workflowDonSize-1)/3)}

			capInfo := commoncap.CapabilityInfo{
				ID:             "cap_id@1",
				CapabilityType: commoncap.CapabilityTypeTrigger,
				Description:    "Load Test Trigger",
			}

			cfg := &commoncap.RemoteTriggerConfig{
				RegistrationRefresh:     500 * time.Millisecond,
				RegistrationExpiry:      2 * time.Hour,
				MinResponsesToAggregate: 1,
				MessageExpiry:           time.Hour,
			}

			subscriber := remote.NewTriggerSubscriber(capInfo.ID, "LogTrigger", dispatcher, lggr)
			agg := aggregation.NewDefaultModeAggregator(cfg.MinResponsesToAggregate)
			require.NoError(t, subscriber.SetConfig(cfg, capInfo, workflowDon.ID, capDon, agg))

			// Register all triggers before Start — RegisterTrigger does not
			// send to the cap DON until registrationLoop runs after Start.
			for i := range tc.nRegistrations {
				req := commoncap.TriggerRegistrationRequest{
					TriggerID: fmt.Sprintf("trigger_%d", i),
					Metadata:  commoncap.RequestMetadata{WorkflowID: generateWorkflowID(i)},
				}
				_, err := subscriber.RegisterTrigger(testutils.Context(t), req)
				require.NoError(t, err)
			}

			// Reset counters after initial registration sends, then start
			// the loop so the first tick generates a clean measurement.
			dispatcher.Reset()
			require.NoError(t, subscriber.Start(testutils.Context(t)))
			t.Cleanup(func() { subscriber.Close() })

			// Wait for exactly one tick (500ms refresh, wait 700ms to give time)
			expectedSends := int64(tc.nRegistrations * tc.capDonSize)
			deadline := time.Now().Add(2 * time.Second)
			for time.Now().Before(deadline) {
				if dispatcher.Count(remotetypes.MethodRegisterTrigger) >= expectedSends {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}

			snap := dispatcher.Snapshot()
			registerCount := snap[remotetypes.MethodRegisterTrigger]

			t.Logf("Registrations: %d, CapDON peers: %d", tc.nRegistrations, tc.capDonSize)
			t.Logf("Expected P2P sends per tick: %d, Actual: %d", expectedSends, registerCount)
			t.Logf("Full snapshot: %v", snap)

			// At least one tick's worth of sends must have occurred.
			require.GreaterOrEqual(t, registerCount, expectedSends,
				"subscriber registrationLoop should send at least N*capDonMembers RegisterTrigger messages per tick")
			// Verify it's a multiple of expectedSends (each tick sends exactly this many).
			require.Equal(t, int64(0), registerCount%expectedSends,
				"total sends should be an exact multiple of per-tick sends")
		})
	}
}

// --- Test: Publisher sendRegistrationChecks traffic volume ---

func TestRegistrationCheckTrafficVolume(t *testing.T) {
	cases := []struct {
		nRegistrations  int
		capDonSize      int
		workflowDonSize int
	}{
		{100, 4, 7},
		{500, 4, 7},
		{1000, 4, 7},
		{2500, 4, 7},
		{5000, 4, 7},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("N=%d", tc.nRegistrations), func(t *testing.T) {
			ctx := testutils.Context(t)
			lggr := logger.Test(t)
			dispatcher := newCountingDispatcher()

			capDonMembers := generatePeers(t, tc.capDonSize)
			capDon := commoncap.DON{ID: 1, Members: capDonMembers, F: donFaultTolerance(t, tc.capDonSize/3)}

			workflowDonMembers := generatePeers(t, tc.workflowDonSize)
			workflowDon := commoncap.DON{ID: 2, Members: workflowDonMembers, F: donFaultTolerance(t, (tc.workflowDonSize-1)/3)}

			capInfo := commoncap.CapabilityInfo{
				ID:             "cap_id@1",
				CapabilityType: commoncap.CapabilityTypeTrigger,
				Description:    "Load Test Trigger",
			}

			underlying := &noopTrigger{info: capInfo}

			// Use the same refresh interval before and after registration so
			// registrationCheckLoop's ticker matches the measurement window (a
			// long initial interval would not pick up a later shorter interval
			// until the next tick of the old ticker).
			cfg := &commoncap.RemoteTriggerConfig{
				RegistrationRefresh:     50 * time.Millisecond,
				RegistrationExpiry:      2 * time.Hour,
				MinResponsesToAggregate: 1,
				MessageExpiry:           time.Hour,
				MaxBatchSize:            100,
				BatchCollectionPeriod:   time.Second,
			}

			workflowDONs := map[uint32]commoncap.DON{workflowDon.ID: workflowDon}
			publisher := remote.NewTriggerPublisher(capInfo.ID, "LogTrigger", dispatcher, lggr)
			require.NoError(t, publisher.SetConfig(cfg, underlying, capDon, workflowDONs))
			require.NoError(t, publisher.Start(ctx))
			t.Cleanup(func() { require.NoError(t, publisher.Close()) })

			// Register N triggers by feeding MethodRegisterTrigger messages.
			// With F=0 on the workflow DON side used in newServices patterns,
			// a single sender suffices for quorum. Here F>0 so we need 2F+1
			// senders. Use all workflow DON members.
			for i := range tc.nRegistrations {
				triggerID := fmt.Sprintf("trigger_%d", i)
				wfID := generateWorkflowID(i)

				for _, sender := range workflowDonMembers {
					req := commoncap.TriggerRegistrationRequest{
						TriggerID: triggerID,
						Metadata:  commoncap.RequestMetadata{WorkflowID: wfID},
					}
					marshaled, err := pb.MarshalTriggerRegistrationRequest(req)
					require.NoError(t, err)

					msg := &remotetypes.MessageBody{
						Sender:      sender[:],
						Method:      remotetypes.MethodRegisterTrigger,
						CallerDonId: workflowDon.ID,
						Payload:     marshaled,
					}
					publisher.Receive(ctx, msg)
				}
			}

			dispatcher.Reset()

			// Count registration checks after registrations are in place (ignore
			// any check traffic during the registration burst above).
			time.Sleep(200 * time.Millisecond)

			snap := dispatcher.Snapshot()
			checkCount := snap[remotetypes.MethodTriggerRegistrationCheck]

			// With MaxBatchSize=100, each registration-check tick sends ceil(N/100) messages per
			// workflow-DON peer (chunked metadata), not one giant payload.
			chunksPerTick := (tc.nRegistrations + 99) / 100
			minPerFullTick := int64(chunksPerTick * tc.workflowDonSize)

			t.Logf("Registrations: %d, WorkflowDON peers: %d, chunksPerTick: %d", tc.nRegistrations, tc.workflowDonSize, chunksPerTick)
			t.Logf("Expected min (>= one tick): %d, Actual: %d", minPerFullTick, checkCount)
			t.Logf("Full snapshot: %v", snap)

			// Allow multiple ticks in the 200ms window (50ms refresh → several ticks).
			require.GreaterOrEqual(t, checkCount, minPerFullTick,
				"publisher should send at least chunksPerTick*peers TriggerRegistrationCheck messages per tick")
			// Total sends stay far below one message per registration (worst case would be N×peers).
			require.Less(t, checkCount, int64(tc.nRegistrations*tc.workflowDonSize),
				"registration check traffic must remain far below per-registration fan-out")
		})
	}
}

// --- Benchmark: Publisher processing rate for duplicate registrations ---

func BenchmarkRegistrationProcessing(b *testing.B) {
	ctx := context.Background()
	lggr, _ := logger.New()
	dispatcher := newCountingDispatcher()

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

	for i := 0; i < b.N; i++ {
		publisher.Receive(ctx, regMsg)
	}

	b.StopTimer()
	elapsed := b.Elapsed()
	rate := float64(b.N) / elapsed.Seconds()
	b.ReportMetric(rate, "msgs/sec")
}

// --- Test: Subscriber RLock hold duration during registration loop ---

func TestRegistrationLoopLockDuration(t *testing.T) {
	cases := []int{100, 500, 1000, 2500}

	for _, n := range cases {
		t.Run(fmt.Sprintf("N=%d", n), func(t *testing.T) {
			lggr := logger.Test(t)
			dispatcher := newCountingDispatcher()

			capDonMembers := generatePeers(t, 4)
			capDon := commoncap.DON{ID: 1, Members: capDonMembers, F: 1}

			workflowDonMembers := generatePeers(t, 7)
			workflowDon := commoncap.DON{ID: 2, Members: workflowDonMembers, F: 2}

			capInfo := commoncap.CapabilityInfo{
				ID:             "cap_id@1",
				CapabilityType: commoncap.CapabilityTypeTrigger,
				Description:    "Lock Duration Test",
			}

			cfg := &commoncap.RemoteTriggerConfig{
				RegistrationRefresh:     50 * time.Millisecond,
				RegistrationExpiry:      2 * time.Hour,
				MinResponsesToAggregate: 1,
				MessageExpiry:           time.Hour,
			}

			subscriber := remote.NewTriggerSubscriber(capInfo.ID, "LogTrigger", dispatcher, lggr)
			agg := aggregation.NewDefaultModeAggregator(cfg.MinResponsesToAggregate)
			require.NoError(t, subscriber.SetConfig(cfg, capInfo, workflowDon.ID, capDon, agg))

			for i := range n {
				req := commoncap.TriggerRegistrationRequest{
					TriggerID: fmt.Sprintf("trigger_%d", i),
					Metadata:  commoncap.RequestMetadata{WorkflowID: generateWorkflowID(i)},
				}
				_, err := subscriber.RegisterTrigger(testutils.Context(t), req)
				require.NoError(t, err)
			}

			dispatcher.Reset()

			expectedSends := int64(n * 4) // n registrations * 4 cap DON members

			require.NoError(t, subscriber.Start(testutils.Context(t)))
			t.Cleanup(func() { subscriber.Close() })

			start := time.Now()

			// Wait for the tick to complete — poll until we see the expected sends.
			deadline := start.Add(10 * time.Second)
			for time.Now().Before(deadline) {
				if dispatcher.Count(remotetypes.MethodRegisterTrigger) >= expectedSends {
					break
				}
				time.Sleep(5 * time.Millisecond)
			}
			elapsed := time.Since(start)

			actual := dispatcher.Count(remotetypes.MethodRegisterTrigger)
			require.GreaterOrEqual(t, actual, expectedSends, "should have sent all registrations")

			t.Logf("N=%d: %d P2P sends completed in %v (%.0f sends/sec)",
				n, actual, elapsed, float64(actual)/elapsed.Seconds())

			// Sanity: with a noop dispatcher, even 10,000 sends should
			// complete in well under 1 second.
			require.Less(t, elapsed, 5*time.Second,
				"registration loop tick should complete in reasonable time")
		})
	}
}

// TestTrafficAttribution_RegisterLoopVsChecksVsEventsAndAcks documents where
// outbound P2P traffic is generated: subscriber registrationLoop vs publisher
// registration checks vs one trigger event dispatch vs subscriber ACK fan-out.
// It uses the same topology as staging analysis (4 cap DON peers, 7 workflow
// DON peers, workflow F=2 → quorum 5). This does not simulate dispatcher drops
// or shared-channel saturation; see test log for that limitation.
func TestTrafficAttribution_RegisterLoopVsChecksVsEventsAndAcks(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)

	const (
		nRegistrations      = 120 // subscriber registrationLoop load (keep moderate for CI log volume)
		nPublisherRegs      = 80  // publisher-side registrations (quorum cost only)
		capDonPeerCount     = 4
		wfDonPeerCount      = 7
		wfDonFaultTolerance = 2 // minRequired = 2*F+1 = 5 for reg aggregation and ACK quorum
	)

	capPeers := generatePeers(t, capDonPeerCount)
	capDon := commoncap.DON{ID: 1, Members: capPeers, F: 1}

	wfPeers := generatePeers(t, wfDonPeerCount)
	wfDon := commoncap.DON{ID: 2, Members: wfPeers, F: donFaultTolerance(t, wfDonFaultTolerance)}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Traffic attribution",
	}

	// --- Phase 1: subscriber — one registrationLoop tick (periodic refresh) ---
	subRegDisp := newCountingDispatcher()
	subCfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     500 * time.Millisecond,
		RegistrationExpiry:      2 * time.Hour,
		MinResponsesToAggregate: 1,
		MessageExpiry:           time.Hour,
	}
	sub := remote.NewTriggerSubscriber(capInfo.ID, "LogTrigger", subRegDisp, lggr)
	agg := aggregation.NewDefaultModeAggregator(subCfg.MinResponsesToAggregate)
	require.NoError(t, sub.SetConfig(subCfg, capInfo, wfDon.ID, capDon, agg))
	for i := range nRegistrations {
		req := commoncap.TriggerRegistrationRequest{
			TriggerID: fmt.Sprintf("trigger_%d", i),
			Metadata:  commoncap.RequestMetadata{WorkflowID: generateWorkflowID(i)},
		}
		_, err := sub.RegisterTrigger(ctx, req)
		require.NoError(t, err)
	}
	subRegDisp.Reset()
	require.NoError(t, sub.Start(ctx))
	t.Cleanup(func() { require.NoError(t, sub.Close()) })

	perTickReg := int64(nRegistrations * capDonPeerCount)
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if subRegDisp.Count(remotetypes.MethodRegisterTrigger) >= perTickReg {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	phase1Reg := subRegDisp.Count(remotetypes.MethodRegisterTrigger)
	require.GreaterOrEqual(t, phase1Reg, perTickReg)
	require.Equal(t, int64(0), phase1Reg%perTickReg)

	// --- Phase 2: subscriber — one engine round-trip: deliver event + AckEvent (ACK fan-out to cap DON) ---
	subAckDisp := newCountingDispatcher()
	sub2 := remote.NewTriggerSubscriber(capInfo.ID, "LogTrigger", subAckDisp, lggr)
	require.NoError(t, sub2.SetConfig(subCfg, capInfo, wfDon.ID, capDon, agg))
	require.NoError(t, sub2.Start(ctx))
	t.Cleanup(func() { require.NoError(t, sub2.Close()) })

	const soloTrig = "solo-trigger"
	regReq := commoncap.TriggerRegistrationRequest{
		TriggerID: soloTrig,
		Metadata:  commoncap.RequestMetadata{WorkflowID: workflowID1},
	}
	_, err := sub2.RegisterTrigger(ctx, regReq)
	require.NoError(t, err)
	subAckDisp.Reset()

	ev := buildTriggerEventWithTriggerID(t, capPeers[0][:], workflowID1, soloTrig, "event-attr-1")
	sub2.Receive(ctx, ev)
	require.NoError(t, sub2.AckEvent(ctx, soloTrig, "event-attr-1", "LogTrigger"))

	ackFanout := subAckDisp.Count(remotetypes.MethodTriggerEventAck)
	require.Equal(t, int64(capDonPeerCount), ackFanout,
		"subscriber AckEvent should send one TriggerEventAck to each cap DON peer")

	// --- Phase 3: publisher — registration checks over a short window ---
	// Use noopTrigger here: multiTrigger's registrationsCh (buffer 10) fills and
	// blocks RegisterTrigger after 10 registrations if not drained.
	pubDisp := newCountingDispatcher()
	underlyingNoop := &noopTrigger{info: capInfo}
	pubCfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     50 * time.Millisecond,
		RegistrationExpiry:      2 * time.Hour,
		MinResponsesToAggregate: 1,
		MessageExpiry:           time.Hour,
		MaxBatchSize:            100,
		BatchCollectionPeriod:   time.Second,
	}
	wfDONs := map[uint32]commoncap.DON{wfDon.ID: wfDon}
	pub := remote.NewTriggerPublisher(capInfo.ID, "LogTrigger", pubDisp, lggr)
	require.NoError(t, pub.SetConfig(pubCfg, underlyingNoop, capDon, wfDONs))
	require.NoError(t, pub.Start(ctx))
	t.Cleanup(func() { require.NoError(t, pub.Close()) })

	minRegSenders := 2*int(wfDon.F) + 1
	for i := range nPublisherRegs {
		tid := fmt.Sprintf("pub_trig_%d", i)
		wid := generateWorkflowID(10000 + i)
		for j := range minRegSenders {
			req := commoncap.TriggerRegistrationRequest{
				TriggerID: tid,
				Metadata:  commoncap.RequestMetadata{WorkflowID: wid},
			}
			payload, mErr := pb.MarshalTriggerRegistrationRequest(req)
			require.NoError(t, mErr)
			pub.Receive(ctx, &remotetypes.MessageBody{
				Sender:      wfPeers[j][:],
				Method:      remotetypes.MethodRegisterTrigger,
				CallerDonId: wfDon.ID,
				Payload:     payload,
			})
		}
	}

	pubDisp.Reset()

	time.Sleep(200 * time.Millisecond)
	chunksPerPublisherTick := (nPublisherRegs + int(pubCfg.MaxBatchSize) - 1) / int(pubCfg.MaxBatchSize)
	pubCheckSendsPerTick := int64(chunksPerPublisherTick * wfDonPeerCount)

	phase3Checks := pubDisp.Count(remotetypes.MethodTriggerRegistrationCheck)
	require.GreaterOrEqual(t, phase3Checks, pubCheckSendsPerTick,
		"expect at least one registration-check tick (chunked by MaxBatchSize)")
	// Measurement spans multiple RegistrationRefresh ticks — total sends grow with ticks,
	// but stay far below naive per-registration × peers fan-out.
	require.Less(t, phase3Checks, int64(nPublisherRegs*wfDonPeerCount),
		"registration check traffic must stay below per-registration × workflow peers")

	// --- Phase 4: publisher — one underlying trigger event + workflow ACK quorum ---
	pubEvDisp := newCountingDispatcher()
	underEv := newMultiTrigger(capInfo)
	pub2 := remote.NewTriggerPublisher(capInfo.ID, "LogTrigger", pubEvDisp, lggr)
	evCfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     time.Hour,
		RegistrationExpiry:      2 * time.Hour,
		MinResponsesToAggregate: 1,
		MessageExpiry:           time.Hour,
		// MaxBatchSize must be 1 (or BatchCollectionPeriod < 10ms) so trigger-event batching
		// is disabled; otherwise events wait for the batch window and this phase sees 0 sends.
		MaxBatchSize:          1,
		BatchCollectionPeriod: time.Second,
	}
	require.NoError(t, pub2.SetConfig(evCfg, underEv, capDon, wfDONs))
	require.NoError(t, pub2.Start(ctx))
	t.Cleanup(func() { require.NoError(t, pub2.Close()) })

	const evTrig = "event-path-trigger"
	evWid := workflowID1
	for j := range minRegSenders {
		req := commoncap.TriggerRegistrationRequest{
			TriggerID: evTrig,
			Metadata:  commoncap.RequestMetadata{WorkflowID: evWid},
		}
		payload, mErr := pb.MarshalTriggerRegistrationRequest(req)
		require.NoError(t, mErr)
		pub2.Receive(ctx, &remotetypes.MessageBody{
			Sender:      wfPeers[j][:],
			Method:      remotetypes.MethodRegisterTrigger,
			CallerDonId: wfDon.ID,
			Payload:     payload,
		})
	}
	<-underEv.registrationsCh // wait for publisher to register underlying trigger

	pubEvDisp.Reset()
	underEv.SendEvent(evTrig, commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: "publisher-ev-1"},
	})
	time.Sleep(50 * time.Millisecond)
	afterEvent := pubEvDisp.Count(remotetypes.MethodTriggerEvent)
	require.Equal(t, int64(wfDonPeerCount), afterEvent,
		"one non-batched trigger event should be dispatched to each workflow DON peer")

	beforeAck := pubEvDisp.Snapshot()
	for j := range minRegSenders {
		pub2.Receive(ctx, newAckEventMessage(t, "publisher-ev-1", evTrig, wfDon.ID, wfPeers[j]))
	}
	time.Sleep(20 * time.Millisecond)
	afterAck := pubEvDisp.Snapshot()

	// ACK handling does not add outbound P2P sends on the publisher path.
	for m, v := range afterAck {
		if m == remotetypes.MethodTriggerEvent {
			require.Equal(t, beforeAck[m], v, "ACK quorum should not cause extra trigger event sends")
		}
	}

	// --- Evidence summary for staging / team discussion ---
	t.Logf("=== P2P traffic attribution (same topology: %d cap peers, %d wf peers, F=%d) ===",
		capDonPeerCount, wfDonPeerCount, wfDonFaultTolerance)
	t.Logf("Phase 1 — subscriber registrationLoop (one tick, N=%d): RegisterTrigger total=%d (expected per tick: %d = N×capPeers)",
		nRegistrations, phase1Reg, perTickReg)
	t.Logf("Phase 2 — subscriber AckEvent after one event: TriggerEventAck total=%d (expected: %d = capPeers)",
		ackFanout, capDonPeerCount)
	t.Logf("Phase 3 — publisher sendRegistrationChecks (~200ms, N=%d, MaxBatchSize=%d): TriggerRegistrationCheck total=%d (min one tick ≈ %d = ceil(N/batch)×wfPeers; window spans multiple ticks)",
		nPublisherRegs, pubCfg.MaxBatchSize, phase3Checks, pubCheckSendsPerTick)
	t.Logf("Phase 4 — publisher one event then %d ACKs: TriggerEvent sends=%d (expected: %d = wfPeers); snapshot after ACK: %s",
		minRegSenders, afterEvent, wfDonPeerCount, formatP2PSnapshot(afterAck))

	// Compare single-tick orders of magnitude: phase1 is one subscriber refresh tick;
	// phase3 denominator is one publisher registration-check tick (chunked), not the
	// summed multi-tick window total (phase3Checks).
	regToCheckRatio := float64(perTickReg) / float64(pubCheckSendsPerTick)
	t.Logf("--- Ratios (illustrative; phase windows differ) ---")
	t.Logf("registrationLoop RegisterTrigger (one tick) / publisher registration-check sends per tick ≈ %.1fx",
		regToCheckRatio)
	ackDenom := ackFanout
	if ackDenom < 1 {
		ackDenom = 1
	}
	t.Logf("registrationLoop per tick / one subscriber AckEvent round ≈ %.1fx",
		float64(perTickReg)/float64(ackDenom))
	evDenom := afterEvent
	if evDenom < 1 {
		evDenom = 1
	}
	t.Logf("registrationLoop per tick / one trigger event dispatch ≈ %.1fx",
		float64(perTickReg)/float64(evDenom))

	require.Greater(t, regToCheckRatio, 10.0,
		"evidence: periodic registration refresh traffic should dwarf registration-check traffic")

	t.Logf("NOTE: this test uses instant no-op dispatchers (no drops). To evidence drop/retry " +
		"amplification, add a bounded-queue or failing Send mock in a separate test.")
}
