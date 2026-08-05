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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"dario.cat/mergo"
	"github.com/moby/moby/client"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	// creSettingsExternalJobID is the fixed external job id shared by every CRESettings
	// job (at most one per node). Mirrors deployment/cre/jobs/pkg/cre_settings_job.go so
	// cancel-by-external-id targets the right proposal.
	creSettingsExternalJobID = "8561c20c-7d06-421e-a155-3baf21b1622b"

	// creSettingsUpdateLogMarker is what the node's cresettings delegate logs when it
	// applies an update (core/services/cresettings/delegate.go: `Updated settings`).
	// We scan for it plus the doc hash as best-effort proof that a node converged.
	creSettingsUpdateLogMarker = "Updated settings"

	// how long to wait (best effort) for every targeted node to log the applied hash.
	creSettingsConvergenceTimeout = 60 * time.Second

	// per-delivery timeout for the propose+approve round trip.
	creSettingsDeliveryTimeout = 2 * time.Minute
)

// Only one CRE settings override may be active at a time. Overrides mutate settings on
// the SHARED Local CRE environment, so two override tests running concurrently would
// clobber each other's document. These package-level vars enforce that and fail the
// second test fast with an actionable message. Re-application by the SAME test is
// allowed (keyed on the owning test name), so apply -> Reset -> apply within one test,
// or two applies from the same test, are both fine.
var (
	creSettingsActiveMu    sync.Mutex
	creSettingsActiveOwner string // t.Name() of the test currently holding an override; "" if none
)

func claimCRESettingsOverride(owner string) error {
	creSettingsActiveMu.Lock()
	defer creSettingsActiveMu.Unlock()
	if creSettingsActiveOwner != "" && creSettingsActiveOwner != owner {
		return fmt.Errorf(
			"a CRE settings override from %q is still active on the shared environment.\n"+
				"Settings-override tests mutate shared state and must run serially: remove t.Parallel() "+
				"from this test (and do not add it to the CRE_TEST_PARALLEL_ENABLED set). If you meant to "+
				"change settings within one test, call handle.Reset(t) before applying again.\n"+
				"See core/scripts/cre/environment/docs/cresettings-override-README.md",
			creSettingsActiveOwner)
	}
	creSettingsActiveOwner = owner
	return nil
}

func releaseCRESettingsOverride(owner string) {
	creSettingsActiveMu.Lock()
	defer creSettingsActiveMu.Unlock()
	if creSettingsActiveOwner == owner {
		creSettingsActiveOwner = ""
	}
}

