package utils

import "time"

// Nower is an interface that fulfills the Now method,
// following the behavior of time.Now.
type Nower interface {
	Now() time.Time
}

// Afterer is an interface that fulfills the After method,
// following the behavior of time.After.
type Afterer interface {
	After(d time.Duration) <-chan time.Time
}

// AfterNower is an interface that fulfills the `After()` and `Now()`
// methods.
//go:generate mockery --name AfterNower --output ../internal/mocks/ --case=underscore
type AfterNower interface {
	After(d time.Duration) <-chan time.Time
	Now() time.Time
}

// Clock is a basic type for scheduling events in the application.
type Clock struct{}

// Now returns the current time.
func (Clock) Now() time.Time {
	return time.Now()
}

// After returns the current time if the given duration has elapsed.
func (Clock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
