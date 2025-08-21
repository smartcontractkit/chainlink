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

// ValidateTemplateSubstitution checks that all template variables have been properly substituted
// Returns an error if any {{.Variable}} patterns are found in the rendered output
//
// This function helps catch configuration issues early by ensuring all template variables
// have been provided and substituted. Common causes of unsubstituted variables:
// - Missing fields in templateData map
// - Typos in template variable names
// - Missing values in runtimeValues
//
// Usage:
//
//	configStr := configBuffer.String()
//	if err := ValidateTemplateSubstitution(configStr, "capability-name"); err != nil {
//	    return nil, errors.Wrap(err, "template validation failed")
//	}
func ValidateTemplateSubstitution(rendered string, templateName string) error {
	// Regex to find unsubstituted template variables like {{.Variable}}
	templateVarRegex := regexp.MustCompile(`\{\{\s*\.[A-Za-z_][A-Za-z0-9_]*\s*\}\}`)

	matches := templateVarRegex.FindAllString(rendered, -1)
	if len(matches) > 0 {
		return errors.Errorf("template '%s' contains unsubstituted variables: %s",
			templateName, strings.Join(matches, ", "))
	}

	return nil
}
