package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"text/template"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg/templates"
)

const (
	templateName = "cre-settings.tmpl"
	// We expect there to only be 1 CRESettings job per node, and they share a fixed UUID for clarity.
	externalJobUUID = "8561c20c-7d06-421e-a155-3baf21b1622b"
	shardAssignmentJobUUID = "a3f8c1d2-4e5b-6a7c-8d9e-0f1a2b3c4d5e"
)

type CRESettingsJob struct {
	Settings string `yaml:"settings"` // toml
}

func (j CRESettingsJob) ResolveJob() (string, error) {
	t, err := template.New("s").ParseFS(templates.FS, templateName)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", templateName, err)
	}

	shaSum := sha256.Sum256([]byte(j.Settings))
	data := map[string]any{
		"ExternalJobID": externalJobUUID,
		"Hash":          hex.EncodeToString(shaSum[:]),
		"Settings":      j.Settings,
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, templateName, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return b.String(), nil
}

type ShardAssignmentJob struct {
	ShardAssignment string `yaml:"shard_assignment"` // toml
}

func (j ShardAssignmentJob) ResolveJob() (string, error) {
	settings := "config_type = \"shard_assignment\"\n" + j.ShardAssignment
	t, err := template.New("s").ParseFS(templates.FS, templateName)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", templateName, err)
	}

	shaSum := sha256.Sum256([]byte(settings))
	data := map[string]any{
		"ExternalJobID": shardAssignmentJobUUID,
		"Hash":          hex.EncodeToString(shaSum[:]),
		"Settings":      settings,
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, templateName, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return b.String(), nil
}
