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
	bytes, err := yaml.Marshal(convertJSONNumbers(map[string]any(j)))
	if err != nil {
		return fmt.Errorf("failed to marshal job spec input to json: %w", err)
	}

	return yaml.Unmarshal(bytes, target)
}

// convertJSONNumbers recursively converts json.Number values to native Go
// numeric types (int64 or float64). This prevents yaml.Marshal from emitting
// them as quoted strings which would then fail to unmarshal into numeric struct
// fields.
func convertJSONNumbers(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, elem := range val {
			out[k] = convertJSONNumbers(elem)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, elem := range val {
			out[i] = convertJSONNumbers(elem)
		}
		return out
	case json.Number:
		if n, err := val.Int64(); err == nil {
			return n
		}
		if f, err := val.Float64(); err == nil {
			return f
		}
		return val.String()
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

func (j JobSpecInput) ToStandardCapabilityJob(jobName string) (pkg.StandardCapabilityJob, error) {
	out := pkg.StandardCapabilityJob{
		JobName: jobName,
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
