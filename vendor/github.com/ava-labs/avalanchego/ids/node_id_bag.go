// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"fmt"
	"strings"
)

// NodeIDBag is a multiset of NodeIDs.
type NodeIDBag struct {
	counts map[NodeID]int
	size   int
}

func (b *NodeIDBag) init() {
	if b.counts == nil {
		b.counts = make(map[NodeID]int, minBagSize)
	}
}

// Add increases the number of times each id has been seen by one.
func (b *NodeIDBag) Add(ids ...NodeID) {
	for _, id := range ids {
		b.AddCount(id, 1)
	}
}

// AddCount increases the nubmer of times the id has been seen by count.
//
// count must be >= 0
func (b *NodeIDBag) AddCount(id NodeID, count int) {
	if count <= 0 {
		return
	}

	b.init()

	totalCount := b.counts[id] + count
	b.counts[id] = totalCount
	b.size += count
}

// Count returns the number of times the id has been added.
func (b *NodeIDBag) Count(id NodeID) int {
	b.init()
	return b.counts[id]
}

// Remove sets the count of the provided ID to zero.
func (b *NodeIDBag) Remove(id NodeID) {
	b.init()
	count := b.counts[id]
	delete(b.counts, id)
	b.size -= count
}

// Len returns the number of times an id has been added.
func (b *NodeIDBag) Len() int { return b.size }

// List returns a list of all IDs that have been added,
// without duplicates.
// e.g. a bag with {ID1, ID1, ID2} returns ids.ShortID[]{ID1, ID2}
func (b *NodeIDBag) List() []NodeID {
	idList := make([]NodeID, len(b.counts))
	i := 0
	for id := range b.counts {
		idList[i] = id
		i++
	}
	return idList
}

// Equals returns true if the bags contain the same elements
func (b *NodeIDBag) Equals(oIDs NodeIDBag) bool {
	if b.Len() != oIDs.Len() {
		return false
	}
	for key, value := range b.counts {
		if value != oIDs.counts[key] {
			return false
		}
	}
	return true
}

func (b *NodeIDBag) PrefixedString(prefix string) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Bag: (Size = %d)", b.Len()))
	for id, count := range b.counts {
		sb.WriteString(fmt.Sprintf("\n%s    ID[%s]: Count = %d", prefix, id, count))
	}

	return sb.String()
}

func (b *NodeIDBag) String() string { return b.PrefixedString("") }
