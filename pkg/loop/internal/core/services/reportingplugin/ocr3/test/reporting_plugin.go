package ocr3_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	outcomeCtx = ocr3types.OutcomeContext{
		SeqNr:           1,
		PreviousOutcome: ocr3types.Outcome([]byte{1, 2, 3}),
	}

	query = libocr.Query{1, 2, 3}
)
var ReportingPlugin = ocr3staticReportingPlugin{
	ocr3staticReportingPluginConfig: ocr3staticReportingPluginConfig{
		expectedOutcomeContext: outcomeCtx,
		queryRequest:           queryRequest{outcomeCtx: outcomeCtx},
		queryResponse:          queryResponse{query: query},
		observationRequest: observationRequest{
			outcomeCtx: outcomeCtx,
			query:      query,
		},
		observationResponse: observationResponse{
			observation: libocr.Observation{1, 2, 3},
		},
		observationQuorumRequest: observationQuorumRequest{
			outcomeCtx: outcomeCtx,
			query:      query,
		},
		observationQuorumResponse: observationQuorumResponse{
			quorum: 1,
		},
		validateObservationRequest: validateObservationRequest{
			outcomeCtx:            outcomeCtx,
			query:                 query,
			attributedObservation: libocr.AttributedObservation{Observer: 1, Observation: []byte{1, 2, 3}},
		},
		outcomeRequest: outcomeRequest{
			outcomeCtx:   outcomeCtx,
			query:        query,
			observations: []libocr.AttributedObservation{{Observer: 1, Observation: []byte{1, 2, 3}}},
		},
		outcomeResponse: outcomeResponse{
			outcome: ocr3types.Outcome{1, 2, 3},
		},
		reportsRequest: reportsRequest{
			seq:     1,
			outcome: ocr3types.Outcome{1, 2, 3},
		},
		reportsResponse: reportsResponse{
			reportWithInfo: []ocr3types.ReportWithInfo[[]byte]{
				{Report: []byte{1, 2, 3}, Info: []byte{1, 2, 3}},
			},
		},
		shouldAcceptAttestedReportRequest: shouldAcceptAttestedReportRequest{
			seq: 1,
			r:   ocr3types.ReportWithInfo[[]byte]{Report: []byte{1, 2, 3}, Info: []byte{1, 2, 3}},
		},
		shouldAcceptAttestedReportResponse: shouldAcceptAttestedReportResponse{
			shouldAccept: true,
		},
		shouldTransmitAcceptedReportRequest: shouldTransmitAcceptedReportRequest{
			seq: 1,

			r: ocr3types.ReportWithInfo[[]byte]{Report: []byte{1, 2, 3}, Info: []byte{1, 2, 3}},
		},
		shouldTransmitAcceptedReportResponse: shouldTransmitAcceptedReportResponse{
			shouldTransmit: true,
		},
	},
}

type queryRequest struct {
	outcomeCtx ocr3types.OutcomeContext
}

type queryResponse struct {
	query libocr.Query
}

type observationRequest struct {
	outcomeCtx ocr3types.OutcomeContext
	query      libocr.Query
}

type observationResponse struct {
	observation libocr.Observation
}

type observationQuorumRequest struct {
	outcomeCtx ocr3types.OutcomeContext
	query      libocr.Query
}

type observationQuorumResponse struct {
	quorum ocr3types.Quorum
}

type validateObservationRequest struct {
	outcomeCtx            ocr3types.OutcomeContext
	query                 libocr.Query
	attributedObservation libocr.AttributedObservation
}

type outcomeRequest struct {
	outcomeCtx   ocr3types.OutcomeContext
	query        libocr.Query
	observations []libocr.AttributedObservation
}

type outcomeResponse struct {
	outcome ocr3types.Outcome
}

type reportsRequest struct {
	seq     uint64
	outcome ocr3types.Outcome
}

type reportsResponse struct {
	reportWithInfo []ocr3types.ReportWithInfo[[]byte]
}

type shouldAcceptAttestedReportRequest struct {
	seq uint64
	r   ocr3types.ReportWithInfo[[]byte]
}

type shouldAcceptAttestedReportResponse struct {
	shouldAccept bool
}

type shouldTransmitAcceptedReportRequest struct {
	seq uint64
	r   ocr3types.ReportWithInfo[[]byte]
}

type shouldTransmitAcceptedReportResponse struct {
	shouldTransmit bool
}

