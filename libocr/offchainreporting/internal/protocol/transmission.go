package protocol

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/permutation"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"
	"golang.org/x/crypto/sha3"
)

func RunTransmission(
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	config config.SharedConfig,
	chReportGenerationToTransmission <-chan EventToTransmission,
	database types.Database,
	id types.OracleID,
	localConfig types.LocalConfig,
	logger types.Logger,
	transmitter types.ContractTransmitter,
) {
	t := transmissionState{
		ctx:          ctx,
		subprocesses: subprocesses,

		config:                           config,
		chReportGenerationToTransmission: chReportGenerationToTransmission,
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

	config                           config.SharedConfig
	chReportGenerationToTransmission <-chan EventToTransmission
	database                         types.Database
	id                               types.OracleID
	localConfig                      types.LocalConfig
	logger                           types.Logger
	transmitter                      types.ContractTransmitter

	latestEpochRound EpochRound
	latestMedian     observation.Observation
	latestReport     ContractReportWithSignatures
	times            MinHeapTimeToContractReport
	tTransmit        <-chan time.Time
}

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
		t.logger.Error("Error fetching pending transmissions from database", types.LogFields{"error": err})
		return
	}

	now := time.Now()

	for key, trans := range pending {
		if now.Before(trans.Time) {
			t.times.Push(MinHeapTimeToContractReportItem{
				key,
				trans,
			})
		}
	}

	latestExpiredTransmissionKey := types.PendingTransmissionKey{}
	latestExpiredTransmission := (*types.PendingTransmission)(nil)
	for key, trans := range pending {
		if trans.Time.Before(now) && (EpochRound{latestExpiredTransmissionKey.Epoch, latestExpiredTransmissionKey.Round}).Less(EpochRound{key.Epoch, key.Round}) {
			latestExpiredTransmissionKey = key
			latestExpiredTransmission = &trans
		}
	}
	if latestExpiredTransmission != nil {
		t.times.Push(MinHeapTimeToContractReportItem{
			latestExpiredTransmissionKey,
			*latestExpiredTransmission,
		})
	}

	if t.times.Len() != 0 {
		t.tTransmit = time.After(now.Sub(t.times.Peek().Time))
	}
}

func (t *transmissionState) eventTransmit(ev EventTransmit) {
	t.logger.Debug("Received transmit event", types.LogFields{
		"event": ev,
	})

	{
		contractConfigDigest, contractEpochRound, err := t.contractState()
		if err != nil {
			t.logger.Error("contractEpoch() failed during eventTransmit", types.LogFields{"error": err})
			return
		}

		if contractConfigDigest != t.config.ConfigDigest {
			t.logger.Info("eventTransmit(ev): discarding ev because contractConfigDigest != configDigest", types.LogFields{
				"ev":                   ev,
				"contractConfigDigest": contractConfigDigest,
				"configDigest":         t.config.ConfigDigest,
			})
			return
		}

		if !t.shouldTransmit(ev.ContractReportWithSignatures, contractEpochRound) {
			t.logger.Info("eventTransmit(ev): discarding ev because shouldTransmit returned false", types.LogFields{
				"ev":                   ev,
				"contractConfigDigest": contractConfigDigest,
				"contractEpochRound":   contractEpochRound,
			})
			return
		}
	}

	var err error
	t.latestEpochRound = EpochRound{ev.Ctx.Epoch, ev.Ctx.Round}
	t.latestMedian, err = ev.Values.Median()
	if err != nil {
		t.logger.Error("could not compute median", types.LogFields{"error": err})
	}

	now := time.Now()
	delayMaybe := t.transmitDelay(ev.Ctx.Epoch, ev.Ctx.Round)
	if delayMaybe == nil {
		return
	}
	delay := *delayMaybe
	serializedReport, rs, ss, vs, err := ev.ContractReportWithSignatures.TransmissionArgs()
	if err != nil {
		t.logger.Error("Failed to serialize contract report", types.LogFields{"error": err})
		return
	}

	key := types.PendingTransmissionKey{
		ConfigDigest: t.config.ConfigDigest,
		Epoch:        ev.Ctx.Epoch,
		Round:        ev.Ctx.Round,
	}
	median, err := ev.ContractReportWithSignatures.Values.Median()
	if err != nil {
		t.logger.Error("could not take median of observations",
			types.LogFields{"error": err})
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
				t.logger.Error("Error while persisting pending transmission to database", types.LogFields{"error": err})
			}
		},
	)
	if !ok {
		t.logger.Error("Database.StorePendingTransmission timed out", types.LogFields{
			"timeout": t.localConfig.DatabaseTimeout,
		})
	}
	t.times.Push(MinHeapTimeToContractReportItem{key, transmission})

	next := t.times.Peek()
	if (EpochRound{ev.Ctx.Epoch, ev.Ctx.Round} == EpochRound{next.Epoch, next.Round}) {
		t.tTransmit = time.After(delay)
	}
}

