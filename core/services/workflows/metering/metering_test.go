package metering

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
)

const (
	testAccountID           = "accountId"
	testWorkflowID          = "workflowId"
	testWorkflowExecutionID = "workflowExecutionId"
)

type mockBillingClient struct {
	submitWorkflowReceiptResponse *billing.SubmitWorkflowReceiptResponse
	submitWorkflowReceiptError    error
	reserveCreditResponse         *billing.ReserveCreditsResponse
	reserveCreditError            error
}

type MockBillingClient interface {
	SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error)
	SetSubmitWorkflowReceipt(*billing.SubmitWorkflowReceiptResponse, error)
	ReserveCredits(context.Context, *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error)
	SetReserveCredits(*billing.ReserveCreditsResponse, error)
}

func (m *mockBillingClient) SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error) {
	return m.submitWorkflowReceiptResponse, m.submitWorkflowReceiptError
}
func (m *mockBillingClient) SetSubmitWorkflowReceipt(response *billing.SubmitWorkflowReceiptResponse, err error) {
	m.submitWorkflowReceiptResponse = response
	m.submitWorkflowReceiptError = err
}
func (m *mockBillingClient) ReserveCredits(context.Context, *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error) {
	return m.reserveCreditResponse, m.reserveCreditError
}
func (m *mockBillingClient) SetReserveCredits(response *billing.ReserveCreditsResponse, err error) {
	m.reserveCreditResponse = response
	m.reserveCreditError = err
}

func newMockBillingClient() MockBillingClient {
	return &mockBillingClient{
		submitWorkflowReceiptResponse: &billing.SubmitWorkflowReceiptResponse{Success: true},
		submitWorkflowReceiptError:    nil,
		reserveCreditResponse:         &billing.ReserveCreditsResponse{Success: true},
		reserveCreditError:            nil,
	}
}

func Test_medianSpend(t *testing.T) {
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
			assert.Equal(t, tc.expected, medianSpend(tc.input).String())
		})
	}

}

func Test_Report_Reserve(t *testing.T) {
	t.Parallel()

	t.Run("Reserve returns an error if no billing client is given", func(t *testing.T) {
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		err := report.Reserve(t.Context())
		require.ErrorIs(t, err, ErrNoBillingClient)
	})

	t.Run("Reserve allows negative balances if the billing client cannot be communicated with", func(t *testing.T) {
		billingClient := newMockBillingClient()
		billingClient.SetReserveCredits(nil, errors.New("some err"))
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.True(t, report.balance.allowNegative)
		require.NoError(t, err)
		err = report.Deduct("ref1", 1)
		require.NoError(t, err)
		require.Negative(t, report.balance.balance)
	})

	t.Run("Reserve returns an error if insufficient funding", func(t *testing.T) {
		billingClient := newMockBillingClient()
		billingClient.SetReserveCredits(&billing.ReserveCreditsResponse{Success: false}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.ErrorIs(t, err, ErrInsufficientFunding)
	})
}

func Test_Report_Deduct(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)

		err := report.Deduct("ref1", 1)
		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("returns an error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)
		err = report.Deduct("ref1", 2)
		require.NoError(t, err)
		err = report.Deduct("ref1", 1)
		require.ErrorIs(t, err, ErrStepDeductExists)
	})
}

func Test_Report_Settle(t *testing.T) {
	t.Parallel()
	testUnitA := "a"

	t.Run("Settle returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		err := report.Settle("ref1", steps)
		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("Settle returns an error if Deduct is not called first", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		require.ErrorIs(t, report.Settle("ref1", steps), ErrNoDeduct)
	})

	t.Run("Settle returns an error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}

		err = report.Deduct("ref1", 1)
		require.NoError(t, err)
		require.NoError(t, report.Settle("ref1", steps))
		require.ErrorIs(t, report.Settle("ref1", steps), ErrStepSpendExists)
	})
}

func Test_Report_SendReceipt(t *testing.T) {
	t.Parallel()

	t.Run("SendReceipt returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.SendReceipt(t.Context())
		require.ErrorIs(t, err, ErrNoReserve)
	})
}

type LabeledError struct {
	err   error
	label string
}

// Test_MeterReports tests the Add, Get, Delete, and Len methods of a MeterReports.
// It also tests concurrent safe access.
func Test_MeterReports(t *testing.T) {
	billingClient := newMockBillingClient()
	mrs := NewReports(billingClient, "", "", logger.Test(t))
	assert.Equal(t, 0, mrs.Len())

	workflowExecutionID1 := "exec1"
	workflowExecutionID2 := "exec2"

	goRoutines := 3
	actions := 4
	possibleErrs := goRoutines * actions
	errs := make(chan LabeledError, possibleErrs)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		err = r.Deduct("ref1", 1)
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		err = r.Settle("ref1", []capabilities.MeteringNodeDetail{})
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		err = mrs.End(t.Context(), workflowExecutionID1)
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		wg.Done()
	}()

	go func() {
		r, err := mrs.Start(t.Context(), workflowExecutionID2)
		errs <- LabeledError{err: err, label: workflowExecutionID2}
		err = r.Deduct("ref1", 1)
		errs <- LabeledError{err: err, label: workflowExecutionID2}
		err = r.Settle("ref1", []capabilities.MeteringNodeDetail{})
		errs <- LabeledError{err: err, label: workflowExecutionID2}
		err = mrs.End(t.Context(), workflowExecutionID2)
		errs <- LabeledError{err: err, label: workflowExecutionID2}
		wg.Done()
	}()

	go func() {
		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		err = r.Deduct("ref1", 1)
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		err = r.Settle("ref1", []capabilities.MeteringNodeDetail{})
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		err = mrs.End(t.Context(), workflowExecutionID1)
		errs <- LabeledError{err: err, label: workflowExecutionID1}
		wg.Done()
	}()

	wg.Wait()

	// // wait for all to finish
	for i := 0; i < possibleErrs; i++ {
		lErr := <-errs
		require.NoError(t, lErr.err)
	}
}

func Test_MeterReports_Length(t *testing.T) {
	billingClient := newMockBillingClient()
	mrs := NewReports(billingClient, "", "", logger.Test(t))

	_, err := mrs.Start(t.Context(), "exec1")
	require.NoError(t, err)
	_, err = mrs.Start(t.Context(), "exec2")
	require.NoError(t, err)
	_, err = mrs.Start(t.Context(), "exec3")
	require.NoError(t, err)
	assert.Equal(t, 3, mrs.Len())

	err = mrs.End(t.Context(), "exec2")
	require.NoError(t, err)
	assert.Equal(t, 2, mrs.Len())
}

func Test_MeterReports_Start(t *testing.T) {
	t.Run("can only start report once", func(t *testing.T) {
		billingClient := newMockBillingClient()
		mrs := NewReports(billingClient, "", "", logger.Test(t))
		_, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		_, err = mrs.Start(t.Context(), "exec1")
		require.ErrorIs(t, err, ErrReportExists)
	})
}

func Test_MeterReports_End(t *testing.T) {
	t.Run("can only end existing report", func(t *testing.T) {
		billingClient := newMockBillingClient()
		mrs := NewReports(billingClient, "", "", logger.Test(t))
		_, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		err = mrs.End(t.Context(), "exec1")
		require.NoError(t, err)
		err = mrs.End(t.Context(), "exec1")
		require.ErrorIs(t, err, ErrReportNotFound)
	})
}
