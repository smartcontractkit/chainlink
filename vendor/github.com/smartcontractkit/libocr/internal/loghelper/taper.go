package loghelper

// LogarithmicTaper provides logarithmic tapering of an event sequence.
// For example, if the taper is Triggered 50 times with a function that
// simply prints the provided count, the output would be 1,2,4,8,16,32.
type LogarithmicTaper struct {
	count uint64
}

// Trigger increments a count and calls f iff the new count is a power of two
func (tap *LogarithmicTaper) Trigger(f func(newCount uint64)) {
	tap.count++
	if f != nil && isPowerOfTwo(tap.count) {
		f(tap.count)
	}
}

// Count returns the internal count of the taper
func (tap *LogarithmicTaper) Count() uint64 {
	return tap.count
}

// Reset resets the count to 0 and then calls f with the previous count
// iff it wasn't already 0
func (tap *LogarithmicTaper) Reset(f func(oldCount uint64)) {
	if tap.count != 0 {
		oldCount := tap.count
		tap.count = 0
		f(oldCount)
	}
}

func isPowerOfTwo(num uint64) bool {
	return num != 0 && (num&(num-1)) == 0
}
