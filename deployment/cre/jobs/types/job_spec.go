package job_types

import (
	"errors"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

type JobSpecInput map[string]interface{}

func (j JobSpecInput) ToStandardCapabilityJob(jobName string) (pkg.StandardCapabilityJob, error) {
	cmd, ok := j["command"].(string)
	if !ok || cmd == "" {
		return pkg.StandardCapabilityJob{}, errors.New("command is required and must be a string")
	}

	config, ok := j["config"].(string)
	if ok && config == "" {
		return pkg.StandardCapabilityJob{}, errors.New("config cannot be an empty string")
	}

	externalJobID, ok := j["externalJobID"].(string)
	if ok && externalJobID == "" {
		return pkg.StandardCapabilityJob{}, errors.New("externalJobID cannot be an empty string")
	}

	// oracleFactory is optional; only validate type if provided.
	var oracleFactory pkg.OracleFactory
	if rawOF, exists := j["oracleFactory"]; exists {
		castOF, ok := rawOF.(pkg.OracleFactory)
		if !ok {
			return pkg.StandardCapabilityJob{}, errors.New("oracleFactory must be of type OracleFactory")
		}
		oracleFactory = castOF
	}

	return pkg.StandardCapabilityJob{
		JobName:       jobName,
		Command:       cmd,
		Config:        config,
		ExternalJobID: externalJobID,
		OracleFactory: oracleFactory,
	}, nil
}
