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
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type ConfigureOCR3Deps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
	Registry             *capabilities_registry_v2.CapabilitiesRegistry
}

type ConfigureOCR3Input struct {
	ContractAddress  *common.Address
	RegistryChainSel uint64
	DONs             []ConfigureCREDON
	Config           *ocr3.OracleConfig
	DryRun           bool

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

		var nodeIDs []string
		for _, don := range input.DONs {
			donConfig := RegisteredDonConfig{
				NodeIDs:          don.NodeIDs,
				Name:             don.Name,
				RegistryChainSel: input.RegistryChainSel,
				Registry:         deps.Registry,
			}
			d, err := newRegisteredDon(*deps.Env, donConfig)
			if err != nil {
				return ConfigureOCR3OpOutput{}, fmt.Errorf("configure-ocr3-op failed: failed to create registered DON %s: %w", don.Name, err)
			}

			// We double-check that the DON accepts workflows...
			if d.Info.AcceptsWorkflows {
				for _, node := range d.Nodes {
					nodeIDs = append(nodeIDs, node.NodeID)
				}
			}
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
		if !ok {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("chain %d not found in environment", input.RegistryChainSel)
		}

		contract, err := contracts.GetOwnedContractV2[*ocr3_capability.OCR3Capability](deps.Env.DataStore.Addresses(), chain, input.ContractAddress.Hex())
		if err != nil {
			return ConfigureOCR3OpOutput{}, fmt.Errorf("failed to get OCR3 contract: %w", err)
		}

		resp, err := ocr3.ConfigureOCR3ContractFromJD(deps.Env, ocr3.ConfigureOCR3Config{
			ChainSel:   input.RegistryChainSel,
			NodeIDs:    nodeIDs,
			OCR3Config: input.Config,
			Contract:   contract.Contract,
			DryRun:     input.DryRun,
			UseMCMS:    input.UseMCMS(),
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
				input.RegistryChainSel: contract.McmsContracts.Timelock.Address().Hex(),
			}
			proposerMCMSes := map[uint64]string{
				input.RegistryChainSel: contract.McmsContracts.ProposerMcm.Address().Hex(),
			}

			inspector, err := proposalutils.McmsInspectorForChain(*deps.Env, input.RegistryChainSel)
			if err != nil {
				return ConfigureOCR3OpOutput{}, err
			}
			inspectorPerChain := map[uint64]sdk.Inspector{
				input.RegistryChainSel: inspector,
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
