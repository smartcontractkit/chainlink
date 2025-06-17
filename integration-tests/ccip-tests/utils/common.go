package utils

import (
	"runtime"
	"sync"
)

// DeleteNilEntriesFromMap checks for nil entry in map, store all not-nil entries to another map and deallocates previous map
// Deleting keys from a map actually does not delete the key, It just sets the corresponding value to nil.
func DeleteNilEntriesFromMap(inputMap *sync.Map) *sync.Map {
	newMap := &sync.Map{}
	foundNil := false
	inputMap.Range(func(key, value any) bool {
		if value != nil {
			newMap.Store(key, value)
		}
		if value == nil {
			foundNil = true
		}
		return true
	})
	if foundNil {
		runtime.GC()
	}
	return newMap
}
