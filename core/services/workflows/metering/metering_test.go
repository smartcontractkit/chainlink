package metering

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	eventspb "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

const (
	testAccountID           = "accountId"
	testWorkflowID          = "workflowId"
	testWorkflowExecutionID = "workflowExecutionId"
)

var (
	successReserveResponse = billing.ReserveCreditsResponse{
		Success: true,
		Credits: 10_000,
	}
	successReserveResponseWithRates = billing.ReserveCreditsResponse{
		Success: true,
		Entries: []*billing.RateCardEntry{
			{
				ResourceType:    billing.ResourceType_RESOURCE_TYPE_COMPUTE,
				MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS,
				UnitsPerCredit:  "2",
			},
		},
		Credits: 10_000,
	}
	successReserveResponseWithMultiRates = billing.ReserveCreditsResponse{Success: true, Entries: []*billing.RateCardEntry{
		{
			ResourceType:    billing.ResourceType_RESOURCE_TYPE_COMPUTE,
			MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS,
			UnitsPerCredit:  "2",
		},
		{
			ResourceType:    billing.ResourceType_RESOURCE_TYPE_GAS,
			MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_COST,
			UnitsPerCredit:  "3",
		},
	}, Credits: 10_000}
	failureReserveResponse = billing.ReserveCreditsResponse{
		Success: false,
	}
	defaultLabels = map[string]string{
		platform.KeyWorkflowOwner:       "accountId",
		platform.KeyWorkflowID:          "workflowId",
		platform.KeyWorkflowExecutionID: "workflowExecutionId",
	}
	testUnitA      = billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_COMPUTE)]
	testUnitB      = billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_UNSPECIFIED)]
	testUnitC      = billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_GAS)]
	validConfig, _ = values.NewMap(map[string]any{
		RatiosKey: map[string]string{
			testUnitA: "0.4",
			testUnitB: "0.6",
		},
	})
)

func Test_Report(t *testing.T) {
	t.Parallel()

	t.Run("error if incorrect labels", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		_, err := NewReport(map[string]string{}, logger.Nop(), billingClient)
		require.ErrorIs(t, err, ErrMissingLabels)
	})
}

