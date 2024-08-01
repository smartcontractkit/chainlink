package protocol

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/persist"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
	"golang.org/x/crypto/sha3"
)

// The capacity of the channel we use for persisting states. If more than this
// many state updates are pending, we will begin to drop updates.
const chPersistCapacity = 256

// Pacemaker keeps track of the state and message handling for an oracle
// participating in the off-chain reporting protocol
func RunPacemaker(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chNetToPacemaker <-chan MessageToPacemakerWithSender,
	chNetToReportGeneration <-chan MessageToReportGenerationWithSender,
	chPacemakerToOracle chan<- uint32,
	chReportGenerationToTransmission chan<- EventToTransmission,
	config config.SharedConfig,
	configOverrider types.ConfigOverrider,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	datasource types.DataSource,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	netSender NetworkSender,
	privateKeys types.PrivateKeys,
	telemetrySender TelemetrySender,
) {
	pace := makePacemakerState(
		ctx, subprocesses, chNetToPacemaker, chNetToReportGeneration, chPacemakerToOracle,
		chReportGenerationToTransmission, config, configOverrider, contractTransmitter, database,
		datasource, id, localConfig, logger, netSender, privateKeys,
		telemetrySender,
	)
	pace.run()
}

func makePacemakerState(ctx context.Context,
	subprocesses *subprocesses.Subprocesses,
	chNetToPacemaker <-chan MessageToPacemakerWithSender,
	chNetToReportGeneration <-chan MessageToReportGenerationWithSender,
	chPacemakerToOracle chan<- uint32,
	chReportGenerationToTransmission chan<- EventToTransmission,
	config config.SharedConfig, configOverrider types.ConfigOverrider,
	contractTransmitter types.ContractTransmitter,
	database types.Database, datasource types.DataSource, id commontypes.OracleID,
	localConfig types.LocalConfig, logger loghelper.LoggerWithContext,
	netSender NetworkSender, privateKeys types.PrivateKeys,
	telemetrySender TelemetrySender,
) pacemakerState {
	return pacemakerState{
		ctx:          ctx,
		subprocesses: subprocesses,

		chNetToPacemaker:                 chNetToPacemaker,
		chNetToReportGeneration:          chNetToReportGeneration,
		chPacemakerToOracle:              chPacemakerToOracle,
		chReportGenerationToTransmission: chReportGenerationToTransmission,
		config:                           config,
		configOverrider:                  configOverrider,
		contractTransmitter:              contractTransmitter,
		database:                         database,
		datasource:                       datasource,
		id:                               id,
		localConfig:                      localConfig,
		logger:                           logger,
		netSender:                        netSender,
		privateKeys:                      privateKeys,
		telemetrySender:                  telemetrySender,

		newepoch: make([]uint32, config.N()),
	}
}

type pacemakerState struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chNetToPacemaker                 <-chan MessageToPacemakerWithSender
	chNetToReportGeneration          <-chan MessageToReportGenerationWithSender
	chReportGenerationToPacemaker    <-chan EventToPacemaker
	chPacemakerToOracle              chan<- uint32
	chReportGenerationToTransmission chan<- EventToTransmission
	config                           config.SharedConfig
	configOverrider                  types.ConfigOverrider
	contractTransmitter              types.ContractTransmitter
	database                         types.Database
	datasource                       types.DataSource
	id                               commontypes.OracleID
	localConfig                      types.LocalConfig
	logger                           loghelper.LoggerWithContext
	netSender                        NetworkSender
	privateKeys                      types.PrivateKeys
	telemetrySender                  TelemetrySender
	// Test use only: send testBlocker an event to halt the pacemaker event loop,
	// send testUnblocker an event to resume it.
	testBlocker   chan eventTestBlock
	testUnblocker chan eventTestUnblock

	reportGenerationSubprocess subprocesses.Subprocesses
	cancelReportGeneration     context.CancelFunc

	chPersist chan<- types.PersistentState

	// ne is the highest epoch number this oracle has broadcast in a newepoch
	// message, during the current epoch
	ne uint32

	// e is the number of the current epoch
	e uint32

	// l is the index of the leader for the current epoch
	l commontypes.OracleID

	// newepoch[j] is the highest epoch number oracle j has sent in a newepoch
	// message, during the current epoch.
	newepoch []uint32

	// tResend is a timeout used by the leader-election protocol to
	// periodically resend the latest Newepoch message in order to
	// guard against unreliable network conditions
	tResend <-chan time.Time

	// tProgress is a timeout used by the leader-election protocol to track
	// whether the current leader is making adequate progress.
	tProgress <-chan time.Time

	notifyOracleOfNewEpoch bool
}

