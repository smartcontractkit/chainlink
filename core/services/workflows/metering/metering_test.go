package metering

import (
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
)

const (
	testAccountID           = "accountId"
	testWorkflowID          = "workflowId"
	testWorkflowExecutionID = "workflowExecutionId"
	testUnitA               = "a"
	testUnitB               = "b"
)

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

	t.Run("Reserve returns an error if no billing client is given", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		err := report.Reserve(t.Context())
		require.ErrorIs(t, err, ErrNoBillingClient)
	})

	t.Run("Reserve turns on metering mode if the billing client cannot be communicated with", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(nil, errors.New("some err"))
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		require.True(t, report.balance.allowNegative)
	})

	t.Run("Reserve returns an error if insufficient funding", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: false}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.ErrorIs(t, err, ErrInsufficientFunding)
	})
}

func Test_Report_ConvertFromBalance(t *testing.T) {
	t.Parallel()

	t.Run("error if reserve is not called first", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		_, err := report.ConvertFromBalance("ref1", 1)
		require.ErrorIs(t, ErrNoReserve, err)
	})

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true, Rates: []*billing.ResourceUnitRate{{ResourceUnit: testUnitA, ConversionRate: "0.5"}}}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		amount, err := report.ConvertFromBalance(testUnitA, 1)
		require.NoError(t, err)
		require.Equal(t, int64(2), amount)
	})

	t.Run("falls back to 1:1 when rate is not found", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true, Rates: []*billing.ResourceUnitRate{{ResourceUnit: testUnitB, ConversionRate: "10"}}}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		amount, err := report.ConvertFromBalance(testUnitA, 1)
		require.NoError(t, err)
		require.Equal(t, int64(1), amount)
	})
}

func Test_Report_ConvertToBalance(t *testing.T) {
	t.Parallel()

	t.Run("error if reserve is not called first", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		_, err := report.ConvertToBalance("ref1", 1)
		require.ErrorIs(t, ErrNoReserve, err)
	})

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true, Rates: []*billing.ResourceUnitRate{{ResourceUnit: testUnitA, ConversionRate: "2"}}}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		amount, err := report.ConvertToBalance(testUnitA, 1)
		require.NoError(t, err)
		require.Equal(t, int64(2), amount)
	})

	t.Run("falls back to 1:1 when rate is not found", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true, Rates: []*billing.ResourceUnitRate{{ResourceUnit: testUnitB, ConversionRate: "10"}}}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		amount, err := report.ConvertToBalance(testUnitA, 1)
		require.NoError(t, err)
		require.Equal(t, int64(1), amount)
	})
}

func Test_Report_GetAvailableForInvocation(t *testing.T) {
	t.Parallel()

	t.Run("error if open slots is 0", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		_, err := report.GetAvailableForInvocation(0)
		require.ErrorIs(t, ErrNoOpenCalls, err)
	})

	t.Run("error if reserve is not called first", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		_, err := report.GetAvailableForInvocation(1)
		require.ErrorIs(t, ErrNoReserve, err)
	})

	t.Run("returns maxint64 in metering mode", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(nil, errors.New("nope"))
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		available, err := report.GetAvailableForInvocation(1)
		require.NoError(t, err)
		require.Equal(t, int64(math.MaxInt64), available)
	})

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true, Rates: []*billing.ResourceUnitRate{{ResourceUnit: testUnitA, ConversionRate: "2"}}}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		// 1 slot = all of available balance
		available, err := report.GetAvailableForInvocation(1)
		require.NoError(t, err)
		// TODO: https://smartcontract-it.atlassian.net/browse/CRE-290 once billing client response contains balance take out dummy balance
		require.Equal(t, int64(10000), available)
	})
}

func Test_Report_Deduct(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		err := report.Deduct("ref1", 1)
		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("returns an error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.Deduct("ref1", 2)
		require.NoError(t, err)
		err = report.Deduct("ref1", 1)
		require.ErrorIs(t, err, ErrStepDeductExists)
	})

	t.Run("does not modify local balance in metering mode", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(nil, errors.New("everything is on fire"))
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		balanceBefore := report.balance.balance
		err = report.Deduct("ref1", 2)
		require.NoError(t, err)
		balanceAfter := report.balance.balance
		require.Equal(t, balanceBefore, balanceAfter)
	})
}

