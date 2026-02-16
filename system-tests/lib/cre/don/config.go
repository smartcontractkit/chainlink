package don

import (
	"maps"
	"regexp"
	"strings"

	"dario.cat/mergo"
	"github.com/cockroachdb/errors"
)

const (
	OCRPeeringPort          = 5001
	CapabilitiesPeeringPort = 6690
	GatewayIncomingPort     = 5002
	GatewayOutgoingPort     = 5003
)

// ApplyRuntimeValues fills in any missing config values with runtime-generated values
func ApplyRuntimeValues(userConfig map[string]any, runtimeValues map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	maps.Copy(result, userConfig)

	// Merge runtime fallbacks without overriding existing user values
	// By default, mergo.Merge won't override existing keys (no WithOverride flag)
	err := mergo.Merge(&result, runtimeValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to merge runtime values")
	}

	return result, nil
}

// ValidateTemplateSubstitution checks that all template variables have been properly substituted.
// Returns an error if any unsubstituted patterns are found:
//   - {{.Variable}} - missing field in templateData
//   - <nil> - field exists but has nil value
//
// Usage:
//
//	configStr := configBuffer.String()
//	if err := ValidateTemplateSubstitution(configStr, "capability-name"); err != nil {
//	    return nil, errors.Wrap(err, "template validation failed")
//	}
func ValidateTemplateSubstitution(rendered string, templateName string) error {
	var problems []string

	// Check for unsubstituted template variables like {{.Variable}}
	templateVarRegex := regexp.MustCompile(`\{\{\s*\.[A-Za-z_][A-Za-z0-9_]*\s*\}\}`)
	if matches := templateVarRegex.FindAllString(rendered, -1); len(matches) > 0 {
		problems = append(problems, "unsubstituted variables: "+strings.Join(matches, ", "))
	}

	// Check for nil values rendered as <nil>
	nilRegex := regexp.MustCompile(`<nil>`)
	if nilCount := len(nilRegex.FindAllString(rendered, -1)); nilCount > 0 {
		problems = append(problems, "nil values found (check templateData for missing or nil fields)")
	}

	if len(problems) > 0 {
		return errors.Errorf("template '%s' has issues: %s", templateName, strings.Join(problems, "; "))
	}

	return nil
}
