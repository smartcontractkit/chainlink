package protocol

import (
	"context"
	"math/big"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)


func (repgen *reportGenerationState) messageObserveReq(msg MessageObserveReq, sender types.OracleID) {
	dropPrefix := "messageObserveReq: dropping MessageObserveReq from "
			if msg.Epoch != repgen.e {
		repgen.logger.Debug(dropPrefix+"wrong epoch",
			types.LogFields{"epoch": repgen.e, "msgEpoch": msg.Epoch},
		)
		return
	}
	if sender != repgen.l {
				repgen.logger.Warn(dropPrefix+"non-leader",
			types.LogFields{"sender": sender, "leader": repgen.l})
		return
	}
	if msg.Round <= repgen.followerState.r {
				repgen.logger.Debug(dropPrefix+"earlier round",
			types.LogFields{"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if int64(repgen.config.RMax)+1 < int64(msg.Round) {
																repgen.logger.Warn(dropPrefix+"out of bounds round",
			types.LogFields{"rMax": repgen.config.RMax, "msgRound": msg.Round})
		return
	}

				repgen.followerState.r = msg.Round

	if repgen.followerState.r > repgen.config.RMax {
		repgen.logger.Debug(
			"messageReportReq: leader sent MessageObserveReq past its expiration "+
				"round. Time to change leader",
			types.LogFields{
				"messageRound": msg.Round,
				"roundMax":     repgen.config.RMax,
			})
		select {
		case repgen.chReportGenerationToPacemaker <- EventChangeLeader{}:
		case <-repgen.ctx.Done():
		}

										return
	}
					repgen.followerState.sentEcho = false
	repgen.followerState.sentReport = false
	repgen.followerState.completedRound = false
	repgen.followerState.finalEcho = make([]*MessageFinalEcho, repgen.config.N())

	value := repgen.observeValue()
	if value.IsMissingValue() {
						return
	}

	m := MessageObserve{
		Epoch: repgen.e,
		Round: repgen.followerState.r,
		Obs: Observation{
			Ctx:      repgen.NewReportingContext(),
			Value:    value,
			OracleID: repgen.id,
		},
	}
	wireMsg := m.Obs.wireMessage()
	var err error
	m.Obs.Sig, err = repgen.privateKeys.SignOffChain(wireMsg)
	if err != nil {
		repgen.logger.Error("messageReportReq: could not sign observation: %s", types.LogFields{
			"error": err,
			"round": repgen.followerState.r,
		})
		return
	}
	key := signature.OffchainPublicKey(repgen.privateKeys.PublicKeyOffChain())
	if !key.Verify(wireMsg, m.Obs.Sig) {
		repgen.logger.Error("sign produced invalid signature:", types.LogFields{
			"error": err,
			"round": repgen.followerState.r,
		})
		return
	}
	repgen.logger.Debug("sent observation to leader", types.LogFields{
		"epoch": repgen.e, "round": repgen.followerState.r, "leader": repgen.l,
		"observation": m.Obs.Value,
	})
	repgen.netSender.SendTo(m, repgen.l)
}

func (repgen *reportGenerationState) messageReportReq(msg MessageReportReq, sender types.OracleID) {
			if repgen.e != msg.Epoch {
		repgen.logger.Debug("messageReportReq from wrong epoch", types.LogFields{
			"epoch": repgen.e, "msgEpoch": msg.Epoch})
		return
	}
	if sender != repgen.l {
				repgen.logger.Warn("messageReportReq from non-leader", types.LogFields{
			"sender": sender, "leader": repgen.l})
		return
	}
	if repgen.followerState.r != msg.Round {
						repgen.logger.Debug("messageReportReq from wrong round", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.followerState.sentReport {
		repgen.logger.Warn("messageReportReq after report sent", nil)
		return
	}
	if repgen.followerState.completedRound {
		repgen.logger.Warn("messageReportReq after round completed", nil)
		return
	}
	err := repgen.verifyReportReq(msg)
	if err != nil {
		repgen.logger.Error("messageReportReq: could not validate report sent by leader", types.LogFields{
			"error": err,
			"msg":   msg,
		})
		return
	}

	if repgen.shouldReport(msg.Observations) {
		attributedValues := make([]OracleValue, len(msg.Observations))
		for i, obs := range msg.Observations {
			attributedValues[i] = OracleValue{
								ID:    obs.OracleID,
				Value: obs.Value,
			}
		}

		report := ContractReport{
			Ctx:    repgen.NewReportingContext(),
			Values: attributedValues,
		}
		if err := report.Sign(repgen.privateKeys.SignOnChain); err != nil {
									repgen.logger.Error("messageReportReq: failed to sign report",
				types.LogFields{"error": err, "id": repgen.id, "report": report,
					"pubkey": repgen.privateKeys.PublicKeyAddressOnChain()})
			return
		}

		{
			err := report.verify(repgen.privateKeys.PublicKeyAddressOnChain())
			if err != nil {
				repgen.logger.Error("could not verify my own signature", types.LogFields{
					"error": err, "id": repgen.id, "report": report, 					"pubkey": repgen.privateKeys.PublicKeyAddressOnChain()})
				return
			}
		}

		repgen.followerState.sentReport = true
		repgen.netSender.SendTo(
			MessageReport{
				Epoch:          repgen.e,
				Round:          repgen.followerState.r,
				ContractReport: report,
			},
			repgen.l,
		)
	} else {
		repgen.completeRound()
	}
}

func (repgen *reportGenerationState) messageFinal(
	msg MessageFinal, sender types.OracleID,
) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("wrong epoch from MessageFinal", types.LogFields{
			"epoch": repgen.e, "msgEpoch": msg.Epoch, "sender": sender})
		return
	}
	if msg.Round != repgen.followerState.r {
		repgen.logger.Debug("wrong round from MessageFinal", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if sender != repgen.l {
		repgen.logger.Warn("MessageFinal from non-leader", types.LogFields{
			"epoch": repgen.e, "msgEpoch": msg.Epoch, "sender": sender,
			"round": repgen.followerState.r, "msgRound": msg.Round})
		return
	}
	if repgen.followerState.sentEcho {
		repgen.logger.Debug("MessageFinal after already sent MessageFinalEcho", nil)
		return
	}
	if !repgen.verifyAttestedReport(msg.Report, sender) {
		return
	}
	repgen.followerState.sentEcho = true
	repgen.netSender.Broadcast(MessageFinalEcho{MessageFinal: msg})
}

func (repgen *reportGenerationState) messageFinalEcho(msg MessageFinalEcho,
	sender types.OracleID,
) {
	if msg.Epoch != repgen.e {
		repgen.logger.Debug("wrong epoch from MessageFinalEcho", types.LogFields{
			"epoch": repgen.e, "msgEpoch": msg.Epoch, "sender": sender})
		return
	}
	if msg.Round != repgen.followerState.r {
		repgen.logger.Debug("wrong round from MessageFinalEcho", types.LogFields{
			"round": repgen.followerState.r, "msgRound": msg.Round, "sender": sender})
		return
	}
	if repgen.followerState.completedRound {
		repgen.logger.Debug("received final echo after round completion", nil)
		return
	}
	if !repgen.verifyAttestedReport(msg.Report, sender) {
		return
	}
	repgen.followerState.finalEcho[sender] = &msg

	if !repgen.followerState.sentEcho {
		repgen.followerState.sentEcho = true
		repgen.netSender.Broadcast(msg)
	}

		{
		count := 0 		for _, msgFe := range repgen.followerState.finalEcho {
			if msgFe != nil && msgFe.Report.Equals(msg.Report) {
				count++
			}
		}
		if repgen.config.F < count {
			select {
			case repgen.chReportGenerationToTransmission <- EventTransmit{msg.Report}:
			case <-repgen.ctx.Done():
			}
			repgen.completeRound()
		}
	}

}

func (repgen *reportGenerationState) observeValue() observation.Observation {
	var value observation.Observation
	var err error
					ok := repgen.subprocesses.BlockForAtMost(
		repgen.ctx,
		repgen.localConfig.DataSourceTimeout,
		func(ctx context.Context) {
			var rawValue types.Observation
			rawValue, err = repgen.datasource.Observe(ctx)
			if err != nil {
				return
			}
			value, err = observation.MakeObservation((*big.Int)(rawValue))
		},
	)

	if !ok {
		repgen.logger.Error("DataSource timed out", types.LogFields{
			"timeout": repgen.localConfig.DataSourceTimeout,
		})
		return observation.Observation{}
	}

	if err != nil {
		repgen.logger.Error("DataSource errored", types.LogFields{"error": err})
		return observation.Observation{}
	}

	return value
}

func (repgen *reportGenerationState) shouldReport(observations []Observation) bool {
	ctx, cancel := context.WithTimeout(repgen.ctx, repgen.localConfig.BlockchainTimeout)
	defer cancel()
	contractConfigDigest, contractEpoch, contractRound, rawAnswer, timestamp,
		err := repgen.contractTransmitter.LatestTransmissionDetails(ctx)
	if err != nil {
		repgen.logger.Error("shouldReport: Error during LatestTransmissionDetails", types.LogFields{
			"error": err,
		})
								return true
	}

	answer, err := observation.MakeObservation(rawAnswer)
	if err != nil {
		repgen.logger.Error("shouldReport: Error during observation.NewObservation", types.LogFields{
			"error": err,
		})
		return false
	}

	initialRound := contractConfigDigest == repgen.config.ConfigDigest && contractEpoch == 0 && contractRound == 0
	deviation := observations[len(observations)/2].Value.Deviates(answer, repgen.config.Alpha)
	deltaCTimeout := timestamp.Add(repgen.config.DeltaC).Before(time.Now())
	result := initialRound || deviation || deltaCTimeout

	repgen.logger.Info("shouldReport: returning result", types.LogFields{
		"result":        result,
		"initialRound":  initialRound,
		"deviation":     deviation,
		"deltaCTimeout": deltaCTimeout,
		"round":         repgen.followerState.r,
	})

	return result
}

func (repgen *reportGenerationState) completeRound() {
	repgen.logger.Debug("ReportGeneration: completed round", types.LogFields{
		"round": repgen.followerState.r,
	})
	repgen.followerState.completedRound = true

	select {
	case repgen.chReportGenerationToPacemaker <- EventProgress{}:
	case <-repgen.ctx.Done():
	}
}

func (repgen *reportGenerationState) verifyReportReq(msg MessageReportReq) error {
		if !sort.SliceIsSorted(msg.Observations,
		func(i, j int) bool {
			return msg.Observations[i].Value.Less(msg.Observations[j].Value)
		}) {
		return errors.Errorf("messages not sorted by value")
	}

		{
		counted := map[types.OracleID]bool{}
		for _, obs := range msg.Observations {
						numOracles := len(repgen.config.OracleIdentities)
			if int(obs.OracleID) < 0 || numOracles <= int(obs.OracleID) {
				return errors.Errorf("given oracle ID of %v is out of bounds (only "+
					"have %v public keys)", obs.OracleID, numOracles)
			}
			if counted[obs.OracleID] {
				return errors.Errorf("duplicate observation by oracle id %v", obs.OracleID)
			} else {
				counted[obs.OracleID] = true
			}
			observerOffchainIdentity := repgen.config.OracleIdentities[obs.OracleID].OffchainPublicKey
			verificationKey := signature.OffchainPublicKey(observerOffchainIdentity)
			if !verificationKey.Verify(obs.wireMessage(), obs.Sig) {
				return errors.Errorf("invalid signature on %+v", obs)
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

func (repgen *reportGenerationState) verifyAttestedReport(
	report ContractReportWithSignatures, sender types.OracleID,
) bool {
	if len(report.Signatures) <= repgen.config.F {
		repgen.logger.Warn("verifyAttestedReport: dropping final report because "+
			"it has too few signatures", types.LogFields{"sender": sender,
			"numSignatures": len(report.Signatures), "F": repgen.config.F})
		return false
	}

	if err := repgen.verifyReportingContext(report.Ctx); err != nil {
		repgen.logger.Warn("invalid ReportingContext on MessageFinal", types.LogFields{
			"error":  err,
			"report": report,
		})
		return false
	}

	keys := make(signature.EthAddresses)
	for oid, id := range repgen.config.OracleIdentities {
		keys[types.OnChainSigningAddress(id.OnChainSigningAddress)] =
			types.OracleID(oid)
	}

	err := report.VerifySignatures(keys)
	if err != nil {
		repgen.logger.Error("could not validate signatures on final report",
			types.LogFields{"error": err, "report": report, "sender": sender})
		return false
	}
	return true
}

func (repgen *reportGenerationState) verifyReportingContext(reportCtx signature.ReportingContext) error {
	localCtx := repgen.NewReportingContext()
	if !localCtx.Equal(reportCtx) {
		return errors.New("contexts not equal")
	}
	return nil
}
