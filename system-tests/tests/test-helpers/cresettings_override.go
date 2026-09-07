package helpers

// Runtime CRE-settings overrides for Local CRE e2e tests.
//
// A node applies CRE settings live — without a restart — when it receives a job of
// type `cresettings`: the node's delegate calls loop.AtomicSettings.Store, which
// hot-swaps the in-memory settings getter (see core/services/cresettings/delegate.go
// and chainlink-common/pkg/loop/settings.go). This is the same mechanism prod uses to
// roll out limit changes (the durable-pipeline `cre-settings` job proposed to
// all-nodes). Here we reuse the deployment changeset to propose that job to every node
// of a DON at test time, and restore the pre-test baseline on cleanup.
//
// Design notes / gotchas (see the usage guide in
// core/scripts/cre/environment/docs/cresettings-override-README.md):
//   - Each Store fully REPLACES the settings getter (updates are not cumulative), so
//     the delivered document must be the DON's boot baseline MERGED with the overrides.
//     A bare-diff document would silently drop the DON's boot CL_CRE_SETTINGS.
//   - Deleting the settings job does NOT revert the settings, so cleanup must re-apply
//     the captured baseline explicitly.
//   - The document must be TOML (JSON is rejected), scoped ([global]/[org.<id>]/...),
//     with all values quoted strings.
//   - Overrides must reach ALL nodes of an OCR DON; a partial rollout makes nodes
//     disagree. Approve() below only returns once every targeted node accepted the job.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/moby/moby/client"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	// creSettingsUpdateLogMarker is what the node's cresettings delegate logs when it
	// applies an update (core/services/cresettings/delegate.go: `Updated settings`).
	// We scan for it plus the doc hash as best-effort proof that a node converged.
	creSettingsUpdateLogMarker = "Updated settings"

	// how long to wait (best effort) for every targeted node to log the applied hash.
	creSettingsConvergenceTimeout = 60 * time.Second
)

// Only one CRE settings override may be active at a time. Overrides mutate settings on
// the SHARED Local CRE environment, so a second active override — whether from another
// test running concurrently or from the same test applying again before reverting —
// would clobber the first one's document. These package-level vars enforce that and fail
// fast with an actionable message. To change settings within one test, call
// handle.Reset(t) (which releases the claim) before applying again; to apply several
// scopes together, compose them in one ApplyCRESettings call.
var (
	creSettingsActiveMu    sync.Mutex
	creSettingsActiveOwner string // t.Name() of the test currently holding an override; "" if none
)

func claimCRESettingsOverride(owner string) error {
	creSettingsActiveMu.Lock()
	defer creSettingsActiveMu.Unlock()
	switch creSettingsActiveOwner {
	case "":
		creSettingsActiveOwner = owner
		return nil
	case owner:
		// Same test applying again without reverting first. Each ApplyCRESettings rebuilds
		// from the boot baseline and REPLACES the previous override, so a silent second apply
		// would drop the first. Make it a loud error instead.
		return errors.New(
			"this test already has an active CRE settings override; call handle.Reset(t) before " +
				"applying again — each ApplyCRESettings replaces the previous override rather than " +
				"stacking on it. To apply several scopes together, compose them in one " +
				"ApplyCRESettings call (e.g. Global(...), Org(id, ...), Workflow(id, ...)).\n" +
				"See core/scripts/cre/environment/docs/cresettings-override-README.md")
	default:
		return fmt.Errorf(
			"a CRE settings override from %q is still active on the shared environment.\n"+
				"Settings-override tests mutate shared state and must run serially: remove t.Parallel() "+
				"from this test (and do not add it to the CRE_TEST_PARALLEL_ENABLED set).\n"+
				"See core/scripts/cre/environment/docs/cresettings-override-README.md",
			creSettingsActiveOwner)
	}
}

func releaseCRESettingsOverride(owner string) {
	creSettingsActiveMu.Lock()
	defer creSettingsActiveMu.Unlock()
	if creSettingsActiveOwner == owner {
		creSettingsActiveOwner = ""
	}
}

// CRESettingsHandle tracks an applied override so a test can revert it. ApplyCRESettings
// already registers a t.Cleanup that reverts automatically; Reset lets a test revert
// early (e.g. to assert post-revert behaviour within the test body).
type CRESettingsHandle struct {
	env      *ttypes.TestEnvironment
	owner    string // t.Name() of the owning test; released back to the guard on restore
	targets  []creSettingsTarget
	reverted bool
}

type creSettingsTarget struct {
	don          *cre.Don
	baselineTOML string
	baselineHash string
	appliedTOML  string
	appliedHash  string
}

