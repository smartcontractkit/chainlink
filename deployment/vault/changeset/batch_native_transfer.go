package changeset

import (
	"fmt"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var BatchNativeTransferChangeset = cldf.CreateChangeSet(batchNativeTransferLogic, batchNativeTransferPrecondition)

func batchNativeTransferPrecondition(e cldf.Environment, cfg types.BatchNativeTransferConfig) error {
	ctx := e.GetContext()
	return ValidateBatchNativeTransferConfig(ctx, e, cfg)
}

func batchNativeTransferLogic(e cldf.Environment, cfg types.BatchNativeTransferConfig) (cldf.ChangesetOutput, error) {
	lggr := e.Logger

	lggr.Infow("Starting batch native transfer",
		"chains", len(cfg.TransfersByChain),
		"mcms_mode", cfg.MCMSConfig != nil,
		"description", cfg.Description)

	mutableDS := datastore.NewMemoryDataStore()
	if e.DataStore != nil {
		if err := mutableDS.Merge(e.DataStore); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge existing datastore: %w", err)
		}
	}

	evmChains := e.BlockChains.EVMChains()

	for chainSelector := range cfg.TransfersByChain {
		if _, exists := evmChains[chainSelector]; !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}
	}

	// Pick the first chain for deps (only needed for direct execution, not MCMS)
	var primaryChain cldf_evm.Chain
	for chainSelector := range cfg.TransfersByChain {
		primaryChain = evmChains[chainSelector]
		break
	}

	deps := VaultDeps{
		Chain:       primaryChain,
		Auth:        primaryChain.DeployerKey,
		DataStore:   mutableDS,
		Environment: e,
	}

	seqInput := BatchNativeTransferSequenceInput{
		TransfersByChain: cfg.TransfersByChain,
		MCMSConfig:       cfg.MCMSConfig,
		Description:      cfg.Description,
	}

	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, BatchNativeTransferSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute batch native transfer sequence: %w", err)
	}

	lggr.Infow("batch native transfer completed successfully",
		"chains", len(cfg.TransfersByChain),
		"mcms_proposals", len(seqReport.Output.MCMSTimelockProposals),
		"execution_reports", len(seqReport.ExecutionReports))

	return cldf.ChangesetOutput{
		DataStore:             mutableDS,
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
		Reports:               seqReport.ExecutionReports,
	}, nil
}
