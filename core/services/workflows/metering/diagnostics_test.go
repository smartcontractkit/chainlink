package metering

import (
	"errors"
	"maps"
	"strconv"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
)

// Synthetic fixtures mirroring the target incident: GAS.4741433654826277614,
// SpendValue 0.001158, aggregate native spend 1158000000000000 and a gas rate
// of ~156881 gas tokens per credit.
const (
	incidentSpendValue        = "0.001158"
	incidentSpendValueNative  = "1158000000000000" // 0.001158 shifted by 18
	incidentMisScaledNative   = "1158000000000"    // 1000x too small
	incidentGasRatePerCredit  = "156881"           // gas tokens per credit
	incidentGasUnit           = "GAS.4741433654826277614"
	incidentComputeUnit       = "RESOURCE_TYPE_COMPUTE"
	incidentComputeSpendValue = "100"
	incidentComputeRate       = "2"
	targetOrgDonN             = "5"
	targetOrgPeerID           = "p2p_target"
)

func targetOrgLabels() map[string]string {
	labels := map[string]string{}
	maps.Copy(labels, defaultLabels)
	labels[platform.KeyOrganizationID] = diagnosticTargetOrgID
	labels[platform.KeyDonN] = targetOrgDonN
	labels[platform.KeyP2PID] = targetOrgPeerID
	return labels
}

func incidentRateResponse(orgID string) *billing.GetWorkflowExecutionRatesResponse {
	return &billing.GetWorkflowExecutionRatesResponse{
		OrganizationId: orgID,
		RateCards: []*billing.RateCard{
			{
				ResourceType:    billing.ResourceType_RESOURCE_TYPE_COMPUTE,
				MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS,
				UnitsPerCredit:  incidentComputeRate,
			},
		},
		GasTokensPerCredit: map[uint64]string{
			4741433654826277614: incidentGasRatePerCredit,
		},
	}
}

func settleIncidentGasStep(t *testing.T, report *Report, nodes []capabilities.MeteringNodeDetail) {
	t.Helper()

	_, err := report.Deduct("step1", ByResource(incidentGasUnit, "target-capability", decimal.NewFromInt(1)))
	require.NoError(t, err)
	require.NoError(t, report.Settle("step1", capabilities.ResponseMetadata{Metering: nodes}))
}

func countScopedDiagnosticsLogs(entries []observer.LoggedEntry) int {
	count := 0
	for _, e := range entries {
		if strings.Contains(e.Message, "scoped metering diagnostics") {
			count++
		}
	}
	return count
}

func diagnosticsPayloadFromLog(t *testing.T, entries []observer.LoggedEntry) *meteringDiagnostics {
	t.Helper()

	for _, e := range entries {
		if !strings.Contains(e.Message, "scoped metering diagnostics") {
			continue
		}
		for _, f := range e.Context {
			if f.Key == "diagnostics" {
				payload, ok := f.Interface.(*meteringDiagnostics)
				require.True(t, ok, "diagnostics field should be *meteringDiagnostics")
				return payload
			}
		}
	}

	t.Fatal("no scoped metering diagnostics log found")
	return nil
}

func TestResolveOrgID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		billingOrgID string
		labels       map[string]string
		wantOrgID    string
		wantSource   string
	}{
		{
			name:         "billing response is authoritative",
			billingOrgID: "org-billing",
			labels:       map[string]string{platform.KeyOrganizationID: "org-label"},
			wantOrgID:    "org-billing",
			wantSource:   diagOrgSourceBilling,
		},
		{
			name:         "label fallback when billing response omits org",
			billingOrgID: "",
			labels:       map[string]string{platform.KeyOrganizationID: "org-label"},
			wantOrgID:    "org-label",
			wantSource:   diagOrgSourceLabel,
		},
		{
			name:         "unknown when both omit org",
			billingOrgID: "",
			labels:       map[string]string{},
			wantOrgID:    "",
			wantSource:   diagOrgSourceUnknown,
		},
		{
			name:         "empty label value is not a fallback",
			billingOrgID: "",
			labels:       map[string]string{platform.KeyOrganizationID: ""},
			wantOrgID:    "",
			wantSource:   diagOrgSourceUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orgID, source := resolveOrgID(tc.billingOrgID, tc.labels)
			assert.Equal(t, tc.wantOrgID, orgID)
			assert.Equal(t, tc.wantSource, source)
		})
	}
}

