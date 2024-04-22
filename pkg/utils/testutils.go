package utils

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

// Deprecated: use tests.Context
func Context(t *testing.T) context.Context {
	return tests.Context(t)
}

// AssertJSONEqual is a helper function to assert that two JSON objects
// are equal.
//
// When they are not equal, it fails the test and provides a helpful diff.
func AssertJSONEqual(t tests.TestingT, x []byte, y []byte) {
	var TransformJSON = cmp.FilterValues(func(x, y []byte) bool {
		return json.Valid(x) && json.Valid(y)
	}, cmp.Transformer("ParseJSON", func(in []byte) (out interface{}) {
		if err := json.Unmarshal(in, &out); err != nil {
			panic(err) // should never occur given previous filter to ensure valid JSON
		}
		return out
	}))

	diff := cmp.Diff(x, y, TransformJSON)

	if diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
		t.FailNow()
	}
}
