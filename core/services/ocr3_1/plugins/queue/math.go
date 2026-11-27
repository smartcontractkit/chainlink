package queue

import (
	"slices"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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
