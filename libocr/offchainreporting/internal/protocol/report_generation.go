package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"
)

func RunReportGeneration(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chNetToReportGeneration <-chan MessageToReportGenerationWithSender,
	chReportGenerationToPacemaker chan<- EventToPacemaker,
	chReportGenerationToTransmission chan<- EventToTransmission,
	config config.SharedConfig,
	contractTransmitter types.ContractTransmitter,
	datasource types.DataSource,
	e uint32,
	id types.OracleID,
	l types.OracleID,
	localConfig types.LocalConfig,
	logger types.Logger,
	netSender NetworkSender,
	privateKeys types.PrivateKeys,
) {
	repgen := reportGenerationState{
		ctx:          ctx,
		subprocesses: subprocesses,

		chNetToReportGeneration:          chNetToReportGeneration,
		chReportGenerationToPacemaker:    chReportGenerationToPacemaker,
		chReportGenerationToTransmission: chReportGenerationToTransmission,
		config:                           config,
		contractTransmitter:              contractTransmitter,
		datasource:                       datasource,
		e:                                e,
		id:                               id,
		l:                                l,
		localConfig:                      localConfig,
		logger:                           loghelper.MakeLoggerWithContext(logger, types.LogFields{"epoch": e, "leader": l}),
		netSender:                        netSender,
		privateKeys:                      privateKeys,
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
	contractTransmitter              types.ContractTransmitter
	datasource                       types.DataSource
	e                                uint32 	id                               types.OracleID
	l                                types.OracleID 	localConfig                      types.LocalConfig
	logger                           types.Logger
	netSender                        NetworkSender
	privateKeys                      types.PrivateKeys

	leaderState   leaderState
	followerState followerState
}

type leaderState struct {
		r uint8

		observe []*MessageObserve

		report []*MessageReport

			tRound <-chan time.Time

				tGrace <-chan time.Time

	phase phase
}

type followerState struct {
		r uint8

				finalEcho []*MessageFinalEcho

			sentEcho bool

			sentReport bool

			completedRound bool
}

func (repgen *reportGenerationState) NewReportingContext() signature.ReportingContext {
	return signature.NewReportingContext(repgen.config.ConfigDigest, repgen.e, repgen.followerState.r)
}

func (repgen *reportGenerationState) run() {
	repgen.logger.Info("Running ReportGeneration", nil)

		repgen.leaderState.r = 0
	repgen.leaderState.report = make([]*MessageReport, repgen.config.N())
	repgen.followerState.r = 0
	repgen.followerState.finalEcho = make([]*MessageFinalEcho, repgen.config.N())
	repgen.followerState.sentEcho = false
	repgen.followerState.completedRound = false

		if repgen.id == repgen.l {
		repgen.startRound()
	}

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

				select {
		case <-chDone:
			repgen.logger.Info("ReportGeneration: exiting", types.LogFields{
				"e": repgen.e,
				"l": repgen.l,
			})
			return
		default:
		}
	}
}
