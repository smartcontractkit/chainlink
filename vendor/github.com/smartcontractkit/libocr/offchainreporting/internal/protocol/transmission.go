package protocol

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/permutation"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// TransmissionProtocol tracks the local oracle process's role in the transmission of a
// report to the on-chain oracle contract.
//
// Note: The transmission protocol doesn't clean up pending transmissions
// when it is terminated. This is by design, but means that old pending
// transmissions may accumulate in the database. They should be garbage
// collected once in a while.
func RunTransmission(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	config config.SharedConfig,
	configOverrider types.ConfigOverrider,
	chReportGenerationToTransmission <-chan EventToTransmission,
	database types.Database,
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	transmitter types.ContractTransmitter,
) {
	t := transmissionState{
		ctx:          ctx,
		subprocesses: subprocesses,

		chReportGenerationToTransmission: chReportGenerationToTransmission,
		config:                           config,
		configOverrider:                  configOverrider,
		database:                         database,
		id:                               id,
		localConfig:                      localConfig,
		logger:                           logger,
		transmitter:                      transmitter,
	}
	t.run()
}

type transmissionState struct {
	ctx          context.Context
	subprocesses *subprocesses.Subprocesses

	chReportGenerationToTransmission <-chan EventToTransmission
	config                           config.SharedConfig
	configOverrider                  types.ConfigOverrider
	database                         types.Database
	id                               commontypes.OracleID
	localConfig                      types.LocalConfig
	logger                           loghelper.LoggerWithContext
	transmitter                      types.ContractTransmitter

	latestEpochRound EpochRound
	latestMedian     observation.Observation
	times            MinHeapTimeToPendingTransmission
	tTransmit        <-chan time.Time
}

// run runs the event loop for the local transmission protocol
func (t *transmissionState) run() {
	t.restoreFromDatabase()

	chDone := t.ctx.Done()
	for {
		select {
		case ev := <-t.chReportGenerationToTransmission:
			ev.processTransmission(t)
		case <-t.tTransmit:
			t.eventTTransmitTimeout()
		case <-chDone:
		}

		// ensure prompt exit
		select {
		case <-chDone:
			t.logger.Info("Transmission: exiting", nil)
			return
		default:
		}
	}
}

func (t *transmissionState) restoreFromDatabase() {
	childCtx, childCancel := context.WithTimeout(t.ctx, t.localConfig.DatabaseTimeout)
	defer childCancel()
	pending, err := t.database.PendingTransmissionsWithConfigDigest(childCtx, t.config.ConfigDigest)
	if err != nil {
		t.logger.ErrorIfNotCanceled("Transmission: error fetching pending transmissions from database", childCtx, commontypes.LogFields{"error": err})
		return
	}

	now := time.Now()

	// insert non-expired transmissions into queue
	for key, trans := range pending {
		if now.Before(trans.Time) {
			t.times.Push(MinHeapTimeToPendingTransmissionItem{
				key,
				trans,
			})
		}
	}

	// find logically latest expired transmission and insert into queue
	latestExpiredTransmissionKey := types.ReportTimestamp{}
	latestExpiredTransmission := (*types.PendingTransmission)(nil)
	for key, trans := range pending {
		if trans.Time.Before(now) && (EpochRound{latestExpiredTransmissionKey.Epoch, latestExpiredTransmissionKey.Round}).Less(EpochRound{key.Epoch, key.Round}) {
			latestExpiredTransmissionKey = key
			transCopy := trans // prevent aliasing of loop var
			latestExpiredTransmission = &transCopy
		}
	}
	if latestExpiredTransmission != nil {
		t.times.Push(MinHeapTimeToPendingTransmissionItem{
			latestExpiredTransmissionKey,
			*latestExpiredTransmission,
		})
	}

	// if queue isn't empty, set tTransmit to expire at next transmission time
	if t.times.Len() != 0 {
		t.tTransmit = time.After(time.Until(t.times.Peek().Time))
	}
}