func (pace *pacemakerState) run() {
	pace.logger.Info("Running Pacemaker", nil)

	// Initialization

	{
		chPersist := make(chan types.PersistentState, chPersistCapacity)
		pace.chPersist = chPersist
		pace.subprocesses.Go(func() {
			persist.Persist(
				pace.ctx,
				chPersist,
				pace.config.ConfigDigest,
				pace.database,
				pace.localConfig.DatabaseTimeout,
				pace.logger,
			)
		})
	}

	// rounds start with 1, so let's make epochs also start with 1
	// this also gives us cleaner behavior for the initial epoch, which is otherwise
	// immediately terminated and superseded due to restoreNeFromTransmitter below
	pace.e = 1
	pace.l = Leader(pace.e, pace.config.N(), pace.config.LeaderSelectionKey())

	// Attempt to restore state from database. This is implicit in the
	// design document.
	pace.restoreStateFromDatabase()

	pace.restoreNeFromTransmitter()

	pace.spawnReportGeneration()

	pace.tProgress = time.After(pace.config.DeltaProgress)

	pace.sendNewepoch(pace.ne)

	pace.notifyOracleOfNewEpoch = true

	// Initialization complete

	// Take a reference to the ctx.Done channel once, here, to avoid taking the
	// context lock below.
	chDone := pace.ctx.Done()

	// Event Loop
	for {
		var nilOrChPacemakerToOracle chan<- uint32
		if pace.notifyOracleOfNewEpoch {
			nilOrChPacemakerToOracle = pace.chPacemakerToOracle
		} else {
			nilOrChPacemakerToOracle = nil
		}

		select {
		case nilOrChPacemakerToOracle <- pace.e:
			pace.notifyOracleOfNewEpoch = false
		case msg := <-pace.chNetToPacemaker:
			msg.msg.processPacemaker(pace, msg.sender)
		case ev := <-pace.chReportGenerationToPacemaker:
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
			pace.reportGenerationSubprocess.Wait()
			pace.logger.Info("Pacemaker: exiting", nil)
			return
		default:
		}
	}
}

func (pace *pacemakerState) restoreStateFromDatabase() {
	var state *types.PersistentState
	var err error
	ok := pace.subprocesses.BlockForAtMost(
		pace.ctx,
		pace.localConfig.DatabaseTimeout,
		func(ctx context.Context) {
			state, err = pace.database.ReadState(ctx, pace.config.ConfigDigest)
		},
	)

	if !ok {
		pace.logger.Error("Pacemaker: Timeout while restoring state from database", commontypes.LogFields{
			"timeout": pace.localConfig.DatabaseTimeout,
		})
		return
	}

	if err != nil {
		pace.logger.ErrorIfNotCanceled("Pacemaker: error while restoring state from database", pace.ctx, commontypes.LogFields{
			"error": err,
		})
		return
	}

	if state == nil {
		pace.logger.Info("Pacemaker: Database contains no state to restore", nil)
		return
	}

	if err := pace.sanityCheckState(state); err != nil {
		pace.logger.Error("Pacemaker: Ignoring state from database because it is corrupted", commontypes.LogFields{
			"error": err,
		})
		return
	}

	if state.Epoch < pace.e {
		pace.logger.Info("Skipped restore state from database because it was stale", commontypes.LogFields{
			"databaseEpoch": state.Epoch,
			"epoch":         pace.e,
		})
		return
	}

	pace.e = state.Epoch
	pace.ne = state.HighestSentEpoch
	copy(pace.newepoch, state.HighestReceivedEpoch)
	pace.l = Leader(pace.e, pace.config.N(), pace.config.LeaderSelectionKey())
	pace.logger.Info("Restored state from database", commontypes.LogFields{
		"epoch":  pace.e,
		"leader": pace.l,
	})
}

