package changeset

import (
	"fmt"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var TransferERC20Changeset cldf.ChangeSetV2[types.TransferERC20Config] = transferERC20Changeset{}

type transferERC20Changeset struct{}

func (t transferERC20Changeset) VerifyPreconditions(e cldf.Environment, cfg types.TransferERC20Config) error {
	return ValidateTransferERC20Config(e.GetContext(), e, cfg)
}

func (t transferERC20Changeset) Apply(e cldf.Environment, cfg types.TransferERC20Config) (cldf.ChangesetOutput, error) {
	lggr := e.Logger

	lggr.Infow("Starting batch ERC20 transfer",
		"chains", len(cfg.TransfersByChain),
		"description", cfg.Description)

	evmChains := e.BlockChains.EVMChains()

	for chainSelector := range cfg.TransfersByChain {
		if _, exists := evmChains[chainSelector]; !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}
	}

	var primaryChain cldf_evm.Chain
	for chainSelector := range cfg.TransfersByChain {
		primaryChain = evmChains[chainSelector]
		break
	}

	deps := VaultDeps{
		Chain:       primaryChain,
		Auth:        primaryChain.DeployerKey,
		DataStore:   e.DataStore,
		Environment: e,
	}

	seqInput := TransferERC20SequenceInput{
		TransfersByChain: cfg.TransfersByChain,
		MCMSConfig:       cfg.MCMSConfig,
		Description:      cfg.Description,
	}

	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, TransferERC20Sequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute ERC20 transfer sequence: %w", err)
	}

	lggr.Infow("batch ERC20 transfer completed successfully",
		"chains", len(cfg.TransfersByChain),
		"mcms_proposals", len(seqReport.Output.MCMSTimelockProposals),
		"execution_reports", len(seqReport.ExecutionReports))

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
		Reports:               seqReport.ExecutionReports,
	}, nil
}
