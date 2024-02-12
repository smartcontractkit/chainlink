package test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
)

type ocr3staticReportingPlugin struct{}

func (s ocr3staticReportingPlugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (libocr.Query, error) {
	err := checkOutCtx(outctx)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func (s ocr3staticReportingPlugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, q libocr.Query) (libocr.Observation, error) {
	err := checkOutCtx(outctx)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(query, q) {
		return nil, fmt.Errorf("expected %x but got %x", query, q)
	}

	return observation, nil
}

func (s ocr3staticReportingPlugin) ValidateObservation(outctx ocr3types.OutcomeContext, q libocr.Query, a libocr.AttributedObservation) error {
	err := checkOutCtx(outctx)
	if err != nil {
		return err
	}
	if !bytes.Equal(query, q) {
		return fmt.Errorf("expected %x but got %x", query, q)
	}

	if a.Observer != ao.Observer {
		return fmt.Errorf("expected %x but got %x", a.Observer, ao.Observer)
	}

	if !bytes.Equal(a.Observation, ao.Observation) {
		return fmt.Errorf("expected %x but got %x", a.Observation, ao.Observation)
	}

	return nil
}

func (s ocr3staticReportingPlugin) ObservationQuorum(outctx ocr3types.OutcomeContext, q libocr.Query) (ocr3types.Quorum, error) {
	err := checkOutCtx(outctx)
	if err != nil {
		return ocr3types.Quorum(0), err
	}
	if !bytes.Equal(q, query) {
		return ocr3types.Quorum(0), fmt.Errorf("expected %x but got %x", q, query)
	}

	return quorum, nil
}

func (s ocr3staticReportingPlugin) Outcome(outctx ocr3types.OutcomeContext, q libocr.Query, aos []libocr.AttributedObservation) (ocr3types.Outcome, error) {
	err := checkOutCtx(outctx)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(q, query) {
		return nil, fmt.Errorf("expected %x but got %x", q, query)
	}

	if !assert.ObjectsAreEqual(aos, obs) {
		return nil, fmt.Errorf("expected %v but got %v", aos, obs)
	}
	return outcome, nil
}

func (s ocr3staticReportingPlugin) Reports(seq uint64, o ocr3types.Outcome) ([]ocr3types.ReportWithInfo[any], error) {
	if seq != seqNr {
		return nil, fmt.Errorf("expected %x but got %x", seq, seqNr)
	}

	if !bytes.Equal(o, outcome) {
		return nil, fmt.Errorf("expected %x but got %x", o, outcome)
	}

	return RIs, nil
}

func (s ocr3staticReportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[any]) (bool, error) {
	if u != seqNr {
		return false, fmt.Errorf("expected %x but got %x", u, seqNr)
	}
	if !assert.ObjectsAreEqual(r, RI) {
		return false, fmt.Errorf("expected %x but got %x", r, RI)
	}
	return true, nil
}

func (s ocr3staticReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[any]) (bool, error) {
	if u != seqNr {
		return false, fmt.Errorf("expected %x but got %x", u, seqNr)
	}
	if !assert.ObjectsAreEqual(r, RI) {
		return false, fmt.Errorf("expected %x but got %x", r, RI)
	}
	return true, nil
}

func (s ocr3staticReportingPlugin) Close() error { return nil }

func checkOutCtx(outctx ocr3types.OutcomeContext) error {
	//nolint:all
	if outctx.Epoch != outcomeContext.Epoch {
		return fmt.Errorf("expected %v but got %v", outcomeContext.Epoch, outctx.Epoch)
	}
	//nolint:all
	if outctx.Round != outcomeContext.Round {
		return fmt.Errorf("expected %v but got %v", outcomeContext.Round, outctx.Round)
	}
	if outctx.SeqNr != outcomeContext.SeqNr {
		return fmt.Errorf("expected %v but got %v", outcomeContext.SeqNr, outctx.SeqNr)
	}
	if !bytes.Equal(outctx.PreviousOutcome, outcomeContext.PreviousOutcome) {
		return fmt.Errorf("expected %x but got %x", outctx.PreviousOutcome, outcomeContext.PreviousOutcome)
	}
	return nil
}
