package workflows_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
)

func TestMeteringReport(t *testing.T) {
	t.Parallel()

	t.Run("AddStep returns error for invalid value", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		step := workflows.MeteringReportStep{
			Peer2PeerID: "abc",
			SpendUnit:   "a",
			SpendValue:  "not a value",
		}

		err := report.AddStep(workflows.MeteringReportStepRef("42"), step)

		require.ErrorIs(t, err, workflows.ErrInvalidMeteringSpendValue)
	})

	t.Run("MedianSpend returns median for multiple spend units", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", "a", "1.0"},
			{"xyz", "a", "2.0"},
			{"abc", "a", "3.0"},
			{"abc", "b", "0.1"},
			{"xyz", "b", "0.2"},
			{"abc", "b", "0.3"},
		}

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			"a": "2",
			"b": "0.2",
		}

		assert.Equal(t, expected, report.MedianSpend())
	})

	t.Run("MedianSpend returns median single spend value", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", "a", "1.0"},
		}

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			"a": "1",
		}

		assert.Equal(t, expected, report.MedianSpend())
	})

	t.Run("MedianSpend returns median odd number of spend values", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"abc", "a", "1.0"},
			{"abc", "a", "3.0"},
			{"xyz", "a", "2.0"},
		}

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			"a": "2",
		}

		assert.Equal(t, expected, report.MedianSpend())
	})

	t.Run("MedianSpend returns median as average for even number of spend values", func(t *testing.T) {
		t.Parallel()

		report := workflows.NewMeteringReport()
		steps := []workflows.MeteringReportStep{
			{"xyz", "a", "42.0"},
			{"abc", "a", "1.0"},
			{"abc", "a", "3.0"},
			{"xyz", "a", "2.0"},
		}

		for idx, step := range steps {
			require.NoError(t, report.AddStep(workflows.MeteringReportStepRef(strconv.Itoa(idx)), step))
		}

		expected := map[workflows.MeteringSpendUnit]workflows.MeteringSpendValue{
			"a": "2.5",
		}

		assert.Equal(t, expected, report.MedianSpend())
	})
}
