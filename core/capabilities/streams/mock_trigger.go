package streams

// NOTE: this file is an amalgamation of MercuryTrigger and the streams trigger load tests
// the mercury trigger was modified to contain non-empty meta and sign the report with mock keys

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

const (
	baseTimestamp = 1000000000
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

	producer := NewMockDataProducer(trigger, lggr)
	if err := producer.Start(ctx); err != nil {
		return nil, err
	}

	return trigger, nil
}

// NOTE: duplicated from trigger_test.go
func newReport(lggr logger.Logger, feedID [32]byte, price *big.Int, timestamp int64) []byte {
	v3Codec := reportcodec.NewReportCodec(feedID, lggr)
	raw, err := v3Codec.BuildReport(v3.ReportFields{
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

const triggerID = "mock-streams-trigger@1.0.0"

var capInfo = capabilities.MustNewCapabilityInfo(
	triggerID,
	capabilities.CapabilityTypeTrigger,
	"Mock Streams Trigger",
)

const defaultTickerResolutionMs = 1000

// TODO pending capabilities configuration implementation - this should be configurable with a sensible default
const defaultSendChannelBufferSize = 1000

type config struct {
	// strings should be hex-encoded 32-byte values, prefixed with "0x", all lowercase, minimum 1 item
	FeedIDs []string `json:"feedIds" jsonschema:"pattern=^0x[0-9a-f]{64}$,minItems=1"`
	// must be greater than 0
	MaxFrequencyMs int `json:"maxFrequencyMs" jsonschema:"minimum=1"`
}

type inputs struct {
	TriggerID string `json:"triggerId"`
}

var mercuryTriggerValidator = capabilities.NewValidator[config, inputs, capabilities.TriggerEvent](capabilities.ValidatorArgs{Info: capInfo})

// This Trigger Service allows for the registration and deregistration of triggers. You can also send reports to the service.
type MockTriggerService struct {
	capabilities.Validator[config, inputs, capabilities.TriggerEvent]
	capabilities.CapabilityInfo
	tickerResolutionMs int64
	subscribers        map[string]*subscriber
	latestReports      map[datastreams.FeedID]datastreams.FeedReport
	mu                 sync.Mutex
	stopCh             services.StopChan
	wg                 sync.WaitGroup
	lggr               logger.Logger

	//
	meta    datastreams.SignersMetadata
	signers []*ecdsa.PrivateKey
	//
}

var _ capabilities.TriggerCapability = (*MockTriggerService)(nil)
var _ services.Service = &MockTriggerService{}

type subscriber struct {
	ch         chan<- capabilities.CapabilityResponse
	workflowID string
	config     config
}

// Mock Trigger will send events to each subscriber every MaxFrequencyMs (configurable per subscriber).
// Event generation happens whenever local unix time is a multiple of tickerResolutionMs. Therefore,
// all subscribers' MaxFrequencyMs values need to be a multiple of tickerResolutionMs.
func NewMockTriggerService(tickerResolutionMs int64, lggr logger.Logger) *MockTriggerService {
	if tickerResolutionMs == 0 {
		tickerResolutionMs = defaultTickerResolutionMs
	}
	//
	f := 1
	meta := datastreams.SignersMetadata{MinRequiredSignatures: 2*f + 1}
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
	//
	return &MockTriggerService{
		Validator:          mercuryTriggerValidator,
		CapabilityInfo:     capInfo,
		tickerResolutionMs: tickerResolutionMs,
		subscribers:        make(map[string]*subscriber),
		latestReports:      make(map[datastreams.FeedID]datastreams.FeedReport),
		stopCh:             make(services.StopChan),
		lggr:               lggr.Named("MockTriggerService"),
		meta:               meta,
		signers:            signers}
}

func (o *MockTriggerService) ProcessReport(reports []datastreams.FeedReport) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.lggr.Debugw("ProcessReport", "nReports", len(reports))
	for _, report := range reports {
		feedID := datastreams.FeedID(report.FeedID)
		o.latestReports[feedID] = report
	}
	return nil
}

func (o *MockTriggerService) RegisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	config, err := o.ValidateConfig(req.Config)
	if err != nil {
		return nil, err
	}

	inputs, err := o.ValidateInputs(req.Inputs)
	if err != nil {
		return nil, err
	}

	triggerID := o.getTriggerID(inputs.TriggerID, wid)
	// If triggerId is already registered, return an error
	if _, ok := o.subscribers[triggerID]; ok {
		return nil, fmt.Errorf("triggerId %s already registered", triggerID)
	}

	if int64(config.MaxFrequencyMs)%o.tickerResolutionMs != 0 {
		return nil, fmt.Errorf("MaxFrequencyMs must be a multiple of %d", o.tickerResolutionMs)
	}

	ch := make(chan capabilities.CapabilityResponse, defaultSendChannelBufferSize)
	o.subscribers[triggerID] =
		&subscriber{
			ch:         ch,
			workflowID: wid,
			config:     *config,
		}
	return ch, nil
}

func (o *MockTriggerService) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	inputs, err := o.ValidateInputs(req.Inputs)
	if err != nil {
		return err
	}
	triggerID := o.getTriggerID(inputs.TriggerID, wid)

	subscriber, ok := o.subscribers[triggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", triggerID)
	}
	close(subscriber.ch)
	delete(o.subscribers, triggerID)
	return nil
}

