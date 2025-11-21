package changeset

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	changesetstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[*DeployRequestV2] = DeployOCR3V2

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
	OCR3Config           *ocr3.OracleConfig
	DryRun               bool
	WriteGeneratedConfig io.Writer // if not nil, write the generated config to this writer as JSON [OCR2OracleConfig]

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
}

func (cfg ConfigureOCR3Config) UseMCMS() bool {
	return cfg.MCMSConfig != nil
}

func ConfigureOCR3Contract(env cldf.Environment, cfg ConfigureOCR3Config) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[cfg.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSel)
	}

	if cfg.Address == nil {
		return cldf.ChangesetOutput{}, errors.New("address of OCR3 contract to configure is required")
	}

	contract, err := contracts.GetOwnedContractV2[*ocr3_capability.OCR3Capability](env.DataStore.Addresses(), chain, cfg.Address.Hex())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get OCR3 contract: %w", err)
	}

	var mcmsContracts *changesetstate.MCMSWithTimelockState
	if cfg.UseMCMS() {
		var mcmsErr error
		mcmsContracts, mcmsErr = strategies.GetMCMSContracts(env, cfg.ChainSel, *cfg.MCMSConfig)
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", mcmsErr)
		}
	}

	strategy, err := strategies.CreateStrategy(
		chain,
		env,
		cfg.MCMSConfig,
		mcmsContracts,
		contract.Contract.Address(),
		"Configure OCR3 contract",
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	resp, err := ocr3.ConfigureOCR3ContractFromJD(&env, ocr3.ConfigureOCR3Config{
		ChainSel:   cfg.ChainSel,
		NodeIDs:    cfg.NodeIDs,
		OCR3Config: cfg.OCR3Config,
		Contract:   contract.Contract,
		DryRun:     cfg.DryRun,
		UseMCMS:    cfg.UseMCMS(),
		Strategy:   strategy,
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

		if contract.McmsContracts == nil {
			return out, fmt.Errorf("expected OCR3 capabilty contract %s to be owned by MCMS", contract.Contract.Address().String())
		}

		proposal, err := strategy.BuildProposal([]mcmstypes.BatchOperation{*resp.Ops})
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}
	return out, nil
}
