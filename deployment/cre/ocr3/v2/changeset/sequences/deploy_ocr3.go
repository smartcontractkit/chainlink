package sequences

import (
	"fmt"
	"io"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

type DeployOCR3Deps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
}

type DeployOCR3Input struct {
	ChainSelector uint64
	Qualifier     string

	DON          contracts.DonNodeSet
	OracleConfig *ocr3.OracleConfig
	DryRun       bool

	MCMSConfig *ocr3.MCMSConfig
}

func (c DeployOCR3Input) Validate() error {
	return nil
}

type DeployOCR3Output struct {
	ChainSelector uint64
	Address       string
	Type          string
	Version       string
	Labels        []string
}

var DeployOCR3 = operations.NewSequence(
	"deploy-ocr3",
	semver.MustParse("1.0.0"),
	"Deploys the OCR3 contract",
	func(b operations.Bundle, deps DeployOCR3Deps, input DeployOCR3Input) (DeployOCR3Output, error) {
		// Step 1: Deploy OCR3 Contract for Consensus Capability
		ocr3DeploymentReport, err := operations.ExecuteOperation(b, contracts.DeployOCR3, contracts.DeployOCR3Deps{Env: deps.Env}, contracts.DeployOCR3Input{
			ChainSelector: input.ChainSelector,
			Qualifier:     input.Qualifier,
		})
		if err != nil {
			return DeployOCR3Output{}, err
		}

		ocr3ContractAddress := common.HexToAddress(ocr3DeploymentReport.Output.Address)

		// Update the environment datastore to include the newly deployed OCR3 contract
		deps.Env.DataStore = ocr3DeploymentReport.Output.Datastore

		// Step 3: Configure OCR3 Contract with DONs
		deps.Env.Logger.Infow("Configuring OCR3 contract with DON",
			"nodes", input.DON.NodeIDs,
			"dryRun", input.DryRun)

		_, err = operations.ExecuteOperation(b, contracts.ConfigureOCR3, contracts.ConfigureOCR3Deps{
			Env:                  deps.Env,
			WriteGeneratedConfig: deps.WriteGeneratedConfig,
		}, contracts.ConfigureOCR3Input{
			ContractAddress: &ocr3ContractAddress,
			ChainSelector:   input.ChainSelector,
			DON:             input.DON,
			Config:          input.OracleConfig,
			DryRun:          input.DryRun,
			MCMSConfig:      input.MCMSConfig,
		})
		if err != nil {
			return DeployOCR3Output{}, fmt.Errorf("failed to configure OCR3 contract: %w", err)
		}

		return DeployOCR3Output{
			ChainSelector: ocr3DeploymentReport.Output.ChainSelector,
			Address:       ocr3DeploymentReport.Output.Address,
			Type:          ocr3DeploymentReport.Output.Type,
			Version:       ocr3DeploymentReport.Output.Version,
			Labels:        ocr3DeploymentReport.Output.Labels,
		}, nil
	},
)
