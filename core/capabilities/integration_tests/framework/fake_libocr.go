package framework

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type oracleContext struct {
	t                     *testing.T
	lggr                  logger.Logger
	key                   ocr2key.KeyBundle
	N                     int
	F                     uint8
	protocolRoundInterval time.Duration
	mux                   sync.Mutex
	pluginNameToFakeOcr   map[string]*FakeLibOCR
}

func (m *oracleContext) addPlugin(ctx context.Context, info ocr3types.ReportingPluginInfo, plugin ocr3types.ReportingPlugin[[]byte],
	args coretypes.OracleArgs) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	libOcr := m.pluginNameToFakeOcr[info.Name]
	if libOcr == nil {
		libOcr = NewFakeLibOCR(m.t, m.lggr, m.F, m.protocolRoundInterval)
		m.pluginNameToFakeOcr[info.Name] = libOcr
	}

	libOcr.AddNode(plugin, args.ContractTransmitter, m.key)

	if libOcr.GetNodeCount() == m.N {
		err := libOcr.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start fake lib ocr: %w", err)
		}
	}
	return nil
}

func (m *oracleContext) Close() error {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, libOcr := range m.pluginNameToFakeOcr {
		if err := libOcr.Close(); err != nil {
			return fmt.Errorf("failed to close fake lib ocr: %w", err)
		}
	}
	return nil
}

type oracleFactoryFactory struct {
	oracleContext *oracleContext
}

func newFakeOracleFactoryFactory(t *testing.T, lggr logger.Logger, key ocr2key.KeyBundle, n int, f uint8, protocolRoundInterval time.Duration) *oracleFactoryFactory {
	return &oracleFactoryFactory{
		oracleContext: &oracleContext{
			t:                     t,
			lggr:                  lggr,
			key:                   key,
			N:                     n,
			F:                     f,
			protocolRoundInterval: protocolRoundInterval,
			pluginNameToFakeOcr:   make(map[string]*FakeLibOCR),
		},
	}
}

func (o *oracleFactoryFactory) NewOracleFactory(params generic.OracleFactoryParams) (coretypes.OracleFactory, error) {
	return &fakeOracleFactory{o.oracleContext}, nil
}

type fakeOracleFactory struct {
	oracleContext *oracleContext
}

func (m *fakeOracleFactory) NewOracle(ctx context.Context, args coretypes.OracleArgs) (coretypes.Oracle, error) {
	return &fakeOracle{oracleContext: m.oracleContext, args: args}, nil
}

type fakeOracle struct {
	oracleContext *oracleContext
	args          coretypes.OracleArgs
}

func (m *fakeOracle) Start(ctx context.Context) error {
	plugin, info, err := m.args.ReportingPluginFactoryService.NewReportingPlugin(ctx, ocr3types.ReportingPluginConfig{
		F: int(m.oracleContext.F),
		N: m.oracleContext.N,
	})

	if err != nil {
		return fmt.Errorf("failed to create reporting plugin: %w", err)
	}

	if err = m.oracleContext.addPlugin(ctx, info, plugin, m.args); err != nil {
		return fmt.Errorf("failed to add plugin: %w", err)
	}

	return nil
}

func (m *fakeOracle) Close(ctx context.Context) error {
	return m.oracleContext.Close()
}

type libocrNode struct {
	id string
	ocr3types.ReportingPlugin[[]byte]
	ocr3types.ContractTransmitter[[]byte]
	key ocr2key.KeyBundle
}

// FakeLibOCR is a fake libocr implementation for testing purposes that simulates libocr protocol rounds without having
// to setup the libocr network
type FakeLibOCR struct {
	services.StateMachine
	t    *testing.T
	lggr logger.Logger

	nodes                 []*libocrNode
	f                     uint8
	protocolRoundInterval time.Duration

	seqNr      uint64
	outcomeCtx ocr3types.OutcomeContext

	mux    sync.Mutex
	stopCh services.StopChan
	wg     sync.WaitGroup
}

