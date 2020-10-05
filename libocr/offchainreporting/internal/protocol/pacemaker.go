package protocol

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"
	"golang.org/x/crypto/sha3"
)

func RunPacemaker(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chNetToPacemaker <-chan MessageToPacemakerWithSender,
	chNetToReportGeneration <-chan MessageToReportGenerationWithSender,
	chReportGenerationToTransmission chan<- EventToTransmission,
	config config.SharedConfig,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	datasource types.DataSource,
	id types.OracleID,
	localConfig types.LocalConfig,
	logger types.Logger,
	netSender NetworkSender,
	privateKeys types.PrivateKeys,
) {
	pace := pacemakerState{
		ctx:                              ctx,
		subprocesses:                     subprocesses,
		chNetToPacemaker:                 chNetToPacemaker,
		chNetToReportGeneration:          chNetToReportGeneration,
		chReportGenerationToTransmission: chReportGenerationToTransmission,
		config:                           config,
		contractTransmitter:              contractTransmitter,
		database:                         database,
		datasource:                       datasource,
		id:                               id,
		localConfig:                      localConfig,
		logger:                           logger,
		netSender:                        netSender,
		privateKeys:                      privateKeys,

		newepoch: make([]uint32, config.N()),
	}
	pace.run()
}

type pacemakerState struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chNetToPacemaker                 <-chan MessageToPacemakerWithSender
	chNetToReportGeneration          <-chan MessageToReportGenerationWithSender
	chReportGenerationToPacemaker    <-chan EventToPacemaker
	chReportGenerationToTransmission chan<- EventToTransmission
	config                           config.SharedConfig
	contractTransmitter              types.ContractTransmitter
	database                         types.Database
	datasource                       types.DataSource
	id                               types.OracleID
	localConfig                      types.LocalConfig
	logger                           types.Logger
	netSender                        NetworkSender
	privateKeys                      types.PrivateKeys

	cancelReportGeneration context.CancelFunc

			ne uint32

		e uint32

		l types.OracleID

			newepoch []uint32

				tResend <-chan time.Time

			tProgress <-chan time.Time
}

