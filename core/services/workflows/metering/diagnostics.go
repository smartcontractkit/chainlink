package metering

import (
	"sort"

	"github.com/shopspring/decimal"

	eventspb "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

// TEMPORARY incident diagnostics - DO NOT EXTEND.
//
// diagnosticTargetOrgID hard-codes the single organization under investigation
// for a native gas spend metering incident (GAS.4741433654826277614 /
// SpendValue 0.001158 / aggregate native spend 1158000000000000 at billing
// ingress). While the engine consumes SpendValueInGasUnits directly when
// populated (otherwise legacy Shift(18)), the conversion path actually taken
// per execution is not visible in production logs, which prevents determining
// whether the native-unit field was absent, correct, or mis-scaled, and
// whether rates/metering-mode contributed.
//
// This constant scopes a single structured Info-level log event, emitted once
// per finished workflow execution by Reports.End, to that organization only.
// It is deliberately NOT configurable: an incident diagnostic does not warrant
// a configuration framework, and an exact hard-coded match guarantees the
// logging cannot silently widen to other organizations.
//
// REMOVAL: revert the commit that introduced this file (and the Report/
// Reports.End hooks) once the incident is resolved or the diagnostic build is
// retired from staging. No other behavior depends on it.
const diagnosticTargetOrgID = "org_9PndwXNviSVigJtI"

// source of the resolved organization identity.
const (
	diagOrgSourceBilling = "billing_response"
	diagOrgSourceLabel   = "engine_label"
	diagOrgSourceUnknown = "unknown"
)

// conversion paths for spend values, mirroring parseSpendValue.
const (
	diagPathNative     = "native_units"          // GAS: SpendValueInGasUnits populated and used directly
	diagPathLegacy     = "legacy_spend_shift18"  // GAS: SpendValueInGasUnits absent; SpendValue.Shift(18) used
	diagPathSkipped    = "skipped_invalid_value" // value unparseable; excluded from settlement aggregation
	diagPathNotApplied = "not_gas"               // non-gas spend unit; no native conversion applies
)

// resolveOrgID resolves the organization identity for an execution. The
// billing service's GetWorkflowExecutionRates response is authoritative; the
// engine-provided label (platform.KeyOrganizationID) is a fallback because
// report labels may omit organization_id. An empty resolution must never
// enable diagnostics.
func resolveOrgID(billingOrgID string, labels map[string]string) (string, string) {
	if billingOrgID != "" {
		return billingOrgID, diagOrgSourceBilling
	}
	if label, ok := labels[platform.KeyOrganizationID]; ok && label != "" {
		return label, diagOrgSourceLabel
	}
	return "", diagOrgSourceUnknown
}

// meteringDiagnostics is the structured payload logged for the target
// organization. It contains only metering/billing metadata: no capability
// inputs/outputs, workflow secrets, credentials or request payloads.
type meteringDiagnostics struct {
	OrgID          string            `json:"orgID"`
	OrgIDSource    string            `json:"orgIDSource"`
	WorkflowID     string            `json:"workflowID"`
	ExecutionID    string            `json:"executionID"`
	WorkflowOwner  string            `json:"workflowOwner"`
	LocalPeerID    string            `json:"localPeerID"`
	EngineVersion  string            `json:"engineVersion"`
	DonID          string            `json:"donID"`
	DonN           string            `json:"donN"`
	DonF           string            `json:"donF"`
	MeteringMode   bool              `json:"meteringMode"`
	MeteringReason string            `json:"meteringReason,omitempty"`
	RateCard       map[string]string `json:"rateCard"`
	Steps          []stepDiagnostics `json:"steps,omitempty"`
	SpentCredits   string            `json:"spentCredits"`
	EmitReceiptErr string            `json:"emitReceiptErr,omitempty"`
	SendReceiptErr string            `json:"sendReceiptErr,omitempty"`
}

// stepDiagnostics captures a single metered step (capability invocation).
type stepDiagnostics struct {
	Ref          string            `json:"ref"`
	CapabilityID string            `json:"capabilityID,omitempty"`
	CapDONN      uint32            `json:"capdonN"`
	Units        []unitDiagnostics `json:"units,omitempty"`
}

// unitDiagnostics captures one spend unit of a step. NodesReported vs DonN
// (top-level) reveals missing per-node spend responses for the unit.
type unitDiagnostics struct {
	SpendUnit     string            `json:"spendUnit"`
	RatePerCredit string            `json:"ratePerCredit,omitempty"` // rate-card units per credit; empty => no rate-card entry for this unit
	NodesReported int               `json:"nodesReported"`
	NodeDetails   []nodeDiagnostics `json:"nodeDetails,omitempty"`
	AggSpendValue string            `json:"aggSpendValue,omitempty"` // aggregated spend (native units for GAS.*)
	AggCREValue   string            `json:"aggCREValue,omitempty"`   // aggregated spend converted to universal credits
}

// nodeDiagnostics captures one node's spend response for a unit and how the
// engine converted it. For GAS.* units with both values parseable,
// ScaleDelta = parsed native value - SpendValue.Shift(18); a non-zero delta
// indicates a mis-scaled native field.
type nodeDiagnostics struct {
	PeerID               string `json:"peerID,omitempty"`
	SpendValue           string `json:"spendValue"`           // legacy field, as reported
	SpendValueInGasUnits string `json:"spendValueInGasUnits"` // native field, as reported ("" when absent)
	NativePresent        bool   `json:"nativePresent"`
	ConversionPath       string `json:"conversionPath"`
	ParsedValue          string `json:"parsedValue,omitempty"` // value used for settlement aggregation (native units for GAS.*)
	LegacyImpliedNative  string `json:"legacyImpliedNative,omitempty"`
	ScaleDelta           string `json:"scaleDelta,omitempty"`
	CREValue             string `json:"creValue,omitempty"`
	ParseErr             string `json:"parseErr,omitempty"`
}

// diagnosticsEnabled reports whether incident diagnostics apply to this
// report. Strict exact-match on the resolved organization identity; unknown
// or mismatched identity keeps diagnostics off.
func (r *Report) diagnosticsEnabled() bool {
	orgID, _ := r.resolveOrgID()
	return orgID == diagnosticTargetOrgID
}

func (r *Report) resolveOrgID() (string, string) {
	return resolveOrgID(r.billingOrgID, r.labels)
}

// logMeteringDiagnostics emits the single per-execution incident diagnostics
// event at Info level (operationally visible). Called once from Reports.End
// (not from the SubmitWorkflowReceipt retry loop) so the event is logged once
// per logical event, not once per receipt retry. Failures here can never fail
// an execution: the payload is derived from report state via FormatReport
// under a read lock, and no billing/metering state is modified.
func (r *Report) logMeteringDiagnostics(emitErr, sendErr error) {
	if !r.diagnosticsEnabled() {
		return
	}

	orgID, orgSource := r.resolveOrgID()
	r.lggr.Infow("scoped metering diagnostics",
		"orgID", orgID,
		"diagnostics", r.buildMeteringDiagnostics(orgID, orgSource, emitErr, sendErr),
	)
}

// buildMeteringDiagnostics derives the diagnostics payload from the final
// report state (the same data submitted to the billing service) plus the rate
// card and total credits. It never mutates report state.
func (r *Report) buildMeteringDiagnostics(orgID, orgSource string, emitErr, sendErr error) *meteringDiagnostics {
	rpt := r.FormatReport()

	diag := &meteringDiagnostics{
		OrgID:          orgID,
		OrgIDSource:    orgSource,
		WorkflowID:     rpt.Metadata.WorkflowID,
		ExecutionID:    rpt.Metadata.WorkflowExecutionID,
		WorkflowOwner:  rpt.Metadata.WorkflowOwner,
		LocalPeerID:    rpt.Metadata.P2PID,
		EngineVersion:  rpt.Metadata.EngineVersion,
		DonID:          r.labels[platform.KeyDonID],
		DonN:           r.labels[platform.KeyDonN],
		DonF:           r.labels[platform.KeyDonF],
		MeteringMode:   rpt.MeteringMode,
		MeteringReason: rpt.Message,
		RateCard:       make(map[string]string, len(r.rateCard)),
		SpentCredits:   r.balance.GetSpent().String(),
	}

	for unit, rate := range r.rateCard {
		diag.RateCard[unit] = rate.String()
	}

	// order steps and units for deterministic output, like FormatReport
	refs := make([]string, 0, len(rpt.Steps))
	for ref := range rpt.Steps {
		refs = append(refs, ref)
	}
	sort.Strings(refs)

	for _, ref := range refs {
		step := rpt.Steps[ref]
		stepDiag := stepDiagnostics{
			Ref:          ref,
			CapabilityID: r.stepCapabilityID(ref),
			CapDONN:      step.CapdonN,
		}

		// group node details by spend unit; includes units excluded from
		// aggregation (e.g. RPC_EVM) so their evidence is not lost
		nodesByUnit := map[string][]nodeDiagnostics{}
		for _, node := range step.Nodes {
			nodesByUnit[node.SpendUnit] = append(nodesByUnit[node.SpendUnit], newNodeDiagnostics(node))
		}

		units := make([]string, 0, len(nodesByUnit))
		for unit := range nodesByUnit {
			units = append(units, unit)
		}
		sort.Strings(units)

		aggByUnit := map[string]*eventspb.AggregatedSpendDetail{}
		for _, agg := range step.AggSpend {
			aggByUnit[agg.SpendUnit] = agg
		}

		for _, unit := range units {
			unitDiag := unitDiagnostics{
				SpendUnit:     unit,
				NodesReported: len(nodesByUnit[unit]),
				NodeDetails:   nodesByUnit[unit],
			}
			if rate, ok := r.rateCard[unit]; ok {
				unitDiag.RatePerCredit = rate.String()
			}
			if agg, ok := aggByUnit[unit]; ok {
				unitDiag.AggSpendValue = agg.SpendValue
				unitDiag.AggCREValue = agg.SpendValueCre
			}

			stepDiag.Units = append(stepDiag.Units, unitDiag)
		}

		diag.Steps = append(diag.Steps, stepDiag)
	}

	if emitErr != nil {
		diag.EmitReceiptErr = emitErr.Error()
	}
	if sendErr != nil {
		diag.SendReceiptErr = sendErr.Error()
	}

	return diag
}

// newNodeDiagnostics reconstructs how parseSpendValue treated one node detail.
// It re-derives the conversion decision from the raw reported fields without
// re-running settlement, so diagnostics cannot perturb billing arithmetic.
func newNodeDiagnostics(node *eventspb.MeteringReportNodeDetail) nodeDiagnostics {
	diag := nodeDiagnostics{
		PeerID:               node.Peer_2PeerId,
		SpendValue:           node.SpendValue,
		SpendValueInGasUnits: node.SpendValueInGasUnits,
		NativePresent:        isGasSpendType(node.SpendUnit) && node.SpendValueInGasUnits != "",
		CREValue:             node.SpendValueCre,
		ConversionPath:       diagPathNotApplied,
	}

	if !isGasSpendType(node.SpendUnit) {
		if v, err := decimal.NewFromString(node.SpendValue); err == nil {
			diag.ParsedValue = v.String()
		} else {
			diag.ConversionPath = diagPathSkipped
			diag.ParseErr = err.Error()
		}
		return diag
	}

	legacy, legacyErr := decimal.NewFromString(node.SpendValue)
	if legacyErr != nil {
		diag.ParseErr = legacyErr.Error()
	} else {
		diag.LegacyImpliedNative = legacy.Shift(18).String()
	}

	if !diag.NativePresent {
		if legacyErr != nil {
			diag.ConversionPath = diagPathSkipped
			return diag
		}
		diag.ConversionPath = diagPathLegacy
		diag.ParsedValue = diag.LegacyImpliedNative
		return diag
	}

	native, nativeErr := decimal.NewFromString(node.SpendValueInGasUnits)
	if nativeErr != nil {
		diag.ParseErr = nativeErr.Error()
		diag.ConversionPath = diagPathSkipped
		return diag
	}

	diag.ConversionPath = diagPathNative
	diag.ParsedValue = native.String()
	if legacyErr == nil {
		diag.ScaleDelta = native.Sub(legacy.Shift(18)).String()
	}

	return diag
}

// stepCapabilityID returns the capability id recorded for a step ref, if any.
func (r *Report) stepCapabilityID(ref string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if step, ok := r.steps[ref]; ok {
		return step.CapabilityID
	}
	return ""
}
