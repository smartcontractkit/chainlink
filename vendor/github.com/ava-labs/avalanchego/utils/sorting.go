// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"bytes"
	"sort"
)

// IsSortedAndUnique returns true if the elements in the data are unique and sorted.
func IsSortedAndUnique(data sort.Interface) bool {
	for i := data.Len() - 2; i >= 0; i-- {
		if !data.Less(i, i+1) {
			return false
		}
	}
	return true
}

type innerSortUint32 []uint32

func (su32 innerSortUint32) Less(i, j int) bool { return su32[i] < su32[j] }
func (su32 innerSortUint32) Len() int           { return len(su32) }
func (su32 innerSortUint32) Swap(i, j int)      { su32[j], su32[i] = su32[i], su32[j] }

// SortUint32 sorts an uint32 array
func SortUint32(u32 []uint32) { sort.Sort(innerSortUint32(u32)) }

// IsSortedAndUniqueUint32 returns true if the array of uint32s are sorted and unique
func IsSortedAndUniqueUint32(arr []uint32) bool {
	for i := 0; i < len(arr)-1; i++ {
		if arr[i] >= arr[i+1] {
			return false
		}
	}
	return true
}

type innerSortUint64 []uint64

func (su64 innerSortUint64) Less(i, j int) bool { return su64[i] < su64[j] }
func (su64 innerSortUint64) Len() int           { return len(su64) }
func (su64 innerSortUint64) Swap(i, j int)      { su64[j], su64[i] = su64[i], su64[j] }

// SortUint64 sorts an uint64 array
func SortUint64(u64 []uint64) { sort.Sort(innerSortUint64(u64)) }

// IsSortedAndUniqueUint64 returns true if the array of uint64s are sorted and unique
func IsSortedAndUniqueUint64(u64 []uint64) bool { return IsSortedAndUnique(innerSortUint64(u64)) }

type innerSortBytes [][]byte

func (arr innerSortBytes) Less(i, j int) bool {
	return bytes.Compare(arr[i], arr[j]) == -1
}

func (arr innerSortBytes) Len() int      { return len(arr) }
func (arr innerSortBytes) Swap(i, j int) { arr[j], arr[i] = arr[i], arr[j] }

// Sort2DBytes sorts a 2D byte array
// Each byte array is not sorted internally; the byte arrays are sorted relative to another.
func Sort2DBytes(arr [][]byte) { sort.Sort(innerSortBytes(arr)) }

// IsSorted2DBytes returns true iff [arr] is sorted
func IsSorted2DBytes(arr [][]byte) bool { return sort.IsSorted(innerSortBytes(arr)) }