func (pace *pacemakerState) run() {
	pace.logger.Info("Running Pacemaker", nil)

					pace.e = 1
	pace.l = leader(pace.e, pace.config.N(), pace.config.LeaderSelectionKey())

			pace.restoreStateFromDatabase()

								pace.restoreNeFromTransmitter()

	pace.spawnReportGeneration()

	pace.tProgress = time.After(pace.config.DeltaProgress)

				pace.sendNewepoch(pace.ne)

	
			chDone := pace.ctx.Done()

		for {
		select {
		case msg := <-pace.chNetToPacemaker:
			msg.msg.processPacemaker(pace, msg.sender)
		case ev := <-pace.chReportGenerationToPacemaker:
			ev.processPacemaker(pace)
		case <-pace.tResend:
			pace.eventTResendTimeout()
		case <-pace.tProgress:
			pace.eventTProgressTimeout()
		case <-chDone:
		}

				select {
		case <-chDone:
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
		pace.logger.Error("Pacemaker: Timeout while restoring state from database", types.LogFields{
			"timeout": pace.localConfig.DatabaseTimeout,
		})
		return
	}

	if err != nil {
		pace.logger.Error("Pacemaker: Unexpected error while restoring state from database", types.LogFields{
			"error": err,
		})
		return
	}

	if state == nil {
		pace.logger.Info("Pacemaker: Database contains no state to restore", nil)
		return
	}

	if err := pace.sanityCheckState(state); err != nil {
		pace.logger.Error("Pacemaker: Ignoring state from database because it is corrupted", types.LogFields{
			"error": err,
		})
		return
	}

	if state.Epoch < pace.e {
		pace.logger.Info("Skipped restore state from database because it was stale", types.LogFields{
			"databaseEpoch": state.Epoch,
			"epoch":         pace.e,
		})
		return
	}

	pace.e = state.Epoch
	pace.ne = state.HighestSentEpoch
	for i, e := range state.HighestReceivedEpoch {
		pace.newepoch[i] = e
	}
	pace.l = leader(pace.e, pace.config.N(), pace.config.LeaderSelectionKey())
	pace.logger.Info("Restored state from database", types.LogFields{
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
		pace.logger.Error("Pacemaker: latestTransmissionDetails timed out while restoring ne", types.LogFields{
			"timeout": pace.localConfig.BlockchainTimeout,
		})
		return
	}

	if err != nil {
		pace.logger.Error("Pacemaker: latestTransmissionDetails returned error while restoring ne", types.LogFields{
			"error": err,
		})
		return
	}

	if pace.config.ConfigDigest != configDigest {
		pace.logger.Info("Pacemaker: ConfigDigest differs from contract. Cannot restore ne", types.LogFields{
			"pacemakerConfigDigest": pace.config.ConfigDigest,
			"contractConfigDigest":  configDigest,
		})
		return
	}

			if pace.ne < epoch+1 {
		pace.logger.Info("Pacemaker: Restored ne from contract", types.LogFields{
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
		highestReceivedEpoch := make([]uint32, pace.config.N())
	copy(highestReceivedEpoch, pace.newepoch)

	var err error
	ok := pace.subprocesses.BlockForAtMost(pace.ctx, pace.localConfig.DatabaseTimeout,
		func(ctx context.Context) {
			err = pace.database.WriteState(
				ctx,
				pace.config.ConfigDigest,
				types.PersistentState{
					pace.e,
					pace.ne,
					highestReceivedEpoch,
				},
			)
		},
	)

	if !ok {
		pace.logger.Error("Timeout while persisting state to database: %v", types.LogFields{
			"timeout": pace.localConfig.DatabaseTimeout,
		})
		return
	}

	if err != nil {
		pace.logger.Error("Unexpected error while persisting state to database: %v", types.LogFields{
			"error": err,
		})
	}
}

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

func (pace *pacemakerState) messageNewepoch(msg MessageNewEpoch, sender types.OracleID) {
		if int(sender) < 0 || int(sender) >= len(pace.newepoch) {
		pace.logger.Error("Pacemaker: dropping NewEpoch message from invalid sender", types.LogFields{
			"sender": sender,
			"N":      len(pace.newepoch),
		})
		return
	}

	if pace.newepoch[sender] < msg.Epoch {
		pace.newepoch[sender] = msg.Epoch
		pace.persist()
	} else {
				return
	}

		{
		candidateEpochs := sortedGreaterThan(pace.newepoch, pace.ne)
		if len(candidateEpochs) > pace.config.F {
						newEpoch := candidateEpochs[len(candidateEpochs)-(pace.config.F+1)]
			pace.sendNewepoch(newEpoch)
		}
	}

		{
		candidateEpochs := sortedGreaterThan(pace.newepoch, pace.e)
		if len(candidateEpochs) > 2*pace.config.F {
																		newEpoch := candidateEpochs[len(candidateEpochs)-(2*pace.config.F+1)]
			pace.logger.Debug("Moving to epoch, based on candidateEpochs", types.LogFields{
				"newEpoch":        newEpoch,
				"candidateEpochs": candidateEpochs,
			})
			l := leader(newEpoch, pace.config.N(), pace.config.LeaderSelectionKey())
			pace.e, pace.l = newEpoch, l 			if pace.ne < pace.e {        				pace.ne = pace.e
			}
			pace.persist()

						pace.spawnReportGeneration()

			pace.tProgress = time.After(pace.config.DeltaProgress) 		}
	}
}

func (pace *pacemakerState) spawnReportGeneration() {
	if pace.cancelReportGeneration != nil {
		pace.cancelReportGeneration()
	}

	chReportGenerationToPacemaker := make(chan EventToPacemaker)
	pace.chReportGenerationToPacemaker = chReportGenerationToPacemaker

	ctxReportGeneration, cancelReportGeneration := context.WithCancel(pace.ctx)
	pace.subprocesses.Go(func() {
		defer cancelReportGeneration()
		RunReportGeneration(
			ctxReportGeneration,
			pace.subprocesses,

			pace.chNetToReportGeneration,
			chReportGenerationToPacemaker,
			pace.chReportGenerationToTransmission,
			pace.config,
			pace.contractTransmitter,
			pace.datasource,
			pace.e,
			pace.id,
			pace.l,
			pace.localConfig,
			pace.logger,
			pace.netSender,
			pace.privateKeys,
		)
	})
	pace.cancelReportGeneration = cancelReportGeneration

}

func sortedGreaterThan(xs []uint32, y uint32) (rv []uint32) {
	for _, x := range xs {
		if x > y {
			rv = append(rv, x)
		}
	}
	sort.Slice(rv, func(i, j int) bool { return rv[i] < rv[j] })
	return rv
}

func leader(epoch uint32, n int, key [16]byte) (leader types.OracleID) {
			h := sha3.NewLegacyKeccak256()
	h.Write(key[:])
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))
	h.Write(b)

	result := big.NewInt(0)
	r := big.NewInt(0).SetBytes(h.Sum(nil))
			result.Mod(r, big.NewInt(int64(n)))
	return types.OracleID(result.Int64())
}
