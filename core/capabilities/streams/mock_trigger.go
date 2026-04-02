package streams

import (
	"context"
	"crypto/ecdsa"
	"maps"
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

	"github.com/smartcontractkit/chainlink-evm/pkg/mercury/v3/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func RegisterMockTrigger(lggr logger.Logger, capRegistry core.CapabilitiesRegistry) (*MockTriggerService, error) {
	ctx := context.TODO()
	trigger := NewMockTriggerService(100, lggr)
	if err := trigger.Start(ctx); err != nil {
		return nil, err
	}
	if err := capRegistry.Add(ctx, trigger); err != nil {
		return nil, err
	}

	return trigger, nil
}

const MockTriggerCapabilityID = "mock-streams-trigger@1.0.0"

const triggerID = MockTriggerCapabilityID

var capInfo = capabilities.MustNewCapabilityInfo(
	triggerID,
	capabilities.CapabilityTypeTrigger,
	"Mock Streams Trigger",
)

// Wraps the MercuryTriggerService to produce a trigger with mocked data
type MockTriggerService struct {
	*triggers.MercuryTriggerService
	meta          datastreams.Metadata
	signers       []*ecdsa.PrivateKey
	stopCh        services.StopChan
	wg            sync.WaitGroup
	loopInterval  time.Duration
	subscribers   map[string][]streams.FeedId
	subscribersMu sync.Mutex
	lggr          logger.Logger
}

func NewMockTriggerService(tickerResolutionMs int64, lggr logger.Logger) *MockTriggerService {
	trigger, err := triggers.NewMercuryTriggerService(tickerResolutionMs, "mock-streams-trigger", "1.0.0", lggr)
	if err != nil {
		panic(err)
	}
	trigger.CapabilityInfo = capInfo

	if tickerResolutionMs == 0 {
		tickerResolutionMs = 1000
	}

	f := 1
	meta := datastreams.Metadata{MinRequiredSignatures: 2*f + 1}
	// gen private keys for MinRequiredSignatures
	signers := []*ecdsa.PrivateKey{}
	for i := 0; i < meta.MinRequiredSignatures; i++ {
		// test keys: need to be the same across nodes
		bytes := make([]byte, 32)
		bytes[31] = mustUint8(i + 1)

		privKey, err := crypto.ToECDSA(bytes)
		if err != nil {
			panic(err)
		}
		signers = append(signers, privKey)

		signerAddr := crypto.PubkeyToAddress(privKey.PublicKey).Bytes()
		meta.Signers = append(meta.Signers, signerAddr)
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
	}
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
	close(m.stopCh)
	m.wg.Wait()
	return m.MercuryTriggerService.Close()
}

func (m *MockTriggerService) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	ch, err := m.MercuryTriggerService.RegisterTrigger(ctx, req)
	if err != nil {
		return nil, err
	}

	config, _ := m.ValidateConfig(req.Config)
	m.subscribersMu.Lock()
	defer m.subscribersMu.Unlock()
	m.subscribers[req.Metadata.WorkflowID] = config.FeedIds
	return ch, nil
}

func (m *MockTriggerService) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	err := m.MercuryTriggerService.UnregisterTrigger(ctx, req)
	m.subscribersMu.Lock()
	defer m.subscribersMu.Unlock()
	delete(m.subscribers, req.Metadata.WorkflowID)
	return err
}

const baseTimestamp = 1000000000

// NOTE: duplicated from trigger_test.go
func newReport(lggr logger.Logger, feedID [32]byte, price *big.Int, timestamp int64) []byte {
	v3Codec := reportcodec.NewReportCodec(feedID, lggr)
	raw, err := v3Codec.BuildReport(context.Background(), v3.ReportFields{
		BenchmarkPrice:     price,
		Timestamp:          mustUint32(timestamp),
		ValidFromTimestamp: mustUint32(timestamp),
		Bid:                price,
		Ask:                price,
		LinkFee:            price,
		NativeFee:          price,
		ExpiresAt:          mustUint32(timestamp + 1000000),
	})
	if err != nil {
		panic(err)
	}
	return raw
}

func mustUint32(v int64) uint32 {
	if v < 0 || v > int64(^uint32(0)) {
		panic("timestamp out of uint32 range")
	}
	return uint32(v)
}

func mustUint8(v int) uint8 {
	if v < 0 || v > int(^uint8(0)) {
		panic("value out of uint8 range")
	}
	return uint8(v)
}

func rawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}

func (m *MockTriggerService) loop() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.loopInterval)
	defer ticker.Stop()

	prices := []int64{300000, 40000, 5000000}

	j := 0

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
		}

		// TODO: properly close
		for i := range prices {
			prices[i]++
		}
		j++

		// https://github.com/smartcontractkit/chainlink/blob/41f9428c3aa8231e8834a230fca4c2ccffd4e6c3/core/capabilities/streams/trigger_test.go#L117-L122

		timestamp := time.Now().Unix()
		// TODO: shouldn't we increment round rather than epoch?
		reportCtx := ocrTypes.ReportContext{ReportTimestamp: ocrTypes.ReportTimestamp{Epoch: mustUint32(int64(baseTimestamp + j))}}

		reports := []datastreams.FeedReport{}
		subscribers := map[string][]streams.FeedId{}
		m.subscribersMu.Lock()
		maps.Copy(subscribers, m.subscribers)
		m.subscribersMu.Unlock()
		for _, feedIDs := range subscribers {
			for _, feedID := range feedIDs {
				feedID := string(feedID)
				report := datastreams.FeedReport{
					FeedID:               feedID,
					FullReport:           newReport(m.lggr, common.HexToHash(feedID), big.NewInt(prices[0]), timestamp),
					ReportContext:        rawReportContext(reportCtx),
					ObservationTimestamp: timestamp,
				}
				// sign report with mock signers
				sigData := append(crypto.Keccak256(report.FullReport), report.ReportContext...)
				hash := crypto.Keccak256(sigData)
				for n := 0; n < m.meta.MinRequiredSignatures; n++ {
					sig, err := crypto.Sign(hash, m.signers[n])
					if err != nil {
						panic(err)
					}
					report.Signatures = append(report.Signatures, sig)
				}

				reports = append(reports, report)
			}
		}

		m.lggr.Infow("New set of Mock reports", "timestamp", time.Now().Unix(), "payload", reports)
		err := m.ProcessReport(reports)
		if err != nil {
			m.lggr.Errorw("failed to process Mock reports", "err", err, "timestamp", time.Now().Unix(), "payload", reports)
		}
	}
}
