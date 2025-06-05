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
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
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
			_, err := report.DeductByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
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
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: "a", SpendValue: "1"},
		}

		for idx := range steps {
			_, err := report.DeductByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
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
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "3"},
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "2"},
		}

		for idx := range steps {
			_, err := report.DeductByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
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
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
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
			_, err := report.DeductByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
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
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		require.ErrorIs(t, report.StartAndReserve(t.Context()), ErrNoBillingClient)
	})

	t.Run("Initialize returns an error if no owner is given", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		billingClient := newMockBillingClient()
		report.client = billingClient
		require.ErrorIs(t, report.StartAndReserve(t.Context()), ErrNoOwner)
	})

	t.Run("Initialize returns an error if no workflow ID is given", func(t *testing.T) {
		t.Parallel()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		billingClient := newMockBillingClient()
		report.client = billingClient
		report.owner = testAccountID
		require.ErrorIs(t, report.StartAndReserve(t.Context()), ErrNoWorkflowID)
	})

	t.Run("Initialize allows negative balances if the billing client cannot be communicated with", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		billingClient.SetReserveCredits(nil, errors.New("some err"))
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.True(t, report.balance.allowNegative)
		require.NoError(t, err)
		_, err = report.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.NoError(t, err)
		require.Negative(t, report.balance.balance)
	})

	t.Run("Initialize returns an error if insufficient funding", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		billingClient.SetReserveCredits(&billing.ReserveCreditsResponse{Success: false}, nil)
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.ErrorIs(t, err, ErrInsufficientFunding)
	})

	t.Run("DeductByLimits returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID

		_, err := report.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.ErrorIs(t, err, ErrUninitializedReport)
	})

	t.Run("DeductByLimits returns an error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)
		_, err = report.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.NoError(t, err)
		_, err = report.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.ErrorIs(t, err, ErrStepDeductExists)
	})

	t.Run("DeductByAvailability returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID

		_, err := report.DeductByAvailability("ref1", capabilities.CapabilityInfo{}, 2, 0)
		require.ErrorIs(t, err, ErrUninitializedReport)
	})

	t.Run("DeductByAvailability returns an error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)
		_, err = report.DeductByAvailability("ref1", capabilities.CapabilityInfo{}, 2, 0)
		require.NoError(t, err)
		_, err = report.DeductByAvailability("ref1", capabilities.CapabilityInfo{}, 1, 0)
		require.ErrorIs(t, err, ErrStepDeductExists)
	})

	t.Run("SetStep returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
		}

		err := report.SetStep("ref1", steps)
		require.ErrorIs(t, err, ErrUninitializedReport)
	})

	t.Run("SetStep returns an error if Deduct is not called first", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
		}

		require.ErrorIs(t, report.SetStep("ref1", steps), ErrNoDeduct)
	})

	t.Run("SetStep returns an error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		err := report.StartAndReserve(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
		}

		_, err = report.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.NoError(t, err)
		require.NoError(t, report.SetStep("ref1", steps))
		require.ErrorIs(t, report.SetStep("ref1", steps), ErrStepSpendExists)
	})

	t.Run("SendReceipt returns an error if not initialized", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID

		err := report.SendReceipt(t.Context())
		require.ErrorIs(t, err, ErrUninitializedReport)
	})
}

// Test_MeterReports tests the Add, Get, Delete, and Len methods of a MeterReports.
// It also tests concurrent safe access.
func Test_MeterReports(t *testing.T) {
	mr := NewReports(nil, "", "")
	billingClient := newMockBillingClient()
	assert.Equal(t, 0, mr.Len())
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		report.StartAndReserve(t.Context())
		mr.Add("exec1", report)
		r, ok := mr.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		report.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.SetStep("ref1", []capabilities.MeteringNodeDetail{})
		mr.Delete("exec1")
	}()
	go func() {
		defer wg.Done()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		report.StartAndReserve(t.Context())
		report.balance.Add(10)
		mr.Add("exec2", report)
		r, ok := mr.Get("exec2")
		assert.True(t, ok)
		_, err := r.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		assert.NoError(t, err)
		err = r.SetStep("ref1", []capabilities.MeteringNodeDetail{})
		assert.NoError(t, err)
		mr.Delete("exec2")
	}()
	go func() {
		defer wg.Done()
		report := NewReport(testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		report.owner = testAccountID
		report.workflowID = testWorkflowID
		report.StartAndReserve(t.Context())
		mr.Add("exec1", report)
		r, ok := mr.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.DeductByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.SetStep("ref1", []capabilities.MeteringNodeDetail{})
		mr.Delete("exec1")
	}()

	wg.Wait()
	assert.Equal(t, 0, mr.Len())
}

func Test_MeterReportsLength(t *testing.T) {
	mr := NewReports(nil, "", "")

	mr.Add("exec1", NewReport(testWorkflowExecutionID, logger.TestSugared(t)))
	mr.Add("exec2", NewReport(testWorkflowExecutionID, logger.TestSugared(t)))
	mr.Add("exec3", NewReport(testWorkflowExecutionID, logger.TestSugared(t)))
	assert.Equal(t, 3, mr.Len())

	mr.Delete("exec2")
	assert.Equal(t, 2, mr.Len())
}