func Test_Report_Settle(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		spendsByNode := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		err := report.Settle("ref1", spendsByNode)
		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("returns an error if Deduct is not called first", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		spendsByNode := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		require.ErrorIs(t, report.Settle("ref1", spendsByNode), ErrNoDeduct)
	})

	t.Run("returns an error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		require.NoError(t, report.Deduct("ref1", 2))
		require.NoError(t, report.Settle("ref1", steps))
		err = report.Settle("ref1", steps)
		require.ErrorIs(t, err, ErrStepSpendExists)
	})

	t.Run("ignores invalid spend values", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "????"},
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		err = report.Deduct("ref1", 2)
		require.NoError(t, err)
		require.NoError(t, report.Settle("ref1", steps))
	})

	t.Run("does not error when spend exceeds reservation", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "2"},
		}

		err = report.Deduct("ref1", 1)
		require.NoError(t, err)
		require.NoError(t, report.Settle("ref1", steps))
	})

	t.Run("does not modify local balance in metering mode", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(nil, errors.New("everything is still on fire"))
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		balanceBefore := report.balance.balance
		err = report.Deduct("ref1", 2)
		require.NoError(t, err)
		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "2"},
		}
		err = report.Settle("ref1", steps)
		require.NoError(t, err)
		balanceAfter := report.balance.balance
		require.Equal(t, balanceBefore, balanceAfter)
	})
}

func Test_Report_FormatReport(t *testing.T) {
	t.Parallel()

	t.Run("does not contain metadata", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		meteringReport := report.FormatReport()
		require.Equal(t, &events.WorkflowMetadata{}, meteringReport.Metadata)
	})

	t.Run("contains all step data", func(t *testing.T) {
		t.Parallel()
		numSteps := 100
		billingClient := mocks.NewBillingClient(t)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		err := report.Reserve(t.Context())
		require.NoError(t, err)

		expected := map[string]*events.MeteringReportStep{}

		for i := range numSteps {
			stepRef := strconv.Itoa(i)
			err := report.Deduct(stepRef, 1)
			require.NoError(t, err)
			spendsByNode := []capabilities.MeteringNodeDetail{
				{Peer2PeerID: "xyz", SpendUnit: "a", SpendValue: "42"},
			}
			err = report.Settle(stepRef, spendsByNode)
			require.NoError(t, err)
			expected[stepRef] = &events.MeteringReportStep{Nodes: []*events.MeteringReportNodeDetail{
				{
					Peer_2PeerId: "xyz",
					SpendUnit:    "a",
					SpendValue:   "42",
				},
			}}
		}

		require.Equal(t, expected, report.FormatReport().Steps)
	})
}

func Test_Report_SendReceipt(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.SendReceipt(t.Context())
		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("returns an error if unable to call billing client", func(t *testing.T) {
		t.Parallel()
		someErr := errors.New("error")
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(nil, someErr)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.SendReceipt(t.Context())
		require.ErrorIs(t, err, someErr)
	})

	t.Run("returns an error if billing client call is unsuccessful", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)

		// errors on nil response
		billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(nil, nil)
		err = report.SendReceipt(t.Context())
		require.ErrorIs(t, err, ErrReceiptFailed)

		// errors on unsuccessful response
		billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(&billing.SubmitWorkflowReceiptResponse{Success: false}, nil)
		err = report.SendReceipt(t.Context())
		require.ErrorIs(t, err, ErrReceiptFailed)
	})
}

