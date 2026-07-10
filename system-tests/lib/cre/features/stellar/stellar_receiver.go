package stellar

// Test-receiver helpers for the Stellar write smoke test. These live in the lib
// (which already depends on chainlink-stellar) so the tests module doesn't import
// the deployer/receiver packages directly. The receiver records on_report calls;
// the write test deploys it, targets it from the workflow, then asserts ReportCount.

import (
	"context"
	"fmt"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"

	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	stellarreceiver "github.com/smartcontractkit/chainlink-stellar/deployment/cre/receiver"

	"github.com/smartcontractkit/chainlink/deployment/helpers"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
)

// DeployStellarTestReceiver deploys the CRE test receiver (on_report recorder) on
// the given Stellar chain and returns its C-address. The workflow's WriteReport
// targets this address; the forwarder dispatches on_report to it.
func DeployStellarTestReceiver(ctx context.Context, chain *stellchain.Blockchain) (string, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return "", err
	}
	owner := stellarChain.Signer.Address()
	if fundErr := chain.Fund(ctx, owner, 0); fundErr != nil {
		return "", fmt.Errorf("failed to fund stellar deployer %s via friendbot: %w", owner, fundErr)
	}
	deployer, err := stellardeployment.NewDeployerFromChain(stellarChain)
	if err != nil {
		return "", fmt.Errorf("failed to build stellar deployer: %w", err)
	}

	// Source the receiver WASM at deploy time (build or download) — no committed binary.
	buildCfg, err := stellarBuildConfig(helpers.ReceiverWASMFile)
	if err != nil {
		return "", fmt.Errorf("failed to resolve stellar receiver WASM source: %w", err)
	}
	wasm, err := helpers.BuildStellar(ctx, buildCfg)
	if err != nil {
		return "", fmt.Errorf("failed to source stellar receiver WASM: %w", err)
	}

	var salt [32]byte
	return stellarreceiver.DeployReceiver(ctx, deployer, wasm, salt)
}

// StellarReceiverReportCount reads the receiver's report_count (read-only simulate,
// no funding needed) so the write test can assert a report was delivered on-chain.
func StellarReceiverReportCount(ctx context.Context, chain *stellchain.Blockchain, contractID string) (uint32, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return 0, err
	}
	deployer, err := stellardeployment.NewDeployerFromChain(stellarChain)
	if err != nil {
		return 0, fmt.Errorf("failed to build stellar deployer: %w", err)
	}
	return stellarreceiver.ReportCount(ctx, deployer, contractID)
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
