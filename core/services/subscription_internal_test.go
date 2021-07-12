package services

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/services/eth/mocks"
	"github.com/stretchr/testify/mock"
)

var _ ethereum.Subscription = EthSub{}

type EthSub struct {
	err chan error
	t   *testing.T
}

func (es EthSub) Unsubscribe() {
	es.t.Log("unsubscribe called")
}

func (es EthSub) Err() <-chan error {
	es.t.Log("err chan called")
	return es.err
}

func TestSubscriptionListenToLogs(t *testing.T) {
	c := make(chan types.Log)
	done := make(chan struct{})
	err := make(chan error)
	ethClient := new(mocks.Client)
	s := InitiatorSubscription{
		done:          done,
		logSubscriber: ethClient,
		logs:          c,
		ethSubscription: EthSub{
			err: err,
			t:   t,
		},
	}
	done2 := make(chan struct{})
	go func() {
		s.Start()
		done2 <- struct{}{}
	}()

	// Force a reconnect
	err2 := make(chan error)
	ethSub2 := EthSub{
		err: err2,
		t:   t,
	}
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Return(ethSub2, nil)
	err <- errors.New("aahh websocket down")

	// Wait for reconnect to complete
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return s.logs != c
	}, time.Second, 10*time.Millisecond).Should(gomega.BeTrue())

	// Ensure we can handle logs after reconnecting
	a := common.HexToAddress("0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c")
	select {
	case s.logs <- types.Log{Address: a}:
		break
	case <-time.After(5 * time.Second):
		t.Error("log listener did not read logs after reconnecting")
	}

	// Unsubscribe and we expect the goroutine to close.
	s.Unsubscribe()
	select {
	case <-done:
		break
	case <-time.After(5 * time.Second):
		t.Error("log listener did not close as expected")
	}
}
