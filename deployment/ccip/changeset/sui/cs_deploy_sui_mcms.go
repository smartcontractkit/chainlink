package sui

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	suistate "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type DeploySuiMCMSConfig struct {
	ChainSelectors         []uint64
	MCMSConfigPerChain     map[uint64]commonTypes.MCMSWithTimelockConfigV2
	TimelockConfigPerChain map[uint64]proposalutils.TimelockConfig
}

var _ cldf.ChangeSetV2[DeploySuiMCMSConfig] = DeploySuiMCMS{}

type DeploySuiMCMS struct{}

func (cs DeploySuiMCMS) VerifyPreconditions(env cldf.Environment, cfg DeploySuiMCMSConfig) error {
	suiChains := env.BlockChains.SuiChains()
	var errs []error

	for _, chainSel := range cfg.ChainSelectors {
		if _, ok := suiChains[chainSel]; !ok {
			errs = append(errs, fmt.Errorf("SUI chain %d not found in environment", chainSel))
		}
		mcmsCfg, ok := cfg.MCMSConfigPerChain[chainSel]
		if !ok {
			errs = append(errs, fmt.Errorf("MCMS config not provided for SUI chain %d", chainSel))
			continue
		}
		for _, roleCfg := range []mcmstypes.Config{mcmsCfg.Canceller, mcmsCfg.Bypasser, mcmsCfg.Proposer} {
			if err := roleCfg.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("invalid MCMS role config for SUI chain %d: %w", chainSel, err))
			}
		}
	}

	return errors.Join(errs...)
}

func (cs DeploySuiMCMS) Apply(env cldf.Environment, cfg DeploySuiMCMSConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	suiStates, err := suistate.LoadOnchainStatesui(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load SUI onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)
	proposals := make([]mcms.TimelockProposal, 0)
	suiChains := env.BlockChains.SuiChains()

	for _, chainSel := range cfg.ChainSelectors {
		suiChain := suiChains[chainSel]

		// Skip if MCMS is already deployed
		if chainState, ok := suiStates[chainSel]; ok && chainState.MCMSPackageID != "" {
			env.Logger.Infow("MCMS already deployed on SUI chain, skipping", "chainSelector", chainSel, "packageId", chainState.MCMSPackageID)
			continue
		}

		deps := Deps{
			AB: ab,
			SuiChain: sui_ops.OpTxDeps{
				Client: suiChain.Client,
				Signer: suiChain.Signer,
				GetCallOpts: func() *bind.CallOpts {
					gasBudget := uint64(400_000_000)
					return &bind.CallOpts{WaitForExecution: true, GasBudget: &gasBudget}
				},
				SuiRPC: suiChain.URL,
			},
			ChainSelector:    chainSel,
			CCIPOnChainState: state,
		}

		mcmsConfig := cfg.MCMSConfigPerChain[chainSel]

		seqInput := seq.DeployMCMSSeqInput{
			ChainSelector: chainSel,
			MCMSConfig:    mcmsConfig,
		}

		mcmsSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployMCMSSequence, deps.SuiChain, seqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS on SUI chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, mcmsSeqReport.ExecutionReports...)

		err = storeMCMSInAddressBook(ab, chainSel, mcmsSeqReport.Output)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store MCMS in address book for SUI chain %d: %w", chainSel, err)
		}

		// Generate MCMS proposal if timelock config is provided
		if timelockCfg, ok := cfg.TimelockConfigPerChain[chainSel]; ok {
			updatedSuiStates, err := suistate.LoadOnchainStatesui(env)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to reload SUI state after deploy for chain %d: %w", chainSel, err)
			}
			chainState := updatedSuiStates[chainSel]

			mcmsOperations := []mcmstypes.BatchOperation{}
			proposal, err := mcmsutil.GenerateProposal(
				env,
				chainState,
				chainSel,
				mcmsOperations,
				"Deploy SUI MCMS",
				timelockCfg,
			)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for SUI chain %d: %w", chainSel, err)
			}
			proposals = append(proposals, *proposal)
		}

		// Always include the accept ownership proposal from the deploy sequence
		proposals = append(proposals, mcmsSeqReport.Output.AcceptOwnershipProposal)
	}

	ds, err := shared.PopulateDataStore(ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate in-memory DataStore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		MCMSTimelockProposals: proposals,
		Reports:               seqReports,
	}, nil
}

func storeMCMSInAddressBook(ab *cldf.AddressBookMap, chainSelector uint64, output seq.DeployMCMSSeqOutput) error {
	typeAndVersion := cldf.NewTypeAndVersion(suistate.SuiMcmsPackageIDType, suistate.Version1_0_0)
	if err := ab.Save(chainSelector, output.PackageID, typeAndVersion); err != nil {
		return fmt.Errorf("failed to save MCMS package ID %s: %w", output.PackageID, err)
	}

	typeAndVersion = cldf.NewTypeAndVersion(suistate.SuiMcmsObjectIDType, suistate.Version1_0_0)
	if err := ab.Save(chainSelector, output.Objects.McmsMultisigStateObjectId, typeAndVersion); err != nil {
		return fmt.Errorf("failed to save MCMS MultisigState object ID %s: %w", output.Objects.McmsMultisigStateObjectId, err)
	}

	typeAndVersion = cldf.NewTypeAndVersion(suistate.SuiMcmsRegistryObjectIDType, suistate.Version1_0_0)
	if err := ab.Save(chainSelector, output.Objects.McmsRegistryObjectId, typeAndVersion); err != nil {
		return fmt.Errorf("failed to save MCMS Registry object ID %s: %w", output.Objects.McmsRegistryObjectId, err)
	}

	typeAndVersion = cldf.NewTypeAndVersion(suistate.SuiMcmsAccountStateObjectIDType, suistate.Version1_0_0)
	if err := ab.Save(chainSelector, output.Objects.McmsAccountStateObjectId, typeAndVersion); err != nil {
		return fmt.Errorf("failed to save MCMS AccountState object ID %s: %w", output.Objects.McmsAccountStateObjectId, err)
	}

	typeAndVersion = cldf.NewTypeAndVersion(suistate.SuiMcmsAccountOwnerCapObjectIDType, suistate.Version1_0_0)
	if err := ab.Save(chainSelector, output.Objects.McmsAccountOwnerCapObjectId, typeAndVersion); err != nil {
		return fmt.Errorf("failed to save MCMS AccountOwnerCap object ID %s: %w", output.Objects.McmsAccountOwnerCapObjectId, err)
	}

	typeAndVersion = cldf.NewTypeAndVersion(suistate.SuiMcmsTimelockObjectIDType, suistate.Version1_0_0)
	if err := ab.Save(chainSelector, output.Objects.TimelockObjectId, typeAndVersion); err != nil {
		return fmt.Errorf("failed to save MCMS Timelock object ID %s: %w", output.Objects.TimelockObjectId, err)
	}

	typeAndVersion = cldf.NewTypeAndVersion(suistate.SuiMcmsDeployerObjectIDType, suistate.Version1_0_0)
	if err := ab.Save(chainSelector, output.Objects.McmsDeployerStateObjectId, typeAndVersion); err != nil {
		return fmt.Errorf("failed to save MCMS Deployer object ID %s: %w", output.Objects.McmsDeployerStateObjectId, err)
	}

	return nil
}
