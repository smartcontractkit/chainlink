package changeset

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[ConfigureVaultDKGInput] = ConfigureVaultDKG{}

type ConfigureVaultDKGInput struct {
	ContractChainSelector uint64 `json:"contractChainSelector" yaml:"contractChainSelector"`
	ContractQualifier     string `json:"contractQualifier" yaml:"contractQualifier"`

	DON          DKGDon             `json:"don" yaml:"don"`
	OracleConfig *ocr3.OracleConfig `json:"oracleConfig" yaml:"oracleConfig"`
	DryRun       bool               `json:"dryRun" yaml:"dryRun"`

	MCMSConfig *ocr3.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type DKGDon struct {
	contracts.DonNodeSet
	RecipientPublicKeys []string `json:"recipientPublicKeys" yaml:"recipientPublicKeys"`
}

type ConfigureVaultDKG struct{}

func (l ConfigureVaultDKG) VerifyPreconditions(_ cldf.Environment, input ConfigureVaultDKGInput) error {
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
	if len(input.DON.RecipientPublicKeys) == 0 {
		return errors.New("at least one recipient public key is required")
	}
	if len(input.DON.NodeIDs) != len(input.DON.RecipientPublicKeys) {
		return errors.New("the number of don node IDs must match the number of recipient public keys")
	}
	if input.OracleConfig == nil {
		return errors.New("oracle config is required")
	}
	return nil
}

func (l ConfigureVaultDKG) Apply(e cldf.Environment, input ConfigureVaultDKGInput) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Configuring Vault DKG contract with DON", "donName", input.DON.Name, "nodes", input.DON.NodeIDs, "dryRun", input.DryRun)

	contractRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.ContractChainSelector, input.ContractQualifier)
	contractAddrRef, err := e.DataStore.Addresses().Get(contractRefKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", input.ContractChainSelector, input.ContractQualifier, err)
	}
	contractAddr := common.HexToAddress(contractAddrRef.Address)
	cfg, err := dkgReportingPluginConfig(input.DON, input.OracleConfig.MaxFaultyOracles+1) // validate config can be created
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create DKG reporting plugin config: %w", err)
	}
	report, err := operations.ExecuteOperation(e.OperationsBundle, contracts.ConfigureDKG, contracts.ConfigureDKGDeps{
		WriteGeneratedConfig: io.Discard,
		Env:                  &e,
	}, contracts.ConfigureDKGInput{
		ContractAddress:       &contractAddr,
		ChainSelector:         input.ContractChainSelector,
		DON:                   input.DON.DonNodeSet,
		Config:                input.OracleConfig,
		DryRun:                input.DryRun,
		MCMSConfig:            input.MCMSConfig,
		ReportingPluginConfig: cfg,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3 contract: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: report.Output.MCMSTimelockProposals,
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

func dkgReportingPluginConfig(don DKGDon, threshold int) (dkgocrtypes.ReportingPluginConfig, error) {
	keys := []dkgocrtypes.P256ParticipantPublicKey{}
	for _, k := range don.RecipientPublicKeys {
		bk, err := hex.DecodeString(k)
		if err != nil {
			return dkgocrtypes.ReportingPluginConfig{}, fmt.Errorf("failed to decode recipient public key %s: %w", k, err)
		}
		keys = append(keys, dkgocrtypes.P256ParticipantPublicKey(bk))
	}
	cfg := dkgocrtypes.ReportingPluginConfig{
		T:                   threshold,
		DealerPublicKeys:    keys,
		RecipientPublicKeys: keys,
	}

	return cfg, nil
}
