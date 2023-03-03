package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

///////////////////////////////////////////////////////////
// Report Generation Leader (Algorithm 3)
///////////////////////////////////////////////////////////

type phase int

const (
	phaseObserve phase = iota
	phaseGrace
	phaseReport
	phaseFinal
)

var englishPhase = map[phase]string{
	phaseObserve: "observe",
	phaseGrace:   "grace",
	phaseReport:  "report",
	phaseFinal:   "final",
}

func (repgen *reportGenerationState) leaderReportTimestamp() types.ReportTimestamp {
	return types.ReportTimestamp{repgen.config.ConfigDigest, repgen.e, repgen.leaderState.r}
}

///////////////////////////////////////////////////////////
// Report Generation Leader (Algorithm 4)
///////////////////////////////////////////////////////////

func (repgen *reportGenerationState) eventTRoundTimeout() {
	repgen.startRound()
}

// startRound is called upon initialization of the leaders' report-generation
// protocol instance, or when the round timer expires, indicating that it
// should start a new round.
//
// It broadcasts an observe-req message to all participants, and restarts the
// round timer.
func (repgen *reportGenerationState) startRound() {
	if !repgen.leaderState.readyToStartRound {
		repgen.leaderState.readyToStartRound = true
		return
	}

	if repgen.leaderState.r > repgen.config.RMax {
		repgen.logger.Warn("ReportGeneration: new round number would be larger than RMax + 1. Looks like your connection to more than f other nodes is not working.", commontypes.LogFields{
			"round": repgen.leaderState.r,
			"f":     repgen.config.F,
			"RMax":  repgen.config.RMax,
		})
		return
	}
	rPlusOne := repgen.leaderState.r + 1
	if rPlusOne <= repgen.leaderState.r {
		repgen.logger.Error("ReportGeneration: round overflows, cannot start new round", commontypes.LogFields{
			"round": repgen.leaderState.r,
		})
		return
	}
	repgen.leaderState.r = rPlusOne
	repgen.leaderState.observe = make([]*SignedObservation, repgen.config.N())
	repgen.leaderState.report = make([]*AttestedReportOne, repgen.config.N())
	repgen.leaderState.tRound = time.After(repgen.config.DeltaRound)
	repgen.leaderState.readyToStartRound = false
	var query types.Query
	{
		ctx, cancel := context.WithTimeout(repgen.ctx, repgen.config.MaxDurationQuery)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			repgen.config.MaxDurationQuery+ReportingPluginTimeoutWarningGracePeriod,
			func() {
				repgen.logger.Error("ReportGeneration: ReportingPlugin.Query is taking too long", commontypes.LogFields{
					"round": repgen.leaderState.r, "maxDuration": repgen.config.MaxDurationQuery,
				})
			},
		)

		var err error
		query, err = repgen.reportingPlugin.Query(ctx, repgen.leaderReportTimestamp())

		ins.Stop()

		if err != nil {
			repgen.logger.Error("ReportGeneration: error while calling ReportingPlugin.Query. cannot start new round", commontypes.LogFields{
				"round": repgen.leaderState.r,
				"error": err,
			})
			return
		}

	}
	repgen.leaderState.q = query
	repgen.leaderState.h = [32]byte{} // Not strictly necessary, but makes testing cleaner
	repgen.leaderState.phase = phaseObserve
	repgen.netSender.Broadcast(MessageObserveReq{
		repgen.e,
		repgen.leaderState.r,
		query,
	})
}

// messageObserve is called when the current leader has received an "observe"
// message. If the leader has enough observations to construct a report, given
// this message, it kicks off the T_observe grace period, to allow slower
// oracles time to submit their observations. It only responds to these messages
// when in the observe or grace phases
func (repgen *reportGenerationState) messageObserve(msg MessageObserve, sender commontypes.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageObserve for wrong epoch", commontypes.LogFields{
			"round":    repgen.leaderState.r,
			"sender":   sender,
			"msgEpoch": msg.Epoch,
			"msgRound": msg.Round,
		})
		return
	}

	if repgen.l != repgen.id {
		repgen.logger.Warn("Non-leader received MessageObserve", commontypes.LogFields{
			"round":  repgen.leaderState.r,
			"sender": sender,
			"msg":    msg,
		})
		return
	}

	if msg.Round != repgen.leaderState.r {
		repgen.logger.Debug("Got MessageObserve for wrong round", commontypes.LogFields{
			"round":    repgen.leaderState.r,
			"sender":   sender,
			"msgEpoch": msg.Epoch,
			"msgRound": msg.Round,
		})
		return
	}

	if repgen.leaderState.phase != phaseObserve && repgen.leaderState.phase != phaseGrace {
		repgen.logger.Debug("received MessageObserve after grace phase", commontypes.LogFields{
			"round": repgen.leaderState.r,
		})
		return
	}

	if repgen.leaderState.observe[sender] != nil {
		repgen.logger.Debug("already sent an observation", commontypes.LogFields{
			"round":  repgen.leaderState.r,
			"sender": sender,
		})
		return
	}

	if err := msg.SignedObservation.Verify(repgen.leaderReportTimestamp(), repgen.leaderState.q, repgen.config.OracleIdentities[sender].OffchainPublicKey); err != nil {
		repgen.logger.Warn("MessageObserve carries invalid SignedObservation", commontypes.LogFields{
			"round":  repgen.leaderState.r,
			"sender": sender,
			"msg":    msg,
			"error":  err,
		})
		return
	}

	repgen.logger.Debug("MessageObserve has valid SignedObservation", commontypes.LogFields{
		"round":    repgen.leaderState.r,
		"sender":   sender,
		"msgEpoch": msg.Epoch,
		"msgRound": msg.Round,
	})

	repgen.leaderState.observe[sender] = &msg.SignedObservation

	//upon (|{p_j ∈ P| observe[j] != ⊥}| > 2f) ∧ (phase = OBSERVE)
	switch repgen.leaderState.phase {
	case phaseObserve:
		observationCount := 0 // FUTUREWORK: Make this count constant-time with state counter
		for _, so := range repgen.leaderState.observe {
			if so != nil {
				observationCount++
			}
		}
		repgen.logger.Debug("One more observation", commontypes.LogFields{
			"round":                    repgen.leaderState.r,
			"observationCount":         observationCount,
			"requiredObservationCount": (2 * repgen.config.F) + 1,
		})
		if observationCount > 2*repgen.config.F {
			// Start grace period, to allow slower oracles to contribute observations
			repgen.logger.Debug("starting observation grace period", commontypes.LogFields{
				"round": repgen.leaderState.r,
			})
			repgen.leaderState.tGrace = time.After(repgen.config.DeltaGrace)
			repgen.leaderState.phase = phaseGrace
		}
	case phaseGrace:
		repgen.logger.Debug("accepted extra observation during grace period", nil)
	case phaseFinal:
		repgen.logger.Error("unexpected phase phaseFinal", commontypes.LogFields{"round": repgen.leaderState.r})
	case phaseReport:
		repgen.logger.Error("unexpected phase phaseReport", commontypes.LogFields{"round": repgen.leaderState.r})
	}
}

