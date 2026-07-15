package solana

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	solstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana"
	pdasol "github.com/smartcontractkit/cld-changesets/pkg/family/solana"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
)

type CCIPSolanaContractVersion string

const (
	SolanaContractV0_1_0 CCIPSolanaContractVersion = "v0.1.0"
	SolanaContractV0_1_1 CCIPSolanaContractVersion = "v0.1.1"
)

var ContractVersionShortSha = map[CCIPSolanaContractVersion]string{
	SolanaContractV0_1_0: "0ee732e80586",
	SolanaContractV0_1_1: "7f8a0f403c3a",
}

func ValidateMCMSConfigSolana(
	e cldf.Environment,
	mcms *cldfproposalutils.TimelockConfig,
	chain cldf_solana.Chain,
	chainState solanastateview.CCIPChainState,
	tokenAddress solana.PublicKey,
	tokenPoolMetadata string,
	contractsToValidate map[cldf.ContractType]bool) error {
	if mcms != nil {
		if err := mcms.ValidateSolana(e, chain.Selector); err != nil {
			return fmt.Errorf("failed to validate MCMS config: %w", err)
		}
	}
	if contractsToValidate[shared.FeeQuoter] {
		if err := solanastateview.ValidateOwnershipSolana(&e, chain, mcms != nil, chainState.FeeQuoter, shared.FeeQuoter, tokenAddress); err != nil {
			return fmt.Errorf("failed to validate ownership for fee quoter: %w", err)
		}
	}
	if contractsToValidate[shared.Router] {
		if err := solanastateview.ValidateOwnershipSolana(&e, chain, mcms != nil, chainState.Router, shared.Router, tokenAddress); err != nil {
			return fmt.Errorf("failed to validate ownership for router: %w", err)
		}
	}
	if contractsToValidate[shared.OffRamp] {
		if err := solanastateview.ValidateOwnershipSolana(&e, chain, mcms != nil, chainState.OffRamp, shared.OffRamp, tokenAddress); err != nil {
			return fmt.Errorf("failed to validate ownership for off ramp: %w", err)
		}
	}
	if contractsToValidate[shared.RMNRemote] {
		if err := solanastateview.ValidateOwnershipSolana(&e, chain, mcms != nil, chainState.RMNRemote, shared.RMNRemote, tokenAddress); err != nil {
			return fmt.Errorf("failed to validate ownership for rmnremote: %w", err)
		}
	}
	if !tokenAddress.IsZero() {
		metadata := shared.CLLMetadata
		if tokenPoolMetadata != "" {
			metadata = tokenPoolMetadata
		}
		if contractsToValidate[shared.BurnMintTokenPool] {
			if err := solanastateview.ValidateOwnershipSolana(&e, chain, mcms != nil, chainState.BurnMintTokenPools[metadata], shared.BurnMintTokenPool, tokenAddress); err != nil {
				return fmt.Errorf("failed to validate ownership for burnmint: %w", err)
			}
		}
		if contractsToValidate[shared.LockReleaseTokenPool] {
			if err := solanastateview.ValidateOwnershipSolana(&e, chain, mcms != nil, chainState.LockReleaseTokenPools[metadata], shared.LockReleaseTokenPool, tokenAddress); err != nil {
				return fmt.Errorf("failed to validate ownership for lockrelease: %w", err)
			}
		}
		if contractsToValidate[shared.CCTPTokenPool] {
			if err := solanastateview.ValidateOwnershipSolana(&e, chain, mcms != nil, chainState.CCTPTokenPool, shared.CCTPTokenPool, tokenAddress); err != nil {
				return fmt.Errorf("failed to validate ownership for cctp token pool: %w", err)
			}
		}
	}

	return nil
}

