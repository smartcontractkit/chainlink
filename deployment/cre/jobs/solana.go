package jobs

import (
	"errors"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
)

func verifySolanaJobSpecInputs(inputs job_types.JobSpecInput) error {
	scj := &pkg.StandardCapabilityJob{}
	if err := inputs.UnmarshalTo(scj); err != nil {
		return errors.New("failed to unmarshal job spec input to StandardCapabilityJob: " + err.Error())
	}

	if strings.TrimSpace(scj.Command) == "" {
		return errors.New("command is required and must be a string")
	}

	if strings.TrimSpace(scj.Config) == "" {
		return errors.New("config is required and must be a string")
	}

	return nil
}
