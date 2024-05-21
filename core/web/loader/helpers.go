package loader

import (
	"strconv"

	"github.com/graph-gophers/dataloader"
)

// keyOrderInt64 returns the keys cast to int64 and a mapping of each key to
// their index order.
func keyOrderInt64(keys dataloader.Keys) ([]int64, map[string]int) {
	keyOrder := make(map[string]int, len(keys))

	var ids []int64
	for ix, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err == nil {
			ids = append(ids, id)
		}

		keyOrder[key.String()] = ix
	}

	return ids, keyOrder
}