func Test_Report_DiagnosticsEnabled(t *testing.T) {
	t.Parallel()

	newReportWithOrg := func(t *testing.T, respOrgID string, respErr error, labelOrgID string) *Report {
		t.Helper()

		labels := targetOrgLabels()
		labels[platform.KeyOrganizationID] = labelOrgID

		billingClient := mocks.NewBillingClient(t)
		if respErr != nil {
			billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
				Return(nil, respErr)
		} else {
			billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
				Return(incidentRateResponse(respOrgID), nil)
		}

		report, err := NewReport(t.Context(), labels, logger.Nop(), billingClient, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		require.NoError(t, err)
		return report
	}

	t.Run("enabled for target org from billing response", func(t *testing.T) {
		t.Parallel()

		report := newReportWithOrg(t, diagnosticTargetOrgID, nil, "org-unrelated")
		assert.True(t, report.diagnosticsEnabled())
	})

	t.Run("enabled for target org from label fallback", func(t *testing.T) {
		t.Parallel()

		report := newReportWithOrg(t, "", nil, diagnosticTargetOrgID)
		assert.True(t, report.diagnosticsEnabled())
	})

	t.Run("enabled for target org label when rates call fails", func(t *testing.T) {
		t.Parallel()

		report := newReportWithOrg(t, "", errors.New("billing unavailable"), diagnosticTargetOrgID)
		assert.True(t, report.diagnosticsEnabled())
	})

	t.Run("disabled for non-target billing org even if label matches", func(t *testing.T) {
		t.Parallel()

		report := newReportWithOrg(t, "org-other", nil, diagnosticTargetOrgID)
		assert.False(t, report.diagnosticsEnabled())
	})

	t.Run("disabled for non-target org", func(t *testing.T) {
		t.Parallel()

		report := newReportWithOrg(t, "org-other", nil, "org-unrelated")
		assert.False(t, report.diagnosticsEnabled())
	})

	t.Run("disabled when identity is unknown", func(t *testing.T) {
		t.Parallel()

		report := newReportWithOrg(t, "", nil, "")
		assert.False(t, report.diagnosticsEnabled())
	})
}

func Test_Report_LogMeteringDiagnostics_Scoping(t *testing.T) {
	t.Parallel()

	t.Run("logs exactly one info event for target org resolved from billing response", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(incidentRateResponse(diagnosticTargetOrgID), nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

		// no org label: identity must resolve from the billing response alone
		labels := targetOrgLabels()
		delete(labels, platform.KeyOrganizationID)

		mrs := NewReports(billingClient, testAccountID, testWorkflowID, lggr, labels, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		r, err := mrs.Start(t.Context(), testWorkflowExecutionID)
		require.NoError(t, err)
		require.NoError(t, r.Reserve(t.Context()))
		require.NoError(t, mrs.End(t.Context(), testWorkflowExecutionID))

		all := logs.All()
		assert.Equal(t, 1, countScopedDiagnosticsLogs(all), "expected exactly one diagnostics event")

		payload := diagnosticsPayloadFromLog(t, all)
		assert.Equal(t, diagnosticTargetOrgID, payload.OrgID)
		assert.Equal(t, diagOrgSourceBilling, payload.OrgIDSource)
	})

	t.Run("does not log for non-target org", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(incidentRateResponse("org-somebody-else"), nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

		mrs := NewReports(billingClient, testAccountID, testWorkflowID, lggr, targetOrgLabels(), defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		r, err := mrs.Start(t.Context(), testWorkflowExecutionID)
		require.NoError(t, err)
		require.NoError(t, r.Reserve(t.Context()))
		require.NoError(t, mrs.End(t.Context(), testWorkflowExecutionID))

		assert.Zero(t, countScopedDiagnosticsLogs(logs.All()))
	})

	t.Run("does not log when identity is unknown", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(&billing.GetWorkflowExecutionRatesResponse{}, nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

		// neither billing response nor label carries an org id
		labels := targetOrgLabels()
		delete(labels, platform.KeyOrganizationID)

		mrs := NewReports(billingClient, testAccountID, testWorkflowID, lggr, labels, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		r, err := mrs.Start(t.Context(), testWorkflowExecutionID)
		require.NoError(t, err)
		require.NoError(t, r.Reserve(t.Context()))
		require.NoError(t, mrs.End(t.Context(), testWorkflowExecutionID))

		assert.Zero(t, countScopedDiagnosticsLogs(logs.All()))
	})

	t.Run("mismatched billing org overrides matching label", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(incidentRateResponse("org-other"), nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

		mrs := NewReports(billingClient, testAccountID, testWorkflowID, lggr, targetOrgLabels(), defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		r, err := mrs.Start(t.Context(), testWorkflowExecutionID)
		require.NoError(t, err)
		require.NoError(t, r.Reserve(t.Context()))
		require.NoError(t, mrs.End(t.Context(), testWorkflowExecutionID))

		assert.Zero(t, countScopedDiagnosticsLogs(logs.All()))
	})

	t.Run("logs once per execution across receipt retries", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(incidentRateResponse(diagnosticTargetOrgID), nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
			Return(nil, status.Error(codes.Unavailable, "transient")).Times(1)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
			Return(&emptypb.Empty{}, nil).Times(1)

		mrs := NewReports(billingClient, testAccountID, testWorkflowID, lggr, targetOrgLabels(), defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		r, err := mrs.Start(t.Context(), testWorkflowExecutionID)
		require.NoError(t, err)
		require.NoError(t, r.Reserve(t.Context()))
		require.NoError(t, mrs.End(t.Context(), testWorkflowExecutionID))

		assert.Equal(t, 1, countScopedDiagnosticsLogs(logs.All()),
			"diagnostics must be logged once per execution, not once per receipt retry")
	})
}

func Test_Report_BuildMeteringDiagnostics(t *testing.T) {
	t.Parallel()

	newSettledGasReport := func(t *testing.T, nodes []capabilities.MeteringNodeDetail) *Report {
		t.Helper()

		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(incidentRateResponse(diagnosticTargetOrgID), nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)

		report, err := NewReport(t.Context(), targetOrgLabels(), logger.Nop(), billingClient, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		require.NoError(t, err)
		require.NoError(t, report.Reserve(t.Context()))
		settleIncidentGasStep(t, report, nodes)
		return report
	}

	payloadOf := func(t *testing.T, report *Report) *meteringDiagnostics {
		t.Helper()
		return report.buildMeteringDiagnostics(diagnosticTargetOrgID, diagOrgSourceBilling, nil, nil)
	}

	firstGasNode := func(t *testing.T, diag *meteringDiagnostics) (nodeDiagnostics, unitDiagnostics, stepDiagnostics) {
		t.Helper()

		require.Len(t, diag.Steps, 1)
		step := diag.Steps[0]
		for _, u := range step.Units {
			if u.SpendUnit == incidentGasUnit {
				require.NotEmpty(t, u.NodeDetails, "gas unit should have node details")
				return u.NodeDetails[0], u, step
			}
		}
		t.Fatal("no gas unit in diagnostics")
		return nodeDiagnostics{}, unitDiagnostics{}, stepDiagnostics{}
	}

	incidentNativeCredits := decimal.RequireFromString(incidentSpendValueNative).
		Div(decimal.RequireFromString(incidentGasRatePerCredit)).
		Round(defaultDecimalPrecision)

	t.Run("identifies populated native units with correct scale", func(t *testing.T) {
		t.Parallel()

		report := newSettledGasReport(t, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue, SpendValueInGasUnits: incidentSpendValueNative},
		})

		diag := payloadOf(t, report)
		node, unit, step := firstGasNode(t, diag)

		assert.True(t, node.NativePresent)
		assert.Equal(t, diagPathNative, node.ConversionPath)
		assert.Equal(t, incidentSpendValueNative, node.ParsedValue)
		assert.Equal(t, incidentSpendValueNative, node.LegacyImpliedNative)
		assert.Equal(t, "0", node.ScaleDelta)
		assert.Empty(t, node.ParseErr)
		assert.Equal(t, incidentGasRatePerCredit, unit.RatePerCredit)
		assert.Equal(t, incidentSpendValueNative+".0000000000", unit.AggSpendValue)
		assert.Equal(t, incidentNativeCredits.StringFixed(defaultDecimalPrecision), unit.AggCREValue)
		assert.Equal(t, incidentNativeCredits.StringFixed(defaultDecimalPrecision), node.CREValue)
		assert.Equal(t, 1, unit.NodesReported)
		assert.Equal(t, "target-capability", step.CapabilityID)
		assert.Equal(t, uint32(1), step.CapDONN)
		assert.Equal(t, targetOrgDonN, diag.DonN)
		assert.Equal(t, targetOrgPeerID, diag.LocalPeerID)
		assert.Equal(t, workflowV2, diag.EngineVersion)
		assert.Equal(t, testWorkflowID, diag.WorkflowID)
		assert.Equal(t, testWorkflowExecutionID, diag.ExecutionID)
		assert.Equal(t, testAccountID, diag.WorkflowOwner)
		assert.False(t, diag.MeteringMode)
		assert.Equal(t, incidentNativeCredits.String(), diag.SpentCredits)
		assert.Contains(t, diag.RateCard, incidentGasUnit)
		assert.Equal(t, incidentGasRatePerCredit, diag.RateCard[incidentGasUnit])
	})

	t.Run("identifies absent native units on legacy shift path", func(t *testing.T) {
		t.Parallel()

		report := newSettledGasReport(t, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue},
		})

		node, unit, _ := firstGasNode(t, payloadOf(t, report))

		assert.False(t, node.NativePresent)
		assert.Empty(t, node.SpendValueInGasUnits)
		assert.Equal(t, diagPathLegacy, node.ConversionPath)
		assert.Equal(t, incidentSpendValueNative, node.ParsedValue)
		assert.Equal(t, incidentSpendValueNative, node.LegacyImpliedNative)
		assert.Empty(t, node.ScaleDelta)
		assert.Equal(t, incidentGasRatePerCredit, unit.RatePerCredit)
	})

	t.Run("identifies mis-scaled native units via non-zero scale delta", func(t *testing.T) {
		t.Parallel()

		report := newSettledGasReport(t, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue, SpendValueInGasUnits: incidentMisScaledNative},
		})

		node, _, _ := firstGasNode(t, payloadOf(t, report))

		assert.True(t, node.NativePresent)
		assert.Equal(t, diagPathNative, node.ConversionPath)
		assert.Equal(t, incidentMisScaledNative, node.ParsedValue)
		wantDelta := decimal.RequireFromString(incidentMisScaledNative).
			Sub(decimal.RequireFromString(incidentSpendValue).Shift(18))
		assert.Equal(t, wantDelta.String(), node.ScaleDelta)
		assert.NotEqual(t, "0", node.ScaleDelta)
	})

	t.Run("identifies unparseable native units as skipped", func(t *testing.T) {
		t.Parallel()

		report := newSettledGasReport(t, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue, SpendValueInGasUnits: "not-a-number"},
		})

		node, unit, _ := firstGasNode(t, payloadOf(t, report))

		assert.True(t, node.NativePresent)
		assert.Equal(t, diagPathSkipped, node.ConversionPath)
		assert.NotEmpty(t, node.ParseErr)
		assert.Empty(t, node.ParsedValue)
		// invalid values are excluded from settlement aggregation
		assert.Equal(t, "0.0000000000", unit.AggSpendValue)
	})

	t.Run("records zero gas spend without native units", func(t *testing.T) {
		t.Parallel()

		report := newSettledGasReport(t, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: "0"},
		})

		node, unit, _ := firstGasNode(t, payloadOf(t, report))

		assert.False(t, node.NativePresent)
		assert.Equal(t, diagPathLegacy, node.ConversionPath)
		assert.Equal(t, "0", node.ParsedValue)
		assert.Equal(t, "0.0000000000", unit.AggSpendValue)
		assert.Equal(t, "0.0000000000", unit.AggCREValue)
	})

	t.Run("flags missing gas rate card entry", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(&billing.GetWorkflowExecutionRatesResponse{
				OrganizationId: diagnosticTargetOrgID,
				RateCards:      successRates, // compute only; no gas rates
			}, nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)

		report, err := NewReport(t.Context(), targetOrgLabels(), logger.Nop(), billingClient, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		require.NoError(t, err)
		require.NoError(t, report.Reserve(t.Context()))
		settleIncidentGasStep(t, report, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue},
		})

		node, unit, _ := firstGasNode(t, payloadOf(t, report))

		assert.Equal(t, diagPathLegacy, node.ConversionPath)
		assert.Empty(t, unit.RatePerCredit, "gas unit without a rate-card entry must show an empty rate")
		assert.Equal(t, "0.0000000000", unit.AggCREValue)
		assert.True(t, report.isMeteringMode(), "missing gas rate entry switches the report to metering mode")
		assert.NotContains(t, payloadOf(t, report).RateCard, incidentGasUnit)
	})

	t.Run("covers non-gas units without native conversion", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(incidentRateResponse(diagnosticTargetOrgID), nil)
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)

		report, err := NewReport(t.Context(), targetOrgLabels(), logger.Nop(), billingClient, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		require.NoError(t, err)
		require.NoError(t, report.Reserve(t.Context()))

		_, err = report.Deduct("step1", ByResource(incidentComputeUnit, "compute-capability", decimal.NewFromInt(1)))
		require.NoError(t, err)
		require.NoError(t, report.Settle("step1", capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentComputeUnit, SpendValue: incidentComputeSpendValue},
		}}))

		diag := payloadOf(t, report)
		require.Len(t, diag.Steps, 1)

		var computeUnit unitDiagnostics
		for _, u := range diag.Steps[0].Units {
			if u.SpendUnit == incidentComputeUnit {
				computeUnit = u
			}
		}
		require.NotEmpty(t, computeUnit.NodeDetails)

		node := computeUnit.NodeDetails[0]
		assert.Equal(t, diagPathNotApplied, node.ConversionPath)
		assert.False(t, node.NativePresent)
		assert.Equal(t, incidentComputeSpendValue, node.ParsedValue)
		assert.Empty(t, node.LegacyImpliedNative)
		assert.Empty(t, node.ScaleDelta)
		assert.Equal(t, incidentComputeRate, computeUnit.RatePerCredit)
	})

	t.Run("captures metering mode, reason and receipt errors", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
			Return(nil, errors.New("rates unavailable"))
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)

		report, err := NewReport(t.Context(), targetOrgLabels(), logger.Nop(), billingClient, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		require.NoError(t, err)
		require.NoError(t, report.Reserve(t.Context()))
		settleIncidentGasStep(t, report, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue},
		})

		emitErr := errors.New("emit failed")
		sendErr := errors.New("send failed")
		diag := report.buildMeteringDiagnostics(diagnosticTargetOrgID, diagOrgSourceLabel, emitErr, sendErr)

		assert.True(t, diag.MeteringMode)
		assert.Contains(t, diag.MeteringReason, "rates unavailable")
		assert.Equal(t, emitErr.Error(), diag.EmitReceiptErr)
		assert.Equal(t, sendErr.Error(), diag.SendReceiptErr)
		assert.Empty(t, diag.RateCard, "failed rates call leaves an empty rate card")
		assert.Equal(t, "0", diag.SpentCredits, "metering mode skips local credit accounting")
		assert.Equal(t, diagOrgSourceLabel, diag.OrgIDSource)
	})

	t.Run("missing gas node responses are visible against DonN", func(t *testing.T) {
		t.Parallel()

		// only 1 of 5 DON nodes reports gas spend for the step
		report := newSettledGasReport(t, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue},
		})

		diag := payloadOf(t, report)

		assert.Equal(t, targetOrgDonN, diag.DonN)
		_, unit, _ := firstGasNode(t, diag)
		assert.Equal(t, 1, unit.NodesReported)
	})

	t.Run("includes units excluded from aggregation such as RPC_EVM", func(t *testing.T) {
		t.Parallel()

		report := newSettledGasReport(t, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue},
			{Peer2PeerID: "peer1", SpendUnit: "RPC_EVM", SpendValue: "7"},
		})

		diag := payloadOf(t, report)

		var rpcUnit *unitDiagnostics
		for i, u := range diag.Steps[0].Units {
			if u.SpendUnit == "RPC_EVM" {
				rpcUnit = &diag.Steps[0].Units[i]
			}
		}
		require.NotNil(t, rpcUnit, "RPC_EVM unit should appear in diagnostics")
		assert.Equal(t, 1, rpcUnit.NodesReported)
		require.Len(t, rpcUnit.NodeDetails, 1)
		assert.Equal(t, "7", rpcUnit.NodeDetails[0].SpendValue)
		assert.Equal(t, diagPathNotApplied, rpcUnit.NodeDetails[0].ConversionPath)
		// Settle skips RPC_EVM aggregation entirely
		assert.Empty(t, rpcUnit.AggSpendValue)
		assert.Empty(t, rpcUnit.AggCREValue)
		assert.Empty(t, rpcUnit.RatePerCredit)
	})
}

