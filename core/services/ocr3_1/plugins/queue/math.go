package queue

import (
	"slices"
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// median returns the median value of a slice of integers.
// len must be greater than 0.
func median[N constraints.Integer, S []N](s S) N {
	switch len(s) {
	case 0:
		panic("median: len must be greater than 0")
	case 1:
		return s[0]
	}
	slices.Sort(s)
	if len(s)%2 == 1 {
		return s[len(s)/2]
	}
	// avg of two
	a, b := s[len(s)/2], s[len(s)/2-1]
	diff := a - b
	return b + diff/2
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