func (t *transmissionState) eventTTransmitTimeout() {
	defer func() {
		if t.times.Len() != 0 {
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
			if err := t.database.DeletePendingTransmission(ctx, types.PendingTransmissionKey{
				ConfigDigest: t.config.ConfigDigest,
				Epoch:        item.Epoch,
				Round:        item.Round,
			}); err != nil {
				t.logger.Error("eventTTransmitTimeout: Error while deleting pending transmission from database", types.LogFields{"error": err})
			}
		},
	)
	if !ok {
		t.logger.Error("Database.DeletePendingTransmission timed out", types.LogFields{
			"timeout": t.localConfig.DatabaseTimeout,
		})
	}

	contractConfigDigest, contractEpochRound, err := t.contractState()
	if err != nil {
		t.logger.Error("eventTTransmitTimeout: contractState() failed", types.LogFields{"error": err})
		return
	}

	if item.ConfigDigest != contractConfigDigest {
		t.logger.Info("eventTTransmitTimeout: configDigest doesn't match, discarding transmission", types.LogFields{
			"contractConfigDigest": contractConfigDigest,
			"configDigest":         item.ConfigDigest,
			"median":               item.Median,
			"epoch":                item.Epoch,
			"round":                item.Round,
		})
		return
	}

	if !contractEpochRound.Less(itemEpochRound) {
		t.logger.Info("eventTTransmitTimeout: Skipping transmission because report is stale", types.LogFields{
			"contractEpochRound": contractEpochRound,
			"median":             item.Median,
			"epoch":              item.Epoch,
			"round":              item.Round,
		})
		return
	}

	t.logger.Info("eventTTransmitTimeout: Transmitting with median", types.LogFields{
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
		t.logger.Error("eventTTransmitTimeout: Transmit timed out", types.LogFields{
			"timeout": t.localConfig.ContractTransmitterTransmitTimeout,
		})
		return
	}
	if err != nil {
		t.logger.Error("eventTTransmitTimeout: Error while transmitting report on-chain", types.LogFields{"error": err})
		return
	}

	t.logger.Info("eventTTransmitTimeout:❗️successfully transmitted report on-chain", types.LogFields{
		"median": item.Median,
		"epoch":  item.Epoch,
		"round":  item.Round,
	})
}

func (t *transmissionState) shouldTransmit(o ContractReportWithSignatures, contractEpochRound EpochRound) bool {
	oEpochRound := EpochRound{o.Ctx.Epoch, o.Ctx.Round}
	if !contractEpochRound.Less(oEpochRound) {
		t.logger.Debug("shouldTransmit() = false, report is stale", types.LogFields{
			"contractEpochRound": contractEpochRound,
			"epochRound":         oEpochRound,
		})
		return false
	}
	if t.latestEpochRound == (EpochRound{}) {
		t.logger.Debug("shouldTransmit() = true, latestEpochRound is empty", types.LogFields{
			"contractEpochRound": contractEpochRound,
			"epochRound":         oEpochRound,
			"latestEpochRound":   t.latestEpochRound,
		})
		return true
	}
	if oEpochRound.Less(t.latestEpochRound) || oEpochRound == t.latestEpochRound {
		t.logger.Debug("shouldTransmit() = false, report is older than latest report", types.LogFields{
			"contractEpochRound": contractEpochRound,
			"epochRound":         oEpochRound,
			"latestEpochRound":   t.latestEpochRound,
		})
		return false
	}

	oMedian, err := o.Values.Median()
	if err != nil {
		t.logger.Error("could not compute median", types.LogFields{
			"error": err,
		})
		return false
	}

	deviates := t.latestMedian.Deviates(oMedian, t.config.Alpha)
	nothingPending := t.latestEpochRound.Less(contractEpochRound) || t.latestEpochRound == contractEpochRound
	result := deviates || nothingPending

	t.logger.Debug("shouldTransmit() = result", types.LogFields{
		"contractEpochRound": contractEpochRound,
		"epochRound":         oEpochRound,
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
