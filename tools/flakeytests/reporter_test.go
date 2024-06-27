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
	ft := map[string]map[string]struct{}{
		"core/assets": map[string]struct{}{
			"TestLink": struct{}{},
		},
	}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(ft)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, "{\"package\":\"core/assets\",\"test_name\":\"TestLink\",\"fq_test_name\":\"core/assets:TestLink\"}"},
		{ts, "{\"num_flakes\":1}"},
	})
}

func TestMakeRequest_MultipleTests(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	ft := map[string]map[string]struct{}{
		"core/assets": map[string]struct{}{
			"TestLink": struct{}{},
			"TestCore": struct{}{},
		},
	}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(ft)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})

	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, "{\"package\":\"core/assets\",\"test_name\":\"TestLink\",\"fq_test_name\":\"core/assets:TestLink\"}"},
		{ts, "{\"package\":\"core/assets\",\"test_name\":\"TestCore\",\"fq_test_name\":\"core/assets:TestCore\"}"},
		{ts, "{\"num_flakes\":2}"},
	})
}

func TestMakeRequest_NoTests(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	ft := map[string]map[string]struct{}{}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(ft)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, "{\"num_flakes\":0}"},
	})
}

func TestMakeRequest_WithContext(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixNano())
	ft := map[string]map[string]struct{}{}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }, ctx: Context{CommitSHA: "42"}}
	pr, err := lr.createRequest(ft)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.ElementsMatch(t, pr.Streams[0].Values, [][]string{
		{ts, "{\"num_flakes\":0,\"commit_sha\":\"42\"}"},
	})
}
