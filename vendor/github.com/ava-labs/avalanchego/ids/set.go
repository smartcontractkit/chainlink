// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"encoding/json"
	"strings"
)

const (
	// The minimum capacity of a set
	minSetSize = 16

	// If a set has more than this many keys, it will be cleared by setting the map to nil
	// rather than iteratively deleting
	clearSizeThreshold = 512
)

// Set is a set of IDs
type Set map[ID]struct{}

// Return a new set with initial capacity [size].
// More or less than [size] elements can be added to this set.
// Using NewSet() rather than ids.Set{} is just an optimization that can
// be used if you know how many elements will be put in this set.
func NewSet(size int) Set {
	if size < 0 {
		return Set{}
	}
	return make(map[ID]struct{}, size)
}

func (ids *Set) init(size int) {
	if *ids == nil {
		if minSetSize > size {
			size = minSetSize
		}
		*ids = make(map[ID]struct{}, size)
	}
}

// Add all the ids to this set, if the id is already in the set, nothing happens
func (ids *Set) Add(idList ...ID) {
	ids.init(2 * len(idList))
	for _, id := range idList {
		(*ids)[id] = struct{}{}
	}
}

// Union adds all the ids from the provided set to this set.
func (ids *Set) Union(set Set) {
	ids.init(2 * set.Len())
	for id := range set {
		(*ids)[id] = struct{}{}
	}
}

// Difference removes all the ids from the provided set to this set.
func (ids *Set) Difference(set Set) {
	for id := range set {
		delete(*ids, id)
	}
}

// Contains returns true if the set contains this id, false otherwise
func (ids *Set) Contains(id ID) bool {
	_, contains := (*ids)[id]
	return contains
}

// Overlaps returns true if the intersection of the set is non-empty
func (ids *Set) Overlaps(big Set) bool {
	small := *ids
	if small.Len() > big.Len() {
		small, big = big, small
	}

	for id := range small {
		if _, ok := big[id]; ok {
			return true
		}
	}
	return false
}

// Len returns the number of ids in this set
func (ids Set) Len() int { return len(ids) }

// Remove all the id from this set, if the id isn't in the set, nothing happens
func (ids *Set) Remove(idList ...ID) {
	for _, id := range idList {
		delete(*ids, id)
	}
}

// Clear empties this set
func (ids *Set) Clear() {
	if len(*ids) > clearSizeThreshold {
		*ids = nil
		return
	}
	for key := range *ids {
		delete(*ids, key)
	}
}

// List converts this set into a list
func (ids Set) List() []ID {
	idList := make([]ID, ids.Len())
	i := 0
	for id := range ids {
		idList[i] = id
		i++
	}
	return idList
}

// SortedList returns this set as a sorted list
func (ids Set) SortedList() []ID {
	lst := ids.List()
	SortIDs(lst)
	return lst
}

// CappedList returns a list of length at most [size].
// Size should be >= 0. If size < 0, returns nil.
func (ids Set) CappedList(size int) []ID {
	if size < 0 {
		return nil
	}
	if l := ids.Len(); l < size {
		size = l
	}
	i := 0
	idList := make([]ID, size)
	for id := range ids {
		if i >= size {
			break
		}
		idList[i] = id
		i++
	}
	return idList
}

// Equals returns true if the sets contain the same elements
func (ids Set) Equals(oIDs Set) bool {
	if ids.Len() != oIDs.Len() {
		return false
	}
	for key := range oIDs {
		if _, contains := ids[key]; !contains {
			return false
		}
	}
	return true
}

// String returns the string representation of a set
func (ids Set) String() string {
	sb := strings.Builder{}
	sb.WriteString("{")
	first := true
	for id := range ids {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(id.String())
	}
	sb.WriteString("}")
	return sb.String()
}

// Removes and returns an element. If the set is empty, does nothing and returns
// false.
func (ids *Set) Pop() (ID, bool) {
	for id := range *ids {
		delete(*ids, id)
		return id, true
	}
	return ID{}, false
}

func (ids *Set) MarshalJSON() ([]byte, error) {
	idsList := ids.List()
	SortIDs(idsList)
	return json.Marshal(idsList)
}
