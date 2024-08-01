package protocol

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type outgenLeaderPhase string

const (
	outgenLeaderPhaseUnknown        outgenLeaderPhase = "unknown"
	outgenLeaderPhaseNewEpoch       outgenLeaderPhase = "newEpoch"
	outgenLeaderPhaseSentEpochStart outgenLeaderPhase = "sentEpochStart"
	outgenLeaderPhaseSentRoundStart outgenLeaderPhase = "sentRoundStart"
	outgenLeaderPhaseGrace          outgenLeaderPhase = "grace"
	outgenLeaderPhaseSentProposal   outgenLeaderPhase = "sentProposal"
)

func (outgen *outcomeGenerationState[RI]) messageEpochStartRequest(msg MessageEpochStartRequest[RI], sender commontypes.OracleID) {
	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageEpochStartRequest for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if outgen.sharedState.l != outgen.id {
		outgen.logger.Warn("dropping MessageEpochStartRequest to non-leader", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseNewEpoch {
		outgen.logger.Debug("dropping MessageEpochStartRequest for wrong phase", commontypes.LogFields{
			"sender": sender,
			"phase":  outgen.leaderState.phase,
		})
		return
	}

	if outgen.leaderState.epochStartRequests[sender] != nil {
		outgen.logger.Warn("dropping duplicate MessageEpochStartRequest", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	outgen.leaderState.epochStartRequests[sender] = &epochStartRequest[RI]{}

	if err := msg.SignedHighestCertifiedTimestamp.Verify(
		outgen.ID(),
		outgen.config.OracleIdentities[sender].OffchainPublicKey,
	); err != nil {
		outgen.leaderState.epochStartRequests[sender].bad = true
		outgen.logger.Warn("MessageEpochStartRequest.SignedHighestCertifiedTimestamp is invalid", commontypes.LogFields{
			"sender": sender,
			"error":  err,
		})
		return
	}

	// Note that the MessageEpochStartRequest might still be invalid, e.g. if its HighestCertified is invalid.
	outgen.logger.Debug("got MessageEpochStartRequest with valid SignedHighestCertifiedTimestamp", commontypes.LogFields{
		"sender":                       sender,
		"msgHighestCertifiedTimestamp": msg.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp,
	})

	outgen.leaderState.epochStartRequests[sender].message = msg

	if len(outgen.leaderState.epochStartRequests) < outgen.config.ByzQuorumSize() {
		return
	}

	goodCount := 0
	var maxSender *commontypes.OracleID
	for sender, epochStartRequest := range outgen.leaderState.epochStartRequests {
		if epochStartRequest.bad {
			continue
		}
		goodCount++

		if maxSender == nil || outgen.leaderState.epochStartRequests[*maxSender].message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp.Less(epochStartRequest.message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp) {
			sender := sender
			maxSender = &sender
		}
	}

	if maxSender == nil || goodCount < outgen.config.ByzQuorumSize() {
		return
	}

	maxRequest := outgen.leaderState.epochStartRequests[*maxSender]

	if maxRequest.message.HighestCertified.Timestamp() != maxRequest.message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp {
		maxRequest.bad = true
		outgen.logger.Warn("timestamp mismatch in MessageEpochStartRequest", commontypes.LogFields{
			"sender":                          *maxSender,
			"highestCertified.Timestamp":      maxRequest.message.HighestCertified.Timestamp(),
			"signedHighestCertifiedTimestamp": maxRequest.message.SignedHighestCertifiedTimestamp.HighestCertifiedTimestamp,
		})
		return
	}

	if err := maxRequest.message.HighestCertified.Verify(
		outgen.config.ConfigDigest,
		outgen.config.OracleIdentities,
		outgen.config.ByzQuorumSize(),
	); err != nil {
		maxRequest.bad = true
		outgen.logger.Warn("MessageEpochStartRequest.HighestCertified is invalid", commontypes.LogFields{
			"sender": *maxSender,
			"error":  err,
		})
		return
	}

	highestCertifiedProof := make([]AttributedSignedHighestCertifiedTimestamp, 0, outgen.config.ByzQuorumSize())
	contributors := make([]commontypes.OracleID, 0, outgen.config.ByzQuorumSize())
	for sender, epochStartRequest := range outgen.leaderState.epochStartRequests {
		if epochStartRequest.bad {
			continue
		}
		highestCertifiedProof = append(highestCertifiedProof, AttributedSignedHighestCertifiedTimestamp{
			epochStartRequest.message.SignedHighestCertifiedTimestamp,
			sender,
		})
		contributors = append(contributors, sender)
		// not necessary, but hopefully helps with readability
		if len(highestCertifiedProof) == outgen.config.ByzQuorumSize() {
			break
		}
	}

	epochStartProof := EpochStartProof{
		maxRequest.message.HighestCertified,
		highestCertifiedProof,
	}

	// This is a sanity check to ensure that we only construct epochStartProofs that are actually valid.
	// This should never fail.
	if err := epochStartProof.Verify(outgen.ID(), outgen.config.OracleIdentities, outgen.config.ByzQuorumSize()); err != nil {
		outgen.logger.Critical("EpochStartProof is invalid, very surprising!", commontypes.LogFields{
			"proof": epochStartProof,
		})
		return
	}

	outgen.leaderState.phase = outgenLeaderPhaseSentEpochStart

	outgen.logger.Info("broadcasting MessageEpochStart", commontypes.LogFields{
		"contributors": contributors,
	})

	outgen.netSender.Broadcast(MessageEpochStart[RI]{
		outgen.sharedState.e,
		epochStartProof,
	})

	if epochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentLeaderRound()
	} else if commitQC, ok := epochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.commit(*commitQC)
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentLeaderRound()
	} else {
		prepareQc := epochStartProof.HighestCertified.(*CertifiedPrepare)
		outgen.sharedState.firstSeqNrOfEpoch = prepareQc.SeqNr + 1
		// We're dealing with a re-proposal from a failed epoch based on a
		// prepare qc.
		// We don't want to send MessageRoundStart.
	}
}

func (outgen *outcomeGenerationState[RI]) eventTRoundTimeout() {
	outgen.logger.Debug("TRound fired", commontypes.LogFields{
		"seqNr":          outgen.sharedState.seqNr,
		"committedSeqNr": outgen.sharedState.committedSeqNr,
		"deltaRound":     outgen.config.DeltaRound.String(),
	})
	outgen.startSubsequentLeaderRound()
}

func (outgen *outcomeGenerationState[RI]) startSubsequentLeaderRound() {
	if !outgen.leaderState.readyToStartRound {
		outgen.leaderState.readyToStartRound = true
		return
	}
	outgen.leaderState.readyToStartRound = false

	query, ok := callPluginFromOutcomeGeneration[types.Query](
		outgen,
		"Query",
		outgen.config.MaxDurationQuery,
		outgen.OutcomeCtx(outgen.sharedState.committedSeqNr+1),
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
			return outgen.reportingPlugin.Query(ctx, outctx)
		},
	)
	if !ok {
		return
	}

	outgen.leaderState.query = query

	outgen.leaderState.observations = map[commontypes.OracleID]*SignedObservation{}

	outgen.leaderState.tRound = time.After(outgen.config.DeltaRound)

	outgen.leaderState.phase = outgenLeaderPhaseSentRoundStart
	outgen.logger.Debug("broadcasting MessageRoundStart", commontypes.LogFields{
		"seqNr": outgen.sharedState.committedSeqNr + 1,
	})
	outgen.netSender.Broadcast(MessageRoundStart[RI]{
		outgen.sharedState.e,
		outgen.sharedState.committedSeqNr + 1,
		query,
	})
}

func (outgen *outcomeGenerationState[RI]) messageObservation(msg MessageObservation[RI], sender commontypes.OracleID) {

	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageObservation for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if outgen.sharedState.l != outgen.id {
		outgen.logger.Warn("dropping MessageObservation to non-leader", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if outgen.leaderState.phase != outgenLeaderPhaseSentRoundStart && outgen.leaderState.phase != outgenLeaderPhaseGrace {
		outgen.logger.Debug("dropping MessageObservation for wrong phase", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"phase":    outgen.leaderState.phase,
		})
		return
	}

	if msg.SeqNr != outgen.sharedState.seqNr {
		outgen.logger.Debug("dropping MessageObservation with invalid SeqNr", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if outgen.leaderState.observations[sender] != nil {
		outgen.logger.Warn("dropping duplicate MessageObservation", commontypes.LogFields{
			"sender": sender,
			"seqNr":  outgen.sharedState.seqNr,
		})
		return
	}

	if err := msg.SignedObservation.Verify(outgen.ID(), outgen.sharedState.seqNr, outgen.leaderState.query, outgen.config.OracleIdentities[sender].OffchainPublicKey); err != nil {
		outgen.logger.Warn("dropping MessageObservation carrying invalid SignedObservation", commontypes.LogFields{
			"sender": sender,
			"seqNr":  outgen.sharedState.seqNr,
			"error":  err,
		})
		return
	}

	err, ok := callPluginFromOutcomeGeneration[error](
		outgen,
		"ValidateObservation",
		0, // ValidateObservation is a pure function and should finish "instantly"
		outgen.OutcomeCtx(outgen.sharedState.seqNr),
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (error, error) {
			return outgen.reportingPlugin.ValidateObservation(
				outctx,
				outgen.leaderState.query,
				types.AttributedObservation{msg.SignedObservation.Observation, sender},
			), nil
		},
	)
	if !ok {
		outgen.logger.Error("dropping MessageObservation that could not be validated", commontypes.LogFields{
			"sender": sender,
			"seqNr":  outgen.sharedState.seqNr,
		})
		return
	}
	if err != nil {
		outgen.logger.Warn("dropping MessageObservation carrying invalid Observation", commontypes.LogFields{
			"sender": sender,
			"seqNr":  outgen.sharedState.seqNr,
			"error":  err,
		})
		return
	}

	quorum, ok := outgen.ObservationQuorum(outgen.leaderState.query)
	if !ok {
		return
	}

	outgen.logger.Debug("got valid MessageObservation", commontypes.LogFields{
		"sender": sender,
		"seqNr":  outgen.sharedState.seqNr,
	})

	outgen.leaderState.observations[sender] = &msg.SignedObservation

	observationCount := 0
	for _, so := range outgen.leaderState.observations {
		if so != nil {
			observationCount++
		}
	}
	if observationCount == quorum {
		outgen.logger.Debug("reached observation quorum, starting observation grace period", commontypes.LogFields{
			"seqNr":             outgen.sharedState.seqNr,
			"deltaGrace":        outgen.config.DeltaGrace.String(),
			"observationQuorum": quorum,
		})
		outgen.leaderState.phase = outgenLeaderPhaseGrace
		outgen.leaderState.tGrace = time.After(outgen.config.DeltaGrace)
	}
}

func (outgen *outcomeGenerationState[RI]) eventTGraceTimeout() {
	if outgen.leaderState.phase != outgenLeaderPhaseGrace {
		outgen.logger.Error("leader's phase conflicts TGrace timeout", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.leaderState.phase,
		})
		return
	}
	asos := make([]AttributedSignedObservation, 0, outgen.config.N())
	contributors := make([]commontypes.OracleID, 0, outgen.config.N())
	for oid, so := range outgen.leaderState.observations {
		if so != nil {
			asos = append(asos, AttributedSignedObservation{
				*so,
				commontypes.OracleID(oid),
			})
			contributors = append(contributors, commontypes.OracleID(oid))
		}
	}

	outgen.leaderState.phase = outgenLeaderPhaseSentProposal

	outgen.logger.Debug("broadcasting MessageProposal after TGrace fired", commontypes.LogFields{
		"seqNr":        outgen.sharedState.seqNr,
		"contributors": contributors,
		"deltaGrace":   outgen.config.DeltaGrace.String(),
	})
	outgen.netSender.Broadcast(MessageProposal[RI]{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		asos,
	})
}
