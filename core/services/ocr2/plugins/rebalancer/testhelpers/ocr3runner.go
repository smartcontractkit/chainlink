package testhelpers

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	ErrQuery                        = errors.New("error in query phase")
	ErrObservation                  = errors.New("error in observation phase")
	ErrValidateObservation          = errors.New("error in validate observation phase")
	ErrOutcome                      = errors.New("error in outcome phase")
	ErrReports                      = errors.New("error in reports phase")
	ErrShouldAcceptAttestedReport   = errors.New("error in should accept attested report phase")
	ErrShouldTransmitAcceptedReport = errors.New("error in should transmit accepted report phase")
)

type OCR3Runner[RI any] struct {
	nodes           []ocr3types.ReportingPlugin[RI]
	round           int
	previousOutcome ocr3types.Outcome
}

func NewOCR3Runner[RI any](nodes []ocr3types.ReportingPlugin[RI]) *OCR3Runner[RI] {
	return &OCR3Runner[RI]{
		nodes: nodes,
		round: 0,
	}
}

// RunRound will run some basic steps of an OCR3 flow.
// This is not a full OCR3 round but only the bare minimum.
func (r *OCR3Runner[RI]) RunRound(ctx context.Context) (transmitted, notAccepted, notTransmitted []ocr3types.ReportWithInfo[RI], outcome []byte, err error) {
	r.round++
	seqNr := uint64(r.round)

	leaderNode := r.selectLeader()

	outcomeCtx := ocr3types.OutcomeContext{SeqNr: seqNr, PreviousOutcome: r.previousOutcome}

	q, err := leaderNode.Query(ctx, outcomeCtx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s: %w", err, ErrQuery)
	}

	attributedObservations := make([]types.AttributedObservation, len(r.nodes))
	for i, n := range r.nodes {
		obs, err2 := n.Observation(ctx, outcomeCtx, q)
		if err2 != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s: %w", err2, ErrObservation)
		}

		attrObs := types.AttributedObservation{Observation: obs, Observer: 0}
		err = leaderNode.ValidateObservation(outcomeCtx, q, attrObs)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s: %w", err, ErrValidateObservation)
		}

		attributedObservations[i] = attrObs
	}

	// todo: run the outcome phase on all nodes and make sure their results are equal.
	outcome, err = leaderNode.Outcome(outcomeCtx, q, attributedObservations)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s: %w", err, ErrOutcome)
	}
	r.previousOutcome = outcome

	reportsWithInfo, err := leaderNode.Reports(seqNr, outcome)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s: %w", err, ErrReports)
	}

	transmitted = make([]ocr3types.ReportWithInfo[RI], 0)
	notAccepted = make([]ocr3types.ReportWithInfo[RI], 0)
	notTransmitted = make([]ocr3types.ReportWithInfo[RI], 0)

	for _, report := range reportsWithInfo {
		shouldAccept, err := leaderNode.ShouldAcceptAttestedReport(ctx, seqNr, report)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s: %w", err, ErrShouldAcceptAttestedReport)
		}
		if !shouldAccept {
			notAccepted = append(notAccepted, report)
			continue
		}

		shouldTransmit, err := leaderNode.ShouldTransmitAcceptedReport(ctx, seqNr, report)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s: %w", err, ErrShouldTransmitAcceptedReport)
		}
		if !shouldTransmit {
			notTransmitted = append(notTransmitted, report)
			continue
		}

		transmitted = append(transmitted, report)
	}

	return transmitted, notAccepted, notTransmitted, outcome, nil // reports transmitted on-chain
}

func (r *OCR3Runner[RI]) selectLeader() ocr3types.ReportingPlugin[RI] {
	numNodes := len(r.nodes)
	if numNodes == 0 {
		return nil
	}
	return r.nodes[rand.Intn(numNodes)]
}
