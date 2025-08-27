package sequences

import (
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

type DeployOCR3Deps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
}

type DeployOCR3Input struct {
	RegistryChainSel uint64
	Qualifier        string

	DONs         []contracts.ConfigureCREDON
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
			ChainSelector: input.RegistryChainSel,
			Qualifier:     input.Qualifier,
		})
		if err != nil {
			return DeployOCR3Output{}, err
		}

		ocr3ContractAddress := common.HexToAddress(ocr3DeploymentReport.Output.Address)

		// Update the environment datastore to include the newly deployed OCR3 contract
		deps.Env.DataStore = ocr3DeploymentReport.Output.Datastore

		// Step 2: Get the capabilities registry contract
		capabilitiesRegistry, err := getCapabilitiesRegistryContract(deps, input)
		if err != nil {
			return DeployOCR3Output{}, fmt.Errorf("failed to get capabilities registry: %w", err)
		}

		// Step 3: Configure OCR3 Contract with DONs
		deps.Env.Logger.Infow("Configuring OCR3 contract with DONs",
			"numDONs", len(input.DONs),
			"dryRun", input.DryRun)

		_, err = operations.ExecuteOperation(b, contracts.ConfigureOCR3, contracts.ConfigureOCR3Deps{
			Env:                  deps.Env,
			WriteGeneratedConfig: deps.WriteGeneratedConfig,
			Registry:             capabilitiesRegistry,
		}, contracts.ConfigureOCR3Input{
			ContractAddress:  &ocr3ContractAddress,
			RegistryChainSel: input.RegistryChainSel,
			DONs:             input.DONs,
			Config:           input.OracleConfig,
			DryRun:           input.DryRun,
			MCMSConfig:       input.MCMSConfig,
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

func getCapabilitiesRegistryContract(deps DeployOCR3Deps, input DeployOCR3Input) (*capabilities_registry_v2.CapabilitiesRegistry, error) {
	refs := deps.Env.DataStore.Addresses().Filter(
		datastore.AddressRefByType("CapabilitiesRegistry"),
		datastore.AddressRefByChainSelector(input.RegistryChainSel))
	if len(refs) == 0 {
		return nil, errors.New("failed to get capabilities registry ref")
	}
	capabilitiesRegistryRef := refs[0]

	chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
	if !ok {
		return nil, fmt.Errorf("chain not found for selector %d", input.RegistryChainSel)
	}
	capabilitiesRegistry, err := crecontracts.GetOwnedContractV2[*capabilities_registry_v2.CapabilitiesRegistry](deps.Env.DataStore.Addresses(), chain, capabilitiesRegistryRef.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned contract: %w", err)
	}
	return capabilitiesRegistry.Contract, nil
}
