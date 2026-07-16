package streams

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"

	"github.com/smartcontractkit/chainlink-data-streams/mercury/v3/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	mockTriggerCapabilityName           = "mock-streams-trigger"
	mockTriggerCapabilityVersion        = "1.0.0"
	mockTriggerRegisterResolution       = 100
	defaultLoopIntervalMs         int64 = 1000
	mockSignerFaultTolerance            = 1
	mockSignerKeyLength                 = 32
	mockSignerKeyLastByteIndex          = mockSignerKeyLength - 1
	reportExpiryOffsetSeconds           = 1_000_000
	initialMockPriceA             int64 = 300_000
	initialMockPriceB             int64 = 40_000
	initialMockPriceC             int64 = 5_000_000
)

func RegisterMockTrigger(lggr logger.Logger, capRegistry core.CapabilitiesRegistry) (*MockTriggerService, error) {
	ctx := context.Background()
	trigger, err := NewMockTriggerService(mockTriggerRegisterResolution, lggr)
	if err != nil {
		return nil, err
	}
	if err := trigger.Start(ctx); err != nil {
		return nil, err
	}
	if err := capRegistry.Add(ctx, trigger); err != nil {
		_ = trigger.Close()
		return nil, err
	}

	return trigger, nil
}

const MockTriggerCapabilityID = mockTriggerCapabilityName + "@" + mockTriggerCapabilityVersion

var capInfo = capabilities.MustNewCapabilityInfo(
	MockTriggerCapabilityID,
	capabilities.CapabilityTypeTrigger,
	"Mock Streams Trigger",
)

// Wraps the MercuryTriggerService to produce a trigger with mocked data
type MockTriggerService struct {
	*triggers.MercuryTriggerService
	meta          datastreams.Metadata
	signers       []*ecdsa.PrivateKey
	stopCh        services.StopChan
	closeOnce     sync.Once
	wg            sync.WaitGroup
	loopInterval  time.Duration
	subscribers   map[string][]streams.FeedId
	subscribersMu sync.Mutex
	lggr          logger.Logger
}

func NewMockTriggerService(tickerResolutionMs int64, lggr logger.Logger) (*MockTriggerService, error) {
	trigger, err := triggers.NewMercuryTriggerService(tickerResolutionMs, mockTriggerCapabilityName, mockTriggerCapabilityVersion, lggr)
	if err != nil {
		return nil, err
	}
	trigger.CapabilityInfo = capInfo

	if tickerResolutionMs <= 0 {
		tickerResolutionMs = defaultLoopIntervalMs
	}

	meta, signers, err := newMockMetadataAndSigners()
	if err != nil {
		return nil, err
	}

	// MercuryTrigger is typically wrapped by other modules that ignore the trigger's meta and provide a different one.
	// Since we're skipping those wrappers we need to provide our own meta here.
	trigger.SetMetaOverride(meta)

	return &MockTriggerService{
		MercuryTriggerService: trigger,
		meta:                  meta,
		signers:               signers,
		stopCh:                make(services.StopChan),
		loopInterval:          time.Duration(tickerResolutionMs) * time.Millisecond,
		subscribers:           make(map[string][]streams.FeedId),
		lggr:                  lggr,
	}, nil
}

func newMockMetadataAndSigners() (datastreams.Metadata, []*ecdsa.PrivateKey, error) {
	meta := datastreams.Metadata{MinRequiredSignatures: 2*mockSignerFaultTolerance + 1}
	signers := make([]*ecdsa.PrivateKey, 0, meta.MinRequiredSignatures)
	for i := 0; i < meta.MinRequiredSignatures; i++ {
		privKey, err := newMockSigner(i + 1)
		if err != nil {
			return datastreams.Metadata{}, nil, err
		}
		signers = append(signers, privKey)
		meta.Signers = append(meta.Signers, crypto.PubkeyToAddress(privKey.PublicKey).Bytes())
	}
	return meta, signers, nil
}

func newMockSigner(index int) (*ecdsa.PrivateKey, error) {
	bytes := make([]byte, mockSignerKeyLength)
	lastByte, err := toUint8(index)
	if err != nil {
		return nil, err
	}
	bytes[mockSignerKeyLastByteIndex] = lastByte
	privKey, err := crypto.ToECDSA(bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func (m *MockTriggerService) Start(ctx context.Context) error {
	if err := m.MercuryTriggerService.Start(ctx); err != nil {
		return err
	}
	m.wg.Add(1)
	go m.loop()
	return nil
}

func (m *MockTriggerService) Close() error {
	m.closeOnce.Do(func() {
		close(m.stopCh)
	})
	m.wg.Wait()
	return m.MercuryTriggerService.Close()
}

func (m *MockTriggerService) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	ch, err := m.MercuryTriggerService.RegisterTrigger(ctx, req)
	if err != nil {
		return nil, err
	}

	config, err := m.ValidateConfig(req.Config)
	if err != nil {
		_ = m.MercuryTriggerService.UnregisterTrigger(ctx, req)
		return nil, err
	}
	m.subscribersMu.Lock()
	defer m.subscribersMu.Unlock()
	m.subscribers[req.Metadata.WorkflowID] = config.FeedIds
	return ch, nil
}

func (m *MockTriggerService) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	err := m.MercuryTriggerService.UnregisterTrigger(ctx, req)
	if err != nil {
		return err
	}

	m.subscribersMu.Lock()
	defer m.subscribersMu.Unlock()
	delete(m.subscribers, req.Metadata.WorkflowID)
	return nil
}

const baseTimestamp = 1000000000