func (pace *pacemakerState) restoreNeFromTransmitter() {
	var configDigest types.ConfigDigest
	var epoch uint32
	var err error
	ok := pace.subprocesses.BlockForAtMost(
		pace.ctx,
		pace.localConfig.BlockchainTimeout,
		func(ctx context.Context) {
			configDigest, epoch, _, _, _, err = pace.contractTransmitter.LatestTransmissionDetails(ctx)
		},
	)

	if !ok {
		pace.logger.Error("Pacemaker: latestTransmissionDetails timed out while restoring ne", commontypes.LogFields{
			"timeout": pace.localConfig.BlockchainTimeout,
		})
		return
	}

	if err != nil {
		pace.logger.Error("Pacemaker: latestTransmissionDetails returned error while restoring ne", commontypes.LogFields{
			"error": err,
		})
		return
	}

	if pace.config.ConfigDigest != configDigest {
		pace.logger.Info("Pacemaker: ConfigDigest differs from contract. Cannot restore ne", commontypes.LogFields{
			"pacemakerConfigDigest": pace.config.ConfigDigest,
			"contractConfigDigest":  configDigest,
		})
		return
	}

	// epoch + 1 can overflow and the condition will be false -- that's fine
	// since we cannot proceed beyond epoch anyways at that point
	if pace.ne < epoch+1 {
		pace.logger.Info("Pacemaker: Restored ne from contract", commontypes.LogFields{
			"previousNe": pace.ne,
			"ne":         epoch + 1,
		})
		pace.ne = epoch + 1
	}
}

func (pace *pacemakerState) sanityCheckState(state *types.PersistentState) error {
	if state.HighestSentEpoch < state.Epoch {
		return fmt.Errorf("HighestSentEpoch < Epoch: %v < %v", state.HighestSentEpoch, state.Epoch)
	}

	if len(state.HighestReceivedEpoch) != pace.config.N() {
		return fmt.Errorf("len(HighestReceivedEpoch) != N: %v != %v", len(state.HighestReceivedEpoch), pace.config.N())
	}

	return nil
}

func (pace *pacemakerState) persist() {
	// send copy to be safe against mutations
	state := types.PersistentState{
		pace.e,
		pace.ne,
		append([]uint32{}, pace.newepoch...),
	}

	select {
	case pace.chPersist <- state:
	default:
		pace.logger.Warn("Pacemaker: chPersist is backed up, discarding state", commontypes.LogFields{
			"state":    state,
			"capacity": chPersistCapacity,
		})
	}
}

// eventProgress is called when a "progress" event is emitted by the reporting
// prototol. It resets the timer which will trigger the oracle to broadcast a
// "newepoch" message, if it runs out.
func (pace *pacemakerState) eventProgress() {
	pace.tProgress = time.After(pace.config.DeltaProgress)
}

func (pace *pacemakerState) sendNewepoch(newEpoch uint32) {
	pace.netSender.Broadcast(MessageNewEpoch{newEpoch})
	if pace.ne != newEpoch {
		pace.ne = newEpoch
		pace.persist()
	}
	pace.tResend = time.After(pace.config.DeltaResend)
}

func (pace *pacemakerState) eventTResendTimeout() {
	pace.sendNewepoch(pace.ne)
}

func (pace *pacemakerState) eventTProgressTimeout() {
	pace.eventChangeLeader()
}

func (pace *pacemakerState) eventChangeLeader() {
	pace.tProgress = nil
	sendEpoch := pace.ne
	epochPlusOne := pace.e + 1
	if epochPlusOne <= pace.e {
		pace.logger.Error("Pacemaker: epoch overflows, cannot change leader", nil)
		return
	}

	if sendEpoch < epochPlusOne {
		sendEpoch = epochPlusOne
	}
	pace.sendNewepoch(sendEpoch)
}

