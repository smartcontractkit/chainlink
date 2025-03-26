package workflows

import (
	"strconv"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func TestMeteringReport(t *testing.T) {
	t.Parallel()

	testUnitA := MeteringSpendUnit("a")
	testUnitB := MeteringSpendUnit("b")

	t.Run("MedianSpend returns median for multiple spend units", func(t *testing.T) {
		t.Parallel()

		report := NewMeteringReport()
		steps := []MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
			{"xyz", testUnitA, testUnitA.IntToSpendValue(2)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(3)},
			{"abc", testUnitB, testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.1))},
			{"xyz", testUnitB, testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.2))},
			{"abc", testUnitB, testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.3))},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[MeteringSpendUnit]MeteringSpendValue{
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

		report := NewMeteringReport()
		steps := []MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[MeteringSpendUnit]MeteringSpendValue{
			testUnitA: testUnitA.IntToSpendValue(1),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median odd number of spend values", func(t *testing.T) {
		t.Parallel()

		report := NewMeteringReport()
		steps := []MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(3)},
			{"xyz", testUnitA, testUnitA.IntToSpendValue(2)},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[MeteringSpendUnit]MeteringSpendValue{
			testUnitA: testUnitA.IntToSpendValue(2),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median as average for even number of spend values", func(t *testing.T) {
		t.Parallel()

		report := NewMeteringReport()
		steps := []MeteringReportStep{
			{"xyz", testUnitA, testUnitA.IntToSpendValue(42)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(3)},
			{"xyz", testUnitA, testUnitA.IntToSpendValue(2)},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[MeteringSpendUnit]MeteringSpendValue{
			testUnitA: testUnitA.DecimalToSpendValue(decimal.NewFromFloat(2.5)),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("SetStep returns error if step already exists", func(t *testing.T) {
		t.Parallel()
		report := NewMeteringReport()
		steps := []MeteringReportStep{
			{"xyz", testUnitA, testUnitA.IntToSpendValue(42)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
		}

		require.NoError(t, report.SetStep("ref1", steps))
		require.Error(t, report.SetStep("ref1", steps))
	})
}

// Test_MeterReports tests the Add, Get, Delete, and Len methods of a MeterReports.
// It also tests concurrent safe access.
func Test_MeterReports(t *testing.T) {
	mr := NewMeterReports()
	assert.Equal(t, 0, mr.Len())
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		mr.Add("exec1", NewMeteringReport())
		r, ok := mr.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.SetStep("ref1", []MeteringReportStep{})
		mr.Delete("exec1")
	}()
	go func() {
		defer wg.Done()
		mr.Add("exec2", NewMeteringReport())
		r, ok := mr.Get("exec2")
		assert.True(t, ok)
		err := r.SetStep("ref1", []MeteringReportStep{})
		assert.NoError(t, err)
		mr.Delete("exec2")
	}()
	go func() {
		defer wg.Done()
		mr.Add("exec1", NewMeteringReport())
		r, ok := mr.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.SetStep("ref1", []MeteringReportStep{})
		mr.Delete("exec1")
	}()

	wg.Wait()
	assert.Equal(t, 0, mr.Len())
}

func Test_MeterReportsLength(t *testing.T) {
	mr := NewMeterReports()

	mr.Add("exec1", NewMeteringReport())
	mr.Add("exec2", NewMeteringReport())
	mr.Add("exec3", NewMeteringReport())
	assert.Equal(t, 3, mr.Len())

	mr.Delete("exec2")
	assert.Equal(t, 2, mr.Len())
}
