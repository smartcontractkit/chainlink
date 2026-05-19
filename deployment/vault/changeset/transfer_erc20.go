package changeset

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var TransferERC20Changeset cldf.ChangeSetV2[types.TransferERC20Config] = transferERC20Changeset{}

type transferERC20Changeset struct{}

func (t transferERC20Changeset) VerifyPreconditions(e cldf.Environment, cfg types.TransferERC20Config) error {
	return ValidateTransferERC20Config(e, cfg)
}

func (t transferERC20Changeset) Apply(e cldf.Environment, cfg types.TransferERC20Config) (cldf.ChangesetOutput, error) {
	lggr := e.Logger

	lggr.Infow("Starting ERC20 transfer from timelock",
		"chain", cfg.ChainSelector,
		"timelock_id", cfg.TimelockIdentifier,
		"transfers", len(cfg.Transfers),
		"description", cfg.Description)

	evmChains := e.BlockChains.EVMChains()
	chain, exists := evmChains[cfg.ChainSelector]
	if !exists {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	deps := VaultDeps{
		Chain:       chain,
		Auth:        chain.DeployerKey,
		DataStore:   e.DataStore,
		Environment: e,
	}

	seqInput := TransferERC20SequenceInput{
		ChainSelector:      cfg.ChainSelector,
		TimelockIdentifier: cfg.TimelockIdentifier,
		Transfers:          cfg.Transfers,
		MCMSConfig:         cfg.MCMSConfig,
		Description:        cfg.Description,
	}

	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, TransferERC20Sequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute ERC20 transfer sequence: %w", err)
	}

	lggr.Infow("ERC20 transfer completed successfully",
		"chain", cfg.ChainSelector,
		"transfers", len(cfg.Transfers),
		"mcms_proposals", len(seqReport.Output.MCMSTimelockProposals),
		"execution_reports", len(seqReport.ExecutionReports))

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
		Reports:               seqReport.ExecutionReports,
	}, nil
}