// NOTE: duplicated from codec_test.go
func newReport(lggr logger.Logger, feedID [32]byte, price *big.Int, timestamp int64) ([]byte, error) {
	uintTimestamp, err := toUint32(timestamp)
	if err != nil {
		return nil, err
	}
	expiresAt, err := toUint32(timestamp + reportExpiryOffsetSeconds)
	if err != nil {
		return nil, err
	}
	v3Codec := reportcodec.NewReportCodec(feedID, lggr)
	raw, err := v3Codec.BuildReport(context.Background(), v3.ReportFields{
		BenchmarkPrice:     price,
		Timestamp:          uintTimestamp,
		ValidFromTimestamp: uintTimestamp,
		Bid:                price,
		Ask:                price,
		LinkFee:            price,
		NativeFee:          price,
		ExpiresAt:          expiresAt,
	})
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func toUint32(v int64) (uint32, error) {
	if v < 0 || v > int64(^uint32(0)) {
		return 0, fmt.Errorf("value %d out of uint32 range", v)
	}
	return uint32(v), nil
}

func toUint8(v int) (uint8, error) {
	if v < 0 || v > int(^uint8(0)) {
		return 0, fmt.Errorf("value %d out of uint8 range", v)
	}
	return uint8(v), nil
}

func rawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := make([]byte, 0, len(rc)*32)
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}

func (m *MockTriggerService) loop() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.loopInterval)
	defer ticker.Stop()

	prices := []int64{initialMockPriceA, initialMockPriceB, initialMockPriceC}
	iteration := 0

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
		}

		incrementPrices(prices)
		iteration++

		timestamp := time.Now().Unix()
		reportCtx, err := newReportContext(iteration)
		if err != nil {
			m.lggr.Errorw("failed to build Mock report context", "err", err, "timestamp", timestamp)
			continue
		}
		reports, err := m.buildReports(timestamp, prices[0], reportCtx)
		if err != nil {
			m.lggr.Errorw("failed to build Mock reports", "err", err, "timestamp", timestamp)
			continue
		}
		if len(reports) == 0 {
			continue
		}

		m.lggr.Infow("New set of Mock reports", "timestamp", timestamp, "payload", reports)
		if err := m.ProcessReport(reports); err != nil {
			m.lggr.Errorw("failed to process Mock reports", "err", err, "timestamp", timestamp, "payload", reports)
		}
	}
}

func incrementPrices(prices []int64) {
	for i := range prices {
		prices[i]++
	}
}

func newReportContext(iteration int) (ocrTypes.ReportContext, error) {
	epoch, err := toUint32(int64(baseTimestamp + iteration))
	if err != nil {
		return ocrTypes.ReportContext{}, err
	}
	return ocrTypes.ReportContext{
		ReportTimestamp: ocrTypes.ReportTimestamp{Epoch: epoch},
	}, nil
}

func (m *MockTriggerService) buildReports(timestamp, price int64, reportCtx ocrTypes.ReportContext) ([]datastreams.FeedReport, error) {
	subscribers := m.snapshotSubscribers()
	reports := make([]datastreams.FeedReport, 0, subscriberCount(subscribers))
	for _, feedIDs := range subscribers {
		for _, feedID := range feedIDs {
			report, err := m.newSignedReport(string(feedID), price, timestamp, reportCtx)
			if err != nil {
				return nil, err
			}
			reports = append(reports, report)
		}
	}
	return reports, nil
}

func (m *MockTriggerService) snapshotSubscribers() map[string][]streams.FeedId {
	m.subscribersMu.Lock()
	defer m.subscribersMu.Unlock()

	snapshot := make(map[string][]streams.FeedId, len(m.subscribers))
	for workflowID, feedIDs := range m.subscribers {
		snapshot[workflowID] = cloneFeedIDs(feedIDs)
	}
	return snapshot
}

func subscriberCount(subscribers map[string][]streams.FeedId) int {
	total := 0
	for _, feedIDs := range subscribers {
		total += len(feedIDs)
	}
	return total
}

func cloneFeedIDs(feedIDs []streams.FeedId) []streams.FeedId {
	cloned := make([]streams.FeedId, len(feedIDs))
	copy(cloned, feedIDs)
	return cloned
}

func (m *MockTriggerService) newSignedReport(feedID string, price, timestamp int64, reportCtx ocrTypes.ReportContext) (datastreams.FeedReport, error) {
	fullReport, err := newReport(m.lggr, common.HexToHash(feedID), big.NewInt(price), timestamp)
	if err != nil {
		return datastreams.FeedReport{}, fmt.Errorf("build report for feed %s: %w", feedID, err)
	}

	report := datastreams.FeedReport{
		FeedID:               feedID,
		FullReport:           fullReport,
		ReportContext:        rawReportContext(reportCtx),
		ObservationTimestamp: timestamp,
	}

	if err := m.signReport(&report); err != nil {
		return datastreams.FeedReport{}, fmt.Errorf("sign report for feed %s: %w", feedID, err)
	}

	return report, nil
}

func (m *MockTriggerService) signReport(report *datastreams.FeedReport) error {
	sigData := append(crypto.Keccak256(report.FullReport), report.ReportContext...)
	hash := crypto.Keccak256(sigData)
	report.Signatures = make([][]byte, 0, m.meta.MinRequiredSignatures)
	for n := 0; n < m.meta.MinRequiredSignatures; n++ {
		sig, err := crypto.Sign(hash, m.signers[n])
		if err != nil {
			return err
		}
		report.Signatures = append(report.Signatures, sig)
	}
	return nil
}
