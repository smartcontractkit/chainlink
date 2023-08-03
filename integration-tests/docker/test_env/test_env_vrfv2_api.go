package test_env

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var (
	ErrCreateVRFv2Jobs = "failed to create VRFv2 jobs"
)

func (m *CLClusterTestEnv) CreateVRFv2Jobs(coord contracts.VRFCoordinatorV2) ([]*VRFV2JobInfo, error) {
	jobs := make([]*VRFV2JobInfo, 0)
	for _, n := range m.CLNodes {
		ji, err := n.CreateVRFv2Job(coord, m.Geth.EthClient)
		if err != nil {
			return nil, errors.Wrap(err, ErrCreateVRFv2Jobs)
		}
		jobs = append(jobs, ji)
	}
	return jobs, nil
}
