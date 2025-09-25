package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type ConfigureOCR3Deps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
}

type ConfigureOCR3Input struct {
	ContractAddress *common.Address
	ChainSelector   uint64
	DON             DonNodeSet
	Config          *ocr3.OracleConfig
	DryRun          bool

	ReportingPluginConfigOverride []byte

	MCMSConfig *ocr3.MCMSConfig
}

func (i ConfigureOCR3Input) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ConfigureOCR3OpOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var ConfigureOCR3 = operations.NewOperation[ConfigureOCR3Input, ConfigureOCR3OpOutput, ConfigureOCR3Deps](
	"configure-ocr3-op",
	semver.MustParse("1.0.0"),
	"Configure OCR3 Contract",
	func(b operations.Bundle, deps ConfigureOCR3Deps, input ConfigureOCR3Input) (ConfigureOCR3OpOutput, error) {
		if input.ContractAddress == nil {
			return ConfigureOCR3OpOutput{}, errors.New("ContractAddress is required")
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("chain %d not found in environment", input.ChainSelector)
		}

		contract, err := contracts.GetOwnedContractV2[*ocr3_capability.OCR3Capability](deps.Env.DataStore.Addresses(), chain, input.ContractAddress.Hex())
		if err != nil {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("failed to get OCR3 contract: %w", err)
		}

		resp, err := ocr3.ConfigureOCR3ContractFromJD(deps.Env, ocr3.ConfigureOCR3Config{
			ChainSel:                      input.ChainSelector,
			NodeIDs:                       input.DON.NodeIDs,
			OCR3Config:                    input.Config,
			Contract:                      contract.Contract,
			DryRun:                        input.DryRun,
			UseMCMS:                       input.UseMCMS(),
			ReportingPluginConfigOverride: input.ReportingPluginConfigOverride,
		})
		if err != nil {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("failed to configure OCR3Capability: %w", err)
		}
		if w := deps.WriteGeneratedConfig; w != nil {
			b, err := json.MarshalIndent(&resp.OCR2OracleConfig, "", "  ")
			if err != nil {
				return ConfigureOCR3OpOutput{}, fmt.Errorf("failed to marshal response output: %w", err)
			}
			deps.Env.Logger.Infof("Generated OCR3 config: %s", string(b))
			n, err := w.Write(b)
			if err != nil {
				return ConfigureOCR3OpOutput{}, fmt.Errorf("failed to write response output: %w", err)
			}
			if n != len(b) {
				return ConfigureOCR3OpOutput{}, errors.New("failed to write all bytes")
			}
		}

		// does not create any new addresses
		var out ConfigureOCR3OpOutput
		if input.UseMCMS() {
			if resp.Ops == nil {
				return out, errors.New("expected MCMS operation to be non-nil")
			}

			if contract.McmsContracts == nil {
				return out, fmt.Errorf("expected OCR3 capabilty contract %s to be owned by MCMS", contract.Contract.Address().String())
			}

			timelocksPerChain := map[uint64]string{
				input.ChainSelector: contract.McmsContracts.Timelock.Address().Hex(),
			}
			proposerMCMSes := map[uint64]string{
				input.ChainSelector: contract.McmsContracts.ProposerMcm.Address().Hex(),
			}

			inspector, err := proposalutils.McmsInspectorForChain(*deps.Env, input.ChainSelector)
			if err != nil {
				return ConfigureOCR3OpOutput{}, err
			}
			inspectorPerChain := map[uint64]sdk.Inspector{
				input.ChainSelector: inspector,
			}
			proposal, err := proposalutils.BuildProposalFromBatchesV2(
				*deps.Env,
				timelocksPerChain,
				proposerMCMSes,
				inspectorPerChain,
				[]mcmstypes.BatchOperation{*resp.Ops},
				"proposal to set OCR3 config",
				proposalutils.TimelockConfig{MinDelay: input.MCMSConfig.MinDuration},
			)
			if err != nil {
				return out, fmt.Errorf("failed to build proposal: %w", err)
			}
			out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
		}
		return out, nil
	},
)
