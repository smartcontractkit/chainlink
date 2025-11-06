package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

var _ cldf.ChangeSetV2[ConfigureOCR3Input] = ConfigureOCR3{}

type ConfigureOCR3Input struct {
	ContractChainSelector uint64 `json:"contractChainSelector" yaml:"contractChainSelector"`
	ContractQualifier     string `json:"contractQualifier" yaml:"contractQualifier"`

	DON          contracts.DonNodeSet `json:"don" yaml:"don"`
	OracleConfig *ocr3.OracleConfig   `json:"oracleConfig" yaml:"oracleConfig"`
	DryRun       bool                 `json:"dryRun" yaml:"dryRun"`

	MCMSConfig *ocr3.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type ConfigureOCR3 struct{}

func (l ConfigureOCR3) VerifyPreconditions(_ cldf.Environment, input ConfigureOCR3Input) error {
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
	return nil
}

func (l ConfigureOCR3) Apply(e cldf.Environment, input ConfigureOCR3Input) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Configuring OCR3 contract with DON", "donName", input.DON.Name, "nodes", input.DON.NodeIDs, "dryRun", input.DryRun)

	contractRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.ContractChainSelector, input.ContractQualifier)
	contractAddrRef, err := e.DataStore.Addresses().Get(contractRefKey)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get OCR3 contract address for chain selector %d and qualifier %s: %w", input.ContractChainSelector, input.ContractQualifier, err)
	}
	contractAddr := common.HexToAddress(contractAddrRef.Address)

	report, err := operations.ExecuteOperation(e.OperationsBundle, contracts.ConfigureOCR3, contracts.ConfigureOCR3Deps{
		Env: &e,
	}, contracts.ConfigureOCR3Input{
		ContractAddress: &contractAddr,
		ChainSelector:   input.ContractChainSelector,
		DON:             input.DON,
		Config:          input.OracleConfig,
		DryRun:          input.DryRun,
		MCMSConfig:      input.MCMSConfig,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3 contract: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: report.Output.MCMSTimelockProposals,
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
