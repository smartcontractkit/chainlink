package job_types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type JobSpecTemplate int64

const (
	Cron JobSpecTemplate = iota
	BootstrapOCR3
	OCR3
	HTTPTrigger
	HTTPAction
	EVM
	Gateway
	BootstrapVault
	Consensus
)

func (jt JobSpecTemplate) String() string {
	switch jt {
	case Cron:
		return "cron"
	case BootstrapOCR3:
		return "bootstrap-ocr3"
	case OCR3:
		return "ocr3"
	case HTTPTrigger:
		return "http-trigger"
	case HTTPAction:
		return "http-action"
	case EVM:
		return "evm"
	case Gateway:
		return "gateway"
	case BootstrapVault:
		return "bootstrap-vault"
	case Consensus:
		return "consensus"
	default:
		return "unknown"
	}
}

// parseJobSpecTemplate converts a (case-insensitive) string into a JobSpecTemplate enum value.
func parseJobSpecTemplate(s string) (JobSpecTemplate, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "cron":
		return Cron, nil
	case "bootstrap-ocr3":
		return BootstrapOCR3, nil
	case "ocr3":
		return OCR3, nil
	case "http-trigger":
		return HTTPTrigger, nil
	case "http-action":
		return HTTPAction, nil
	case "evm":
		return EVM, nil
	case "gateway":
		return Gateway, nil
	case "bootstrap-vault":
		return BootstrapVault, nil
	case "consensus":
		return Consensus, nil
	case "", "unknown":
		return 0, errors.New("job spec template cannot be empty")
	default:
		return 0, fmt.Errorf("unsupported job spec template: %s", s)
	}
}

// UnmarshalJSON allows JobSpecTemplate to be provided as a JSON string.
func (jt *JobSpecTemplate) UnmarshalJSON(b []byte) error {
	var raw any
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal JobSpecTemplate: %w", err)
	}
	switch v := raw.(type) {
	case string:
		parsed, err := parseJobSpecTemplate(v)
		if err != nil {
			return err
		}
		*jt = parsed
		return nil
	case float64: // allow legacy numeric enum if ever passed
		iv := int64(v)
		if iv == int64(Cron) {
			*jt = Cron
			return nil
		}
		return fmt.Errorf("unsupported numeric job spec template: %v", v)
	default:
		return fmt.Errorf("expected string for JobSpecTemplate, got %T", raw)
	}
}

// UnmarshalYAML allows JobSpecTemplate to be provided as a YAML string.
func (jt *JobSpecTemplate) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err == nil {
		parsed, pErr := parseJobSpecTemplate(s)
		if pErr != nil {
			return pErr
		}
		*jt = parsed
		return nil
	}
	// fallback: try numeric
	var num int64
	if err := node.Decode(&num); err == nil {
		if JobSpecTemplate(num) == Cron {
			*jt = Cron
			return nil
		}
		return fmt.Errorf("unsupported numeric job spec template: %d", num)
	}
	return fmt.Errorf("failed to decode JobSpecTemplate from YAML: %s", node.Value)
}
