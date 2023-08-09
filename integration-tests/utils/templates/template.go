package templates

import (
	"bytes"
	"errors"
	"text/template"
)

var (
	ErrParsingTemplate = errors.New("failed to parse Go text template")
)

// MarshalTemplate Helper to marshal templates
func MarshalTemplate(jobSpec interface{}, name, templateString string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", errors.Join(err, ErrParsingTemplate)
	}
	err = tmpl.Execute(&buf, jobSpec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
