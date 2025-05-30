package changeset

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"

	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[uint64] = DeployOCR3

// Deprecated: use DeployOCR3V2 instead
func DeployOCR3(env cldf.Environment, registryChainSel uint64) (cldf.ChangesetOutput, error) {
	return DeployOCR3V2(env, &DeployRequestV2{
		ChainSel: registryChainSel,
	})
}

var _ cldf.ChangeSet[ConfigureOCR3Config] = ConfigureOCR3Contract

func DeployOCR3V2(env cldf.Environment, req *DeployRequestV2) (cldf.ChangesetOutput, error) {
	req.deployFn = internal.DeployOCR3
	return deploy(env, req)
}

type ConfigureOCR3Config struct {
	ChainSel             uint64
	NodeIDs              []string
	Address              *common.Address // address of the OCR3 contract to configure
	OCR3Config           *internal.OracleConfig
	DryRun               bool
	WriteGeneratedConfig io.Writer // if not nil, write the generated config to this writer as JSON [OCR2OracleConfig]

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
}

func (cfg ConfigureOCR3Config) UseMCMS() bool {
	return cfg.MCMSConfig != nil
}

func ConfigureOCR3Contract(env cldf.Environment, cfg ConfigureOCR3Config) (cldf.ChangesetOutput, error) {
	resp, err := internal.ConfigureOCR3ContractFromJD(&env, internal.ConfigureOCR3Config{
		ChainSel:   cfg.ChainSel,
		NodeIDs:    cfg.NodeIDs,
		OCR3Config: cfg.OCR3Config,
		Address:    cfg.Address,
		DryRun:     cfg.DryRun,
		UseMCMS:    cfg.UseMCMS(),
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3Capability: %w", err)
	}
	if w := cfg.WriteGeneratedConfig; w != nil {
		b, err := json.MarshalIndent(&resp.OCR2OracleConfig, "", "  ")
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal response output: %w", err)
		}
		env.Logger.Infof("Generated OCR3 config: %s", string(b))
		n, err := w.Write(b)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to write response output: %w", err)
		}
		if n != len(b) {
			return cldf.ChangesetOutput{}, errors.New("failed to write all bytes")
		}
	}
	// does not create any new addresses
	var out cldf.ChangesetOutput
	if cfg.UseMCMS() {
		if resp.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}

		chain, ok := env.BlockChains.EVMChains()[cfg.ChainSel]
		if !ok {
			return out, fmt.Errorf("chain %d not found in environment", cfg.ChainSel)
		}

		contract, err := GetOwnedContractV2[*ocr3_capability.OCR3Capability](env.DataStore.Addresses(), chain, cfg.Address.Hex())
		if err != nil {
			return out, fmt.Errorf("failed to get OCR3 contract: %w", err)
		}

		if contract.McmsContracts == nil {
			return out, fmt.Errorf("expected OCR3 capabilty contract %s to be owned by MCMS", contract.Contract.Address().String())
		}

		timelocksPerChain := map[uint64]string{
			cfg.ChainSel: contract.McmsContracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			cfg.ChainSel: contract.McmsContracts.ProposerMcm.Address().Hex(),
		}

		inspector, err := proposalutils.McmsInspectorForChain(env, cfg.ChainSel)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		inspectorPerChain := map[uint64]sdk.Inspector{
			cfg.ChainSel: inspector,
		}
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocksPerChain,
			proposerMCMSes,
			inspectorPerChain,
			[]mcmstypes.BatchOperation{*resp.Ops},
			"proposal to set OCR3 config",
			proposalutils.TimelockConfig{MinDelay: cfg.MCMSConfig.MinDuration},
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}
	return out, nil
}
