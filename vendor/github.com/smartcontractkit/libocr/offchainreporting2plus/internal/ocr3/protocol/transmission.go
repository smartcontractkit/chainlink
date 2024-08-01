package protocol

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/scheduler"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/permutation"
	"github.com/smartcontractkit/libocr/subprocesses"
)

const ContractTransmitterTimeoutWarningGracePeriod = 50 * time.Millisecond

func RunTransmission[RI any](
	ctx context.Context,
	subprocesses *subprocesses.Subprocesses,

	chReportAttestationToTransmission <-chan EventToTransmission[RI],
	config ocr3config.SharedConfig,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	id commontypes.OracleID,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	reportingPlugin ocr3types.ReportingPlugin[RI],
) {
	sched := scheduler.NewScheduler[EventAttestedReport[RI]]()
	defer sched.Close()

	t := transmissionState[RI]{
		ctx,
		subprocesses,

		chReportAttestationToTransmission,
		config,
		contractTransmitter,
		id,
		localConfig,
		logger.MakeUpdated(commontypes.LogFields{"proto": "transmission"}),
		reportingPlugin,

		sched,
	}
	t.run()
}

type transmissionState[RI any] struct {
	ctx context.Context

	subprocesses *subprocesses.Subprocesses

	chReportAttestationToTransmission <-chan EventToTransmission[RI]
	config                            ocr3config.SharedConfig
	contractTransmitter               ocr3types.ContractTransmitter[RI]
	id                                commontypes.OracleID
	localConfig                       types.LocalConfig
	logger                            loghelper.LoggerWithContext
	reportingPlugin                   ocr3types.ReportingPlugin[RI]

	scheduler *scheduler.Scheduler[EventAttestedReport[RI]]
}

// run runs the event loop for the local transmission protocol
func (t *transmissionState[RI]) run() {
	t.logger.Info("Transmission: running", nil)

	chDone := t.ctx.Done()
	for {
		select {
		case ev := <-t.chReportAttestationToTransmission:
			ev.processTransmission(t)
		case ev := <-t.scheduler.Scheduled():
			t.scheduled(ev)
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

func (t *transmissionState[RI]) eventAttestedReport(ev EventAttestedReport[RI]) {
	now := time.Now()

	shouldAccept, ok := callPlugin[bool](
		t.ctx,
		t.logger,
		commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		},
		"ShouldAcceptAttestedReport",
		t.config.MaxDurationShouldAcceptAttestedReport,
		func(ctx context.Context) (bool, error) {
			return t.reportingPlugin.ShouldAcceptAttestedReport(
				ctx,
				ev.SeqNr,
				ev.AttestedReport.ReportWithInfo,
			)
		},
	)
	if !ok {
		return
	}

	if !shouldAccept {
		t.logger.Debug("ReportingPlugin.ShouldAcceptAttestedReport returned false", commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		})
		return
	}

	delayMaybe := t.transmitDelay(ev.SeqNr, ev.Index)
	if delayMaybe == nil {
		t.logger.Debug("dropping EventAttestedReport because we're not included in transmission schedule", commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		})
		return
	}
	delay := *delayMaybe

	t.logger.Debug("accepted AttestedReport for transmission", commontypes.LogFields{
		"seqNr": ev.SeqNr,
		"index": ev.Index,
		"delay": delay.String(),
	})
	t.scheduler.ScheduleDeadline(ev, now.Add(delay))
}

func (t *transmissionState[RI]) scheduled(ev EventAttestedReport[RI]) {
	shouldTransmit, ok := callPlugin[bool](
		t.ctx,
		t.logger,
		commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		},
		"ShouldTransmitAcceptedReport",
		t.config.MaxDurationShouldTransmitAcceptedReport,
		func(ctx context.Context) (bool, error) {
			return t.reportingPlugin.ShouldTransmitAcceptedReport(
				ctx,
				ev.SeqNr,
				ev.AttestedReport.ReportWithInfo,
			)
		},
	)
	if !ok {
		return
	}

	if !shouldTransmit {
		t.logger.Info("ReportingPlugin.ShouldTransmitAcceptedReport returned false", commontypes.LogFields{
			"seqNr": ev.SeqNr,
			"index": ev.Index,
		})
		return
	}

	t.logger.Debug("transmitting report", commontypes.LogFields{
		"seqNr": ev.SeqNr,
		"index": ev.Index,
	})

	{
		ctx, cancel := context.WithTimeout(
			t.ctx,
			t.localConfig.ContractTransmitterTransmitTimeout,
		)
		defer cancel()

		ins := loghelper.NewIfNotStopped(
			t.localConfig.ContractTransmitterTransmitTimeout+ContractTransmitterTimeoutWarningGracePeriod,
			func() {
				t.logger.Error("ContractTransmitter.Transmit is taking too long", commontypes.LogFields{
					"maxDuration": t.localConfig.ContractTransmitterTransmitTimeout.String(),
					"seqNr":       ev.SeqNr,
					"index":       ev.Index,
				})
			},
		)

		err := t.contractTransmitter.Transmit(
			ctx,
			t.config.ConfigDigest,
			ev.SeqNr,
			ev.AttestedReport.ReportWithInfo,
			ev.AttestedReport.AttributedSignatures,
		)

		ins.Stop()

		if err != nil {
			t.logger.Error("ContractTransmitter.Transmit error", commontypes.LogFields{"error": err})
			return
		}

	}

	t.logger.Info("🚀 successfully invoked ContractTransmitter.Transmit", commontypes.LogFields{
		"seqNr": ev.SeqNr,
		"index": ev.Index,
	})
}

func (t *transmissionState[RI]) transmitDelay(seqNr uint64, index int) *time.Duration {
	transmissionOrderKey := t.config.TransmissionOrderKey()
	mac := hmac.New(sha256.New, transmissionOrderKey[:])
	_ = binary.Write(mac, binary.BigEndian, seqNr)
	_ = binary.Write(mac, binary.BigEndian, uint64(index))

	var key [16]byte
	_ = copy(key[:], mac.Sum(nil))
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
