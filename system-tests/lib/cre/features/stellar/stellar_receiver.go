package stellar

// Test-receiver helpers for the Stellar write smoke test. These live in the lib
// (which already depends on chainlink-stellar) so the tests module doesn't import
// the deployer/receiver packages directly. The receiver records on_report calls;
// the write test deploys it, targets it from the workflow, then asserts ReportCount.

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment/cre/stellar"
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

	buildCfg, err := helpers.BuildStellarConfigFor(ctx, stellar.ReceiverWasm)
	if err != nil {
		return "", fmt.Errorf("failed to resolve stellar receiver WASM source: %w", err)
	}

	var salt [32]byte
	return stellar.DeployReceiverForChain(ctx, stellarChain, buildCfg, salt)
}

// ReceiverReportCount reads the receiver's report_count (read-only simulate,
// no funding needed) so the write test can assert a report was delivered on-chain.
func ReceiverReportCount(ctx context.Context, chain *stellchain.Blockchain, contractID string) (uint32, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return 0, err
	}
	return stellar.ReportCountForChain(ctx, stellarChain, contractID)
}

// ReceiverLastValueU64 reads the receiver's last_value_u64 (read-only simulate)
// so the write→read roundtrip test can assert payload integrity.
func ReceiverLastValueU64(ctx context.Context, chain *stellchain.Blockchain, contractID string) (uint64, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return 0, err
	}
	return stellar.ReceiverLastValueU64ForChain(ctx, stellarChain, contractID)
}

// DeployStellarRejectingReceiver deploys the CRE rejecting test receiver (always
// returns Err from on_report) on the given Stellar chain and returns its C-address.
// Used by write regression tests to exercise the forwarder's Failed transmission state.
func DeployStellarRejectingReceiver(ctx context.Context, chain *stellchain.Blockchain) (string, error) {
	stellarChain, err := stellarCldfChain(chain)
	if err != nil {
		return "", err
	}
	owner := stellarChain.Signer.Address()
	if fundErr := chain.Fund(ctx, owner, 0); fundErr != nil {
		return "", fmt.Errorf("failed to fund stellar deployer %s via friendbot: %w", owner, fundErr)
	}

	buildCfg, err := helpers.BuildStellarConfigFor(ctx, stellar.RejectingReceiverWasm)
	if err != nil {
		return "", fmt.Errorf("failed to resolve stellar rejecting receiver WASM source: %w", err)
	}

	var salt [32]byte
	salt[0] = 0x52 // 'R' — distinct from the cooperative receiver's all-zero salt
	return stellar.DeployRejectingReceiverForChain(ctx, stellarChain, buildCfg, salt)
}
