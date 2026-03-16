package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
)

var _ cldf.ChangeSetV2[config.DeployCurseMCMSConfig] = DeployCurseMCMS{}

// DeployCurseMCMS deploys and configures the CurseMCMS contract on Aptos chains.
type DeployCurseMCMS struct{}

func (cs DeployCurseMCMS) VerifyPreconditions(env cldf.Environment, cfg config.DeployCurseMCMSConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	state, err := aptosstate.LoadOnchainStateAptos(env)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	aptosChains := env.BlockChains.AptosChains()
	var errs []error
	for chainSel := range cfg.CurseMCMSConfigPerChain {
		if _, ok := aptosChains[chainSel]; !ok {
			errs = append(errs, fmt.Errorf("aptos chain %d not found in env", chainSel))
			continue
		}
		chainState, ok := state[chainSel]
		if !ok {
			errs = append(errs, fmt.Errorf("aptos chain %d not found in state", chainSel))
			continue
		}
		if chainState.MCMSAddress == (aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("MCMS not deployed for Aptos chain %d", chainSel))
		}
		if chainState.CCIPAddress == (aptos.AccountAddress{}) {
			errs = append(errs, fmt.Errorf("CCIP not deployed for Aptos chain %d", chainSel))
		}
	}

	return errors.Join(errs...)
}

func (cs DeployCurseMCMS) Apply(env cldf.Environment, cfg config.DeployCurseMCMSConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)
	proposals := make([]mcms.TimelockProposal, 0)

	aptosChains := env.BlockChains.AptosChains()
	for chainSel, curseMCMSConfig := range cfg.CurseMCMSConfigPerChain {
		aptosChain := aptosChains[chainSel]
		chainState := state.AptosChains[chainSel]

		deps := dependency.AptosDeps{
			AB:               ab,
			AptosChain:       aptosChain,
			CCIPOnChainState: state,
		}

		seqInput := seq.DeployCurseMCMSSeqInput{
			MCMSAddress: chainState.MCMSAddress,
			CCIPAddress: chainState.CCIPAddress,
			CurseMCMS:   curseMCMSConfig,
		}

		curseMCMSSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployCurseMCMSSequence, deps, seqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CurseMCMS for Aptos chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, curseMCMSSeqReport.ExecutionReports...)

		// Skip address saving and proposal generation when sequence was a no-op (already deployed).
		if curseMCMSSeqReport.Output.CurseMCMSAddress == (aptos.AccountAddress{}) {
			continue
		}

		typeAndVersion := cldf.NewTypeAndVersion(shared.AptosCurseMCMSType, deployment.Version1_6_0)
		if err := ab.Save(chainSel, curseMCMSSeqReport.Output.CurseMCMSAddress.StringLong(), typeAndVersion); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CurseMCMS address for Aptos chain %d: %w", chainSel, err)
		}

		// Generate a CurseMCMS proposal for self-governance operations (AcceptOwnership + SetMinDelay).
		proposal, err := utils.GenerateCurseMCMSProposal(
			env,
			curseMCMSSeqReport.Output.CurseMCMSAddress,
			chainSel,
			[]mcmstypes.BatchOperation{curseMCMSSeqReport.Output.CurseMCMSOperation},
			"CurseMCMS accept ownership and set timelock min delay",
			cfg.MCMSTimelockConfigPerChain[chainSel],
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate CurseMCMS proposal for Aptos chain %d: %w", chainSel, err)
		}
		proposals = append(proposals, *proposal)
	}

	ds, err := shared.PopulateDataStore(ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate DataStore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		MCMSTimelockProposals: proposals,
		Reports:               seqReports,
	}, nil
}
