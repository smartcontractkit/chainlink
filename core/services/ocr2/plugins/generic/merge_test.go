package generic

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	vars := map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID": "some-job-id",
		},
	}
	addedVars := map[string]interface{}{
		"jb": map[string]interface{}{
			"some-other-var": "foo",
		},
		"val": 0,
	}

	merge(vars, addedVars)

	assert.True(t, reflect.DeepEqual(vars, map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":     "some-job-id",
			"some-other-var": "foo",
		},
		"val": 0,
	}), vars)
}
