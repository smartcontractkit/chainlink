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

// Config defaults
const (
	logTriggerPollInterval        = 1500 * time.Millisecond
	receiverGasMinimum     uint64 = 500
	logTriggerSendChanBuf  uint64 = 3000
	network                       = "evm"
)

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

type ProposeEVMCapJobSpecInput struct {
	Environment string `json:"environment" yaml:"environment"`
	Zone        string `json:"zone" yaml:"zone"`
	Domain      string `json:"domain" yaml:"domain"`
	DONName     string `json:"donName" yaml:"donName"`

	ChainSelector        uint64               `json:"chainSelector" yaml:"chainSelector"`
	BootstrapperOCR3Urls []string             `json:"bootstrapperOCR3Urls" yaml:"bootstrapperOCR3Urls"`
	OCRContractQualifier string               `json:"ocrContractQualifier" yaml:"ocrContractQualifier"`
	ForwardersQualifier  string               `json:"forwardersContractQualifier" yaml:"forwardersContractQualifier"`
	EVMCapabilityInputs  []EVMCapabilityInput `json:"evmCapabilityInputs" yaml:"evmCapabilityInputs"`
}

type ProposeEVMCapJobSpec struct{}

func (u ProposeEVMCapJobSpec) VerifyPreconditions(e cldf.Environment, input ProposeEVMCapJobSpecInput) error {
	if input.Environment == "" {
		return errors.New("environment is required")
	}
	if input.Domain == "" {
		return errors.New("domain is required")
	}
	if input.Zone == "" {
		return errors.New("zone is required")
	}
	if input.DONName == "" {
		return errors.New("donName is required")
	}
	if input.ChainSelector == 0 {
		return errors.New("chain selector is required")
	}
	if len(input.EVMCapabilityInputs) == 0 {
		return errors.New("at least one evm capability input is required")
	}
	if len(input.BootstrapperOCR3Urls) == 0 {
		return errors.New("at least one bootstrapper OCR3 URL is required")
	}
	for i, u := range input.BootstrapperOCR3Urls {
		if u == "" {
			return fmt.Errorf("bootstrapper OCR3 URL at index %d is empty", i)
		}
	}

	if input.OCRContractQualifier == "" {
		return errors.New("ocr contract qualifier is required")
	}
	if input.ForwardersQualifier == "" {
		return errors.New("cre forwarder qualifier is required")
	}

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chainID from selector: %w", err)
	}

	ocrAddrRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.ChainSelector, input.OCRContractQualifier)
	if _, err := e.DataStore.Addresses().Get(ocrAddrRefKey); err != nil {
		return fmt.Errorf("failed to get OCR contract address for ref key %s: %w", ocrAddrRefKey, err)
	}

	fwdAddrRefKey := pkg.GetKeystoneForwarderCapabilityAddressRefKey(input.ChainSelector, input.ForwardersQualifier)
	fwdAddress, err := e.DataStore.Addresses().Get(fwdAddrRefKey)
	if err != nil {
		return fmt.Errorf("failed to get CRE forwarder address for chain selector %d and qualifier %s: %w", input.ChainSelector, input.ForwardersQualifier, err)
	}

	for _, evmCapInput := range input.EVMCapabilityInputs {
		ov := evmCapInput.OverrideDefaultCfg
		if evmCapInput.NodeID == "" {
			return errors.New("nodeID in evm capability input is required")
		}

		// If user provided ChainID, ensure it matches selector-derived chain id.
		if ov.ChainID != 0 && chainIDStr != strconv.FormatUint(ov.ChainID, 10) {
			return fmt.Errorf(
				"chainID in override config (%d) does not match chainID from chain selector (%s) for node %s; "+
					"this field is auto-populated and can be omitted",
				ov.ChainID, chainIDStr, evmCapInput.NodeID,
			)
		}

		// If user set Network, it must be "evm".
		if ov.Network != "" && ov.Network != network {
			return fmt.Errorf("network in override config must be %q if set; got %q for node %s", network, ov.Network, evmCapInput.NodeID)
		}

		// If user set CREForwarderAddress, ensure it matches computed value.
		if ov.CREForwarderAddress != "" && fwdAddress.Address != ov.CREForwarderAddress {
			return fmt.Errorf(
				"CRE forwarder address in override config (%s) does not match address from data store (%s) for node %s; "+
					"this field is auto-populated and can be omitted",
				ov.CREForwarderAddress, fwdAddress.Address, evmCapInput.NodeID,
			)
		}

		if ov.LogTriggerPollInterval != 0 && ov.LogTriggerPollInterval < logTriggerPollInterval {
			return fmt.Errorf("logTriggerPollInterval (%s) is below minimum (%s) for node %s",
				ov.LogTriggerPollInterval, logTriggerPollInterval, evmCapInput.NodeID)
		}
		if ov.ReceiverGasMinimum != 0 && ov.ReceiverGasMinimum < receiverGasMinimum {
			return fmt.Errorf("receiverGasMinimum (%d) is below minimum (%d) for node %s",
				ov.ReceiverGasMinimum, receiverGasMinimum, evmCapInput.NodeID)
		}
		if ov.LogTriggerSendChannelBufferSize != 0 && ov.LogTriggerSendChannelBufferSize < logTriggerSendChanBuf {
			return fmt.Errorf("logTriggerSendChannelBufferSize (%d) is below minimum (%d) for node %s",
				ov.LogTriggerSendChannelBufferSize, logTriggerSendChanBuf, evmCapInput.NodeID)
		}
	}

	return nil
}

