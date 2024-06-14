package validation

import (
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type counter[T any] struct {
	data  T
	count int
}

// MinObservationFilter provides a way to ensure a minimum number of observations for
// some piece of data have occurred. It maintains an internal cache and provides a list
// of valid or invalid data points.
type MinObservationFilter[T any] interface {
	Add(data T) error
	GetValid() ([]T, error)
}

// minObservationValidator is a helper object to validate reports for a single chain.
// It keeps track of all reports and determines if they observations are consistent
// with one another and whether they meet the required fChain threshold.
type minObservationValidator[T any] struct {
	minObservation int
	cache          map[cciptypes.Bytes32]*counter[T]
	idFunc         func(T) [32]byte
}

// NewMinObservationValidator constructs a concrete MinObservationFilter object. The
// supplied idFunc is used to generate a uniqueID for the type being observed.
func NewMinObservationValidator[T any](min int, idFunc func(T) [32]byte) MinObservationFilter[T] {
	return &minObservationValidator[T]{
		minObservation: min,
		cache:          make(map[cciptypes.Bytes32]*counter[T]),
		idFunc:         idFunc,
	}
}

func (cv *minObservationValidator[T]) Add(data T) error {
	id := cv.idFunc(data)
	if _, ok := cv.cache[id]; ok {
		cv.cache[id].count++
	} else {
		cv.cache[id] = &counter[T]{data: data, count: 1}
	}
	return nil
}

func (cv *minObservationValidator[T]) GetValid() ([]T, error) {
	var validated []T
	for _, rc := range cv.cache {
		if rc.count >= cv.minObservation {
			rc := rc
			validated = append(validated, rc.data)
		}
	}
	return validated, nil
}
