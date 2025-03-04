package changeset

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[uint64] = DeployOCR3

func DeployOCR3(env deployment.Environment, registryChainSel uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	// ocr3 only deployed on registry chain
	c, ok := env.Chains[registryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ocr3Resp, err := kslib.DeployOCR3(c, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", ocr3Resp.Tv.String(), c.Selector, ocr3Resp.Address.String())
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

var _ deployment.ChangeSet[ConfigureOCR3Config] = ConfigureOCR3Contract

type ConfigureOCR3Config struct {
	ChainSel             uint64
	NodeIDs              []string
	Address              *common.Address // address of the OCR3 contract to configure
	OCR3Config           *kslib.OracleConfig
	DryRun               bool
	WriteGeneratedConfig io.Writer // if not nil, write the generated config to this writer as JSON [OCR2OracleConfig]

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
}

func (cfg ConfigureOCR3Config) UseMCMS() bool {
	return cfg.MCMSConfig != nil
}

func ConfigureOCR3Contract(env deployment.Environment, cfg ConfigureOCR3Config) (deployment.ChangesetOutput, error) {
	resp, err := kslib.ConfigureOCR3ContractFromJD(&env, kslib.ConfigureOCR3Config{
		ChainSel:   cfg.ChainSel,
		NodeIDs:    cfg.NodeIDs,
		OCR3Config: cfg.OCR3Config,
		Address:    cfg.Address,
		DryRun:     cfg.DryRun,
		UseMCMS:    cfg.UseMCMS(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3Capability: %w", err)
	}
	if w := cfg.WriteGeneratedConfig; w != nil {
		b, err := json.MarshalIndent(&resp.OCR2OracleConfig, "", "  ")
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal response output: %w", err)
		}
		env.Logger.Infof("Generated OCR3 config: %s", string(b))
		n, err := w.Write(b)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to write response output: %w", err)
		}
		if n != len(b) {
			return deployment.ChangesetOutput{}, errors.New("failed to write all bytes")
		}
	}
	// does not create any new addresses
	var out deployment.ChangesetOutput
	if cfg.UseMCMS() {
		if resp.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		r, err := GetContractSets(env.Logger, &GetContractSetsRequest{
			Chains:      env.Chains,
			AddressBook: env.ExistingAddresses,
		})
		if err != nil {
			return out, fmt.Errorf("failed to get contract sets: %w", err)
		}
		contracts := r.ContractSets[cfg.ChainSel]
		timelocksPerChain := map[uint64]string{
			cfg.ChainSel: contracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			cfg.ChainSel: contracts.ProposerMcm.Address().Hex(),
		}

		inspector, err := proposalutils.McmsInspectorForChain(env, cfg.ChainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, err
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
			cfg.MCMSConfig.MinDuration,
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}
	return out, nil
}