func Test_MeterReports(t *testing.T) {
	t.Parallel()

	workflowExecutionID1 := "exec1"
	capabilityCall1 := "ref1"

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(&billing.SubmitWorkflowReceiptResponse{Success: true}, nil)
		mrs := NewReports(billingClient, testAccountID, testWorkflowID, logger.Test(t))
		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		require.NoError(t, err)
		err = r.Reserve(t.Context())
		require.NoError(t, err)
		err = r.Deduct(capabilityCall1, 1)
		require.NoError(t, err)
		err = r.Settle(capabilityCall1, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "1", SpendUnit: testUnitA, SpendValue: "0.8"},
			{Peer2PeerID: "2", SpendUnit: testUnitA, SpendValue: "0.9"},
			{Peer2PeerID: "3", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "4", SpendUnit: testUnitA, SpendValue: "1"},
		})
		require.NoError(t, err)
		err = mrs.End(t.Context(), workflowExecutionID1)
		require.NoError(t, err)
	})

	t.Run("happy path in metering mode", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(nil, errors.New("cannot"))
		billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(&billing.SubmitWorkflowReceiptResponse{Success: true}, nil)
		mrs := NewReports(billingClient, testAccountID, testWorkflowID, logger.Test(t))
		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		require.NoError(t, err)
		err = r.Reserve(t.Context())
		require.NoError(t, err)
		err = r.Deduct(capabilityCall1, 1)
		require.NoError(t, err)
		err = r.Settle(capabilityCall1, []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "1", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "2", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "3", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "4", SpendUnit: testUnitA, SpendValue: "1"},
		})
		require.NoError(t, err)
		err = mrs.End(t.Context(), workflowExecutionID1)
		require.NoError(t, err)
	})
}

func Test_MeterReports_Length(t *testing.T) {
	t.Parallel()

	billingClient := mocks.NewBillingClient(t)
	billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
	billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(&billing.SubmitWorkflowReceiptResponse{Success: true}, nil)
	mrs := NewReports(billingClient, "", "", logger.Test(t))

	_, err := mrs.Start(t.Context(), "exec1")
	require.NoError(t, err)
	mr, err := mrs.Start(t.Context(), "exec2")
	require.NoError(t, err)
	_, err = mrs.Start(t.Context(), "exec3")
	require.NoError(t, err)
	assert.Equal(t, 3, mrs.Len())

	err = mr.Reserve(t.Context())
	require.NoError(t, err)
	err = mrs.End(t.Context(), "exec2")
	require.NoError(t, err)
	assert.Equal(t, 2, mrs.Len())
}

func Test_MeterReports_Start(t *testing.T) {
	t.Parallel()

	t.Run("can only start report once", func(t *testing.T) {
		t.Parallel()

		billingClient := mocks.NewBillingClient(t)
		mrs := NewReports(billingClient, "", "", logger.Test(t))
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
		lggr := logger.Test(t)
		mrs := NewReports(billingClient, "", "", lggr)
		_, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		report, exists := mrs.Get("exec1")
		require.True(t, exists)
		require.NotEmpty(t, report)
	})
	t.Run("returns when no report exists", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		mrs := NewReports(billingClient, "", "", logger.Test(t))
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
		mrs := NewReports(billingClient, "", "", logger.Test(t))
		err := mrs.End(t.Context(), "exec1")
		require.ErrorIs(t, err, ErrReportNotFound)
	})

	t.Run("cleans up report on successful transmission to billing client", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(&billing.SubmitWorkflowReceiptResponse{Success: true}, nil)
		mrs := NewReports(billingClient, "", "", logger.Test(t))
		mr, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		require.Len(t, mrs.reports, 1)
		err = mr.Reserve(t.Context())
		require.NoError(t, err)
		err = mrs.End(t.Context(), "exec1")
		require.NoError(t, err)
		require.Empty(t, mrs.reports)
	})

	t.Run("cleans up report on failed transmission to billing client", func(t *testing.T) {
		t.Parallel()
		billingClient := mocks.NewBillingClient(t)
		billingClient.On("ReserveCredits", mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: true}, nil)
		billingClient.On("SubmitWorkflowReceipt", mock.Anything, mock.Anything).Return(nil, errors.New("errrrr"))
		mrs := NewReports(billingClient, "", "", logger.Test(t))
		mr, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		require.Len(t, mrs.reports, 1)
		err = mr.Reserve(t.Context())
		require.NoError(t, err)
		err = mrs.End(t.Context(), "exec1")
		require.Error(t, err)
		require.Empty(t, mrs.reports)
	})
}
