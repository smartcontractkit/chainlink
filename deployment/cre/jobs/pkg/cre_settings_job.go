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
	data := map[string]interface{}{
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