func Test_Report_MeteringMode(t *testing.T) {
	t.Parallel()

	t.Run("Reserve switches to metering mode", func(t *testing.T) {
		t.Parallel()

		t.Run("if billing client is nil", func(t *testing.T) {
			t.Parallel()

			report := newTestReport(t, logger.Nop(), nil)

			require.NoError(t, report.Reserve(t.Context()))
			assert.True(t, report.meteringMode)
		})

		t.Run("if billing client returns an error", func(t *testing.T) {
			t.Parallel()

			billingClient := mocks.NewBillingClient(t)
			report := newTestReport(t, logger.Nop(), billingClient)

			billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(nil, errors.New("some err"))
			require.NoError(t, report.Reserve(t.Context()))
			require.True(t, report.meteringMode)
			billingClient.AssertExpectations(t)
		})

		t.Run("if rate card contains invalid entry", func(t *testing.T) {
			t.Parallel()

			lggr, logs := logger.TestObserved(t, zapcore.WarnLevel)
			billingClient := mocks.NewBillingClient(t)
			report := newTestReport(t, lggr, billingClient)

			billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
				Return(&billing.ReserveCreditsResponse{Success: true, Entries: []*billing.RateCardEntry{
					{ResourceType: billing.ResourceType_RESOURCE_TYPE_COMPUTE, UnitsPerCredit: "invalid"},
				}, Credits: 10_000}, nil)
			require.NoError(t, report.Reserve(t.Context()))
			require.True(t, report.meteringMode)
			assert.Len(t, logs.All(), 1)
			billingClient.AssertExpectations(t)
		})
	})

	t.Run("ConvertToBalance falls back to 1:1 when rate is not found and switches to metering mode", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
		report := newTestReport(t, lggr, billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&billing.ReserveCreditsResponse{Success: true, Entries: []*billing.RateCardEntry{
				{ResourceType: billing.ResourceType_RESOURCE_TYPE_UNSPECIFIED, MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS, UnitsPerCredit: "10"},
			}}, nil)
		require.NoError(t, report.Reserve(t.Context()))

		amount, err := report.ConvertToBalance(testUnitA, decimal.NewFromInt(1))
		require.NoError(t, err)
		assert.True(t, amount.Equal(decimal.NewFromInt(1)))
		require.True(t, report.meteringMode)
		assert.Len(t, logs.All(), 1)
		billingClient.AssertExpectations(t)
	})

	t.Run("GetMaxSpendForInvocation returns null decimal in metering mode", func(t *testing.T) {
		t.Parallel()

		emptyUserSpendLimit := decimal.NewNullDecimal(decimal.Zero)
		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(nil, errors.New("nope"))

		require.NoError(t, report.Reserve(t.Context()))

		available, err := report.GetMaxSpendForInvocation(emptyUserSpendLimit, 1)
		require.NoError(t, err)
		assert.False(t, available.Valid)
		billingClient.AssertExpectations(t)
	})

	two := decimal.NewFromInt(2)

	t.Run("Deduct does not modify local balance in metering mode", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(nil, errors.New("everything is on fire"))
		require.NoError(t, report.Reserve(t.Context()))

		balanceBefore := report.balance.balance
		require.NoError(t, report.Deduct("ref1", two))

		balanceAfter := report.balance.balance
		assert.Equal(t, balanceBefore, balanceAfter)
		billingClient.AssertExpectations(t)
	})

	t.Run("Settle does not modify local balance in metering mode", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		// trigger metering mode with a billing reserve error
		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(nil, errors.New("everything is still on fire"))
		require.NoError(t, report.Reserve(t.Context()))

		balanceBefore := report.balance.balance

		require.NoError(t, report.Deduct("ref1", two))

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "2"},
		}
		require.NoError(t, report.Settle("ref1", steps))

		balanceAfter := report.balance.balance
		require.Equal(t, balanceBefore, balanceAfter)
	})

	t.Run("CreditToSpendingLimits switches to metering mode", func(t *testing.T) {
		t.Parallel()

		t.Run("if only one spend type and rate does not exist", func(t *testing.T) {
			t.Parallel()

			lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
			billingClient := mocks.NewBillingClient(t)
			report := newTestReport(t, lggr, billingClient)

			billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
				Return(&successReserveResponseWithMultiRates, nil)
			require.NoError(t, report.Reserve(t.Context()))

			// ratios and spend types should match
			config, _ := values.NewMap(map[string]any{
				RatiosKey: map[string]string{
					testUnitB: "1",
				},
			})

			// trigger metering mode spending type that doesn't match rates in reserve response
			limits := report.CreditToSpendingLimits(capabilities.CapabilityInfo{
				SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitB)},
			}, config, decimal.NewFromInt(1_000))

			assert.Empty(t, limits)
			assert.True(t, report.meteringMode)
			assert.Len(t, logs.All(), 1)
			billingClient.AssertExpectations(t)
		})

		t.Run("if ratio and spend type lengths do not match", func(t *testing.T) {
			t.Parallel()

			lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
			billingClient := mocks.NewBillingClient(t)
			report := newTestReport(t, lggr, billingClient)

			billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
				Return(&successReserveResponseWithRates, nil)
			require.NoError(t, report.Reserve(t.Context()))

			// 3 spend types and 2 ratios creates the mismatch
			limits := report.CreditToSpendingLimits(capabilities.CapabilityInfo{
				SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA), capabilities.CapabilitySpendType(testUnitB), capabilities.CapabilitySpendType(testUnitC)},
			}, validConfig, decimal.NewFromInt(1_000))

			assert.Empty(t, limits)
			assert.True(t, report.meteringMode)
			assert.Len(t, logs.All(), 1)
			billingClient.AssertExpectations(t)
		})

		t.Run("if multiple spend types and ratio does not exist", func(t *testing.T) {
			t.Parallel()

			lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
			billingClient := mocks.NewBillingClient(t)
			report := newTestReport(t, lggr, billingClient)

			billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
				Return(&successReserveResponseWithMultiRates, nil)
			require.NoError(t, report.Reserve(t.Context()))

			// spend types and rates should match
			// spend types and ratios should not match and return an error
			limits := report.CreditToSpendingLimits(capabilities.CapabilityInfo{
				SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA), capabilities.CapabilitySpendType(testUnitC)},
			}, validConfig, decimal.NewFromInt(1_000))

			assert.Empty(t, limits)
			assert.True(t, report.meteringMode)
			assert.Len(t, logs.All(), 1)
			billingClient.AssertExpectations(t)
		})

		t.Run("if multiple spend types and rate does not exist", func(t *testing.T) {
			t.Parallel()

			lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
			billingClient := mocks.NewBillingClient(t)
			report := newTestReport(t, lggr, billingClient)

			billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
				Return(&successReserveResponseWithMultiRates, nil)
			require.NoError(t, report.Reserve(t.Context()))

			// ratios for spend types should match
			// rates for spend types should not match
			limits := report.CreditToSpendingLimits(capabilities.CapabilityInfo{
				SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA), capabilities.CapabilitySpendType(testUnitB)},
			}, validConfig, decimal.NewFromInt(1_000))

			assert.Empty(t, limits)
			assert.True(t, report.meteringMode)
			assert.Len(t, logs.All(), 1)
			billingClient.AssertExpectations(t)
		})
	})
}