// TOMLFile wraps a TOML settings fragment as an *fstest.MapFile. It is used internally by
// the scope options below; exported for building a settings tree by hand if ever needed.
func TOMLFile(content string) *fstest.MapFile { return &fstest.MapFile{Data: []byte(content)} }

// Option contributes one scope's settings to an ApplyCRESettings call. Compose several in a
// single call to apply multiple scopes together (see ApplyCRESettings).
type Option func(fstest.MapFS) error

// Global applies the TOML fragment to the [global] scope (every org/owner/workflow).
func Global(toml string) Option {
	return func(files fstest.MapFS) error { files["global.toml"] = TOMLFile(toml); return nil }
}

// Org applies the TOML fragment to the [org.<id>] scope. id is the org id (no 0x prefix).
func Org(id, toml string) Option { return scopedOption("org", id, toml) }

// Owner applies the TOML fragment to the [owner.<id>] scope. id is the workflow-owner hex (no 0x prefix).
func Owner(id, toml string) Option { return scopedOption("owner", id, toml) }

// Workflow applies the TOML fragment to the [workflow.<id>] scope. id is the workflow id hex (no 0x prefix).
func Workflow(id, toml string) Option { return scopedOption("workflow", id, toml) }

func scopedOption(scope, id, toml string) Option {
	return func(files fstest.MapFS) error {
		if id == "" {
			return fmt.Errorf("%s scope requires a non-empty id", scope)
		}
		files[scope+"/"+id+".toml"] = TOMLFile(toml)
		return nil
	}
}

// FromFS applies a whole tree of settings files laid out like prod (global.toml,
// org/<id>.toml, owner/<id>.toml, workflow/<id>.toml) — for on-disk fixtures via
// os.DirFS("testdata/...") or an embed.FS. It composes with the other options.
func FromFS(src fs.FS) Option {
	return func(files fstest.MapFS) error {
		return fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			data, readErr := fs.ReadFile(src, path)
			if readErr != nil {
				return readErr
			}
			files[path] = &fstest.MapFile{Data: data}
			return nil
		})
	}
}

// ApplyCRESettings overrides CRE settings across the whole environment at runtime, without
// restarting the topology, and registers a t.Cleanup that restores the pre-test baseline
// when the test finishes.
//
// Compose one or more scope options in a SINGLE call — Global, Org, Owner, Workflow (or
// FromFS for an on-disk tree):
//
//	ApplyCRESettings(t, env,
//	    Global("[PerWorkflow.HTTPAction]\nCallLimit = '9'"),
//	    Workflow(wfID, "[PerWorkflow]\nExecutionConcurrencyLimit = '1'"),
//	)
//
// Each option's string is that scope's settings TOML with no scope prefix (the scope and id
// come from the option). Internally the options build a prod-style settings tree that is
// combined by the deployment CombineCRESettingsFiles helper and merged onto each DON's
// baseline.
//
// IMPORTANT: pass every scope you want in ONE call. Calling ApplyCRESettings again REPLACES
// the previous settings (each call rebuilds from the boot baseline) rather than adding to
// them — use Reset then apply again only when you deliberately want to swap.
//
// Settings are applied to every DON that has worker (plugin) nodes — the workflow and
// capabilities DONs — since CRE settings are enforced there and the delivery changeset
// targets type=plugin nodes. Bootstrap-only DONs (e.g. bootstrap-gateway) are skipped.
//
// It fails the test if an option, combining, or proposing fails; application is confirmed
// best-effort via the nodes' "Updated settings" logs.
func ApplyCRESettings(t *testing.T, env *ttypes.TestEnvironment, opts ...Option) *CRESettingsHandle {
	t.Helper()

	require.NotNil(t, env, "test environment must not be nil")
	require.NotNil(t, env.CreEnvironment, "CreEnvironment must not be nil")
	require.NotNil(t, env.CreEnvironment.CldfEnvironment, "CldfEnvironment must not be nil")
	require.NotNil(t, env.Dons, "Dons must not be nil")
	require.NotEmpty(t, opts, "at least one settings option is required (e.g. helpers.Global(...))")

	// Build the prod-style settings tree from the options, then combine it exactly like prod
	// and layer the result onto each DON's baseline below.
	files := fstest.MapFS{}
	for _, opt := range opts {
		require.NoError(t, opt(files), "invalid CRE settings option")
	}
	combined, combineErr := cre_jobs.CombineCRESettingsFiles(t.TempDir(), files)
	require.NoError(t, combineErr, "failed to combine CRE settings files")
	overrideDoc := map[string]any{}
	require.NoError(t, toml.Unmarshal(combined, &overrideDoc), "combined CRE settings is not valid TOML")
	require.NotEmpty(t, overrideDoc, "no settings provided")

	// CRE settings are enforced on worker (plugin) nodes, and the delivery changeset filters
	// proposals to type=plugin — so bootstrap-only DONs (e.g. bootstrap-gateway) have no
	// matching nodes. Deliver only to DONs that have worker nodes.
	require.NotEmpty(t, env.Dons.List(), "no DONs found in the environment")
	targets := make([]*cre.Don, 0)
	for _, don := range env.Dons.List() {
		if don.WorkersCount() > 0 {
			targets = append(targets, don)
		} else {
			t.Logf("[cresettings] skipping DON %q (no worker nodes)", don.Name)
		}
	}
	require.NotEmpty(t, targets, "no DONs with worker nodes found in the environment")

	// Claim the single-active-override slot before touching anything. Fails fast (with an
	// actionable message) if another test's override is still active on the shared env.
	owner := t.Name()
	if err := claimCRESettingsOverride(owner); err != nil {
		require.FailNow(t, err.Error())
	}

	h := &CRESettingsHandle{env: env, owner: owner}
	// Register cleanup up front, so the override is reverted and the slot released even if
	// a delivery below fails partway through.
	t.Cleanup(func() { h.restore(t, false /* not fatal: the test already finished */) })

	for _, don := range targets {
		baselineJSON := bootSettingsForDON(env, don.Name)
		baselineTOML, baselineHash := renderSettings(t, baselineJSON, nil)
		appliedTOML, appliedHash := renderSettings(t, baselineJSON, overrideDoc)

		t.Logf("[cresettings] DON %q: applying override (hash %s) over baseline (hash %s)",
			don.Name, shortHash(appliedHash), shortHash(baselineHash))

		err := deliverCRESettings(env, don, appliedTOML)
		require.NoErrorf(t, err, "failed to deliver CRE settings override to DON %q", don.Name)

		h.targets = append(h.targets, creSettingsTarget{
			don:          don,
			baselineTOML: baselineTOML,
			baselineHash: baselineHash,
			appliedTOML:  appliedTOML,
			appliedHash:  appliedHash,
		})
	}

	// Best-effort confirmation that every node actually logged the applied settings.
	for _, tg := range h.targets {
		logSettingsConvergence(t, tg.don, tg.appliedHash, creSettingsConvergenceTimeout)
	}

	return h
}

