package jobhelpers

import (
	"fmt"
	"runtime"

	"dario.cat/mergo"
)

func MergeSpecsByIndex(results []map[string][]string) (map[string][]string, error) {
	merged := make(map[string][]string)
	for i, result := range results {
		if result == nil {
			continue
		}
		if err := mergo.Merge(&merged, result, mergo.WithAppendSlice); err != nil {
			return nil, fmt.Errorf("failed to merge proposal result %d: %w", i, err)
		}
	}

	return merged, nil
}

func Parallelism(workItems int) int {
	if workItems <= 1 {
		return 1
	}

	limit := runtime.GOMAXPROCS(0)
	if workItems < limit {
		return workItems
	}

	return limit
}
