package protocol

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol/pool"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Identifies an instance of the outcome generation protocol
type OutcomeGenerationID struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint64
}

const futureMessageBufferSize = 10 // big enough for a couple of full rounds of outgen protocol
const poolSize = 3

func RunOutcomeGeneration[RI any](
	ctx context.Context,

	chNetToOutcomeGeneration <-chan MessageToOutcomeGenerationWithSender[RI],
	chPacemakerToOutcomeGeneration <-chan EventToOutcomeGeneration[RI],
	chOutcomeGenerationToPacemaker chan<- EventToPacemaker[RI],
	chOutcomeGenerationToReportAttestation chan<- EventToReportAttestation[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	reportingPlugin ocr3types.ReportingPlugin[RI],
	telemetrySender TelemetrySender,

	restoredCert CertifiedPrepareOrCommit,
) {
	outgen := outcomeGenerationState[RI]{
		ctx: ctx,

		chNetToOutcomeGeneration:               chNetToOutcomeGeneration,
		chPacemakerToOutcomeGeneration:         chPacemakerToOutcomeGeneration,
		chOutcomeGenerationToPacemaker:         chOutcomeGenerationToPacemaker,
		chOutcomeGenerationToReportAttestation: chOutcomeGenerationToReportAttestation,
		config:                                 config,
		database:                               database,
		id:                                     id,
		localConfig:                            localConfig,
		logger:                                 logger.MakeUpdated(commontypes.LogFields{"proto": "outgen"}),
		metrics:                                newOutcomeGenerationMetrics(metricsRegisterer, logger),
		netSender:                              netSender,
		offchainKeyring:                        offchainKeyring,
		reportingPlugin:                        reportingPlugin,
		telemetrySender:                        telemetrySender,
	}
	outgen.run(restoredCert)
}

type outcomeGenerationState[RI any] struct {
	ctx context.Context

	chNetToOutcomeGeneration               <-chan MessageToOutcomeGenerationWithSender[RI]
	chPacemakerToOutcomeGeneration         <-chan EventToOutcomeGeneration[RI]
	chOutcomeGenerationToPacemaker         chan<- EventToPacemaker[RI]
	chOutcomeGenerationToReportAttestation chan<- EventToReportAttestation[RI]
	config                                 ocr3config.SharedConfig
	database                               Database
	id                                     commontypes.OracleID
	localConfig                            types.LocalConfig
	logger                                 loghelper.LoggerWithContext
	metrics                                *outcomeGenerationMetrics
	netSender                              NetworkSender[RI]
	offchainKeyring                        types.OffchainKeyring
	reportingPlugin                        ocr3types.ReportingPlugin[RI]
	telemetrySender                        TelemetrySender

	bufferedMessages []*MessageBuffer[RI]
	leaderState      leaderState[RI]
	followerState    followerState[RI]
	sharedState      sharedState
}

type leaderState[RI any] struct {
	phase outgenLeaderPhase

	epochStartRequests map[commontypes.OracleID]*epochStartRequest[RI]

	readyToStartRound bool // TODO: explain meaning of this vs design doc
	tRound            <-chan time.Time

	query        types.Query
	observations map[commontypes.OracleID]*SignedObservation
	tGrace       <-chan time.Time
}

type epochStartRequest[RI any] struct {
	message MessageEpochStartRequest[RI]
	bad     bool
}

type followerState[RI any] struct {
	phase outgenFollowerPhase

	tInitial <-chan time.Time

	roundStartPool *pool.Pool[MessageRoundStart[RI]]

	query *types.Query

	proposalPool *pool.Pool[MessageProposal[RI]]

	outcome outcomeAndDigests

	// lock
	cert CertifiedPrepareOrCommit

	preparePool *pool.Pool[PrepareSignature]
	commitPool  *pool.Pool[CommitSignature]
}

type outcomeAndDigests struct {
	Outcome      ocr3types.Outcome
	InputsDigest OutcomeInputsDigest
	Digest       OutcomeDigest
}

type sharedState struct {
	e uint64               // Current epoch number
	l commontypes.OracleID // Current leader number

	firstSeqNrOfEpoch uint64
	seqNr             uint64
	observationQuorum *int
	committedSeqNr    uint64
	committedOutcome  ocr3types.Outcome
}

func (outgen *outcomeGenerationState[RI]) run(restoredCert CertifiedPrepareOrCommit) {
	outgen.logger.Info("OutcomeGeneration: running", nil)

	for i := 0; i < outgen.config.N(); i++ {
		outgen.bufferedMessages = append(outgen.bufferedMessages, NewMessageBuffer[RI](futureMessageBufferSize))
	}

	// Initialization
	outgen.leaderState = leaderState[RI]{
		outgenLeaderPhaseUnknown,
		map[commontypes.OracleID]*epochStartRequest[RI]{},
		false,
		nil,
		nil,
		nil,
		nil,
	}

	outgen.followerState = followerState[RI]{
		outgenFollowerPhaseUnknown,
		nil,
		nil,
		nil,
		nil,
		outcomeAndDigests{},
		restoredCert,
		nil,
		nil,
	}

	outgen.sharedState = sharedState{
		0,
		0,

		0,
		0,
		nil,
		0,
		nil,
	}

	// Event Loop
	chDone := outgen.ctx.Done()
	for {
		select {
		case msg := <-outgen.chNetToOutcomeGeneration:
			outgen.messageToOutcomeGeneration(msg)
		case ev := <-outgen.chPacemakerToOutcomeGeneration:
			ev.processOutcomeGeneration(outgen)
		case <-outgen.followerState.tInitial:
			outgen.eventTInitialTimeout()
		case <-outgen.leaderState.tGrace:
			outgen.eventTGraceTimeout()
		case <-outgen.leaderState.tRound:
			outgen.eventTRoundTimeout()
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			outgen.logger.Info("OutcomeGeneration: winding down", commontypes.LogFields{
				"e": outgen.sharedState.e,
				"l": outgen.sharedState.l,
			})
			outgen.metrics.Close()
			outgen.logger.Info("OutcomeGeneration: exiting", commontypes.LogFields{
				"e": outgen.sharedState.e,
				"l": outgen.sharedState.l,
			})
			return
		default:
		}
	}
}

func (outgen *outcomeGenerationState[RI]) messageToOutcomeGeneration(msg MessageToOutcomeGenerationWithSender[RI]) {
	msgEpoch := msg.msg.epoch()
	if msgEpoch < outgen.sharedState.e {
		// drop
		outgen.logger.Debug("dropping message for past epoch", commontypes.LogFields{
			"epoch":    outgen.sharedState.e,
			"msgEpoch": msgEpoch,
			"sender":   msg.sender,
		})
	} else if msgEpoch == outgen.sharedState.e {
		msg.msg.processOutcomeGeneration(outgen, msg.sender)
	} else {
		outgen.bufferedMessages[msg.sender].Push(msg.msg)
		outgen.logger.Trace("buffering message for future epoch", commontypes.LogFields{
			"msgEpoch": msgEpoch,
			"sender":   msg.sender,
		})
	}
}

func (outgen *outcomeGenerationState[RI]) unbufferMessages() {
	outgen.logger.Trace("getting messages for new epoch", nil)
	for i, buffer := range outgen.bufferedMessages {
		sender := commontypes.OracleID(i)
		for buffer.Length() > 0 {
			msg := buffer.Peek()
			msgEpoch := msg.epoch()
			if msgEpoch < outgen.sharedState.e {
				buffer.Pop()
				outgen.logger.Debug("unbuffered and dropped message", commontypes.LogFields{
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
			} else if msgEpoch == outgen.sharedState.e {
				buffer.Pop()
				outgen.logger.Trace("unbuffered message for new epoch", commontypes.LogFields{
					"msgEpoch": msgEpoch,
					"sender":   sender,
				})
				msg.processOutcomeGeneration(outgen, sender)
			} else { // msgEpoch > e
				// this and all subsequent messages are for future epochs
				// leave them in the buffer
				break
			}
		}
	}
	outgen.logger.Trace("done unbuffering messages for new epoch", nil)
}

func (outgen *outcomeGenerationState[RI]) eventNewEpochStart(ev EventNewEpochStart[RI]) {
	// Initialization
	outgen.logger.Info("starting new epoch", commontypes.LogFields{
		"epoch": ev.Epoch,
	})

	outgen.sharedState.e = ev.Epoch
	outgen.sharedState.l = Leader(outgen.sharedState.e, outgen.config.N(), outgen.config.LeaderSelectionKey())

	outgen.logger = outgen.logger.MakeUpdated(commontypes.LogFields{
		"e": outgen.sharedState.e,
		"l": outgen.sharedState.l,
	})

	outgen.sharedState.firstSeqNrOfEpoch = 0
	outgen.sharedState.seqNr = 0

	outgen.followerState.phase = outgenFollowerPhaseNewEpoch
	outgen.followerState.tInitial = time.After(outgen.config.DeltaInitial)
	outgen.followerState.outcome = outcomeAndDigests{}

	outgen.followerState.roundStartPool = pool.NewPool[MessageRoundStart[RI]](poolSize)
	outgen.followerState.proposalPool = pool.NewPool[MessageProposal[RI]](poolSize)
	outgen.followerState.preparePool = pool.NewPool[PrepareSignature](poolSize)
	outgen.followerState.commitPool = pool.NewPool[CommitSignature](poolSize)

	outgen.leaderState.phase = outgenLeaderPhaseNewEpoch
	outgen.leaderState.epochStartRequests = map[commontypes.OracleID]*epochStartRequest[RI]{}
	outgen.leaderState.readyToStartRound = false
	outgen.leaderState.tGrace = nil

	var highestCertified CertifiedPrepareOrCommit
	var highestCertifiedTimestamp HighestCertifiedTimestamp
	highestCertified = outgen.followerState.cert
	highestCertifiedTimestamp = outgen.followerState.cert.Timestamp()

	signedHighestCertifiedTimestamp, err := MakeSignedHighestCertifiedTimestamp(
		outgen.ID(),
		highestCertifiedTimestamp,
		outgen.offchainKeyring.OffchainSign,
	)
	if err != nil {
		outgen.logger.Error("error signing timestamp", commontypes.LogFields{
			"error": err,
		})
		return
	}

	outgen.logger.Info("sending MessageEpochStartRequest to leader", commontypes.LogFields{
		"highestCertifiedTimestamp": highestCertifiedTimestamp,
	})
	outgen.netSender.SendTo(MessageEpochStartRequest[RI]{
		outgen.sharedState.e,
		highestCertified,
		signedHighestCertifiedTimestamp,
	}, outgen.sharedState.l)

	if outgen.id == outgen.sharedState.l {
		outgen.leaderState.tRound = time.After(outgen.config.DeltaRound)
	}

	outgen.unbufferMessages()
}

func (outgen *outcomeGenerationState[RI]) ID() OutcomeGenerationID {
	return OutcomeGenerationID{outgen.config.ConfigDigest, outgen.sharedState.e}
}

func (outgen *outcomeGenerationState[RI]) OutcomeCtx(seqNr uint64) ocr3types.OutcomeContext {
	if seqNr != outgen.sharedState.committedSeqNr+1 {
		outgen.logger.Critical("assumption violation, seqNr isn't successor to committedSeqNr", commontypes.LogFields{
			"seqNr":          seqNr,
			"committedSeqNr": outgen.sharedState.committedSeqNr,
		})
		panic("")
	}
	return ocr3types.OutcomeContext{
		seqNr,
		outgen.sharedState.committedOutcome,
		uint64(outgen.sharedState.e),
		seqNr - outgen.sharedState.firstSeqNrOfEpoch + 1,
	}
}

func (outgen *outcomeGenerationState[RI]) ObservationQuorum(query types.Query) (quorum int, ok bool) {
	if outgen.sharedState.observationQuorum != nil {
		return *outgen.sharedState.observationQuorum, true
	}

	observationQuorum, ok := callPluginFromOutcomeGeneration[ocr3types.Quorum](
		outgen,
		"ObservationQuorum",
		0, // pure function
		outgen.OutcomeCtx(outgen.sharedState.seqNr),
		func(ctx context.Context, outctx ocr3types.OutcomeContext) (ocr3types.Quorum, error) {
			return outgen.reportingPlugin.ObservationQuorum(outctx, query)
		},
	)

	if !ok {
		return 0, false
	}

	nMinusF := outgen.config.N() - outgen.config.F

	switch observationQuorum {
	case ocr3types.QuorumFPlusOne, ocr3types.OldQuorumFPlusOne:
		quorum = outgen.config.F + 1
	case ocr3types.QuorumTwoFPlusOne, ocr3types.OldQuorumTwoFPlusOne:
		quorum = 2*outgen.config.F + 1
	case ocr3types.QuorumByzQuorum, ocr3types.OldQuorumByzQuorum:
		quorum = outgen.config.ByzQuorumSize()
	case ocr3types.QuorumNMinusF, ocr3types.OldQuorumNMinusF:
		quorum = nMinusF
	default:
		quorum = int(observationQuorum)
	}

	if !(0 < quorum && int(quorum) <= nMinusF) {
		outgen.logger.Error("invalid observation quorum", commontypes.LogFields{
			"quorum":  quorum,
			"n":       outgen.config.N(),
			"f":       outgen.config.F,
			"nMinusF": nMinusF,
		})
		return 0, false
	}

	outgen.sharedState.observationQuorum = &quorum

	return quorum, true
}

func callPluginFromOutcomeGeneration[T any, RI any](
	outgen *outcomeGenerationState[RI],
	name string,
	maxDuration time.Duration,
	outctx ocr3types.OutcomeContext,
	f func(context.Context, ocr3types.OutcomeContext) (T, error),
) (T, bool) {
	return callPlugin[T](
		outgen.ctx,
		outgen.logger,
		commontypes.LogFields{
			"seqNr": outctx.SeqNr,
			"round": outctx.Round, // nolint: staticcheck
		},
		name,
		maxDuration,
		func(ctx context.Context) (T, error) {
			return f(ctx, outctx)
		},
	)
}
