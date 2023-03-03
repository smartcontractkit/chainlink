// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"fmt"
	"strings"
)

const (
	minBagSize = 16
)

// Bag is a multiset of IDs.
//
// A bag has the ability to split and filter on its bits for ease of use for
// binary voting.
type Bag struct {
	counts map[ID]int
	size   int

	mode     ID
	modeFreq int

	threshold    int
	metThreshold Set
}

func (b *Bag) init() {
	if b.counts == nil {
		b.counts = make(map[ID]int, minBagSize)
	}
}

// SetThreshold sets the number of times an ID must be added to be contained in
// the threshold set.
func (b *Bag) SetThreshold(threshold int) {
	if b.threshold == threshold {
		return
	}

	b.threshold = threshold
	b.metThreshold.Clear()
	for vote, count := range b.counts {
		if count >= threshold {
			b.metThreshold.Add(vote)
		}
	}
}

// Add increases the number of times each id has been seen by one.
func (b *Bag) Add(ids ...ID) {
	for _, id := range ids {
		b.AddCount(id, 1)
	}
}

// AddCount increases the number of times the id has been seen by count.
//
// count must be >= 0
func (b *Bag) AddCount(id ID, count int) {
	if count <= 0 {
		return
	}

	b.init()

	totalCount := b.counts[id] + count
	b.counts[id] = totalCount
	b.size += count

	if totalCount > b.modeFreq {
		b.mode = id
		b.modeFreq = totalCount
	}
	if totalCount >= b.threshold {
		b.metThreshold.Add(id)
	}
}

// Count returns the number of times the id has been added.
func (b *Bag) Count(id ID) int {
	return b.counts[id]
}

// Len returns the number of times an id has been added.
func (b *Bag) Len() int { return b.size }

// List returns a list of all ids that have been added.
func (b *Bag) List() []ID {
	idList := make([]ID, len(b.counts))
	i := 0
	for id := range b.counts {
		idList[i] = id
		i++
	}
	return idList
}

// Equals returns true if the bags contain the same elements
func (b *Bag) Equals(oIDs Bag) bool {
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

// Mode returns the id that has been seen the most and the number of times it
// has been seen. Ties are broken by the first id to be seen the reported number
// of times.
func (b *Bag) Mode() (ID, int) { return b.mode, b.modeFreq }

// Threshold returns the ids that have been seen at least threshold times.
func (b *Bag) Threshold() Set { return b.metThreshold }

// Filter returns the bag of ids with the same counts as this bag, except all
// the ids in the returned bag must have the same bits in the range [start, end)
// as id.
func (b *Bag) Filter(start, end int, id ID) Bag {
	newBag := Bag{}
	for vote, count := range b.counts {
		if EqualSubset(start, end, id, vote) {
			newBag.AddCount(vote, count)
		}
	}
	return newBag
}

// Split returns the bags of ids with the same counts a this bag, except all ids
// in the 0th index have a 0 at bit [index], and all ids in the 1st index have a
// 1 at bit [index].
func (b *Bag) Split(index uint) [2]Bag {
	splitVotes := [2]Bag{}
	for vote, count := range b.counts {
		bit := vote.Bit(index)
		splitVotes[bit].AddCount(vote, count)
	}
	return splitVotes
}

func (b *Bag) PrefixedString(prefix string) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Bag: (Size = %d)", b.Len()))
	for id, count := range b.counts {
		sb.WriteString(fmt.Sprintf("\n%s    ID[%s]: Count = %d", prefix, id, count))
	}

	return sb.String()
}

func (b *Bag) String() string { return b.PrefixedString("") }
