package protocol

import (
	"context"
	"math/big"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

func (repgen *reportGenerationState) followerReportContext() ReportContext {
	return ReportContext{repgen.config.ConfigDigest, repgen.e, repgen.followerState.r}
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
	repgen.followerState.sentEcho = nil
	repgen.followerState.sentReport = false
	repgen.followerState.completedRound = false
	repgen.followerState.receivedEcho = make([]bool, repgen.config.N())

	repgen.telemetrySender.RoundStarted(
		repgen.config.ConfigDigest,
		repgen.e,
		repgen.followerState.r,
		repgen.l,
	)

	value := repgen.observeValue()
	if value.IsMissingValue() {
		// Failed to get data from API, nothing to be done...
		// No need to log because observeValue already does
		return
	}

	so, err := MakeSignedObservation(value, repgen.followerReportContext(), repgen.privateKeys.SignOffChain)
	if err != nil {
		repgen.logger.Error("messageObserveReq: could not make SignedObservation observation", commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return
	}

	if err := so.Verify(repgen.followerReportContext(), repgen.privateKeys.PublicKeyOffChain()); err != nil {
		repgen.logger.Error("MakeSignedObservation produced invalid signature:", commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return
	}

	repgen.logger.Debug("sent observation to leader", commontypes.LogFields{
		"round":       repgen.followerState.r,
		"observation": value,
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

	if repgen.shouldReport(msg.AttributedSignedObservations) {
		attributedValues := make([]AttributedObservation, len(msg.AttributedSignedObservations))
		for i, aso := range msg.AttributedSignedObservations {
			// Observation/Observer attribution is verified by checking signature in verifyReportReq
			attributedValues[i] = AttributedObservation{
				aso.SignedObservation.Observation,
				aso.Observer,
			}
		}
		report, err := MakeAttestedReportOne(
			attributedValues,
			repgen.followerReportContext(),
			repgen.privateKeys.SignOnChain,
		)
		if err != nil {
			// Can't really do much here except logging as much detail as possible to
			// aid reproduction, and praying it won't happen again
			repgen.logger.Error("messageReportReq: failed to sign report", commontypes.LogFields{
				"round":  repgen.followerState.r,
				"error":  err,
				"id":     repgen.id,
				"report": report,
				"pubkey": repgen.privateKeys.PublicKeyAddressOnChain(),
			})
			return
		}

		{
			err := report.Verify(repgen.followerReportContext(), repgen.privateKeys.PublicKeyAddressOnChain())
			if err != nil {
				repgen.logger.Error("could not verify my own signature", commontypes.LogFields{
					"round":  repgen.followerState.r,
					"error":  err,
					"id":     repgen.id,
					"report": report, // includes sig
					"pubkey": repgen.privateKeys.PublicKeyAddressOnChain()})
				return
			}
		}

		repgen.followerState.sentReport = true
		repgen.netSender.SendTo(
			MessageReport{
				repgen.e,
				repgen.followerState.r,
				report,
			},
			repgen.l,
		)
	} else {
		repgen.completeRound()
	}
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
	if repgen.followerState.sentEcho != nil {
		repgen.logger.Debug("MessageFinal after already sent MessageFinalEcho", nil)
		return
	}
	if !repgen.verifyAttestedReport(msg.Report, sender) {
		return
	}
	repgen.followerState.sentEcho = &msg.Report
	repgen.netSender.Broadcast(MessageFinalEcho{MessageFinal: msg})
}

// messageFinalEcho is called when the local oracle process receives a
// "final-echo" message. If the report it contains is valid and the round is not
// yet complete, it keeps track of how many such echos have been received, and
// invokes the "transmit" event when enough echos have been seen to ensure that
// at least one (other?) honest node is broadcasting this report. This completes
// the round, from the local oracle's perspective.
func (repgen *reportGenerationState) messageFinalEcho(msg MessageFinalEcho,
	sender commontypes.OracleID,
) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("wrong epoch from MessageFinalEcho", commontypes.LogFields{
			"round": repgen.followerState.r, "msgEpoch": msg.Epoch, "sender": sender})
		return
	}
	if msg.Round != repgen.followerState.r {
		repgen.logger.Debug("wrong round from MessageFinalEcho", commontypes.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round, "sender": sender})
		return
	}
	if repgen.followerState.receivedEcho[sender] {
		repgen.logger.Warn("extra MessageFinalEcho received", commontypes.LogFields{
			"round": repgen.followerState.r, "sender": sender})
		return
	}
	if repgen.followerState.completedRound {
		repgen.logger.Debug("received final echo after round completion", nil)
		return
	}
	if !repgen.verifyAttestedReport(msg.Report, sender) { // if verify-attested-report(O) then
		// log messages are in verifyAttestedReport
		return
	}
	repgen.followerState.receivedEcho[sender] = true // receivedecho[j] ← true

	if repgen.followerState.sentEcho == nil { // if sentecho = ⊥ then
		repgen.followerState.sentEcho = &msg.Report // sentecho ← O
		repgen.netSender.Broadcast(msg)             // send [ FINALECHO , r, O] to all p_j ∈ P
	}

	// upon {p j ∈ P | receivedecho[j] = true} > f ∧ ¬completedround do
	{
		count := 0 // FUTUREWORK: Make this constant-time with a stateful counter
		for _, receivedEcho := range repgen.followerState.receivedEcho {
			if receivedEcho {
				count++
			}
		}
		if repgen.config.F < count {
			select {
			case repgen.chReportGenerationToTransmission <- EventTransmit{
				repgen.e,
				repgen.followerState.r,
				*repgen.followerState.sentEcho,
			}:
			case <-repgen.ctx.Done():
			}
			repgen.completeRound()
		}
	}

}

// observeValue is called when the oracle needs to gather a fresh observation to
// send back to the current leader.
func (repgen *reportGenerationState) observeValue() observation.Observation {
	var value observation.Observation
	var err error
	// We don't trust datasource.Observe(ctx) to actually exit after the context
	// deadline. We want to make sure we don't wait too long in order to not
	// drop out of the protocol. Even if an instance cannot make observations,
	// it can still be useful, e.g. by signing reports.
	//
	// We pass a context with timeout DataSourceTimeout to datasource.Observe().
	// However, we block for DataSourceTimeout *plus* DataSourceGracePeriod to
	// provide the DataSource with time to do a little bit of work after we
	// cancel the context (e.g. computing a median). Expect the
	// DataSourceGracePeriod to be on the order of 100ms for real-world
	// deployments.
	ok := repgen.subprocesses.BlockForAtMost(
		repgen.ctx,
		repgen.localConfig.DataSourceTimeout+repgen.localConfig.DataSourceGracePeriod,
		func(ctx context.Context) {
			// As mentioned above, datasource still has a grace period to return a result
			// when this context is cancelled.
			warnCtx, cancel := context.WithTimeout(ctx, repgen.localConfig.DataSourceTimeout)
			defer cancel()
			var rawValue types.Observation
			rawValue, err = repgen.datasource.Observe(warnCtx)
			if err != nil {
				return
			}
			value, err = observation.MakeObservation((*big.Int)(rawValue))
		},
	)

	if !ok {
		repgen.logger.Warn("DataSource timed out", commontypes.LogFields{
			"round":   repgen.followerState.r,
			"timeout": repgen.localConfig.DataSourceTimeout,
		})
		return observation.Observation{}
	}

	if err != nil {
		repgen.logger.ErrorIfNotCanceled("ReportGeneration: DataSource errored", repgen.ctx, commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return observation.Observation{}
	}

	return value
}

func (repgen *reportGenerationState) shouldReport(observations []AttributedSignedObservation) bool {
	var resultTransmissionDetails struct {
		configDigest    types.ConfigDigest
		epoch           uint32
		round           uint8
		latestAnswer    types.Observation
		latestTimestamp time.Time
		err             error
	}
	var resultRoundRequested struct {
		configDigest types.ConfigDigest
		epoch        uint32
		round        uint8
		err          error
	}
	ok, oks := repgen.subprocesses.BlockForAtMostMany(repgen.ctx, repgen.localConfig.BlockchainTimeout,
		func(ctx context.Context) {
			resultTransmissionDetails.configDigest,
				resultTransmissionDetails.epoch,
				resultTransmissionDetails.round,
				resultTransmissionDetails.latestAnswer,
				resultTransmissionDetails.latestTimestamp,
				resultTransmissionDetails.err =
				repgen.contractTransmitter.LatestTransmissionDetails(ctx)
		}, func(ctx context.Context) {
			resultRoundRequested.configDigest,
				resultRoundRequested.epoch,
				resultRoundRequested.round,
				resultRoundRequested.err =
				repgen.contractTransmitter.LatestRoundRequested(
					ctx,
					repgen.config.DeltaC,
				)
		},
	)
	if !ok {
		fnNames := []string{}
		if !oks[0] {
			fnNames = append(fnNames, "LatestTransmissionDetails()")
		}
		if !oks[1] {
			fnNames = append(fnNames, "LatestRoundRequested()")
		}
		repgen.logger.Error("shouldReport: blockchain interaction timed out, returning true", commontypes.LogFields{
			"round":             repgen.followerState.r,
			"timedOutFunctions": fnNames,
		})
		// Err on the side of creating too many reports. For instance, the Ethereum node
		// might be down, but that need not prevent us from still contributing to the
		// protocol.
		return true
	}

	if resultTransmissionDetails.err != nil {
		repgen.logger.ErrorIfNotCanceled("shouldReport: Error during LatestTransmissionDetails", repgen.ctx, commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": resultTransmissionDetails.err,
		})
		// Err on the side of creating too many reports. For instance, the Ethereum node
		// might be down, but that need not prevent us from still contributing to the
		// protocol.
		return true
	}
	if resultRoundRequested.err != nil {
		repgen.logger.Error("shouldReport: Error during LatestRoundRequested", commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": resultRoundRequested.err,
		})
		// Err on the side of creating too many reports. For instance, the Ethereum node
		// might be down, but that need not prevent us from still contributing to the
		// protocol.
		return true
	}

	answer, err := observation.MakeObservation(resultTransmissionDetails.latestAnswer)
	if err != nil {
		repgen.logger.Error("shouldReport: Error during observation.MakeObservation", commontypes.LogFields{
			"round": repgen.followerState.r,
			"error": err,
		})
		return false
	}

	alphaPPB, deltaC := repgen.config.AlphaPPB, repgen.config.DeltaC
	if override := repgen.configOverrider.ConfigOverride(); override != nil {
		repgen.logger.Debug("shouldReport: using overrides for alphaPPB and deltaC", commontypes.LogFields{
			"round":              repgen.followerState.r,
			"alphaPPB":           alphaPPB,
			"deltaC":             deltaC,
			"overriddenAlphaPPB": override.AlphaPPB,
			"overriddenDeltaC":   override.DeltaC,
		})
		alphaPPB = override.AlphaPPB
		deltaC = override.DeltaC
	}

	initialRound := // Is this the first round for this configuration?
		resultTransmissionDetails.configDigest == repgen.config.ConfigDigest &&
			resultTransmissionDetails.epoch == 0 &&
			resultTransmissionDetails.round == 0
	deviation := // Has the result changed enough to merit a new report?
		observations[len(observations)/2].SignedObservation.Observation.
			Deviates(answer, alphaPPB)
	deltaCTimeout := // Has enough time passed since the last report, to merit a new one?
		resultTransmissionDetails.latestTimestamp.Add(deltaC).
			Before(time.Now())
	unfulfilledRequest := // Has a new report been requested explicitly?
		resultRoundRequested.configDigest == repgen.config.ConfigDigest &&
			!(EpochRound{resultRoundRequested.epoch, resultRoundRequested.round}).
				Less(EpochRound{resultTransmissionDetails.epoch, resultTransmissionDetails.round})

	logger := repgen.logger.MakeChild(commontypes.LogFields{
		"round":                     repgen.followerState.r,
		"initialRound":              initialRound,
		"deviation":                 deviation,
		"deltaC":                    deltaC,
		"deltaCTimeout":             deltaCTimeout,
		"lastTransmissionTimestamp": resultTransmissionDetails.latestTimestamp,
		"unfulfilledRequest":        unfulfilledRequest,
	})

	// The following is more succinctly expressed as a disjunction, but breaking
	// the branches up into their own conditions makes it easier to check that
	// each branch is tested, and also allows for more expressive log messages
	if initialRound {
		logger.Info("shouldReport: yes, because it's the first round of the first epoch", commontypes.LogFields{
			"result": true,
		})
		return true
	}
	if deviation {
		logger.Info("shouldReport: yes, because new median deviates sufficiently from current onchain value", commontypes.LogFields{
			"result": true,
		})
		return true
	}
	if deltaCTimeout {
		logger.Info("shouldReport: yes, because deltaC timeout since last onchain report", commontypes.LogFields{
			"result": true,
		})
		return true
	}
	if unfulfilledRequest {
		logger.Info("shouldReport: yes, because a new report has been explicitly requested", commontypes.LogFields{
			"result": true,
		})
		return true
	}
	logger.Info("shouldReport: no", commontypes.LogFields{"result": false})
	return false
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

// verifyReportReq errors unless the reports observations are sorted, its
// signatures are all correct given the current round/epoch/config, and from
// distinct oracles, and there are more than 2f observations.
func (repgen *reportGenerationState) verifyReportReq(msg MessageReportReq) error {
	// check sortedness
	if !sort.SliceIsSorted(msg.AttributedSignedObservations,
		func(i, j int) bool {
			return msg.AttributedSignedObservations[i].SignedObservation.Observation.Less(msg.AttributedSignedObservations[j].SignedObservation.Observation)
		}) {
		return errors.Errorf("messages not sorted by value")
	}

	// check signatures and signature distinctness
	{
		counted := map[commontypes.OracleID]bool{}
		for _, obs := range msg.AttributedSignedObservations {
			// NOTE: OracleID is untrusted, therefore we _must_ bounds check it first
			numOracles := len(repgen.config.OracleIdentities)
			if int(obs.Observer) < 0 || numOracles <= int(obs.Observer) {
				return errors.Errorf("given oracle ID of %v is out of bounds (only "+
					"have %v public keys)", obs.Observer, numOracles)
			}
			if counted[obs.Observer] {
				return errors.Errorf("duplicate observation by oracle id %v", obs.Observer)
			} else {
				counted[obs.Observer] = true
			}
			observerOffchainPublicKey := repgen.config.OracleIdentities[obs.Observer].OffchainPublicKey
			if err := obs.SignedObservation.Verify(repgen.followerReportContext(), observerOffchainPublicKey); err != nil {
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

// verifyAttestedReport returns true iff the signatures on msg are valid
// signatures by oracle participants
func (repgen *reportGenerationState) verifyAttestedReport(
	report AttestedReportMany, sender commontypes.OracleID,
) bool {
	if len(report.Signatures) <= repgen.config.F {
		repgen.logger.Warn("verifyAttestedReport: dropping final report because "+
			"it has too few signatures", commontypes.LogFields{"sender": sender,
			"numSignatures": len(report.Signatures), "F": repgen.config.F})
		return false
	}

	keys := make(signature.EthAddresses)
	for oid, id := range repgen.config.OracleIdentities {
		keys[types.OnChainSigningAddress(id.OnChainSigningAddress)] =
			commontypes.OracleID(oid)
	}

	err := report.VerifySignatures(repgen.followerReportContext(), keys)
	if err != nil {
		repgen.logger.Error("could not validate signatures on final report",
			commontypes.LogFields{
				"round":  repgen.followerState.r,
				"error":  err,
				"report": report,
				"sender": sender,
			})
		return false
	}
	return true
}