// eventTransmit is called when the local process sends a transmit event
func (t *transmissionState) eventTransmit(ev EventTransmit) {
	t.logger.Debug("Received transmit event", commontypes.LogFields{
		"event": ev,
	})

	{
		contractConfigDigest, contractEpochRound, err := t.contractState()
		if err != nil {
			t.logger.Error("contractEpoch() failed during eventTransmit", commontypes.LogFields{"error": err})
			return
		}

		if contractConfigDigest != t.config.ConfigDigest {
			t.logger.Info("eventTransmit(ev): discarding ev because contractConfigDigest != configDigest", commontypes.LogFields{
				"ev":                   ev,
				"contractConfigDigest": contractConfigDigest,
				"configDigest":         t.config.ConfigDigest,
			})
			return
		}

		if !t.shouldTransmit(ev, contractEpochRound) {
			t.logger.Info("eventTransmit(ev): discarding ev because shouldTransmit returned false", commontypes.LogFields{
				"ev":                   ev,
				"contractConfigDigest": contractConfigDigest,
				"contractEpochRound":   contractEpochRound,
			})
			return
		}
	}

	var err error
	t.latestEpochRound = EpochRound{ev.Epoch, ev.Round}
	t.latestMedian, err = ev.Report.AttributedObservations.Median()
	if err != nil {
		t.logger.Error("could not compute median", commontypes.LogFields{"error": err})
	}

	now := time.Now()
	delayMaybe := t.transmitDelay(ev.Epoch, ev.Round)
	if delayMaybe == nil {
		return
	}
	delay := *delayMaybe
	serializedReport, rs, ss, vs, err := ev.Report.TransmissionArgs(ReportContext{
		t.config.ConfigDigest,
		ev.Epoch,
		ev.Round,
	})
	if err != nil {
		t.logger.Error("Failed to serialize contract report", commontypes.LogFields{"error": err})
		return
	}

	key := types.ReportTimestamp{
		ConfigDigest: t.config.ConfigDigest,
		Epoch:        ev.Epoch,
		Round:        ev.Round,
	}
	median, err := ev.Report.AttributedObservations.Median()
	if err != nil {
		t.logger.Error("could not take median of observations",
			commontypes.LogFields{"error": err})
	}
	transmission := types.PendingTransmission{
		Time:             now.Add(delay),
		Median:           median.RawObservation(),
		SerializedReport: serializedReport,
		Rs:               rs, Ss: ss, Vs: vs,
	}

	ok := t.subprocesses.BlockForAtMost(
		t.ctx,
		t.localConfig.DatabaseTimeout,
		func(ctx context.Context) {
			if err := t.database.StorePendingTransmission(ctx, key, transmission); err != nil {
				t.logger.Error("Error while persisting pending transmission to database", commontypes.LogFields{"error": err})
			}
		},
	)
	if !ok {
		t.logger.Error("Database.StorePendingTransmission timed out", commontypes.LogFields{
			"timeout": t.localConfig.DatabaseTimeout,
		})
	}
	t.times.Push(MinHeapTimeToPendingTransmissionItem{key, transmission})

	next := t.times.Peek()
	if (EpochRound{ev.Epoch, ev.Round}) == (EpochRound{next.Epoch, next.Round}) {
		t.tTransmit = time.After(delay)
	}
}

func (t *transmissionState) eventTTransmitTimeout() {
	defer func() {
		if t.times.Len() != 0 { // If there's other transmissions due later...
			// ...reset timer to expire when the next one is due
			item := t.times.Peek()
			t.tTransmit = time.After(time.Until(item.Time))
		}
	}()

	if t.times.Len() == 0 {
		return
	}
	item := t.times.Pop()
	itemEpochRound := EpochRound{item.Epoch, item.Round}

	ok := t.subprocesses.BlockForAtMost(
		t.ctx,
		t.localConfig.DatabaseTimeout,
		func(ctx context.Context) {
			if err := t.database.DeletePendingTransmission(ctx, types.ReportTimestamp{
				ConfigDigest: t.config.ConfigDigest,
				Epoch:        item.Epoch,
				Round:        item.Round,
			}); err != nil {
				t.logger.Error("eventTTransmitTimeout: Error while deleting pending transmission from database", commontypes.LogFields{"error": err})
			}
		},
	)
	if !ok {
		t.logger.Error("Database.DeletePendingTransmission timed out", commontypes.LogFields{
			"timeout": t.localConfig.DatabaseTimeout,
		})
		// carry on
	}

	contractConfigDigest, contractEpochRound, err := t.contractState()
	if err != nil {
		t.logger.Error("eventTTransmitTimeout: contractState() failed", commontypes.LogFields{"error": err})
		return
	}

	if item.ConfigDigest != contractConfigDigest {
		t.logger.Info("eventTTransmitTimeout: configDigest doesn't match, discarding transmission", commontypes.LogFields{
			"contractConfigDigest": contractConfigDigest,
			"configDigest":         item.ConfigDigest,
			"median":               item.Median,
			"epoch":                item.Epoch,
			"round":                item.Round,
		})
		return
	}

	if !contractEpochRound.Less(itemEpochRound) {
		t.logger.Info("eventTTransmitTimeout: Skipping transmission because report is stale", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"median":             item.Median,
			"epoch":              item.Epoch,
			"round":              item.Round,
		})
		return
	}

	t.logger.Info("eventTTransmitTimeout: Transmitting with median", commontypes.LogFields{
		"median": item.Median,
		"epoch":  item.Epoch,
		"round":  item.Round,
	})

	ok = t.subprocesses.BlockForAtMost(
		t.ctx,
		t.localConfig.ContractTransmitterTransmitTimeout,
		func(ctx context.Context) {
			err = t.transmitter.Transmit(ctx, item.SerializedReport, item.Rs, item.Ss, item.Vs)
		},
	)
	if !ok {
		t.logger.Error("eventTTransmitTimeout: Transmit timed out", commontypes.LogFields{
			"timeout": t.localConfig.ContractTransmitterTransmitTimeout,
		})
		return
	}
	if err != nil {
		t.logger.Error("eventTTransmitTimeout: Error while transmitting report on-chain", commontypes.LogFields{"error": err})
		return
	}

	t.logger.Info("eventTTransmitTimeout:❗️successfully transmitted report on-chain", commontypes.LogFields{
		"median": item.Median,
		"epoch":  item.Epoch,
		"round":  item.Round,
	})
}

