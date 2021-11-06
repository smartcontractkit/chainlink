package utils_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
)

func TestMailbox(t *testing.T) {
	m := utils.NewMailbox(10)

	var (
		expected  = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		toDeliver = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		recvd     []int
	)

	chDone := make(chan struct{})
	go func() {
		defer close(chDone)
		for _ = range m.Notify() {
			for {
				x, exists := m.Retrieve()
				if !exists {
					break
				}
				recvd = append(recvd, x.(int))
			}
		}
	}()

	for _, i := range toDeliver {
		m.Deliver(i)
	}
	close(m.Notify())

	<-chDone

	if len(recvd) > 10 {
		t.Fatal("received too many")
	} else if len(recvd) < 10 {
		t.Fatal("received too few")
	}
	require.Equal(t, expected, recvd)
}

func TestMailbox_EmptyReceivesWhenCapacityIsOne(t *testing.T) {
	m := utils.NewMailbox(1)

	var (
		recvd         []int
		emptyReceives []int
	)

	chDone := make(chan struct{})
	go func() {
		defer close(chDone)
		for _ = range m.Notify() {
			x, exists := m.Retrieve()
			if !exists {
				emptyReceives = append(emptyReceives, recvd[len(recvd)-1])
			} else {
				recvd = append(recvd, x.(int))
			}

		}
	}()

	for i := 0; i < 100000; i++ {
		m.Deliver(i)
	}
	close(m.Notify())

	<-chDone
	require.Greater(t, len(emptyReceives), 0)
}

func TestMailbox_NoEmptyReceivesWhenCapacityIsTwo(t *testing.T) {
	m := utils.NewMailbox(2)

	var (
		recvd         []int
		emptyReceives []int
	)

	chDone := make(chan struct{})
	go func() {
		defer close(chDone)
		for _ = range m.Notify() {
			x, exists := m.Retrieve()
			if !exists {
				emptyReceives = append(emptyReceives, recvd[len(recvd)-1])
			} else {
				recvd = append(recvd, x.(int))
			}
		}
	}()

	for i := 0; i < 100000; i++ {
		m.Deliver(i)
	}
	close(m.Notify())

	<-chDone
	require.Len(t, emptyReceives, 0)
}