func Test_medianSpend(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    []decimal.Decimal
		expected string
	}{
		{
			name: "MedianSpend returns median for a list of int spend values",
			input: []decimal.Decimal{
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			expected: "2",
		},
		{
			name: "MedianSpend returns median for a list of float spend values",
			input: []decimal.Decimal{
				decimal.NewFromFloat(0.1),
				decimal.NewFromFloat(0.2),
				decimal.NewFromFloat(0.3),
			},
			expected: "0.2",
		},
		{
			name: "MedianSpend returns median single spend value",
			input: []decimal.Decimal{
				decimal.NewFromInt(1),
			},
			expected: "1",
		},
		{
			name: "MedianSpend returns median even number of spend values",
			input: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(2),
				decimal.NewFromInt(4),
				decimal.NewFromInt(4),
			},
			expected: "3",
		},
		{
			name: "MedianSpend returns median odd number of spend values",
			input: []decimal.Decimal{
				decimal.NewFromInt(1),
				decimal.NewFromInt(13),
				decimal.NewFromInt(50),
				decimal.NewFromInt(51),
				decimal.NewFromInt(100),
			},
			expected: "50",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, medianSpend(tc.input).String())
		})
	}
}

func Test_Report_Reserve(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if insufficient funding", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		lggr, logs := logger.TestObserved(t, zapcore.WarnLevel)
		report := newTestReport(t, lggr, billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&failureReserveResponse, nil)
		require.ErrorIs(t, report.Reserve(t.Context()), ErrInsufficientFunding)
		assert.False(t, report.meteringMode)
		assert.Empty(t, logs.All())
		billingClient.AssertExpectations(t)
	})

	t.Run("returns no error on success", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		lggr, logs := logger.TestObserved(t, zapcore.WarnLevel)
		report := newTestReport(t, lggr, billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)
		require.NoError(t, report.Reserve(t.Context()))
		assert.False(t, report.meteringMode)
		assert.Empty(t, logs.All())
		billingClient.AssertExpectations(t)
	})
}

func Test_Report_ConvertToBalance(t *testing.T) {
	t.Parallel()

	one, two := decimal.NewFromInt(1), decimal.NewFromInt(2)

	t.Run("error if reserve is not called first", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)
		_, err := report.ConvertToBalance("ref1", one)

		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&billing.ReserveCreditsResponse{Success: true, Entries: []*billing.RateCardEntry{
				{ResourceType: billing.ResourceType_RESOURCE_TYPE_COMPUTE, MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS, UnitsPerCredit: "2"},
			}}, nil)

		require.NoError(t, report.Reserve(t.Context()))

		amount, err := report.ConvertToBalance(testUnitA, one)
		require.NoError(t, err)
		assert.True(t, amount.Equal(two))
		billingClient.AssertExpectations(t)
	})
}