func (o *MockTriggerService) getTriggerID(triggerID string, wid string) string {
	tid := wid + "|" + triggerID
	return tid
}

func (o *MockTriggerService) loop() {
	defer o.wg.Done()
	now := time.Now().UnixMilli()
	nextWait := o.tickerResolutionMs - now%o.tickerResolutionMs

	for {
		select {
		case <-o.stopCh:
			return
		case <-time.After(time.Duration(nextWait) * time.Millisecond):
			startTs := time.Now().UnixMilli()
			// find closest timestamp that is a multiple of o.tickerResolutionMs
			aligned := (startTs + o.tickerResolutionMs/2) / o.tickerResolutionMs * o.tickerResolutionMs
			o.process(aligned)
			endTs := time.Now().UnixMilli()
			if endTs-startTs > o.tickerResolutionMs {
				o.lggr.Errorw("processing took longer than ticker resolution", "duration", endTs-startTs, "tickerResolutionMs", o.tickerResolutionMs)
			}
			nextWait = getNextWaitIntervalMs(aligned, o.tickerResolutionMs, endTs)
		}
	}
}

func getNextWaitIntervalMs(lastTs, tickerResolutionMs, currentTs int64) int64 {
	desiredNext := lastTs + tickerResolutionMs
	nextWait := desiredNext - currentTs
	if nextWait <= 0 {
		nextWait = 0
	}
	return nextWait
}

func (o *MockTriggerService) process(timestamp int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, sub := range o.subscribers {
		if timestamp%int64(sub.config.MaxFrequencyMs) == 0 {
			reportList := make([]datastreams.FeedReport, 0)
			for _, feedID := range sub.config.FeedIDs {
				if latest, ok := o.latestReports[datastreams.FeedID(feedID)]; ok {
					reportList = append(reportList, latest)
				}
			}

			// use 32-byte-padded timestamp as EventID (human-readable)
			eventID := fmt.Sprintf("streams_%024s", strconv.FormatInt(timestamp, 10))
			// ---
			// sign reports with mock signers
			for i := range reportList {
				report := reportList[i]
				sigData := append(crypto.Keccak256(report.FullReport), report.ReportContext...)
				hash := crypto.Keccak256(sigData)
				for n := 0; n < o.meta.MinRequiredSignatures; n++ {
					sig, err := crypto.Sign(hash, o.signers[n])
					if err != nil {
						panic(err)
					}
					reportList[i].Signatures = append(reportList[i].Signatures, sig)
				}
			}
			// ---
			capabilityResponse, err := wrapReports(reportList, eventID, timestamp, o.meta)

			if err != nil {
				o.lggr.Errorw("error wrapping reports", "err", err)
				continue
			}

			o.lggr.Debugw("ProcessReport pushing event", "nReports", len(reportList), "eventID", eventID)
			select {
			case sub.ch <- capabilityResponse:
			default:
				o.lggr.Errorw("subscriber channel full, dropping event", "eventID", eventID, "workflowID", sub.workflowID)
			}
		}
	}
}

