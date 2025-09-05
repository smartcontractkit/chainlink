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
	if !ok || config == "" {
		return pkg.StandardCapabilityJob{}, errors.New("config is required and must be a string")
	}

	externalJobID, ok := j["externalJobID"].(string)
	if !ok || externalJobID == "" {
		return pkg.StandardCapabilityJob{}, errors.New("externalJobID is required and must be a string")
	}

	oracleFactory, ok := j["oracleFactory"].(pkg.OracleFactory)
	if !ok {
		return pkg.StandardCapabilityJob{}, errors.New("oracleFactory is required and must be of type OracleFactory")
	}

	return pkg.StandardCapabilityJob{
		JobName:       jobName,
		Command:       cmd,
		Config:        config,
		ExternalJobID: externalJobID,
		OracleFactory: oracleFactory,
	}, nil
}