// buildProposalCommonWithConfig builds a Solana MCMS timelock proposal honoring the MCMS
// action (schedule/bypass/cancel) carried in mcmsCfg. It selects the MCM signer group that
// matches the action (proposer for schedule, bypasser for bypass, canceller for cancel) so
// the resulting proposal's Merkle root is anchored to the correct MCM contract and its op
// count. This selection must happen here, at generation time: the root is cryptographically
// bound to the chosen MCM, so a proposal cannot be reliably retargeted to a different signer
// group after it has been generated. A nil mcmsCfg defaults to a schedule proposal against
// the proposer MCM.
func buildProposalCommonWithConfig(
	e cldf.Environment,
	chainSelector uint64,
	description string,
	mcmsCfg *cldfproposalutils.TimelockConfig,
	batches []mcmsTypes.BatchOperation) (*mcms.TimelockProposal, error) {
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}

	chain := e.BlockChains.SolanaChains()[chainSelector]
	addresses, _ := e.ExistingAddresses.AddressesForChain(chainSelector)
	mcmState, _ := solstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)

	cfg := cldfproposalutils.TimelockConfig{}
	if mcmsCfg != nil {
		cfg = *mcmsCfg
	}

	timelocks[chainSelector] = mcmsSolana.ContractAddress(
		mcmState.TimelockProgram,
		mcmsSolana.PDASeed(mcmState.TimelockSeed),
	)
	// Select the MCM signer group matching the action (proposer for schedule/empty).
	proposers[chainSelector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmSeedForAction(mcmState, cfg.MCMSAction))
	inspectors[chainSelector] = mcmsSolana.NewInspector(chain.Client)

	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		description,
		cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build proposal: %w", err)
	}
	return proposal, nil
}

// mcmSeedForAction returns the MCM PDA seed for the signer group matching the MCMS action:
// bypasser for bypass, canceller for cancel, and proposer otherwise (schedule or empty).
// Callers building a proposal for a specific action must anchor it to the matching MCM,
// because the proposal's Merkle root is bound to the chosen MCM contract.
func mcmSeedForAction(mcmState *solstate.MCMSWithTimelockState, action mcmsTypes.TimelockAction) mcmsSolana.PDASeed {
	switch action {
	case mcmsTypes.TimelockActionBypass:
		return mcmsSolana.PDASeed(mcmState.BypasserMcmSeed)
	case mcmsTypes.TimelockActionCancel:
		return mcmsSolana.PDASeed(mcmState.CancellerMcmSeed)
	default:
		return mcmsSolana.PDASeed(mcmState.ProposerMcmSeed)
	}
}

// BuildProposalsForTxnsWithConfig wraps the given Solana transactions in a single batch and
// builds an MCMS timelock proposal, honoring the MCMS action (schedule/bypass/cancel) in
// mcmsCfg and selecting the matching signer group.
func BuildProposalsForTxnsWithConfig(
	e cldf.Environment,
	chainSelector uint64,
	description string,
	mcmsCfg *cldfproposalutils.TimelockConfig,
	txns []mcmsTypes.Transaction) (*mcms.TimelockProposal, error) {
	batches := []mcmsTypes.BatchOperation{
		{
			ChainSelector: mcmsTypes.ChainSelector(chainSelector),
			Transactions:  txns,
		},
	}
	return buildProposalCommonWithConfig(e, chainSelector, description, mcmsCfg, batches)
}

// BuildProposalsForBatchesWithConfig builds an MCMS timelock proposal from the given batch
// operations, honoring the MCMS action (schedule/bypass/cancel) in mcmsCfg and selecting the
// matching signer group.
func BuildProposalsForBatchesWithConfig(
	e cldf.Environment,
	chainSelector uint64,
	description string,
	mcmsCfg *cldfproposalutils.TimelockConfig,
	batches []mcmsTypes.BatchOperation) (*mcms.TimelockProposal, error) {
	return buildProposalCommonWithConfig(e, chainSelector, description, mcmsCfg, batches)
}

type MCMSTxParams struct {
	Ix           solana.Instruction
	ProgramID    string
	ContractType cldf.ContractType
}

