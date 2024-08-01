package go_heaps

// Interface is basic interface that all Heaps implement.
type Interface interface {
	// Inserts an element to the heap and returns it
	Insert(v Item) Item

	// DeleteMin deletes and returns the smallest element
	DeleteMin() Item

	// FindMin returns the minimum element
	FindMin() Item

	// Removes all items
	Clear()
}

// Extended adds operations on heaps are often useful.
type Extended interface {
	Interface
	// Return the heap formed by taking the union of the item disjoint
	// current heap and a
	Meld(a Interface) Interface

	// Adjusts the key of item old in heap h to new
	Adjust(old, new Item) Item

	// Delete arbitrary item from heap h.
	Delete(item Item) Item
}

// Item is the basic element that is inserted in a heap
type Item interface {
	// Should return a number:
	//    negative , if a < b
	//    zero     , if a == b
	//    positive , if a > b
	Compare(than Item) int
}

// ItemIterator allows callers of Do to iterate in-order over portions of
// the tree.  When this function returns false, iteration will stop and the
// function will immediately return.
type ItemIterator func(item Item) bool

// String implements the Item interface
type String string

// Integer implements the Item interface
type Integer int

func (a String) Compare(b Item) int {
	s1 := a
	s2 := b.(String)
	min := len(s2)
	if len(s1) < len(s2) {
		min = len(s1)
	}
	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(s1[i]) - int(s2[i])
	}
	if diff == 0 {
		diff = len(s1) - len(s2)
	}
	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}

func (a Integer) Compare(b Item) int {
	a1 := a
	a2 := b.(Integer)
	switch {
	case a1 > a2:
		return 1
	case a1 < a2:
		return -1
	default:
		return 0
	}
}