func wrapReports(reportList []datastreams.FeedReport, eventID string, timestamp int64, meta datastreams.SignersMetadata) (capabilities.CapabilityResponse, error) {
	val, err := values.Wrap(reportList)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	metaVal, err := values.Wrap(meta)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	triggerEvent := capabilities.TriggerEvent{
		TriggerType: triggerID,
		ID:          eventID,
		Timestamp:   strconv.FormatInt(timestamp, 10),
		Metadata:    metaVal,
		Payload:     val,
	}

	eventVal, err := values.Wrap(triggerEvent)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	// Create a new CapabilityResponse with the MockTriggerEvent
	return capabilities.CapabilityResponse{
		Value: eventVal.(*values.Map),
	}, nil
}

func (o *MockTriggerService) Start(ctx context.Context) error {
	o.wg.Add(1)
	go o.loop()
	o.lggr.Info("MockTriggerService started")
	return nil
}

func (o *MockTriggerService) Close() error {
	close(o.stopCh)
	o.wg.Wait()
	o.lggr.Info("MockTriggerService closed")
	return nil
}

func (o *MockTriggerService) Ready() error {
	return nil
}

func (o *MockTriggerService) HealthReport() map[string]error {
	return nil
}

func (o *MockTriggerService) Name() string {
	return "MockTriggerService"
}

type mockDataProducer struct {
	trigger *MockTriggerService
	wg      sync.WaitGroup
	closeCh chan struct{}
	lggr    logger.Logger
}

var _ services.Service = &mockDataProducer{}

func NewMockDataProducer(trigger *MockTriggerService, lggr logger.Logger) *mockDataProducer {
	return &mockDataProducer{
		trigger: trigger,
		closeCh: make(chan struct{}),
		lggr:    lggr,
	}
}

func (m *mockDataProducer) Start(ctx context.Context) error {
	m.wg.Add(1)
	go m.loop()
	return nil
}

func (m *mockDataProducer) loop() {
	defer m.wg.Done()

	sleepSec := 15
	ticker := time.NewTicker(time.Duration(sleepSec) * time.Second)
	defer ticker.Stop()

	prices := []int64{300000, 40000, 5000000}

	j := 0

	for range ticker.C {
		for i := range prices {
			prices[i] = prices[i] + 1
		}
		j += 1

		// https://github.com/smartcontractkit/chainlink/blob/41f9428c3aa8231e8834a230fca4c2ccffd4e6c3/core/capabilities/streams/trigger_test.go#L117-L122

		timestamp := time.Now().Unix()
		// TODO: shouldn't we increment round rather than epoch?
		reportCtx := ocrTypes.ReportContext{ReportTimestamp: ocrTypes.ReportTimestamp{Epoch: uint32(baseTimestamp + j)}}

		reports := []datastreams.FeedReport{
			{
				FeedID:               "0x1111111111111111111100000000000000000000000000000000000000000000",
				FullReport:           newReport(m.lggr, common.HexToHash("0x1111111111111111111100000000000000000000000000000000000000000000"), big.NewInt(prices[0]), timestamp),
				ReportContext:        rawReportContext(reportCtx),
				ObservationTimestamp: timestamp,
			},
			{
				FeedID:               "0x2222222222222222222200000000000000000000000000000000000000000000",
				FullReport:           newReport(m.lggr, common.HexToHash("0x2222222222222222222200000000000000000000000000000000000000000000"), big.NewInt(prices[1]), timestamp),
				ReportContext:        rawReportContext(reportCtx),
				ObservationTimestamp: timestamp,
			},
			{
				FeedID:               "0x3333333333333333333300000000000000000000000000000000000000000000",
				FullReport:           newReport(m.lggr, common.HexToHash("0x3333333333333333333300000000000000000000000000000000000000000000"), big.NewInt(prices[2]), timestamp),
				ReportContext:        rawReportContext(reportCtx),
				ObservationTimestamp: timestamp,
			},
		}

		m.lggr.Infow("New set of Mock reports", "timestamp", time.Now().Unix(), "payload", reports)
		err := m.trigger.ProcessReport(reports)
		if err != nil {
			m.lggr.Errorw("failed to process Mock reports", "err", err, "timestamp", time.Now().Unix(), "payload", reports)
		}
	}
}

func (m *mockDataProducer) Close() error {
	close(m.closeCh)
	m.wg.Wait()
	return nil
}

func (m *mockDataProducer) HealthReport() map[string]error {
	return nil
}

func (m *mockDataProducer) Ready() error {
	return nil
}

func (m *mockDataProducer) Name() string {
	return "mockDataProducer"
}
