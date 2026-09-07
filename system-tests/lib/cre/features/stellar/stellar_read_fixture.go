package stellar

import (
	"context"
	"fmt"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"

	"github.com/smartcontractkit/chainlink/deployment/cre/stellar"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
)

// DeployStellarReadFixture deploys the CRE read fixture contract and returns its
// C-address. Used by ReadContract smoke cases (constants, echoes, intentional trap).
func DeployStellarReadFixture(ctx context.Context, chain *stellchain.Blockchain) (string, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return "", err
	}
	owner := stellarChain.Signer.Address()
	if fundErr := chain.Fund(ctx, owner, 0); fundErr != nil {
		return "", fmt.Errorf("failed to fund stellar deployer %s via friendbot: %w", owner, fundErr)
	}

	var salt [32]byte
	// Distinct salt from the receiver deploy (all-zero) when both run in one suite.
	salt[0] = 0x52 // 'R' — distinct from the all-zero receiver salt
	return stellar.DeployReadFixtureForChain(ctx, stellarChain, salt)
}

// stellarCldfChain builds the cldf stellar chain from the environment blockchain.
// NOTE: ToCldfChain() mints a fresh random deployer keypair per call, so callers
// that need a single funded deployer for multiple ops must not re-derive it.
func stellarCldfChain(chain *stellchain.Blockchain) (cldfstellar.Chain, error) {
	cldfChain, err := chain.ToCldfChain()
	if err != nil {
		return cldfstellar.Chain{}, fmt.Errorf("failed to build cldf stellar chain (selector %d): %w", chain.ChainSelector(), err)
	}
	// RPCChainProvider.Initialize returns *stellar.Chain (pointer); deref to the
	// value type callers and NewDeployerFromChain expect.
	sc, ok := cldfChain.(*cldfstellar.Chain)
	if !ok || sc == nil {
		return cldfstellar.Chain{}, fmt.Errorf("expected cldf stellar chain, got %T", cldfChain)
	}
	return *sc, nil
}
