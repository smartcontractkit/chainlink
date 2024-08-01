package protocol

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type outgenFollowerPhase string

const (
	outgenFollowerPhaseUnknown         outgenFollowerPhase = "unknown"
	outgenFollowerPhaseNewEpoch        outgenFollowerPhase = "newEpoch"
	outgenFollowerPhaseNewRound        outgenFollowerPhase = "newRound"
	outgenFollowerPhaseSentObservation outgenFollowerPhase = "sentObservation"
	outgenFollowerPhaseSentPrepare     outgenFollowerPhase = "sentPrepare"
	outgenFollowerPhaseSentCommit      outgenFollowerPhase = "sentCommit"
)

func (outgen *outcomeGenerationState[RI]) eventTInitialTimeout() {
	outgen.logger.Debug("TInitial fired", commontypes.LogFields{
		"seqNr":        outgen.sharedState.seqNr,
		"deltaInitial": outgen.config.DeltaInitial.String(),
	})
	select {
	case outgen.chOutcomeGenerationToPacemaker <- EventNewEpochRequest[RI]{}:
	case <-outgen.ctx.Done():
		return
	}
}

func (outgen *outcomeGenerationState[RI]) messageEpochStart(msg MessageEpochStart[RI], sender commontypes.OracleID) {
	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageEpochStart for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"msgEpoch": msg.Epoch,
		})
		return
	}

	if sender != outgen.sharedState.l {
		outgen.logger.Warn("dropping MessageEpochStart from non-leader", commontypes.LogFields{
			"sender": sender,
		})
		return
	}

	if outgen.followerState.phase != outgenFollowerPhaseNewEpoch {
		outgen.logger.Warn("dropping MessageEpochStart for wrong phase", commontypes.LogFields{
			"phase": outgen.followerState.phase,
		})
		return
	}

	{
		err := msg.EpochStartProof.Verify(
			outgen.ID(),
			outgen.config.OracleIdentities,
			outgen.config.ByzQuorumSize(),
		)
		if err != nil {
			outgen.logger.Warn("dropping MessageEpochStart containing invalid StartRoundQuorumCertificate", commontypes.LogFields{
				"error": err,
			})
			return
		}
	}

	outgen.followerState.tInitial = nil

	if msg.EpochStartProof.HighestCertified.IsGenesis() {
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentFollowerRound()
	} else if commitQC, ok := msg.EpochStartProof.HighestCertified.(*CertifiedCommit); ok {
		outgen.commit(*commitQC)
		outgen.sharedState.firstSeqNrOfEpoch = outgen.sharedState.committedSeqNr + 1
		outgen.startSubsequentFollowerRound()
	} else {
		// We're dealing with a re-proposal from a failed epoch

		prepareQc := msg.EpochStartProof.HighestCertified.(*CertifiedPrepare)

		// We don't know the actual inputs, so we always use the empty OutcomeInputsDigest
		// in case of a re-proposal.
		outcomeInputsDigest := OutcomeInputsDigest{}

		outcomeDigest := MakeOutcomeDigest(prepareQc.Outcome)

		prepareSignature, err := MakePrepareSignature(
			outgen.ID(),
			prepareQc.SeqNr,
			outcomeInputsDigest,
			outcomeDigest,
			outgen.offchainKeyring.OffchainSign,
		)
		if err != nil {
			outgen.logger.Critical("failed to sign Prepare", commontypes.LogFields{
				"seqNr": outgen.sharedState.seqNr,
				"error": err,
			})
			return
		}

		outgen.sharedState.firstSeqNrOfEpoch = prepareQc.SeqNr + 1
		outgen.sharedState.seqNr = prepareQc.SeqNr
		outgen.sharedState.observationQuorum = nil

		outgen.followerState.phase = outgenFollowerPhaseSentPrepare
		outgen.followerState.outcome = outcomeAndDigests{
			prepareQc.Outcome,
			outcomeInputsDigest,
			outcomeDigest,
		}
		outgen.logger.Debug("broadcasting MessagePrepare (reproposal)", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		outgen.netSender.Broadcast(MessagePrepare[RI]{
			outgen.sharedState.e,
			prepareQc.SeqNr,
			prepareSignature,
		})
	}
}

