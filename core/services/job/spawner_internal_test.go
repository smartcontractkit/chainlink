package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitPrerequisiteJobs(t *testing.T) {
	t.Parallel()

	ids := func(jbs []Job) []int32 {
		out := make([]int32, len(jbs))
		for i, jb := range jbs {
			out[i] = jb.ID
		}
		return out
	}

	t.Run("separates prerequisite jobs first, preserving order", func(t *testing.T) {
		t.Parallel()

		jbs := []Job{
			{ID: 1, Type: Type("offchainreporting2")},
			{ID: 2, Type: CRESettings},
			{ID: 3, Type: Type("cron")},
			{ID: 4, Type: CRESettings},
		}
		prerequisite, remaining := splitPrerequisiteJobs(jbs)
		assert.Equal(t, []int32{2, 4}, ids(prerequisite))
		assert.Equal(t, []int32{1, 3}, ids(remaining))
	})

	t.Run("no prerequisite jobs leaves all as remaining", func(t *testing.T) {
		t.Parallel()

		jbs := []Job{
			{ID: 1, Type: Type("offchainreporting2")},
			{ID: 2, Type: Type("cron")},
		}
		prerequisite, remaining := splitPrerequisiteJobs(jbs)
		assert.Empty(t, prerequisite)
		assert.Equal(t, []int32{1, 2}, ids(remaining))
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()

		prerequisite, remaining := splitPrerequisiteJobs(nil)
		assert.Empty(t, prerequisite)
		assert.Empty(t, remaining)
	})
}
