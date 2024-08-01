package protocol

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/permutation"
)

func RunPacemaker[RI any](
	ctx context.Context,

	chNetToPacemaker <-chan MessageToPacemakerWithSender[RI],
	chPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI],
	chOutcomeGenerationToPacemaker <-chan EventToPacemaker[RI],
	config ocr3config.SharedConfig,
	database Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,

	restoredState PacemakerState,
) {
	pace := makePacemakerState[RI](
		ctx, chNetToPacemaker,
		chPacemakerToOutcomeGeneration, chOutcomeGenerationToPacemaker,
		config, database,
		id, localConfig, logger, metricsRegisterer, netSender, offchainKeyring,
		telemetrySender,
	)
	pace.run(restoredState)
}

func makePacemakerState[RI any](
	ctx context.Context,
	chNetToPacemaker <-chan MessageToPacemakerWithSender[RI],
	chPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI],
	chOutcomeGenerationToPacemaker <-chan EventToPacemaker[RI],
	config ocr3config.SharedConfig,
	database Database, id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	netSender NetworkSender[RI],
	offchainKeyring types.OffchainKeyring,
	telemetrySender TelemetrySender,
) pacemakerState[RI] {
	return pacemakerState[RI]{
		ctx: ctx,

		chNetToPacemaker:               chNetToPacemaker,
		chPacemakerToOutcomeGeneration: chPacemakerToOutcomeGeneration,
		chOutcomeGenerationToPacemaker: chOutcomeGenerationToPacemaker,
		config:                         config,
		database:                       database,
		id:                             id,
		localConfig:                    localConfig,
		logger:                         logger.MakeUpdated(commontypes.LogFields{"proto": "pacemaker"}),
		metrics:                        newPacemakerMetrics(metricsRegisterer, logger),
		netSender:                      netSender,
		offchainKeyring:                offchainKeyring,
		telemetrySender:                telemetrySender,

		newEpochWishes: make([]uint64, config.N()),
	}
}

type pacemakerState[RI any] struct {
	ctx context.Context

	chNetToPacemaker               <-chan MessageToPacemakerWithSender[RI]
	chPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI]
	chOutcomeGenerationToPacemaker <-chan EventToPacemaker[RI]
	config                         ocr3config.SharedConfig
	database                       Database
	id                             commontypes.OracleID
	localConfig                    types.LocalConfig
	logger                         loghelper.LoggerWithContext
	metrics                        *pacemakerMetrics
	netSender                      NetworkSender[RI]
	offchainKeyring                types.OffchainKeyring
	telemetrySender                TelemetrySender
	// Test use only: send testBlocker an event to halt the pacemaker event loop,
	// send testUnblocker an event to resume it.
	testBlocker   chan eventTestBlock
	testUnblocker chan eventTestUnblock

	// ne is the highest epoch number this oracle has broadcast in a
	// NewEpochWish message
	ne uint64

	// e is the number of the current epoch
	e uint64

	// l is the index of the leader for the current epoch
	l commontypes.OracleID

	// newEpochWishes[j] is the highest epoch number oracle j has sent in a
	// NewEpochWish message
	newEpochWishes []uint64

	// tResend is a timeout used to periodically resend the latest NewEpochWish
	// message in order to guard against unreliable network conditions
	tResend <-chan time.Time

	// tProgress is a timeout used by the protocol to track whether the current
	// leader/epoch is making adequate progress.
	tProgress <-chan time.Time

	notifyOutcomeGenerationOfNewEpoch bool
}

