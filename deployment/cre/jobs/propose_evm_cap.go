package jobs

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"gopkg.in/yaml.v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

var _ cldf.ChangeSetV2[ProposeEVMCapJobSpecInput] = ProposeEVMCapJobSpec{}

type OverrideDefaultCfg struct {
	ChainID                         uint64        `json:"chainID" yaml:"chainID"`
	Network                         string        `json:"network" yaml:"network"`
	LogTriggerPollInterval          time.Duration `json:"logTriggerPollInterval" yaml:"logTriggerPollInterval"`
	LogTriggerSendChannelBufferSize uint64        `json:"logTriggerSendChannelBufferSize" yaml:"logTriggerSendChannelBufferSize"`
	LogTriggerLimitQueryLogSize     uint64        `json:"logTriggerLimitQueryLogSize" yaml:"logTriggerLimitQueryLogSize"`
	BlockDepth                      int64         `json:"blockDepth" yaml:"blockDepth"`
	CREForwarderAddress             string        `json:"creForwarderAddress" yaml:"creForwarderAddress"`
	// The minimum amount of gas that the receiver contract must get to process the forwarder report.
	// This is the default value used when the user doesn't specify a gas limit when invoking WriteReport.
	ReceiverGasMinimum            uint64        `json:"receiverGasMinimum" yaml:"receiverGasMinimum"`
	NodeAddress                   string        `json:"nodeAddress" yaml:"nodeAddress"`
	ObservationPollerWorkersCount uint          `json:"observationPollerWorkersCount" yaml:"observationPollerWorkersCount"`
	ObservationPollPeriod         time.Duration `json:"observationPollPeriod" yaml:"observationPollPeriod"`
	ChainHeightPollPeriod         time.Duration `json:"chainHeightPollPeriod" yaml:"chainHeightPollPeriod"`
	UnknownRequestsTTL            time.Duration `json:"unknownRequestsTTL" yaml:"unknownRequestsTTL"`
}

type EVMCapabilityInput struct {
	NodeID             string             `json:"nodeID" yaml:"nodeID"`
	OverrideDefaultCfg OverrideDefaultCfg `json:"overrideDefaultCfg" yaml:"overrideDefaultCfg"`
}

type RequiredDonInput struct {
	Environment          string   `json:"environment" yaml:"environment"`
	Zone                 string   `json:"zone" yaml:"zone"`
	Domain               string   `json:"domain" yaml:"domain"`
	DONName              string   `json:"donName" yaml:"donName"`
	ChainSelector        uint64   `json:"chainSelector" yaml:"chainSelector"`
	BootstrapperOCR3Urls []string `json:"bootstrapperOCR3Urls" yaml:"bootstrapperOCR3Urls"`
	OCRContractQualifier string   `json:"ocrContractQualifier" yaml:"ocrContractQualifier"`
	ForwardersQualifier  string   `json:"forwardersContractQualifier" yaml:"forwardersContractQualifier"`
}
type ProposeEVMCapJobSpecInput struct {
	RequiredDonInput    RequiredDonInput     `json:"requiredDonInput" yaml:"requiredDonInput"`
	EVMCapabilityInputs []EVMCapabilityInput `json:"evmCapabilityInputs" yaml:"evmCapabilityInputs"`
}

type ProposeEVMCapJobSpec struct{}

