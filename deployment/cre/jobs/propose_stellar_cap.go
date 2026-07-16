package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

var _ cldf.ChangeSetV2[ProposeStellarCapJobSpecInput] = ProposeStellarCapJobSpec{}

const stellarNetwork = "stellar"

type StellarOverrideDefaultCfg struct {
	CREForwarderAddress           string        `json:"creForwarderAddress,omitempty" yaml:"creForwarderAddress,omitempty"`
	Network                       string        `json:"network,omitempty" yaml:"network,omitempty"`
	ChainID                       string        `json:"chainId,omitempty" yaml:"chainId,omitempty"`
	ForwarderLookbackLedgers      int64         `json:"forwarderLookbackLedgers,omitempty" yaml:"forwarderLookbackLedgers,omitempty"`
	DeltaStage                    time.Duration `json:"deltaStage,omitempty" yaml:"deltaStage,omitempty"`
	IsLocal                       bool          `json:"isLocal,omitempty" yaml:"isLocal,omitempty"`
	ObservationPollerWorkersCount uint          `json:"observationPollerWorkersCount,omitempty" yaml:"observationPollerWorkersCount,omitempty"`
	ObservationPollPeriod         time.Duration `json:"observationPollPeriod,omitempty" yaml:"observationPollPeriod,omitempty"`
	UnknownRequestsTTL            time.Duration `json:"unknownRequestsTTL,omitempty" yaml:"unknownRequestsTTL,omitempty"`
}

type StellarCapabilityInput struct {
	NodeID             string                    `json:"nodeID" yaml:"nodeID"`
	OverrideDefaultCfg StellarOverrideDefaultCfg `json:"overrideDefaultCfg" yaml:"overrideDefaultCfg"`
}

type ProposeStellarCapJobSpecInput struct {
	Environment string `json:"environment" yaml:"environment"`
	Zone        string `json:"zone" yaml:"zone"`
	Domain      string `json:"domain" yaml:"domain"`
	DONName     string `json:"donName" yaml:"donName"`

	ChainSelector        uint64        `json:"chainSelector" yaml:"chainSelector"`
	BootstrapperOCR3Urls []string      `json:"bootstrapperOCR3Urls" yaml:"bootstrapperOCR3Urls"`
	OCRContractQualifier string        `json:"ocrContractQualifier" yaml:"ocrContractQualifier"`
	OCRChainSelector     uint64        `json:"ocrChainSelector" yaml:"ocrChainSelector"`
	ForwardersQualifier  string        `json:"forwardersContractQualifier" yaml:"forwardersContractQualifier"`
	DeltaStage           time.Duration `json:"deltaStage" yaml:"deltaStage,omitempty"`

	StellarCapabilityInputs []StellarCapabilityInput `json:"stellarCapabilityInputs" yaml:"stellarCapabilityInputs"`
}

type ProposeStellarCapJobSpec struct{}

func (u ProposeStellarCapJobSpec) VerifyPreconditions(e cldf.Environment, input ProposeStellarCapJobSpecInput) error {
	if len(input.StellarCapabilityInputs) == 0 {
		return errors.New("at least one stellar capability input is required")
	}
	for i, capInput := range input.StellarCapabilityInputs {
		if capInput.NodeID == "" {
			return fmt.Errorf("nodeID is required for stellar capability input at index %d", i)
		}
	}

	if err := validateCommonFields(commonCapFields{
		Environment:          input.Environment,
		Domain:               input.Domain,
		Zone:                 input.Zone,
		DONName:              input.DONName,
		ChainSelector:        input.ChainSelector,
		OCRChainSelector:     input.OCRChainSelector,
		BootstrapperOCR3Urls: input.BootstrapperOCR3Urls,
		OCRContractQualifier: input.OCRContractQualifier,
		DeltaStage:           input.DeltaStage,
	}); err != nil {
		return err
	}

	if input.ForwardersQualifier == "" {
		return errors.New("cre forwarder qualifier is required")
	}

	if err := resolveCapRegAddress(e, input.OCRChainSelector, input.OCRContractQualifier); err != nil {
		return err
	}

	family, err := chainselectors.GetSelectorFamily(input.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get family for chain selector %d: %w", input.ChainSelector, err)
	}
	if family != chainselectors.FamilyStellar {
		return fmt.Errorf("chain selector %d belongs to family %q, expected %q", input.ChainSelector, family, chainselectors.FamilyStellar)
	}

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chainID from selector: %w", err)
	}

	for _, capInput := range input.StellarCapabilityInputs {
		ov := capInput.OverrideDefaultCfg
		if ov.ChainID != "" && ov.ChainID != chainIDStr {
			return fmt.Errorf(
				"chainID in override config (%s) does not match chainID from chain selector (%s) for node %s; "+
					"this field is auto-populated and can be omitted",
				ov.ChainID, chainIDStr, capInput.NodeID,
			)
		}
		if err = validateOverrideNetwork(ov.Network, stellarNetwork, capInput.NodeID); err != nil {
			return err
		}
	}

	return nil
}

func (u ProposeStellarCapJobSpec) Apply(e cldf.Environment, input ProposeStellarCapJobSpecInput) (cldf.ChangesetOutput, error) {
	chainName, err := chainselectors.GetChainNameFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain name from selector: %w", err)
	}

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	jobName := fmt.Sprintf("stellar-cap-v2-%s-%s", chainName, input.Zone)
	job := pkg.StandardCapabilityJob{
		JobName:               jobName,
		Command:               "/usr/local/bin/stellar",
		GenerateOracleFactory: true,
		UseCapRegOCRConfig:    true,
		ContractQualifier:     input.OCRContractQualifier,
		CapRegVersion:         stellarCapRegVersion,
		OCRSigningStrategy:    "multi-chain",
		OCRChainSelector:      pkg.ChainSelector(input.OCRChainSelector),
		ChainSelectorEVM:      pkg.ChainSelector(input.OCRChainSelector),
		ChainSelectorStellar:  pkg.ChainSelector(input.ChainSelector),
		BootstrapPeers:        input.BootstrapperOCR3Urls,
	}

	nodeIDToConfig := make(map[string]string, len(input.StellarCapabilityInputs))
	for _, capInput := range input.StellarCapabilityInputs {
		if _, exists := nodeIDToConfig[capInput.NodeID]; exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("duplicate nodeID %q in stellarCapabilityInputs", capInput.NodeID)
		}

		cfg := capInput.OverrideDefaultCfg
		cfg.ChainID = chainIDStr
		cfg.Network = stellarNetwork
		cfg.DeltaStage = input.DeltaStage
		if capInput.OverrideDefaultCfg.DeltaStage != 0 {
			cfg.DeltaStage = capInput.OverrideDefaultCfg.DeltaStage
		}

		enc, err := json.Marshal(cfg)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal stellar cap config: %w", err)
		}

		nodeIDToConfig[capInput.NodeID] = string(enc)
	}

	return proposeAndReport(e, job, nodeIDToConfig, input.Domain, input.Environment, input.DONName, input.Zone)
}
