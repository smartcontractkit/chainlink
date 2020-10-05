package protocol

import (
	"sort"
	"time"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)


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

func (repgen *reportGenerationState) eventTRoundTimeout() {
	repgen.startRound()
}

func (repgen *reportGenerationState) startRound() {
	rPlusOne := repgen.leaderState.r + 1
	if rPlusOne <= repgen.leaderState.r {
		repgen.logger.Error("ReportGeneration: round overflows, cannot start new round", nil)
		return
	}
	repgen.leaderState.r = rPlusOne
	repgen.leaderState.observe = make([]*MessageObserve, repgen.config.N())
	repgen.leaderState.phase = phaseObserve
	repgen.netSender.Broadcast(MessageObserveReq{Epoch: repgen.e, Round: repgen.leaderState.r})
	repgen.leaderState.tRound = time.After(repgen.config.DeltaRound)
}

func (repgen *reportGenerationState) messageObserve(msg MessageObserve, sender types.OracleID) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("Got MessageObserve for wrong epoch", types.LogFields{
			"epoch":    repgen.e,
			"round":    repgen.leaderState.r,
			"sender":   sender,
			"msgEpoch": msg.Epoch,
			"msgRound": msg.Round,
		})
		return
	}

	if repgen.l != repgen.id {
		repgen.logger.Warn("Non-leader received MessageObserve", types.LogFields{
			"sender": sender,
			"msg":    msg,
		})
		return
	}

	if msg.Round != repgen.leaderState.r {
		repgen.logger.Debug("Got MessageObserve for wrong round", types.LogFields{
			"epoch":    repgen.e,
			"round":    repgen.leaderState.r,
			"sender":   sender,
			"msgEpoch": msg.Epoch,
			"msgRound": msg.Round,
		})
		return
	}

	if repgen.leaderState.phase != phaseObserve && repgen.leaderState.phase != phaseGrace {
		repgen.logger.Debug("received MessageObserve after grace phase", nil)
		return
	}

	if repgen.leaderState.observe[sender] != nil {
						repgen.logger.Debug("already sent an observation", types.LogFields{
			"sender": sender})
		return
	}

	senderPublicKey := signature.OffchainPublicKey(
		repgen.config.OracleIdentities[sender].OffchainPublicKey)
	if !senderPublicKey.Verify(msg.Obs.wireMessage(), msg.Obs.Sig) {
		repgen.logger.Warn("MessageObserve has invalid signature", types.LogFields{
			"round":  repgen.leaderState.r,
			"sender": sender,
			"msg":    msg,
		})
		return
	}

	repgen.logger.Debug("MessageObserve has valid signature", types.LogFields{
		"round":    repgen.leaderState.r,
		"sender":   sender,
		"msgEpoch": msg.Epoch,
		"msgRound": msg.Round,
	})

	repgen.leaderState.observe[sender] = &msg

		switch repgen.leaderState.phase {
	case phaseObserve:
		observationCount := 0 		for _, obs := range repgen.leaderState.observe {
			if obs != nil {
				observationCount++
			}
		}
		repgen.logger.Debug("One more observation", types.LogFields{
			"observationCount":         observationCount,
			"requiredObservationCount": (2 * repgen.config.F) + 1,
		})
		if observationCount > 2*repgen.config.F {
						repgen.logger.Debug("starting observation grace period", nil)
			repgen.leaderState.tGrace = time.After(repgen.config.DeltaGrace)
			repgen.leaderState.phase = phaseGrace
		}
	case phaseGrace:
		repgen.logger.Debug("accepted extra observation during grace period", nil)
	}
}

func (repgen *reportGenerationState) eventTGraceTimeout() {
	if repgen.leaderState.phase != phaseGrace {
		repgen.logger.Error("leader's phase conflicts tGrace timeout", types.LogFields{
			"phase": englishPhase[repgen.leaderState.phase],
		})
		return
	}
	observations := []Observation{}
	for _, msgObs := range repgen.leaderState.observe {
		if msgObs != nil {
			observations = append(observations, msgObs.Obs)
		}
	}
	sort.Slice(observations, func(i, j int) bool {
		return observations[i].Value.Less(observations[j].Value)
	})
	repgen.netSender.Broadcast(MessageReportReq{
		Epoch:        repgen.e,
		Round:        repgen.leaderState.r,
		Observations: observations,
	})
	repgen.leaderState.phase = phaseReport
}

func (repgen *reportGenerationState) messageReport(msg MessageReport, sender types.OracleID) {
	dropPrefix := "messageReport: dropping MessageReport due to "
	if msg.Epoch != repgen.e {
		repgen.logger.Debug(dropPrefix+"wrong epoch",
			types.LogFields{"epoch": repgen.e, "msgEpoch": msg.Epoch})
		return
	}
	if repgen.l != repgen.id {
		repgen.logger.Warn(dropPrefix+"not being leader of the current epoch",
			types.LogFields{"leader": repgen.l})
		return
	}
	if msg.Round != repgen.leaderState.r {
		repgen.logger.Debug(dropPrefix+"wrong round",
			types.LogFields{"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.leaderState.phase != phaseReport {
		repgen.logger.Debug(dropPrefix+"not being in report phase",
			types.LogFields{"currentPhase": englishPhase[repgen.leaderState.phase]})
		return
	}

	a := types.OnChainSigningAddress(repgen.config.OracleIdentities[sender].OnChainSigningAddress)
	err := msg.ContractReport.verify(a)
	if err != nil {
		repgen.logger.Error("could not validate signature", types.LogFields{
			"error": err,
			"msg":   msg,
		})
		return
	}

	repgen.leaderState.report[sender] = &msg

		{ 		sigs := [][]byte{}
		for _, msgR := range repgen.leaderState.report {
			if msgR == nil {
				continue
			}
			if msgR.ContractReport.Equal(msg.ContractReport) {
				sigs = append(sigs, msgR.ContractReport.Sig)
			} else {
				repgen.logger.Warn("received disparate contractReport messages", types.LogFields{
					"msg1": msgR,
					"msg2": msg,
				})
			}
		}

		if repgen.config.F < len(sigs) {
			repgen.netSender.Broadcast(MessageFinal{
				repgen.e,
				repgen.l,
				repgen.leaderState.r,
				ContractReportWithSignatures{
					ContractReport: msg.ContractReport,
					Signatures:     sigs,
				},
			})
			repgen.leaderState.phase = phaseFinal
		}
	}
}
