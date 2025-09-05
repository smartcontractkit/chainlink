package contracts

import (
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type ConfigureOCR3OpDeps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
	Registry             *capabilities_registry.CapabilitiesRegistry
}

type ConfigureOCR3OpInput struct {
	ContractAddress  *common.Address
	RegistryChainSel uint64
	DON              ConfigureKeystoneDON
	Config           *ocr3.OracleConfig
	DryRun           bool

	MCMSConfig *changeset.MCMSConfig
}

func (i ConfigureOCR3OpInput) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ConfigureOCR3OpOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var ConfigureOCR3Op = operations.NewOperation[ConfigureOCR3OpInput, ConfigureOCR3OpOutput, ConfigureOCR3OpDeps](
	"configure-ocr3-op",
	semver.MustParse("1.0.0"),
	"Configure OCR3 Contract",
	func(b operations.Bundle, deps ConfigureOCR3OpDeps, input ConfigureOCR3OpInput) (ConfigureOCR3OpOutput, error) {
		if input.ContractAddress == nil {
			return ConfigureOCR3OpOutput{}, errors.New("ContractAddress is required")
		}

		donConfig := internal.RegisteredDonConfig{
			NodeIDs:          input.DON.NodeIDs,
			Name:             input.DON.Name,
			RegistryChainSel: input.RegistryChainSel,
			Registry:         deps.Registry,
		}
		d, err := internal.NewRegisteredDon(*deps.Env, donConfig)
		if err != nil {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("configure-ocr3-op failed: failed to create registered DON %s: %w", input.DON.Name, err)
		}

		nodeIDs := make([]string, 0, len(d.Nodes))
		for _, node := range d.Nodes {
			nodeIDs = append(nodeIDs, node.NodeID)
		}

		resp, err := changeset.ConfigureOCR3Contract(*deps.Env, changeset.ConfigureOCR3Config{
			ChainSel:             input.RegistryChainSel,
			NodeIDs:              nodeIDs,
			Address:              input.ContractAddress,
			OCR3Config:           input.Config,
			DryRun:               input.DryRun,
			WriteGeneratedConfig: deps.WriteGeneratedConfig,
			MCMSConfig:           input.MCMSConfig,
		})
		if err != nil {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("configure-ocr3-op failed: %w", err)
		}

		return ConfigureOCR3OpOutput{MCMSTimelockProposals: resp.MCMSTimelockProposals}, nil
	},
)
