package cre

import (
	"testing"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

// Test_CRE_CRESettings_Override exercises the runtime CRE-settings override helper at
// every scope — workflow, org, global, and several scopes merged — without restarting
// the topology, and reverts each one. It doubles as executable documentation for
// system-tests/tests/test-helpers/cresettings_override.go.
// See core/scripts/cre/environment/docs/cresettings-override-README.md.
//
// Settings are passed as scope options (Global / Org / Owner / Workflow); each option's
// string is that scope's settings TOML with no scope prefix. Compose several in one call
// to apply multiple scopes together.
//
// Scope guidance (mirrored in the README): prefer the narrowest scope. Workflow/Owner
// scope isolates the override to your own workflow; Global changes it for the WHOLE
// environment and is the deliberate env-wide escape hatch.
//
// It is a single test (one environment spin-up) whose sections run SERIALLY: overrides
// mutate settings on the shared environment, so they must not run concurrently. The
// helper enforces this and fails fast if two overrides overlap.
//
// Requirements: a running Local CRE environment (default topology). Run with:
//
//	go test ./system-tests/tests/smoke/cre -run '^Test_CRE_CRESettings_Override$' -timeout 20m -v
//
// NOTE: the org/workflow IDs are illustrative. In a real test you use the actual org id /
// workflow id your deployed workflow runs under (obtainable from the deployment step). The
// delivery + cleanup mechanics are identical regardless of the ID value; only the isolation
// guarantee depends on using a real ID.

// creSettingsDON is only the DON whose applied/baseline document the test prints for
// readability — overrides are applied to ALL DONs. It is the workflow DON of the default
// topology (configs/workflow-gateway-capabilities-don.toml).
const creSettingsDON = "workflow"

// Illustrative scoped identifiers (must not be 0x-prefixed).
const (
	exampleOrgID      = "cresettingstestorg01"
	exampleWorkflowID = "abababababababababababababababababababababababababababababababab" // 62-char hex-like id
)

//nolint:paralleltest // mutates settings on the shared environment; must run serially
func Test_CRE_CRESettings_Override(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

	// 1) Workflow scope (preferred) — isolated to a single workflow. Reverted explicitly
	//    so the full apply -> revert cycle is visible.
	t.Log("=== workflow-scoped override ===")
	wf := t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Workflow(exampleWorkflowID, `
[PerWorkflow.HTTPAction]
CallLimit = '9'
[PerWorkflow.HTTPTrigger]
RateLimit = 'every5s:2'`))
	t.Logf("[workflow] applied for DON %q:\n%s", creSettingsDON, wf.AppliedTOML(creSettingsDON))
	wf.Reset(t)
	t.Logf("[workflow] reverted DON %q to baseline:\n%s", creSettingsDON, wf.BaselineTOML(creSettingsDON))

	// 2) Org scope.
	t.Log("=== org-scoped override ===")
	org := t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Org(exampleOrgID, `
[PerOrg]
BaseTriggerRetransmitEnabled = 'false'
WorkflowExecutionConcurrencyLimit = '42'`))
	t.Logf("[org=%s] applied for DON %q:\n%s", exampleOrgID, creSettingsDON, org.AppliedTOML(creSettingsDON))
	org.Reset(t)

	// 3) Global scope — the env-wide escape hatch (affects every workflow).
	t.Log("=== global-scoped override (env-wide) ===")
	global := t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(`
[PerWorkflow.HTTPAction]
CallLimit = '7'
[PerWorkflow]
ExecutionConcurrencyLimit = '3'`))
	t.Logf("[global] applied for DON %q:\n%s", creSettingsDON, global.AppliedTOML(creSettingsDON))
	global.Reset(t)

	// 4) Multiple scopes at once — compose the options in one call. Left to auto-revert via
	//    t.Cleanup, to exercise the automatic cleanup path as well.
	t.Log("=== multi-scope override (merged; auto-revert on cleanup) ===")
	multi := t_helpers.ApplyCRESettings(t, testEnv,
		t_helpers.Global("[PerWorkflow.HTTPAction]\nCallLimit = '11'"),
		t_helpers.Org(exampleOrgID, "[PerOrg]\nBaseTriggerRetransmitEnabled = 'false'"),
		t_helpers.Workflow(exampleWorkflowID, "[PerWorkflow]\nExecutionConcurrencyLimit = '1'"),
	)
	t.Logf("[multi-scope] applied for DON %q:\n%s", creSettingsDON, multi.AppliedTOML(creSettingsDON))
}