// Reset restores the baseline settings on all targeted DONs immediately and waits
// (best effort) for convergence. Safe to call multiple times; the automatic cleanup
// becomes a no-op afterwards.
func (h *CRESettingsHandle) Reset(t *testing.T) {
	t.Helper()
	h.restore(t, true /* fatal: called from the test body */)
}

// AppliedTOML returns the settings document applied to the named DON (its boot
// baseline merged with the overrides). Empty if the DON was not targeted. Useful for
// logging exactly what a test applied at each scope.
func (h *CRESettingsHandle) AppliedTOML(donName string) string {
	for _, tg := range h.targets {
		if tg.don.Name == donName {
			return tg.appliedTOML
		}
	}
	return ""
}

// BaselineTOML returns the pre-test baseline document captured for the named DON, i.e.
// what the settings are restored to on cleanup. Empty if the DON was not targeted.
func (h *CRESettingsHandle) BaselineTOML(donName string) string {
	for _, tg := range h.targets {
		if tg.don.Name == donName {
			return tg.baselineTOML
		}
	}
	return ""
}

func (h *CRESettingsHandle) restore(t *testing.T, fatal bool) {
	t.Helper()
	if h.reverted {
		return
	}
	h.reverted = true
	// Release the single-active-override slot once we've reverted, so the next test can
	// claim it. Runs exactly once (guarded by h.reverted above).
	defer releaseCRESettingsOverride(h.owner)

	for _, tg := range h.targets {
		t.Logf("[cresettings] DON %q: reverting to baseline (hash %s)", tg.don.Name, shortHash(tg.baselineHash))
		err := deliverCRESettings(h.env, tg.don, tg.baselineTOML)
		if err != nil {
			if fatal {
				require.NoErrorf(t, err, "failed to revert CRE settings on DON %q", tg.don.Name)
			} else {
				t.Errorf("[cresettings] failed to revert CRE settings on DON %q: %v", tg.don.Name, err)
			}
			continue
		}
		logSettingsConvergence(t, tg.don, tg.baselineHash, creSettingsConvergenceTimeout)
	}
}