func (u ProposeEVMCapJobSpec) Apply(e cldf.Environment, input ProposeEVMCapJobSpecInput) (cldf.ChangesetOutput, error) {
	chainName, err := chainselectors.GetChainNameFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain name from selector: %w", err)
	}

	networkEnv, err := chainselectors.ExtractNetworkEnvName(chainName)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to extract network env from chain name: %w", err)
	}

	jobName := fmt.Sprintf("evm-capabilities-v2-%s-%s-%s", chainName, networkEnv, input.Zone)

	job := pkg.StandardCapabilityJob{
		JobName:               jobName,
		Command:               "/usr/local/bin/evm",
		GenerateOracleFactory: true,
		ContractQualifier:     input.OCRContractQualifier,
		ChainSelectorEVM:      pkg.ChainSelector(input.ChainSelector),
		BootstrapPeers:        input.BootstrapperOCR3Urls,
	}

	fwdAddrRefKey := pkg.GetKeystoneForwarderCapabilityAddressRefKey(input.ChainSelector, input.ForwardersQualifier)
	fwdAddress, err := e.DataStore.Addresses().Get(fwdAddrRefKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get CRE forwarder address for chain selector %d and qualifier %s: %w", input.ChainSelector, input.ForwardersQualifier, err)
	}

	nodeIDToConfig := make(map[string]string, len(input.EVMCapabilityInputs))
	for _, evmCapInput := range input.EVMCapabilityInputs {
		if _, exists := nodeIDToConfig[evmCapInput.NodeID]; exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("duplicate nodeID %q in evmCapabilityInputs", evmCapInput.NodeID)
		}

		cfg := evmCapInput.OverrideDefaultCfg

		// Canonical values derived from inputs.
		cfg.Network = network
		cfg.CREForwarderAddress = fwdAddress.Address

		// Apply defaults if unset (zero-values). Values provided and validated in VerifyPreconditions are preserved.
		if cfg.LogTriggerPollInterval == 0 {
			cfg.LogTriggerPollInterval = logTriggerPollInterval
		}
		if cfg.ReceiverGasMinimum == 0 {
			cfg.ReceiverGasMinimum = receiverGasMinimum
		}
		if cfg.LogTriggerSendChannelBufferSize == 0 {
			cfg.LogTriggerSendChannelBufferSize = logTriggerSendChanBuf
		}

		enc, err := yaml.Marshal(cfg)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal evm cap config: %w", err)
		}

		nodeIDToConfig[evmCapInput.NodeID] = string(enc)
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		operations2.ProposeStandardCapabilityJob,
		operations2.ProposeStandardCapabilityJobDeps{Env: e},
		operations2.ProposeStandardCapabilityJobInput{
			Job:            job,
			NodeIDToConfig: nodeIDToConfig,
			Domain:         input.Domain,
			DONName:        input.DONName,
			DONFilters: []offchain.TargetDONFilter{
				{Key: "product", Value: offchain.ProductLabel},
				{Key: "environment", Value: input.Environment},
				{Key: "zone", Value: input.Zone},
				// Preserved: sequence may handle this specially.
				{Key: input.DONName, Value: ""},
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
