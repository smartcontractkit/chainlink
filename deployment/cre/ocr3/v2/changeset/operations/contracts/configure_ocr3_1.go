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

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

type ConfigureOCR3_1Deps struct {
	Env                  *cldf.Environment
	WriteGeneratedConfig io.Writer
}

type ConfigureOCR3_1Input struct {
	ContractAddress *common.Address
	ChainSelector   uint64
	DON             DonNodeSet
	Config          *ocr3.V3_1OracleConfig
	DryRun          bool

	ReportingPluginConfigOverride []byte

	MCMSConfig *ocr3.MCMSConfig
}

func (i ConfigureOCR3_1Input) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ConfigureOCR3_1OpOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var ConfigureOCR3_1 = operations.NewOperation[ConfigureOCR3_1Input, ConfigureOCR3_1OpOutput, ConfigureOCR3_1Deps](
	"configure-ocr3-1-op",
	semver.MustParse("1.0.0"),
	"Configure OCR3.1 Contract",
	func(b operations.Bundle, deps ConfigureOCR3_1Deps, input ConfigureOCR3_1Input) (ConfigureOCR3_1OpOutput, error) {
		if input.ContractAddress == nil {
			return ConfigureOCR3_1OpOutput{}, errors.New("ContractAddress is required")
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return ConfigureOCR3_1OpOutput{}, fmt.Errorf("chain %d not found in environment", input.ChainSelector)
		}

		contract, err := contracts.GetOwnedContractV2[*ocr3_capability.OCR3Capability](deps.Env.DataStore.Addresses(), chain, input.ContractAddress.Hex())
		if err != nil {
			return ConfigureOCR3_1OpOutput{}, fmt.Errorf("failed to get OCR3 contract: %w", err)
		}

		nodes, err := deployment.NodeInfo(input.DON.NodeIDs, deps.Env.Offchain)
		if err != nil {
			return ConfigureOCR3_1OpOutput{}, err
		}

		config, err := ocr3.GenerateOCR3_1ConfigFromNodes(
			*input.Config,
			nodes,
			input.ChainSelector,
			deps.Env.OCRSecrets,
			input.ReportingPluginConfigOverride,
		)
		if err != nil {
			return ConfigureOCR3_1OpOutput{}, fmt.Errorf("failed to generate DKG config: %w", err)
		}
		resp, err := ocr3.ConfigureOCR3contract(ocr3.ConfigureOCR3Request{
			Config:   config,
			Chain:    chain,
			Contract: contract.Contract,
			DryRun:   input.DryRun,
			UseMCMS:  input.UseMCMS(),
		})
		if err != nil {
			return ConfigureOCR3_1OpOutput{}, err
		}
		if w := deps.WriteGeneratedConfig; w != nil {
			b, err := json.MarshalIndent(&resp.OcrConfig, "", "  ")
			if err != nil {
				return ConfigureOCR3_1OpOutput{}, fmt.Errorf("failed to marshal response output: %w", err)
			}
			deps.Env.Logger.Infof("Generated OCR3 config: %s", string(b))
			n, err := w.Write(b)
			if err != nil {
				return ConfigureOCR3_1OpOutput{}, fmt.Errorf("failed to write response output: %w", err)
			}
			if n != len(b) {
				return ConfigureOCR3_1OpOutput{}, errors.New("failed to write all bytes")
			}
		}

		// does not create any new addresses
		var out ConfigureOCR3_1OpOutput
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
				return ConfigureOCR3_1OpOutput{}, err
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
				"proposal to set OCR3.1 config",
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