// deliverCRESettings proposes a `cresettings` job carrying settingsTOML to the DON's
// worker nodes. CRE nodes auto-approve the settings job and apply it live, so proposing
// is sufficient — we deliberately do NOT cancel the previous job or explicitly approve
// the new one. With the fixed settings-job UUID, an explicit cancel/approve corrupts the
// JD proposal history on repeated deliveries (e.g. reverting to the same baseline twice,
// which failed with "no job proposal found"). Application is confirmed best-effort via
// the nodes' "Updated settings" logs (see logSettingsConvergence).
func deliverCRESettings(env *ttypes.TestEnvironment, don *cre.Don, settingsTOML string) error {
	input := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: env.CreEnvironment.CldfEnvironment.Name,
		DONName:     don.Name,
		JobName:     "cre-settings",
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: "cre-settings-override"},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.CRESettings,
		Inputs:   job_types.JobSpecInput{"settings": settingsTOML},
	}

	if err := (cre_jobs.ProposeJobSpec{}).VerifyPreconditions(*env.CreEnvironment.CldfEnvironment, input); err != nil {
		return fmt.Errorf("verify settings job preconditions: %w", err)
	}
	if _, err := (cre_jobs.ProposeJobSpec{}).Apply(*env.CreEnvironment.CldfEnvironment, input); err != nil {
		return fmt.Errorf("propose settings job: %w", err)
	}
	return nil
}

// renderSettings merges overrideDoc (may be nil) onto the DON's boot baseline (a
// CL_CRE_SETTINGS JSON string, may be empty) and returns the resulting settings TOML
// plus its sha256 hash (the same hash the node computes and logs).
func renderSettings(t *testing.T, baselineJSON string, overrideDoc map[string]any) (string, string) {
	t.Helper()
	base := map[string]any{}
	if strings.TrimSpace(baselineJSON) != "" {
		require.NoErrorf(t, json.Unmarshal([]byte(baselineJSON), &base),
			"invalid baseline CL_CRE_SETTINGS json: %s", baselineJSON)
	}
	if overrideDoc != nil {
		deepMergeInto(base, overrideDoc)
	}
	if len(base) == 0 {
		// Empty document => node resets to compiled defaults.
		return "", hashString("")
	}
	b, err := toml.Marshal(base)
	require.NoError(t, err, "failed to marshal settings toml")
	s := string(b)
	return s, hashString(s)
}

// deepMergeInto recursively merges src into dst (src wins on leaves). Both are nested
// map[string]any documents; overlapping sub-tables are merged, not replaced.
func deepMergeInto(dst, src map[string]any) {
	for k, sv := range src {
		if sm, ok := sv.(map[string]any); ok {
			if dm, ok := dst[k].(map[string]any); ok {
				deepMergeInto(dm, sm)
				continue
			}
		}
		dst[k] = sv
	}
}

func bootSettingsForDON(env *ttypes.TestEnvironment, donName string) string {
	if env.Config == nil {
		return ""
	}
	for _, nodeSet := range env.Config.NodeSets {
		if nodeSet != nil && nodeSet.Input != nil && nodeSet.Name == donName {
			return nodeSet.EnvVars["CL_CRE_SETTINGS"]
		}
	}
	return ""
}

// logSettingsConvergence polls container logs (best effort) and reports how many nodes
// of the DON have logged the given settings hash. It never fails the test — the
// authoritative signal that the settings were delivered is that Approve succeeded.
func logSettingsConvergence(t *testing.T, don *cre.Don, hash string, timeout time.Duration) {
	t.Helper()
	want := len(don.Nodes)
	deadline := time.Now().Add(timeout)
	for {
		got := countContainersWithSettingsHash(hash)
		if got >= want {
			t.Logf("[cresettings] DON %q: %d/%d nodes logged settings hash %s", don.Name, got, want, shortHash(hash))
			return
		}
		if time.Now().After(deadline) {
			t.Logf("[cresettings] DON %q: only %d/%d nodes logged settings hash %s within %s "+
				"(delivery via Approve already succeeded; log-scan is best-effort)",
				don.Name, got, want, shortHash(hash), timeout)
			return
		}
		time.Sleep(3 * time.Second)
	}
}

// countContainersWithSettingsHash returns the number of containers whose logs contain
// the settings-update marker together with the given hash.
func countContainersWithSettingsHash(hash string) int {
	logStreams, err := framework.StreamContainerLogs(
		client.ContainerListOptions{All: true},
		client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	if err != nil {
		framework.L.Warn().Err(err).Msg("[cresettings] could not stream container logs")
		return 0
	}
	count := 0
	for _, reader := range logStreams {
		content, readErr := readContainerLogs(reader) // closes reader
		if readErr != nil {
			continue
		}
		if strings.Contains(content, creSettingsUpdateLogMarker) && strings.Contains(content, hash) {
			count++
		}
	}
	return count
}

func hashString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func shortHash(h string) string {
	if len(h) <= 12 {
		return h
	}
	return h[:12]
}
