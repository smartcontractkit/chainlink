package commit

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/ccipocr3/internal/libs/slicelib"
	"github.com/smartcontractkit/ccipocr3/internal/model"
	"github.com/smartcontractkit/ccipocr3/internal/reader"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Plugin implements the main ocr3 plugin logic.
type Plugin struct {
	nodeID     model.NodeID
	cfg        model.CommitPluginConfig
	ccipReader reader.CCIP
	lggr       logger.Logger
}

func NewPlugin(
	_ context.Context,
	nodeID model.NodeID,
	cfg model.CommitPluginConfig,
	ccipReader reader.CCIP,
) *Plugin {
	return &Plugin{
		nodeID:     nodeID,
		cfg:        cfg,
		ccipReader: ccipReader,
	}
}

func (p *Plugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	return types.Query{}, nil
}

func (p *Plugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	scanFrom := time.Now().Add(-p.cfg.NewMsgScanDuration)
	newMsgs, err := p.ccipReader.MsgsAfterTimestamp(ctx, p.cfg.Reads, scanFrom, p.cfg.NewMsgScanLimit)
	if err != nil {
		p.lggr.Errorw("get new ccip messages", "err", err)
		return types.Observation{}, err
	}
	p.lggr.Debugf("found %d new ccip messages", len(newMsgs))

	newMsgDetails := make([]model.CCIPMsgBaseDetails, 0, len(newMsgs))
	for _, msg := range newMsgs {
		p.lggr.Debugf("new msg: %s", msg)
		newMsgDetails = append(newMsgDetails, msg.CCIPMsgBaseDetails)
	}

	return model.NewCommitPluginObservation(p.nodeID, newMsgDetails).Encode()
}

func (p *Plugin) ValidateObservation(outctx ocr3types.OutcomeContext, _ types.Query, ao types.AttributedObservation) error {
	_, err := model.DecodeCommitPluginObservation(ao.Observation)
	return err
}

func (p *Plugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query types.Query) (ocr3types.Quorum, error) {
	// across all chains we require at least 2f+1 observations.
	return ocr3types.QuorumTwoFPlusOne, nil
}

func (p *Plugin) Outcome(outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	msgsFromObservations := make([]model.CCIPMsgBaseDetails, 0)
	for _, ao := range aos {
		parsedObservation, err := model.DecodeCommitPluginObservation(ao.Observation)
		if err != nil {
			p.lggr.Errorw("decode commit plugin observation", "err", err)
			return ocr3types.Outcome{}, err
		}
		msgsFromObservations = append(msgsFromObservations, parsedObservation.NewMsgs...)
	}

	sourceChains, groupedMsgs := slicelib.GroupBy(
		msgsFromObservations, func(msg model.CCIPMsgBaseDetails) model.ChainSelector { return msg.SourceChain })
	for _, sourceChain := range sourceChains {
		p.lggr.Debugf("for source chain %d we got %d msg observations", len(groupedMsgs[sourceChain]))
	}

	return ocr3types.Outcome{}, fmt.Errorf("implement me")
}

func (p *Plugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[[]byte], error) {
	/*
		# Only a single report is needed containing a batch of gas price updates for the chain.
		# and a list of roots by source. We only include the gas price batch update
		# if a timer has expired - ie periodically we batch write all gas prices.
	*/
	panic("implement me")
}

func (p *Plugin) ShouldAcceptAttestedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	panic("implement me")
}

func (p *Plugin) ShouldTransmitAcceptedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	/*
		if dest not in self.chains:
			# Can't write, skip
			return False
		return self.chains[dest].is_stale(report):
	*/
	panic("implement me")
}

func (p *Plugin) Close() error {
	panic("implement me")
}

// Interface compatibility checks.
var _ ocr3types.ReportingPlugin[[]byte] = &Plugin{}
