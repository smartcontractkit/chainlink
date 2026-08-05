package stellar

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	stellarprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar/provider"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// TestStellarDataFeedsE2E is the Definition-of-Done gate: it stands up a local
// CTF Soroban devnet, deploys both Data Feeds contracts, and invokes EVERY
// exported function on-chain — the 26 cache fns and the 14 proxy fns — through
// the real changesets/operations/bindings (no fakes).
//
// Function-coverage checklist. Outcome (a) = on-chain success; outcome (b) = a
// documented, asserted, expected contract-domain/host error (never ignored).
//
// DataFeedsCache (26):
//   __constructor        (a) DeployCache changeset
//   upgrade              (a) Upgrade changeset
//   recover_tokens       (b) RecoverTokens changeset — SAC transfer(from=self)
//                            re-enters the contract => host trap Error(Context,
//                            InvalidAction), asserted
//   add_feed_admin       (a) AddFeedAdmin changeset
//   remove_feed_admin    (a) RemoveFeedAdmin changeset
//   set_feed_configs     (a) SetFeedConfigs changeset
//   remove_feed_configs  (a) RemoveFeedConfigs changeset
//   transfer_ownership   (a) TransferOwnership changeset (self-transfer)
//   accept_ownership     (a) AcceptOwnership changeset
//   renounce_ownership   (a) RenounceOwnership changeset (LAST — irreversible)
//   version              (a) LoadCacheClient direct
//   type_and_version     (a) LoadCacheClient direct
//   get_owner            (a) LoadCacheClient direct
//   decimals             (a) LoadCacheClient direct — Some(0) while configured,
//                            None after remove_feed_configs
//   description          (a) LoadCacheClient direct — Some while configured,
//                            None after remove_feed_configs
//   is_configured        (a) LoadCacheClient direct — true, false after removal
//   is_frozen            (a) LoadCacheClient direct — false / true while frozen
//   set_feed_frozen      (a) SetFeedFrozen changeset — freeze/unfreeze
//                            round-trip; frozen decimals asserts the
//                            FeedFrozen error (b)
//   is_feed_admin        (a) LoadCacheClient direct
//   has_permission       (a) LoadCacheClient direct
//   get_feed_permissions (a) LoadCacheClient direct
//   on_report            (a) LoadCacheClient direct — valid metadata+XDR report from the
//                            configured AllowedSender; stores a round
//   latest_round         (a) LoadCacheClient direct — returns the stored round
//   get_round            (a) LoadCacheClient direct — returns the stored round
//   round_range          (a) LoadCacheClient direct
//   find_round           (a) LoadCacheClient direct
//
// DataFeedsProxy (14):
//   __constructor        (a) DeployProxy changeset
//   set_cache            (a) SetProxyCache changeset
//   transfer_ownership   (a) TransferOwnership changeset (self-transfer)
//   accept_ownership     (a) AcceptOwnership changeset
//   renounce_ownership   (a) RenounceOwnership changeset (LAST — irreversible)
//   version              (a) LoadProxyClient direct
//   type_and_version     (a) LoadProxyClient direct
//   get_owner            (a) LoadProxyClient direct
//   decimals             (a) LoadProxyClient direct
//   description          (a) LoadProxyClient direct
//   latest_round         (a) LoadProxyClient direct — reads the round stored via on_report
//   get_round            (a) LoadProxyClient direct — reads the round stored via on_report
//   upgrade              (a) Upgrade changeset
//   recover_tokens       (b) RecoverTokens changeset — SAC transfer(from=self)
//                            re-enters the contract => host trap Error(Context,
//                            InvalidAction), asserted

// e2eOnce guards the CTF container lifecycle; the provider requires a sync.Once.
var e2eOnce sync.Once