func (pace *pacemakerState[RI]) run(restoredState PacemakerState) {
	pace.logger.Info("Pacemaker: running", nil)

	// Initialization

	if restoredState == (PacemakerState{}) {
		// seqNrs start with 1, so let's make epochs also start with 1
		pace.ne = 1
		pace.e = 1
	} else {
		pace.ne = restoredState.HighestSentNewEpochWish
		pace.e = restoredState.Epoch
	}
	pace.l = Leader(pace.e, pace.config.N(), pace.config.LeaderSelectionKey())

	pace.tProgress = time.After(pace.config.DeltaProgress)

	pace.sendNewEpochWish()

	pace.notifyOutcomeGenerationOfNewEpoch = true

	// Initialization complete

	// Take a reference to the ctx.Done channel once, here, to avoid taking the
	// context lock below.
	chDone := pace.ctx.Done()

	// Event Loop
	for {
		var nilOrChPacemakerToOutcomeGeneration chan<- EventToOutcomeGeneration[RI]
		if pace.notifyOutcomeGenerationOfNewEpoch {
			nilOrChPacemakerToOutcomeGeneration = pace.chPacemakerToOutcomeGeneration
		} else {
			nilOrChPacemakerToOutcomeGeneration = nil
		}

		select {
		case nilOrChPacemakerToOutcomeGeneration <- EventNewEpochStart[RI]{pace.e}:
			pace.notifyOutcomeGenerationOfNewEpoch = false
		case msg := <-pace.chNetToPacemaker:
			msg.msg.processPacemaker(pace, msg.sender)
		case ev := <-pace.chOutcomeGenerationToPacemaker:
			ev.processPacemaker(pace)
		case <-pace.tResend:
			pace.eventTResendTimeout()
		case <-pace.tProgress:
			pace.eventTProgressTimeout()
		case <-pace.testBlocker:
			<-pace.testUnblocker
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			pace.logger.Info("Pacemaker: winding down", nil)
			pace.metrics.Close()
			pace.logger.Info("Pacemaker: exiting", nil)
			return
		default:
		}
	}
}

func (pace *pacemakerState[RI]) eventProgress() {
	pace.tProgress = time.After(pace.config.DeltaProgress)
}

func (pace *pacemakerState[RI]) sendNewEpochWish() {
	pace.netSender.Broadcast(MessageNewEpochWish[RI]{pace.ne})
	pace.tResend = time.After(pace.config.DeltaResend)
}

func (pace *pacemakerState[RI]) eventTResendTimeout() {
	pace.sendNewEpochWish()
}

func (pace *pacemakerState[RI]) eventTProgressTimeout() {
	pace.logger.Debug("TProgress fired", commontypes.LogFields{
		"deltaProgress": pace.config.DeltaProgress.String(),
	})
	pace.eventNewEpochRequest()
}

func (pace *pacemakerState[RI]) eventNewEpochRequest() {
	pace.tProgress = nil
	epochPlusOne := pace.e + 1
	if epochPlusOne <= pace.e {
		pace.logger.Critical("epoch overflows, cannot change leader", nil)
		return
	}

	if pace.ne < epochPlusOne { // ne ← max{e + 1, ne}
		if err := pace.persist(PacemakerState{pace.e, epochPlusOne}); err != nil {
			pace.logger.Error("could not persist pacemaker state in eventNewEpochRequest", commontypes.LogFields{
				"error": err,
			})
		}

		pace.ne = epochPlusOne
	}
	pace.sendNewEpochWish()
}

