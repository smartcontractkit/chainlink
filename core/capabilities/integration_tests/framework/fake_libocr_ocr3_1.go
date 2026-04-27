package framework

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocrintegrationtesthelpers"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
)

// stubBlobBroadcastFetcher is the v1 placeholder for the integration test
// harness. It satisfies ocr3_1types.BlobBroadcastFetcher but never actually
// round-trips a payload.
//
// The blob path is intentionally out of scope for v1 of FakeLibOCR_OCR3_1:
// libocr's BlobHandle wraps an internal LightCertifiedBlob (signed merkle
// commitments live in libocr/.../internal/ocr3_1/blobtypes/), so the only
// way to get a wire-survivable handle is to round-trip through libocr's
// real blob transport. That's tracked as plan §3.7 preflight item 4.
//
// Plugins that hit Broadcast/Fetch in this simulator will panic, surfacing
// the limit clearly rather than silently producing handles that won't
// unmarshal downstream.
type stubBlobBroadcastFetcher struct{}

func (stubBlobBroadcastFetcher) BroadcastBlob(context.Context, []byte, ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	panic("FakeLibOCR_OCR3_1: BroadcastBlob not implemented in v1 — see plan §3.7 preflight item 4")
}

func (stubBlobBroadcastFetcher) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	panic("FakeLibOCR_OCR3_1: FetchBlob not implemented in v1 — see plan §3.7 preflight item 4")
}

// libocrNodeOCR3_1 mirrors libocrNode for the OCR3_1 plugin shape. Each node
// owns its own KeyValueDatabase (keyed by a per-node synthetic ConfigDigest)
// to model the real per-node KV state libocr provides at runtime.
type libocrNodeOCR3_1 struct {
	id     string
	plugin ocr3_1types.ReportingPlugin[[]byte]
	kvDB   ocr3_1types.KeyValueDatabase
	key    ocr2key.KeyBundle
}

// FakeLibOCR_OCR3_1 is a minimal fake libocr harness for the OCR3_1 plugin,
// parallel to FakeLibOCR.
//
// v1 scope: the simulator only drives Query. OCR3_1's
// Observation/ValidateObservation/StateTransition phases require a real
// BlobHandle round-trip, and BlobHandle has no public constructor (its
// concrete LightCertifiedBlob lives in a libocr internal package). Driving
// the full protocol cycle is therefore deferred until we adopt libocr's
// real blob transport.
//
// What v1 still validates end-to-end:
//   - DON.AddOCR3_1NonStandardCapability wiring
//   - chainlink-common's NewOCR3_1 / NewReportingPluginFactoryOCR3_1 path
//   - The OCR3_1 reporting plugin construction and capability registry
//     registration on each node
//
// Future versions should extend simulateProtocolRound to drive Observation,
// ValidateObservation, ObservationQuorum, StateTransition, Committed, and
// Reports once the blob handle round-trip is solved (plan §3.7 preflight 4).
type FakeLibOCR_OCR3_1 struct {
	services.StateMachine
	t    *testing.T
	lggr logger.Logger

	nodes                 []*libocrNodeOCR3_1
	f                     uint8
	protocolRoundInterval time.Duration

	seqNr uint64

	kvFactory *ocrintegrationtesthelpers.StatefulInMemoryKeyValueDatabaseFactory

	mux    sync.Mutex
	stopCh services.StopChan
	wg     sync.WaitGroup
}

func NewFakeLibOCR_OCR3_1(t *testing.T, lggr logger.Logger, f uint8, protocolRoundInterval time.Duration) *FakeLibOCR_OCR3_1 {
	return &FakeLibOCR_OCR3_1{
		t:                     t,
		lggr:                  logger.Named(lggr, "FakeLibOCR_OCR3_1"),
		f:                     f,
		protocolRoundInterval: protocolRoundInterval,
		seqNr:                 1,
		kvFactory:             ocrintegrationtesthelpers.NewStatefulInMemoryKeyValueDatabaseFactory(),
		stopCh:                make(services.StopChan),
	}
}

func (m *FakeLibOCR_OCR3_1) Start(ctx context.Context) error {
	return m.StartOnce("FakeLibOCR_OCR3_1", func() error {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			ticker := time.NewTicker(m.protocolRoundInterval)
			defer ticker.Stop()
			for {
				select {
				case <-m.stopCh:
					return
				case <-ticker.C:
					serviceCtx, cancel := m.stopCh.NewCtx()
					if err := m.simulateProtocolRound(serviceCtx); err != nil {
						m.lggr.Errorf("OCR3_1 simulating protocol round: %v", err)
					}
					cancel()
				}
			}
		}()
		return nil
	})
}

func (m *FakeLibOCR_OCR3_1) Close() error {
	return m.StopOnce("FakeLibOCR_OCR3_1", func() error {
		close(m.stopCh)
		m.wg.Wait()
		m.mux.Lock()
		defer m.mux.Unlock()
		for _, n := range m.nodes {
			if n.kvDB != nil {
				_ = n.kvDB.Close()
			}
		}
		return nil
	})
}

// AddNode registers an OCR3_1 reporting plugin under a fresh per-node KV
// database. The node id is returned so the caller can later RemoveNode it.
func (m *FakeLibOCR_OCR3_1) AddNode(plugin ocr3_1types.ReportingPlugin[[]byte], key ocr2key.KeyBundle) (string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	nodeID := uuid.New().String()
	var configDigest types.ConfigDigest
	copy(configDigest[:], nodeID)

	kvDB, err := m.kvFactory.NewKeyValueDatabase(configDigest)
	if err != nil {
		return "", fmt.Errorf("create KV database for node %s: %w", nodeID, err)
	}

	m.nodes = append(m.nodes, &libocrNodeOCR3_1{
		id:     nodeID,
		plugin: plugin,
		kvDB:   kvDB,
		key:    key,
	})
	return nodeID, nil
}

func (m *FakeLibOCR_OCR3_1) GetNodeCount() int {
	m.mux.Lock()
	defer m.mux.Unlock()
	return len(m.nodes)
}

// simulateProtocolRound is the v1 round-driver. It exercises Query only on a
// randomly-selected leader. Other plugin methods are skipped — see type-level
// doc.
func (m *FakeLibOCR_OCR3_1) simulateProtocolRound(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if len(m.nodes) == 0 {
		return nil
	}

	leader := m.nodes[rand.Intn(len(m.nodes))]

	rdrTx, err := leader.kvDB.NewReadTransaction()
	if err != nil {
		return fmt.Errorf("open KV read tx: %w", err)
	}
	defer rdrTx.Discard()

	if _, err := leader.plugin.Query(ctx, m.seqNr, rdrTx, stubBlobBroadcastFetcher{}); err != nil {
		return fmt.Errorf("Query: %w", err)
	}

	m.seqNr++
	return nil
}
