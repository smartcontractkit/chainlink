package job_types

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/v2/core/config/parse"
)

type JobSpecInput map[string]interface{}

func (j JobSpecInput) ToStandardCapabilityJob(jobName string) (pkg.StandardCapabilityJob, error) {
	cmd, ok := j["command"].(string)
	if !ok || cmd == "" {
		return pkg.StandardCapabilityJob{}, errors.New("command is required and must be a string")
	}

	// config is optional; only validate type if provided.
	var config string
	if rawCfg, exists := j["config"]; exists {
		castCfg, ok := rawCfg.(string)
		if !ok {
			return pkg.StandardCapabilityJob{}, errors.New("config must be a string")
		}
		if castCfg == "" {
			return pkg.StandardCapabilityJob{}, errors.New("config cannot be an empty string")
		}
		config = castCfg
	}

	// externalJobID is optional; only validate type if provided.
	var externalJobID string
	if rawEJID, exists := j["externalJobID"]; exists {
		castEJID, ok := rawEJID.(string)
		if !ok {
			return pkg.StandardCapabilityJob{}, errors.New("externalJobID must be a string")
		}
		if castEJID == "" {
			return pkg.StandardCapabilityJob{}, errors.New("externalJobID cannot be an empty string")
		}
		externalJobID = castEJID
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

func (j JobSpecInput) ToOCR3BootstrapJobInput() (pkg.BootstrapJobInput, error) {
	qualifier, ok := j["contract_qualifier"].(string)
	if !ok || qualifier == "" {
		return pkg.BootstrapJobInput{}, errors.New("contract_qualifier is required and must be a string")
	}

	chainSelector, ok := j["chain_selector"].(string)
	if !ok {
		return pkg.BootstrapJobInput{}, errors.New("chain_selector is required and must be a string")
	}

	chainSel, err := parse.Uint64(chainSelector)
	if err != nil {
		return pkg.BootstrapJobInput{}, fmt.Errorf("failed to parse chain_selector: %w", err)
	}

	return pkg.BootstrapJobInput{
		ContractQualifier: qualifier,
		ChainSelector:     chainSel,
	}, nil
}
