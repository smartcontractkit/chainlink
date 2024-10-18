package flakeytests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeRequest_SingleTest(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	r := &Report{
		tests: map[string]map[string]int{
			"core/assets": map[string]int{
				"TestLink": 1,
			},
		},
	}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(r)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, `{"message_type":"flakey_test","commit_sha":"","repository":"","event_type":"","package":"core/assets","test_name":"TestLink","fq_test_name":"core/assets:TestLink"}`},
		{ts, `{"message_type":"run_report","commit_sha":"","repository":"","event_type":"","num_package_panics":0,"num_flakes":1,"num_combined":1}`},
	})
}

func TestMakeRequest_MultipleTests(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	r := &Report{
		tests: map[string]map[string]int{
			"core/assets": map[string]int{
				"TestLink": 1,
				"TestCore": 1,
			},
		},
	}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(r)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})

	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, `{"message_type":"flakey_test","commit_sha":"","repository":"","event_type":"","package":"core/assets","test_name":"TestLink","fq_test_name":"core/assets:TestLink"}`},
		{ts, `{"message_type":"flakey_test","commit_sha":"","repository":"","event_type":"","package":"core/assets","test_name":"TestCore","fq_test_name":"core/assets:TestCore"}`},
		{ts, `{"message_type":"run_report","commit_sha":"","repository":"","event_type":"","num_package_panics":0,"num_flakes":2,"num_combined":2}`},
	})
}

func TestMakeRequest_NoTests(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	r := NewReport()
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(r)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, `{"message_type":"run_report","commit_sha":"","repository":"","event_type":"","num_package_panics":0,"num_flakes":0,"num_combined":0}`},
	})
}

func TestMakeRequest_WithContext(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	r := NewReport()
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }, ctx: Context{CommitSHA: "42"}}
	pr, err := lr.createRequest(r)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, `{"message_type":"run_report","commit_sha":"42","repository":"","event_type":"","num_package_panics":0,"num_flakes":0,"num_combined":0}`},
	})
}

func TestMakeRequest_Panics(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	r := &Report{
		tests: map[string]map[string]int{
			"core/assets": map[string]int{
				"TestLink": 1,
			},
		},
		packagePanics: map[string]int{
			"core/assets": 1,
		},
	}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(r)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})

	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, `{"message_type":"flakey_test","commit_sha":"","repository":"","event_type":"","package":"core/assets","test_name":"TestLink","fq_test_name":"core/assets:TestLink"}`},
		{ts, `{"message_type":"package_panic","commit_sha":"","repository":"","event_type":"","package":"core/assets"}`},
		{ts, `{"message_type":"run_report","commit_sha":"","repository":"","event_type":"","num_package_panics":1,"num_flakes":1,"num_combined":2}`},
	})
}

func TestDedupeEntries(t *testing.T) {
	r := &Report{
		tests: map[string]map[string]int{
			"core/assets": map[string]int{
				"TestSomethingAboutAssets/test_1": 2,
				"TestSomethingAboutAssets":        4,
				"TestSomeOtherTest":               1,
				"TestSomethingAboutAssets/test_2": 2,
				"TestFinalTest/test_1":            1,
			},
			"core/services/important_service": map[string]int{
				"TestAnImportantService/a_subtest": 1,
			},
		},
	}

	otherReport, err := dedupeEntries(r)
	require.NoError(t, err)

	expectedMap := map[string]map[string]int{
		"core/assets": map[string]int{
			"TestSomethingAboutAssets/test_1": 2,
			"TestSomeOtherTest":               1,
			"TestSomethingAboutAssets/test_2": 2,
			"TestFinalTest/test_1":            1,
		},
		"core/services/important_service": map[string]int{
			"TestAnImportantService/a_subtest": 1,
		},
	}
	assert.Equal(t, expectedMap, otherReport.tests)
}
