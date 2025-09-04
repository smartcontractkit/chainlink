package contracts

import (
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

/*
type ConfigureOCR3OpDeps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
}
*/
/*
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
*/
//type ConfigureOCR3OpInput contracts.ConfigureOCR3Input

//type ConfigureOCR3OpOutput = contracts.ConfigureOCR3OpOutput

var ConfigureOCR3Op = contracts.ConfigureOCR3

/*
var ConfigureOCR3Op = operations.NewOperation[ConfigureOCR3OpInput, ConfigureOCR3OpOutput, ConfigureOCR3OpDeps](
	"configure-ocr3-op",
	semver.MustParse("1.0.0"),
	"Configure OCR3 Contract",
	func(b operations.Bundle, deps ConfigureOCR3OpDeps, input ConfigureOCR3OpInput) (ConfigureOCR3OpOutput, error) {
		if input.ContractAddress == nil {
			return ConfigureOCR3OpOutput{}, errors.New("ContractAddress is required")
		}

		deps.Env.Logger.Infow("Configuring OCR3 contract with DON",
			"nodes", input.DON.NodeIDs,
			"dryRun", input.DryRun)
		resp, err := changeset.ConfigureOCR3Contract(*deps.Env, changeset.ConfigureOCR3Config{
			ChainSel:             input.RegistryChainSel,
			NodeIDs:              input.DON.NodeIDs,
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
*/
//
