package testhelpers

import (
	"errors"
	"testing"

	solana "github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// DebugLogSolanaFeeTokenPrices is a TEMPORARY diagnostic for the Solana->Sui
// InvalidTokenPrice (8023) investigation. It reads the wSOL and LINK BillingTokenConfig
// PDAs on the Solana fee-quoter (the accounts get_fee reads at fee_for_msg:185 /
// fee_juels->convert:140) and logs, for each: the FeeQuoter program id used to derive the
// PDA, account existence/owner, mint, and usd_per_token {value, timestamp}.
//
// It reuses the exact PDA derivation (solState.FindFqBillingTokenConfigPDA) and account type
// (solFeeQuoter.BillingTokenConfigWrapper) that UpdatePrices.Validate uses, so the addresses
// inspected match the ones seeded by seedSolanaSourcePricesEOA and read by SendRequestSol.
//
// Call at checkpoints (label = e.g. "C1 post-deploy", "C2 post-AddLane", "C3 pre-send") to
// see which token reads {0,0} and where, and whether the FeeQuoter program id is consistent
// across checkpoints. DELETE this file once the root cause is fixed.
func DebugLogSolanaFeeTokenPrices(t *testing.T, e *DeployedEnv, solSel uint64, state stateview.CCIPOnChainState, label string) {
	t.Helper()
	solChainState, ok := state.SolChains[solSel]
	if !ok {
		t.Logf("[PROBE %s] no Solana chain state for selector %d", label, solSel)
		return
	}
	chain := e.Env.BlockChains.SolanaChains()[solSel]
	feeQuoter := solChainState.FeeQuoter
	t.Logf("[PROBE %s] Solana chain=%d FeeQuoter program=%s", label, solSel, feeQuoter.String())

	tokens := []struct {
		name string
		mint solana.PublicKey
	}{
		{"wSOL (fee token)", solChainState.WSOL},
		{"LINK (fee_juels)", solChainState.LinkToken},
	}
	ctx := e.Env.GetContext()
	for _, tk := range tokens {
		pda, _, err := solState.FindFqBillingTokenConfigPDA(tk.mint, feeQuoter)
		if err != nil {
			t.Logf("[PROBE %s] %s: FindFqBillingTokenConfigPDA err=%v (mint=%s fq=%s)",
				label, tk.name, err, tk.mint.String(), feeQuoter.String())
			continue
		}

		var wrapper solFeeQuoter.BillingTokenConfigWrapper
		err = chain.GetAccountDataBorshInto(ctx, pda, &wrapper)
		if err != nil {
			if errors.Is(err, solRpc.ErrNotFound) {
				t.Logf("[PROBE %s] %s: BillingTokenConfig PDA NOT FOUND (mint=%s fq=%s pda=%s) "+
					"-> on-chain get_fee would see a system-owned/uninitialized account",
					label, tk.name, tk.mint.String(), feeQuoter.String(), pda.String())
			} else {
				t.Logf("[PROBE %s] %s: GetAccountDataBorshInto err=%v (pda=%s)",
					label, tk.name, err, pda.String())
			}
			continue
		}

		owner := "<unknown>"
		if info, infoErr := chain.Client.GetAccountInfo(ctx, pda); infoErr == nil && info != nil && info.Value != nil {
			owner = info.Value.Owner.String()
		}
		zeroValue := allZero28(wrapper.Config.UsdPerToken.Value)
		t.Logf("[PROBE %s] %s: owner=%s version=%d enabled=%v mint=%s "+
			"usd_per_token.value=%x timestamp=%d zeroValue=%v (pda=%s)",
			label, tk.name, owner, wrapper.Version, wrapper.Config.Enabled,
			wrapper.Config.Mint.String(), wrapper.Config.UsdPerToken.Value,
			wrapper.Config.UsdPerToken.Timestamp, zeroValue, pda.String())
	}
}

// allZero28 reports whether a [28]uint8 is all zero (i.e. usd_per_token.value == 0, which
// together with timestamp==0 is the get_validated_token_price -> InvalidTokenPrice (8023)
// precondition).
func allZero28(v [28]uint8) bool {
	for _, b := range v {
		if b != 0 {
			return false
		}
	}
	return true
}