const (
	e2eQualifier = "e2e"
	e2eVersion   = "1.0.0"
	// e2eDataID is a canonical left-aligned feed id (BTC/USD-shaped). After
	// dataIDsToBytes left-justifies it into [16]byte, byte[7] is 0x00, which is
	// outside the cache's decimals window [0x20,0x60], so cache.decimals returns
	// Some(0) while configured. Non-zero id => a valid feed id (is_valid_id passes).
	e2eDataID       = "0x018e16c39e00032000000"
	e2eDescription  = "BTC/USD"
	e2eWorkflowName = "e2e"
	// e2eWorkflowOwner is the 20-byte workflow owner shared by the feed config
	// permission and the on_report metadata (they must match for the report to
	// be accepted and stored).
	e2eWorkflowOwner = "0x0102030405060708090a0b0c0d0e0f1011121314"
	// e2eReportTS must be > 0 (the initial stored timestamp) so the report is
	// recorded as a new round rather than rejected as stale.
	e2eReportTS = uint64(1_700_000_000)
)

func TestStellarDataFeedsE2E(t *testing.T) {
	skipInCI(t)

	sel := chainsel.STELLAR_LOCALNET.Selector
	ctx := t.Context()

	// --- stand up / attach to the local Soroban devnet ---------------------
	bc := newE2EChain(t, sel)

	e, err := environment.New(ctx,
		environment.WithChains(bc),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)
	env := *e

	ch, ok := env.BlockChains.StellarChains()[sel]
	require.True(t, ok, "stellar chain %d must be present in the environment", sel)
	deployer := ch.Signer.Address()

	// The random deployer key has no balance; friendbot funds it before any fee-
	// bearing operation. Tolerates "already funded".
	fundViaFriendbot(t, ch.FriendbotURL, deployer)

	const cacheWasm = "testdata/data_feeds_cache.wasm"
	const proxyWasm = "testdata/data_feeds_proxy.wasm"

	// --- deploy + configure (changeset pipeline) ---------------------------
	env = apply(t, env, DeployCache{}, &DeployCacheRequest{
		ChainSel: sel, WasmPath: cacheWasm, Qualifier: e2eQualifier, Version: e2eVersion,
	})
	env = apply(t, env, DeployProxy{}, &DeployProxyRequest{
		ChainSel: sel, WasmPath: proxyWasm, CacheQualifier: e2eQualifier,
		Qualifier: e2eQualifier, Version: e2eVersion,
	})
	env = apply(t, env, AddFeedAdmin{}, &FeedAdminRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Admin: deployer,
	})
	env = apply(t, env, SetFeedConfigs{}, &SetFeedConfigsRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Admin: deployer,
		DataIDs: []string{e2eDataID}, Descriptions: []string{e2eDescription},
		Permissions: []FeedPermission{{
			AllowedSender:        deployer,
			AllowedWorkflowOwner: e2eWorkflowOwner,
			AllowedWorkflowName:  e2eWorkflowName,
		}},
	})
	env = apply(t, env, SetProxyCache{}, &SetProxyCacheRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, CacheQualifier: e2eQualifier,
	})

	// Binding-typed feed id and permission tuple for the direct calls below.
	ids, err := dataIDsToBytes([]string{e2eDataID})
	require.NoError(t, err)
	did := [16]byte(ids[0])

	ownerBytes, err := workflowOwnerToBytes(e2eWorkflowOwner)
	require.NoError(t, err)
	nameBytes, err := workflowNameToBytes(e2eWorkflowName)
	require.NoError(t, err)

	// --- cache: on_report (stores a round) + every read -------------------
	cacheClient, cacheRef, err := LoadCacheClient(env, sel, e2eQualifier, e2eVersion)
	require.NoError(t, err)
	cacheAddr := cacheRef.Address

	metadata := buildReportMetadata([10]byte(nameBytes), [20]byte(ownerBytes))
	report := buildReport(t, ids[0], big.NewInt(123456), e2eReportTS)
	require.NoError(t, cacheClient.OnReport(ctx, deployer, metadata, report),
		"on_report with a valid encoding from the configured AllowedSender must succeed")

	// Prove the report was stored; the proxy read-through depends on it.
	stored, err := cacheClient.LatestRound(ctx, did)
	require.NoError(t, err)
	require.NotNil(t, stored, "on_report should have stored a round")
	require.Equalf(t, 0, stored.Answer.Cmp(big.NewInt(123456)),
		"stored answer mismatch: got %v", stored.Answer)
	require.Equal(t, uint64(1), stored.RoundId)

	mustCall(t, "cache.version", func() error { _, e := cacheClient.Version(ctx); return e })
	mustCall(t, "cache.type_and_version", func() error { _, e := cacheClient.TypeAndVersion(ctx); return e })
	mustCall(t, "cache.get_owner", func() error { _, e := cacheClient.GetOwner(ctx); return e })
	mustCall(t, "cache.decimals", func() error {
		dec, e := cacheClient.Decimals(ctx, did)
		if e != nil {
			return e
		}
		require.NotNil(t, dec, "decimals must be present for a configured feed")
		require.Equal(t, uint32(0), *dec)
		return nil
	})
	mustCall(t, "cache.description", func() error {
		desc, e := cacheClient.Description(ctx, did)
		if e != nil {
			return e
		}
		require.NotNil(t, desc, "description must be present for a configured feed")
		require.Equal(t, e2eDescription, *desc)
		return nil
	})
	mustCall(t, "cache.is_configured", func() error {
		configured, e := cacheClient.IsConfigured(ctx, did)
		if e != nil {
			return e
		}
		require.True(t, configured, "feed was configured via SetFeedConfigs")
		return nil
	})
	mustCall(t, "cache.is_frozen", func() error {
		frozen, e := cacheClient.IsFrozen(ctx, did)
		if e != nil {
			return e
		}
		require.False(t, frozen, "feed starts unfrozen")
		return nil
	})
	mustCall(t, "cache.is_feed_admin", func() error { _, e := cacheClient.IsFeedAdmin(ctx, deployer); return e })
	mustCall(t, "cache.has_permission", func() error {
		_, e := cacheClient.HasPermission(ctx, did, deployer, [20]byte(ownerBytes), [10]byte(nameBytes))
		return e
	})
	mustCall(t, "cache.get_feed_permissions", func() error { _, e := cacheClient.GetFeedPermissions(ctx, did); return e })
	mustCall(t, "cache.get_round", func() error {
		r, e := cacheClient.GetRound(ctx, did, 1)
		if e != nil {
			return e
		}
		require.NotNil(t, r, "get_round(did,1) must return the stored round, not void")
		return nil
	})
	mustCall(t, "cache.round_range", func() error {
		rounds, e := cacheClient.RoundRange(ctx, did, 1, 1)
		if e != nil {
			return e
		}
		require.NotEmpty(t, rounds, "round_range should return the stored round")
		return nil
	})
	mustCall(t, "cache.find_round", func() error {
		r, e := cacheClient.FindRound(ctx, did, e2eReportTS, cache.BoundAtOrBefore)
		if e != nil {
			return e
		}
		require.NotNil(t, r, "find_round must return the stored round, not void")
		return nil
	})

	// --- set_feed_frozen: freeze -> reads rejected -> unfreeze -------------
	env = apply(t, env, SetFeedFrozen{}, &SetFeedFrozenRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Admin: deployer,
		DataIDs: []string{e2eDataID}, Frozen: true,
	})
	frozen, err := cacheClient.IsFrozen(ctx, did)
	require.NoError(t, err)
	require.True(t, frozen, "is_frozen must report the freeze")
	_, frozenErr := cacheClient.Decimals(ctx, did)
	require.Error(t, frozenErr, "decimals on a frozen feed must be rejected")
	require.ErrorContains(t, frozenErr, "#109", "expected CacheError::FeedFrozen")
	t.Logf("(b) cache.decimals while frozen error: %v", frozenErr)
	env = apply(t, env, SetFeedFrozen{}, &SetFeedFrozenRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Admin: deployer,
		DataIDs: []string{e2eDataID}, Frozen: false,
	})

	// --- proxy reads (read through to the cache's stored round) -----------
	proxyClient, _, err := LoadProxyClient(env, sel, e2eQualifier, e2eVersion)
	require.NoError(t, err)

	mustCall(t, "proxy.version", func() error { _, e := proxyClient.Version(ctx); return e })
	mustCall(t, "proxy.type_and_version", func() error { _, e := proxyClient.TypeAndVersion(ctx); return e })
	mustCall(t, "proxy.get_owner", func() error { _, e := proxyClient.GetOwner(ctx); return e })
	mustCall(t, "proxy.decimals", func() error { _, e := proxyClient.Decimals(ctx, did); return e })
	mustCall(t, "proxy.description", func() error { _, e := proxyClient.Description(ctx, did); return e })
	mustCall(t, "proxy.latest_round", func() error { _, e := proxyClient.LatestRound(ctx, did); return e })
	mustCall(t, "proxy.get_round", func() error { _, e := proxyClient.GetRound(ctx, did, 1); return e })

	// --- proxy upgrade + recover_tokens ------------------------------------
	env = apply(t, env, Upgrade{}, &UpgradeRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion,
		Contract: ProxyContract, WasmPath: proxyWasm,
	})

	// recover_tokens is exercised with the contract's own address as the token:
	// the SAC transfer sub-call re-enters the contract and Soroban traps with
	// Error(Context, InvalidAction). The call is really invoked on-chain; a
	// funded token balance for a success path is disproportionate here.
	proxyRecoverErr := applyErr(t, env, RecoverTokens{}, &RecoverTokensRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion,
		Contract: ProxyContract, Token: proxyClient.ContractID(), To: deployer, Amount: 1,
	})
	require.Error(t, proxyRecoverErr, "proxy recover_tokens must fail: self-transfer re-enters the contract")
	require.ErrorContains(t, proxyRecoverErr, "InvalidAction", "expected a contract re-entry host trap")
	t.Logf("(b) proxy.recover_tokens error: %v", proxyRecoverErr)

	// --- proxy ownership: transfer -> accept -> renounce (renounce LAST) ---
	proxyLiveUntil := nextLiveUntil(t, ch)
	env = apply(t, env, TransferOwnership{}, &OwnershipRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion,
		Contract: ProxyContract, NewOwner: deployer, LiveUntilLedger: proxyLiveUntil,
	})
	env = apply(t, env, AcceptOwnership{}, &OwnershipRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Contract: ProxyContract,
	})
	env = apply(t, env, RenounceOwnership{}, &OwnershipRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Contract: ProxyContract,
	})

	// --- cache upgrade + recover + config/admin removal --------------------
	env = apply(t, env, Upgrade{}, &UpgradeRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion,
		Contract: CacheContract, WasmPath: cacheWasm,
	})

	// recover_tokens: same re-entry trap as the proxy. Runs while the deployer
	// is still owner; renounce happens last.
	cacheRecoverErr := applyErr(t, env, RecoverTokens{}, &RecoverTokensRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion,
		Contract: CacheContract, Token: cacheAddr, To: deployer, Amount: 1,
	})
	require.Error(t, cacheRecoverErr, "cache recover_tokens must fail: self-transfer re-enters the contract")
	require.ErrorContains(t, cacheRecoverErr, "InvalidAction", "expected a contract re-entry host trap")
	t.Logf("(b) cache.recover_tokens error: %v", cacheRecoverErr)

	env = apply(t, env, RemoveFeedConfigs{}, &RemoveFeedConfigsRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Admin: deployer,
		DataIDs: []string{e2eDataID},
	})

	// The optional reads must decode None once the feed config is gone.
	dec, err := cacheClient.Decimals(ctx, did)
	require.NoError(t, err)
	require.Nil(t, dec, "decimals must be absent after remove_feed_configs")
	desc, err := cacheClient.Description(ctx, did)
	require.NoError(t, err)
	require.Nil(t, desc, "description must be absent after remove_feed_configs")
	configured, err := cacheClient.IsConfigured(ctx, did)
	require.NoError(t, err)
	require.False(t, configured, "is_configured must be false after remove_feed_configs")

	env = apply(t, env, RemoveFeedAdmin{}, &FeedAdminRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Admin: deployer,
	})

	// --- cache ownership: transfer -> accept -> renounce (renounce LAST) ---
	cacheLiveUntil := nextLiveUntil(t, ch)
	env = apply(t, env, TransferOwnership{}, &OwnershipRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion,
		Contract: CacheContract, NewOwner: deployer, LiveUntilLedger: cacheLiveUntil,
	})
	env = apply(t, env, AcceptOwnership{}, &OwnershipRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Contract: CacheContract,
	})
	env = apply(t, env, RenounceOwnership{}, &OwnershipRequest{
		ChainSel: sel, Qualifier: e2eQualifier, Version: e2eVersion, Contract: CacheContract,
	})
}

