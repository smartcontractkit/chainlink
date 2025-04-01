package operations

import (
	"errors"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MemoryReporter(t *testing.T) {
	t.Parallel()

	now := time.Now()
	existingReport := Report[any, any]{
		ID:                    "1",
		Def:                   Definition{},
		Output:                "2",
		Input:                 1,
		Timestamp:             &now,
		ChildOperationReports: []string{uuid.New().String()},
	}

	reporter := NewMemoryReporter(WithReports([]Report[any, any]{existingReport}))

	reports, err := reporter.GetReports()
	require.NoError(t, err)
	assert.Len(t, reports, 1)
	assert.Equal(t, existingReport, reports[0])

	report, err := reporter.GetReport("1")
	require.NoError(t, err)
	assert.Equal(t, report, reports[0])

	newReport := Report[any, any]{
		ID:        "2",
		Def:       Definition{},
		Output:    "3",
		Input:     2,
		Timestamp: &now,
	}
	err = reporter.AddReport(newReport)
	require.NoError(t, err)
	reports, err = reporter.GetReports()
	require.NoError(t, err)
	assert.Len(t, reports, 2)
	assert.Equal(t, newReport, reports[1])

	// get non-existing report
	_, err = reporter.GetReport("100")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrReportNotFound)
	assert.ErrorContains(t, err, "report_id 100: report not found")
}

func Test_NewReport(t *testing.T) {
	t.Parallel()

	version := semver.MustParse("1.0.0")
	description := "test operation"
	handler := func(b Bundle, deps any, input int) (output int, err error) {
		return input + 1, nil
	}
	op := NewOperation("plus1", version, description, handler)

	testErr := errors.New("test error")
	childOperationID := uuid.New().String()
	report := NewReport[int, int](op.def, 1, 2, testErr, childOperationID)
	assert.NotEmpty(t, report.ID)
	assert.Equal(t, op.def, report.Def)
	assert.Equal(t, 1, report.Input)
	assert.Equal(t, 2, report.Output)
	assert.NotEmpty(t, report.Timestamp)
	require.ErrorContains(t, report.Err, testErr.Error())
	assert.Len(t, report.ChildOperationReports, 1)
	assert.Equal(t, childOperationID, report.ChildOperationReports[0])
}

func Test_RecentReporter(t *testing.T) {
	t.Parallel()

	now := time.Now()
	existingReport := Report[any, any]{
		ID:                    "1",
		Def:                   Definition{},
		Output:                "2",
		Input:                 1,
		Timestamp:             &now,
		ChildOperationReports: []string{uuid.New().String()},
	}

	reporter := NewMemoryReporter(WithReports([]Report[any, any]{existingReport}))
	recentReporter := NewRecentMemoryReporter(reporter)

	// no new reports added since creation of recentReporter
	reports := recentReporter.GetRecentReports()
	assert.Empty(t, reports)

	newReport := Report[any, any]{
		ID:        "2",
		Def:       Definition{},
		Output:    "3",
		Input:     2,
		Timestamp: &now,
	}
	err := recentReporter.AddReport(newReport)
	require.NoError(t, err)
	reports = recentReporter.GetRecentReports()
	assert.Len(t, reports, 1)
	assert.Equal(t, newReport, reports[0])
}

func Test_typeReport(t *testing.T) {
	t.Parallel()

	now := time.Now()
	report := Report[any, any]{
		ID:                    "1",
		Def:                   Definition{},
		Output:                2,
		Input:                 1,
		Timestamp:             &now,
		Err:                   nil,
		ChildOperationReports: []string{uuid.New().String()},
	}

	_, ok := typeReport[int, int](report)
	assert.True(t, ok)

	// incorrect input type
	_, ok = typeReport[string, int](report)
	assert.False(t, ok)

	// incorrect output type
	_, ok = typeReport[int, string](report)
	assert.False(t, ok)
}
