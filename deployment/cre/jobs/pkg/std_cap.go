package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg/templates"
)

const (
	ErrorEmptyJobName = "job name cannot be empty"
)

type StandardCapabilityJob struct {
	JobName string `json:"jobName" yaml:"jobName"` // Must be alphanumeric, with _, -, ., no spaces.
	Command string `json:"command" yaml:"command"`
	Config  string `json:"config" yaml:"config"`

	// If not provided, ExternalJobID is automatically filled in by calling `externalJobIDHashFunc`
	ExternalJobID string `json:"externalJobID" yaml:"externalJobID"`
	// OracleFactory is the configuration for the Oracle Factory job.
	OracleFactory *OracleFactory `json:"oracleFactory" yaml:"oracleFactory"`

	// Additional fields used to drive oracle factory creation/config
	GenerateOracleFactory bool          `json:"generateOracleFactory" yaml:"generateOracleFactory"` // if true, an oracle factory will be generated using the fields below
	OCRSigningStrategy    string        `json:"ocrSigningStrategy" yaml:"ocrSigningStrategy"`       // used to set the signing strategy in the oracle factory
	ContractQualifier     string        `json:"contractQualifier" yaml:"contractQualifier"`         // qualifier for the OCR3 contract or CapabilitiesRegistry (when capRegVersion is set)
	OCRChainSelector      ChainSelector `json:"ocrChainSelector" yaml:"ocrChainSelector"`           // contract chain selector, doesn't have to live on the same chain as the evm selector
	UseCapRegOCRConfig    bool          `json:"useCapRegOCRConfig" yaml:"useCapRegOCRConfig"`       // if true, use CapabilitiesRegistry instead of legacy OCR3 contract for oracle factory config
	CapRegVersion         string        `json:"capRegVersion" yaml:"capRegVersion"`                 // CapabilitiesRegistry contract version (e.g. "2.0.0"); required when useCapRegOCRConfig is true

	ChainSelectorEVM    ChainSelector `json:"chainSelectorEVM" yaml:"chainSelectorEVM"`       // used to fetch OCR EVM configs from nodes
	ChainSelectorAptos  ChainSelector `json:"chainSelectorAptos" yaml:"chainSelectorAptos"`    // used to fetch OCR Aptos configs from nodes - optional
	ChainSelectorSolana ChainSelector `json:"chainSelectorSolana" yaml:"chainSelectorSolana"`  // used to fetch OCR Solana configs from nodes - optional
	BootstrapPeers      []string      `json:"bootstrapPeers" yaml:"bootstrapPeers"`            // set as value in the oracle factory
}

func (s *StandardCapabilityJob) Resolve() (string, error) {
	if s.ExternalJobID == "" {
		// We expect there to only be 1 instance of a standard capability per node
		// This is because adding duplicate capabilities to the registry will typically fail due to an ID clash.
		// Some capabilities, such as contract read and write, are unique per their config
		externalJobID, err := externalJobIDHashFunc([]byte(s.Command), []byte(s.Config))
		if err != nil {
			return "", fmt.Errorf("failed to create external job id: %w", err)
		}
		s.ExternalJobID = externalJobID.String()
	}

	t, err := template.New("s").ParseFS(templates.FS, "stdcap.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse stdcap.tmpl: %w", err)
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, "stdcap.tmpl", s)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return b.String(), nil
}

func (s *StandardCapabilityJob) Validate() error {
	if s.JobName == "" {
		return errors.New(ErrorEmptyJobName)
	}

	if !s.GenerateOracleFactory {
		// If not generating the oracle factory, no further validation is needed
		return nil
	}

	if !s.UseCapRegOCRConfig && s.ContractQualifier == "" {
		return errors.New("contract qualifier cannot be empty")
	}

	if s.UseCapRegOCRConfig && s.CapRegVersion == "" {
		return errors.New("capRegVersion is required when useCapRegOCRConfig is true")
	}

	if s.ChainSelectorEVM == 0 {
		return errors.New("chain selector EVM cannot be zero")
	}

	if len(s.BootstrapPeers) == 0 {
		return errors.New("bootstrap peers cannot be empty")
	}

	return nil
}

func externalJobIDHashFunc(command, config []byte) (uuid.UUID, error) {
	var externalJobID uuid.UUID
	if len(config) > 0 {
		externalJobID = uuid.NewSHA1(uuid.Nil, append(command, config...))
		return externalJobID, nil
	}
	externalJobID = uuid.NewSHA1(uuid.Nil, command)
	return externalJobID, nil
}
