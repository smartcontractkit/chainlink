package ocr

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/wasp"
	"go.uber.org/ratelimit"

	client2 "github.com/smartcontractkit/chainlink-testing-framework/client"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// VU is a virtual user for the OCR load test
// it creates a feed and triggers new rounds
type VU struct {
	rl            ratelimit.Limiter
	rate          int
	rateUnit      time.Duration
	roundNum      atomic.Int64
	cc            blockchain.EVMClient
	lt            contracts.LinkToken
	cd            contracts.ContractDeployer
	bootstrapNode *client.ChainlinkK8sClient
	workerNodes   []*client.ChainlinkK8sClient
	msClient      *client2.MockserverClient
	l             zerolog.Logger
	ocrInstances  []contracts.OffchainAggregator
	stop          chan struct{}
}

func NewVU(
	l zerolog.Logger,
	rate int,
	rateUnit time.Duration,
	cc blockchain.EVMClient,
	lt contracts.LinkToken,
	cd contracts.ContractDeployer,
	bootstrapNode *client.ChainlinkK8sClient,
	workerNodes []*client.ChainlinkK8sClient,
	msClient *client2.MockserverClient,
) *VU {
	return &VU{
		rl:            ratelimit.New(rate, ratelimit.Per(rateUnit)),
		rate:          rate,
		rateUnit:      rateUnit,
		l:             l,
		cc:            cc,
		lt:            lt,
		cd:            cd,
		msClient:      msClient,
		bootstrapNode: bootstrapNode,
		workerNodes:   workerNodes,
	}
}

func (m *VU) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &VU{
		stop:          make(chan struct{}, 1),
		rl:            ratelimit.New(m.rate, ratelimit.Per(m.rateUnit)),
		rate:          m.rate,
		rateUnit:      m.rateUnit,
		l:             m.l,
		cc:            m.cc,
		lt:            m.lt,
		cd:            m.cd,
		msClient:      m.msClient,
		bootstrapNode: m.bootstrapNode,
		workerNodes:   m.workerNodes,
	}
}

func (m *VU) Setup(_ *wasp.Generator) error {
	ocrInstances, err := actions.DeployOCRContracts(1, m.lt, m.cd, m.workerNodes, m.cc)
	if err != nil {
		return err
	}
	err = actions.CreateOCRJobs(ocrInstances, m.bootstrapNode, m.workerNodes, 5, m.msClient, m.cc.GetChainID().String())
	if err != nil {
		return err
	}
	m.ocrInstances = ocrInstances
	return nil
}

func (m *VU) Teardown(_ *wasp.Generator) error {
	return nil
}

func (m *VU) Call(l *wasp.Generator) {
	m.rl.Take()
	m.roundNum.Add(1)
	requestedRound := m.roundNum.Load()
	m.l.Info().
		Int64("RoundNum", requestedRound).
		Str("FeedID", m.ocrInstances[0].Address()).
		Msg("starting new round")
	err := m.ocrInstances[0].RequestNewRound()
	if err != nil {
		l.ResponsesChan <- &wasp.Response{Error: err.Error(), Failed: true}
	}
	for {
		time.Sleep(5 * time.Second)
		lr, err := m.ocrInstances[0].GetLatestRound(context.Background())
		if err != nil {
			l.ResponsesChan <- &wasp.Response{Error: err.Error(), Failed: true}
		}
		m.l.Info().Interface("LatestRound", lr).Msg("latest round")
		if lr.RoundId.Int64() >= requestedRound {
			l.ResponsesChan <- &wasp.Response{}
		}
	}
}

func (m *VU) Stop(_ *wasp.Generator) {
	m.stop <- struct{}{}
}

func (m *VU) StopChan() chan struct{} {
	return m.stop
}