// CRESettingsOverrides is a scoped set of CRE settings overrides to apply for the
// duration of a test. All values are strings, as required by the settings schema.
// Keys are dotted setting paths, e.g. "PerWorkflow.HTTPAction.CallLimit" or
// "PerOrg.BaseTriggerRetransmitEnabled".
type CRESettingsOverrides struct {
	// Global applies to every org/owner/workflow (the [global] scope).
	Global map[string]string
	// Org is keyed by org id (the [org.<id>] scope).
	Org map[string]map[string]string
	// Owner is keyed by workflow-owner hex WITHOUT the 0x prefix (the [owner.<id>] scope).
	Owner map[string]map[string]string
	// Workflow is keyed by workflow id hex WITHOUT the 0x prefix (the [workflow.<id>] scope).
	Workflow map[string]map[string]string
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

// ApplyCRESettings overrides CRE settings across the whole environment at runtime,
// without restarting the topology, and registers a t.Cleanup that restores the pre-test
// baseline when the test finishes.
//
// Settings are applied to every DON that has worker (plugin) nodes — the workflow and
// capabilities DONs — since CRE settings are enforced there and the delivery changeset
// targets type=plugin nodes. Bootstrap-only DONs (e.g. bootstrap-gateway) have no such
// nodes and are skipped. The user does not choose which DONs are targeted.
//
// For each DON it captures the DON's boot CL_CRE_SETTINGS as the baseline, merges the
// overrides on top, and proposes+approves a `cresettings` job to every node of the DON.
// Approve only returns once every node accepted the job, so a successful call means the
// environment converged on the new settings.
//
// It fails the test (require) if delivery to any DON fails. A best-effort log-scan
// confirmation is emitted via t.Logf for visibility.
func ApplyCRESettings(t *testing.T, env *ttypes.TestEnvironment, o CRESettingsOverrides) *CRESettingsHandle {
	t.Helper()

	require.NotNil(t, env, "test environment must not be nil")
	require.NotNil(t, env.CreEnvironment, "CreEnvironment must not be nil")
	require.NotNil(t, env.CreEnvironment.CldfEnvironment, "CldfEnvironment must not be nil")
	require.NotNil(t, env.Dons, "Dons must not be nil")

	// The same overrides are layered onto every targeted DON; only the baseline differs.
	overrideDoc := overridesToDoc(o)
	require.NotEmpty(t, overrideDoc, "no overrides provided")

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

// deliverCRESettings cancels any active settings job on every node of the DON, then
// proposes and approves a new one carrying settingsTOML. Approve returns only after
// every targeted node accepted the proposal.
func deliverCRESettings(env *ttypes.TestEnvironment, don *cre.Don, settingsTOML string) error {
	ctx, cancel := context.WithTimeout(context.Background(), creSettingsDeliveryTimeout)
	defer cancel()

	// At most one CRESettings job may be active per node; cancel any existing proposal
	// first (mirrors `env swap capability`). Best-effort: on the first apply there is
	// nothing to cancel.
	for _, node := range don.Nodes {
		if _, err := node.CancelProposalsByExternalJobID(ctx, []string{creSettingsExternalJobID}); err != nil {
			framework.L.Warn().
				Str("don", don.Name).
				Str("node", node.JobDistributorDetails.NodeID).
				Err(err).
				Msg("[cresettings] could not cancel existing settings proposal (continuing)")
		}
	}

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

	out, err := (cre_jobs.ProposeJobSpec{}).Apply(*env.CreEnvironment.CldfEnvironment, input)
	if err != nil {
		return fmt.Errorf("propose settings job: %w", err)
	}

	// Collect the per-node proposed specs so we can approve them on each node.
	specs := make(map[string][]string)
	for _, r := range out.Reports {
		o, ok := r.Output.(cre_jobs_ops.ProposeCRESettingsJobsOutput)
		if !ok {
			return fmt.Errorf("unexpected settings job report output type: %T", r.Output)
		}
		if mErr := mergo.Merge(&specs, o.Specs, mergo.WithAppendSlice); mErr != nil {
			return fmt.Errorf("merge settings job specs: %w", mErr)
		}
	}
	if len(specs) == 0 {
		return fmt.Errorf("settings job proposal produced no specs for DON %q", don.Name)
	}

	if err := jobs.Approve(ctx, env.CreEnvironment.CldfEnvironment.Offchain, env.Dons, specs); err != nil {
		return fmt.Errorf("approve settings job: %w", err)
	}
	return nil
}

// overridesToDoc converts the scoped overrides into a nested map matching the settings
// document shape: {global:{...}, org:{<id>:{...}}, owner:{...}, workflow:{...}}. Dotted
// keys ("PerWorkflow.HTTPAction.CallLimit") are expanded into nested tables so the
// marshalled TOML is scoped correctly.
func overridesToDoc(o CRESettingsOverrides) map[string]any {
	doc := map[string]any{}
	for k, v := range o.Global {
		setNested(doc, append([]string{"global"}, strings.Split(k, ".")...), v)
	}
	addScoped := func(scope string, byID map[string]map[string]string) {
		for id, m := range byID {
			for k, v := range m {
				setNested(doc, append([]string{scope, id}, strings.Split(k, ".")...), v)
			}
		}
	}
	addScoped("org", o.Org)
	addScoped("owner", o.Owner)
	addScoped("workflow", o.Workflow)
	return doc
}

func setNested(m map[string]any, path []string, val string) {
	for i := 0; i < len(path)-1; i++ {
		child, ok := m[path[i]].(map[string]any)
		if !ok {
			child = map[string]any{}
			m[path[i]] = child
		}
		m = child
	}
	m[path[len(path)-1]] = val
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
