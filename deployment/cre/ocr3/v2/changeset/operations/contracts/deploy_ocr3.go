package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
)

type DeployOCR3Deps struct {
	Env *cldf.Environment
}

type DeployOCR3Input struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployOCR3Output struct {
	Address       string
	ChainSelector uint64
	Qualifier     string
	Type          string
	Version       string
	Labels        []string
}

// DeployOCR3 is an operation that deploys the OCR3 contract.
// This atomic operation performs the single side effect of deploying and registering the contract.
var DeployOCR3 = operations.NewOperation[DeployOCR3Input, DeployOCR3Output, DeployOCR3Deps](
	"deploy-ocr3-op",
	semver.MustParse("1.0.0"),
	"Deploy OCR3 Contract",
	func(b operations.Bundle, deps DeployOCR3Deps, input DeployOCR3Input) (DeployOCR3Output, error) {
		lggr := deps.Env.Logger

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployOCR3Output{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Deploy the OCR3 contract
		ocr3Addr, tx, ocr3, err := ocr3_capability.DeployOCR3Capability(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return DeployOCR3Output{}, fmt.Errorf("failed to deploy OCR3: %w", err)
		}

		// Wait for deployment confirmation
		_, err = chain.Confirm(tx)
		if err != nil {
			return DeployOCR3Output{}, fmt.Errorf("failed to confirm OCR3 deployment: %w", err)
		}

		// Get type and version from the deployed contract
		tvStr, err := ocr3.TypeAndVersion(&bind.CallOpts{})
		if err != nil {
			return DeployOCR3Output{}, fmt.Errorf("failed to get type and version: %w", err)
		}

		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return DeployOCR3Output{}, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}

		lggr.Infof("Deployed %s on chain selector %d at address %s", tv.String(), chain.Selector, ocr3Addr.String())

		return DeployOCR3Output{
			Address:       ocr3Addr.String(),
			ChainSelector: input.ChainSelector,
			Qualifier:     input.Qualifier,
			Type:          string(tv.Type),
			Version:       tv.Version.String(),
			Labels:        tv.Labels.List(),
		}, nil
	},
)