func (pace *pacemakerState[RI]) messageNewEpochWish(msg MessageNewEpochWish[RI], sender commontypes.OracleID) {
	if pace.newEpochWishes[sender] < msg.Epoch {
		pace.newEpochWishes[sender] = msg.Epoch
	}

	var wishForEpoch uint64

	// upon |{p_j ∈ P | newEpochWishes[j] > ne}| > f do
	{
		candidateEpochs := sortedGreaterThan(pace.newEpochWishes, pace.ne)
		if len(candidateEpochs) > pace.config.F {
			// ē ← max {e' | {p_j ∈ P | newEpochWishes[j] ≥ e' } > f}
			wishForEpoch = candidateEpochs[len(candidateEpochs)-(pace.config.F+1)]
			// ne ← max(ne, ē) is superfluous because ē is always greater or
			// equal ne: this rule is only triggered if there are at least f+1
			// wishes greater than ne. ē is the greatest wish such that f+1
			// wishes are greater or equal ē.

			// see "if wishForEpoch != 0 {" for continuation below
		}
	}

	var switchToEpoch uint64

	// upon |{p_j ∈ P | newEpochWishes[j] > e}| > 2f do
	{
		candidateEpochs := sortedGreaterThan(pace.newEpochWishes, pace.e)
		if len(candidateEpochs) > 2*pace.config.F {
			// ē ← max {e' | {p_j ∈ P | newEpochWishes[j] ≥ e' } > 2f}
			//
			// since candidateEpochs contains, in increasing order, the epochs
			// from the received NewEpochWish messages, this value was sent by
			// at least 2F+1 processes
			switchToEpoch = candidateEpochs[len(candidateEpochs)-(2*pace.config.F+1)]
			// see "if switchToEpoch != 0 {" for continuation below
		}
	}

	// persist wishForEpoch and switchToEpoch
	if wishForEpoch != 0 || switchToEpoch != 0 {
		persistState := PacemakerState{}
		if wishForEpoch == 0 {
			persistState.HighestSentNewEpochWish = pace.ne
		} else {
			persistState.HighestSentNewEpochWish = wishForEpoch
		}
		if switchToEpoch == 0 {
			persistState.Epoch = pace.e
		} else {
			persistState.Epoch = switchToEpoch
		}

		// needed so that persisted state is consistent with "ne ← max{ne, e}"
		// statement in agreement rule
		if persistState.HighestSentNewEpochWish < persistState.Epoch {
			persistState.HighestSentNewEpochWish = persistState.Epoch
		}

		if err := pace.persist(persistState); err != nil {
			pace.logger.Error("could not persist pacemaker state in messageNewEpochWish", commontypes.LogFields{
				"error": err,
			})
		}
	}

	if wishForEpoch != 0 {
		pace.ne = wishForEpoch
		pace.sendNewEpochWish()
	}

	if switchToEpoch != 0 {
		pace.logger.Debug("moving to new epoch", commontypes.LogFields{
			"newEpoch": switchToEpoch,
		})
		l := Leader(switchToEpoch, pace.config.N(), pace.config.LeaderSelectionKey())
		pace.e, pace.l = switchToEpoch, l // (e, l) ← (ē, leader(ē))
		if pace.ne < pace.e {             // ne ← max{ne, e}
			pace.ne = pace.e
		}
		pace.metrics.epoch.Set(float64(pace.e))
		pace.metrics.leader.Set(float64(pace.l))
		pace.tProgress = time.After(pace.config.DeltaProgress) // restart timer T_{progress}

		pace.notifyOutcomeGenerationOfNewEpoch = true // invoke event newEpochStart(e, l)
	}
}

func (pace *pacemakerState[RI]) persist(state PacemakerState) error {
	writeCtx, writeCancel := context.WithTimeout(pace.ctx, pace.localConfig.DatabaseTimeout)
	defer writeCancel()
	err := pace.database.WritePacemakerState(
		writeCtx,
		pace.config.ConfigDigest,
		state,
	)
	if err != nil {
		return fmt.Errorf("error while persisting pacemaker state: %w", err)
	}
	return nil
}

// sortedGreaterThan returns the *sorted* elements of xs which are greater than y
func sortedGreaterThan(xs []uint64, y uint64) (rv []uint64) {
	for _, x := range xs {
		if x > y {
			rv = append(rv, x)
		}
	}
	sort.Slice(rv, func(i, j int) bool { return rv[i] < rv[j] })
	return rv
}

// Leader will produce an oracle id for the given epoch.
func Leader(epoch uint64, n int, key [16]byte) (leader commontypes.OracleID) {
	span := epoch / uint64(n)
	epochInSpan := epoch % uint64(n)

	mac := hmac.New(sha256.New, key[:])
	_ = binary.Write(mac, binary.BigEndian, span)

	var permutationKey [16]byte
	copy(permutationKey[:], mac.Sum(nil))
	pi := permutation.Permutation(n, permutationKey)
	return commontypes.OracleID(pi[epochInSpan])
}

type eventTestBlock struct{}
type eventTestUnblock struct{}
