package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

var _ deployment.ChangeSet[v1_6.SetOCR3OffRampConfig] = SetOCR3ConfigSolana
var _ deployment.ChangeSet[AddRemoteChainToRouterConfig] = AddRemoteChainToRouter
var _ deployment.ChangeSet[AddRemoteChainToOffRampConfig] = AddRemoteChainToOffRamp
var _ deployment.ChangeSet[AddRemoteChainToFeeQuoterConfig] = AddRemoteChainToFeeQuoter
var _ deployment.ChangeSet[DisableRemoteChainConfig] = DisableRemoteChain
var _ deployment.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset
var _ deployment.ChangeSet[BillingTokenForRemoteChainConfig] = AddBillingTokenForRemoteChain
var _ deployment.ChangeSet[RegisterTokenAdminRegistryConfig] = RegisterTokenAdminRegistry
var _ deployment.ChangeSet[TransferAdminRoleTokenAdminRegistryConfig] = TransferAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[AcceptAdminRoleTokenAdminRegistryConfig] = AcceptAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[SetFeeAggregatorConfig] = SetFeeAggregator
var _ deployment.ChangeSet[BillingTokenConfig] = AddBillingTokenChangeset
var _ deployment.ChangeSet[BillingTokenForRemoteChainConfig] = AddBillingTokenForRemoteChain
var _ deployment.ChangeSet[DeployTestRouterConfig] = DeployTestRouter
var _ deployment.ChangeSet[OffRampRefAddressesConfig] = UpdateOffRampRefAddresses

type MCMSConfigSolana struct {
	MCMS *ccipChangeset.MCMSConfig
	// Public key of program authorities. Depending on when this changeset is called, some may be under
	// the control of the deployer, and some may be under the control of the timelock. (e.g. during new offramp deploy)
	RouterOwnedByTimelock    bool
	FeeQuoterOwnedByTimelock bool
	OffRampOwnedByTimelock   bool
	// Assumes whatever token pool we're operating on
	TokenPoolPDAOwnedByTimelock bool
}

// HELPER FUNCTIONS
// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName deployment.ContractType) (solana.PublicKey, error) {
	tokenPrograms := map[deployment.ContractType]solana.PublicKey{
		ccipChangeset.SPLTokens:     solana.TokenProgramID,
		ccipChangeset.SPL2022Tokens: solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: %s, %s", programName, ccipChangeset.SPLTokens, ccipChangeset.SPL2022Tokens)
	}
	return programID, nil
}

func commonValidation(e deployment.Environment, selector uint64, tokenPubKey solana.PublicKey) error {
	chain, ok := e.SolChains[selector]
	if !ok {
		return fmt.Errorf("chain selector %d not found in environment", selector)
	}
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if tokenPubKey.Equals(chainState.LinkToken) || tokenPubKey.Equals(chainState.WSOL) {
		return nil
	}
	exists := false
	allTokens := chainState.SPL2022Tokens
	allTokens = append(allTokens, chainState.SPLTokens...)
	for _, token := range allTokens {
		if token.Equals(tokenPubKey) {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("token %s not found in existing state, deploy the token first", tokenPubKey.String())
	}
	return nil
}

func validateRouterConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState, testRouter bool) error {
	_, routerConfigPDA, err := chainState.GetRouterInfo(testRouter)
	if err != nil {
		return err
	}
	var routerConfigAccount solRouter.Config
	err = chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, initialize the router first %d", chain.Selector)
	}
	return nil
}

func validateFeeQuoterConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.FeeQuoter.IsZero() {
		return fmt.Errorf("fee quoter not found in existing state, deploy the fee quoter first for chain %d", chain.Selector)
	}
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(chainState.FeeQuoter)
	err := chain.GetAccountDataBorshInto(context.Background(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		return fmt.Errorf("fee quoter config not found in existing state, initialize the fee quoter first %d", chain.Selector)
	}
	return nil
}

func validateOffRampConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.OffRamp.IsZero() {
		return fmt.Errorf("offramp not found in existing state, deploy the offramp first for chain %d", chain.Selector)
	}
	var offRampConfig solOffRamp.Config
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(chainState.OffRamp)
	err := chain.GetAccountDataBorshInto(context.Background(), offRampConfigPDA, &offRampConfig)
	if err != nil {
		return fmt.Errorf("offramp config not found in existing state, initialize the offramp first %d", chain.Selector)
	}
	return nil
}

// The user is not required to provide all the addresses, only the ones they want to update
type OffRampRefAddressesConfig struct {
	ChainSelector      uint64
	Router             solana.PublicKey
	FeeQuoter          solana.PublicKey
	AddressLookupTable solana.PublicKey
	MCMSSolana         *MCMSConfigSolana
}

func (cfg OffRampRefAddressesConfig) Validate(e deployment.Environment) error {
	chain := e.SolChains[cfg.ChainSelector]
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if err := ValidateMCMSConfigSolana(e, cfg.ChainSelector, cfg.MCMSSolana); err != nil {
		return err
	}
	offRampUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.OffRampOwnedByTimelock
	if err := ccipChangeset.ValidateOwnershipSolana(&e, chain, offRampUsingMCMS, chainState.OffRamp, ccipChangeset.OffRamp); err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}
	return nil
}

func UpdateOffRampRefAddresses(
	e deployment.Environment,
	config OffRampRefAddressesConfig,
) (deployment.ChangesetOutput, error) {
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	chain := e.SolChains[config.ChainSelector]
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if chainState.OffRamp.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get offramp address for chain %s", chain.String())
	}

	var referenceAddressesAccount solOffRamp.ReferenceAddresses
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(chainState.OffRamp)
	if err = chain.GetAccountDataBorshInto(e.GetContext(), offRampReferenceAddressesPDA, &referenceAddressesAccount); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get offramp reference addresses: %w", err)
	}
	routerToSet := referenceAddressesAccount.Router
	if !config.Router.IsZero() {
		e.Logger.Infof("setting router on offramp to %s", config.Router.String())
		routerToSet = config.Router
	}
	feeQuoterToSet := referenceAddressesAccount.FeeQuoter
	if !config.FeeQuoter.IsZero() {
		e.Logger.Infof("setting fee quoter on offramp to %s", config.FeeQuoter.String())
		feeQuoterToSet = config.FeeQuoter
	}
	addressLookupTableToSet := referenceAddressesAccount.OfframpLookupTable
	if !config.AddressLookupTable.IsZero() {
		e.Logger.Infof("setting address lookup table on offramp to %s", config.AddressLookupTable.String())
		addressLookupTableToSet = config.AddressLookupTable
	}

	offRampUsingMCMS := config.MCMSSolana != nil && config.MCMSSolana.OffRampOwnedByTimelock
	timelockSigner, err := FetchTimelockSigner(e, chain.Selector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}
	var authority solana.PublicKey
	if offRampUsingMCMS {
		authority = timelockSigner
	} else {
		authority = chain.DeployerKey.PublicKey()
	}

	solOffRamp.SetProgramID(chainState.OffRamp)
	ix, err := solOffRamp.NewUpdateReferenceAddressesInstruction(
		routerToSet,
		feeQuoterToSet,
		addressLookupTableToSet,
		chainState.OffRampConfigPDA,
		offRampReferenceAddressesPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}
