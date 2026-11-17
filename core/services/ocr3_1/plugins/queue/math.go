package queue

import (
	"slices"
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// median returns the median value of a slice of integers, as well as a count of the number of occurrences of the median.
// len must be greater than 0.
func median[N constraints.Integer, S []N](s S) (N, int) {
	switch len(s) {
	case 0:
		panic("median: len must be greater than 0")
	case 1:
		return s[0], 1
	}
	slices.Sort(s)
	if len(s)%2 == 1 {
		m := s[len(s)/2]
		return m, count(s, m)
	}
	// avg of two
	a, b := s[len(s)/2], s[len(s)/2-1]
	diff := a - b
	m := b + diff/2
	return m, count(s, m)
}

// count returns the number of occurrences of e in s, but s must be sorted.
func count[N comparable, S []N](s S, e N) int {
	firstIndex := slices.Index(s, e)
	if firstIndex == -1 {
		return 0
	}
	var c int
	for i := firstIndex; i < len(s) && s[i] == e; i++ {
		c++
	}
	return c
}

// medianTimestamp returns the median value of a slice of timestamps.
// len must be greater than 0.
func medianTimestamp(s []*timestamppb.Timestamp) *timestamppb.Timestamp {
	switch len(s) {
	case 0:
		panic("median: len must be greater than 0")
	case 1:
		return s[0]
	}
	times := make([]time.Time, len(s))
	for i, t := range s {
		times[i] = time.Unix(t.Seconds, int64(t.Nanos))
	}
	slices.SortFunc(times, time.Time.Compare)

	if len(times)%2 == 1 {
		return timestamppb.New(times[len(times)/2])
	}
	// avg of two
	a, b := times[len(times)/2], times[len(times)/2-1]
	diff := a.Sub(b)
	return timestamppb.New(b.Add(diff / 2))
}