type ocr3staticReportingPluginConfig struct {
	expectedOutcomeContext ocr3types.OutcomeContext
	// Query method
	queryRequest  queryRequest
	queryResponse queryResponse
	// Observation method
	observationRequest  observationRequest
	observationResponse observationResponse
	// ObservationQuorum method
	observationQuorumRequest  observationQuorumRequest
	observationQuorumResponse observationQuorumResponse
	// ValidateObservation method
	validateObservationRequest validateObservationRequest
	// Outcome method
	outcomeRequest  outcomeRequest
	outcomeResponse outcomeResponse
	// Reports method
	reportsRequest  reportsRequest
	reportsResponse reportsResponse
	// ShouldAcceptAttestedReport method
	shouldAcceptAttestedReportRequest  shouldAcceptAttestedReportRequest
	shouldAcceptAttestedReportResponse shouldAcceptAttestedReportResponse
	// ShouldTransmitAcceptedReport method
	shouldTransmitAcceptedReportRequest  shouldTransmitAcceptedReportRequest
	shouldTransmitAcceptedReportResponse shouldTransmitAcceptedReportResponse
}

type ocr3staticReportingPlugin struct {
	ocr3staticReportingPluginConfig
}

func (s ocr3staticReportingPlugin) Query(ctx context.Context, outcomeCtx ocr3types.OutcomeContext) (libocr.Query, error) {
	err := s.checkOutCtx(outcomeCtx)
	if err != nil {
		return nil, err
	}
	return s.queryResponse.query, nil
}

func (s ocr3staticReportingPlugin) Observation(ctx context.Context, outcomeCtx ocr3types.OutcomeContext, q libocr.Query) (libocr.Observation, error) {
	err := s.checkOutCtx(outcomeCtx)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(s.observationRequest.query, q) {
		return nil, fmt.Errorf("expected %x but got %x", s.observationRequest.query, q)
	}

	return s.observationResponse.observation, nil
}

func (s ocr3staticReportingPlugin) ValidateObservation(ctx context.Context, outcomeCtx ocr3types.OutcomeContext, q libocr.Query, a libocr.AttributedObservation) error {
	err := s.checkOutCtx(outcomeCtx)
	if err != nil {
		return err
	}
	if !bytes.Equal(s.validateObservationRequest.query, q) {
		return fmt.Errorf("expected %x but got %x", s.validateObservationRequest.query, q)
	}

	if a.Observer != s.validateObservationRequest.attributedObservation.Observer {
		return fmt.Errorf("expected %x but got %x", s.validateObservationRequest.attributedObservation.Observer, a.Observer)
	}

	if !bytes.Equal(a.Observation, s.validateObservationRequest.attributedObservation.Observation) {
		return fmt.Errorf("expected %x but got %x", s.validateObservationRequest.attributedObservation.Observation, a.Observation)
	}

	return nil
}

func (s ocr3staticReportingPlugin) ObservationQuorum(ctx context.Context, outcomeCtx ocr3types.OutcomeContext, q libocr.Query) (ocr3types.Quorum, error) {
	err := s.checkOutCtx(outcomeCtx)
	if err != nil {
		return ocr3types.Quorum(0), err
	}
	if !bytes.Equal(q, s.observationQuorumRequest.query) {
		return ocr3types.Quorum(0), fmt.Errorf("expected %x but got %x", s.observationQuorumRequest.query, q)
	}

	return s.observationQuorumResponse.quorum, nil
}

func (s ocr3staticReportingPlugin) Outcome(ctx context.Context, outcomeCtx ocr3types.OutcomeContext, q libocr.Query, aos []libocr.AttributedObservation) (ocr3types.Outcome, error) {
	err := s.checkOutCtx(outcomeCtx)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(q, s.outcomeRequest.query) {
		return nil, fmt.Errorf("expected %x but got %x", s.outcomeRequest.query, q)
	}

	if !assert.ObjectsAreEqual(aos, s.outcomeRequest.observations) {
		return nil, fmt.Errorf("expected %v but got %v", s.outcomeRequest.observations, aos)
	}
	return s.outcomeResponse.outcome, nil
}

