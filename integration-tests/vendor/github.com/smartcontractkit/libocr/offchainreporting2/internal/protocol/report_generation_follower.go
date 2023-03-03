package protocol

import (
	"context"
	"crypto/sha256"
	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func reportContextHash(query types.Query, aos []types.AttributedObservation) [32]byte {
	h := sha256.New()

	_ = binary.Write(h, binary.BigEndian, uint64(len(query)))
	_, _ = h.Write(query)

	_ = binary.Write(h, binary.BigEndian, uint64(len(aos)))
	for _, ao := range aos {
		_ = binary.Write(h, binary.BigEndian, uint64(len(ao.Observation)))
		_, _ = h.Write(ao.Observation)
		_, _ = h.Write([]byte{byte(ao.Observer)})
	}

	hash := h.Sum(nil)
	var output [32]byte
	if len(hash) != len(output) {
		// assertion
		panic("sha256 size mismatch")
	}
	copy(output[:], hash)
	return output
}

func (repgen *reportGenerationState) followerReportTimestamp() types.ReportTimestamp {
	return types.ReportTimestamp{repgen.config.ConfigDigest, repgen.e, repgen.followerState.r}
}

///////////////////////////////////////////////////////////
// Report Generation Follower (Algorithm 2)
///////////////////////////////////////////////////////////

// messageObserveReq is called when the oracle receives an observe-req message
// from the current leader. It responds with a message to the leader
// containing a fresh observation, as long as the message comes from the
// designated leader, pertains to the current valid round/epoch. It sets up the
// follower state used to track which the protocol is at in view of this
// follower.
func (repgen *reportGenerationState) messageObserveReq(msg MessageObserveReq, sender commontypes.OracleID) {
	dropPrefix := "messageObserveReq: dropping MessageObserveReq from "
	// Each of these guards get their own if statement, to ease test-coverage
	// verification
	if msg.Epoch != repgen.e {
		repgen.logger.Debug(dropPrefix+"wrong epoch",
			commontypes.LogFields{"round": repgen.followerState.r, "msgEpoch": msg.Epoch},
		)
		return
	}
	if sender != repgen.l {
		// warn because someone *from this epoch* is trying to usurp the lead
		repgen.logger.Warn(dropPrefix+"non-leader",
			commontypes.LogFields{"round": repgen.followerState.r, "sender": sender})
		return
	}
	if msg.Round <= repgen.followerState.r {
		// this can happen due to network delays, so it's only a debug output
		repgen.logger.Debug(dropPrefix+"earlier round",
			commontypes.LogFields{"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if int64(repgen.config.RMax)+1 < int64(msg.Round) {
		// This check prevents the leader from triggering the changeleader behavior
		// an arbitrary number of times (with round=RMax+2, RMax+3, ...) until
		// consensus on the next epoch has developed. Since advancing to the next
		// epoch involves broadcast network messages from all participants, a
		// malicious leader could otherwise potentially trigger a network flood.
		//
		// Warn because the leader should never send a round value this high
		repgen.logger.Warn(dropPrefix+"out of bounds round",
			commontypes.LogFields{"round": repgen.followerState.r, "rMax": repgen.config.RMax, "msgRound": msg.Round})
		return
	}

	repgen.followerState.r = msg.Round

	// msg.Round>0, because msg.Round>repgen.followerState.r, and the initial
	// value of repgen.followerState.r is zero. msg.Round<=repgen.config.RMax
	// thus ensures that at most RMax rounds are possible for the current leader.
	if repgen.followerState.r > repgen.config.RMax {
		repgen.logger.Debug(
			"messageReportReq: leader sent MessageObserveReq past its expiration "+
				"round. Time to change leader",
			commontypes.LogFields{
				"round":        repgen.followerState.r,
				"messageRound": msg.Round,
				"roundMax":     repgen.config.RMax,
			})
		select {
		case repgen.chReportGenerationToPacemaker <- EventChangeLeader{}:
		case <-repgen.ctx.Done():
		}

		return
	}
	// Re-initialize follower state, in preparation for the next round
	//
	// A malicious leader could reset these values by sending an observeReq later
	// in the protocol, but they would only harm themselves, because that would
	// advance the follower's view of the current epoch's round, which only
	// reduces the number of rounds the current leader has left to report in
	// without influencing the transmitted report in any way. (A valid observeReq
	// after the report has been passed to the transmission machinery is expected,
	// and has no impact on the transmission process.)
	repgen.followerState.sentReport = false
	repgen.followerState.completedRound = false

	repgen.telemetrySender.RoundStarted(
		repgen.config.ConfigDigest,
		repgen.e,
		repgen.followerState.r,
		repgen.l,
	)

	var o types.Observation
	{
		ctx, cancel := context.WithTimeout(repgen.ctx, repgen.config.MaxDurationObservation)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			repgen.config.MaxDurationObservation+ReportingPluginTimeoutWarningGracePeriod,
			func() {
				repgen.logger.Error("ReportGeneration: ReportingPlugin.Observation is taking too long", commontypes.LogFields{
					"round": repgen.followerState.r, "maxDuration": repgen.config.MaxDurationObservation,
				})
			},
		)

		var err error
		o, err = repgen.reportingPlugin.Observation(ctx, repgen.followerReportTimestamp(), msg.Query)

		ins.Stop()

		if err != nil {
			repgen.logger.ErrorIfNotCanceled("ReportGeneration: ReportingPlugin.Observation errored", repgen.ctx, commontypes.LogFields{
				"round": repgen.followerState.r,
				"error": err,
			})
			// failed to get data, nothing to be done
			return
		}
	}

	so, err := MakeSignedObservation(repgen.followerReportTimestamp(), msg.Query, o, repgen.offchainKeyring.OffchainSign)
	if err != nil {
		repgen.logger.Error("messageObserveReq: could not make SignedObservation observation", commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return
	}

	if err := so.Verify(repgen.followerReportTimestamp(), msg.Query, repgen.offchainKeyring.OffchainPublicKey()); err != nil {
		repgen.logger.Error("MakeSignedObservation produced invalid signature:", commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return
	}

	repgen.logger.Debug("sent observation to leader", commontypes.LogFields{
		"round":       repgen.followerState.r,
		"observation": o,
	})
	repgen.netSender.SendTo(MessageObserve{
		repgen.e,
		repgen.followerState.r,
		so,
	}, repgen.l)
}

// messageReportReq is called when an oracle receives a report-req message from
// the current leader. If the contained report validates, the oracle signs it
// and sends it back to the leader.
func (repgen *reportGenerationState) messageReportReq(msg MessageReportReq, sender commontypes.OracleID) {
	// Each of these guards get their own if statement, to ease test-coverage
	// verification
	if repgen.e != msg.Epoch {
		repgen.logger.Debug("messageReportReq from wrong epoch", commontypes.LogFields{
			"round":    repgen.followerState.r,
			"msgEpoch": msg.Epoch})
		return
	}
	if sender != repgen.l {
		// warn because someone *from this epoch* is trying to usurp the lead
		repgen.logger.Warn("messageReportReq from non-leader", commontypes.LogFields{
			"round": repgen.followerState.r, "sender": sender})
		return
	}
	if repgen.followerState.r != msg.Round {
		// too low a round can happen due to network delays, too high if the local
		// oracle loses network connectivity. So this is only debug-level
		repgen.logger.Debug("messageReportReq from wrong round", commontypes.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.followerState.sentReport {
		repgen.logger.Warn("messageReportReq after report sent", commontypes.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.followerState.completedRound {
		repgen.logger.Warn("messageReportReq after round completed", commontypes.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	err := repgen.verifyReportReq(msg)
	if err != nil {
		repgen.logger.Error("messageReportReq: could not validate report sent by leader", commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": err,
			"msg":   msg,
		})
		return
	}

	aos := []types.AttributedObservation{}
	for _, aso := range msg.AttributedSignedObservations {
		aos = append(aos, types.AttributedObservation{
			aso.SignedObservation.Observation,
			aso.Observer,
		})
	}

	var shouldReport bool
	var report types.Report
	{
		ctx, cancel := context.WithTimeout(repgen.ctx, repgen.config.MaxDurationReport)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			repgen.config.MaxDurationReport+ReportingPluginTimeoutWarningGracePeriod,
			func() {
				repgen.logger.Error("ReportGeneration: ReportingPlugin.Report is taking too long", commontypes.LogFields{
					"round": repgen.followerState.r, "maxDuration": repgen.config.MaxDurationReport,
				})
			},
		)

		var err error
		shouldReport, report, err = repgen.reportingPlugin.Report(
			ctx,
			repgen.followerReportTimestamp(),
			msg.Query,
			aos,
		)

		ins.Stop()

		if err != nil {
			repgen.logger.Error("messageReportReq: error in ReportingPlugin.Report", commontypes.LogFields{
				"round": repgen.followerState.r,
				"error": err,
				"id":    repgen.id,
				"msg":   msg,
			})
			return
		}
	}

	var attestedReport AttestedReportOne
	if shouldReport {
		attestedReport, err = MakeAttestedReportOneNoskip(
			types.ReportContext{
				repgen.followerReportTimestamp(),
				reportContextHash(msg.Query, aos),
			},
			report,
			repgen.onchainKeyring.Sign,
		)

		if err != nil {
			// Can't really do much here except logging as much detail as possible to
			// aid reproduction, and praying it won't happen again
			repgen.logger.Error("messageReportReq: failed to sign report", commontypes.LogFields{
				"round":          repgen.followerState.r,
				"error":          err,
				"id":             repgen.id,
				"attestedReport": attestedReport,
				"pubkey":         repgen.onchainKeyring.PublicKey(),
			})
			return
		}
	} else {
		attestedReport = MakeAttestedReportOneSkip()
		repgen.completeRound()
	}

	{
		err := attestedReport.Verify(
			repgen.onchainKeyring,
			repgen.onchainKeyring.PublicKey(),
			types.ReportContext{
				repgen.followerReportTimestamp(),
				reportContextHash(msg.Query, aos),
			},
		)

		if err != nil {
			repgen.logger.Error("could not verify my own signature", commontypes.LogFields{
				"round":          repgen.followerState.r,
				"error":          err,
				"id":             repgen.id,
				"attestedReport": attestedReport, // includes sig
				"pubkey":         repgen.onchainKeyring.PublicKey()})
			return
		}
	}

	repgen.followerState.sentReport = true
	repgen.netSender.SendTo(
		MessageReport{
			repgen.e,
			repgen.followerState.r,
			attestedReport,
		},
		repgen.l,
	)
}

// messageFinal is called when a "final" message is received for the local
// oracle process. If the report in the msg is valid, the oracle broadcasts it
// in a "final-echo" message.
func (repgen *reportGenerationState) messageFinal(
	msg MessageFinal, sender commontypes.OracleID,
) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("wrong epoch from MessageFinal", commontypes.LogFields{
			"round": repgen.followerState.r, "msgEpoch": msg.Epoch, "sender": sender})
		return
	}
	if msg.Round != repgen.followerState.r {
		repgen.logger.Debug("wrong round from MessageFinal", commontypes.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if sender != repgen.l {
		repgen.logger.Warn("MessageFinal from non-leader", commontypes.LogFields{
			"msgEpoch": msg.Epoch, "sender": sender,
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if err := msg.AttestedReport.VerifySignatures(
		repgen.reportQuorum,
		repgen.onchainKeyring,
		repgen.config.OracleIdentities,
		types.ReportContext{repgen.followerReportTimestamp(), msg.H},
	); err != nil {
		repgen.logger.Error("could not validate signatures on attested report in MessageFinal",
			commontypes.LogFields{
				"error":  err,
				"msg":    msg,
				"sender": sender,
			})
		return
	}

	select {
	case repgen.chReportGenerationToReportFinalization <- EventFinal{msg}:
	case <-repgen.ctx.Done():
	}
	repgen.completeRound()
}

// completeRound is called by the local report-generation process when the
// current round has been completed by either concluding that the report sent by
// the current leader should not be transmitted to the on-chain smart contract,
// or by initiating the transmission protocol with this report.
func (repgen *reportGenerationState) completeRound() {
	repgen.logger.Debug("ReportGeneration: completed round", commontypes.LogFields{
		"round": repgen.followerState.r,
	})
	repgen.followerState.completedRound = true

	select {
	case repgen.chReportGenerationToPacemaker <- EventProgress{}:
	case <-repgen.ctx.Done():
	}
}

// verifyReportReq errors unless its signatures are all correct given the
// current round/epoch/config, and from distinct oracles, and there are more
// than 2f observations.
func (repgen *reportGenerationState) verifyReportReq(msg MessageReportReq) error {
	// check signatures and signature distinctness
	{
		counted := map[commontypes.OracleID]bool{}
		for _, obs := range msg.AttributedSignedObservations {
			// NOTE: OracleID is untrusted, therefore we _must_ bounds check it first
			if int(obs.Observer) < 0 || repgen.config.N() <= int(obs.Observer) {
				return errors.Errorf("given oracle ID of %v is out of bounds (only "+
					"have %v public keys)", obs.Observer, repgen.config.N())
			}
			if counted[obs.Observer] {
				return errors.Errorf("duplicate observation by oracle id %v", obs.Observer)
			} else {
				counted[obs.Observer] = true
			}
			observerOffchainPublicKey := repgen.config.OracleIdentities[obs.Observer].OffchainPublicKey
			if err := obs.SignedObservation.Verify(repgen.followerReportTimestamp(), msg.Query, observerOffchainPublicKey); err != nil {
				return errors.Errorf("invalid signed observation: %s", err)
			}
		}
		bound := 2 * repgen.config.F
		if len(counted) <= bound {
			return errors.Errorf("not enough observations in report; got %d, "+
				"need more than %d", len(counted), bound)
		}
	}
	return nil
}
