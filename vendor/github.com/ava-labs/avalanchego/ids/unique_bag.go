// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"fmt"
	"strings"
)

const (
	minUniqueBagSize = 16
)

type UniqueBag map[ID]BitSet64

func (b *UniqueBag) init() {
	if *b == nil {
		*b = make(map[ID]BitSet64, minUniqueBagSize)
	}
}

func (b *UniqueBag) Add(setID uint, idSet ...ID) {
	bs := BitSet64(0)
	bs.Add(setID)

	for _, id := range idSet {
		b.UnionSet(id, bs)
	}
}

func (b *UniqueBag) UnionSet(id ID, set BitSet64) {
	b.init()

	previousSet := (*b)[id]
	previousSet.Union(set)
	(*b)[id] = previousSet
}

func (b *UniqueBag) DifferenceSet(id ID, set BitSet64) {
	b.init()

	previousSet := (*b)[id]
	previousSet.Difference(set)
	(*b)[id] = previousSet
}

func (b *UniqueBag) Difference(diff *UniqueBag) {
	b.init()

	for id, previousSet := range *b {
		if previousSetDiff, exists := (*diff)[id]; exists {
			previousSet.Difference(previousSetDiff)
		}
		(*b)[id] = previousSet
	}
}

func (b *UniqueBag) GetSet(id ID) BitSet64 { return (*b)[id] }

func (b *UniqueBag) RemoveSet(id ID) { delete(*b, id) }

func (b *UniqueBag) List() []ID {
	idList := make([]ID, len(*b))
	i := 0
	for id := range *b {
		idList[i] = id
		i++
	}
	return idList
}

func (b *UniqueBag) Bag(alpha int) Bag {
	bag := Bag{
		counts: make(map[ID]int, len(*b)),
	}
	bag.SetThreshold(alpha)
	for id, bs := range *b {
		bag.AddCount(id, bs.Len())
	}
	return bag
}

func (b *UniqueBag) PrefixedString(prefix string) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("UniqueBag: (Size = %d)", len(*b)))
	for id, set := range *b {
		sb.WriteString(fmt.Sprintf("\n%s    ID[%s]: Members = %s", prefix, id, set))
	}

	return sb.String()
}

func (b *UniqueBag) String() string { return b.PrefixedString("") }

func (b *UniqueBag) Clear() {
	for id := range *b {
		delete(*b, id)
	}
}