// eventTGraceTimeout is called by the leader when the grace period
// is over. It collates the signed observations it has received so far, and
// sends out a request for participants' signatures on the final report.
func (repgen *reportGenerationState) eventTGraceTimeout() {
	if repgen.leaderState.phase != phaseGrace {
		repgen.logger.Error("leader's phase conflicts tGrace timeout", commontypes.LogFields{
			"round": repgen.leaderState.r,
			"phase": englishPhase[repgen.leaderState.phase],
		})
		return
	}
	asos := []AttributedSignedObservation{}
	aos := []types.AttributedObservation{}
	for oid, so := range repgen.leaderState.observe {
		if so != nil {
			asos = append(asos, AttributedSignedObservation{
				*so,
				commontypes.OracleID(oid),
			})
			aos = append(aos, types.AttributedObservation{
				so.Observation,
				commontypes.OracleID(oid),
			})
		}
	}

	repgen.netSender.Broadcast(MessageReportReq{
		repgen.e,
		repgen.leaderState.r,
		repgen.leaderState.q,
		asos,
	})

	repgen.leaderState.h = reportContextHash(repgen.leaderState.q, aos)

	repgen.leaderState.phase = phaseReport
}

func (repgen *reportGenerationState) messageReport(msg MessageReport, sender commontypes.OracleID) {
	dropPrefix := "messageReport: dropping MessageReport due to "
	if msg.Epoch != repgen.e {
		repgen.logger.Debug(dropPrefix+"wrong epoch",
			commontypes.LogFields{"round": repgen.leaderState.r, "msgEpoch": msg.Epoch})
		return
	}
	if repgen.l != repgen.id {
		repgen.logger.Warn(dropPrefix+"not being leader of the current epoch",
			commontypes.LogFields{"round": repgen.leaderState.r})
		return
	}
	if msg.Round != repgen.leaderState.r {
		repgen.logger.Debug(dropPrefix+"wrong round",
			commontypes.LogFields{"round": repgen.leaderState.r, "msgRound": msg.Round})
		return
	}
	if repgen.leaderState.phase != phaseReport {
		repgen.logger.Debug(dropPrefix+"not being in report phase",
			commontypes.LogFields{"round": repgen.leaderState.r, "currentPhase": englishPhase[repgen.leaderState.phase]})
		return
	}
	if repgen.leaderState.report[sender] != nil {
		repgen.logger.Warn(dropPrefix+"having already received sender's report",
			commontypes.LogFields{"round": repgen.leaderState.r, "sender": sender, "msg": msg})
		return
	}

	err := msg.AttestedReport.Verify(
		repgen.onchainKeyring,
		repgen.config.OracleIdentities[sender].OnchainPublicKey,
		types.ReportContext{
			repgen.leaderReportTimestamp(),
			repgen.leaderState.h,
		},
	)
	if err != nil {
		repgen.logger.Error("could not validate signature", commontypes.LogFields{
			"round": repgen.leaderState.r,
			"error": err,
			"msg":   msg,
		})
		return
	}

	repgen.leaderState.report[sender] = &msg.AttestedReport

	// upon exists R s.t. |{p_j ∈ P | report[j]=(R,·)}| > f ∧ phase = REPORT
	{ // FUTUREWORK: make it non-quadratic time
		ass := []types.AttributedOnchainSignature{}
		for id, report := range repgen.leaderState.report {
			if report == nil {
				continue
			}
			if report.EqualExceptSignature(msg.AttestedReport) {
				ass = append(ass, types.AttributedOnchainSignature{
					report.Signature,
					commontypes.OracleID(id),
				})
			} else if !report.Skip && !msg.AttestedReport.Skip { // oracles may commonly disagree on whether to skip, no need to warn about that
				repgen.logger.Warn("received disparate reports messages", commontypes.LogFields{
					"round":          repgen.leaderState.r,
					"previousReport": report,
					"msgReport":      msg,
				})
			}
		}

		if repgen.reportQuorum <= len(ass) {
			if !msg.AttestedReport.Skip {
				repgen.netSender.Broadcast(MessageFinal{
					repgen.e,
					repgen.leaderState.r,
					repgen.leaderState.h,
					AttestedReportMany{
						msg.AttestedReport.Report,
						ass,
					},
				})
			}
			repgen.leaderState.phase = phaseFinal
			repgen.startRound()
		}
	}
}