func Test_Report_GetAvailableForInvocation(t *testing.T) {
	t.Parallel()

	emptyUserSpendLimit := decimal.NewNullDecimal(decimal.Zero)
	emptyUserSpendLimit.Valid = false

	t.Run("error if open slots is 0", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)
		_, err := report.GetMaxSpendForInvocation(emptyUserSpendLimit, 0)

		require.ErrorIs(t, ErrNoOpenCalls, err)
	})

	t.Run("error if reserve is not called first", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)
		_, err := report.GetMaxSpendForInvocation(emptyUserSpendLimit, 1)

		require.ErrorIs(t, ErrNoReserve, err)
	})

	t.Run("happy path without user-defined spending limit", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithRates, nil)

		require.NoError(t, report.Reserve(t.Context()))

		// 1 slot = all of available balance
		available, err := report.GetMaxSpendForInvocation(emptyUserSpendLimit, 1)
		require.NoError(t, err)

		// TODO: https://smartcontract-it.atlassian.net/browse/CRE-290 once billing client response contains balance take out dummy balance
		assert.True(t, available.Decimal.Equal(decimal.NewFromInt(10_000)),
			"unexpected available balance %s", available.Decimal.String())
		assert.True(t, available.Valid, "available value should be a valid amount")
		billingClient.AssertExpectations(t)
	})

	t.Run("happy path with user-defined spending limit", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithRates, nil)

		require.NoError(t, report.Reserve(t.Context()))

		// 1 slot = all of available balance
		nonEmptyUserSpendLimit := decimal.NewNullDecimal(decimal.NewFromInt(5_000))
		nonEmptyUserSpendLimit.Valid = true
		available, err := report.GetMaxSpendForInvocation(nonEmptyUserSpendLimit, 1)
		require.NoError(t, err)

		// TODO: https://smartcontract-it.atlassian.net/browse/CRE-290 once billing client response contains balance take out dummy balance
		assert.True(t, available.Decimal.Equal(decimal.NewFromInt(5_000)), available.Decimal.String())
		assert.True(t, available.Valid, "available value should be a valid amount")
		billingClient.AssertExpectations(t)
	})
}

func Test_Report_Deduct(t *testing.T) {
	t.Parallel()

	one := decimal.NewFromInt(1)
	two := decimal.NewFromInt(2)

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)

		require.ErrorIs(t, report.Deduct("ref1", one), ErrNoReserve)
	})

	t.Run("returns an error if step already exists", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)
		require.NoError(t, report.Reserve(t.Context()))

		require.NoError(t, report.Deduct("ref1", two))
		require.ErrorIs(t, report.Deduct("ref1", one), ErrStepDeductExists)
		billingClient.AssertExpectations(t)
	})

	t.Run("returns insufficient balance when not in metering mode", func(t *testing.T) {
		t.Parallel()

		deductValue := decimal.NewFromInt(11_000)
		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithRates, nil)
		require.NoError(t, report.Reserve(t.Context()))

		require.ErrorIs(t, report.Deduct("ref1", deductValue), ErrInsufficientBalance)
		billingClient.AssertExpectations(t)
	})
}

func Test_Report_Settle(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)

		require.ErrorIs(t, report.Settle("ref1", []capabilities.MeteringNodeDetail{}), ErrNoReserve)
	})

	t.Run("returns an error if Deduct is not called first", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)
		require.NoError(t, report.Reserve(t.Context()))
		require.ErrorIs(t, report.Settle("ref1", []capabilities.MeteringNodeDetail{}), ErrNoDeduct)
		billingClient.AssertExpectations(t)
	})

	t.Run("returns an error if step already exists", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)
		require.NoError(t, report.Reserve(t.Context()))

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		require.NoError(t, report.Deduct("ref1", decimal.NewFromInt(2)))
		require.NoError(t, report.Settle("ref1", steps))
		require.ErrorIs(t, report.Settle("ref1", steps), ErrStepSpendExists)
		billingClient.AssertExpectations(t)
	})

	t.Run("ignores invalid spend values", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
		report := newTestReport(t, lggr, billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithRates, nil)
		require.NoError(t, report.Reserve(t.Context()))

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "????"},
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		require.NoError(t, report.Deduct("ref1", decimal.NewFromInt(2)))
		require.NoError(t, report.Settle("ref1", steps))
		assert.Len(t, logs.All(), 1)
		billingClient.AssertExpectations(t)
	})

	t.Run("does not error when spend exceeds reservation", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
		report := newTestReport(t, lggr, billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithRates, nil)
		require.NoError(t, report.Reserve(t.Context()))

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "2"},
		}

		require.NoError(t, report.Deduct("ref1", decimal.NewFromInt(1)))
		require.NoError(t, report.Settle("ref1", steps))
		assert.Len(t, logs.All(), 1)
		billingClient.AssertExpectations(t)
	})
}

