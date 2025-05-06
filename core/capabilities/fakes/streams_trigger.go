package fakes

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	commonStreams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

type fakeStreamsTrigger struct {
	services.Service
	eng  *services.Engine
	lggr logger.Logger

	signers []ocr2key.KeyBundle
	codec   datastreams.ReportCodec
	meta    datastreams.Metadata

	registrations map[string]*regState
	lastPrice     map[string]int64
	lastEventTs   int64
	mu            sync.Mutex
}

type regState struct {
	cfg     commonStreams.TriggerConfig
	eventCh chan commonCap.TriggerResponse
}

var _ services.Service = (*fakeStreamsTrigger)(nil)
var _ commonCap.TriggerCapability = (*fakeStreamsTrigger)(nil)

const (
	triggerID    = "streams-trigger@1.0.0"
	resolutionMs = 5000
)

func (st *fakeStreamsTrigger) Info(ctx context.Context) (commonCap.CapabilityInfo, error) {
	return commonCap.CapabilityInfo{
		ID:             triggerID,
		CapabilityType: commonCap.CapabilityTypeTrigger,
		Description:    "Fake Streams Trigger",
		DON:            &commonCap.DON{},
		IsLocal:        true,
	}, nil
}

func (st *fakeStreamsTrigger) RegisterTrigger(ctx context.Context, request commonCap.TriggerRegistrationRequest) (<-chan commonCap.TriggerResponse, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	reg, found := st.registrations[request.Metadata.WorkflowID]
	if found {
		return reg.eventCh, nil
	}
	state := &regState{
		eventCh: make(chan commonCap.TriggerResponse, 1000),
	}
	if err := request.Config.UnwrapTo(&state.cfg); err != nil {
		return nil, err
	}
	for _, feed := range state.cfg.FeedIds {
		if _, found := st.lastPrice[string(feed)]; !found {
			st.lastPrice[string(feed)] = 1
		}
	}

	st.registrations[request.Metadata.WorkflowID] = state
	st.eng.Infow("Registered to Fake Streams Trigger", "workflowID", request.Metadata.WorkflowID)
	return state.eventCh, nil
}

func (st *fakeStreamsTrigger) UnregisterTrigger(ctx context.Context, request commonCap.TriggerRegistrationRequest) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	reg, found := st.registrations[request.Metadata.WorkflowID]
	if !found {
		return errors.New("not registered")
	}
	close(reg.eventCh)
	delete(st.registrations, request.Metadata.WorkflowID)
	st.eng.Infow("Unregistered from Fake Streams Trigger", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func NewFakeStreamsTrigger(lggr logger.Logger, nSigners int) *fakeStreamsTrigger {
	signers := make([]ocr2key.KeyBundle, nSigners)
	rawSigners := make([][]byte, nSigners)
	for i := range nSigners {
		signers[i], _ = ocr2key.New(chaintype.EVM)
		rawSigners[i] = signers[i].PublicKey()
	}

	st := &fakeStreamsTrigger{
		signers: signers,
		meta: datastreams.Metadata{
			Signers:               rawSigners,
			MinRequiredSignatures: nSigners,
		},
		codec:         streams.NewCodec(lggr),
		registrations: make(map[string]*regState),
		lastPrice:     make(map[string]int64),
		lggr:          lggr,
	}
	st.Service, st.eng = services.Config{
		Name:  "fakeStreamsTrigger",
		Start: st.start,
	}.NewServiceEngine(lggr)
	return st
}

func (st *fakeStreamsTrigger) start(_ context.Context) error {
	ticker := services.TickerConfig{
		Initial:   1000 * time.Millisecond,
		JitterPct: 0.0,
	}.NewTicker(100 * time.Millisecond)
	st.eng.GoTick(ticker, st.emitEvent)
	return nil
}

func (st *fakeStreamsTrigger) emitEvent(ctx context.Context) {
	st.mu.Lock()
	defer st.mu.Unlock()

	nowMs := time.Now().UnixMilli()
	newTs := (nowMs / resolutionMs) * resolutionMs
	if newTs == st.lastEventTs {
		return
	}
	st.lastEventTs = newTs
	tsSec := uint32(newTs / 1000) //nolint:gosec // G115

	// generate new reports for all feeds
	reports := make(map[string]datastreams.FeedReport)
	for feed, lastPrice := range st.lastPrice {
		st.lastPrice[feed] = lastPrice + 1
		idBytes, err := hex.DecodeString(feed[2:])
		if err != nil {
			st.eng.Errorw("Failed to decode feed ID", "error", err)
			continue
		}
		idBytes32 := [32]byte{}
		copy(idBytes32[:], idBytes)
		fullReport := newReport(ctx, st.lggr, idBytes32, lastPrice, tsSec)

		reportCtx := ocrTypes.ReportContext{}
		rawCtx := rawReportContext(reportCtx)

		feedReport := datastreams.FeedReport{
			FeedID:        feed,
			FullReport:    fullReport,
			ReportContext: rawCtx,
			Signatures:    [][]byte{},
		}
		for _, signer := range st.signers {
			signature, err := signer.Sign(reportCtx, fullReport)
			if err != nil {
				st.eng.Errorw("Failed to sign report", "error", err)
				continue
			}
			feedReport.Signatures = append(feedReport.Signatures, signature)
		}
		reports[feed] = feedReport
	}

	eventID := fmt.Sprintf("streams_%d", newTs)

	for wf, reg := range st.registrations {
		if reg.cfg.MaxFrequencyMs > 0 && newTs%int64(reg.cfg.MaxFrequencyMs) != 0 { //nolint:gosec // G115
			continue
		}
		requestedReports := []datastreams.FeedReport{}
		for _, feedID := range reg.cfg.FeedIds {
			if report, found := reports[string(feedID)]; found {
				requestedReports = append(requestedReports, report)
			}
		}
		event, err := triggers.WrapReports(requestedReports, eventID, nowMs, st.meta, triggerID)
		if err != nil {
			st.eng.Errorw("Failed to wrap reports", "error", err)
			return
		}
		st.eng.Infow("Sending event to workflow", "eventID", eventID, "workflowID", wf)
		reg.eventCh <- event
	}
}

func newReport(ctx context.Context, lggr logger.Logger, feedID [32]byte, price int64, timestampSec uint32) []byte {
	priceB := big.NewInt(price)
	v3Codec := reportcodec.NewReportCodec(feedID, lggr)
	raw, _ := v3Codec.BuildReport(ctx, v3.ReportFields{
		BenchmarkPrice:     priceB,
		Timestamp:          timestampSec,
		ValidFromTimestamp: timestampSec,
		Bid:                priceB,
		Ask:                priceB,
		LinkFee:            priceB,
		NativeFee:          priceB,
		ExpiresAt:          timestampSec + 3600,
	})
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
