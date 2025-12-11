package generic

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestMerge(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	vars := map[string]any{
		"jb": map[string]any{
			"databaseID": "some-job-id",
		},
	}
	addedVars := map[string]any{
		"jb": map[string]any{
			"some-other-var": "foo",
		},
		"val": 0,
	}

	merge(vars, addedVars)

	assert.True(t, reflect.DeepEqual(vars, map[string]any{
		"jb": map[string]any{
			"databaseID":     "some-job-id",
			"some-other-var": "foo",
		},
		"val": 0,
	}), vars)
}