// skipInCI mirrors the Solana suite convention: container-backed e2e tests do
// not run in CI (they need a local Docker daemon).
func skipInCI(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping container-backed e2e test in CI")
	}
}

// newE2EChain provisions the Stellar devnet chain.
//
// The contracts under test are built with soroban-sdk 26.1.0, i.e. their wasm
// env-meta requires ledger protocol 26. quickstart:latest arms its local
// network at protocol 25 (uploads fail with `HostError: Error(Context,
// InternalError)`), so the CTF path pins stellar/quickstart:future, whose
// local network defaults to protocol 26.
//
// Default path: boot a CTF quickstart container via NewCTFChainProvider.
//
// Override path (STELLAR_E2E_RPC_URL set): attach to an externally started
// quickstart container via NewRPCChainProvider. This exists because the CTF
// framework hardcodes the container platform to linux/amd64
// (chainlink-testing-framework framework/components/blockchain/stellar.go) and
// CLDF's CTFChainProviderConfig does not expose the ImagePlatform knob — so on
// a native linux/arm64 host without amd64 emulation the CTF path fails with
// "exec /start: exec format error" even though the quickstart images ship an
// arm64 variant. On such hosts run the same image natively and point the test
// at it:
//
//	docker run -d -p 8000:8000 stellar/quickstart:latest \
//	  --local --enable-soroban-rpc --protocol-version 26
//	STELLAR_E2E_RPC_URL=http://127.0.0.1:8000/rpc \
//	STELLAR_E2E_FRIENDBOT_URL=http://127.0.0.1:8000/friendbot \
//	go test ./data-feeds/changeset/stellar/ -run TestStellarDataFeedsE2E -v -timeout 25m
func newE2EChain(t *testing.T, sel uint64) cldfchain.BlockChain {
	t.Helper()

	if rpcURL := os.Getenv("STELLAR_E2E_RPC_URL"); rpcURL != "" {
		// Passphrase for the quickstart `--local` standalone network — the same
		// one the CTF path uses (blockchain.DefaultStellarNetworkPassphrase).
		const localPassphrase = "Standalone Network ; February 2017"
		p := stellarprovider.NewRPCChainProvider(sel, stellarprovider.RPCChainProviderConfig{
			SorobanRPCURL:      rpcURL,
			NetworkPassphrase:  localPassphrase,
			FriendbotURL:       os.Getenv("STELLAR_E2E_FRIENDBOT_URL"),
			DeployerKeypairGen: stellarprovider.KeypairRandom(),
		})
		bc, err := p.Initialize(t.Context())
		require.NoError(t, err)
		return bc
	}

	p := stellarprovider.NewCTFChainProvider(t, sel, stellarprovider.CTFChainProviderConfig{
		DeployerKeypairGen: stellarprovider.KeypairRandom(),
		Once:               &e2eOnce,
		// protocol 26 local network; see doc comment above.
		Image: "stellar/quickstart:future",
	})
	bc, err := p.Initialize(t.Context())
	require.NoError(t, err)
	return bc
}

