package protocol

import (
	"context"
	"math"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func RunReportFinalization(
	ctx context.Context,

	chNetToReportFinalization <-chan MessageToReportFinalizationWithSender,
	chReportFinalizationToTransmission chan<- EventToTransmission,
	chReportGenerationToReportFinalization <-chan EventToReportFinalization,
	config config.SharedConfig,
	contractSigner types.OnchainKeyring,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender,
	reportQuorum int,
) {
	newReportFinalizationState(ctx, chNetToReportFinalization,
		chReportFinalizationToTransmission, chReportGenerationToReportFinalization,
		config, contractSigner, logger, netSender, reportQuorum).run()
}

const minExpirationAgeRounds int = 10
const expirationAgeDuration = 10 * time.Minute
const maxExpirationAgeRounds int = 1_000

type reportFinalizationState struct {
	ctx context.Context

	chNetToReportFinalization              <-chan MessageToReportFinalizationWithSender
	chReportFinalizationToTransmission     chan<- EventToTransmission
	chReportGenerationToReportFinalization <-chan EventToReportFinalization
	config                                 config.SharedConfig
	contractSigner                         types.OnchainKeyring
	logger                                 loghelper.LoggerWithContext
	netSender                              NetworkSender
	reportQuorum                           int

	// reap() is used to prevent unbounded state growth of finalized
	finalized       map[EpochRound]struct{}
	finalizedLatest EpochRound
}

func (repfin *reportFinalizationState) run() {
	for {
		select {
		case msg := <-repfin.chNetToReportFinalization:
			msg.msg.processReportFinalization(repfin, msg.sender)
		case ev := <-repfin.chReportGenerationToReportFinalization:
			ev.processReportFinalization(repfin)
		case <-repfin.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-repfin.ctx.Done():
			repfin.logger.Info("ReportFinalization: exiting", nil)
			return
		default:
		}
	}
}

// messageFinalEcho is called when the local oracle process receives a
// "final-echo" message.
func (repfin *reportFinalizationState) messageFinalEcho(
	msg MessageFinalEcho,
	sender commontypes.OracleID,
) {
	if msg.Round > repfin.config.RMax {
		repfin.logger.Debug("ignoring MessageFinalEcho for round larger than rMax", commontypes.LogFields{
			"epoch": msg.Epoch, "round": msg.Round, "rMax": repfin.config.RMax, "sender": sender})
		return
	}

	epochRound := EpochRound{msg.Epoch, msg.Round}
	if repfin.isExpired(epochRound) {
		repfin.logger.Debug("ignoring MessageFinalEcho for expired epoch and round", commontypes.LogFields{
			"epoch": msg.Epoch, "round": msg.Round, "sender": sender})
		return
	}
	if _, ok := repfin.finalized[epochRound]; ok {
		repfin.logger.Debug("ignoring MessageFinalEcho for already finalized epoch and round", commontypes.LogFields{
			"epoch": msg.Epoch, "round": msg.Round, "sender": sender})
		return
	}

	err := msg.AttestedReport.VerifySignatures(
		repfin.reportQuorum,
		repfin.contractSigner,
		repfin.config.OracleIdentities,
		types.ReportContext{
			types.ReportTimestamp{repfin.config.ConfigDigest, msg.Epoch, msg.Round},
			msg.H,
		},
	)
	if err != nil {
		repfin.logger.Warn("error while verifying signatures on attested report", commontypes.LogFields{
			"msg":    msg,
			"sender": sender,
			"error":  err,
		})
		return
	}

	repfin.finalize(msg.MessageFinal)
}

func (repfin *reportFinalizationState) eventFinal(ev EventFinal) {
	epochRound := EpochRound{ev.Epoch, ev.Round}
	if repfin.isExpired(epochRound) {
		repfin.logger.Debug("ignoring EventFinal for expired epoch and round", commontypes.LogFields{
			"epoch": ev.Epoch, "round": ev.Round})
		return
	}
	if _, ok := repfin.finalized[epochRound]; ok {
		repfin.logger.Debug("ignoring EventFinal for already finalized epoch and round", commontypes.LogFields{
			"epoch": ev.Epoch, "round": ev.Round})
		return
	}

	repfin.finalize(ev.MessageFinal)
}

func (repfin *reportFinalizationState) finalize(msg MessageFinal) {
	repfin.logger.Debug("finalizing report", commontypes.LogFields{
		"epoch": msg.Epoch,
		"round": msg.Round,
	})

	epochRound := EpochRound{msg.Epoch, msg.Round}

	repfin.finalized[epochRound] = struct{}{}
	if repfin.finalizedLatest.Less(epochRound) {
		repfin.finalizedLatest = epochRound
	}

	repfin.netSender.Broadcast(MessageFinalEcho{msg}) // send [ FINALECHO, e, r, O] to all p_j ∈ P

	select {
	case repfin.chReportFinalizationToTransmission <- EventTransmit(msg):
	case <-repfin.ctx.Done():
	}

	repfin.reap()
}

func (repfin *reportFinalizationState) isExpired(er EpochRound) bool {
	latestIndex := repfin.epochRoundIndex(repfin.finalizedLatest)
	expiredIndex := latestIndex - int64(repfin.expirationAgeRounds())
	return repfin.epochRoundIndex(er) <= expiredIndex
}

// reap expired entries from repfin.finalized to prevent unbounded state growth
func (repfin *reportFinalizationState) reap() {
	if len(repfin.finalized) <= 2*repfin.expirationAgeRounds() {
		return
	}
	// A long time ago in a galaxy far, far away, Go used to leak memory when
	// repeatedly adding and deleting from the same map without ever exceeding
	// some maximum length. Fortunately, this is no longer the case
	// https://go-review.googlesource.com/c/go/+/25049/
	for er := range repfin.finalized {
		if repfin.isExpired(er) {
			delete(repfin.finalized, er)
		}
	}
}

// The age (denoted in rounds) after which a report is considered expired and
// will automatically be dropped
func (repfin *reportFinalizationState) expirationAgeRounds() int {
	// number of rounds in a window of duration expirationAgeDuration
	age := math.Ceil(expirationAgeDuration.Seconds() / repfin.config.DeltaRound.Seconds())

	if age < float64(minExpirationAgeRounds) {
		age = float64(minExpirationAgeRounds)
	}
	if math.IsNaN(age) || age > float64(maxExpirationAgeRounds) {
		age = float64(maxExpirationAgeRounds)
	}

	return int(age)
}

// Assigns consecutive indexes to consecutive (epoch, round) tuples
func (repfin *reportFinalizationState) epochRoundIndex(er EpochRound) int64 {
	// safe from overflow. Epoch is a uint32 and Round is a uint8, so this will
	// always fit in a uint40
	return int64(repfin.config.RMax)*int64(er.Epoch) + int64(er.Round)
}

func newReportFinalizationState(
	ctx context.Context,

	chNetToReportFinalization <-chan MessageToReportFinalizationWithSender,
	chReportFinalizationToTransmission chan<- EventToTransmission,
	chReportGenerationToReportFinalization <-chan EventToReportFinalization,
	config config.SharedConfig,
	contractSigner types.OnchainKeyring,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender,
	reportQuorum int,
) *reportFinalizationState {
	return &reportFinalizationState{
		ctx,

		chNetToReportFinalization,
		chReportFinalizationToTransmission,
		chReportGenerationToReportFinalization,
		config,
		contractSigner,
		logger,
		netSender,
		reportQuorum,

		map[EpochRound]struct{}{},
		EpochRound{},
	}
}