func Test_Report_FormatReport(t *testing.T) {
	t.Parallel()

	t.Run("does not contain metadata", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
		require.NoError(t, report.Reserve(t.Context()))

		meteringReport := report.FormatReport()
		require.Equal(t, &eventspb.WorkflowMetadata{}, meteringReport.Metadata)
		billingClient.AssertExpectations(t)
	})

	t.Run("contains all step data", func(t *testing.T) {
		t.Parallel()

		numSteps := 100
		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&successReserveResponse, nil)
		require.NoError(t, report.Reserve(t.Context()))

		expected := map[string]*eventspb.MeteringReportStep{}

		for i := range numSteps {
			stepRef := strconv.Itoa(i)

			require.NoError(t, report.Deduct(stepRef, decimal.NewFromInt(1)))
			require.NoError(t, report.Settle(stepRef, []capabilities.MeteringNodeDetail{
				{Peer2PeerID: "xyz", SpendUnit: "a", SpendValue: "42"},
			}))

			expected[stepRef] = &eventspb.MeteringReportStep{Nodes: []*eventspb.MeteringReportNodeDetail{
				{
					Peer_2PeerId: "xyz",
					SpendUnit:    "a",
					SpendValue:   "42",
				},
			}}
		}

		assert.Equal(t, expected, report.FormatReport().Steps)
		billingClient.AssertExpectations(t)
	})
}

func Test_Report_SendReceipt(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)

		require.ErrorIs(t, report.SendReceipt(t.Context()), ErrNoReserve)
	})

	t.Run("returns an error billing client not set", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)

		require.NoError(t, report.Reserve(t.Context()))
		require.ErrorIs(t, report.SendReceipt(t.Context()), ErrNoBillingClient)
	})

	t.Run("returns an error if unable to call billing client", func(t *testing.T) {
		t.Parallel()

		someErr := errors.New("error")
		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(nil, someErr)

		require.NoError(t, report.Reserve(t.Context()))
		require.ErrorIs(t, report.SendReceipt(t.Context()), someErr)
		billingClient.AssertExpectations(t)
	})

	t.Run("returns an error if billing client call is unsuccessful", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)

		require.NoError(t, report.Reserve(t.Context()))

		// errors on nil response
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(nil, nil)
		require.ErrorIs(t, report.SendReceipt(t.Context()), ErrReceiptFailed)

		// errors on unsuccessful response
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
			Return(&emptypb.Empty{}, nil)
		require.ErrorIs(t, report.SendReceipt(t.Context()), ErrReceiptFailed)

		billingClient.AssertExpectations(t)
	})
}

func Test_Report_EmitReceipt(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// No parallel
		beholderTester := beholdertest.NewObserver(t)
		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithRates, nil)
		require.NoError(t, report.Reserve(t.Context()))

		require.NoError(t, report.EmitReceipt(t.Context()))

		assert.Equal(t, 1, beholderTester.Len(t, "beholder_entity", fmt.Sprintf("%s.%s", events.ProtoPkg, events.MeteringReportEntity)))

		messages := beholderTester.Messages(t, "beholder_entity", fmt.Sprintf("%s.%s", events.ProtoPkg, events.MeteringReportEntity))

		for _, msg := range messages {
			entity := msg.Attrs["beholder_entity"]
			if entity == fmt.Sprintf("%s.%s", events.ProtoPkg, events.MeteringReportEntity) {
				var report eventspb.MeteringReport
				require.NoError(t, proto.Unmarshal(msg.Body, &report))
				assert.Equal(t, testWorkflowID, report.Metadata.WorkflowID)
				assert.NotEmpty(t, report.Metadata.WorkflowExecutionID)
				assert.Equal(t, testAccountID, report.Metadata.WorkflowOwner)
			}
		}
	})

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		require.ErrorIs(t, report.EmitReceipt(t.Context()), ErrNoReserve)
	})
}

