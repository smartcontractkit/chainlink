package jobs

import (
	"errors"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

func verifyAptosJobSpecInputs(inputs job_types.JobSpecInput) error {
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

	if scj.ChainSelectorEVM == 0 {
		return errors.New("chainSelectorEVM is required")
	}

	if scj.ChainSelectorAptos == 0 {
		return errors.New("chainSelectorAptos is required")
	}

	if len(scj.BootstrapPeers) == 0 {
		return errors.New("bootstrapPeers is required")
	}
	if _, err := ocrcommon.ParseBootstrapPeers(scj.BootstrapPeers); err != nil {
		return errors.New("bootstrapPeers is invalid: " + err.Error())
	}

	return nil
}