func BuildManyMCMSTxsFrom(input []MCMSTxParams) ([]*mcmsTypes.Transaction, error) {
	mcmsTxs := []*mcmsTypes.Transaction{}
	for _, params := range input {
		tx, err := BuildMCMSTxn(params.Ix, params.ProgramID, params.ContractType)
		if err != nil {
			return []*mcmsTypes.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		mcmsTxs = append(mcmsTxs, tx)
	}
	return mcmsTxs, nil
}

func BuildMCMSTxn(ixn solana.Instruction, programID string, contractType cldf.ContractType) (*mcmsTypes.Transaction, error) {
	data, err := ixn.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to extract data: %w", err)
	}
	for _, account := range ixn.Accounts() {
		if account.IsSigner {
			account.IsSigner = false
		}
	}
	tx, err := mcmsSolana.NewTransaction(
		programID,
		data,
		big.NewInt(0),        // e.g. value
		ixn.Accounts(),       // pass along needed accounts
		string(contractType), // some string identifying the target
		[]string{},           // any relevant metadata
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return &tx, nil
}

func FetchTimelockSigner(e cldf.Environment, chainSelector uint64) (solana.PublicKey, error) {
	addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to load addresses for chain %d: %w", chainSelector, err)
	}
	mcmState, err := solstate.MaybeLoadMCMSWithTimelockChainState(e.BlockChains.SolanaChains()[chainSelector], addresses)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to load mcm state: %w", err)
	}
	timelockSignerPDA := pdasol.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	return timelockSignerPDA, nil
}

// returns the authority for the given ixn based on program ownership
func GetAuthorityForIxn(
	e *cldf.Environment,
	chain cldf_solana.Chain,
	chainState solanastateview.CCIPChainState,
	contractType cldf.ContractType,
	tokenAddress solana.PublicKey, // used for burnmint and lockrelease
	tokenMetadata string, // used for burnmint and lockrelease
) solana.PublicKey {
	timelockSigner, err := FetchTimelockSigner(*e, chain.Selector)
	if err != nil {
		return chain.DeployerKey.PublicKey()
	}
	if solanastateview.IsSolanaProgramOwnedByTimelock(e, chain, chainState, contractType, tokenAddress, tokenMetadata) {
		return timelockSigner
	}
	return chain.DeployerKey.PublicKey()
}

// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName cldf.ContractType) (solana.PublicKey, error) {
	tokenPrograms := map[cldf.ContractType]solana.PublicKey{
		shared.SPLTokens:     solana.TokenProgramID,
		shared.SPL2022Tokens: solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: %s, %s", programName, shared.SPLTokens, shared.SPL2022Tokens)
	}
	return programID, nil
}

func generateProposalIfMCMS(e cldf.Environment, chainSelector uint64, mcmsCfg *cldfproposalutils.TimelockConfig, mcmsTxs []mcmsTypes.Transaction) (cldf.ChangesetOutput, error) {
	if len(mcmsTxs) > 0 {
		if mcmsCfg == nil {
			return cldf.ChangesetOutput{}, errors.New("MCMS txn detected but no MCMS config provided. Please re-run with mcms specified")
		}
		proposal, err := BuildProposalsForTxnsWithConfig(
			e, chainSelector, "proposal to upgrade CCIP contracts", mcmsCfg, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

type ExecuteConfig struct {
	ChainSelector uint64
	MCMS          *cldfproposalutils.TimelockConfig
	Chain         cldf_solana.Chain
}

func ExecuteInstructionsAndBuildProposals(e cldf.Environment, cfg ExecuteConfig, instructions [][]solana.Instruction, mcmsTxs []mcmsTypes.Transaction) (cldf.ChangesetOutput, error) {
	for _, instructionSet := range instructions {
		if err := cfg.Chain.Confirm(instructionSet); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxnsWithConfig(
			e, cfg.ChainSelector, "proposal in Solana", cfg.MCMS, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	return cldf.ChangesetOutput{}, nil
}
