package job_types

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

type JobSpecInput map[string]any

func (j JobSpecInput) UnmarshalTo(target any) error {
	bytes, err := yaml.Marshal(j)
	if err != nil {
		return fmt.Errorf("failed to marshal job spec input to json: %w", err)
	}

	return yaml.Unmarshal(bytes, target)
}

func (j JobSpecInput) ToStandardCapabilityJob(jobName string, generateOracleFactory bool) (pkg.StandardCapabilityJob, error) {
	out := pkg.StandardCapabilityJob{
		JobName:               jobName,
		GenerateOracleFactory: generateOracleFactory,
	}
	err := j.UnmarshalTo(&out)
	if err != nil {
		return pkg.StandardCapabilityJob{}, fmt.Errorf("failed to unmarshal job spec input to StandardCapabilityJob: %w", err)
	}

	if out.Command == "" {
		return pkg.StandardCapabilityJob{}, errors.New("command is required and must be a string")
	}

	return out, nil
}

func (j JobSpecInput) ToOCR3JobConfigInput() (pkg.OCR3JobConfigInput, error) {
	out := pkg.OCR3JobConfigInput{}
	err := j.UnmarshalTo(&out)
	if err != nil {
		return pkg.OCR3JobConfigInput{}, fmt.Errorf("failed to unmarshal job spec input to OCR3JobConfigInput: %w", err)
	}

	if out.TemplateName == "" || strings.TrimSpace(out.TemplateName) == "" {
		return pkg.OCR3JobConfigInput{}, errors.New("templateName is required and must be a non-empty string")
	}

	if out.ContractQualifier == "" || strings.TrimSpace(out.ContractQualifier) == "" {
		return pkg.OCR3JobConfigInput{}, errors.New("contractQualifier is required and must be a non-empty string")
	}

	if len(out.BootstrapperOCR3Urls) == 0 {
		return pkg.OCR3JobConfigInput{}, errors.New("bootstrapperOCR3Urls is required and cannot be empty")
	}

	return out, nil
}
