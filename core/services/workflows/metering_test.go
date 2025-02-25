package workflows_test

import (
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			testUnitA: testUnitB.IntToSpendValue(2),
			testUnitB: testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.2)),
		}

		assert.Equal(t, expected, report.MedianSpend())
	})

	t.Run("MedianSpend returns median single spend value", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
		}

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			testUnitA: testUnitA.IntToSpendValue(1),
		}

		assert.Equal(t, expected, report.MedianSpend())
	})

	t.Run("MedianSpend returns median odd number of spend values", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", testUnitA, testUnitA.IntToSpendValue(1)},
			{"abc", testUnitA, testUnitA.IntToSpendValue(3)},
			{"xyz", testUnitA, testUnitA.IntToSpendValue(2)},
		}

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			testUnitA: testUnitA.IntToSpendValue(2),
		}

		assert.Equal(t, expected, report.MedianSpend())
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

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			testUnitA: testUnitA.DecimalToSpendValue(decimal.NewFromFloat(2.5)),
		}

		assert.Equal(t, expected, report.MedianSpend())
	})
}