// apply configures and applies a single changeset, requiring success, and
// returns the merged environment (datastore + address book threaded forward).
func apply[C any](t *testing.T, env cldf.Environment, cs cldf.ChangeSetV2[C], cfg C) cldf.Environment {
	t.Helper()
	out, err := commonchangeset.Apply(t, env, commonchangeset.Configure(cs, cfg))
	require.NoError(t, err)
	return out
}

// applyErr configures and applies a single changeset expected to fail, and
// returns the error (the environment is unchanged on failure).
func applyErr[C any](t *testing.T, env cldf.Environment, cs cldf.ChangeSetV2[C], cfg C) error {
	t.Helper()
	_, err := commonchangeset.Apply(t, env, commonchangeset.Configure(cs, cfg))
	return err
}

// mustCall runs a read/query closure and requires it to return no error.
func mustCall(t *testing.T, name string, fn func() error) {
	t.Helper()
	require.NoErrorf(t, fn(), "call %s", name)
}

// nextLiveUntil returns a pending-transfer expiry ledger valid for the current
// chain: strictly greater than the current ledger and comfortably below
// max_live_until_ledger. u32::MAX would trip InvalidLiveUntilLedger.
func nextLiveUntil(t *testing.T, ch cldfstellar.Chain) uint32 {
	t.Helper()
	resp, err := ch.Client.GetLatestLedger(t.Context())
	require.NoError(t, err)
	return resp.Sequence + 100_000
}

