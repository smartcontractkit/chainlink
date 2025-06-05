package metering

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

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

func TestReport(t *testing.T) {
	t.Parallel()

	testA := "a"
	testUnitA := SpendUnit(testA)
	testB := "b"
	testUnitB := SpendUnit(testB)

	t.Run("MedianSpend returns median for multiple spend units", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "2"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "3"},
			{Peer2PeerID: "abc", SpendUnit: testB, SpendValue: "0.1"},
			{Peer2PeerID: "xyz", SpendUnit: testB, SpendValue: "0.2"},
			{Peer2PeerID: "abc", SpendUnit: testB, SpendValue: "0.3"},
		}

		for idx := range steps {
			err := report.Deduct(strconv.Itoa(idx), 1)
			require.NoError(t, err)
			require.NoError(t, report.Settle(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitB.IntToSpendValue(2),
			testUnitB: testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.2)),
		}

		median := report.MedianSpend()

		require.Len(t, median, 2)
		require.Contains(t, maps.Keys(median), testUnitA)
		require.Contains(t, maps.Keys(median), testUnitB)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
		assert.Equal(t, expected[testUnitB].String(), median[testUnitB].String())
	})

	t.Run("MedianSpend returns median single spend value", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: "a", SpendValue: "1"},
		}

		for idx := range steps {
			err := report.Deduct(strconv.Itoa(idx), 1)
			require.NoError(t, err)
			require.NoError(t, report.Settle(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitA.IntToSpendValue(1),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median odd number of spend values", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "3"},
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "2"},
		}

		for idx := range steps {
			err := report.Deduct(strconv.Itoa(idx), 1)
			require.NoError(t, err)
			require.NoError(t, report.Settle(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitA.IntToSpendValue(2),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median as average for even number of spend values", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "3"},
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "2"},
		}

		for idx := range steps {
			err := report.Deduct(strconv.Itoa(idx), 1)
			require.NoError(t, err)
			require.NoError(t, report.Settle(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitA.DecimalToSpendValue(decimal.NewFromFloat(2.5)),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("Initialize returns an error if no billing client is given", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), nil)
		report.Reserve(t.Context())
	})

	t.Run("Initialize allows negative balances if the billing client cannot be communicated with", func(t *testing.T) {
		t.Parallel()
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

	t.Run("Initialize returns an error if insufficient funding", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		billingClient.SetReserveCredits(&billing.ReserveCreditsResponse{Success: false}, nil)
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.Reserve(t.Context())
		require.ErrorIs(t, err, ErrInsufficientFunding)
	})

	t.Run("Deduct returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)

		err := report.Deduct("ref1", 1)
		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("Deduct returns an error if step already exists", func(t *testing.T) {
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

	t.Run("Settle returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
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
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
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
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
		}

		err = report.Deduct("ref1", 1)
		require.NoError(t, err)
		require.NoError(t, report.Settle("ref1", steps))
		require.ErrorIs(t, report.Settle("ref1", steps), ErrStepSpendExists)
	})

	t.Run("SendReceipt returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t), billingClient)
		err := report.SendReceipt(t.Context())
		require.ErrorIs(t, err, ErrNoReserve)
	})
}

// Test_MeterReports tests the Add, Get, Delete, and Len methods of a MeterReports.
// It also tests concurrent safe access.
func Test_MeterReports(t *testing.T) {
	billingClient := newMockBillingClient()
	mrs := NewReports(billingClient, "", "", logger.Test(t))
	assert.Equal(t, 0, mrs.Len())
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		mrs.Start(t.Context(), "exec1")
		r, ok := mrs.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.Deduct("ref1", 1)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.Settle("ref1", []capabilities.MeteringNodeDetail{})
		mrs.End(t.Context(), "exec1")
	}()
	go func() {
		defer wg.Done()
		report, err := mrs.Start(t.Context(), "exec2")
		require.NoError(t, err)
		r, ok := mrs.Get("exec2")
		assert.True(t, ok)
		err = report.Deduct("ref1", 1)
		require.NoError(t, err)
		err = r.Settle("ref1", []capabilities.MeteringNodeDetail{})
		require.NoError(t, err)
		mrs.End(t.Context(), "exec2")
	}()
	go func() {
		defer wg.Done()
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		mrs.Start(t.Context(), "exec1")
		r, ok := mrs.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.Deduct("ref1", 1)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.Settle("ref1", []capabilities.MeteringNodeDetail{})
		mrs.End(t.Context(), "exec1")
	}()

	wg.Wait()
	assert.Equal(t, 0, mrs.Len())
}

func Test_MeterReportsLength(t *testing.T) {
	mrs := NewReports(nil, "", "", logger.Test(t))

	mrs.Start(t.Context(), "exec1")
	mrs.Start(t.Context(), "exec2")
	mrs.Start(t.Context(), "exec3")
	assert.Equal(t, 3, mrs.Len())

	mrs.End(t.Context(), "exec2")
	assert.Equal(t, 2, mrs.Len())
}
