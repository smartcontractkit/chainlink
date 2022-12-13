package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMailbox(t *testing.T) {
	var (
		expected  = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		toDeliver = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	)

	const capacity = 10
	m := NewMailbox[int](capacity)

	// Queue deliveries
	for i, d := range toDeliver {
		atCapacity := m.Deliver(d)
		if atCapacity && i < capacity {
			t.Errorf("mailbox at capacity %d", i)
		} else if !atCapacity && i >= capacity {
			t.Errorf("mailbox below capacity %d", i)
		}
	}

	// Retrieve them
	var recvd []int
	chDone := make(chan struct{})
	go func() {
		defer close(chDone)
		for range m.Notify() {
			for {
				x, exists := m.Retrieve()
				if !exists {
					break
				}
				recvd = append(recvd, x)
			}
		}
	}()

	close(m.chNotify)
	<-chDone

	require.Equal(t, expected, recvd)
}

func TestMailbox_RetrieveAll(t *testing.T) {
	var (
		expected  = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		toDeliver = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	)

	const capacity = 10
	m := NewMailbox[int](capacity)

	// Queue deliveries
	for i, d := range toDeliver {
		atCapacity := m.Deliver(d)
		if atCapacity && i < capacity {
			t.Errorf("mailbox at capacity %d", i)
		} else if !atCapacity && i >= capacity {
			t.Errorf("mailbox below capacity %d", i)
		}
	}

	require.Equal(t, expected, m.RetrieveAll())
}

func TestMailbox_RetrieveLatestAndClear(t *testing.T) {
	var (
		expected  = 11
		toDeliver = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	)

	const capacity = 10
	m := NewMailbox[int](capacity)

	// Queue deliveries
	for i, d := range toDeliver {
		atCapacity := m.Deliver(d)
		if atCapacity && i < capacity {
			t.Errorf("mailbox at capacity %d", i)
		} else if !atCapacity && i >= capacity {
			t.Errorf("mailbox below capacity %d", i)
		}
	}

	require.Equal(t, expected, m.RetrieveLatestAndClear())
	require.Len(t, m.RetrieveAll(), 0)
}

func TestMailbox_NoEmptyReceivesWhenCapacityIsTwo(t *testing.T) {
	m := NewMailbox[int](2)

	var (
		recvd         []int
		emptyReceives []int
	)

	chDone := make(chan struct{})
	go func() {
		defer close(chDone)
		for range m.Notify() {
			x, exists := m.Retrieve()
			if !exists {
				emptyReceives = append(emptyReceives, recvd[len(recvd)-1])
			} else {
				recvd = append(recvd, x)
			}
		}
	}()

	for i := 0; i < 100000; i++ {
		m.Deliver(i)
	}
	close(m.chNotify)

	<-chDone
	require.Len(t, emptyReceives, 0)
}

func TestMailbox_load(t *testing.T) {
	for _, tt := range []struct {
		name     string
		capacity uint64
		deliver  []int
		exp      float64

		retrieve int
		exp2     float64

		all bool
	}{
		{"single-all", 1, []int{1}, 100, 0, 100, true},
		{"single-latest", 1, []int{1}, 100, 0, 100, false},
		{"ten-low", 10, []int{1}, 10, 1, 0.0, false},
		{"ten-full-all", 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 100, 5, 50, true},
		{"ten-full-latest", 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 100, 5, 50, false},
		{"ten-overflow", 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 100, 5, 50, false},
		{"nine", 9, []int{1, 2, 3}, 100.0 / 3.0, 2, 100.0 / 9.0, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMailbox[int](tt.capacity)

			// Queue deliveries
			for i, d := range tt.deliver {
				atCapacity := m.Deliver(d)
				if atCapacity && i < int(tt.capacity) {
					t.Errorf("mailbox at capacity %d", i)
				} else if !atCapacity && i >= int(tt.capacity) {
					t.Errorf("mailbox below capacity %d", i)
				}
			}
			gotCap, gotLoad := m.load()
			require.Equal(t, gotCap, tt.capacity)
			require.Equal(t, gotLoad, tt.exp)

			// Retrieve some
			for i := 0; i < tt.retrieve; i++ {
				_, ok := m.Retrieve()
				require.True(t, ok)
			}
			gotCap, gotLoad = m.load()
			require.Equal(t, gotCap, tt.capacity)
			require.Equal(t, gotLoad, tt.exp2)

			// Drain it
			if tt.all {
				m.RetrieveAll()
			} else {
				m.RetrieveLatestAndClear()
			}
			gotCap, gotLoad = m.load()
			require.Equal(t, gotCap, tt.capacity)
			require.Equal(t, gotLoad, 0.0)
		})
	}
}
