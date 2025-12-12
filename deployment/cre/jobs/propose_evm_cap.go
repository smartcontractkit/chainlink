package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

var _ cldf.ChangeSetV2[ProposeEVMCapJobSpecInput] = ProposeEVMCapJobSpec{}

// Config defaults
const (
	logTriggerPollInterval        = 1500000000
	receiverGasMinimum     uint64 = 500
	logTriggerSendChanBuf  uint64 = 3000
	network                       = "evm"
)

type OverrideDefaultCfg struct {
	ChainID                         uint64 `json:"chainId,omitempty" yaml:"chainId,omitempty"`
	Network                         string `json:"network,omitempty" yaml:"network,omitempty"`
	LogTriggerPollInterval          uint64 `json:"logTriggerPollInterval,omitempty" yaml:"logTriggerPollInterval,omitempty"`
	LogTriggerSendChannelBufferSize uint64 `json:"logTriggerSendChannelBufferSize,omitempty" yaml:"logTriggerSendChannelBufferSize,omitempty"`
	LogTriggerLimitQueryLogSize     uint64 `json:"logTriggerLimitQueryLogSize,omitempty" yaml:"logTriggerLimitQueryLogSize,omitempty"`
	// ForwarderLookbackBlocks defines how many blocks back to search for the ReportProcessed event (default 100).
	ForwarderLookbackBlocks int64  `json:"forwarderLookbackBlocks" yaml:"forwarderLookbackBlocks,omitempty"`
	CREForwarderAddress     string `json:"creForwarderAddress,omitempty" yaml:"creForwarderAddress,omitempty"`
	// ReceiverGasMinimum is the minimum amount of gas that the receiver contract must get to process the forwarder report.
	// This is the default value used when the user doesn't specify a gas limit when invoking WriteReport.
	ReceiverGasMinimum            uint64        `json:"receiverGasMinimum,omitempty" yaml:"receiverGasMinimum,omitempty"`
	NodeAddress                   string        `json:"nodeAddress,omitempty" yaml:"nodeAddress,omitempty"`
	ObservationPollerWorkersCount uint          `json:"observationPollerWorkersCount,omitempty" yaml:"observationPollerWorkersCount,omitempty"`
	ObservationPollPeriod         time.Duration `json:"observationPollPeriod,omitempty" yaml:"observationPollPeriod,omitempty"`
	ChainHeightPollPeriod         time.Duration `json:"chainHeightPollPeriod,omitempty" yaml:"chainHeightPollPeriod,omitempty"`
	UnknownRequestsTTL            time.Duration `json:"unknownRequestsTTL,omitempty" yaml:"unknownRequestsTTL,omitempty"`
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

	ChainSelector        uint64   `json:"chainSelector" yaml:"chainSelector"`
	BootstrapperOCR3Urls []string `json:"bootstrapperOCR3Urls" yaml:"bootstrapperOCR3Urls"`
	OCRContractQualifier string   `json:"ocrContractQualifier" yaml:"ocrContractQualifier"`
	OCRChainSelector     uint64   `json:"ocrChainSelector" yaml:"ocrChainSelector"`
	ForwardersQualifier  string   `json:"forwardersContractQualifier" yaml:"forwardersContractQualifier"`
	// ForwarderLookbackBlocks defines how many blocks back to search for the ReportProcessed event (default 100)
	ForwarderLookbackBlocks int64                `json:"forwarderLookbackBlocks" yaml:"forwarderLookbackBlocks,omitempty"`
	EVMCapabilityInputs     []EVMCapabilityInput `json:"evmCapabilityInputs" yaml:"evmCapabilityInputs"`
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
	if input.OCRChainSelector == 0 {
		return errors.New("ocr chain selector is required")
	}
	if len(input.EVMCapabilityInputs) == 0 {
		return errors.New("at least one evm capability input is required")
	}
	for i, evmCapInput := range input.EVMCapabilityInputs {
		if evmCapInput.NodeID == "" {
			return fmt.Errorf("nodeID is required for evm capability input at index %d", i)
		}
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

	ocrAddrRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.OCRChainSelector, input.OCRContractQualifier)
	if _, err := e.DataStore.Addresses().Get(ocrAddrRefKey); err != nil {
		return fmt.Errorf("failed to get OCR contract address for ref key %s: %w", ocrAddrRefKey, err)
	}

	fwdAddrRefKey := pkg.GetKeystoneForwarderCapabilityAddressRefKey(input.ChainSelector, input.ForwardersQualifier)
	fwdAddress, err := e.DataStore.Addresses().Get(fwdAddrRefKey)
	if err != nil {
		return fmt.Errorf("failed to get CRE forwarder address for ref key %q: %w", fwdAddrRefKey, err)
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

		if ov.Network != "" && ov.Network != network {
			return fmt.Errorf("network in override config must be %q if set; got %q for node %s", network, ov.Network, evmCapInput.NodeID)
		}

		if ov.CREForwarderAddress != "" && fwdAddress.Address != ov.CREForwarderAddress {
			return fmt.Errorf(
				"CRE forwarder address in override config (%s) does not match address from data store (%s) for node %s; "+
					"this field is auto-populated and can be omitted",
				ov.CREForwarderAddress, fwdAddress.Address, evmCapInput.NodeID,
			)
		}

		if ov.LogTriggerPollInterval != 0 && ov.LogTriggerPollInterval < logTriggerPollInterval {
			return fmt.Errorf("logTriggerPollInterval (%d) is below minimum (%d) for node %s",
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

	chainIDStr, err := chainselectors.GetChainIDFromSelector(input.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	chainID, err := strconv.ParseUint(chainIDStr, 10, 64)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to parse chain ID %s: %w", chainIDStr, err)
	}

	jobName := fmt.Sprintf("evm-capabilities-v2-%s-%s", chainName, input.Zone)
	job := pkg.StandardCapabilityJob{
		JobName:               jobName,
		Command:               "/usr/local/bin/evm",
		GenerateOracleFactory: true,
		ContractQualifier:     input.OCRContractQualifier,
		OCRSigningStrategy:    "single-chain",
		OCRChainSelector:      pkg.ChainSelector(input.OCRChainSelector),
		ChainSelectorEVM:      pkg.ChainSelector(input.ChainSelector),
		BootstrapPeers:        input.BootstrapperOCR3Urls,
	}

	fwdAddrRefKey := pkg.GetKeystoneForwarderCapabilityAddressRefKey(input.ChainSelector, input.ForwardersQualifier)
	fwdAddress, err := e.DataStore.Addresses().Get(fwdAddrRefKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get CRE forwarder address for chain selector %d and ref key %s: %w", input.ChainSelector, fwdAddrRefKey, err)
	}

	nodeIDToConfig := make(map[string]string, len(input.EVMCapabilityInputs))
	for _, evmCapInput := range input.EVMCapabilityInputs {
		if _, exists := nodeIDToConfig[evmCapInput.NodeID]; exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("duplicate nodeID %q in evmCapabilityInputs", evmCapInput.NodeID)
		}

		nodeInfos, err := deployment.NodeInfo([]string{evmCapInput.NodeID}, e.Offchain)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get node info for node %s: %w", evmCapInput.NodeID, err)
		}
		if len(nodeInfos) == 0 {
			return cldf.ChangesetOutput{}, fmt.Errorf("no node info for node %s", evmCapInput.NodeID)
		}

		evmOCRConfig, ok := nodeInfos[0].OCRConfigForChainSelector(input.ChainSelector)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("no evm ocr config for node %s and chain selector %d", evmCapInput.NodeID, input.ChainSelector)
		}

		cfg := evmCapInput.OverrideDefaultCfg
		cfg.ChainID = chainID
		cfg.Network = network
		cfg.NodeAddress = string(evmOCRConfig.TransmitAccount)
		cfg.CREForwarderAddress = fwdAddress.Address
		if cfg.ForwarderLookbackBlocks == 0 {
			cfg.ForwarderLookbackBlocks = input.ForwarderLookbackBlocks
		}

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

		enc, err := json.Marshal(cfg)
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
				{Key: "zone", Value: input.Zone},
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
