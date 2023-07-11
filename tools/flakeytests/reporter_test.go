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
	ft := map[string][]string{
		"core/assets": {"TestLink"},
	}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(ft)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.Len(t, pr.Streams[0].Values, 1)
	assert.Equal(t, pr.Streams[0].Values[0], []string{fmt.Sprintf("%d", now.UnixNano()), "{\"package\":\"core/assets\",\"test_name\":\"TestLink\",\"fq_test_name\":\"core/assets:TestLink\"}"})
}

func TestMakeRequest_MultipleTests(t *testing.T) {
	now := time.Now()
	ft := map[string][]string{
		"core/assets": {"TestLink", "TestCore"},
	}
	lr := &LokiReporter{auth: "bla", host: "bla", command: "go_core_tests", now: func() time.Time { return now }}
	pr, err := lr.createRequest(ft)
	require.NoError(t, err)
	assert.Len(t, pr.Streams, 1)
	assert.Equal(t, pr.Streams[0].Stream, map[string]string{"command": "go_core_tests", "app": "flakey-test-reporter"})
	assert.Len(t, pr.Streams[0].Values, 2)

	ts := fmt.Sprintf("%d", now.UnixNano())
	assert.Equal(t, pr.Streams[0].Values, [][]string{
		{ts, "{\"package\":\"core/assets\",\"test_name\":\"TestLink\",\"fq_test_name\":\"core/assets:TestLink\"}"},
		{ts, "{\"package\":\"core/assets\",\"test_name\":\"TestCore\",\"fq_test_name\":\"core/assets:TestCore\"}"},
	})
}
