package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg/templates"
)

const bootstrapPth = "ocr3_bootstrap.tmpl"

type BootstrapJobInput struct {
	ContractQualifier string        `json:"contractQualifier" yaml:"contractQualifier"` // OCR contract address qualifier
	ChainSelector     ChainSelector `json:"chainSelector" yaml:"chainSelector"`
}

type BootstrapCfg struct {
	JobName       string
	ExternalJobID string // If empty, will be generated
	ContractID    string // OCR contract address
	ChainID       string
}

func (cfg BootstrapCfg) Validate() error {
	if cfg.JobName == "" {
		return errors.New("ocr3 bootstrap job name cannot be empty")
	}

	if cfg.ContractID == "" {
		return errors.New("ocr3 bootstrap contract ID cannot be empty")
	}

	if cfg.ChainID == "" {
		return errors.New("ocr3 bootstrap chain ID cannot be empty")
	}

	return nil
}

func (cfg BootstrapCfg) ResolveSpec() (string, error) {
	t, err := template.New("s").ParseFS(templates.FS, bootstrapPth)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", bootstrapPth, err)
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, bootstrapPth, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return b.String(), nil
}

func BootstrapExternalJobID(donName, contractID string, evmChainSel uint64) (string, error) {
	return ExternalJobID(donName+"-bootstrap", contractID, "ocr3_bootstrap", evmChainSel)
}
