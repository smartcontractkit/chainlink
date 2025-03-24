package workflows_test

import (
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
)

func TestMeteringReport(t *testing.T) {
	t.Parallel()

	testUnitA := workflows.MeteringSpendUnit("a")
	testUnitB := workflows.MeteringSpendUnit("b")

	t.Run("MedianSpend returns median for multiple spend units", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
			{"xyz", testUnitA, testUnitA.IntToSpendValue(2)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(3)},
			{"abc", testUnitB, testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.1))},
			{"xyz", testUnitB, testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.2))},
			{"abc", testUnitB, testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.3))},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
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

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			testUnitA: testUnitA.IntToSpendValue(1),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median odd number of spend values", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(3)},
			{"xyz", testUnitA, testUnitA.IntToSpendValue(2)},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			testUnitA: testUnitA.IntToSpendValue(2),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median as average for even number of spend values", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"xyz", testUnitA, testUnitA.IntToSpendValue(42)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(3)},
			{"xyz", testUnitA, testUnitA.IntToSpendValue(2)},
		}

		for idx := range steps {
			require.NoError(t, report.SetStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), steps))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			testUnitA: testUnitA.DecimalToSpendValue(decimal.NewFromFloat(2.5)),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})
}
