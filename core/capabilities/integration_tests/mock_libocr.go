package integration_tests

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

type libocrNode struct {
	ocr3types.ReportingPlugin[[]byte]
	*ocr3.ContractTransmitter
	key ocr2key.KeyBundle
}

// mockLibOCR is a mock libocr implementation for testing purposes that simulates libocr protocol rounds without having
// to setup the libocr network
type mockLibOCR struct {
	services.StateMachine
	t *testing.T

	nodes                 []*libocrNode
	f                     uint8
	protocolRoundInterval time.Duration

	seqNr      uint64
	outcomeCtx ocr3types.OutcomeContext

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func newMockLibOCR(t *testing.T, f uint8, protocolRoundInterval time.Duration) *mockLibOCR {
	return &mockLibOCR{
		t: t,
		f: f, outcomeCtx: ocr3types.OutcomeContext{
			SeqNr:           0,
			PreviousOutcome: nil,
			Epoch:           0,
			Round:           0,
		},
		protocolRoundInterval: protocolRoundInterval,
		stopCh:                make(services.StopChan),
	}
}

func (m *mockLibOCR) Start(ctx context.Context) error {
	return m.StartOnce("mockLibOCR", func() error {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()

			ticker := time.NewTicker(m.protocolRoundInterval)
			defer ticker.Stop()

			for {
				select {
				case <-m.stopCh:
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					err := m.simulateProtocolRound(ctx)
					if err != nil {
						require.FailNow(m.t, err.Error())
					}
				}
			}
		}()
		return nil
	})
}

func (m *mockLibOCR) Close() error {
	return m.StopOnce("mockLibOCR", func() error {
		close(m.stopCh)
		m.wg.Wait()
		return nil
	})
}

func (m *mockLibOCR) AddNode(plugin ocr3types.ReportingPlugin[[]byte], transmitter *ocr3.ContractTransmitter, key ocr2key.KeyBundle) {
	m.nodes = append(m.nodes, &libocrNode{plugin, transmitter, key})
}

func (m *mockLibOCR) simulateProtocolRound(ctx context.Context) error {
	// randomly select a leader
	leader := m.nodes[rand.Intn(len(m.nodes))]

	// get the query
	query, err := leader.Query(ctx, m.outcomeCtx)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	var observations []types.AttributedObservation
	for oracleID, node := range m.nodes {
		obs, err2 := node.Observation(ctx, m.outcomeCtx, query)
		if err2 != nil {
			return fmt.Errorf("failed to get observation: %w", err)
		}

		observations = append(observations, types.AttributedObservation{
			Observation: obs,
			Observer:    commontypes.OracleID(oracleID),
		})
	}

	var outcomes []ocr3types.Outcome
	for _, node := range m.nodes {
		outcome, err2 := node.Outcome(m.outcomeCtx, query, observations)
		if err2 != nil {
			return fmt.Errorf("failed to get outcome: %w", err)
		}

		if len(outcome) == 0 {
			return nil // wait until all nodes have an outcome for testing purposes
		}

		outcomes = append(outcomes, outcome)
	}

	// if all outcomes are equal proceed to reports
	for _, outcome := range outcomes {
		if !bytes.Equal(outcome, outcomes[0]) {
			return nil
		}
	}

	reports, err := leader.Reports(0, outcomes[0])
	if err != nil {
		return fmt.Errorf("failed to get reports: %w", err)
	}
	for _, report := range reports {
		// create signatures
		var signatures []types.AttributedOnchainSignature
		for i, node := range m.nodes {
			sig, err := node.key.Sign(types.ReportContext{}, report.Report)
			if err != nil {
				return fmt.Errorf("failed to sign report: %w", err)
			}

			signatures = append(signatures, types.AttributedOnchainSignature{
				Signer:    commontypes.OracleID(i),
				Signature: sig,
			})
		}

		for _, node := range m.nodes {
			accept, err := node.ShouldAcceptAttestedReport(ctx, m.seqNr, report)
			if err != nil {
				return fmt.Errorf("failed to check if report should be accepted: %w", err)
			}
			if !accept {
				continue
			}

			transmit, err := node.ShouldTransmitAcceptedReport(ctx, m.seqNr, report)
			if err != nil {
				return fmt.Errorf("failed to check if report should be transmitted: %w", err)
			}

			if !transmit {
				continue
			}

			// For each node select a random set of f+1 signatures to mimic libocr behaviour
			s := rand.NewSource(time.Now().UnixNano())
			r := rand.New(s)
			indices := r.Perm(len(signatures))
			selectedSignatures := make([]types.AttributedOnchainSignature, m.f+1)
			for i := 0; i < int(m.f+1); i++ {
				selectedSignatures[i] = signatures[indices[i]]
			}

			err = node.Transmit(ctx, types.ConfigDigest{}, 0, report, selectedSignatures)
			if err != nil {
				return fmt.Errorf("failed to transmit report: %w", err)
			}
		}

		m.seqNr++
		m.outcomeCtx = ocr3types.OutcomeContext{
			SeqNr:           0,
			PreviousOutcome: outcomes[0],
		}
	}

	return nil
}