// buildReportMetadata lays out the 64-byte on_report metadata header the cache
// decodes: [0:32) workflow_cid | [32:42) workflow_name | [42:62) workflow_owner
// | [62:64) report_id. Only workflow_name and workflow_owner participate in the
// permission hash, so the cid/report_id bytes are left zero.
func buildReportMetadata(name [10]byte, owner [20]byte) []byte {
	md := make([]byte, 64)
	copy(md[32:42], name[:])
	copy(md[42:62], owner[:])
	return md
}

// buildReport builds the on_report `report` payload: the XDR encoding of a
// Soroban `Vec<ReportEntry>` with a single entry. The cache truncates the
// 32-byte wire data id to its high 16 bytes, so wireHi16 is placed in the top
// half and the rest left zero.
func buildReport(t *testing.T, wireHi16 [16]byte, answer *big.Int, ts uint64) []byte {
	t.Helper()
	var wire [32]byte
	copy(wire[0:16], wireHi16[:])

	// #[contracttype] structs serialize as an ScMap with symbol keys sorted
	// alphabetically (answer < data_id < timestamp) — BuildStructScVal sorts.
	entry, err := scval.BuildStructScVal(map[string]xdr.ScVal{
		"data_id":   scval.Bytes32ToScVal(wire),
		"answer":    scval.MustToScVal(scval.I256ToScVal(answer)),
		"timestamp": scval.Uint64ToScVal(ts),
	})
	require.NoError(t, err)

	vec := xdr.ScVec{entry}
	vecPtr := &vec
	reportVal := xdr.ScVal{Type: xdr.ScValTypeScvVec, Vec: &vecPtr}
	b, err := reportVal.MarshalBinary()
	require.NoError(t, err)
	return b
}

// fundViaFriendbot funds addr through the CTF chain's friendbot, tolerating an
// already-funded account. Friendbot may not be ready the instant the container
// starts, so it retries.
func fundViaFriendbot(t *testing.T, friendbotURL, addr string) {
	t.Helper()
	require.NotEmpty(t, friendbotURL, "friendbot URL must be set to fund the deployer")

	target := friendbotURL + "?addr=" + url.QueryEscape(addr)
	var lastErr error
	for i := 0; i < 15; i++ {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, target, nil) //nolint:gosec // local CTF devnet URL
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return
		}
		// Friendbot returns a 400 with op_already_exists when the account is
		// already funded — that's a success for our purposes.
		if bytes.Contains(body, []byte("already")) || bytes.Contains(body, []byte("exists")) {
			return
		}
		lastErr = fmt.Errorf("friendbot status %d: %s", resp.StatusCode, string(body))
		time.Sleep(2 * time.Second)
	}
	// Don't fail here — let the first fee-bearing deploy surface a precise
	// underfunded error if funding genuinely failed.
	t.Logf("warning: friendbot funding did not confirm success: %v", lastErr)
}