func (outgen *outcomeGenerationState[RI]) startSubsequentFollowerRound() {
	outgen.sharedState.seqNr = outgen.sharedState.committedSeqNr + 1
	outgen.sharedState.observationQuorum = nil

	outgen.followerState.phase = outgenFollowerPhaseNewRound
	outgen.followerState.query = nil
	outgen.followerState.outcome = outcomeAndDigests{}

	outgen.tryProcessRoundStartPool()
}

func (outgen *outcomeGenerationState[RI]) messageRoundStart(msg MessageRoundStart[RI], sender commontypes.OracleID) {
	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageRoundStart for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if sender != outgen.sharedState.l {
		outgen.logger.Warn("dropping MessageRoundStart from non-leader", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.roundStartPool.Put(msg.SeqNr, sender, msg); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessageRoundStart", commontypes.LogFields{
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"reason":   putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessageRoundStart", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})

	outgen.tryProcessRoundStartPool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessRoundStartPool() {
	if outgen.followerState.phase != outgenFollowerPhaseNewRound {
		outgen.logger.Debug("cannot process RoundStartPool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.roundStartPool.Entries(outgen.sharedState.seqNr)

	if poolEntries == nil || poolEntries[outgen.sharedState.l] == nil {

		outgen.logger.Debug("cannot process RoundStartPool, it's empty", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		return
	}

	if outgen.followerState.query != nil {
		outgen.logger.Warn("cannot process RoundStartPool, query already set", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		return
	}

	msg := poolEntries[outgen.sharedState.l].Item

	outgen.followerState.query = &msg.Query

	outctx := outgen.OutcomeCtx(outgen.sharedState.seqNr)

	outgen.telemetrySender.RoundStarted(
		outgen.config.ConfigDigest,
		outctx.Epoch,
		outctx.SeqNr,
		outctx.Round,
		outgen.sharedState.l,
	)

	o, ok := callPluginFromOutcomeGeneration[types.Observation](
		outgen,
		"Observation",
		outgen.config.MaxDurationObservation,
		outctx,
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Observation, error) {
			return outgen.reportingPlugin.Observation(ctx, outctx, *outgen.followerState.query)
		},
	)
	if !ok {
		return
	}

	so, err := MakeSignedObservation(outgen.ID(), outgen.sharedState.seqNr, msg.Query, o, outgen.offchainKeyring.OffchainSign)
	if err != nil {
		outgen.logger.Error("MakeSignedObservation returned error", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	if err := so.Verify(outgen.ID(), outgen.sharedState.seqNr, msg.Query, outgen.offchainKeyring.OffchainPublicKey()); err != nil {
		outgen.logger.Error("MakeSignedObservation produced invalid signature", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseSentObservation
	outgen.metrics.sentObservationsTotal.Inc()
	outgen.logger.Debug("sent MessageObservation to leader", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})
	outgen.netSender.SendTo(MessageObservation[RI]{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		so,
	}, outgen.sharedState.l)

	outgen.tryProcessProposalPool()
}

func (outgen *outcomeGenerationState[RI]) messageProposal(msg MessageProposal[RI], sender commontypes.OracleID) {
	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageProposal for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if sender != outgen.sharedState.l {
		outgen.logger.Warn("dropping MessageProposal from non-leader", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.proposalPool.Put(msg.SeqNr, sender, msg); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessageProposal", commontypes.LogFields{
			"seqNr":        outgen.sharedState.seqNr,
			"messageSeqNr": msg.SeqNr,
			"reason":       putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessageProposal", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})

	outgen.tryProcessProposalPool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessProposalPool() {
	if outgen.followerState.phase != outgenFollowerPhaseSentObservation {
		outgen.logger.Debug("cannot process ProposalPool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.proposalPool.Entries(outgen.sharedState.seqNr)

	if poolEntries == nil || poolEntries[outgen.sharedState.l] == nil {

		return
	}

	msg := poolEntries[outgen.sharedState.l].Item

	if msg.SeqNr <= outgen.sharedState.committedSeqNr {
		outgen.logger.Critical("MessageProposal contains invalid SeqNr", commontypes.LogFields{
			"msgSeqNr":       msg.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		return
	}

	attributedObservations := []types.AttributedObservation{}
	{
		quorum, ok := outgen.ObservationQuorum(*outgen.followerState.query)
		if !ok {
			return
		}

		if len(msg.AttributedSignedObservations) < quorum {
			outgen.logger.Warn("dropping MessageProposal that contains too few signed observations", commontypes.LogFields{
				"seqNr":                             outgen.sharedState.seqNr,
				"attributedSignedObservationsCount": len(msg.AttributedSignedObservations),
				"quorum":                            quorum,
			})
			return
		}
		seen := map[commontypes.OracleID]bool{}
		for _, aso := range msg.AttributedSignedObservations {
			if !(0 <= int(aso.Observer) && int(aso.Observer) <= outgen.config.N()) {
				outgen.logger.Warn("dropping MessageProposal that contains signed observation with invalid observer", commontypes.LogFields{
					"seqNr":           outgen.sharedState.seqNr,
					"invalidObserver": aso.Observer,
				})
				return
			}

			if seen[aso.Observer] {
				outgen.logger.Warn("dropping MessageProposal that contains duplicate signed observation", commontypes.LogFields{
					"seqNr": outgen.sharedState.seqNr,
				})
				return
			}

			seen[aso.Observer] = true

			if err := aso.SignedObservation.Verify(outgen.ID(), outgen.sharedState.seqNr, *outgen.followerState.query, outgen.config.OracleIdentities[aso.Observer].OffchainPublicKey); err != nil {
				outgen.logger.Warn("dropping MessageProposal that contains signed observation with invalid signature", commontypes.LogFields{
					"seqNr": outgen.sharedState.seqNr,
					"error": err,
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
						*outgen.followerState.query,
						types.AttributedObservation{aso.SignedObservation.Observation, aso.Observer},
					), nil
				},
			)
			if !ok {
				outgen.logger.Error("dropping MessageProposal containing observation that could not be validated", commontypes.LogFields{
					"seqNr":    outgen.sharedState.seqNr,
					"observer": aso.Observer,
				})
				return
			}
			if err != nil {
				outgen.logger.Warn("dropping MessageProposal that contains an invalid observation", commontypes.LogFields{
					"seqNr":    outgen.sharedState.seqNr,
					"error":    err,
					"observer": aso.Observer,
				})
				return
			}

			if aso.Observer == outgen.id {
				outgen.metrics.includedObservationsTotal.Inc()
			}

			attributedObservations = append(attributedObservations, types.AttributedObservation{
				aso.SignedObservation.Observation,
				aso.Observer,
			})
		}
	}

	outcomeInputsDigest := MakeOutcomeInputsDigest(
		outgen.ID(),
		outgen.sharedState.committedOutcome,
		outgen.sharedState.seqNr,
		*outgen.followerState.query,
		attributedObservations,
	)

	outcome, ok := callPluginFromOutcomeGeneration[ocr3types.Outcome](
		outgen,
		"Outcome",
		0, // Outcome is a pure function and should finish "instantly"
		outgen.OutcomeCtx(outgen.sharedState.seqNr),
		func(_ context.Context, outctx ocr3types.OutcomeContext) (ocr3types.Outcome, error) {
			return outgen.reportingPlugin.Outcome(outctx, *outgen.followerState.query, attributedObservations)
		},
	)
	if !ok {
		return
	}

	outcomeDigest := MakeOutcomeDigest(outcome)

	prepareSignature, err := MakePrepareSignature(
		outgen.ID(),
		msg.SeqNr,
		outcomeInputsDigest,
		outcomeDigest,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Critical("failed to sign Prepare", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseSentPrepare
	outgen.followerState.outcome = outcomeAndDigests{
		outcome,
		outcomeInputsDigest,
		outcomeDigest,
	}

	outgen.logger.Debug("broadcasting MessagePrepare", commontypes.LogFields{
		"seqNr": msg.SeqNr,
	})
	outgen.netSender.Broadcast(MessagePrepare[RI]{
		outgen.sharedState.e,
		msg.SeqNr,
		prepareSignature,
	})
}

func (outgen *outcomeGenerationState[RI]) messagePrepare(msg MessagePrepare[RI], sender commontypes.OracleID) {
	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessagePrepare for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.preparePool.Put(msg.SeqNr, sender, msg.Signature); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessagePrepare", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"reason":   putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessagePrepare", commontypes.LogFields{
		"sender":   sender,
		"seqNr":    outgen.sharedState.seqNr,
		"msgSeqNr": msg.SeqNr,
	})

	outgen.tryProcessPreparePool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessPreparePool() {
	if outgen.followerState.phase != outgenFollowerPhaseSentPrepare {
		outgen.logger.Debug("cannot process PreparePool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.preparePool.Entries(outgen.sharedState.seqNr)
	if len(poolEntries) < outgen.config.ByzQuorumSize() {

		return
	}

	for sender, preparePoolEntry := range poolEntries {
		if preparePoolEntry.Verified != nil {
			continue
		}
		err := preparePoolEntry.Item.Verify(
			outgen.ID(),
			outgen.sharedState.seqNr,
			outgen.followerState.outcome.InputsDigest,
			outgen.followerState.outcome.Digest,
			outgen.config.OracleIdentities[sender].OffchainPublicKey,
		)
		ok := err == nil
		outgen.followerState.preparePool.StoreVerified(outgen.sharedState.seqNr, sender, ok)
		if !ok {
			outgen.logger.Warn("dropping invalid MessagePrepare", commontypes.LogFields{
				"sender": sender,
				"seqNr":  outgen.sharedState.seqNr,
				"error":  err,
			})
		}
	}

	var prepareQuorumCertificate []AttributedPrepareSignature
	for sender, preparePoolEntry := range poolEntries {
		if preparePoolEntry.Verified != nil && *preparePoolEntry.Verified {
			prepareQuorumCertificate = append(prepareQuorumCertificate, AttributedPrepareSignature{
				preparePoolEntry.Item,
				sender,
			})
			if len(prepareQuorumCertificate) == outgen.config.ByzQuorumSize() {
				break
			}
		}
	}

	if len(prepareQuorumCertificate) < outgen.config.ByzQuorumSize() {
		return
	}

	commitSignature, err := MakeCommitSignature(
		outgen.ID(),
		outgen.sharedState.seqNr,
		outgen.followerState.outcome.Digest,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Critical("failed to sign Commit", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return
	}

	outgen.followerState.cert = &CertifiedPrepare{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		outgen.followerState.outcome.InputsDigest,
		outgen.followerState.outcome.Outcome,
		prepareQuorumCertificate,
	}
	if !outgen.persistCert() {
		return
	}

	outgen.followerState.phase = outgenFollowerPhaseSentCommit

	outgen.logger.Debug("broadcasting MessageCommit", commontypes.LogFields{
		"seqNr": outgen.sharedState.seqNr,
	})
	outgen.netSender.Broadcast(MessageCommit[RI]{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		commitSignature,
	})
}

func (outgen *outcomeGenerationState[RI]) messageCommit(msg MessageCommit[RI], sender commontypes.OracleID) {
	if msg.Epoch != outgen.sharedState.e {
		outgen.logger.Debug("dropping MessageCommit for wrong epoch", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgEpoch": msg.Epoch,
			"msgSeqNr": msg.SeqNr,
		})
		return
	}

	if putResult := outgen.followerState.commitPool.Put(msg.SeqNr, sender, msg.Signature); putResult != pool.PutResultOK {
		outgen.logger.Debug("dropping MessageCommit", commontypes.LogFields{
			"sender":   sender,
			"seqNr":    outgen.sharedState.seqNr,
			"msgSeqNr": msg.SeqNr,
			"reason":   putResult,
		})
		return
	}

	outgen.logger.Debug("pooled MessageCommit", commontypes.LogFields{
		"sender":   sender,
		"seqNr":    outgen.sharedState.seqNr,
		"msgSeqNr": msg.SeqNr,
	})

	outgen.tryProcessCommitPool()
}

func (outgen *outcomeGenerationState[RI]) tryProcessCommitPool() {
	if outgen.followerState.phase != outgenFollowerPhaseSentCommit {
		outgen.logger.Debug("cannot process CommitPool, wrong phase", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"phase": outgen.followerState.phase,
		})
		return
	}

	poolEntries := outgen.followerState.commitPool.Entries(outgen.sharedState.seqNr)
	if len(poolEntries) < outgen.config.ByzQuorumSize() {

		return
	}

	for sender, commitPoolEntry := range poolEntries {
		if commitPoolEntry.Verified != nil {
			continue
		}
		err := commitPoolEntry.Item.Verify(
			outgen.ID(),
			outgen.sharedState.seqNr,
			outgen.followerState.outcome.Digest,
			outgen.config.OracleIdentities[sender].OffchainPublicKey,
		)
		ok := err == nil
		commitPoolEntry.Verified = &ok
		if !ok {
			outgen.logger.Warn("dropping invalid MessageCommit", commontypes.LogFields{
				"sender": sender,
				"seqNr":  outgen.sharedState.seqNr,
				"error":  err,
			})
		}
	}

	var commitQuorumCertificate []AttributedCommitSignature
	for sender, commitPoolEntry := range poolEntries {
		if commitPoolEntry.Verified != nil && *commitPoolEntry.Verified {
			commitQuorumCertificate = append(commitQuorumCertificate, AttributedCommitSignature{
				commitPoolEntry.Item,
				sender,
			})
			if len(commitQuorumCertificate) == outgen.config.ByzQuorumSize() {
				break
			}
		}
	}

	if len(commitQuorumCertificate) < outgen.config.ByzQuorumSize() {
		return
	}

	outgen.commit(CertifiedCommit{
		outgen.sharedState.e,
		outgen.sharedState.seqNr,
		outgen.followerState.outcome.Outcome,
		commitQuorumCertificate,
	})
	if outgen.id == outgen.sharedState.l {
		outgen.metrics.ledCommittedRoundsTotal.Inc()
	}

	if uint64(outgen.config.RMax) <= outgen.sharedState.seqNr-outgen.sharedState.firstSeqNrOfEpoch+1 {
		outgen.logger.Debug("epoch has been going on for too long, sending EventChangeLeader to Pacemaker", commontypes.LogFields{
			"firstSeqNrOfEpoch": outgen.sharedState.firstSeqNrOfEpoch,
			"seqNr":             outgen.sharedState.seqNr,
			"rMax":              outgen.config.RMax,
		})
		select {
		case outgen.chOutcomeGenerationToPacemaker <- EventNewEpochRequest[RI]{}:
		case <-outgen.ctx.Done():
			return
		}
		return
	} else {
		outgen.logger.Debug("sending EventProgress to Pacemaker", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
		})
		select {
		case outgen.chOutcomeGenerationToPacemaker <- EventProgress[RI]{}:
		case <-outgen.ctx.Done():
			return
		}
	}

	outgen.startSubsequentFollowerRound()
	if outgen.id == outgen.sharedState.l {
		outgen.startSubsequentLeaderRound()
	}

	outgen.tryProcessRoundStartPool()
}

func (outgen *outcomeGenerationState[RI]) commit(commit CertifiedCommit) {
	if commit.SeqNr < outgen.sharedState.committedSeqNr {
		outgen.logger.Critical("assumption violation, commitSeqNr is less than committedSeqNr", commontypes.LogFields{
			"commitSeqNr":    commit.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		return
	}

	if commit.SeqNr <= outgen.sharedState.committedSeqNr {

		outgen.logger.Debug("skipping commit of already committed outcome", commontypes.LogFields{
			"commitSeqNr ":   commit.SeqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
	} else {
		outgen.followerState.cert = &commit
		if !outgen.persistCert() {
			return
		}

		outgen.sharedState.committedSeqNr = commit.SeqNr
		outgen.sharedState.committedOutcome = commit.Outcome
		outgen.metrics.committedSeqNr.Set(float64(commit.SeqNr))

		outgen.logger.Debug("✅ committed outcome", commontypes.LogFields{
			"seqNr": commit.SeqNr,
		})

		select {
		case outgen.chOutcomeGenerationToReportAttestation <- EventCommittedOutcome[RI]{commit}:
		case <-outgen.ctx.Done():
			return
		}
	}

	outgen.followerState.roundStartPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.proposalPool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.preparePool.ReapCompleted(outgen.sharedState.committedSeqNr)
	outgen.followerState.commitPool.ReapCompleted(outgen.sharedState.committedSeqNr)
}

func (outgen *outcomeGenerationState[RI]) persistCert() (ok bool) {
	ctx, cancel := context.WithTimeout(outgen.ctx, outgen.localConfig.DatabaseTimeout)
	defer cancel()
	if err := outgen.database.WriteCert(ctx, outgen.config.ConfigDigest, outgen.followerState.cert); err != nil {
		outgen.logger.Error("error persisting cert to database, cannot safely continue current round", commontypes.LogFields{
			"seqNr": outgen.sharedState.seqNr,
			"error": err,
		})
		return false
	}
	return true
}
