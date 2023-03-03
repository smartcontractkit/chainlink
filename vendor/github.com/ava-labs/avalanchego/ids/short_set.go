// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import "strings"

const (
	minShortSetSize = 16
)

// ShortSet is a set of ShortIDs
type ShortSet map[ShortID]struct{}

// Return a new ShortSet with initial capacity [size].
// More or less than [size] elements can be added to this set.
// Using NewShortSet() rather than ids.ShortSet{} is just an optimization that can
// be used if you know how many elements will be put in this set.
func NewShortSet(size int) ShortSet {
	if size < 0 {
		return ShortSet{}
	}
	return make(map[ShortID]struct{}, size)
}

func (ids *ShortSet) init(size int) {
	if *ids == nil {
		if minShortSetSize > size {
			size = minShortSetSize
		}
		*ids = make(map[ShortID]struct{}, size)
	}
}

// Add all the ids to this set, if the id is already in the set, nothing happens
func (ids *ShortSet) Add(idList ...ShortID) {
	ids.init(2 * len(idList))
	for _, id := range idList {
		(*ids)[id] = struct{}{}
	}
}

// Union adds all the ids from the provided set to this set.
func (ids *ShortSet) Union(idSet ShortSet) {
	ids.init(2 * idSet.Len())
	for id := range idSet {
		(*ids)[id] = struct{}{}
	}
}

// Difference removes all the ids from the provided set to this set.
func (ids *ShortSet) Difference(idSet ShortSet) {
	for id := range idSet {
		delete(*ids, id)
	}
}

// Contains returns true if the set contains this id, false otherwise
func (ids *ShortSet) Contains(id ShortID) bool {
	ids.init(1)
	_, contains := (*ids)[id]
	return contains
}

// Len returns the number of ids in this set
func (ids ShortSet) Len() int { return len(ids) }

// Remove all the id from this set, if the id isn't in the set, nothing happens
func (ids *ShortSet) Remove(idList ...ShortID) {
	ids.init(1)
	for _, id := range idList {
		delete(*ids, id)
	}
}

// Clear empties this set
func (ids *ShortSet) Clear() { *ids = nil }

// CappedList returns a list of length at most [size].
// Size should be >= 0. If size < 0, returns nil.
func (ids ShortSet) CappedList(size int) []ShortID {
	if size < 0 {
		return nil
	}
	if l := ids.Len(); l < size {
		size = l
	}
	i := 0
	idList := make([]ShortID, size)
	for id := range ids {
		if i >= size {
			break
		}
		idList[i] = id
		i++
	}
	return idList
}

// List converts this set into a list
func (ids ShortSet) List() []ShortID {
	idList := make([]ShortID, len(ids))
	i := 0
	for id := range ids {
		idList[i] = id
		i++
	}
	return idList
}

// SortedList returns this set as a sorted list
func (ids ShortSet) SortedList() []ShortID {
	lst := ids.List()
	SortShortIDs(lst)
	return lst
}

// Equals returns true if the sets contain the same elements
func (ids ShortSet) Equals(oIDs ShortSet) bool {
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
func (ids ShortSet) String() string {
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

// Returns an element. If the set is empty, returns false
func (ids *ShortSet) Peek() (ShortID, bool) {
	for id := range *ids {
		return id, true
	}
	return ShortID{}, false
}

// Removes and returns an element. If the set is empty, does nothing and returns
// false
func (ids *ShortSet) Pop() (ShortID, bool) {
	for id := range *ids {
		delete(*ids, id)
		return id, true
	}
	return ShortID{}, false
}
