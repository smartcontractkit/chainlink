package types

import (
	"errors"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

type JobSpecInput map[string]interface{}

func (j JobSpecInput) ToStandardCapabilityJob(jobName string) (pkg.StandardCapabilityJob, error) {
	var cmd, config, externalJobID string
	var oracleFactory pkg.OracleFactory

	if c, ok := j["command"].(string); ok {
		cmd = c
	} else {
		return pkg.StandardCapabilityJob{}, errors.New("command is required and must be a string")
	}

	if c, ok := j["config"].(string); ok {
		config = c
	} else {
		return pkg.StandardCapabilityJob{}, errors.New("config is required and must be a string")
	}

	if ej, ok := j["externalJobID"].(string); ok {
		externalJobID = ej
	} else {
		return pkg.StandardCapabilityJob{}, errors.New("externalJobID is required and must be a string")
	}

	if of, ok := j["oracleFactory"].(pkg.OracleFactory); ok {
		oracleFactory = of
	} else {
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
