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

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
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

const triggerID = "mock-streams-trigger@1.0.0"

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
	subscribers   map[string][]streams.FeedId
	subscribersMu sync.Mutex
	lggr          logger.Logger
}

func NewMockTriggerService(tickerResolutionMs int64, lggr logger.Logger) *MockTriggerService {
	trigger := triggers.NewMercuryTriggerService(tickerResolutionMs, lggr)
	trigger.CapabilityInfo = capInfo

	f := 1
	meta := datastreams.Metadata{MinRequiredSignatures: 2*f + 1}
	// gen private keys for MinRequiredSignatures
	signers := []*ecdsa.PrivateKey{}
	for i := 0; i < meta.MinRequiredSignatures; i++ {
		// test keys: need to be the same across nodes
		bytes := make([]byte, 32)
		bytes[31] = uint8(i + 1)

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
		subscribers:           make(map[string][]streams.FeedId),
		lggr:                  lggr}
}

func (m *MockTriggerService) Start(ctx context.Context) error {
	if err := m.MercuryTriggerService.Start(ctx); err != nil {
		return err
	}
	go m.loop()
	return nil
}

func (m *MockTriggerService) Close() error {
	close(m.stopCh)
	m.wg.Wait()
	return m.MercuryTriggerService.Close()
}

func (o *MockTriggerService) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	ch, err := o.MercuryTriggerService.RegisterTrigger(ctx, req)
	if err != nil {
		return nil, err
	}

	config, _ := o.MercuryTriggerService.ValidateConfig(req.Config)
	o.subscribersMu.Lock()
	defer o.subscribersMu.Unlock()
	o.subscribers[req.Metadata.WorkflowID] = config.FeedIds
	return ch, nil
}

func (o *MockTriggerService) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	err := o.MercuryTriggerService.UnregisterTrigger(ctx, req)
	o.subscribersMu.Lock()
	defer o.subscribersMu.Unlock()
	delete(o.subscribers, req.Metadata.WorkflowID)
	return err
}

const baseTimestamp = 1000000000

// NOTE: duplicated from trigger_test.go
func newReport(lggr logger.Logger, feedID [32]byte, price *big.Int, timestamp int64) []byte {
	ctx := context.Background()
	v3Codec := reportcodec.NewReportCodec(feedID, lggr)
	raw, err := v3Codec.BuildReport(ctx, v3.ReportFields{
		BenchmarkPrice:     price,
		Timestamp:          uint32(timestamp),
		ValidFromTimestamp: uint32(timestamp),
		Bid:                price,
		Ask:                price,
		LinkFee:            price,
		NativeFee:          price,
		ExpiresAt:          uint32(timestamp + 1000000),
	})
	if err != nil {
		panic(err)
	}
	return raw
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
	sleepSec := 15
	ticker := time.NewTicker(time.Duration(sleepSec) * time.Second)
	defer ticker.Stop()

	prices := []int64{300000, 40000, 5000000}

	j := 0

	for range ticker.C {
		// TODO: properly close
		for i := range prices {
			prices[i] = prices[i] + 1
		}
		j++

		// https://github.com/smartcontractkit/chainlink/blob/41f9428c3aa8231e8834a230fca4c2ccffd4e6c3/core/capabilities/streams/trigger_test.go#L117-L122

		timestamp := time.Now().Unix()
		// TODO: shouldn't we increment round rather than epoch?
		reportCtx := ocrTypes.ReportContext{ReportTimestamp: ocrTypes.ReportTimestamp{Epoch: uint32(baseTimestamp + j)}}

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
		err := m.MercuryTriggerService.ProcessReport(reports)
		if err != nil {
			m.lggr.Errorw("failed to process Mock reports", "err", err, "timestamp", time.Now().Unix(), "payload", reports)
		}
	}
}
