package job_types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

type JobSpecInput map[string]any

func (j JobSpecInput) UnmarshalTo(target any) error {
	sanitized := convertJSONNumbers(map[string]any(j))
	bytes, err := yaml.Marshal(sanitized)
	if err != nil {
		return fmt.Errorf("failed to marshal job spec input to yaml: %w", err)
	}
	return yaml.Unmarshal(bytes, target)
}

func convertJSONNumbers(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = convertValue(v)
	}
	return out
}

func convertValue(v any) any {
	switch val := v.(type) {
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return i
		}
		if f, err := val.Float64(); err == nil {
			return f
		}
		return val.String()
	case map[string]any:
		return convertJSONNumbers(val)
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = convertValue(item)
		}
		return out
	default:
		return v
	}
}

func (j JobSpecInput) UnmarshalFrom(source any) error {
	bytes, err := yaml.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal source to json: %w", err)
	}

	return yaml.Unmarshal(bytes, &j)
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

	if out.CapRegVersion == "" && (out.ContractQualifier == "" || strings.TrimSpace(out.ContractQualifier) == "") {
		return pkg.OCR3JobConfigInput{}, errors.New("contractQualifier is required and must be a non-empty string (for legacy OCR3 config contracts)")
	}

	if len(out.BootstrapperOCR3Urls) == 0 {
		return pkg.OCR3JobConfigInput{}, errors.New("bootstrapperOCR3Urls is required and cannot be empty")
	}

	return out, nil
}
