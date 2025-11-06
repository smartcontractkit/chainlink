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
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type ConfigureDKGDeps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
}

type ConfigureDKGInput struct {
	ContractAddress *common.Address
	ChainSelector   uint64
	DON             DonNodeSet
	Config          *ocr3.V3_1OracleConfig
	DryRun          bool

	MCMSConfig            *ocr3.MCMSConfig
	ReportingPluginConfig dkgocrtypes.ReportingPluginConfig
}

func (i ConfigureDKGInput) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ConfigureDKGOpOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var ConfigureDKG = operations.NewOperation(
	"configure-dkg-op",
	semver.MustParse("1.0.0"),
	"Configure DKG Contract",
	func(b operations.Bundle, deps ConfigureDKGDeps, input ConfigureDKGInput) (ConfigureDKGOpOutput, error) {
		if input.ContractAddress == nil {
			return ConfigureDKGOpOutput{}, errors.New("ContractAddress is required")
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return ConfigureDKGOpOutput{}, fmt.Errorf("chain %d not found in environment", input.ChainSelector)
		}

		contract, err := contracts.GetOwnedContractV2[*ocr3_capability.OCR3Capability](deps.Env.DataStore.Addresses(), chain, input.ContractAddress.Hex())
		if err != nil {
			return ConfigureDKGOpOutput{}, fmt.Errorf("failed to get DKG contract: %w", err)
		}

		nodes, err := deployment.NodeInfo(input.DON.NodeIDs, deps.Env.Offchain)
		if err != nil {
			return ConfigureDKGOpOutput{}, err
		}

		config, err := ocr3.GenerateDKGConfigFromNodes(
			*input.Config,
			nodes,
			input.ChainSelector,
			deps.Env.OCRSecrets,
			input.ReportingPluginConfig,
		)
		if err != nil {
			return ConfigureDKGOpOutput{}, fmt.Errorf("failed to generate DKG config: %w", err)
		}

		resp, err := ocr3.ConfigureOCR3contract(ocr3.ConfigureOCR3Request{
			Config:   config,
			Chain:    chain,
			Contract: contract.Contract,
			DryRun:   input.DryRun,
			UseMCMS:  input.UseMCMS(),
		})
		if err != nil {
			return ConfigureDKGOpOutput{}, err
		}
		if w := deps.WriteGeneratedConfig; w != nil {
			b, err := json.MarshalIndent(&resp.OcrConfig, "", "  ")
			if err != nil {
				return ConfigureDKGOpOutput{}, fmt.Errorf("failed to marshal response output: %w", err)
			}
			deps.Env.Logger.Infof("Generated DKG config: %s", string(b))
			n, err := w.Write(b)
			if err != nil {
				return ConfigureDKGOpOutput{}, fmt.Errorf("failed to write response output: %w", err)
			}
			if n != len(b) {
				return ConfigureDKGOpOutput{}, errors.New("failed to write all bytes")
			}
		}

		// does not create any new addresses
		var out ConfigureDKGOpOutput
		if input.UseMCMS() {
			if resp.Ops == nil {
				return out, errors.New("expected MCMS operation to be non-nil")
			}

			if contract.McmsContracts == nil {
				return out, fmt.Errorf("expected DKG capabilty contract %s to be owned by MCMS", contract.Contract.Address().String())
			}

			timelocksPerChain := map[uint64]string{
				input.ChainSelector: contract.McmsContracts.Timelock.Address().Hex(),
			}
			proposerMCMSes := map[uint64]string{
				input.ChainSelector: contract.McmsContracts.ProposerMcm.Address().Hex(),
			}

			inspector, err := proposalutils.McmsInspectorForChain(*deps.Env, input.ChainSelector)
			if err != nil {
				return ConfigureDKGOpOutput{}, err
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
				"proposal to set DKG config",
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
