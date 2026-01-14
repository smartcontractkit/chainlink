package jobs

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
)

// verifyRingJobSpecInputs validates the inputs for a Ring job spec.
func verifyRingJobSpecInputs(inputs job_types.JobSpecInput) error {
	ringInput := &pkg.RingJobConfigInput{}
	if err := inputs.UnmarshalTo(ringInput); err != nil {
		return errors.New("failed to unmarshal job spec input to RingJobConfigInput: " + err.Error())
	}

	if strings.TrimSpace(ringInput.ContractQualifier) == "" {
		return errors.New("contractQualifier is required")
	}

	if ringInput.ChainSelectorEVM == 0 {
		return errors.New("chainSelectorEVM is required")
	}

	if strings.TrimSpace(ringInput.ShardConfigAddr) == "" {
		return errors.New("shardConfigAddr is required")
	}

	if !common.IsHexAddress(ringInput.ShardConfigAddr) {
		return errors.New("shardConfigAddr is invalid: not a valid hex address")
	}

	if len(ringInput.BootstrapperRingUrls) == 0 {
		return errors.New("bootstrapperRingUrls is required")
	}

	return nil
}