func Test_Report_CreditToSpendingLimits(t *testing.T) {
	t.Parallel()

	t.Run("happy path splits spend types per provided ratios", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		report := newTestReport(t, logger.Nop(), billingClient)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithMultiRates, nil)

		require.NoError(t, report.Reserve(t.Context()))

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]string{
				testUnitA: "0.4",
				testUnitC: "0.6",
			},
		})

		limits := report.CreditToSpendingLimits(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA), capabilities.CapabilitySpendType(testUnitC)},
		}, config, decimal.NewFromInt(1_000))

		require.NotNil(t, limits)
		require.Len(t, limits, 2)
		assert.Equal(t, testUnitA, string(limits[0].SpendType))
		assert.Equal(t, testUnitC, string(limits[1].SpendType))
		assert.Equal(t, "200.000", limits[0].Limit) // conversion rate of 2 at 40% ratio
		assert.Equal(t, "200.000", limits[1].Limit) // conversion rate of 3 at 60% ratio
		assert.False(t, report.meteringMode)
	})

	t.Run("empty limits for empty spend types", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop(), nil)
		limits := report.CreditToSpendingLimits(
			capabilities.CapabilityInfo{},
			validConfig,
			decimal.NewFromInt(1_000),
		)

		assert.Empty(t, limits)
	})
}

func Test_MeterReports(t *testing.T) {
	t.Parallel()

	workflowExecutionID1 := "exec1"
	capabilityCall1 := "ref1"

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, testAccountID, testWorkflowID, logger.Nop(), defaultLabels, metrics)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponseWithRates, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
			Return(&emptypb.Empty{}, nil)

		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		require.NoError(t, err)

		require.NoError(t, r.Reserve(t.Context()))
		require.NoError(t, r.Deduct(capabilityCall1, decimal.NewFromInt(1)))
		require.NoError(t, r.Settle(capabilityCall1, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "1", SpendUnit: testUnitA, SpendValue: "0.8"},
			{Peer2PeerID: "2", SpendUnit: testUnitA, SpendValue: "0.9"},
			{Peer2PeerID: "3", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "4", SpendUnit: testUnitA, SpendValue: "1"},
		}))
		require.NoError(t, mrs.End(t.Context(), workflowExecutionID1))
		billingClient.AssertExpectations(t)
	})

	t.Run("happy path in metering mode", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, testAccountID, testWorkflowID, logger.Nop(), defaultLabels, metrics)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(nil, errors.New("cannot"))
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
			Return(&emptypb.Empty{}, nil)

		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		require.NoError(t, err)

		require.NoError(t, r.Reserve(t.Context()))
		require.NoError(t, r.Deduct(capabilityCall1, decimal.NewFromInt(1)))
		require.NoError(t, r.Settle(capabilityCall1, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "1", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "2", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "3", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "4", SpendUnit: testUnitA, SpendValue: "1"},
		}))
		require.NoError(t, mrs.End(t.Context(), workflowExecutionID1))
		billingClient.AssertExpectations(t)
	})
}

func Test_MeterReports_Length(t *testing.T) {
	t.Parallel()

	billingClient := mocks.NewBillingClient(t)
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	metrics := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
	mrs := NewReports(billingClient, "", "", logger.Nop(), defaultLabels, metrics)

	billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
		Return(&successReserveResponse, nil)
	billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
		Return(&emptypb.Empty{}, nil)

	_, err = mrs.Start(t.Context(), "exec1")
	require.NoError(t, err)

	mr, err := mrs.Start(t.Context(), "exec2")
	require.NoError(t, err)

	_, err = mrs.Start(t.Context(), "exec3")
	require.NoError(t, err)
	assert.Equal(t, 3, mrs.Len())

	require.NoError(t, mr.Reserve(t.Context()))
	require.NoError(t, mrs.End(t.Context(), "exec2"))
	assert.Equal(t, 2, mrs.Len())
}