func (pace *pacemakerState) messageNewepoch(msg MessageNewEpoch, sender commontypes.OracleID) {

	if int(sender) < 0 || int(sender) >= len(pace.newepoch) {
		pace.logger.Error("Pacemaker: dropping NewEpoch message from invalid sender", commontypes.LogFields{
			"sender": sender,
			"N":      len(pace.newepoch),
		})
		return
	}

	if pace.newepoch[sender] < msg.Epoch {
		pace.newepoch[sender] = msg.Epoch
		pace.persist()
	} else {
		// neither of the following two "upon" handlers can be triggered
		return
	}

	// upon |{p_j ∈ P | newepoch[j] > ne}| > f do
	{
		candidateEpochs := sortedGreaterThan(pace.newepoch, pace.ne)
		if len(candidateEpochs) > pace.config.F {
			// ē ← max {e' | {p_j ∈ P | newepoch[j] ≥ e' } > f}
			newEpoch := candidateEpochs[len(candidateEpochs)-(pace.config.F+1)]
			pace.sendNewepoch(newEpoch)
		}
	}

	// upon |{p_j ∈ P | newepoch[j] > e}| > 2f do
	{
		candidateEpochs := sortedGreaterThan(pace.newepoch, pace.e)
		if len(candidateEpochs) > 2*pace.config.F {
			// ē ← max {e' | {p_j ∈ P | newepoch[j] ≥ e' } > 2f}
			//
			// since candidateEpochs contains, in increasing order, the epochs from
			// the received newepoch messages, this value of newEpoch was sent by at
			// least 2F+1 processes
			newEpoch := candidateEpochs[len(candidateEpochs)-(2*pace.config.F+1)]
			pace.logger.Debug("Moving to epoch, based on candidateEpochs", commontypes.LogFields{
				"newEpoch":        newEpoch,
				"candidateEpochs": candidateEpochs,
			})
			l := Leader(newEpoch, pace.config.N(), pace.config.LeaderSelectionKey())
			pace.e, pace.l = newEpoch, l // (e, l) ← (ē, leader(ē))
			if pace.ne < pace.e {        // ne ← max{ne, e}
				pace.ne = pace.e
			}
			pace.persist()

			// abort instance [...], initialize instance (e,l) of report generation
			pace.spawnReportGeneration()

			pace.notifyOracleOfNewEpoch = true

			pace.tProgress = time.After(pace.config.DeltaProgress) // restart timer T_{progress}
		}
	}
}

func (pace *pacemakerState) spawnReportGeneration() {
	if pace.cancelReportGeneration != nil {
		pace.cancelReportGeneration()
		pace.reportGenerationSubprocess.Wait()
	}

	chReportGenerationToPacemaker := make(chan EventToPacemaker)
	pace.chReportGenerationToPacemaker = chReportGenerationToPacemaker

	ctxReportGeneration, cancelReportGeneration := context.WithCancel(pace.ctx)

	{
		// Take a copy of the relevant pacemaker fields, to avoid a race condition between the
		// following go func and the agreement section of messageNewepoch, which
		// assigns new values to some pace attributes. This race condition will never
		// happen, given a reasonable value for DeltaProgress, but
		// TestPacemakerNodesEventuallyReachEpochConsensus has an unreasonable value.
		subprocesses,
			chNetToReportGeneration,
			chReportGenerationToTransmission,
			config,
			configOverrider,
			contractTransmitter,
			datasource,
			e,
			id,
			l,
			localConfig,
			logger,
			netSender,
			privateKeys,
			telemetrySender := pace.subprocesses,
			pace.chNetToReportGeneration,
			pace.chReportGenerationToTransmission,
			pace.config,
			pace.configOverrider,
			pace.contractTransmitter,
			pace.datasource,
			pace.e,
			pace.id,
			pace.l,
			pace.localConfig,
			pace.logger,
			pace.netSender,
			pace.privateKeys,
			pace.telemetrySender

		pace.reportGenerationSubprocess.Go(func() {
			defer cancelReportGeneration()
			RunReportGeneration(
				ctxReportGeneration,
				subprocesses,

				chNetToReportGeneration,
				chReportGenerationToPacemaker,
				chReportGenerationToTransmission,
				config,
				configOverrider,
				contractTransmitter,
				datasource,
				e,
				id,
				l,
				localConfig,
				logger,
				netSender,
				privateKeys,
				telemetrySender,
			)
		})
	}
	pace.cancelReportGeneration = cancelReportGeneration
}

// sortedGreaterThan returns the *sorted* elements of xs which are greater than y
func sortedGreaterThan(xs []uint32, y uint32) (rv []uint32) {
	for _, x := range xs {
		if x > y {
			rv = append(rv, x)
		}
	}
	sort.Slice(rv, func(i, j int) bool { return rv[i] < rv[j] })
	return rv
}

// Leader will produce an oracle id for the given epoch.
func Leader(epoch uint32, n int, key [16]byte) (leader commontypes.OracleID) {
	// No need for HMAC. Since we use Keccak256, prepending
	// with key gives us a PRF already.
	h := sha3.NewLegacyKeccak256()
	h.Write(key[:])
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))
	h.Write(b)

	result := big.NewInt(0)
	r := big.NewInt(0).SetBytes(h.Sum(nil))
	// This is biased, but we don't care because the prob of us hitting the bias are
	// less than 2**5/2**256 = 2**-251.
	result.Mod(r, big.NewInt(int64(n)))
	return commontypes.OracleID(result.Int64())
}

type eventTestBlock struct{}
type eventTestUnblock struct{}