func (u ProposeEVMCapJobSpec) VerifyPreconditions(e cldf.Environment, input ProposeEVMCapJobSpecInput) error {
	requiredDonInput := input.RequiredDonInput
	if requiredDonInput.Environment == "" {
		return errors.New("environment is required")
	}
	if requiredDonInput.Domain == "" {
		return errors.New("domain is required")
	}
	if requiredDonInput.DONName == "" {
		return errors.New("donName is required for node")
	}
	if requiredDonInput.ChainSelector == 0 {
		return errors.New("chain selector is required")
	}
	if input.EVMCapabilityInputs == nil || len(input.EVMCapabilityInputs) == 0 {
		return errors.New("at least one evm capability input is required")
	}
	if requiredDonInput.BootstrapperOCR3Urls == nil || len(requiredDonInput.BootstrapperOCR3Urls) == 0 {
		return errors.New("at least one bootstrapper OCR3 URL is required")
	}
	if requiredDonInput.OCRContractQualifier == "" {
		return errors.New("ocr contract qualifier is required")
	}
	if requiredDonInput.ForwardersQualifier == "" {
		return errors.New("cre forwarder qualifier is required")
	}

	chainID, err := chainselectors.GetChainIDFromSelector(requiredDonInput.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chainID from selector: %w", err)
	}

	ocrAddrRefKey := pkg.GetOCR3CapabilityAddressRefKey(requiredDonInput.ChainSelector, requiredDonInput.OCRContractQualifier)
	_, err = e.DataStore.Addresses().Get(ocrAddrRefKey)
	if err != nil {
		return fmt.Errorf("failed to get OCR contract address for ref key %s: %w", ocrAddrRefKey, err)
	}

	fwdAddrRefKey := pkg.GetKeystoneForwarderCapabilityAddressRefKey(requiredDonInput.ChainSelector, requiredDonInput.ForwardersQualifier)
	fwdAddress, err := e.DataStore.Addresses().Get(fwdAddrRefKey)
	if err != nil {
		return fmt.Errorf("failed to get CRE forwarder address for chain selector %d and qualifier %s: %w", requiredDonInput.ChainSelector, requiredDonInput.ForwardersQualifier, err)
	}

	for _, evmCapInput := range input.EVMCapabilityInputs {
		overrideDefaultCfg := evmCapInput.OverrideDefaultCfg
		if evmCapInput.NodeID == "" {
			return fmt.Errorf("nodeID in evm cap config is required")
		}
		if chainID != strconv.FormatUint(overrideDefaultCfg.ChainID, 10) {
			return fmt.Errorf("chainID in override config (%d) does not match chainID from chain selector (%s) for node %s, this field gets auto generated from the qualifier and doesn't need to be filled out", overrideDefaultCfg.ChainID, chainID, evmCapInput.NodeID)
		}
		if overrideDefaultCfg.Network != "evm" {
			return fmt.Errorf("network in override config (%s) must be 'evm' for node %s, this field is auto generated and doesn't need to be filled out", overrideDefaultCfg.Network, evmCapInput.NodeID)
		}
		if fwdAddress.Address != overrideDefaultCfg.CREForwarderAddress {
			return fmt.Errorf("CRE forwarder address in override config (%s) does not match address from data store (%s) for node %s, this field gets auto generated from the qualifier and doesn't need to be filled out", overrideDefaultCfg.CREForwarderAddress, fwdAddress.Address, evmCapInput.NodeID)
		}
	}

	return nil
}

func (u ProposeEVMCapJobSpec) Apply(e cldf.Environment, input ProposeEVMCapJobSpecInput) (cldf.ChangesetOutput, error) {
	requiredDonInput := input.RequiredDonInput

	chainName, err := chainselectors.GetChainNameFromSelector(requiredDonInput.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain name from selector: %w", err)
	}

	networkEnv, err := chainselectors.ExtractNetworkEnvName(chainName)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to extract network env from chain name: %w", err)
	}

	jobName := fmt.Sprintf("evm-capabilities-v2-%s-%s-%s", chainName, networkEnv, requiredDonInput.Zone)

	job := pkg.StandardCapabilityJob{
		JobName: jobName,
		Command: "/usr/local/bin/evm",
		// OCR
		GenerateOracleFactory: true,
		ContractQualifier:     requiredDonInput.OCRContractQualifier,
		ChainSelectorEVM:      pkg.ChainSelector(requiredDonInput.ChainSelector),
		BootstrapPeers:        requiredDonInput.BootstrapperOCR3Urls,
	}

	fwdAddrRefKey := pkg.GetKeystoneForwarderCapabilityAddressRefKey(requiredDonInput.ChainSelector, requiredDonInput.ForwardersQualifier)
	fwdAddress, err := e.DataStore.Addresses().Get(fwdAddrRefKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get CRE forwarder address for chain selector %d and qualifier %s: %w", requiredDonInput.ChainSelector, requiredDonInput.ForwardersQualifier, err)
	}

	nodeIDToConfig := make(map[string]string)
	for _, evmCapInput := range input.EVMCapabilityInputs {
		capConfig := evmCapInput.OverrideDefaultCfg
		capConfig.Network = "evm"

		// set EVM Config defaults
		capConfig.CREForwarderAddress = fwdAddress.Address
		if capConfig.LogTriggerPollInterval == 0 {
			capConfig.LogTriggerPollInterval = 1500 * time.Millisecond
		}
		if capConfig.ReceiverGasMinimum == 0 {
			capConfig.ReceiverGasMinimum = 500
		}
		if capConfig.LogTriggerSendChannelBufferSize == 0 {
			capConfig.LogTriggerSendChannelBufferSize = 3000
		}

		evmCapConfig, err := yaml.Marshal(capConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal evm cap config: %w", err)
		}

		nodeIDToConfig[evmCapInput.NodeID] = string(evmCapConfig)
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		operations2.ProposeStandardCapabilityJob,
		operations2.ProposeStandardCapabilityJobDeps{Env: e},
		operations2.ProposeStandardCapabilityJobInput{
			Job:            job,
			NodeIDToConfig: nodeIDToConfig,
			Domain:         requiredDonInput.Domain,
			DONName:        requiredDonInput.DONName,
			DONFilters: []offchain.TargetDONFilter{
				{
					Key:   "product",
					Value: offchain.ProductLabel,
				},
				{
					Key:   "environment",
					Value: requiredDonInput.Environment,
				},
				{
					Key:   "zone",
					Value: requiredDonInput.Zone,
				},
				// TODO? I think this is resolved in the sequence later
				{
					Key:   requiredDonInput.DONName,
					Value: "",
				},
			},
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose standard capability job: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
