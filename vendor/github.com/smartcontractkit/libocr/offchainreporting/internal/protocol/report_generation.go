package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// Report Generation protocol corresponding to alg. 2 & 3.
func RunReportGeneration(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chNetToReportGeneration <-chan MessageToReportGenerationWithSender,
	chReportGenerationToPacemaker chan<- EventToPacemaker,
	chReportGenerationToTransmission chan<- EventToTransmission,
	config config.SharedConfig,
	configOverrider types.ConfigOverrider,
	contractTransmitter types.ContractTransmitter,
	datasource types.DataSource,
	e uint32,
	id commontypes.OracleID,
	l commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender,
	privateKeys types.PrivateKeys,
	telemetrySender TelemetrySender,
) {
	repgen := reportGenerationState{
		ctx:          ctx,
		subprocesses: subprocesses,

		chNetToReportGeneration:          chNetToReportGeneration,
		chReportGenerationToPacemaker:    chReportGenerationToPacemaker,
		chReportGenerationToTransmission: chReportGenerationToTransmission,
		config:                           config,
		configOverrider:                  configOverrider,
		contractTransmitter:              contractTransmitter,
		datasource:                       datasource,
		e:                                e,
		id:                               id,
		l:                                l,
		localConfig:                      localConfig,
		logger:                           logger.MakeChild(commontypes.LogFields{"epoch": e, "leader": l}),
		netSender:                        netSender,
		privateKeys:                      privateKeys,
		telemetrySender:                  telemetrySender,
	}
	repgen.run()
}

type reportGenerationState struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chNetToReportGeneration          <-chan MessageToReportGenerationWithSender
	chReportGenerationToPacemaker    chan<- EventToPacemaker
	chReportGenerationToTransmission chan<- EventToTransmission
	config                           config.SharedConfig
	configOverrider                  types.ConfigOverrider
	contractTransmitter              types.ContractTransmitter
	datasource                       types.DataSource
	e                                uint32 // Current epoch number
	id                               commontypes.OracleID
	l                                commontypes.OracleID // Current leader number
	localConfig                      types.LocalConfig
	logger                           loghelper.LoggerWithContext
	netSender                        NetworkSender
	privateKeys                      types.PrivateKeys
	telemetrySender                  TelemetrySender

	leaderState   leaderState
	followerState followerState
}

type leaderState struct {
	// r is the current round within the epoch
	r uint8

	// observe contains the observations received so far
	observe []*SignedObservation

	// report contains the signed reports received so far
	report []*AttestedReportOne

	// tRound is a heartbeat indicating when the current leader should start a new
	// round.
	tRound <-chan time.Time

	// tGrace is a grace period the leader waits for after it has achieved
	// quorum on "observe" messages, to allow slower oracles time to submit their
	// observations.
	tGrace <-chan time.Time

	phase phase
}

type followerState struct {
	// r is the current round within the epoch
	r uint8

	// receivedEcho's j-th entry indicates whether a valid final echo has been received
	// from the j-th oracle
	receivedEcho []bool

	// sentEcho tracks the report the current oracle has final-echoed during
	// this round.
	sentEcho *AttestedReportMany

	// sentReport tracks whether the current oracles has sent a report during
	// this round
	sentReport bool

	// completedRound tracks whether the current oracle has completed the current
	// round
	completedRound bool
}

// Run starts the event loop for the report-generation protocol
func (repgen *reportGenerationState) run() {
	repgen.logger.Info("Running ReportGeneration", nil)

	// Initialization
	repgen.leaderState.r = 0
	repgen.leaderState.report = make([]*AttestedReportOne, repgen.config.N())
	repgen.followerState.r = 0
	repgen.followerState.receivedEcho = make([]bool, repgen.config.N())
	repgen.followerState.sentEcho = nil
	repgen.followerState.completedRound = false

	// kick off the protocol
	if repgen.id == repgen.l {
		repgen.startRound()
	}

	// Event Loop
	chDone := repgen.ctx.Done()
	for {
		select {
		case msg := <-repgen.chNetToReportGeneration:
			msg.msg.processReportGeneration(repgen, msg.sender)
		case <-repgen.leaderState.tGrace:
			repgen.eventTGraceTimeout()
		case <-repgen.leaderState.tRound:
			repgen.eventTRoundTimeout()
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			repgen.logger.Info("ReportGeneration: exiting", commontypes.LogFields{
				"e": repgen.e,
				"l": repgen.l,
			})
			return
		default:
		}
	}
}
