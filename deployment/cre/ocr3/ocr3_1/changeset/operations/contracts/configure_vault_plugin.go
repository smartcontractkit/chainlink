package contracts

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	changesetstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[ConfigureVaultPluginInput] = ConfigureVaultPlugin{}

type InstanceIDComponents struct {
	DKGContractQualifier string `json:"dkgContractQualifier" yaml:"dkgContractQualifier"`
	ConfigDigest         string `json:"configDigest" yaml:"configDigest"`
}

type ConfigureVaultPluginInput struct {
	ContractChainSelector uint64 `json:"contractChainSelector" yaml:"contractChainSelector"`
	ContractQualifier     string `json:"contractQualifier" yaml:"contractQualifier"`

	DON                   contracts.DonNodeSet         `json:"don" yaml:"don"`
	OracleConfig          *ocr3_1.V3_1OracleConfig     `json:"oracleConfig" yaml:"oracleConfig"`
	DryRun                bool                         `json:"dryRun" yaml:"dryRun"`
	InstanceID            InstanceIDComponents         `json:"instanceID" yaml:"instanceID"`
	ReportingPluginConfig *vault.ReportingPluginConfig `json:"reportingPluginConfig,omitempty" yaml:"reportingPluginConfig,omitempty"`

	MCMSConfig *crecontracts.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type ConfigureVaultPlugin struct{}

func (l ConfigureVaultPlugin) VerifyPreconditions(_ cldf.Environment, input ConfigureVaultPluginInput) error {
	if input.ContractChainSelector == 0 {
		return errors.New("contract chain selector is required")
	}
	if input.ContractQualifier == "" {
		return errors.New("contract qualifier is required")
	}
	if input.DON.Name == "" {
		return errors.New("don name is required")
	}
	if len(input.DON.NodeIDs) == 0 {
		return errors.New("at least one don node ID is required")
	}
	if input.OracleConfig == nil {
		return errors.New("oracle config is required")
	}
	_, _, err := ocr3_1.VerifyAndExtractOCR3_1Fields(input.OracleConfig.PrevConfigDigest, input.OracleConfig.PrevSeqNr, input.OracleConfig.PrevHistoryDigest)
	if err != nil {
		return errors.New("verifyAndExtractOCR3_1Fields failed verification: " + err.Error())
	}

	if input.InstanceID.DKGContractQualifier == "" {
		return errors.New("instanceID.dkgContractQualifier is required")
	}
	if input.InstanceID.ConfigDigest == "" {
		return errors.New("instanceID.config_digest is required")
	}
	_, err = ocr3_1.HexStringTo32ByteArray(input.InstanceID.ConfigDigest)
	if err != nil {
		return fmt.Errorf("invalid instanceID.configDigest, should be hex encoded 32 bytes: %w", err)
	}
	return nil
}

func (l ConfigureVaultPlugin) Apply(e cldf.Environment, input ConfigureVaultPluginInput) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Configuring VaultPlugin contract with DON", "donName", input.DON.Name, "nodes", input.DON.NodeIDs, "dryRun", input.DryRun)

	var mcmsContracts *changesetstate.MCMSWithTimelockState
	if input.MCMSConfig != nil {
		var mcmsErr error
		mcmsContracts, mcmsErr = strategies.GetMCMSContracts(e, input.ContractChainSelector, *input.MCMSConfig)
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", mcmsErr)
		}
	}

	chain, ok := e.BlockChains.EVMChains()[input.ContractChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found in environment", input.ContractChainSelector)
	}

	contractRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.ContractChainSelector, input.ContractQualifier)
	contractAddrRef, err := e.DataStore.Addresses().Get(contractRefKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get VaultPlugin contract address for chain selector %d and qualifier %s: %w", input.ContractChainSelector, input.ContractQualifier, err)
	}
	contractAddr := common.HexToAddress(contractAddrRef.Address)

	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		input.MCMSConfig,
		mcmsContracts,
		common.HexToAddress(contractAddrRef.Address),
		changeset.ConfigureOCR3Description,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	dkgRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.ContractChainSelector, input.InstanceID.DKGContractQualifier)
	dkgAddrRef, err := e.DataStore.Addresses().Get(dkgRefKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get VaultPlugin contract address for chain selector %d and qualifier %s: %w", input.ContractChainSelector, input.ContractQualifier, err)
	}
	dkgAddr := common.HexToAddress(dkgAddrRef.Address)

	configDigestBytes, err := hex.DecodeString(input.InstanceID.ConfigDigest)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to decode config digest: %w", err)
	}
	instanceID := string(dkgocrtypes.MakeInstanceID(dkgAddr, [32]byte(configDigestBytes)))
	input.ReportingPluginConfig.DKGInstanceID = &instanceID

	cfgb, err := proto.Marshal(input.ReportingPluginConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal VaultPlugin reporting plugin config: %w", err)
	}

	report, err := operations.ExecuteOperation(e.OperationsBundle, ConfigureOCR3_1, ConfigureOCR3_1Deps{
		Env:      &e,
		Strategy: strategy,
	}, ConfigureOCR3_1Input{
		ContractAddress:               &contractAddr,
		ChainSelector:                 input.ContractChainSelector,
		DON:                           input.DON,
		Config:                        input.OracleConfig,
		DryRun:                        input.DryRun,
		MCMSConfig:                    input.MCMSConfig,
		ReportingPluginConfigOverride: cfgb,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3 contract: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: report.Output.MCMSTimelockProposals,
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