func Test_MeterReports_Start(t *testing.T) {
	t.Parallel()

	t.Run("can only start report once", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, "", "", logger.Nop(), defaultLabels, metrics)

		_, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)

		_, err = mrs.Start(t.Context(), "exec1")
		require.ErrorIs(t, err, ErrReportExists)
	})
}

func Test_MeterReports_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns when report exists", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, "", "", logger.Nop(), defaultLabels, metrics)

		_, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)

		report, exists := mrs.Get("exec1")
		require.True(t, exists)
		require.NotEmpty(t, report)
	})

	t.Run("returns when no report exists", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, "", "", logger.Nop(), defaultLabels, metrics)

		report, exists := mrs.Get("exec1")
		require.False(t, exists)
		require.Nil(t, report)
	})
}

func Test_MeterReports_End(t *testing.T) {
	t.Parallel()

	t.Run("can only end existing report", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, "", "", logger.Nop(), defaultLabels, metrics)

		require.ErrorIs(t, mrs.End(t.Context(), "exec1"), ErrReportNotFound)
	})

	t.Run("cleans up report on successful transmission to billing client", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, "", "", logger.Nop(), defaultLabels, metrics)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
			Return(&emptypb.Empty{}, nil)

		mr, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		assert.Len(t, mrs.reports, 1)

		require.NoError(t, mr.Reserve(t.Context()))
		require.NoError(t, mrs.End(t.Context(), "exec1"))
		assert.Empty(t, mrs.reports)
		billingClient.AssertExpectations(t)
	})

	t.Run("cleans up report on failed transmission to billing client", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		metrics := defaultMetrics(t)
		mrs := NewReports(billingClient, "", "", logger.Nop(), defaultLabels, metrics)

		billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).
			Return(&successReserveResponse, nil)
		billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).
			Return(nil, errors.New("errrrr"))

		mr, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		assert.Len(t, mrs.reports, 1)

		require.NoError(t, mr.Reserve(t.Context()))
		require.Error(t, mrs.End(t.Context(), "exec1"))
		assert.Empty(t, mrs.reports)
		billingClient.AssertExpectations(t)
	})
}

func TestRatiosFromConfig(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitB),
			},
		}, validConfig)

		require.NoError(t, err)
		require.Len(t, ratios, 2)

		assert.Contains(t, ratios, capabilities.CapabilitySpendType(testUnitA))
		assert.Contains(t, ratios, capabilities.CapabilitySpendType(testUnitB))
	})

	t.Run("automatic ratio of 1 for single spend type", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA)},
		}, nil)

		require.NoError(t, err)
		require.Len(t, ratios, 1)

		assert.Contains(t, ratios, capabilities.CapabilitySpendType(testUnitA))
	})

	t.Run("error when spend ratios key does not exist", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{}, new(values.Map))
		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when spend ratios fails to unwrap to map", func(t *testing.T) {
		t.Parallel()

		config := &values.Map{
			Underlying: map[string]values.Value{
				"spendRatios": &values.String{Underlying: "invalid"},
			},
		}

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{}, config)
		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when spend type is not in ratios map", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitC),
			},
		}, validConfig)

		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when ratio contains invalid data type", func(t *testing.T) {
		t.Parallel()

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]any{
				testUnitA: "0.2",
				testUnitB: 5,
			},
		})

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitB),
			},
		}, config)

		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when ratio contains invalid decimal", func(t *testing.T) {
		t.Parallel()

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]any{
				testUnitA: "invalid",
				testUnitB: "0.2",
			},
		})

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitB),
			},
		}, config)

		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})
}

func newTestReport(t *testing.T, lggr logger.Logger, client *mocks.BillingClient) *Report {
	t.Helper()

	if client == nil {
		meteringReport, err := NewReport(defaultLabels, lggr, nil)
		require.NoError(t, err)

		return meteringReport
	}

	meteringReport, err := NewReport(defaultLabels, lggr, client)
	require.NoError(t, err)

	return meteringReport
}

func defaultMetrics(t *testing.T) *monitoring.WorkflowsMetricLabeler {
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)

	return monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
}
