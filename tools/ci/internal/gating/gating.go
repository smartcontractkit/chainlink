package gating

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	eventWorkflowDispatch = "workflow_dispatch"
	eventPush             = "push"
	eventMergeGroup       = "merge_group"
	eventPullRequest      = "pull_request"
	refTypeTag            = "tag"
	defaultBranch         = "develop"
)

// Inputs are the signals that drive the integration-tests gating decisions.
type Inputs struct {
	EventName           string
	RefName             string
	RefType             string
	CREChanges          bool
	CCIPChanges         bool
	RunE2ELabel         bool
	SkipRegressionLabel bool
	SkipMixedEnvLabel   bool
}

// Decisions holds the computed integration-tests gates and derived image builds.
type Decisions struct {
	CREShouldRun      bool
	CREWithRegression bool
	CRERunMixedEnv    bool
	CCIPShouldRun     bool
	BuildCoreImage    bool
	BuildPluginsImage bool
}

// Evaluate computes every gating decision from the given inputs.
func Evaluate(in Inputs) Decisions {
	cre := in.EventName == eventWorkflowDispatch ||
		(in.EventName == eventPush && (in.CREChanges || in.RefType == refTypeTag)) ||
		(in.EventName == eventMergeGroup && in.CREChanges) ||
		(in.EventName == eventPullRequest && (in.CREChanges || in.RunE2ELabel))

	ccip := in.EventName == eventWorkflowDispatch ||
		(in.EventName == eventPush && (in.CCIPChanges || in.RefType == refTypeTag)) ||
		(in.EventName == eventMergeGroup && in.CCIPChanges)

	regression := in.EventName == eventWorkflowDispatch ||
		(in.EventName == eventPush && in.RefName == defaultBranch) ||
		(in.EventName == eventPullRequest && !in.SkipRegressionLabel)

	mixedEnv := in.EventName == eventWorkflowDispatch ||
		(in.EventName == eventPush && in.RefName == defaultBranch) ||
		(in.EventName == eventPullRequest && !in.SkipMixedEnvLabel)

	return Decisions{
		CREShouldRun:      cre,
		CREWithRegression: regression,
		CRERunMixedEnv:    mixedEnv,
		CCIPShouldRun:     ccip,
		BuildCoreImage:    cre || ccip,
		BuildPluginsImage: cre,
	}
}

// OutputVars renders the decisions as GitHub Actions output variables.
func (d Decisions) OutputVars() map[string]string {
	return map[string]string{
		"cre-should-run":      strconv.FormatBool(d.CREShouldRun),
		"cre-with-regression": strconv.FormatBool(d.CREWithRegression),
		"cre-run-mixed-env":   strconv.FormatBool(d.CRERunMixedEnv),
		"ccip-should-run":     strconv.FormatBool(d.CCIPShouldRun),
		"build-core-image":    strconv.FormatBool(d.BuildCoreImage),
		"build-plugins-image": strconv.FormatBool(d.BuildPluginsImage),
	}
}

// SummaryTable renders the step-summary markdown table for a set of decisions.
func (d Decisions) SummaryTable(in Inputs) string {
	var b strings.Builder
	b.WriteString("### Integration Test Gating Decisions\n\n")
	b.WriteString("| Gate | Triggered? | Context |\n|:---|:---:|:---|\n")
	fmt.Fprintf(&b, "| **Build Core Image** | `%t` | CCIP or CRE tests requested |\n", d.BuildCoreImage)
	fmt.Fprintf(&b, "| **Build Plugins Image** | `%t` | CRE tests requested |\n", d.BuildPluginsImage)
	fmt.Fprintf(&b, "| **Core CRE Smoke Tests** | `%t` | Event: `%s`, changes: `%t`, label: `%t` |\n",
		d.CREShouldRun, in.EventName, in.CREChanges, in.RunE2ELabel)
	fmt.Fprintf(&b, "| **Core CRE Regression Tests** | `%t` | skip-regression: `%t` |\n",
		d.CREWithRegression, in.SkipRegressionLabel)
	fmt.Fprintf(&b, "| **Core CRE Mixed-Env Tests** | `%t` | skip-mixed-env: `%t` |\n",
		d.CRERunMixedEnv, in.SkipMixedEnvLabel)
	fmt.Fprintf(&b, "| **CCIP v1.6 E2E Tests** | `%t` | Event: `%s`, changes: `%t` |\n",
		d.CCIPShouldRun, in.EventName, in.CCIPChanges)
	return b.String()
}