func Test_Report_LogMeteringDiagnostics_MeteringModeEnd(t *testing.T) {
	t.Parallel()

	lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
	billingClient := mocks.NewBillingClient(t)
	billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
		Return(nil, errors.New("rates unavailable"))
	billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
	billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

	mrs := NewReports(billingClient, testAccountID, testWorkflowID, lggr, targetOrgLabels(), defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
	r, err := mrs.Start(t.Context(), testWorkflowExecutionID)
	require.NoError(t, err)
	require.NoError(t, r.Reserve(t.Context()))
	require.NoError(t, mrs.End(t.Context(), testWorkflowExecutionID))

	all := logs.All()
	assert.Equal(t, 1, countScopedDiagnosticsLogs(all))

	payload := diagnosticsPayloadFromLog(t, all)
	assert.True(t, payload.MeteringMode)
	assert.Contains(t, payload.MeteringReason, "rates unavailable")
}

func Test_Report_BuildMeteringDiagnostics_ConcurrentWithSettle(t *testing.T) {
	t.Parallel()

	billingClient := mocks.NewBillingClient(t)
	billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).
		Return(incidentRateResponse(diagnosticTargetOrgID), nil).Maybe()
	billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
		Return(&successReserveResponse, nil).Maybe()
	billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
		Return(&emptypb.Empty{}, nil).Maybe()

	report, err := NewReport(t.Context(), targetOrgLabels(), logger.Nop(), billingClient, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
	require.NoError(t, err)
	require.NoError(t, report.Reserve(t.Context()))

	numRefs := 10
	for i := range numRefs {
		ref := "ref" + strconv.Itoa(i)
		_, err := report.Deduct(ref, ByResource(incidentGasUnit, "target-capability", decimal.NewFromInt(1)))
		require.NoError(t, err)
	}

	// diagnostics reads (FormatReport, rateCard, balance) must not race with
	// concurrent Settle writes; guarded by the existing report/balance locks.
	done := make(chan struct{})
	go func() {
		for range 100 {
			_ = report.buildMeteringDiagnostics(diagnosticTargetOrgID, diagOrgSourceBilling, nil, nil)
		}
		close(done)
	}()

	for i := range numRefs {
		go func(idx int) {
			ref := "ref" + strconv.Itoa(idx)
			_ = report.Settle(ref, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
				{Peer2PeerID: "peer1", SpendUnit: incidentGasUnit, SpendValue: incidentSpendValue},
			}})
		}(i)
	}

	<-done
}