func (s ocr3staticReportingPlugin) Reports(ctx context.Context, seq uint64, o ocr3types.Outcome) ([]ocr3types.ReportWithInfo[[]byte], error) {
	if seq != s.reportsRequest.seq {
		return nil, fmt.Errorf("expected %x but got %x", s.reportsRequest.seq, seq)
	}

	if !bytes.Equal(o, s.reportsRequest.outcome) {
		return nil, fmt.Errorf("expected %x but got %x", s.reportsRequest.outcome, o)
	}

	return s.reportsResponse.reportWithInfo, nil
}

func (s ocr3staticReportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	if u != s.shouldAcceptAttestedReportRequest.seq {
		return false, fmt.Errorf("expected %x but got %x", s.shouldAcceptAttestedReportRequest.seq, u)
	}
	if !assert.ObjectsAreEqual(r, s.shouldAcceptAttestedReportRequest.r) {
		return false, fmt.Errorf("expected %x but got %x", s.shouldAcceptAttestedReportRequest.r, r)
	}
	return s.shouldAcceptAttestedReportResponse.shouldAccept, nil
}

func (s ocr3staticReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	if u != s.shouldTransmitAcceptedReportRequest.seq {
		return false, fmt.Errorf("expected %x but got %x", s.shouldTransmitAcceptedReportRequest.seq, u)
	}
	if !assert.ObjectsAreEqual(r, s.shouldTransmitAcceptedReportRequest.r) {
		return false, fmt.Errorf("expected %x but got %x", s.shouldTransmitAcceptedReportRequest.r, r)
	}
	return s.shouldTransmitAcceptedReportResponse.shouldTransmit, nil
}

func (s ocr3staticReportingPlugin) Close() error { return nil }

func (s ocr3staticReportingPlugin) checkOutCtx(outcomeCtx ocr3types.OutcomeContext) error {
	if outcomeCtx.SeqNr != s.expectedOutcomeContext.SeqNr {
		return fmt.Errorf("expected %v but got %v", s.expectedOutcomeContext.SeqNr, outcomeCtx.SeqNr)
	}
	if !bytes.Equal(outcomeCtx.PreviousOutcome, s.expectedOutcomeContext.PreviousOutcome) {
		return fmt.Errorf("expected %x but got %x", outcomeCtx.PreviousOutcome, s.expectedOutcomeContext.PreviousOutcome)
	}
	return nil
}

func (s ocr3staticReportingPlugin) AssertEqual(ctx context.Context, t *testing.T, rp ocr3types.ReportingPlugin[[]byte]) {
	gotQuery, err := rp.Query(ctx, s.queryRequest.outcomeCtx)
	require.NoError(t, err)
	assert.Equal(t, s.queryResponse.query, gotQuery)

	gotObs, err := rp.Observation(ctx, s.observationRequest.outcomeCtx, s.observationRequest.query)
	require.NoError(t, err)
	assert.Equal(t, s.observationResponse.observation, gotObs)

	err = rp.ValidateObservation(ctx, s.validateObservationRequest.outcomeCtx, s.validateObservationRequest.query, s.validateObservationRequest.attributedObservation)
	require.NoError(t, err)

	gotQuorum, err := rp.ObservationQuorum(ctx, s.observationQuorumRequest.outcomeCtx, s.observationQuorumRequest.query)
	require.NoError(t, err)
	assert.Equal(t, s.observationQuorumResponse.quorum, gotQuorum)

	gotOutcome, err := rp.Outcome(ctx, s.outcomeRequest.outcomeCtx, s.outcomeRequest.query, s.outcomeRequest.observations)
	require.NoError(t, err)
	assert.Equal(t, s.outcomeResponse.outcome, gotOutcome)

	gotRI, err := rp.Reports(ctx, s.reportsRequest.seq, s.reportsRequest.outcome)
	require.NoError(t, err)
	assert.Equal(t, s.reportsResponse.reportWithInfo, gotRI)

	gotShouldAccept, err := rp.ShouldAcceptAttestedReport(ctx, s.shouldAcceptAttestedReportRequest.seq, s.shouldAcceptAttestedReportRequest.r)
	require.NoError(t, err)
	assert.Equal(t, s.shouldAcceptAttestedReportResponse.shouldAccept, gotShouldAccept)

	gotShouldTransmit, err := rp.ShouldTransmitAcceptedReport(ctx, s.shouldTransmitAcceptedReportRequest.seq, s.shouldTransmitAcceptedReportRequest.r)
	require.NoError(t, err)
	assert.Equal(t, s.shouldTransmitAcceptedReportResponse.shouldTransmit, gotShouldTransmit)
}