func NewFakeLibOCR(t *testing.T, lggr logger.Logger, f uint8, protocolRoundInterval time.Duration) *FakeLibOCR {
	return &FakeLibOCR{
		t:    t,
		lggr: lggr,
		f:    f, outcomeCtx: ocr3types.OutcomeContext{
			SeqNr:           1,
			PreviousOutcome: nil,
			Epoch:           0,
			Round:           0,
		},
		protocolRoundInterval: protocolRoundInterval,
		stopCh:                make(services.StopChan),
	}
}

func (m *FakeLibOCR) Start(ctx context.Context) error {
	return m.StartOnce("FakeLibOCR", func() error {
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
					serviceCtx, cancel := m.stopCh.NewCtx()
					err := m.simulateProtocolRound(serviceCtx)
					cancel()
					if err != nil {
						m.lggr.Errorf("simulating protocol round: %v", err)
					}
				}
			}
		}()
		return nil
	})
}

func (m *FakeLibOCR) Close() error {
	return m.StopOnce("FakeLibOCR", func() error {
		close(m.stopCh)
		m.wg.Wait()
		return nil
	})
}

func (m *FakeLibOCR) AddNode(plugin ocr3types.ReportingPlugin[[]byte], transmitter ocr3types.ContractTransmitter[[]byte], key ocr2key.KeyBundle) string {
	m.mux.Lock()
	defer m.mux.Unlock()
	node := &libocrNode{uuid.New().String(), plugin, transmitter, key}
	m.nodes = append(m.nodes, node)
	return node.id
}

func (m *FakeLibOCR) GetNodeCount() int {
	m.mux.Lock()
	defer m.mux.Unlock()
	return len(m.nodes)
}

func (m *FakeLibOCR) RemoveNode(id string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	var updatedNodes []*libocrNode
	for _, node := range m.nodes {
		if node.id != id {
			updatedNodes = append(updatedNodes, node)
		}
	}

	m.nodes = updatedNodes
}

func (m *FakeLibOCR) simulateProtocolRound(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if len(m.nodes) == 0 {
		return nil
	}

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
			return fmt.Errorf("failed to get observation: %w", err2)
		}

		observations = append(observations, types.AttributedObservation{
			Observation: obs,
			Observer:    commontypes.OracleID(oracleID),
		})
	}

	var outcomes []ocr3types.Outcome
	for _, node := range m.nodes {
		outcome, err2 := node.Outcome(ctx, m.outcomeCtx, query, observations)
		if err2 != nil {
			return fmt.Errorf("failed to get outcome: %w", err2)
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

	reports, err := leader.Reports(ctx, 0, outcomes[0])
	if err != nil {
		return fmt.Errorf("failed to get reports: %w", err)
	}
	for _, report := range reports {
		// create signatures
		var signatures []types.AttributedOnchainSignature
		for i, node := range m.nodes {
			sig, err := node.key.Sign(types.ReportContext{}, report.ReportWithInfo.Report)
			if err != nil {
				return fmt.Errorf("failed to sign report: %w", err)
			}

			signatures = append(signatures, types.AttributedOnchainSignature{
				Signer:    commontypes.OracleID(i),
				Signature: sig,
			})
		}

		for _, node := range m.nodes {
			accept, err := node.ShouldAcceptAttestedReport(ctx, m.seqNr, report.ReportWithInfo)
			if err != nil {
				return fmt.Errorf("failed to check if report should be accepted: %w", err)
			}
			if !accept {
				continue
			}

			transmit, err := node.ShouldTransmitAcceptedReport(ctx, m.seqNr, report.ReportWithInfo)
			if err != nil {
				return fmt.Errorf("failed to check if report should be transmitted: %w", err)
			}

			if !transmit {
				continue
			}

			// For each node select a random set of F+1 signatures to mimic libocr behaviour
			s := rand.NewSource(time.Now().UnixNano())
			r := rand.New(s)
			indices := r.Perm(len(signatures))
			selectedSignatures := make([]types.AttributedOnchainSignature, m.f+1)
			for i := 0; i < int(m.f+1); i++ {
				selectedSignatures[i] = signatures[indices[i]]
			}

			err = node.Transmit(ctx, types.ConfigDigest{}, 0, report.ReportWithInfo, selectedSignatures)
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
