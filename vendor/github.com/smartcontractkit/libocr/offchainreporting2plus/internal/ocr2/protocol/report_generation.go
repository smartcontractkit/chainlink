package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr2config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// Report Generation protocol corresponding to alg. 2 & 3.
func RunReportGeneration(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chNetToReportGeneration <-chan MessageToReportGenerationWithSender,
	chReportGenerationToPacemaker chan<- EventToPacemaker,
	chReportGenerationToReportFinalization chan<- EventToReportFinalization,
	config ocr2config.SharedConfig,
	contractTransmitter types.ContractTransmitter,
	e uint32,
	id commontypes.OracleID,
	l commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring types.OnchainKeyring,
	reportGenerationMetrics *reportGenerationMetrics,
	reportingPlugin types.ReportingPlugin,
	reportQuorum int,
	telemetrySender TelemetrySender,
) {
	repgen := reportGenerationState{
		ctx:          ctx,
		subprocesses: subprocesses,

		chNetToReportGeneration:                chNetToReportGeneration,
		chReportGenerationToPacemaker:          chReportGenerationToPacemaker,
		chReportGenerationToReportFinalization: chReportGenerationToReportFinalization,
		config:                                 config,
		contractTransmitter:                    contractTransmitter,
		e:                                      e,
		id:                                     id,
		l:                                      l,
		localConfig:                            localConfig,
		logger:                                 logger.MakeChild(commontypes.LogFields{"epoch": e, "leader": l}),
		netSender:                              netSender,
		offchainKeyring:                        offchainKeyring,
		onchainKeyring:                         onchainKeyring,
		reportGenerationMetrics:                reportGenerationMetrics,
		reportingPlugin:                        reportingPlugin,
		reportQuorum:                           reportQuorum,
		telemetrySender:                        telemetrySender,
	}
	repgen.run()
}

type reportGenerationState struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chNetToReportGeneration                <-chan MessageToReportGenerationWithSender
	chReportGenerationToPacemaker          chan<- EventToPacemaker
	chReportGenerationToReportFinalization chan<- EventToReportFinalization
	config                                 ocr2config.SharedConfig
	contractTransmitter                    types.ContractTransmitter
	e                                      uint32 // Current epoch number
	id                                     commontypes.OracleID
	l                                      commontypes.OracleID // Current leader number
	localConfig                            types.LocalConfig
	logger                                 loghelper.LoggerWithContext
	netSender                              NetworkSender
	offchainKeyring                        types.OffchainKeyring
	onchainKeyring                         types.OnchainKeyring
	reportGenerationMetrics                *reportGenerationMetrics
	reportingPlugin                        types.ReportingPlugin
	reportQuorum                           int
	telemetrySender                        TelemetrySender

	leaderState   leaderState
	followerState followerState
}

type leaderState struct {
	// r is the current round within the epoch
	r uint8

	q types.Query

	// observe contains the observations received so far
	observe []*SignedObservation

	// report contains the signed reports received so far
	report []*AttestedReportOne

	// tRound is a heartbeat indicating when the current leader should start a new
	// round.
	tRound <-chan time.Time

	// a round is only ready to start once the previous round has finished AND
	// tRound has fired
	readyToStartRound bool

	// tGrace is a grace period the leader waits for after it has achieved
	// quorum on "observe" messages, to allow slower oracles time to submit their
	// observations.
	tGrace <-chan time.Time

	phase phase

	// H(Q, B)
	h [32]byte
}

type followerState struct {
	// r is the current round within the epoch
	r uint8

	// sentReport tracks whether the current oracles has sent a report during
	// this round
	sentReport bool
	// sentReportAttributedObservations []types.AttributedObservation

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
	repgen.leaderState.readyToStartRound = false
	repgen.followerState.r = 0
	repgen.followerState.completedRound = false

	// kick off the protocol
	if repgen.id == repgen.l {
		repgen.startRound()
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
