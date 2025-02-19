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
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[cs.SetOCR3OffRampConfig] = SetOCR3ConfigSolana
var _ deployment.ChangeSet[AddRemoteChainToSolanaConfig] = AddRemoteChainToSolana
var _ deployment.ChangeSet[BillingTokenConfig] = AddBillingToken
var _ deployment.ChangeSet[BillingTokenForRemoteChainConfig] = AddBillingTokenForRemoteChain
var _ deployment.ChangeSet[RegisterTokenAdminRegistryConfig] = RegisterTokenAdminRegistry
var _ deployment.ChangeSet[TransferAdminRoleTokenAdminRegistryConfig] = TransferAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[AcceptAdminRoleTokenAdminRegistryConfig] = AcceptAdminRoleTokenAdminRegistry
var _ deployment.ChangeSet[SetFeeAggregatorConfig] = SetFeeAggregator

// HELPER FUNCTIONS
// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName string) (solana.PublicKey, error) {
	tokenPrograms := map[string]solana.PublicKey{
		deployment.SPLTokens:     solana.TokenProgramID, // not used yet
		deployment.SPL2022Tokens: solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: %s, %s", programName, deployment.SPLTokens, deployment.SPL2022Tokens)
	}
	return programID, nil
}

func commonValidation(e deployment.Environment, selector uint64, tokenPubKey solana.PublicKey) error {
	chain, ok := e.SolChains[selector]
	if !ok {
		return fmt.Errorf("chain selector %d not found in environment", selector)
	}
	state, err := cs.LoadOnchainState(e)
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
	for _, token := range chainState.SPL2022Tokens {
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

func validateRouterConfig(chain deployment.SolChain, chainState cs.SolCCIPChainState) error {
	if chainState.Router.IsZero() {
		return fmt.Errorf("router not found in existing state, deploy the router first for chain %d", chain.Selector)
	}
	// addressing errcheck in the next PR
	var routerConfigAccount solRouter.Config
	err := chain.GetAccountDataBorshInto(context.Background(), chainState.RouterConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, initialize the router first %d", chain.Selector)
	}
	return nil
}

func validateFeeQuoterConfig(chain deployment.SolChain, chainState cs.SolCCIPChainState) error {
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

func validateOffRampConfig(chain deployment.SolChain, chainState cs.SolCCIPChainState) error {
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
