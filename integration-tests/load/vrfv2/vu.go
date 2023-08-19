package loadvrfv2

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions"
	vrfConst "github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/wasp"
	"time"
)

/* JobVolumeVU is a "virtual user" that creates a VRFv2 job and constantly requesting new randomness only for this job instance  */

type JobVolumeVU struct {
	pace                     time.Duration
	minIncomingConfirmations uint16
	nodes                    []*client.ChainlinkClient
	bc                       blockchain.EVMClient
	contracts                *vrfv2_actions.VRFV2Contracts
	jobs                     []vrfv2_actions.VRFV2JobInfo
	keyHash                  [32]byte
	stop                     chan struct{}
}

func NewJobVolumeVU(
	pace time.Duration,
	confirmations uint16,
	nodes []*client.ChainlinkClient,
	bc blockchain.EVMClient,
	contracts *vrfv2_actions.VRFV2Contracts,
) *JobVolumeVU {
	return &JobVolumeVU{
		pace:                     pace,
		minIncomingConfirmations: confirmations,
		nodes:                    nodes,
		bc:                       bc,
		contracts:                contracts,
		stop:                     make(chan struct{}, 1),
	}
}

func (m *JobVolumeVU) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &JobVolumeVU{
		pace:                     m.pace,
		minIncomingConfirmations: m.minIncomingConfirmations,
		nodes:                    m.nodes,
		bc:                       m.bc,
		contracts:                m.contracts,
		stop:                     make(chan struct{}, 1),
	}
}

func (m *JobVolumeVU) Setup(_ *wasp.Generator) error {
	jobs, err := vrfv2_actions.CreateVRFV2Jobs(m.nodes, m.contracts.Coordinator, m.bc, m.minIncomingConfirmations)
	if err != nil {
		return errors.Wrap(err, "failed to create VRFv2 jobs in setup")
	}
	m.jobs = jobs
	m.keyHash = jobs[0].KeyHash
	return nil
}

func (m *JobVolumeVU) Teardown(_ *wasp.Generator) error {
	return nil
}

func (m *JobVolumeVU) Call(l *wasp.Generator) {
	time.Sleep(m.pace)
	tn := time.Now()
	err := m.contracts.LoadTestConsumer.RequestRandomness(
		m.keyHash,
		vrfConst.SubID,
		vrfConst.MinimumConfirmations,
		vrfConst.CallbackGasLimit,
		vrfConst.NumberOfWords,
		vrfConst.RandomnessRequestCountPerRequest,
	)
	if err != nil {
		l.ResponsesChan <- &wasp.CallResult{Duration: time.Since(tn), Error: err.Error(), Failed: true}
		return
	}
	l.ResponsesChan <- &wasp.CallResult{Duration: time.Since(tn)}
}

func (m *JobVolumeVU) Stop(_ *wasp.Generator) {
	m.stop <- struct{}{}
}

func (m *JobVolumeVU) StopChan() chan struct{} {
	return m.stop
}
