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

var _ cldf.ChangeSetV2[ProposeAptosCapJobSpecInput] = ProposeAptosCapJobSpec{}

const aptosNetwork = "aptos"

type AptosOverrideDefaultCfg struct {
	CREForwarderAddress           string            `json:"creForwarderAddress,omitempty" yaml:"creForwarderAddress,omitempty"`
	Network                       string            `json:"network,omitempty" yaml:"network,omitempty"`
	ChainID                       string            `json:"chainId,omitempty" yaml:"chainId,omitempty"`
	ObservationPollerWorkersCount uint              `json:"observationPollerWorkersCount,omitempty" yaml:"observationPollerWorkersCount,omitempty"`
	ObservationPollPeriod         time.Duration     `json:"observationPollPeriod,omitempty" yaml:"observationPollPeriod,omitempty"`
	ChainHeightPollPeriod         time.Duration     `json:"chainHeightPollPeriod,omitempty" yaml:"chainHeightPollPeriod,omitempty"`
	UnknownRequestsTTL            time.Duration     `json:"unknownRequestsTTL,omitempty" yaml:"unknownRequestsTTL,omitempty"`
	DeltaStage                    time.Duration     `json:"deltaStage" yaml:"deltaStage,omitempty"`
	TxSearchStartingBuffer        time.Duration     `json:"txSearchStartingBuffer" yaml:"txSearchStartingBuffer,omitempty"`
	P2PToTransmitterMap           map[string]string `json:"p2pToTransmitterMap,omitempty" yaml:"p2pToTransmitterMap,omitempty"`
}

type AptosCapabilityInput struct {
	NodeID             string                  `json:"nodeID" yaml:"nodeID"`
	OverrideDefaultCfg AptosOverrideDefaultCfg `json:"overrideDefaultCfg" yaml:"overrideDefaultCfg"`
}

type ProposeAptosCapJobSpecInput struct {
	Environment string `json:"environment" yaml:"environment"`
	Zone        string `json:"zone" yaml:"zone"`
	Domain      string `json:"domain" yaml:"domain"`
	DONName     string `json:"donName" yaml:"donName"`

	ChainSelector        uint64   `json:"chainSelector" yaml:"chainSelector"`
	BootstrapperOCR3Urls []string `json:"bootstrapperOCR3Urls" yaml:"bootstrapperOCR3Urls"`
	OCRContractQualifier string   `json:"ocrContractQualifier" yaml:"ocrContractQualifier"`
	OCRChainSelector     uint64   `json:"ocrChainSelector" yaml:"ocrChainSelector"`

	DeltaStage             time.Duration          `json:"deltaStage" yaml:"deltaStage,omitempty"`
	TxSearchStartingBuffer time.Duration          `json:"txSearchStartingBuffer" yaml:"txSearchStartingBuffer,omitempty"`
	CREForwarderAddress    string                 `json:"creForwarderAddress" yaml:"creForwarderAddress,omitempty"`
	P2PToTransmitterMap    map[string]string      `json:"p2pToTransmitterMap,omitempty" yaml:"p2pToTransmitterMap,omitempty"`
	AptosCapabilityInputs  []AptosCapabilityInput `json:"aptosCapabilityInputs" yaml:"aptosCapabilityInputs"`
}

type ProposeAptosCapJobSpec struct{}

func (u ProposeAptosCapJobSpec) VerifyPreconditions(e cldf.Environment, input ProposeAptosCapJobSpecInput) error {
	if len(input.AptosCapabilityInputs) == 0 {
		return errors.New("at least one aptos capability input is required")
	}
	for i, aptosCapInput := range input.AptosCapabilityInputs {
		if aptosCapInput.NodeID == "" {
			return fmt.Errorf("nodeID is required for aptos capability input at index %d", i)
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

	// PLEX-2797: accept forwarder address directly instead of deriving from datastore,
	// since Aptos forwarder addresses are not yet managed in the catalog.
	if input.CREForwarderAddress == "" {
		return errors.New("cre forwarder address is required")
	}

	family, err := chainselectors.GetSelectorFamily(input.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get family for chain selector %d: %w", input.ChainSelector, err)
	}
	if family != chainselectors.FamilyAptos {
		return fmt.Errorf("chain selector %d belongs to family %q, expected %q", input.ChainSelector, family, chainselectors.FamilyAptos)
	}

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chainID from selector: %w", err)
	}

	ocrAddrRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.OCRChainSelector, input.OCRContractQualifier)
	if _, err := e.DataStore.Addresses().Get(ocrAddrRefKey); err != nil {
		return fmt.Errorf("failed to get OCR contract address for ref key %s: %w", ocrAddrRefKey, err)
	}

	for _, aptosCapInput := range input.AptosCapabilityInputs {
		ov := aptosCapInput.OverrideDefaultCfg

		if ov.ChainID != "" && ov.ChainID != chainIDStr {
			return fmt.Errorf(
				"chainID in override config (%s) does not match chainID from chain selector (%s) for node %s; "+
					"this field is auto-populated and can be omitted",
				ov.ChainID, chainIDStr, aptosCapInput.NodeID,
			)
		}

		if err := validateOverrideNetwork(ov.Network, aptosNetwork, aptosCapInput.NodeID); err != nil {
			return err
		}
	}

	return nil
}

func (u ProposeAptosCapJobSpec) Apply(e cldf.Environment, input ProposeAptosCapJobSpecInput) (cldf.ChangesetOutput, error) {
	chainName, err := chainselectors.GetChainNameFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain name from selector: %w", err)
	}

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	jobName := fmt.Sprintf("aptos-cap-v2-%s-%s", chainName, input.Zone)
	job := pkg.StandardCapabilityJob{
		JobName:               jobName,
		Command:               "/usr/local/bin/aptos",
		GenerateOracleFactory: true,
		ContractQualifier:     input.OCRContractQualifier,
		OCRSigningStrategy:    "single-chain",
		OCRChainSelector:      pkg.ChainSelector(input.OCRChainSelector),
		ChainSelectorEVM:      pkg.ChainSelector(input.OCRChainSelector),
		ChainSelectorAptos:    pkg.ChainSelector(input.ChainSelector),
		BootstrapPeers:        input.BootstrapperOCR3Urls,
	}

	nodeIDToConfig := make(map[string]string, len(input.AptosCapabilityInputs))
	for _, aptosCapInput := range input.AptosCapabilityInputs {
		if _, exists := nodeIDToConfig[aptosCapInput.NodeID]; exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("duplicate nodeID %q in aptosCapabilityInputs", aptosCapInput.NodeID)
		}

		cfg := aptosCapInput.OverrideDefaultCfg
		cfg.ChainID = chainIDStr
		cfg.Network = aptosNetwork
		cfg.CREForwarderAddress = input.CREForwarderAddress // PLEX-2797
		cfg.P2PToTransmitterMap = input.P2PToTransmitterMap
		cfg.DeltaStage = input.DeltaStage
		cfg.TxSearchStartingBuffer = input.TxSearchStartingBuffer
		enc, err := json.Marshal(cfg)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal aptos cap config: %w", err)
		}

		nodeIDToConfig[aptosCapInput.NodeID] = string(enc)
	}

	return proposeAndReport(e, job, nodeIDToConfig, input.Domain, input.DONName, input.Zone)
}