func (t *transmissionState) shouldTransmit(ev EventTransmit, contractEpochRound EpochRound) bool {
	reportEpochRound := EpochRound{ev.Epoch, ev.Round}
	if !contractEpochRound.Less(reportEpochRound) {
		t.logger.Debug("shouldTransmit() = false, report is stale", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"epochRound":         reportEpochRound,
		})
		return false
	}
	if t.latestEpochRound == (EpochRound{}) {
		t.logger.Debug("shouldTransmit() = true, latestEpochRound is empty", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"epochRound":         reportEpochRound,
			"latestEpochRound":   t.latestEpochRound,
		})
		return true
	}
	if reportEpochRound.Less(t.latestEpochRound) || reportEpochRound == t.latestEpochRound {
		t.logger.Debug("shouldTransmit() = false, report is older than latest report", commontypes.LogFields{
			"contractEpochRound": contractEpochRound,
			"epochRound":         reportEpochRound,
			"latestEpochRound":   t.latestEpochRound,
		})
		return false
	}

	reportMedian, err := ev.Report.AttributedObservations.Median()
	if err != nil {
		t.logger.Error("could not compute median", commontypes.LogFields{
			"error": err,
		})
		return false
	}

	alphaPPB := t.config.AlphaPPB
	if override := t.configOverrider.ConfigOverride(); override != nil {
		t.logger.Debug("shouldTransmit: using override for alphaPPB", commontypes.LogFields{
			"epochRound":         reportEpochRound,
			"alphaPPB":           alphaPPB,
			"overriddenAlphaPPB": override.AlphaPPB,
		})
		alphaPPB = override.AlphaPPB
	}

	deviates := t.latestMedian.Deviates(reportMedian, alphaPPB)
	nothingPending := t.latestEpochRound.Less(contractEpochRound) || t.latestEpochRound == contractEpochRound
	result := deviates || nothingPending

	t.logger.Debug("shouldTransmit() = result", commontypes.LogFields{
		"contractEpochRound": contractEpochRound,
		"epochRound":         reportEpochRound,
		"latestEpochRound":   t.latestEpochRound,
		"deviates":           deviates,
		"result":             result,
	})

	return result
}

func (t *transmissionState) contractState() (
	types.ConfigDigest,
	EpochRound,
	error,
) {
	var configDigest types.ConfigDigest
	var epoch uint32
	var round uint8
	var err error
	ok := t.subprocesses.BlockForAtMost(
		t.ctx,
		t.localConfig.BlockchainTimeout,
		func(ctx context.Context) {
			configDigest, epoch, round, _, _, err = t.transmitter.LatestTransmissionDetails(ctx)
		},
	)

	if !ok {
		return types.ConfigDigest{}, EpochRound{}, fmt.Errorf("LatestTransmissionDetails timed out. Timeout: %v", t.localConfig.BlockchainTimeout)
	}

	if err != nil {
		return types.ConfigDigest{}, EpochRound{}, errors.Wrap(err, "Error during LatestTransmissionDetails in Transmission")
	}

	return configDigest, EpochRound{epoch, round}, nil
}

func (t *transmissionState) transmitDelay(epoch uint32, round uint8) *time.Duration {
	// No need for HMAC. Since we use Keccak256, prepending
	// with key gives us a PRF already.
	hash := sha3.NewLegacyKeccak256()
	transmissionOrderKey := t.config.TransmissionOrderKey()
	hash.Write(transmissionOrderKey[:])
	hash.Write(t.config.ConfigDigest[:])
	temp := make([]byte, 8)
	binary.LittleEndian.PutUint64(temp, uint64(epoch))
	hash.Write(temp)
	binary.LittleEndian.PutUint64(temp, uint64(round))
	hash.Write(temp)

	var key [16]byte
	copy(key[:], hash.Sum(nil))
	pi := permutation.Permutation(t.config.N(), key)

	sum := 0
	for i, s := range t.config.S {
		sum += s
		if pi[t.id] < sum {
			result := time.Duration(i) * t.config.DeltaStage
			return &result
		}
	}
	return nil
}
