// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import "strings"

// NodeIDSet is a set of NodeIDs
type NodeIDSet map[NodeID]struct{}

// Return a new NodeIDSet with initial capacity [size].
// More or less than [size] elements can be added to this set.
// Using NewNodeIDSet() rather than ids.NodeIDSet{} is just an optimization that can
// be used if you know how many elements will be put in this set.
func NewNodeIDSet(size int) NodeIDSet {
	if size < 0 {
		return NodeIDSet{}
	}
	return make(map[NodeID]struct{}, size)
}

func (ids *NodeIDSet) init(size int) {
	if *ids == nil {
		if minShortSetSize > size {
			size = minShortSetSize
		}
		*ids = make(map[NodeID]struct{}, size)
	}
}

// Add all the ids to this set, if the id is already in the set, nothing happens
func (ids *NodeIDSet) Add(idList ...NodeID) {
	ids.init(2 * len(idList))
	for _, id := range idList {
		(*ids)[id] = struct{}{}
	}
}

// Union adds all the ids from the provided set to this set.
func (ids *NodeIDSet) Union(idSet NodeIDSet) {
	ids.init(2 * idSet.Len())
	for id := range idSet {
		(*ids)[id] = struct{}{}
	}
}

// Difference removes all the ids from the provided set to this set.
func (ids *NodeIDSet) Difference(idSet NodeIDSet) {
	for id := range idSet {
		delete(*ids, id)
	}
}

// Contains returns true if the set contains this id, false otherwise
func (ids *NodeIDSet) Contains(id NodeID) bool {
	ids.init(1)
	_, contains := (*ids)[id]
	return contains
}

// Len returns the number of ids in this set
func (ids NodeIDSet) Len() int { return len(ids) }

// Remove all the id from this set, if the id isn't in the set, nothing happens
func (ids *NodeIDSet) Remove(idList ...NodeID) {
	ids.init(1)
	for _, id := range idList {
		delete(*ids, id)
	}
}

// Clear empties this set
func (ids *NodeIDSet) Clear() { *ids = nil }

// CappedList returns a list of length at most [size].
// Size should be >= 0. If size < 0, returns nil.
func (ids NodeIDSet) CappedList(size int) []NodeID {
	if size < 0 {
		return nil
	}
	if l := ids.Len(); l < size {
		size = l
	}
	i := 0
	idList := make([]NodeID, size)
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
func (ids NodeIDSet) List() []NodeID {
	idList := make([]NodeID, len(ids))
	i := 0
	for id := range ids {
		idList[i] = id
		i++
	}
	return idList
}

// SortedList returns this set as a sorted list
func (ids NodeIDSet) SortedList() []NodeID {
	lst := ids.List()
	SortNodeIDs(lst)
	return lst
}

// Equals returns true if the sets contain the same elements
func (ids NodeIDSet) Equals(oIDs NodeIDSet) bool {
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
func (ids NodeIDSet) String() string {
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
func (ids *NodeIDSet) Peek() (NodeID, bool) {
	for id := range *ids {
		return id, true
	}
	return NodeID{}, false
}

// Removes and returns an element. If the set is empty, does nothing and returns
// false
func (ids *NodeIDSet) Pop() (NodeID, bool) {
	for id := range *ids {
		delete(*ids, id)
		return id, true
	}
	return NodeID{}, false
}
