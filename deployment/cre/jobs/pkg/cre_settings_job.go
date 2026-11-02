package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg/templates"
)

const (
	templateName = "cre-settings.tmpl"
	// We expect there to only be 1 CRESettings job per node, and they share a fixed UUID for clarity.
	externalJobUUID             = "8561c20c-7d06-421e-a155-3baf21b1622b"
	shardAssignmentJobUUID      = "a3f8c1d2-4e5b-6a7c-8d9e-0f1a2b3c4d5e"
	capabilitiesRegistryJobUUID = "c4a1e8f7-2b3d-4e5f-9a6b-7c8d9e0f1a2b"
)

type CRESettingsJob struct {
	Settings string `yaml:"settings"` // toml
}

// SettingsWithConfigType returns settings prefixed with the explicit settings config_type.
func SettingsWithConfigType(settings string) string {
	return "config_type = \"settings\"\n" + settings
}

func (j CRESettingsJob) ResolveJob() (string, error) {
	settings := SettingsWithConfigType(j.Settings)
	t, err := template.New("s").ParseFS(templates.FS, templateName)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", templateName, err)
	}

	shaSum := sha256.Sum256([]byte(settings))
	data := map[string]any{
		"ExternalJobID": externalJobUUID,
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

// CapabilitiesRegistryJob delivers the offchain capabilities registry payload (see
// core/capabilities/globalconfig and core/services/cresettings). OffchainConfig is the
// proto-JSON body (capabilitiespb.OffchainCapabilitiesRegistry), without a "version" field:
// the node enforces that it be strictly increasing, so callers must set it to a value greater
// than whatever is currently applied.
type CapabilitiesRegistryJob struct {
	OffchainConfig string `yaml:"offchain_config"` // proto-JSON
}

func (j CapabilitiesRegistryJob) ResolveJob() (string, error) {
	// cre-settings.tmpl wraps Settings in a '''...''' TOML multi-line literal string. That
	// scanner has no concept of "nested TOML": it ends Settings at the first literal "'''" it
	// finds anywhere in the raw text, even one that's part of a properly-escaped inner string
	// (e.g. a JSON string value that happens to contain "'''"). Encoding offchain_config through
	// a TOML encoder (below) only protects against OffchainConfig's own quoting/escaping; it
	// can't protect against this outer collision, so guard it explicitly instead of producing
	// TOML that fails to parse later, far from this root cause.
	if strings.Contains(j.OffchainConfig, "'''") {
		return "", errors.New("offchain config must not contain \"'''\": it would prematurely end the settings TOML literal string")
	}

	body, err := toml.Marshal(struct {
		ConfigType     string `toml:"config_type"`
		OffchainConfig string `toml:"offchain_config"`
	}{
		ConfigType:     "capabilities_registry",
		OffchainConfig: j.OffchainConfig,
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode capabilities registry settings: %w", err)
	}
	settings := string(body)

	t, err := template.New("s").ParseFS(templates.FS, templateName)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", templateName, err)
	}

	shaSum := sha256.Sum256([]byte(settings))
	data := map[string]any{
		"ExternalJobID": capabilitiesRegistryJobUUID,
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
