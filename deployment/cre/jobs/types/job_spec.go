package job_types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/v2/core/config/parse"
)

type JobSpecInput map[string]interface{}

func (j JobSpecInput) ToStandardCapabilityJob(jobName string) (pkg.StandardCapabilityJob, error) {
	cmd, ok := j["command"].(string)
	if !ok || cmd == "" {
		return pkg.StandardCapabilityJob{}, errors.New("command is required and must be a string")
	}

	// config is optional; only validate type if provided.
	var config string
	if rawCfg, exists := j["config"]; exists {
		castCfg, ok := rawCfg.(string)
		if !ok {
			return pkg.StandardCapabilityJob{}, errors.New("config must be a string")
		}
		if castCfg == "" {
			return pkg.StandardCapabilityJob{}, errors.New("config cannot be an empty string")
		}
		config = castCfg
	}

	// externalJobID is optional; only validate type if provided.
	var externalJobID string
	if rawEJID, exists := j["externalJobID"]; exists {
		castEJID, ok := rawEJID.(string)
		if !ok {
			return pkg.StandardCapabilityJob{}, errors.New("externalJobID must be a string")
		}
		if castEJID == "" {
			return pkg.StandardCapabilityJob{}, errors.New("externalJobID cannot be an empty string")
		}
		externalJobID = castEJID
	}

	// oracleFactory is optional; only validate type if provided.
	var oracleFactory pkg.OracleFactory
	if rawOF, exists := j["oracleFactory"]; exists {
		castOF, ok := rawOF.(pkg.OracleFactory)
		if !ok {
			return pkg.StandardCapabilityJob{}, errors.New("oracleFactory must be of type OracleFactory")
		}
		oracleFactory = castOF
	}

	return pkg.StandardCapabilityJob{
		JobName:       jobName,
		Command:       cmd,
		Config:        config,
		ExternalJobID: externalJobID,
		OracleFactory: oracleFactory,
	}, nil
}

func (j JobSpecInput) ToOCR3BootstrapJobInput() (pkg.BootstrapJobInput, error) {
	qualifier, ok := j["contract_qualifier"].(string)
	if !ok || qualifier == "" {
		return pkg.BootstrapJobInput{}, errors.New("contract_qualifier is required and must be a string")
	}

	chainSelector, ok := j["chain_selector"].(string)
	if !ok {
		return pkg.BootstrapJobInput{}, errors.New("chain_selector is required and must be a string")
	}

	chainSel, err := parse.Uint64(chainSelector)
	if err != nil {
		return pkg.BootstrapJobInput{}, fmt.Errorf("failed to parse chain_selector: %w", err)
	}

	return pkg.BootstrapJobInput{
		ContractQualifier: qualifier,
		ChainSelector:     chainSel,
	}, nil
}

func (j JobSpecInput) ToOCR3JobConfigInput() (pkg.OCR3JobConfigInput, error) {
	// Required: template_name
	// TODO: validate all supported templates
	rawTemplate, ok := j["template_name"].(string)
	if !ok || strings.TrimSpace(rawTemplate) == "" {
		return pkg.OCR3JobConfigInput{}, errors.New("template_name is required and must be a non-empty string")
	}

	// Required: contract_qualifier
	rawQualifier, ok := j["contract_qualifier"].(string)
	if !ok || strings.TrimSpace(rawQualifier) == "" {
		return pkg.OCR3JobConfigInput{}, errors.New("contract_qualifier is required and must be a non-empty string")
	}

	// Required: chain_selector_evm (as string)
	rawEVM, ok := j["chain_selector_evm"].(string)
	if !ok {
		return pkg.OCR3JobConfigInput{}, errors.New("chain_selector_evm is required and must be a string")
	}
	evnSel, err := parse.Uint64(rawEVM)
	if err != nil {
		return pkg.OCR3JobConfigInput{}, fmt.Errorf("failed to parse chain_selector_evm: %w", err)
	}

	// Required: chain_selector_aptos (as string)
	rawAptos, ok := j["chain_selector_aptos"].(string)
	if !ok {
		return pkg.OCR3JobConfigInput{}, errors.New("chain_selector_aptos is required and must be a string")
	}
	aptSel, err := parse.Uint64(rawAptos)
	if err != nil {
		return pkg.OCR3JobConfigInput{}, fmt.Errorf("failed to parse chain_selector_aptos: %w", err)
	}

	// Required: bootstrapper_ocr3_urls (slice of strings)
	rawURLs, exists := j["bootstrapper_ocr3_urls"]
	if !exists {
		return pkg.OCR3JobConfigInput{}, errors.New("bootstrapper_ocr3_urls is required")
	}
	urls, err := toStringSlice(rawURLs)
	if err != nil {
		return pkg.OCR3JobConfigInput{}, fmt.Errorf("bootstrapper_ocr3_urls must be an array of strings: %w", err)
	}
	if len(urls) == 0 {
		return pkg.OCR3JobConfigInput{}, errors.New("bootstrapper_ocr3_urls cannot be empty")
	}

	// Optional: master_public_key
	var masterPub string
	if v, ok := j["master_public_key"]; ok {
		mpk, ok := v.(string)
		if !ok {
			return pkg.OCR3JobConfigInput{}, errors.New("master_public_key must be a string")
		}
		masterPub = mpk
	}

	// Optional: encrypted_private_key_share
	var encShare string
	if v, ok := j["encrypted_private_key_share"]; ok {
		eps, ok := v.(string)
		if !ok {
			return pkg.OCR3JobConfigInput{}, errors.New("encrypted_private_key_share must be a string")
		}
		encShare = eps
	}

	return pkg.OCR3JobConfigInput{
		TemplateName:             strings.TrimSpace(rawTemplate),
		ContractQualifier:        strings.TrimSpace(rawQualifier),
		ChainSelectorEVM:         evnSel,
		ChainSelectorAptos:       aptSel,
		BootstrapperOCR3Urls:     urls,
		MasterPublicKey:          masterPub,
		EncryptedPrivateKeyShare: encShare,
	}, nil
}

// toStringSlice attempts to coerce v into []string, supporting []string and []any with string elements.
func toStringSlice(v any) ([]string, error) {
	switch s := v.(type) {
	case []string:
		return s, nil
	case []any:
		out := make([]string, 0, len(s))
		for i, el := range s {
			str, ok := el.(string)
			if !ok {
				return nil, fmt.Errorf("element %d is %T, expected string", i, el)
			}
			out = append(out, str)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}
